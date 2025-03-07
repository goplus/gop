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
	"strings"

	"github.com/goplus/gogen"
	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/ast/fromgo"
	"github.com/goplus/gop/cl/internal/typesutil"
	"github.com/goplus/gop/token"
)

type goxRecorder struct {
	Recorder
	checkTypes   map[ast.Expr]types.TypeAndValue
	checkBuiltin map[*ast.Ident]types.Object
	referDefs    map[*ast.Ident]ast.Node
	referUses    map[string][]*ast.Ident
	builtin      *types.Package
}

func newRecorder(rec Recorder) *goxRecorder {
	checkTypes := make(map[ast.Expr]types.TypeAndValue)
	checkBuiltin := make(map[*ast.Ident]types.Object)
	referDefs := make(map[*ast.Ident]ast.Node)
	referUses := make(map[string][]*ast.Ident)
	builtin := types.NewPackage("gop/builtin", "builtin")
	return &goxRecorder{Recorder: rec, checkTypes: checkTypes, checkBuiltin: checkBuiltin,
		referDefs: referDefs, referUses: referUses, builtin: builtin}
}

// Refer uses maps identifiers to name for ast.OverloadFuncDecl.
func (p *goxRecorder) ReferUse(ident *ast.Ident, name string) {
	p.referUses[name] = append(p.referUses[name], ident)
}

// Refer def maps for ast.FuncLit or ast.OverloadFuncDecl.
func (p *goxRecorder) ReferDef(ident *ast.Ident, node ast.Node) {
	p.referDefs[ident] = node
}

// Complete computes the types record.
func (p *goxRecorder) Complete(pkg *gogen.Package) {
	scope := pkg.Types.Scope()
	for id, node := range p.referDefs {
		switch fn := node.(type) {
		case *ast.FuncLit:
			if obj := scope.Lookup(id.Name); obj != nil {
				p.recordFuncLit(fn, obj.Type())
				p.Implicit(node, obj)
			}
		case *ast.OverloadFuncDecl:
			if fn.Recv == nil {
				if obj := scope.Lookup(id.Name); obj != nil {
					p.Def(id, obj)
				}
			} else {
				if obj := scope.Lookup(fn.Recv.List[0].Type.(*ast.Ident).Name); obj != nil {
					if named, ok := obj.Type().(*types.Named); ok {
						n := named.NumMethods()
						for i := 0; i < n; i++ {
							if m := named.Method(i); m.Name() == id.Name {
								p.Def(id, m)
								break
							}
						}
					}
				}
			}
		}
	}
	for name, idents := range p.referUses {
		pos := strings.Index(name, ".")
		if pos == -1 {
			if obj := scope.Lookup(name); obj != nil {
				for _, id := range idents {
					p.Use(id, obj)
				}
			}
			continue
		}
		if obj := scope.Lookup(name[:pos]); obj != nil {
			if named, ok := obj.Type().(*types.Named); ok {
				n := named.NumMethods()
				for i := 0; i < n; i++ {
					if m := named.Method(i); m.Name() == name[pos+1:] {
						for _, id := range idents {
							p.Use(id, m)
						}
						break
					}
				}
			}
		}
	}
	p.checkTypes = nil
	p.checkBuiltin = nil
	p.referDefs = nil
	p.referUses = nil

	p.Builtin(func(universe *types.Scope) {
		scope := pkg.Builtin().Types.Scope()
		names := scope.Names()
		for _, name := range names {
			obj := scope.Lookup(name)
			if ast.IsExported(name) {
				continue
			}
			if isGopBuiltinFunc(name) {
				id := strings.Title(name)
				o := p.builtin.Scope().Lookup(id)
				if o == nil {
					if fns, ok := gogen.CheckOverloadFunc(obj.Type().(*types.Signature)); ok {
						o = types.NewFunc(token.NoPos, p.builtin, id, fns[0].Type().(*types.Signature))
					}
				}
				universe.Insert(o)
			} else if _, ok := obj.(*types.TypeName); ok && !gogen.IsTypeEx(obj.Type()) {
				universe.Insert(obj)
			}
		}
	})
}

func isGopBuiltinFunc(name string) bool {
	switch name {
	case "blines", "lines", "create", "open", "type", "echo", "errorf",
		"print", "println", "printf",
		"fprint", "fprintln", "fprintf",
		"sprint", "sprintln", "sprintf":
		return true
	}
	return false
}

// Member maps identifiers to the objects they denote.
func (p *goxRecorder) Member(id ast.Node, obj types.Object) {
	switch v := id.(type) {
	case *ast.SelectorExpr:
		sel := v.Sel
		// TODO: record event for a Go ident
		if _, ok := fromgo.CheckIdent(sel); !ok {
			var tv types.TypeAndValue
			// check v.X call result by value
			if f, ok := obj.(*types.Var); ok && f.IsField() && p.checkExprByValue(v.X) {
				tv = typesutil.NewTypeAndValueForValue(obj.Type(), nil, typesutil.Value)
			} else {
				tv = typesutil.NewTypeAndValueForObject(obj)
			}
			p.Use(sel, obj)
			p.Type(v, tv)
		}
	case *ast.Ident: // it's in a classfile and impossible converted from Go
		p.recordIdent(v, obj)
	}
}

func (p *goxRecorder) Call(id ast.Node, obj types.Object) {
	switch v := id.(type) {
	case *ast.Ident:
		p.recordIdent(v, obj)
	case *ast.SelectorExpr:
		p.Use(v.Sel, obj)
		p.Type(v, typesutil.NewTypeAndValueForObject(obj))
	case *ast.CallExpr:
		switch id := v.Fun.(type) {
		case *ast.Ident:
			p.recordIdent(id, obj)
			return
		case *ast.SelectorExpr:
			p.Use(id.Sel, obj)
		}
		p.Type(v.Fun, typesutil.NewTypeAndValueForObject(obj))
	}
}

func (rec *goxRecorder) checkExprByValue(v ast.Expr) bool {
	if tv, ok := rec.checkTypes[v]; ok {
		switch v.(type) {
		case *ast.CallExpr:
			if _, ok := tv.Type.(*types.Pointer); !ok {
				return true
			}
		default:
			if tv, ok := rec.checkTypes[v]; ok {
				return !tv.Addressable()
			}
		}
	}
	return false
}

func (rec *goxRecorder) Type(expr ast.Expr, tv types.TypeAndValue) {
	rec.checkTypes[expr] = tv
	rec.Recorder.Type(expr, tv)
}

func (rec *goxRecorder) instantiate(expr ast.Expr, _, typ types.Type) {
	// check gox TyOverloadNamed
	if tv, ok := rec.checkTypes[expr]; ok {
		tv.Type = typ
		rec.Recorder.Type(expr, tv)
	}
	var ident *ast.Ident
	switch id := expr.(type) {
	case *ast.Ident:
		ident = id
	case *ast.SelectorExpr:
		ident = id.Sel
	}
	if ident != nil {
		if named, ok := typ.(*types.Named); ok {
			rec.Use(ident, named.Obj())
		}
	}
}

func (rec *goxRecorder) recordTypeValue(ctx *blockCtx, expr ast.Expr, mode typesutil.OperandMode) {
	e := ctx.cb.Get(-1)
	t, _ := gogen.DerefType(e.Type)
	rec.Type(expr, typesutil.NewTypeAndValueForValue(t, e.CVal, mode))
}

func (rec *goxRecorder) indexExpr(ctx *blockCtx, expr *ast.IndexExpr) {
	if tv, ok := rec.checkTypes[expr.X]; ok {
		switch tv.Type.(type) {
		case *types.Map:
			rec.recordTypeValue(ctx, expr, typesutil.MapIndex)
			return
		case *types.Slice:
			rec.recordTypeValue(ctx, expr, typesutil.Variable)
			return
		}
	}
	op := typesutil.Variable
	switch e := expr.X.(type) {
	case *ast.CompositeLit:
		op = typesutil.Value
	case *ast.SelectorExpr:
		if rec.checkExprByValue(e.X) {
			op = typesutil.Value
		}
	}
	rec.recordTypeValue(ctx, expr, op)
}

func (rec *goxRecorder) unaryExpr(ctx *blockCtx, expr *ast.UnaryExpr) {
	switch expr.Op {
	case token.ARROW:
		rec.recordTypeValue(ctx, expr, typesutil.CommaOK)
	default:
		rec.recordTypeValue(ctx, expr, typesutil.Value)
	}
}

func (rec *goxRecorder) recordCallExpr(ctx *blockCtx, v *ast.CallExpr, fnt types.Type) {
	e := ctx.cb.Get(-1)
	if _, ok := rec.checkTypes[v.Fun]; !ok {
		rec.Type(v.Fun, typesutil.NewTypeAndValueForValue(fnt, nil, typesutil.Value))
	}
	rec.Type(v, typesutil.NewTypeAndValueForCallResult(e.Type, e.CVal))
}

func (rec *goxRecorder) recordCompositeLit(v *ast.CompositeLit, typ types.Type) {
	rec.Type(v.Type, typesutil.NewTypeAndValueForType(typ))
	rec.Type(v, typesutil.NewTypeAndValueForValue(typ, nil, typesutil.Value))
}

func (rec *goxRecorder) recordFuncLit(v *ast.FuncLit, typ types.Type) {
	rec.Type(v.Type, typesutil.NewTypeAndValueForType(typ))
	rec.Type(v, typesutil.NewTypeAndValueForValue(typ, nil, typesutil.Value))
}

func (rec *goxRecorder) recordType(typ ast.Expr, t types.Type) {
	rec.Type(typ, typesutil.NewTypeAndValueForType(t))
}

func isGopBuiltin(obj types.Object) bool {
	pkg := obj.Pkg()
	return pkg == nil || (pkg.Path() == "" && pkg.Name() == "")
}

func (rec *goxRecorder) recordIdent(ident *ast.Ident, obj types.Object) {
	// check go+ builtin
	if _, ok := obj.Type().(*gogen.TemplateSignature); ok {
		if o := types.Universe.Lookup(ident.Name); o != nil {
			obj = o
		}
	} else {
		switch ident.Name {
		case "complex", "real", "imag", "recover", "panic":
			if t, ok := rec.checkTypes[ident]; ok && t.IsBuiltin() {
				return
			}
			if isGopBuiltin(obj) {
				if o := types.Universe.Lookup(ident.Name); o != nil {
					obj = o
				}
			}
		default:
			if isGopBuiltinFunc(ident.Name) {
				if isGopBuiltin(obj) {
					rec.checkBuiltin[ident] = obj
					return
				} else if _, ok := rec.checkBuiltin[ident]; ok {
					id := strings.Title(ident.Name)
					o := rec.builtin.Scope().Lookup(id)
					if o == nil {
						o = types.NewFunc(token.NoPos, rec.builtin, id, obj.Type().(*types.Signature))
						rec.builtin.Scope().Insert(o)
					}
					obj = o
				}
			}
		}
	}
	rec.Use(ident, obj)
	rec.Type(ident, typesutil.NewTypeAndValueForObject(obj))
}

func (rec *goxRecorder) recordExpr(ctx *blockCtx, expr ast.Expr, _ bool) {
	switch v := expr.(type) {
	case *ast.Ident:
	case *ast.BasicLit:
		rec.recordTypeValue(ctx, v, typesutil.Value)
	case *ast.CallExpr:
	case *ast.SelectorExpr:
		if _, ok := rec.checkTypes[v]; !ok {
			rec.recordTypeValue(ctx, v, typesutil.Variable)
		}
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
