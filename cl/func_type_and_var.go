/*
 Copyright 2021 The GoPlus Authors (goplus.org)

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package cl

import (
	"go/types"
	"strconv"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/token"
	"github.com/goplus/gox"
)

// -----------------------------------------------------------------------------

func toFuncType(pkg *gox.Package, typ *ast.FuncType) *types.Signature {
	params, variadic := toParams(pkg, typ.Params.List)
	results := toResults(pkg, typ.Results)
	return types.NewSignature(nil, params, results, variadic)
}

func toResults(pkg *gox.Package, in *ast.FieldList) *types.Tuple {
	if in == nil {
		return nil
	}
	flds := in.List
	n := len(flds)
	args := make([]*types.Var, 0, n)
	for _, fld := range flds {
		args = toParam(pkg, fld, args)
	}
	return types.NewTuple(args...)
}

func toParams(pkg *gox.Package, flds []*ast.Field) (typ *types.Tuple, variadic bool) {
	n := len(flds)
	if n == 0 {
		return nil, false
	}
	args := make([]*types.Var, 0, n)
	for _, fld := range flds {
		args = toParam(pkg, fld, args)
	}
	_, ok := flds[n-1].Type.(*ast.Ellipsis)
	return types.NewTuple(args...), ok
}

func toParam(pkg *gox.Package, fld *ast.Field, args []*types.Var) []*gox.Param {
	typ := toType(fld.Type)
	for _, name := range fld.Names {
		args = append(args, pkg.NewParam(name.Name, typ))
	}
	return args
}

func toType(t ast.Expr) types.Type {
	panic("TODO: toType")
}

// -----------------------------------------------------------------------------

func toString(l *ast.BasicLit) string {
	if l.Kind == token.STRING {
		s, err := strconv.Unquote(l.Value)
		if err == nil {
			return s
		}
	}
	panic("TODO: toString - convert ast.BasicLit to string failed")
}

// -----------------------------------------------------------------------------
