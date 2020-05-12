package exec

import (
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

func execJmp(i Instr, ctx *Context) {
	delta := int32(i&bitsOperand) << bitsOp >> bitsOp
	ctx.ip += int(delta)
}

func execJmpIfFalse(i Instr, ctx *Context) {
	if !ctx.Pop().(bool) {
		execJmp(i, ctx)
	}
}

func execCaseNE(i Instr, ctx *Context) {
	n := len(ctx.data)
	if ctx.data[n-2] != ctx.data[n-1] {
		ctx.data = ctx.data[:n-1]
		execJmp(i, ctx)
	} else {
		ctx.data = ctx.data[:n-2]
	}
}

// -----------------------------------------------------------------------------

// Label represents a label.
type Label struct {
	anyUnresolved
	Name string
}

// NewLabel creates a label object.
func NewLabel(name string) *Label {
	return &Label{Name: name}
}

func (p *Builder) resolveLabels() {
	data := p.code.data
	for l, pos := range p.labels {
		if pos < 0 {
			log.Panicln("resolveLabels failed: label is not defined -", l.Name)
		}
		for _, off := range l.offs {
			data[off] |= uint32(pos-(off+1)) & bitsOperand
		}
		l.offs = nil
	}
}

func (p *Builder) labelOp(op int, l *Label) *Builder {
	if _, ok := p.labels[l]; !ok {
		p.labels[l] = -1
	}
	code := p.code
	l.offs = append(l.offs, code.Len())
	code.data = append(code.data, uint32(op)<<bitsOpShift)
	return p
}

// Label defines a label to jmp here.
func (p *Builder) Label(l *Label) *Builder {
	if v, ok := p.labels[l]; ok && v >= 0 {
		log.Panicln("Label failed: label is defined already -", l.Name)
	}
	p.labels[l] = p.code.Len()
	return p
}

// Jmp instr
func (p *Builder) Jmp(l *Label) *Builder {
	return p.labelOp(opJmp, l)
}

// JmpIfFalse instr
func (p *Builder) JmpIfFalse(l *Label) *Builder {
	return p.labelOp(opJmpIfFalse, l)
}

// CaseNE instr
func (p *Builder) CaseNE(l *Label) *Builder {
	return p.labelOp(opCaseNE, l)
}

// Default instr
func (p *Builder) Default() *Builder {
	return p.Pop(1)
}

// -----------------------------------------------------------------------------
