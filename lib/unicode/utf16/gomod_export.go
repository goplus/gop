// Package utf16 provide Go+ "unicode/utf16" package, as "unicode/utf16" package in Go.
package utf16

import (
	utf16 "unicode/utf16"

	gop "github.com/goplus/gop"
)

func execDecode(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := utf16.Decode(args[0].([]uint16))
	p.Ret(1, ret0)
}

func execDecodeRune(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := utf16.DecodeRune(args[0].(rune), args[1].(rune))
	p.Ret(2, ret0)
}

func execEncode(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := utf16.Encode(args[0].([]rune))
	p.Ret(1, ret0)
}

func execEncodeRune(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := utf16.EncodeRune(args[0].(rune))
	p.Ret(1, ret0, ret1)
}

func execIsSurrogate(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := utf16.IsSurrogate(args[0].(rune))
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("unicode/utf16")

func init() {
	I.RegisterFuncs(
		I.Func("Decode", utf16.Decode, execDecode),
		I.Func("DecodeRune", utf16.DecodeRune, execDecodeRune),
		I.Func("Encode", utf16.Encode, execEncode),
		I.Func("EncodeRune", utf16.EncodeRune, execEncodeRune),
		I.Func("IsSurrogate", utf16.IsSurrogate, execIsSurrogate),
	)
}
