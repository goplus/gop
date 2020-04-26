package gopkg

import (
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
	"strings"

	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

type fileLoader struct {
	pkg     *Package
	imports map[string]string
}

func newFileLoader(pkg *Package) *fileLoader {
	return &fileLoader{pkg: pkg}
}

func (p *fileLoader) load(f *ast.File) {
	p.imports = make(map[string]string)
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

func (p *fileLoader) loadImport(spec *ast.ImportSpec) {
	var path = ToString(spec.Path)
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
		name = p.pkg.GetPkgName(path)
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
	//fmt.Println("type:", *spec)
}

func (p *fileLoader) loadTypes(d *ast.GenDecl) {
	for _, item := range d.Specs {
		p.loadType(item.(*ast.TypeSpec))
	}
}

func (p *fileLoader) loadVal(spec *ast.ValueSpec, tok token.Token) {
	//fmt.Println("tok:", tok, "val:", *spec)
}

func (p *fileLoader) loadVals(d *ast.GenDecl, tok token.Token) {
	for _, item := range d.Specs {
		p.loadVal(item.(*ast.ValueSpec), tok)
	}
}

func (p *fileLoader) loadFunc(d *ast.FuncDecl) {
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

// ToString converts a ast.BasicLit to string value.
func ToString(l *ast.BasicLit) string {
	if l.Kind == token.STRING {
		s, err := strconv.Unquote(l.Value)
		if err == nil {
			return s
		}
	}
	panic("ToString: convert ast.BasicLit to string failed")
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
