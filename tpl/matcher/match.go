/*
 * Copyright (c) 2025 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package matcher

import (
	"errors"
	"fmt"
	"log"

	"github.com/goplus/gop/tpl/token"
	"github.com/goplus/gop/tpl/types"
)

var (
	// ErrVarAssigned error
	ErrVarAssigned = errors.New("variable is already assigned")

	errNoWhitespace  = errors.New("no whitespace")
	errAdjoinEmpty   = errors.New("adjoin empty")
	errMultiMismatch = errors.New("multiple mismatch")
)

// -----------------------------------------------------------------------------

type dbgFlags int

const (
	DbgFlagMatchVar dbgFlags = 1 << iota
	DbgFlagAll               = DbgFlagMatchVar
)

var (
	enableMatchVar bool
)

func SetDebug(flags dbgFlags) {
	enableMatchVar = (flags & DbgFlagMatchVar) != 0
}

// -----------------------------------------------------------------------------

// Error represents a matching error.
type Error struct {
	Fset *token.FileSet
	Pos  token.Pos
	Msg  string

	// is a runtime error
	Dyn bool
}

func (p *Error) Error() string {
	pos := p.Fset.Position(p.Pos)
	return fmt.Sprintf("%v: %s", pos, p.Msg)
}

func isDyn(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Dyn
	}
	return false
}

// RecursiveError represents a recursive error.
type RecursiveError struct {
	*Var
}

func (e RecursiveError) Error() string {
	return "recursive variable " + e.Name
}

// -----------------------------------------------------------------------------

// Context represents the context of a matching process.
type Context struct {
	Fset    *token.FileSet
	FileEnd token.Pos
	toks    []*types.Token

	Left    int
	LastErr error
}

// NewContext creates a new matching context.
func NewContext(fset *token.FileSet, fileEnd token.Pos, toks []*types.Token) *Context {
	return &Context{
		Fset:    fset,
		FileEnd: fileEnd,
		toks:    toks,
		Left:    len(toks),
	}
}

// SetLastError sets the last error.
func (p *Context) SetLastError(left int, err error) {
	if left < p.Left {
		p.Left, p.LastErr = left, err
	}
}

// NewError creates a new error.
func (p *Context) NewError(pos token.Pos, msg string) *Error {
	return &Error{p.Fset, pos, msg, false}
}

// NewErrorf creates a new error with a format string.
func (p *Context) NewErrorf(pos token.Pos, format string, args ...any) error {
	return &Error{p.Fset, pos, fmt.Sprintf(format, args...), false}
}

// -----------------------------------------------------------------------------

// MatchToken represents a matching literal.
type MatchToken struct {
	Tok token.Token
	Lit string
}

func (p *MatchToken) String() string {
	return p.Lit
}

func hasConflictToken(me token.Token, next []any) bool {
	for _, n := range next {
		switch n := n.(type) {
		case *MatchToken:
			if n.Tok == me {
				return true
			}
		case token.Token:
			if n == me {
				return true
			}
		default:
			panic("unreachable")
		}
	}
	return false
}

func hasConflictMatchToken(me *MatchToken, next []any) bool {
	for _, n := range next {
		switch n := n.(type) {
		case *MatchToken:
			if n.Tok == me.Tok && n.Lit == me.Lit {
				return true
			}
		case token.Token:
		default:
			panic("unreachable")
		}
	}
	return false
}

func hasConflictMe(me any, next []any) bool {
	switch me := me.(type) {
	case token.Token:
		return hasConflictToken(me, next)
	case *MatchToken:
		return hasConflictMatchToken(me, next)
	}
	panic("unreachable")
}

func hasConflict(me []any, next []any) bool {
	for _, m := range me {
		if hasConflictMe(m, next) {
			return true
		}
	}
	return false
}

func conflictWith(me []any, next [][]any, from int) int {
	for i, n := from, len(next); i < n; i++ {
		if hasConflict(me, next[i]) {
			return i
		}
	}
	return -1
}

// -----------------------------------------------------------------------------
// Matcher

// Matcher represents a matcher.
type Matcher interface {
	Match(src []*types.Token, ctx *Context) (n int, result any, err error)
	First(in []any) (first []any, mayEmpty bool) // can be token.Token or *MatchToken
	IsList() bool
}

// -----------------------------------------------------------------------------

type gTrue struct{}

func (p gTrue) Match(src []*types.Token, ctx *Context) (n int, result any, err error) {
	return 0, nil, nil
}

func (p gTrue) First(in []any) (first []any, mayEmpty bool) {
	return in, true
}

func (p gTrue) IsList() bool {
	return false
}

// True returns a matcher that always succeeds.
func True() Matcher {
	return gTrue{}
}

// -----------------------------------------------------------------------------

type gWS struct{}

func (p gWS) Match(src []*types.Token, ctx *Context) (n int, result any, err error) {
	if left := len(src); left > 0 {
		toks := ctx.toks
		if n := len(toks); n > left {
			last := n - left - 1
			if toks[last].End() != src[0].Pos {
				return 0, nil, nil
			}
		}
	}
	return 0, nil, errNoWhitespace
}

func (p gWS) First(in []any) (first []any, mayEmpty bool) {
	return in, true
}

func (p gWS) IsList() bool {
	return false
}

// WhiteSpace returns a matcher that matches whitespace.
func WhiteSpace() Matcher {
	return gWS{}
}

// -----------------------------------------------------------------------------

type gString byte

func (p gString) Match(src []*types.Token, ctx *Context) (n int, result any, err error) {
	if len(src) == 0 {
		return 0, nil, ctx.NewErrorf(ctx.FileEnd, "expect `%s`, but got EOF", stringType(p))
	}
	t := src[0]
	if t.Tok != token.STRING || t.Lit[0] != byte(p) {
		return 0, nil, ctx.NewErrorf(t.Pos, "expect `%s`, but got `%v`", stringType(p), t)
	}
	return 1, t, nil
}

func (p gString) First(in []any) (first []any, mayEmpty bool) {
	return append(in, token.STRING), false
}

func (p gString) IsList() bool {
	return false
}

// String returns a matcher that matches a string literal.
func String(quoteCh byte) Matcher {
	return gString(quoteCh)
}

func stringType(quoteCh gString) string {
	if quoteCh == '"' {
		return "QSTRING"
	}
	return "RAWSTRING"
}

// -----------------------------------------------------------------------------

type gToken struct {
	tok token.Token
}

func (p *gToken) Match(src []*types.Token, ctx *Context) (n int, result any, err error) {
	if len(src) == 0 {
		return 0, nil, ctx.NewErrorf(ctx.FileEnd, "expect `%s`, but got EOF", p.tok)
	}
	t := src[0]
	if t.Tok != p.tok {
		return 0, nil, ctx.NewErrorf(t.Pos, "expect `%s`, but got `%s`", p.tok, t.Tok)
	}
	return 1, t, nil
}

func (p *gToken) First(in []any) (first []any, mayEmpty bool) {
	return append(in, p.tok), false
}

func (p *gToken) IsList() bool {
	return false
}

// Token: ADD, SUB, IDENT, INT, FLOAT, CHAR, STRING, etc.
func Token(tok token.Token) Matcher {
	return &gToken{tok}
}

// -----------------------------------------------------------------------------

type gLiteral MatchToken

func (p *gLiteral) Match(src []*types.Token, ctx *Context) (n int, result any, err error) {
	if len(src) == 0 {
		return 0, nil, ctx.NewErrorf(ctx.FileEnd, "expect `%s`, but got EOF", p.Lit)
	}
	t := src[0]
	if t.Tok != p.Tok || t.Lit != p.Lit {
		return 0, nil, ctx.NewErrorf(t.Pos, "expect `%s`, but got `%v`", p.Lit, t)
	}
	return 1, t, nil
}

func (p *gLiteral) First(in []any) (first []any, mayEmpty bool) {
	return append(in, (*MatchToken)(p)), false
}

func (p *gLiteral) IsList() bool {
	return false
}

// Literal: "abc", 'a', 123, 1.23, etc.
func Literal(tok token.Token, lit string) Matcher {
	return &gLiteral{tok, lit}
}

// -----------------------------------------------------------------------------

// Choices represents a choice matcher.
type Choices struct {
	options []Matcher
	stops   []bool
}

func (p *Choices) CheckConflicts(conflict func(firsts [][]any, i, at int)) {
	options := p.options
	n := len(options)
	firsts := make([][]any, n)
	for i, g := range options {
		firsts[i], _ = g.First(nil)
	}
	stops := make([]bool, n)
	for i, me := range firsts {
		at := conflictWith(me, firsts, i+1)
		if at >= 0 {
			conflict(firsts, i, at)
		} else {
			stops[i] = true
		}
	}
	p.stops = stops
}

func (p *Choices) Match(src []*types.Token, ctx *Context) (n int, result any, err error) {
	var nMax = -1
	var errMax error
	var multiErr = true

	stops := p.stops // be set by CheckConflicts
	for i, g := range p.options {
		if n, result, err = g.Match(src, ctx); err == nil || (n > 0 && stops[i]) {
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
		errMax = errMultiMismatch
	}
	return nMax, nil, errMax
}

func (p *Choices) First(in []any) (first []any, mayEmpty bool) {
	for _, g := range p.options {
		var me bool
		if in, me = g.First(in); me {
			mayEmpty = true
		}
	}
	first = in
	return
}

func (p *Choices) IsList() bool {
	return false
}

// Choice: R1 | R2 | ... | Rn
// Should be used with CheckConflicts.
func Choice(options ...Matcher) *Choices {
	return &Choices{options, nil}
}

// -----------------------------------------------------------------------------

type gSequence struct {
	items []Matcher
}

func (p *gSequence) Match(src []*types.Token, ctx *Context) (n int, result any, err error) {
	nitems := len(p.items)
	rets := make([]any, nitems)
	for i, g := range p.items {
		n1, ret1, err1 := g.Match(src[n:], ctx)
		if err1 != nil {
			if isDyn(err1) {
				err = err1
			} else {
				return n + n1, nil, err1
			}
		}
		rets[i] = ret1
		n += n1
	}
	result = rets
	return
}

func (p *gSequence) First(in []any) (first []any, mayEmpty bool) {
	for _, g := range p.items {
		if in, mayEmpty = g.First(in); !mayEmpty {
			break
		}
	}
	first = in
	return
}

func (p *gSequence) IsList() bool {
	return true
}

// Sequence: R1 R2 ... Rn
// Should length of items > 0
func Sequence(items ...Matcher) Matcher {
	return &gSequence{items}
}

// -----------------------------------------------------------------------------

type gRepeat0 struct {
	r Matcher
}

func (p *gRepeat0) Match(src []*types.Token, ctx *Context) (n int, result any, err error) {
	g := p.r
	rets := make([]any, 0, 2)
	for {
		n1, ret1, err1 := g.Match(src, ctx)
		if err1 != nil {
			if isDyn(err1) {
				err = err1
			} else {
				ctx.SetLastError(len(src)-n1, err1)
				result = rets
				return
			}
		}
		rets = append(rets, ret1)
		n += n1
		src = src[n1:]
	}
}

func (p *gRepeat0) First(in []any) (first []any, mayEmpty bool) {
	first, _ = p.r.First(in)
	mayEmpty = true
	return
}

func (p *gRepeat0) IsList() bool {
	return true
}

// Repeat0: *R
func Repeat0(r Matcher) Matcher {
	return &gRepeat0{r}
}

// -----------------------------------------------------------------------------

type gRepeat1 struct {
	r Matcher
}

func (p *gRepeat1) Match(src []*types.Token, ctx *Context) (n int, result any, err error) {
	g := p.r
	n, ret0, err := g.Match(src, ctx)
	if err != nil {
		return
	}

	rets := make([]any, 1, 2)
	rets[0] = ret0
	for {
		n1, ret1, err1 := g.Match(src[n:], ctx)
		if err1 != nil {
			if isDyn(err1) {
				err = err1
			} else {
				ctx.SetLastError(len(src)-n-n1, err1)
				result = rets
				return
			}
		}
		rets = append(rets, ret1)
		n += n1
	}
}

func (p *gRepeat1) First(in []any) (first []any, mayEmpty bool) {
	return p.r.First(in)
}

func (p *gRepeat1) IsList() bool {
	return true
}

// Repeat1: +R
func Repeat1(r Matcher) Matcher {
	return &gRepeat1{r}
}

// -----------------------------------------------------------------------------

type gRepeat01 struct {
	r Matcher
}

func (p *gRepeat01) Match(src []*types.Token, ctx *Context) (n int, result any, err error) {
	n, result, err = p.r.Match(src, ctx)
	if err != nil {
		return 0, nil, nil
	}
	return
}

func (p *gRepeat01) First(in []any) (first []any, mayEmpty bool) {
	first, _ = p.r.First(in)
	mayEmpty = true
	return
}

func (p *gRepeat01) IsList() bool {
	return false
}

// Repeat01: ?R
func Repeat01(r Matcher) Matcher {
	return &gRepeat01{r}
}

// -----------------------------------------------------------------------------

// List: R1 % R2 is equivalent to R1 *(R2 R1)
func List(a, b Matcher) Matcher {
	return Sequence(a, Repeat0(Sequence(b, a)))
}

// -----------------------------------------------------------------------------

type gAdjoin struct {
	a, b Matcher
}

func (p *gAdjoin) Match(src []*types.Token, ctx *Context) (n int, result any, err error) {
	n, ret0, err := p.a.Match(src, ctx)
	if err != nil {
		return
	}
	if n == 0 {
		err = errAdjoinEmpty
		return
	}
	n1, ret1, err := p.b.Match(src[n:], ctx)
	if err != nil && !isDyn(err) {
		return
	}
	if n1 == 0 {
		err = errAdjoinEmpty
		return
	}
	if src[n-1].End() != src[n].Pos {
		err = ctx.NewError(src[n].Pos, "not adjoin")
		return
	}
	n += n1
	result = []any{ret0, ret1}
	return
}

func (p *gAdjoin) First(in []any) (first []any, mayEmpty bool) {
	first, _ = p.a.First(in)
	return
}

func (p *gAdjoin) IsList() bool {
	return true
}

// Adjoin: R1 ++ R2
func Adjoin(a, b Matcher) Matcher {
	return &gAdjoin{a, b}
}

// -----------------------------------------------------------------------------

type RetProc = func(any) any
type ListRetProc = func([]any) any

type Var struct {
	Elem Matcher
	Name string
	Pos  token.Pos

	RetProc any
}

func (p *Var) Match(src []*types.Token, ctx *Context) (n int, result any, err error) {
	g := p.Elem
	if g == nil {
		return 0, nil, ctx.NewErrorf(p.Pos, "variable `%s` not assigned", p.Name)
	}
	if enableMatchVar && len(src) > 0 {
		log.Println("==> Match", p.Name, src[0])
	}
	n, result, err = g.Match(src, ctx)
	if err == nil {
		if retProc := p.RetProc; retProc != nil {
			defer func() {
				if e := recover(); e != nil {
					switch e := e.(type) {
					case *Error:
						if e.Fset == nil {
							e.Fset = ctx.Fset
						}
						err = e
					case string:
						err = &Error{
							Fset: ctx.Fset,
							Pos:  src[0].Pos,
							Msg:  e,
							Dyn:  true,
						}
					default:
						err = e.(error)
					}
				}
			}()
			if g.IsList() {
				result = retProc.(ListRetProc)(result.([]any))
			} else {
				result = retProc.(RetProc)(result)
			}
		}
	} else if err == errMultiMismatch {
		var posErr token.Pos
		var tokErr any
		if len(src) > 0 {
			posErr, tokErr = src[0].Pos, src[0]
		} else {
			posErr, tokErr = ctx.FileEnd, "EOF"
		}
		err = ctx.NewErrorf(posErr, "expect `%s`, but got `%s`", p.Name, tokErr)
	}
	return
}

func (p *Var) First(in []any) (first []any, mayEmpty bool) {
	elem := p.Elem
	if elem != nil {
		p.Elem = nil // to stop recursion
		first, mayEmpty = elem.First(in)
		p.Elem = elem
	} else {
		panic(RecursiveError{p})
	}
	return
}

func (p *Var) IsList() bool {
	return false
}

// Assign assigns a value to this variable.
func (p *Var) Assign(elem Matcher) error {
	if p.Elem != nil {
		return ErrVarAssigned
	}
	p.Elem = elem
	return nil
}

// NewVar creates a new Var instance.
func NewVar(pos token.Pos, name string) *Var {
	return &Var{Pos: pos, Name: name}
}

// -----------------------------------------------------------------------------
