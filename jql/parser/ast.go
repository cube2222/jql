package parser

import (
	"fmt"

	"github.com/cube2222/jql/jql"
)

type Query struct {
	Expression Expression
}

type ExpressionConstructorContext struct {
	Functions          map[string]func(...jql.Expression) (jql.Expression, error)
	ConstantExpression func(interface{}) jql.Expression
}

type Expression interface {
	IExpression()
	GetExecutionExpression(eCtx ExpressionConstructorContext) (jql.Expression, error)
}

type SExpression struct {
	Name string
	Args []Expression
}

func (e *SExpression) GetExecutionExpression(eCtx ExpressionConstructorContext) (jql.Expression, error) {
	f, ok := eCtx.Functions[e.Name]
	if !ok {
		return nil, fmt.Errorf("no such function: %s", e.Name)
	}

	arguments := make([]jql.Expression, len(e.Args))
	for i := range e.Args {
		expr, err := e.Args[i].GetExecutionExpression(eCtx)
		if err != nil {
			return nil, fmt.Errorf("couldn't get argument expression with index %d: %w", i, err)
		}

		arguments[i] = expr
	}

	expr, err := f(arguments...)
	if err != nil {
		return nil, fmt.Errorf("couldn't get expression for function %s: %w", e.Name, err)
	}

	return expr, nil
}

func (e *SExpression) IExpression() {}

type Expressions []Expression

type Constant struct {
	Value interface{}
}

func (e *Constant) IExpression() {}

func (e *Constant) GetExecutionExpression(eCtx ExpressionConstructorContext) (jql.Expression, error) {
	return eCtx.ConstantExpression(e.Value), nil
}
