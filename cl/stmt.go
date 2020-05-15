package cl

import (
	"reflect"

	"github.com/qiniu/qlang/v6/ast"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

func compileBlockStmt(ctx *blockCtx, body *ast.BlockStmt) {
	for _, stmt := range body.List {
		switch v := stmt.(type) {
		case *ast.ExprStmt:
			compileExprStmt(ctx, v)
		case *ast.AssignStmt:
			compileAssignStmt(ctx, v)
		case *ast.ReturnStmt:
			compileReturnStmt(ctx, v)
		default:
			log.Panicln("compileBlockStmt failed: unknown -", reflect.TypeOf(v))
		}
	}
}

func compileReturnStmt(ctx *blockCtx, expr *ast.ReturnStmt) {
	fun := ctx.fun
	if fun == nil {
		log.Panicln("compileReturnStmt failed: return statement not in a function.")
	}
	rets := expr.Results
	if rets == nil {
		if fun.IsUnnamedOut() {
			log.Panicln("compileReturnStmt failed: return without values -", fun.Name)
		}
		ctx.out.Return(-1)
		return
	}
	for _, ret := range rets {
		compileExpr(ctx, ret, 0)
	}
	n := len(rets)
	if fun.NumOut() != n {
		log.Panicln("compileReturnStmt failed: mismatched count of return values -", fun.Name)
	}
	if ctx.infer.Len() != n {
		log.Panicln("compileReturnStmt failed: can't use multi values funcation result as return values -", fun.Name)
	}
	results := ctx.infer.GetArgs(uint32(n))
	for i, result := range results {
		v := fun.Out(i)
		checkType(v.Type, result, ctx.out)
	}
	ctx.infer.SetLen(0)
	ctx.out.Return(int32(n))
}

func compileExprStmt(ctx *blockCtx, expr *ast.ExprStmt) {
	compileExpr(ctx, expr.X, 0)
	ctx.infer.PopN(1)
}

func compileAssignStmt(ctx *blockCtx, expr *ast.AssignStmt) {
	if ctx.infer.Len() != 0 {
		log.Panicln("compileAssignStmt internal error: infer stack is not empty.")
	}
	if len(expr.Rhs) == 1 {
		compileExpr(ctx, expr.Rhs[0], 0)
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
			compileExpr(ctx, item, 0)
			if ctx.infer.Get(-1).(iValue).NumValues() != 1 {
				log.Panicln("compileAssignStmt failed: expr has multiple values.")
			}
		}
	}
	if ctx.infer.Len() != len(expr.Lhs) {
		log.Panicln("compileAssignStmt: assign statment has mismatched variables count -", ctx.infer.Len())
	}
	for i := len(expr.Lhs) - 1; i >= 0; i-- {
		compileExpr(ctx, expr.Lhs[i], expr.Tok)
	}
}

// -----------------------------------------------------------------------------
