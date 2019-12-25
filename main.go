package main

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/cube2222/jql/jql"
	"github.com/cube2222/jql/jql/parser"
)

func main() {
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

	input := json.NewDecoder(os.Stdin)
	output := json.NewEncoder(os.Stdout)
	output.SetIndent("", "  ")

	for {
		var inObject map[string]interface{}
		err := input.Decode(&inObject)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("couldn't decode json: %v", err)
		}

		outObject, err := expr.Get(inObject)
		if err != nil {
			log.Fatalf("error getting expression value for object: %v", err)
		}

		err = output.Encode(&outObject)
		if err != nil {
			log.Fatalf("couldn't encode json: %v", err)
		}
	}
}
