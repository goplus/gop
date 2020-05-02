package cl

import (
	"path"
	"reflect"

	"github.com/qiniu/qlang/ast"
	"github.com/qiniu/qlang/ast/astutil"
	"github.com/qiniu/qlang/exec"
	"github.com/qiniu/qlang/token"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

// A Compiler represents a qlang compiler.
type Compiler struct {
	out     *exec.Builder
	imports map[string]string
}

// New returns a qlang compiler instance.
func New() *Compiler {
	return &Compiler{}
}

// Compile compiles a qlang package.
func (p *Compiler) Compile(out *exec.Builder, pkg *ast.Package) (err error) {
	p.out = out
	return
}

func (p *Compiler) loadFile(f *ast.File) {
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
				log.Fatalln("tok:", d.Tok, "spec:", reflect.TypeOf(d.Specs).Elem())
			}
		default:
			log.Fatalln("gopkg.Package.load: unknown decl -", reflect.TypeOf(decl))
		}
	}
}

func (p *Compiler) loadImports(d *ast.GenDecl) {
	for _, item := range d.Specs {
		p.loadImport(item.(*ast.ImportSpec))
	}
}

func (p *Compiler) loadImport(spec *ast.ImportSpec) {
	var pkgPath = astutil.ToString(spec.Path)
	var name string
	if spec.Name != nil {
		name = spec.Name.Name
		switch name {
		case "_", ".":
			panic("not impl")
		}
	} else {
		name = path.Base(pkgPath)
	}
	p.imports[name] = pkgPath
	log.Debug("import:", name, pkgPath)
}

func (p *Compiler) loadTypes(d *ast.GenDecl) {
	for _, item := range d.Specs {
		p.loadType(item.(*ast.TypeSpec))
	}
}

func (p *Compiler) loadType(spec *ast.TypeSpec) {
}

func (p *Compiler) loadConsts(d *ast.GenDecl) {
}

func (p *Compiler) loadVars(d *ast.GenDecl) {
	for _, item := range d.Specs {
		p.loadVar(item.(*ast.ValueSpec))
	}
}

func (p *Compiler) loadVar(spec *ast.ValueSpec) {
}

func (p *Compiler) loadFunc(d *ast.FuncDecl) {
}

// -----------------------------------------------------------------------------
