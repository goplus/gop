package exec

import (
	"reflect"

	"github.com/qiniu/qlang/v6/ast/spec"
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

func execMakeArray(i Instr, p *Context) {
	typSlice := getType(i&bitsOpMakeArrayOperand, p)
	arity := int((i >> bitsOpMakeArrayShift) & bitsFuncvArityOperand)
	if arity == bitsFuncvArityVar { // args...
		v := reflect.ValueOf(p.Get(-1))
		n := v.Len()
		ret := reflect.MakeSlice(typSlice, n, n)
		reflect.Copy(ret, v)
		p.Ret(1, ret.Interface())
	} else {
		if arity == bitsFuncvArityMax {
			arity = p.Pop().(int) + bitsFuncvArityMax
		}
		makeArray(typSlice, arity, p)
	}
}

func makeArray(typSlice reflect.Type, arity int, p *Context) {
	args := p.GetArgs(arity)
	var ret reflect.Value
	if typSlice.Kind() == reflect.Slice {
		ret = reflect.MakeSlice(typSlice, arity, arity)
	} else {
		ret = reflect.New(typSlice).Elem()
	}
	for i, arg := range args {
		ret.Index(i).Set(getElementOf(arg, typSlice))
	}
	p.Ret(arity, ret.Interface())
}

func execMakeMap(i Instr, p *Context) {
	typMap := getType(i&bitsOpMakeArrayOperand, p)
	arity := int((i >> bitsOpMakeArrayShift) & bitsFuncvArityOperand)
	if arity == bitsFuncvArityMax {
		arity = p.Pop().(int) + bitsFuncvArityMax
	}
	makeMap(typMap, arity, p)
}

func makeMap(typMap reflect.Type, arity int, p *Context) {
	n := arity << 1
	args := p.GetArgs(n)
	ret := reflect.MakeMapWithSize(typMap, arity)
	for i := 0; i < n; i += 2 {
		key := getKeyOf(args[i], typMap)
		val := getElementOf(args[i+1], typMap)
		ret.SetMapIndex(key, val)
	}
	p.Ret(n, ret.Interface())
}

func execZero(i Instr, p *Context) {
	typ := getType(i&bitsOpMakeArrayOperand, p)
	p.Push(reflect.Zero(typ).Interface())
}

// -----------------------------------------------------------------------------

// A ConstKind represents the specific kind of type that a Type represents.
// The zero Kind is not a valid kind.
type ConstKind = spec.ConstKind

const (
	// ConstBoundRune - bound type: rune
	ConstBoundRune = spec.ConstBoundRune
	// ConstBoundString - bound type: string
	ConstBoundString = spec.ConstBoundString
	// ConstUnboundInt - unbound int type
	ConstUnboundInt = spec.ConstUnboundInt
	// ConstUnboundFloat - unbound float type
	ConstUnboundFloat = spec.ConstUnboundFloat
	// ConstUnboundComplex - unbound complex type
	ConstUnboundComplex = spec.ConstUnboundComplex
	// ConstUnboundPtr - nil: unbound ptr
	ConstUnboundPtr = spec.ConstUnboundPtr
)

// SymbolKind represents symbol kind.
type SymbolKind uint32

const (
	// SymbolFunc - function
	SymbolFunc SymbolKind = opCallGoFunc
	// SymbolFuncv - variadic function
	SymbolFuncv SymbolKind = opCallGoFuncv
	// SymbolVar - variable
	SymbolVar SymbolKind = 0
)

// GoPackage represents a Go package.
type GoPackage struct {
	PkgPath string
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
		PkgPath: pkgPath,
		syms:    make(map[string]uint32),
		types:   make(map[string]reflect.Type),
		consts:  make(map[string]*GoConstInfo),
	}
	gopkgs[pkgPath] = pkg
	return pkg
}

// FindGoPackage lookups a Go package by pkgPath. It returns nil if not found.
func FindGoPackage(pkgPath string) *GoPackage {
	return gopkgs[pkgPath]
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
			log.Panicln("RegisterTypes failed: register an existed type -", p.PkgPath, ti.Name)
		}
		p.types[ti.Name] = ti.Type
	}
}

// -----------------------------------------------------------------------------

var (
	gopkgs  = make(map[string]*GoPackage)
	gofuns  []GoFuncInfo
	gofunvs []GoFuncvInfo
	govars  []GoVarInfo
)

// GoFuncAddr represents a Go function address.
type GoFuncAddr uint32

// GoFuncvAddr represents a variadic Go function address.
type GoFuncvAddr uint32

// GoVarAddr represents a variadic Go variable address.
type GoVarAddr uint32

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
type GoConstInfo struct {
	Pkg   *GoPackage
	Name  string
	Kind  ConstKind
	Value interface{}
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
func (i GoFuncvAddr) GetInfo() *GoFuncInfo {
	if i < GoFuncvAddr(len(gofunvs)) {
		return &gofunvs[i].GoFuncInfo
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
func (p *Builder) CallGoFuncv(fun GoFuncvAddr, arity int) *Builder {
	code := p.code
	if arity < 0 {
		arity = bitsFuncvArityVar
	} else if arity >= bitsFuncvArityMax {
		p.Push(arity - bitsFuncvArityMax)
		arity = bitsFuncvArityMax
	}
	i := (opCallGoFuncv << bitsOpShift) | (uint32(arity) << bitsOpCallFuncvShift) | uint32(fun)
	code.data = append(code.data, i)
	return p
}

// MakeArray instr
func (p *Builder) MakeArray(typ reflect.Type, arity int) *Builder {
	if arity < 0 {
		if typ.Kind() == reflect.Array {
			log.Panicln("MakeArray failed: can't be variadic.")
		}
		arity = bitsFuncvArityVar
	} else if arity >= bitsFuncvArityMax {
		p.Push(arity - bitsFuncvArityMax)
		arity = bitsFuncvArityMax
	}
	code := p.code
	i := (opMakeArray << bitsOpShift) | (uint32(arity) << bitsOpMakeArrayShift) | p.newType(typ)
	code.data = append(code.data, i)
	return p
}

// MakeMap instr
func (p *Builder) MakeMap(typ reflect.Type, arity int) *Builder {
	if arity < 0 {
		log.Panicln("MakeMap failed: can't be variadic.")
	} else if arity >= bitsFuncvArityMax {
		p.Push(arity - bitsFuncvArityMax)
		arity = bitsFuncvArityMax
	}
	code := p.code
	i := (opMakeMap << bitsOpShift) | (uint32(arity) << bitsOpMakeMapShift) | p.newType(typ)
	code.data = append(code.data, i)
	return p
}

// Zero instr
func (p *Builder) Zero(typ reflect.Type) *Builder {
	code := p.code
	i := (opZero << bitsOpShift) | p.requireType(typ)
	code.data = append(code.data, i)
	return p
}

func (p *Builder) requireType(typ reflect.Type) uint32 {
	kind := typ.Kind()
	bt := builtinTypes[kind]
	if bt.size > 0 {
		return uint32(kind)
	}
	return p.newType(typ)
}

func (p *Builder) newType(typ reflect.Type) uint32 {
	if ityp, ok := p.types[typ]; ok {
		return ityp
	}
	code := p.code
	ityp := uint32(len(code.types) + len(builtinTypes))
	code.types = append(code.types, typ)
	p.types[typ] = ityp
	return ityp
}

func getType(ityp uint32, ctx *Context) reflect.Type {
	if ityp < uint32(len(builtinTypes)) {
		return builtinTypes[ityp].typ
	}
	return ctx.code.types[ityp-uint32(len(builtinTypes))]
}

// -----------------------------------------------------------------------------
