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

// tStackAddr represents a stack address.
type tStackAddr uint32

// makeStackAddr creates a stack address.
func makeStackAddr(scope, idx uint32) tStackAddr {
	if scope >= (1<<bitsStackScope) || idx > bitsOpStackOperand {
		log.Panicln("makeStackAddr failed: invalid scope or stack index -", scope, idx)
	}
	return tStackAddr((scope << bitsOpStackShift) | idx)
}

func getStackScope(p *Context, addr tStackAddr) (*varScope, uint32) {
	depth := uint32(addr) >> bitsOpStackShift
	idx := uint32(addr & bitsOpStackOperand)
	scope := &p.varScope
	for depth > 0 {
		scope = scope.parent
		depth--
	}
	return scope, idx
}

func execLoad(i Instr, p *Context) {
	idx := i & bitsOperand
	scope, id := getStackScope(p, tStackAddr(idx))
	index := len(scope.args) - int(id)
	p.Push(scope.args[index].Interface())
}

func execAddr(i Instr, p *Context) {
	idx := i & bitsOperand
	scope, id := getStackScope(p, tStackAddr(idx))
	index := len(scope.args) - int(id)
	v := scope.args[index]
	if v.Kind() != reflect.Ptr {
		p.Push(v.Addr().Interface())
	} else {
		p.Push(v.Interface())
	}
}

func execStore(i Instr, p *Context) {
	idx := i & bitsOperand
	scope, id := getStackScope(p, tStackAddr(idx))
	index := len(scope.args) - int(id)
	scope.args[index].Set(reflect.ValueOf(p.Pop()))
}

const (
	closureVariadicFlag = (1 << bitsOpClosureShift)
)

func makeClosure(i Instr, p *Context) Closure {
	idx := i & bitsOpClosureOperand
	var fun *FuncInfo
	if (i & closureVariadicFlag) != 0 {
		fun = p.code.funvs[idx]
	} else {
		fun = p.code.funs[idx]
	}
	if len(p.blockScope) > 0 {
		bs := p.blockScope[len(p.blockScope)-1]
		scope := bs
		scope.vars = make([]reflect.Value, len(bs.vars))
		for i := 0; i < len(scope.vars); i++ {
			scope.vars[i] = bs.vars[i]
		}
		return Closure{fun: fun, parent: &scope}
	}
	return Closure{fun: fun, parent: p.getScope(fun.nestDepth > 1)}
}

func execGoClosure(i Instr, p *Context) {
	closure := makeClosure(i, p)
	v := reflect.MakeFunc(closure.fun.Type(), closure.Call)
	p.Push(v.Interface())
}

func execCallGoClosure(i Instr, p *Context) {
	arity := int(i & bitsOperand)
	fn := reflect.ValueOf(p.Pop())
	t := fn.Type()
	var out []reflect.Value
	if t.IsVariadic() && arity == bitsOperand {
		arity = t.NumIn()
		args := p.GetArgs(arity)
		in := make([]reflect.Value, arity)
		for i, arg := range args {
			in[i] = getArgOf(arg, t, i)
		}
		out = fn.CallSlice(in)
	} else {
		args := p.GetArgs(arity)
		in := make([]reflect.Value, arity)
		for i, arg := range args {
			in[i] = getArgOf(arg, t, i)
		}
		out = fn.Call(in)
	}
	p.PopN(int(arity))
	for _, v := range out {
		p.Push(v.Interface())
	}
}

func execClosure(i Instr, ctx *Context) {
	closure := makeClosure(i, ctx)
	ctx.Push(&closure)
}

func execCallClosure(i Instr, ctx *Context) {
	arity := i & bitsOperand
	c := ctx.Pop().(*Closure)
	fun, parent := c.fun, c.parent
	if fun.IsVariadic() && arity != bitsOperand { // not is: args...
		fun.execVariadic(arity, ctx, parent)
	} else {
		fun.exec(ctx, parent)
	}
}

func execFuncv(i Instr, ctx *Context) {
	idx := i & bitsOpCallFuncvOperand
	arity := (i >> bitsOpCallFuncvShift) & bitsFuncvArityOperand
	fun := ctx.code.funvs[idx]
	parent := ctx.getScope(fun.nestDepth > 1)
	if arity == bitsFuncvArityVar { // args...
		fun.exec(ctx, parent)
	} else {
		if arity == bitsFuncvArityMax {
			arity = uint32(ctx.Pop().(int) + bitsFuncvArityMax)
		}
		fun.execVariadic(arity, ctx, parent)
	}
}

func execFunc(i Instr, ctx *Context) {
	fun := ctx.code.funs[i&bitsOperand]
	fun.exec(ctx, ctx.getScope(fun.nestDepth > 1))
}

// Call calls a function.
func (ctx *Context) Call(f exec.FuncInfo) {
	fun := (*FuncInfo)(f.(*iFuncInfo))
	fun.exec(ctx, ctx.getScope(fun.nestDepth > 1))
}

// -----------------------------------------------------------------------------

// Closure represents a Go+ closure.
type Closure struct {
	fun    *FuncInfo
	recv   interface{}
	parent *varScope
}

// Call calls a closure.
func (p *Closure) Call(in []reflect.Value) (out []reflect.Value) {
	ctx := NewContext(p.fun.Pkg.code)
	for _, v := range in {
		ctx.Push(v.Interface())
	}
	fun := p.fun
	fun.exec(ctx, p.parent)
	n := len(ctx.data)
	if n > 0 {
		out = make([]reflect.Value, n)
		for i, ret := range ctx.data {
			out[i] = getRetOf(ret, fun, i)
		}
	}
	return
}

const (
	nVariadicInvalid      = 0
	nVariadicFixedArgs    = 1
	nVariadicVariadicArgs = 2
)

// FuncInfo represents a Go+ function information.
type FuncInfo struct {
	Pkg      *Package
	name     string
	funEntry int
	funEnd   int
	t        reflect.Type
	in       []reflect.Type
	anyUnresolved
	numOut int
	varManager
	nVariadic uint16
}

// NewFunc create a Go+ function.
func NewFunc(name string, nestDepth uint32) *FuncInfo {
	f := &FuncInfo{
		name:       name,
		varManager: varManager{nestDepth: nestDepth},
	}
	return f
}

func newFuncWith(pkg *Package, name string, nestDepth uint32) *FuncInfo {
	f := &FuncInfo{
		Pkg:        pkg,
		name:       name,
		varManager: varManager{nestDepth: nestDepth},
	}
	return f
}

// Name returns the function name.
func (p *FuncInfo) Name() string {
	return p.name
}

// NumIn returns a function's input parameter count.
func (p *FuncInfo) NumIn() int {
	return len(p.in)
}

// NumOut returns a function's output parameter count.
func (p *FuncInfo) NumOut() int {
	return p.numOut
}

// Out returns the type of a function type's i'th output parameter.
func (p *FuncInfo) Out(i int) *Var {
	if i >= p.numOut {
		log.Panicln("FuncInfo.Out: out of range -", i, "func:", p.name)
	}
	return p.vlist[i]
}

// IsUnnamedOut returns if function results unnamed or not.
func (p *FuncInfo) IsUnnamedOut() bool {
	if p.numOut > 0 {
		return p.vlist[0].IsUnnamedOut()
	}
	return false
}

// Args sets argument types of a Go+ function.
func (p *FuncInfo) Args(in ...reflect.Type) *FuncInfo {
	p.in = in
	p.setVariadic(nVariadicFixedArgs)
	return p
}

// Vargs sets argument types of a variadic Go+ function.
func (p *FuncInfo) Vargs(in ...reflect.Type) *FuncInfo {
	if in[len(in)-1].Kind() != reflect.Slice {
		log.Panicln("Vargs failed: last argument must be a slice.")
	}
	p.in = in
	p.setVariadic(nVariadicVariadicArgs)
	return p
}

// Return sets return types of a Go+ function.
func (p *FuncInfo) Return(out ...*Var) *FuncInfo {
	if p.vlist != nil {
		log.Panicln("don't call DefineVar before calling Return.")
	}
	p.addVar(out...)
	p.numOut = len(out)
	return p
}

// IsVariadic returns if this function is variadic or not.
func (p *FuncInfo) IsVariadic() bool {
	if p.nVariadic == 0 {
		log.Panicln("FuncInfo is unintialized.")
	}
	return p.nVariadic == nVariadicVariadicArgs
}

func (p *FuncInfo) setVariadic(nVariadic uint16) {
	if p.nVariadic == 0 {
		p.nVariadic = nVariadic
	} else if p.nVariadic != nVariadic {
		log.Panicln("setVariadic failed: unmatched -", p.name)
	}
}

// Type returns type of this function.
func (p *FuncInfo) Type() reflect.Type {
	if p.t == nil {
		out := make([]reflect.Type, p.numOut)
		for i := 0; i < p.numOut; i++ {
			out[i] = p.vlist[i].typ
		}
		p.t = reflect.FuncOf(p.in, out, p.IsVariadic())
	}
	return p.t
}

func (p *FuncInfo) execFunc(ctx *Context) {
	oldDefers := ctx.defers
	ctx.defers = nil
	defer func() {
		ctx.execDefers()
		ctx.defers = oldDefers
	}()
	ctx.Exec(p.funEntry, p.funEnd)
	if ctx.ip == ipReturnN { // TODO: optimize
		if ctx.defers != nil {
			ctx.ip = ipInvalid
			rets := ctx.GetArgs(p.numOut)
			for i, v := range rets {
				ctx.setVar(uint32(i), v)
			}
		} else {
			n := len(ctx.data)
			ctx.data = append(ctx.data[:ctx.base-len(p.in)], ctx.data[n-p.numOut:]...)
		}
	}
}

func (p *FuncInfo) exec(ctx *Context, parent *varScope) {
	old := ctx.switchScope(parent, &p.varManager, p.in)
	p.execFunc(ctx)
	if ctx.ip != ipReturnN {
		ctx.data = ctx.data[:ctx.base-len(p.in)]
		n := uint32(p.numOut)
		for i := uint32(0); i < n; i++ {
			ctx.data = append(ctx.data, ctx.getVar(i))
		}
	}
	ctx.restoreScope(old)
}

func (p *FuncInfo) execVariadic(arity uint32, ctx *Context, parent *varScope) {
	var n = uint32(len(p.in) - 1)
	if arity > n {
		tVariadic := p.in[n]
		nVariadic := int(arity - n)
		if tVariadic == exec.TyEmptyInterfaceSlice {
			var empty []interface{}
			ctx.Ret(nVariadic, append(empty, ctx.GetArgs(nVariadic)...))
		} else {
			variadic := reflect.MakeSlice(tVariadic, nVariadic, nVariadic)
			items := ctx.GetArgs(nVariadic)
			for i, item := range items {
				setValue(variadic.Index(i), item)
			}
			ctx.Ret(nVariadic, variadic.Interface())
		}
	}
	p.exec(ctx, parent)
}

// -----------------------------------------------------------------------------

func (p *Builder) resolveFuncs() {
	data := p.code.data
	for fun, pos := range p.funcs {
		if pos < 0 {
			log.Panicln("resolveFuncs failed: func is not defined -", fun.name)
		}
		for _, off := range fun.offs {
			if isClosure(data[off]>>bitsOpShift) && fun.IsVariadic() {
				data[off] |= closureVariadicFlag | uint32(pos)
			} else {
				data[off] |= uint32(pos)
			}
		}
		fun.offs = nil
	}
}

func isClosure(op uint32) bool {
	return op == opClosure || op == opGoClosure
}

// DefineFunc instr
func (p *Builder) DefineFunc(fun *FuncInfo) *Builder {
	if idx, ok := p.funcs[fun]; ok && idx >= 0 {
		log.Panicln("DefineFunc failed: func is defined already -", fun.name)
	}
	p.varManager = &fun.varManager
	fun.funEntry = len(p.code.data)
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
	if p.varManager != &fun.varManager {
		log.Panicln("EndFunc failed: doesn't match with DefineFunc -", fun.name)
	}
	fun.funEnd = len(p.code.data)
	p.varManager = &p.code.varManager
	return p
}

// Closure instr
func (p *Builder) Closure(fun *FuncInfo) *Builder {
	code := p.code
	fun.offs = append(fun.offs, len(code.data))
	code.data = append(code.data, opClosure<<bitsOpShift)
	return p
}

// GoClosure instr
func (p *Builder) GoClosure(fun *FuncInfo) *Builder {
	code := p.code
	fun.offs = append(fun.offs, len(code.data))
	code.data = append(code.data, opGoClosure<<bitsOpShift)
	return p
}

// CallClosure instr
func (p *Builder) CallClosure(arity int) *Builder {
	p.code.data = append(p.code.data, (opCallClosure<<bitsOpShift)|(uint32(arity)&bitsOperand))
	return p
}

// CallGoClosure instr
func (p *Builder) CallGoClosure(arity int) *Builder {
	p.code.data = append(p.code.data, (opCallGoClosure<<bitsOpShift)|(uint32(arity)&bitsOperand))
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
func (p *Builder) Return(n int32) *Builder {
	p.code.data = append(p.code.data, opReturn<<bitsOpShift|(uint32(n)&bitsOperand))
	return p
}

// Load instr
func (p *Builder) Load(fun *FuncInfo, idx int32) *Builder {
	addr := makeStackAddr(p.nestDepth-fun.nestDepth, uint32(-idx))
	p.code.data = append(p.code.data, (opLoad<<bitsOpShift)|uint32(addr))
	return p
}

// Addr instr
func (p *Builder) Addr(fun *FuncInfo, idx int32) *Builder {
	addr := makeStackAddr(p.nestDepth-fun.nestDepth, uint32(-idx))
	p.code.data = append(p.code.data, (opAddr<<bitsOpShift)|uint32(addr))
	return p
}

// Store instr
func (p *Builder) Store(fun *FuncInfo, idx int32) *Builder {
	addr := makeStackAddr(p.nestDepth-fun.nestDepth, uint32(-idx))
	p.code.data = append(p.code.data, (opStore<<bitsOpShift)|uint32(addr))
	return p
}

// -----------------------------------------------------------------------------
