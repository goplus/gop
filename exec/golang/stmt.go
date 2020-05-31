package golang

import (
	"go/ast"

	"github.com/qiniu/qlang/v6/exec.spec"
	"github.com/qiniu/x/log"
)

// ----------------------------------------------------------------------------

// Return instr
func (p *Builder) Return(n int32) *Builder {
	var results []ast.Expr
	if n > 0 {
		log.Panicln("todo")
	}
	p.rhs.Push(&ast.ReturnStmt{Results: results})
	return p
}

// Label defines a label to jmp here.
func (p *Builder) Label(l exec.Label) *Builder {
	log.Panicln("todo")
	return p
}

// Jmp instr
func (p *Builder) Jmp(l exec.Label) *Builder {
	log.Panicln("todo")
	return p
}

// JmpIf instr
func (p *Builder) JmpIf(zeroOrOne uint32, l exec.Label) *Builder {
	log.Panicln("todo")
	return p
}

// CaseNE instr
func (p *Builder) CaseNE(l exec.Label, arity int) *Builder {
	log.Panicln("todo")
	return p
}

// Default instr
func (p *Builder) Default() *Builder {
	log.Panicln("todo")
	return p
}

// ForPhrase instr
func (p *Builder) ForPhrase(f exec.ForPhrase, key, val exec.Var, hasExecCtx ...bool) *Builder {
	log.Panicln("todo")
	return p
}

// FilterForPhrase instr
func (p *Builder) FilterForPhrase(f exec.ForPhrase) *Builder {
	log.Panicln("todo")
	return p
}

// EndForPhrase instr
func (p *Builder) EndForPhrase(f exec.ForPhrase) *Builder {
	log.Panicln("todo")
	return p
}

// ListComprehension instr
func (p *Builder) ListComprehension(c exec.Comprehension) *Builder {
	log.Panicln("todo")
	return p
}

// MapComprehension instr
func (p *Builder) MapComprehension(c exec.Comprehension) *Builder {
	log.Panicln("todo")
	return p
}

// EndComprehension instr
func (p *Builder) EndComprehension(c exec.Comprehension) *Builder {
	log.Panicln("todo")
	return p
}

// ----------------------------------------------------------------------------
