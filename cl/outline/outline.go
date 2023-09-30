/*
 * Copyright (c) 2023 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package outline

import (
	"go/types"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/token"
	"github.com/goplus/gox"
	"github.com/goplus/mod/modfile"
)

// -----------------------------------------------------------------------------

type Project = modfile.Project

type Config struct {
	// Fset provides source position information for syntax trees and types.
	// If Fset is nil, Load will use a new fileset, but preserve Fset's value.
	Fset *token.FileSet

	// WorkingDir is the directory in which to run gop compiler.
	WorkingDir string

	// C2goBase specifies base of standard c2go packages.
	// Default is github.com/goplus/.
	C2goBase string

	// LookupPub lookups the c2go package pubfile (named c2go.a.pub).
	LookupPub func(pkgPath string) (pubfile string, err error)

	// LookupClass lookups a class by specified file extension.
	LookupClass func(ext string) (c *Project, ok bool)

	// An Importer resolves import paths to Packages.
	Importer types.Importer
}

type Package struct {
	objs []types.Object
	Pkg  *types.Package
}

// NewPackage creates a Go/Go+ outline package.
func NewPackage(pkgPath string, pkg *ast.Package, conf *Config) (_ Package, err error) {
	ret, err := cl.NewPackage(pkgPath, pkg, &cl.Config{
		Fset:           conf.Fset,
		WorkingDir:     conf.WorkingDir,
		C2goBase:       conf.C2goBase,
		LookupPub:      conf.LookupPub,
		LookupClass:    conf.LookupClass,
		Importer:       conf.Importer,
		NoFileLine:     true,
		NoAutoGenMain:  true,
		NoSkipConstant: true,
		Outline:        true,
	})
	if err != nil {
		return
	}
	return From(ret), nil
}

func From(pkg *gox.Package) Package {
	scope := pkg.Types.Scope()
	names := scope.Names()
	objs := make([]types.Object, len(names))
	for i, name := range names {
		objs[i] = scope.Lookup(name)
	}
	return Package{objs, pkg.Types}
}

func (p Package) Valid() bool {
	return p.Pkg != nil
}

// -----------------------------------------------------------------------------

type All struct {
	Consts []Const
	Vars   []Var
	Funcs  []Func
	Types  []*TypeName

	named map[*types.TypeName]*TypeName
}

func (p *All) initNamed(objs []types.Object, all bool) {
	for _, o := range objs {
		if t, ok := o.(*types.TypeName); ok {
			named := &TypeName{TypeName: t}
			p.named[t] = named
			if all || o.Exported() {
				p.Types = append(p.Types, named)
			}
		}
	}
}

func (p Package) Outline(withUnexported ...bool) (ret *All) {
	ret = &All{
		named: make(map[*types.TypeName]*TypeName),
	}
	all := (withUnexported != nil && withUnexported[0])
	ret.initNamed(p.objs, all)
	for _, o := range p.objs {
		if _, ok := o.(*types.TypeName); ok || !(all || o.Exported()) {
			continue
		}
		switch v := o.(type) {
		case *gox.Func:
			sig := v.Type().(*types.Signature)
			kind, t := sigKind(p.Pkg, sig)
			switch kind {
			case sigNormal:
				ret.Funcs = append(ret.Funcs, Func{v})
			case sigCreator:
				if named, ok := ret.named[t.Obj()]; ok {
					named.Creators = append(named.Creators, Func{v})
				}
			}
		case *types.Const:
			if t := checkLocal(p.Pkg, v.Type()); t != nil {
				if named, ok := ret.named[t.Obj()]; ok {
					named.Consts = append(named.Consts, Const{v})
				}
			} else {
				ret.Consts = append(ret.Consts, Const{v})
			}
		case *types.Var:
			ret.Vars = append(ret.Vars, Var{v})
		}
	}
	return
}

// -----------------------------------------------------------------------------

type sigKindType int

const (
	sigNormal sigKindType = iota
	sigCreator
)

func sigKind(pkg *types.Package, sig *types.Signature) (kind sigKindType, t *types.Named) {
	rets := sig.Results()
	if rets.Len() > 0 {
		if t := checkLocal(pkg, rets.At(0).Type()); t != nil {
			return sigCreator, t
		}
	}
	return sigNormal, nil
}

func checkLocal(pkg *types.Package, first types.Type) *types.Named {
	if t, ok := indirect(first).(*types.Named); ok {
		if t.Obj().Pkg() == pkg {
			return t
		}
	}
	return nil
}

func indirect(typ types.Type) types.Type {
	if t, ok := typ.(*types.Pointer); ok {
		return t.Elem()
	}
	return typ
}

// -----------------------------------------------------------------------------

type Const struct {
	*types.Const
}

type Var struct {
	*types.Var
}

type Func struct {
	*gox.Func
}

func (p Func) Obj() *types.Func {
	return &p.Func.Func
}

func (p Func) Doc() string {
	return p.Comments().Text()
}

// -----------------------------------------------------------------------------

type TypeName struct {
	*types.TypeName
	Creators []Func
	Consts   []Const
}

func (p *TypeName) Type() Type {
	return Type{p.TypeName.Type()}
}

// -----------------------------------------------------------------------------

type Type struct {
	types.Type
}

func (p Type) CheckNamed() (_ Named, ok bool) {
	ret, ok := p.Type.(*types.Named)
	if ok {
		return Named{ret}, true
	}
	return
}

// -----------------------------------------------------------------------------

type Named struct {
	*types.Named
}

func (p Named) Methods() []Func {
	n := p.NumMethods()
	ret := make([]Func, n)
	for i := 0; i < n; i++ {
		fn := p.Method(i)
		ret[i] = Func{gox.MethodFrom(fn)}
	}
	return ret
}

// -----------------------------------------------------------------------------
