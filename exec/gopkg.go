package exec

import (
	"reflect"

	"github.com/qiniu/x/log"
)

func execGoFunc(i Instr, p *Context) {
	idx := i & bitsOperand
	gofuns[idx].exec(0, p)
}

func execGoFuncv(i Instr, p *Context) {
	idx := i & bitsOpCallFuncvOperand
	arity := (i >> bitsOpCallFuncvShift) & bitsFuncvArityMax
	if arity == bitsFuncvArityMax {
		arity = uint32(p.Pop().(int) + bitsFuncvArityMax)
	}
	gofunvs[idx].exec(arity, p)
}

// -----------------------------------------------------------------------------

// SymbolKind represents symbol kind.
type SymbolKind uint32

const (
	// SymbolFunc - function
	SymbolFunc SymbolKind = opCallGoFunc
	// SymbolVariadicFunc - variadic function
	SymbolVariadicFunc SymbolKind = opCallGoFuncv
	// SymbolVar - variable
	SymbolVar SymbolKind = 0
)

// GoPackage represents a Go package.
type GoPackage struct {
	PkgPath string
	syms    map[string]uint32
}

// NewGoPackage creates a new builtin Go Package.
func NewGoPackage(pkgPath string) *GoPackage {
	if _, ok := gopkgs[pkgPath]; ok {
		log.Fatalln("NewPackage failed: package exists -", pkgPath)
	}
	pkg := &GoPackage{PkgPath: pkgPath, syms: make(map[string]uint32)}
	gopkgs[pkgPath] = pkg
	return pkg
}

// FindGoPackage lookups a Go package by pkgPath. It returns nil if not found.
func FindGoPackage(pkgPath string) *GoPackage {
	return gopkgs[pkgPath]
}

// Find lookups a symbol by specified its name.
func (p *GoPackage) Find(name string) (addr uint32, kind SymbolKind, ok bool) {
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

// FindVariadicFunc lookups a Go function by name.
func (p *GoPackage) FindVariadicFunc(name string) (addr GoVariadicFuncAddr, ok bool) {
	if v, ok := p.syms[name]; ok {
		if (v >> bitsOpShift) == opCallGoFuncv {
			return GoVariadicFuncAddr(v & bitsOperand), true
		}
	}
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
		p.syms[v.Name] = (uint32(base) + uint32(i)) | (opCallGoFunc << bitsOpShift)
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
		p.syms[v.Name] = (uint32(base) + uint32(i)) | (opCallGoFuncv << bitsOpShift)
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

// CallGoFunc instr
func (p *Builder) CallGoFunc(fun GoFuncAddr) *Builder {
	p.code.data = append(p.code.data, (opCallGoFunc<<bitsOpShift)|uint32(fun))
	return p
}

// CallGoFuncv instr
func (p *Builder) CallGoFuncv(fun GoVariadicFuncAddr, arity int) *Builder {
	code := p.code
	if arity >= bitsFuncvArityMax {
		p.Push(arity - bitsFuncvArityMax)
		arity = bitsFuncvArityMax
	}
	i := (opCallGoFuncv << bitsOpShift) | (uint32(arity) << bitsOpCallFuncvShift) | uint32(fun)
	code.data = append(code.data, i)
	return p
}

// -----------------------------------------------------------------------------
