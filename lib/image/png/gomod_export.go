// Package png provide Go+ "image/png" package, as "image/png" package in Go.
package png

import (
	image "image"
	png "image/png"
	io "io"
	reflect "reflect"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func toType0(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execDecode(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := png.Decode(toType0(args[0]))
	p.Ret(1, ret0, ret1)
}

func execDecodeConfig(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := png.DecodeConfig(toType0(args[0]))
	p.Ret(1, ret0, ret1)
}

func toType1(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func toType2(v interface{}) image.Image {
	if v == nil {
		return nil
	}
	return v.(image.Image)
}

func execEncode(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := png.Encode(toType1(args[0]), toType2(args[1]))
	p.Ret(2, ret0)
}

func execmEncoderEncode(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*png.Encoder).Encode(toType1(args[1]), toType2(args[2]))
	p.Ret(3, ret0)
}

func execiEncoderBufferPoolGet(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(png.EncoderBufferPool).Get()
	p.Ret(1, ret0)
}

func execiEncoderBufferPoolPut(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(png.EncoderBufferPool).Put(args[1].(*png.EncoderBuffer))
	p.PopN(2)
}

func execmFormatErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(png.FormatError).Error()
	p.Ret(1, ret0)
}

func execmUnsupportedErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(png.UnsupportedError).Error()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("image/png")

func init() {
	I.RegisterFuncs(
		I.Func("Decode", png.Decode, execDecode),
		I.Func("DecodeConfig", png.DecodeConfig, execDecodeConfig),
		I.Func("Encode", png.Encode, execEncode),
		I.Func("(*Encoder).Encode", (*png.Encoder).Encode, execmEncoderEncode),
		I.Func("(EncoderBufferPool).Get", (png.EncoderBufferPool).Get, execiEncoderBufferPoolGet),
		I.Func("(EncoderBufferPool).Put", (png.EncoderBufferPool).Put, execiEncoderBufferPoolPut),
		I.Func("(FormatError).Error", (png.FormatError).Error, execmFormatErrorError),
		I.Func("(UnsupportedError).Error", (png.UnsupportedError).Error, execmUnsupportedErrorError),
	)
	I.RegisterTypes(
		I.Type("CompressionLevel", reflect.TypeOf((*png.CompressionLevel)(nil)).Elem()),
		I.Type("Encoder", reflect.TypeOf((*png.Encoder)(nil)).Elem()),
		I.Type("EncoderBuffer", reflect.TypeOf((*png.EncoderBuffer)(nil)).Elem()),
		I.Type("EncoderBufferPool", reflect.TypeOf((*png.EncoderBufferPool)(nil)).Elem()),
		I.Type("FormatError", reflect.TypeOf((*png.FormatError)(nil)).Elem()),
		I.Type("UnsupportedError", reflect.TypeOf((*png.UnsupportedError)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("BestCompression", qspec.Int, png.BestCompression),
		I.Const("BestSpeed", qspec.Int, png.BestSpeed),
		I.Const("DefaultCompression", qspec.Int, png.DefaultCompression),
		I.Const("NoCompression", qspec.Int, png.NoCompression),
	)
}
