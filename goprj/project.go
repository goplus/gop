package goprj

import (
	"go/ast"

	"github.com/qiniu/qlang/modutil"
)

// -----------------------------------------------------------------------------

// TypeInferrer represents a TypeInferrer who can infer type from a ast.Expr.
type TypeInferrer interface {
	InferType(expr ast.Expr) (typ Type)
	InferConst(expr ast.Expr, i int) (typ Type, val interface{})
}

type nilTypeInferer struct {
}

func (p *nilTypeInferer) InferType(expr ast.Expr) Type {
	return &UninferedType{expr}
}

func (p *nilTypeInferer) InferConst(expr ast.Expr, i int) (typ Type, val interface{}) {
	return &UninferedType{expr}, expr
}

// -----------------------------------------------------------------------------

// Project represents a new Go project.
type Project struct {
	prjMod modutil.Module
	names  PkgNames
	types  map[string]Type
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
	pi, err := p.prjMod.Lookup(pkgPath)
	if err != nil {
		return
	}
	return openPackage(pi.PkgPath, pi.Location, p)
}

// LookupPkgName lookups a package name by specified PkgPath.
func (p *Project) LookupPkgName(pkg string) string {
	return p.names.LookupPkgName(p.prjMod, pkg)
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
