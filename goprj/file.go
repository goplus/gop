package goprj

import (
	"go/ast"
	"go/token"
	"reflect"
	"sort"
	"strconv"

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
				log.Debug("tok:", d.Tok, "spec:", reflect.TypeOf(d.Specs).Elem())
			}
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
		name = p.prj.LookupPkgName(path)
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

func (p *fileLoader) loadVar(spec *ast.ValueSpec) {
	if spec.Type != nil {
		typ := p.ToType(spec.Type)
		for _, name := range spec.Names {
			p.pkg.insertSym(name.Name, &VarSym{typ})
			if log.Ldebug >= log.Std.Level {
				log.Debug("var:", name.Name, "-", typ.Unique())
			}
		}
	} else {
		for i, name := range spec.Names {
			typ := p.prj.InferType(spec.Values[i])
			p.pkg.insertSym(name.Name, &VarSym{typ})
			if log.Ldebug >= log.Std.Level {
				log.Debug("var:", name.Name, "-", typ.Unique())
			}
		}
	}
}

func (p *fileLoader) loadVars(d *ast.GenDecl) {
	for _, item := range d.Specs {
		p.loadVar(item.(*ast.ValueSpec))
	}
}

func (p *fileLoader) loadConst(spec *ast.ValueSpec, idx int, last []ast.Expr) []ast.Expr {
	var typ Type
	if spec.Type != nil {
		typ = p.ToType(spec.Type)
	}
	vals := spec.Values
	if vals == nil {
		if last == nil {
			panic("const: no value?")
		}
		vals = last
	}
	for i, name := range spec.Names {
		typInfer, val := p.prj.InferConst(vals[i], idx)
		if typ != nil {
			typInfer = typ
		}
		p.pkg.insertSym(name.Name, &ConstSym{typInfer, val})
		if log.Ldebug >= log.Std.Level {
			log.Debug("const:", name.Name, "-", typ.Unique(), "-", val)
		}
	}
	return vals
}

func (p *fileLoader) loadConsts(d *ast.GenDecl) {
	var last []ast.Expr
	for i, item := range d.Specs {
		last = p.loadConst(item.(*ast.ValueSpec), i, last)
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
		x, ok := v.X.(*ast.Ident)
		if !ok {
			log.Fatalln("ToType: SelectorExpr isn't *ast.Ident -", reflect.TypeOf(v.X))
		}
		pkgPath, ok := p.imports[x.Name]
		if !ok {
			log.Fatalln("ToType: PkgName isn't imported -", x.Name)
		}
		return p.prj.UniqueType(&NamedType{PkgPath: pkgPath, Name: v.Sel.Name})
	case *ast.StarExpr:
		elem := p.ToType(v.X)
		return p.prj.UniqueType(&PointerType{elem})
	case *ast.ArrayType:
		n := ToLen(v.Len)
		elem := p.ToType(v.Elt)
		return p.prj.UniqueType(&ArrayType{Len: n, Elem: elem})
	case *ast.InterfaceType:
		return p.InterfaceType(v)
	}
	log.Fatalln("ToType: unknown -", reflect.TypeOf(typ))
	return nil
}

func (p *fileLoader) InterfaceType(v *ast.InterfaceType) Type {
	methods := v.Methods.List
	out := make([]Member, 0, len(methods))
	for _, field := range v.Methods.List {
		if field.Names == nil {
			panic("todo: embbed member")
		}
		typ := p.ToType(field.Type)
		for _, name := range field.Names {
			out = append(out, Member{name.Name, typ})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return p.prj.UniqueType(&InterfaceType{out})
}

// IdentType converts an ident to a Type.
func (p *fileLoader) IdentType(ident string) Type {
	return p.prj.UniqueType(&NamedType{Name: ident})
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
