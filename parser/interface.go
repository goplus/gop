// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains the exported entry points for invoking the parser.

package parser

import (
	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/token"
)

// A Mode value is a set of flags (or 0).
// They control the amount of source code parsed and other optional
// parser functionality.
//
type Mode uint

const (
	// PackageClauseOnly - stop parsing after package clause
	PackageClauseOnly Mode = 1 << iota
	// ImportsOnly - stop parsing after import declarations
	ImportsOnly
	// ParseComments - parse comments and add them to AST
	ParseComments
	// Trace - print a trace of parsed productions
	Trace
	// DeclarationErrors - report declaration errors
	DeclarationErrors
	// AllErrors - report all errors (not just the first 10 on different lines)
	AllErrors
)

// ParseFile parses the source code of a single Go source file and returns
// the corresponding ast.File node. The source code may be provided via
// the filename of the source file, or via the src parameter.
//
// If src != nil, ParseFile parses the source from src and the filename is
// only used when recording position information. The type of the argument
// for the src parameter must be string, []byte, or io.Reader.
// If src == nil, ParseFile parses the file specified by filename.
//
// The mode parameter controls the amount of source text parsed and other
// optional parser functionality. Position information is recorded in the
// file set fset, which must not be nil.
//
// If the source couldn't be read, the returned AST is nil and the error
// indicates the specific failure. If the source was read but syntax
// errors were found, the result is a partial AST (with ast.Bad* nodes
// representing the fragments of erroneous source code). Multiple errors
// are returned via a scanner.ErrorList which is sorted by source position.
//
func parseFile(fset *token.FileSet, filename string, src interface{}, mode Mode) (f *ast.File, err error) {
	if fset == nil {
		panic("parser.ParseFile: no token.FileSet provided (fset == nil)")
	}

	// get source
	text, err := readSource(src)
	if err != nil {
		return nil, err
	}

	var p parser
	defer func() {
		if e := recover(); e != nil {
			// resume same panic if it's not a bailout
			if _, ok := e.(bailout); !ok {
				panic(e)
			}
		}

		// set result values
		if f == nil {
			// source is not a valid Go source file - satisfy
			// ParseFile API and return a valid (but) empty
			// *ast.File
			f = &ast.File{
				Name:  new(ast.Ident),
				Scope: ast.NewScope(nil),
			}
		}
		f.Code = text

		p.errors.Sort()
		err = p.errors.Err()
	}()

	// parse source
	p.init(fset, filename, text, mode)
	f = p.parseFile()

	return
}

/*
// ParseExprFrom is a convenience function for parsing an expression.
// The arguments have the same meaning as for ParseFile, but the source must
// be a valid Go (type or value) expression. Specifically, fset must not
// be nil.
//
// If the source couldn't be read, the returned AST is nil and the error
// indicates the specific failure. If the source was read but syntax
// errors were found, the result is a partial AST (with ast.Bad* nodes
// representing the fragments of erroneous source code). Multiple errors
// are returned via a scanner.ErrorList which is sorted by source position.
//
func parseExprFrom(fset *token.FileSet, filename string, src interface{}, mode Mode) (expr ast.Expr, err error) {
	if fset == nil {
		panic("parser.ParseExprFrom: no token.FileSet provided (fset == nil)")
	}

	// get source
	text, err := readSource(src)
	if err != nil {
		return nil, err
	}

	var p parser
	defer func() {
		if e := recover(); e != nil {
			// resume same panic if it's not a bailout
			if _, ok := e.(bailout); !ok {
				panic(e)
			}
		}
		p.errors.Sort()
		err = p.errors.Err()
	}()

	// parse expr
	p.init(fset, filename, text, mode)
	// Set up pkg-level scopes to avoid nil-pointer errors.
	// This is not needed for a correct expression x as the
	// parser will be ok with a nil topScope, but be cautious
	// in case of an erroneous x.
	p.openScope()
	p.pkgScope = p.topScope
	expr = p.parseRHSOrType()
	p.closeScope()
	assert(p.topScope == nil, "unbalanced scopes")

	// If a semicolon was inserted, consume it;
	// report an error if there's more tokens.
	if p.tok == token.SEMICOLON && p.lit == "\n" {
		p.next()
	}
	p.expect(token.EOF)

	return
}

// ParseExpr is a convenience function for obtaining the AST of an expression x.
// The position information recorded in the AST is undefined. The filename used
// in error messages is the empty string.
//
// If syntax errors were found, the result is a partial AST (with ast.Bad* nodes
// representing the fragments of erroneous source code). Multiple errors are
// returned via a scanner.ErrorList which is sorted by source position.
//
func parseExpr(x string) (ast.Expr, error) {
	return parseExprFrom(token.NewFileSet(), "", x, 0)
}
*/
