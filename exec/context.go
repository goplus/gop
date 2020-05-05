package exec

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

// -----------------------------------------------------------------------------

// A Context represents the context of an executor.
type Context struct {
	Stack
	code   *Code
	parent *Context
	vars   varsContext
	ip     int
}

func newSimpleContext(data []interface{}) *Context {
	return &Context{Stack: Stack{data: data}}
}

// NewContext returns a new context of an executor.
func NewContext(code *Code, vars ...*Var) *Context {
	p := &Context{
		code: code,
	}
	if len(vars) > 0 {
		p.vars = makeVarsContext(vars, p)
	}
	p.Stack.Init()
	return p
}

// NewNest creates a nest closure context, with some local variables.
func (ctx *Context) NewNest(vars ...*Var) *Context {
	p := &Context{
		code:   ctx.code,
		parent: ctx,
	}
	if len(vars) > 0 {
		p.vars = makeVarsContext(vars, p)
	}
	p.Stack.Init()
	return p
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
