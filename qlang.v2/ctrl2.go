package qlang

import (
	"qlang.io/exec.v2"
	"qiniupkg.com/text/tpl.v1/interpreter.util"
)

// -----------------------------------------------------------------------------

func (p *Compiler) Or(e interpreter.Engine) {

	reserved := p.code.Reserve()
	expr, _ := p.gstk.Pop()
	if err := e.EvalCode(p, "term4", expr); err != nil {
		panic(err)
	}
	reserved.Set(exec.Or(p.code.Len() - reserved.Next()))
}

func (p *Compiler) And(e interpreter.Engine) {

	reserved := p.code.Reserve()
	expr, _ := p.gstk.Pop()
	if err := e.EvalCode(p, "term3", expr); err != nil {
		panic(err)
	}
	reserved.Set(exec.And(p.code.Len() - reserved.Next()))
}

func (p *Compiler) If(e interpreter.Engine) {

	var elseCode interface{}

	elseArity := p.popArity()
	if elseArity == 1 {
		elseCode, _ = p.gstk.Pop()
	}

	condArity := p.popArity()
	condArity++

	ifbr := p.gstk.PopNArgs(condArity << 1) // 2 * (condArity + 1)
	p.doIf(e, ifbr, elseCode, condArity)
}

func (p *Compiler) doIf(e interpreter.Engine, ifbr []interface{}, elseCode interface{}, condArity int) {

	reserved2 := make([]exec.ReservedInstr, condArity)

	for i := 0; i < condArity; i++ {
		condCode := ifbr[i<<1]
		if err := e.EvalCode(p, "expr", condCode); err != nil {
			panic(err)
		}
		p.CodeLine(condCode)
		reserved1 := p.code.Reserve()
		bodyCode := ifbr[(i<<1)+1]
		bctx := evalDocCode(e, p, bodyCode)
		bctx.MergeTo(&p.bctx)
		reserved2[i] = p.code.Reserve()
		reserved1.Set(exec.JmpIfFalse(reserved2[i].Delta(reserved1)))
	}

	bctx := evalDocCode(e, p, elseCode)
	bctx.MergeTo(&p.bctx)

	end := p.code.Len()
	for i := 0; i < condArity; i++ {
		reserved2[i].Set(exec.Jmp(end - reserved2[i].Next()))
	}
}

func (p *Compiler) Switch(e interpreter.Engine) {

	var defaultCode interface{}

	defaultArity := p.popArity()
	if defaultArity == 1 {
		defaultCode, _ = p.gstk.Pop()
	}

	caseArity := p.popArity()
	casebr := p.gstk.PopNArgs(caseArity << 1) // 2 * caseArity
	switchCode, _ := p.gstk.Pop()

	old := p.bctx
	p.bctx = blockCtx{}
	if switchCode == nil {
		p.doIf(e, casebr, defaultCode, caseArity)
		p.bctx.MergeSw(&old, p.code.Len())
		return
	}

	reserved2 := make([]exec.ReservedInstr, caseArity)
	if err := e.EvalCode(p, "expr", switchCode); err != nil {
		panic(err)
	}
	p.CodeLine(switchCode)
	for i := 0; i < caseArity; i++ {
		caseCode := casebr[i<<1]
		if err := e.EvalCode(p, "expr", caseCode); err != nil {
			panic(err)
		}
		p.CodeLine(caseCode)
		reserved1 := p.code.Reserve()
		bodyCode := casebr[(i<<1)+1]
		bctx := evalDocCode(e, p, bodyCode)
		bctx.MergeTo(&p.bctx)
		reserved2[i] = p.code.Reserve()
		reserved1.Set(exec.Case(reserved2[i].Delta(reserved1)))
	}

	p.code.Block(exec.Default)
	bctx := evalDocCode(e, p, defaultCode)
	bctx.MergeTo(&p.bctx)

	end := p.code.Len()
	for i := 0; i < caseArity; i++ {
		reserved2[i].Set(exec.Jmp(end - reserved2[i].Next()))
	}
	p.bctx.MergeSw(&old, end)
}

func (p *Compiler) For(e interpreter.Engine) {

	bodyCode, _ := p.gstk.Pop()

	var condCode interface{}
	var stepCode interface{}

	arity := p.popArity()
	switch arity {
	case 1, 3:
		forCode := p.gstk.PopNArgs(arity)
		if arity == 3 {
			initCode := forCode[0]
			if initCode != nil {
				if err := e.EvalCode(p, "s", initCode); err != nil {
					panic(err)
				}
			}
			stepCode = forCode[2]
		}
		condCode = forCode[arity>>1]
	case 0:
	default:
		panic("for statement is in illegal form")
	}

	loop := p.code.Len()

	if condCode != nil {
		if err := e.EvalCode(p, "expr", condCode); err != nil {
			panic(err)
		}
		p.CodeLine(condCode)
		reserved := p.code.Reserve()
		defer func() {
			reserved.Set(exec.JmpIfFalse(p.code.Len() - reserved.Next()))
		}()
	}

	bctx := evalDocCode(e, p, bodyCode)
	if stepCode != nil {
		bctx.conts.JmpTo(p.code.Len())
		if err := e.EvalCode(p, "s", stepCode); err != nil {
			panic(err)
		}
	} else {
		bctx.conts.JmpTo(loop)
	}

	p.code.Block(exec.Jmp(loop - (p.code.Len() + 1)))
	bctx.brks.JmpTo(p.code.Len())
}

func (p *Compiler) Break() {

	instr := p.code.Reserve()
	p.bctx.brks = &instrNode{
		prev:  p.bctx.brks,
		instr: instr,
	}
}

func (p *Compiler) Continue() {

	instr := p.code.Reserve()
	p.bctx.conts = &instrNode{
		prev:  p.bctx.conts,
		instr: instr,
	}
}

func evalDocCode(e interpreter.Engine, p *Compiler, code interface{}) (bctx blockCtx) {

	if code == nil {
		return
	}

	old := p.bctx
	p.bctx = bctx

	err := e.EvalCode(p, "doc", code)

	bctx = p.bctx
	p.bctx = old

	if err != nil {
		panic(err)
	}
	return
}

// -----------------------------------------------------------------------------

