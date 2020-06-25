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

package goprj

import (
	"go/ast"
	"go/token"
	"reflect"
	"sort"

	"github.com/qiniu/goplus/ast/astutil"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

type fileLoader struct {
	pkg     *Package
	prj     *Project
	imports map[string]string
}

func newFileLoader(pkg *Package, prj *Project) *fileLoader {
	return &fileLoader{pkg: pkg, prj: prj}
}

func (p *fileLoader) load(f *ast.File) {
	p.imports = make(map[string]string)
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			p.loadFunc(d)
		case *ast.GenDecl:
			switch d.Tok {
			case token.IMPORT:
				p.loadImports(d)
			case token.TYPE:
				p.loadTypes(d)
			case token.CONST:
				p.loadConsts(d)
			case token.VAR:
				p.loadVars(d)
			default:
				log.Fatalln("tok:", d.Tok, "spec:", reflect.TypeOf(d.Specs).Elem())
			}
		default:
			log.Fatalln("gopkg.Package.load: unknown decl -", reflect.TypeOf(decl))
		}
	}
}

func (p *fileLoader) loadImport(spec *ast.ImportSpec) {
	var path = astutil.ToString(spec.Path)
	var name string
	if spec.Name != nil {
		name = spec.Name.Name
		switch name {
		case "_":
			return
		case ".":
			panic("not impl")
		}
	} else {
		name = p.pkg.LookupPkgName(path)
	}
	p.imports[name] = path
	log.Debug("import:", name, path)
}

func (p *fileLoader) loadImports(d *ast.GenDecl) {
	for _, item := range d.Specs {
		p.loadImport(item.(*ast.ImportSpec))
	}
}

func (p *fileLoader) loadType(spec *ast.TypeSpec) {
	name := spec.Name.Name
	alias := spec.Assign.IsValid()
	typ := p.ToType(spec.Type)
	p.pkg.insertSymbol(name, &TypeSym{typ, alias})
	if log.Ldebug >= log.Std.Level {
		log.Debug("type:", name, typ, "alias:", alias)
	}
}

func (p *fileLoader) loadTypes(d *ast.GenDecl) {
	for _, item := range d.Specs {
		p.loadType(item.(*ast.TypeSpec))
	}
}

func (p *fileLoader) loadVar(spec *ast.ValueSpec) {
	if spec.Type != nil {
		typ := p.ToType(spec.Type)
		for _, name := range spec.Names {
			if name.Name == "_" {
				continue
			}
			p.pkg.insertVar(name.Name, typ)
		}
		return
	}
	if len(spec.Names) == len(spec.Values) {
		for i, name := range spec.Names {
			if name.Name == "_" {
				continue
			}
			typ := p.InferType(spec.Values[i])
			p.pkg.insertVar(name.Name, typ)
		}
		return
	}
	if len(spec.Values) != 1 {
		log.Fatalln("loadVar: unexpected -", *spec)
	}
	typ := p.InferType(spec.Values[0])
	switch ret := typ.(type) {
	case *RetType:
		for i, name := range spec.Names {
			if name.Name == "_" {
				continue
			}
			typ := ret.Results[i]
			p.pkg.insertVar(name.Name, typ)
		}
		return
	case *UninferedRetType:
		for i, name := range spec.Names {
			if name.Name == "_" {
				continue
			}
			typ := &UninferedRetType{Fun: ret.Fun, Nth: i}
			p.pkg.insertVar(name.Name, typ)
		}
		return
	default:
		log.Fatalln("loadVar: InferType return unknown type -", reflect.TypeOf(typ))
	}
}

func (p *fileLoader) loadVars(d *ast.GenDecl) {
	for _, item := range d.Specs {
		p.loadVar(item.(*ast.ValueSpec))
	}
}

func (p *fileLoader) loadConst(spec *ast.ValueSpec, idx int64, tlast *Type, last *[]ast.Expr) {
	var typ Type
	if spec.Type == nil {
		typ = *tlast
	} else {
		typ = p.ToType(spec.Type)
		*tlast = typ
	}
	vals := spec.Values
	if vals == nil {
		vals = *last
		if vals == nil {
			panic("const: no value?")
		}
	} else {
		*last = vals
	}
	for i, name := range spec.Names {
		typInfer, val := p.ToConst(vals[i], idx)
		if typ != nil {
			if false {
				t, ok := checkType(typ.(AtomType), typInfer.(AtomType))
				if ok {
					assertNotOverflow(t, val)
				} else {
					log.Fatalln("loadConst: checkType failed:", typ, typInfer)
				}
				typInfer = t
			}
			typInfer = typ
		}
		p.pkg.insertSymbol(name.Name, &ConstSym{typInfer, val})
		if log.Ldebug >= log.Std.Level {
			log.Debug("const:", name.Name, "-", typInfer, "-", val)
		}
	}
}

func (p *fileLoader) loadConsts(d *ast.GenDecl) {
	var tlast Type
	var last []ast.Expr
	for i, item := range d.Specs {
		p.loadConst(item.(*ast.ValueSpec), int64(i), &tlast, &last)
	}
}

func (p *fileLoader) loadFunc(d *ast.FuncDecl) {
	var name = d.Name.Name
	var recv Type
	if d.Recv != nil {
		fields := d.Recv.List
		if len(fields) != 1 {
			panic("loadFunc: multi recv object?")
		}
		field := fields[0]
		recv = p.ToType(field.Type)
	} else if name == "init" { // skip init function
		return
	}
	typ := p.ToFuncType(d.Type)
	p.pkg.insertFunc(name, recv, typ)
	if log.Ldebug >= log.Std.Level {
		log.Debug("func:", getFuncName(name, recv), "-", typ)
	}
}

func getFuncName(name string, recv Type) string {
	if recv == nil {
		return name
	}
	return "(" + recv.String() + ")." + name
}

// -----------------------------------------------------------------------------

// ToFuncType converts ast.FuncDecl to a FuncType.
func (p *fileLoader) ToFuncType(t *ast.FuncType) *FuncType {
	params := p.ToTypes(t.Params)
	results := p.ToTypes(t.Results)
	return p.prj.UniqueType(&FuncType{Params: params, Results: results}).(*FuncType)
}

// ToTypes converts ast.FieldList to types []Type.
func (p *fileLoader) ToTypes(fields *ast.FieldList) (types []Type) {
	if fields == nil {
		return
	}
	for _, field := range fields.List {
		n := len(field.Names)
		if n == 0 {
			n = 1
		}
		typ := p.ToType(field.Type)
		for i := 0; i < n; i++ {
			types = append(types, typ)
		}
	}
	return
}

// ToType converts ast.Expr to a Type.
func (p *fileLoader) ToType(typ ast.Expr) Type {
	switch v := typ.(type) {
	case *ast.Ident:
		return p.IdentType(v.Name)
	case *ast.SelectorExpr:
		return p.ExternalType(v)
	case *ast.StarExpr:
		elem := p.ToType(v.X)
		return p.prj.UniqueType(&PointerType{elem})
	case *ast.ArrayType:
		n := p.ToLen(v.Len)
		elem := p.ToType(v.Elt)
		if n < 0 {
			return &EllipsisType{elem}
		}
		return p.prj.UniqueType(&ArrayType{Len: n, Elem: elem})
	case *ast.FuncType:
		return p.ToFuncType(v)
	case *ast.InterfaceType:
		return p.InterfaceType(v)
	case *ast.StructType:
		return p.StructType(v)
	case *ast.MapType:
		key := p.ToType(v.Key)
		val := p.ToType(v.Value)
		return p.prj.UniqueType(&MapType{key, val})
	case *ast.ChanType:
		val := p.ToType(v.Value)
		return p.prj.UniqueType(&ChanType{val, v.Dir})
	case *ast.Ellipsis:
		elem := p.ToType(v.Elt)
		return p.prj.UniqueType(&EllipsisType{elem})
	}
	log.Fatalln("ToType: unknown -", reflect.TypeOf(typ))
	return nil
}

func (p *fileLoader) StructType(v *ast.StructType) Type {
	fields := v.Fields.List
	out := make([]Field, 0, len(fields))
	for _, field := range fields {
		typ := p.ToType(field.Type)
		if field.Names == nil { // embbed
			out = append(out, Field{"", typ})
			continue
		}
		for _, name := range field.Names {
			out = append(out, Field{name.Name, typ})
		}
	}
	return p.prj.UniqueType(&StructType{out})
}

func (p *fileLoader) InterfaceType(v *ast.InterfaceType) Type {
	methods := v.Methods.List
	out := make([]Field, 0, len(methods))
	for _, field := range methods {
		typ := p.ToType(field.Type)
		if field.Names == nil { // embbed
			out = append(out, Field{"", typ})
		}
		for _, name := range field.Names {
			out = append(out, Field{name.Name, typ})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return p.prj.UniqueType(&InterfaceType{out})
}

func (p *fileLoader) ExternalType(v *ast.SelectorExpr) Type {
	x, ok := v.X.(*ast.Ident)
	if !ok {
		log.Fatalln("ExternalType: SelectorExpr isn't *ast.Ident -", reflect.TypeOf(v.X))
	}
	pkgPath, ok := p.imports[x.Name]
	if !ok {
		log.Fatalln("ExternalType: PkgName isn't imported -", x.Name)
	}
	typ, err := p.pkg.FindPackageType(pkgPath, v.Sel.Name)
	if err != nil {
		log.Fatalln("ExternalType - FindPackageType failed:", err, "-", pkgPath, v.Sel.Name)
	}
	return typ
}

// IdentType converts an ident to a Type.
func (p *fileLoader) IdentType(ident string) Type {
	if t, ok := builtinTypes[ident]; ok {
		return t
	}
	return p.prj.UniqueType(&NamedType{Name: ident})
}

var builtinTypes = map[string]AtomType{
	"bool":       Bool,
	"int":        Int,
	"int8":       Int8,
	"int16":      Int16,
	"int32":      Int32,
	"int64":      Int64,
	"uint":       Uint,
	"uint8":      Uint8,
	"uint16":     Uint16,
	"uint32":     Uint32,
	"uint64":     Uint64,
	"uintptr":    Uintptr,
	"float32":    Float32,
	"float64":    Float64,
	"complex64":  Complex64,
	"complex128": Complex128,
	"string":     String,
	"rune":       Rune,
	"byte":       Byte,
}

// -----------------------------------------------------------------------------
