package tpl

import (
	"strconv"
	"strings"
	"testing"
)

// -----------------------------------------------------------------------------

func TestCompiler(t *testing.T) {

	const grammar = `

term = factor *(('*' factor)/mul | ('/' factor)/div)

doc = term *(('+' term)/add | ('-' term)/sub)

factor =
	FLOAT/push |
	('-' factor)/neg |
	'(' doc ')' |
	(IDENT '(' doc % ','/arity ')')/call |
	'+' factor
`
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

	marker := func(g Grammar, mark string) Grammar {
		if mark == "arity" {
			args = g
			return g
		} else {
			return Action(g, push)
		}
	}

	compiler := &Compiler{
		Grammar: []byte(grammar),
		Marker:  marker,
	}
	m, err := compiler.Cl()
	if err != nil {
		t.Fatal("compiler.Cl failed:", err)
	}

	err = m.MatchExactly([]byte(`max(1.2 + sin(.3) * 2, cos(3), pow(5, 6), 7)`), "")
	if err != nil {
		t.Fatal("MatchExactly failed:\n", err)
	}
	text := strings.Join(stk, " ")
	if text != "1.2 .3 sin/1 2 '*' '+' 3 cos/1 5 6 pow/2 7 max/4" {
		t.Fatal("MatchExactly failed:", text)
	}
}

// -----------------------------------------------------------------------------
