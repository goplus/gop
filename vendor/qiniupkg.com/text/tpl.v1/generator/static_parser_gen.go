package generator

import (
	"fmt"
	"strings"

	"qiniupkg.com/text/tpl.v1"
)

// GenStaticCode 生成StaticCompiler的Cl() method.
func GenStaticCode(source string) (string, error) {
	var stk []string
	grammars := make(map[string]string)
	vars := make(map[string]string)
	scanner := new(tpl.AutoKwScanner)
	list := func(tokens []tpl.Token, g tpl.Grammar) {
		n := len(stk)
		a, b := stk[n-2], stk[n-1]
		stk[n-2] = fmt.Sprintf("tpl.List(%s,%s)", a, b)
		stk = stk[:n-1]
	}
	list0 := func(tokens []tpl.Token, g tpl.Grammar) {
		n := len(stk)
		a, b := stk[n-2], stk[n-1]
		stk[n-2] = fmt.Sprintf("tpl.List0(%s,%s)", a, b)
		stk = stk[:n-1]
	}

	mark := func(tokens []tpl.Token, g tpl.Grammar) {
		n := len(stk)
		a := stk[n-1]
		stk[n-1] = fmt.Sprintf("cmplr.Marker(%s,\"%s\")", a, tokens[0].Literal)
	}
	and := func(tokens []tpl.Token, g tpl.Grammar) {
		m := g.Len()
		if m == 1 {
			return
		}
		n := len(stk)
		stk[n-m] = fmt.Sprintf("tpl.And(%s)", strings.Join(stk[n-m:], ","))
		stk = stk[:n-m+1]
	}
	or := func(tokens []tpl.Token, g tpl.Grammar) {
		m := g.Len()
		if m == 1 {
			return
		}
		n := len(stk)
		stk[n-m] = fmt.Sprintf("tpl.Or(%s)", strings.Join(stk[n-m:], ","))
		stk = stk[:n-m+1]
	}
	assign := func(tokens []tpl.Token, g tpl.Grammar) {
		n := len(stk)
		a := stk[n-1]
		name := tokens[0].Literal
		if _, ok := vars[name]; ok {
			vars[name] = fmt.Sprintf("%s.Assign(%s)", name, a)
		} else if _, ok := grammars[name]; ok {
			panic("grammar already exists: " + name)
		} else {
			grammars[name] = a
		}
		stk = stk[:n-1]
	}
	ident := func(tokens []tpl.Token, g tpl.Grammar) {
		name := tokens[0].Literal
		ch := name[0]
		var s string
		if ch >= 'A' && ch <= 'Z' {
			tok := scanner.Ltot(name)
			if tok == tpl.ILLEGAL {
				panic("illegal token: " + name)
			}
			if strings.HasPrefix(name, "'") {
				s = fmt.Sprintf("tpl.Gr(%s)", name)
			} else {
				s = fmt.Sprintf("tpl.Gr(tpl.%s)", name)
			}
		} else {
			_, ok := grammars[name]
			if ok {
				s = name
			} else if s, ok = vars[name]; !ok {
				if name == "true" {
					vars[name] = "tpl.GrTrue"
					s = "tpl.GrTrue"
				} else {
					vars[name] = name
					s = name
				}
			}
		}
		stk = append(stk, s)

	}
	gr := func(tokens []tpl.Token, g tpl.Grammar) {
		name := tokens[0].Literal
		tok := scanner.Ltot(name)
		if tok == tpl.ILLEGAL {
			panic("illegal token: " + name)
		}
		stk = append(stk, fmt.Sprintf("tpl.Gr(%s)", name))
	}
	grString := func(tokens []tpl.Token, g tpl.Grammar) {
		name := tokens[0].Literal
		tok := scanner.Ltot(name)
		if tok == tpl.ILLEGAL {
			panic("illegal token: " + name)
		}
		stk = append(stk, "tpl.Gr(tpl.STRING)")
	}
	grTrue := func(tokens []tpl.Token, g tpl.Grammar) {
		if tokens[0].Literal != "1" {
			panic("illegal token: " + tokens[0].Literal)
		}
		stk = append(stk, "tpl.GrTrue")
	}

	grNil := func(tokens []tpl.Token, g tpl.Grammar) {
		stk = append(stk, "nil")
	}

	repeat0 := func(tokens []tpl.Token, g tpl.Grammar) {
		n := len(stk)
		a := stk[n-1]
		stk[n-1] = fmt.Sprintf("tpl.Repeat0(%s)", a)
	}

	repeat1 := func(tokens []tpl.Token, g tpl.Grammar) {
		n := len(stk)
		a := stk[n-1]
		stk[n-1] = fmt.Sprintf("tpl.tpl.Repeat1(%s)", a)
	}

	repeat01 := func(tokens []tpl.Token, g tpl.Grammar) {
		n := len(stk)
		a := stk[n-1]
		stk[n-1] = fmt.Sprintf("tpl.Repeat01(%s)", a)
	}

	not := func(tokens []tpl.Token, g tpl.Grammar) {
		n := len(stk)
		a := stk[n-1]
		stk[n-1] = fmt.Sprintf("tpl.Not(%s)", a)
	}

	peek := func(tokens []tpl.Token, g tpl.Grammar) {
		n := len(stk)
		a := stk[n-1]
		stk[n-1] = fmt.Sprintf("tpl.Peek(%s)", a)
	}

	factor := tpl.Var("factor")

	term := tpl.And(factor, tpl.Repeat0(tpl.Or(
		tpl.And(tpl.Gr('%'), tpl.Action(factor, list)),
		tpl.And(tpl.Gr(tpl.REM_ASSIGN), tpl.Action(factor, list0)),
		tpl.And(tpl.Gr('/'), tpl.Action(tpl.Gr(tpl.IDENT), mark)),
	)))

	expr := tpl.Action(tpl.Repeat1(tpl.Or(term, tpl.Action(tpl.Gr('!'), grNil))), and)

	grammar := tpl.Action(tpl.List(expr, tpl.Gr('|')), or)

	doc := tpl.Repeat1(
		tpl.Action(tpl.And(tpl.Gr(tpl.IDENT), tpl.Gr('='), grammar, tpl.Gr(';')), assign),
	)

	factor.Assign(tpl.Or(
		tpl.Action(tpl.Gr(tpl.IDENT), ident),
		tpl.Action(tpl.Gr(tpl.CHAR), gr),
		tpl.Action(tpl.Gr(tpl.STRING), grString),
		tpl.Action(tpl.Gr(tpl.INT), grTrue),
		tpl.And(tpl.Gr('*'), tpl.Action(factor, repeat0)),
		tpl.And(tpl.Gr('+'), tpl.Action(factor, repeat1)),
		tpl.And(tpl.Gr('?'), tpl.Action(factor, repeat01)),
		tpl.And(tpl.Gr('~'), tpl.Action(factor, not)),
		tpl.And(tpl.Gr('@'), tpl.Action(factor, peek)),
		tpl.And(tpl.Gr('('), grammar, tpl.Gr(')')),
	))

	m := &tpl.Matcher{
		Grammar:  doc,
		Scanner:  scanner,
		ScanMode: tpl.InsertSemis,
	}
	err := m.MatchExactly([]byte(source), "")
	if err != nil {
		return "", err
	}
	var (
		declarations string = "var (\n"
		assignments  string
	)
	for k, v := range vars {
		declarations += fmt.Sprintf("%s = new(tpl.GrVar)\n", k)
		assignments += fmt.Sprintf("// GrVar %s\n%s\n", k, v)
	}
	for k, v := range grammars {
		declarations += fmt.Sprintf("%s = new(tpl.GrNamed)\n", k)
		assignments += fmt.Sprintf("// GrNamed %s\n%s.Assign(\"%s\", %s)\n", k, k, k, v)
	}
	declarations += ")\n"
	if _, ok := grammars["doc"]; ok {
		output := fmt.Sprintf("\n"+
			"type StaticCompiler struct {\n"+
			"	Grammar  string\n"+
			"	Marker   func(g tpl.Grammar, mark string) tpl.Grammar\n"+
			"	Init     func()\n"+
			"	Scanner  tpl.Tokener\n"+
			"	ScanMode tpl.ScanMode\n"+
			"}\n"+
			"func(cmplr *StaticCompiler)Cl()tpl.CompileRet{\n"+
			"%s"+
			"%s"+
			"ret := tpl.CompileRet{}\n"+
			"ret.Matcher = &tpl.Matcher{\n"+
			"		Grammar: doc,\n"+
			"		Scanner: cmplr.Scanner,\n"+
			"		Init: cmplr.Init,\n"+
			"	}\n"+
			"return ret}\n",
			declarations,
			assignments,
		)
		return output, nil
	}
	return "", tpl.ErrNoDoc
}
