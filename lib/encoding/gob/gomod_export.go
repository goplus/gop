// Package gob provide Go+ "encoding/gob" package, as "encoding/gob" package in Go.
package gob

import (
	gob "encoding/gob"
	io "io"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func execmDecoderDecode(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*gob.Decoder).Decode(args[1])
	p.Ret(2, ret0)
}

func execmDecoderDecodeValue(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*gob.Decoder).DecodeValue(args[1].(reflect.Value))
	p.Ret(2, ret0)
}

func execmEncoderEncode(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*gob.Encoder).Encode(args[1])
	p.Ret(2, ret0)
}

func execmEncoderEncodeValue(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*gob.Encoder).EncodeValue(args[1].(reflect.Value))
	p.Ret(2, ret0)
}

func execiGobDecoderGobDecode(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(gob.GobDecoder).GobDecode(args[1].([]byte))
	p.Ret(2, ret0)
}

func execiGobEncoderGobEncode(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(gob.GobEncoder).GobEncode()
	p.Ret(1, ret0, ret1)
}

func toType0(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execNewDecoder(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := gob.NewDecoder(toType0(args[0]))
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
	ret0 := gob.NewEncoder(toType1(args[0]))
	p.Ret(1, ret0)
}

func execRegister(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	gob.Register(args[0])
	p.PopN(1)
}

func execRegisterName(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	gob.RegisterName(args[0].(string), args[1])
	p.PopN(2)
}

// I is a Go package instance.
var I = gop.NewGoPackage("encoding/gob")

func init() {
	I.RegisterFuncs(
		I.Func("(*Decoder).Decode", (*gob.Decoder).Decode, execmDecoderDecode),
		I.Func("(*Decoder).DecodeValue", (*gob.Decoder).DecodeValue, execmDecoderDecodeValue),
		I.Func("(*Encoder).Encode", (*gob.Encoder).Encode, execmEncoderEncode),
		I.Func("(*Encoder).EncodeValue", (*gob.Encoder).EncodeValue, execmEncoderEncodeValue),
		I.Func("(GobDecoder).GobDecode", (gob.GobDecoder).GobDecode, execiGobDecoderGobDecode),
		I.Func("(GobEncoder).GobEncode", (gob.GobEncoder).GobEncode, execiGobEncoderGobEncode),
		I.Func("NewDecoder", gob.NewDecoder, execNewDecoder),
		I.Func("NewEncoder", gob.NewEncoder, execNewEncoder),
		I.Func("Register", gob.Register, execRegister),
		I.Func("RegisterName", gob.RegisterName, execRegisterName),
	)
	I.RegisterTypes(
		I.Type("CommonType", reflect.TypeOf((*gob.CommonType)(nil)).Elem()),
		I.Type("Decoder", reflect.TypeOf((*gob.Decoder)(nil)).Elem()),
		I.Type("Encoder", reflect.TypeOf((*gob.Encoder)(nil)).Elem()),
		I.Type("GobDecoder", reflect.TypeOf((*gob.GobDecoder)(nil)).Elem()),
		I.Type("GobEncoder", reflect.TypeOf((*gob.GobEncoder)(nil)).Elem()),
	)
}
