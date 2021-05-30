// Package utf8 provide Go+ "unicode/utf8" package, as "unicode/utf8" package in Go.
package utf8

import (
	utf8 "unicode/utf8"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execDecodeLastRune(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := utf8.DecodeLastRune(args[0].([]byte))
	p.Ret(1, ret0, ret1)
}

func execDecodeLastRuneInString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := utf8.DecodeLastRuneInString(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execDecodeRune(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := utf8.DecodeRune(args[0].([]byte))
	p.Ret(1, ret0, ret1)
}

func execDecodeRuneInString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := utf8.DecodeRuneInString(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execEncodeRune(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := utf8.EncodeRune(args[0].([]byte), args[1].(rune))
	p.Ret(2, ret0)
}

func execFullRune(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := utf8.FullRune(args[0].([]byte))
	p.Ret(1, ret0)
}

func execFullRuneInString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := utf8.FullRuneInString(args[0].(string))
	p.Ret(1, ret0)
}

func execRuneCount(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := utf8.RuneCount(args[0].([]byte))
	p.Ret(1, ret0)
}

func execRuneCountInString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := utf8.RuneCountInString(args[0].(string))
	p.Ret(1, ret0)
}

func execRuneLen(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := utf8.RuneLen(args[0].(rune))
	p.Ret(1, ret0)
}

func execRuneStart(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := utf8.RuneStart(args[0].(byte))
	p.Ret(1, ret0)
}

func execValid(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := utf8.Valid(args[0].([]byte))
	p.Ret(1, ret0)
}

func execValidRune(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := utf8.ValidRune(args[0].(rune))
	p.Ret(1, ret0)
}

func execValidString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := utf8.ValidString(args[0].(string))
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("unicode/utf8")

func init() {
	I.RegisterFuncs(
		I.Func("DecodeLastRune", utf8.DecodeLastRune, execDecodeLastRune),
		I.Func("DecodeLastRuneInString", utf8.DecodeLastRuneInString, execDecodeLastRuneInString),
		I.Func("DecodeRune", utf8.DecodeRune, execDecodeRune),
		I.Func("DecodeRuneInString", utf8.DecodeRuneInString, execDecodeRuneInString),
		I.Func("EncodeRune", utf8.EncodeRune, execEncodeRune),
		I.Func("FullRune", utf8.FullRune, execFullRune),
		I.Func("FullRuneInString", utf8.FullRuneInString, execFullRuneInString),
		I.Func("RuneCount", utf8.RuneCount, execRuneCount),
		I.Func("RuneCountInString", utf8.RuneCountInString, execRuneCountInString),
		I.Func("RuneLen", utf8.RuneLen, execRuneLen),
		I.Func("RuneStart", utf8.RuneStart, execRuneStart),
		I.Func("Valid", utf8.Valid, execValid),
		I.Func("ValidRune", utf8.ValidRune, execValidRune),
		I.Func("ValidString", utf8.ValidString, execValidString),
	)
	I.RegisterConsts(
		I.Const("MaxRune", qspec.ConstBoundRune, utf8.MaxRune),
		I.Const("RuneError", qspec.ConstBoundRune, utf8.RuneError),
		I.Const("RuneSelf", qspec.ConstUnboundInt, utf8.RuneSelf),
		I.Const("UTFMax", qspec.ConstUnboundInt, utf8.UTFMax),
	)
}
