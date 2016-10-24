package main

import (
	"strconv"
	"strings"
	"testing"

	"qiniupkg.com/text/tpl.v1"
)

//go:generate tplgen -f static_compiler.go -g grammar
func TestCompiler(t *testing.T) {
	var stk []string
	var args tpl.Grammar

	var scanner = new(tpl.Scanner)
	grammar := `

term = factor *(('*' factor)/mul | ('/' factor)/div)

doc = term *(('+' term)/add | ('-' term)/sub)

factor =
	FLOAT/push |
	('-' factor)/neg |
	'(' doc ')' |
	(IDENT '(' doc % ','/arity ')')/call |
	'+' factor
`
	push := func(tokens []tpl.Token, g tpl.Grammar) {
		v := tokens[0].Literal
		if v == "" {
			v = scanner.Ttol(tokens[0].Kind)
		} else if tokens[0].Kind == tpl.IDENT {
			v += "/" + strconv.Itoa(args.Len())
		}
		stk = append(stk, v)
	}

	marker := func(g tpl.Grammar, mark string) tpl.Grammar {
		if mark == "arity" {
			args = g
			return g
		} else {
			return tpl.Action(g, push)
		}
	}

	compiler := StaticCompiler{
		Grammar: grammar,
		Marker:  marker,
	}
	m := compiler.Cl()

	err := m.MatchExactly([]byte(`max(1.2 + sin(.3) * 2, cos(3), pow(5, 6), 7)`), "")
	if err != nil {
		t.Fatal(err)
	}
	text := strings.Join(stk, " ")
	if text != "1.2 .3 sin/1 2 '*' '+' 3 cos/1 5 6 pow/2 7 max/4" {
		t.Fatal("MatchExactly failed: " + text)
	}

}
