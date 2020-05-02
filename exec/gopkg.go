package exec

import (
	"reflect"

	"github.com/qiniu/x/log"
)

func execGoFun(i Instr, p *Context) {
	idx := i & bitsOperand
	gofuns[idx].exec(0, p)
}

func execGoFunv(i Instr, p *Context) {
	idx := i & bitsOpCallGoFunvOperand
	arity := (i >> bitsOpCallGoFunvShift) & bitsGoFunvArityMax
	if arity == bitsGoFunvArityMax {
		arity = uint32(p.pop().(int) + bitsGoFunvArityMax)
	}
	gofunvs[idx].exec(arity, p)
}

// -----------------------------------------------------------------------------

// GoPackage represents a Go package.
type GoPackage struct {
	PkgPath string
	syms    map[string]uint32
}

// NewPackage creates a new builtin Go Package.
func NewPackage(pkgPath string) *GoPackage {
	if _, ok := gopkgs[pkgPath]; ok {
		log.Fatalln("NewPackage failed: package exists -", pkgPath)
	}
	pkg := &GoPackage{PkgPath: pkgPath, syms: make(map[string]uint32)}
	gopkgs[pkgPath] = pkg
	return pkg
}

// Package lookups a Go package by pkgPath. It returns nil if not found.
func Package(pkgPath string) *GoPackage {
	return gopkgs[pkgPath]
}

// FindFunc lookups a Go function by name.
func (p *GoPackage) FindFunc(name string) (addr GoFuncAddr, ok bool) {
	if v, ok := p.syms[name]; ok {
		if (v >> bitsOpShift) == opCallGoFun {
			return GoFuncAddr(v & bitsOperand), true
		}
	}
	return
}

// FindVariadicFunc lookups a Go function by name.
func (p *GoPackage) FindVariadicFunc(name string) (addr GoVariadicFuncAddr, ok bool) {
	if v, ok := p.syms[name]; ok {
		if (v >> bitsOpShift) == opCallGoFunv {
			return GoVariadicFuncAddr(v & bitsOperand), true
		}
	}
	return
}

// FindVar lookups a Go variable by name.
func (p *GoPackage) FindVar(name string) (addr GoVarAddr, ok bool) {
	if v, ok := p.syms[name]; ok {
		if (v >> bitsOpShift) == 0 {
			return GoVarAddr(v & bitsOperand), true
		}
	}
	return
}

// Var creates a GoVarInfo instance.
func (p *GoPackage) Var(name string, addr interface{}) GoVarInfo {
	if log.CanOutput(log.Ldebug) {
		if reflect.TypeOf(addr).Kind() != reflect.Ptr {
			log.Fatalln("variable address isn't a pointer?")
		}
	}
	return GoVarInfo{Pkg: p, Name: name, Addr: addr}
}

// Func creates a GoFuncInfo instance.
func (p *GoPackage) Func(name string, fn interface{}, exec func(i Instr, p *Context)) GoFuncInfo {
	return GoFuncInfo{Pkg: p, Name: name, This: fn, exec: exec}
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

// RegisterFuncs registers all exported Go functions of this package.
func (p *GoPackage) RegisterFuncs(funs ...GoFuncInfo) (base GoFuncAddr) {
	if log.CanOutput(log.Ldebug) {
		for _, v := range funs {
			if v.Pkg != p {
				log.Fatalln("function doesn't belong to this package:", v.Name)
			}
			if v.This != nil && reflect.TypeOf(v.This).IsVariadic() {
				log.Fatalln("function is variadic? -", v.Name)
			}
		}
	}
	base = GoFuncAddr(len(gofuns))
	gofuns = append(gofuns, funs...)
	for i, v := range funs {
		p.syms[v.Name] = (uint32(base) + uint32(i)) | (opCallGoFun << bitsOpShift)
	}
	return
}

// RegisterVariadicFuncs registers all exported Go functions with variadic arguments of this package.
func (p *GoPackage) RegisterVariadicFuncs(funs ...GoFuncInfo) (base GoVariadicFuncAddr) {
	if log.CanOutput(log.Ldebug) {
		for _, v := range funs {
			if v.Pkg != p {
				log.Fatalln("function doesn't belong to this package:", v.Name)
			}
			if v.This != nil && !reflect.TypeOf(v.This).IsVariadic() {
				log.Fatalln("function isn't variadic? -", v.Name)
			}
		}
	}
	base = GoVariadicFuncAddr(len(gofunvs))
	gofunvs = append(gofunvs, funs...)
	for i, v := range funs {
		p.syms[v.Name] = (uint32(base) + uint32(i)) | (opCallGoFunv << bitsOpShift)
	}
	return
}

// -----------------------------------------------------------------------------

var (
	gopkgs  = make(map[string]*GoPackage)
	gofuns  []GoFuncInfo
	gofunvs []GoFuncInfo
	govars  []GoVarInfo
)

// GoFuncAddr represents a Go function address.
type GoFuncAddr uint32

// GoVariadicFuncAddr represents a variadic Go function address.
type GoVariadicFuncAddr uint32

// GoVarAddr represents a variadic Go variable address.
type GoVarAddr uint32

// GoFuncInfo represents a Go function information.
type GoFuncInfo struct {
	Pkg  *GoPackage
	Name string
	This interface{}
	exec func(i Instr, p *Context)
}

// GoVarInfo represents a Go variable information.
type GoVarInfo struct {
	Pkg  *GoPackage
	Name string
	Addr interface{}
}

// GetInfo retuns a Go function info.
func (i GoFuncAddr) GetInfo() *GoFuncInfo {
	if i < GoFuncAddr(len(gofuns)) {
		return &gofuns[i]
	}
	return nil
}

// GetInfo retuns a Go function info.
func (i GoVariadicFuncAddr) GetInfo() *GoFuncInfo {
	if i < GoVariadicFuncAddr(len(gofunvs)) {
		return &gofunvs[i]
	}
	return nil
}

// GetInfo retuns a Go variable info.
func (i GoVarAddr) GetInfo() *GoVarInfo {
	if i < GoVarAddr(len(govars)) {
		return &govars[i]
	}
	return nil
}

// CallGoFun instr
func (p *Builder) CallGoFun(fun GoFuncAddr) *Builder {
	p.code.data = append(p.code.data, (opCallGoFun<<bitsOpShift)|uint32(fun))
	return p
}

// CallGoFunv instr
func (p *Builder) CallGoFunv(fun GoVariadicFuncAddr, arity int) *Builder {
	code := p.code
	if arity >= bitsGoFunvArityMax {
		p.Push(arity - bitsGoFunvArityMax)
		arity = bitsGoFunvArityMax
	}
	i := (opCallGoFunv << bitsOpShift) | (uint32(arity) << bitsOpCallGoFunvShift) | uint32(fun)
	code.data = append(code.data, i)
	return p
}

const (
	bitsGoFunvArityMax = (1 << bitsGoFunvArity) - 1
)

// -----------------------------------------------------------------------------
