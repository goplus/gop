package cl

import (
	"reflect"

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

type varDecl struct {
	Type reflect.Type
	Idx  uint32
}

type blockCtx struct {
	*pkgCtx
	file   *fileCtx
	parent *blockCtx
	vars   map[string]*varDecl
	rhs    exec.Stack // only used by assign statement
}

func newBlockCtx(file *fileCtx, parent *blockCtx) *blockCtx {
	return &blockCtx{
		pkgCtx: file.pkg,
		file:   file,
		parent: parent,
		vars:   make(map[string]*varDecl),
	}
}

func (p *blockCtx) findVar(name string) (addr exec.Address, typ reflect.Type, ok bool) {
	var scope uint32
	for p != nil {
		if v, ok := p.vars[name]; ok {
			return exec.MakeAddr(scope, v.Idx), v.Type, true
		}
		p = p.parent
		scope++
	}
	return
}

func (p *blockCtx) requireVar(name string, typ iType) (addr exec.Address, ok bool) {
	addr, typExist, ok := p.findVar(name)
	if ok {
		if addr.Scope() > 0 {
			log.Fatalln("requireVar failed: variable is shadowed -", name)
		}
		if !checkType(typExist, typ) {
			log.Fatalln("requireVar failed: mismatched type -", name)
		}
		return
	}
	idx := uint32(len(p.vars))
	p.vars[name] = &varDecl{Type: boundType(typ), Idx: idx}
	return exec.MakeAddr(0, idx), true
}

// -----------------------------------------------------------------------------

func (p *Package) compileBlockStmt(ctx *blockCtx, body *ast.BlockStmt) {
	for _, stmt := range body.List {
		switch v := stmt.(type) {
		case *ast.ExprStmt:
			p.compileExprStmt(ctx, v)
		case *ast.AssignStmt:
			p.compileAssignStmt(ctx, v)
		default:
			log.Fatalln("compileBlockStmt failed: unknown -", reflect.TypeOf(v))
		}
	}
}

func (p *Package) compileExprStmt(ctx *blockCtx, expr *ast.ExprStmt) {
	p.compileExpr(ctx, expr.X, 0)
}

func (p *Package) compileAssignStmt(ctx *blockCtx, expr *ast.AssignStmt) {
	/*
		ctx.rhs.Init()
		for _, rhs := range expr.Rhs {
			p.compileExpr(ctx, rhs, 0)
			t := ctx.infer.Get(-1).(iValue).TypeOf()
			n := t.NumValues()
			for i := 0; i < n; i++ {
				ctx.rhs.Push(t.TypeValue(i))
			}
		}
		if len(typs) != len(expr.Lhs) {
			log.Fatalln("compileAssignStmt failed: assign statment has mismatched variables count.")
		}
		ctx.rhs = typs
		for i := len(typs) - 1; i >= 0; i-- {
			p.compileExpr(ctx, expr.Lhs[i], expr.Tok)
		}
	*/
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
	default:
		log.Fatalln("compileExpr failed: unknown -", reflect.TypeOf(v))
	}
}

func (p *Package) compileIdent(ctx *blockCtx, name string, mode compleMode) {
	if mode > lhsBase {
		/*		var addr exec.Address
				var typ reflect.Type
				var ok bool
				if mode == lhsAssign {
					if addr, type, ok = ctx.findVar(name); !ok {
						log.Fatalln("compileIdent failed: symbol not found -", name)
					}
				} else {
					ctx.requireVar(name, val.(iValue).TypeOf())
				}
		*/
		log.Fatalln("compileIdent: can't be lhs (left hand side) expr.")
	} else {
		addr, kind, ok := ctx.builtin.Find(name)
		if !ok {
			log.Fatalln("compileIdent failed: unknown -", name)
		}
		switch kind {
		case exec.SymbolVar:
		case exec.SymbolFunc, exec.SymbolVariadicFunc:
			ctx.infer.Push(newGoFunc(addr, kind))
			if mode == inferOnly {
				return
			}
		}
		log.Fatalln("compileIdent failed: unknown -", kind, addr)
	}
}

func (p *Package) compileBasicLit(ctx *blockCtx, v *ast.BasicLit, mode compleMode) {
	if mode > lhsBase {
		log.Fatalln("compileBasicLit: can't be lhs (left hand side) expr.")
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
		p.out.Push(n)
	} else {
		ret.reserve = p.out.Reserve()
	}
}

func (p *Package) compileBinaryExpr(ctx *blockCtx, v *ast.BinaryExpr, mode compleMode) {
	if mode > lhsBase {
		log.Fatalln("compileBinaryExpr: can't be lhs (left hand side) expr.")
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
			ret.reserve = p.out.Reserve()
		}
		return
	}
	log.Println("compileBinaryExpr: not impl")
	/*	kind, op := inferBinaryOp(v.Op, x, y)
		ret := &operatorResult{kind: kind, op: op}
		if mode == inferOnly {
			ctx.infer.Ret(2, ret)
			return
		}
		p.compileExpr(ctx, v.X, 0)
		p.compileExpr(ctx, v.Y, 0)
		x = ctx.infer.Get(-2)
		y = ctx.infer.Get(-1)
		checkBinaryOp(kind, op, x, y, p.out)
		p.out.BuiltinOp(kind, op)
		ctx.infer.Ret(4, ret)
	*/
}

/*
func inferBinaryOp(tok token.Token, x, y interface{}) (kind exec.Kind, op exec.Operator) {
	tx := x.(iValue).TypeOf()
	ty := y.(iValue).TypeOf()
	op = binaryOps[tok]
	if tx.NumValues() != 1 || ty.NumValues() != 1 {
		log.Fatalln("inferBinaryOp failed: argument isn't an expr.")
	}
	kind = tx.TypeValue(0).Kind()
	if !astutil.IsConstBound(kind) {
		kind = ty.TypeValue(0).Kind()
		if !astutil.IsConstBound(kind) {
			log.Fatalln("inferBinaryOp failed: expect x, y aren't const values either.")
		}
	}
	return
}
*/

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
		log.Fatalln("compileCallExpr: can't be lhs (left hand side) expr.")
	}
	p.compileExpr(ctx, v.Fun, inferOnly)
	fn := ctx.infer.Get(-1)
	switch vfn := fn.(type) {
	case *goFunc:
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
		arity := checkFuncCall(vfn.Proto(), args, p.out)
		switch vfn.kind {
		case exec.SymbolFunc:
			p.out.CallGoFun(exec.GoFuncAddr(vfn.addr))
		case exec.SymbolVariadicFunc:
			p.out.CallGoFunv(exec.GoVariadicFuncAddr(vfn.addr), arity)
		}
		ctx.infer.Ret(uint32(len(v.Args)+1), ret)
		return
	}
	log.Fatalln("compileCallExpr failed: unknown -", reflect.TypeOf(fn))
}

// -----------------------------------------------------------------------------
