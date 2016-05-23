package main

import (
	"go/ast"
	"go/build"
	"go/doc"
	"go/importer"
	"go/parser"
	"go/token"
	"path/filepath"
	"sort"
	"strings"
)

type DocType int

const (
	NilType DocType = iota
	Interface
	Struct
	Type
	Var
	Const
	Factor
	Func
)

type Package struct {
	pkg      string
	files    map[string]*ast.File
	dpkg     map[string]*doc.Package
	keys     map[DocType]map[string]map[string]interface{}
	bps      map[string]*build.Package
	defctx   bool
	contexts []*build.Context
}

func NewPackage(pkg string, defctx bool) (*Package, error) {
	_, err := build.Import(pkg, "", 0)
	if err != nil {
		return nil, err
	}
	p := new(Package)
	p.pkg = pkg
	p.files = make(map[string]*ast.File)
	p.dpkg = make(map[string]*doc.Package)
	p.keys = make(map[DocType]map[string]map[string]interface{})
	p.bps = make(map[string]*build.Package)
	p.defctx = defctx
	p.contexts = contexts
	return p, nil
}

func (p *Package) BuildPackage() *build.Package {
	for _, v := range p.bps {
		if v != nil {
			return v
		}
	}
	return nil
}

func poorMansImporter(imports map[string]*ast.Object, path string) (*ast.Object, error) {
	pkg := imports[path]
	if pkg == nil {
		//fmt.Println("import", path)
		importer.Default()
		// note that strings.LastIndex returns -1 if there is no "/"
		pkg = ast.NewObj(ast.Pkg, path[strings.LastIndex(path, "/")+1:])
		pkg.Data = ast.NewScope(nil) // required by ast.NewPackage for dot-import
		imports[path] = pkg
	}
	return pkg, nil
}

func (p *Package) Parser() {
	if p.defctx {
		p.parser(build.Default)
	} else {
		p.parserAllContext()
	}
}

func (p *Package) parserAllContext() {
	for _, c := range p.contexts {
		ctx := build.Default
		ctx.GOOS = c.GOOS
		ctx.GOARCH = c.GOARCH
		ctx.CgoEnabled = c.CgoEnabled
		p.parser(ctx)
	}
}

func (p *Package) parser(ctx build.Context) (*doc.Package, error) {
	bp, err := ctx.Import(p.pkg, "", 0)
	if err != nil {
		return nil, err
	}
	ctxName := contextName(&ctx)
	p.bps[ctxName] = bp
	var files []string
	files = append(files, bp.GoFiles...)
	files = append(files, bp.CgoFiles...)
	fset := token.NewFileSet()
	fs := make(map[string]*ast.File)
	var afiles []*ast.File
	for _, file := range files {
		fileName := filepath.Join(bp.Dir, file)
		f, ok := p.files[fileName]
		if !ok {
			f, err = parser.ParseFile(fset, fileName, nil, parser.ParseComments)
			if err != nil {
				return nil, err
			}
			p.files[fileName] = f
		}
		pkgName := f.Name.Name
		if pkgName == bp.Name {
			fs[fileName] = f
		}
		afiles = append(afiles, f)
	}

	apkg, err := ast.NewPackage(fset, fs, nil, nil)
	if err != nil {
		//fmt.Println(err)
	}
	dpkg := doc.New(apkg, bp.Name, doc.AllMethods)
	p.dpkg[contextName(&ctx)] = dpkg

	updata := func(typ DocType, ctx string, key string, value interface{}) {
		if p.keys[typ] == nil {
			p.keys[typ] = make(map[string]map[string]interface{})
		}
		if p.keys[typ][key] == nil {
			p.keys[typ][key] = make(map[string]interface{})
		}
		p.keys[typ][key][ctx] = value
	}

	//types
	for _, v := range dpkg.Types {
		typ := parserDeclType(v.Decl)
		if typ == NilType {
			continue
		}
		updata(typ, ctxName, v.Name, v)
		for _, f := range v.Funcs {
			updata(Factor, ctxName, f.Name, f)
		}
	}
	//funcs
	for _, v := range dpkg.Funcs {
		updata(Func, ctxName, v.Name, v)
	}
	//consts
	for _, v := range dpkg.Consts {
		for _, name := range v.Names {
			updata(Const, ctxName, name, v)
		}
	}
	//vars
	for _, v := range dpkg.Vars {
		for _, name := range v.Names {
			updata(Const, ctxName, name, v)
		}
	}
	return dpkg, nil
}

func parserDeclType(d *ast.GenDecl) DocType {
	switch d.Tok {
	case token.TYPE:
		for _, sp := range d.Specs {
			ts := sp.(*ast.TypeSpec)
			switch ts.Type.(type) {
			case *ast.InterfaceType:
				return Interface
			case *ast.StructType:
				return Struct
			default:
				return Type
			}
		}
	}
	return NilType
}

func (p *Package) Filter(typ DocType) ([]string, map[string]map[string]interface{}) {
	rm := make(map[string]map[string]interface{})
	if m, ok := p.keys[typ]; ok {
		var key []string
		for k, v := range m {
			if ast.IsExported(k) {
				key = append(key, k)
				rm[k] = v
			}
		}
		sort.Strings(key)
		return key, rm
	}
	return nil, nil
}

func (p *Package) FilterCommon(typ DocType) ([]string, map[string]interface{}) {
	var ctxsize int
	if p.defctx {
		ctxsize = 1
	} else {
		ctxsize = len(p.contexts)
	}
	if m, ok := p.keys[typ]; ok {
		cm := make(map[string]interface{})
		def := contextName(&build.Default)
		var key []string
		for k, v := range m {
			if ast.IsExported(k) && len(v) == ctxsize {
				key = append(key, k)
				cm[k] = v[def]
			}
		}
		sort.Strings(key)
		return key, cm
	}
	return nil, nil
}

func (p *Package) CommonCount() int {
	var count int
	var typs = []DocType{Func, Factor, Const, Var, Struct}
	for _, typ := range typs {
		ks, _ := p.FilterCommon(typ)
		count += len(ks)
	}
	return count
}
