// Package mime provide Go+ "mime" package, as "mime" package in Go.
package mime

import (
	mime "mime"
	reflect "reflect"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execAddExtensionType(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := mime.AddExtensionType(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execExtensionsByType(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := mime.ExtensionsByType(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execFormatMediaType(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := mime.FormatMediaType(args[0].(string), args[1].(map[string]string))
	p.Ret(2, ret0)
}

func execParseMediaType(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := mime.ParseMediaType(args[0].(string))
	p.Ret(1, ret0, ret1, ret2)
}

func execTypeByExtension(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := mime.TypeByExtension(args[0].(string))
	p.Ret(1, ret0)
}

func execmWordDecoderDecode(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*mime.WordDecoder).Decode(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmWordDecoderDecodeHeader(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*mime.WordDecoder).DecodeHeader(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmWordEncoderEncode(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(mime.WordEncoder).Encode(args[1].(string), args[2].(string))
	p.Ret(3, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("mime")

func init() {
	I.RegisterFuncs(
		I.Func("AddExtensionType", mime.AddExtensionType, execAddExtensionType),
		I.Func("ExtensionsByType", mime.ExtensionsByType, execExtensionsByType),
		I.Func("FormatMediaType", mime.FormatMediaType, execFormatMediaType),
		I.Func("ParseMediaType", mime.ParseMediaType, execParseMediaType),
		I.Func("TypeByExtension", mime.TypeByExtension, execTypeByExtension),
		I.Func("(*WordDecoder).Decode", (*mime.WordDecoder).Decode, execmWordDecoderDecode),
		I.Func("(*WordDecoder).DecodeHeader", (*mime.WordDecoder).DecodeHeader, execmWordDecoderDecodeHeader),
		I.Func("(WordEncoder).Encode", (mime.WordEncoder).Encode, execmWordEncoderEncode),
	)
	I.RegisterVars(
		I.Var("ErrInvalidMediaParameter", &mime.ErrInvalidMediaParameter),
	)
	I.RegisterTypes(
		I.Type("WordDecoder", reflect.TypeOf((*mime.WordDecoder)(nil)).Elem()),
		I.Type("WordEncoder", reflect.TypeOf((*mime.WordEncoder)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("BEncoding", qspec.Uint8, mime.BEncoding),
		I.Const("QEncoding", qspec.Uint8, mime.QEncoding),
	)
}
