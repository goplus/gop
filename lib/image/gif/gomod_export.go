// Package gif provide Go+ "image/gif" package, as "image/gif" package in Go.
package gif

import (
	image "image"
	gif "image/gif"
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
	ret0, ret1 := gif.Decode(toType0(args[0]))
	p.Ret(1, ret0, ret1)
}

func execDecodeAll(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := gif.DecodeAll(toType0(args[0]))
	p.Ret(1, ret0, ret1)
}

func execDecodeConfig(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := gif.DecodeConfig(toType0(args[0]))
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
	ret0 := gif.Encode(toType1(args[0]), toType2(args[1]), args[2].(*gif.Options))
	p.Ret(3, ret0)
}

func execEncodeAll(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := gif.EncodeAll(toType1(args[0]), args[1].(*gif.GIF))
	p.Ret(2, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("image/gif")

func init() {
	I.RegisterFuncs(
		I.Func("Decode", gif.Decode, execDecode),
		I.Func("DecodeAll", gif.DecodeAll, execDecodeAll),
		I.Func("DecodeConfig", gif.DecodeConfig, execDecodeConfig),
		I.Func("Encode", gif.Encode, execEncode),
		I.Func("EncodeAll", gif.EncodeAll, execEncodeAll),
	)
	I.RegisterTypes(
		I.Type("GIF", reflect.TypeOf((*gif.GIF)(nil)).Elem()),
		I.Type("Options", reflect.TypeOf((*gif.Options)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("DisposalBackground", qspec.ConstUnboundInt, gif.DisposalBackground),
		I.Const("DisposalNone", qspec.ConstUnboundInt, gif.DisposalNone),
		I.Const("DisposalPrevious", qspec.ConstUnboundInt, gif.DisposalPrevious),
	)
}
