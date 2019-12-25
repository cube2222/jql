package parser

import (
	"log"
	"regexp"
	"strconv"
	"unicode"
)

type Tokenizer struct {
	queryText string
	index     int
	query     *Query
}

func setQuery(tokenizer interface{}, query *Query) {
	tokenizer.(*Tokenizer).query = query
}

var integerRegexp = regexp.MustCompile("[0-9]+")
var identifierRegexp = regexp.MustCompile("[a-zA-Z][a-zA-Z0-9]*")
var stringRegexp = regexp.MustCompile(`"(?:[^"\\]|\\.)*"`)

func (t *Tokenizer) Lex(lval *yySymType) int {
	if t.index == len(t.queryText) {
		return -1
	}
	for unicode.IsSpace(rune(t.queryText[t.index])) {
		t.index++
	}

	ch := t.queryText[t.index]

	if ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') {
		indices := identifierRegexp.FindStringIndex(t.queryText[t.index:])

		lval.bytes = []byte(t.queryText[t.index : t.index+indices[1]])
		t.index += indices[1]
		return ID
	}
	if ch >= '0' && ch <= '9' {
		indices := integerRegexp.FindStringIndex(t.queryText[t.index:])

		var err error
		lval.int, err = strconv.Atoi(t.queryText[t.index : t.index+indices[1]])
		if err != nil {
			return -1
		}
		t.index += indices[1]
		return INTEGER
	}

	switch ch {
	case '"':
		indices := stringRegexp.FindStringIndex(t.queryText[t.index:])
		lval.string = t.queryText[t.index+1 : t.index+indices[1]-1]
		t.index += indices[1]
		return STRING
	case '(', ')':
		t.index++
		return int(ch)
	}

	return -1
}

func (t Tokenizer) Error(s string) {
	log.Fatalf("error at %d: %s", t.index, s)
}
