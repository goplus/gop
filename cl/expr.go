/*
 Copyright 2020 The GoPlus Authors (goplus.org)

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
	"reflect"
	"strconv"
	"strings"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/ast/astutil"
	"github.com/goplus/gop/constant"
	"github.com/goplus/gop/exec.spec"
	"github.com/goplus/gop/token"
	"github.com/qiniu/x/ctype"
	"github.com/qiniu/x/errors"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

type compileMode = token.Token

const (
	lhsAssign compileMode = token.ASSIGN // leftHandSide = ...
	lhsDefine compileMode = token.DEFINE // leftHandSide := ...
)

// -----------------------------------------------------------------------------

func compileExprLHS(ctx *blockCtx, expr ast.Expr, mode compileMode) {
	switch v := expr.(type) {
	case *ast.Ident:
		compileIdentLHS(ctx, v.Name, mode)
	case *ast.IndexExpr:
		compileIndexExprLHS(ctx, v, mode)
	case *ast.SelectorExpr:
		compileSelectorExprLHS(ctx, v, mode)
	case *ast.StarExpr:
		compileStarExprLHS(ctx, v, mode)
	default:
		log.Panicln("compileExpr failed: unknown -", reflect.TypeOf(v))
	}
}

func compileExpr(ctx *blockCtx, expr ast.Expr) func() {
	switch v := expr.(type) {
	case *ast.Ident:
		return compileIdent(ctx, v, false)
	case *ast.BasicLit:
		return compileBasicLit(ctx, v)
	case *ast.CallExpr:
		return compileCallExpr(ctx, v, 0)
	case *ast.BinaryExpr:
		return compileBinaryExpr(ctx, v)
	case *ast.UnaryExpr:
		return compileUnaryExpr(ctx, v)
	case *ast.SelectorExpr:
		return compileSelectorExpr(ctx, nil, v, false)
	case *ast.ErrWrapExpr:
		return compileErrWrapExpr(ctx, v)
	case *ast.IndexExpr:
		return compileIndexExpr(ctx, v, false)
	case *ast.TwoValueIndexExpr:
		return compileIndexExpr(ctx, v.IndexExpr, true)
	case *ast.SliceExpr:
		return compileSliceExpr(ctx, v)
	case *ast.CompositeLit:
		return compileCompositeLit(ctx, v)
	case *ast.SliceLit:
		return compileSliceLit(ctx, v)
	case *ast.FuncLit:
		return compileFuncLit(ctx, v)
	case *ast.ParenExpr:
		return compileExpr(ctx, v.X)
	case *ast.ListComprehensionExpr:
		return compileListComprehensionExpr(ctx, v)
	case *ast.MapComprehensionExpr:
		return compileMapComprehensionExpr(ctx, v)
	case *ast.ArrayType:
		return compileArrayType(ctx, v)
	case *ast.Ellipsis:
		return compileEllipsis(ctx, v)
	case *ast.StarExpr:
		return compileStarExpr(ctx, v)
	case *ast.KeyValueExpr:
		panic("compileExpr: ast.KeyValueExpr unexpected")
	default:
		log.Panicln("compileExpr failed: unknown -", reflect.TypeOf(v))
		return nil
	}
}

func compileIdentLHS(ctx *blockCtx, name string, mode compileMode) {
	in := ctx.infer.Get(-1)
	var addr iVar
	if name == "_" {
		ctx.underscore++
		typ := boundType(in.(iValue))
		addr = ctx.insertVar(name, typ)
	} else {
		var err error
		addr, err = ctx.findVar(name)
		if mode == lhsDefine {
			addr, err = ctx.getCtxVar(name)
			if addr != nil {
				log.Panicf("compileIdentLHS failed: %s redeclared in this block\n", name)
			}
		}
		if err == nil {
			if mode == lhsDefine && !addr.inCurrentCtx(ctx) {
				log.Warn("requireVar: variable is shadowed -", name)
			}
		} else if mode == lhsAssign || err != ErrNotFound {
			log.Panicln("compileIdentLHS failed:", err, "-", name)
		} else {
			typ := boundType(in.(iValue))
			addr = ctx.insertVar(name, typ)
		}
	}

	typ := addr.getType()
	if ctx.indirect {
		typ = typ.Elem()
	}
	checkType(typ, in, ctx.out)
	ctx.infer.PopN(1)
	if v, ok := addr.(*execVar); ok {
		if mode == token.ASSIGN || mode == token.DEFINE {
			ctx.out.StoreVar(v.v)
		} else if op, ok := addrops[mode]; ok {
			typ := v.v.Type()
			if typ.Kind() == reflect.Ptr {
				ctx.out.LoadVar(v.v).AddrOp(kindOf(typ), op)
			} else {
				ctx.out.AddrVar(v.v).AddrOp(kindOf(typ), op)
			}
		} else {
			log.Panicln("compileIdentLHS failed: unknown op -", mode)
		}
	} else {
		if mode == token.ASSIGN || mode == token.DEFINE {
			if ctx.indirect {
				ctx.out.Load(addr.(*stackVar).index).AddrOp(kindOf(addr.(*stackVar).getType()), exec.OpAssign)
			} else {
				ctx.out.Store(addr.(*stackVar).index)
			}
		} else {
			panic("compileIdentLHS: todo")
		}
	}
}

var addrops = map[token.Token]exec.AddrOperator{
	token.ASSIGN:         exec.OpAssign,
	token.ADD_ASSIGN:     exec.OpAddAssign,
	token.SUB_ASSIGN:     exec.OpSubAssign,
	token.MUL_ASSIGN:     exec.OpMulAssign,
	token.QUO_ASSIGN:     exec.OpQuoAssign,
	token.REM_ASSIGN:     exec.OpModAssign,
	token.AND_ASSIGN:     exec.OpAndAssign,
	token.OR_ASSIGN:      exec.OpOrAssign,
	token.XOR_ASSIGN:     exec.OpXorAssign,
	token.SHL_ASSIGN:     exec.OpLshAssign,
	token.SHR_ASSIGN:     exec.OpRshAssign,
	token.AND_NOT_ASSIGN: exec.OpAndNotAssign,
	token.INC:            exec.OpInc,
	token.DEC:            exec.OpDec,
}

func compileIdent(ctx *blockCtx, ident *ast.Ident, compileByCallExpr bool) func() {
	name := ident.Name
	if name == "_" {
		log.Panicln("cannot use _ as value")
	}
	if sym, ok := ctx.find(name); ok {
		switch v := sym.(type) {
		case *constVal:
			c := newConstVal(v.v, v.Kind())
			ctx.infer.Push(c)
			return func() {
				pushConstVal(ctx.out, c)
			}
		case *execVar:
			typ := v.v.Type()
			kind := typ.Kind()
			ctx.infer.Push(&goValue{t: typ})
			ctx.resetFieldIndex()
			return func() {
				if ctx.takeAddr || (ctx.checkLoadAddr && kind != reflect.Slice && kind != reflect.Map) {
					ctx.out.AddrVar(v.v)
				} else {
					ctx.out.LoadVar(v.v)
				}
			}
		case *stackVar:
			ctx.infer.Push(&goValue{t: v.typ})
			kind := v.typ.Kind()
			ctx.resetFieldIndex()
			return func() {
				if ctx.takeAddr || (ctx.checkLoadAddr && kind != reflect.Slice && kind != reflect.Map) {
					ctx.out.Addr(v.index)
				} else {
					ctx.out.Load(v.index)
				}
			}
		case string: // pkgPath
			pkg := ctx.FindGoPackage(v)
			if pkg == nil {
				log.Panicln("compileIdent failed: package not found -", v)
			}
			ctx.infer.Push(&nonValue{pkg})
			return nil
		case *funcDecl:
			fn := newQlFunc(v)
			ctx.use(v)
			ctx.infer.Push(fn)
			return func() { // TODO: maybe slowly, use Closure instead of GoClosure
				ctx.out.GoClosure(fn.fi)
			}
		case *typeDecl:
			ctx.infer.Push(&nonValue{v.Type})
			return nil
		default:
			log.Panicln("compileIdent failed: unknown -", reflect.TypeOf(sym))
		}
	} else {
		if name == "iota" {
			c := newIotaValue()
			ctx.infer.Push(c)
			return nil
		}
		if addr, kind, ok := ctx.builtin.Find(name); ok {
			switch kind {
			case exec.SymbolVar:
			case exec.SymbolFunc, exec.SymbolFuncv:
				if compileByCallExpr {
					fn := newGoFunc(addr, kind, 0, ctx)
					ctx.infer.Push(fn)
					return nil
				} else {
					fn := newGoFunc(addr, kind, 0, ctx)
					ftyp := astutil.FuncType(fn.t)
					decl := funcToClosure(ctx, ident, ftyp)
					ctx.use(decl)
					ctx.infer.Push(newQlFunc(decl))
					return func() {
						ctx.out.GoClosure(decl.fi)
					}
				}
			}
			log.Panicln("compileIdent todo: var -", kind, addr)
		}
		if typ, ok := ctx.builtin.FindType(name); ok {
			ctx.infer.Push(&nonValue{typ})
			return nil
		}
		if ci, ok := ctx.builtin.FindConst(name); ok {
			return compileConst(ctx, ci.Kind, ci.Value)
		}
		if gi, ok := goinstrs[name]; ok {
			ctx.infer.Push(&nonValue{gi.instr})
			return nil
		}
		log.Panicln("compileIdent failed: unknown -", name)
	}
	return nil
}

func compileArrayType(ctx *blockCtx, v *ast.ArrayType) func() {
	typ := toArrayType(ctx, v)
	ctx.infer.Push(&nonValue{typ})
	return nil
}

func compileEllipsis(ctx *blockCtx, v *ast.Ellipsis) func() {
	if v.Elt != nil {
		log.Panicln("compileEllipsis: todo")
	}
	ctx.infer.Push(&constVal{v: constant.MakeInt64(-1), kind: astutil.ConstUnboundInt})
	return nil
}

func compileCompositeLit(ctx *blockCtx, v *ast.CompositeLit) func() {
	if v.Type == nil {
		return compileMapLit(ctx, v)
	}
	typ := toType(ctx, v.Type)

	switch kind := typ.Kind(); kind {
	case reflect.Slice, reflect.Array:
		var typSlice reflect.Type
		if t, ok := typ.(*unboundArrayType); ok {
			n := toBoundArrayLen(ctx, v)
			typSlice = reflect.ArrayOf(n, t.elem)
		} else {
			typSlice = typ.(reflect.Type)
		}
		ctx.infer.Push(&goValue{t: typSlice})
		return func() {
			nLen := toBoundArrayLen(ctx, v)
			n := -1
			elts := make([]ast.Expr, nLen)
			for _, elt := range v.Elts {
				switch e := elt.(type) {
				case *ast.KeyValueExpr:
					n = toInt(ctx, e.Key)
					elts[n] = e.Value
				default:
					n++
					elts[n] = e
				}
			}
			typElem := typSlice.Elem()
			for _, elt := range elts {
				if elt != nil {
					compileExpr(ctx, elt)()
					checkType(typElem, ctx.infer.Pop(), ctx.out)
				} else {
					ctx.out.Zero(typElem)
				}
			}
			ctx.out.MakeArray(typSlice, nLen)
		}
	case reflect.Map:
		typMap := typ.(reflect.Type)
		ctx.infer.Push(&goValue{t: typMap})
		return func() {
			typKey := typMap.Key()
			typVal := typMap.Elem()
			for _, elt := range v.Elts {
				switch e := elt.(type) {
				case *ast.KeyValueExpr:
					compileExpr(ctx, e.Key)()
					checkType(typKey, ctx.infer.Pop(), ctx.out)
					compileExpr(ctx, e.Value)()
					checkType(typVal, ctx.infer.Pop(), ctx.out)
				default:
					log.Panicln("compileCompositeLit: map requires key-value expr.")
				}
			}
			ctx.out.MakeMap(typMap, len(v.Elts))
		}
	case reflect.Struct:
		typStruct := typ.(reflect.Type)
		ctx.infer.Push(&goValue{t: typStruct})
		return func() {
			old := ctx.takeAddr
			ctx.takeAddr = false
			for i, elt := range v.Elts {
				switch e := elt.(type) {
				case *ast.KeyValueExpr:
					fieldName := e.Key.(*ast.Ident).Name
					field, _ := typStruct.FieldByName(fieldName)
					ctx.out.Push(field.Index[0])
					typVal := field.Type
					compileExpr(ctx, e.Value)()
					checkType(typVal, ctx.infer.Pop(), ctx.out)
				default:
					// ast.Expr
					field := typStruct.Field(i)
					ctx.out.Push(field.Index[0])
					typVal := field.Type
					compileExpr(ctx, elt)()
					checkType(typVal, ctx.infer.Pop(), ctx.out)
				}
			}
			ctx.takeAddr = old
			if ctx.takeAddr {
				ctx.out.Struct(reflect.PtrTo(typStruct), len(v.Elts))
			} else {
				ctx.out.Struct(typStruct, len(v.Elts))
			}
		}
	default:
		log.Panicln("compileCompositeLit failed: unknown -", reflect.TypeOf(typ))
		return nil
	}
}

func compileSliceLit(ctx *blockCtx, v *ast.SliceLit) func() {
	n := len(v.Elts)
	if n == 0 {
		ctx.infer.Push(&goValue{t: exec.TyEmptyInterfaceSlice})
		return func() {
			log.Debug("compileSliceLit:", exec.TyEmptyInterfaceSlice)
			ctx.out.MakeArray(exec.TyEmptyInterfaceSlice, 0)
		}
	}
	fnElts := make([]func(), n)
	elts := make([]interface{}, n)
	for i, elt := range v.Elts {
		fnElts[i] = compileExpr(ctx, elt)
		elts[i] = ctx.infer.Get(-1)
	}
	typElem := boundElementType(elts, 0, n, 1)
	if typElem == nil {
		typElem = exec.TyEmptyInterface
	}
	typSlice := reflect.SliceOf(typElem)
	ctx.infer.Ret(n, &goValue{t: typSlice})
	return func() {
		log.Debug("compileSliceLit:", typSlice)
		for _, fnElt := range fnElts {
			fnElt()
		}
		checkElementType(typElem, elts, 0, n, 1, ctx.out)
		ctx.out.MakeArray(typSlice, len(v.Elts))
	}
}

func compileForPhrase(parent *blockCtx, f ast.ForPhrase, noExecCtx bool) (*blockCtx, func(exprElt func())) {
	var typKey, typVal reflect.Type
	var varKey, varVal exec.Var
	var ctx = newNormBlockCtxEx(parent, noExecCtx)

	exprX := compileExpr(parent, f.X)
	typData := boundType(ctx.infer.Pop().(iValue))
	if f.Key != nil {
		switch kind := typData.Kind(); kind {
		case reflect.Slice, reflect.Array:
			typKey = exec.TyInt
		case reflect.Map:
			typKey = typData.Key()
		default:
			log.Panicln("compileListComprehensionExpr: require slice, array or map")
		}
		varKey = ctx.insertVar(f.Key.Name, typKey, true).v
	}
	if f.Value != nil {
		typVal = typData.Elem()
		varVal = ctx.insertVar(f.Value.Name, typVal, true).v
	}
	return ctx, func(exprElt func()) {
		ctx.nextFlow(nil, nil, "")
		defer func() {
			ctx.currentFlow = ctx.currentFlow.parent
		}()
		exprX()
		out := ctx.out
		c := ctx.NewForPhrase(typData)
		out.ForPhrase(c, varKey, varVal, !noExecCtx)
		if f.Cond != nil {
			compileExpr(ctx, f.Cond)()
			checkBool(ctx.infer.Pop())
			out.FilterForPhrase(c)
		}
		exprElt()
		out.EndForPhrase(c)
	}
}

func compileForPhrases(ctx *blockCtx, fors []ast.ForPhrase) (*blockCtx, []func(exprElt func())) {
	n := len(fors)
	fns := make([]func(exprElt func()), n)
	for i := n - 1; i >= 0; i-- {
		ctx, fns[i] = compileForPhrase(ctx, fors[i], true)
	}
	return ctx, fns
}

func compileListComprehensionExpr(parent *blockCtx, v *ast.ListComprehensionExpr) func() {
	ctx, fns := compileForPhrases(parent, v.Fors)
	exprElt := compileExpr(ctx, v.Elt)
	typElem := boundType(ctx.infer.Get(-1).(iValue))
	typSlice := reflect.SliceOf(typElem)
	ctx.infer.Ret(1, &goValue{t: typSlice})
	return func() {
		for _, v := range fns {
			e, fn := exprElt, v
			exprElt = func() { fn(e) }
		}
		c := ctx.NewComprehension(typSlice)
		ctx.out.ListComprehension(c)
		exprElt()
		ctx.out.EndComprehension(c)
	}
}

func compileMapComprehensionExpr(parent *blockCtx, v *ast.MapComprehensionExpr) func() {
	ctx, fns := compileForPhrases(parent, v.Fors)
	exprEltKey := compileExpr(ctx, v.Elt.Key)
	exprEltVal := compileExpr(ctx, v.Elt.Value)
	typEltKey := boundType(ctx.infer.Get(-2).(iValue))
	typEltVal := boundType(ctx.infer.Get(-1).(iValue))
	typMap := reflect.MapOf(typEltKey, typEltVal)
	exprElt := func() {
		exprEltKey()
		exprEltVal()
	}
	ctx.infer.Ret(2, &goValue{t: typMap})
	return func() {
		for _, v := range fns {
			e, fn := exprElt, v
			exprElt = func() { fn(e) }
		}
		c := ctx.NewComprehension(typMap)
		ctx.out.MapComprehension(c)
		exprElt()
		ctx.out.EndComprehension(c)
	}
}

func compileMapLit(ctx *blockCtx, v *ast.CompositeLit) func() {
	n := len(v.Elts) << 1
	if n == 0 {
		typMap := reflect.MapOf(exec.TyString, exec.TyEmptyInterface)
		ctx.infer.Push(&goValue{t: typMap})
		return func() {
			log.Debug("compileMapLit:", typMap)
			ctx.out.MakeMap(typMap, 0)
		}
	}
	fnElts := make([]func(), n)
	elts := make([]interface{}, n)
	for i, elt := range v.Elts {
		switch e := elt.(type) {
		case *ast.KeyValueExpr:
			fnElts[i<<1] = compileExpr(ctx, e.Key)
			elts[i<<1] = ctx.infer.Get(-1)
			fnElts[(i<<1)+1] = compileExpr(ctx, e.Value)
			elts[(i<<1)+1] = ctx.infer.Get(-1)
		default:
			log.Panicln("compileMapLit: map requires key-value expr.")
		}
	}
	typKey := boundElementType(elts, 0, n, 2)
	if typKey == nil {
		log.Panicln("compileMapLit: mismatched key type.")
	}
	typVal := boundElementType(elts, 1, n, 2)
	if typVal == nil {
		typVal = exec.TyEmptyInterface
	}
	typMap := reflect.MapOf(typKey, typVal)
	ctx.infer.Ret(n, &goValue{t: typMap})
	return func() {
		log.Debug("compileMapLit:", typMap)
		for _, fnElt := range fnElts {
			fnElt()
		}
		out := ctx.out
		checkElementType(typKey, elts, 0, n, 2, out)
		checkElementType(typVal, elts, 1, n, 2, out)
		out.MakeMap(typMap, len(v.Elts))
	}
}

func compileFuncLit(ctx *blockCtx, v *ast.FuncLit) func() {
	funCtx := newExecBlockCtx(ctx)
	decl := newFuncDecl("", nil, v.Type, v.Body, funCtx)
	ctx.use(decl)
	ctx.infer.Push(newQlFunc(decl))
	return func() { // TODO: maybe slowly, use Closure instead of GoClosure
		ctx.out.GoClosure(decl.fi)
	}
}

func compileBasicLit(ctx *blockCtx, v *ast.BasicLit) func() {
	kind, n := astutil.ToConst(v)
	return compileConst(ctx, kind, n)
}

func compileConst(ctx *blockCtx, kind astutil.ConstKind, n interface{}) func() {
	ret := newConstVal(n, kind)
	ctx.infer.Push(ret)
	return func() {
		pushConstVal(ctx.out, ret)
	}
}

func pushConstVal(b exec.Builder, c *constVal) {
	c.reserve = b.Reserve()
	if isConstBound(c.kind) {
		v := boundConst(c, exec.TypeFromKind(c.kind))
		c.reserve.Push(b, v)
	}
}

func compileUnaryExpr(ctx *blockCtx, v *ast.UnaryExpr) func() {
	exprX := compileExpr(ctx, v.X)
	x := ctx.infer.Get(-1)
	op := unaryOps[v.Op]
	if op == 0 {
		if v.Op == token.ADD { // +x
			return exprX
		}
		if v.Op == token.AND {
			vx := x.(iValue)
			t := reflect.TypeOf(reflect.New(vx.Type()).Interface())
			ret := &goValue{t: t}
			ctx.infer.Ret(1, ret)
			return func() {
				ctx.takeAddr = true
				exprX()
				ctx.takeAddr = false
			}
		}
		if v.Op == token.ARROW { // <- x
			vx := x.(iValue)
			if vx.Type().Kind() == reflect.Chan {
				ret := &goValue{t: vx.Type().Elem()}
				ctx.infer.Ret(1, ret)
				return func() {
					exprX()
					ctx.out.Recv()
				}
			}
		}
	}
	xcons, xok := x.(*constVal)
	if xok { // op <const>
		ret := unaryOp(v.Op, op, xcons)
		ctx.infer.Ret(1, ret)
		return func() {
			ret.reserve = ctx.out.Reserve()
		}
	}
	kind, ret := unaryOpResult(op, x)
	ctx.infer.Ret(1, ret)
	return func() {
		exprX()
		checkUnaryOp(kind, op, x, ctx.out)
		ctx.out.BuiltinOp(kind, op)
	}
}

func unaryOpResult(op exec.Operator, x interface{}) (exec.Kind, iValue) {
	vx := x.(iValue)
	if vx.NumValues() != 1 {
		log.Panicln("unaryOp: argument isn't an expr.")
	}
	kind := vx.Kind()
	if !isConstBound(kind) {
		log.Panicln("unaryOp: expect x aren't const values.")
	}
	i := op.GetInfo()
	kindRet := kind
	if i.Out != exec.SameAsFirst {
		kindRet = i.Out
	}
	return kind, &goValue{t: exec.TypeFromKind(kindRet)}
}

var unaryOps = [...]exec.Operator{
	token.SUB: exec.OpNeg,
	token.NOT: exec.OpLNot,
	token.XOR: exec.OpBitNot,
}

func isUnboundNumberType(kind exec.Kind) bool {
	return kind == exec.ConstUnboundInt ||
		kind == exec.ConstUnboundFloat ||
		kind == exec.ConstUnboundComplex
}

func isBoundNumberType(kind exec.Kind) bool {
	return kind >= exec.Int && kind <= exec.Complex128
}

func isCompareToken(op token.Token) bool {
	return op == token.EQL || op == token.NEQ ||
		op == token.LSS || op == token.LEQ ||
		op == token.GTR || op == token.GEQ
}

func compileBinaryExpr(ctx *blockCtx, v *ast.BinaryExpr) func() {
	exprX := compileExpr(ctx, v.X)
	exprY := compileExpr(ctx, v.Y)
	op := binaryOps[v.Op]
	x := ctx.infer.Get(-2)
	y := ctx.infer.Get(-1)
	xcons, xok := x.(*constVal)
	ycons, yok := y.(*constVal)
	opShift := (op == exec.OpLsh || op == exec.OpRsh)
	if xok && yok { // <const> op <const>
		if opShift && xcons.kind == exec.ConstUnboundFloat {
			xcons.v = extractUnboundInt(xcons.v.(constant.Value), xcons.kind)
			xcons.kind = exec.ConstUnboundInt
		}
		ret := binaryOp(v.Op, op, xcons, ycons)
		ctx.infer.Ret(2, ret)
		return func() {
			ret.reserve = ctx.out.Reserve()
		}
	} else if isCompareToken(v.Op) {
		xkind := x.(iValue).Kind()
		ykind := y.(iValue).Kind()
		if xok && isUnboundNumberType(xcons.kind) && isBoundNumberType(ykind) {
			xcons.kind = ykind
		} else if yok && isUnboundNumberType(ycons.kind) && isBoundNumberType(xkind) {
			ycons.kind = xkind
		}
	}

	var kind iKind
	var shift *shiftValue
	if opShift && xok && !isConstBound(xcons.kind) {
		kind = xcons.kind
		shift = &shiftValue{x: xcons, r: exec.InvalidReserved}
		ctx.infer.Ret(2, shift)
	} else {
		var ret iValue
		kind, ret = binaryOpResult(op, x, y)
		ctx.infer.Ret(2, ret)
	}
	return func() {
		var label exec.Label
		exprX()
		if b := (op == exec.OpLAnd); b || op == exec.OpLOr { // TODO: optimize to rm calling BuiltinOp
			label = ctx.NewLabel("")
			if b {
				ctx.out.JmpIf(exec.JcFalse|exec.JcNotPopMask, label)
			} else {
				ctx.out.JmpIf(exec.JcTrue|exec.JcNotPopMask, label)
			}
		}
		exprY()
		if shift != nil {
			if !isConstBound(xcons.kind) {
				shift.r = ctx.out.ReserveOpShift()
				shift.fnCheckType = func(typ reflect.Type) {
					kind := typ.Kind()
					xcons.v = boundConst(xcons, typ)
					checkBinaryOp(kind, op, x, y, ctx.out)
					if err := checkOpMatchType(op, x, y); err != nil {
						log.Panicf("invalid operator: %v (%v)", ctx.code(v), err)
					}
					ctx.out.ReservedAsOpShift(shift.r, kind, op)
				}
				if label != nil {
					ctx.out.Label(label)
				}
				return
			}
			kind = xcons.kind
		}

		checkBinaryOp(kind, op, x, y, ctx.out)
		if err := checkOpMatchType(op, x, y); err != nil {
			log.Panicf("invalid operator: %v (%v)", ctx.code(v), err)
		}
		ctx.out.BuiltinOp(kind, op)
		if label != nil {
			ctx.out.Label(label)
		}
	}
}

func binaryOpResult(op exec.Operator, x, y interface{}) (exec.Kind, iValue) {
	vx := x.(iValue)
	vy := y.(iValue)
	if vx.NumValues() != 1 || vy.NumValues() != 1 {
		log.Panicln("binaryOp: argument isn't an expr.")
	}
	kind := vx.Kind()
	if !isConstBound(kind) {
		kind = vy.Kind()
		if xshv, xok := x.(*shiftValue); xok {
			var unbound bool
			if !isConstBound(kind) {
				if kind != xshv.Kind() {
					unbound = true
				}
				if c, ok := y.(*constVal); ok {
					kind = c.boundKind()
					c.v = boundConst(c, c.boundType())
					c.kind = kind
				} else if shv, ok := y.(*shiftValue); ok && shv.Kind() == exec.ConstUnboundInt {
					kind = exec.Int
					shv.bound(exec.TyInt)
				}
			}
			if !unbound {
				xshv.bound(vy.Type())
			}
		} else if !isConstBound(kind) {
			log.Panicln("binaryOp: expect x, y aren't const values either.")
		}
	}
	i := op.GetInfo()
	kindRet := kind
	if i.Out != exec.SameAsFirst {
		kindRet = i.Out
	}
	return kind, &goValue{t: exec.TypeFromKind(kindRet)}
}

var binaryOps = [...]exec.Operator{
	token.ADD:     exec.OpAdd,
	token.SUB:     exec.OpSub,
	token.MUL:     exec.OpMul,
	token.QUO:     exec.OpQuo,
	token.REM:     exec.OpMod,
	token.AND:     exec.OpAnd,
	token.OR:      exec.OpOr,
	token.XOR:     exec.OpXor,
	token.AND_NOT: exec.OpAndNot,
	token.SHL:     exec.OpLsh,
	token.SHR:     exec.OpRsh,
	token.LSS:     exec.OpLT,
	token.LEQ:     exec.OpLE,
	token.GTR:     exec.OpGT,
	token.GEQ:     exec.OpGE,
	token.EQL:     exec.OpEQ,
	token.NEQ:     exec.OpNE,
	token.LAND:    exec.OpLAnd,
	token.LOR:     exec.OpLOr,
}

func compileCallExpr(ctx *blockCtx, v *ast.CallExpr, ct callType) func() {
	var exprFun func()
	switch f := v.Fun.(type) {
	case *ast.Ident:
		exprFun = compileIdent(ctx, f, true)
	case *ast.SelectorExpr:
		exprFun = compileSelectorExpr(ctx, v, f, true)
	default:
		exprFun = compileExpr(ctx, f)
	}
	return compileCallExprCall(ctx, exprFun, v, ct)
}

func compileCallExprCall(ctx *blockCtx, exprFun func(), v *ast.CallExpr, ct callType) func() {
	fn := ctx.infer.Pop()
	ctx.resetFieldIndex()
	switch vfn := fn.(type) {
	case *qlFunc:
		if ct == callExpr {
			ret := vfn.Results()
			ctx.infer.Push(ret)
		}
		return func() {
			var isMethod int
			if vfn.recv != nil {
				isMethod = 1
				exprX := compileExpr(ctx, v.Fun.(*ast.SelectorExpr).X)
				recv := ctx.infer.Get(-1).(*goValue)

				if astutil.ToRecv(vfn.recv).Pointer == 0 {
					exprX()
					if recv.Kind() == reflect.Ptr {
						recv.t = recv.t.Elem()
						ctx.infer.Ret(1, recv)
						ctx.out.AddrOp(recv.t.Kind(), exec.OpAddrVal) // Ptr => Elem()
					}
				} else {
					ctx.checkLoadAddr = true
					exprX()
					ctx.checkLoadAddr = false
					if recv.Kind() != reflect.Ptr {
						recv.t = reflect.PtrTo(recv.t)
						ctx.infer.Ret(1, recv)
					}
				}
			}
			for _, arg := range v.Args {
				compileExpr(ctx, arg)()
			}
			arity := checkFuncCall(vfn.Proto(), isMethod, v, ctx)
			fun := vfn.FuncInfo()
			if fun.IsVariadic() {
				builder(ctx, ct).CallFuncv(fun, len(v.Args), arity)
			} else {
				builder(ctx, ct).CallFunc(fun, len(v.Args))
			}
		}
	case *goFunc:
		var ret iValue
		if ct == callExpr {
			ret = vfn.Results()
			ctx.infer.Push(ret)
		}
		return func() {
			if vfn.isMethod != 0 {
				exprX := compileExpr(ctx, v.Fun.(*ast.SelectorExpr).X)
				recv := ctx.infer.Get(-1).(*goValue)
				if vfn.Type().In(0).Kind() != reflect.Ptr {
					exprX()
					if recv.Kind() == reflect.Ptr {
						recv.t = recv.t.Elem()
						ctx.infer.Ret(1, recv)
						ctx.out.AddrOp(recv.t.Kind(), exec.OpAddrVal) // Ptr => Elem()
					}
				} else {
					if recv.Kind() == reflect.Ptr {
						exprX()
					} else {
						ctx.checkLoadAddr = true
						exprX()
						ctx.checkLoadAddr = false
						recv.t = reflect.PtrTo(recv.t)
						ctx.infer.Ret(1, recv)
					}
				}
			}
			for _, arg := range v.Args {
				compileExpr(ctx, arg)()
			}
			nexpr := len(v.Args) + vfn.isMethod
			arity := checkFuncCall(vfn.Proto(), vfn.isMethod, v, ctx)
			switch vfn.kind {
			case exec.SymbolFunc:
				builder(ctx, ct).CallGoFunc(exec.GoFuncAddr(vfn.addr), nexpr)
			case exec.SymbolFuncv:
				builder(ctx, ct).CallGoFuncv(exec.GoFuncvAddr(vfn.addr), nexpr, arity)
			}
			if ct == callExpr {
			}
		}
	case *goValue:
		if vfn.t.Kind() != reflect.Func {
			log.Panicln("compileCallExpr failed: call a non function.")
		}
		if ct == callExpr {
			ret := newFuncResults(vfn.t)
			ctx.infer.Push(ret)
		}
		return func() {
			for _, arg := range v.Args {
				compileExpr(ctx, arg)()
			}
			exprFun()
			arity, ellipsis := checkFuncCall(vfn.t, 0, v, ctx), false
			if arity == -1 {
				arity, ellipsis = len(v.Args), true
			}
			builder(ctx, ct).CallGoClosure(len(v.Args), arity, ellipsis)
		}
	case *nonValue:
		switch nv := vfn.v.(type) {
		case goInstr:
			return nv(ctx, v, ct)
		case reflect.Type:
			if ct != callExpr {
				log.Panicf("%s requires function call, not conversion\n", gCallTypes[ct])
			}
			return compileTypeCast(nv, ctx, v)
		}
	}
	log.Panicln("compileCallExpr failed: unknown -", reflect.TypeOf(fn))
	return nil
}

func builder(ctx *blockCtx, ct callType) (out exec.Builder) {
	switch out = ctx.out; ct {
	case callByDefer:
		return out.Defer()
	case callByGo:
		return out.Go()
	}
	return
}

func compileIndexExprLHS(ctx *blockCtx, v *ast.IndexExpr, mode compileMode) {
	if mode == lhsDefine {
		log.Panicln("compileIndexExprLHS: `:=` can't be used for index expression")
	}
	val := ctx.infer.Get(-1)

	exprX := compileExpr(ctx, v.X)

	typ := ctx.infer.Get(-1).(iValue).Type()
	typElem := typ.Elem()
	if typ.Kind() == reflect.Array {
		ctx.checkLoadAddr = true
	}
	if typ.Kind() == reflect.Ptr {
		if typElem.Kind() != reflect.Array {
			logPanic(ctx, v, `type %v does not support indexing`, typ)
		}
		typ = typElem
		typElem = typElem.Elem()
	}
	exprX()
	ctx.checkLoadAddr = false
	if ctx.indirect {
		typElem = typElem.Elem()
	}
	if cons, ok := val.(*constVal); ok {
		cons.bound(typElem, ctx.out)
	} else if t := val.(iValue).Type(); t != typElem && typElem.Kind() != reflect.Interface {
		log.Panicf("compileIndexExprLHS: can't assign `%v`[i] = `%v`\n", typ, t)
	}
	exprIdx := compileExpr(ctx, v.Index)
	i := ctx.infer.Get(-1)
	ctx.infer.PopN(3)

	switch typ.Kind() {
	case reflect.Slice, reflect.Array:
		if cons, ok := i.(*constVal); ok {
			n := boundConst(cons, exec.TyInt)
			if ctx.indirect {
				ctx.out.Index(n.(int)).AddrOp(kindOf(typElem), exec.OpAssign)
			} else {
				ctx.out.SetIndex(n.(int))
			}
			return
		}
		exprIdx()
		if typIdx := i.(iValue).Type(); typIdx != exec.TyInt {
			if typIdx.ConvertibleTo(exec.TyInt) {
				ctx.out.TypeCast(typIdx, exec.TyInt)
			} else {
				log.Panicln("compileIndexExprLHS: index expression value type is invalid")
			}
		}
		if ctx.indirect {
			ctx.out.Index(-1).AddrOp(kindOf(typElem), exec.OpAssign)
		} else {
			ctx.out.SetIndex(-1)
		}
	case reflect.Map:
		exprIdx()
		typIdx := typ.Key()
		if cons, ok := i.(*constVal); ok {
			cons.bound(typIdx, ctx.out)
		}
		if t := i.(iValue).Type(); t != typIdx {
			logIllTypeMapIndexPanic(ctx, v, t, typIdx)
		}
		if ctx.indirect {
			ctx.out.MapIndex(false).AddrOp(kindOf(typElem), exec.OpAssign)
		} else {
			ctx.out.SetMapIndex()
		}
	default:
		log.Panicln("compileIndexExprLHS: unknown -", typ)
	}
}

func compileSliceExpr(ctx *blockCtx, v *ast.SliceExpr) func() { // x[i:j:k]
	var kind reflect.Kind
	exprX := compileExpr(ctx, v.X)
	x := ctx.infer.Get(-1)
	typ := x.(iValue).Type()
	if kind = typ.Kind(); kind == reflect.Array {
		typ = reflect.SliceOf(typ.Elem())
		ctx.infer.Ret(1, &goValue{typ})
	}
	return func() {
		if kind == reflect.Array {
			ctx.checkLoadAddr = true
		}
		exprX()
		ctx.checkLoadAddr = false
		i, j, k := exec.SliceDefaultIndex, exec.SliceDefaultIndex, exec.SliceDefaultIndex
		if v.Low != nil {
			i = compileIdx(ctx, v.Low, exec.SliceConstIndexLast, kind)
		}
		if v.High != nil {
			j = compileIdx(ctx, v.High, exec.SliceConstIndexLast, kind)
		}
		if v.Max != nil {
			k = compileIdx(ctx, v.Max, exec.SliceConstIndexLast, kind)
		}
		if v.Slice3 {
			ctx.out.Slice3(i, j, k)
		} else {
			ctx.out.Slice(i, j)
		}
	}
}

func compileIdx(ctx *blockCtx, v ast.Expr, nlast int, kind reflect.Kind) int {
	expr := compileExpr(ctx, v)
	i := ctx.infer.Pop()
	if cons, ok := i.(*constVal); ok {
		nv := boundConst(cons, exec.TyInt)
		n := nv.(int)
		if n <= nlast {
			return n
		}
		ctx.out.Push(n)
		return -1
	} else if shv, ok := i.(*shiftValue); ok {
		shv.bound(exec.TyInt)
	}

	expr()
	if typIdx := i.(iValue).Type(); typIdx != exec.TyInt {
		if typIdx.ConvertibleTo(exec.TyInt) {
			ctx.out.TypeCast(typIdx, exec.TyInt)
		} else {
			logNonIntegerIdxPanic(ctx, v, kind)
		}
	}
	return -1
}

func compileIndexExpr(ctx *blockCtx, v *ast.IndexExpr, twoValue bool) func() { // x[i]
	var kind reflect.Kind
	var typElem reflect.Type
	exprX := compileExpr(ctx, v.X)
	x := ctx.infer.Get(-1)
	typ := x.(iValue).Type()
	kind = typ.Kind()
	if kind == reflect.Ptr {
		typ = typ.Elem()
		if kind = typ.Kind(); kind != reflect.Array {
			logPanic(ctx, v, `type *%v does not support indexing`, typ)
		}
	}
	if kind == reflect.String {
		typElem = exec.TyByte
	} else {
		typElem = typ.Elem()
	}
	ctx.infer.Ret(1, &goValue{typElem})
	ctx.resetFieldIndex()
	return func() {
		if ctx.takeAddr {
			if kind == reflect.String || kind == reflect.Map {
				_, info := ctx.getCodeInfo(v)
				log.Panicf("cannot take the address of %v\n", info)
			} else if kind == reflect.Slice {
				ctx.takeAddr = false
				exprX()
				ctx.takeAddr = true
			} else {
				exprX()
			}
		} else {
			exprX()
		}
		switch kind {
		case reflect.String, reflect.Slice, reflect.Array:
			n := compileIdx(ctx, v.Index, 1<<30, kind)
			if ctx.takeAddr || ctx.checkLoadAddr {
				ctx.out.AddrIndex(n)
			} else {
				ctx.out.Index(n)
			}
		case reflect.Map:
			typIdx := typ.Key()
			compileExpr(ctx, v.Index)()
			i := ctx.infer.Pop()
			if cons, ok := i.(*constVal); ok {
				cons.bound(typIdx, ctx.out)
			}
			if t := i.(iValue).Type(); t != typIdx {
				logIllTypeMapIndexPanic(ctx, v, t, typIdx)
			}
			if twoValue {
				ctx.infer.Push(&goValue{t: exec.TyBool})
			}
			ctx.out.MapIndex(twoValue)
		default:
			log.Panicln("compileIndexExpr: unknown -", typ)
		}
	}
}

func compileErrWrapExpr(ctx *blockCtx, v *ast.ErrWrapExpr) func() {
	exprX := compileExpr(ctx, v.X)
	x := ctx.infer.Get(-1).(iValue)
	nx := x.NumValues()
	if nx < 1 || !x.Value(nx-1).Type().Implements(exec.TyError) {
		log.Panicln("last output parameter doesn't implement `error` interface")
	}
	ctx.infer.Ret(1, &wrapValue{x})
	return func() {
		exprX()
		if v.Default == nil { // expr? or expr!
			var fun = ctx.fun
			var ok bool
			var retErr exec.Var
			if v.Tok == token.QUESTION {
				if retErr, ok = returnErr(fun); !ok {
					log.Panicln("used `expr?` in a function that last output parameter is not an error")
				}
			} else if v.Tok == token.NOT {
				retErr = nil
			}
			pos, code := ctx.getCodeInfo(v)
			fn, narg := getFuncInfo(fun)
			frame := &errors.Frame{
				Pkg:  ctx.pkg.Name,
				Func: fn,
				Code: code,
				File: pos.Filename,
				Line: pos.Line,
			}
			ctx.out.ErrWrap(nx, retErr, frame, narg)
			return
		}
		if nx != 2 {
			log.Panicln("compileErrWrapExpr: output parameters count must be 2")
		}
		label := ctx.NewLabel("")
		ctx.out.WrapIfErr(nx, label)
		compileExpr(ctx, v.Default)()
		checkType(x.Value(0).Type(), ctx.infer.Pop(), ctx.out)
		ctx.out.Label(label)
	}
}

func returnErr(fun exec.FuncInfo) (retErr exec.Var, ok bool) {
	if fun == nil {
		return
	}
	n := fun.NumOut()
	if n == 0 {
		return
	}
	retErr = fun.Out(n - 1)
	ok = retErr.Type() == exec.TyError
	return
}

func getFuncInfo(fun exec.FuncInfo) (name string, narg int) {
	if fun != nil {
		return fun.Name(), fun.NumIn()
	}
	return "main", 0
}

func compileSelectorExprLHS(ctx *blockCtx, v *ast.SelectorExpr, mode compileMode) {
	if mode == lhsDefine {
		log.Panicln("compileSelectorExprLHS: `:=` can't be used for index expression")
	}
	in := ctx.infer.Get(-1)
	exprX := compileExpr(ctx, v.X)

	x := ctx.infer.Get(-1)
	ctx.infer.PopN(2)
	switch vx := x.(type) {
	case *nonValue:
		switch nv := vx.v.(type) {
		case exec.GoPackage:
			if c, ok := nv.FindConst(v.Sel.Name); ok {
				log.Panicln("cannot assign to ", c.Pkg.PkgPath()+"."+c.Name)
			}
			addr, kind, ok := nv.Find(v.Sel.Name)
			if !ok {
				log.Panicln("compileSelectorExprLHS: not found -", nv.PkgPath(), v.Sel.Name)
			}
			switch kind {
			case exec.SymbolVar:
				info := ctx.GetGoVarInfo(exec.GoVarAddr(addr))
				t := reflect.TypeOf(info.This).Elem()
				checkType(t, in, ctx.out)
				ctx.out.StoreGoVar(exec.GoVarAddr(addr))
			default:
				log.Panicln("compileSelectorExprLHS: unknown GoPackage symbol kind -", kind)
			}
		default:
			log.Panicln("compileSelectorExprLHS: unknown nonValue -", reflect.TypeOf(nv))
		}
	case *goValue:
		_, t := countPtr(vx.t)
		name := v.Sel.Name
		if t.PkgPath() != "" && ast.IsExported(name) || t.PkgPath() == "" {
			if t.Kind() == reflect.Struct {
				if sf, ok := t.FieldByName(name); ok {
					typ := sf.Type
					if ctx.indirect {
						typ = typ.Elem()
					}
					checkType(typ, in, ctx.out)
					if ctx.fieldIndex == nil {
						ctx.fieldStructType = vx.t
					}
					fieldIndex := append(ctx.fieldIndex, sf.Index...)
					fieldStructType := ctx.fieldStructType
					ctx.checkLoadAddr = true
					if ctx.fieldExprX != nil {
						ctx.fieldExprX()
					} else {
						exprX()
					}
					ctx.checkLoadAddr = false
					if ctx.indirect {
						ctx.out.LoadField(fieldStructType, fieldIndex).AddrOp(kindOf(typ), exec.OpAssign)
					} else {
						ctx.out.StoreField(fieldStructType, fieldIndex)
					}
				}
			}
		} else if t.PkgPath() != "" && !ast.IsExported(name) {
			log.Panicf("%s.%s.%s undefined (cannot refer to unexported field or method %s)\n", t.PkgPath(), t.Name(), name, name)
		}
	default:
		log.Panicln("compileSelectorExprLHS failed: unknown -", reflect.TypeOf(vx))
	}
	_ = exprX
}

func compileStarExprLHS(ctx *blockCtx, v *ast.StarExpr, mode compileMode) {
	ctx.indirect = true
	compileExprLHS(ctx, v.X, mode)
	ctx.indirect = false
}

func funcToClosure(ctx *blockCtx, fun ast.Expr, ftyp *ast.FuncType) *funcDecl {
	typ := &ast.FuncType{Params: &ast.FieldList{}, Results: ftyp.Results}
	var args []ast.Expr
	var ellipsis bool
	for i, field := range ftyp.Params.List {
		if _, ok := field.Type.(*ast.Ellipsis); ok {
			ellipsis = true
		}
		if field.Names != nil {
			for _, name := range field.Names {
				args = append(args, name)
			}
			typ.Params.List = append(typ.Params.List, field)
		} else {
			ident := &ast.Ident{Name: strconv.Itoa(i)}
			args = append(args, ident)
			typ.Params.List = append(typ.Params.List, &ast.Field{
				Names: []*ast.Ident{ident},
				Type:  field.Type,
			})
		}
	}
	call := &ast.CallExpr{Fun: fun, Args: args}
	if ellipsis {
		call.Ellipsis++
	}
	var body *ast.BlockStmt
	if typ.Results == nil {
		body = &ast.BlockStmt{List: []ast.Stmt{&ast.ExprStmt{X: call}}}
	} else {
		body = &ast.BlockStmt{List: []ast.Stmt{&ast.ReturnStmt{Return: fun.Pos(), Results: []ast.Expr{call}}}}
	}
	funCtx := newExecBlockCtx(ctx)
	return newFuncDecl("", nil, typ, body, funCtx)
}

func toElem(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func isUserStruct(t reflect.Type) bool {
	t = toElem(t)
	return t.Kind() == reflect.Struct && (t.PkgPath() == "" || t.PkgPath() == "main")
}

func findUserStructAnonymous(ctx *blockCtx, t reflect.Type, name string) []string {
	t = toElem(t)
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if sf.Anonymous {
			var found bool
			if isUserStruct(sf.Type) {
				_, found = ctx.findMethod(sf.Type, name)
				if !found {
					if names := findUserStructAnonymous(ctx, sf.Type, name); names != nil {
						return append([]string{sf.Name}, names...)
					}
				}
			} else {
				_, found = sf.Type.MethodByName(name)
			}
			if found {
				return []string{sf.Name}
			}
		}
	}
	return nil
}

func compileSelectorExpr(ctx *blockCtx, call *ast.CallExpr, v *ast.SelectorExpr, compileByCallExpr bool) func() {
	exprX := compileExpr(ctx, v.X)
	if v.Sel == nil {
		return exprX
	}
	x := ctx.infer.Get(-1)
	switch vx := x.(type) {
	case *nonValue:
		switch nv := vx.v.(type) {
		case exec.GoPackage:
			name := strings.Title(v.Sel.Name)
			if c, ok := nv.FindConst(name); ok {
				ret := newConstVal(c.Value, c.Kind)
				ctx.infer.Ret(1, ret)
				return func() {
					pushConstVal(ctx.out, ret)
				}
			}
			if t, ok := nv.FindType(name); ok {
				ctx.infer.Ret(1, &nonValue{t})
				return nil
			}
			addr, kind, ok := nv.Find(name)
			if !ok {
				log.Panicln("compileSelectorExpr: not found -", nv.PkgPath(), name)
			}
			switch kind {
			case exec.SymbolFunc, exec.SymbolFuncv:
				if nv.PkgPath() == "unsafe" {
					if gi, ok := goinstrs["unsafe."+name]; ok {
						ctx.infer.Push(&nonValue{gi.instr})
						return nil
					}
				}
				if compileByCallExpr {
					fn := newGoFunc(addr, kind, 0, ctx)
					ctx.infer.Ret(1, fn)
					return nil
				} else {
					ctx.infer.Pop()
					fn := newGoFunc(addr, kind, 0, ctx)
					ftyp := astutil.FuncType(fn.t)
					decl := funcToClosure(ctx, v, ftyp)
					ctx.use(decl)
					ctx.infer.Push(newQlFunc(decl))
					return func() {
						ctx.out.GoClosure(decl.fi)
					}
				}
			case exec.SymbolVar:
				info := ctx.GetGoVarInfo(exec.GoVarAddr(addr))
				vt := reflect.ValueOf(info.This)
				typ := vt.Elem().Type()
				kind := typ.Kind()
				ctx.infer.Ret(1, &goValue{t: typ})
				ctx.resetFieldIndex()
				return func() {
					if ctx.takeAddr || (ctx.checkLoadAddr && kind != reflect.Slice && kind != reflect.Map) {
						ctx.out.AddrGoVar(exec.GoVarAddr(addr))
					} else {
						ctx.out.LoadGoVar(exec.GoVarAddr(addr))
					}
				}
			default:
				log.Panicln("compileSelectorExpr: unknown GoPackage symbol kind -", kind)
			}
		default:
			log.Panicln("compileSelectorExpr: unknown nonValue -", reflect.TypeOf(nv))
		}
	case *goValue:
		n, t := countPtr(vx.t)
		autoCall := false
		name := v.Sel.Name
		if t.PkgPath() != "" && ast.IsExported(name) || t.PkgPath() == "" {
			if t.Kind() == reflect.Struct {
				if sf, ok := t.FieldByName(name); ok {
					ctx.infer.Ret(1, &goValue{t: sf.Type})
					if ctx.fieldIndex == nil {
						ctx.fieldExprX = exprX
						ctx.fieldStructType = vx.t
					}
					ctx.fieldIndex = append(ctx.fieldIndex, sf.Index...)
					fieldIndex := ctx.fieldIndex
					fieldExprX := ctx.fieldExprX
					fieldStructType := ctx.fieldStructType
					return func() {
						if fieldExprX != nil {
							fieldExprX()
						}
						if ctx.takeAddr || ctx.checkLoadAddr {
							ctx.out.AddrField(fieldStructType, fieldIndex)
						} else {
							ctx.out.LoadField(fieldStructType, fieldIndex)
						}
					}
				}
			}
			if fDecl, ok := ctx.findMethod(t, name); ok {
				if compileByCallExpr {
					ctx.infer.Pop()
					fn := newQlFunc(fDecl)
					ctx.use(fDecl)
					ctx.infer.Push(fn)
					return nil
				} else {
					ctx.infer.Pop()
					decl := funcToClosure(ctx, v, fDecl.typ)
					ctx.use(decl)
					ctx.infer.Push(newQlFunc(decl))
					return func() {
						ctx.out.GoClosure(decl.fi)
					}
				}
			}
			if call != nil && isUserStruct(t) {
				if names := findUserStructAnonymous(ctx, t, name); names != nil {
					ctx.infer.Pop()
					x := &ast.SelectorExpr{X: v.X}
					for i := 0; i < len(names)-1; i++ {
						x.X = &ast.SelectorExpr{X: x.X, Sel: &ast.Ident{Name: names[i]}}
					}
					x.Sel = &ast.Ident{Name: names[len(names)-1]}
					fun := &ast.SelectorExpr{X: x, Sel: v.Sel}
					call.Fun = fun
					return compileSelectorExpr(ctx, call, fun, compileByCallExpr)
				}
			}
		}
		_, toptr, ok := findMethod(t, name)
		if !ok && isLower(name) {
			name = strings.Title(name)
			if _, toptr, ok = findMethod(t, name); ok {
				v.Sel.Name = name
				autoCall = !compileByCallExpr
			}
		}
		if !ok {
			log.Panicln("compileSelectorExpr: symbol not found -", vx.t, v.Sel.Name)
		}
		if n > 1 {
			log.Panicf("calling method %v with receiver %v (type %v) requires explicit dereference.", v.Sel.Name, ctx.code(v.X), vx.t)
		}
		pkgPath := t.PkgPath()
		pkg := ctx.FindGoPackage(pkgPath)
		if pkg == nil {
			log.Panicln("compileSelectorExpr failed: package not found -", pkgPath)
		}
		var fnname string
		if toptr {
			fnname = "(*" + t.Name() + ")." + name
		} else {
			fnname = "(" + t.Name() + ")." + name
		}
		addr, kind, ok := pkg.Find(fnname)
		if !ok {
			log.Panicln("compileSelectorExpr: method not found -", fnname)
		}
		if !compileByCallExpr && !autoCall {
			ctx.infer.Pop()
			fn := newGoFunc(addr, kind, 1, ctx)
			ftyp := astutil.FuncType(fn.t)
			ftyp.Params.List = ftyp.Params.List[1:]
			decl := funcToClosure(ctx, v, ftyp)
			ctx.use(decl)
			ctx.infer.Push(newQlFunc(decl))
			return func() {
				ctx.out.GoClosure(decl.fi)
			}
		}
		ctx.infer.Ret(1, newGoFunc(addr, kind, 1, ctx))
		if autoCall { // change AST tree
			copy := *v
			call := &ast.CallExpr{Fun: &copy}
			v.X = call
			v.Sel = nil
			return compileCallExprCall(ctx, nil, call, 0)
		}
		return nil
	default:
		log.Panicln("compileSelectorExpr failed: unknown -", reflect.TypeOf(vx))
	}
	return nil
}

func compileStarExpr(ctx *blockCtx, v *ast.StarExpr) func() {
	exprX := compileExpr(ctx, v.X)
	x := ctx.infer.Get(-1)
	switch vx := x.(type) {
	case *nonValue:
		switch t := vx.v.(type) {
		case reflect.Type:
			ctx.infer.Ret(1, &nonValue{reflect.PtrTo(t)})
			return nil
		}
	case *goValue:
		if vx.Kind() == reflect.Ptr {
			ctx.infer.Ret(1, &goValue{vx.t.Elem()})
		}
		return func() {
			exprX()
			ctx.out.AddrOp(kindOf(vx.t), exec.OpAddrVal)
		}
	}
	log.Panicln("compileStarExpr failed: unknown -", reflect.TypeOf(x))
	return nil
}

func isLower(name string) bool {
	for _, c := range name {
		return ctype.Is(ctype.LOWER, c)
	}
	return false
}

func countPtr(t reflect.Type) (int, reflect.Type) {
	n := 0
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		n++
	}
	return n, t
}

func findMethod(t reflect.Type, name string) (method reflect.Method, toptr bool, found bool) {
	method, found = t.MethodByName(name)
	if !found && t.Kind() == reflect.Struct {
		t = reflect.PtrTo(t)
		toptr = true
		method, found = t.MethodByName(name)
	}
	return
}

// -----------------------------------------------------------------------------
