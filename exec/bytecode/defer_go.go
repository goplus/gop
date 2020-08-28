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

type theDefer struct {
	next  *theDefer
	state varScope // block state
	n     int      // stack len
	i     Instr
}

func (ctx *Context) execDefers() {
	d := ctx.defers
	ctx.defers = nil
	for d != nil {
		ctx.data = ctx.data[:d.n]
		ctx.varScope = d.state
		execTable[d.i>>bitsOpShift](d.i, ctx)
		d = d.next
	}
}

func execDefer(i Instr, ctx *Context) {
	i = ctx.code.data[ctx.ip]
	ctx.ip++
	ctx.defers = &theDefer{
		next:  ctx.defers,
		state: ctx.varScope,
		n:     len(ctx.data),
		i:     i,
	}
}

// Defer instr
func (p *Builder) Defer() *Builder {
	p.code.data = append(p.code.data, opDefer<<bitsOpShift)
	return p
}

// -----------------------------------------------------------------------------

func execGo(i Instr, ctx *Context) {
	arity := int(i & bitsOperand)
	i = ctx.code.data[ctx.ip]
	ctx.ip++
	ctx.Go(arity, func(goctx *Context) {
		execTable[i>>bitsOpShift](i, goctx)
	})
}

type goBuilder struct {
	*iBuilder
}

func (p goBuilder) builder(arity int) *iBuilder {
	p.code.data = append(p.code.data, (opGo<<bitsOpShift)|Instr(arity))
	return p.iBuilder
}

// CallClosure instr
func (p goBuilder) CallClosure(nexpr, arity int, ellipsis bool) exec.Builder {
	return p.builder(arity).CallClosure(nexpr, arity, ellipsis)
}

// CallGoClosure instr
func (p goBuilder) CallGoClosure(nexpr, arity int, ellipsis bool) exec.Builder {
	return p.builder(arity).CallGoClosure(nexpr, arity, ellipsis)
}

// CallFunc instr
func (p goBuilder) CallFunc(fun exec.FuncInfo, nexpr int) exec.Builder {
	return p.builder(fun.NumIn()).CallFunc(fun, nexpr)
}

// CallFuncv instr
func (p goBuilder) CallFuncv(fun exec.FuncInfo, nexpr, arity int) exec.Builder {
	return p.builder(arity).CallFuncv(fun, nexpr, arity)
}

// CallGoFunc instr
func (p goBuilder) CallGoFunc(fun GoFuncAddr, nexpr int) exec.Builder {
	arity := NewPackage(p.code).GetGoFuncType(fun).NumIn()
	return p.builder(arity).CallGoFunc(fun, nexpr)
}

// CallGoFuncv instr
func (p goBuilder) CallGoFuncv(fun GoFuncvAddr, nexpr, arity int) exec.Builder {
	return p.builder(arity).CallGoFuncv(fun, nexpr, arity)
}

// GoBuiltin instr
func (p goBuilder) GoBuiltin(typ reflect.Type, op exec.GoBuiltin) exec.Builder {
	arity := goBuiltinArity(op)
	return p.builder(arity).GoBuiltin(typ, op)
}

func goBuiltinArity(op exec.GoBuiltin) int {
	switch op {
	case exec.GobCopy, exec.GobDelete:
		return 2
	case exec.GobClose:
		return 1
	}
	panic("no reach here")
}

// Go instr
func (p *Builder) Go() exec.Builder {
	return goBuilder{(*iBuilder)(p)}
}

// -----------------------------------------------------------------------------
