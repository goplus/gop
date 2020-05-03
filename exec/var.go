package exec

import (
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

func execLoadVar(i Instr, p *Context) {
	idx := i & bitsOperand
	if idx <= bitsOpVarOperand {
		p.Push(p.vars[idx])
	}
	scope := idx >> bitsOpVarShift
	for scope > 0 {
		p = p.parent
		scope--
	}
	p.Push(p.vars[idx&bitsOpVarOperand])
}

func execStoreVar(i Instr, p *Context) {
	idx := i & bitsOperand
	if idx <= bitsOpVarOperand {
		p.vars[idx] = p.pop()
	}
	scope := idx >> bitsOpVarShift
	for scope > 0 {
		p = p.parent
		scope--
	}
	p.vars[idx&bitsOpVarOperand] = p.pop()
}

// -----------------------------------------------------------------------------

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

// -----------------------------------------------------------------------------
