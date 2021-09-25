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
	"fmt"
	"log"
	"path/filepath"
	"reflect"

	goast "go/ast"
	gotoken "go/token"
	"go/types"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/token"
)

func relFile(dir string, file string) string {
	if rel, err := filepath.Rel(dir, file); err == nil {
		if rel[0] == '.' {
			return rel
		}
		return "./" + rel
	}
	return file
}

func commentStmt(ctx *blockCtx, stmt ast.Stmt) {
	if ctx.fileLine {
		start := stmt.Pos()
		pos := ctx.fset.Position(start)
		if ctx.relativePath {
			pos.Filename = relFile(ctx.targetDir, pos.Filename)
		}
		line := fmt.Sprintf("\n//line %s:%d", pos.Filename, pos.Line)
		comments := &goast.CommentGroup{
			List: []*goast.Comment{{Text: line}},
		}
		ctx.cb.SetComments(comments, false)
	}
}

func compileStmts(ctx *blockCtx, body []ast.Stmt) {
	for _, stmt := range body {
		compileStmt(ctx, stmt)
	}
}

func compileStmt(ctx *blockCtx, stmt ast.Stmt) {
	if enableRecover {
		defer func() {
			if e := recover(); e != nil {
				ctx.handleRecover(e)
				ctx.cb.ResetStmt()
			}
		}()
	}
	commentStmt(ctx, stmt)
	switch v := stmt.(type) {
	case *ast.ExprStmt:
		compileExpr(ctx, v.X)
		if _, ok := v.X.(*ast.Ident); ok && isFunc(ctx.cb.InternalStack().Get(-1).Type) {
			ctx.cb.Call(0)
		}
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
		ctx.cb.Block()
		compileStmts(ctx, v.List)
		ctx.cb.End()
		return
	case *ast.SelectStmt:
		compileSelectStmt(ctx, v)
	case *ast.EmptyStmt:
		// do nothing
	default:
		log.Panicln("TODO - compileStmt failed: unknown -", reflect.TypeOf(v))
	}
	ctx.cb.EndStmt()
}

func compileReturnStmt(ctx *blockCtx, expr *ast.ReturnStmt) {
	var n = -1
	var results *types.Tuple
	for i, ret := range expr.Results {
		if c, ok := ret.(*ast.CompositeLit); ok && c.Type == nil {
			if n < 0 {
				results = ctx.cb.Func().Type().(*types.Signature).Results()
				n = results.Len()
			}
			var typ types.Type
			if i < n {
				typ = results.At(i).Type()
			}
			compileCompositeLit(ctx, c, typ, true)
		} else {
			twoValue := false
			if len(expr.Results) == 1 {
				if _, ok := ret.(*ast.ComprehensionExpr); ok {
					results = ctx.cb.Func().Type().(*types.Signature).Results()
					twoValue = (results.Len() == 2)
				}
			}
			compileExpr(ctx, ret, twoValue)
		}
	}
	ctx.cb.Return(len(expr.Results), expr)
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
		ctx.cb.DefineVarStart(expr.Pos(), names...)
		if enableRecover {
			defer func() {
				if e := recover(); e != nil {
					ctx.cb.ResetInit()
					panic(e)
				}
			}()
		}
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
		ctx.cb.AssignWith(len(expr.Lhs), len(expr.Rhs), expr)
		return
	}
	if len(expr.Lhs) != 1 || len(expr.Rhs) != 1 {
		panic("TODO: invalid syntax of assign by operator")
	}
	ctx.cb.AssignOp(gotoken.Token(tok), expr)
}

// forRange(names...) x rangeAssignThen
//    body
// end
// forRange k v x rangeAssignThen
//    body
// end
func compileRangeStmt(ctx *blockCtx, v *ast.RangeStmt) {
	cb := ctx.cb
	comments := cb.Comments()
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
	pos := v.TokPos
	if pos == 0 {
		pos = v.For
	}
	cb.RangeAssignThen(pos)
	compileStmts(ctx, v.Body.List)
	cb.SetComments(comments, true)
	setBodyHandler(ctx)
	cb.End()
}

func compileForPhraseStmt(ctx *blockCtx, v *ast.ForPhraseStmt) {
	cb := ctx.cb
	comments := cb.Comments()
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
	cb.RangeAssignThen(v.TokPos)
	if v.Cond != nil {
		cb.If()
		compileExpr(ctx, v.Cond)
		cb.Then()
		compileStmts(ctx, v.Body.List)
		cb.SetComments(comments, true)
		cb.End()
	} else {
		compileStmts(ctx, v.Body.List)
	}
	cb.SetComments(comments, true)
	setBodyHandler(ctx)
	cb.End()
}

// for init; cond then
//    body
//    post
// end
func compileForStmt(ctx *blockCtx, v *ast.ForStmt) {
	cb := ctx.cb
	comments := cb.Comments()
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
	cb.SetComments(comments, true)
	setBodyHandler(ctx)
	cb.End()
}

// if init; cond then
//    body
// end
func compileIfStmt(ctx *blockCtx, v *ast.IfStmt) {
	cb := ctx.cb
	comments := cb.Comments()
	cb.If()
	if v.Init != nil {
		compileStmt(ctx, v.Init)
	}
	compileExpr(ctx, v.Cond)
	cb.Then()
	compileStmts(ctx, v.Body.List)
	if e := v.Else; e != nil {
		cb.Else()
		if stmts, ok := e.(*ast.BlockStmt); ok {
			compileStmts(ctx, stmts.List)
		} else {
			compileStmt(ctx, e)
		}
	}
	cb.SetComments(comments, true)
	cb.End()
}

// typeSwitch(name) init; expr typeAssertThen()
// type1 type2 ... typeN typeCase(N)
//    ...
//    end
// type1 type2 ... typeM typeCase(M)
//    ...
//    end
// end
func compileTypeSwitchStmt(ctx *blockCtx, v *ast.TypeSwitchStmt) {
	var cb = ctx.cb
	comments := cb.Comments()
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
		commentStmt(ctx, stmt)
		cb.End()
	}
	cb.SetComments(comments, true)
	cb.End()
}

// switch init; tag then
// expr1 expr2 ... exprN case(N)
//    ...
//    end
// expr1 expr2 ... exprM case(M)
//    ...
//    end
// end
func compileSwitchStmt(ctx *blockCtx, v *ast.SwitchStmt) {
	cb := ctx.cb
	comments := cb.Comments()
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
		commentStmt(ctx, stmt)
		cb.End()
	}
	cb.SetComments(comments, true)
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

// select
// stmt1 commCase(1)
//    ...
//    end
// stmt2 commCase(1)
//    ...
//    end
// ...
// commCase(0)
//    ...
//    end
// end
func compileSelectStmt(ctx *blockCtx, v *ast.SelectStmt) {
	cb := ctx.cb
	comments := cb.Comments()
	cb.Select()
	for _, stmt := range v.Body.List {
		c, ok := stmt.(*ast.CommClause)
		if !ok {
			log.Panicln("TODO: compile SelectStmt failed - comm clause expected.")
		}
		var n int
		if c.Comm != nil {
			compileStmt(ctx, c.Comm)
			n = 1
		}
		cb.CommCase(n) // CommCase(0) means default case
		compileStmts(ctx, c.Body)
		commentStmt(ctx, stmt)
		cb.End()
	}
	cb.SetComments(comments, true)
	cb.End()
}

func compileBranchStmt(ctx *blockCtx, v *ast.BranchStmt) {
	label := v.Label
	switch v.Tok {
	case token.GOTO:
		if label == nil {
			log.Panicln("TODO: label not defined")
		}
		ctx.cb.Goto(label.Name, label)
	case token.BREAK:
		var name string
		if label != nil {
			name = label.Name
		}
		ctx.cb.Break(name, label)
	case token.CONTINUE:
		var name string
		if label != nil {
			name = label.Name
		}
		ctx.cb.Continue(name, label)
	case token.FALLTHROUGH:
		panic("TODO: fallthrough statement out of place")
	default:
		panic("TODO: compileBranchStmt - unknown")
	}
}

func compileLabeledStmt(ctx *blockCtx, v *ast.LabeledStmt) {
	ctx.cb.Label(v.Label.Name, v.Label)
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
				compileType(ctx, spec.(*ast.TypeSpec))
			}
		case token.CONST:
			cdecl := ctx.pkg.NewConstDecl(ctx.cb.Scope())
			loadConstSpecs(ctx, cdecl, d.Specs)
		case token.VAR:
			for _, spec := range d.Specs {
				v := spec.(*ast.ValueSpec)
				loadVars(ctx, v, false)
			}
		default:
			log.Panicln("TODO: compileDeclStmt - unknown")
		}
	}
}

func compileType(ctx *blockCtx, t *ast.TypeSpec) {
	name := t.Name.Name
	if t.Assign != token.NoPos { // alias type
		ctx.cb.AliasType(name, toType(ctx, t.Type))
	} else {
		ctx.cb.NewType(name).InitType(ctx.pkg, toType(ctx, t.Type))
	}
}

// -----------------------------------------------------------------------------
