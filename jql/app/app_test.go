package app

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testJson = `{
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
}`

func BenchmarkStdLib(b *testing.B) {
	var buf bytes.Buffer
	for i := 0; i < 100000; i++ {
		buf.WriteString(testJson)
	}
	for i := 0; i < b.N; i++ {
		input := json.NewDecoder(bytes.NewReader(buf.Bytes()))
		output := json.NewEncoder(ioutil.Discard)
		query := `("countries" ((keys) ("name")))`
		if err := NewApp(query, input, output).Run(); err != nil {
			log.Fatal(err)
		}
	}
}

/*func BenchmarkStdLib2(b *testing.B) {
	var buf bytes.Buffer
	for i := 0; i < 100000; i++ {
		buf.WriteString(testJson)
	}
	for i := 0; i < b.N; i++ {
		input := json.NewDecoder(bytes.NewReader(buf.Bytes()))
		output := json.NewEncoder(ioutil.Discard)
		query := `("countries" ((keys) ("name")))`
		if err := NewApp(query, input, output).Run(); err != nil {
			log.Fatal(err)
		}
	}
}*/

func TestApp_Run(t *testing.T) {
	input := `{
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
}`

	tests := []struct {
		query  string
		output string
	}{
		{
			query: `(elem "countries")`,
			output: `[
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
]`,
		},
		{
			query: `(elem "countries" (elem 0))`,
			output: `{
  "eu_since": "2004",
  "european": true,
  "name": "Poland",
  "population": 38000000
}`,
		},
		{
			query: `(elem "countries" (elem (array 0 2)))`,
			output: `[
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
]`,
		},
		{
			query: `(elem "countries" (elem (keys) (elem "name")))`,
			output: `[
  "Poland",
  "United States",
  "Germany"
]`,
		},
		{
			query: `("countries" ((keys) ("name")))`,
			output: `[
  "Poland",
  "United States",
  "Germany"
]`,
		},
		{
			query: `("countries" ((array (array 0 (array 0 (array 0 (array 0 2)))) 1 (object "key1" 1 "key2" (array 0 (object "key1" 1 "key2" (array 0 2))))) ("population")))`,
			output: `[
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
]`,
		},
		{
			query: `("countries" ((range 1 3) ("name")))`,
			output: `[
  "United States",
  "Germany"
]`,
		},
		{
			query: `("countries" ((keys) (array ("name") ("population"))))`,
			output: `[
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
]`,
		},
		{
			query: `(object
                            "names" ("countries" ((keys) ("name")))
                            "populations" ("countries" ((array 0 0 1) ("population"))))`,
			output: `{
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
}`,
		},
		{
			query: `("countries" ((keys) (join (array ("name") ("population") ("european")))))`,
			output: `[
  "Poland3.8e+07true",
  "United States3.27e+08false",
  "Germany8.3e+07true"
]`,
		},
		{
			query: `("countries" ((keys) (join (array ("name") ("population") ("european")) ", ")))`,
			output: `[
  "Poland, 3.8e+07, true",
  "United States, 3.27e+08, false",
  "Germany, 8.3e+07, true"
]`,
		},
		{
			query: `("countries" ((keys) (sprintf "%s population: %.0f" ("name") ("population"))))`,
			output: `[
  "Poland population: 38000000",
  "United States population: 327000000",
  "Germany population: 83000000"
]`,
		},
		{
			query:  `(eq "test" "test")`,
			output: `true`,
		},
		{
			query:  `(eq "test" "test2")`,
			output: `false`,
		},
		{
			query:  `(lt "a" "b")`,
			output: `true`,
		},
		{
			query:  `(lt "b" "a")`,
			output: `false`,
		},
		{
			query:  `(gt 5 4)`,
			output: `true`,
		},
		{
			query:  `(and true true true)`,
			output: `true`,
		},
		{
			query:  `(and true true false)`,
			output: `false`,
		},
		{
			query:  `(and true true null)`,
			output: `false`,
		},
		{
			query:  `(or true true false)`,
			output: `true`,
		},
		{
			query:  `(or)`,
			output: `false`,
		},
		{
			query:  `(and)`,
			output: `true`,
		},
		{
			query:  `(or true (error "bad"))`,
			output: `true`,
		},
		{
			query:  `(and false (error "bad"))`,
			output: `false`,
		},
		{
			query:  `(not true)`,
			output: `false`,
		},
		{
			query:  `(not false)`,
			output: `true`,
		},
		{
			query:  `(not null)`,
			output: `true`,
		},
		{
			query:  `(not (array false))`,
			output: `false`,
		},
		{
			query:  `(ifte true "true" "false")`,
			output: `"true"`,
		},
		{
			query:  `(ifte true "true" (error ":("))`,
			output: `"true"`,
		},
		{
			query: `("countries" (filter (gt ("population") 50000000)))`,
			output: `[
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
]`,
		},
		{
			query: `("countries" (filter ("eu_since")))`,
			output: `[
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
]`,
		},
		{
			query: `(pipe
                           ("countries")
                           ((range 2))
                           ((keys) ("name")))`,
			output: `[
  "Poland",
  "United States"
]`,
		},
		{
			query: `("countries" ((keys) (recover (ifte ("european") (id) (error "not european")))))`,
			output: `[
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
]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			input := json.NewDecoder(strings.NewReader(input))
			var buf bytes.Buffer
			output := json.NewEncoder(&buf)
			app := &App{
				query:  tt.query,
				input:  input,
				output: output,
			}
			if err := app.Run(); err != nil {
				t.Errorf("Run() error = %v", err)
			}
			assert.JSONEq(t, tt.output, string(buf.Bytes()))
		})
	}
}
