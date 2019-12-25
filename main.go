package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/cube2222/jql/jql"
	"github.com/cube2222/jql/jql/parser"
)

func main() {
	var object map[string]interface{}
	err := json.Unmarshal([]byte(`{
  "cities": [
    {
      "name": "Berlin",
      "population": 8000000
    },
    {
      "name": "Warsaw",
      "population": 4000000
    }
  ]
}`), &object)
	if err != nil {
		log.Fatal(err)
	}

	parsed := parser.Parse(os.Args[1])
	expr, err := parsed.GetExecutionExpression(parser.ExpressionConstructorContext{
		Functions: jql.Functions,
		ConstantExpression: func(value interface{}) jql.Expression {
			return jql.NewConstant(value)
		},
	})
	if err != nil {
		log.Fatalf("couldn't get execution expression from AST: %v", err)
	}

	out, err := expr.Get(object)
	if err != nil {
		log.Fatalf("error getting expression value for object: %v", err)
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "  ")
	err = e.Encode(&out)
	if err != nil {
		log.Fatal(err)
	}
}
