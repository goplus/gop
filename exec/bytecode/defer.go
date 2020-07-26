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
)

type theDefer struct {
	next *theDefer
	Stack
	start int
	end   int
}

const (
	bitsOpDefer        = bitsOp + 2
	bitsOpDeferOperand = bitsOpDefer - 1
)

func execDefers(ctx *Context) {
	d := ctx.defers
	if d == nil {
		return
	}

	ctx.defers = nil
	for {
		ctx.Stack = d.Stack
		ctx.Exec(d.start, d.end)
		d = d.next
		if d == nil {
			break
		}
	}
}

func execDeferOp(i Instr, p *Context) {
	end := i & bitsOpDeferOperand << bitsOpDefer >> bitsOpDefer
	def := &theDefer{
		Stack: p.Stack,
		next:  p.defers,
		start: p.ip,
	}

	p.ip += int(end)
	def.end = p.ip
	p.defers = def
}

// Defer instr
func (p *Builder) Defer(start, end *Label) exec.Instr {
	return &iDefer{start: p.labels[start], end: p.labels[end]}
}

type iDefer struct {
	start int
	end   int
}

func (c *iDefer) Val() interface{} {
	instr := opDeferOp<<bitsOpShift | uint32(c.end-c.start)&bitsOpDeferOperand
	return Instr(instr)
}
