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
func New(filename string, src any, fset *token.FileSet) (ret Compiler, err error) {
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
	next, result, err := p.Match(filename, src, conf)
	if err != nil {
		return
	}
	if len(next) == 0 || isEOL(next[0].Tok) {
		return
	}
	t := next[0]
	err = &matcher.Error{Fset: p.Fset, Pos: t.Pos, Msg: fmt.Sprintf("unexpected token: %v", t)}
	return
}

// Parse parses a source file.
func (p *Compiler) Parse(filename string, src any, conf *Config) (result any, err error) {
	next, result, err := p.Match(filename, src, conf)
	if err != nil {
		return
	}
	if len(next) > 0 {
		t := next[0]
		err = &matcher.Error{Fset: p.Fset, Pos: t.Pos, Msg: fmt.Sprintf("unexpected token: %v", t)}
	}
	return
}

// Match matches a source file.
func (p *Compiler) Match(filename string, src any, conf *Config) (next []*Token, result any, err error) {
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
	toks := make([]*Token, 0, len(b)>>3)
	for {
		t := s.Scan()
		if t.Tok == token.EOF {
			break
		}
		toks = append(toks, &t)
	}
	ctx := &matcher.Context{
		Fset:    fset,
		FileEnd: token.Pos(f.Base() + len(b)),
	}
	n, result, err := p.Doc.Match(toks, ctx)
	if err != nil {
		return
	}
	next = toks[n:]
	return
}

// -----------------------------------------------------------------------------

func isEOL(tok token.Token) bool {
	return tok == token.SEMICOLON || tok == token.EOF
}

// -----------------------------------------------------------------------------

func Dump(result any) {
	Fdump(os.Stdout, result, "", "  ")
}

func Fdump(w io.Writer, result any, prefix, indent string) {
	switch result := result.(type) {
	case *Token:
		fmt.Fprint(w, prefix, result, "\n")
	case []any:
		fmt.Print(prefix, "[\n")
		for _, v := range result {
			Fdump(w, v, prefix+indent, indent)
		}
		fmt.Print(prefix, "]\n")
	default:
		panic("unexpected node")
	}
}

// -----------------------------------------------------------------------------
