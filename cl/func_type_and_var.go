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

package cl

import (
	"go/constant"
	"go/types"
	"log"
	"math/big"
	"reflect"
	"strconv"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/token"
	"github.com/goplus/gox"
)

// -----------------------------------------------------------------------------

func toFuncType(ctx *blockCtx, typ *ast.FuncType, recv *types.Var) *types.Signature {
	params, variadic := toParams(ctx, typ.Params.List)
	results := toResults(ctx, typ.Results)
	return types.NewSignature(recv, params, results, variadic)
}

func toRecv(ctx *blockCtx, recv *ast.FieldList) *types.Var {
	v := recv.List[0]
	var name string
	if len(v.Names) > 0 {
		name = v.Names[0].Name
	}
	return ctx.pkg.NewParam(v.Pos(), name, toType(ctx, v.Type))
}

func getRecvTypeName(ctx *pkgCtx, recv *ast.FieldList, handleErr bool) (string, bool) {
	typ := recv.List[0].Type
	if t, ok := typ.(*ast.StarExpr); ok {
		typ = t.X
	}
	if t, ok := typ.(*ast.Ident); ok {
		return t.Name, true
	}
	if handleErr {
		src, pos := ctx.LoadExpr(typ)
		ctx.handleCodeErrorf(&pos, "invalid receiver type %v (%v is not a defined type)", src, src)
	}
	return "", false
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
		return append(args, pkg.NewParam(fld.Pos(), "", typ))
	}
	for _, name := range fld.Names {
		args = append(args, pkg.NewParam(name.Pos(), name.Name, typ))
	}
	return args
}

// -----------------------------------------------------------------------------

func toType(ctx *blockCtx, typ ast.Expr) types.Type {
	switch v := typ.(type) {
	case *ast.Ident:
		return toIdentType(ctx, v)
	case *ast.StarExpr:
		elem := toType(ctx, v.X)
		return types.NewPointer(elem)
	case *ast.ArrayType:
		return toArrayType(ctx, v)
	case *ast.InterfaceType:
		return toInterfaceType(ctx, v)
	case *ast.Ellipsis:
		elem := toType(ctx, v.Elt)
		return types.NewSlice(elem)
	case *ast.MapType:
		return toMapType(ctx, v)
	case *ast.StructType:
		return toStructType(ctx, v)
	case *ast.ChanType:
		return toChanType(ctx, v)
	case *ast.FuncType:
		return toFuncType(ctx, v, nil)
	case *ast.SelectorExpr:
		return toExternalType(ctx, v)
	case *ast.ParenExpr:
		return toType(ctx, v.X)
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

func toChanType(ctx *blockCtx, v *ast.ChanType) *types.Chan {
	return types.NewChan(typesChanDirs[v.Dir], toType(ctx, v.Value))
}

func toExternalType(ctx *blockCtx, v *ast.SelectorExpr) types.Type {
	name := v.X.(*ast.Ident).Name
	if pkgRef, ok := ctx.imports[name]; ok {
		o := pkgRef.TryRef(v.Sel.Name)
		if t, ok := o.(*types.TypeName); ok {
			return t.Type()
		}
		panic(ctx.newCodeErrorf(v.Pos(), "%s.%s is not a type", name, v.Sel.Name))
	}
	panic(ctx.newCodeErrorf(v.Pos(), "undefined: %s", name))
}

/*-----------------------------------------------------------------------------

Name context:
- type
- pkgRef.type
- spx.type

// ---------------------------------------------------------------------------*/

func toIdentType(ctx *blockCtx, ident *ast.Ident) types.Type {
	v, builtin := lookupType(ctx, ident.Name)
	if isBuiltin(builtin) {
		panic(ctx.newCodeErrorf(ident.Pos(), "use of builtin %s not in function call", ident.Name))
	}
	if t, ok := v.(*types.TypeName); ok {
		return t.Type()
	}
	if v, _ := lookupPkgRef(ctx, nil, ident); v != nil {
		if t, ok := v.(*types.TypeName); ok {
			return t.Type()
		}
	}
	panic(ctx.newCodeErrorf(ident.Pos(), "%s is not a type", ident.Name))
}

// TODO: optimization
func lookupType(ctx *blockCtx, name string) (types.Object, types.Object) {
	at, o := ctx.cb.Scope().LookupParent(name, token.NoPos)
	if o != nil && at != types.Universe {
		if debugLookup {
			log.Println("==> LookupParent", name, "=>", o)
		}
		return o, nil
	}
	if ctx.loadSymbol(name) {
		if v := ctx.pkg.Types.Scope().Lookup(name); v != nil {
			if debugLookup {
				log.Println("==> Lookup (LoadSymbol)", name, "=>", v)
			}
			return v, nil
		}
	}
	if obj := ctx.pkg.Builtin().TryRef(name); obj != nil {
		return obj, o
	}
	return o, o
}

func toStructType(ctx *blockCtx, v *ast.StructType) *types.Struct {
	pkg := ctx.pkg.Types
	fieldList := v.Fields.List
	fields := make([]*types.Var, 0, len(fieldList))
	tags := make([]string, 0, len(fieldList))
	for _, field := range fieldList {
		typ := toType(ctx, field.Type)
		if field.Names == nil { // embedded
			fld := types.NewField(token.NoPos, pkg, getTypeName(typ), typ, true)
			fields = append(fields, fld)
			tags = append(tags, toFieldTag(field.Tag))
			continue
		}
		for _, name := range field.Names {
			fld := types.NewField(token.NoPos, pkg, name.Name, typ, false)
			fields = append(fields, fld)
			tags = append(tags, toFieldTag(field.Tag))
		}
	}
	return types.NewStruct(fields, tags)
}

func toFieldTag(v *ast.BasicLit) string {
	if v != nil {
		tag, err := strconv.Unquote(v.Value)
		if err != nil {
			log.Panicln("TODO: toFieldTag -", err)
		}
		return tag
	}
	return ""
}

func getTypeName(typ types.Type) string {
	if t, ok := typ.(*types.Pointer); ok {
		typ = t.Elem()
	}
	switch t := typ.(type) {
	case *types.Named:
		return t.Obj().Name()
	case *types.Basic:
		return t.Name()
	default:
		panic("TODO: getTypeName")
	}
}

func toMapType(ctx *blockCtx, v *ast.MapType) *types.Map {
	key := toType(ctx, v.Key)
	val := toType(ctx, v.Value)
	return types.NewMap(key, val)
}

func toArrayType(ctx *blockCtx, v *ast.ArrayType) types.Type {
	elem := toType(ctx, v.Elt)
	if v.Len == nil {
		return types.NewSlice(elem)
	}
	if _, ok := v.Len.(*ast.Ellipsis); ok {
		return types.NewArray(elem, -1) // A negative length indicates an unknown length
	}
	return types.NewArray(elem, toInt64(ctx, v.Len, "non-constant array bound %s"))
}

func toInt64(ctx *blockCtx, e ast.Expr, emsg string) int64 {
	cb := ctx.pkg.ConstStart()
	compileExpr(ctx, e)
	tv := cb.EndConst()
	if val := tv.CVal; val != nil {
		if val.Kind() == constant.Float {
			if v, ok := constant.Val(val).(*big.Rat); ok && v.IsInt() {
				return v.Num().Int64()
			}
		} else if v, ok := constant.Int64Val(val); ok {
			return v
		}
	}
	src, pos := ctx.LoadExpr(e)
	panic(newCodeErrorf(&pos, emsg, src))
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
		if m.Names == nil { // embedded
			embeddeds = append(embeddeds, toType(ctx, m.Type))
			continue
		}
		name := m.Names[0].Name
		sig := toFuncType(ctx, m.Type.(*ast.FuncType), nil)
		methods = append(methods, types.NewFunc(token.NoPos, pkg, name, sig))
	}
	intf := types.NewInterfaceType(methods, embeddeds).Complete()
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
