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

func toFuncType(ctx *blockCtx, typ *ast.FuncType) *types.Signature {
	params, variadic := toParams(ctx, typ.Params.List)
	results := toResults(ctx, typ.Results)
	return types.NewSignature(nil, params, results, variadic)
}

func toResults(ctx *blockCtx, in *ast.FieldList) *types.Tuple {
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

func toParams(ctx *blockCtx, flds []*ast.Field) (typ *types.Tuple, variadic bool) {
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

func toParam(ctx *blockCtx, fld *ast.Field, args []*gox.Param) []*gox.Param {
	typ := toType(ctx, fld.Type)
	pkg := ctx.pkg
	if len(fld.Names) == 0 {
		return append(args, pkg.NewParam("", typ))
	}
	for _, name := range fld.Names {
		args = append(args, pkg.NewParam(name.Name, typ))
	}
	return args
}

// -----------------------------------------------------------------------------

func toType(ctx *blockCtx, typ ast.Expr) types.Type {
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
	case *ast.StarExpr:
		elem := toType(ctx, v.X)
		return types.NewPointer(elem)
	case *ast.MapType:
		key := toType(ctx, v.Key)
		elem := toType(ctx, v.Value)
		return types.NewMap(key, elem)
	case *ast.ChanType:
		return types.NewChan(typesChanDirs[v.Dir], toType(ctx, v.Value))
	case *ast.FuncType:
		return toFuncType(ctx, v)
		/*	case *ast.SelectorExpr:
				return toExternalType(ctx, v)
			case *ast.StructType:
				return toStructType(ctx, v)
		*/
	}
	log.Panicln("toType: unknown -", reflect.TypeOf(typ))
	return nil
}

var (
	typesChanDirs = [...]types.ChanDir{
		ast.RECV:            types.RecvOnly,
		ast.SEND:            types.SendOnly,
		ast.SEND | ast.RECV: types.SendRecv,
	}
)

func toIdentType(ctx *blockCtx, ident string) types.Type {
	if _, v := ctx.pkg.Types.Scope().LookupParent(ident, token.NoPos); v != nil {
		if t, ok := v.(*types.TypeName); ok {
			return t.Type()
		}
	}
	panic("TODO: not a type")
}

func toArrayType(ctx *blockCtx, v *ast.ArrayType) types.Type {
	elem := toType(ctx, v.Elt)
	if v.Len == nil {
		return types.NewSlice(elem)
	}
	if _, ok := v.Len.(*ast.Ellipsis); ok {
		return types.NewArray(elem, -1) // A negative length indicates an unknown length
	}
	compileExpr(ctx, v.Elt)
	panic("TODO: array")
}

func toInterfaceType(ctx *blockCtx, v *ast.InterfaceType) types.Type {
	methodsList := v.Methods.List
	if methodsList == nil {
		return types.NewInterfaceType(nil, nil)
	}
	var pkg = ctx.pkg.Types
	var methods []*types.Func
	var embeddeds []types.Type
	for _, m := range methodsList {
		if m.Type == nil { // embedded
			panic("TODO: embedded")
		}
		name := m.Names[0].Name
		typ, ok := m.Type.(*ast.FuncType)
		if !ok {
			panic("TODO: not function type")
		}
		sig := toFuncType(ctx, typ)
		methods = append(methods, types.NewFunc(token.NoPos, pkg, name, sig))
	}
	intf := types.NewInterfaceType(methods, embeddeds)
	intf.Complete()
	return intf
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
