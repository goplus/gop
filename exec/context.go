package exec

// -----------------------------------------------------------------------------

// A Stack represents a FILO container.
type Stack struct {
	data []interface{}
}

func (p *Stack) init() {
	p.data = make([]interface{}, 0, 64)
}

// Get returns the value at specified index.
func (p *Stack) Get(idx int) interface{} {
	return p.data[len(p.data)-idx]
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
	Code   *Code
	parent *Context
	ip     int
}

// NewContext returns a new context of an executor.
func NewContext(code *Code) *Context {
	p := &Context{
		Code: code,
	}
	p.Stack.init()
	return p
}

// Exec executes a code block from ip to ipEnd.
func (ctx *Context) Exec(ip, ipEnd int) {
	data := ctx.Code.data
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
	opPushFloat:   execPushFloat,
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
}

// -----------------------------------------------------------------------------
