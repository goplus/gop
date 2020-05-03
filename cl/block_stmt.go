package cl

import (
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
	ctx.infer.Push(&constVal{v: n, kind: kind})
	if mode == inferOnly {
		return
	}
	if !astutil.IsConstBound(kind) {
		log.Fatalln("compileBasicLit: todo - how to support unbound const?")
	} else {
		p.out.Push(n)
	}
}

func (p *Package) compileCallExpr(ctx *blockCtx, v *ast.CallExpr, mode int) {
	p.compileExpr(ctx, v.Fun, inferOnly)
	fn := ctx.infer.Get(-1)
	switch vfn := fn.(type) {
	case *goFunc:
		ret := &exprResult{results: vfn.TypesOut(), expr: v}
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
		arity := checkFuncCall(tfn, args)
		switch vfn.kind {
		case exec.SymbolFunc:
			p.out.CallGoFun(exec.GoFuncAddr(vfn.addr))
		case exec.SymbolVariadicFunc:
			p.out.CallGoFunv(exec.GoVariadicFuncAddr(vfn.addr), arity)
		}
		ctx.infer.Ret(1, ret)
		return
	}
	log.Fatalln("compileCallExpr failed: unknown -", reflect.TypeOf(fn))
}

// -----------------------------------------------------------------------------
