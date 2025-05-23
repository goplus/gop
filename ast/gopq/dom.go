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

package gopq

import (
	"github.com/goplus/xgo/ast"
	"github.com/goplus/xgo/token"
)

// -----------------------------------------------------------------------------

type astPackages map[string]*ast.Package

func (p astPackages) Pos() token.Pos { return token.NoPos }
func (p astPackages) End() token.Pos { return token.NoPos }

func (p astPackages) ForEach(filter func(node Node) error) error {
	for _, pkg := range p {
		node := astPackage{pkg}
		if err := filter(node); err == ErrBreak {
			return err
		}
	}
	return nil
}

func (p astPackages) Obj() any {
	return p
}

// -----------------------------------------------------------------------------

type astPackage struct {
	*ast.Package
}

func (p astPackage) ForEach(filter func(node Node) error) error {
	for _, file := range p.Files {
		node := astFile{file}
		if err := filter(node); err == ErrBreak {
			return err
		}
	}
	return nil
}

func (p astPackage) Obj() any {
	return p.Package
}

// -----------------------------------------------------------------------------

type astFile struct {
	*ast.File
}

func (p astFile) ForEach(filter func(node Node) error) error {
	for _, decl := range p.Decls {
		node := &astDecl{decl}
		if err := filter(node); err == ErrBreak {
			return err
		}
	}
	return nil
}

func (p astFile) Obj() any {
	return p.File
}

// -----------------------------------------------------------------------------

type astDecl struct {
	ast.Decl
}

func (p *astDecl) ForEach(filter func(node Node) error) error {
	if decl, ok := p.Decl.(*ast.GenDecl); ok {
		for _, spec := range decl.Specs {
			node := &astSpec{spec}
			if err := filter(node); err == ErrBreak {
				return err
			}
		}
	}
	return nil
}

func (p *astDecl) Obj() any {
	return p.Decl
}

// -----------------------------------------------------------------------------

type astSpec struct {
	ast.Spec
}

func (p *astSpec) ForEach(filter func(node Node) error) error {
	return nil
}

func (p *astSpec) Obj() any {
	return p.Spec
}

// -----------------------------------------------------------------------------

func visitStmt(stmt ast.Stmt, filter func(node Node) error) error {
	if stmt != nil {
		return filter(&astStmt{stmt})
	}
	return nil
}

type astStmt struct {
	ast.Stmt
}

func (p *astStmt) ForEach(filter func(node Node) error) error {
	switch stmt := p.Stmt.(type) {
	case *ast.IfStmt:
		if err := visitStmt(stmt.Init, filter); err == ErrBreak {
			return err
		}
		if err := filter(&astStmt{stmt.Body}); err == ErrBreak {
			return err
		}
		if err := visitStmt(stmt.Else, filter); err == ErrBreak {
			return err
		}
	case *ast.BlockStmt:
		for _, stmt := range stmt.List {
			node := &astStmt{stmt}
			if err := filter(node); err == ErrBreak {
				return err
			}
		}
	}
	// TODO(xsw): visit other stmts
	return nil
}

func (p *astStmt) Obj() any {
	return p.Stmt
}

// -----------------------------------------------------------------------------

type astExpr struct {
	ast.Expr
}

func (p *astExpr) ForEach(filter func(node Node) error) error {
	return nil
}

func (p *astExpr) Obj() any {
	return p.Expr
}

// -----------------------------------------------------------------------------
