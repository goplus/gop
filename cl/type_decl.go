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
	"github.com/goplus/reflectx"
	"github.com/qiniu/x/log"
)

type iType interface {
	Kind() reflect.Kind
	Elem() reflect.Type
}

// -----------------------------------------------------------------------------

type unboundArrayType struct {
	elem reflect.Type
}

func (p *unboundArrayType) Kind() reflect.Kind {
	return reflect.Array
}

func (p *unboundArrayType) Elem() reflect.Type {
	return p.elem
}

// -----------------------------------------------------------------------------

func toTypeEx(ctx *blockCtx, typ ast.Expr) (t iType, variadic bool) {
	if t = toType(ctx, typ); t != nil {
		return
	}
	if v, ok := typ.(*ast.Ellipsis); ok {
		elem := toType(ctx, v.Elt)
		return reflect.SliceOf(elem.(reflect.Type)), true
	}
	log.Panicln("toType: unknown -", reflect.TypeOf(typ))
	return nil, false
}

func toType(ctx *blockCtx, typ ast.Expr) iType {
	switch v := typ.(type) {
	case *ast.Ident:
		return toIdentType(ctx, v.Name)
	case *ast.SelectorExpr:
		return toExternalType(ctx, v)
	case *ast.StarExpr:
		elem := toType(ctx, v.X)
		return reflect.PtrTo(elem.(reflect.Type))
	case *ast.ArrayType:
		return toArrayType(ctx, v)
	case *ast.FuncType:
		return toFuncType(ctx, v)
	case *ast.InterfaceType:
		return toInterfaceType(ctx, v)
	case *ast.StructType:
		return toStructType(ctx, v)
	case *ast.MapType:
		key := toType(ctx, v.Key)
		elem := toType(ctx, v.Value)
		return reflect.MapOf(key.(reflect.Type), elem.(reflect.Type))
	case *ast.ChanType:
		val := toType(ctx, v.Value)
		return reflect.ChanOf(toChanDir(v.Dir), val.(reflect.Type))
	case *ast.Ellipsis:
		return nil
	}
	log.Panicln("toType: unknown -", reflect.TypeOf(typ))
	return nil
}

func toChanDir(dir ast.ChanDir) (ret reflect.ChanDir) {
	if (dir & ast.SEND) != 0 {
		ret |= reflect.SendDir
	}
	if (dir & ast.RECV) != 0 {
		ret |= reflect.RecvDir
	}
	return
}

func toFuncType(ctx *blockCtx, t *ast.FuncType) iType {
	in, variadic := toTypesEx(ctx, t.Params)
	out := toTypes(ctx, t.Results)
	return reflect.FuncOf(in, out, variadic)
}

func buildFuncType(recv *ast.FieldList, fi exec.FuncInfo, ctx *blockCtx, t *ast.FuncType) {
	in, args, variadic := toArgTypes(ctx, recv, t.Params)
	rets := toReturnTypes(ctx, t.Results)
	if variadic {
		fi.Vargs(in...)
	} else {
		fi.Args(in...)
	}
	fi.Return(rets...)
	ctx.insertFuncVars(fi, in, args, rets)
}

func toTypes(ctx *blockCtx, fields *ast.FieldList) (types []reflect.Type) {
	if fields == nil {
		return
	}
	for _, field := range fields.List {
		n := len(field.Names)
		if n == 0 {
			n = 1
		}
		typ := toType(ctx, field.Type)
		if typ == nil {
			log.Panicln("toType: unknown -", reflect.TypeOf(field.Type))
		}
		for i := 0; i < n; i++ {
			types = append(types, typ.(reflect.Type))
		}
	}
	return
}

func toTypesEx(ctx *blockCtx, fields *ast.FieldList) ([]reflect.Type, bool) {
	var types []reflect.Type
	last := len(fields.List) - 1
	for i := 0; i <= last; i++ {
		field := fields.List[i]
		n := len(field.Names)
		if n == 0 {
			n = 1
		}
		typ, variadic := toTypeEx(ctx, field.Type)
		for i := 0; i < n; i++ {
			types = append(types, typ.(reflect.Type))
		}
		if variadic {
			if i != last {
				log.Panicln("toTypes failed: the variadic type isn't last argument?")
			}
			return types, true
		}
	}
	return types, false
}

func toReturnTypes(ctx *blockCtx, fields *ast.FieldList) (vars []exec.Var) {
	if fields == nil {
		return
	}
	index := 0
	for _, field := range fields.List {
		n := len(field.Names)
		typ := toType(ctx, field.Type)
		if typ == nil {
			log.Panicln("toType: unknown -", reflect.TypeOf(field.Type))
		}
		if n == 0 {
			index++
			vars = append(vars, ctx.NewVar(typ.(reflect.Type), strconv.Itoa(index)))
		} else {
			for i := 0; i < n; i++ {
				vars = append(vars, ctx.NewVar(typ.(reflect.Type), field.Names[i].Name))
			}
			index += n
		}
	}
	return
}

func toArgTypes(ctx *blockCtx, recv, fields *ast.FieldList) ([]reflect.Type, []string, bool) {
	var types []reflect.Type
	var names []string
	if recv != nil {
		for _, fld := range recv.List[0].Names {
			names = append(names, fld.Name)
			break
		}
		typ, _ := toTypeEx(ctx, recv.List[0].Type)
		types = append(types, typ.(reflect.Type))
	}
	last := len(fields.List) - 1
	for i := 0; i <= last; i++ {
		field := fields.List[i]
		n := len(field.Names)
		if n == 0 {
			names = append(names, "")
			n = 1
		} else {
			for _, fld := range field.Names {
				names = append(names, fld.Name)
			}
		}
		typ, variadic := toTypeEx(ctx, field.Type)
		for i := 0; i < n; i++ {
			types = append(types, typ.(reflect.Type))
		}
		if variadic {
			if i != last {
				log.Panicln("toTypes failed: the variadic type isn't last argument?")
			}
			return types, names, true
		}
	}
	return types, names, false
}

func toStructType(ctx *blockCtx, v *ast.StructType) iType {
	var fields []reflect.StructField
	for _, field := range v.Fields.List {
		fields = append(fields, toStructField(ctx, field)...)
	}
	return reflectx.StructOf(fields)
}

func toInterfaceType(ctx *blockCtx, v *ast.InterfaceType) iType {
	methods := v.Methods.List
	if methods == nil {
		return exec.TyEmptyInterface
	}
	panic("toInterfaceType: todo")
}

func toExternalType(ctx *blockCtx, v *ast.SelectorExpr) iType {
	if ident, ok := v.X.(*ast.Ident); ok {
		if sym, ok := ctx.find(ident.Name); ok {
			switch t := sym.(type) {
			case string:
				pkg := ctx.FindGoPackage(t)
				if pkg == nil {
					log.Panicln("toExternalType failed: package not found -", v)
				}
				if typ, ok := pkg.FindType(v.Sel.Name); ok {
					return typ
				}
			}
		}
	}
	panic("toExternalType: todo")
}

func toIdentType(ctx *blockCtx, ident string) iType {
	if typ, ok := ctx.builtin.FindType(ident); ok {
		return typ
	}
	typ, err := ctx.findType(ident)
	if err != nil {
		log.Panicln("toIdentType failed: findType error", err)
		return nil
	}
	return typ.Type
}

func toArrayType(ctx *blockCtx, v *ast.ArrayType) iType {
	elem := toType(ctx, v.Elt)
	if v.Len == nil {
		return reflect.SliceOf(elem.(reflect.Type))
	}
	compileExpr(ctx, v.Len)
	n := ctx.infer.Pop()
	if nv, ok := n.(iValue).(*constVal); ok {
		if iv, ok := nv.v.(int64); ok {
			if iv < 0 {
				return &unboundArrayType{elem: elem.(reflect.Type)}
			}
			return reflect.ArrayOf(int(iv), elem.(reflect.Type))
		}
	}
	log.Panicln("toArrayType failed: unknown -", reflect.TypeOf(n))
	return nil
}

func toBoundArrayLen(ctx *blockCtx, v *ast.CompositeLit) int {
	n := -1
	max := -1
	for _, elt := range v.Elts {
		if e, ok := elt.(*ast.KeyValueExpr); ok {
			i := toInt(ctx, e.Key)
			if i < n {
				max = n
			}
			n = i
		} else {
			n++
		}
	}
	if n < max {
		n = max
	}
	return n + 1
}

func toInt(ctx *blockCtx, e ast.Expr) int {
	compileExpr(ctx, e)
	nv, ok := ctx.infer.Pop().(*constVal)
	if !ok {
		log.Panicln("toInt: require constant expr.")
	}
	iv, ok := nv.v.(int64)
	if !ok {
		log.Panicln("toInt: constant expr isn't an integer.")
	}
	return int(iv)
}

func toStructField(ctx *blockCtx, field *ast.Field) []reflect.StructField {
	var fields []reflect.StructField
	if field.Names == nil {
		fields = append(fields, buildField(ctx, field, true, ""))
	} else {
		for _, name := range field.Names {
			fields = append(fields, buildField(ctx, field, false, name.Name))
		}
	}
	return fields
}

func typeName(typ ast.Expr) string {
	switch v := typ.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		return typeName(v.Sel)
	case *ast.StarExpr:
		return typeName(v.X)
	}
	return ""
}

func buildField(ctx *blockCtx, field *ast.Field, anonymous bool, fieldName string) reflect.StructField {
	var f = reflect.StructField{}
	if !anonymous {
		f = reflect.StructField{}
	}
	if anonymous {
		fieldName = typeName(field.Type)
	}
	f = reflect.StructField{
		Name:      fieldName,
		Type:      toType(ctx, field.Type).(reflect.Type),
		Anonymous: anonymous,
	}
	if fieldName != "" {
		c := fieldName[0]
		if 'a' <= c && c <= 'z' || c == '_' {
			f.PkgPath = ctx.pkg.Name
		}
	}
	if field.Tag != nil {
		tag, _ := strconv.Unquote(field.Tag.Value)
		f.Tag = reflect.StructTag(tag)
	}
	return f
}

// -----------------------------------------------------------------------------

type typeDecl struct {
	Methods map[string]*funcDecl
	Type    reflect.Type
	Name    string
	Alias   bool
}

type methodDecl struct {
	recv    string // recv object name
	pointer int
	typ     *ast.FuncType
	body    *ast.BlockStmt
	file    *fileCtx
}

// FuncDecl represents a function declaration.
type FuncDecl struct {
	typ      *ast.FuncType
	body     *ast.BlockStmt
	ctx      *blockCtx
	recv     *ast.FieldList
	fi       exec.FuncInfo
	used     bool
	cached   bool
	compiled bool
	reserved bool
}

type funcDecl = FuncDecl

func newFuncDecl(name string, recv *ast.FieldList, typ *ast.FuncType, body *ast.BlockStmt, ctx *blockCtx) *FuncDecl {
	nestDepth := ctx.getNestDepth()
	var isMethod int
	if recv != nil {
		isMethod = 1
	}
	fi := ctx.NewFunc(name, nestDepth, isMethod)
	return &FuncDecl{typ: typ, recv: recv, body: body, ctx: ctx, fi: fi}
}

// Get returns function information.
func (p *FuncDecl) Get() exec.FuncInfo {
	if !p.cached {
		buildFuncType(p.recv, p.fi, p.ctx, p.typ)
		p.cached = true
	}
	return p.fi
}

// Type returns the type of this function.
func (p *FuncDecl) Type() reflect.Type {
	return p.Get().Type()
}

// Compile compiles this function
func (p *FuncDecl) Compile() exec.FuncInfo {
	fun := p.Get()
	if !p.compiled {
		ctx := p.ctx
		out := ctx.out
		out.DefineFunc(fun)
		ctx.funcCtx = newFuncCtx(fun)
		compileBlockStmtWithout(ctx, p.body)
		ctx.funcCtx.checkLabels()
		ctx.funcCtx = nil
		out.EndFunc(fun)
		p.compiled = true
	}
	return fun
}

// -----------------------------------------------------------------------------
