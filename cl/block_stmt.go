package cl

import (
	"go/token"
	"reflect"

	"github.com/qiniu/qlang/ast"
	"github.com/qiniu/qlang/ast/astutil"
	"github.com/qiniu/qlang/exec"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

const (
	inferOnly = 1 // don't generate any code.
)

type blockCtx struct {
	*pkgCtx
	file   *fileCtx
	parent *blockCtx
}

func newBlockCtx(file *fileCtx, parent *blockCtx) *blockCtx {
	return &blockCtx{pkgCtx: file.pkg, file: file, parent: parent}
}

func (p *Package) compileBlockStmt(ctx *blockCtx, body *ast.BlockStmt) {
	for _, stmt := range body.List {
		switch v := stmt.(type) {
		case *ast.ExprStmt:
			p.compileExprStmt(ctx, v)
		default:
			log.Fatalln("compileBlockStmt failed: unknown -", reflect.TypeOf(v))
		}
	}
}

func (p *Package) compileExprStmt(ctx *blockCtx, expr *ast.ExprStmt) {
	p.compileExpr(ctx, expr.X, 0)
}

func (p *Package) compileExpr(ctx *blockCtx, expr ast.Expr, mode int) {
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

func (p *Package) compileIdent(ctx *blockCtx, name string, mode int) {
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

func (p *Package) compileBasicLit(ctx *blockCtx, v *ast.BasicLit, mode int) {
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

func (p *Package) compileBinaryExpr(ctx *blockCtx, v *ast.BinaryExpr, mode int) {
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
	kind, op := inferBinaryOp(v.Op, x, y)
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

func (p *Package) compileCallExpr(ctx *blockCtx, v *ast.CallExpr, mode int) {
	p.compileExpr(ctx, v.Fun, inferOnly)
	fn := ctx.infer.Get(-1)
	switch vfn := fn.(type) {
	case *goFunc:
		ret := &exprResult{results: vfn.TypesOut()}
		if mode == inferOnly {
			ctx.infer.Ret(1, ret)
			return
		}
		for _, arg := range v.Args {
			p.compileExpr(ctx, arg, 0)
		}
		nargs := uint32(len(v.Args))
		args := ctx.infer.GetArgs(nargs)
		tfn := vfn.Proto()
		arity := checkFuncCall(tfn, args, p.out)
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
