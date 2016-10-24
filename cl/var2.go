package qlang

import (
	"qlang.io/exec"
	qlang "qlang.io/spec"
)

// -----------------------------------------------------------------------------

func (p *Compiler) inc() {

	p.code.Block(exec.IncEx)
}

func (p *Compiler) dec() {

	p.code.Block(exec.DecEx)
}

func (p *Compiler) addAssign() {

	p.code.Block(exec.AddAssignEx)
}

func (p *Compiler) subAssign() {

	p.code.Block(exec.SubAssignEx)
}

func (p *Compiler) mulAssign() {

	p.code.Block(exec.MulAssignEx)
}

func (p *Compiler) quoAssign() {

	p.code.Block(exec.QuoAssignEx)
}

func (p *Compiler) modAssign() {

	p.code.Block(exec.ModAssignEx)
}

func (p *Compiler) xorAssign() {

	p.code.Block(exec.XorAssignEx)
}

func (p *Compiler) bitandAssign() {

	p.code.Block(exec.BitAndAssignEx)
}

func (p *Compiler) bitorAssign() {

	p.code.Block(exec.BitOrAssignEx)
}

func (p *Compiler) andnotAssign() {

	p.code.Block(exec.AndNotAssignEx)
}

func (p *Compiler) lshrAssign() {

	p.code.Block(exec.LshrAssignEx)
}

func (p *Compiler) rshrAssign() {

	p.code.Block(exec.RshrAssignEx)
}

func (p *Compiler) multiAssign() {

	nval := p.popArity()
	arity := p.popArity() + 1
	if nval == 1 {
		p.code.Block(exec.MultiAssignFromSliceEx(arity))
	} else if arity != nval {
		panic("argument count of multi assignment doesn't match")
	} else {
		p.code.Block(exec.MultiAssignEx(arity))
	}
}

func (p *Compiler) assign() {

	p.code.Block(exec.AssignEx)
}

func (p *funcCtx) requireSymbol(name string) (id int, fnew bool) {
	id, ok := p.getSymbol(name)
	if ok {
		return
	}
	return p.newSymbol(name), true
}

func (p *Compiler) ref(name string) {

	if val, ok := p.gvars[name]; ok {
		p.code.Block(exec.Push(val))
		return
	}

	fnctx := p.fnctx
	id, ok := fnctx.getSymbol(name)
	if !ok {
		if val, ok := qlang.Fntable[name]; ok {
			p.code.Block(exec.GfnRef(val, func() exec.Instr {
				id := fnctx.newSymbol(name)
				return exec.Var(id)
			}))
			return
		}
		id = fnctx.newSymbol(name)
	}
	p.code.Block(exec.Ref(id))
}

func (p *Compiler) index() {

	arity2 := p.popArity()
	arityMid := p.popArity()
	arity1 := p.popArity()

	if arityMid == 0 {
		if arity1 == 0 {
			panic("call operator[] without index")
		}
		p.code.Block(exec.Get)
	} else {
		p.code.Block(exec.Op3(qlang.SubSlice, arity1 != 0, arity2 != 0))
	}
}

func (p *Compiler) toVar() {

	p.code.ToVar()
}

// -----------------------------------------------------------------------------
