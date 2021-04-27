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

package astutil

import (
	"path"
	"reflect"
	"strconv"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/token"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

// StructType instr
func StructType(typ reflect.Type) ast.Expr {
	var fields = &ast.FieldList{}
	for i := 0; i < typ.NumField(); i++ {
		fields.List = append(fields.List, toField(typ.Field(i)))
	}
	return &ast.StructType{Fields: fields}
}

func toField(f reflect.StructField) *ast.Field {
	var field = &ast.Field{}
	if !f.Anonymous {
		field.Names = []*ast.Ident{
			Ident(f.Name),
		}
	}
	field.Type = Type(f.Type)
	return field
}

// ChanType instr
func ChanType(typ reflect.Type) *ast.ChanType {
	var dir ast.ChanDir
	switch typ.ChanDir() {
	case reflect.RecvDir:
		dir = ast.RECV
	case reflect.SendDir:
		dir = ast.SEND
	case reflect.BothDir:
		dir = ast.RECV | ast.SEND
	}
	return &ast.ChanType{Dir: dir, Value: Ident(typ.Elem().String())}
}

// MapType instr
func MapType(typ reflect.Type) *ast.MapType {
	key := Type(typ.Key())
	val := Type(typ.Elem())
	return &ast.MapType{Key: key, Value: val}
}

// ArrayType instr
func ArrayType(typ reflect.Type) *ast.ArrayType {
	var len ast.Expr
	var kind = typ.Kind()
	if kind == reflect.Array {
		len = IntConst(int64(typ.Len()))
	}
	return &ast.ArrayType{Len: len, Elt: Type(typ.Elem())}
}

// PtrType instr
func PtrType(typElem reflect.Type) ast.Expr {
	elt := Type(typElem)
	return &ast.StarExpr{X: elt}
}

// Field instr
func Field(name string, typ reflect.Type, tag string, ellipsis bool) *ast.Field {
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
		typExpr = &ast.Ellipsis{Elt: Type(typ.Elem())}
	} else {
		typExpr = Type(typ)
	}
	return &ast.Field{Names: names, Type: typExpr, Tag: ftag}
}

// Methods instr
func Methods(typ reflect.Type) []*ast.Field {
	n := typ.NumMethod()
	if n == 0 {
		return nil
	}
	fields := make([]*ast.Field, n, n)
	for i := 0; i < n; i++ {
		m := typ.Method(i)
		fields[i] = &ast.Field{
			Names: []*ast.Ident{&ast.Ident{Name: m.Name}},
			Type:  FuncType(m.Type),
		}
	}
	return fields
}

// InterfaceType instr
func InterfaceType(typ reflect.Type) *ast.InterfaceType {
	return &ast.InterfaceType{
		Methods: &ast.FieldList{List: Methods(typ)},
	}
}

// FuncType instr
func FuncType(typ reflect.Type) *ast.FuncType {
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
			params[i] = Field("", typ.In(i), "", false)
		}
		if variadic {
			params[numIn] = Field("", typ.In(numIn), "", true)
		}
	}
	if numOut > 0 {
		results = make([]*ast.Field, numOut)
		for i := 0; i < numOut; i++ {
			results[i] = Field("", typ.Out(i), "", false)
		}
		opening++
	}
	return &ast.FuncType{
		Params:  &ast.FieldList{Opening: 1, List: params, Closing: 1},
		Results: &ast.FieldList{Opening: opening, List: results, Closing: opening},
	}
}

// Type instr
func Type(typ reflect.Type) ast.Expr {
	pkgPath, name := typ.PkgPath(), typ.Name()
	if name != "" {
		if pkgPath != "" {
			pkg := path.Base(pkgPath)
			return &ast.SelectorExpr{X: Ident(pkg), Sel: Ident(name)}
		}
		return Ident(name)
	}
	kind := typ.Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		return ArrayType(typ)
	case reflect.Map:
		return MapType(typ)
	case reflect.Ptr:
		return PtrType(typ.Elem())
	case reflect.Func:
		return FuncType(typ)
	case reflect.Interface:
		return InterfaceType(typ)
	case reflect.Chan:
		return ChanType(typ)
	case reflect.Struct:
		return StructType(typ)
	}
	log.Panicln("Type: unknown type -", typ)
	return nil
}

// Ident - ast.Ident
func Ident(name string) *ast.Ident {
	return &ast.Ident{Name: name}
}

// StringConst - ast.BasicLit
func StringConst(v string) *ast.BasicLit {
	return &ast.BasicLit{
		Kind:  token.STRING,
		Value: strconv.Quote(v),
	}
}

// IntConst - ast.BasicLit
func IntConst(v int64) *ast.BasicLit {
	return &ast.BasicLit{
		Kind:  token.INT,
		Value: strconv.FormatInt(v, 10),
	}
}
