// Package base64 provide Go+ "encoding/base64" package, as "encoding/base64" package in Go.
package base64

import (
	base64 "encoding/base64"
	io "io"
	reflect "reflect"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execmCorruptInputErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(base64.CorruptInputError).Error()
	p.Ret(1, ret0)
}

func execmEncodingWithPadding(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(base64.Encoding).WithPadding(args[1].(rune))
	p.Ret(2, ret0)
}

func execmEncodingStrict(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(base64.Encoding).Strict()
	p.Ret(1, ret0)
}

func execmEncodingEncode(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*base64.Encoding).Encode(args[1].([]byte), args[2].([]byte))
	p.PopN(3)
}

func execmEncodingEncodeToString(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*base64.Encoding).EncodeToString(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmEncodingEncodedLen(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*base64.Encoding).EncodedLen(args[1].(int))
	p.Ret(2, ret0)
}

func execmEncodingDecodeString(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*base64.Encoding).DecodeString(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmEncodingDecode(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*base64.Encoding).Decode(args[1].([]byte), args[2].([]byte))
	p.Ret(3, ret0, ret1)
}

func execmEncodingDecodedLen(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*base64.Encoding).DecodedLen(args[1].(int))
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
	ret0 := base64.NewDecoder(args[0].(*base64.Encoding), toType0(args[1]))
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
	ret0 := base64.NewEncoder(args[0].(*base64.Encoding), toType1(args[1]))
	p.Ret(2, ret0)
}

func execNewEncoding(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := base64.NewEncoding(args[0].(string))
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("encoding/base64")

func init() {
	I.RegisterFuncs(
		I.Func("(CorruptInputError).Error", (base64.CorruptInputError).Error, execmCorruptInputErrorError),
		I.Func("(Encoding).WithPadding", (base64.Encoding).WithPadding, execmEncodingWithPadding),
		I.Func("(Encoding).Strict", (base64.Encoding).Strict, execmEncodingStrict),
		I.Func("(*Encoding).Encode", (*base64.Encoding).Encode, execmEncodingEncode),
		I.Func("(*Encoding).EncodeToString", (*base64.Encoding).EncodeToString, execmEncodingEncodeToString),
		I.Func("(*Encoding).EncodedLen", (*base64.Encoding).EncodedLen, execmEncodingEncodedLen),
		I.Func("(*Encoding).DecodeString", (*base64.Encoding).DecodeString, execmEncodingDecodeString),
		I.Func("(*Encoding).Decode", (*base64.Encoding).Decode, execmEncodingDecode),
		I.Func("(*Encoding).DecodedLen", (*base64.Encoding).DecodedLen, execmEncodingDecodedLen),
		I.Func("NewDecoder", base64.NewDecoder, execNewDecoder),
		I.Func("NewEncoder", base64.NewEncoder, execNewEncoder),
		I.Func("NewEncoding", base64.NewEncoding, execNewEncoding),
	)
	I.RegisterVars(
		I.Var("RawStdEncoding", &base64.RawStdEncoding),
		I.Var("RawURLEncoding", &base64.RawURLEncoding),
		I.Var("StdEncoding", &base64.StdEncoding),
		I.Var("URLEncoding", &base64.URLEncoding),
	)
	I.RegisterTypes(
		I.Type("CorruptInputError", reflect.TypeOf((*base64.CorruptInputError)(nil)).Elem()),
		I.Type("Encoding", reflect.TypeOf((*base64.Encoding)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("NoPadding", qspec.Uint32, base64.NoPadding),
		I.Const("StdPadding", qspec.Uint32, base64.StdPadding),
	)
}
