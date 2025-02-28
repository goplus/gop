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
	"github.com/goplus/gop/tpl/token"
	"github.com/goplus/gop/tpl/types"
)

// -----------------------------------------------------------------------------
// Matcher

// Context represents the context of a matching process.
type Context struct {
}

// Matcher represents a matcher.
type Matcher interface {
	Match(src []*types.Token, ctx *Context) (n int, result any, err error)
}

// -----------------------------------------------------------------------------

// Error represents a matching error.
type Error struct {
}

func (p *Error) Error() string {
	panic("todo")
}

// -----------------------------------------------------------------------------

type gToken struct {
	tok token.Token
}

func (p *gToken) Match(src []*types.Token, ctx *Context) (n int, result any, err error) {
	if len(src) == 0 {
		return 0, nil, &Error{} // TODO(xsw): err
	}
	t := src[0]
	if t.Tok != p.tok {
		return 0, nil, &Error{} // TODO(xsw): err
	}
	return 1, t, nil
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
		return 0, nil, &Error{} // TODO(xsw): err
	}
	t := src[0]
	if t.Tok != p.tok || t.Lit != p.lit {
		return 0, nil, &Error{} // TODO(xsw): err
	}
	return 1, t, nil
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
	var nMax int
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
		errMax = &Error{} // TODO(xsw): err
	}
	return nMax, nil, errMax
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
	rets := make([]any, len(p.items))
	for i, g := range p.items {
		n1, ret1, err1 := g.Match(src[n:], ctx)
		if err1 != nil {
			return n + n1, nil, err1
		}
		rets[i] = ret1
		n += n1
	}
	result = rets
	return
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
			result = rets
			return
		}
		rets = append(rets, ret1)
		n += n1
		src = src[n1:]
	}
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
			result = rets
			return
		}
		rets = append(rets, ret1)
		n += n1
	}
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
