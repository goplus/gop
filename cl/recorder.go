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
	goast "go/ast"
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

func findExprObject(expr goast.Expr, objs ...types.Object) (types.Object, bool) {
	switch id := expr.(type) {
	case *goast.Ident:
		for _, obj := range objs {
			if obj.Name() == id.Name {
				return obj, true
			}
		}
	case *goast.SelectorExpr:
		for _, obj := range objs {
			if obj.Name() == id.Sel.Name {
				return obj, true
			}
		}
	}
	return nil, false
}

func findOverloadObject(fnt types.Type, expr goast.Expr) (types.Object, bool) {
	if call, ok := expr.(*goast.CallExpr); ok {
		switch v := fnt.(type) {
		case *types.Signature:
			if recv := v.Recv(); recv != nil {
				switch t := recv.Type().(type) {
				case *gox.TyOverloadFunc:
					return findExprObject(call.Fun, t.Funcs...)
				case *gox.TyOverloadMethod:
					return findExprObject(call.Fun, t.Methods...)
				case *gox.TyTemplateRecvMethod:
					if tsig, ok := t.Func.Type().(*types.Signature); ok {
						if trecv := tsig.Recv(); trecv != nil {
							if t, ok := trecv.Type().(*gox.TyOverloadFunc); ok {
								return findExprObject(call.Fun, t.Funcs...)
							}
						}
					}
					return t.Func, true
				}
			}
		}
	}
	return nil, false
}

func (rec *typesRecorder) recordCallExpr(ctx *blockCtx, v *ast.CallExpr, fnt types.Type) {
	e := ctx.cb.Get(-1)
	obj, found := findOverloadObject(fnt, e.Val)
	if found {
		switch id := v.Fun.(type) {
		case *ast.Ident:
			rec.Use(id, obj)
		case *ast.SelectorExpr:
			rec.Use(id.Sel, obj)
		}
		rec.Type(v.Fun, typesutil.NewTypeAndValueForObject(obj))
	} else {
		if _, ok := rec.types[v.Fun]; !ok {
			rec.Type(v.Fun, typesutil.NewTypeAndValueForValue(fnt, nil, typesutil.Value))
		}
	}
	rec.Type(v, typesutil.NewTypeAndValueForCallResult(e.Type, e.CVal))
}

func (rec *typesRecorder) recordCompositeLit(ctx *blockCtx, v *ast.CompositeLit, typ types.Type) {
	rec.Type(v.Type, typesutil.NewTypeAndValueForType(typ))
	rec.Type(v, typesutil.NewTypeAndValueForValue(typ, nil, typesutil.Value))
}

func (rec *typesRecorder) recordFuncLit(ctx *blockCtx, v *ast.FuncLit, typ types.Type) {
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
		rec.recordTypeValue(ctx, v, typesutil.TypExpr)
	case *ast.MapType:
		rec.recordTypeValue(ctx, v, typesutil.TypExpr)
	case *ast.StructType:
		rec.recordTypeValue(ctx, v, typesutil.TypExpr)
	case *ast.ChanType:
		rec.recordTypeValue(ctx, v, typesutil.TypExpr)
	case *ast.InterfaceType:
		rec.recordTypeValue(ctx, v, typesutil.TypExpr)
	case *ast.ComprehensionExpr:
	case *ast.TypeAssertExpr:
		rec.recordTypeValue(ctx, v, typesutil.CommaOK)
	case *ast.ParenExpr:
		rec.recordTypeValue(ctx, v, typesutil.Value)
	case *ast.ErrWrapExpr:
	case *ast.FuncType:
		rec.recordTypeValue(ctx, v, typesutil.TypExpr)
	case *ast.Ellipsis:
	case *ast.KeyValueExpr:
	default:
	}
}
