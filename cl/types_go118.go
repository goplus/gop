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
	if named, ok := typ.(*types.Named); ok {
		ctx.loadNamed(ctx.pkg, named)
	}
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
	var typeParams []*types.TypeParam
	if recv != nil {
		if tparams := recv.Type().(*types.Named).TypeParams(); tparams != nil {
			n := tparams.Len()
			typeParams = make([]*types.TypeParam, n)
			for i := 0; i < n; i++ {
				tp := tparams.At(i)
				typeParams[i] = types.NewTypeParam(tp.Obj(), tp.Constraint())
			}
		}
	} else {
		typeParams = toTypeParams(ctx, typ.TypeParams)
	}
	if len(typeParams) > 0 {
		ctx.tlookup = &typeParamLookup{typeParams}
		defer func() {
			ctx.tlookup = nil
		}()
	}
	params, variadic := toParams(ctx, typ.Params.List)
	results := toResults(ctx, typ.Results)
	if recv != nil {
		return types.NewSignatureType(recv, typeParams, nil, params, results, variadic)
	}
	return types.NewSignatureType(recv, nil, typeParams, params, results, variadic)
}

type typeParamLookup struct {
	typeParams []*types.TypeParam
}

func (p *typeParamLookup) Lookup(name string) types.Type {
	for _, t := range p.typeParams {
		if name == t.Obj().Name() {
			return t
		}
	}
	return nil
}

func initType(ctx *blockCtx, named *types.Named, spec *ast.TypeSpec) {
	typeParams := toTypeParams(ctx, spec.TypeParams)
	if len(typeParams) > 0 {
		named.SetTypeParams(typeParams)
		ctx.tlookup = &typeParamLookup{typeParams}
		defer func() {
			ctx.tlookup = nil
		}()
	}
	typ := toType(ctx, spec.Type)
	if named, ok := typ.(*types.Named); ok {
		typ = getUnderlying(ctx, named)
	}
	named.SetUnderlying(typ)
}

func getRecvType(typ ast.Expr) ast.Expr {
	if t, ok := typ.(*ast.StarExpr); ok {
		typ = t.X
	}
	switch t := typ.(type) {
	case *ast.IndexExpr:
		return t.X
	case *ast.IndexListExpr:
		return t.X
	}
	return typ
}
