package goprj

// -----------------------------------------------------------------------------

// Project represents a new Go project.
type Project struct {
	types      map[string]Type
	openedPkgs map[string]*Package // dir => Package
}

// NewProject creates a new Project.
func NewProject() *Project {
	return &Project{
		types:      make(map[string]Type),
		openedPkgs: make(map[string]*Package),
	}
}

// OpenPackage open a package by specified directory.
func (p *Project) OpenPackage(dir string) (pkg *Package, err error) {
	if pkg, ok := p.openedPkgs[dir]; ok {
		return pkg, nil
	}
	pkg, err = openPackage(dir, p)
	if err != nil {
		return
	}
	p.openedPkgs[dir] = pkg
	return
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
