package bytecode

import (
	"reflect"
)

// Send instr
func (p *Builder) Send() *Builder {
	p.code.data = append(p.code.data, opSend<<bitsOpShift)
	return p
}

// Recv instr
func (p *Builder) Recv() *Builder {
	p.code.data = append(p.code.data, opRecv<<bitsOpShift)
	return p
}

func execSend(i Instr, ctx *Context) {
	args := ctx.GetArgs(2)
	v := reflect.ValueOf(args[0])
	v.Send(reflect.ValueOf(args[1]))
	ctx.PopN(2)
}

func execRecv(i Instr, p *Context) {
	n := len(p.data)
	v := reflect.ValueOf(p.data[n-1])
	x, ok := v.Recv()
	if ok {
		p.Ret(1, x.Interface())
	}
}
