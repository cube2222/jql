# jql
Hey there!

You're probably here cause you're fed up with other json query processors being too complicated to use for anything surpassing simple single field selection.

Well, at least that's what led me here. And that's why I've written **jql**, a json query processor with an even more cryptic - lispy - syntax (Maybe not that cryptic after all? I quite like it :) )

**jql** aims to be friendly, helpful, pat you on the back when you're faced with a monstrosity of a JSON blob. Help you mold it into something useful, step by step.

If you want to see **benchmarks**, they're at [the end of the README.](#Benchmarks)

Ok, let's check it out now, but first things first, you have to install it:
```
go get github.com/cube2222/jql
```

Ok. Done.

Let's check out a few _simple_ examples. (remember? That's explicitly **not** why we're here. But it aids understanding of the more complex examples, so stay with me just a little bit longer!)

We'll be working with this piece of json:
```json
{ 
  "count": 3,
  "countries": [
    {
      "name": "Poland",
      "population": 38000000,
      "european": true,
      "eu_since": "2004"
    },
    {
      "name": "United States",
      "population": 327000000,
      "european": false
    },
    {
      "name": "Germany",
      "population": 83000000,
      "european": true,
      "eu_since": "1993"
    }
  ]
}
```
To start with, let's get the countries array only.
```
> cat test.json | jql '(elem "countries")'
[
  {
    "eu_since": "2004",
    "european": true,
    "name": "Poland",
    "population": 38000000
  },
  {
    "european": false,
    "name": "United States",
    "population": 327000000
  },
  {
    "eu_since": "1993",
    "european": true,
    "name": "Germany",
    "population": 83000000
  }
]
```

Whoa, whoa, whoa! What's up with the parentheses?!

Remember what I said before? Lispy syntax. In short, whenever you see (add 1 2), it's basically the same as add(1, 2). I like that syntax for big hierarchical expressions, if you really really _really_ don't like it, then you can probably stop right here, though you'll be losing out!

I've warned you.

What if we wanted only the first country from that list?
```
> cat test.json | jql '(elem "countries" (elem 0))'
{
  "eu_since": "2004",
  "european": true,
  "name": "Poland",
  "population": 38000000
}
```
Let's break down what happened here. First we took the "countries" field. elem takes an additional argument, it says "how to transform the element" and is also an expression. Here we say we want to take the first element of the countries array. The default function is _id_, which stands for identity.

#### array

We can also pass an array of positions to elem, to get more than one country:
```
> cat test.json | jql '(elem "countries" (elem (array 0 2)))'
[
  {
    "eu_since": "2004",
    "european": true,
    "name": "Poland",
    "population": 38000000
  },
  {
    "eu_since": "1993",
    "european": true,
    "name": "Germany",
    "population": 83000000
  }
]
```

elem can work with single strings, single integers, arrays of those, and objects with them as values (but we won't cover those now to keep things simple).

#### keys

What if we want to get all the country names? A new friend - _keys_ - is key in this situation. 🗝
```
cat test.json | jql '(elem "countries" (elem (keys) (elem "name")))'
[
  "Poland",
  "United States",
  "Germany"
]
```

It returns an array of all the keys of the given collection. Fields and Indices for Objects and Arrays respectively.

To illustrate, here's _keys_ used on an object:
```
> cat test.json | jql '("countries" (0 (keys)))'
[
  "name",
  "population",
  "european",
  "eu_since"
]
```

---
### Attention

Now we have to understand a very important mechanism underlying **jql**. All functions operate in the context of the JSON we're operating on.

Some functions, like elem, will cut down the context for expressions it evaluates. The first argument - which should evaluate to the positions we need - gets evaluated in the context of the entire array, that's why _keys_ returns all the indices. The second one on the other hand, operates in the context of a single element.

In theory we're really just creating and composing a big function - pipeline, so to say - which gets applied to our JSON blob.

This may sound complicated, but I find it becomes intuitive quite quickly.

---

You can see that elem is the most used function, and in fact it's what you'll usually be using when munging data, so there's a shortcut. If you put a value in function name position, it implicitly converts it to an elem.

This way we can rewrite the previous query to be much shorter, and better match the shape of the data.
```
> cat test.json | jql '("countries" ((keys) ("name")))'
[
  "Poland",
  "United States",
  "Germany"
]
```

---
### Little showcase
Remember when I said you can use integers, strings, arrays and objects as positions?
```
> cat test.json | jql '("countries" ((array (array 0 (array 0 (array 0 (array 0 2)))) 1 (object "key1" 1 "key2" (array 0 (object "key1" 1 "key2" (array 0 2))))) ("population")))'
[
  [
    38000000,
    [
      38000000,
      [
        38000000,
        [
          38000000,
          83000000
        ]
      ]
    ]
  ],
  327000000,
  {
    "key1": 327000000,
    "key2": [
      38000000,
      {
        "key1": 327000000,
        "key2": [
          38000000,
          83000000
        ]
      }
    ]
  }
]
```
🔥🔥🔥

Don't do this.

---

#### range

We can also select a range of elements, using the... you guessed it - _range_ function.
```
> cat test.json | jql '("countries" ((range 1 3) ("name")))'
[
  "United States",
  "Germany"
]
```

You can use the _array_ function in value position too obviously. If you want a list of name-population tuples you can just
```
> cat test.json | jql '("countries" ((keys) (array ("name") ("population"))))'
[
  [
    "Poland",
    38000000
  ],
  [
    "United States",
    327000000
  ],
  [
    "Germany",
    83000000
  ]
]
```

Here you can see that _array_ passes the whole context given to it to each of its arguments. (Reminder: We're using "name" and "population" as elem shortcuts here.)

Most functions work like this. Only elem is the "context-cutdowner", so to say.

#### object

You can also use _object_ to create objects, with arguments alternating keys and values.
```
> cat test.json | jql '(object
                            "names" ("countries" ((keys) ("name")))
                            "populations" ("countries" ((array 0 0 1) ("population"))))'
{
  "names": [
    "Poland",
    "United States",
    "Germany"
  ],
  "populations": [
    38000000,
    38000000,
    327000000
  ]
}
```

Now we're done with the **core** functionality of jql. The stuff so far will probably suffice for most use cases and even very complex data structures.

However, here come more functions:

### String manipulation 🎻

#### join
If you ever need to _join_ an array of expressions into a string, _join_'s the mate you're looking for! _join_ will also stringify anything it meets.
Without separator:
```
> cat test.json | jql '("countries" ((keys) (join (array ("name") ("population") ("european")))))'
[
  "Poland3.8e+07true",
  "United States3.27e+08false",
  "Germany8.3e+07true"
]
```

With separator:
```
> cat test.json | jql '("countries" ((keys) (join (array ("name") ("population") ("european")) ", ")))'
[
  "Poland, 3.8e+07, true",
  "United States, 3.27e+08, false",
  "Germany, 8.3e+07, true"
]
```

#### sprintf
Whenever I learn a new language, I feel much more comfortable when I know there's a sprintf function. (and how to use it)

Anyways, here you go, the syntax is the same as that of the go standard library [fmt.Sprintf](https://golang.org/pkg/fmt/) function:
```
> cat test.json | jql '("countries" ((keys) (sprintf "%s population: %.0f" ("name") ("population"))))'
[
  "Poland population: 38000000",
  "United States population: 327000000",
  "Germany population: 83000000"
]
```
Hope you're feeling comfortable 🛋 now :)

### error
There's a little helper function - _error_ - for those times when you're debugging your queries.

It's an expression which errors on evaluation and gives you a stack trace. It can also evaluate and print any expression you want in it's context.

```
> cat test.json | jql '("countries" ((keys) (sprintf "%s population: %.0f" ("name") (error "test message"))))'
2019/12/26 00:17:04 error getting expression value for object: couldn't get transformed value for field countries with value [map[eu_since:2004 european:true name:Poland population:3.8e+07] map[european:false name:United States population:3.27e+08] map[eu_since:1993 european:true name:Germany population:8.3e+07]]: couldn't get element using position at array index 0: couldn't get transformed value for index 0 with value map[eu_since:2004 european:true name:Poland population:3.8e+07]: couldn't evaluate sprintf argument with index 1: Message: test message
goroutine 1 [running]:
runtime/debug.Stack(0xc0000723f0, 0x1114e00, 0xc0000744e0)
        /usr/local/Cellar/go/1.13.4/libexec/src/runtime/debug/stack.go:24 +0x9d
github.com/cube2222/jql/jql.Error.Get(0x1164680, 0xc0000723f0, 0x1114e00, 0xc0000744e0, 0x110c5c0, 0xc000072440, 0x0, 0x0)
        /Users/jakub/Projects/Go/src/github.com/cube2222/jql/jql/functions.go:715 +0x102
github.com/cube2222/jql/jql.Sprintf.Get(0x1164680, 0xc0000723b0, 0xc000074460, 0x2, 0x2, 0x1114e00, 0xc0000744e0, 0x100b136, 0xc000072490, 0x10, ...)
        /Users/jakub/Projects/Go/src/github.com/cube2222/jql/jql/functions.go:299 +0x16a
github.com/cube2222/jql/jql.GetElement(0x110bb80, 0x12479e0, 0x1109ec0, 0xc00008a1e0, 0x1164a20, 0xc000074480, 0x122c560, 0x128b6d0, 0x0, 0x100ba08)
        /Users/jakub/Projects/Go/src/github.com/cube2222/jql/jql/functions.go:101 +0xa96
github.com/cube2222/jql/jql.GetElement(0x1109ec0, 0xc00008a220, 0x1109ec0, 0xc00008a1e0, 0x1164a20, 0xc000074480, 0x0, 0x15fffff, 0xc0000afba0, 0x194)
        /Users/jakub/Projects/Go/src/github.com/cube2222/jql/jql/functions.go:71 +0x46d
github.com/cube2222/jql/jql.Element.Get(0x1164980, 0x12470c8, 0x1164a20, 0xc000074480, 0x1109ec0, 0xc00008a1e0, 0xc0000c6018, 0x1ee, 0xc0000c6000, 0x103301c)
        /Users/jakub/Projects/Go/src/github.com/cube2222/jql/jql/functions.go:135 +0x12f
github.com/cube2222/jql/jql.GetElement(0x110c5c0, 0xc0000722b0, 0x1114e00, 0xc0000744b0, 0x11648c0, 0xc00008a160, 0x0, 0x10, 0xc0000724b0, 0xa)
        /Users/jakub/Projects/Go/src/github.com/cube2222/jql/jql/functions.go:118 +0x862
github.com/cube2222/jql/jql.Element.Get(0x1164680, 0xc0000723a0, 0x11648c0, 0xc00008a160, 0x1114e00, 0xc0000744b0, 0x10, 0x1114860, 0x1, 0xc0000724b0)
        /Users/jakub/Projects/Go/src/github.com/cube2222/jql/jql/functions.go:135 +0x12f
main.main()
        /Users/jakub/Projects/Go/src/github.com/cube2222/jql/main.go:39 +0x2b4
```

### Logic
Yey, back to the basics!

You've got _eq_, _lt_ and _gt_, all working as you'd probably expect:
```
> cat test.json | jql '(eq "test" "test")'
true
> cat test.json | jql '(eq "test" "test2")'
false
> cat test.json | jql '(lt "a" "b")'
true
> cat test.json | jql '(lt "b" "a")'
false
> cat test.json | jql '(gt 5 4)'
true
```
In case you're wondering, _eq_ does a [reflect.DeepEqual](https://golang.org/pkg/reflect/#DeepEqual) on both arguments.

You've also got _and_, _or_, _not_, to cover your back when tackling those primal and primitive (some would say **fundamental**) problems you may encounter:
```
> cat test.json | jql '(and true true true)'
true
> cat test.json | jql '(and true true false)'
false
> cat test.json | jql '(and true true null)'
false
> cat test.json | jql '(or true true false)'
true
> cat test.json | jql '(or)'
false
> cat test.json | jql '(and)'
true
> cat test.json | jql '(not true)'
false
> cat test.json | jql '(not false)'
true
> cat test.json | jql '(not null)'
true
> cat test.json | jql '(not (array false))'
false
```
_and_ and _or_ are both lazy.

#### Truthiness

This brings us to the topic of truthiness. What does _and_ consider to be "true"? Well, it's quite simple actually.

* null **is not** truthy.
* false **is not** truthy.
* anything else **is** truthy.

#### ifte

ifte sounds kinda fluffy, but unfortunately it only stands for If Then Else.
```
> cat test.json | jql '(ifte true "true" "false")'
"true"
```

It's lazy too. If it weren't, this would error:
```
l> cat test.json | jql '(ifte true "true" (error ":("))'
"true"
```

Fluffy and lazy. Like a cat. Who doesn't like cats? Who doesn't like ifte? 🐈

#### filter 🍰
Sometimes you want just part of the cake, the part with no \<insert disliked fruit here\>.

I've got no data set on cakes though, so let's get back to our beloved countries:
```
> cat test.json | jql '("countries" (filter (gt ("population") 50000000)))'
[
  {
    "european": false,
    "name": "United States",
    "population": 327000000
  },
  {
    "eu_since": "1993",
    "european": true,
    "name": "Germany",
    "population": 83000000
  }
]
```

Also, because _null_ is not truthy, and _elem_ returns _false_ if it encounters missing fields or indices, you can use _filter_ to get rid of wrong-schema data:
```
> cat test.json | jql '("countries" (filter ("eu_since")))'
[
  {
    "eu_since": "2004",
    "european": true,
    "name": "Poland",
    "population": 38000000
  },
  {
    "eu_since": "1993",
    "european": true,
    "name": "Germany",
    "population": 83000000
  }
]
```

### pipe
_pipe_ is a fairly useless function because you can just use a bash pipe. But if for some reason you want to save cpu cycles:
```
> cat test.json | jql '(pipe
                           ("countries")
                           ((range 2))
                           ((keys) ("name")))'
[
  "Poland",
  "United States"
]
```
is equal to
```
> cat test.json | jql '("countries")' | jql '((range 2))' | jql '((keys) ("name"))'
[
  "Poland",
  "United States"
]
```

### recover
Finally, the one whose name shall not be spoken out loud. 👹 Needing him means you either encountered a bug in **jql**, or that your dataset is seriously botched.

He'll help you catch errors and panics, nullyfing the expression, leaving only void behind.
```
> cat test.json | jql '("countries" ((keys) (rec*** (ifte ("european") (id) (error "not european")))))'
[
  {
    "eu_since": "2004",
    "european": true,
    "name": "Poland",
    "population": 38000000
  },
  null,
  {
    "eu_since": "1993",
    "european": true,
    "name": "Germany",
    "population": 83000000
  }
]
```
Combine him with _filter_ and even the void will be gone!

In practice you obviously have to spell out his name, otherwise it won't work, but that's on you!

# Summary
Hope you enjoyed this **incredible** journey!

Moreover, I hope it's not the end of it! Hordes of JSON blobs still await and I hope **jql** will become your weapon of choice for dealing with them from now on! ⚔️

Issues, comments, messages, reviews, benchmarks, you name it! - all are very appreciated! 😉

# Benchmarks

This is a fairly synthetic benchmark (it's the first that came to my mind), so I won't say **jql** is faster than jq. Especially since jq has various much more advanced functionalities.

What I will say though, is that **jql** is definitely not slow and you can freely cut your way through gigabytes of data quickly, as in the following benchmarks it was 2-3 times faster than jq.

The benchbig.json file contains 2^20 repetitions of the json we've been using in the examples so far, and the benchmedium.json file contains 20.

| Command | Mean [s] | Min [s] | Max [s] | Relative |
|:---|---:|---:|---:|---:|
| `cat benchbig.json \| jql '("countries" ((keys) ("name")))' >> out.json` | 5.923 ± 0.034 | 5.890 | 6.003 | 797.90 ± 52.78 |
| `cat benchbig.json \| jq '.countries[].name' >> out.json` | 14.960 ± 0.047 | 14.906 | 15.047 | 2015.24 ± 132.94 |
| `cat benchmedium.json \| jql '("countries" ((keys) ("name")))' >> out.json` | 0.007 ± 0.000 | 0.007 | 0.010 | 1.00 |
| `cat benchmedium.json \| jq '.countries[].name' >> out.json` | 0.024 ± 0.001 | 0.022 | 0.032 | 3.23 ± 0.26 |

The benchmarks were run using [hyperfine](https://github.com/sharkdp/hyperfine) on a preheated (very loud) Macbook Pro 13 mid-2019 i7 2.8GHz 16GB 256GB.

You can generate the benchmark data with benchmarks/gen.sh.