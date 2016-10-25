package tpl

import (
	"errors"
	"fmt"
	"go/token"
)

var (
	ErrNoDoc = errors.New("no doc")
)

// -----------------------------------------------------------------------------

type AutoKwScanner struct {
	Scanner
	tokens   []string
	keywords map[string]uint
}

func (p *AutoKwScanner) Init(file *token.File, src []byte, err ScanErrorHandler, mode ScanMode) {

	p.Scanner.Init(file, src, err, mode)
	if p.keywords == nil {
		p.keywords = make(map[string]uint)
	}
}

func (p *AutoKwScanner) Scan() (t Token) {

	t = p.Scanner.Scan()
	if t.Kind == IDENT {
		if tok, ok := p.keywords[t.Literal]; ok {
			t.Kind = tok
		}
	}
	return
}

func (p *AutoKwScanner) Ttol(tok uint) (lit string) {

	if tok >= USER_TOKEN_BEGIN {
		return p.tokens[tok-USER_TOKEN_BEGIN]
	}
	return p.Scanner.Ttol(tok)
}

func (p *AutoKwScanner) Ltot(lit string) (tok uint) { // "keyword"

	if tok = p.Scanner.Ltot(lit); tok != ILLEGAL {
		return
	}
	if lit[0] == '"' {
		kw := lit[1 : len(lit)-1]
		if tok, ok := p.keywords[kw]; ok {
			return tok
		}
		tok = USER_TOKEN_BEGIN + uint(len(p.tokens))
		p.keywords[kw] = tok
		p.tokens = append(p.tokens, lit)
	}
	return
}

// -----------------------------------------------------------------------------

type CompileRet struct {
	*Matcher
	Grammars map[string]Grammar
	Vars     map[string]*GrVar
}

func (p CompileRet) EvalSub(name string, src interface{}) error {

	g, ok := p.Grammars[name]
	if !ok {
		return ErrNoGrammar
	}

	var text []byte
	switch doc := src.(type) {
	case []Token:
		if p.Init != nil { // 因为下面调 Grammar.Match，需要先主动调用 Init
			p.Init()
		}
		matched, err := g.Match(doc, p.Matcher)
		if err != nil {
			return err
		}
		if len(doc) != matched {
			return &MatchError{GrEOF, p.Matcher, doc, nil, 0}
		}
		return nil
	case []byte:
		text = doc
	case string:
		text = []byte(doc)
	default:
		return fmt.Errorf("unsupported source type: `%v`", src)
	}

	m := &Matcher{
		Grammar:  g,
		Scanner:  p.Scanner,
		ScanMode: p.ScanMode,
	}
	return m.MatchExactly(text, "")
}

// -----------------------------------------------------------------------------

func nilMarker(g Grammar, mark string) Grammar {
	return g
}

type Compiler struct {
	Grammar  []byte
	Marker   func(g Grammar, mark string) Grammar
	Init     func()
	Scanner  Tokener
	ScanMode ScanMode
}

/*
	term = factor *(
		'%' factor/list |
		"%=" factor/list0 |
		'/' IDENT/mark
		)

	expr = +(term | '!'/nil)/and

	grammar = expr % '|'/or

	doc = +((IDENT '=' grammar ';')/assign)

	factor =
		IDENT/ident |
		CHAR/gr |
		STRING/gr |
		INT/true |
		'*' factor/repeat0 |
		'+' factor/repeat1 |
		'?' factor/repeat01 |
		'~' factor/not |
		'@' factor/peek |
		'(' grammar ')'
*/
func (p *Compiler) Cl() (ret CompileRet, err error) {

	var stk []Grammar

	if p.Scanner == nil {
		p.Scanner = new(AutoKwScanner)
	}
	scanner := p.Scanner
	grammars := make(map[string]Grammar)
	vars := make(map[string]*GrVar)

	list := func(tokens []Token, g Grammar) {
		n := len(stk)
		a, b := stk[n-2], stk[n-1]
		stk[n-2] = List(a, b)
		stk = stk[:n-1]
	}

	list0 := func(tokens []Token, g Grammar) {
		n := len(stk)
		a, b := stk[n-2], stk[n-1]
		stk[n-2] = List0(a, b)
		stk = stk[:n-1]
	}

	fnMark := p.Marker
	if fnMark == nil {
		fnMark = nilMarker
	}

	mark := func(tokens []Token, g Grammar) {
		n := len(stk)
		a := stk[n-1]
		stk[n-1] = fnMark(a, tokens[0].Literal)
	}

	and := func(tokens []Token, g Grammar) {
		m := g.Len()
		if m == 1 {
			return
		}
		n := len(stk)
		stk[n-m] = And(Clone(stk[n-m:])...)
		stk = stk[:n-m+1]
	}

	or := func(tokens []Token, g Grammar) {
		m := g.Len()
		if m == 1 {
			return
		}
		n := len(stk)
		stk[n-m] = Or(Clone(stk[n-m:])...)
		stk = stk[:n-m+1]
	}

	assign := func(tokens []Token, g Grammar) {
		n := len(stk)
		a := stk[n-1]
		name := tokens[0].Literal
		if v, ok := vars[name]; ok {
			if err := v.Assign(a); err != nil {
				panic(err)
			}
		} else if _, ok := grammars[name]; ok {
			panic("grammar already exists: " + name)
		} else {
			grammars[name] = a
		}
		stk = stk[:n-1]
	}

	ident := func(tokens []Token, g Grammar) {
		name := tokens[0].Literal
		ch := name[0]
		if ch >= 'A' && ch <= 'Z' {
			tok := scanner.Ltot(name)
			if tok == ILLEGAL {
				panic("illegal token: " + name)
			}
			g = Gr(tok)
		} else {
			if g2, ok := grammars[name]; ok {
				g = Named(name, g2)
			} else if g, ok = vars[name]; !ok {
				if name == "true" {
					g = GrTrue
				} else {
					v := Var(name)
					vars[name] = v
					g = v
				}
			}
		}
		stk = append(stk, g)
	}

	gr := func(tokens []Token, g Grammar) {
		name := tokens[0].Literal
		tok := scanner.Ltot(name)
		if tok == ILLEGAL {
			panic("illegal token: " + name)
		}
		stk = append(stk, Gr(tok))
	}

	grTrue := func(tokens []Token, g Grammar) {
		if tokens[0].Literal != "1" {
			panic("illegal token: " + tokens[0].Literal)
		}
		stk = append(stk, GrTrue)
	}

	grNil := func(tokens []Token, g Grammar) {
		stk = append(stk, nil)
	}

	repeat0 := func(tokens []Token, g Grammar) {
		n := len(stk)
		a := stk[n-1]
		stk[n-1] = Repeat0(a)
	}

	repeat1 := func(tokens []Token, g Grammar) {
		n := len(stk)
		a := stk[n-1]
		stk[n-1] = Repeat1(a)
	}

	repeat01 := func(tokens []Token, g Grammar) {
		n := len(stk)
		a := stk[n-1]
		stk[n-1] = Repeat01(a)
	}

	not := func(tokens []Token, g Grammar) {
		n := len(stk)
		a := stk[n-1]
		stk[n-1] = Not(a)
	}

	peek := func(tokens []Token, g Grammar) {
		n := len(stk)
		a := stk[n-1]
		stk[n-1] = Peek(a)
	}

	factor := Var("factor")

	term := And(factor, Repeat0(Or(
		And(Gr('%'), Action(factor, list)),
		And(Gr(REM_ASSIGN), Action(factor, list0)),
		And(Gr('/'), Action(Gr(IDENT), mark)),
	)))

	expr := Action(Repeat1(Or(term, Action(Gr('!'), grNil))), and)

	grammar := Action(List(expr, Gr('|')), or)

	doc := Repeat1(
		Action(And(Gr(IDENT), Gr('='), grammar, Gr(';')), assign),
	)

	factor.Assign(Or(
		Action(Gr(IDENT), ident),
		Action(Gr(CHAR), gr),
		Action(Gr(STRING), gr),
		Action(Gr(INT), grTrue),
		And(Gr('*'), Action(factor, repeat0)),
		And(Gr('+'), Action(factor, repeat1)),
		And(Gr('?'), Action(factor, repeat01)),
		And(Gr('~'), Action(factor, not)),
		And(Gr('@'), Action(factor, peek)),
		And(Gr('('), grammar, Gr(')')),
	))

	m := &Matcher{
		Grammar:  doc,
		Scanner:  scanner,
		ScanMode: InsertSemis,
	}
	err = m.MatchExactly(p.Grammar, "")
	if err != nil {
		return
	}

	root, ok := grammars["doc"]
	if !ok {
		if e, ok := vars["doc"]; ok {
			root = e.Elem
		} else {
			err = ErrNoDoc
			return
		}
	}
	for name, v := range vars {
		if v.Elem == nil {
			err = fmt.Errorf("variable `%s` is not assigned", name)
			return
		}
	}
	ret.Matcher = &Matcher{
		Grammar:  root,
		Scanner:  scanner,
		ScanMode: p.ScanMode,
		Init:     p.Init,
	}
	ret.Grammars, ret.Vars = grammars, vars
	return
}

// -----------------------------------------------------------------------------
