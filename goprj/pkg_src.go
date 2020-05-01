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
	pkgs, err := parser.ParseDir(fset, dir, filterTest, parser.ParseComments)
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

func (p *Package) requireSource() (err error) {
	if p.src != nil {
		return nil
	}
	p.src, err = OpenGoPackage(p.dir)
	if err != nil {
		return
	}
	p.busy = true
	p.syms = make(map[string]Symbol)
	loader := newFileLoader(p, p.prj)
	for name, f := range p.src.Files {
		log.Info("file:", name)
		loader.load(f)
	}
	p.busy = false
	return
}

func (p *Package) insertFunc(name string, recv Type, typ *FuncType) {
	if recv == nil {
		p.insertSymbol(name, typ)
	} else {
		// TODO
	}
}

func (p *Package) insertVar(name string, typ Type) {
	if _, ok := p.syms[name]; ok {
		log.Debug(p.mod.VersionPkgPath(), "insertVar failed: exists -", name)
		return // We don't support condition compile, and it's normal we get an existed var.
	}
	p.syms[name] = &VarSym{typ}
	log.Debug("var:", name, "-", typ)
}

func (p *Package) insertSymbol(name string, sym Symbol) {
	if _, ok := p.syms[name]; ok {
		log.Debug(p.mod.VersionPkgPath(), "insertSymbol failed: exists -", name)
		return // We don't support condition compile, and it's normal we get an existed symbol.
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

// AtomType represents a basic Go type.
type AtomType reflect.Kind

func (p AtomType) String() string {
	return p.ID()
}

// ID returns a unique id of this type.
func (p AtomType) ID() string {
	return reflect.Kind(p).String()
}

// Sizeof returns size of this type.
func (p AtomType) Sizeof(prj *Project) (n uintptr) {
	if p == Unbound {
		log.Fatalln("Sizeof(Unbound): don't call me")
	}
	return sizeofAtomTypes[int(p)]
}

const (
	// Unbound - unbound
	Unbound = AtomType(reflect.Invalid)
	// Bool - bool
	Bool = AtomType(reflect.Bool)
	// Int - int
	Int = AtomType(reflect.Int)
	// Int8 - int8
	Int8 = AtomType(reflect.Int8)
	// Int16 - int16
	Int16 = AtomType(reflect.Int16)
	// Int32 - int32
	Int32 = AtomType(reflect.Int32)
	// Int64 - int64
	Int64 = AtomType(reflect.Int64)
	// Uint - uint
	Uint = AtomType(reflect.Uint)
	// Uint8 - uint8
	Uint8 = AtomType(reflect.Uint8)
	// Uint16 - uint16
	Uint16 = AtomType(reflect.Uint16)
	// Uint32 - uint32
	Uint32 = AtomType(reflect.Uint32)
	// Uint64 - uint64
	Uint64 = AtomType(reflect.Uint64)
	// Uintptr - uintptr
	Uintptr = AtomType(reflect.Uintptr)
	// Float32 - float32
	Float32 = AtomType(reflect.Float32)
	// Float64 - float64
	Float64 = AtomType(reflect.Float64)
	// Complex64 - complex64
	Complex64 = AtomType(reflect.Complex64)
	// Complex128 - complex128
	Complex128 = AtomType(reflect.Complex128)
	// String - string
	String = AtomType(reflect.String)
	// UnsafePointer - unsafe.Pointer
	UnsafePointer = AtomType(reflect.UnsafePointer)
	// Rune - rune
	Rune = Int32
	// Byte - byte
	Byte = Uint8
)

var sizeofAtomTypes = [...]uintptr{
	reflect.Bool:          1,
	reflect.Int:           unsafe.Sizeof(int(0)),
	reflect.Int8:          1,
	reflect.Int16:         2,
	reflect.Int32:         4,
	reflect.Int64:         8,
	reflect.Uint:          unsafe.Sizeof(uint(0)),
	reflect.Uint8:         1,
	reflect.Uint16:        2,
	reflect.Uint32:        4,
	reflect.Uint64:        8,
	reflect.Uintptr:       unsafe.Sizeof(uintptr(0)),
	reflect.Float32:       4,
	reflect.Float64:       8,
	reflect.Complex64:     8,
	reflect.Complex128:    16,
	reflect.String:        unsafe.Sizeof(string('0')),
	reflect.UnsafePointer: unsafe.Sizeof(uintptr(0)),
}

// NamedType represents a named type.
type NamedType struct {
	VersionPkgPath string
	Name           string
}

func (p *NamedType) String() string {
	return p.ID()
}

// ID returns a unique id of this type.
func (p *NamedType) ID() string {
	if p.VersionPkgPath == "" {
		return p.Name
	}
	return p.VersionPkgPath + "." + p.Name
}

// Sizeof returns size of this type.
func (p *NamedType) Sizeof(prj *Project) uintptr {
	t := prj.FindVersionType(p.VersionPkgPath, p.Name)
	return t.Sizeof(prj)
}

// ChanType represents chan type.
type ChanType struct {
	Value Type        // value type
	Dir   ast.ChanDir // channel direction
}

// Sizeof returns size of this type.
func (p *ChanType) Sizeof(prj *Project) (n uintptr) {
	return unsafe.Sizeof(uintptr(0))
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

// Sizeof returns size of this type.
func (p *EllipsisType) Sizeof(prj *Project) (n uintptr) {
	return unsafe.Sizeof([]int{})
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
	Len  int64 // Len=0 for slice type
	Elem Type
}

// Sizeof returns size of this type.
func (p *ArrayType) Sizeof(prj *Project) (n uintptr) {
	if p.Len > 0 {
		return p.Elem.Sizeof(prj) * uintptr(p.Len)
	}
	if p.Len == 0 {
		return unsafe.Sizeof([]int{})
	}
	log.Fatalln("ArrayType.Sizeof: len < 0?")
	return 0
}

func (p *ArrayType) String() string {
	elem := p.Elem.String()
	if p.Len > 0 {
		len := strconv.FormatInt(p.Len, 10)
		return "[" + len + "]" + elem
	}
	return "[]" + elem
}

// ID returns a unique id of this type.
func (p *ArrayType) ID() string {
	elem := pointer(p.Elem)
	if p.Len > 0 {
		len := strconv.FormatInt(p.Len, 10)
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

// Sizeof returns size of this type.
func (p *PointerType) Sizeof(prj *Project) (n uintptr) {
	return unsafe.Sizeof(uintptr(0))
}

func (p *PointerType) String() string {
	return "*" + p.Elem.String()
}

// ID returns a unique id of this type.
func (p *PointerType) ID() string {
	return "*" + pointer(p.Elem)
}

// RetType represents types of function return value.
type RetType struct {
	Results []Type
}

// Sizeof returns size of this type.
func (p *RetType) Sizeof(prj *Project) uintptr {
	log.Fatalln("RetType.Sizeof: don't call me")
	return 0
}

func (p *RetType) String() string {
	return listTypeList(p.Results, false)
}

// ID returns a unique id of this type.
func (p *RetType) ID() string {
	items := make([]string, len(p.Results))
	for i, ret := range p.Results {
		items[i] = pointer(ret)
	}
	return strings.Join(items, " ")
}

// FuncType represents a function type.
type FuncType struct {
	Params  []Type
	Results []Type
}

// Sizeof returns size of this type.
func (p *FuncType) Sizeof(prj *Project) (n uintptr) {
	return unsafe.Sizeof(uintptr(0))
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

// Sizeof returns size of this type.
func (p *InterfaceType) Sizeof(prj *Project) (n uintptr) {
	return unsafe.Sizeof(uintptr(0)) * 2
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

// Sizeof returns size of this type.
func (p *StructType) Sizeof(prj *Project) (n uintptr) {
	for _, field := range p.Fields {
		n += field.Type.Sizeof(prj)
	}
	return
}

// TypeField gets field type by specified name.
func TypeField(t Type, name string) (ft Type, ok bool) {
	switch v := t.(type) {
	case *StructType:
		for _, field := range v.Fields {
			if field.Name == name {
				return field.Type, true
			}
			if field.Name == "" { // embbed
				if ft, ok = TypeField(field.Type, name); ok {
					return
				}
			}
		}
	case *InterfaceType:
		for _, field := range v.Methods {
			if field.Name == name {
				return field.Type, true
			}
		}
	}
	return
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

// Sizeof returns size of this type.
func (p *MapType) Sizeof(prj *Project) uintptr {
	return unsafe.Sizeof(uintptr(0))
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
	return "u:" + strconv.Itoa(pos) + "," + strconv.Itoa(end)
}

// Sizeof returns size of this type.
func (p *UninferedType) Sizeof(prj *Project) uintptr {
	log.Fatalln("UninferedType.Sizeof: don't call me")
	return 0
}

// UninferedRetType represents a function return types that needs to be infered.
type UninferedRetType struct {
	Fun string
	Nth int
}

func (p *UninferedRetType) String() string {
	return p.ID()
}

// ID returns a unique id of this type.
func (p *UninferedRetType) ID() string {
	return "uret:" + p.Fun + "," + strconv.Itoa(p.Nth)
}

// Sizeof returns size of this type.
func (p *UninferedRetType) Sizeof(prj *Project) uintptr {
	log.Fatalln("UninferedRetType.Sizeof: don't call me")
	return 0
}

// Type represents a Go type.
type Type interface {
	// Sizeof returns size of this type.
	Sizeof(prj *Project) uintptr
	// ID returns a unique id of this type.
	ID() string
	String() string
}

// -----------------------------------------------------------------------------
