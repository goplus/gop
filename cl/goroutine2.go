package qlang

import (
	"github.com/qiniu/text/tpl/interpreter.util"
	"github.com/qiniu/qlang/exec"
)

// -----------------------------------------------------------------------------

func (p *Compiler) fnGo(e interpreter.Engine) {

	src, _ := p.gstk.Pop()
	instr := p.code.Reserve()
	fnctx := p.fnctx
	p.exits = append(p.exits, func() {
		start, end := p.clBlock(e, "expr", src, fnctx)
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
