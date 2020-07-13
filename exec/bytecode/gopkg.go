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
	"reflect"

	"github.com/goplus/gop/exec.spec"
	"github.com/qiniu/x/log"
)

func execGoFunc(i Instr, p *Context) {
	idx := i & bitsOperand
	gofuns[idx].exec(0, p)
}

func execGoFuncv(i Instr, p *Context) {
	idx := i & bitsOpCallFuncvOperand
	arity := int((i >> bitsOpCallFuncvShift) & bitsFuncvArityOperand)
	fun := gofunvs[idx]
	if arity == bitsFuncvArityVar {
		v := p.Pop()
		args := reflect.ValueOf(v)
		n := args.Len()
		for i := 0; i < n; i++ {
			p.Push(args.Index(i).Interface())
		}
		arity = fun.getNumIn() - 1 + n
	} else if arity == bitsFuncvArityMax {
		arity = p.Pop().(int) + bitsFuncvArityMax
	}
	fun.exec(arity, p)
}

func execLoadGoVar(i Instr, p *Context) {
	idx := i & bitsOperand
	v := reflect.ValueOf(govars[idx].Addr).Elem()
	p.Push(v.Interface())
}

func execStoreGoVar(i Instr, p *Context) {
	idx := i & bitsOperand
	v := reflect.ValueOf(govars[idx].Addr).Elem()
	v.Set(reflect.ValueOf(p.Pop()))
}

func execAddrGoVar(i Instr, p *Context) {
	idx := i & bitsOperand
	p.Push(govars[idx].Addr)
}

func execLoadGoField(i Instr, p *Context) {
	index := p.Pop()
	v := reflect.ValueOf(p.Pop())
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	p.Push(v.FieldByIndex(index.([]int)).Interface())
}

func execStoreGoField(i Instr, p *Context) {
	index := p.Pop()
	v := reflect.ValueOf(p.Pop())
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	v.FieldByIndex(index.([]int)).Set(reflect.ValueOf(p.Pop()))
}

func execAddrGoField(i Instr, p *Context) {
	index := p.Pop()
	v := reflect.ValueOf(p.Pop()).Elem()
	p.Push(v.FieldByIndex(index.([]int)).Addr().Interface())
}

// -----------------------------------------------------------------------------

// A ConstKind represents the specific kind of type that a Type represents.
// The zero Kind is not a valid kind.
type ConstKind = exec.ConstKind

const (
	// ConstBoundRune - bound type: rune
	ConstBoundRune = exec.ConstBoundRune
	// ConstBoundString - bound type: string
	ConstBoundString = exec.ConstBoundString
	// ConstUnboundInt - unbound int type
	ConstUnboundInt = exec.ConstUnboundInt
	// ConstUnboundFloat - unbound float type
	ConstUnboundFloat = exec.ConstUnboundFloat
	// ConstUnboundComplex - unbound complex type
	ConstUnboundComplex = exec.ConstUnboundComplex
	// ConstUnboundPtr - nil: unbound ptr
	ConstUnboundPtr = exec.ConstUnboundPtr
)

// SymbolKind represents symbol kind.
type SymbolKind = exec.SymbolKind

const (
	// SymbolVar - variable
	SymbolVar = exec.SymbolVar
	// SymbolFunc - function
	SymbolFunc = exec.SymbolFunc
	// SymbolFuncv - variadic function
	SymbolFuncv = exec.SymbolFuncv
)

// GoPackage represents a Go package.
type GoPackage struct {
	pkgPath string
	syms    map[string]uint32
	types   map[string]reflect.Type
	consts  map[string]*GoConstInfo
}

// NewGoPackage creates a new builtin Go Package.
func NewGoPackage(pkgPath string) *GoPackage {
	if _, ok := gopkgs[pkgPath]; ok {
		log.Panicln("NewPackage failed: package exists -", pkgPath)
	}
	pkg := &GoPackage{
		pkgPath: pkgPath,
		syms:    make(map[string]uint32),
		types:   make(map[string]reflect.Type),
		consts:  make(map[string]*GoConstInfo),
	}
	gopkgs[pkgPath] = pkg
	return pkg
}

// FindGoPackage lookups a Go package by pkgPath. It returns nil if not found.
func FindGoPackage(pkgPath string) exec.GoPackage {
	return gopkgs[pkgPath]
}

// PkgPath returns the package path for importing.
func (p *GoPackage) PkgPath() string {
	return p.pkgPath
}

// Find lookups a symbol by specified its name.
func (p *GoPackage) Find(name string) (addr uint32, kind SymbolKind, ok bool) {
	if p == nil {
		return
	}
	if v, ok := p.syms[name]; ok {
		return v & bitsOperand, SymbolKind(v >> bitsOpShift), true
	}
	return
}

// FindFunc lookups a Go function by name.
func (p *GoPackage) FindFunc(name string) (addr GoFuncAddr, ok bool) {
	if v, ok := p.syms[name]; ok {
		if (v >> bitsOpShift) == opCallGoFunc {
			return GoFuncAddr(v & bitsOperand), true
		}
	}
	return
}

// FindFuncv lookups a Go function by name.
func (p *GoPackage) FindFuncv(name string) (addr GoFuncvAddr, ok bool) {
	if v, ok := p.syms[name]; ok {
		if (v >> bitsOpShift) == opCallGoFuncv {
			return GoFuncvAddr(v & bitsOperand), true
		}
	}
	return
}

// FindConst lookups a Go constant by name.
func (p *GoPackage) FindConst(name string) (ci *GoConstInfo, ok bool) {
	ci, ok = p.consts[name]
	return
}

// FindVar lookups a Go variable by name.
func (p *GoPackage) FindVar(name string) (addr GoVarAddr, ok bool) {
	if v, ok := p.syms[name]; ok {
		if (v >> bitsOpShift) == 0 {
			return GoVarAddr(v), true
		}
	}
	return
}

// FindType lookups a Go type by name.
func (p *GoPackage) FindType(name string) (typ reflect.Type, ok bool) {
	typ, ok = p.types[name]
	return
}

// Const creates a GoConstInfo instance.
func (p *GoPackage) Const(name string, kind ConstKind, val interface{}) GoConstInfo {
	return GoConstInfo{Pkg: p, Name: name, Kind: kind, Value: val}
}

// Var creates a GoVarInfo instance.
func (p *GoPackage) Var(name string, addr interface{}) GoVarInfo {
	if log.CanOutput(log.Ldebug) {
		if reflect.TypeOf(addr).Kind() != reflect.Ptr {
			log.Panicln("variable address isn't a pointer?")
		}
	}
	return GoVarInfo{Pkg: p, Name: name, Addr: addr}
}

// Func creates a GoFuncInfo instance.
func (p *GoPackage) Func(name string, fn interface{}, exec func(i int, p *Context)) GoFuncInfo {
	return GoFuncInfo{Pkg: p, Name: name, This: fn, exec: exec}
}

// Funcv creates a GoFuncvInfo instance.
func (p *GoPackage) Funcv(name string, fn interface{}, exec func(i int, p *Context)) GoFuncvInfo {
	return GoFuncvInfo{GoFuncInfo{Pkg: p, Name: name, This: fn, exec: exec}, 0}
}

// Type creates a GoTypeInfo instance.
func (p *GoPackage) Type(name string, typ reflect.Type) GoTypeInfo {
	return GoTypeInfo{Pkg: p, Name: name, Type: typ}
}

// Rtype gets the real type information.
func (p *GoPackage) Rtype(typ reflect.Type) GoTypeInfo {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return GoTypeInfo{Pkg: p, Name: typ.Name(), Type: typ}
}

// RegisterVars registers all exported Go variables of this package.
func (p *GoPackage) RegisterVars(vars ...GoVarInfo) (base GoVarAddr) {
	base = GoVarAddr(len(govars))
	govars = append(govars, vars...)
	for i, v := range vars {
		p.syms[v.Name] = uint32(base) + uint32(i)
	}
	return
}

// RegisterConsts registers all exported Go constants of this package.
func (p *GoPackage) RegisterConsts(consts ...GoConstInfo) {
	for i := range consts {
		ci := &consts[i]
		if ci.Kind == ConstUnboundInt { // TODO
			ci.Value = reflect.ValueOf(ci.Value).Int()
		}
		p.consts[ci.Name] = ci
	}
}

// RegisterFuncs registers all exported Go functions of this package.
func (p *GoPackage) RegisterFuncs(funs ...GoFuncInfo) (base GoFuncAddr) {
	if log.CanOutput(log.Ldebug) {
		for _, v := range funs {
			if v.Pkg != p {
				log.Panicln("function doesn't belong to this package:", v.Name)
			}
			if v.This != nil && reflect.TypeOf(v.This).IsVariadic() {
				log.Panicln("function is variadic? -", v.Name)
			}
		}
	}
	base = GoFuncAddr(len(gofuns))
	gofuns = append(gofuns, funs...)
	for i, v := range funs {
		p.syms[v.Name] = (uint32(base) + uint32(i)) | (opCallGoFunc << bitsOpShift)
	}
	return
}

// RegisterFuncvs registers all exported Go functions with variadic arguments of this package.
func (p *GoPackage) RegisterFuncvs(funs ...GoFuncvInfo) (base GoFuncvAddr) {
	if log.CanOutput(log.Ldebug) {
		for _, v := range funs {
			if v.Pkg != p {
				log.Panicln("function doesn't belong to this package:", v.Name)
			}
			if v.This != nil && !reflect.TypeOf(v.This).IsVariadic() {
				log.Panicln("function isn't variadic? -", v.Name)
			}
		}
	}
	base = GoFuncvAddr(len(gofunvs))
	gofunvs = append(gofunvs, funs...)
	for i, v := range funs {
		p.syms[v.Name] = (uint32(base) + uint32(i)) | (opCallGoFuncv << bitsOpShift)
	}
	return
}

// RegisterTypes registers all exported Go types defined by this package.
func (p *GoPackage) RegisterTypes(typinfos ...GoTypeInfo) {
	for _, ti := range typinfos {
		if p != ti.Pkg {
			log.Panicln("RegisterTypes failed: unmatched package instance.")
		}
		if ti.Name == "" {
			log.Panicln("RegisterTypes failed: unnamed type? -", ti.Type)
		}
		if _, ok := p.types[ti.Name]; ok {
			log.Panicln("RegisterTypes failed: register an existed type -", p.pkgPath, ti.Name)
		}
		p.types[ti.Name] = ti.Type
	}
}

// -----------------------------------------------------------------------------

var (
	gopkgs  = make(map[string]exec.GoPackage)
	gofuns  []GoFuncInfo
	gofunvs []GoFuncvInfo
	govars  []GoVarInfo
)

// GoFuncAddr represents a Go function address.
type GoFuncAddr = exec.GoFuncAddr

// GoFuncvAddr represents a variadic Go function address.
type GoFuncvAddr = exec.GoFuncvAddr

// GoVarAddr represents a variadic Go variable address.
type GoVarAddr = exec.GoVarAddr

// GoFuncInfo represents a Go function information.
type GoFuncInfo struct {
	Pkg  *GoPackage
	Name string
	This interface{}
	exec func(arity int, p *Context)
}

// GoFuncvInfo represents a Go function information.
type GoFuncvInfo struct {
	GoFuncInfo
	numIn int // cache
}

func (p *GoFuncvInfo) getNumIn() int {
	if p.numIn == 0 {
		p.numIn = reflect.TypeOf(p.This).NumIn()
	}
	return p.numIn
}

// GoTypeInfo represents a Go type information.
type GoTypeInfo struct {
	Pkg  *GoPackage
	Name string
	Type reflect.Type
}

// GoConstInfo represents a Go constant information.
type GoConstInfo = exec.GoConstInfo

// GoVarInfo represents a Go variable information.
type GoVarInfo struct {
	Pkg  *GoPackage
	Name string
	Addr interface{}
}

// CallGoFunc instr
func (p *Builder) CallGoFunc(fun GoFuncAddr) *Builder {
	p.code.data = append(p.code.data, (opCallGoFunc<<bitsOpShift)|uint32(fun))
	return p
}

// CallGoFuncv instr
func (p *Builder) CallGoFuncv(fun GoFuncvAddr, arity int) *Builder {
	if arity < 0 {
		arity = bitsFuncvArityVar
	} else if arity >= bitsFuncvArityMax {
		p.Push(arity - bitsFuncvArityMax)
		arity = bitsFuncvArityMax
	}
	i := (opCallGoFuncv << bitsOpShift) | (uint32(arity) << bitsOpCallFuncvShift) | uint32(fun)
	p.code.data = append(p.code.data, i)
	return p
}

// LoadGoVar instr
func (p *Builder) LoadGoVar(addr GoVarAddr) *Builder {
	i := (opLoadGoVar << bitsOpShift) | uint32(addr)
	p.code.data = append(p.code.data, i)
	return p
}

// StoreGoVar instr
func (p *Builder) StoreGoVar(addr GoVarAddr) *Builder {
	i := (opStoreGoVar << bitsOpShift) | uint32(addr)
	p.code.data = append(p.code.data, i)
	return p
}

// AddrGoVar instr
func (p *Builder) AddrGoVar(addr GoVarAddr) *Builder {
	i := (opAddrGoVar << bitsOpShift) | uint32(addr)
	p.code.data = append(p.code.data, i)
	return p
}

// LoadGoField instr
func (p *Builder) LoadGoField(index []int) *Builder {
	p.Push(index)
	i := (opLoadGoField << bitsOpShift) | uint32(len(index))
	p.code.data = append(p.code.data, uint32(i))
	return p
}

// StoreGoField instr
func (p *Builder) StoreGoField(index []int) *Builder {
	p.Push(index)
	i := (opStoreGoField << bitsOpShift) | (uint32(len(index)))
	p.code.data = append(p.code.data, i)
	return p
}

// AddrGoField instr
func (p *Builder) AddrGoField(index []int) *Builder {
	p.Push(index)
	i := (opAddrGoField << bitsOpShift) | uint32(len(index))
	p.code.data = append(p.code.data, i)
	return p
}

// -----------------------------------------------------------------------------
