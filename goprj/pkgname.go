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

// LookupPkgName lookups a package name by specified PkgPath.
func LookupPkgName(prjMod modutil.Module, pkg string) (name string, err error) {
	pi, err := prjMod.Lookup(pkg)
	if err != nil {
		return
	}
	if pi.Type == modutil.PkgTypeStd {
		return path.Base(pkg), nil
	}
	return GetPkgName(pi.Location)
}

// -----------------------------------------------------------------------------

// PkgNames keeps a package name cache: pkg => pkgName
type PkgNames struct {
	names map[string]string // pkg => pkgName
}

// NewPkgNames creates a new PkgNames object.
func NewPkgNames() PkgNames {
	return PkgNames{
		names: make(map[string]string),
	}
}

// LookupPkgName lookups a package name by specified PkgPath.
func (p PkgNames) LookupPkgName(prjMod modutil.Module, pkg string) string {
	if name, ok := p.names[pkg]; ok {
		return name
	}
	name, err := LookupPkgName(prjMod, pkg)
	if err != nil {
		panic(err)
	}
	p.names[pkg] = name
	return name
}

// -----------------------------------------------------------------------------
