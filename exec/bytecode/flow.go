/*
 Copyright 2020 The GoPlus Authors (goplus.org)

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

package bytecode

import (
	"github.com/goplus/gop/exec.spec"
	"github.com/qiniu/x/errors"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

const (
	bitsOpJmp              = bitsOp + 4
	bitsOpJmpIfCond        = bitsInstr - bitsOpJmp
	bitsOpJmpBoolCondFlag  = 1 << bitsOpJmpIfCond
	bitsOpJmpPtrCondFlag   = 2 << bitsOpJmpIfCond
	bitsOpJmpNoPopMaskFlag = 4 << bitsOpJmpIfCond
	bitsOpJmpOperand       = bitsOpJmpBoolCondFlag - 1
)

func execJmp(i Instr, ctx *Context) {
	delta := int32(i&bitsOpJmpOperand) << bitsOpJmp >> bitsOpJmp
	ctx.ip += int(delta)
}

func execWrapIfErr(i Instr, ctx *Context) {
	v := ctx.Pop()
	if v == nil {
		execJmp(i, ctx)
		return
	}
	ctx.PopN(1)
}

func execJmpIf(i Instr, ctx *Context) {
	var v interface{}
	if i&bitsOpJmpNoPopMaskFlag != 0 {
		v = ctx.Get(-1)
	} else {
		v = ctx.Pop()
	}
	if (i & bitsOpJmpPtrCondFlag) != 0 {
		if (i & bitsOpJmpBoolCondFlag) == 0 {
			if v != nil {
				return
			}
		} else if v == nil {
			return
		}
	} else {
		cond := v.(bool)
		if (i & bitsOpJmpBoolCondFlag) == 0 {
			if cond {
				return
			}
		} else if !cond {
			return
		}
	}
	execJmp(i, ctx)
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
	name string
}

// NewLabel creates a label object.
func NewLabel(name string) *Label {
	return &Label{name: name}
}

// Name returns the label name.
func (p *Label) Name() string {
	return p.name
}

func (p *Builder) resolveLabels() {
	data := p.code.data
	for l, pos := range p.labels {
		if pos < 0 {
			log.Panicln("resolveLabels failed: label is not defined -", l.name)
		}
		for _, off := range l.offs {
			if (data[off] >> bitsOpShift) == opCaseNE {
				data[off] |= uint32(int16(pos - (off + 1)))
			} else {
				data[off] |= uint32(pos-(off+1)) & bitsOpJmpOperand
			}
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
		log.Panicln("Label failed: label is defined already -", l.name)
	}
	p.labels[l] = p.code.Len()
	return p
}

// Jmp instr
func (p *Builder) Jmp(l *Label) *Builder {
	return p.labelOp(opJmp<<bitsOpShift, l)
}

// JmpIf instr
func (p *Builder) JmpIf(cond exec.JmpCondFlag, l *Label) *Builder {
	return p.labelOp((opJmpIf<<bitsOpShift)|(uint32(cond)<<bitsOpJmpIfCond), l)
}

// WrapIfErr instr
func (p *Builder) WrapIfErr(nret int, l *Label) *Builder {
	return p.labelOp(opWrapIfErr<<bitsOpShift, l)
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

type errWrap struct {
	retErr exec.Var
	frame  *errors.Frame
	nret   int
	narg   int
}

func execErrWrap(i Instr, ctx *Context) {
	val := ctx.Pop()
	if val == nil { // success
		return
	}
	idx := i & bitsOperand
	ew := ctx.code.errWraps[idx]
	frame := *ew.frame
	frame.Err = val.(error)
	if ew.narg > 0 {
		frame.Args = make([]interface{}, ew.narg)
		for i := 0; i < ew.narg; i++ {
			frame.Args[i] = ctx.data[ctx.base+i-ew.narg]
		}
	}
	if ew.retErr == nil {
		panic(&frame)
	}
	ctx.setVar(ew.retErr.(*Var).idx, &frame)
	ctx.ip = ipInvalid
}

// ErrWrap instr
func (p *Builder) ErrWrap(nret int, retErr exec.Var, frame *errors.Frame, narg int) *Builder {
	code := p.code
	idx := len(code.errWraps)
	code.errWraps = append(
		code.errWraps,
		errWrap{nret: nret, retErr: retErr, frame: frame, narg: narg},
	)
	code.data = append(code.data, (opErrWrap<<bitsOpShift)|uint32(idx))
	return p
}

// -----------------------------------------------------------------------------
