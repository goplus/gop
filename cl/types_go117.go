//go:build !go1.18
// +build !go1.18

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
	"go/types"

	"github.com/goplus/gop/ast"
)

func toBinaryExprType(ctx *blockCtx, v *ast.BinaryExpr) types.Type {
	panic("type parameters are unsupported at this go version")
}

func toUnaryExprType(ctx *blockCtx, v *ast.UnaryExpr) types.Type {
	panic("type parameters are unsupported at this go version")
}

type typeParamLookup struct {
}

func (p *typeParamLookup) Lookup(name string) types.Type {
	return nil
}

func toFuncType(ctx *blockCtx, typ *ast.FuncType, recv *types.Var) *types.Signature {
	params, variadic := toParams(ctx, typ.Params.List)
	results := toResults(ctx, typ.Results)
	return types.NewSignature(recv, params, results, variadic)
}

func initType(ctx *blockCtx, named *types.Named, spec *ast.TypeSpec) {
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
	return typ
}
