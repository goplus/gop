// Package base32 provide Go+ "encoding/base32" package, as "encoding/base32" package in Go.
package base32

import (
	base32 "encoding/base32"
	io "io"
	reflect "reflect"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execmCorruptInputErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(base32.CorruptInputError).Error()
	p.Ret(1, ret0)
}

func execmEncodingWithPadding(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(base32.Encoding).WithPadding(args[1].(rune))
	p.Ret(2, ret0)
}

func execmEncodingEncode(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*base32.Encoding).Encode(args[1].([]byte), args[2].([]byte))
	p.PopN(3)
}

func execmEncodingEncodeToString(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*base32.Encoding).EncodeToString(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmEncodingEncodedLen(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*base32.Encoding).EncodedLen(args[1].(int))
	p.Ret(2, ret0)
}

func execmEncodingDecode(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*base32.Encoding).Decode(args[1].([]byte), args[2].([]byte))
	p.Ret(3, ret0, ret1)
}

func execmEncodingDecodeString(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*base32.Encoding).DecodeString(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmEncodingDecodedLen(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*base32.Encoding).DecodedLen(args[1].(int))
	p.Ret(2, ret0)
}

func toType0(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execNewDecoder(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := base32.NewDecoder(args[0].(*base32.Encoding), toType0(args[1]))
	p.Ret(2, ret0)
}

func toType1(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execNewEncoder(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := base32.NewEncoder(args[0].(*base32.Encoding), toType1(args[1]))
	p.Ret(2, ret0)
}

func execNewEncoding(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := base32.NewEncoding(args[0].(string))
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("encoding/base32")

func init() {
	I.RegisterFuncs(
		I.Func("(CorruptInputError).Error", (base32.CorruptInputError).Error, execmCorruptInputErrorError),
		I.Func("(Encoding).WithPadding", (base32.Encoding).WithPadding, execmEncodingWithPadding),
		I.Func("(*Encoding).Encode", (*base32.Encoding).Encode, execmEncodingEncode),
		I.Func("(*Encoding).EncodeToString", (*base32.Encoding).EncodeToString, execmEncodingEncodeToString),
		I.Func("(*Encoding).EncodedLen", (*base32.Encoding).EncodedLen, execmEncodingEncodedLen),
		I.Func("(*Encoding).Decode", (*base32.Encoding).Decode, execmEncodingDecode),
		I.Func("(*Encoding).DecodeString", (*base32.Encoding).DecodeString, execmEncodingDecodeString),
		I.Func("(*Encoding).DecodedLen", (*base32.Encoding).DecodedLen, execmEncodingDecodedLen),
		I.Func("NewDecoder", base32.NewDecoder, execNewDecoder),
		I.Func("NewEncoder", base32.NewEncoder, execNewEncoder),
		I.Func("NewEncoding", base32.NewEncoding, execNewEncoding),
	)
	I.RegisterVars(
		I.Var("HexEncoding", &base32.HexEncoding),
		I.Var("StdEncoding", &base32.StdEncoding),
	)
	I.RegisterTypes(
		I.Type("CorruptInputError", reflect.TypeOf((*base32.CorruptInputError)(nil)).Elem()),
		I.Type("Encoding", reflect.TypeOf((*base32.Encoding)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("NoPadding", qspec.Uint32, base32.NoPadding),
		I.Const("StdPadding", qspec.Uint32, base32.StdPadding),
	)
}
