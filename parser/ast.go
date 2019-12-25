package parser

type Query struct {
	Expression Expression
}

type Expression interface {
	IExpression()
}

type SExpression struct {
	Name string
	Args []Expression
}

func (S *SExpression) IExpression() {}

type Expressions []Expression

type Constant struct {
	Value interface{}
}

func (c *Constant) IExpression() {}
