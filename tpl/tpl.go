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
	return
}

// -----------------------------------------------------------------------------

type Scanner interface {
	Scan() types.Token
	Init(file *token.File, src []byte, err scanner.ScanErrorHandler, mode scanner.ScanMode)
}

type Config struct {
	Fset             *token.FileSet
	Scanner          Scanner
	ScanMode         scanner.ScanMode
	ScanErrorHandler scanner.ScanErrorHandler
}

func (p *Compiler) Eval(filename string, src any, conf *Config) (result any, err error) {
	b, err := iox.ReadSourceLocal(filename, src)
	if err != nil {
		return
	}
	if conf == nil {
		conf = &Config{}
	}
	fset := conf.Fset
	if fset == nil {
		fset = token.NewFileSet()
	}
	s := conf.Scanner
	if s == nil {
		s = new(scanner.Scanner)
	}
	f := fset.AddFile(filename, fset.Base(), len(b))
	s.Init(f, b, conf.ScanErrorHandler, conf.ScanMode)
	toks := make([]*types.Token, 0, len(b)>>3)
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
	if n < len(toks) {
		err = ctx.NewErrorf(toks[n].Pos, "unexpected token: %v", toks[n])
	}
	return
}

// -----------------------------------------------------------------------------
