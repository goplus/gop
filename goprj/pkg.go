package goprj

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
	"strings"

	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

// GoPackage represents a set of source files collectively building a Go package.
type GoPackage struct {
	fset *token.FileSet
	impl *ast.Package
	name string
}

// LoadGoPackage loads a Go package.
func LoadGoPackage(dir string) (pkg *GoPackage, err error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, filterTest, 0)
	if err != nil {
		return
	}
	if len(pkgs) != 1 {
		return nil, ErrMultiPackages
	}
	for name, impl := range pkgs {
		return &GoPackage{fset: fset, impl: impl, name: name}, nil
	}
	panic("not here")
}

// -----------------------------------------------------------------------------

// Package represents a set of source files collectively building a Go package.
type Package struct {
	name string
	syms map[string]Symbol
}

// NewPackage creates a Go package.
func NewPackage(name string) *Package {
	return &Package{
		name: name,
		syms: make(map[string]Symbol),
	}
}

// NewPackageFrom creates a package from a Go package.
func NewPackageFrom(p *GoPackage, prj *Project) *Package {
	pkg := NewPackage(p.name)
	loader := newFileLoader(pkg, prj)
	for name, f := range p.impl.Files {
		log.Debug("file:", name)
		loader.load(f)
	}
	return pkg
}

// Name returns the package name.
func (p *Package) Name() string {
	return p.name
}

func (p *Package) insertFunc(name string, recv Type, typ *FuncType) {
	if recv == nil {
		p.insertSym(name, &FuncSym{typ})
	} else {
		// TODO
	}
}

func (p *Package) insertSym(name string, sym Symbol) {
	if _, ok := p.syms[name]; ok {
		log.Fatalln(p.name, "package insert symbol failed: exists -", name)
	}
	p.syms[name] = sym
}

// -----------------------------------------------------------------------------

// TypeSym represents a type symbol.
type TypeSym struct {
	Type  Type
	Alias bool
}

func (p *TypeSym) String() string {
	alias := ""
	if p.Alias {
		alias = "= "
	}
	return alias + p.Type.Unique()
}

// FuncSym represents a function symbol.
type FuncSym struct {
	Type *FuncType
}

func (p *FuncSym) String() string {
	return p.Type.Unique()
}

// VarSym represents a variable symbol.
type VarSym struct {
	Type Type
}

func (p *VarSym) String() string {
	return p.Type.Unique()
}

// ConstSym represents a const symbol.
type ConstSym struct {
	Type  Type
	Value interface{}
}

func (p *ConstSym) String() string {
	if p.Value == nil {
		return p.Type.Unique()
	}
	return fmt.Sprintf("%s = %v", p.Type.Unique(), p.Value)
}

// Symbol represents a Go symbol.
type Symbol interface {
	String() string
}

// -----------------------------------------------------------------------------

// AtomKind represents kind of a basic Go type.
type AtomKind = reflect.Kind

const (
	// Bool - bool
	Bool AtomKind = reflect.Bool
	// Int - int
	Int AtomKind = reflect.Int
	// Int8 - int8
	Int8 AtomKind = reflect.Int8
	// Int16 - int16
	Int16 AtomKind = reflect.Int16
	// Int32 - int32
	Int32 AtomKind = reflect.Int32
	// Int64 - int64
	Int64 AtomKind = reflect.Int64
	// Uint - uint
	Uint AtomKind = reflect.Uint
	// Uint8 - uint8
	Uint8 AtomKind = reflect.Uint8
	// Uint16 - uint16
	Uint16 AtomKind = reflect.Uint16
	// Uint32 - uint32
	Uint32 AtomKind = reflect.Uint32
	// Uint64 - uint64
	Uint64 AtomKind = reflect.Uint64
	// Uintptr - uintptr
	Uintptr AtomKind = reflect.Uintptr
	// Float32 - float32
	Float32 AtomKind = reflect.Float32
	// Float64 - float64
	Float64 AtomKind = reflect.Float64
	// Complex64 - complex64
	Complex64 AtomKind = reflect.Complex64
	// Complex128 - complex128
	Complex128 AtomKind = reflect.Complex128
	// String - string
	String AtomKind = reflect.String
)

// AtomType represents a basic Go type.
type AtomType AtomKind

// Unique returns a unique id of this type.
func (p AtomType) Unique() string {
	return AtomKind(p).String()
}

// NamedType represents a named type.
type NamedType struct {
	PkgPath string
	Name    string
}

// Unique returns a unique id of this type.
func (p *NamedType) Unique() string {
	if p.PkgPath == "" {
		return p.Name
	}
	return p.PkgPath + "." + p.Name
}

// ArrayType represents a array/slice type.
type ArrayType struct {
	Len  int // Len=0 for slice type
	Elem Type
}

// Unique returns a unique id of this type.
func (p *ArrayType) Unique() string {
	len := ""
	if p.Len > 0 {
		len = strconv.Itoa(p.Len)
	}
	return "[" + len + "]" + p.Elem.Unique()
}

// PointerType represents a pointer type.
type PointerType struct {
	Elem Type
}

// Unique returns a unique id of this type.
func (p *PointerType) Unique() string {
	return "*" + p.Elem.Unique()
}

// FuncType represents a function type.
type FuncType struct {
	Params  []Type
	Results []Type
}

// Unique returns a unique id of this type.
func (p *FuncType) Unique() string {
	params := UniqueTypeList(p.Params, true)
	results := UniqueTypeList(p.Results, false)
	return "func" + params + results
}

// UniqueTypeList returns unique id of a type list.
func UniqueTypeList(types []Type, noEmpty bool) string {
	if types == nil {
		if noEmpty {
			return "()"
		}
		return ""
	}
	items := make([]string, len(types))
	for i, typ := range types {
		items[i] = typ.Unique()
	}
	return "(" + strings.Join(items, ",") + ")"
}

// Member represents a struct member or an interface method.
type Member struct {
	Name string
	Type Type
}

// InterfaceType represents a Go interface type.
type InterfaceType struct {
	Methods []Member
}

// Unique returns a unique id of this type.
func (p *InterfaceType) Unique() string {
	items := make([]string, len(p.Methods)<<1)
	for i, method := range p.Methods {
		items[i<<1] = method.Name
		items[(i<<1)+1] = method.Type.Unique()
	}
	return "i{" + strings.Join(items, " ") + "}"
}

// UninferedType represents a type that needs to be infered.
type UninferedType struct {
	Expr ast.Expr
}

// Unique returns a unique id of this type.
func (p *UninferedType) Unique() string {
	pos := int(p.Expr.Pos())
	end := int(p.Expr.End())
	return "uninfer:" + strconv.Itoa(pos) + "," + strconv.Itoa(end)
}

// Type represents a Go type.
type Type interface {
	// Unique returns a unique id of this type.
	Unique() string
}

// -----------------------------------------------------------------------------
