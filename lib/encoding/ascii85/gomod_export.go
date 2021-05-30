// Package ascii85 provide Go+ "encoding/ascii85" package, as "encoding/ascii85" package in Go.
package ascii85

import (
	ascii85 "encoding/ascii85"
	io "io"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func execmCorruptInputErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(ascii85.CorruptInputError).Error()
	p.Ret(1, ret0)
}

func execDecode(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1, ret2 := ascii85.Decode(args[0].([]byte), args[1].([]byte), args[2].(bool))
	p.Ret(3, ret0, ret1, ret2)
}

func execEncode(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := ascii85.Encode(args[0].([]byte), args[1].([]byte))
	p.Ret(2, ret0)
}

func execMaxEncodedLen(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := ascii85.MaxEncodedLen(args[0].(int))
	p.Ret(1, ret0)
}

func toType0(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execNewDecoder(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := ascii85.NewDecoder(toType0(args[0]))
	p.Ret(1, ret0)
}

func toType1(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execNewEncoder(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := ascii85.NewEncoder(toType1(args[0]))
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("encoding/ascii85")

func init() {
	I.RegisterFuncs(
		I.Func("(CorruptInputError).Error", (ascii85.CorruptInputError).Error, execmCorruptInputErrorError),
		I.Func("Decode", ascii85.Decode, execDecode),
		I.Func("Encode", ascii85.Encode, execEncode),
		I.Func("MaxEncodedLen", ascii85.MaxEncodedLen, execMaxEncodedLen),
		I.Func("NewDecoder", ascii85.NewDecoder, execNewDecoder),
		I.Func("NewEncoder", ascii85.NewEncoder, execNewEncoder),
	)
	I.RegisterTypes(
		I.Type("CorruptInputError", reflect.TypeOf((*ascii85.CorruptInputError)(nil)).Elem()),
	)
}
