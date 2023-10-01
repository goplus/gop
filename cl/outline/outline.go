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
	"strings"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/token"
	"github.com/goplus/gox"
	"github.com/goplus/mod/modfile"
	"golang.org/x/tools/go/types/typeutil"
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

	pkg   *types.Package
	named map[*types.TypeName]*TypeName
}

// aliasr typeutil.Map // types.Type => *TypeName
func checkAlias(aliasr *typeutil.Map, t types.Type) *TypeName {
	if v := aliasr.At(t); v != nil {
		return v.(*TypeName)
	}
	return nil
}

func setAlias(aliasr *typeutil.Map, t types.Type, named *TypeName) {
	real := indirect(t)
	if aliasr.Set(real, named) != nil { // conflict: has old value
		aliasr.Set(real, nil)
	}
}

func (p *All) initNamed(aliasr *typeutil.Map, objs []types.Object, all bool) {
	for _, o := range objs {
		if t, ok := o.(*types.TypeName); ok {
			named := &TypeName{TypeName: t}
			p.named[t] = named
			if all || o.Exported() {
				p.Types = append(p.Types, named)
			}
			if t.IsAlias() {
				setAlias(aliasr, t.Type(), named)
			}
		}
	}
}

func (p *All) lookupNamed(pkg *types.Package, name string) (_ *TypeName, ok bool) {
	o := pkg.Scope().Lookup(name)
	if o == nil {
		return
	}
	t, ok := o.(*types.TypeName)
	if !ok {
		return
	}
	return p.getNamed(t), true
}

func (p *All) getNamed(t *types.TypeName) *TypeName {
	if named, ok := p.named[t]; ok {
		return named
	}
	panic("getNamed: type not found - " + t.Name())
}

func (p Package) Outline(withUnexported ...bool) (ret *All) {
	ret = &All{
		pkg:   p.Pkg,
		named: make(map[*types.TypeName]*TypeName),
	}
	all := (withUnexported != nil && withUnexported[0])
	aliasr := &typeutil.Map{}
	ret.initNamed(aliasr, p.objs, all)
	for _, o := range p.objs {
		if _, ok := o.(*types.TypeName); ok || !(all || o.Exported()) {
			continue
		}
		switch v := o.(type) {
		case *gox.Func:
			sig := v.Type().(*types.Signature)
			if sig.Recv() == nil {
				if name, ok := checkGoptFunc(o.Name()); ok {
					if named, ok := ret.lookupNamed(p.Pkg, name); ok {
						named.GoptFuncs = append(named.GoptFuncs, Func{v})
						continue
					}
				}
			}
			kind, named := ret.sigKind(aliasr, sig)
			switch kind {
			case sigNormal:
				ret.Funcs = append(ret.Funcs, Func{v})
			case sigCreator:
				named.Creators = append(named.Creators, Func{v})
			case sigHelper:
				named.Helpers = append(named.Helpers, Func{v})
			}
		case *types.Const:
			if named := ret.checkLocal(aliasr, v.Type()); named != nil {
				named.Consts = append(named.Consts, Const{v})
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
	sigHelper
)

func (p *All) sigKind(aliasr *typeutil.Map, sig *types.Signature) (sigKindType, *TypeName) {
	rets := sig.Results()
	if rets.Len() > 0 {
		if t := p.checkLocal(aliasr, rets.At(0).Type()); t != nil {
			return sigCreator, t
		}
	}
	params := sig.Params()
	if params.Len() > 0 {
		if t := p.checkLocal(aliasr, params.At(0).Type()); t != nil {
			return sigHelper, t
		}
	}
	return sigNormal, nil
}

func (p *All) checkLocal(aliasr *typeutil.Map, first types.Type) *TypeName {
	first = indirect(first)
	if t, ok := first.(*types.Named); ok {
		o := t.Obj()
		if o.Pkg() == p.pkg {
			return p.getNamed(o)
		}
	}
	return checkAlias(aliasr, first)
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

func (p Const) Obj() types.Object {
	return p.Const
}

func (p Const) Doc() string {
	return ""
}

type Var struct {
	*types.Var
}

func (p Var) Obj() types.Object {
	return p.Var
}

func (p Var) Doc() string {
	return ""
}

type Func struct {
	*gox.Func
}

func (p Func) Obj() types.Object {
	return &p.Func.Func
}

func (p Func) Doc() string {
	return p.Comments().Text()
}

func CheckOverload(obj types.Object) (name string, fn *types.Func, ok bool) {
	if fn, ok = obj.(*types.Func); ok {
		name, ok = checkOverloadFunc(fn.Name())
	}
	return
}

const (
	goptPrefix = "Gopt_"
)

func isGoptFunc(name string) bool {
	return strings.HasPrefix(name, goptPrefix)
}

func isOverloadFunc(name string) bool {
	n := len(name)
	return n > 3 && name[n-3:n-1] == "__"
}

func checkGoptFunc(name string) (string, bool) {
	if isGoptFunc(name) {
		name = name[len(goptPrefix):]
		if pos := strings.IndexByte(name, '_'); pos > 0 {
			return name[:pos], true
		}
	}
	return "", false
}

func checkOverloadFunc(name string) (string, bool) {
	if isOverloadFunc(name) {
		return name[:len(name)-3], true
	}
	return "", false
}

// -----------------------------------------------------------------------------

type TypeName struct {
	*types.TypeName
	Consts    []Const
	Creators  []Func
	GoptFuncs []Func
	Helpers   []Func
}

func (p *TypeName) Obj() types.Object {
	return p.TypeName
}

func (p *TypeName) Doc() string {
	return ""
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
