/*
 * Copyright (c) 2021 The GoPlus Authors (goplus.org). All rights reserved.
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
	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/token"
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

func (p astPackages) Obj() interface{} {
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

func (p astPackage) Obj() interface{} {
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

func (p astFile) Obj() interface{} {
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

func (p *astDecl) Obj() interface{} {
	return p.Decl
}

// -----------------------------------------------------------------------------

type astSpec struct {
	ast.Spec
}

func (p *astSpec) ForEach(filter func(node Node) error) error {
	return nil
}

func (p *astSpec) Obj() interface{} {
	return p.Spec
}

// -----------------------------------------------------------------------------

// NameOf returns name of an ast node.
func NameOf(node Node) string {
	switch v := node.Obj().(type) {
	case *ast.FuncDecl:
		return v.Name.Name
	case *ast.ImportSpec:
		n := v.Name
		if n == nil {
			return ""
		}
		return n.Name
	default:
		panic("node doesn't contain the `name` property")
	}
}

// -----------------------------------------------------------------------------
