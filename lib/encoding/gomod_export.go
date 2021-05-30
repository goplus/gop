// Package encoding provide Go+ "encoding" package, as "encoding" package in Go.
package encoding

import (
	encoding "encoding"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func execiBinaryMarshalerMarshalBinary(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(encoding.BinaryMarshaler).MarshalBinary()
	p.Ret(1, ret0, ret1)
}

func execiBinaryUnmarshalerUnmarshalBinary(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(encoding.BinaryUnmarshaler).UnmarshalBinary(args[1].([]byte))
	p.Ret(2, ret0)
}

func execiTextMarshalerMarshalText(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(encoding.TextMarshaler).MarshalText()
	p.Ret(1, ret0, ret1)
}

func execiTextUnmarshalerUnmarshalText(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(encoding.TextUnmarshaler).UnmarshalText(args[1].([]byte))
	p.Ret(2, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("encoding")

func init() {
	I.RegisterFuncs(
		I.Func("(BinaryMarshaler).MarshalBinary", (encoding.BinaryMarshaler).MarshalBinary, execiBinaryMarshalerMarshalBinary),
		I.Func("(BinaryUnmarshaler).UnmarshalBinary", (encoding.BinaryUnmarshaler).UnmarshalBinary, execiBinaryUnmarshalerUnmarshalBinary),
		I.Func("(TextMarshaler).MarshalText", (encoding.TextMarshaler).MarshalText, execiTextMarshalerMarshalText),
		I.Func("(TextUnmarshaler).UnmarshalText", (encoding.TextUnmarshaler).UnmarshalText, execiTextUnmarshalerUnmarshalText),
	)
	I.RegisterTypes(
		I.Type("BinaryMarshaler", reflect.TypeOf((*encoding.BinaryMarshaler)(nil)).Elem()),
		I.Type("BinaryUnmarshaler", reflect.TypeOf((*encoding.BinaryUnmarshaler)(nil)).Elem()),
		I.Type("TextMarshaler", reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()),
		I.Type("TextUnmarshaler", reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()),
	)
}
