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
	"log"
	"reflect"

	goast "go/ast"
	gotoken "go/token"
	"go/types"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/token"
	"github.com/goplus/gox"
)

// -----------------------------------------------------------------------------

func compileExprLHS(ctx *blockCtx, expr ast.Expr) {
	switch v := expr.(type) {
	case *ast.Ident:
		compileIdentLHS(ctx, v.Name)
	case *ast.IndexExpr:
		compileIndexExprLHS(ctx, v)
	case *ast.SelectorExpr:
		compileSelectorExprLHS(ctx, v)
	case *ast.StarExpr:
		compileStarExprLHS(ctx, v)
	default:
		log.Panicln("compileExpr failed: unknown -", reflect.TypeOf(v))
	}
}

func compileExpr(ctx *blockCtx, expr ast.Expr, twoValue ...bool) {
	switch v := expr.(type) {
	case *ast.Ident:
		compileIdent(ctx, v, false)
	case *ast.BasicLit:
		compileBasicLit(ctx, v)
	case *ast.CallExpr:
		compileCallExpr(ctx, v)
	case *ast.SelectorExpr:
		compileSelectorExpr(ctx, v, true)
	case *ast.BinaryExpr:
		compileBinaryExpr(ctx, v)
	case *ast.UnaryExpr:
		compileUnaryExpr(ctx, v, twoValue != nil && twoValue[0])
	case *ast.FuncLit:
		compileFuncLit(ctx, v)
	case *ast.CompositeLit:
		compileCompositeLit(ctx, v)
	case *ast.SliceLit:
		compileSliceLit(ctx, v)
	case *ast.IndexExpr:
		compileIndexExpr(ctx, v, twoValue != nil && twoValue[0])
	case *ast.SliceExpr:
		compileSliceExpr(ctx, v)
	case *ast.StarExpr:
		compileStarExpr(ctx, v)
	case *ast.ArrayType:
		ctx.cb.Typ(toArrayType(ctx, v))
	case *ast.MapType:
		ctx.cb.Typ(toMapType(ctx, v))
	case *ast.StructType:
		ctx.cb.Typ(toStructType(ctx, v))
	case *ast.ChanType:
		ctx.cb.Typ(toChanType(ctx, v))
	case *ast.InterfaceType:
		ctx.cb.Typ(toInterfaceType(ctx, v))
	case *ast.ComprehensionExpr:
		compileComprehensionExpr(ctx, v, twoValue != nil && twoValue[0])
	case *ast.TypeAssertExpr:
		compileTypeAssertExpr(ctx, v, twoValue != nil && twoValue[0])
	case *ast.ParenExpr:
		compileExpr(ctx, v.X)
	case *ast.ErrWrapExpr:
		compileErrWrapExpr(ctx, v)
	case *ast.FuncType:
		ctx.cb.Typ(toFuncType(ctx, v, nil))
	case *ast.Ellipsis:
		panic("compileEllipsis: ast.Ellipsis unexpected")
	case *ast.KeyValueExpr:
		panic("compileExpr: ast.KeyValueExpr unexpected")
	default:
		log.Panicln("compileExpr failed: unknown -", reflect.TypeOf(v))
	}
}

func compileExprOrNone(ctx *blockCtx, expr ast.Expr) {
	if expr != nil {
		compileExpr(ctx, expr)
	} else {
		ctx.cb.None()
	}
}

func compileUnaryExpr(ctx *blockCtx, v *ast.UnaryExpr, twoValue bool) {
	compileExpr(ctx, v.X)
	ctx.cb.UnaryOp(gotoken.Token(v.Op), twoValue)
}

func compileBinaryExpr(ctx *blockCtx, v *ast.BinaryExpr) {
	compileExpr(ctx, v.X)
	compileExpr(ctx, v.Y)
	ctx.cb.BinaryOp(gotoken.Token(v.Op), v)
}

func compileIndexExprLHS(ctx *blockCtx, v *ast.IndexExpr) {
	compileExpr(ctx, v.X)
	compileExpr(ctx, v.Index)
	ctx.cb.IndexRef(1, v)
}

func compileStarExprLHS(ctx *blockCtx, v *ast.StarExpr) { // *x = ...
	compileExpr(ctx, v.X)
	ctx.cb.ElemRef()
}

func compileStarExpr(ctx *blockCtx, v *ast.StarExpr) { // ... = *x
	compileExpr(ctx, v.X)
	ctx.cb.Star()
}

func compileTypeAssertExpr(ctx *blockCtx, v *ast.TypeAssertExpr, twoValue bool) {
	compileExpr(ctx, v.X)
	if v.Type == nil {
		panic("TODO: x.(type) is only used in type switch")
	}
	typ := toType(ctx, v.Type)
	ctx.cb.TypeAssert(typ, twoValue)
}

func compileIndexExpr(ctx *blockCtx, v *ast.IndexExpr, twoValue bool) { // x[i]
	compileExpr(ctx, v.X)
	compileExpr(ctx, v.Index)
	ctx.cb.Index(1, twoValue, v)
}

func compileSliceExpr(ctx *blockCtx, v *ast.SliceExpr) { // x[i:j:k]
	compileExpr(ctx, v.X)
	compileExprOrNone(ctx, v.Low)
	compileExprOrNone(ctx, v.High)
	if v.Slice3 {
		compileExprOrNone(ctx, v.Max)
	}
	ctx.cb.Slice(v.Slice3, v)
}

func compileSelectorExprLHS(ctx *blockCtx, v *ast.SelectorExpr) {
	switch x := v.X.(type) {
	case *ast.Ident:
		o, at := lookupParent(ctx, x.Name)
		if at != nil {
			ctx.cb.VarRef(at.Ref(v.Sel.Name))
			return
		}
		if o == nil {
			log.Panicln("TODO: ident not found -", x.Name)
		}
		ctx.cb.Val(o)
	default:
		compileExpr(ctx, v.X)
	}
	ctx.cb.MemberRef(v.Sel.Name)
}

func compileSelectorExpr(ctx *blockCtx, v *ast.SelectorExpr, autoCall bool) {
	cb := ctx.cb
	switch x := v.X.(type) {
	case *ast.Ident:
		o, at := lookupParent(ctx, x.Name)
		if at != nil {
			cb.Val(at.Ref(v.Sel.Name))
			return
		}
		if o == nil {
			log.Panicln("TODO: ident not found -", x.Name)
		}
		cb.Val(o)
	default:
		compileExpr(ctx, v.X)
	}
	name := v.Sel.Name
	kind, err := cb.Member(name, v)
	if kind != 0 {
		return
	}
	if c := name[0]; c >= 'a' && c <= 'z' {
		name = string(rune(c)+('A'-'a')) + name[1:]
		switch kind, _ = cb.Member(name, v); kind {
		case gox.MemberMethod:
			if autoCall {
				cb.Call(0)
			}
			return
		case gox.MemberField:
			return
		}
	}
	panic(err)
}

func compileCallExpr(ctx *blockCtx, v *ast.CallExpr) {
	switch fn := v.Fun.(type) {
	case *ast.Ident:
		compileIdent(ctx, fn, true)
	case *ast.SelectorExpr:
		compileSelectorExpr(ctx, fn, false)
	default:
		compileExpr(ctx, fn)
	}
	for _, arg := range v.Args {
		compileExpr(ctx, arg)
	}
	ctx.cb.Call(len(v.Args), v.Ellipsis != gotoken.NoPos)
}

func compileFuncLit(ctx *blockCtx, v *ast.FuncLit) {
	cb := ctx.cb
	comments := cb.Comments()
	sig := toFuncType(ctx, v.Type, nil)
	fn := cb.NewClosureWith(sig)
	if body := v.Body; body != nil {
		loadFuncBody(ctx, fn, body)
		cb.SetComments(comments, false)
	}
}

func compileIdentLHS(ctx *blockCtx, name string) {
	if name == "_" {
		ctx.cb.VarRef(nil)
	} else {
		_, o := ctx.cb.Scope().LookupParent(name, gotoken.NoPos)
		if o == nil {
			panic("TODO: var not found")
		}
		ctx.cb.VarRef(o)
	}
}

func compileIdent(ctx *blockCtx, ident *ast.Ident, allowBuiltin bool) {
	name := ident.Name
	if name == "_" {
		log.Panicln("TODO: cannot use _ as value")
	}
	if o, _ := lookupParent(ctx, name); o != nil {
		if !allowBuiltin {
			if _, ok := o.(*types.Builtin); ok {
				panic("unexpected builtin: " + name)
			}
		}
		ctx.cb.Val(o, ident)
	} else {
		log.Panicln("TODO: var not found -", name)
	}
}

func lookupParent(ctx *blockCtx, name string) (types.Object, *gox.PkgRef) {
	at, o := ctx.cb.Scope().LookupParent(name, token.NoPos)
	if o != nil && at != types.Universe {
		return o, nil
	}
	if ctx.loadSymbol(name) {
		return ctx.pkg.Types.Scope().Lookup(name), nil
	}
	if pkgRef, ok := ctx.imports[name]; ok {
		return nil, pkgRef
	}
	if obj := ctx.pkg.Builtin().Ref(name); obj != nil {
		return obj, nil
	}
	return o, nil
}

func compileBasicLit(ctx *blockCtx, v *ast.BasicLit) {
	if v.Kind == token.RAT {
		panic("TODO: rational constant")
	}
	ctx.cb.Val(&goast.BasicLit{
		Kind:  gotoken.Token(v.Kind),
		Value: v.Value,
	}, v)
}

const (
	compositeLitVal    = 0
	compositeLitKeyVal = 1
)

func checkCompositeLitElts(ctx *blockCtx, elts []ast.Expr) (kind int) {
	for _, elt := range elts {
		if _, ok := elt.(*ast.KeyValueExpr); ok {
			return compositeLitKeyVal
		}
	}
	return compositeLitVal
}

func compileCompositeLitElts(ctx *blockCtx, elts []ast.Expr, kind int) {
	for _, elt := range elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			compileExpr(ctx, kv.Key)
			compileExpr(ctx, kv.Value)
		} else {
			if kind == compositeLitKeyVal {
				ctx.cb.None()
			}
			compileExpr(ctx, elt)
		}
	}
}

func compileStructLitInKeyVal(ctx *blockCtx, elts []ast.Expr, t *types.Struct, typ types.Type) {
	for _, elt := range elts {
		kv := elt.(*ast.KeyValueExpr)
		name := kv.Key.(*ast.Ident).Name
		if idx := lookupField(t, name); idx >= 0 {
			ctx.cb.Val(idx)
		} else {
			log.Panicln("TODO: struct member not found -", name)
		}
		compileExpr(ctx, kv.Value)
	}
	ctx.cb.StructLit(typ, len(elts)<<1, true)
}

func lookupField(t *types.Struct, name string) int {
	for i, n := 0, t.NumFields(); i < n; i++ {
		if fld := t.Field(i); fld.Name() == name {
			return i
		}
	}
	return -1
}

func compileCompositeLit(ctx *blockCtx, v *ast.CompositeLit) {
	var typ, underlying types.Type
	if v.Type != nil {
		typ = toType(ctx, v.Type)
		underlying = getUnderlying(ctx, typ)
	}
	kind := checkCompositeLitElts(ctx, v.Elts)
	if t, ok := underlying.(*types.Struct); ok && kind == compositeLitKeyVal {
		compileStructLitInKeyVal(ctx, v.Elts, t, typ)
		return
	}
	compileCompositeLitElts(ctx, v.Elts, kind)
	n := len(v.Elts)
	if typ == nil {
		ctx.cb.MapLit(nil, n<<1)
		return
	}
	switch underlying.(type) {
	case *types.Slice:
		ctx.cb.SliceLit(typ, n<<kind, kind == compositeLitKeyVal)
	case *types.Array:
		ctx.cb.ArrayLit(typ, n<<kind, kind == compositeLitKeyVal)
	case *types.Map:
		ctx.cb.MapLit(typ, n<<1)
	case *types.Struct:
		ctx.cb.StructLit(typ, n, false)
	default:
		log.Panicln("compileCompositeLit: unknown type -", reflect.TypeOf(underlying))
	}
}

func compileSliceLit(ctx *blockCtx, v *ast.SliceLit) {
	n := len(v.Elts)
	for _, elt := range v.Elts {
		compileExpr(ctx, elt)
	}
	ctx.cb.SliceLit(nil, n)
}

const (
	comprehensionInvalid = iota
	comprehensionList
	comprehensionMap
	comprehensionSelect
)

func comprehensionKind(v *ast.ComprehensionExpr) int {
	switch v.Tok {
	case token.LBRACK: // [
		return comprehensionList
	case token.LBRACE: // {
		if _, ok := v.Elt.(*ast.KeyValueExpr); ok {
			return comprehensionMap
		}
		return comprehensionSelect
	}
	panic("TODO: invalid comprehensionExpr")
}

func compileComprehensionExpr(ctx *blockCtx, v *ast.ComprehensionExpr, twoValue bool) {
	kind := comprehensionKind(v)
	pkg, cb := ctx.pkg, ctx.cb
	ret := pkg.NewAutoParam("_gop_ret")
	var results *types.Tuple
	if kind == comprehensionSelect && twoValue {
		boolean := pkg.NewParam(token.NoPos, "_gop_ok", types.Typ[types.Bool])
		results = types.NewTuple(ret, boolean)
	} else {
		results = types.NewTuple(ret)
	}
	cb.NewClosure(nil, results, false).BodyStart(pkg)
	if kind == comprehensionMap {
		cb.VarRef(ret).ZeroLit(ret.Type()).Assign(1)
	}
	end := 0
	for i := len(v.Fors) - 1; i >= 0; i-- {
		names := make([]string, 0, 2)
		forStmt := v.Fors[i]
		if forStmt.Key != nil {
			names = append(names, forStmt.Key.Name)
		} else {
			names = append(names, "_")
		}
		names = append(names, forStmt.Value.Name)
		cb.ForRange(names...)
		compileExpr(ctx, forStmt.X)
		cb.RangeAssignThen()
		if forStmt.Cond != nil {
			cb.If()
			compileExpr(ctx, forStmt.Cond)
			cb.Then()
			end++
		}
		end++
	}
	switch kind {
	case comprehensionList: // _gop_ret = append(_gop_ret, elt)
		cb.VarRef(ret)
		cb.Val(pkg.Builtin().Ref("append"))
		cb.Val(ret)
		compileExpr(ctx, v.Elt)
		cb.Call(2).Assign(1)
	case comprehensionMap: // _gop_ret[key] = val
		cb.Val(ret)
		kv := v.Elt.(*ast.KeyValueExpr)
		compileExpr(ctx, kv.Key)
		cb.IndexRef(1)
		compileExpr(ctx, kv.Value)
		cb.Assign(1)
	default: // return elt, true
		compileExpr(ctx, v.Elt)
		n := 1
		if twoValue {
			cb.Val(true)
			n++
		}
		cb.Return(n)
	}
	for i := 0; i < end; i++ {
		cb.End()
	}
	cb.Return(0).End().Call(0)
}

var (
	tyError = types.Universe.Lookup("error").Type()
)

func compileErrWrapExpr(ctx *blockCtx, v *ast.ErrWrapExpr) {
	pkg, cb := ctx.pkg, ctx.cb
	inGlobal := (cb.Scope().Parent() == types.Universe) // not in a function body
	ret := pkg.NewAutoParam("_gop_ret")
	sig := types.NewSignature(nil, nil, types.NewTuple(ret), false)
	if inGlobal {
		cb.NewClosureWith(sig).BodyStart(pkg)
	} else {
		cb.CallInlineClosureStart(sig, 0, false)
	}
	cb.NewVar(tyError, "_gop_err")
	err := cb.Scope().Lookup("_gop_err")
	cb.VarRef(ret).VarRef(err)
	compileExpr(ctx, v.X)
	cb.Assign(2, 1)
	cb.If().Val(err).CompareNil(gotoken.NEQ).Then()
	if v.Tok == token.NOT { // expr!
		cb.Val(pkg.Builtin().Ref("panic")).Val(err).Call(1).EndStmt() // TODO: wrap err
	} else if v.Default == nil { // expr?
		if inGlobal {
			panic("TODO: can't use expr? in global")
		}
		cb.Val(err).ReturnErr(true) // TODO: wrap err & return err
	} else { // expr?:val
		compileExpr(ctx, v.Default)
		cb.Return(1)
	}
	cb.End().Return(0).End()
	if inGlobal {
		cb.Call(0)
	}
}

// -----------------------------------------------------------------------------
