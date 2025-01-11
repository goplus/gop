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

package formatutil

import (
	"github.com/goplus/gop/format"
	"github.com/goplus/gop/scanner"
	"github.com/goplus/gop/token"
)

// RearrangeFuncs rearranges functions in src.
func RearrangeFuncs(src []byte, filename ...string) ([]byte, error) {
	var fname string
	if filename != nil {
		fname = filename[0]
	}

	fset := token.NewFileSet()
	base := fset.Base()
	f := fset.AddFile(fname, base, len(src))

	var s scanner.Scanner
	s.Init(f, src, nil, scanner.ScanComments)
	stmts := splitStmts(&s)
	first := firstNonDecl(stmts)
	if first < 0 { // no non-decl stmt
		return src, nil
	}

	return src, nil
}

type aWord struct {
	pos token.Pos
	tok token.Token
}

type aStmt struct {
	words []aWord
}

func (s aStmt) tok() (tok token.Token, at int) {
	for i, w := range s.words {
		if w.tok != token.COMMENT {
			return w.tok, i
		}
	}
	return s.words[0].tok, 0
}

func splitStmts(s *scanner.Scanner) (stmts []aStmt) {
	var level int
	var stmt aStmt
	for {
		pos, tok, _ := s.Scan()
		if tok == token.EOF {
			if len(stmt.words) > 0 {
				stmts = append(stmts, stmt)
			}
			return
		}
		stmt.words = append(stmt.words, aWord{pos, tok})
		switch tok {
		case token.LBRACE:
			level++
		case token.RBRACE:
			level--
		}
		if tok == token.SEMICOLON && level == 0 {
			stmts = append(stmts, stmt)
			stmt = aStmt{}
			continue
		}
	}
}

func firstNonDecl(stmts []aStmt) int {
	for i, stmt := range stmts {
		switch tok, at := stmt.tok(); tok {
		case token.CONST, token.TYPE, token.VAR:
			continue
		case token.FUNC:
			if isFuncDecl(stmt.words[at+1:]) {
				continue
			}
		}
		return i
	}
	return -1
}

func isFuncDecl(words []aWord) bool {
	if startWith(words, token.LPAREN) { // func (
		words = seekAfter(words[1:], token.RPAREN, token.LPAREN) // func (...)
		if startWith(words, token.LBRACE) {                      // func (...) {
			return false
		}
	}
	return true
}

func seekAfter(words []aWord, tokR, tokL token.Token) []aWord {
	level := 0
	for i, w := range words {
		switch w.tok {
		case tokR:
			if level == 0 {
				return words[i+1:]
			}
			level--
		case tokL:
			level++
		}
	}
	return nil
}

func startWith(words []aWord, tok token.Token) bool {
	for _, w := range words {
		switch w.tok {
		case token.COMMENT:
			continue
		case tok:
			return true
		}
		break
	}
	return false
}

// SourceEx formats src in canonical gopfmt style and returns the result
// or an (I/O or syntax) error. src is expected to be a syntactically
// correct Go+ source file, or a list of Go+ declarations or statements.
//
// If src is a partial source file, the leading and trailing space of src
// is applied to the result (such that it has the same leading and trailing
// space as src), and the result is indented by the same amount as the first
// line of src containing code. Imports are not sorted for partial source files.
func SourceEx(src []byte, class bool, filename ...string) (formatted []byte, err error) {
	formatted, err = format.Source(src, class, filename...)
	if err == nil {
		return
	}
	src, err = RearrangeFuncs(src, filename...)
	if err != nil {
		return
	}
	return format.Source(src, class, filename...)
}
