package app

import (
	"fmt"
	"io"

	"github.com/cube2222/jql/jql"
	"github.com/cube2222/jql/jql/functions"
	"github.com/cube2222/jql/jql/parser"
)

type Input interface {
	Decode(v interface{}) error
}

type Output interface {
	Encode(v interface{}) error
}

type App struct {
	query  string
	input  Input
	output Output
}

func NewApp(query string, input Input, output Output) *App {
	return &App{
		query:  query,
		input:  input,
		output: output,
	}
}

func (app *App) Run() error {
	parsed := parser.Parse(app.query)
	expr, err := parsed.GetExecutionExpression(parser.ExpressionConstructorContext{
		Functions: functions.Functions,
		ConstantExpression: func(value interface{}) jql.Expression {
			return jql.NewConstant(value)
		},
	})
	if err != nil {
		return fmt.Errorf("couldn't get execution expression from AST: %w", err)
	}

	for {
		var inObject interface{}
		err := app.input.Decode(&inObject)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("couldn't decode json: %w", err)
		}

		outObject, err := expr.Get(inObject)
		if err != nil {
			return fmt.Errorf("couldn't get expression value for object: %w", err)
		}

		err = app.output.Encode(outObject)
		if err != nil {
			return fmt.Errorf("couldn't encode json: %w", err)
		}
	}

	return nil
}
