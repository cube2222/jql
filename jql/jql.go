package jql

type Expression interface {
	Get(interface{}) (interface{}, error)
}
