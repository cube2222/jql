package parser

func Parse(query string) Expression {
	tokenizer := &Tokenizer{
		queryText: query,
		index:     0,
	}
	yyParse(tokenizer)

	return tokenizer.query.Expression
}
