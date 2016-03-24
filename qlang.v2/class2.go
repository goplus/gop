package qlang

import (
	"qlang.io/exec.v2"
	"qiniupkg.com/text/tpl.v1/interpreter.util"
)

// -----------------------------------------------------------------------------

type FunctionInfo struct {
	args     []string // args[0] => function name
	fnb      interface{}
	variadic bool
}

// -----------------------------------------------------------------------------

func (p *Compiler) MemberFuncDecl() {

	fnb, _ := p.gstk.Pop()
	variadic := p.popArity()
	arity := p.popArity()
	args := p.gstk.PopFnArgs(arity + 1)
	fn := &FunctionInfo{
		args:     args,
		fnb:      fnb,
		variadic: variadic != 0,
	}
	p.gstk.Push(fn)
}

func (p *Compiler) newClass(e interpreter.Engine, members []interface{}) *exec.Class {

	cls := exec.Class_()
	for _, val := range members {
		v := val.(*FunctionInfo)
		name := v.args[0]
		v.args[0] = "this"
		start, end := p.cl(e, "doc", v.fnb)
		fn := exec.NewFunction(cls, start, end, v.args, v.variadic)
		cls.Fns[name] = fn
	}
	return cls
}

func (p *Compiler) Class(e interpreter.Engine) {

	arity := p.popArity()
	members := p.gstk.PopNArgs(arity)
	instr := p.code.Reserve()
	p.exits = append(p.exits, func() {
		cls := p.newClass(e, members)
		instr.Set(cls)
	})
}

func (p *Compiler) New() {

	nArgs := p.popArity()
	if nArgs != 0 {
		nArgs = p.popArity()
	}
	p.code.Block(exec.New_(nArgs))
}

func (p *Compiler) MemberRef(name string) {

	p.code.Block(exec.MemberRef(name))
}

// -----------------------------------------------------------------------------
