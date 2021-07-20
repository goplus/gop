/*
 Copyright 2021 The GoPlus Authors (goplus.org)

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
	"log"
	"reflect"

	gotoken "go/token"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/token"
)

func compileStmts(ctx *blockCtx, body []ast.Stmt) {
	for _, stmt := range body {
		compileStmt(ctx, stmt)
	}
}

func compileStmt(ctx *blockCtx, stmt ast.Stmt) {
	switch v := stmt.(type) {
	case *ast.ExprStmt:
		compileExpr(ctx, v.X)
	case *ast.AssignStmt:
		compileAssignStmt(ctx, v)
	case *ast.ReturnStmt:
		compileReturnStmt(ctx, v)
	case *ast.IfStmt:
		compileIfStmt(ctx, v)
	case *ast.SwitchStmt:
		compileSwitchStmt(ctx, v)
	case *ast.RangeStmt:
		compileRangeStmt(ctx, v)
	case *ast.ForStmt:
		compileForStmt(ctx, v)
	case *ast.ForPhraseStmt:
		compileForPhraseStmt(ctx, v)
	case *ast.IncDecStmt:
		compileIncDecStmt(ctx, v)
	case *ast.DeferStmt:
		compileDeferStmt(ctx, v)
	case *ast.GoStmt:
		compileGoStmt(ctx, v)
	case *ast.DeclStmt:
		compileDeclStmt(ctx, v)
	case *ast.TypeSwitchStmt:
		compileTypeSwitchStmt(ctx, v)
	case *ast.SendStmt:
		compileSendStmt(ctx, v)
	case *ast.BranchStmt:
		compileBranchStmt(ctx, v)
	case *ast.LabeledStmt:
		compileLabeledStmt(ctx, v)
	case *ast.BlockStmt:
		compileStmts(ctx, v.List)
	case *ast.EmptyStmt:
		// do nothing
	default:
		log.Panicln("TODO - compileStmt failed: unknown -", reflect.TypeOf(v))
	}
	ctx.cb.EndStmt()
}

func compileReturnStmt(ctx *blockCtx, expr *ast.ReturnStmt) {
	for _, ret := range expr.Results {
		compileExpr(ctx, ret)
	}
	ctx.cb.Return(len(expr.Results))
}

func compileIncDecStmt(ctx *blockCtx, expr *ast.IncDecStmt) {
	compileExprLHS(ctx, expr.X)
	ctx.cb.IncDec(gotoken.Token(expr.Tok))
}

func compileSendStmt(ctx *blockCtx, expr *ast.SendStmt) {
	compileExpr(ctx, expr.Chan)
	compileExpr(ctx, expr.Value)
	ctx.cb.Send()
}

func compileAssignStmt(ctx *blockCtx, expr *ast.AssignStmt) {
	tok := expr.Tok
	twoValue := (len(expr.Lhs) == 2 && len(expr.Rhs) == 1)
	if tok == token.DEFINE {
		names := make([]string, len(expr.Lhs))
		for i, lhs := range expr.Lhs {
			if v, ok := lhs.(*ast.Ident); ok {
				names[i] = v.Name
			} else {
				log.Panicln("TODO: non-name $v on left side of :=")
			}
		}
		ctx.cb.DefineVarStart(names...)
		for _, rhs := range expr.Rhs {
			compileExpr(ctx, rhs, twoValue)
		}
		ctx.cb.EndInit(len(expr.Rhs))
		return
	}
	for _, lhs := range expr.Lhs {
		compileExprLHS(ctx, lhs)
	}
	for _, rhs := range expr.Rhs {
		compileExpr(ctx, rhs, twoValue)
	}
	if tok == token.ASSIGN {
		ctx.cb.Assign(len(expr.Lhs), len(expr.Rhs))
		return
	}
	if len(expr.Lhs) != 1 || len(expr.Rhs) != 1 {
		panic("TODO: invalid syntax of assign by operator")
	}
	ctx.cb.AssignOp(gotoken.Token(tok))
}

// forRange names... := range x {...}
// forRange k, v = range x {...}
func compileRangeStmt(ctx *blockCtx, v *ast.RangeStmt) {
	cb := ctx.cb
	if v.Tok == token.DEFINE {
		names := make([]string, 1, 2)
		if v.Key == nil {
			names[0] = "_"
		} else {
			names[0] = v.Key.(*ast.Ident).Name
		}
		if v.Value != nil {
			names = append(names, v.Value.(*ast.Ident).Name)
		}
		cb.ForRange(names...)
		compileExpr(ctx, v.X)
	} else {
		cb.ForRange()
		n := 0
		if v.Key == nil {
			if v.Value != nil {
				ctx.cb.VarRef(nil) // underscore
				n++
			}
		} else {
			compileExprLHS(ctx, v.Key)
			n++
		}
		if v.Value != nil {
			compileExprLHS(ctx, v.Value)
			n++
		}
		compileExpr(ctx, v.X)
	}
	cb.RangeAssignThen()
	compileStmts(ctx, v.Body.List)
	cb.End()
}

func compileForPhraseStmt(ctx *blockCtx, v *ast.ForPhraseStmt) {
	cb := ctx.cb
	names := make([]string, 1, 2)
	if v.Key == nil {
		names[0] = "_"
	} else {
		names[0] = v.Key.Name
	}
	if v.Value != nil {
		names = append(names, v.Value.Name)
	}
	cb.ForRange(names...)
	compileExpr(ctx, v.X)
	cb.RangeAssignThen()
	if v.Cond != nil {
		cb.If()
		compileExpr(ctx, v.Cond)
		cb.Then()
		compileStmts(ctx, v.Body.List)
		cb.End()
	} else {
		compileStmts(ctx, v.Body.List)
	}
	cb.End()
}

// for init; cond then
//    body
//    post
// end
func compileForStmt(ctx *blockCtx, v *ast.ForStmt) {
	cb := ctx.cb
	cb.For()
	if v.Init != nil {
		compileStmt(ctx, v.Init)
	}
	if v.Cond != nil {
		compileExpr(ctx, v.Cond)
	} else {
		cb.None()
	}
	cb.Then()
	compileStmts(ctx, v.Body.List)
	if v.Post != nil {
		cb.Post()
		compileStmt(ctx, v.Post)
	}
	cb.End()
}

func compileIfStmt(ctx *blockCtx, v *ast.IfStmt) {
	cb := ctx.cb
	cb.If()
	if v.Init != nil {
		compileStmt(ctx, v.Init)
	}
	compileExpr(ctx, v.Cond)
	cb.Then()
	compileStmts(ctx, v.Body.List)
	if v.Else != nil {
		cb.Else()
		compileStmt(ctx, v.Else)
	}
	cb.End()
}

// typeSwitch(name) init; expr typeAssertThen()
// type1, type2, ... typeN typeCase(N)
//    ...
//    end
// type1, type2, ... typeM typeCase(M)
//    ...
//    end
// end
func compileTypeSwitchStmt(ctx *blockCtx, v *ast.TypeSwitchStmt) {
	var cb = ctx.cb
	var name string
	var ta *ast.TypeAssertExpr
	switch stmt := v.Assign.(type) {
	case *ast.AssignStmt:
		if stmt.Tok != token.DEFINE || len(stmt.Lhs) != 1 || len(stmt.Rhs) != 1 {
			panic("TODO: type switch syntax error")
		}
		name = stmt.Lhs[0].(*ast.Ident).Name
		ta = stmt.Rhs[0].(*ast.TypeAssertExpr)
	case *ast.ExprStmt:
		ta = stmt.X.(*ast.TypeAssertExpr)
	}
	if ta.Type != nil {
		panic("TODO: type switch syntax error, please use x.(type)")
	}
	cb.TypeSwitch(name)
	if v.Init != nil {
		compileStmt(ctx, v.Init)
	}
	compileExpr(ctx, ta.X)
	cb.TypeAssertThen()
	for _, stmt := range v.Body.List {
		c, ok := stmt.(*ast.CaseClause)
		if !ok {
			log.Panicln("TODO: compile TypeSwitchStmt failed - case clause expected.")
		}
		for _, citem := range c.List {
			compileExpr(ctx, citem)
		}
		cb.TypeCase(len(c.List)) // TypeCase(0) means default case
		compileStmts(ctx, c.Body)
		cb.End()
	}
	cb.End()
}

func compileSwitchStmt(ctx *blockCtx, v *ast.SwitchStmt) {
	cb := ctx.cb
	cb.Switch()
	if v.Init != nil {
		compileStmt(ctx, v.Init)
	}
	if v.Tag != nil { // switch tag {....}
		compileExpr(ctx, v.Tag)
	} else {
		cb.None() // switch {...}
	}
	cb.Then()
	for _, stmt := range v.Body.List {
		c, ok := stmt.(*ast.CaseClause)
		if !ok {
			log.Panicln("TODO: compile SwitchStmt failed - case clause expected.")
		}
		for _, citem := range c.List {
			compileExpr(ctx, citem)
		}
		cb.Case(len(c.List)) // Case(0) means default case
		body, has := hasFallthrough(c.Body)
		compileStmts(ctx, body)
		if has {
			cb.Fallthrough()
		}
		cb.End()
	}
	cb.End()
}

func hasFallthrough(body []ast.Stmt) ([]ast.Stmt, bool) {
	if n := len(body); n > 0 {
		if bs, ok := body[n-1].(*ast.BranchStmt); ok && bs.Tok == token.FALLTHROUGH {
			return body[:n-1], true
		}
	}
	return body, false
}

func compileBranchStmt(ctx *blockCtx, v *ast.BranchStmt) {
	switch v.Tok {
	case token.GOTO:
		if v.Label == nil {
			log.Panicln("TODO: label not defined")
		}
		ctx.cb.Goto(v.Label.Name)
	case token.BREAK:
		var name string
		if v.Label != nil {
			name = v.Label.Name
		}
		ctx.cb.Break(name)
	case token.CONTINUE:
		var name string
		if v.Label != nil {
			name = v.Label.Name
		}
		ctx.cb.Continue(name)
	case token.FALLTHROUGH:
		panic("TODO: fallthrough statement out of place")
	default:
		panic("TODO: compileBranchStmt - unknown")
	}
}

func compileLabeledStmt(ctx *blockCtx, v *ast.LabeledStmt) {
	ctx.cb.Label(v.Label.Name)
	compileStmt(ctx, v.Stmt)
}

func compileGoStmt(ctx *blockCtx, v *ast.GoStmt) {
	compileCallExpr(ctx, v.Call)
	ctx.cb.Go()
}

func compileDeferStmt(ctx *blockCtx, v *ast.DeferStmt) {
	compileCallExpr(ctx, v.Call)
	ctx.cb.Defer()
}

func compileDeclStmt(ctx *blockCtx, expr *ast.DeclStmt) {
	switch d := expr.Decl.(type) {
	case *ast.GenDecl:
		switch d.Tok {
		case token.TYPE:
			for _, spec := range d.Specs {
				loadType(ctx, spec.(*ast.TypeSpec))
			}
		case token.CONST:
			for _, spec := range d.Specs {
				v := spec.(*ast.ValueSpec)
				names := makeNames(v.Names)
				loadConsts(ctx, names, v)
			}
		case token.VAR:
			for _, spec := range d.Specs {
				v := spec.(*ast.ValueSpec)
				names := makeNames(v.Names)
				loadVars(ctx, names, v)
			}
		default:
			log.Panicln("TODO: compileDeclStmt - unknown")
		}
	}
}

// -----------------------------------------------------------------------------
