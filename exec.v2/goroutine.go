package exec

import (
	"qlang.io/qlang.spec.v1"
)

// -----------------------------------------------------------------------------

type iGo struct {
	start int
	end   int
}

func (p *iGo) Exec(stk *Stack, ctx *Context) {

	go func() {
		fn := NewFunction(nil, p.start, p.end, nil, false)
		fn.parent = ctx
		fn.ExtCall(nil)
	}()
	stk.Push(nil)
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

