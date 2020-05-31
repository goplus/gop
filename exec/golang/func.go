package golang

import (
	"github.com/qiniu/qlang/v6/exec.spec"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

// Closure instr
func (p *Builder) Closure(fun exec.FuncInfo) *Builder {
	log.Panicln("todo")
	return p
}

// GoClosure instr
func (p *Builder) GoClosure(fun exec.FuncInfo) *Builder {
	log.Panicln("todo")
	return p
}

// CallClosure instr
func (p *Builder) CallClosure(arity int) *Builder {
	log.Panicln("todo")
	return p
}

// CallGoClosure instr
func (p *Builder) CallGoClosure(arity int) *Builder {
	log.Panicln("todo")
	return p
}

// CallFunc instr
func (p *Builder) CallFunc(fun exec.FuncInfo) *Builder {
	log.Panicln("todo")
	return p
}

// CallFuncv instr
func (p *Builder) CallFuncv(fun exec.FuncInfo, arity int) *Builder {
	log.Panicln("todo")
	return p
}

// DefineFunc instr
func (p *Builder) DefineFunc(fun exec.FuncInfo) *Builder {
	log.Panicln("todo")
	return p
}

// Load instr
func (p *Builder) Load(idx int32) *Builder {
	log.Panicln("todo")
	return p
}

// Store instr
func (p *Builder) Store(idx int32) *Builder {
	log.Panicln("todo")
	return p
}

// EndFunc instr
func (p *Builder) EndFunc(fun exec.FuncInfo) *Builder {
	log.Panicln("todo")
	return p
}
