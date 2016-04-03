package qlang

import (
	"strconv"

	"qlang.io/exec.v2"
)

// -----------------------------------------------------------------------------

func (p *Compiler) PushInt(v int) {

	p.code.Block(exec.Push(v))
}

func (p *Compiler) PushFloat(v float64) {

	p.code.Block(exec.Push(v))
}

func (p *Compiler) PushByte(lit string) {

	v, multibyte, tail, err := strconv.UnquoteChar(lit[1:len(lit)-1], '\'')
	if err != nil {
		panic("invalid char `" + lit + "`: " + err.Error())
	}
	if tail != "" || multibyte {
		panic("invalid char: " + lit)
	}
	p.code.Block(exec.Push(byte(v)))
}

func (p *Compiler) PushString(lit string) {

	v, err := strconv.Unquote(lit)
	if err != nil {
		panic("invalid string `" + lit + "`: " + err.Error())
	}
	p.code.Block(exec.Push(v))
}

func (p *Compiler) PushNil() {

	p.code.Block(exec.Nil)
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

func (p *Compiler) PushCode(code interface{}) {

	p.gstk.Push(code)
}

func (p *Compiler) Arity(arity int) {

	p.gstk.Push(arity)
}

func (p *Compiler) PushName(name string) {

	p.gstk.Push(name)
}

// -----------------------------------------------------------------------------

