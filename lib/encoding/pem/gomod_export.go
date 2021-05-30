// Package pem provide Go+ "encoding/pem" package, as "encoding/pem" package in Go.
package pem

import (
	pem "encoding/pem"
	io "io"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func execDecode(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := pem.Decode(args[0].([]byte))
	p.Ret(1, ret0, ret1)
}

func toType0(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execEncode(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := pem.Encode(toType0(args[0]), args[1].(*pem.Block))
	p.Ret(2, ret0)
}

func execEncodeToMemory(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := pem.EncodeToMemory(args[0].(*pem.Block))
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("encoding/pem")

func init() {
	I.RegisterFuncs(
		I.Func("Decode", pem.Decode, execDecode),
		I.Func("Encode", pem.Encode, execEncode),
		I.Func("EncodeToMemory", pem.EncodeToMemory, execEncodeToMemory),
	)
	I.RegisterTypes(
		I.Type("Block", reflect.TypeOf((*pem.Block)(nil)).Elem()),
	)
}
