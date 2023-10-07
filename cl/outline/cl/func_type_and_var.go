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

func toRecv(ctx *blockCtx, recv *ast.FieldList) (ret *types.Var, ok bool) {
	v := recv.List[0]
	var name string
	if len(v.Names) > 0 {
		name = v.Names[0].Name
	}
	t := toType(ctx, v.Type)
	if ok = checkRecvType(ctx, t, v.Type); ok {
		ret = ctx.pkg.NewParam(v.Pos(), name, t)
	}
	return
}

func checkRecvType(ctx *blockCtx, typ types.Type, recv ast.Expr) (ok bool) {
	t := indirect(typ)
	if _, ok = t.(*types.Named); !ok {
		pos := ctx.Position(recv.Pos())
		ctx.handleCodeErrorf(&pos, "invalid receiver type %v", typ)
	}
	return
}

func indirect(typ types.Type) types.Type {
	if t, ok := typ.(*types.Pointer); ok {
		return t.Elem()
	}
	return typ
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

func toTypeInited(ctx *blockCtx, typ ast.Expr) types.Type {
	t := toType(ctx, typ)
	if named, ok := t.(*types.Named); ok {
		ctx.ensureLoaded(named)
	}
	return t
}

func toType(ctx *blockCtx, typ ast.Expr) types.Type {
	switch v := typ.(type) {
	case *ast.Ident:
		/*
			ctx.idents = append(ctx.idents, v)
			defer func() {
				ctx.idents = ctx.idents[:len(ctx.idents)-1]
			}()
		*/
		typ := toIdentType(ctx, v)
		/*
			if ctx.inInst == 0 {
				if t, ok := typ.(*types.Named); ok {
					if withTypeParams(ctx, t) {
						pos := ctx.idents[0].Pos()
						for _, i := range ctx.idents {
							if i.Name == v.Name {
								pos = i.Pos()
								break
							}
						}
						panic(ctx.newCodeErrorf(pos, "cannot use generic type %v without instantiation", t.Obj().Type()))
					}
				}
			}
		*/
		return typ
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
		return toFuncType(ctx, v, nil, nil)
	case *ast.SelectorExpr:
		typ := toExternalType(ctx, v)
		if ctx.inInst == 0 {
			if t, ok := typ.(*types.Named); ok {
				if withTypeParams(ctx, t) {
					panic(ctx.newCodeErrorf(v.Pos(), "cannot use generic type %v without instantiation", t.Obj().Type()))
				}
			}
		}
		return typ
	case *ast.ParenExpr:
		return toType(ctx, v.X)
	case *ast.BinaryExpr:
		return toBinaryExprType(ctx, v)
	case *ast.UnaryExpr:
		return toUnaryExprType(ctx, v)
	case *ast.IndexExpr:
		return toIndexType(ctx, v)
	case *ast.IndexListExpr:
		return toIndexListType(ctx, v)
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
	if pr, ok := ctx.findImport(name); ok {
		o := pr.TryRef(v.Sel.Name)
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
	if ctx.tlookup != nil {
		if typ := ctx.tlookup.Lookup(ident.Name); typ != nil {
			return typ
		}
	}
	v, builtin := lookupType(ctx, ident.Name)
	if isBuiltin(builtin) {
		panic(ctx.newCodeErrorf(ident.Pos(), "use of builtin %s not in function call", ident.Name))
	}
	if t, ok := v.(*types.TypeName); ok {
		return t.Type()
	}
	if v, _ := lookupPkgRef(ctx, nil, ident, objPkgRef); v != nil {
		if t, ok := v.(*types.TypeName); ok {
			return t.Type()
		}
	}
	panic(ctx.newCodeErrorf(ident.Pos(), "%s is not a type", ident.Name))
}

func lookupType(ctx *blockCtx, name string) (obj types.Object, builtin types.Object) {
	at, o := ctx.cb.Scope().LookupParent(name, token.NoPos)
	if o != nil {
		if at != types.Universe {
			if debugLookup {
				log.Println("==> LookupParent", name, "=>", o)
			}
			return o, nil
		}
	}
	if obj := ctx.pkg.Builtin().TryRef(name); obj != nil {
		return obj, o
	}
	return o, o
}

type checkRedecl struct {
	// ctx *blockCtx
	names map[string]token.Pos
}

func newCheckRedecl() *checkRedecl {
	p := &checkRedecl{names: make(map[string]token.Pos)}
	return p
}

func (p *checkRedecl) chkRedecl(ctx *blockCtx, name string, pos token.Pos) bool {
	if name == "_" {
		return false
	}
	if opos, ok := p.names[name]; ok {
		npos := ctx.Position(pos)
		ctx.handleCodeErrorf(&npos, "%v redeclared\n\t%v other declaration of %v",
			name, ctx.Position(opos), name)
		return true
	}
	p.names[name] = pos
	return false
}

func toStructType(ctx *blockCtx, v *ast.StructType) *types.Struct {
	pkg := ctx.pkg.Types
	fieldList := v.Fields.List
	fields := make([]*types.Var, 0, len(fieldList))
	tags := make([]string, 0, len(fieldList))
	chk := newCheckRedecl()
	for _, field := range fieldList {
		typ := toTypeInited(ctx, field.Type)
		if len(field.Names) == 0 { // embedded
			name := getTypeName(typ)
			if chk.chkRedecl(ctx, name, field.Type.Pos()) {
				continue
			}
			fld := types.NewField(token.NoPos, pkg, name, typ, true)
			fields = append(fields, fld)
			tags = append(tags, toFieldTag(field.Tag))
			continue
		}
		for _, name := range field.Names {
			if chk.chkRedecl(ctx, name.Name, name.NamePos) {
				continue
			}
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

func toInterfaceType(ctx *blockCtx, v *ast.InterfaceType) *types.Interface {
	methodsList := v.Methods.List
	if methodsList == nil {
		return gox.TyEmptyInterface
	}
	var pkg = ctx.pkg.Types
	var methods []*types.Func
	var embeddeds []types.Type
	for _, m := range methodsList {
		if len(m.Names) == 0 { // embedded
			typ := toTypeInited(ctx, m.Type)
			embeddeds = append(embeddeds, typ)
			continue
		}
		name := m.Names[0].Name
		sig := toFuncType(ctx, m.Type.(*ast.FuncType), nil, nil)
		methods = append(methods, types.NewFunc(token.NoPos, pkg, name, sig))
	}
	intf := types.NewInterfaceType(methods, embeddeds).Complete()
	return intf
}

func toIndexType(ctx *blockCtx, v *ast.IndexExpr) types.Type {
	ctx.inInst++
	defer func() {
		ctx.inInst--
	}()
	ctx.cb.Typ(toTypeInited(ctx, v.X), v.X)
	ctx.cb.Typ(toTypeInited(ctx, v.Index), v.Index)
	ctx.cb.Index(1, false, v)
	return ctx.cb.InternalStack().Pop().Type.(*gox.TypeType).Type()
}

func toIndexListType(ctx *blockCtx, v *ast.IndexListExpr) types.Type {
	ctx.inInst++
	defer func() {
		ctx.inInst--
	}()
	ctx.cb.Typ(toTypeInited(ctx, v.X), v.X)
	for _, index := range v.Indices {
		ctx.cb.Typ(toTypeInited(ctx, index), index)
	}
	ctx.cb.Index(len(v.Indices), false, v)
	return ctx.cb.InternalStack().Pop().Type.(*gox.TypeType).Type()
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
