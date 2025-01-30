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

package format

import (
	"go/token"
	"log"
	"reflect"

	"github.com/goplus/gop/ast"
)

// -----------------------------------------------------------------------------

func formatType(ctx *formatCtx, typ ast.Expr, ref *ast.Expr) {
	switch t := typ.(type) {
	case *ast.Ident, nil:
	case *ast.SelectorExpr:
		formatSelectorExpr(ctx, t, ref)
	case *ast.StarExpr:
		formatType(ctx, t.X, &t.X)
	case *ast.MapType:
		formatType(ctx, t.Key, &t.Key)
		formatType(ctx, t.Value, &t.Value)
	case *ast.StructType:
		formatFields(ctx, t.Fields)
	case *ast.ArrayType:
		formatExpr(ctx, t.Len, &t.Len)
		formatType(ctx, t.Elt, &t.Elt)
	case *ast.ChanType:
		formatType(ctx, t.Value, &t.Value)
	case *ast.InterfaceType:
		formatFields(ctx, t.Methods)
	case *ast.FuncType:
		formatFuncType(ctx, t)
	case *ast.Ellipsis:
		formatType(ctx, t.Elt, &t.Elt)
	default:
		log.Panicln("TODO: format -", reflect.TypeOf(typ))
	}
}

func formatFuncType(ctx *formatCtx, t *ast.FuncType) {
	formatFields(ctx, t.Params)
	formatFields(ctx, t.Results)
}

func formatFields(ctx *formatCtx, flds *ast.FieldList) {
	if flds != nil {
		for _, fld := range flds.List {
			formatField(ctx, fld)
		}
	}
}

func formatField(ctx *formatCtx, fld *ast.Field) {
	formatType(ctx, fld.Type, &fld.Type)
}

// -----------------------------------------------------------------------------

func formatExprs(ctx *formatCtx, exprs []ast.Expr) {
	for i, expr := range exprs {
		formatExpr(ctx, expr, &exprs[i])
	}
}

func formatExpr(ctx *formatCtx, expr ast.Expr, ref *ast.Expr) {
	switch v := expr.(type) {
	case *ast.Ident, *ast.BasicLit, *ast.BadExpr, nil:
	case *ast.BinaryExpr:
		formatExpr(ctx, v.X, &v.X)
		formatExpr(ctx, v.Y, &v.Y)
	case *ast.UnaryExpr:
		formatExpr(ctx, v.X, &v.X)
	case *ast.CallExpr:
		formatCallExpr(ctx, v)
	case *ast.SelectorExpr:
		formatSelectorExpr(ctx, v, ref)
	case *ast.SliceExpr:
		formatSliceExpr(ctx, v)
	case *ast.IndexExpr:
		formatExpr(ctx, v.X, &v.X)
		formatExpr(ctx, v.Index, &v.Index)
	case *ast.SliceLit:
		formatExprs(ctx, v.Elts)
	case *ast.CompositeLit:
		formatType(ctx, v.Type, &v.Type)
		formatExprs(ctx, v.Elts)
	case *ast.StarExpr:
		formatExpr(ctx, v.X, &v.X)
	case *ast.KeyValueExpr:
		formatExpr(ctx, v.Key, &v.Key)
		formatExpr(ctx, v.Value, &v.Value)
	case *ast.FuncLit:
		formatFuncType(ctx, v.Type)
		formatBlockStmt(ctx, v.Body)
	case *ast.TypeAssertExpr:
		formatExpr(ctx, v.X, &v.X)
		formatType(ctx, v.Type, &v.Type)
	case *ast.LambdaExpr:
		formatExprs(ctx, v.Rhs)
	case *ast.LambdaExpr2:
		formatBlockStmt(ctx, v.Body)
	case *ast.RangeExpr:
		formatRangeExpr(ctx, v)
	case *ast.ComprehensionExpr:
		formatComprehensionExpr(ctx, v)
	case *ast.ErrWrapExpr:
		formatExpr(ctx, v.X, &v.X)
		formatExpr(ctx, v.Default, &v.Default)
	case *ast.ParenExpr:
		formatExpr(ctx, v.X, &v.X)
	case *ast.Ellipsis:
		formatExpr(ctx, v.Elt, &v.Elt)
	default:
		formatType(ctx, expr, ref)
	}
}

func formatRangeExpr(ctx *formatCtx, v *ast.RangeExpr) {
	formatExpr(ctx, v.First, &v.First)
	formatExpr(ctx, v.Last, &v.Last)
	formatExpr(ctx, v.Expr3, &v.Expr3)
}

func formatComprehensionExpr(ctx *formatCtx, v *ast.ComprehensionExpr) {
	old := ctx.enterBlock()
	defer ctx.leaveBlock(old)

	formatForPhrases(ctx, v.Fors)
	formatExpr(ctx, v.Elt, &v.Elt)
}

func formatForPhrases(ctx *formatCtx, fors []*ast.ForPhrase) {
	for _, f := range fors {
		formatForPhrase(ctx, f)
	}
}

func formatForPhrase(ctx *formatCtx, v *ast.ForPhrase) {
	formatExpr(ctx, v.X, &v.X)
	formatStmt(ctx, v.Init)
	formatExpr(ctx, v.Cond, &v.Cond)
}

func formatSliceExpr(ctx *formatCtx, v *ast.SliceExpr) {
	formatExpr(ctx, v.X, &v.X)
	formatExpr(ctx, v.Low, &v.Low)
	formatExpr(ctx, v.High, &v.High)
	formatExpr(ctx, v.Max, &v.Max)
}

func formatCallExpr(ctx *formatCtx, v *ast.CallExpr) {
	formatExpr(ctx, v.Fun, &v.Fun)
	fncallStartingLowerCase(v)
	for i, arg := range v.Args {
		if fn, ok := arg.(*ast.FuncLit); ok {
			funcLitToLambdaExpr(ctx, fn, &v.Args[i])
		}
	}
	formatExprs(ctx, v.Args)
}

func funcLitToLambdaExpr(ctx *formatCtx, v *ast.FuncLit, ret *ast.Expr) {
	nres, named := checkResult(v.Type.Results)
	if len(named) > 0 {
		return
	}
	var lsh []*ast.Ident
	for _, p := range v.Type.Params.List {
		if p.Names == nil {
			lsh = append(lsh, ast.NewIdent("_"))
		} else {
			lsh = append(lsh, p.Names...)
		}
	}
	if len(lsh) > 0 && len(v.Body.List) == 1 {
		if stmt, ok := v.Body.List[0].(*ast.ReturnStmt); ok && len(stmt.Results) == nres {
			*ret = &ast.LambdaExpr{First: v.Pos(), Last: v.Pos(), Lhs: lsh, Rhs: stmt.Results, LhsHasParen: len(lsh) > 1, RhsHasParen: len(stmt.Results) > 1}
			return
		}
	}
	*ret = &ast.LambdaExpr2{Lhs: lsh, Body: v.Body, LhsHasParen: len(lsh) > 1}
}

func checkResult(v *ast.FieldList) (nres int, named []*ast.Ident) {
	if v != nil {
		for _, f := range v.List {
			if f.Names == nil {
				nres++
			} else {
				nres += len(f.Names)
				named = append(named, f.Names...)
			}
		}
	}
	return
}

func formatSelectorExpr(ctx *formatCtx, v *ast.SelectorExpr, ref *ast.Expr) {
	switch x := v.X.(type) {
	case *ast.Ident:
		if _, o := ctx.scope.LookupParent(x.Name, token.NoPos); o != nil {
			break
		}
		if imp, ok := ctx.imports[x.Name]; ok {
			if !fmtToBuiltin(imp, v.Sel, ref) {
				imp.isUsed = true
			}
		}
	default:
		formatExpr(ctx, x, &v.X)
	}
}

// -----------------------------------------------------------------------------

func formatBlockStmt(ctx *formatCtx, stmt *ast.BlockStmt) {
	if stmt != nil {
		old := ctx.enterBlock()
		defer ctx.leaveBlock(old)
		formatStmts(ctx, stmt.List)
	}
}

func formatStmts(ctx *formatCtx, stmts []ast.Stmt) {
	for _, stmt := range stmts {
		formatStmt(ctx, stmt)
	}
}

func formatStmt(ctx *formatCtx, stmt ast.Stmt) {
	switch v := stmt.(type) {
	case *ast.ExprStmt:
		formatExprStmt(ctx, v)
	case *ast.AssignStmt:
		formatAssignStmt(ctx, v)
	case *ast.IncDecStmt:
		formatExpr(ctx, v.X, &v.X)
	case *ast.ForStmt:
		formatForStmt(ctx, v)
	case *ast.RangeStmt:
		formatRangeStmt(ctx, v)
	case *ast.ForPhraseStmt:
		formatForPhraseStmt(ctx, v)
	case *ast.IfStmt:
		formatIfStmt(ctx, v)
	case *ast.CaseClause:
		formatExprs(ctx, v.List)
		formatStmts(ctx, v.Body)
	case *ast.SwitchStmt:
		formatSwitchStmt(ctx, v)
	case *ast.TypeSwitchStmt:
		formatTypeSwitchStmt(ctx, v)
	case *ast.CommClause:
		formatStmt(ctx, v.Comm)
		formatStmts(ctx, v.Body)
	case *ast.SelectStmt:
		formatBlockStmt(ctx, v.Body)
	case *ast.DeclStmt:
		formatDeclStmt(ctx, v)
	case *ast.ReturnStmt:
		formatExprs(ctx, v.Results)
	case *ast.BlockStmt:
		formatBlockStmt(ctx, v)
	case *ast.DeferStmt:
		formatCallExpr(ctx, v.Call)
	case *ast.GoStmt:
		formatCallExpr(ctx, v.Call)
	case *ast.SendStmt:
		formatExpr(ctx, v.Chan, &v.Chan)
		formatExpr(ctx, v.Value, &v.Value)
	case *ast.LabeledStmt:
		formatStmt(ctx, v.Stmt)
	case *ast.BranchStmt, *ast.EmptyStmt, nil, *ast.BadStmt:
	default:
		log.Panicln("TODO: formatStmt -", reflect.TypeOf(stmt))
	}
}

func formatExprStmt(ctx *formatCtx, v *ast.ExprStmt) {
	switch x := v.X.(type) {
	case *ast.CallExpr:
		commandStyleFirst(x)
	}
	formatExpr(ctx, v.X, &v.X)
}

func formatAssignStmt(ctx *formatCtx, v *ast.AssignStmt) {
	formatExprs(ctx, v.Lhs)
	formatExprs(ctx, v.Rhs)
}

func formatSwitchStmt(ctx *formatCtx, v *ast.SwitchStmt) {
	old := ctx.enterBlock()
	defer ctx.leaveBlock(old)

	formatStmt(ctx, v.Init)
	formatExpr(ctx, v.Tag, &v.Tag)
	formatBlockStmt(ctx, v.Body)
}

func formatTypeSwitchStmt(ctx *formatCtx, v *ast.TypeSwitchStmt) {
	old := ctx.enterBlock()
	defer ctx.leaveBlock(old)

	formatStmt(ctx, v.Init)
	formatStmt(ctx, v.Assign)
	formatBlockStmt(ctx, v.Body)
}

func formatIfStmt(ctx *formatCtx, v *ast.IfStmt) {
	old := ctx.enterBlock()
	defer ctx.leaveBlock(old)

	formatStmt(ctx, v.Init)
	formatExpr(ctx, v.Cond, &v.Cond)
	formatBlockStmt(ctx, v.Body)
	formatStmt(ctx, v.Else)
}

func formatRangeStmt(ctx *formatCtx, v *ast.RangeStmt) {
	old := ctx.enterBlock()
	defer ctx.leaveBlock(old)

	formatExpr(ctx, v.Key, &v.Key)
	formatExpr(ctx, v.Value, &v.Value)
	formatExpr(ctx, v.X, &v.X)
	formatBlockStmt(ctx, v.Body)
}

func formatForPhraseStmt(ctx *formatCtx, v *ast.ForPhraseStmt) {
	old := ctx.enterBlock()
	defer ctx.leaveBlock(old)

	formatForPhrase(ctx, v.ForPhrase)
	formatBlockStmt(ctx, v.Body)
}

func formatForStmt(ctx *formatCtx, v *ast.ForStmt) {
	old := ctx.enterBlock()
	defer ctx.leaveBlock(old)

	formatStmt(ctx, v.Init)
	formatExpr(ctx, v.Cond, &v.Cond)
	formatBlockStmt(ctx, v.Body)
}

func formatDeclStmt(ctx *formatCtx, v *ast.DeclStmt) {
	if decl, ok := v.Decl.(*ast.GenDecl); ok {
		formatGenDecl(ctx, decl)
	}
}

// -----------------------------------------------------------------------------
