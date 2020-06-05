package exec

import (
	"errors"
	"fmt"
	"os"
)

var (
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
	Exit Instr = Call(doExit)
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
	Recover Instr = iRecover(0)
)

// -----------------------------------------------------------------------------
// NewFunction

type Function struct {
	Cls      *Class
	Parent   *Context
	start    int
	end      int
	Args     []string
	Variadic bool
}

func NewFunction(cls *Class, start, end int, args []string, variadic bool) *Function {

	return &Function{cls, nil, start, end, args, variadic}
}

func (p *Function) Call(args ...interface{}) (ret interface{}) {

	return p.ExtCall(nil, args...)
}

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

	var base int
	var vars map[string]interface{}
	var stk *Stack

	if ctx == nil {
		parent := p.Parent
		vars = make(map[string]interface{})
		stk = parent.Stack
		base = stk.BaseFrame()
		ctx = &Context{
			parent: parent,
			Stack:  stk,
			Code:   parent.Code,
			modmgr: parent.modmgr,
			vars:   vars,
			base:   base,
		}
	} else {
		vars = ctx.vars
		stk = ctx.Stack
		base = stk.BaseFrame()
	}

	if p.Variadic {
		for i := 0; i < n-1; i++ {
			vars[p.Args[i]] = args[i]
		}
		vars[p.Args[n-1]] = args[n-1:]
	} else {
		for i, arg := range args {
			vars[p.Args[i]] = arg
		}
	}

	defer func() {
		if e := recover(); e != ErrReturn { // 正常 return 导致，见 (*iReturn).Exec 函数
			ctx.Recov = e
		}
		ret = ctx.ret
		ctx.ExecDefers()
		stk.SetFrame(base)
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

func Func(cls *Class, start, end int, args []string, variadic bool) Instr {
	f := NewFunction(cls, start, end, args, variadic)
	return (*iFunc)(f)
}

// -----------------------------------------------------------------------------
