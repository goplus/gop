package exec

import (
	"reflect"

	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

// AddrOperator type.
type AddrOperator Operator

const (
	// OpAddrVal `*addr`
	OpAddrVal = AddrOperator(OpInvalid)
	// OpAddAssign `+=`
	OpAddAssign = AddrOperator(OpAdd)
	// OpSubAssign `-=`
	OpSubAssign = AddrOperator(OpSub)
	// OpMulAssign `*=`
	OpMulAssign = AddrOperator(OpMul)
	// OpDivAssign `/=`
	OpDivAssign = AddrOperator(OpDiv)
	// OpModAssign `%=`
	OpModAssign = AddrOperator(OpMod)

	// OpBitAndAssign '&='
	OpBitAndAssign = AddrOperator(OpBitAnd)
	// OpBitOrAssign '|='
	OpBitOrAssign = AddrOperator(OpBitOr)
	// OpBitXorAssign '^='
	OpBitXorAssign = AddrOperator(OpBitXor)
	// OpBitAndNotAssign '&^='
	OpBitAndNotAssign = AddrOperator(OpBitAndNot)
	// OpBitSHLAssign '<<='
	OpBitSHLAssign = AddrOperator(OpBitSHL)
	// OpBitSHRAssign '>>='
	OpBitSHRAssign = AddrOperator(OpBitSHR)

	// OpAssign `=`
	OpAssign AddrOperator = iota
	// OpInc '++'
	OpInc AddrOperator = iota
	// OpDec '--'
	OpDec
)

// -----------------------------------------------------------------------------

func execOpAddrVal(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = reflect.ValueOf(p.data[n-1]).Elem().Interface()
}

func execOpAssign(i Instr, p *Context) {
	n := len(p.data)
	v := reflect.ValueOf(p.data[n-1])
	reflect.ValueOf(p.data[n-2]).Elem().Set(v)
	p.data = p.data[:n-2]
}

func execOpAddAssign(i Instr, p *Context) {
}

func execOpSubAssign(i Instr, p *Context) {
}

func execOpMulAssign(i Instr, p *Context) {
}

func execOpDivAssign(i Instr, p *Context) {
}

func execOpModAssign(i Instr, p *Context) {
}

func execOpBitAndAssign(i Instr, p *Context) {
}

func execOpBitOrAssign(i Instr, p *Context) {
}

func execOpBitXorAssign(i Instr, p *Context) {
}

func execOpBitAndNotAssign(i Instr, p *Context) {
}

func execOpBitSHLAssign(i Instr, p *Context) {
}

func execOpBitSHRAssign(i Instr, p *Context) {
}

func execOpInc(i Instr, p *Context) {
}

func execOpDec(i Instr, p *Context) {
}

func execAddrOp(i Instr, p *Context) {
	op := (i & bitsOperand) >> bitsKind
	builtinAssignOps[op](i, p)
}

var builtinAssignOps = [...]func(i Instr, p *Context){
	OpAddrVal:         execOpAddrVal,
	OpAssign:          execOpAssign,
	OpAddAssign:       execOpAddAssign,
	OpSubAssign:       execOpSubAssign,
	OpMulAssign:       execOpMulAssign,
	OpDivAssign:       execOpDivAssign,
	OpModAssign:       execOpModAssign,
	OpBitAndAssign:    execOpBitAndAssign,
	OpBitOrAssign:     execOpBitOrAssign,
	OpBitXorAssign:    execOpBitXorAssign,
	OpBitAndNotAssign: execOpBitAndNotAssign,
	OpBitSHLAssign:    execOpBitSHLAssign,
	OpBitSHRAssign:    execOpBitSHRAssign,
	OpInc:             execOpInc,
	OpDec:             execOpDec,
}

// -----------------------------------------------------------------------------

func execAddrVar(i Instr, p *Context) {
	idx := i & bitsOperand
	if idx <= bitsOpVarOperand {
		p.Push(p.AddrVar(idx))
		return
	}
	pp := p.parent
	scope := idx >> bitsOpVarShift
	for scope > 1 {
		pp = pp.parent
		scope--
	}
	p.Push(pp.AddrVar(idx & bitsOpVarOperand))
}

func execLoadVar(i Instr, p *Context) {
	idx := i & bitsOperand
	if idx <= bitsOpVarOperand {
		p.Push(p.GetVar(idx))
		return
	}
	pp := p.parent
	scope := idx >> bitsOpVarShift
	for scope > 1 {
		pp = pp.parent
		scope--
	}
	p.Push(pp.GetVar(idx & bitsOpVarOperand))
}

func execStoreVar(i Instr, p *Context) {
	idx := i & bitsOperand
	if idx <= bitsOpVarOperand {
		p.SetVar(idx, p.pop())
		return
	}
	pp := p.parent
	scope := idx >> bitsOpVarShift
	for scope > 1 {
		pp = pp.parent
		scope--
	}
	pp.SetVar(idx&bitsOpVarOperand, p.pop())
}

// -----------------------------------------------------------------------------

// AddrVar instr
func (p *Builder) AddrVar(scope, idx uint32) *Builder {
	if scope >= (1<<bitsVarScope) || idx > bitsOpVarOperand {
		log.Fatalln("AddrVar failed: invalid scope or variable index -", scope, idx)
	}
	i := (opAddrVar << bitsOpShift) | (scope << bitsOpVarShift) | idx
	p.code.data = append(p.code.data, i)
	return p
}

// LoadVar instr
func (p *Builder) LoadVar(scope, idx uint32) *Builder {
	if scope >= (1<<bitsVarScope) || idx > bitsOpVarOperand {
		log.Fatalln("LoadVar failed: invalid scope or variable index -", scope, idx)
	}
	i := (opLoadVar << bitsOpShift) | (scope << bitsOpVarShift) | idx
	p.code.data = append(p.code.data, i)
	return p
}

// StoreVar instr
func (p *Builder) StoreVar(scope, idx uint32) *Builder {
	if scope >= (1<<bitsVarScope) || idx > bitsOpVarOperand {
		log.Fatalln("LoadVar failed: invalid scope or variable index -", scope, idx)
	}
	i := (opStoreVar << bitsOpShift) | (scope << bitsOpVarShift) | idx
	p.code.data = append(p.code.data, i)
	return p
}

// AddrOp instr
func (p *Builder) AddrOp(kind Kind, op AddrOperator) *Builder {
	i := (int(op) << bitsKind) | int(kind)
	p.code.data = append(p.code.data, (opAddrOp<<bitsOpShift)|uint32(i))
	return p
}

// -----------------------------------------------------------------------------
