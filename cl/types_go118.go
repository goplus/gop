//go:build go1.18
// +build go1.18

/*
 * Copyright (c) 2022 The GoPlus Authors (goplus.org). All rights reserved.
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
	"go/types"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/token"
	"github.com/goplus/gox"
)

func toTermList(ctx *blockCtx, expr ast.Expr) []*types.Term {
retry:
	switch v := expr.(type) {
	case *ast.UnaryExpr:
		if v.Op != token.TILDE {
			panic(ctx.newCodeErrorf(v.Pos(), "invalid op %v must ~", v.Op))
		}
		return []*types.Term{types.NewTerm(true, toType(ctx, v.X))}
	case *ast.BinaryExpr:
		if v.Op != token.OR {
			panic(ctx.newCodeErrorf(v.Pos(), "invalid op %v must |", v.Op))
		}
		return append(toTermList(ctx, v.X), toTermList(ctx, v.Y)...)
	case *ast.ParenExpr:
		expr = v.X
		goto retry
	}
	return []*types.Term{types.NewTerm(false, toType(ctx, expr))}
}

func toBinaryExprType(ctx *blockCtx, v *ast.BinaryExpr) types.Type {
	return types.NewInterfaceType(nil, []types.Type{types.NewUnion(toTermList(ctx, v))})
}

func toUnaryExprType(ctx *blockCtx, v *ast.UnaryExpr) types.Type {
	return types.NewInterfaceType(nil, []types.Type{types.NewUnion(toTermList(ctx, v))})
}

func toTypeParam(ctx *blockCtx, fld *ast.Field) *types.TypeParam {
	typ := toType(ctx, fld.Type)
	obj := types.NewTypeName(fld.Pos(), ctx.pkg.Types, fld.Names[0].Name, nil)
	return types.NewTypeParam(obj, typ)
}

func toTypeParams(ctx *blockCtx, params *ast.FieldList) []*types.TypeParam {
	if params == nil {
		return nil
	}
	n := len(params.List)
	list := make([]*types.TypeParam, n, n)
	for i, fld := range params.List {
		list[i] = toTypeParam(ctx, fld)
	}
	return list
}

func toFuncType(ctx *blockCtx, typ *ast.FuncType, recv *types.Var) *types.Signature {
	typeParams := toTypeParams(ctx, typ.TypeParams)
	params, variadic := toParams(ctx, typ.Params.List, typeParams)
	results := toResults(ctx, typ.Results, typeParams)
	return types.NewSignatureType(recv, nil, typeParams, params, results, variadic)
}

func toResults(ctx *blockCtx, in *ast.FieldList, typeParams []*types.TypeParam) *types.Tuple {
	if in == nil {
		return nil
	}
	flds := in.List
	n := len(flds)
	args := make([]*types.Var, 0, n)
	for _, fld := range flds {
		args = toParam(ctx, fld, args, typeParams)
	}
	return types.NewTuple(args...)
}

func toParams(ctx *blockCtx, flds []*ast.Field, typeParams []*types.TypeParam) (typ *types.Tuple, variadic bool) {
	n := len(flds)
	if n == 0 {
		return nil, false
	}
	args := make([]*types.Var, 0, n)
	for _, fld := range flds {
		args = toParam(ctx, fld, args, typeParams)
	}
	_, ok := flds[n-1].Type.(*ast.Ellipsis)
	return types.NewTuple(args...), ok
}

func toTypeCheckTypeParams(ctx *blockCtx, typ ast.Expr, typeParams []*types.TypeParam) types.Type {
	if ident, ok := typ.(*ast.Ident); ok {
		for _, t := range typeParams {
			if ident.Name == t.Obj().Name() {
				return t
			}
		}
		return toIdentType(ctx, ident)
	}
	return toType(ctx, typ)
}

func toParam(ctx *blockCtx, fld *ast.Field, args []*gox.Param, typeParams []*types.TypeParam) []*gox.Param {
	typ := toTypeCheckTypeParams(ctx, fld.Type, typeParams)
	pkg := ctx.pkg
	if len(fld.Names) == 0 {
		return append(args, pkg.NewParam(fld.Pos(), "", typ))
	}
	for _, name := range fld.Names {
		args = append(args, pkg.NewParam(name.Pos(), name.Name, typ))
	}
	return args
}
