package tpl

import (
	"strconv"
	"strings"
	"testing"
)

// -----------------------------------------------------------------------------

func TestMatch(t *testing.T) {

	/*	term = factor *('*' factor/mul | '/' factor/div)

		expr = term *('+' term/add | '-' term/sub)

		factor =
			FLOAT/push |
			'-' factor/neg |
			'(' expr ')' |
			(IDENT '(' expr % ','/arity ')')/call |
			'+' factor
	*/

	var stk []string
	var args Grammar

	scanner := new(Scanner)
	push := func(tokens []Token, g Grammar) {
		v := tokens[0].Literal
		if v == "" {
			v = scanner.Ttol(tokens[0].Kind)
		} else if tokens[0].Kind == IDENT {
			v += "/" + strconv.Itoa(args.Len())
		}
		stk = append(stk, v)
	}

	factor := Var("factor")

	term := And(factor, Repeat0(Or(
		Action(And(Gr('*'), factor), push),
		Action(And(Gr('/'), factor), push),
	)))

	expr := And(term, Repeat0(Or(
		Action(And(Gr('+'), term), push),
		Action(And(Gr('-'), term), push),
	)))

	args = List(expr, Gr(','))

	factor.Assign(Or(
		Action(Gr(FLOAT), push),
		Action(And(Gr('-'), factor), push),
		And(Gr('('), expr, Gr(')')),
		Action(And(Gr(IDENT), Gr('('), args, Gr(')')), push),
		And(Gr('+'), factor),
	))

	m := Matcher{Grammar: expr, Scanner: scanner}
	err := m.MatchExactly([]byte(`max(1.2 + sin(.3) * 2, cos(3), pow(5, 6), 7)`), "")
	if err != nil {
		t.Fatal("MatchExactly failed:", err)
	}
	text := strings.Join(stk, " ")
	if text != "1.2 .3 sin/1 2 '*' '+' 3 cos/1 5 6 pow/2 7 max/4" {
		t.Fatal("MatchExactly failed:", text)
	}
}

// -----------------------------------------------------------------------------
