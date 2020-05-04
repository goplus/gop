package exec

import (
	"reflect"
)

// -----------------------------------------------------------------------------

// A Stack represents a FILO container.
type Stack struct {
	data []interface{}
}

// Init initializes this Stack object.
func (p *Stack) Init() {
	p.data = make([]interface{}, 0, 64)
}

// Get returns the value at specified index.
func (p *Stack) Get(idx int) interface{} {
	return p.data[len(p.data)+idx]
}

// GetArgs returns all arguments of a function.
func (p *Stack) GetArgs(arity uint32) []interface{} {
	return p.data[len(p.data)-int(arity):]
}

// Ret pops n values from this stack, and then pushes results.
func (p *Stack) Ret(arity uint32, results ...interface{}) {
	p.data = append(p.data[:len(p.data)-int(arity)], results...)
}

// Push pushes a value into this stack.
func (p *Stack) Push(v interface{}) {
	p.data = append(p.data, v)
}

// Top returns the last pushed value, if it exists.
func (p *Stack) Top() (v interface{}, ok bool) {
	n := len(p.data)
	if n > 0 {
		v, ok = p.data[n-1], true
	}
	return
}

// Pop pops a value from this stack.
func (p *Stack) Pop() (v interface{}, ok bool) {
	n := len(p.data)
	if n > 0 {
		v, ok = p.pop(), true
	}
	return
}

func (p *Stack) pop() interface{} {
	n := len(p.data)
	v := p.data[n-1]
	p.data = p.data[:n-1]
	return v
}

// -----------------------------------------------------------------------------

// A Context represents the context of an executor.
type Context struct {
	Stack
	code   *Code
	parent *Context
	vars   reflect.Value
	ip     int
}

func newSimpleContext(data []interface{}) *Context {
	return &Context{Stack: Stack{data: data}}
}

// NewContext returns a new context of an executor.
func NewContext(code *Code) *Context {
	p := &Context{
		code: code,
	}
	p.Stack.Init()
	return p
}

// NewContextEx returns a new context of an executor, with n local variables.
func NewContextEx(code *Code, si *StructInfo, parent *Context) *Context {
	p := &Context{
		code:   code,
		parent: parent,
	}
	if si != nil {
		p.vars = reflect.New(si.Type()).Elem()
	}
	p.Stack.Init()
	return p
}

// AddrVar returns a variable address by index.
func (ctx *Context) AddrVar(idx uint32) interface{} {
	return ctx.vars.Field(int(idx)).Addr().Interface()
}

// GetVar returns a variable value by index.
func (ctx *Context) GetVar(idx uint32) interface{} {
	return ctx.vars.Field(int(idx)).Interface()
}

// SetVar returns a variable value by index.
func (ctx *Context) SetVar(idx uint32, v interface{}) {
	ctx.vars.Field(int(idx)).Set(reflect.ValueOf(v))
}

// Exec executes a code block from ip to ipEnd.
func (ctx *Context) Exec(ip, ipEnd int) {
	data := ctx.code.data
	ctx.ip = ip
	for ctx.ip != ipEnd {
		i := data[ctx.ip]
		ctx.ip++
		switch i >> bitsOpShift {
		case opBuiltinOp:
			execBuiltinOp(i, ctx)
		case opPushInt:
			execPushInt(i, ctx)
		case opPushStringR:
			execPushStringR(i, ctx)
		case opLoadVar:
			execLoadVar(i, ctx)
		case opStoreVar:
			execStoreVar(i, ctx)
		case opCallGoFun:
			execGoFun(i, ctx)
		case opCallGoFunv:
			execGoFunv(i, ctx)
		case opPushUint:
			execPushUint(i, ctx)
		default:
			execTable[i>>bitsOpShift](i, ctx)
		}
	}
}

var execTable = [...]func(i Instr, p *Context){
	opPushInt:     execPushInt,
	opPushUint:    execPushUint,
	opPushValSpec: execPushValSpec,
	opPushStringR: execPushStringR,
	opPushIntR:    execPushIntR,
	opPushUintR:   execPushUintR,
	opPushFloatR:  execPushFloatR,
	opBuiltinOp:   execBuiltinOp,
	opBuiltinOpS:  execBuiltinOpS,
	opJmp:         execJmp,
	opJmpIfFalse:  execJmpIfFalse,
	opCaseNE:      execCaseNE,
	opPop:         execPop,
	opCallGoFun:   execGoFun,
	opCallGoFunv:  execGoFunv,
	opLoadVar:     execLoadVar,
	opStoreVar:    execStoreVar,
	opAddrVar:     execAddrVar,
	opAddrOp:      execAddrOp,
}

// -----------------------------------------------------------------------------
