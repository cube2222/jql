package jql

type Expression interface {
	Get(interface{}) (interface{}, error)
}

type Constant struct {
	Value interface{}
}

func NewConstant(value interface{}) Expression {
	return &Constant{Value: value}
}

func (s Constant) Get(input interface{}) (interface{}, error) {
	return s.Value, nil
}
