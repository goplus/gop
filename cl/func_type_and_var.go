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
	"log"
	"reflect"
	"strconv"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/token"
	"github.com/goplus/gox"
)

// -----------------------------------------------------------------------------

func toFuncType(ctx *loadCtx, typ *ast.FuncType) *types.Signature {
	params, variadic := toParams(ctx, typ.Params.List)
	results := toResults(ctx, typ.Results)
	return types.NewSignature(nil, params, results, variadic)
}

func toResults(ctx *loadCtx, in *ast.FieldList) *types.Tuple {
	if in == nil {
		return nil
	}
	flds := in.List
	n := len(flds)
	args := make([]*types.Var, 0, n)
	for _, fld := range flds {
		args = toParam(ctx, fld, args)
	}
	return types.NewTuple(args...)
}

func toParams(ctx *loadCtx, flds []*ast.Field) (typ *types.Tuple, variadic bool) {
	n := len(flds)
	if n == 0 {
		return nil, false
	}
	args := make([]*types.Var, 0, n)
	for _, fld := range flds {
		args = toParam(ctx, fld, args)
	}
	_, ok := flds[n-1].Type.(*ast.Ellipsis)
	return types.NewTuple(args...), ok
}

func toParam(ctx *loadCtx, fld *ast.Field, args []*types.Var) []*gox.Param {
	typ := toType(ctx, fld.Type)
	pkg := ctx.pkg
	for _, name := range fld.Names {
		args = append(args, pkg.NewParam(name.Name, typ))
	}
	return args
}

// -----------------------------------------------------------------------------

func toType(ctx *loadCtx, typ ast.Expr) types.Type {
	switch v := typ.(type) {
	case *ast.Ident:
		return toIdentType(ctx, v.Name)
	case *ast.ArrayType:
		return toArrayType(ctx, v)
	case *ast.InterfaceType:
		return toInterfaceType(ctx, v)
	case *ast.Ellipsis:
		elem := toType(ctx, v.Elt)
		return types.NewSlice(elem)
		/*	case *ast.SelectorExpr:
				return toExternalType(ctx, v)
			case *ast.StarExpr:
				elem := toType(ctx, v.X)
				return reflect.PtrTo(elem.(reflect.Type))
			case *ast.FuncType:
				return toFuncType(ctx, v)
			case *ast.StructType:
				return toStructType(ctx, v)
			case *ast.MapType:
				key := toType(ctx, v.Key)
				elem := toType(ctx, v.Value)
				return reflect.MapOf(key.(reflect.Type), elem.(reflect.Type))
			case *ast.ChanType:
				val := toType(ctx, v.Value)
				return reflect.ChanOf(toChanDir(v.Dir), val.(reflect.Type))
		*/
	}
	log.Panicln("toType: unknown -", reflect.TypeOf(typ))
	return nil
}

func toIdentType(ctx *loadCtx, ident string) types.Type {
	if _, v := ctx.pkg.Types.Scope().LookupParent(ident, token.NoPos); v != nil {
		if t, ok := v.(*types.TypeName); ok {
			return t.Type()
		}
	}
	panic("TODO: not a type")
}

func toArrayType(ctx *loadCtx, v *ast.ArrayType) types.Type {
	elem := toType(ctx, v.Elt)
	if v.Len == nil {
		return types.NewSlice(elem)
	}
	panic("TODO: array type")
}

func toInterfaceType(ctx *loadCtx, v *ast.InterfaceType) types.Type {
	methods := v.Methods.List
	if methods == nil {
		return types.NewInterfaceType(nil, nil)
	}
	panic("toInterfaceType: todo")
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
