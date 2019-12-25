package main

import (
	"fmt"
	"strings"
)

type NodeType string

const (
	List  NodeType = "list"
	Value NodeType = "value"
)

type Ast struct {
	NodeType NodeType
	List     []Ast
	Value    string
}

func ParseQuery(query string) (Ast, error) {
	query = strings.TrimSpace(query)
	if query[0] != '(' {
		return Ast{}, fmt.Errorf("invalid query: %s", query)
	}

	for i := range query {

	}
}
