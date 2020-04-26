package gopkg

import (
	"go/ast"
	"go/token"
	"reflect"
	"strconv"

	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

type fileLoader struct {
	pkg     *Package
	names   *PkgNames
	types   *UniqueTypes
	imports map[string]string
}

func newFileLoader(pkg *Package, names *PkgNames, types *UniqueTypes) *fileLoader {
	return &fileLoader{pkg: pkg, names: names, types: types}
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
		name = p.names.GetPkgName(path)
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
	var name = d.Name.Name
	var recv Type
	if d.Recv != nil {
		fields := d.Recv.List
		if len(fields) != 1 {
			panic("loadFunc: multi recv object?")
		}
		field := fields[0]
		recv = p.ToType(field.Type)
	}
	typ := p.ToFuncType(d)
	p.pkg.insertFunc(name, recv, typ)
	if log.Ldebug >= log.Std.Level {
		log.Debug("func:", getFuncName(name, recv), "-", typ.Unique())
	}
}

func getFuncName(name string, recv Type) string {
	if recv == nil {
		return name
	}
	return "(" + recv.Unique() + ")." + name
}

// -----------------------------------------------------------------------------

// ToFuncType converts ast.FuncDecl to a FuncType.
func (p *fileLoader) ToFuncType(d *ast.FuncDecl) *FuncType {
	params := p.ToTypes(d.Type.Params)
	results := p.ToTypes(d.Type.Results)
	return p.types.Unique(&FuncType{Params: params, Results: results}).(*FuncType)
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
	case *ast.StarExpr:
		elem := p.ToType(v.X)
		return p.types.Unique(&PointerType{elem})
	case *ast.SelectorExpr:
		x, ok := v.X.(*ast.Ident)
		if !ok {
			log.Fatalln("ToType: SelectorExpr isn't *ast.Ident -", reflect.TypeOf(v.X))
		}
		pkgPath, ok := p.imports[x.Name]
		if !ok {
			log.Fatalln("ToType: PkgName isn't imported -", x.Name)
		}
		return p.types.Unique(&NamedType{PkgPath: pkgPath, Name: v.Sel.Name})
	case *ast.ArrayType:
		n := ToLen(v.Len)
		elem := p.ToType(v.Elt)
		return p.types.Unique(&ArrayType{Len: n, Elem: elem})
	case *ast.Ident:
		return p.IdentType(v.Name)
	}
	log.Fatalln("ToType: unknown -", reflect.TypeOf(typ))
	return nil
}

// IdentType converts an ident to a Type.
func (p *fileLoader) IdentType(ident string) Type {
	return p.types.Unique(&NamedType{Name: ident})
}

// ToLen converts ast.Expr to a Len.
func ToLen(e ast.Expr) int {
	if e != nil {
		log.Fatalln("not impl - ToLen:", reflect.TypeOf(e))
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
