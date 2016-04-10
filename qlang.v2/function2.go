package qlang

import (
	"qlang.io/exec.v2"
	"qiniupkg.com/text/tpl.v1/interpreter.util"
)

// -----------------------------------------------------------------------------

func (p *Compiler) Function(e interpreter.Engine) {

	fnb, _ := p.gstk.Pop()
	variadic := p.popArity()
	arity := p.popArity()
	args := p.gstk.PopFnArgs(arity)
	instr := p.code.Reserve()
	p.exits = append(p.exits, func() {
		start, end := p.cl(e, "doc", fnb)
		instr.Set(exec.Func(nil, start, end, args, variadic != 0))
	})
}

func (p *Compiler) AnonymFn(e interpreter.Engine) {

	fnb, _ := p.gstk.Pop()
	instr := p.code.Reserve()
	p.exits = append(p.exits, func() {
		start, end := p.cl(e, "doc", fnb)
		instr.Set(exec.AnonymFn(start, end))
	})
}

func (p *Compiler) Return(e interpreter.Engine) {

	arity := p.popArity()
	p.code.Block(exec.Return(arity))
}

func (p *Compiler) Done() {

	for {
		n := len(p.exits)
		if n == 0 {
			break
		}
		onExit := p.exits[n-1]
		p.exits = p.exits[:n-1]
		onExit()
	}
}

func (p *Compiler) Defer(e interpreter.Engine) {

	src, _ := p.gstk.Pop()
	instr := p.code.Reserve()
	p.exits = append(p.exits, func() {
		start, end := p.cl(e, "expr", src)
		p.CodeLine(src)
		instr.Set(exec.Defer(start, end))
	})
}

func (p *Compiler) Recover() {

	p.code.Block(exec.Recover)
}

func (p *Compiler) cl(e interpreter.Engine, g string, src interface{}) (start, end int) {

	start = p.code.Len()
	if src != nil {
		if err := e.EvalCode(p, g, src); err != nil {
			panic(err)
		}
	}
	end = p.code.Len()
	return
}

// -----------------------------------------------------------------------------
