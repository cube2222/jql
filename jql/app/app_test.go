package app

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"
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
