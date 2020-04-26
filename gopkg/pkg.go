package gopkg

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
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

// GetPkgName returns package name of a module by specified directory.
func GetPkgName(dir string) (string, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, filterTest, parser.PackageClauseOnly)
	if err != nil {
		return "", err
	}
	if len(pkgs) != 1 {
		return "", ErrMultiPackages
	}
	for name := range pkgs {
		return name, nil
	}
	panic("not here")
}

// -----------------------------------------------------------------------------

// A Package node represents a set of source files collectively building a Go package.
type Package struct {
	fset *token.FileSet
	impl *ast.Package
	name string

	mod   modutil.Module
	names map[string]string // pkg => pkgName
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
		return &Package{
			fset: fset, impl: impl, name: name,
			mod: mod, names: make(map[string]string),
		}, nil
	}
	panic("not here")
}

// Name returns the package name.
func (p *Package) Name() string {
	return p.name
}

// GetPkgName returns package name of a module.
func (p *Package) GetPkgName(pkg string) string {
	if name, ok := p.names[pkg]; ok {
		return name
	}
	name, err := p.getPkgName(pkg)
	if err != nil {
		panic(err)
	}
	p.names[pkg] = name
	return name
}

func (p *Package) getPkgName(pkg string) (name string, err error) {
	pi, err := p.mod.Lookup(pkg)
	if err != nil {
		return
	}
	if pi.Type == modutil.PkgTypeStd {
		return path.Base(pkg), nil
	}
	return GetPkgName(pi.Location)
}

// Export returns a Go package's export info, for example, functions/variables/etc.
func (p *Package) Export() (exp *Export) {
	loader := newFileLoader(p)
	for name, f := range p.impl.Files {
		log.Debug("file:", name)
		loader.load(f)
	}
	exp = &Export{}
	return
}

// Export represents a Go package's export info.
type Export struct {
}

// -----------------------------------------------------------------------------
