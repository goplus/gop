/*
 * Copyright (c) 2024 The XGo Authors (xgo.dev). All rights reserved.
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
	"strconv"

	"github.com/goplus/xgo/ast"
	"github.com/goplus/xgo/token"
)

// -----------------------------------------------------------------------------

// Funcs returns all `*ast.FuncDecl` nodes.
func (p NodeSet) Funcs() NodeSet {
	return p.Any().FuncDecl__0().Cache()
}

// -----------------------------------------------------------------------------

func (p NodeSet) UnquotedString__1(exactly bool) (ret string, err error) {
	item, err := p.CollectOne__1(exactly)
	if err != nil {
		return
	}
	if lit, ok := item.Obj().(*ast.BasicLit); ok && lit.Kind == token.STRING {
		return strconv.Unquote(lit.Value)
	}
	return "", ErrUnexpectedNode
}

func (p NodeSet) UnquotedString__0() (ret string, err error) {
	return p.UnquotedString__1(false)
}

// -----------------------------------------------------------------------------

func (p NodeSet) UnquotedStringElts__1(exactly bool) (ret []string, err error) {
	item, err := p.CollectOne__1(exactly)
	if err != nil {
		return
	}
	if lit, ok := item.Obj().(*ast.CompositeLit); ok {
		ret = make([]string, len(lit.Elts))
		for i, elt := range lit.Elts {
			if lit, ok := elt.(*ast.BasicLit); ok && lit.Kind == token.STRING {
				if ret[i], err = strconv.Unquote(lit.Value); err != nil {
					return
				}
			} else {
				return nil, ErrUnexpectedNode
			}
		}
		return
	}
	return nil, ErrUnexpectedNode
}

func (p NodeSet) UnquotedStringElts__0() (ret []string, err error) {
	return p.UnquotedStringElts__1(false)
}

// -----------------------------------------------------------------------------

func (p NodeSet) Positions__1(exactly bool) (ret []token.Pos, err error) {
	item, err := p.CollectOne__1(exactly)
	if err != nil {
		return
	}
	switch o := item.Obj().(type) {
	case *ast.CompositeLit:
		return []token.Pos{o.Pos(), o.End(), o.Lbrace, o.Rbrace}, nil
	case ast.Node:
		return []token.Pos{o.Pos(), o.End()}, nil
	}
	return nil, ErrUnexpectedNode
}

func (p NodeSet) Positions__0() (ret []token.Pos, err error) {
	return p.Positions__1(false)
}

// -----------------------------------------------------------------------------

func (p NodeSet) EltLen__1(exactly bool) (ret int, err error) {
	item, err := p.CollectOne__1(exactly)
	if err != nil {
		return
	}
	if lit, ok := item.Obj().(*ast.CompositeLit); ok {
		return len(lit.Elts), nil
	}
	return 0, ErrUnexpectedNode
}

func (p NodeSet) EltLen__0() (ret int, err error) {
	return p.EltLen__1(false)
}

// -----------------------------------------------------------------------------

func (p NodeSet) Ident__1(exactly bool) (ret string, err error) {
	item, err := p.CollectOne__1(exactly)
	if err != nil {
		return
	}
	if ident, ok := item.Obj().(*ast.Ident); ok {
		return ident.Name, nil
	}
	return "", ErrUnexpectedNode
}

func (p NodeSet) Ident__0() (ret string, err error) {
	return p.Ident__1(false)
}

// -----------------------------------------------------------------------------

func getElt(elts []ast.Expr, name string) (ast.Expr, bool) {
	for _, elt := range elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			if ident, ok := kv.Key.(*ast.Ident); ok && ident.Name == name {
				return kv.Value, true
			}
		}
	}
	return nil, false
}

// -----------------------------------------------------------------------------

// NameOf returns name of an ast node.
func NameOf(node Node) string {
	return getName(node.Obj(), false)
}

func getName(v any, useEmpty bool) string {
	switch v := v.(type) {
	case *ast.FuncDecl:
		return v.Name.Name
	case *ast.ImportSpec:
		n := v.Name
		if n == nil {
			return ""
		}
		return n.Name
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		return getName(v.X, useEmpty) + "." + v.Sel.Name
	}
	if useEmpty {
		return ""
	}
	panic("node doesn't contain the `name` property")
}

// -----------------------------------------------------------------------------

func CodeOf(fset *token.FileSet, f *ast.File, start, end token.Pos) string {
	pos := fset.Position(start)
	n := int(end - start)
	return string(f.Code[pos.Offset : pos.Offset+n])
}

// -----------------------------------------------------------------------------
