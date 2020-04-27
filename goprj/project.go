package goprj

import (
	"github.com/qiniu/qlang/modutil"
)

// -----------------------------------------------------------------------------

// Project represents a new Go project.
type Project struct {
	prjMod modutil.Module
	names  PkgNames
	types  map[string]Type
}

// Open loads module from `dir` and creates a new Project object.
func Open(dir string) (*Project, error) {
	mod, err := modutil.LoadModule(dir)
	if err != nil {
		return nil, err
	}
	return &Project{
		prjMod: mod,
		names:  NewPkgNames(),
		types:  make(map[string]Type),
	}, nil
}

// LoadPackage loads a package.
func (p *Project) LoadPackage(dir string) (pkg *Package, err error) {
	gopkg, err := LoadGoPackage(dir)
	if err != nil {
		return nil, err
	}
	return NewPackageFrom(gopkg, p), nil
}

// LookupPkgName lookups a package name by specified PkgPath.
func (p *Project) LookupPkgName(pkg string) string {
	return p.names.LookupPkgName(p.prjMod, pkg)
}

// UniqueType returns the unique instance of a type.
func (p *Project) UniqueType(t Type) Type {
	id := t.Unique()
	if v, ok := p.types[id]; ok {
		return v
	}
	p.types[id] = t
	return t
}

// -----------------------------------------------------------------------------
