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

	"github.com/qiniu/goplus/ast"
	"github.com/qiniu/goplus/exec.spec"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

func compileBlockStmtWith(ctx *blockCtx, body *ast.BlockStmt) {
	compileBodyWith(ctx, body.List)
}

func compileBlockStmtWithout(ctx *blockCtx, body *ast.BlockStmt) {
	for _, stmt := range body.List {
		compileStmt(ctx, stmt)
	}
}

func compileBodyWith(ctx *blockCtx, body []ast.Stmt) {
	ctxWith := newNormBlockCtx(ctx)
	for _, stmt := range body {
		compileStmt(ctxWith, stmt)
	}
}

func compileStmt(ctx *blockCtx, stmt ast.Stmt) {
	start := ctx.out.StartStmt(stmt)
	switch v := stmt.(type) {
	case *ast.ExprStmt:
		compileExprStmt(ctx, v)
	case *ast.AssignStmt:
		compileAssignStmt(ctx, v)
	case *ast.IfStmt:
		compileIfStmt(ctx, v)
	case *ast.SwitchStmt:
		compileSwitchStmt(ctx, v)
	case *ast.ForPhraseStmt:
		compileForPhraseStmt(ctx, v)
	case *ast.RangeStmt:
		compileRangeStmt(ctx, v)
	case *ast.BlockStmt:
		compileBlockStmtWith(ctx, v)
	case *ast.ReturnStmt:
		compileReturnStmt(ctx, v)
	case *ast.IncDecStmt:
		compileIncDecStmt(ctx, v)
	default:
		log.Panicln("compileStmt failed: unknown -", reflect.TypeOf(v))
	}
	ctx.out.EndStmt(stmt, start)
}

func compileForPhraseStmt(parent *blockCtx, v *ast.ForPhraseStmt) {
	noExecCtx := isNoExecCtx(parent, v.Body)
	ctx, exprFor := compileForPhrase(parent, v.ForPhrase, noExecCtx)
	exprFor(func() {
		compileBlockStmtWithout(ctx, v.Body)
	})
}

func compileRangeStmt(parent *blockCtx, v *ast.RangeStmt) {
	log.Panicln("compileRangeStmt: todo")
}

func compileSwitchStmt(ctx *blockCtx, v *ast.SwitchStmt) {
	var defaultBody []ast.Stmt
	var ctxSw *blockCtx
	if v.Init != nil {
		ctxSw = newNormBlockCtx(ctx)
		compileStmt(ctxSw, v.Init)
	} else {
		ctxSw = ctx
	}
	out := ctx.out
	done := ctx.NewLabel("")
	hasTag := v.Tag != nil
	if hasTag {
		if len(v.Body.List) == 0 {
			return
		}
		compileExpr(ctxSw, v.Tag)()
		tag := ctx.infer.Pop()
		for _, item := range v.Body.List {
			c, ok := item.(*ast.CaseClause)
			if !ok {
				log.Panicln("compile SwitchStmt failed: case clause expected.")
			}
			if c.List == nil { // default
				defaultBody = c.Body
				continue
			}
			for _, caseExp := range c.List {
				compileExpr(ctxSw, caseExp)()
				checkCaseCompare(tag, ctx.infer.Pop(), out)
			}
			next := ctx.NewLabel("")
			out.CaseNE(next, len(c.List))
			compileBodyWith(ctxSw, c.Body)
			out.Jmp(done)
			out.Label(next)
		}
		out.Default()
	} else {
		for _, item := range v.Body.List {
			c, ok := item.(*ast.CaseClause)
			if !ok {
				log.Panicln("compile SwitchStmt failed: case clause expected.")
			}
			if c.List == nil { // default
				defaultBody = c.Body
				continue
			}
			next := ctx.NewLabel("")
			last := len(c.List) - 1
			if last == 0 {
				compileExpr(ctxSw, c.List[0])()
				checkBool(ctxSw.infer.Pop())
				out.JmpIf(0, next)
			} else {
				start := ctx.NewLabel("")
				for i := 0; i < last; i++ {
					compileExpr(ctxSw, c.List[i])()
					checkBool(ctxSw.infer.Pop())
					out.JmpIf(1, start)
				}
				compileExpr(ctxSw, c.List[last])()
				checkBool(ctxSw.infer.Pop())
				out.JmpIf(0, next)
				out.Label(start)
			}
			compileBodyWith(ctxSw, c.Body)
			out.Jmp(done)
			out.Label(next)
		}
	}
	if defaultBody != nil {
		compileBodyWith(ctxSw, defaultBody)
	}
	out.Label(done)
}

func compileIfStmt(ctx *blockCtx, v *ast.IfStmt) {
	var done exec.Label
	var ctxIf *blockCtx
	if v.Init != nil {
		ctxIf = newNormBlockCtx(ctx)
		compileStmt(ctxIf, v.Init)
	} else {
		ctxIf = ctx
	}
	compileExpr(ctxIf, v.Cond)()
	checkBool(ctx.infer.Pop())
	out := ctx.out
	label := ctx.NewLabel("")
	hasElse := v.Else != nil
	out.JmpIf(0, label)
	compileBlockStmtWith(ctxIf, v.Body)
	if hasElse {
		done = ctx.NewLabel("")
		out.Jmp(done)
	}
	out.Label(label)
	if hasElse {
		compileStmt(ctxIf, v.Else)
		out.Label(done)
	}
}

func compileReturnStmt(ctx *blockCtx, expr *ast.ReturnStmt) {
	fun := ctx.fun
	if fun == nil {
		if expr.Results == nil { // return in main
			ctx.out.Return(0)
			return
		}
		log.Panicln("compileReturnStmt failed: return statement not in a function.")
	}
	rets := expr.Results
	if rets == nil {
		if fun.IsUnnamedOut() {
			log.Panicln("compileReturnStmt failed: return without values -", fun.Name())
		}
		ctx.out.Return(-1)
		return
	}
	for _, ret := range rets {
		compileExpr(ctx, ret)()
	}
	n := len(rets)
	if fun.NumOut() != n {
		log.Panicln("compileReturnStmt failed: mismatched count of return values -", fun.Name())
	}
	if ctx.infer.Len() != n {
		log.Panicln("compileReturnStmt failed: can't use multi values funcation result as return values -", fun.Name())
	}
	results := ctx.infer.GetArgs(n)
	for i, result := range results {
		v := fun.Out(i)
		checkType(v.Type(), result, ctx.out)
	}
	ctx.infer.SetLen(0)
	ctx.out.Return(int32(n))
}

func compileExprStmt(ctx *blockCtx, expr *ast.ExprStmt) {
	compileExpr(ctx, expr.X)()
	ctx.infer.PopN(1)
}

func compileIncDecStmt(ctx *blockCtx, expr *ast.IncDecStmt) {
	compileExpr(ctx, expr.X)()
	compileExprLHS(ctx, expr.X, expr.Tok)
}

func compileAssignStmt(ctx *blockCtx, expr *ast.AssignStmt) {
	if ctx.infer.Len() != 0 {
		log.Panicln("compileAssignStmt internal error: infer stack is not empty.")
	}
	if len(expr.Rhs) == 1 {
		compileExpr(ctx, expr.Rhs[0])()
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
			compileExpr(ctx, item)()
			if ctx.infer.Get(-1).(iValue).NumValues() != 1 {
				log.Panicln("compileAssignStmt failed: expr has multiple values.")
			}
		}
	}
	if ctx.infer.Len() != len(expr.Lhs) {
		log.Panicln("compileAssignStmt: assign statement has mismatched variables count -", ctx.infer.Len())
	}
	for i := len(expr.Lhs) - 1; i >= 0; i-- {
		compileExprLHS(ctx, expr.Lhs[i], expr.Tok)
	}
}

// -----------------------------------------------------------------------------
