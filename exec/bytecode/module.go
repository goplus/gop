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

package bytecode

import (
	"github.com/goplus/gop/exec.spec"
	"github.com/goplus/gop/reflect"
)

// -----------------------------------------------------------------------------

// Package represents a Go+ package.
type Package struct {
	code *Code
}

// NewPackage creates a Go+ package instance.
func NewPackage(code *Code) *Package {
	return &Package{code: code}
}

// NewVar creates a variable instance.
func (p *Package) NewVar(typ reflect.Type, name string) exec.Var {
	return NewVar(typ, name)
}

// NewLabel creates a label object.
func (p *Package) NewLabel(name string) exec.Label {
	return NewLabel(name)
}

// NewForPhrase creates a new ForPhrase instance.
func (p *Package) NewForPhrase(in reflect.Type) exec.ForPhrase {
	return NewForPhrase(in)
}

// NewComprehension creates a new Comprehension instance.
func (p *Package) NewComprehension(out reflect.Type) exec.Comprehension {
	return NewComprehension(out)
}

// NewFunc create a Go+ function.
func (p *Package) NewFunc(name string, nestDepth uint32, funcType ...int) exec.FuncInfo {
	return (*iFuncInfo)(newFuncWith(p, name, nestDepth))
}

// FindGoPackage lookups a Go package by pkgPath. It returns nil if not found.
func (p *Package) FindGoPackage(pkgPath string) exec.GoPackage {
	return FindGoPackage(pkgPath)
}

// GetGoFuncType returns a Go function's type.
func (p *Package) GetGoFuncType(addr exec.GoFuncAddr) reflect.Type {
	return reflect.TypeOf(gofuns[addr].This)
}

// GetGoFuncvType returns a Go function's type.
func (p *Package) GetGoFuncvType(addr exec.GoFuncvAddr) reflect.Type {
	return reflect.TypeOf(gofunvs[addr].This)
}

// GetGoFuncInfo returns a Go function's information.
func (p *Package) GetGoFuncInfo(addr exec.GoFuncAddr) *exec.GoFuncInfo {
	gfi := &gofuns[addr]
	return &exec.GoFuncInfo{Pkg: gfi.Pkg, Name: gfi.Name, This: gfi.This}
}

// GetGoFuncvInfo returns a Go function's information.
func (p *Package) GetGoFuncvInfo(addr exec.GoFuncvAddr) *exec.GoFuncInfo {
	gfi := &gofunvs[addr]
	return &exec.GoFuncInfo{Pkg: gfi.Pkg, Name: gfi.Name, This: gfi.This}
}

// GetGoVarInfo returns a Go variable's information.
func (p *Package) GetGoVarInfo(addr exec.GoVarAddr) *exec.GoVarInfo {
	gvi := &govars[addr]
	return &exec.GoVarInfo{Pkg: gvi.Pkg, Name: gvi.Name, This: gvi.Addr}
}

// -----------------------------------------------------------------------------
