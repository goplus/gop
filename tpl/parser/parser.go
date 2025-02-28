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

package parser

import (
	"bytes"
	"errors"
	"io"
	"os"

	"github.com/goplus/gop/tpl/ast"
	"github.com/goplus/gop/tpl/scanner"
	"github.com/goplus/gop/tpl/token"
)

// -----------------------------------------------------------------------------

// A Mode value is a set of flags (or 0).
// They control the amount of source code parsed and other optional
// parser functionality.
type Mode uint

// ParseFile parses a file and returns the AST.
func ParseFile(fset *token.FileSet, filename string, src any, _ Mode) (f *ast.File, err error) {
	b, err := readSourceLocal(filename, src)
	if err != nil {
		return nil, err
	}
	var p parser
	p.init(fset, filename, b)
	f = p.parseFile()
	switch p.errors.Len() {
	case 0:
	case 1:
		err = p.errors[0]
	default:
		p.errors.Sort()
		err = p.errors
	}
	return
}

// -----------------------------------------------------------------------------

var (
	errInvalidSource = errors.New("invalid source")
)

func readSource(src any) ([]byte, error) {
	switch s := src.(type) {
	case string:
		return []byte(s), nil
	case []byte:
		return s, nil
	case *bytes.Buffer:
		// is io.Reader, but src is already available in []byte form
		if s != nil {
			return s.Bytes(), nil
		}
	case io.Reader:
		return io.ReadAll(s)
	}
	return nil, errInvalidSource
}

// If src != nil, readSource converts src to a []byte if possible;
// otherwise it returns an error. If src == nil, readSource returns
// the result of reading the file specified by filename.
func readSourceLocal(filename string, src any) ([]byte, error) {
	if src != nil {
		return readSource(src)
	}
	return os.ReadFile(filename)
}

// -----------------------------------------------------------------------------

// parser represents a parser.
type parser struct {
	scanner scanner.Scanner
	file    *token.File

	// Current token
	pos token.Pos
	tok token.Token
	lit string

	// Error handling
	errors scanner.ErrorList
}

func (p *parser) init(fset *token.FileSet, filename string, src []byte) {
	p.file = fset.AddFile(filename, -1, len(src))
	eh := func(pos token.Position, msg string) { p.errors.Add(pos, msg) }
	p.scanner.Init(p.file, src, eh, scanner.InsertSemis)
	p.next() // initialize first token
}

// next advances to the next token.
func (p *parser) next() {
	t := p.scanner.Scan()
	p.pos, p.tok, p.lit = t.Pos, t.Tok, t.Lit
}

func (p *parser) errorExpected(pos token.Pos, msg string) {
	msg = "expected " + msg
	if pos == p.pos {
		// the error happened at the current position;
		// make the error message more specific
		switch {
		case p.tok == token.SEMICOLON && p.lit == "\n":
			msg += ", found newline"
		case len(p.lit) > 0:
			// print 123 rather than 'INT', etc.
			msg += ", found " + p.lit
		default:
			msg += ", found '" + p.tok.String() + "'"
		}
	}
	p.error(pos, msg)
}

// expect consumes the current token if it matches the expected token.
// If not, it adds an error.
func (p *parser) expect(tok token.Token) token.Pos {
	pos := p.pos
	if p.tok != tok {
		p.errorExpected(pos, "'"+tok.String()+"'")
	}
	p.next() // make progress
	return pos
}

func (p *parser) error(pos token.Pos, msg string) {
	epos := p.file.Position(pos)
	p.errors.Add(epos, msg)
}

// parseFile parses a file and returns the AST.
func (p *parser) parseFile() *ast.File {
	file := &ast.File{
		FileStart: p.pos,
	}

	for p.tok != token.EOF {
		rule := p.parseRule()
		file.Decls = append(file.Decls, rule)
	}

	return file
}

func (p *parser) parseIdent() *ast.Ident {
	pos := p.pos
	name := "_"
	if p.tok == token.IDENT {
		name = p.lit
		p.next()
	} else {
		p.errorExpected(p.pos, "'IDENT'")
	}
	return &ast.Ident{NamePos: pos, Name: name}
}

// parseRule parses a rule: IDENT '=' expr ';'
func (p *parser) parseRule() *ast.Rule {
	if p.tok != token.IDENT {
		p.errorExpected(p.pos, "'IDENT'")
		return nil
	}

	name := p.parseIdent()
	tokPos := p.expect(token.ASSIGN)
	expr := p.parseExpr()
	if expr == nil {
		return nil
	}

	p.expect(token.SEMICOLON)
	return &ast.Rule{
		Name:   name,
		TokPos: tokPos,
		Expr:   expr,
	}
}

// parseExpr parses an expression: +term % '|'
func (p *parser) parseExpr() ast.Expr {
	termList := p.parseTermList()
	for p.tok != token.OR {
		return termList
	}

	options := make([]ast.Expr, 0, 4)
	options = append(options, termList)
	for p.tok == token.OR {
		p.next()
		termList = p.parseTermList()
		options = append(options, termList)
	}
	return &ast.Choice{Options: options}
}

// parseTermList parses a list of terms: +term
func (p *parser) parseTermList() ast.Expr {
	terms := make([]ast.Expr, 0, 1)
	for {
		term, ok := p.parseTerm()
		if !ok {
			break
		}
		terms = append(terms, term)
	}
	switch n := len(terms); n {
	case 1:
		return terms[0]
	case 0:
		p.error(p.pos, "expected factor")
		fallthrough // TODO(xsw): BadExpr
	default:
		return &ast.Sequence{Items: terms}
	}
}

// parseTerm parses a term: factor % '%'
func (p *parser) parseTerm() (ast.Expr, bool) {
	x, ok := p.parseFactor()
	if !ok {
		return nil, false
	}

	for p.tok == token.REM {
		opPos := p.pos
		p.next()
		y, ok := p.parseFactor()
		if !ok {
			p.error(p.pos, "expected factor")
		}
		x = &ast.BinaryExpr{
			X:     x,
			OpPos: opPos,
			Op:    token.REM,
			Y:     y,
		}
	}
	return x, true
}

// parseFactor: IDENT | CHAR | STRING | ('*' | '+' | '?') factor | '(' expr ')'
func (p *parser) parseFactor() (ast.Expr, bool) {
	switch tok := p.tok; tok {
	case token.IDENT:
		ident := &ast.Ident{
			NamePos: p.pos,
			Name:    p.lit,
		}
		p.next()
		return ident, true

	case token.CHAR, token.STRING:
		lit := &ast.BasicLit{
			ValuePos: p.pos,
			Kind:     tok,
			Value:    p.lit,
		}
		p.next()
		return lit, true

	case token.MUL, token.ADD, token.QUESTION:
		opPos := p.pos
		p.next()

		factor, ok := p.parseFactor()
		if !ok {
			p.error(p.pos, "expected factor")
		}
		ret := &ast.UnaryExpr{
			OpPos: opPos,
			Op:    tok,
			X:     factor,
		}
		return ret, true

	case token.LPAREN:
		p.next()
		expr := p.parseExpr()
		p.expect(token.RPAREN)
		return expr, true

	default:
		return nil, false
	}
}

// -----------------------------------------------------------------------------
