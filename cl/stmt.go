/*
 * Copyright (c) 2021 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cl

import (
	"fmt"
	"go/constant"
	"log"
	"path/filepath"
	"reflect"

	goast "go/ast"
	gotoken "go/token"
	"go/types"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/token"
	"github.com/goplus/gox"
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

func offsetFileLine(ctx *blockCtx, file string) string {
	if !filepath.IsAbs(file) {
		file = filepath.Join(ctx.fileLineWorkingDir, file)
	}
	if ret, err := filepath.Rel(ctx.fileLineRoot, file); err == nil {
		return filepath.ToSlash(ret)
	}
	return file
}

func commentStmt(ctx *blockCtx, stmt ast.Stmt) {
	if ctx.fileLine {
		start := stmt.Pos()
		pos := ctx.fset.Position(start)
		if ctx.fileLineRoot != "" {
			pos.Filename = offsetFileLine(ctx, pos.Filename)
		} else if ctx.relativePath {
			pos.Filename = relFile(ctx.targetDir, pos.Filename)
		}
		line := fmt.Sprintf("\n//line %s:%d:1", pos.Filename, pos.Line)
		comments := &goast.CommentGroup{
			List: []*goast.Comment{{Text: line}},
		}
		ctx.cb.SetComments(comments, false)
	}
}

func commentFunc(ctx *blockCtx, fn *gox.Func, decl *ast.FuncDecl) {
	start := decl.Name.Pos()
	if ctx.fileLine && start != token.NoPos {
		pos := ctx.fset.Position(start)
		if ctx.fileLineRoot != "" {
			pos.Filename = offsetFileLine(ctx, pos.Filename)
		} else if ctx.relativePath {
			pos.Filename = relFile(ctx.targetDir, pos.Filename)
		}
		var line string
		if decl.Shadow {
			line = fmt.Sprintf("//line %s:%d", pos.Filename, pos.Line)
		} else {
			line = fmt.Sprintf("//line %s:%d:1", pos.Filename, pos.Line)
		}
		doc := &goast.CommentGroup{}
		if decl.Doc != nil {
			doc.List = append(doc.List, decl.Doc.List...)
			doc.List = append(doc.List, &goast.Comment{Text: "//"})
		}
		doc.List = append(doc.List, &goast.Comment{Text: line})
		fn.SetComments(ctx.pkg, doc)
	} else if decl.Doc != nil {
		fn.SetComments(ctx.pkg, decl.Doc)
	}
}

func compileStmts(ctx *blockCtx, body []ast.Stmt) {
	for _, stmt := range body {
		if v, ok := stmt.(*ast.LabeledStmt); ok {
			expr := v.Label
			ctx.cb.NewLabel(expr.Pos(), expr.Name)
		}
	}
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
		inFlags := 0
		if isCommandWithoutArgs(v.X) {
			inFlags = clCommandWithoutArgs
		}
		compileExpr(ctx, v.X, inFlags)
		if inFlags != 0 && gox.IsFunc(ctx.cb.InternalStack().Get(-1).Type) {
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

func isCommandWithoutArgs(x ast.Expr) bool {
retry:
	switch t := x.(type) {
	case *ast.Ident:
		return true
	case *ast.SelectorExpr:
		x = t.X
		goto retry
	}
	return false
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
			inFlags := 0
			if len(expr.Results) == 1 {
				if _, ok := ret.(*ast.ComprehensionExpr); ok {
					results = ctx.cb.Func().Type().(*types.Signature).Results()
					if results.Len() == 2 {
						inFlags = clCallWithTwoValue
					}
				}
			}
			switch v := ret.(type) {
			case *ast.LambdaExpr, *ast.LambdaExpr2:
				rtyp := ctx.cb.Func().Type().(*types.Signature).Results().At(i).Type()
				sig, ok := rtyp.(*types.Signature)
				if !ok {
					panic(ctx.newCodeErrorf(
						ret.Pos(), "cannot use lambda expression as type %v in return statement", rtyp))
				}
				compileLambda(ctx, v, sig)
			default:
				compileExpr(ctx, ret, inFlags)
			}
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
	inFlags := 0
	if len(expr.Lhs) == 2 && len(expr.Rhs) == 1 {
		inFlags = clCallWithTwoValue
	}
	if tok == token.DEFINE {
		names := make([]string, len(expr.Lhs))
		for i, lhs := range expr.Lhs {
			if v, ok := lhs.(*ast.Ident); ok {
				names[i] = v.Name
			} else {
				compileExprLHS(ctx, lhs) // only for typesutil.Check
				log.Panicln("TODO: non-name $v on left side of :=")
			}
		}
		if rec := ctx.recorder(); rec != nil {
			newNames := make([]*ast.Ident, 0, len(names))
			scope := ctx.cb.Scope()
			for _, lhs := range expr.Lhs {
				v := lhs.(*ast.Ident)
				if scope.Lookup(v.Name) == nil {
					newNames = append(newNames, v)
				}
			}
			defer defNames(ctx, newNames, scope)
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
			compileExpr(ctx, rhs, inFlags)
		}
		ctx.cb.EndInit(len(expr.Rhs))
		return
	}
	for _, lhs := range expr.Lhs {
		compileExprLHS(ctx, lhs)
	}
	for _, rhs := range expr.Rhs {
		switch e := rhs.(type) {
		case *ast.LambdaExpr, *ast.LambdaExpr2:
			if len(expr.Lhs) == 1 && len(expr.Rhs) == 1 {
				typ := ctx.cb.Get(-1).Type.(interface{ Elem() types.Type }).Elem()
				sig := checkLambdaFuncType(ctx, e, typ, clLambaAssign, expr.Lhs[0])
				compileLambda(ctx, e, sig)
			} else {
				panic(ctx.newCodeErrorf(e.Pos(), "lambda unsupport multiple assignment"))
			}
		default:
			compileExpr(ctx, rhs, inFlags)
		}
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
//
//	body
//
// end
// forRange k v x rangeAssignThen
//
//	body
//
// end
func compileRangeStmt(ctx *blockCtx, v *ast.RangeStmt) {
	if re, ok := v.X.(*ast.RangeExpr); ok {
		tok := token.DEFINE
		if v.Tok == token.ASSIGN {
			tok = v.Tok
		}
		compileForStmt(ctx, toForStmt(v.For, v.Key, v.Body, re, tok, nil))
		return
	}
	cb := ctx.cb
	comments, once := cb.BackupComments()
	var defineNames []*ast.Ident
	if v.Tok == token.DEFINE {
		names := make([]string, 1, 2)
		defineNames = make([]*ast.Ident, 1, 2)
		if v.Key == nil {
			names[0] = "_"
		} else {
			key := v.Key.(*ast.Ident)
			names[0] = key.Name
			defineNames[0] = key
		}
		if v.Value != nil {
			value := v.Value.(*ast.Ident)
			names = append(names, value.Name)
			defineNames = append(defineNames, value)
		}
		cb.ForRangeEx(names, v)
		compileExpr(ctx, v.X)
	} else {
		cb.ForRangeEx(nil, v)
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
	cb.RangeAssignThen(pos) // TODO: need NewScope for body
	if len(defineNames) > 0 {
		defNames(ctx, defineNames, cb.Scope())
	}
	compileStmts(ctx, v.Body.List)
	cb.SetComments(comments, once)
	setBodyHandler(ctx)
	cb.End(v)
}

func compileForPhraseStmt(ctx *blockCtx, v *ast.ForPhraseStmt) {
	if re, ok := v.X.(*ast.RangeExpr); ok {
		compileForStmt(ctx, toForStmt(v.For, v.Value, v.Body, re, token.DEFINE, v.ForPhrase))
		return
	}
	cb := ctx.cb
	comments, once := cb.BackupComments()
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
		cb.SetComments(comments, once)
		cb.End()
	} else {
		compileStmts(ctx, v.Body.List)
		cb.SetComments(comments, once)
	}
	setBodyHandler(ctx)
	cb.End()
}

func toForStmt(forPos token.Pos, value ast.Expr, body *ast.BlockStmt, re *ast.RangeExpr, tok token.Token, fp *ast.ForPhrase) *ast.ForStmt {
	nilIdent := value == nil
	if !nilIdent {
		if v, ok := value.(*ast.Ident); ok && v.Name == "_" {
			nilIdent = true
		}
	}
	if nilIdent {
		value = &ast.Ident{NamePos: forPos, Name: "_gop_k"}
	}
	first := re.First
	if first == nil {
		first = &ast.BasicLit{ValuePos: forPos, Kind: token.INT, Value: "0"}
	}
	initLhs := []ast.Expr{value}
	initRhs := []ast.Expr{first}
	replaceValue := false
	var cond ast.Expr
	var post ast.Expr
	switch re.Last.(type) {
	case *ast.Ident, *ast.BasicLit:
		cond = re.Last
	default:
		replaceValue = true
		cond = &ast.Ident{NamePos: forPos, Name: "_gop_end"}
		initLhs = append(initLhs, cond)
		initRhs = append(initRhs, re.Last)
	}
	if re.Expr3 == nil {
		post = &ast.BasicLit{ValuePos: forPos, Kind: token.INT, Value: "1"}
	} else {
		switch re.Expr3.(type) {
		case *ast.Ident, *ast.BasicLit:
			post = re.Expr3
		default:
			replaceValue = true
			post = &ast.Ident{NamePos: forPos, Name: "_gop_step"}
			initLhs = append(initLhs, post)
			initRhs = append(initRhs, re.Expr3)
		}
	}
	if tok == token.ASSIGN && replaceValue {
		oldValue := value
		value = &ast.Ident{NamePos: forPos, Name: "_gop_k"}
		initLhs[0] = value
		body.List = append([]ast.Stmt{&ast.AssignStmt{
			Lhs:    []ast.Expr{oldValue},
			TokPos: forPos,
			Tok:    token.ASSIGN,
			Rhs:    []ast.Expr{value},
		}}, body.List...)
		tok = token.DEFINE
	}
	if fp != nil && fp.Cond != nil {
		condStmt := &ast.IfStmt{
			If:   fp.IfPos,
			Init: fp.Init,
			Cond: fp.Cond,
			Body: body,
		}
		body = &ast.BlockStmt{
			List: []ast.Stmt{condStmt},
		}
	}
	return &ast.ForStmt{
		For: forPos,
		Init: &ast.AssignStmt{
			Lhs:    initLhs,
			TokPos: re.To,
			Tok:    tok,
			Rhs:    initRhs,
		},
		Cond: &ast.BinaryExpr{
			X:     value,
			OpPos: re.To,
			Op:    token.LSS,
			Y:     cond,
		},
		Post: &ast.AssignStmt{
			Lhs:    []ast.Expr{value},
			TokPos: re.Colon2,
			Tok:    token.ADD_ASSIGN,
			Rhs:    []ast.Expr{post},
		},
		Body: body,
	}
}

// for init; cond then
//
//	body
//	post
//
// end
func compileForStmt(ctx *blockCtx, v *ast.ForStmt) {
	cb := ctx.cb
	comments, once := cb.BackupComments()
	cb.For(v)
	if v.Init != nil {
		compileStmt(ctx, v.Init)
	}
	if v.Cond != nil {
		compileExpr(ctx, v.Cond)
	} else {
		cb.None()
	}
	cb.Then(v.Body)
	compileStmts(ctx, v.Body.List)
	if v.Post != nil {
		cb.Post()
		compileStmt(ctx, v.Post)
	}
	cb.SetComments(comments, once)
	setBodyHandler(ctx)
	cb.End(v)
}

// if init; cond then
//
//	body
//
// end
func compileIfStmt(ctx *blockCtx, v *ast.IfStmt) {
	cb := ctx.cb
	comments, once := cb.BackupComments()
	cb.If(v)
	if v.Init != nil {
		compileStmt(ctx, v.Init)
	}
	compileExpr(ctx, v.Cond)
	cb.Then(v.Body)
	compileStmts(ctx, v.Body.List)
	if e := v.Else; e != nil {
		cb.Else(e)
		if stmts, ok := e.(*ast.BlockStmt); ok {
			compileStmts(ctx, stmts.List)
		} else {
			compileStmt(ctx, e)
		}
	}
	cb.SetComments(comments, once)
	cb.End(v)
}

// typeSwitch(name) init; expr typeAssertThen()
// type1 type2 ... typeN typeCase(N)
//
//	...
//	end
//
// type1 type2 ... typeM typeCase(M)
//
//	...
//	end
//
// end
func compileTypeSwitchStmt(ctx *blockCtx, v *ast.TypeSwitchStmt) {
	var cb = ctx.cb
	comments, once := cb.BackupComments()
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
	cb.TypeSwitch(name, v)
	if v.Init != nil {
		compileStmt(ctx, v.Init)
	}
	compileExpr(ctx, ta.X)
	cb.TypeAssertThen()
	seen := make(map[types.Type]ast.Expr)
	var firstDefault ast.Stmt
	for _, stmt := range v.Body.List {
		c, ok := stmt.(*ast.CaseClause)
		if !ok {
			log.Panicln("TODO: compile TypeSwitchStmt failed - case clause expected.")
		}
		for _, citem := range c.List {
			compileExpr(ctx, citem)
			T := cb.Get(-1).Type
			if tt, ok := T.(*gox.TypeType); ok {
				T = tt.Type()
			}
			var haserr bool
			for t, other := range seen {
				if T == nil && t == nil || T != nil && t != nil && types.Identical(T, t) {
					haserr = true
					pos := citem.Pos()
					if T == types.Typ[types.UntypedNil] {
						ctx.handleErrorf(
							pos, "multiple nil cases in type switch (first at %v)",
							ctx.Position(other.Pos()))
					} else {
						ctx.handleErrorf(
							pos, "duplicate case %s in type switch\n\tprevious case at %v",
							T, ctx.Position(other.Pos()))
					}
				}
			}
			if !haserr {
				seen[T] = citem
			}
		}
		if c.List == nil {
			if firstDefault != nil {
				ctx.handleErrorf(
					c.Pos(), "multiple defaults in type switch (first at %v)", ctx.Position(firstDefault.Pos()))
			} else {
				firstDefault = c
			}
		}
		cb.TypeCase(len(c.List), c) // TypeCase(0) means default case
		compileStmts(ctx, c.Body)
		commentStmt(ctx, stmt)
		cb.End(c)
	}
	cb.SetComments(comments, once)
	cb.End(v)
}

// switch init; tag then
// expr1 expr2 ... exprN case(N)
//
//	...
//	end
//
// expr1 expr2 ... exprM case(M)
//
//	...
//	end
//
// end
func compileSwitchStmt(ctx *blockCtx, v *ast.SwitchStmt) {
	cb := ctx.cb
	comments, once := cb.BackupComments()
	cb.Switch(v)
	if v.Init != nil {
		compileStmt(ctx, v.Init)
	}
	if v.Tag != nil { // switch tag {....}
		compileExpr(ctx, v.Tag)
	} else {
		cb.None() // switch {...}
	}
	cb.Then(v.Body)
	seen := make(valueMap)
	var firstDefault ast.Stmt
	for _, stmt := range v.Body.List {
		c, ok := stmt.(*ast.CaseClause)
		if !ok {
			log.Panicln("TODO: compile SwitchStmt failed - case clause expected.")
		}
		for _, citem := range c.List {
			compileExpr(ctx, citem)
			v := cb.Get(-1)
			if val := goVal(v.CVal); val != nil {
				// look for duplicate types for a given value
				// (quadratic algorithm, but these lists tend to be very short)
				typ := types.Default(v.Type)
				var haserr bool
				for _, vt := range seen[val] {
					if types.Identical(typ, vt.typ) {
						haserr = true
						src := ctx.LoadExpr(v.Src)
						if lit, ok := v.Src.(*ast.BasicLit); ok {
							ctx.handleErrorf(lit.Pos(), "duplicate case %s in switch\n\tprevious case at %v",
								src, ctx.Position(vt.pos))
						} else {
							ctx.handleErrorf(v.Src.Pos(), "duplicate case %s (value %#v) in switch\n\tprevious case at %v",
								src, val, ctx.Position(vt.pos))
						}
					}
				}
				if !haserr {
					seen[val] = append(seen[val], valueType{v.Src.Pos(), typ})
				}
			}
		}
		if c.List == nil {
			if firstDefault != nil {
				ctx.handleErrorf(c.Pos(), "multiple defaults in switch (first at %v)", ctx.Position(firstDefault.Pos()))
			} else {
				firstDefault = c
			}
		}
		cb.Case(len(c.List), c) // Case(0) means default case
		body, has := hasFallthrough(c.Body)
		compileStmts(ctx, body)
		if has {
			cb.Fallthrough()
		}
		commentStmt(ctx, stmt)
		cb.End(c)
	}
	cb.SetComments(comments, once)
	cb.End(v)
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
//
//	...
//	end
//
// stmt2 commCase(1)
//
//	...
//	end
//
// ...
// commCase(0)
//
//	...
//	end
//
// end
func compileSelectStmt(ctx *blockCtx, v *ast.SelectStmt) {
	cb := ctx.cb
	comments, once := cb.BackupComments()
	cb.Select(v)
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
		cb.CommCase(n, c) // CommCase(0) means default case
		compileStmts(ctx, c.Body)
		commentStmt(ctx, stmt)
		cb.End(c)
	}
	cb.SetComments(comments, once)
	cb.End(v)
}

func compileBranchStmt(ctx *blockCtx, v *ast.BranchStmt) {
	label := v.Label
	switch v.Tok {
	case token.GOTO:
		cb := ctx.cb
		if l, ok := cb.LookupLabel(label.Name); ok {
			cb.Goto(l)
			return
		}
		compileCallExpr(ctx, &ast.CallExpr{
			Fun:  &ast.Ident{NamePos: v.TokPos, Name: "goto", Obj: &ast.Object{Data: label}},
			Args: []ast.Expr{label},
		}, clIdentGoto)
	case token.BREAK:
		ctx.cb.Break(getLabel(ctx, label))
	case token.CONTINUE:
		ctx.cb.Continue(getLabel(ctx, label))
	case token.FALLTHROUGH:
		ctx.handleErrorf(v.Pos(), "fallthrough statement out of place")
	default:
		panic("unknown branch statement")
	}
}

func getLabel(ctx *blockCtx, label *ast.Ident) *gox.Label {
	if label != nil {
		if l, ok := ctx.cb.LookupLabel(label.Name); ok {
			return l
		}
		ctx.handleErrorf(label.Pos(), "label %v is not defined", label.Name)
	}
	return nil
}

func compileLabeledStmt(ctx *blockCtx, v *ast.LabeledStmt) {
	l, _ := ctx.cb.LookupLabel(v.Label.Name)
	ctx.cb.Label(l)
	compileStmt(ctx, v.Stmt)
}

func compileGoStmt(ctx *blockCtx, v *ast.GoStmt) {
	compileCallExpr(ctx, v.Call, 0)
	ctx.cb.Go()
}

func compileDeferStmt(ctx *blockCtx, v *ast.DeferStmt) {
	compileCallExpr(ctx, v.Call, 0)
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
			cdecl := ctx.pkg.NewConstDefs(ctx.cb.Scope())
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
	if name == "_" {
		return
	}
	if t.Assign != token.NoPos { // alias type
		ctx.cb.AliasType(name, toType(ctx, t.Type))
	} else {
		ctx.cb.NewType(name).InitType(ctx.pkg, toType(ctx, t.Type))
	}
}

type (
	valueMap  map[interface{}][]valueType // underlying Go value -> valueType
	valueType struct {
		pos token.Pos
		typ types.Type
	}
)

// goVal returns the Go value for val, or nil.
func goVal(val constant.Value) interface{} {
	// val should exist, but be conservative and check
	if val == nil {
		return nil
	}
	// Match implementation restriction of other compilers.
	// gc only checks duplicates for integer, floating-point
	// and string values, so only create Go values for these
	// types.
	switch val.Kind() {
	case constant.Int:
		if x, ok := constant.Int64Val(val); ok {
			return x
		}
		if x, ok := constant.Uint64Val(val); ok {
			return x
		}
	case constant.Float:
		if x, ok := constant.Float64Val(val); ok {
			return x
		}
	case constant.String:
		return constant.StringVal(val)
	}
	return nil
}

// -----------------------------------------------------------------------------
