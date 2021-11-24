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
	goast "go/ast"
	gotoken "go/token"
	"go/types"
	"log"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/token"
	"github.com/goplus/gox"
)

/*-----------------------------------------------------------------------------

Name context:
- varVal               (ident)
- varRef = expr        (identLHS)
- pkgRef.member        (selectorExpr)
- pkgRef.member = expr (selectorExprLHS)
- pkgRef.fn(args)      (callExpr)
- fn(args)             (callExpr)
- spx.fn(args)         (callExpr)
- this.member          (classMember)
- this.method(args)    (classMember)

Name lookup:
- local variables
- $recv members (only in class files)
- package globals (variables, constants, types, imported packages etc.)
- $spx package exports (only in class files)
- $universe package exports (including builtins)

// ---------------------------------------------------------------------------*/

const (
	clIdentAutoCall = 1 << iota
	clIdentAllowBuiltin
	clIdentLHS
	clIdentSelectorExpr
	clIdentGoto
)

func compileIdent(ctx *blockCtx, ident *ast.Ident, flags int) *gox.PkgRef {
	fvalue := (flags&clIdentSelectorExpr) != 0 || (flags&clIdentLHS) == 0
	name := ident.Name
	if name == "_" {
		if fvalue {
			panic(ctx.newCodeError(ident.Pos(), "cannot use _ as value"))
		}
		ctx.cb.VarRef(nil)
		return nil
	}

	scope := ctx.pkg.Types.Scope()
	at, o := ctx.cb.Scope().LookupParent(name, token.NoPos)
	if o != nil {
		if at != scope && at != types.Universe { // local object
			goto find
		}
	}

	if ctx.fileType > 0 { // in a Go+ class file
		if fn := ctx.cb.Func(); fn != nil {
			sig := fn.Ancestor().Type().(*types.Signature)
			if recv := sig.Recv(); recv != nil {
				ctx.cb.Val(recv)
				if compileMember(ctx, ident, name, flags) == nil { // class member object
					return nil
				}
				ctx.cb.InternalStack().PopN(1)
			}
		}
	}

	// global object
	if ctx.loadSymbol(name) {
		o, at = scope.Lookup(name), scope
	}
	if o != nil && at != types.Universe {
		goto find
	}

	// pkgRef object
	if (flags & clIdentSelectorExpr) != 0 {
		if pkgRef, ok := ctx.imports[name]; ok {
			return pkgRef
		}
	}

	// object from import . "xxx"
	if compilePkgRef(ctx, nil, ident, flags) {
		return nil
	}

	// universe object
	if (flags&clIdentAllowBuiltin) == 0 && isBuiltin(o) {
		panic(ctx.newCodeErrorf(ident.Pos(), "use of builtin %s not in function call", name))
	}
	if obj := ctx.pkg.Builtin().TryRef(name); obj != nil {
		o = obj
	} else if o == nil {
		if (clIdentGoto & flags) != 0 {
			l := ident.Obj.Data.(*ast.Ident)
			panic(ctx.newCodeErrorf(l.Pos(), "label %v is not defined", l.Name))
		}
		panic(ctx.newCodeErrorf(ident.Pos(), "undefined: %s", name))
	}

find:
	if fvalue {
		ctx.cb.Val(o, ident)
	} else {
		ctx.cb.VarRef(o, ident)
	}
	return nil
}

func isBuiltin(o types.Object) bool {
	if _, ok := o.(*types.Builtin); ok {
		return ok
	}
	return false
}

func compileMember(ctx *blockCtx, v ast.Node, name string, flags int) error {
	cb := ctx.cb
	lhs := (flags & clIdentLHS) != 0
	kind, err := cb.Member(name, lhs, v)
	if kind != 0 {
		return nil
	}
	if c := name[0]; c >= 'a' && c <= 'z' {
		name = string(rune(c)+('A'-'a')) + name[1:]
		switch kind, _ = cb.Member(name, lhs, v); kind {
		case gox.MemberMethod:
			if (flags & clIdentAutoCall) != 0 {
				cb.Call(0)
			}
			return nil
		case gox.MemberField:
			return nil
		}
	}
	return err
}

func compileExprLHS(ctx *blockCtx, expr ast.Expr) {
	switch v := expr.(type) {
	case *ast.Ident:
		compileIdent(ctx, v, clIdentLHS)
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
		compileIdent(ctx, v, clIdentAutoCall)
	case *ast.BasicLit:
		compileBasicLit(ctx, v)
	case *ast.CallExpr:
		compileCallExpr(ctx, v, 0)
	case *ast.SelectorExpr:
		compileSelectorExpr(ctx, v, clIdentAutoCall)
	case *ast.BinaryExpr:
		compileBinaryExpr(ctx, v)
	case *ast.UnaryExpr:
		compileUnaryExpr(ctx, v, twoValue != nil && twoValue[0])
	case *ast.FuncLit:
		compileFuncLit(ctx, v)
	case *ast.CompositeLit:
		compileCompositeLit(ctx, v, nil, false)
	case *ast.SliceLit:
		compileSliceLit(ctx, v)
	case *ast.RangeExpr:
		compileRangeExpr(ctx, v)
	case *ast.IndexExpr:
		compileIndexExpr(ctx, v, twoValue != nil && twoValue[0])
	case *ast.SliceExpr:
		compileSliceExpr(ctx, v)
	case *ast.StarExpr:
		compileStarExpr(ctx, v)
	case *ast.ArrayType:
		ctx.cb.Typ(toArrayType(ctx, v), v)
	case *ast.MapType:
		ctx.cb.Typ(toMapType(ctx, v), v)
	case *ast.StructType:
		ctx.cb.Typ(toStructType(ctx, v), v)
	case *ast.ChanType:
		ctx.cb.Typ(toChanType(ctx, v), v)
	case *ast.InterfaceType:
		ctx.cb.Typ(toInterfaceType(ctx, v), v)
	case *ast.ComprehensionExpr:
		compileComprehensionExpr(ctx, v, twoValue != nil && twoValue[0])
	case *ast.TypeAssertExpr:
		compileTypeAssertExpr(ctx, v, twoValue != nil && twoValue[0])
	case *ast.ParenExpr:
		compileExpr(ctx, v.X, twoValue...)
	case *ast.ErrWrapExpr:
		compileErrWrapExpr(ctx, v)
	case *ast.FuncType:
		ctx.cb.Typ(toFuncType(ctx, v, nil), v)
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
	ctx.cb.TypeAssert(typ, twoValue, v)
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
		if at := compileIdent(ctx, x, clIdentLHS|clIdentSelectorExpr); at != nil {
			ctx.cb.VarRef(at.Ref(v.Sel.Name))
			return
		}
	default:
		compileExpr(ctx, v.X)
	}
	ctx.cb.MemberRef(v.Sel.Name, v)
}

func compileSelectorExpr(ctx *blockCtx, v *ast.SelectorExpr, flags int) {
	switch x := v.X.(type) {
	case *ast.Ident:
		if at := compileIdent(ctx, x, flags|clIdentSelectorExpr); at != nil {
			if compilePkgRef(ctx, at, v.Sel, flags) {
				return
			}
			if token.IsExported(v.Sel.Name) {
				panic(ctx.newCodeErrorf(x.Pos(), "undefined: %s.%s", x.Name, v.Sel.Name))
			}
			panic(ctx.newCodeErrorf(x.Pos(), "cannot refer to unexported name %s.%s", x.Name, v.Sel.Name))
		}
	default:
		compileExpr(ctx, v.X)
	}
	if err := compileMember(ctx, v, v.Sel.Name, flags); err != nil {
		panic(err)
	}
}

func pkgRef(at *gox.PkgRef, name string) (o types.Object, canAutoCall bool) {
	if c := name[0]; c >= 'a' && c <= 'z' {
		name = string(rune(c)+('A'-'a')) + name[1:]
		return at.TryRef(name), true
	}
	return at.TryRef(name), false
}

func lookupPkgRef(ctx *blockCtx, pkg *gox.PkgRef, x *ast.Ident) (o types.Object, canAutoCall bool) {
	if pkg != nil {
		return pkgRef(pkg, x.Name)
	}
	for _, at := range ctx.lookups {
		if o2, canAutoCall2 := pkgRef(at, x.Name); o2 != nil {
			if o != nil {
				panic(ctx.newCodeErrorf(
					x.Pos(), "confliction: %s declared both in \"%s\" and \"%s\"",
					x.Name, at.Types.Path(), pkg.Types.Path()))
			}
			pkg, o, canAutoCall = at, o2, canAutoCall2
		}
	}
	return
}

func compilePkgRef(ctx *blockCtx, at *gox.PkgRef, x *ast.Ident, flags int) bool {
	if v, canAutoCall := lookupPkgRef(ctx, at, x); v != nil {
		cb := ctx.cb
		if (flags & clIdentLHS) != 0 {
			cb.VarRef(v, x)
		} else {
			cb.Val(v, x)
			if canAutoCall && (flags&clIdentAutoCall) != 0 && isFunc(v.Type()) {
				cb.Call(0)
			}
		}
		return true
	}
	return false
}

func isFunc(typ types.Type) bool {
	switch typ.(type) {
	case *types.Signature:
		return true
	}
	return false
}

type fnType struct {
	params   *types.Tuple
	n1       int
	variadic bool
	inited   bool
}

func (p *fnType) arg(i int, ellipsis bool) types.Type {
	if i < p.n1 {
		return p.params.At(i).Type()
	}
	if p.variadic {
		t := p.params.At(p.n1).Type()
		if ellipsis {
			return t
		}
		return t.(*types.Slice).Elem()
	}
	return nil
}

func (p *fnType) init(t *types.Signature) {
	p.params, p.variadic = t.Params(), t.Variadic()
	p.n1 = p.params.Len()
	if p.variadic {
		p.n1--
	}
}

func (p *fnType) initWith(fnt types.Type, idx, nin int) {
	if p.inited {
		return
	}
	p.inited = true
	if t := gox.CheckSignature(fnt, idx, nin); t != nil {
		p.init(t)
	}
}

func compileCallExpr(ctx *blockCtx, v *ast.CallExpr, flags int) {
	switch fn := v.Fun.(type) {
	case *ast.Ident:
		compileIdent(ctx, fn, clIdentAllowBuiltin|flags)
	case *ast.SelectorExpr:
		compileSelectorExpr(ctx, fn, 0)
	default:
		compileExpr(ctx, fn)
	}
	var fn fnType
	fnt := ctx.cb.Get(-1).Type
	ellipsis := (v.Ellipsis != gotoken.NoPos)
	for i, arg := range v.Args {
		if l, ok := arg.(*ast.LambdaExpr); ok {
			fn.initWith(fnt, i, len(l.Lhs))
			if sig, ok := fn.arg(i, true).(*types.Signature); ok {
				compileLambdaExpr(ctx, l, sig.Params())
				continue
			}
		}
		if l, ok := arg.(*ast.LambdaExpr2); ok {
			fn.initWith(fnt, i, len(l.Lhs))
			if sig, ok := fn.arg(i, true).(*types.Signature); ok {
				compileLambdaExpr2(ctx, l, sig.Params(), sig.Results())
				continue
			}
		}
		if c, ok := arg.(*ast.CompositeLit); ok && c.Type == nil {
			fn.initWith(fnt, i, -1)
			compileCompositeLit(ctx, c, fn.arg(i, ellipsis), true)
		} else {
			compileExpr(ctx, arg)
		}
	}
	ctx.cb.CallWith(len(v.Args), ellipsis, v)
}

func compileLambdaParams(ctx *blockCtx, pos token.Pos, lhs []*ast.Ident, in *types.Tuple) []*types.Var {
	pkg := ctx.pkg
	n := len(lhs)
	if nin := in.Len(); n != nin {
		fewOrMany := "few"
		if n > nin {
			fewOrMany = "many"
		}
		has := make([]string, n)
		for i, v := range lhs {
			has[i] = v.Name
		}
		panic(ctx.newCodeErrorf(
			pos, "too %s arguments in lambda expression\n\thave (%s)\n\twant %v", fewOrMany, strings.Join(has, ", "), in))
	}
	params := make([]*types.Var, n)
	for i, name := range lhs {
		params[i] = pkg.NewParam(name.Pos(), name.Name, in.At(i).Type())
	}
	return params
}

func compileLambdaExpr(ctx *blockCtx, v *ast.LambdaExpr, in *types.Tuple) {
	pkg := ctx.pkg
	params := compileLambdaParams(ctx, v.Pos(), v.Lhs, in)
	nout := len(v.Rhs)
	results := make([]*types.Var, nout)
	for i := 0; i < nout; i++ {
		results[i] = pkg.NewAutoParam("")
	}
	ctx.cb.NewClosure(types.NewTuple(params...), types.NewTuple(results...), false).BodyStart(pkg)
	for _, v := range v.Rhs {
		compileExpr(ctx, v)
	}
	ctx.cb.Return(nout).End()
}

func compileLambdaExpr2(ctx *blockCtx, v *ast.LambdaExpr2, in *types.Tuple, out *types.Tuple) {
	pkg := ctx.pkg
	params := compileLambdaParams(ctx, v.Pos(), v.Lhs, in)
	cb := ctx.cb
	comments := cb.Comments()
	nout := out.Len()
	results := make([]*types.Var, nout)
	for i := 0; i < nout; i++ {
		results[i] = pkg.NewParam(token.NoPos, "", out.At(i).Type())
	}
	fn := cb.NewClosure(types.NewTuple(params...), types.NewTuple(results...), false)
	loadFuncBody(ctx, fn, v.Body)
	cb.SetComments(comments, false)
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

func compileBasicLit(ctx *blockCtx, v *ast.BasicLit) {
	if v.Kind == token.RAT {
		val := v.Value
		bi, _ := new(big.Int).SetString(val[:len(val)-1], 10) // remove r suffix
		ctx.cb.UntypedBigInt(bi, v)
		return
	}
	ctx.cb.Val(&goast.BasicLit{Kind: gotoken.Token(v.Kind), Value: v.Value}, v)
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

func compileCompositeLitElts(ctx *blockCtx, elts []ast.Expr, kind int, expected *kvType) {
	for _, elt := range elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			if key, ok := kv.Key.(*ast.CompositeLit); ok && key.Type == nil {
				compileCompositeLit(ctx, key, expected.Key(), false)
			} else {
				compileExpr(ctx, kv.Key)
			}
			if val, ok := kv.Value.(*ast.CompositeLit); ok && val.Type == nil {
				compileCompositeLit(ctx, val, expected.Elem(), false)
			} else {
				compileExpr(ctx, kv.Value)
			}
		} else {
			if kind == compositeLitKeyVal {
				ctx.cb.None()
			}
			if val, ok := elt.(*ast.CompositeLit); ok && val.Type == nil {
				compileCompositeLit(ctx, val, expected.Elem(), false)
			} else {
				compileExpr(ctx, elt)
			}
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

type kvType struct {
	underlying types.Type
	key, val   types.Type
	cached     bool
}

func (p *kvType) required() *kvType {
	if !p.cached {
		p.cached = true
		switch t := p.underlying.(type) {
		case *types.Slice:
			p.key, p.val = types.Typ[types.Int], t.Elem()
		case *types.Array:
			p.key, p.val = types.Typ[types.Int], t.Elem()
		case *types.Map:
			p.key, p.val = t.Key(), t.Elem()
		}
	}
	return p
}

func (p *kvType) Key() types.Type {
	return p.required().key
}

func (p *kvType) Elem() types.Type {
	return p.required().val
}

func getUnderlying(ctx *blockCtx, typ types.Type) types.Type {
	u := typ.Underlying()
	if u == nil {
		if t, ok := typ.(*types.Named); ok {
			ctx.loadNamed(ctx.pkg, t)
			u = t.Underlying()
		}
	}
	return u
}

func compileCompositeLit(ctx *blockCtx, v *ast.CompositeLit, expected types.Type, onlyStruct bool) {
	var hasPtr bool
	var typ, underlying types.Type
	var kind = checkCompositeLitElts(ctx, v.Elts)
	if v.Type != nil {
		typ = toType(ctx, v.Type)
		underlying = getUnderlying(ctx, typ)
	} else if expected != nil {
		if t, ok := expected.(*types.Pointer); ok {
			expected, hasPtr = t.Elem(), true
		}
		if onlyStruct {
			if kind == compositeLitKeyVal {
				t := getUnderlying(ctx, expected)
				if _, ok := t.(*types.Struct); ok { // can't omit non-struct type
					typ, underlying = expected, t
				}
			}
		} else {
			typ, underlying = expected, getUnderlying(ctx, expected)
		}
	}
	if t, ok := underlying.(*types.Struct); ok && kind == compositeLitKeyVal {
		compileStructLitInKeyVal(ctx, v.Elts, t, typ)
		if hasPtr {
			ctx.cb.UnaryOp(gotoken.AND)
		}
		return
	}
	compileCompositeLitElts(ctx, v.Elts, kind, &kvType{underlying: underlying})
	n := len(v.Elts)
	if typ == nil {
		if kind == compositeLitVal && n > 0 {
			panic("TODO: mapLit should be in {key: val, ...} form")
		}
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
	if hasPtr {
		ctx.cb.UnaryOp(gotoken.AND)
	}
}

func compileSliceLit(ctx *blockCtx, v *ast.SliceLit) {
	n := len(v.Elts)
	for _, elt := range v.Elts {
		compileExpr(ctx, elt)
	}
	ctx.cb.SliceLit(nil, n)
}

func compileRangeExpr(ctx *blockCtx, v *ast.RangeExpr) {
	pkg, cb := ctx.pkg, ctx.cb
	cb.Val(pkg.Builtin().Ref("newRange"))
	if v.First == nil {
		ctx.cb.Val(0, v)
	} else {
		compileExpr(ctx, v.First)
	}
	compileExpr(ctx, v.Last)
	if v.Expr3 == nil {
		ctx.cb.Val(1, v)
	} else {
		compileExpr(ctx, v.Expr3)
	}
	cb.Call(3)
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

// [expr for k, v <- container, cond]
// {for k, v <- container, cond}
// {expr for k, v <- container, cond}
// {kexpr: vexpr for k, v <- container, cond}
func compileComprehensionExpr(ctx *blockCtx, v *ast.ComprehensionExpr, twoValue bool) {
	kind := comprehensionKind(v)
	pkg, cb := ctx.pkg, ctx.cb
	var results *types.Tuple
	var ret *gox.Param
	if v.Elt == nil {
		boolean := pkg.NewParam(token.NoPos, "_gop_ok", types.Typ[types.Bool])
		results = types.NewTuple(boolean)
	} else {
		ret = pkg.NewAutoParam("_gop_ret")
		if kind == comprehensionSelect && twoValue {
			boolean := pkg.NewParam(token.NoPos, "_gop_ok", types.Typ[types.Bool])
			results = types.NewTuple(ret, boolean)
		} else {
			results = types.NewTuple(ret)
		}
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
		cb.RangeAssignThen(forStmt.TokPos)
		if forStmt.Cond != nil {
			cb.If()
			if forStmt.Init != nil {
				compileStmt(ctx, forStmt.Init)
			}
			compileExpr(ctx, forStmt.Cond)
			cb.Then()
			end++
		}
		end++
	}
	switch kind {
	case comprehensionList:
		// _gop_ret = append(_gop_ret, elt)
		cb.VarRef(ret)
		cb.Val(pkg.Builtin().Ref("append"))
		cb.Val(ret)
		compileExpr(ctx, v.Elt)
		cb.Call(2).Assign(1)
	case comprehensionMap:
		// _gop_ret[key] = val
		cb.Val(ret)
		kv := v.Elt.(*ast.KeyValueExpr)
		compileExpr(ctx, kv.Key)
		cb.IndexRef(1)
		compileExpr(ctx, kv.Value)
		cb.Assign(1)
	default:
		if v.Elt == nil {
			// return true
			cb.Val(true)
			cb.Return(1)
		} else {
			// return elt, true
			compileExpr(ctx, v.Elt)
			n := 1
			if twoValue {
				cb.Val(true)
				n++
			}
			cb.Return(n)
		}
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
	useClosure := (v.Tok == token.NOT || v.Default != nil)
	if !useClosure && (cb.Scope().Parent() == types.Universe) {
		panic("TODO: can't use expr? in global")
	}

	compileExpr(ctx, v.X)
	x := cb.InternalStack().Pop()
	n := 0
	results, ok := x.Type.(*types.Tuple)
	if ok {
		n = results.Len() - 1
	}

	var ret []*types.Var
	if n > 0 {
		i, retName := 0, "_gop_ret"
		ret = make([]*gox.Param, n)
		for {
			ret[i] = pkg.NewAutoParam(retName)
			i++
			if i >= n {
				break
			}
			retName = "_gop_ret" + strconv.Itoa(i+1)
		}
	}
	sig := types.NewSignature(nil, nil, types.NewTuple(ret...), false)
	if useClosure {
		cb.NewClosureWith(sig).BodyStart(pkg)
	} else {
		cb.CallInlineClosureStart(sig, 0, false)
	}

	cb.NewVar(tyError, "_gop_err")
	err := cb.Scope().Lookup("_gop_err")

	for _, retVar := range ret {
		cb.VarRef(retVar)
	}
	cb.VarRef(err)
	cb.InternalStack().Push(x)
	cb.Assign(n+1, 1)

	cb.If().Val(err).CompareNil(gotoken.NEQ).Then()
	if v.Tok == token.NOT { // expr!
		cb.Val(pkg.Builtin().Ref("panic")).Val(err).Call(1).EndStmt() // TODO: wrap err
	} else if v.Default == nil { // expr?
		cb.Val(err).ReturnErr(true) // TODO: wrap err & return err
	} else { // expr?:val
		compileExpr(ctx, v.Default)
		cb.Return(1)
	}
	cb.End().Return(0).End()
	if useClosure {
		cb.Call(0)
	}
}

// -----------------------------------------------------------------------------
