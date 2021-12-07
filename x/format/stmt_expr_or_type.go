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
	"github.com/goplus/gop/ast"
)

// -----------------------------------------------------------------------------

func formatType(ctx *formatCtx, typ ast.Expr, ref *ast.Expr) {
	switch t := typ.(type) {
	case *ast.SelectorExpr:
		formatSelectorExpr(ctx, t, ref)
	case *ast.StructType:
		formatStructType(ctx, t)
	case *ast.FuncType:
		formatFuncType(ctx, t)
	}
}

func formatFuncType(ctx *formatCtx, t *ast.FuncType) {
	formatFields(ctx, t.Params)
	formatFields(ctx, t.Results)
}

func formatStructType(ctx *formatCtx, t *ast.StructType) {
	formatFields(ctx, t.Fields)
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
	case *ast.CallExpr:
		formatCallExpr(ctx, v)
	case *ast.SelectorExpr:
		formatSelectorExpr(ctx, v, ref)
	}
}

func formatCallExpr(ctx *formatCtx, v *ast.CallExpr) {
	formatExpr(ctx, v.Fun, &v.Fun)
	formatExprs(ctx, v.Args)
}

func formatSelectorExpr(ctx *formatCtx, v *ast.SelectorExpr, ref *ast.Expr) {
	switch x := v.X.(type) {
	case *ast.Ident:
		// TODO: maybe have bugs: `x.Name` may be a variable, not a import.
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
	case *ast.DeclStmt:
		formatDeclStmt(ctx, v)
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

func formatDeclStmt(ctx *formatCtx, v *ast.DeclStmt) {
	if decl, ok := v.Decl.(*ast.GenDecl); ok {
		formatGenDecl(ctx, decl)
	}
}

// -----------------------------------------------------------------------------
