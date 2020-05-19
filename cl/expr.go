package cl

import (
	"reflect"
	"strings"
	"syscall"

	"github.com/qiniu/qlang/v6/ast"
	"github.com/qiniu/qlang/v6/ast/astutil"
	"github.com/qiniu/qlang/v6/exec"
	"github.com/qiniu/qlang/v6/token"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

type compleMode = token.Token

const (
	inferOnly compleMode = 1 // don't generate any code.
	lhsBase   compleMode = 10
	lhsAssign compleMode = token.ASSIGN // leftHandSide = ...
	lhsDefine compleMode = token.DEFINE // leftHandSide := ...
)

// -----------------------------------------------------------------------------

func compileExpr(ctx *blockCtx, expr ast.Expr, mode compleMode) {
	switch v := expr.(type) {
	case *ast.Ident:
		compileIdent(ctx, v.Name, mode)
	case *ast.BasicLit:
		compileBasicLit(ctx, v, mode)
	case *ast.CallExpr:
		compileCallExpr(ctx, v, mode)
	case *ast.BinaryExpr:
		compileBinaryExpr(ctx, v, mode)
	case *ast.SelectorExpr:
		compileSelectorExpr(ctx, v, mode)
	case *ast.CompositeLit:
		compileCompositeLit(ctx, v, mode)
	case *ast.SliceLit:
		compileSliceLit(ctx, v, mode)
	case *ast.FuncLit:
		compileFuncLit(ctx, v, mode)
	case *ast.Ellipsis:
		compileEllipsis(ctx, v, mode)
	case *ast.KeyValueExpr:
		panic("compileExpr: ast.KeyValueExpr unexpected")
	default:
		log.Panicln("compileExpr failed: unknown -", reflect.TypeOf(v))
	}
}

func compileEllipsis(ctx *blockCtx, v *ast.Ellipsis, mode compleMode) {
	if mode != inferOnly {
		log.Panicln("compileEllipsis: only support inferOnly mode.")
	}
	if v.Elt != nil {
		log.Panicln("compileEllipsis: todo")
	}
	ctx.infer.Push(&constVal{v: int64(-1), kind: astutil.ConstUnboundInt})
}

func compileIdent(ctx *blockCtx, name string, mode compleMode) {
	if mode > lhsBase {
		in := ctx.infer.Get(-1)
		addr, err := ctx.findVar(name)
		if err == nil {
			if mode == lhsDefine && !addr.inCurrentCtx(ctx) {
				log.Warn("requireVar: variable is shadowed -", name)
			}
		} else if mode == lhsAssign || err != syscall.ENOENT {
			log.Panicln("compileIdent failed:", err, "-", name)
		} else {
			typ := boundType(in.(iValue))
			addr = ctx.insertVar(name, typ)
		}
		checkType(addr.getType(), in, ctx.out)
		ctx.infer.PopN(1)
		if v, ok := addr.(*execVar); ok {
			ctx.out.StoreVar((*exec.Var)(v))
		} else {
			ctx.out.Store(addr.(*stackVar).index)
		}
	} else if sym, ok := ctx.find(name); ok {
		switch v := sym.(type) {
		case *execVar:
			ctx.infer.Push(&goValue{t: v.Type})
			if mode == inferOnly {
				return
			}
			ctx.out.LoadVar((*exec.Var)(v))
		case *stackVar:
			ctx.infer.Push(&goValue{t: v.typ})
			if mode == inferOnly {
				return
			}
			ctx.out.Load(v.index)
		case string: // pkgPath
			pkg := exec.FindGoPackage(v)
			if pkg == nil {
				log.Panicln("compileIdent failed: package not found -", v)
			}
			ctx.infer.Push(&nonValue{pkg})
		case *funcDecl:
			ctx.use(v)
			ctx.infer.Push(newQlFunc(v))
			if mode == inferOnly {
				return
			}
			log.Panicln("compileIdent failed: todo - funcDecl")
		default:
			log.Panicln("compileIdent failed: unknown -", reflect.TypeOf(sym))
		}
	} else {
		if addr, kind, ok := ctx.builtin.Find(name); ok {
			switch kind {
			case exec.SymbolVar:
			case exec.SymbolFunc, exec.SymbolFuncv:
				ctx.infer.Push(newGoFunc(addr, kind, 0))
				if mode == inferOnly {
					return
				}
			}
			log.Panicln("compileIdent failed: unknown -", kind, addr)
		}
		if ci, ok := ctx.builtin.FindConst(name); ok {
			compileConst(ctx, ci.Kind, ci.Value, mode)
			return
		}
		log.Panicln("compileIdent failed: unknown -", name)
	}
}

func compileCompositeLit(ctx *blockCtx, v *ast.CompositeLit, mode compleMode) {
	if mode > lhsBase {
		log.Panicln("compileCompositeLit: can't be lhs (left hand side) expr.")
	}
	if v.Type == nil {
		compileMapLit(ctx, v, mode)
		return
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
		if mode == inferOnly {
			ctx.infer.Push(&goValue{t: typSlice})
			return
		}
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
				compileExpr(ctx, elt, 0)
				checkType(typElem, ctx.infer.Pop(), ctx.out)
			} else {
				ctx.out.Zero(typElem)
			}
		}
		ctx.out.MakeArray(typSlice, n)
		ctx.infer.Push(&goValue{t: typSlice})
	case reflect.Map:
		typMap := typ.(reflect.Type)
		if mode == inferOnly {
			ctx.infer.Push(&goValue{t: typMap})
			return
		}
		typKey := typMap.Key()
		typVal := typMap.Elem()
		for _, elt := range v.Elts {
			switch e := elt.(type) {
			case *ast.KeyValueExpr:
				compileExpr(ctx, e.Key, 0)
				checkType(typKey, ctx.infer.Pop(), ctx.out)
				compileExpr(ctx, e.Value, 0)
				checkType(typVal, ctx.infer.Pop(), ctx.out)
			default:
				log.Panicln("compileCompositeLit: map requires key-value expr.")
			}
		}
		ctx.out.MakeMap(typMap, len(v.Elts))
		ctx.infer.Push(&goValue{t: typMap})
	default:
		log.Panicln("compileCompositeLit failed: unknown -", reflect.TypeOf(typ))
	}
}

func compileSliceLit(ctx *blockCtx, v *ast.SliceLit, mode compleMode) {
	for _, elt := range v.Elts {
		compileExpr(ctx, elt, mode)
	}
	n := len(v.Elts)
	if n == 0 {
		log.Debug("compileSliceLit:", exec.TyEmptyInterfaceSlice)
		ctx.infer.Push(&goValue{t: exec.TyEmptyInterfaceSlice})
		if mode != inferOnly {
			ctx.out.MakeArray(exec.TyEmptyInterfaceSlice, 0)
		}
		return
	}
	elts := ctx.infer.GetArgs(uint32(n))
	typElem := boundElementType(elts, 0, n, 1)
	if typElem == nil {
		typElem = exec.TyEmptyInterface
	}
	typSlice := reflect.SliceOf(typElem)
	if mode == inferOnly {
		ctx.infer.Ret(uint32(n), &goValue{t: typSlice})
		return
	}
	out := ctx.out
	log.Debug("compileSliceLit:", typSlice)
	checkElementType(typElem, elts, 0, n, 1, out)
	out.MakeArray(typSlice, len(v.Elts))
	ctx.infer.Ret(uint32(n), &goValue{t: typSlice})
}

func compileMapLit(ctx *blockCtx, v *ast.CompositeLit, mode compleMode) {
	for _, elt := range v.Elts {
		switch e := elt.(type) {
		case *ast.KeyValueExpr:
			compileExpr(ctx, e.Key, mode)
			compileExpr(ctx, e.Value, mode)
		default:
			log.Panicln("compileMapLit: map requires key-value expr.")
		}
	}
	n := len(v.Elts) << 1
	if n == 0 {
		typMap := reflect.MapOf(exec.TyString, exec.TyEmptyInterface)
		log.Debug("compileMapLit:", typMap)
		ctx.infer.Push(&goValue{t: typMap})
		if mode != inferOnly {
			ctx.out.MakeMap(typMap, 0)
		}
		return
	}
	elts := ctx.infer.GetArgs(uint32(n))
	typKey := boundElementType(elts, 0, n, 2)
	if typKey == nil {
		log.Panicln("compileMapLit: mismatched key type.")
	}
	typVal := boundElementType(elts, 1, n, 2)
	if typVal == nil {
		typVal = exec.TyEmptyInterface
	}
	typMap := reflect.MapOf(typKey, typVal)
	if mode == inferOnly {
		ctx.infer.Ret(uint32(n), &goValue{t: typMap})
		return
	}
	out := ctx.out
	log.Debug("compileMapLit:", typMap)
	checkElementType(typKey, elts, 0, n, 2, out)
	checkElementType(typVal, elts, 1, n, 2, out)
	out.MakeMap(typMap, len(v.Elts))
	ctx.infer.Ret(uint32(n), &goValue{t: typMap})
}

func compileFuncLit(ctx *blockCtx, v *ast.FuncLit, mode compleMode) {
	if mode > lhsBase {
		log.Panicln("compileFuncLit: can't be lhs (left hand side) expr.")
	}
	funCtx := newBlockCtx(ctx, false)
	decl := newFuncDecl("", v.Type, v.Body, funCtx)
	ctx.use(decl)
	ctx.infer.Push(newQlFunc(decl))
	if mode == inferOnly {
		return
	}
	ctx.out.GoClosure(decl.fi)
}

func compileBasicLit(ctx *blockCtx, v *ast.BasicLit, mode compleMode) {
	if mode > lhsBase {
		log.Panicln("compileBasicLit: can't be lhs (left hand side) expr.")
	}
	kind, n := astutil.ToConst(v)
	compileConst(ctx, kind, n, mode)
}

func compileConst(ctx *blockCtx, kind astutil.ConstKind, n interface{}, mode compleMode) {
	ret := newConstVal(n, kind)
	ctx.infer.Push(ret)
	if mode == inferOnly {
		return
	}
	if isConstBound(kind) {
		if kind == astutil.ConstBoundRune {
			n = rune(n.(int64))
		}
		ctx.out.Push(n)
	} else {
		ret.reserve = ctx.out.Reserve()
	}
}

func compileBinaryExpr(ctx *blockCtx, v *ast.BinaryExpr, mode compleMode) {
	if mode > lhsBase {
		log.Panicln("compileBinaryExpr: can't be lhs (left hand side) expr.")
	}
	compileExpr(ctx, v.X, inferOnly)
	compileExpr(ctx, v.Y, inferOnly)
	x := ctx.infer.Get(-2)
	y := ctx.infer.Get(-1)
	op := binaryOps[v.Op]
	xcons, xok := x.(*constVal)
	ycons, yok := y.(*constVal)
	if xok && yok { // <const> op <const>
		ret := binaryOp(op, xcons, ycons)
		ctx.infer.Ret(2, ret)
		if mode != inferOnly {
			ret.reserve = ctx.out.Reserve()
		}
		return
	}
	kind, ret := binaryOpResult(op, x, y)
	if mode == inferOnly {
		ctx.infer.Ret(2, ret)
		return
	}
	compileExpr(ctx, v.X, 0)
	compileExpr(ctx, v.Y, 0)
	x = ctx.infer.Get(-2)
	y = ctx.infer.Get(-1)
	checkBinaryOp(kind, op, x, y, ctx.out)
	ctx.out.BuiltinOp(kind, op)
	ctx.infer.Ret(4, ret)
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
	token.QUO:     exec.OpDiv,
	token.REM:     exec.OpMod,
	token.AND:     exec.OpBitAnd,
	token.OR:      exec.OpBitOr,
	token.XOR:     exec.OpBitXor,
	token.AND_NOT: exec.OpBitAndNot,
	token.SHL:     exec.OpBitSHL,
	token.SHR:     exec.OpBitSHR,
	token.LSS:     exec.OpLT,
	token.LEQ:     exec.OpLE,
	token.GTR:     exec.OpGT,
	token.GEQ:     exec.OpGE,
	token.EQL:     exec.OpEQ,
	token.NEQ:     exec.OpNE,
	token.LAND:    exec.OpLAnd,
	token.LOR:     exec.OpLOr,
}

func compileCallExpr(ctx *blockCtx, v *ast.CallExpr, mode compleMode) {
	if mode > lhsBase {
		log.Panicln("compileCallExpr: can't be lhs (left hand side) expr.")
	}
	compileExpr(ctx, v.Fun, inferOnly)
	fn := ctx.infer.Get(-1)
	switch vfn := fn.(type) {
	case *qlFunc:
		ret := vfn.Results()
		if mode == inferOnly {
			ctx.infer.Ret(1, ret)
			return
		}
		for _, arg := range v.Args {
			compileExpr(ctx, arg, 0)
		}
		out := ctx.out
		nargs := uint32(len(v.Args))
		args := ctx.infer.GetArgs(nargs)
		arity := checkFuncCall(vfn.Proto(), 0, args, out)
		fun := vfn.FuncInfo()
		if fun.IsVariadic() {
			out.CallFuncv(fun, arity)
		} else {
			out.CallFunc(fun)
		}
		ctx.infer.Ret(uint32(len(v.Args)+1), ret)
		return
	case *goFunc:
		ret := vfn.Results()
		if mode == inferOnly {
			ctx.infer.Ret(1, ret)
			return
		}
		if vfn.isMethod != 0 {
			compileExpr(ctx, v.Fun.(*ast.SelectorExpr).X, 0)
		}
		for _, arg := range v.Args {
			compileExpr(ctx, arg, 0)
		}
		nargs := uint32(len(v.Args))
		args := ctx.infer.GetArgs(nargs)
		out := ctx.out
		arity := checkFuncCall(vfn.Proto(), vfn.isMethod, args, out)
		switch vfn.kind {
		case exec.SymbolFunc:
			out.CallGoFunc(exec.GoFuncAddr(vfn.addr))
		case exec.SymbolFuncv:
			out.CallGoFuncv(exec.GoFuncvAddr(vfn.addr), arity)
		}
		ctx.infer.Ret(uint32(len(v.Args)+1+vfn.isMethod), ret)
		return
	case *goValue:
		if vfn.t.Kind() != reflect.Func {
			log.Panicln("compileCallExpr failed: call a non function.")
		}
		ret := newFuncResults(vfn.t)
		if mode == inferOnly {
			ctx.infer.Ret(1, ret)
			return
		}
		for _, arg := range v.Args {
			compileExpr(ctx, arg, 0)
		}
		compileExpr(ctx, v.Fun, 0)
		nargs := uint32(len(v.Args))
		args := ctx.infer.GetArgs(nargs)
		arity := checkFuncCall(vfn.t, 0, args, ctx.out)
		ctx.out.CallGoClosure(arity)
		ctx.infer.Ret(uint32(len(v.Args)+2), ret)
		return
	}
	log.Panicln("compileCallExpr failed: unknown -", reflect.TypeOf(fn))
}

func compileSelectorExpr(ctx *blockCtx, v *ast.SelectorExpr, mode compleMode) {
	compileExpr(ctx, v.X, inferOnly)
	x := ctx.infer.Get(-1)
	switch vx := x.(type) {
	case *nonValue:
		switch nv := vx.v.(type) {
		case *exec.GoPackage:
			addr, kind, ok := nv.Find(v.Sel.Name)
			if !ok {
				log.Panicln("compileSelectorExpr: not found -", nv.PkgPath, v.Sel.Name)
			}
			switch kind {
			case exec.SymbolFunc, exec.SymbolFuncv:
				ctx.infer.Ret(1, newGoFunc(addr, kind, 0))
				if mode == inferOnly {
					return
				}
				log.Panicln("compileSelectorExpr: todo")
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
		pkg := exec.FindGoPackage(pkgPath)
		if pkg == nil {
			log.Panicln("compileSelectorExpr failed: package not found -", pkgPath)
		}
		addr, kind, ok := pkg.Find(method)
		if !ok {
			log.Panicln("compileSelectorExpr: method not found -", method)
		}
		ctx.infer.Ret(1, newGoFunc(addr, kind, 1))
		if mode == inferOnly {
			return
		}
		log.Panicln("compileSelectorExpr: todo")
	default:
		log.Panicln("compileSelectorExpr failed: unknown -", reflect.TypeOf(vx))
	}
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
