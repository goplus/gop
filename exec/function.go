package exec

import (
	"errors"
	"fmt"
	"os"
)

var (
	// ErrReturn is used for `panic(ErrReturn)`
	ErrReturn = errors.New("return")
)

// -----------------------------------------------------------------------------
// Exit

func doExit(ret interface{}) {

	code := 0
	if ret != nil {
		code = ret.(int)
	}
	os.Exit(code)
}

var (
	// Exit is the instruction that executes `exit(<code>)`.
	Exit = Call(doExit)
)

// -----------------------------------------------------------------------------
// Return

type iReturn int

func (arity iReturn) Exec(stk *Stack, ctx *Context) {

	if arity == 0 {
		ctx.ret = nil
	} else if arity == 1 {
		ctx.ret, _ = stk.Pop()
	} else {
		ctx.ret = stk.PopNArgs(int(arity))
	}
	if ctx.parent != nil {
		panic(ErrReturn) // 利用 panic 来实现 return (正常退出)
	}

	ctx.ExecDefers()

	if ctx.ret == nil {
		os.Exit(0)
	}
	if v, ok := ctx.ret.(int); ok {
		os.Exit(v)
	}
	panic("must return `int` for main function")
}

// Return returns an instruction that means `return <expr1>, ..., <exprN>`.
//
func Return(arity int) Instr {
	return iReturn(arity)
}

// -----------------------------------------------------------------------------
// Defer

type iDefer struct {
	start int
	end   int
}

func (p *iDefer) Exec(stk *Stack, ctx *Context) {

	ctx.defers = &theDefer{
		next:  ctx.defers,
		start: p.start,
		end:   p.end,
	}
}

// Defer returns an instruction that executes `defer <expr>`.
//
func Defer(start, end int) Instr {
	return &iDefer{start, end}
}

// -----------------------------------------------------------------------------
// Recover

type iRecover int

func (p iRecover) Exec(stk *Stack, ctx *Context) {

	if parent := ctx.parent; parent != nil {
		e := parent.Recov
		parent.Recov = nil
		stk.Push(e)
	} else {
		stk.Push(nil)
	}
}

var (
	// Recover is the instruction that executes `recover()`.
	Recover Instr = iRecover(0)
)

// -----------------------------------------------------------------------------
// NewFunction

// A Function represents a function object.
//
type Function struct {
	Cls      *Class
	Parent   *Context
	start    int
	end      int
	symtbl   map[string]int
	Args     []string
	Variadic bool
}

// NewFunction creates a new function object.
//
func NewFunction(cls *Class, start, end int, symtbl map[string]int, args []string, variadic bool) *Function {

	return &Function{cls, nil, start, end, symtbl, args, variadic}
}

// Call calls this function with default context.
//
func (p *Function) Call(stk *Stack, args ...interface{}) (ret interface{}) {

	parent := p.Parent
	ctx := &Context{
		parent: parent,
		Stack:  stk,
		Code:   parent.Code,
		modmgr: parent.modmgr,
		base:   stk.BaseFrame(),
	}
	ctx.initVars(p.symtbl)
	return p.ExtCall(ctx, args...)
}

// ExtCall calls this function with a specified context.
//
func (p *Function) ExtCall(ctx *Context, args ...interface{}) (ret interface{}) {

	n := len(p.Args)
	if p.Variadic {
		if len(args) < n-1 {
			panic(fmt.Sprintf("function requires >= %d arguments, but we got %d", n-1, len(args)))
		}
	} else {
		if len(args) != n {
			panic(fmt.Sprintf("function requires %d arguments, but we got %d", n, len(args)))
		}
	}

	if p.start == p.end {
		return nil
	}

	stk := ctx.Stack
	vars := ctx.Vars()
	if p.Variadic {
		for i := 0; i < n-1; i++ {
			vars[i] = args[i]
		}
		vars[n-1] = args[n-1:]
	} else {
		for i, arg := range args {
			vars[i] = arg
		}
	}

	defer func() {
		if e := recover(); e != ErrReturn { // 正常 return 导致，见 (*iReturn).Exec 函数
			ctx.Recov = e
		}
		ret = ctx.ret
		ctx.ExecDefers()
		stk.SetFrame(ctx.base)
		if ctx.Recov != nil {
			panic(ctx.Recov)
		}
	}()
	ctx.Code.Exec(p.start, p.end, stk, ctx)
	return
}

// -----------------------------------------------------------------------------

type iFunc Function

func (p *iFunc) Exec(stk *Stack, ctx *Context) {
	p.Parent = ctx
	stk.Push((*Function)(p))
}

// Func returns an instruction that create a function object.
//
func Func(cls *Class, start, end int, symtbl map[string]int, args []string, variadic bool) Instr {
	f := NewFunction(cls, start, end, symtbl, args, variadic)
	return (*iFunc)(f)
}

// -----------------------------------------------------------------------------
