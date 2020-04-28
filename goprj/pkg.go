package goprj

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

// GoPackage represents a set of source files collectively building a Go package.
type GoPackage struct {
	*ast.Package
	FileSet *token.FileSet
	Name    string
}

// OpenGoPackage opens a Go package from a directory.
func OpenGoPackage(dir string) (pkg *GoPackage, err error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, filterTest, 0)
	if err != nil {
		return
	}
	if len(pkgs) != 1 {
		delete(pkgs, "main")
		if len(pkgs) != 1 {
			log.Debug("GetPkgName failed:", ErrMultiPackages, "-", pkgs)
			return nil, ErrMultiPackages
		}
	}
	for name, impl := range pkgs {
		return &GoPackage{FileSet: fset, Package: impl, Name: name}, nil
	}
	panic("not here")
}

// -----------------------------------------------------------------------------

// Package represents a set of source files collectively building a Go package.
type Package struct {
	syms    map[string]Symbol
	src     *GoPackage
	pkgPath string
}

func openPackage(pkgPath string, dir string, prj *Project) (*Package, error) {
	gopkg, err := OpenGoPackage(dir)
	if err != nil {
		return nil, err
	}
	return newPackageFrom(pkgPath, gopkg, prj), nil
}

// newPackageFrom creates a package from a Go package.
func newPackageFrom(pkgPath string, gopkg *GoPackage, prj *Project) *Package {
	pkg := &Package{
		syms:    make(map[string]Symbol),
		src:     gopkg,
		pkgPath: pkgPath,
	}
	loader := newFileLoader(pkg, prj)
	for name, f := range gopkg.Files {
		log.Debug("file:", name)
		loader.load(f)
	}
	return pkg
}

// Source returns the Go package instance.
func (p *Package) Source() *GoPackage {
	return p.src
}

// PkgPath retuns PkgPath (with version).
func (p *Package) PkgPath() string {
	return p.pkgPath
}

// LookupSymbol lookups symbol info.
func (p *Package) LookupSymbol(name string) (sym Symbol, ok bool) {
	sym, ok = p.syms[name]
	return
}

func (p *Package) insertFunc(name string, recv Type, typ *FuncType) {
	if recv == nil {
		p.insertSym(name, typ)
	} else {
		// TODO
	}
}

func (p *Package) insertSym(name string, sym Symbol) {
	if _, ok := p.syms[name]; ok {
		log.Fatalln(p.pkgPath, "package insert symbol failed: exists -", name)
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
	return alias + p.Type.String()
}

// FuncSym represents a function symbol.
type FuncSym = FuncType

// VarSym represents a variable symbol.
type VarSym struct {
	Type Type
}

func (p *VarSym) String() string {
	return p.Type.String()
}

// ConstSym represents a const symbol.
type ConstSym struct {
	Type  Type
	Value interface{}
}

func (p *ConstSym) String() string {
	if p.Value == nil {
		return p.Type.String()
	}
	return fmt.Sprintf("%s = %v", p.Type.String(), p.Value)
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
	// Rune - rune
	Rune = Int32
	// Byte - byte
	Byte = Uint8
)

// AtomType represents a basic Go type.
type AtomType AtomKind

func (p AtomType) String() string {
	return p.ID()
}

// ID returns a unique id of this type.
func (p AtomType) ID() string {
	return AtomKind(p).String()
}

// NamedType represents a named type.
type NamedType struct {
	PkgPath string
	Name    string
}

func (p *NamedType) String() string {
	return p.ID()
}

// ID returns a unique id of this type.
func (p *NamedType) ID() string {
	if p.PkgPath == "" {
		return p.Name
	}
	return p.PkgPath + "." + p.Name
}

// ChanType represents chan type.
type ChanType struct {
	Value Type        // value type
	Dir   ast.ChanDir // channel direction
}

func (p *ChanType) String() string {
	val := p.Value.String()
	switch p.Dir {
	case ast.SEND:
		return "chan<- " + val
	case ast.RECV:
		return "<-chan " + val
	default:
		return "chan " + val
	}
}

// ID returns a unique id of this type.
func (p *ChanType) ID() string {
	val := pointer(p.Value)
	switch p.Dir {
	case ast.SEND:
		return "chan<- " + val
	case ast.RECV:
		return "<-chan " + val
	default:
		return "chan " + val
	}
}

// EllipsisType represents ...type
type EllipsisType struct {
	Elem Type
}

func (p *EllipsisType) String() string {
	return "..." + p.Elem.String()
}

// ID returns a unique id of this type.
func (p *EllipsisType) ID() string {
	return "..." + pointer(p.Elem)
}

// ArrayType represents a array/slice type.
type ArrayType struct {
	Len  int // Len=0 for slice type
	Elem Type
}

func (p *ArrayType) String() string {
	elem := p.Elem.String()
	if p.Len > 0 {
		len := strconv.Itoa(p.Len)
		return "[" + len + "]" + elem
	}
	return "[]" + elem
}

// ID returns a unique id of this type.
func (p *ArrayType) ID() string {
	elem := pointer(p.Elem)
	if p.Len > 0 {
		len := strconv.Itoa(p.Len)
		return "[" + len + "]" + elem
	}
	return "[]" + elem
}

type interfaceStruct struct {
	itabOrType uintptr
	word       uintptr
}

func pointer(typ Type) string {
	return strconv.FormatInt(int64((*interfaceStruct)(unsafe.Pointer(&typ)).word), 32)
}

// PointerType represents a pointer type.
type PointerType struct {
	Elem Type
}

func (p *PointerType) String() string {
	return "*" + p.Elem.String()
}

// ID returns a unique id of this type.
func (p *PointerType) ID() string {
	return "*" + pointer(p.Elem)
}

// FuncType represents a function type.
type FuncType struct {
	Params  []Type
	Results []Type
}

func (p *FuncType) String() string {
	params := listTypeList(p.Params, true)
	results := listTypeList(p.Results, false)
	return "func" + params + results
}

// ID returns a unique id of this type.
func (p *FuncType) ID() string {
	items := make([]string, 1, len(p.Params)+len(p.Results)+3)
	items[0] = "f{"
	for _, param := range p.Params {
		items = append(items, pointer(param))
	}
	items = append(items, ":")
	for _, ret := range p.Results {
		items = append(items, pointer(ret))
	}
	items = append(items, "}")
	return strings.Join(items, " ")
}

func listTypeList(types []Type, noEmpty bool) string {
	if types == nil {
		if noEmpty {
			return "()"
		}
		return ""
	}
	items := make([]string, len(types))
	for i, typ := range types {
		items[i] = typ.String()
	}
	return "(" + strings.Join(items, ",") + ")"
}

// Field represents a struct field or an interface method.
type Field struct {
	Name string // empty if embbed
	Type Type
}

// InterfaceType represents a Go interface type.
type InterfaceType struct {
	Methods []Field
}

func (p *InterfaceType) String() string {
	return listMembers("i{", p.Methods)
}

// ID returns a unique id of this type.
func (p *InterfaceType) ID() string {
	return uniqueMembers("i{", p.Methods)
}

func listMembers(typeTag string, fields []Field) string {
	items := make([]string, len(fields)<<1)
	for i, method := range fields {
		items[i<<1] = method.Name
		items[(i<<1)+1] = method.Type.String()
	}
	return typeTag + strings.Join(items, " ") + "}"
}

func uniqueMembers(typeTag string, fields []Field) string {
	n := len(fields) << 1
	items := make([]string, n+2)
	items[0] = typeTag
	for i, method := range fields {
		items[(i<<1)+1] = method.Name
		items[(i<<1)+2] = pointer(method.Type)
	}
	items[n+1] = "}"
	return strings.Join(items, "|")
}

// StructType represents a Go struct type.
type StructType struct {
	Fields []Field
}

func (p *StructType) String() string {
	return listMembers("s{", p.Fields)
}

// ID returns a unique id of this type.
func (p *StructType) ID() string {
	return uniqueMembers("s{", p.Fields)
}

// MapType represents a Go map type.
type MapType struct {
	Key   Type
	Value Type
}

func (p *MapType) String() string {
	return "map[" + p.Key.String() + "]" + p.Value.String()
}

// ID returns a unique id of this type.
func (p *MapType) ID() string {
	return "map[" + pointer(p.Key) + "]" + pointer(p.Value)
}

// UninferedType represents a type that needs to be infered.
type UninferedType struct {
	Expr ast.Expr
}

func (p *UninferedType) String() string {
	return p.ID()
}

// ID returns a unique id of this type.
func (p *UninferedType) ID() string {
	pos := int(p.Expr.Pos())
	end := int(p.Expr.End())
	return "uninfer:" + strconv.Itoa(pos) + "," + strconv.Itoa(end)
}

// Type represents a Go type.
type Type interface {
	// ID returns a unique id of this type.
	ID() string
	String() string
}

// -----------------------------------------------------------------------------
