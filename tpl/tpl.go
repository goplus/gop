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

package tpl

import (
	"fmt"
	"io"
	"os"

	"github.com/goplus/gop/parser/iox"
	"github.com/goplus/gop/tpl/cl"
	"github.com/goplus/gop/tpl/matcher"
	"github.com/goplus/gop/tpl/parser"
	"github.com/goplus/gop/tpl/scanner"
	"github.com/goplus/gop/tpl/token"
	"github.com/goplus/gop/tpl/types"
)

// -----------------------------------------------------------------------------

// Compiler represents a TPL compiler.
type Compiler struct {
	cl.Result
	Fset *token.FileSet
}

// New creates a new TPL compiler.
func New(src any) (ret Compiler, err error) {
	return FromFile("", src, nil)
}

// FromFile creates a new TPL compiler from a file.
func FromFile(filename string, src any, fset *token.FileSet) (ret Compiler, err error) {
	if fset == nil {
		fset = token.NewFileSet()
	}
	f, err := parser.ParseFile(fset, filename, src, 0)
	if err != nil {
		return
	}
	ret.Result, err = cl.New(fset, f)
	ret.Fset = fset
	return
}

// -----------------------------------------------------------------------------

// A Token is a lexical unit returned by Scan.
type Token = types.Token

// Scanner represents a TPL scanner.
type Scanner interface {
	Scan() Token
	Init(file *token.File, src []byte, err scanner.ScanErrorHandler, mode scanner.ScanMode)
}

// Config represents a parsing configuration of [Compiler.Parse].
type Config struct {
	Scanner          Scanner
	ScanMode         scanner.ScanMode
	ScanErrorHandler scanner.ScanErrorHandler
}

// ParseExpr parses an expression.
func (p *Compiler) ParseExpr(x string, conf *Config) (result any, err error) {
	return p.ParseExprFrom("", x, conf)
}

// ParseExprFrom parses an expression from a file.
func (p *Compiler) ParseExprFrom(filename string, src any, conf *Config) (result any, err error) {
	ms, result, err := p.Match(filename, src, conf)
	if err != nil {
		return
	}
	if len(ms.Toks) == ms.N || isEOL(ms.Toks[ms.N].Tok) {
		return
	}
	t := ms.Toks[len(ms.Toks)-ms.Ctx.Left]
	err = &matcher.Error{Fset: p.Fset, Pos: t.Pos, Msg: fmt.Sprintf("unexpected token: %v", t)}
	return
}

// Parse parses a source file.
func (p *Compiler) Parse(filename string, src any, conf *Config) (result any, err error) {
	ms, result, err := p.Match(filename, src, conf)
	if err != nil {
		return
	}
	if len(ms.Toks) > ms.N {
		t := ms.Toks[len(ms.Toks)-ms.Ctx.Left]
		err = &matcher.Error{Fset: p.Fset, Pos: t.Pos, Msg: fmt.Sprintf("unexpected token: %v", t)}
	}
	return
}

// MatchState represents a matching state.
type MatchState struct {
	Toks []*Token
	Ctx  *matcher.Context
	N    int
}

// Match matches a source file.
func (p *Compiler) Match(filename string, src any, conf *Config) (ms MatchState, result any, err error) {
	b, err := iox.ReadSourceLocal(filename, src)
	if err != nil {
		return
	}
	if conf == nil {
		conf = &Config{}
	}
	s := conf.Scanner
	if s == nil {
		s = new(scanner.Scanner)
	}
	fset := p.Fset
	f := fset.AddFile(filename, fset.Base(), len(b))
	s.Init(f, b, conf.ScanErrorHandler, conf.ScanMode)
	n := (len(b) >> 3) &^ 7
	if n < 8 {
		n = 8
	}
	toks := make([]*Token, 0, n)
	for {
		t := s.Scan()
		if t.Tok == token.EOF {
			break
		}
		toks = append(toks, &t)
	}
	ms.Ctx = matcher.NewContext(fset, token.Pos(f.Base()+len(b)), len(toks))
	ms.N, result, err = p.Doc.Match(toks, ms.Ctx)
	ms.Ctx.SetLastError(len(toks)-ms.N, err)
	if err != nil {
		return
	}
	ms.Toks = toks
	return
}

// -----------------------------------------------------------------------------

func isEOL(tok token.Token) bool {
	return tok == token.SEMICOLON || tok == token.EOF
}

func isPlain(result []any) bool {
	for _, v := range result {
		if _, ok := v.(*Token); !ok {
			if a, ok := v.([]any); !ok || len(a) != 0 {
				return false
			}
		}
	}
	return true
}

// -----------------------------------------------------------------------------

func Dump(result any, omitSemi ...bool) {
	Fdump(os.Stdout, result, "", "  ", omitSemi != nil && omitSemi[0])
}

func Fdump(w io.Writer, result any, prefix, indent string, omitSemi bool) {
	switch result := result.(type) {
	case *Token:
		if result.Tok != token.SEMICOLON {
			fmt.Fprint(w, prefix, result, "\n")
		} else if !omitSemi {
			fmt.Fprint(w, prefix, ";\n")
		}
	case []any:
		if isPlain(result) {
			fmt.Print(prefix, "[")
			for i, v := range result {
				if i > 0 {
					fmt.Print(" ")
				}
				fmt.Fprint(w, v)
			}
			fmt.Print("]\n")
		} else {
			fmt.Print(prefix, "[\n")
			for _, v := range result {
				Fdump(w, v, prefix+indent, indent, omitSemi)
			}
			fmt.Print(prefix, "]\n")
		}
	default:
		panic("unexpected node")
	}
}

// -----------------------------------------------------------------------------

// List converts the matching result of (R % ",") to a flat list.
// R % "," means R *("," R)
func List(this any) []any {
	in := this.([]any)
	next := in[1].([]any)
	ret := make([]any, len(next)+1)
	ret[0] = in[0]
	for i, v := range next {
		ret[i+1] = v.([]any)[1]
	}
	return ret
}

// -----------------------------------------------------------------------------
