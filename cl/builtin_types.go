package qlang

import (
	"qlang.io/exec"
	qlang "qlang.io/spec"
)

// -----------------------------------------------------------------------------

func (p *Compiler) structInit() {

	arity := p.popArity()
	p.code.Block(exec.StructInit((arity << 1) + 1))
}

func (p *Compiler) mapInit() {

	arity := p.popArity()
	p.code.Block(exec.MapInit((arity << 1) + 1))
}

// -----------------------------------------------------------------------------

func (p *Compiler) tMap() {

	p.code.Block(exec.Map)
}

func (p *Compiler) vMap() {

	arity := p.popArity()
	p.code.Block(exec.Call(qlang.MapFrom, arity<<1))
}

// -----------------------------------------------------------------------------

func (p *Compiler) tSlice() {

	p.code.Block(exec.Slice)
}

func (p *Compiler) vSlice() {

	hasSlice := p.popArity()
	hasInit := 0
	arityInit := 0
	if hasSlice > 0 {
		hasInit = p.popArity()
		if hasInit > 0 {
			arityInit = p.popArity()
		}
	}
	arity := p.popArity()

	if hasSlice > 0 {
		if arity > 0 {
			panic("must be []type")
		}
		if hasInit > 0 { // []T{a1, a2, ...}
			p.code.Block(exec.SliceFromTy(arityInit + 1))
		} else { // []T
			p.code.Block(exec.Slice)
		}
	} else { // [a1, a2, ...]
		p.code.Block(exec.SliceFrom(arity))
	}
}

// -----------------------------------------------------------------------------
