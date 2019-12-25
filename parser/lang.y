%{
package parser

%}

%union {
	bytes []byte
	expression 	Expression
	sexpression *SExpression
	expressions Expressions
	query *Query
	constant *Constant
}
%token <bytes> ID STRING INTEGRAL
%token <empty> '(' ')'

%type <expression> expression
%type <constant> constant
%type <sexpression> sexpr
%type <expressions> args_opt
%type <expressions> args
%type <query> query

%start query

%%

query:
	expression
	{
		$$ = &Query{Expression: $1}
		setQuery(yylex, $$)
	}

expression:
	constant
	{
		$$ = $1
	}
| sexpr
	{
		$$ = $1
	}

constant:
	STRING
	{
		$$ = &Constant{Value: string($1)}
	}

sexpr:
	'(' ID args_opt ')'
	{
		$$ = &SExpression{Name: string($2), Args: $3}
	}

args_opt:
	{
		$$ = nil
	}
| args
	{
		$$ = $1
	}

args:
	expression
	{
		$$ = []Expression{$1}
	}
| args expression
	{
		$$ = append($$, $2)
	}
