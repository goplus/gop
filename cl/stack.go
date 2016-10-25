package qlang

import (
	"strconv"

	"qlang.io/exec"
)

// -----------------------------------------------------------------------------

func (p *Compiler) pushInt(v int) {

	p.code.Block(exec.Push(v))
}

func (p *Compiler) pushFloat(v float64) {

	p.code.Block(exec.Push(v))
}

func (p *Compiler) pushByte(lit string) {

	v, multibyte, tail, err := strconv.UnquoteChar(lit[1:len(lit)-1], '\'')
	if err != nil {
		panic("invalid char `" + lit + "`: " + err.Error())
	}
	if tail != "" || multibyte {
		panic("invalid char: " + lit)
	}
	p.code.Block(exec.Push(byte(v)))
}

func (p *Compiler) pushString(lit string) {

	v, err := strconv.Unquote(lit)
	if err != nil {
		panic("invalid string `" + lit + "`: " + err.Error())
	}
	p.code.Block(exec.Push(v))
}

func (p *Compiler) pushID(name string) {

	p.code.Block(exec.Push(name))
}

// -----------------------------------------------------------------------------

func (p *Compiler) popArity() int {

	if v, ok := p.gstk.Pop(); ok {
		if arity, ok := v.(int); ok {
			return arity
		}
	}
	panic("no arity")
}

func (p *Compiler) popName() string {

	if v, ok := p.gstk.Pop(); ok {
		if name, ok := v.(string); ok {
			return name
		}
	}
	panic("no ident name")
}

func (p *Compiler) pushCode(code interface{}) {

	p.gstk.Push(code)
}

func (p *Compiler) arity(arity int) {

	p.gstk.Push(arity)
}

func (p *Compiler) pushName(name string) {

	p.gstk.Push(name)
}

// -----------------------------------------------------------------------------
