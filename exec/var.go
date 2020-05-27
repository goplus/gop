/*
 Copyright 2020 Qiniu Cloud (qiniu.com)

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

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
	OpAddrVal = AddrOperator(0)
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
	OpInc
	// OpDec '--'
	OpDec
)

// AddrOperatorInfo represents an addr-operator information.
type AddrOperatorInfo struct {
	Lit      string
	InFirst  uint64 // first argument supported types.
	InSecond uint64 // second argument supported types. It may have SameAsFirst flag.
}

var addropInfos = [...]AddrOperatorInfo{
	OpAddAssign:       {"+=", bitsAllNumber | bitString, bitSameAsFirst},
	OpSubAssign:       {"-=", bitsAllNumber, bitSameAsFirst},
	OpMulAssign:       {"*=", bitsAllNumber, bitSameAsFirst},
	OpDivAssign:       {"/=", bitsAllNumber, bitSameAsFirst},
	OpModAssign:       {"%=", bitsAllIntUint, bitSameAsFirst},
	OpBitAndAssign:    {"&=", bitsAllIntUint, bitSameAsFirst},
	OpBitOrAssign:     {"|=", bitsAllIntUint, bitSameAsFirst},
	OpBitXorAssign:    {"^=", bitsAllIntUint, bitSameAsFirst},
	OpBitAndNotAssign: {"&^=", bitsAllIntUint, bitSameAsFirst},
	OpBitSHLAssign:    {"<<=", bitsAllIntUint, bitsAllIntUint},
	OpBitSHRAssign:    {">>=", bitsAllIntUint, bitsAllIntUint},
	OpInc:             {"++", bitsAllNumber, bitNone},
	OpDec:             {"--", bitsAllNumber, bitNone},
}

// GetInfo returns the information of this operator.
func (op AddrOperator) GetInfo() *AddrOperatorInfo {
	return &addropInfos[op]
}

func (op AddrOperator) String() string {
	switch op {
	case OpAddrVal:
		return "*"
	case OpAssign:
		return "="
	default:
		return addropInfos[op].Lit
	}
}

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

func execAddrOp(i Instr, p *Context) {
	op := i & bitsOperand
	if fn := builtinAddrOps[int(op)]; fn != nil {
		fn(0, p)
		return
	}
	switch AddrOperator(op >> bitsKind) {
	case OpAssign:
		execOpAssign(0, p)
	case OpAddrVal:
		execOpAddrVal(0, p)
	default:
		log.Panicln("execAddrOp: invalid instr -", i)
	}
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

// GetVar func.
func (p *Context) GetVar(x *Var) interface{} {
	if x.isGlobal() {
		return p.getVar(x.idx)
	}
	panic("variable not defined, or not a global variable")
}

// SetVar func.
func (p *Context) SetVar(x *Var, v interface{}) {
	if x.isGlobal() {
		p.setVar(x.idx, v)
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

// -----------------------------------------------------------------------------

// Var represents a variable.
type Var struct {
	typ       reflect.Type
	name      string
	nestDepth uint32
	idx       uint32
}

// NewVar creates a variable instance.
func NewVar(typ reflect.Type, name string) *Var {
	return &Var{typ: typ, name: "Q" + name, idx: ^uint32(0)}
}

func (p *Var) isGlobal() bool {
	return p.idx <= bitsOpVarOperand && p.nestDepth == 0
}

// Type returns variable's type.
func (p *Var) Type() reflect.Type {
	return p.typ
}

// Name returns variable's name.
func (p *Var) Name() string {
	return p.name[1:]
}

// IsUnnamedOut returns if variable unnamed or not.
func (p *Var) IsUnnamedOut() bool {
	c := p.name[0]
	return c >= '0' && c <= '9'
}

// SetAddr sets a variable address.
func (p *Var) SetAddr(nestDepth, idx uint32) {
	if p.idx <= bitsOpVarOperand {
		log.Panicln("Var.setAddr failed: the variable is defined already -", p.name[1:])
	}
	p.nestDepth, p.idx = nestDepth, idx
}

// -----------------------------------------------------------------------------

type varManager struct {
	vlist     []*Var
	tcache    reflect.Type
	nestDepth uint32
}

func newVarManager(vars ...*Var) *varManager {
	return &varManager{vlist: vars}
}

func (p *varManager) addVars(vars ...*Var) {
	n := len(p.vlist)
	nestDepth := p.nestDepth
	for i, v := range vars {
		v.SetAddr(nestDepth, uint32(i+n))
		log.Debug("DefineVar:", v.Name(), "-", nestDepth)
	}
	p.vlist = append(p.vlist, vars...)
}

type blockCtx struct {
	varManager
	parent *varManager
}

func newBlockCtx(nestDepth uint32, parent *varManager) *blockCtx {
	return &blockCtx{varManager: varManager{nestDepth: nestDepth}, parent: parent}
}

// -----------------------------------------------------------------------------

func (p *Context) getNestDepth() (nestDepth uint32) {
	for {
		if p = p.parent; p == nil {
			return
		}
		nestDepth++
	}
}

// InCurrentCtx returns if a variable is in current context or not.
func (p *Builder) InCurrentCtx(v *Var) bool {
	return p.nestDepth == v.nestDepth
}

// DefineVar defines variables.
func (p *Builder) DefineVar(vars ...*Var) *Builder {
	p.addVars(vars...)
	return p
}

// LoadVar instr
func (p *Builder) LoadVar(v *Var) *Builder {
	p.loadVar(makeAddr(p.nestDepth-v.nestDepth, v.idx))
	return p
}

// StoreVar instr
func (p *Builder) StoreVar(v *Var) *Builder {
	p.storeVar(makeAddr(p.nestDepth-v.nestDepth, v.idx))
	return p
}

// AddrVar instr
func (p *Builder) AddrVar(v *Var) *Builder {
	p.addrVar(makeAddr(p.nestDepth-v.nestDepth, v.idx))
	return p
}

// AddrOp instr
func (p *Builder) AddrOp(kind Kind, op AddrOperator) *Builder {
	i := (int(op) << bitsKind) | int(kind)
	p.code.data = append(p.code.data, (opAddrOp<<bitsOpShift)|uint32(i))
	return p
}

// CallAddrOp calls AddrOp
func CallAddrOp(kind Kind, op AddrOperator, data ...interface{}) {
	if fn := builtinAddrOps[(int(op)<<bitsKind)|int(kind)]; fn != nil {
		ctx := newSimpleContext(data)
		fn(0, ctx)
		return
	}
	panic("CallAddrOp: invalid addrOp")
}

// -----------------------------------------------------------------------------
