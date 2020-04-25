package gopkg

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/qiniu/qlang/modutil"
	"github.com/qiniu/x/log"
)

var (
	// ErrMultiPackages error.
	ErrMultiPackages = errors.New("multi packages in the same dir")
)

func filterTest(fi os.FileInfo) bool {
	return !strings.HasSuffix(fi.Name(), "_test.go")
}

// -----------------------------------------------------------------------------

// A Package node represents a set of source files collectively building a Go package.
type Package struct {
	fset *token.FileSet
	impl *ast.Package
	name string
	mod  modutil.Module
}

// Load loads a Go package.
func Load(dir string) (pkg *Package, err error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, filterTest, 0)
	if err != nil {
		return
	}
	if len(pkgs) != 1 {
		return nil, ErrMultiPackages
	}
	for name, impl := range pkgs {
		mod, err := modutil.LoadModule(dir)
		if err != nil {
			return nil, err
		}
		return &Package{fset: fset, impl: impl, name: name, mod: mod}, nil
	}
	panic("not here")
}

// Name returns the package name.
func (p *Package) Name() string {
	return p.name
}

// Export returns a Go package's export info, for example, functions/variables/etc.
func (p *Package) Export() (exp *Export) {
	exp = &Export{}
	p.load()
	return
}

// -----------------------------------------------------------------------------

func (p *Package) load() {
	for _, f := range p.impl.Files {
		for _, decl := range f.Decls {
			switch d := decl.(type) {
			case *ast.GenDecl:
				switch d.Tok {
				case token.IMPORT:
					p.loadImports(d)
				case token.TYPE:
					p.loadTypes(d)
				case token.CONST, token.VAR:
					p.loadVals(d, d.Tok)
				default:
					log.Debug("tok:", d.Tok, "spec:", reflect.TypeOf(d.Specs).Elem())
				}
			case *ast.FuncDecl:
				p.loadFunc(d)
			default:
				log.Debug("gopkg.Package.load: unknown decl -", reflect.TypeOf(decl))
			}
		}
	}
}

func (p *Package) loadImport(spec *ast.ImportSpec) {
	//fmt.Println("import:", *spec)
}

func (p *Package) loadImports(d *ast.GenDecl) {
	for _, item := range d.Specs {
		p.loadImport(item.(*ast.ImportSpec))
	}
}

func (p *Package) loadType(spec *ast.TypeSpec) {
	//fmt.Println("type:", *spec)
}

func (p *Package) loadTypes(d *ast.GenDecl) {
	for _, item := range d.Specs {
		p.loadType(item.(*ast.TypeSpec))
	}
}

func (p *Package) loadVal(spec *ast.ValueSpec, tok token.Token) {
	//fmt.Println("tok:", tok, "val:", *spec)
}

func (p *Package) loadVals(d *ast.GenDecl, tok token.Token) {
	for _, item := range d.Specs {
		p.loadVal(item.(*ast.ValueSpec), tok)
	}
}

func (p *Package) loadFunc(d *ast.FuncDecl) {
	log.Debug("func:", d.Name.Name, "-", ToFuncType(d))
}

// -----------------------------------------------------------------------------

// ToFuncType converts ast.FuncDecl to a FuncType.
func ToFuncType(d *ast.FuncDecl) *FuncType {
	var recv Type
	if d.Recv != nil {
		fields := d.Recv.List
		if len(fields) != 1 {
			panic("ToFuncType: multi recv object?")
		}
		field := fields[0]
		recv = ToType(field.Type)
	}
	params := ToTypes(d.Type.Params)
	results := ToTypes(d.Type.Results)
	return &FuncType{Recv: recv, Params: params, Results: results}
}

// ToTypes converts ast.FieldList to types []Type.
func ToTypes(fields *ast.FieldList) (types []Type) {
	if fields == nil {
		return
	}
	for _, field := range fields.List {
		n := len(field.Names)
		if n == 0 {
			n = 1
		}
		typ := ToType(field.Type)
		for i := 0; i < n; i++ {
			types = append(types, typ)
		}
	}
	return
}

// ToType converts ast.Expr to a Type.
func ToType(typ ast.Expr) Type {
	switch v := typ.(type) {
	case *ast.StarExpr:
		elem := ToType(v.X)
		return &PointerType{elem}
	case *ast.SelectorExpr:
		x, ok := v.X.(*ast.Ident)
		if !ok {
			log.Fatal("ToType: SelectorExpr isn't *ast.Ident -", reflect.TypeOf(v.X))
		}
		return &NamedType{PkgName: x.Name, PkgPath: x.Name, Name: v.Sel.Name}
	case *ast.ArrayType:
		n := ToLen(v.Len)
		elem := ToType(v.Elt)
		return &ArrayType{Len: n, Elem: elem}
	case *ast.Ident:
		return IdentType(v.Name)
	}
	log.Fatal("ToType: unknown -", reflect.TypeOf(typ))
	return nil
}

// IdentType converts an ident to a Type.
func IdentType(ident string) Type {
	return &NamedType{Name: ident}
}

// ToLen converts ast.Expr to a Len.
func ToLen(e ast.Expr) int {
	if e != nil {
		log.Debug("ToLen:", reflect.TypeOf(e))
	}
	return 0
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

// Export represents a Go package's export info.
type Export struct {
}

// -----------------------------------------------------------------------------
