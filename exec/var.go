package exec

import (
	"reflect"
	"strings"

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

func getParentCtx(p *Context, idx tAddress) *Context {
	pp := p.parent
	scope := uint32(idx) >> bitsOpVarShift
	for scope > 1 {
		pp = pp.parent
		scope--
	}
	return pp
}

/*
// FastAddrVar func.
func (p *Context) fastAddrVar(idx tAddress) interface{} {
	if idx <= bitsOpVarOperand {
		return p.addrVar(uint32(idx))
	}
	return getParentCtx(p, idx).addrVar(uint32(idx) & bitsOpVarOperand)
}

// FastGetVar func.
func (p *Context) fastGetVar(idx tAddress) interface{} {
	if idx <= bitsOpVarOperand {
		return p.getVar(uint32(idx))
	}
	return getParentCtx(p, idx).getVar(uint32(idx) & bitsOpVarOperand)
}

// FastSetVar func.
func (p *Context) fastSetVar(idx tAddress, v interface{}) {
	if idx <= bitsOpVarOperand {
		p.setVar(uint32(idx), v)
		return
	}
	getParentCtx(p, idx).setVar(uint32(idx)&bitsOpVarOperand, v)
}
*/

// GetVar func.
func (p *Context) GetVar(x *Var) interface{} {
	if x.isGlobal() {
		return p.getVar(uint32(x.idx))
	}
	panic("variable not defined, or not a global variable")
}

// SetVar func.
func (p *Context) SetVar(x *Var, v interface{}) {
	if x.isGlobal() {
		p.setVar(uint32(x.idx), v)
		return
	}
	panic("variable not defined, or not a global variable")
}

func execAddrVar(i Instr, p *Context) {
	idx := i & bitsOperand
	if idx <= bitsOpVarOperand {
		p.Push(p.addrVar(idx))
		return
	}
	p.Push(getParentCtx(p, tAddress(idx)).addrVar(idx & bitsOpVarOperand))
}

func execLoadVar(i Instr, p *Context) {
	idx := i & bitsOperand
	if idx <= bitsOpVarOperand {
		p.Push(p.getVar(idx))
		return
	}
	p.Push(getParentCtx(p, tAddress(idx)).getVar(idx & bitsOpVarOperand))
}

func execStoreVar(i Instr, p *Context) {
	idx := i & bitsOperand
	if idx <= bitsOpVarOperand {
		p.setVar(idx, p.Pop())
		return
	}
	getParentCtx(p, tAddress(idx)).setVar(idx&bitsOpVarOperand, p.Pop())
}

// -----------------------------------------------------------------------------

// Address represents a variable address.
type tAddress uint32

// MakeAddr creates a variable address.
func makeAddr(scope, idx uint32) tAddress {
	if scope >= (1<<bitsVarScope) || idx > bitsOpVarOperand {
		log.Panicln("MakeAddr failed: invalid scope or variable index -", scope, idx)
	}
	return tAddress((scope << bitsOpVarShift) | idx)
}

// AddrVar instr
func (p *Builder) addrVar(addr tAddress) *Builder {
	p.code.data = append(p.code.data, (opAddrVar<<bitsOpShift)|uint32(addr))
	return p
}

// LoadVar instr
func (p *Builder) loadVar(addr tAddress) *Builder {
	p.code.data = append(p.code.data, (opLoadVar<<bitsOpShift)|uint32(addr))
	return p
}

// StoreVar instr
func (p *Builder) storeVar(addr tAddress) *Builder {
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

// Var represents a variable.
type Var struct {
	Type      reflect.Type
	Name      string
	NestDepth uint32
	idx       uint32
}

// NewVar creates a variable instance.
func NewVar(typ reflect.Type, name string) *Var {
	return &Var{Type: typ, Name: strings.ToTitle(name), idx: ^uint32(0)}
}

func (p *Var) isGlobal() bool {
	return p.idx <= bitsOpVarOperand && p.NestDepth == 0
}

// SetAddr sets a variable address.
func (p *Var) SetAddr(nestDepth, idx uint32) {
	if p.idx <= bitsOpVarOperand {
		log.Fatalln("Var.setAddr failed: the variable is defined already -", p.Name)
	}
	p.NestDepth, p.idx = nestDepth, idx
}

func (p *Context) getNestDepth() (nestDepth uint32) {
	for {
		if p = p.parent; p == nil {
			return
		}
		nestDepth++
	}
}

// SetNestDepth sets function nest depth.
func (p *Builder) SetNestDepth(nestDepth uint32) *Builder {
	p.NestDepth = nestDepth
	return p
}

// DefineVar defines all local variable of a function (closure).
func (p *Builder) DefineVar(vars ...*Var) *Builder {
	for i, v := range vars {
		v.SetAddr(p.NestDepth, uint32(i))
	}
	return p
}

// AddrVar instr
func (p *Builder) AddrVar(v *Var) *Builder {
	p.addrVar(makeAddr(p.NestDepth-v.NestDepth, v.idx))
	return p
}

// LoadVar instr
func (p *Builder) LoadVar(v *Var) *Builder {
	p.loadVar(makeAddr(p.NestDepth-v.NestDepth, v.idx))
	return p
}

// StoreVar instr
func (p *Builder) StoreVar(v *Var) *Builder {
	p.storeVar(makeAddr(p.NestDepth-v.NestDepth, v.idx))
	return p
}

// -----------------------------------------------------------------------------
