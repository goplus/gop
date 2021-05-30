// Package hex provide Go+ "encoding/hex" package, as "encoding/hex" package in Go.
package hex

import (
	hex "encoding/hex"
	io "io"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func execDecode(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := hex.Decode(args[0].([]byte), args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execDecodeString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := hex.DecodeString(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execDecodedLen(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := hex.DecodedLen(args[0].(int))
	p.Ret(1, ret0)
}

func execDump(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := hex.Dump(args[0].([]byte))
	p.Ret(1, ret0)
}

func toType0(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execDumper(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := hex.Dumper(toType0(args[0]))
	p.Ret(1, ret0)
}

func execEncode(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := hex.Encode(args[0].([]byte), args[1].([]byte))
	p.Ret(2, ret0)
}

func execEncodeToString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := hex.EncodeToString(args[0].([]byte))
	p.Ret(1, ret0)
}

func execEncodedLen(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := hex.EncodedLen(args[0].(int))
	p.Ret(1, ret0)
}

func execmInvalidByteErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(hex.InvalidByteError).Error()
	p.Ret(1, ret0)
}

func toType1(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execNewDecoder(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := hex.NewDecoder(toType1(args[0]))
	p.Ret(1, ret0)
}

func execNewEncoder(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := hex.NewEncoder(toType0(args[0]))
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("encoding/hex")

func init() {
	I.RegisterFuncs(
		I.Func("Decode", hex.Decode, execDecode),
		I.Func("DecodeString", hex.DecodeString, execDecodeString),
		I.Func("DecodedLen", hex.DecodedLen, execDecodedLen),
		I.Func("Dump", hex.Dump, execDump),
		I.Func("Dumper", hex.Dumper, execDumper),
		I.Func("Encode", hex.Encode, execEncode),
		I.Func("EncodeToString", hex.EncodeToString, execEncodeToString),
		I.Func("EncodedLen", hex.EncodedLen, execEncodedLen),
		I.Func("(InvalidByteError).Error", (hex.InvalidByteError).Error, execmInvalidByteErrorError),
		I.Func("NewDecoder", hex.NewDecoder, execNewDecoder),
		I.Func("NewEncoder", hex.NewEncoder, execNewEncoder),
	)
	I.RegisterVars(
		I.Var("ErrLength", &hex.ErrLength),
	)
	I.RegisterTypes(
		I.Type("InvalidByteError", reflect.TypeOf((*hex.InvalidByteError)(nil)).Elem()),
	)
}
