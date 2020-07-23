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

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/exec.spec"
	"github.com/goplus/gop/token"
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

func compileNewBlock(ctx *blockCtx, block *ast.BlockStmt) {
	ctx.out.DefineBlock()
	compileBlockStmtWith(ctx, block)
	ctx.out.EndBlock()
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
	case *ast.ForStmt:
		compileForStmt(ctx, v)
	case *ast.BlockStmt:
		compileNewBlock(ctx, v)
	case *ast.ReturnStmt:
		compileReturnStmt(ctx, v)
	case *ast.IncDecStmt:
		compileIncDecStmt(ctx, v)
	case *ast.BranchStmt:
		compileBranchStmt(ctx, v)
	case *ast.LabeledStmt:
		compileLabeledStmt(ctx, v)
	case *ast.DeferStmt:
		compileDeferStmt(ctx, v)
	case *ast.EmptyStmt:
		// do nothing
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
	noExecCtx := isNoExecCtx(parent, v.Body)
	f := ast.ForPhrase{
		For:    v.For,
		TokPos: v.TokPos,
		X:      v.X,
	}
	switch v.Tok {
	case token.DEFINE:
		f.Key = toIdent(v.Key)
		f.Value = toIdent(v.Value)
	case token.ASSIGN:
		var lhs, rhs [2]ast.Expr
		var idx int
		if v.Key != nil {
			assign := true
			if id, ok := v.Key.(*ast.Ident); ok && id.Name == "_" {
				assign = false
			}
			if assign {
				k0 := ast.NewObj(ast.Var, "_gop_k")
				f.Key = &ast.Ident{Name: k0.Name, Obj: k0}
				lhs[idx], rhs[idx] = v.Key, f.Key
				idx++
			}
		}
		if v.Value != nil {
			assign := true
			if id, ok := v.Value.(*ast.Ident); ok && id.Name == "_" {
				assign = false
			}
			if assign {
				v0 := ast.NewObj(ast.Var, "_gop_v")
				f.Value = &ast.Ident{Name: v0.Name, Obj: v0}
				lhs[idx], rhs[idx] = v.Value, f.Value
				idx++
			}
		}
		v.Body.List = append([]ast.Stmt{&ast.AssignStmt{
			Lhs: lhs[0:idx],
			Tok: token.ASSIGN,
			Rhs: rhs[0:idx],
		}}, v.Body.List...)
	}
	ctx, exprFor := compileForPhrase(parent, f, noExecCtx)
	exprFor(func() {
		compileBlockStmtWithout(ctx, v.Body)
	})
}

func toIdent(e ast.Expr) *ast.Ident {
	if e == nil {
		return nil
	}
	return e.(*ast.Ident)
}

func compileForStmt(ctx *blockCtx, v *ast.ForStmt) {
	if init := v.Init; init != nil {
		v.Init = nil
		block := &ast.BlockStmt{List: []ast.Stmt{init, v}}
		compileNewBlock(ctx, block)
		return
	}
	out := ctx.out
	start := ctx.NewLabel("")
	post := ctx.NewLabel("")
	done := ctx.NewLabel("")
	labelName := ""
	if ctx.currentLabel != nil && ctx.currentLabel.Stmt == v {
		labelName = ctx.currentLabel.Label.Name
	}
	ctx.nextFlow(post, done, labelName)
	defer func() {
		ctx.currentFlow = ctx.currentFlow.parent
	}()
	out.Label(start)
	compileExpr(ctx, v.Cond)()
	checkBool(ctx.infer.Pop())
	out.JmpIf(0, done)
	noExecCtx := isNoExecCtx(ctx, v.Body)
	ctx = newNormBlockCtxEx(ctx, noExecCtx)
	compileBlockStmtWith(ctx, v.Body)
	out.Jmp(post)
	out.Label(post)
	if v.Post != nil {
		compileStmt(ctx, v.Post)
	}
	out.Jmp(start)
	out.Label(done)
}

func compileBranchStmt(ctx *blockCtx, v *ast.BranchStmt) {
	switch v.Tok {
	case token.FALLTHROUGH:
		log.Panicln("fallthrough statement out of place")
	case token.GOTO:
		if v.Label == nil {
			log.Panicln("label not defined")
		}
		ctx.out.Jmp(ctx.requireLabel(v.Label.Name))
	case token.BREAK:
		var labelName string
		if v.Label != nil {
			labelName = v.Label.Name
		}
		label, rangeFor := ctx.getBreakLabel(labelName)
		if label != nil {
			ctx.out.Jmp(label)
			return
		}
		if rangeFor {
			ctx.out.Return(exec.BreakAsReturn)
			return
		}
		log.Panicln("break statement out of for/switch/select statements")
	case token.CONTINUE:
		var labelName string
		if v.Label != nil {
			labelName = v.Label.Name
		}
		label, rangeFor := ctx.getContinueLabel(labelName)
		if label != nil {
			ctx.out.Jmp(label)
			return
		}
		if rangeFor {
			ctx.out.Return(exec.ContinueAsReturn)
			return
		}
		log.Panicln("continue statement out of for statements")
	}
}

func compileLabeledStmt(ctx *blockCtx, v *ast.LabeledStmt) {
	label := ctx.defineLabel(v.Label.Name)
	// make sure all labels in golang code  will be used
	// TODO improvement exec/bytecode not to jump if delta==0
	ctx.out.Jmp(label)
	ctx.out.Label(label)
	ctx.currentLabel = v
	compileStmt(ctx, v.Stmt)
}

func compileDeferStmt(ctx *blockCtx, v *ast.DeferStmt) {
	var instr exec.Reserved
	out := ctx.out
	start := ctx.NewLabel("")
	end := ctx.NewLabel("")

	var f func()
	exprFun := compileExpr(ctx, v.Call.Fun)
	fn := ctx.infer.Pop()
	switch vfn := fn.(type) {
	case *qlFunc:
		ret := vfn.Results()
		ctx.infer.Push(ret)

		for _, arg := range v.Call.Args {
			compileExpr(ctx, arg)()
		}
		instr = ctx.out.Reserve()
		out.Label(start)
		f = func() {
			arity := checkFuncCall(vfn.Proto(), 0, v.Call, ctx)
			fun := vfn.FuncInfo()
			if fun.IsVariadic() {
				ctx.out.CallFuncv(fun, len(v.Call.Args), arity)
			} else {
				ctx.out.CallFunc(fun, len(v.Call.Args))
			}
		}
	case *goFunc:
		ret := vfn.Results()
		ctx.infer.Push(ret)

		for _, arg := range v.Call.Args {
			compileExpr(ctx, arg)()
		}
		instr = ctx.out.Reserve()
		out.Label(start)
		f = func() {
			if vfn.isMethod != 0 {
				compileExpr(ctx, v.Call.Fun.(*ast.SelectorExpr).X)()
			}
			nexpr := len(v.Call.Args) + vfn.isMethod
			arity := checkFuncCall(vfn.Proto(), vfn.isMethod, v.Call, ctx)
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

		for _, arg := range v.Call.Args {
			compileExpr(ctx, arg)()
		}
		instr = ctx.out.Reserve()
		out.Label(start)
		f = func() {
			exprFun()
			arity, ellipsis := checkFuncCall(vfn.t, 0, v.Call, ctx), false
			if arity == -1 {
				arity, ellipsis = len(v.Call.Args), true
			}
			ctx.out.CallGoClosure(len(v.Call.Args), arity, ellipsis)
		}
	case *nonValue:
		instr = ctx.out.Reserve()
		out.Label(start)
		// TODO compile args before compileCallExpr
		switch nv := vfn.v.(type) {
		case goInstr:
			f = nv(ctx, v.Call)
		case reflect.Type:
			f = compileTypeCast(nv, ctx, v.Call)
		}
	}

	f()
	ctx.infer.Pop()

	out.Label(end)
	instr.Set(out, out.Defer(start, end))
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
	labelName := ""
	if ctx.currentLabel != nil && ctx.currentLabel.Stmt == v {
		labelName = ctx.currentLabel.Label.Name
	}
	ctx.nextFlow(nil, done, labelName)
	defer func() {
		ctx.currentFlow = ctx.currentFlow.parent
	}()
	hasTag := v.Tag != nil
	hasCaseClause := false
	var withoutCheck exec.Label
	if hasTag {
		if len(v.Body.List) == 0 {
			return
		}
		compileExpr(ctxSw, v.Tag)()
		tag := ctx.infer.Pop()
		for idx, item := range v.Body.List {
			c, ok := item.(*ast.CaseClause)
			if !ok {
				log.Panicln("compile SwitchStmt failed: case clause expected.")
			}
			if c.List == nil { // default
				defaultBody = c.Body
				continue
			}
			if idx == len(v.Body.List)-1 {
				checkFinalFallthrough(c.Body)
			}
			hasCaseClause = true
			for _, caseExp := range c.List {
				compileExpr(ctxSw, caseExp)()
				checkCaseCompare(tag, ctx.infer.Pop(), out)
			}
			next := ctx.NewLabel("")
			out.CaseNE(next, len(c.List))
			withoutCheck = compileCaseClause(c, ctxSw, done, next, withoutCheck)
		}
		if withoutCheck != nil {
			out.Label(withoutCheck)
			withoutCheck = nil
		}
		out.Default()
	} else {
		for idx, item := range v.Body.List {
			c, ok := item.(*ast.CaseClause)
			if !ok {
				log.Panicln("compile SwitchStmt failed: case clause expected.")
			}
			if c.List == nil { // default
				defaultBody = c.Body
				continue
			}
			if idx == len(v.Body.List)-1 {
				checkFinalFallthrough(c.Body)
			}
			hasCaseClause = true
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
			withoutCheck = compileCaseClause(c, ctxSw, done, next, withoutCheck)
		}
		if withoutCheck != nil {
			out.Label(withoutCheck)
			withoutCheck = nil
		}
	}
	if defaultBody != nil {
		checkFinalFallthrough(defaultBody)
		compileBodyWith(ctxSw, defaultBody)
		if hasCaseClause {
			out.Jmp(done)
		}
	}
	if hasCaseClause {
		out.Label(done)
	}
}
func checkFinalFallthrough(body []ast.Stmt) {
	if len(body) > 0 {
		bs, ok := body[len(body)-1].(*ast.BranchStmt)
		if ok && bs.Tok == token.FALLTHROUGH {
			log.Panic("cannot fallthrough final case in switch")
		}
	}
}

func compileCaseClause(c *ast.CaseClause, ctxSw *blockCtx, done exec.Label, next exec.Label, withoutCheck exec.Label) exec.Label {
	if withoutCheck != nil {
		ctxSw.out.Label(withoutCheck)
		withoutCheck = nil
	}
	fallNext := false
	if len(c.Body) > 0 {
		bs, ok := c.Body[len(c.Body)-1].(*ast.BranchStmt)
		fallNext = ok && bs.Tok == token.FALLTHROUGH
	}
	if fallNext {
		compileBodyWith(ctxSw, c.Body[0:len(c.Body)-1])
		withoutCheck = ctxSw.NewLabel("")
		ctxSw.out.Jmp(withoutCheck)
	} else {
		compileBodyWith(ctxSw, c.Body)
		ctxSw.out.Jmp(done)
	}
	ctxSw.out.Label(next)
	return withoutCheck
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
		if ve, ok := v.Else.(*ast.BlockStmt); ok {
			compileBlockStmtWithout(ctx, ve)
		} else {
			compileStmt(ctxIf, v.Else)
		}
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
