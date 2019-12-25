%{
package parser

%}

%union {
  int int
  bytes []byte
  string string
  bool bool
  null interface{}
  expression   Expression
  sexpression *SExpression
  expressions Expressions
  query *Query
  constant *Constant
}
%token <bytes> ID
%token <string> STRING
%token <int> INTEGER
%token <bool> BOOLEAN
%token <null> NULL
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
    $$ = &Constant{Value: $1}
  }
| INTEGER
  {
    $$ = &Constant{Value: $1}
  }
| BOOLEAN
	{
		$$ = &Constant{Value: $1}
	}
| NULL
	{
		$$ = &Constant{Value: nil}
	}

sexpr:
  '(' ID args_opt ')'
  {
    $$ = &SExpression{Name: string($2), Args: $3}
  }
| '(' expression args_opt ')'
    {
      $$ = &SExpression{Name: "elem", Args: append([]Expression{$2}, $3...)}
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
