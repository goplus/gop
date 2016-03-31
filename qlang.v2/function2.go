package qlang

import (
	"strconv"

	"qlang.io/exec.v2"
	"qiniupkg.com/text/tpl.v1/interpreter.util"
)

// -----------------------------------------------------------------------------

func (p *Compiler) Clear() {

	p.code.Block(exec.Clear)
}

func (p *Compiler) Unset(name string) {

	p.code.Block(exec.Unset(name))
}

func (p *Compiler) Inc(name string) {

	p.code.Block(exec.Inc(name))
}

func (p *Compiler) Dec(name string) {

	p.code.Block(exec.Dec(name))
}

func (p *Compiler) AddAssign(name string) {

	p.code.Block(exec.AddAssign(name))
}

func (p *Compiler) SubAssign(name string) {

	p.code.Block(exec.SubAssign(name))
}

func (p *Compiler) MulAssign(name string) {

	p.code.Block(exec.MulAssign(name))
}

func (p *Compiler) QuoAssign(name string) {

	p.code.Block(exec.QuoAssign(name))
}

func (p *Compiler) ModAssign(name string) {

	p.code.Block(exec.ModAssign(name))
}

func (p *Compiler) MultiAssign() {

	arity := p.popArity() + 1
	names := p.gstk.PopFnArgs(arity)
	p.code.Block(exec.MultiAssign(names))
}

func (p *Compiler) Assign(name string) {

	p.code.Block(exec.Assign(name))
}

func (p *Compiler) Ref(name string) {

	if val, ok := p.gvars[name]; ok {
		p.code.Block(exec.Push(val))
	} else {
		p.code.Block(exec.Ref(name))
	}
}

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

func (p *Compiler) Main(e interpreter.Engine) {

	fnb, _ := p.gstk.Pop()
	instr := p.code.Reserve()
	p.exits = append(p.exits, func() {
		start, end := p.cl(e, "doc", fnb)
		instr.Set(exec.Func(nil, start, end, nil, false))
	})
	p.code.Block(
		exec.CallFn(0),
		exec.Exit,
	)
}

func (p *Compiler) Include(lit string) {

	file, err := strconv.Unquote(lit)
	if err != nil {
		panic("invalid string `" + lit + "`: " + err.Error())
	}

	code := p.code
	instr := code.Reserve()
	p.exits = append(p.exits, func() {
		start := code.Len()
		end := p.Incl(file)
		instr.Set(exec.Block(start, end))
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
