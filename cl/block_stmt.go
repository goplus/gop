package cl

import (
	"reflect"

	"github.com/qiniu/qlang/ast"
	"github.com/qiniu/qlang/ast/astutil"
	"github.com/qiniu/qlang/exec"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

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
	p.compileExpr(ctx, expr.X)
}

func (p *Package) compileExpr(ctx *blockCtx, expr ast.Expr) {
	switch v := expr.(type) {
	case *ast.Ident:
		p.compileIdent(ctx, v.Name)
	case *ast.BasicLit:
		p.compileBasicLit(ctx, v)
	case *ast.CallExpr:
		p.compileCallExpr(ctx, v)
	default:
		log.Fatalln("compileExpr failed: unknown -", reflect.TypeOf(v))
	}
}

func (p *Package) compileIdent(ctx *blockCtx, name string) {
	addr, kind, ok := ctx.builtin.Find(name)
	if !ok {
		log.Fatalln("compileIdent failed: unknown -", name)
	}
	switch kind {
	case exec.SymbolVar:
	case exec.SymbolFunc, exec.SymbolVariadicFunc:
		ctx.infer.Push(getGoFunc(addr, kind))
		return
	}
	log.Fatalln("compileIdent failed: unknown -", kind, addr)
}

func (p *Package) compileBasicLit(ctx *blockCtx, v *ast.BasicLit) {
	kind, n := astutil.ToConst(v)
	ctx.infer.Push(&constVal{v: n, kind: kind})
}

func (p *Package) compileCallExpr(ctx *blockCtx, v *ast.CallExpr) {
	p.compileExpr(ctx, v.Fun)
	for _, arg := range v.Args {
		p.compileExpr(ctx, arg)
	}
}

// -----------------------------------------------------------------------------
