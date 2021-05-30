// Package subtle provide Go+ "crypto/subtle" package, as "crypto/subtle" package in Go.
package subtle

import (
	subtle "crypto/subtle"

	gop "github.com/goplus/gop"
)

func execConstantTimeByteEq(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := subtle.ConstantTimeByteEq(args[0].(uint8), args[1].(uint8))
	p.Ret(2, ret0)
}

func execConstantTimeCompare(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := subtle.ConstantTimeCompare(args[0].([]byte), args[1].([]byte))
	p.Ret(2, ret0)
}

func execConstantTimeCopy(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	subtle.ConstantTimeCopy(args[0].(int), args[1].([]byte), args[2].([]byte))
	p.PopN(3)
}

func execConstantTimeEq(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := subtle.ConstantTimeEq(args[0].(int32), args[1].(int32))
	p.Ret(2, ret0)
}

func execConstantTimeLessOrEq(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := subtle.ConstantTimeLessOrEq(args[0].(int), args[1].(int))
	p.Ret(2, ret0)
}

func execConstantTimeSelect(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := subtle.ConstantTimeSelect(args[0].(int), args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("crypto/subtle")

func init() {
	I.RegisterFuncs(
		I.Func("ConstantTimeByteEq", subtle.ConstantTimeByteEq, execConstantTimeByteEq),
		I.Func("ConstantTimeCompare", subtle.ConstantTimeCompare, execConstantTimeCompare),
		I.Func("ConstantTimeCopy", subtle.ConstantTimeCopy, execConstantTimeCopy),
		I.Func("ConstantTimeEq", subtle.ConstantTimeEq, execConstantTimeEq),
		I.Func("ConstantTimeLessOrEq", subtle.ConstantTimeLessOrEq, execConstantTimeLessOrEq),
		I.Func("ConstantTimeSelect", subtle.ConstantTimeSelect, execConstantTimeSelect),
	)
}
