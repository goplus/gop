// Package bits provide Go+ "math/bits" package, as "math/bits" package in Go.
package bits

import (
	bits "math/bits"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execAdd(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := bits.Add(args[0].(uint), args[1].(uint), args[2].(uint))
	p.Ret(3, ret0, ret1)
}

func execAdd32(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := bits.Add32(args[0].(uint32), args[1].(uint32), args[2].(uint32))
	p.Ret(3, ret0, ret1)
}

func execAdd64(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := bits.Add64(args[0].(uint64), args[1].(uint64), args[2].(uint64))
	p.Ret(3, ret0, ret1)
}

func execDiv(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := bits.Div(args[0].(uint), args[1].(uint), args[2].(uint))
	p.Ret(3, ret0, ret1)
}

func execDiv32(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := bits.Div32(args[0].(uint32), args[1].(uint32), args[2].(uint32))
	p.Ret(3, ret0, ret1)
}

func execDiv64(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := bits.Div64(args[0].(uint64), args[1].(uint64), args[2].(uint64))
	p.Ret(3, ret0, ret1)
}

func execLeadingZeros(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.LeadingZeros(args[0].(uint))
	p.Ret(1, ret0)
}

func execLeadingZeros16(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.LeadingZeros16(args[0].(uint16))
	p.Ret(1, ret0)
}

func execLeadingZeros32(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.LeadingZeros32(args[0].(uint32))
	p.Ret(1, ret0)
}

func execLeadingZeros64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.LeadingZeros64(args[0].(uint64))
	p.Ret(1, ret0)
}

func execLeadingZeros8(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.LeadingZeros8(args[0].(uint8))
	p.Ret(1, ret0)
}

func execLen(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.Len(args[0].(uint))
	p.Ret(1, ret0)
}

func execLen16(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.Len16(args[0].(uint16))
	p.Ret(1, ret0)
}

func execLen32(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.Len32(args[0].(uint32))
	p.Ret(1, ret0)
}

func execLen64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.Len64(args[0].(uint64))
	p.Ret(1, ret0)
}

func execLen8(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.Len8(args[0].(uint8))
	p.Ret(1, ret0)
}

func execMul(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := bits.Mul(args[0].(uint), args[1].(uint))
	p.Ret(2, ret0, ret1)
}

func execMul32(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := bits.Mul32(args[0].(uint32), args[1].(uint32))
	p.Ret(2, ret0, ret1)
}

func execMul64(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := bits.Mul64(args[0].(uint64), args[1].(uint64))
	p.Ret(2, ret0, ret1)
}

func execOnesCount(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.OnesCount(args[0].(uint))
	p.Ret(1, ret0)
}

func execOnesCount16(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.OnesCount16(args[0].(uint16))
	p.Ret(1, ret0)
}

func execOnesCount32(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.OnesCount32(args[0].(uint32))
	p.Ret(1, ret0)
}

func execOnesCount64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.OnesCount64(args[0].(uint64))
	p.Ret(1, ret0)
}

func execOnesCount8(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.OnesCount8(args[0].(uint8))
	p.Ret(1, ret0)
}

func execRem(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := bits.Rem(args[0].(uint), args[1].(uint), args[2].(uint))
	p.Ret(3, ret0)
}

func execRem32(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := bits.Rem32(args[0].(uint32), args[1].(uint32), args[2].(uint32))
	p.Ret(3, ret0)
}

func execRem64(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := bits.Rem64(args[0].(uint64), args[1].(uint64), args[2].(uint64))
	p.Ret(3, ret0)
}

func execReverse(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.Reverse(args[0].(uint))
	p.Ret(1, ret0)
}

func execReverse16(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.Reverse16(args[0].(uint16))
	p.Ret(1, ret0)
}

func execReverse32(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.Reverse32(args[0].(uint32))
	p.Ret(1, ret0)
}

func execReverse64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.Reverse64(args[0].(uint64))
	p.Ret(1, ret0)
}

func execReverse8(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.Reverse8(args[0].(uint8))
	p.Ret(1, ret0)
}

func execReverseBytes(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.ReverseBytes(args[0].(uint))
	p.Ret(1, ret0)
}

func execReverseBytes16(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.ReverseBytes16(args[0].(uint16))
	p.Ret(1, ret0)
}

func execReverseBytes32(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.ReverseBytes32(args[0].(uint32))
	p.Ret(1, ret0)
}

func execReverseBytes64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.ReverseBytes64(args[0].(uint64))
	p.Ret(1, ret0)
}

func execRotateLeft(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bits.RotateLeft(args[0].(uint), args[1].(int))
	p.Ret(2, ret0)
}

func execRotateLeft16(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bits.RotateLeft16(args[0].(uint16), args[1].(int))
	p.Ret(2, ret0)
}

func execRotateLeft32(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bits.RotateLeft32(args[0].(uint32), args[1].(int))
	p.Ret(2, ret0)
}

func execRotateLeft64(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bits.RotateLeft64(args[0].(uint64), args[1].(int))
	p.Ret(2, ret0)
}

func execRotateLeft8(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bits.RotateLeft8(args[0].(uint8), args[1].(int))
	p.Ret(2, ret0)
}

func execSub(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := bits.Sub(args[0].(uint), args[1].(uint), args[2].(uint))
	p.Ret(3, ret0, ret1)
}

func execSub32(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := bits.Sub32(args[0].(uint32), args[1].(uint32), args[2].(uint32))
	p.Ret(3, ret0, ret1)
}

func execSub64(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := bits.Sub64(args[0].(uint64), args[1].(uint64), args[2].(uint64))
	p.Ret(3, ret0, ret1)
}

func execTrailingZeros(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.TrailingZeros(args[0].(uint))
	p.Ret(1, ret0)
}

func execTrailingZeros16(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.TrailingZeros16(args[0].(uint16))
	p.Ret(1, ret0)
}

func execTrailingZeros32(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.TrailingZeros32(args[0].(uint32))
	p.Ret(1, ret0)
}

func execTrailingZeros64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.TrailingZeros64(args[0].(uint64))
	p.Ret(1, ret0)
}

func execTrailingZeros8(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bits.TrailingZeros8(args[0].(uint8))
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("math/bits")

func init() {
	I.RegisterFuncs(
		I.Func("Add", bits.Add, execAdd),
		I.Func("Add32", bits.Add32, execAdd32),
		I.Func("Add64", bits.Add64, execAdd64),
		I.Func("Div", bits.Div, execDiv),
		I.Func("Div32", bits.Div32, execDiv32),
		I.Func("Div64", bits.Div64, execDiv64),
		I.Func("LeadingZeros", bits.LeadingZeros, execLeadingZeros),
		I.Func("LeadingZeros16", bits.LeadingZeros16, execLeadingZeros16),
		I.Func("LeadingZeros32", bits.LeadingZeros32, execLeadingZeros32),
		I.Func("LeadingZeros64", bits.LeadingZeros64, execLeadingZeros64),
		I.Func("LeadingZeros8", bits.LeadingZeros8, execLeadingZeros8),
		I.Func("Len", bits.Len, execLen),
		I.Func("Len16", bits.Len16, execLen16),
		I.Func("Len32", bits.Len32, execLen32),
		I.Func("Len64", bits.Len64, execLen64),
		I.Func("Len8", bits.Len8, execLen8),
		I.Func("Mul", bits.Mul, execMul),
		I.Func("Mul32", bits.Mul32, execMul32),
		I.Func("Mul64", bits.Mul64, execMul64),
		I.Func("OnesCount", bits.OnesCount, execOnesCount),
		I.Func("OnesCount16", bits.OnesCount16, execOnesCount16),
		I.Func("OnesCount32", bits.OnesCount32, execOnesCount32),
		I.Func("OnesCount64", bits.OnesCount64, execOnesCount64),
		I.Func("OnesCount8", bits.OnesCount8, execOnesCount8),
		I.Func("Rem", bits.Rem, execRem),
		I.Func("Rem32", bits.Rem32, execRem32),
		I.Func("Rem64", bits.Rem64, execRem64),
		I.Func("Reverse", bits.Reverse, execReverse),
		I.Func("Reverse16", bits.Reverse16, execReverse16),
		I.Func("Reverse32", bits.Reverse32, execReverse32),
		I.Func("Reverse64", bits.Reverse64, execReverse64),
		I.Func("Reverse8", bits.Reverse8, execReverse8),
		I.Func("ReverseBytes", bits.ReverseBytes, execReverseBytes),
		I.Func("ReverseBytes16", bits.ReverseBytes16, execReverseBytes16),
		I.Func("ReverseBytes32", bits.ReverseBytes32, execReverseBytes32),
		I.Func("ReverseBytes64", bits.ReverseBytes64, execReverseBytes64),
		I.Func("RotateLeft", bits.RotateLeft, execRotateLeft),
		I.Func("RotateLeft16", bits.RotateLeft16, execRotateLeft16),
		I.Func("RotateLeft32", bits.RotateLeft32, execRotateLeft32),
		I.Func("RotateLeft64", bits.RotateLeft64, execRotateLeft64),
		I.Func("RotateLeft8", bits.RotateLeft8, execRotateLeft8),
		I.Func("Sub", bits.Sub, execSub),
		I.Func("Sub32", bits.Sub32, execSub32),
		I.Func("Sub64", bits.Sub64, execSub64),
		I.Func("TrailingZeros", bits.TrailingZeros, execTrailingZeros),
		I.Func("TrailingZeros16", bits.TrailingZeros16, execTrailingZeros16),
		I.Func("TrailingZeros32", bits.TrailingZeros32, execTrailingZeros32),
		I.Func("TrailingZeros64", bits.TrailingZeros64, execTrailingZeros64),
		I.Func("TrailingZeros8", bits.TrailingZeros8, execTrailingZeros8),
	)
	I.RegisterConsts(
		I.Const("UintSize", qspec.ConstUnboundInt, bits.UintSize),
	)
}
