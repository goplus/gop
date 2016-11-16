package main

import (
	"bytes"
	"go/ast"
	"go/build"
	"go/doc"
	"go/parser"
	"go/printer"
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
	fset     *token.FileSet
	files    map[string]*ast.File
	dpkg     map[string]*doc.Package
	keys     map[DocType]map[string]map[string]interface{}
	bps      map[string]*build.Package
	defctx   bool
	contexts []*build.Context
	fnskip   func(key string) bool
}

func NewPackage(pkg string, defctx bool) (*Package, error) {
	_, err := build.Import(pkg, "", 0)
	if err != nil {
		return nil, err
	}
	p := new(Package)
	p.fset = token.NewFileSet()
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

	fs := make(map[string]*ast.File)
	var afiles []*ast.File
	for _, file := range files {
		fileName := filepath.Join(bp.Dir, file)
		f, ok := p.files[fileName]
		if !ok {
			f, err = parser.ParseFile(p.fset, fileName, nil, parser.ParseComments)
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

	apkg, err := ast.NewPackage(p.fset, fs, nil, nil)

	if err != nil {
		//fmt.Println(err)
	}

	dpkg := doc.New(apkg, bp.Name, doc.AllMethods|doc.AllDecls)

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
		typ := parserTypeDecl(v.Decl)
		if typ == NilType {
			continue
		}
		updata(typ, ctxName, v.Name, v)
		//type func
		for _, f := range v.Funcs {
			updata(Factor, ctxName, f.Name, f)
		}
		//type var
		for _, f := range v.Vars {
			for _, name := range f.Names {
				updata(Var, ctxName, name, f)
			}
		}
		//type const
		for _, f := range v.Consts {
			for _, name := range f.Names {
				updata(Const, ctxName, name, f)
			}
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
			updata(Var, ctxName, name, v)
		}
	}
	return dpkg, nil
}

func parserTypeDecl(d *ast.GenDecl) DocType {
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

func (p *Package) simpleValueDeclType(d *ast.GenDecl) string {
	switch d.Tok {
	case token.VAR:
		for _, sp := range d.Specs {
			ts := sp.(*ast.ValueSpec)
			return p.nodeString(ts.Type)
		}
	}
	return ""
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

func (p *Package) FilterCommon(typs ...DocType) ([]string, map[string]interface{}) {
	var ctxsize int
	if p.defctx {
		ctxsize = 1
	} else {
		ctxsize = len(p.contexts)
	}

	var key []string
	cm := make(map[string]interface{})
	def := contextName(&build.Default)
	for _, typ := range typs {
		if m, ok := p.keys[typ]; ok {
			for k, v := range m {
				if len(v) == ctxsize && (ast.IsExported(k) || typ == Struct) {
					if p.fnskip != nil && p.fnskip(k) {
						continue
					}
					key = append(key, k)
					cm[k] = v[def]
				}
			}
			sort.Strings(key)
		}
	}

	return key, cm
}

func (p *Package) CommonCount() int {
	var typs = []DocType{Func, Factor, Const, Var, Struct}
	ks, _ := p.FilterCommon(typs...)
	return len(ks)
}

func (p *Package) nodeString(node interface{}) string {
	if node == nil {
		return ""
	}
	var b bytes.Buffer
	printer.Fprint(&b, p.fset, node)
	return b.String()
}

// check type field un exported
func (p *Package) CheckTypeFields(d *ast.GenDecl, hasUnexportedField bool, hasUnexportedFieldPtr bool, hasUnexportedInterfaceField bool) bool {
	interfaceNames, _ := p.FilterCommon(Interface)
	switch d.Tok {
	case token.TYPE:
		for _, sp := range d.Specs {
			ts := sp.(*ast.TypeSpec)
			switch t := ts.Type.(type) {
			case *ast.StructType:
				for _, field := range t.Fields.List {
					for _, name := range field.Names {
						if !ast.IsExported(name.Name) {
							if hasUnexportedField {
								return true
							}
							if hasUnexportedFieldPtr {
								typ := p.nodeString(field.Type)
								if strings.HasPrefix(typ, "*") {
									return true
								}
							}
							if hasUnexportedInterfaceField {
								typ := p.nodeString(field.Type)
								for _, iname := range interfaceNames {
									if typ == iname {
										return true
									}
								}
							}
						}
					}
				}
			default:
				break
			}
		}
	}
	return false
}

// check func param type is var,ptr,array
func (p *Package) CheckExportType(typ string) (isVar bool, isPtr bool, isArray bool) {
	fnCheck := func(f *doc.Func) {
		if f.Decl.Type.Params.List == nil {
			return
		}
		for _, field := range f.Decl.Type.Params.List {
			str := p.nodeString(field.Type)
			if str == typ {
				isVar = true
			} else if str == "*"+typ {
				isPtr = true
			} else if str == "[]"+typ {
				isArray = true
			}
		}
	}
	//check func
	_, mf := p.FilterCommon(Func)
	for _, v := range mf {
		f := v.(*doc.Func)
		fnCheck(f)
	}
	//check type
	_, mt := p.FilterCommon(Struct)
	for _, v := range mt {
		t := v.(*doc.Type)
		for _, f := range t.Funcs {
			fnCheck(f)
		}
		for _, f := range t.Methods {
			fnCheck(f)
		}
	}
	return
}
