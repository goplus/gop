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
	"github.com/goplus/gop/ast/fromgo"
	"github.com/goplus/gop/cl/internal/typesutil"
	"github.com/goplus/gop/token"
	"github.com/goplus/gox"
)

type goxRecorder struct {
	rec Recorder
}

// Member maps identifiers to the objects they denote.
func (p *goxRecorder) Member(id ast.Node, obj types.Object) {
	tv := typesutil.NewTypeAndValueForObject(obj)
	switch v := id.(type) {
	case *ast.SelectorExpr:
		sel := v.Sel
		// TODO: record event for a Go ident
		if _, ok := fromgo.CheckIdent(sel); !ok {
			p.rec.Use(sel, obj)
			p.rec.Type(v, tv)
		}
	case *ast.Ident: // it's in a classfile and impossible converted from Go
		p.rec.Use(v, obj)
		p.rec.Type(v, tv)
	}
}

type typesRecorder struct {
	Recorder
	types map[ast.Expr]types.TypeAndValue
}

func (rec *typesRecorder) Type(expr ast.Expr, tv types.TypeAndValue) {
	rec.types[expr] = tv
	rec.Recorder.Type(expr, tv)
}

func newTypeRecord(rec Recorder) *typesRecorder {
	return &typesRecorder{rec, make(map[ast.Expr]types.TypeAndValue)}
}

func (rec *typesRecorder) recordTypeValue(ctx *blockCtx, expr ast.Expr, mode typesutil.OperandMode) {
	e := ctx.cb.Get(-1)
	t, _ := gox.DerefType(e.Type)
	rec.Type(expr, typesutil.NewTypeAndValueForValue(t, e.CVal, mode))
}

func (rec *typesRecorder) indexExpr(ctx *blockCtx, expr *ast.IndexExpr) {
	if tv, ok := rec.types[expr.X]; ok {
		switch tv.Type.(type) {
		case *types.Map:
			rec.recordTypeValue(ctx, expr, typesutil.MapIndex)
			return
		case *types.Slice:
			rec.recordTypeValue(ctx, expr, typesutil.Variable)
			return
		}
	}
	switch expr.X.(type) {
	case *ast.CompositeLit:
		rec.recordTypeValue(ctx, expr, typesutil.Value)
	default:
		rec.recordTypeValue(ctx, expr, typesutil.Variable)
	}
}

func (rec *typesRecorder) unaryExpr(ctx *blockCtx, expr *ast.UnaryExpr) {
	switch expr.Op {
	case token.ARROW:
		rec.recordTypeValue(ctx, expr, typesutil.CommaOK)
	default:
		rec.recordTypeValue(ctx, expr, typesutil.Value)
	}
}

func makeSig(variadic bool, res types.Type, args ...types.Type) *types.Signature {
	list := make([]*types.Var, len(args))
	for i, param := range args {
		list[i] = types.NewVar(token.NoPos, nil, "", types.Default(param))
	}
	params := types.NewTuple(list...)
	var result *types.Tuple
	if res != nil {
		result = types.NewTuple(types.NewVar(token.NoPos, nil, "", res))
	}
	return types.NewSignature(nil, params, result, variadic)
}

func argsToTypes(args []*gox.Element) (typs []types.Type) {
	typs = make([]types.Type, len(args))
	for i, arg := range args {
		t := arg.Type
		switch x := t.(type) {
		case *gox.TypeType:
			t = x.Type()
		}
		typs[i] = t
	}
	return
}

func buildinId(fn ast.Expr) (id string) {
	switch v := fn.(type) {
	case *ast.Ident:
		id = v.Name
	}
	return
}

func typeIsString(t types.Type) bool {
	switch t := t.Underlying().(type) {
	case *types.Basic:
		return t.Info()&types.IsString != 0
	}
	return false
}

func (rec *typesRecorder) recordBuiltin(ctx *blockCtx, v *ast.CallExpr, retType types.Type, args []*gox.Element) {
	id := buildinId(v.Fun)
	var sig *types.Signature
	switch id {
	case "append":
		xtyp := args[0].Type
		switch xtyp.Underlying().(type) {
		case *types.Slice:
			if len(args) == 2 && v.Ellipsis.IsValid() && typeIsString(args[1].Type) {
				sig = makeSig(true, retType, xtyp, args[1].Type)
			} else {
				sig = makeSig(true, retType, xtyp, xtyp)
			}
		}
	default:
		sig = makeSig(false, retType, argsToTypes(args)...)
	}
	rec.Type(v.Fun, typesutil.NewTypeAndValueForValue(sig, nil, typesutil.Builtin))
}

func (rec *typesRecorder) recordCallExpr(ctx *blockCtx, v *ast.CallExpr, fnt types.Type, args []*gox.Element) {
	e := ctx.cb.Get(-1)
	f, ok := rec.types[v.Fun]
	if !ok {
		rec.Type(v.Fun, typesutil.NewTypeAndValueForValue(fnt, nil, typesutil.Value))
	} else if f.IsBuiltin() {
		rec.recordBuiltin(ctx, v, e.Type, args)
	}
	rec.Type(v, typesutil.NewTypeAndValueForCallResult(e.Type, e.CVal))
}

func (rec *typesRecorder) recordCompositeLit(ctx *blockCtx, v *ast.CompositeLit, typ types.Type) {
	rec.Type(v.Type, typesutil.NewTypeAndValueForType(typ))
	rec.Type(v, typesutil.NewTypeAndValueForValue(typ, nil, typesutil.Value))
}

func (rec *typesRecorder) recordType(ctx *blockCtx, typ ast.Expr, t types.Type) {
	rec.Type(typ, typesutil.NewTypeAndValueForType(t))
}

func (rec *typesRecorder) recordIdent(ctx *blockCtx, ident *ast.Ident, obj types.Object) {
	rec.Use(ident, obj)
	rec.Type(ident, typesutil.NewTypeAndValueForObject(obj))
}

func (rec *typesRecorder) recordFuncLit(ctx *blockCtx, v *ast.FuncLit, typ types.Type) {
	rec.Type(v.Type, typesutil.NewTypeAndValueForType(typ))
	rec.Type(v, typesutil.NewTypeAndValueForValue(typ, nil, typesutil.Value))
}

func (rec *typesRecorder) recordExpr(ctx *blockCtx, expr ast.Expr, rhs bool) {
	switch v := expr.(type) {
	case *ast.Ident:
	case *ast.BasicLit:
		rec.recordTypeValue(ctx, v, typesutil.Value)
	case *ast.CallExpr:
	case *ast.SelectorExpr:
		rec.recordTypeValue(ctx, v, typesutil.Variable)
	case *ast.BinaryExpr:
		rec.recordTypeValue(ctx, v, typesutil.Value)
	case *ast.UnaryExpr:
		rec.unaryExpr(ctx, v)
	case *ast.FuncLit:
	case *ast.CompositeLit:
	case *ast.SliceLit:
	case *ast.RangeExpr:
	case *ast.IndexExpr:
		rec.indexExpr(ctx, v)
	case *ast.IndexListExpr:
	case *ast.SliceExpr:
		rec.recordTypeValue(ctx, v, typesutil.Value)
	case *ast.StarExpr:
		rec.recordTypeValue(ctx, v, typesutil.Variable)
	case *ast.ArrayType:
	case *ast.MapType:
	case *ast.StructType:
	case *ast.ChanType:
	case *ast.InterfaceType:
	case *ast.ComprehensionExpr:
	case *ast.TypeAssertExpr:
		rec.recordTypeValue(ctx, v, typesutil.CommaOK)
	case *ast.ParenExpr:
		rec.recordTypeValue(ctx, v, typesutil.Value)
	case *ast.ErrWrapExpr:
	case *ast.FuncType:
	case *ast.Ellipsis:
	case *ast.KeyValueExpr:
	default:
	}
}
