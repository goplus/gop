package qlang

import (
	"qlang.io/exec.v2"
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

// -----------------------------------------------------------------------------
