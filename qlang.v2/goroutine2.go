package qlang

import (
	"qlang.io/exec.v2"
	"qiniupkg.com/text/tpl.v1/interpreter.util"
)

// -----------------------------------------------------------------------------

func (p *Compiler) Go(e interpreter.Engine) {

	src, _ := p.gstk.Pop()
	instr := p.code.Reserve()
	p.exits = append(p.exits, func() {
		start, end := p.cl(e, "expr", src)
		instr.Set(exec.Go(start, end))
	})
}

// -----------------------------------------------------------------------------

func (p *Compiler) ChanIn() {

	p.code.Block(exec.ChanIn)
}

func (p *Compiler) ChanOut() {

	p.code.Block(exec.ChanOut)
}

// -----------------------------------------------------------------------------

