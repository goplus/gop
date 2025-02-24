/*
 * Copyright (c) 2024 The GoPlus Authors (goplus.org). All rights reserved.
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

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/token"
)

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
	return getName(node.Obj())
}

func getName(v interface{}) string {
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
		return getName(v.X) + "." + v.Sel.Name
	}
	panic("node doesn't contain the `name` property")
}

// -----------------------------------------------------------------------------
