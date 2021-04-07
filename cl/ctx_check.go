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

// Debug only
func Debug(pkg *Package) {
	isLower("")
	if _, v, ok := pkg.Find("main"); ok {
		fn := v.(*funcDecl)
		isNoExecCtx(fn.ctx, fn.body)
	}
}

func newBlockCtxWithFlag(parent *blockCtx) *blockCtx {
	ctx := newNormBlockCtx(parent)
	ctx.checkFlag = true
	return ctx
}

func isNoExecCtx(parent *blockCtx, body *ast.BlockStmt) bool {
	ctx := newBlockCtxWithFlag(parent)
	for _, stmt := range body.List {
		if noExecCtx := isNoExecCtxStmt(ctx, stmt); !noExecCtx {
			return false
		}
	}
	return true
}

func isNoExecCtxStmt(ctx *blockCtx, stmt ast.Stmt) bool {
	switch v := stmt.(type) {
	case *ast.ExprStmt:
		return isNoExecCtxExpr(ctx, v.X)
	case *ast.AssignStmt:
		return isNoExecCtxAssignStmt(ctx, v)
	case *ast.IfStmt:
		return isNoExecCtxIfStmt(ctx, v)
	case *ast.ForPhraseStmt:
		return isNoExecCtxForPhraseStmt(ctx, v)
	case *ast.RangeStmt:
		return isNoExecCtxRangeStmt(ctx, v)
	case *ast.ForStmt:
		return isNoExecForStmt(ctx, v)
	case *ast.SwitchStmt:
		return isNoExecCtxSwitchStmt(ctx, v)
	case *ast.TypeSwitchStmt:
		return isNoExecCtxTypeSwitchStmt(ctx, v)
	case *ast.BlockStmt:
		return isNoExecCtx(ctx, v)
	case *ast.ReturnStmt:
		return isNoExecCtxExprs(ctx, v.Results)
	case *ast.IncDecStmt:
		return isNoExecCtxExpr(ctx, v.X)
	case *ast.BranchStmt:
		return true
	case *ast.LabeledStmt:
		return true
	case *ast.DeferStmt:
		return isNoExecCtxCallExpr(ctx, v.Call)
	case *ast.GoStmt:
		return isNoExecCtxCallExpr(ctx, v.Call)
	case *ast.DeclStmt:
		return isNoExecCtxDeclStmt(ctx, v)
	case *ast.SendStmt:
		return true
	default:
		log.Panicln("isNoExecCtxStmt failed: unknown -", reflect.TypeOf(v))
	}
	return true
}

func isNoExecCtxExpr(ctx *blockCtx, expr ast.Expr) bool {
	switch v := expr.(type) {
	case *ast.Ident:
		return true
	case *ast.BasicLit:
		return true
	case *ast.CallExpr:
		return isNoExecCtxCallExpr(ctx, v)
	case *ast.BinaryExpr:
		return isNoExecCtx2nd(ctx, v.X, v.Y)
	case *ast.UnaryExpr:
		return isNoExecCtxExpr(ctx, v.X)
	case *ast.SelectorExpr:
		return isNoExecCtxExpr(ctx, v.X)
	case *ast.ParenExpr:
		return isNoExecCtxExpr(ctx, v.X)
	case *ast.ErrWrapExpr:
		return isNoExecCtx2nd(ctx, v.X, v.Default)
	case *ast.IndexExpr:
		return isNoExecCtx2nd(ctx, v.X, v.Index)
	case *ast.SliceExpr:
		return isNoExecCtxSliceExpr(ctx, v)
	case *ast.CompositeLit:
		return isNoExecCtxExprs(ctx, v.Elts)
	case *ast.SliceLit:
		return isNoExecCtxExprs(ctx, v.Elts)
	case *ast.FuncLit:
		return isNoExecCtxFuncLit(ctx, v)
	case *ast.ListComprehensionExpr:
		return isNoExecCtxListComprehensionExpr(ctx, v)
	case *ast.MapComprehensionExpr:
		return isNoExecCtxMapComprehensionExpr(ctx, v)
	case *ast.StarExpr:
		return isNoExecCtxStarExpr(ctx, v)
	case *ast.TypeAssertExpr:
		return isNoExecCtxExpr(ctx, v.X)
	case *ast.InterfaceType:
		return true
	case *ast.ChanType:
		return true
	case *ast.MapType:
		return true
	case *ast.FuncType:
		return true
	case *ast.ArrayType:
		return true
	case *ast.Ellipsis:
		return true
	case *ast.KeyValueExpr:
		return isNoExecCtx2nd(ctx, v.Key, v.Value)
	default:
		log.Panicln("isNoExecCtxExpr failed: unknown -", reflect.TypeOf(v))
	}
	return true
}

func isNoExecCtxForPhrase(parent *blockCtx, f ast.ForPhrase) (*blockCtx, bool) {
	ctx := newBlockCtxWithFlag(parent)
	if noExecCtx := isNoExecCtxExpr(parent, f.X); !noExecCtx {
		return ctx, false
	}
	if f.Key != nil {
		ctx.insertVar(f.Key.Name, exec.TyEmptyInterface, true)
	}
	if f.Value != nil {
		ctx.insertVar(f.Value.Name, exec.TyEmptyInterface, true)
	}
	if f.Cond != nil {
		return ctx, isNoExecCtxExpr(ctx, f.Cond)
	}
	return ctx, true
}

func isNoExecCtxForPhrases(ctx *blockCtx, fors []ast.ForPhrase) (*blockCtx, bool) {
	var noExecCtx bool
	for i := len(fors) - 1; i >= 0; i-- {
		if ctx, noExecCtx = isNoExecCtxForPhrase(ctx, fors[i]); !noExecCtx {
			return ctx, false
		}
	}
	return ctx, true
}

func isNoExecCtxForPhraseStmt(parent *blockCtx, v *ast.ForPhraseStmt) bool {
	ctx, noExecCtx := isNoExecCtxForPhrase(parent, v.ForPhrase)
	if !noExecCtx {
		return false
	}
	return isNoExecCtx(ctx, v.Body)
}

func isNoExecCtxRangeStmt(parent *blockCtx, v *ast.RangeStmt) bool {
	ctx := newBlockCtxWithFlag(parent)
	for _, e := range []ast.Expr{v.X, v.Key, v.Value} {
		if e == nil {
			continue
		}
		if v.Tok == token.DEFINE {
			if id, ok := e.(*ast.Ident); ok {
				ctx.insertVar(id.Name, exec.TyEmptyInterface, true)
				continue
			}
		}
		if noExecCtx := isNoExecCtxExpr(parent, e); !noExecCtx {
			return false
		}
	}
	return isNoExecCtx(ctx, v.Body)
}

func isNoExecForStmt(parent *blockCtx, v *ast.ForStmt) bool {
	ctx := newBlockCtxWithFlag(parent)
	if v.Cond != nil {
		if noExecCtx := isNoExecCtxExpr(ctx, v.Cond); !noExecCtx {
			return false
		}
	}
	for _, e := range []ast.Stmt{v.Init, v.Post} {
		if e == nil {
			continue
		}
		if noExecCtx := isNoExecCtxStmt(parent, e); !noExecCtx {
			return false
		}
	}
	return isNoExecCtx(ctx, v.Body)
}

func isNoExecCtxListComprehensionExpr(parent *blockCtx, v *ast.ListComprehensionExpr) bool {
	ctx, noExecCtx := isNoExecCtxForPhrases(parent, v.Fors)
	if !noExecCtx {
		return false
	}
	return isNoExecCtxExpr(ctx, v.Elt)
}

func isNoExecCtxMapComprehensionExpr(parent *blockCtx, v *ast.MapComprehensionExpr) bool {
	ctx, noExecCtx := isNoExecCtxForPhrases(parent, v.Fors)
	if !noExecCtx {
		return false
	}
	elt := v.Elt
	return isNoExecCtx2nd(ctx, elt.Key, elt.Value)
}

func isNoExecCtxSliceExpr(ctx *blockCtx, v *ast.SliceExpr) bool {
	if noExecCtx := isNoExecCtxExpr(ctx, v.X); !noExecCtx {
		return false
	}
	if v.Low != nil {
		if noExecCtx := isNoExecCtxExpr(ctx, v.Low); !noExecCtx {
			return false
		}
	}
	if v.High != nil {
		if noExecCtx := isNoExecCtxExpr(ctx, v.High); !noExecCtx {
			return false
		}
	}
	if v.Max != nil {
		if noExecCtx := isNoExecCtxExpr(ctx, v.Max); !noExecCtx {
			return false
		}
	}
	return true
}

func isNoExecCtxFuncLit(ctx *blockCtx, v *ast.FuncLit) bool {
	// TODO: log.Warn("isNoExecCtxFuncLit: to be optimized")
	return false
}

func isNoExecCtx2nd(ctx *blockCtx, a, b ast.Expr) bool {
	if noExecCtx := isNoExecCtxExpr(ctx, a); !noExecCtx {
		return false
	}
	if b == nil {
		return true
	}
	return isNoExecCtxExpr(ctx, b)
}

func isNoExecCtxStarExpr(ctx *blockCtx, v *ast.StarExpr) bool {
	if noExecCtx := isNoExecCtxExpr(ctx, v.X); !noExecCtx {
		return false
	}
	return true
}

func isNoExecCtxCallExpr(ctx *blockCtx, v *ast.CallExpr) bool {
	switch expr := v.Fun.(type) {
	case *ast.Ident:
		switch expr.Name {
		case "make":
			return isNoExecCtxExprs(ctx, v.Args[1:])
		case "new":
			return true
		}
	}
	if noExecCtx := isNoExecCtxExpr(ctx, v.Fun); !noExecCtx {
		return false
	}
	return isNoExecCtxExprs(ctx, v.Args)
}

func isNoExecCtxExprs(ctx *blockCtx, exprs []ast.Expr) bool {
	for _, expr := range exprs {
		if noExecCtx := isNoExecCtxExpr(ctx, expr); !noExecCtx {
			return false
		}
	}
	return true
}

func isNoExecCtxTypeSwitchStmt(ctx *blockCtx, v *ast.TypeSwitchStmt) bool {
	ctxSw := ctx
	if v.Assign != nil {
		if noExecCtx := isNoExecCtxStmt(ctxSw, v.Assign); !noExecCtx {
			return false
		}
	}
	for _, item := range v.Body.List {
		c, ok := item.(*ast.CaseClause)
		if !ok {
			log.Panicln("compile SwitchStmt failed: case clause expected.")
		}
		if noExecCtx := isNoExecCtxExprs(ctxSw, c.List); !noExecCtx {
			return false
		}
		ctxBody := newBlockCtxWithFlag(ctxSw)
		for _, stmt := range c.Body {
			if noExecCtx := isNoExecCtxStmt(ctxBody, stmt); !noExecCtx {
				return false
			}
		}
	}
	return true
}

func isNoExecCtxSwitchStmt(ctx *blockCtx, v *ast.SwitchStmt) bool {
	ctxSw := ctx
	if v.Tag != nil {
		if noExecCtx := isNoExecCtxExpr(ctxSw, v.Tag); !noExecCtx {
			return false
		}
	}
	for _, item := range v.Body.List {
		c, ok := item.(*ast.CaseClause)
		if !ok {
			log.Panicln("compile SwitchStmt failed: case clause expected.")
		}
		if noExecCtx := isNoExecCtxExprs(ctxSw, c.List); !noExecCtx {
			return false
		}
		ctxBody := newBlockCtxWithFlag(ctxSw)
		for _, stmt := range c.Body {
			if noExecCtx := isNoExecCtxStmt(ctxBody, stmt); !noExecCtx {
				return false
			}
		}
	}
	return true
}

func isNoExecCtxIfStmt(ctx *blockCtx, v *ast.IfStmt) bool {
	ctxIf := ctx
	if noExecCtx := isNoExecCtxExpr(ctxIf, v.Cond); !noExecCtx {
		return false
	}
	ctxWith := newBlockCtxWithFlag(ctxIf)
	if noExecCtx := isNoExecCtxStmt(ctxWith, v.Body); !noExecCtx {
		return false
	}
	if v.Else != nil {
		return isNoExecCtxStmt(ctxIf, v.Else)
	}
	return true
}

func isNoExecCtxAssignStmt(ctx *blockCtx, expr *ast.AssignStmt) bool {
	if noExecCtx := isNoExecCtxExprs(ctx, expr.Rhs); !noExecCtx {
		return false
	}
	for i := len(expr.Lhs) - 1; i >= 0; i-- {
		if noExecCtx := isNoExecCtxExprLHS(ctx, expr.Lhs[i], expr.Tok); !noExecCtx {
			return false
		}
	}
	return true
}

func isNoExecCtxDeclStmt(ctx *blockCtx, expr *ast.DeclStmt) bool {
	switch d := expr.Decl.(type) {
	case *ast.GenDecl:
		switch d.Tok {
		case token.VAR, token.CONST:
			for _, spec := range d.Specs {
				vs := spec.(*ast.ValueSpec)
				if vs.Values != nil {
					if noExecCtx := isNoExecCtxExprs(ctx, vs.Values); !noExecCtx {
						return false
					}
				}
				for i := len(vs.Names) - 1; i >= 0; i-- {
					if noExecCtx := isNoExecCtxExprLHS(ctx, vs.Names[i], lhsAssign); !noExecCtx {
						return false
					}
				}
			}
		}
	}
	return true
}

func isNoExecCtxExprLHS(ctx *blockCtx, expr ast.Expr, mode compileMode) bool {
	switch v := expr.(type) {
	case *ast.Ident:
		return isNoExecCtxIdentLHS(ctx, v.Name, mode)
	case *ast.IndexExpr:
		return isNoExecCtxIndexExprLHS(ctx, v, mode)
	case *ast.SelectorExpr:
		return isNoExecCtxSelectorExprLHS(ctx, v, mode)
	case *ast.StarExpr:
		return isNoExecCtxStarExprLHS(ctx, v, mode)
	default:
		log.Panicln("isNoExecCtxExprLHS failed: unknown -", reflect.TypeOf(v))
	}
	return true
}

func isNoExecCtxIndexExprLHS(ctx *blockCtx, v *ast.IndexExpr, mode compileMode) bool {
	if noExecCtx := isNoExecCtxExpr(ctx, v.X); !noExecCtx {
		return false
	}
	return isNoExecCtxExpr(ctx, v.Index)
}

func isNoExecCtxIdentLHS(ctx *blockCtx, name string, mode compileMode) bool {
	if mode == lhsDefine && !ctx.exists(name) {
		ctx.insertVar(name, exec.TyEmptyInterface, true)
	}
	return true
}

func isNoExecCtxSelectorExprLHS(ctx *blockCtx, v *ast.SelectorExpr, mode compileMode) bool {
	if noExecCtx := isNoExecCtxExpr(ctx, v.X); !noExecCtx {
		return false
	}
	return true
}

func isNoExecCtxStarExprLHS(ctx *blockCtx, v *ast.StarExpr, mode compileMode) bool {
	if noExecCtx := isNoExecCtxExpr(ctx, v.X); !noExecCtx {
		return false
	}
	return true
}

// -----------------------------------------------------------------------------
