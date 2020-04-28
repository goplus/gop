package goprj

import (
	"errors"
	"go/ast"
	"syscall"

	"github.com/qiniu/qlang/modutil"
	"github.com/qiniu/x/log"
)

var (
	// ErrSymbolIsNotAType error.
	ErrSymbolIsNotAType = errors.New("symbol isn't a type")
)

// -----------------------------------------------------------------------------

// TypeInferrer represents a TypeInferrer who can infer type from a ast.Expr.
type TypeInferrer interface {
	InferType(pkg *Package, expr ast.Expr, reserved int) (typ Type)
	InferConst(pkg *Package, expr ast.Expr, i int) (typ Type, val interface{})
}

type nilTypeInferer struct {
}

func (p *nilTypeInferer) InferType(pkg *Package, expr ast.Expr, reserved int) Type {
	return &UninferedType{expr}
}

func (p *nilTypeInferer) InferConst(pkg *Package, expr ast.Expr, i int) (typ Type, val interface{}) {
	return &UninferedType{expr}, expr
}

// -----------------------------------------------------------------------------

// Project represents a new Go project.
type Project struct {
	prjMod modutil.Module
	names  PkgNames
	types  map[string]Type
	pkgs   map[string]*Package
	TypeInferrer
}

// Open loads module from `dir` and creates a new Project object.
func Open(dir string) (*Project, error) {
	mod, err := modutil.LoadModule(dir)
	if err != nil {
		return nil, err
	}
	return &Project{
		prjMod:       mod,
		names:        NewPkgNames(),
		types:        make(map[string]Type),
		pkgs:         make(map[string]*Package),
		TypeInferrer: &nilTypeInferer{},
	}, nil
}

// ThisModule returns the Module instance.
func (p *Project) ThisModule() modutil.Module {
	return p.prjMod
}

// Load loads the main package of a Go module.
func (p *Project) Load() (pkg *Package, err error) {
	mod := p.prjMod
	return openPackage(mod.PkgPath(), mod.RootPath(), p)
}

// LoadPackage loads a package.
func (p *Project) LoadPackage(pkgPath string) (pkg *Package, err error) {
	if pkgPath == "C" {
		panic("LoadPackage: \"C\" not allow")
	}
	if pkg, ok := p.pkgs[pkgPath]; ok {
		return pkg, nil
	}
	pi, err := p.prjMod.Lookup(pkgPath)
	if err != nil {
		return
	}
	log.Info("====> LoadPackage:", pi.PkgPath)
	pkg, err = openPackage(pi.PkgPath, pi.Location, p)
	if err != nil {
		return
	}
	p.pkgs[pkgPath] = pkg
	return
}

// LookupType returns the type of symbol pkgPath.name.
func (p *Project) LookupType(pkgPath string, name string) (typ Type, err error) {
	if pkgPath == "C" {
		return p.UniqueType(&NamedType{PkgPath: "C", Name: name}), nil
	}
	pkg, err := p.LoadPackage(pkgPath)
	if err != nil {
		return
	}
	sym, ok := pkg.LookupSymbol(name)
	if !ok {
		return nil, syscall.ENOENT
	}
	tsym, ok := sym.(*TypeSym)
	if !ok {
		return nil, ErrSymbolIsNotAType
	}
	if tsym.Alias {
		return tsym.Type, nil
	}
	return p.UniqueType(&NamedType{PkgPath: pkg.PkgPath(), Name: name}), nil
}

// LookupPkgName lookups a package name by specified PkgPath.
func (p *Project) LookupPkgName(pkgPath string) string {
	return p.names.LookupPkgName(p.prjMod, pkgPath)
}

// UniqueType returns the unique instance of a type.
func (p *Project) UniqueType(t Type) Type {
	id := t.ID()
	if v, ok := p.types[id]; ok {
		return v
	}
	p.types[id] = t
	return t
}

// -----------------------------------------------------------------------------
