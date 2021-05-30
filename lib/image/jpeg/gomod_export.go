// Package jpeg provide Go+ "image/jpeg" package, as "image/jpeg" package in Go.
package jpeg

import (
	image "image"
	jpeg "image/jpeg"
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
	ret0, ret1 := jpeg.Decode(toType0(args[0]))
	p.Ret(1, ret0, ret1)
}

func execDecodeConfig(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := jpeg.DecodeConfig(toType0(args[0]))
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
	args := p.GetArgs(3)
	ret0 := jpeg.Encode(toType1(args[0]), toType2(args[1]), args[2].(*jpeg.Options))
	p.Ret(3, ret0)
}

func execmFormatErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(jpeg.FormatError).Error()
	p.Ret(1, ret0)
}

func execiReaderRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(jpeg.Reader).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execiReaderReadByte(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(jpeg.Reader).ReadByte()
	p.Ret(1, ret0, ret1)
}

func execmUnsupportedErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(jpeg.UnsupportedError).Error()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("image/jpeg")

func init() {
	I.RegisterFuncs(
		I.Func("Decode", jpeg.Decode, execDecode),
		I.Func("DecodeConfig", jpeg.DecodeConfig, execDecodeConfig),
		I.Func("Encode", jpeg.Encode, execEncode),
		I.Func("(FormatError).Error", (jpeg.FormatError).Error, execmFormatErrorError),
		I.Func("(Reader).Read", (jpeg.Reader).Read, execiReaderRead),
		I.Func("(Reader).ReadByte", (jpeg.Reader).ReadByte, execiReaderReadByte),
		I.Func("(UnsupportedError).Error", (jpeg.UnsupportedError).Error, execmUnsupportedErrorError),
	)
	I.RegisterTypes(
		I.Type("FormatError", reflect.TypeOf((*jpeg.FormatError)(nil)).Elem()),
		I.Type("Options", reflect.TypeOf((*jpeg.Options)(nil)).Elem()),
		I.Type("Reader", reflect.TypeOf((*jpeg.Reader)(nil)).Elem()),
		I.Type("UnsupportedError", reflect.TypeOf((*jpeg.UnsupportedError)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("DefaultQuality", qspec.ConstUnboundInt, jpeg.DefaultQuality),
	)
}
