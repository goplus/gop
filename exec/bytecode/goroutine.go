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

// Package bytecode implements a bytecode backend for the Go+ language.
package bytecode

import (
	"github.com/goplus/gop/exec.spec"
)

const (
	bitsOpGoroutine        = bitsOp + 2
	bitsOpGoroutineOperand = bitsOpGoroutine - 1
)

func execGoroutineOp(i Instr, p *Context) {
	end := i & bitsOpGoroutineOperand << bitsOpGoroutine >> bitsOpGoroutine

	data := p.Stack.data

	ctx := newSimpleContext(data)
	ctx.base = p.base
	ctx.code = p.code
	ctx.code.data = ctx.code.data[p.ip : p.ip+int(end)]
	go ctx.Exec(0, ctx.code.Len())
	p.ip += int(end)
}

func (p *Builder) Goroutine(start, end *Label) exec.Instr {
	return &iGoroutine{start: p.labels[start], end: p.labels[end]}
}

type iGoroutine struct {
	start int
	end   int
}

func (c *iGoroutine) Val() interface{} {
	instr := opGoroutineOp<<bitsOpShift | uint32(c.end-c.start)&bitsOpGoroutineOperand
	return Instr(instr)
}
