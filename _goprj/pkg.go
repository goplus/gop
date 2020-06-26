/*
 Copyright 2020 The GoPlus Authors (goplus.org)

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package goprj

import (
	"errors"
	"go/token"
	"syscall"

	"github.com/qiniu/goplus/modutil"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

// Package represents a set of source files collectively building a Go package.
type Package struct {
	prj *Project
	mod *modutil.Module
	dir string

	names PkgNames

	syms map[string]Symbol // cached
	src  *GoPackage
	busy bool
}

func openPackage(dir string, prj *Project) (pkg *Package, err error) {
	mod, err := modutil.LoadModule(dir)
	if err != nil {
		return
	}
	log.Info("====> OpenPackage:", mod.VersionPkgPath(), "-", dir)
	return &Package{prj: prj, mod: mod, dir: dir, names: NewPkgNames()}, nil
}

// Source returns the Go package instance.
func (p *Package) Source() *GoPackage {
	if err := p.requireSource(); err != nil {
		log.Fatal("requireSource failed:", err)
	}
	return p.src
}

// Project returns the project instance.
func (p *Package) Project() *Project {
	return p.prj
}

// ThisModule returns the module instance.
func (p *Package) ThisModule() *modutil.Module {
	return p.mod
}

// LookupPkgName lookups a package name by specified PkgPath.
func (p *Package) LookupPkgName(pkgPath string) string {
	return p.names.LookupPkgName(p.mod, pkgPath)
}

// LoadPackage loads a package.
func (p *Package) LoadPackage(pkgPath string) (pkg *Package, err error) {
	if pkgPath == "C" || pkgPath == "unsafe" {
		log.Fatalln("LoadPackage not allow -", pkgPath)
	}
	pi, err := p.mod.Lookup(pkgPath)
	if err != nil {
		return
	}
	return p.prj.OpenPackage(pi.Location)
}

// FindPackageType returns the type who is named as `pkgPath.name`.
func (p *Package) FindPackageType(pkgPath string, name string) (typ Type, err error) {
	if pkgPath == "C" {
		return p.prj.UniqueType(&NamedType{VersionPkgPath: "C", Name: name}), nil
	}
	if pkgPath == "unsafe" {
		if name == "Pointer" {
			return UnsafePointer, nil
		}
		log.Fatalln("FindPackageType: unknown - unsafe", name)
	}
	pkg, err := p.LoadPackage(pkgPath)
	if err != nil {
		return
	}
	return pkg.FindType(name)
}

// InferPackageValue returns the type of value `pkgPath.name`.
func (p *Package) InferPackageValue(pkgPath string, name string) (typ Type, err error) {
	if pkgPath == "C" {
		log.Fatalln("InferPackageValue: unsupport \"C\" -", name)
	}
	if pkgPath == "unsafe" {
		log.Fatalln("InferPackageValue: unknown - unsafe", name)
	}
	pkg, err := p.LoadPackage(pkgPath)
	if err != nil {
		return
	}
	return pkg.InferValue(name)
}

// FindType returns the type who is named as `name`.
func (p *Package) FindType(name string) (typ Type, err error) {
	sym, err := p.FindSymbol(name)
	if err != nil {
		return
	}
	tsym, ok := sym.(*TypeSym)
	if !ok {
		return nil, ErrSymbolIsNotAType
	}
	if tsym.Alias {
		return tsym.Type, nil
	}
	return p.prj.UniqueType(&NamedType{VersionPkgPath: p.mod.VersionPkgPath(), Name: name}), nil
}

// InferValue returns the type of value `name`.
func (p *Package) InferValue(name string) (typ Type, err error) {
	sym, err := p.FindSymbol(name)
	if err != nil {
		return
	}
	tsym, kind := sym.SymInfo()
	if kind == token.TYPE {
		return nil, ErrSymbolIsAType
	}
	return tsym, nil
}

// FindSymbol lookups symbol info.
func (p *Package) FindSymbol(name string) (sym Symbol, err error) {
	if err = p.requireSource(); err != nil {
		log.Fatal("LookupSymbol failed:", err)
		return
	}
	sym, ok := p.syms[name]
	if ok {
		return
	}
	if p.busy {
		return nil, ErrPackageIsLoading
	}
	return nil, ErrNotFound
}

var (
	// ErrSymbolIsNotAType error.
	ErrSymbolIsNotAType = errors.New("symbol isn't a type")
	// ErrSymbolIsAType error.
	ErrSymbolIsAType = errors.New("symbol is a type")
	// ErrPackageIsLoading error.
	ErrPackageIsLoading = errors.New("package is loading")
	// ErrNotFound error.
	ErrNotFound = syscall.ENOENT
)

// -----------------------------------------------------------------------------
