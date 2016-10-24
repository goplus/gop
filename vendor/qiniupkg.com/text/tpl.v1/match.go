package tpl

import (
	"errors"
	"go/token"
)

var (
	ErrNoGrammar = errors.New("no grammar")
)

// -----------------------------------------------------------------------------

type Matcher struct {
	Grammar  Grammar
	Scanner  Tokener
	ScanMode ScanMode
	Init     func()
}

type TokenizeError struct {
	Pos token.Position
	Msg string
}

func (p *TokenizeError) Error() string {

	return p.Msg
}

func (p *Matcher) Code() string {

	return Code(p.Grammar, p.requireScanner())
}

func (p *Matcher) Tokenize(src []byte, fname string) (tokens []Token, err error) {

	fset := token.NewFileSet()
	file := fset.AddFile(fname, -1, len(src))

	s := p.requireScanner()

	var errTok *TokenizeError
	onError := func(pos token.Position, msg string) {
		errTok = &TokenizeError{pos, msg}
	}

	s.Init(file, src, onError, p.ScanMode)
	for {
		token := s.Scan()
		if errTok != nil {
			return tokens, errTok
		}
		tokens = append(tokens, token)
		if token.Kind == EOF { // 尽管 tokens 里面不包含 EOF，但我们希望有 EOF 过尾值（为了 Source 函数）
			tokens = tokens[:len(tokens)-1]
			break
		}
	}
	return
}

func (p *Matcher) requireScanner() Tokener {

	s := p.Scanner
	if s == nil {
		s = new(Scanner)
		p.Scanner = s
	}
	return s
}

func (p *Matcher) Tokener() Tokener {

	return p.requireScanner()
}

func (p *Matcher) Match(src []byte, fname string) (next []Token, err error) {

	g := p.Grammar
	if g == nil {
		return nil, ErrNoGrammar
	}

	tokens, err := p.Tokenize(src, fname)
	if err != nil {
		return
	}

	if p.Init != nil {
		p.Init()
	}
	n, err := g.Match(tokens, p)
	if err != nil {
		return
	}
	return tokens[n:], nil
}

func (p *Matcher) MatchExactly(src []byte, fname string) (err error) {

	next, err := p.Match(src, fname)
	if err != nil {
		return
	}

	if len(next) != 0 {
		return &MatchError{GrEOF, p, next, nil, 0}
	}
	return
}

func (p *Matcher) Eval(src string) (err error) {

	return p.MatchExactly([]byte(src), "")
}

// -----------------------------------------------------------------------------
