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

// -----------------------------------------------------------------------------

type theDefer struct {
	next  *theDefer
	state varScope // block state
	n     int      // stack len
	i     Instr
}

func (ctx *Context) execDefers() {
	d := ctx.defers
	ctx.defers = nil
	for d != nil {
		ctx.data = ctx.data[:d.n]
		ctx.varScope = d.state
		execTable[d.i>>bitsOpShift](d.i, ctx)
		d = d.next
	}
}

func execDefer(i Instr, ctx *Context) {
	i = ctx.code.data[ctx.ip]
	ctx.ip++
	ctx.defers = &theDefer{
		next:  ctx.defers,
		state: ctx.varScope,
		n:     len(ctx.data),
		i:     i,
	}
}

// Defer instr
func (p *Builder) Defer() *Builder {
	p.code.data = append(p.code.data, opDefer<<bitsOpShift)
	return p
}

// -----------------------------------------------------------------------------
