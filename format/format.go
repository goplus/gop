/*
 * Copyright (c) 2021 The XGo Authors (xgo.dev). All rights reserved.
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

// Package format implements standard formatting of XGo source.
//
// Note that formatting of Go source code changes over time, so tools relying on
// consistent formatting should execute a specific version of the gofmt binary
// instead of using this package. That way, the formatting will be stable, and
// the tools won't need to be recompiled each time gofmt changes.
//
// For example, pre-submit checks that use this package directly would behave
// differently depending on what Go version each developer uses, causing the
// check to be inherently fragile.
package format

import (
	"bytes"
	"io"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/printer"
	"github.com/goplus/gop/token"
)

var config = printer.Config{Mode: printer.UseSpaces | printer.TabIndent, Tabwidth: 8}

const parserMode = parser.ParseComments

// Node formats node in canonical gofmt style and writes the result to dst.
//
// The node type must be *ast.File, *printer.CommentedNode, []ast.Decl,
// []ast.Stmt, or assignment-compatible to ast.Expr, ast.Decl, ast.Spec,
// or ast.Stmt. Node does not modify node. Imports are not sorted for
// nodes representing partial source files (for instance, if the node is
// not an *ast.File or a *printer.CommentedNode not wrapping an *ast.File).
//
// The function may return early (before the entire result is written)
// and return a formatting error, for instance due to an incorrect AST.
func Node(dst io.Writer, fset *token.FileSet, node any) error {
	// Determine if we have a complete source file (file != nil).
	var file *ast.File
	var cnode *printer.CommentedNode
	switch n := node.(type) {
	case *ast.File:
		file = n
	case *printer.CommentedNode:
		if f, ok := n.Node.(*ast.File); ok {
			file = f
			cnode = n
		}
	}

	// Sort imports if necessary.
	if file != nil && hasUnsortedImports(file) {
		// Make a copy of the AST because ast.SortImports is destructive.
		// TODO(gri) Do this more efficiently.
		var buf bytes.Buffer
		err := config.Fprint(&buf, fset, file)
		if err != nil {
			return err
		}
		file, err = parser.ParseFile(fset, "", buf.Bytes(), parserMode)
		if err != nil {
			_, err = dst.Write(buf.Bytes())
			return err
			// We should never get here. If we do, provide good diagnostic.
			// return fmt.Errorf("format.Node internal error (%s)", err)
		}
		ast.SortImports(fset, file)

		// Use new file with sorted imports.
		node = file
		if cnode != nil {
			node = &printer.CommentedNode{Node: file, Comments: cnode.Comments}
		}
	}

	return config.Fprint(dst, fset, node)
}

// Source formats src in canonical gofmt style and returns the result
// or an (I/O or syntax) error. src is expected to be a syntactically
// correct Go source file, or a list of Go declarations or statements.
//
// If src is a partial source file, the leading and trailing space of src
// is applied to the result (such that it has the same leading and trailing
// space as src), and the result is indented by the same amount as the first
// line of src containing code. Imports are not sorted for partial source files.
func Source(src []byte, class bool, filename ...string) ([]byte, error) {
	var fname string
	if filename != nil {
		fname = filename[0]
	}
	fset := token.NewFileSet()
	file, sourceAdj, indentAdj, err := parse(fset, fname, src, class, true)
	if err != nil {
		return nil, err
	}

	if sourceAdj == nil {
		// Complete source file.
		// TODO(gri) consider doing this always.
		ast.SortImports(fset, file)
	}

	return format(fset, file, sourceAdj, indentAdj, src, config)
}

func hasUnsortedImports(file *ast.File) bool {
	for _, d := range file.Decls {
		d, ok := d.(*ast.GenDecl)
		if !ok || d.Tok != token.IMPORT {
			// Not an import declaration, so we're done.
			// Imports are always first.
			return false
		}
		if d.Lparen.IsValid() {
			// For now assume all grouped imports are unsorted.
			// TODO(gri) Should check if they are sorted already.
			return true
		}
		// Ungrouped imports are sorted by default.
	}
	return false
}
