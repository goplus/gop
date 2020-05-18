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

const (
	bitsCaseArityOperand = (1 << (bitsOpCaseNE - bitsOp)) - 1
)

func execCaseNE(i Instr, ctx *Context) {
	arity := (i >> 16) & bitsCaseArityOperand
	n := len(ctx.data)
	if arity == 1 {
		if ctx.data[n-2] != ctx.data[n-1] {
			ctx.data = ctx.data[:n-1]
			ctx.ip += int(int16(i))
		} else {
			ctx.data = ctx.data[:n-2]
		}
	} else {
		itag := n - 1 - int(arity)
		v := ctx.data[itag]
		for i := int(arity); i > 0; i-- {
			if v == ctx.data[n-i] {
				ctx.data = ctx.data[:itag]
				return
			}
		}
		ctx.data = ctx.data[:itag+1]
		ctx.ip += int(int16(i))
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

func (p *Builder) labelOp(op uint32, l *Label) *Builder {
	if _, ok := p.labels[l]; !ok {
		p.labels[l] = -1
	}
	code := p.code
	l.offs = append(l.offs, code.Len())
	code.data = append(code.data, op)
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
	return p.labelOp(opJmp<<bitsOpShift, l)
}

// JmpIfFalse instr
func (p *Builder) JmpIfFalse(l *Label) *Builder {
	return p.labelOp(opJmpIfFalse<<bitsOpShift, l)
}

// CaseNE instr
func (p *Builder) CaseNE(l *Label, arity int) *Builder {
	return p.labelOp((opCaseNE<<bitsOpShift)|(uint32(arity)<<bitsOpCaseNEShift), l)
}

// Default instr
func (p *Builder) Default() *Builder {
	return p.Pop(1)
}

// -----------------------------------------------------------------------------
