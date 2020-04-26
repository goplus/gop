package gopkg

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"

	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

// goPackage represents a set of source files collectively building a Go package.
type goPackage struct {
	fset *token.FileSet
	impl *ast.Package
	name string
}

// loadGoPackage loads a Go package.
func loadGoPackage(dir string) (pkg *goPackage, err error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, filterTest, 0)
	if err != nil {
		return
	}
	if len(pkgs) != 1 {
		return nil, ErrMultiPackages
	}
	for name, impl := range pkgs {
		return &goPackage{fset: fset, impl: impl, name: name}, nil
	}
	panic("not here")
}

// -----------------------------------------------------------------------------

// Package represents a set of source files collectively building a Go package.
type Package struct {
	name string
	syms map[string]Symbol
}

// New creates a Go package.
func New(name string) *Package {
	return &Package{
		name: name,
		syms: make(map[string]Symbol),
	}
}

// Load loads a Go package.
func Load(dir string, names *PkgNames) (*Package, error) {
	p, err := loadGoPackage(dir)
	if err != nil {
		return nil, err
	}
	pkg := New(p.name)
	loader := newFileLoader(pkg, names)
	for name, f := range p.impl.Files {
		log.Debug("file:", name)
		loader.load(f)
	}
	return pkg, nil
}

// Name returns the package name.
func (p *Package) Name() string {
	return p.name
}

func (p *Package) insertFunc(name string, typ *FuncType) {
	if typ.Recv == nil {
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
	return alias + p.Type.String()
}

// FuncSym represents a function symbol.
type FuncSym struct {
	Type *FuncType
}

func (p *FuncSym) String() string {
	return p.Type.String()
}

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
	return fmt.Sprintf("%v = %v", p.Type, p.Value)
}

// Symbol represents a Go symbol.
type Symbol interface {
	String() string
}

// -----------------------------------------------------------------------------

// NamedType represents a named type.
type NamedType struct {
	PkgName string
	PkgPath string
	Name    string
}

func (p *NamedType) String() string {
	pkg := p.PkgName
	if pkg != "" {
		pkg += "."
	}
	return pkg + p.Name
}

// ArrayType represents a array/slice type.
type ArrayType struct {
	Len  int // Len=0 for slice type
	Elem Type
}

func (p *ArrayType) String() string {
	len := ""
	if p.Len > 0 {
		len = strconv.Itoa(p.Len)
	}
	return "[" + len + "]" + p.Elem.String()
}

// PointerType represents a pointer type.
type PointerType struct {
	Elem Type
}

func (p *PointerType) String() string {
	return "*" + p.Elem.String()
}

// FuncType represents a function type.
type FuncType struct {
	Recv    Type
	Params  []Type
	Results []Type
}

func (p *FuncType) String() string {
	recv := ""
	if p.Recv != nil {
		recv = "(" + p.Recv.String() + ")"
	}
	params := TypeListString(p.Params, true)
	results := TypeListString(p.Results, false)
	return recv + "func" + params + results
}

// TypeListString returns friendly info of a type list.
func TypeListString(types []Type, noEmpty bool) string {
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

// Type represents a Go type.
type Type interface {
	String() string
}

// -----------------------------------------------------------------------------
