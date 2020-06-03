/*
 Copyright 2020 Qiniu Cloud (qiniu.com)

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
	exec "github.com/qiniu/goplus/exec.spec"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

// A Stack represents a FILO container.
type Stack struct {
	data []interface{}
}

// NewStack creates a Stack instance.
func NewStack() (p *Stack) {
	return &Stack{data: make([]interface{}, 0, 64)}
}

// Init initializes this Stack object.
func (p *Stack) Init() {
	p.data = make([]interface{}, 0, 64)
}

// Get returns the value at specified index.
func (p *Stack) Get(idx int) interface{} {
	return p.data[len(p.data)+idx]
}

// Set returns the value at specified index.
func (p *Stack) Set(idx int, v interface{}) {
	p.data[len(p.data)+idx] = v
}

// GetArgs returns all arguments of a function.
func (p *Stack) GetArgs(arity int) []interface{} {
	return p.data[len(p.data)-arity:]
}

// Ret pops n values from this stack, and then pushes results.
func (p *Stack) Ret(arity int, results ...interface{}) {
	p.data = append(p.data[:len(p.data)-arity], results...)
}

// Push pushes a value into this stack.
func (p *Stack) Push(v interface{}) {
	p.data = append(p.data, v)
}

// PopN pops n elements.
func (p *Stack) PopN(n int) {
	p.data = p.data[:len(p.data)-n]
}

// Pop pops a value from this stack.
func (p *Stack) Pop() interface{} {
	n := len(p.data)
	v := p.data[n-1]
	p.data = p.data[:n-1]
	return v
}

// Len returns count of stack elements.
func (p *Stack) Len() int {
	return len(p.data)
}

// SetLen sets count of stack elements.
func (p *Stack) SetLen(base int) {
	p.data = p.data[:base]
}

// -----------------------------------------------------------------------------

// A Context represents the context of an executor.
type Context struct {
	*Stack
	code   *Code
	parent *Context
	vars   varsContext
	ip     int
	base   int
}

func newSimpleContext(data []interface{}) *Context {
	return &Context{Stack: &Stack{data: data}}
}

// NewContext returns a new context of an executor.
func NewContext(in exec.Code) *Context {
	code := in.(*Code)
	p := &Context{
		Stack: NewStack(),
		code:  code,
	}
	if len(code.vlist) > 0 {
		p.vars = code.makeVarsContext(p)
	}
	return p
}

// NewContextEx creates a closure context, with some local variables.
func newContextEx(parent *Context, stk *Stack, code *Code, vmgr *varManager) *Context {
	p := &Context{
		Stack:  stk,
		code:   code,
		parent: parent,
		base:   len(stk.data),
	}
	if vmgr != nil {
		p.vars = vmgr.makeVarsContext(p)
	}
	return p
}

func (ctx *Context) globalCtx() *Context {
	for ctx.parent != nil {
		ctx = ctx.parent
	}
	return ctx
}

// Exec executes a code block from ip to ipEnd.
func (ctx *Context) Exec(ip, ipEnd int) {
	data := ctx.code.data
	ctx.ip = ip
	for ctx.ip < ipEnd {
		i := data[ctx.ip]
		ctx.ip++
		switch i >> bitsOpShift {
		case opBuiltinOp:
			execBuiltinOp(i, ctx)
		case opJmp:
			execJmp(i, ctx)
		case opJmpIf:
			execJmpIf(i, ctx)
		case opPushInt:
			execPushInt(i, ctx)
		case opPushConstR:
			execPushConstR(i, ctx)
		case opLoadVar:
			execLoadVar(i, ctx)
		case opStoreVar:
			execStoreVar(i, ctx)
		case opCallFunc:
			execFunc(i, ctx)
		case opCallFuncv:
			execFuncv(i, ctx)
		case opCallGoFunc:
			execGoFunc(i, ctx)
		case opCallGoFuncv:
			execGoFuncv(i, ctx)
		case opReturn:
			if i != iReturn {
				ctx.ip = ipReturnN
			}
			return
		case opPushUint:
			execPushUint(i, ctx)
		default:
			if fn := execTable[i>>bitsOpShift]; fn != nil {
				fn(i, ctx)
			} else {
				log.Panicln("Exec failed: unknown instr -", i>>bitsOpShift, "ip:", ctx.ip-1)
			}
		}
	}
}

var _execTable = [...]func(i Instr, p *Context){
	opCallGoFunc:    execGoFunc,
	opCallGoFuncv:   execGoFuncv,
	opPushInt:       execPushInt,
	opPushUint:      execPushUint,
	opPushValSpec:   execPushValSpec,
	opPushConstR:    execPushConstR,
	opIndex:         execIndex,
	opMake:          execMake,
	opAppend:        execAppend,
	opBuiltinOp:     execBuiltinOp,
	opJmp:           execJmp,
	opJmpIf:         execJmpIf,
	opCaseNE:        execCaseNE,
	opPop:           execPop,
	opLoadVar:       execLoadVar,
	opStoreVar:      execStoreVar,
	opAddrVar:       execAddrVar,
	opAddrOp:        execAddrOp,
	opLoad:          execLoad,
	opStore:         execStore,
	opClosure:       execClosure,
	opCallClosure:   execCallClosure,
	opGoClosure:     execGoClosure,
	opCallGoClosure: execCallGoClosure,
	opMakeArray:     execMakeArray,
	opMakeMap:       execMakeMap,
	opZero:          execZero,
	opForPhrase:     execForPhrase,
	opLstComprehens: execListComprehension,
	opMapComprehens: execMapComprehension,
	opTypeCast:      execTypeCast,
	opSlice:         execSlice,
	opSlice3:        execSlice3,
	opMapIndex:      execMapIndex,
	opGoBuiltin:     execGoBuiltin,
}

var execTable []func(i Instr, p *Context)

func init() {
	execTable = _execTable[:]
}

// -----------------------------------------------------------------------------
