package cl

import (
	"reflect"
	"strings"
	"syscall"

	"github.com/qiniu/qlang/ast"
	"github.com/qiniu/qlang/ast/astutil"
	"github.com/qiniu/qlang/exec"
	"github.com/qiniu/qlang/token"
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

func (p *Package) compileBlockStmt(ctx *blockCtx, body *ast.BlockStmt) {
	for _, stmt := range body.List {
		switch v := stmt.(type) {
		case *ast.ExprStmt:
			p.compileExprStmt(ctx, v)
		case *ast.AssignStmt:
			p.compileAssignStmt(ctx, v)
		default:
			log.Panicln("compileBlockStmt failed: unknown -", reflect.TypeOf(v))
		}
	}
}

func (p *Package) compileExprStmt(ctx *blockCtx, expr *ast.ExprStmt) {
	p.compileExpr(ctx, expr.X, 0)
	ctx.infer.PopN(1)
}

func (p *Package) compileAssignStmt(ctx *blockCtx, expr *ast.AssignStmt) {
	if ctx.infer.Len() != 0 {
		log.Panicln("compileAssignStmt internal error: infer stack is not empty.")
	}
	if len(expr.Rhs) == 1 {
		p.compileExpr(ctx, expr.Rhs[0], 0)
		v := ctx.infer.Get(-1).(iValue)
		n := v.NumValues()
		if n != 1 {
			if n == 0 {
				log.Panicln("compileAssignStmt failed: expr has no return value.")
			}
			rhs := make([]interface{}, n)
			for i := 0; i < n; i++ {
				rhs[i] = v.Value(i)
			}
			ctx.infer.Ret(1, rhs...)
		}
	} else {
		for _, item := range expr.Rhs {
			p.compileExpr(ctx, item, 0)
			if ctx.infer.Get(-1).(iValue).NumValues() != 1 {
				log.Panicln("compileAssignStmt failed: expr has multiple values.")
			}
		}
	}
	if ctx.infer.Len() != len(expr.Lhs) {
		log.Panicln("compileAssignStmt failed: assign statment has mismatched variables count.")
	}
	for i := len(expr.Lhs) - 1; i >= 0; i-- {
		p.compileExpr(ctx, expr.Lhs[i], expr.Tok)
	}
}

func (p *Package) compileExpr(ctx *blockCtx, expr ast.Expr, mode compleMode) {
	switch v := expr.(type) {
	case *ast.Ident:
		p.compileIdent(ctx, v.Name, mode)
	case *ast.BasicLit:
		p.compileBasicLit(ctx, v, mode)
	case *ast.CallExpr:
		p.compileCallExpr(ctx, v, mode)
	case *ast.BinaryExpr:
		p.compileBinaryExpr(ctx, v, mode)
	case *ast.SelectorExpr:
		p.compileSelectorExpr(ctx, v, mode)
	default:
		log.Panicln("compileExpr failed: unknown -", reflect.TypeOf(v))
	}
}

func (p *Package) compileIdent(ctx *blockCtx, name string, mode compleMode) {
	if mode > lhsBase {
		in := ctx.infer.Get(-1)
		addr, err := ctx.findVar(name)
		if err == nil {
			if mode == lhsDefine && addr.NestDepth != ctx.out.NestDepth {
				log.Warn("requireVar: variable is shadowed -", name)
			}
		} else if mode == lhsAssign || err != syscall.ENOENT {
			log.Panicln("compileIdent failed:", err)
		} else {
			typ := boundType(in.(iValue))
			addr = ctx.insertVar(name, typ)
		}
		checkType(addr.Type, in, ctx.out)
		ctx.infer.PopN(1)
		ctx.out.StoreVar(addr)
	} else if sym, ok := ctx.find(name); ok {
		switch v := sym.(type) {
		case *exec.Var:
			ctx.infer.Push(&goValue{t: v.Type})
			if mode == inferOnly {
				return
			}
			ctx.out.LoadVar(v)
		case string: // pkgPath
			pkg := exec.FindGoPackage(v)
			if pkg == nil {
				log.Panicln("compileIdent failed: package not found -", v)
			}
			ctx.infer.Push(&nonValue{pkg})
		case *funcDecl:
			ctx.infer.Push(newQlFunc(v))
			if mode == inferOnly {
				return
			}
			log.Panicln("compileIdent failed: todo - funcDecl")
		default:
			log.Panicln("compileIdent failed: unknown -", reflect.TypeOf(sym))
		}
	} else {
		addr, kind, ok := ctx.builtin.Find(name)
		if !ok {
			log.Panicln("compileIdent failed: unknown -", name)
		}
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
}

func (p *Package) compileBasicLit(ctx *blockCtx, v *ast.BasicLit, mode compleMode) {
	if mode > lhsBase {
		log.Panicln("compileBasicLit: can't be lhs (left hand side) expr.")
	}
	kind, n := astutil.ToConst(v)
	ret := &constVal{v: n, kind: kind, reserve: -1}
	ctx.infer.Push(ret)
	if mode == inferOnly {
		return
	}
	if astutil.IsConstBound(kind) {
		if kind == astutil.ConstBoundRune {
			n = rune(n.(int64))
		}
		ctx.out.Push(n)
	} else {
		ret.reserve = ctx.out.Reserve()
	}
}

func (p *Package) compileBinaryExpr(ctx *blockCtx, v *ast.BinaryExpr, mode compleMode) {
	if mode > lhsBase {
		log.Panicln("compileBinaryExpr: can't be lhs (left hand side) expr.")
	}
	p.compileExpr(ctx, v.X, inferOnly)
	p.compileExpr(ctx, v.Y, inferOnly)
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
	p.compileExpr(ctx, v.X, 0)
	p.compileExpr(ctx, v.Y, 0)
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
	if !astutil.IsConstBound(kind) {
		kind = vy.Kind()
		if !astutil.IsConstBound(kind) {
			log.Panicln("binaryOp: expect x, y aren't const values either.")
		}
	}
	i := op.GetInfo()
	if i.Out != exec.SameAsFirst {
		kind = i.Out
	}
	return kind, &goValue{t: exec.TypeFromKind(kind)}
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

func (p *Package) compileCallExpr(ctx *blockCtx, v *ast.CallExpr, mode compleMode) {
	if mode > lhsBase {
		log.Panicln("compileCallExpr: can't be lhs (left hand side) expr.")
	}
	p.compileExpr(ctx, v.Fun, inferOnly)
	fn := ctx.infer.Get(-1)
	switch vfn := fn.(type) {
	case *qlFunc:
		ret := vfn.Results()
		if mode == inferOnly {
			ctx.infer.Ret(1, ret)
			return
		}
		for _, arg := range v.Args {
			p.compileExpr(ctx, arg, 0)
		}
		nargs := uint32(len(v.Args))
		args := ctx.infer.GetArgs(nargs)
		out := ctx.out
		arity := checkFuncCall(vfn.Proto(), 0, args, out)
		fun := vfn.getFuncInfo()
		if fun.IsVariadic() {
			out.CallFuncv(fun, arity)
		} else {
			out.CallFunc(fun)
		}
		ctx.infer.Ret(uint32(len(v.Args)+1), ret)
	case *goFunc:
		ret := vfn.Results()
		if mode == inferOnly {
			ctx.infer.Ret(1, ret)
			return
		}
		if vfn.isMethod != 0 {
			p.compileExpr(ctx, v.Fun.(*ast.SelectorExpr).X, 0)
		}
		for _, arg := range v.Args {
			p.compileExpr(ctx, arg, 0)
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
	}
	log.Panicln("compileCallExpr failed: unknown -", reflect.TypeOf(fn))
}

func (p *Package) compileSelectorExpr(ctx *blockCtx, v *ast.SelectorExpr, mode compleMode) {
	p.compileExpr(ctx, v.X, inferOnly)
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
