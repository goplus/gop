package qlang

import (
	"errors"

	qlang "github.com/qiniu/qlang/spec"
	"github.com/qiniu/text/tpl/interpreter.util"
)

var (
	ErrIfCondNotBoolean = errors.New("if condition isn't a boolean expression")
	ErrForCodeForm      = errors.New("for statement is in illegal form")
)

// -----------------------------------------------------------------------------

func (p *Interpreter) Or(e interpreter.Engine) {

	bCode, _ := p.stk.Pop()
	a, _ := p.stk.Pop()
	a1, ok := a.(bool)
	if !ok {
		panic("left operand of || operator isn't a boolean expression")
	}
	if a1 {
		p.stk.Push(true)
		return
	}
	evalCode(e, p, "term4", bCode)
}

func (p *Interpreter) And(e interpreter.Engine) {

	bCode, _ := p.stk.Pop()
	a, _ := p.stk.Pop()
	a1, ok := a.(bool)
	if !ok {
		panic("left operand of && operator isn't a boolean expression")
	}
	if !a1 {
		p.stk.Push(false)
		return
	}
	evalCode(e, p, "term3", bCode)
}

func (p *Interpreter) If(e interpreter.Engine) {

	var elseCode interface{}

	stk := p.stk
	elseArity := stk.popArity()
	if elseArity == 1 {
		elseCode, _ = stk.Pop()
	}

	condArity := stk.popArity()
	condArity++

	ifbr := stk.popNArgs(condArity << 1) // 2 * (condArity + 1)
	for i := 0; i < condArity; i++ {
		condCode := ifbr[i<<1]
		evalCode(e, p, "expr", condCode)
		v, _ := stk.Pop()
		cond, ok := v.(bool)
		if !ok {
			panic(ErrIfCondNotBoolean)
		}
		if cond {
			bodyCode := ifbr[(i<<1)+1]
			if bodyCode == nil {
				stk.Push(nil)
			} else {
				evalCode(e, p, "doc", bodyCode)
			}
			return
		}
	}
	if elseCode == nil {
		stk.Push(nil)
	} else {
		evalCode(e, p, "doc", elseCode)
	}
}

func (p *Interpreter) Switch(e interpreter.Engine) {

	var defaultCode interface{}

	stk := p.stk
	defaultArity := stk.popArity()
	if defaultArity == 1 {
		defaultCode, _ = stk.Pop()
	}

	caseArity := stk.popArity()
	n := caseArity << 1
	casebr := stk.popNArgs(n) // 2 * caseArity
	switchCode, _ := stk.Pop()

	var switchCond interface{}
	if switchCode == nil {
		switchCond = true
	} else {
		evalCode(e, p, "expr", switchCode)
		switchCond, _ = stk.Pop()
	}

	for i := 0; i < caseArity; i++ {
		caseCode := casebr[i<<1]
		evalCode(e, p, "expr", caseCode)
		caseCond, _ := stk.Pop()
		vCond := qlang.EQ(switchCond, caseCond)
		if cond, ok := vCond.(bool); ok && cond {
			bodyCode := casebr[(i<<1)+1]
			if bodyCode == nil {
				stk.Push(nil)
			} else {
				evalCode(e, p, "doc", bodyCode)
			}
			return
		}
	}
	if defaultCode == nil {
		stk.Push(nil)
	} else {
		evalCode(e, p, "doc", defaultCode)
	}
}

// -----------------------------------------------------------------------------

func (p *Interpreter) For(e interpreter.Engine) {

	stk := p.stk
	bodyCode, _ := stk.Pop()

	var condCode interface{}
	var stepCode interface{}

	arity := stk.popArity()
	switch arity {
	case 1, 3:
		forCode := stk.popNArgs(arity)
		if arity == 3 {
			initCode := forCode[0]
			if initCode != nil {
				evalCode(e, p, "s", initCode)
				stk.Pop()
			}
			stepCode = forCode[2]
		}
		condCode = forCode[arity>>1]
	case 0:
	default:
		panic(ErrForCodeForm)
	}

	for condCode == nil || p.evalCondS(e, condCode) {
		if bodyCode != nil {
			evalCode(e, p, "doc", bodyCode)
			stk.Pop()
		}
		if stepCode != nil {
			evalCode(e, p, "s", stepCode)
			stk.Pop()
		}
	}
	stk.Push(nil)
}

func (p *Interpreter) evalCondS(e interpreter.Engine, condCode interface{}) bool {

	evalCode(e, p, "s", condCode)

	if v, ok := p.stk.Pop(); ok {
		if cond, ok := v.(bool); ok {
			return cond
		}
		panic("for loop condition isn't a boolean expression")
	}
	panic("unexpected code in for statement")
}

func evalCode(e interpreter.Engine, p *Interpreter, name string, code interface{}) {

	old := p.base
	p.base = p.stk.BaseFrame()
	err := e.EvalCode(p, name, code)
	p.base = old
	if err != nil {
		panic(err)
	}
}

// -----------------------------------------------------------------------------
