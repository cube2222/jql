# jql
Hey there!

You're probably here cause you're fed up with other json query processors being too complicated to use for anything surpassing simple single field selection.

Well, at least that's what led me here. And that's why I've written jql, a json query processor with an even more cryptic - lispy - syntax (Maybe not that cryptic after all? I quite like it :) )

jql aims to be friendly, helpful, pat you on the back when you're faced with a monstrosity of a JSON blob. Help you mold it into something useful, step by step.

Ok, let's check it out now, but first things first, you have to install it:
```
go get github.com/cube2222/jql
```

Ok. Done.

Let's check out a few _simple_ examples. (remember? That's explicitly *not* why we're here. But it aids understanding of the more complex examples, so stay with me just a little bit longer!)

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
Let's break down what happened here. First we took the "countries" field. elem takes an additional argument, which says "how to transform the element", it's also an expression. Here we say we want to take the first element of the countries array. The default function is _id_, which stands for identity.

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

elem can work with single strings, single integers, arrays of those, and objects with them as values (but we won't cover those now to keep things simple, please refer to the documentation of elem).

---
A little showcase:
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
Don't do this.
---

What if we want to get all the country names? A new friend - keys - can help us here.
```
cat test.json | jql '(elem "countries" (elem (keys) (elem "name")))'
[
  "Poland",
  "United States",
  "Germany"
]
```

It returns an array of all the keys of the given collection. Fields and Indices for Objects and Arrays respectively.

---
Now we have to understand a very important mechanism, underlying jql. All functions operate in the context of the JSON we're operating on.

Some functions, like elem, will cut down the context for expressions it evaluates. The first argument - which should evaluate to the positions we need - gets evaluated in the context of the entire array, that's why _keys_ returns all the indices. The second one on the other hand, operates in the context of a single element.

In theory we're really just creating and composing a big function - pipeline, so to say - which gets applied to our JSON blob.

This may sound complicated, but I find it becomes intuitive quite quickly.
---

You can see that elem is the most used function, and in fact that is what you'll usually be using when munging data, so there's a shortcut. If you put a value in function name position, it implicitly converts it to an elem.

This way we can rewrite the previous query to be much shorter, and better match the shape of the data.
```
> cat test.json | jql '("countries" ((keys) ("name")))'
[
  "Poland",
  "United States",
  "Germany"
]
```

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

Here you can see that _array_ passes the whole context given to it to each of its arguments. (Remidner: We're using "name" and "population" as elem shortcuts here.)


TODO: Benchmarks