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

package golang

import (
	"go/ast"
	"go/token"
	"reflect"

	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

// MapType instr
func MapType(p *Builder, typ reflect.Type) *ast.MapType {
	key := Type(p, typ.Key())
	val := Type(p, typ.Elem())
	return &ast.MapType{Key: key, Value: val}
}

// ArrayType instr
func ArrayType(p *Builder, typ reflect.Type) *ast.ArrayType {
	var len ast.Expr
	var kind = typ.Kind()
	if kind == reflect.Array {
		len = IntConst(int64(typ.Len()))
	}
	return &ast.ArrayType{Len: len, Elt: Type(p, typ.Elem())}
}

// PtrType instr
func PtrType(p *Builder, typElem reflect.Type) ast.Expr {
	elt := Type(p, typElem)
	return &ast.StarExpr{X: elt}
}

// Field instr
func Field(p *Builder, name string, typ reflect.Type, tag string, ellipsis bool) *ast.Field {
	var names []*ast.Ident
	var ftag *ast.BasicLit
	var typExpr ast.Expr
	if name != "" {
		names = []*ast.Ident{Ident(name)}
	}
	if tag != "" {
		ftag = StringConst(tag)
	}
	if ellipsis {
		typExpr = &ast.Ellipsis{Elt: Type(p, typ.Elem())}
	} else {
		typExpr = Type(p, typ)
	}
	return &ast.Field{Names: names, Type: typExpr, Tag: ftag}
}

// Methods instr
func Methods(p *Builder, typ reflect.Type) []*ast.Field {
	n := typ.NumMethod()
	if n == 0 {
		return nil
	}
	panic("Methods: todo")
}

// InterfaceType instr
func InterfaceType(p *Builder, typ reflect.Type) *ast.InterfaceType {
	return &ast.InterfaceType{
		Methods: &ast.FieldList{List: Methods(p, typ)},
	}
}

// FuncType instr
func FuncType(p *Builder, typ reflect.Type) *ast.FuncType {
	numIn, numOut := typ.NumIn(), typ.NumOut()
	variadic := typ.IsVariadic()
	var opening token.Pos
	var params, results []*ast.Field
	if numIn > 0 {
		params = make([]*ast.Field, numIn)
		if variadic {
			numIn--
		}
		for i := 0; i < numIn; i++ {
			params[i] = Field(p, "", typ.In(i), "", false)
		}
		if variadic {
			params[numIn] = Field(p, "", typ.In(numIn), "", true)
		}
	}
	if numOut > 0 {
		results = make([]*ast.Field, numOut)
		for i := 0; i < numOut; i++ {
			results[i] = Field(p, "", typ.Out(i), "", false)
		}
		opening++
	}
	return &ast.FuncType{
		Params:  &ast.FieldList{Opening: 1, List: params, Closing: 1},
		Results: &ast.FieldList{Opening: opening, List: results, Closing: opening},
	}
}

// Type instr
func Type(p *Builder, typ reflect.Type) ast.Expr {
	pkgPath, name := typ.PkgPath(), typ.Name()
	log.Debug(typ, "-", "pkgPath:", pkgPath, "name:", name)
	if name != "" {
		if pkgPath != "" {
			pkg := p.Import(pkgPath)
			return &ast.SelectorExpr{X: Ident(pkg), Sel: Ident(name)}
		}
		return Ident(name)
	}
	kind := typ.Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		return ArrayType(p, typ)
	case reflect.Map:
		return MapType(p, typ)
	case reflect.Ptr:
		return PtrType(p, typ.Elem())
	case reflect.Func:
		return FuncType(p, typ)
	case reflect.Interface:
		return InterfaceType(p, typ)
	case reflect.Chan:
	case reflect.Struct:
	}
	log.Panicln("Type: unknown type -", typ)
	return nil
}

// -----------------------------------------------------------------------------
