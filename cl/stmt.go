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
	"strconv"

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
	case *ast.GoStmt:
		compileGoStmt(ctx, v)
	case *ast.DeclStmt:
		compileDeclStmt(ctx, v)
	case *ast.SendStmt:
		compileSendStmt(ctx, v)
	case *ast.TypeSwitchStmt:
		compileTypeSwitchStmt(ctx, v)
	case *ast.EmptyStmt:
		// do nothing
	default:
		log.Panicln("compileStmt failed: unknown -", reflect.TypeOf(v))
	}
	ctx.out.EndStmt(stmt, start)
}

func compileForPhraseStmt(parent *blockCtx, v *ast.ForPhraseStmt) {
	if v.Cond != nil {
		v.Body.List = append([]ast.Stmt{&ast.IfStmt{
			If: v.TokPos,
			Cond: &ast.UnaryExpr{
				Op: token.NOT,
				X:  v.Cond,
			},
			Body: &ast.BlockStmt{List: []ast.Stmt{&ast.BranchStmt{Tok: token.CONTINUE, TokPos: v.TokPos}}},
		}}, v.Body.List...)
	}
	rangeStmt := &ast.RangeStmt{
		For:    v.For,
		Key:    v.Key,
		Value:  v.Value,
		TokPos: v.TokPos,
		Tok:    token.DEFINE,
		X:      v.X,
		Body:   v.Body,
	}
	if parent.currentLabel != nil && parent.currentLabel.Stmt == v {
		parent.currentLabel.Stmt = rangeStmt
	}
	compileRangeStmt(parent, rangeStmt)
}

func compileRangeStmt(parent *blockCtx, v *ast.RangeStmt) {
	kvDef := map[string]reflect.Type{}
	newIterName := "_gop_NewIter"
	nextName := "_gop_Next"
	keyName := "_gop_Key"
	valueName := "_gop_Value"
	iter := &ast.Ident{Name: "_gop_iter", NamePos: v.For}

	var keyIdent, valIdent *ast.Ident
	switch v.Tok {
	case token.DEFINE:
		keyIdent = toIdent(v.Key)
		valIdent = toIdent(v.Value)
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
				keyIdent = &ast.Ident{Name: k0.Name, Obj: k0}
				lhs[idx], rhs[idx] = v.Key, keyIdent
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
				valIdent = &ast.Ident{Name: v0.Name, Obj: v0}
				lhs[idx], rhs[idx] = v.Value, valIdent
				idx++
			}
		}
		v.Body.List = append([]ast.Stmt{&ast.AssignStmt{
			Lhs: lhs[0:idx],
			Tok: token.ASSIGN,
			Rhs: rhs[0:idx],
		}}, v.Body.List...)
	}

	compileExpr(parent, v.X)
	typData := boundType(parent.infer.Pop().(iValue))
	var typKey, typVal reflect.Type
	var forStmts []ast.Stmt
	switch kind := typData.Kind(); kind {
	case reflect.String:
		typKey = exec.TyInt
		typVal = exec.TyByte
	case reflect.Slice, reflect.Array:
		typKey = exec.TyInt
		typVal = typData.Elem()
	case reflect.Map:
		typKey = typData.Key()
		typVal = typData.Elem()
	default:
		log.Panicln("compileRangeStmt: require slice, array or map")
	}
	if keyIdent != nil {
		// Key(iter,&k)
		if id, ok := v.Key.(*ast.Ident); !ok || (ok && id.Name != "_") {
			kvDef[keyIdent.Name] = typKey
			forStmts = append(forStmts, &ast.ExprStmt{X: &ast.CallExpr{
				Fun: &ast.Ident{Name: keyName, NamePos: v.For},
				Args: []ast.Expr{iter, &ast.UnaryExpr{
					Op: token.AND,
					X:  keyIdent,
				}},
			}})
		}
	}
	if valIdent != nil {
		// Value(iter,&v)
		if id, ok := v.Value.(*ast.Ident); !ok || (ok && id.Name != "_") {
			kvDef[valIdent.Name] = typVal
			forStmts = append(forStmts, &ast.ExprStmt{X: &ast.CallExpr{
				Fun: &ast.Ident{Name: valueName, NamePos: v.For},
				Args: []ast.Expr{iter, &ast.UnaryExpr{
					Op: token.AND,
					X:  valIdent,
				}},
			}})
		}
	}
	// iter:=_gop_NewIter(obj)
	init := &ast.AssignStmt{
		Lhs: []ast.Expr{iter},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{&ast.CallExpr{
			Fun:    &ast.Ident{Name: newIterName},
			Lparen: v.For,
			Rparen: v.For,
			Args:   []ast.Expr{v.X},
		}},
		TokPos: v.For,
	}
	// range => for
	fs := &ast.ForStmt{
		For: v.For,
		Cond: &ast.CallExpr{
			Fun:    &ast.Ident{Name: nextName},
			Lparen: v.For,
			Rparen: v.For,
			Args:   []ast.Expr{iter},
		},
		Body: &ast.BlockStmt{
			Lbrace: v.Body.Lbrace,
			List:   append(forStmts, v.Body.List...),
			Rbrace: v.Body.Rbrace,
		},
	}
	parent.out.DefineBlock()
	defer parent.out.EndBlock()
	ctx := newNormBlockCtx(parent)
	if ctx.currentLabel != nil && ctx.currentLabel.Stmt == v {
		ctx.currentLabel.Stmt = fs
	}
	for k, v := range kvDef {
		ctx.insertVar(k, v)
	}
	compileStmt(ctx, init)
	compileStmt(ctx, fs)
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
	if v.Cond != nil {
		compileExpr(ctx, v.Cond)()
		checkBool(ctx.infer.Pop())
		out.JmpIf(0, done)
	}
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
		label := ctx.getBreakLabel(labelName)
		if label != nil {
			ctx.out.Jmp(label)
			return
		}
		log.Panicln("break statement out of for/switch/select statements")
	case token.CONTINUE:
		var labelName string
		if v.Label != nil {
			labelName = v.Label.Name
		}
		label := ctx.getContinueLabel(labelName)
		if label != nil {
			ctx.out.Jmp(label)
			return
		}
		log.Panicln("continue statement out of for statements")
	}
}

func compileLabeledStmt(ctx *blockCtx, v *ast.LabeledStmt) {
	label := ctx.defineLabel(v.Label.Name)
	ctx.out.Label(label)
	ctx.currentLabel = v
	compileStmt(ctx, v.Stmt)
}

type callType int

const (
	callExpr callType = iota
	callByDefer
	callByGo
)

var gCallTypes = []string{
	"",
	"defer",
	"go",
}

func compileGoStmt(ctx *blockCtx, v *ast.GoStmt) {
	compileCallExpr(ctx, v.Call, callByGo)()
}

func compileDeferStmt(ctx *blockCtx, v *ast.DeferStmt) {
	compileCallExpr(ctx, v.Call, callByDefer)()
}

func unusedIdent(ctx *blockCtx, ident string) string {
	name := ident
	var ok bool
	var i int
	for {
		_, ok = ctx.find(name)
		if !ok {
			break
		}
		name = ident + "_" + strconv.Itoa(i)
		i++
	}
	return name
}

func isNilExpr(expr ast.Expr) bool {
	if ident, ok := expr.(*ast.Ident); ok && ident.Name == "nil" {
		return true
	}
	return false
}

func buildTypeSwitchStmtDefault(ctx *blockCtx, c *ast.CaseClause, vExpr, xExpr ast.Expr, hasValue bool) ast.Stmt {
	if hasValue {
		stms := []ast.Stmt{
			&ast.AssignStmt{
				Lhs: []ast.Expr{vExpr},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{xExpr},
			},
			&ast.AssignStmt{
				Lhs: []ast.Expr{ast.NewIdent("_")},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{vExpr},
			},
		}
		return &ast.BlockStmt{
			List: append(stms, &ast.BlockStmt{List: c.Body}),
		}
	} else {
		return &ast.BlockStmt{
			List: c.Body,
		}
	}
}

func buildTypeSwitchStmtCaseStmt(ctx *blockCtx, c *ast.CaseClause, body *ast.BlockStmt, vExpr, xExpr ast.Expr, hasValue bool) *ast.IfStmt {
	vCond := ast.NewIdent(unusedIdent(ctx, "ok"))
	if isNilExpr(c.List[0]) {
		if hasValue {
			return &ast.IfStmt{
				If: c.Pos(),
				Init: &ast.AssignStmt{
					Lhs: []ast.Expr{
						vExpr,
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						xExpr,
					},
				},
				Cond: &ast.BinaryExpr{
					X:  vExpr,
					Op: token.EQL,
					Y:  c.List[0],
				},
				Body: body,
			}
		} else {
			return &ast.IfStmt{
				If: c.Pos(),
				Cond: &ast.BinaryExpr{
					X:  xExpr,
					Op: token.EQL,
					Y:  c.List[0],
				},
				Body: body,
			}
		}
	} else {
		return &ast.IfStmt{
			If: c.Pos(),
			Init: &ast.AssignStmt{
				Lhs: []ast.Expr{
					vExpr,
					vCond,
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.TypeAssertExpr{
						X:    xExpr,
						Type: c.List[0],
					},
				},
			},
			Cond: vCond,
			Body: body,
		}
	}
}

func buildTypeSwitchStmtCaseListStmt(ctx *blockCtx, c *ast.CaseClause, body *ast.BlockStmt, vExpr, xExpr ast.Expr, hasValue bool) *ast.IfStmt {
	cond := ast.NewIdent("ok")
	var stmts []ast.Stmt
	for _, expr := range c.List {
		if isNilExpr(expr) {
			stmts = append(stmts, &ast.IfStmt{
				Cond: &ast.BinaryExpr{
					X:  xExpr,
					Op: token.EQL,
					Y:  expr,
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ReturnStmt{
							Results: []ast.Expr{ast.NewIdent("true")},
						},
					},
				},
			})
		} else {
			stmts = append(stmts, &ast.IfStmt{
				Init: &ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("_"),
						cond,
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.TypeAssertExpr{
							X:    xExpr,
							Type: expr,
						},
					},
				},
				Cond: cond,
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ReturnStmt{},
					},
				},
			})
		}
	}
	funLit := &ast.FuncLit{
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{cond},
						Type:  ast.NewIdent("bool"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: append(stmts, &ast.ReturnStmt{}),
		},
	}
	ifstmt := &ast.IfStmt{
		If: c.Pos(),
		Cond: &ast.CallExpr{
			Fun: funLit,
		},
		Body: body,
	}
	if hasValue {
		ifstmt.Init = &ast.AssignStmt{
			Lhs: []ast.Expr{
				vExpr,
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				xExpr,
			},
		}
	}
	return ifstmt
}

func buildTypeSwitchStmtCase(ctx *blockCtx, c *ast.CaseClause, vExpr, xExpr ast.Expr, hasValue bool) *ast.IfStmt {
	var body *ast.BlockStmt
	if hasValue {
		stmts := []ast.Stmt{
			&ast.AssignStmt{
				Lhs: []ast.Expr{ast.NewIdent("_")},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{vExpr},
			},
		}
		body = &ast.BlockStmt{
			List: append(stmts, &ast.BlockStmt{List: c.Body}),
		}
	} else {
		body = &ast.BlockStmt{
			List: []ast.Stmt{&ast.BlockStmt{List: c.Body}},
		}
	}
	if len(c.List) > 1 {
		return buildTypeSwitchStmtCaseListStmt(ctx, c, body, vExpr, xExpr, hasValue)
	} else {
		return buildTypeSwitchStmtCaseStmt(ctx, c, body, vExpr, xExpr, hasValue)
	}
}

func compileTypeSwitchStmt(ctx *blockCtx, v *ast.TypeSwitchStmt) {
	ctx.out.DefineBlock()
	defer ctx.out.EndBlock()
	if v.Init != nil {
		compileStmt(ctx, v.Init)
	}
	if len(v.Body.List) == 0 {
		return
	}
	uExpr := ast.NewIdent("_")
	var xInitExpr ast.Expr
	var vExpr ast.Expr
	var xExpr ast.Expr
	var vName string
	var hasValue bool
	switch assign := v.Assign.(type) {
	case *ast.AssignStmt:
		hasValue = true
		vExpr = assign.Lhs[0]
		vName = vExpr.(*ast.Ident).Name
		xInitExpr = assign.Rhs[0].(*ast.TypeAssertExpr).X
	case *ast.ExprStmt:
		vExpr = uExpr
		xInitExpr = assign.X.(*ast.TypeAssertExpr).X
	}
	compileExpr(ctx, xInitExpr)
	xtyp := ctx.infer.Pop().(iValue).Type()
	xExpr = ast.NewIdent(unusedIdent(ctx, "_gop_"+vName))
	vinit := &ast.AssignStmt{
		Lhs: []ast.Expr{xExpr},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{xInitExpr},
	}
	vused := &ast.AssignStmt{
		Lhs: []ast.Expr{uExpr},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{xExpr},
	}
	compileStmt(ctx, vinit)
	compileStmt(ctx, vused)
	var ifStmt *ast.IfStmt
	var lastIfStmt *ast.IfStmt
	var defaultStmt ast.Stmt
	dupcheck := make(map[reflect.Type]bool)
	for _, item := range v.Body.List {
		c, _ := item.(*ast.CaseClause)
		if c.List == nil {
			if defaultStmt != nil {
				log.Panicf("multiple defaults in switch")
			}
			defaultStmt = buildTypeSwitchStmtDefault(ctx, c, vExpr, xExpr, hasValue)
			continue
		}
		for _, expr := range c.List {
			if isNilExpr(expr) {
				continue
			}
			typ := toType(ctx, expr).(reflect.Type)
			if xtyp.NumMethod() != 0 && typ.Kind() != reflect.Interface {
				if !typ.Implements(xtyp) {
					log.Panicf("impossible type switch case: %v (type %v) cannot have dynamic type %v",
						ctx.code(xInitExpr), xtyp, typ)
				}
			}
			if _, ok := dupcheck[typ]; ok {
				log.Panicf("duplicate case %v in type switch", typ)
			}
			dupcheck[typ] = true
		}
		stmt := buildTypeSwitchStmtCase(ctx, c, vExpr, xExpr, hasValue)
		if ifStmt == nil {
			ifStmt = stmt
		}
		if lastIfStmt != nil {
			lastIfStmt.Else = stmt
		}
		lastIfStmt = stmt
	}
	if ifStmt == nil {
		if defaultStmt != nil {
			compileStmt(ctx, defaultStmt)
		}
		return
	}
	if defaultStmt != nil {
		lastIfStmt.Else = &ast.BlockStmt{
			List: []ast.Stmt{defaultStmt},
		}
	}
	compileIfStmt(ctx, ifStmt)
}

func compileSwitchStmt(ctx *blockCtx, v *ast.SwitchStmt) {
	if init := v.Init; init != nil {
		v.Init = nil
		block := &ast.BlockStmt{List: []ast.Stmt{init, v}}
		compileNewBlock(ctx, block)
		return
	}
	var defaultBody []ast.Stmt
	ctxSw := ctx
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
				xcase := ctx.infer.Pop()
				checkCaseCompare(tag, xcase, out)
				if err := checkOpMatchType(exec.OpEQ, tag, xcase); err != nil {
					log.Panicf("invalid case %v in switch on %v (%v)", ctx.code(caseExp), ctx.code(v.Tag), err)
				}
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
	if init := v.Init; init != nil {
		v.Init = nil
		block := &ast.BlockStmt{List: []ast.Stmt{init, v}}
		compileNewBlock(ctx, block)
		return
	}
	var done exec.Label
	ctxIf := ctx
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
	if n == 1 && ctx.infer.Len() == 1 {
		if ret, ok := ctx.infer.Get(-1).(*funcResults); ok {
			n := ret.NumValues()
			if fun.NumOut() != n {
				log.Panicln("compileReturnStmt failed: mismatched count of return values -", fun.Name())
			}
			ctx.infer.SetLen(0)
			ctx.out.Return(int32(n))
			return
		}
	}

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
	if ctx.infer.Len() > 0 {
		in := ctx.infer.Get(-1)
		if v, ok := in.(*constVal); ok {
			for i := 0; i < v.NumValues(); i++ {
				checkType(exec.TyEmptyInterface, v.Value(i), ctx.out)
			}
		}
	}
	ctx.infer.PopN(1)
}

func compileDeclStmt(ctx *blockCtx, expr *ast.DeclStmt) {
	switch d := expr.Decl.(type) {
	case *ast.GenDecl:
		switch d.Tok {
		case token.TYPE:
			loadTypes(ctx, d)
		case token.CONST:
			loadConsts(ctx, d)
		case token.VAR:
			loadVars(ctx, d, expr)
		default:
			log.Panicln("tok:", d.Tok, "spec:", reflect.TypeOf(d.Specs).Elem())
		}
	}
}

func compileSendStmt(ctx *blockCtx, expr *ast.SendStmt) {
	compileExpr(ctx, expr.Chan)()
	compileExpr(ctx, expr.Value)()
	checkType(ctx.infer.Get(-2).(iValue).Type().Elem(), ctx.infer.Get(-1), ctx.out)
	ctx.out.Send()
	ctx.infer.PopN(2)
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
		rhsExpr := expr.Rhs[0]
		if len(expr.Lhs) == 2 {
			switch ie := rhsExpr.(type) {
			case *ast.IndexExpr:
				rhsExpr = &ast.TwoValueIndexExpr{IndexExpr: ie}
			case *ast.TypeAssertExpr:
				rhsExpr = &ast.TwoValueTypeAssertExpr{TypeAssertExpr: ie}
			}
		}
		compileExpr(ctx, rhsExpr)()
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
	count := len(expr.Lhs)
	ctx.underscore = 0
	for i := len(expr.Lhs) - 1; i >= 0; i-- {
		compileExprLHS(ctx, expr.Lhs[i], expr.Tok)
	}
	if ctx.underscore == count && expr.Tok == token.DEFINE {
		log.Panicln("no new variables on left side of :=")
	}
}

// -----------------------------------------------------------------------------
