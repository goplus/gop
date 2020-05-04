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
	v := reflect.ValueOf(p.data[n-2])
	reflect.ValueOf(p.data[n-1]).Elem().Set(v)
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

func getParentCtx(p *Context, idx Address) *Context {
	pp := p.parent
	scope := uint32(idx) >> bitsOpVarShift
	for scope > 1 {
		pp = pp.parent
		scope--
	}
	return pp
}

// FastAddrVar func.
func (p *Context) FastAddrVar(idx Address) interface{} {
	if idx <= bitsOpVarOperand {
		return p.addrVar(uint32(idx))
	}
	return getParentCtx(p, idx).addrVar(uint32(idx) & bitsOpVarOperand)
}

// FastGetVar func.
func (p *Context) FastGetVar(idx Address) interface{} {
	if idx <= bitsOpVarOperand {
		return p.getVar(uint32(idx))
	}
	return getParentCtx(p, idx).getVar(uint32(idx) & bitsOpVarOperand)
}

// FastSetVar func.
func (p *Context) FastSetVar(idx Address, v interface{}) {
	if idx <= bitsOpVarOperand {
		p.setVar(uint32(idx), v)
		return
	}
	getParentCtx(p, idx).setVar(uint32(idx)&bitsOpVarOperand, v)
}

func execAddrVar(i Instr, p *Context) {
	idx := i & bitsOperand
	if idx <= bitsOpVarOperand {
		p.Push(p.addrVar(idx))
		return
	}
	p.Push(getParentCtx(p, Address(idx)).addrVar(idx & bitsOpVarOperand))
}

func execLoadVar(i Instr, p *Context) {
	idx := i & bitsOperand
	if idx <= bitsOpVarOperand {
		p.Push(p.getVar(idx))
		return
	}
	p.Push(getParentCtx(p, Address(idx)).getVar(idx & bitsOpVarOperand))
}

func execStoreVar(i Instr, p *Context) {
	idx := i & bitsOperand
	if idx <= bitsOpVarOperand {
		p.setVar(idx, p.Pop())
		return
	}
	getParentCtx(p, Address(idx)).setVar(idx&bitsOpVarOperand, p.Pop())
}

// -----------------------------------------------------------------------------

// Address represents a variable address.
type Address uint32

// Scope returns the scope of this variable. Zero means it is defined in current block.
func (p Address) Scope() int {
	return int(p >> bitsVarScope)
}

// MakeAddr creates a variable address.
func MakeAddr(scope, idx uint32) Address {
	if scope >= (1<<bitsVarScope) || idx > bitsOpVarOperand {
		log.Fatalln("MakeAddr failed: invalid scope or variable index -", scope, idx)
	}
	return Address((scope << bitsOpVarShift) | idx)
}

// AddrVar instr
func (p *Builder) AddrVar(addr Address) *Builder {
	p.code.data = append(p.code.data, (opAddrVar<<bitsOpShift)|uint32(addr))
	return p
}

// LoadVar instr
func (p *Builder) LoadVar(addr Address) *Builder {
	p.code.data = append(p.code.data, (opLoadVar<<bitsOpShift)|uint32(addr))
	return p
}

// StoreVar instr
func (p *Builder) StoreVar(addr Address) *Builder {
	p.code.data = append(p.code.data, (opStoreVar<<bitsOpShift)|uint32(addr))
	return p
}

// AddrOp instr
func (p *Builder) AddrOp(kind Kind, op AddrOperator) *Builder {
	i := (int(op) << bitsKind) | int(kind)
	p.code.data = append(p.code.data, (opAddrOp<<bitsOpShift)|uint32(i))
	return p
}

// -----------------------------------------------------------------------------
