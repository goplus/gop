package qlang

import (
	"qiniupkg.com/text/tpl.v1/interpreter.util"
	"qlang.io/exec.v2"
)

// -----------------------------------------------------------------------------

func (p *Compiler) fnGo(e interpreter.Engine) {

	src, _ := p.gstk.Pop()
	instr := p.code.Reserve()
	p.exits = append(p.exits, func() {
		start, end := p.cl(e, "expr", src)
		instr.Set(exec.Go(start, end))
	})
}

// -----------------------------------------------------------------------------

func (p *Compiler) chanIn() {

	p.code.Block(exec.ChanIn)
}

func (p *Compiler) chanOut() {

	p.code.Block(exec.ChanOut)
}

func (p *Compiler) tChan() {

	p.code.Block(exec.Chan)
}

// -----------------------------------------------------------------------------
