package parser

import (
	"log"
	"strings"
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

func (t *Tokenizer) Lex(lval *yySymType) int {
	if t.index == len(t.queryText) {
		return -1
	}
	for unicode.IsSpace(rune(t.queryText[t.index])) {
		t.index++
	}

	ch := t.queryText[t.index]

	log.Println(t.index)

	if ch >= 'a' && ch <= 'z' {
		endOfToken := strings.Index(t.queryText[t.index:], ")")
		newIndex := t.index + endOfToken
		if nextSpace := strings.Index(t.queryText[t.index:], " "); nextSpace != -1 && nextSpace < endOfToken {
			endOfToken = nextSpace
			newIndex = t.index + endOfToken + 1
		}
		lval.bytes = []byte(t.queryText[t.index : t.index+endOfToken])
		t.index = newIndex
		return ID
	}

	switch ch {
	case '"':
		endOfToken := strings.Index(t.queryText[t.index+1:], `"`)
		lval.bytes = []byte(t.queryText[t.index+1 : t.index+1+endOfToken])
		t.index = t.index + endOfToken + 1 + 1
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
