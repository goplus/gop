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

	"github.com/goplus/gop/tpl/token"
	"github.com/goplus/gop/tpl/types"
)

var (
	errMultiMismatch = errors.New("multiple mismatch")

	// ErrVarAssigned error
	ErrVarAssigned = errors.New("variable is already assigned")
)

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

// -----------------------------------------------------------------------------

// Context represents the context of a matching process.
type Context struct {
	Fset    *token.FileSet
	FileEnd token.Pos

	LastErr error
	Left    int
}

// NewContext creates a new matching context.
func NewContext(fset *token.FileSet, fileEnd token.Pos, n int) *Context {
	return &Context{
		Fset:    fset,
		FileEnd: fileEnd,
		Left:    n,
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
// Matcher

// Matcher represents a matcher.
type Matcher interface {
	Match(src []*types.Token, ctx *Context) (n int, result any, err error)
	IsList() bool
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

func (p *gToken) IsList() bool {
	return false
}

// Token: ADD, SUB, IDENT, INT, FLOAT, CHAR, STRING, etc.
func Token(tok token.Token) Matcher {
	return &gToken{tok}
}

// -----------------------------------------------------------------------------

type gLiteral struct {
	tok token.Token
	lit string
}

func (p *gLiteral) Match(src []*types.Token, ctx *Context) (n int, result any, err error) {
	if len(src) == 0 {
		return 0, nil, ctx.NewErrorf(ctx.FileEnd, "expect `%s`, but got EOF", p.lit)
	}
	t := src[0]
	if t.Tok != p.tok || t.Lit != p.lit {
		return 0, nil, ctx.NewErrorf(t.Pos, "expect `%s`, but got `%s`", p.lit, t.Lit)
	}
	return 1, t, nil
}

func (p *gLiteral) IsList() bool {
	return false
}

// Literal: "abc", 'a', 123, 1.23, etc.
func Literal(tok token.Token, lit string) Matcher {
	return &gLiteral{tok, lit}
}

// -----------------------------------------------------------------------------

type gChoice struct {
	options []Matcher
}

func (p *gChoice) Match(src []*types.Token, ctx *Context) (n int, result any, err error) {
	var nMax = -1
	var errMax error
	var multiErr = true

	for _, g := range p.options {
		if n, result, err = g.Match(src, ctx); err == nil {
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

func (p *gChoice) IsList() bool {
	return false
}

// Choice: R1 | R2 | ... | Rn
func Choice(options ...Matcher) Matcher {
	return &gChoice{options}
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

func (p *gSequence) IsList() bool {
	return true
}

// Sequence: R1 R2 ... Rn
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
