package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
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

	t := Object{
		Names: []string{
			"names",
			"populations",
		},
		ValueTransformers: []Transformer{
			Element{
				PositionSelector: "cities",
				ValueTransformer: Indices{
					IndexSelector: All,
					ValueTransformer: Element{
						PositionSelector: "name",
						ValueTransformer: Identity{},
					},
				},
			},
			Element{
				PositionSelector: "cities",
				ValueTransformer: Indices{
					IndexSelector: All,
					ValueTransformer: Element{
						PositionSelector: "population",
						ValueTransformer: Identity{},
					},
				},
			},
		},
	}

	outObject, err := t.Get(object)
	if err != nil {
		log.Fatal(err)
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "  ")
	err = e.Encode(&outObject)
	if err != nil {
		log.Fatal(err)
	}
}

type Transformer interface {
	Get(interface{}) (interface{}, error)
}

type Element struct {
	PositionSelector Transformer
	ValueTransformer Transformer
}

func (t Element) Get(arg interface{}) (interface{}, error) {
	positions, err := t.PositionSelector.Get(arg)
	if err != nil {
		return nil, fmt.Errorf("couldn't get positions to get: %w", err)
	}

	switch typed := arg.(type) {
	case []interface{}:
		indices, ok := positions.([]int)
		if !ok {
			return nil, fmt.Errorf("invalid positions for array: %v", positions)
		}

		outArray := make([]interface{}, len(indices))
		for i := range indices {
			if len(typed) <= indices[i] {
				return nil, fmt.Errorf("index %d out of range in array of length %d", indices[i], len(typed))
			}
			value := typed[indices[i]]

			outArray[indices[i]], err = t.ValueTransformer.Get(value)
			if err != nil {
				return nil, fmt.Errorf("couldn't get transformed value for index %s with value %v: %w", indices[i], value, err)
			}
		}

		return outArray, nil

	case map[string]interface{}:
		fields, ok := positions.([]string)
		if !ok {
			return nil, fmt.Errorf("invalid fields for object: %v", positions)
		}

		outObject := make(map[string]interface{}, len(fields))
		for i := range fields {
			value, ok := typed[fields[i]]
			if !ok {
				return nil, fmt.Errorf("no such field in object: %s", fields[i])
			}

			outObject[fields[i]], err = t.ValueTransformer.Get(value)
			if err != nil {
				return nil, fmt.Errorf("couldn't get transformed value for field %s with value %v: %w", fields[i], value, err)
			}
		}

		return outObject, nil

	default:
		return nil, fmt.Errorf("can only use element on array or object, used on: %s", reflect.TypeOf(arg))
	}
}

type Const struct {
	value interface{}
}

func (s Const) Get(input []interface{}) (interface{}, error) {
	return s.value, nil
}

type Keys struct {
}

func (s Keys) Get(arg interface{}) (interface{}, error) {
	switch typed := arg.(type) {
	case []interface{}:
		outIndices := make([]int, len(typed))
		for i := range typed {
			outIndices[i] = i
		}

		return outIndices, nil

	case map[string]interface{}:
		outFields := make([]string, 0, len(typed))
		for field := range typed {
			outFields = append(outFields, field)
		}

		return outFields, nil

	default:
		return nil, fmt.Errorf("can only use keys on array or object, used on: %s", reflect.TypeOf(arg))
	}
}

type Identity struct {
}

func (t Identity) Get(arg interface{}) (interface{}, error) {
	return arg, nil
}

type Array struct {
	ValueTransformers []Transformer
}

func (t Array) Get(arg interface{}) (interface{}, error) {
	outArray := make([]interface{}, len(t.ValueTransformers))
	for i := range t.ValueTransformers {
		var err error
		outArray[i], err = t.ValueTransformers[i].Get(arg)
		if err != nil {
			return nil, fmt.Errorf("couldn't construct array index %d: %w", i, err)
		}
	}

	return outArray, nil
}

type Object struct {
	Names             []string
	ValueTransformers []Transformer
}

func (t Object) Get(arg interface{}) (interface{}, error) {
	outObject := make(map[string]interface{}, len(t.ValueTransformers))
	for i, name := range t.Names {
		var err error
		outObject[name], err = t.ValueTransformers[i].Get(arg)
		if err != nil {
			return nil, fmt.Errorf("couldn't construct object field %s: %w", name, err)
		}
	}

	return outObject, nil
}
