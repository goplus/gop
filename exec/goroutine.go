package exec

import (
	qlang "qlang.io/spec"
)

// -----------------------------------------------------------------------------

type iGo struct {
	start int
	end   int
}

func (p *iGo) Exec(stk *Stack, ctx *Context) {

	go func() {
		cloneCtx := *ctx // 每个goroutine需要自己的上下文，有自己的堆栈
		cloneCtx.Stack = NewStack()
		cloneCtx.Code.Exec(p.start, p.end, cloneCtx.Stack, &cloneCtx)
	}()
}

func Go(start, end int) Instr {

	return &iGo{start, end}
}

// -----------------------------------------------------------------------------

type iChanIn int
type iChanOut int

func (p iChanIn) Exec(stk *Stack, ctx *Context) {

	v, _ := stk.Pop()
	ch, _ := stk.Pop()
	ret := qlang.ChanIn(ch, v, ctx.onsel)
	stk.Push(ret)
}

func (p iChanOut) Exec(stk *Stack, ctx *Context) {

	ch, _ := stk.Pop()
	ret := qlang.ChanOut(ch, ctx.onsel)
	stk.Push(ret)
}

var (
	ChanIn  Instr = iChanIn(0)
	ChanOut Instr = iChanOut(0)
)

// -----------------------------------------------------------------------------
