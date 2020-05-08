package exec

import (
	"reflect"

	"github.com/qiniu/x/log"
)

func execLoad(i Instr, p *Context) {
	idx := int32(i) << bitsOp >> bitsOp
	p.Push(p.data[p.base+int(idx)])
}

func execStore(i Instr, p *Context) {
	idx := int32(i) << bitsOp >> bitsOp
	p.data[p.base+int(idx)] = p.Pop()
}

func execFunc(i Instr, p *Context) {
	idx := i & bitsOperand
	p.code.funs[idx].exec(p)
}

func execFuncv(i Instr, p *Context) {
	idx := i & bitsOpCallFuncvOperand
	arity := (i >> bitsOpCallFuncvShift) & bitsFuncvArityOperand
	fun := p.code.funvs[idx]
	if arity == bitsFuncvArityVar {
		fun.exec(p)
	} else {
		if arity == bitsFuncvArityMax {
			arity = uint32(p.Pop().(int) + bitsFuncvArityMax)
		}
		fun.execVariadic(arity, p)
	}
}

// -----------------------------------------------------------------------------

// Package represents a qlang package.
type Package struct {
}

const (
	nVariadicInvalid      = 0
	nVariadicFixedArgs    = 1
	nVariadicVariadicArgs = 2
)

// FuncInfo represents a qlang function information.
type FuncInfo struct {
	Pkg      *Package
	Name     string
	FunEntry int
	FunEnd   int
	Vars     []*Var
	In       []reflect.Type
	Out      []reflect.Type
	anyUnresolved
	nestDepth uint32
	nVariadic uint32
}

// NewFunc create a qlang function.
func NewFunc(name string) *FuncInfo {
	return &FuncInfo{Name: name}
}

// Args sets argument types of a qlang function.
func (p *FuncInfo) Args(in ...reflect.Type) *FuncInfo {
	p.In = in
	p.setVariadic(nVariadicFixedArgs)
	return p
}

// Vargs sets argument types of a variadic qlang function.
func (p *FuncInfo) Vargs(in ...reflect.Type) *FuncInfo {
	if in[len(in)-1].Kind() != reflect.Slice {
		log.Panicln("Vargs failed: last argument must be a slice.")
	}
	p.In = in
	p.setVariadic(nVariadicVariadicArgs)
	return p
}

// Return sets return types of a qlang function.
func (p *FuncInfo) Return(out ...reflect.Type) *FuncInfo {
	p.Out = out
	return p
}

// DefineVar defines variables in this function.
func (p *FuncInfo) DefineVar(vars ...*Var) *FuncInfo {
	p.Vars = vars
	return p
}

// IsVariadic returns if this function is variadic or not.
func (p *FuncInfo) IsVariadic() bool {
	if p.nVariadic == 0 {
		log.Panicln("FuncInfo is unintialized.")
	}
	return p.nVariadic == nVariadicVariadicArgs
}

func (p *FuncInfo) setVariadic(nVariadic uint32) {
	if p.nVariadic == 0 {
		p.nVariadic = nVariadic
	} else if p.nVariadic != nVariadic {
		log.Panicln("setVariadic failed: unmatched -", p.Name)
	}
}

// Type returns type of this function.
func (p *FuncInfo) Type() reflect.Type {
	return reflect.FuncOf(p.In, p.Out, p.IsVariadic())
}

func (p *FuncInfo) exec(ctx *Context) {
	sub := ctx.NewNest(p.Vars...)
	sub.Exec(p.FunEntry, p.FunEnd)
	ctx.PopN(len(p.In))
}

func (p *FuncInfo) execVariadic(arity uint32, ctx *Context) {
	var n = uint32(len(p.In) - 1)
	var last interface{}
	if arity > n {
		tVariadic := p.In[n]
		nVariadic := arity - n
		if tVariadic == TyEmptyInterface {
			last = ctx.GetArgs(nVariadic)
			ctx.Ret(nVariadic, last)
		} else {
			variadic := reflect.MakeSlice(tVariadic, int(nVariadic), int(nVariadic))
			items := ctx.GetArgs(nVariadic)
			for i, item := range items {
				setValue(variadic.Index(i), item)
			}
			ctx.Ret(nVariadic, variadic.Interface())
		}
	}
	p.exec(ctx)
}

// -----------------------------------------------------------------------------

func (p *Builder) resolveFuncs() {
	data := p.code.data
	for fun, pos := range p.funcs {
		if pos < 0 {
			log.Panicln("resolveFuncs failed: func is not defined -", fun.Name)
		}
		for _, off := range fun.offs {
			data[off] |= uint32(pos)
		}
		fun.offs = nil
	}
}

// DefineFunc instr
func (p *Builder) DefineFunc(fun *FuncInfo) *Builder {
	if idx, ok := p.funcs[fun]; ok && idx >= 0 {
		log.Panicln("DefineFunc failed: func is defined already -", fun.Name)
	}
	p.NestDepth++
	fun.FunEntry = len(p.code.data)
	fun.nestDepth = p.NestDepth
	p.DefineVar(fun.Vars...)
	if fun.IsVariadic() {
		p.funcs[fun] = len(p.code.funvs)
		p.code.funvs = append(p.code.funvs, fun)
	} else {
		p.funcs[fun] = len(p.code.funs)
		p.code.funs = append(p.code.funs, fun)
	}
	return p
}

// EndFunc instr
func (p *Builder) EndFunc(fun *FuncInfo) *Builder {
	if fun.nestDepth != p.NestDepth {
		log.Panicln("EndFunc failed: doesn't match with DefineFunc -", fun.Name)
	}
	fun.FunEnd = len(p.code.data)
	p.NestDepth--
	return p
}

// CallFunc instr
func (p *Builder) CallFunc(fun *FuncInfo) *Builder {
	fun.setVariadic(nVariadicFixedArgs)
	if _, ok := p.funcs[fun]; !ok {
		p.funcs[fun] = -1
	}
	code := p.code
	fun.offs = append(fun.offs, len(code.data))
	code.data = append(code.data, opCallFunc<<bitsOpShift)
	return p
}

// CallFuncv instr
func (p *Builder) CallFuncv(fun *FuncInfo, arity int) *Builder {
	fun.setVariadic(nVariadicVariadicArgs)
	if _, ok := p.funcs[fun]; !ok {
		p.funcs[fun] = -1
	}
	if arity < 0 {
		arity = bitsFuncvArityVar
	} else if arity >= bitsFuncvArityMax {
		p.Push(arity - bitsFuncvArityMax)
		arity = bitsFuncvArityMax
	}
	code := p.code
	fun.offs = append(fun.offs, len(code.data))
	i := (opCallFuncv << bitsOpShift) | (uint32(arity) << bitsOpCallFuncvShift)
	code.data = append(code.data, i)
	return p
}

// Return instr
func (p *Builder) Return() *Builder {
	p.code.data = append(p.code.data, opReturn<<bitsOpShift)
	return p
}

// Load instr
func (p *Builder) Load(idx int32) *Builder {
	p.code.data = append(p.code.data, (opLoad<<bitsOpShift)|(uint32(idx)&bitsOperand))
	return p
}

// Store instr
func (p *Builder) Store(idx int32) *Builder {
	p.code.data = append(p.code.data, (opStore<<bitsOpShift)|(uint32(idx)&bitsOperand))
	return p
}

// -----------------------------------------------------------------------------
