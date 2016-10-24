package tpl

import (
	"bytes"
	"errors"
	"fmt"
	"go/token"
)

var (
	ErrVarNotAssigned = errors.New("variable is not assigned")
	ErrVarAssigned    = errors.New("variable is already assigned")
)

/* -----------------------------------------------------------------------------

Repeat: *G +G ?G,  Not: ~G,  Peek: @G
List: G1%G2 G1%=G2
And: G1 G2! ... Gn
Or: G1 | G2 | ... | Gn

// ---------------------------------------------------------------------------*/

const (
	lvlBase = iota
	lvlOr
	lvlAnd
	lvlList
	lvlRepeat
	lvlNot  = lvlRepeat
	lvlPeek = lvlRepeat
)

type TokenSource struct {
	File *token.File
	Src  []byte
}

type Tokener interface {
	Scan() (t Token)
	Source() (src TokenSource)
	Ttol(tok uint) (lit string)
	Ltot(lit string) (tok uint)
	Init(file *token.File, src []byte, err ScanErrorHandler, mode ScanMode)
}

type Context interface {
	Tokener() Tokener
}

type Grammar interface {
	Match(src []Token, ctx Context) (n int, err error)
	Marshal(b []byte, t Tokener, lvlParent int) []byte
	Len() int // 如果这个Grammar不是array型的，返回-1
}

func Source(doc []Token, t Tokener) (text []byte) {

	n := len(doc)
	if n == 0 {
		return nil
	}
	doc = doc[:n+1]
	s := t.Source()
	start := s.File.Offset(doc[0].Pos)
	end := s.File.Offset(doc[n].Pos)
	return s.Src[start:end]
}

func FileLine(doc []Token, t Tokener) (string, int) {

	n := len(doc)
	if n == 0 {
		return "", 0
	}
	f := t.Source().File
	return f.Name(), f.Line(doc[0].Pos)
}

func Code(g Grammar, t Tokener) string {

	return string(g.Marshal(nil, t, 0))
}

func Clone(gs []Grammar) []Grammar {

	dest := make([]Grammar, len(gs))
	for i, g := range gs {
		dest[i] = g
	}
	return dest
}

// -----------------------------------------------------------------------------

type MatchError struct {
	Grammar Grammar
	Ctx     Context
	Src     []Token
	Err     error
	IdxErr  int
}

func (p *MatchError) Error() string {
	s := p.Ctx.Tokener()
	line, text := calcErrorInfo(p.Src, p.IdxErr, s)
	if p.Err == nil {
		return fmt.Sprintf(
			"line %d: match failed: `%s` doesn't match `%s`",
			line, text, Code(p.Grammar, s))
	}
	return fmt.Sprintf(
		"line %d: match failed: `%s` doesn't match `%s`\n%v",
		line, text, Code(p.Grammar, s), p.Err)
}

func calcErrorInfo(doc []Token, idxErr int, t Tokener) (line int, text []byte) {

	n := len(doc)
	if n == 0 {
		return
	}
	doc = doc[:n+1]
	s := t.Source()
	start := s.File.Offset(doc[0].Pos)
	end := s.File.Offset(doc[n].Pos)
	errpos := s.File.Offset(doc[idxErr+1].Pos)
	all := s.Src[start:end]
	text = nextLines(all, 9)
	if len(text) < errpos-start {
		text = text[:errpos-start]
	}
	line = s.File.Line(doc[idxErr].Pos)
	return
}

func nextLines(src []byte, n int) []byte {

	from := 0
	for n > 0 {
		pos := bytes.IndexByte(src[from:], '\n')
		if pos < 0 {
			return src
		}
		from += pos + 1
		n--
	}
	return src[:from]
}

// -----------------------------------------------------------------------------

type grToken struct {
	Kind uint
}

func (p grToken) Match(src []Token, ctx Context) (n int, err error) {

	if len(src) > 0 {
		kind := src[0].Kind
		if kind == p.Kind || (kind == INT && p.Kind == FLOAT) {
			return 1, nil
		}
	}
	return 0, &MatchError{p, ctx, src, nil, 0}
}

func (p grToken) Marshal(b []byte, t Tokener, lvlParent int) []byte {

	return append(b, t.Ttol(p.Kind)...)
}

func (p grToken) Len() int {

	return -1
}

func Gr(tok uint) Grammar {

	return grToken{tok}
}

// -----------------------------------------------------------------------------

type grTrue struct {
}

func (p grTrue) Match(src []Token, ctx Context) (n int, err error) {

	return 0, nil
}

func (p grTrue) Marshal(b []byte, t Tokener, lvlParent int) []byte {

	return append(b, '1')
}

func (p grTrue) Len() int {

	return -1
}

var GrTrue Grammar = grTrue{}

// -----------------------------------------------------------------------------

type grPeek struct {
	g Grammar
}

func (p *grPeek) Match(src []Token, ctx Context) (n int, err error) {

	if _, err1 := p.g.Match(src, ctx); err1 == nil {
		return 0, nil
	}
	return 0, &MatchError{p, ctx, src, nil, 0}
}

func (p *grPeek) Marshal(b []byte, t Tokener, lvlParent int) []byte {

	b = append(b, '@')
	return p.g.Marshal(b, t, lvlPeek)
}

func (p *grPeek) Len() int {

	return -1
}

func Peek(g Grammar) Grammar {

	return &grPeek{g}
}

// -----------------------------------------------------------------------------

type grNot struct {
	g Grammar
}

func (p *grNot) Match(src []Token, ctx Context) (n int, err error) {

	if _, err1 := p.g.Match(src, ctx); err1 != nil {
		return 0, nil
	}
	return 0, &MatchError{p, ctx, src, nil, 0}
}

func (p *grNot) Marshal(b []byte, t Tokener, lvlParent int) []byte {

	b = append(b, '~')
	return p.g.Marshal(b, t, lvlNot)
}

func (p *grNot) Len() int {

	return -1
}

func Not(g Grammar) Grammar {

	return &grNot{g}
}

// -----------------------------------------------------------------------------
// EOF

type grEOF struct {
}

func (p grEOF) Match(src []Token, ctx Context) (n int, err error) {

	if len(src) != 0 {
		err = &MatchError{p, ctx, src, nil, 0}
	}
	return
}

func (p grEOF) Marshal(b []byte, t Tokener, lvlParent int) []byte {

	return append(b, 'E', 'O', 'F')
}

func (p grEOF) Len() int {

	return -1
}

var GrEOF Grammar = grEOF{}

// -----------------------------------------------------------------------------
// G1 G2 ... Gn

type grAnd struct {
	gs []Grammar
}

func (p *grAnd) Match(src []Token, ctx Context) (n int, err error) {

	failFast := false
	for _, g := range p.gs {
		if g == nil {
			failFast = true
			continue
		}
		n1, err1 := g.Match(src[n:], ctx)
		if err1 != nil {
			if failFast {
				err = &MatchError{p, ctx, src, err1, n + n1}
				panic(err)
			}
			return n + n1, err1
		}
		n += n1
	}
	return
}

func (p *grAnd) Marshal(b []byte, t Tokener, lvlParent int) []byte {

	if lvlParent > lvlAnd {
		b = append(b, '(')
	}
	for i, g := range p.gs {
		if g == nil {
			continue
		}
		if i > 0 {
			b = append(b, ' ')
		}
		b = g.Marshal(b, t, lvlAnd)
	}
	if lvlParent > lvlAnd {
		b = append(b, ')')
	}
	return b
}

func (p *grAnd) Len() int {

	return len(p.gs)
}

func And(gs ...Grammar) Grammar {

	return &grAnd{gs}
}

// -----------------------------------------------------------------------------
// G1 | G2 | ... | Gn

type grOr struct {
	gs []Grammar
}

func (p *grOr) Match(src []Token, ctx Context) (n int, err error) {

	var nMax int
	var errMax error
	var multiErr = true

	for _, g := range p.gs {
		if n, err = g.Match(src, ctx); err == nil {
			return
		}
		if n >= nMax {
			if n == nMax {
				multiErr = true
			} else {
				nMax, errMax, multiErr = n, err, false
			}
		}
	}
	if multiErr {
		errMax = &MatchError{p, ctx, src, nil, 0}
	}
	return nMax, errMax
}

func (p *grOr) Marshal(b []byte, t Tokener, lvlParent int) []byte {

	if lvlParent > lvlOr {
		b = append(b, '(')
	}
	for i, g := range p.gs {
		if i > 0 {
			b = append(b, ' ', '|', ' ')
		}
		b = g.Marshal(b, t, lvlOr)
	}
	if lvlParent > lvlOr {
		b = append(b, ')')
	}
	return b
}

func (p *grOr) Len() int {

	return len(p.gs)
}

func Or(gs ...Grammar) Grammar {

	return &grOr{gs}
}

// -----------------------------------------------------------------------------
// *G

type grRepeat0 struct {
	g   Grammar
	len int
}

func (p *grRepeat0) Len() int {

	return p.len
}

func (p *grRepeat0) Match(src []Token, ctx Context) (n int, err error) {

	g := p.g
	len := 0
	for {
		n1, err1 := g.Match(src, ctx)
		if err1 != nil {
			p.len = len
			return
		}
		len++
		n += n1
		src = src[n1:]
	}
}

func (p *grRepeat0) Marshal(b []byte, t Tokener, lvlParent int) []byte {

	b = append(b, '*')
	return p.g.Marshal(b, t, lvlRepeat)
}

func Repeat0(g Grammar) Grammar {

	return &grRepeat0{g, 0}
}

// -----------------------------------------------------------------------------
// +G

type grRepeat1 struct {
	g   Grammar
	len int
}

func (p *grRepeat1) Len() int {

	return p.len
}

func (p *grRepeat1) Match(src []Token, ctx Context) (n int, err error) {

	n, err = p.g.Match(src, ctx)
	if err != nil {
		p.len = 0
		return
	}

	r2 := grRepeat0{p.g, 0}
	n2, _ := r2.Match(src[n:], ctx)

	p.len = r2.len + 1
	return n + n2, nil
}

func (p *grRepeat1) Marshal(b []byte, t Tokener, lvlParent int) []byte {

	b = append(b, '+')
	return p.g.Marshal(b, t, lvlRepeat)
}

func Repeat1(g Grammar) *grRepeat1 {

	return &grRepeat1{g, 0}
}

// -----------------------------------------------------------------------------
// ?G

type grRepeat01 struct {
	g   Grammar
	len int
}

func (p *grRepeat01) Len() int {

	return p.len
}

func (p *grRepeat01) Match(src []Token, ctx Context) (n int, err error) {

	n, err = p.g.Match(src, ctx)
	if err != nil {
		p.len = 0
		return 0, nil
	}
	p.len = 1
	return
}

func (p *grRepeat01) Marshal(b []byte, t Tokener, lvlParent int) []byte {

	b = append(b, '?')
	return p.g.Marshal(b, t, lvlRepeat)
}

func Repeat01(g Grammar) Grammar {

	return &grRepeat01{g, 0}
}

// -----------------------------------------------------------------------------
// G1 % G2 ==> G1 *(G2 G1)
// G1 %= G2 ==> ?(G1 % G2)

type grList struct {
	a, b  Grammar
	len   int
	list0 bool
}

func (p *grList) Len() int {

	return p.len
}

func (p *grList) Match(src []Token, ctx Context) (n int, err error) {

	n, err = p.a.Match(src, ctx)
	if err != nil {
		p.len = 0
		if p.list0 {
			return 0, nil
		}
		return
	}
	src = src[n:]

	len := 1
	g := And(p.b, p.a)
	for {
		n1, err1 := g.Match(src, ctx)
		if err1 != nil {
			p.len = len
			return
		}
		len++
		n += n1
		src = src[n1:]
	}
}

func (p *grList) Marshal(b []byte, t Tokener, lvlParent int) []byte {

	if lvlParent > lvlList {
		b = append(b, '(')
	}
	b = p.a.Marshal(b, t, lvlList)
	b = append(b, ' ', '%')
	if p.list0 {
		b = append(b, '=')
	}
	b = append(b, ' ')
	b = p.b.Marshal(b, t, lvlList)
	if lvlParent > lvlList {
		b = append(b, ')')
	}
	return b
}

func List(a, b Grammar) *grList {

	return &grList{a, b, 0, false}
}

func List0(a, b Grammar) *grList {

	return &grList{a, b, 0, true}
}

// -----------------------------------------------------------------------------

type GrVar struct {
	Elem Grammar
	Name string
}

func (p *GrVar) Assign(g Grammar) error {

	if p.Elem != nil {
		return ErrVarAssigned
	}
	p.Elem = g
	return nil
}

func (p *GrVar) Len() int {

	return -1
}

func (p *GrVar) Match(src []Token, ctx Context) (n int, err error) {

	g := p.Elem
	if g == nil {
		return 0, ErrVarNotAssigned
	}

	return g.Match(src, ctx)
}

func (p *GrVar) Marshal(b []byte, t Tokener, lvlParent int) []byte {

	return append(b, p.Name...)
}

func Var(name string) *GrVar {

	return &GrVar{nil, name}
}

// -----------------------------------------------------------------------------

type GrNamed struct {
	g    Grammar
	name string
}

func (p *GrNamed) Len() int {

	return -1
}

func (p *GrNamed) Match(src []Token, ctx Context) (n int, err error) {

	return p.g.Match(src, ctx)
}

func (p *GrNamed) Marshal(b []byte, t Tokener, lvlParent int) []byte {

	return append(b, p.name...)
}

func (p *GrNamed) Assign(name string, g Grammar) {
	p.name = name
	p.g = g
}

func Named(name string, g Grammar) *GrNamed {

	return &GrNamed{g, name}
}

// -----------------------------------------------------------------------------
// G/Action

type grAction struct {
	g   Grammar
	act func(tokens []Token, g Grammar)
}

func (p *grAction) Len() int {

	return p.g.Len()
}

func (p *grAction) Match(src []Token, ctx Context) (n int, err error) {

	g := p.g
	if n, err = g.Match(src, ctx); err != nil {
		return
	}
	p.act(src[:n], g)
	return
}

func (p *grAction) Marshal(b []byte, t Tokener, lvlParent int) []byte {

	return p.g.Marshal(b, t, lvlParent)
}

func Action(g Grammar, act func(tokens []Token, g Grammar)) Grammar {

	return &grAction{g, act}
}

// -----------------------------------------------------------------------------
// G/Transaction

type grTrans struct {
	g     Grammar
	begin func() interface{}
	end   func(trans interface{}, err error)
}

func (p *grTrans) Len() int {

	return p.g.Len()
}

func (p *grTrans) Match(src []Token, ctx Context) (n int, err error) {

	trans := p.begin()
	n, err = p.g.Match(src, ctx)
	p.end(trans, err)
	return
}

func (p *grTrans) Marshal(b []byte, t Tokener, lvlParent int) []byte {

	return p.g.Marshal(b, t, lvlParent)
}

func Transaction(
	g Grammar, begin func() interface{}, end func(trans interface{}, err error)) Grammar {

	return &grTrans{g, begin, end}
}

// -----------------------------------------------------------------------------
