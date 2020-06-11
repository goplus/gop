/*
 Copyright 2020 Qiniu Cloud (qiniu.com)

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
	"strings"
	"syscall"

	"github.com/qiniu/goplus/ast"
	"github.com/qiniu/goplus/ast/astutil"
	"github.com/qiniu/goplus/exec.spec"
	"github.com/qiniu/goplus/token"
	"github.com/qiniu/x/errors"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

type compleMode = token.Token

const (
	lhsAssign compleMode = token.ASSIGN // leftHandSide = ...
	lhsDefine compleMode = token.DEFINE // leftHandSide := ...
)

// -----------------------------------------------------------------------------

func compileExprLHS(ctx *blockCtx, expr ast.Expr, mode compleMode) {
	switch v := expr.(type) {
	case *ast.Ident:
		compileIdentLHS(ctx, v.Name, mode)
	case *ast.IndexExpr:
		compileIndexExprLHS(ctx, v, mode)
	default:
		log.Panicln("compileExpr failed: unknown -", reflect.TypeOf(v))
	}
}

func compileExpr(ctx *blockCtx, expr ast.Expr) func() {
	switch v := expr.(type) {
	case *ast.Ident:
		return compileIdent(ctx, v.Name)
	case *ast.BasicLit:
		return compileBasicLit(ctx, v)
	case *ast.CallExpr:
		return compileCallExpr(ctx, v)
	case *ast.BinaryExpr:
		return compileBinaryExpr(ctx, v)
	case *ast.UnaryExpr:
		return compileUnaryExpr(ctx, v)
	case *ast.SelectorExpr:
		return compileSelectorExpr(ctx, v)
	case *ast.ErrWrapExpr:
		return compileErrWrapExpr(ctx, v)
	case *ast.IndexExpr:
		return compileIndexExpr(ctx, v)
	case *ast.SliceExpr:
		return compileSliceExpr(ctx, v)
	case *ast.CompositeLit:
		return compileCompositeLit(ctx, v)
	case *ast.SliceLit:
		return compileSliceLit(ctx, v)
	case *ast.FuncLit:
		return compileFuncLit(ctx, v)
	case *ast.ListComprehensionExpr:
		return compileListComprehensionExpr(ctx, v)
	case *ast.MapComprehensionExpr:
		return compileMapComprehensionExpr(ctx, v)
	case *ast.Ellipsis:
		return compileEllipsis(ctx, v)
	case *ast.KeyValueExpr:
		panic("compileExpr: ast.KeyValueExpr unexpected")
	default:
		log.Panicln("compileExpr failed: unknown -", reflect.TypeOf(v))
		return nil
	}
}

func compileIdentLHS(ctx *blockCtx, name string, mode compleMode) {
	in := ctx.infer.Get(-1)
	addr, err := ctx.findVar(name)
	if err == nil {
		if mode == lhsDefine && !addr.inCurrentCtx(ctx) {
			log.Warn("requireVar: variable is shadowed -", name)
		}
	} else if mode == lhsAssign || err != syscall.ENOENT {
		log.Panicln("compileIdentLHS failed:", err, "-", name)
	} else {
		typ := boundType(in.(iValue))
		addr = ctx.insertVar(name, typ)
	}
	checkType(addr.getType(), in, ctx.out)
	ctx.infer.PopN(1)
	if v, ok := addr.(*execVar); ok {
		if mode == token.ASSIGN || mode == token.DEFINE {
			ctx.out.StoreVar(v.v)
		} else if op, ok := addrops[mode]; ok {
			ctx.out.AddrVar(v.v).AddrOp(v.v.Type().Kind(), op)
		} else {
			log.Panicln("compileIdentLHS failed: unknown op -", mode)
		}
	} else {
		if mode == token.ASSIGN || mode == token.DEFINE {
			ctx.out.Store(addr.(*stackVar).index)
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
}

func compileIdent(ctx *blockCtx, name string) func() {
	if sym, ok := ctx.find(name); ok {
		switch v := sym.(type) {
		case *execVar:
			ctx.infer.Push(&goValue{t: v.v.Type()})
			return func() {
				ctx.out.LoadVar(v.v)
			}
		case *stackVar:
			ctx.infer.Push(&goValue{t: v.typ})
			return func() {
				ctx.out.Load(v.index)
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
		default:
			log.Panicln("compileIdent failed: unknown -", reflect.TypeOf(sym))
		}
	} else {
		if addr, kind, ok := ctx.builtin.Find(name); ok {
			switch kind {
			case exec.SymbolVar:
			case exec.SymbolFunc, exec.SymbolFuncv:
				fn := newGoFunc(addr, kind, 0, ctx)
				ctx.infer.Push(fn)
				return func() {
					log.Panicln("compileIdent todo: goFunc")
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

func compileEllipsis(ctx *blockCtx, v *ast.Ellipsis) func() {
	if v.Elt != nil {
		log.Panicln("compileEllipsis: todo")
	}
	ctx.infer.Push(&constVal{v: int64(-1), kind: astutil.ConstUnboundInt})
	return nil
}

func compileCompositeLit(ctx *blockCtx, v *ast.CompositeLit) func() {
	if v.Type == nil {
		return compileMapLit(ctx, v)
	}
	typ := toType(ctx, v.Type)
	switch kind := typ.Kind(); kind {
	case reflect.Slice, reflect.Array:
		var typSlice, typRet reflect.Type
		if t, ok := typ.(*unboundArrayType); ok {
			n := toBoundArrayLen(ctx, v)
			typSlice = reflect.ArrayOf(n, t.elem)
		} else {
			typSlice = typ.(reflect.Type)
		}
		if typSlice.Kind() == reflect.Array {
			typRet = reflect.PtrTo(typSlice)
		} else {
			typRet = typSlice
		}
		ctx.infer.Push(&goValue{t: typRet})
		return func() {
			var nLen int
			if kind == reflect.Array {
				nLen = typSlice.Len()
			} else {
				nLen = toBoundArrayLen(ctx, v)
			}
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
			n++
			typElem := typSlice.Elem()
			for _, elt := range elts {
				if elt != nil {
					compileExpr(ctx, elt)()
					checkType(typElem, ctx.infer.Pop(), ctx.out)
				} else {
					ctx.out.Zero(typElem)
				}
			}
			ctx.out.MakeArray(typSlice, n)
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
	decl := newFuncDecl("", v.Type, v.Body, funCtx)
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
		v := boundConst(c.v, exec.TypeFromKind(c.kind))
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
	}
	xcons, xok := x.(*constVal)
	if xok { // op <const>
		ret := unaryOp(op, xcons)
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

func compileBinaryExpr(ctx *blockCtx, v *ast.BinaryExpr) func() {
	exprX := compileExpr(ctx, v.X)
	exprY := compileExpr(ctx, v.Y)
	x := ctx.infer.Get(-2)
	y := ctx.infer.Get(-1)
	op := binaryOps[v.Op]
	xcons, xok := x.(*constVal)
	ycons, yok := y.(*constVal)
	if xok && yok { // <const> op <const>
		ret := binaryOp(op, xcons, ycons)
		ctx.infer.Ret(2, ret)
		return func() {
			ret.reserve = ctx.out.Reserve()
		}
	}
	kind, ret := binaryOpResult(op, x, y)
	ctx.infer.Ret(2, ret)
	return func() {
		exprX()
		exprY()
		checkBinaryOp(kind, op, x, y, ctx.out)
		ctx.out.BuiltinOp(kind, op)
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
		if !isConstBound(kind) {
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

func compileCallExpr(ctx *blockCtx, v *ast.CallExpr) func() {
	exprFun := compileExpr(ctx, v.Fun)
	fn := ctx.infer.Pop()
	switch vfn := fn.(type) {
	case *qlFunc:
		ret := vfn.Results()
		ctx.infer.Push(ret)
		return func() {
			for _, arg := range v.Args {
				compileExpr(ctx, arg)()
			}
			arity := checkFuncCall(vfn.Proto(), 0, v, ctx)
			fun := vfn.FuncInfo()
			if fun.IsVariadic() {
				ctx.out.CallFuncv(fun, len(v.Args), arity)
			} else {
				ctx.out.CallFunc(fun, len(v.Args))
			}
		}
	case *goFunc:
		ret := vfn.Results()
		ctx.infer.Push(ret)
		return func() {
			if vfn.isMethod != 0 {
				compileExpr(ctx, v.Fun.(*ast.SelectorExpr).X)()
			}
			for _, arg := range v.Args {
				compileExpr(ctx, arg)()
			}
			nexpr := len(v.Args) + vfn.isMethod
			arity := checkFuncCall(vfn.Proto(), vfn.isMethod, v, ctx)
			switch vfn.kind {
			case exec.SymbolFunc:
				ctx.out.CallGoFunc(exec.GoFuncAddr(vfn.addr), nexpr)
			case exec.SymbolFuncv:
				ctx.out.CallGoFuncv(exec.GoFuncvAddr(vfn.addr), nexpr, arity)
			}
		}
	case *goValue:
		if vfn.t.Kind() != reflect.Func {
			log.Panicln("compileCallExpr failed: call a non function.")
		}
		ret := newFuncResults(vfn.t)
		ctx.infer.Push(ret)
		return func() {
			for _, arg := range v.Args {
				compileExpr(ctx, arg)()
			}
			exprFun()
			arity, ellipsis := checkFuncCall(vfn.t, 0, v, ctx), false
			if arity == -1 {
				arity, ellipsis = len(v.Args), true
			}
			ctx.out.CallGoClosure(len(v.Args), arity, ellipsis)
		}
	case *nonValue:
		switch nv := vfn.v.(type) {
		case goInstr:
			return nv(ctx, v)
		case reflect.Type:
			return compileTypeCast(nv, ctx, v)
		}
	}
	log.Panicln("compileCallExpr failed: unknown -", reflect.TypeOf(fn))
	return nil
}

func compileIndexExprLHS(ctx *blockCtx, v *ast.IndexExpr, mode compleMode) {
	if mode == lhsDefine {
		log.Panicln("compileIndexExprLHS: `:=` can't be used for index expression")
	}
	val := ctx.infer.Get(-1)
	compileExpr(ctx, v.X)()
	typ := ctx.infer.Get(-1).(iValue).Type()
	typElem := typ.Elem()
	if typ.Kind() == reflect.Ptr {
		if typElem.Kind() != reflect.Array {
			logPanic(ctx, v, `type %v does not support indexing`, typ)
		}
		typ = typElem
		typElem = typElem.Elem()
	}
	if cons, ok := val.(*constVal); ok {
		cons.bound(typElem, ctx.out)
	} else if t := val.(iValue).Type(); t != typElem {
		log.Panicf("compileIndexExprLHS: can't assign `%v`[i] = `%v`\n", typ, t)
	}
	exprIdx := compileExpr(ctx, v.Index)
	i := ctx.infer.Get(-1)
	ctx.infer.PopN(3)
	switch typ.Kind() {
	case reflect.Slice, reflect.Array:
		if cons, ok := i.(*constVal); ok {
			n := boundConst(cons.v, exec.TyInt)
			ctx.out.SetIndex(n.(int))
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
		ctx.out.SetIndex(-1)
	case reflect.Map:
		exprIdx()
		typIdx := typ.Key()
		if cons, ok := i.(*constVal); ok {
			cons.bound(typIdx, ctx.out)
		}
		if t := i.(iValue).Type(); t != typIdx {
			logIllTypeMapIndexPanic(ctx, v, t, typIdx)
		}
		ctx.out.SetMapIndex()
	default:
		log.Panicln("compileIndexExprLHS: unknown -", typ)
	}
}

func compileSliceExpr(ctx *blockCtx, v *ast.SliceExpr) func() { // x[i:j:k]
	var kind reflect.Kind
	exprX := compileExpr(ctx, v.X)
	x := ctx.infer.Get(-1)
	typ := x.(iValue).Type()
	if kind = typ.Kind(); kind == reflect.Ptr {
		typ = typ.Elem()
		if kind = typ.Kind(); kind != reflect.Array {
			logPanic(ctx, v, `cannot slice a (type *%v)`, typ)
		}
		typ = reflect.SliceOf(typ.Elem())
		ctx.infer.Ret(1, &goValue{typ})
	}
	return func() {
		exprX()
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
		nv := boundConst(cons.v, exec.TyInt)
		n := nv.(int)
		if n <= nlast {
			return n
		}
		ctx.out.Push(n)
		return -1
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

func compileIndexExpr(ctx *blockCtx, v *ast.IndexExpr) func() { // x[i]
	var kind reflect.Kind
	var typElem reflect.Type
	exprX := compileExpr(ctx, v.X)
	x := ctx.infer.Get(-1)
	typ := x.(iValue).Type()
	if kind = typ.Kind(); kind == reflect.Ptr {
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
	return func() {
		exprX()
		switch kind {
		case reflect.String, reflect.Slice, reflect.Array:
			n := compileIdx(ctx, v.Index, 1<<30, kind)
			ctx.out.Index(n)
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
			ctx.out.MapIndex()
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

func compileSelectorExpr(ctx *blockCtx, v *ast.SelectorExpr) func() {
	exprX := compileExpr(ctx, v.X)
	x := ctx.infer.Get(-1)
	switch vx := x.(type) {
	case *nonValue:
		switch nv := vx.v.(type) {
		case exec.GoPackage:
			if c, ok := nv.FindConst(v.Sel.Name); ok {
				ret := newConstVal(c.Value, c.Kind)
				ctx.infer.Ret(1, ret)
				return func() {
					pushConstVal(ctx.out, ret)
				}
			}
			addr, kind, ok := nv.Find(v.Sel.Name)
			if !ok {
				log.Panicln("compileSelectorExpr: not found -", nv.PkgPath(), v.Sel.Name)
			}
			switch kind {
			case exec.SymbolFunc, exec.SymbolFuncv:
				ctx.infer.Ret(1, newGoFunc(addr, kind, 0, ctx))
				return func() {
					log.Panicln("compileSelectorExpr: todo")
				}
			default:
				log.Panicln("compileSelectorExpr: unknown GoPackage symbol kind -", kind)
			}
		default:
			log.Panicln("compileSelectorExpr: unknown nonValue -", reflect.TypeOf(nv))
		}
	case *goValue:
		n, t := countPtr(vx.t)
		name := v.Sel.Name
		if sf, ok := t.FieldByName(name); ok {
			log.Panicln("compileSelectorExpr todo: structField -", t, sf)
		}
		pkgPath, method := normalizeMethod(n, t, name)
		pkg := ctx.FindGoPackage(pkgPath)
		if pkg == nil {
			log.Panicln("compileSelectorExpr failed: package not found -", pkgPath)
		}
		addr, kind, ok := pkg.Find(method)
		if !ok {
			log.Panicln("compileSelectorExpr: method not found -", method)
		}
		ctx.infer.Ret(1, newGoFunc(addr, kind, 1, ctx))
		return func() {
			log.Panicln("compileSelectorExpr: todo")
		}
	default:
		log.Panicln("compileSelectorExpr failed: unknown -", reflect.TypeOf(vx))
	}
	_ = exprX
	return nil
}

func countPtr(t reflect.Type) (int, reflect.Type) {
	n := 0
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		n++
	}
	return n, t
}

func normalizeMethod(n int, t reflect.Type, name string) (pkgPath string, formalName string) {
	typName := t.Name()
	if n > 0 {
		typName = strings.Repeat("*", n) + typName
	}
	return t.PkgPath(), "(" + typName + ")." + name
}

// -----------------------------------------------------------------------------
