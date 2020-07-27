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
)

// -----------------------------------------------------------------------------

type theDefer struct {
	next *theDefer
	n    int // stack len
	i    Instr
}

func (ctx *Context) execDefers() {
	d := ctx.defers
	ctx.defers = nil
	for d != nil {
		ctx.data = ctx.data[:d.n]
		execTable[d.i>>bitsOpShift](d.i, ctx)
		d = d.next
	}
}

func execDefer(i Instr, ctx *Context) {
	i = ctx.code.data[ctx.ip]
	ctx.ip++
	ctx.defers = &theDefer{
		next: ctx.defers,
		n:    len(ctx.data),
		i:    i,
	}
}

// -----------------------------------------------------------------------------

type deferable iBuilder

// CallClosure instr
func (p *deferable) CallClosure(nexpr, arity int, ellipsis bool) exec.Builder {
	return ((*iBuilder)(p)).CallClosure(nexpr, arity, ellipsis)
}

// CallGoClosure instr
func (p *deferable) CallGoClosure(nexpr, arity int, ellipsis bool) exec.Builder {
	return ((*iBuilder)(p)).CallGoClosure(nexpr, arity, ellipsis)
}

// CallFunc instr
func (p *deferable) CallFunc(fun exec.FuncInfo, nexpr int) exec.Builder {
	return ((*iBuilder)(p)).CallFunc(fun, nexpr)
}

// CallFuncv instr
func (p *deferable) CallFuncv(fun exec.FuncInfo, nexpr, arity int) exec.Builder {
	return ((*iBuilder)(p)).CallFuncv(fun, nexpr, arity)
}

// CallGoFunc instr
func (p *deferable) CallGoFunc(fun exec.GoFuncAddr, nexpr int) exec.Builder {
	return ((*iBuilder)(p)).CallGoFunc(fun, nexpr)
}

// CallGoFuncv instr
func (p *deferable) CallGoFuncv(fun exec.GoFuncvAddr, nexpr, arity int) exec.Builder {
	return ((*iBuilder)(p)).CallGoFuncv(fun, nexpr, arity)
}

// Append instr
func (p *deferable) Append(typ reflect.Type, arity int) exec.Builder {
	return ((*iBuilder)(p)).Append(typ, arity)
}

// GoBuiltin instr
func (p *deferable) GoBuiltin(typ reflect.Type, op GoBuiltin) exec.Builder {
	return ((*iBuilder)(p)).GoBuiltin(typ, op)
}

// Defer instr
func (p *Builder) Defer() exec.Deferable {
	p.code.data = append(p.code.data, opDefer<<bitsOpShift)
	return (*deferable)(p)
}

// -----------------------------------------------------------------------------
