package goprj

import (
	"errors"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"

	"github.com/qiniu/qlang/modutil"
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

// PkgNames keeps a package name cache: pkg => pkgName
type PkgNames struct {
	mod   modutil.Module
	names map[string]string // pkg => pkgName
}

// OpenPkgNames loads module from `dir` and creates a new PkgNames object.
func OpenPkgNames(dir string) (*PkgNames, error) {
	mod, err := modutil.LoadModule(dir)
	if err != nil {
		return nil, err
	}
	return NewPkgNames(mod), nil
}

// NewPkgNames creates a new PkgNames object.
func NewPkgNames(mod modutil.Module) *PkgNames {
	return &PkgNames{
		mod:   mod,
		names: make(map[string]string),
	}
}

// GetPkgName returns package name of a module.
func (p *PkgNames) GetPkgName(pkg string) string {
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

func (p *PkgNames) getPkgName(pkg string) (name string, err error) {
	pi, err := p.mod.Lookup(pkg)
	if err != nil {
		return
	}
	if pi.Type == modutil.PkgTypeStd {
		return path.Base(pkg), nil
	}
	return GetPkgName(pi.Location)
}

// -----------------------------------------------------------------------------
