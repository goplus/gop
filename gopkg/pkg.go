package gopkg

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strings"

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
		return &Package{fset: fset, impl: impl, name: name}, nil
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
	log.Debug("func:", d.Name.Name)
}

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

// -----------------------------------------------------------------------------

// Export represents a Go package's export info.
type Export struct {
}

// -----------------------------------------------------------------------------
