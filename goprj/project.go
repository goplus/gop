package goprj

import (
	"errors"
	"go/ast"

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
	types      map[string]Type
	openedPkgs map[string]*Package // dir => Package
	TypeInferrer
}

// NewProject creates a new Project.
func NewProject() *Project {
	return &Project{
		types:        make(map[string]Type),
		openedPkgs:   make(map[string]*Package),
		TypeInferrer: &nilTypeInferer{},
	}
}

// OpenPackage open a package by specified directory.
func (p *Project) OpenPackage(dir string) (pkg *Package, err error) {
	if pkg, ok := p.openedPkgs[dir]; ok {
		return pkg, nil
	}
	log.Info("====> OpenPackage:", dir)
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
