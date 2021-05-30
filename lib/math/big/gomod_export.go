// Package big provide Go+ "math/big" package, as "math/big" package in Go.
package big

import (
	fmt "fmt"
	big "math/big"
	rand "math/rand"
	reflect "reflect"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execmAccuracyString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(big.Accuracy).String()
	p.Ret(1, ret0)
}

func execmErrNaNError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(big.ErrNaN).Error()
	p.Ret(1, ret0)
}

func execmFloatSetPrec(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Float).SetPrec(args[1].(uint))
	p.Ret(2, ret0)
}

func execmFloatSetMode(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Float).SetMode(args[1].(big.RoundingMode))
	p.Ret(2, ret0)
}

func execmFloatPrec(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*big.Float).Prec()
	p.Ret(1, ret0)
}

func execmFloatMinPrec(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*big.Float).MinPrec()
	p.Ret(1, ret0)
}

func execmFloatMode(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*big.Float).Mode()
	p.Ret(1, ret0)
}

func execmFloatAcc(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*big.Float).Acc()
	p.Ret(1, ret0)
}

func execmFloatSign(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*big.Float).Sign()
	p.Ret(1, ret0)
}

func execmFloatMantExp(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Float).MantExp(args[1].(*big.Float))
	p.Ret(2, ret0)
}

func execmFloatSetMantExp(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Float).SetMantExp(args[1].(*big.Float), args[2].(int))
	p.Ret(3, ret0)
}

func execmFloatSignbit(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*big.Float).Signbit()
	p.Ret(1, ret0)
}

func execmFloatIsInf(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*big.Float).IsInf()
	p.Ret(1, ret0)
}

func execmFloatIsInt(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*big.Float).IsInt()
	p.Ret(1, ret0)
}

func execmFloatSetUint64(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Float).SetUint64(args[1].(uint64))
	p.Ret(2, ret0)
}

func execmFloatSetInt64(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Float).SetInt64(args[1].(int64))
	p.Ret(2, ret0)
}

func execmFloatSetFloat64(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Float).SetFloat64(args[1].(float64))
	p.Ret(2, ret0)
}

func execmFloatSetInt(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Float).SetInt(args[1].(*big.Int))
	p.Ret(2, ret0)
}

func execmFloatSetRat(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Float).SetRat(args[1].(*big.Rat))
	p.Ret(2, ret0)
}

func execmFloatSetInf(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Float).SetInf(args[1].(bool))
	p.Ret(2, ret0)
}

func execmFloatSet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Float).Set(args[1].(*big.Float))
	p.Ret(2, ret0)
}

func execmFloatCopy(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Float).Copy(args[1].(*big.Float))
	p.Ret(2, ret0)
}

func execmFloatUint64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*big.Float).Uint64()
	p.Ret(1, ret0, ret1)
}

func execmFloatInt64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*big.Float).Int64()
	p.Ret(1, ret0, ret1)
}

func execmFloatFloat32(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*big.Float).Float32()
	p.Ret(1, ret0, ret1)
}

func execmFloatFloat64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*big.Float).Float64()
	p.Ret(1, ret0, ret1)
}

func execmFloatInt(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*big.Float).Int(args[1].(*big.Int))
	p.Ret(2, ret0, ret1)
}

func execmFloatRat(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*big.Float).Rat(args[1].(*big.Rat))
	p.Ret(2, ret0, ret1)
}

func execmFloatAbs(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Float).Abs(args[1].(*big.Float))
	p.Ret(2, ret0)
}

func execmFloatNeg(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Float).Neg(args[1].(*big.Float))
	p.Ret(2, ret0)
}

func execmFloatAdd(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Float).Add(args[1].(*big.Float), args[2].(*big.Float))
	p.Ret(3, ret0)
}

func execmFloatSub(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Float).Sub(args[1].(*big.Float), args[2].(*big.Float))
	p.Ret(3, ret0)
}

func execmFloatMul(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Float).Mul(args[1].(*big.Float), args[2].(*big.Float))
	p.Ret(3, ret0)
}

func execmFloatQuo(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Float).Quo(args[1].(*big.Float), args[2].(*big.Float))
	p.Ret(3, ret0)
}

func execmFloatCmp(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Float).Cmp(args[1].(*big.Float))
	p.Ret(2, ret0)
}

func execmFloatSetString(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*big.Float).SetString(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmFloatParse(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1, ret2 := args[0].(*big.Float).Parse(args[1].(string), args[2].(int))
	p.Ret(3, ret0, ret1, ret2)
}

func toType0(v interface{}) fmt.ScanState {
	if v == nil {
		return nil
	}
	return v.(fmt.ScanState)
}

func execmFloatScan(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Float).Scan(toType0(args[1]), args[2].(rune))
	p.Ret(3, ret0)
}

func execmFloatGobEncode(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*big.Float).GobEncode()
	p.Ret(1, ret0, ret1)
}

func execmFloatGobDecode(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Float).GobDecode(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmFloatMarshalText(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*big.Float).MarshalText()
	p.Ret(1, ret0, ret1)
}

func execmFloatUnmarshalText(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Float).UnmarshalText(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmFloatText(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Float).Text(args[1].(byte), args[2].(int))
	p.Ret(3, ret0)
}

func execmFloatString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*big.Float).String()
	p.Ret(1, ret0)
}

func execmFloatAppend(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := args[0].(*big.Float).Append(args[1].([]byte), args[2].(byte), args[3].(int))
	p.Ret(4, ret0)
}

func toType1(v interface{}) fmt.State {
	if v == nil {
		return nil
	}
	return v.(fmt.State)
}

func execmFloatFormat(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*big.Float).Format(toType1(args[1]), args[2].(rune))
	p.PopN(3)
}

func execmFloatSqrt(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Float).Sqrt(args[1].(*big.Float))
	p.Ret(2, ret0)
}

func execmIntSign(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*big.Int).Sign()
	p.Ret(1, ret0)
}

func execmIntSetInt64(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Int).SetInt64(args[1].(int64))
	p.Ret(2, ret0)
}

func execmIntSetUint64(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Int).SetUint64(args[1].(uint64))
	p.Ret(2, ret0)
}

func execmIntSet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Int).Set(args[1].(*big.Int))
	p.Ret(2, ret0)
}

func execmIntBits(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*big.Int).Bits()
	p.Ret(1, ret0)
}

func execmIntSetBits(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Int).SetBits(args[1].([]big.Word))
	p.Ret(2, ret0)
}

func execmIntAbs(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Int).Abs(args[1].(*big.Int))
	p.Ret(2, ret0)
}

func execmIntNeg(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Int).Neg(args[1].(*big.Int))
	p.Ret(2, ret0)
}

func execmIntAdd(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Int).Add(args[1].(*big.Int), args[2].(*big.Int))
	p.Ret(3, ret0)
}

func execmIntSub(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Int).Sub(args[1].(*big.Int), args[2].(*big.Int))
	p.Ret(3, ret0)
}

func execmIntMul(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Int).Mul(args[1].(*big.Int), args[2].(*big.Int))
	p.Ret(3, ret0)
}

func execmIntMulRange(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Int).MulRange(args[1].(int64), args[2].(int64))
	p.Ret(3, ret0)
}

func execmIntBinomial(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Int).Binomial(args[1].(int64), args[2].(int64))
	p.Ret(3, ret0)
}

func execmIntQuo(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Int).Quo(args[1].(*big.Int), args[2].(*big.Int))
	p.Ret(3, ret0)
}

func execmIntRem(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Int).Rem(args[1].(*big.Int), args[2].(*big.Int))
	p.Ret(3, ret0)
}

func execmIntQuoRem(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := args[0].(*big.Int).QuoRem(args[1].(*big.Int), args[2].(*big.Int), args[3].(*big.Int))
	p.Ret(4, ret0, ret1)
}

func execmIntDiv(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Int).Div(args[1].(*big.Int), args[2].(*big.Int))
	p.Ret(3, ret0)
}

func execmIntMod(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Int).Mod(args[1].(*big.Int), args[2].(*big.Int))
	p.Ret(3, ret0)
}

func execmIntDivMod(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := args[0].(*big.Int).DivMod(args[1].(*big.Int), args[2].(*big.Int), args[3].(*big.Int))
	p.Ret(4, ret0, ret1)
}

func execmIntCmp(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Int).Cmp(args[1].(*big.Int))
	p.Ret(2, ret0)
}

func execmIntCmpAbs(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Int).CmpAbs(args[1].(*big.Int))
	p.Ret(2, ret0)
}

func execmIntInt64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*big.Int).Int64()
	p.Ret(1, ret0)
}

func execmIntUint64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*big.Int).Uint64()
	p.Ret(1, ret0)
}

func execmIntIsInt64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*big.Int).IsInt64()
	p.Ret(1, ret0)
}

func execmIntIsUint64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*big.Int).IsUint64()
	p.Ret(1, ret0)
}

func execmIntSetString(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*big.Int).SetString(args[1].(string), args[2].(int))
	p.Ret(3, ret0, ret1)
}

func execmIntSetBytes(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Int).SetBytes(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmIntBytes(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*big.Int).Bytes()
	p.Ret(1, ret0)
}

func execmIntBitLen(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*big.Int).BitLen()
	p.Ret(1, ret0)
}

func execmIntTrailingZeroBits(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*big.Int).TrailingZeroBits()
	p.Ret(1, ret0)
}

func execmIntExp(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := args[0].(*big.Int).Exp(args[1].(*big.Int), args[2].(*big.Int), args[3].(*big.Int))
	p.Ret(4, ret0)
}

func execmIntGCD(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	ret0 := args[0].(*big.Int).GCD(args[1].(*big.Int), args[2].(*big.Int), args[3].(*big.Int), args[4].(*big.Int))
	p.Ret(5, ret0)
}

func execmIntRand(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Int).Rand(args[1].(*rand.Rand), args[2].(*big.Int))
	p.Ret(3, ret0)
}

func execmIntModInverse(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Int).ModInverse(args[1].(*big.Int), args[2].(*big.Int))
	p.Ret(3, ret0)
}

func execmIntModSqrt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Int).ModSqrt(args[1].(*big.Int), args[2].(*big.Int))
	p.Ret(3, ret0)
}

func execmIntLsh(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Int).Lsh(args[1].(*big.Int), args[2].(uint))
	p.Ret(3, ret0)
}

func execmIntRsh(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Int).Rsh(args[1].(*big.Int), args[2].(uint))
	p.Ret(3, ret0)
}

func execmIntBit(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Int).Bit(args[1].(int))
	p.Ret(2, ret0)
}

func execmIntSetBit(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := args[0].(*big.Int).SetBit(args[1].(*big.Int), args[2].(int), args[3].(uint))
	p.Ret(4, ret0)
}

func execmIntAnd(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Int).And(args[1].(*big.Int), args[2].(*big.Int))
	p.Ret(3, ret0)
}

func execmIntAndNot(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Int).AndNot(args[1].(*big.Int), args[2].(*big.Int))
	p.Ret(3, ret0)
}

func execmIntOr(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Int).Or(args[1].(*big.Int), args[2].(*big.Int))
	p.Ret(3, ret0)
}

func execmIntXor(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Int).Xor(args[1].(*big.Int), args[2].(*big.Int))
	p.Ret(3, ret0)
}

func execmIntNot(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Int).Not(args[1].(*big.Int))
	p.Ret(2, ret0)
}

func execmIntSqrt(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Int).Sqrt(args[1].(*big.Int))
	p.Ret(2, ret0)
}

func execmIntText(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Int).Text(args[1].(int))
	p.Ret(2, ret0)
}

func execmIntAppend(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Int).Append(args[1].([]byte), args[2].(int))
	p.Ret(3, ret0)
}

func execmIntString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*big.Int).String()
	p.Ret(1, ret0)
}

func execmIntFormat(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*big.Int).Format(toType1(args[1]), args[2].(rune))
	p.PopN(3)
}

func execmIntScan(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Int).Scan(toType0(args[1]), args[2].(rune))
	p.Ret(3, ret0)
}

func execmIntGobEncode(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*big.Int).GobEncode()
	p.Ret(1, ret0, ret1)
}

func execmIntGobDecode(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Int).GobDecode(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmIntMarshalText(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*big.Int).MarshalText()
	p.Ret(1, ret0, ret1)
}

func execmIntUnmarshalText(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Int).UnmarshalText(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmIntMarshalJSON(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*big.Int).MarshalJSON()
	p.Ret(1, ret0, ret1)
}

func execmIntUnmarshalJSON(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Int).UnmarshalJSON(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmIntProbablyPrime(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Int).ProbablyPrime(args[1].(int))
	p.Ret(2, ret0)
}

func execJacobi(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := big.Jacobi(args[0].(*big.Int), args[1].(*big.Int))
	p.Ret(2, ret0)
}

func execNewFloat(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := big.NewFloat(args[0].(float64))
	p.Ret(1, ret0)
}

func execNewInt(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := big.NewInt(args[0].(int64))
	p.Ret(1, ret0)
}

func execNewRat(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := big.NewRat(args[0].(int64), args[1].(int64))
	p.Ret(2, ret0)
}

func execParseFloat(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1, ret2 := big.ParseFloat(args[0].(string), args[1].(int), args[2].(uint), args[3].(big.RoundingMode))
	p.Ret(4, ret0, ret1, ret2)
}

func execmRatSetFloat64(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Rat).SetFloat64(args[1].(float64))
	p.Ret(2, ret0)
}

func execmRatFloat32(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*big.Rat).Float32()
	p.Ret(1, ret0, ret1)
}

func execmRatFloat64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*big.Rat).Float64()
	p.Ret(1, ret0, ret1)
}

func execmRatSetFrac(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Rat).SetFrac(args[1].(*big.Int), args[2].(*big.Int))
	p.Ret(3, ret0)
}

func execmRatSetFrac64(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Rat).SetFrac64(args[1].(int64), args[2].(int64))
	p.Ret(3, ret0)
}

func execmRatSetInt(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Rat).SetInt(args[1].(*big.Int))
	p.Ret(2, ret0)
}

func execmRatSetInt64(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Rat).SetInt64(args[1].(int64))
	p.Ret(2, ret0)
}

func execmRatSetUint64(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Rat).SetUint64(args[1].(uint64))
	p.Ret(2, ret0)
}

func execmRatSet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Rat).Set(args[1].(*big.Rat))
	p.Ret(2, ret0)
}

func execmRatAbs(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Rat).Abs(args[1].(*big.Rat))
	p.Ret(2, ret0)
}

func execmRatNeg(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Rat).Neg(args[1].(*big.Rat))
	p.Ret(2, ret0)
}

func execmRatInv(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Rat).Inv(args[1].(*big.Rat))
	p.Ret(2, ret0)
}

func execmRatSign(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*big.Rat).Sign()
	p.Ret(1, ret0)
}

func execmRatIsInt(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*big.Rat).IsInt()
	p.Ret(1, ret0)
}

func execmRatNum(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*big.Rat).Num()
	p.Ret(1, ret0)
}

func execmRatDenom(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*big.Rat).Denom()
	p.Ret(1, ret0)
}

func execmRatCmp(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Rat).Cmp(args[1].(*big.Rat))
	p.Ret(2, ret0)
}

func execmRatAdd(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Rat).Add(args[1].(*big.Rat), args[2].(*big.Rat))
	p.Ret(3, ret0)
}

func execmRatSub(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Rat).Sub(args[1].(*big.Rat), args[2].(*big.Rat))
	p.Ret(3, ret0)
}

func execmRatMul(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Rat).Mul(args[1].(*big.Rat), args[2].(*big.Rat))
	p.Ret(3, ret0)
}

func execmRatQuo(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Rat).Quo(args[1].(*big.Rat), args[2].(*big.Rat))
	p.Ret(3, ret0)
}

func execmRatScan(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*big.Rat).Scan(toType0(args[1]), args[2].(rune))
	p.Ret(3, ret0)
}

func execmRatSetString(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*big.Rat).SetString(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmRatString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*big.Rat).String()
	p.Ret(1, ret0)
}

func execmRatRatString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*big.Rat).RatString()
	p.Ret(1, ret0)
}

func execmRatFloatString(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Rat).FloatString(args[1].(int))
	p.Ret(2, ret0)
}

func execmRatGobEncode(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*big.Rat).GobEncode()
	p.Ret(1, ret0, ret1)
}

func execmRatGobDecode(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Rat).GobDecode(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmRatMarshalText(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*big.Rat).MarshalText()
	p.Ret(1, ret0, ret1)
}

func execmRatUnmarshalText(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*big.Rat).UnmarshalText(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmRoundingModeString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(big.RoundingMode).String()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("math/big")

func init() {
	I.RegisterFuncs(
		I.Func("(Accuracy).String", (big.Accuracy).String, execmAccuracyString),
		I.Func("(ErrNaN).Error", (big.ErrNaN).Error, execmErrNaNError),
		I.Func("(*Float).SetPrec", (*big.Float).SetPrec, execmFloatSetPrec),
		I.Func("(*Float).SetMode", (*big.Float).SetMode, execmFloatSetMode),
		I.Func("(*Float).Prec", (*big.Float).Prec, execmFloatPrec),
		I.Func("(*Float).MinPrec", (*big.Float).MinPrec, execmFloatMinPrec),
		I.Func("(*Float).Mode", (*big.Float).Mode, execmFloatMode),
		I.Func("(*Float).Acc", (*big.Float).Acc, execmFloatAcc),
		I.Func("(*Float).Sign", (*big.Float).Sign, execmFloatSign),
		I.Func("(*Float).MantExp", (*big.Float).MantExp, execmFloatMantExp),
		I.Func("(*Float).SetMantExp", (*big.Float).SetMantExp, execmFloatSetMantExp),
		I.Func("(*Float).Signbit", (*big.Float).Signbit, execmFloatSignbit),
		I.Func("(*Float).IsInf", (*big.Float).IsInf, execmFloatIsInf),
		I.Func("(*Float).IsInt", (*big.Float).IsInt, execmFloatIsInt),
		I.Func("(*Float).SetUint64", (*big.Float).SetUint64, execmFloatSetUint64),
		I.Func("(*Float).SetInt64", (*big.Float).SetInt64, execmFloatSetInt64),
		I.Func("(*Float).SetFloat64", (*big.Float).SetFloat64, execmFloatSetFloat64),
		I.Func("(*Float).SetInt", (*big.Float).SetInt, execmFloatSetInt),
		I.Func("(*Float).SetRat", (*big.Float).SetRat, execmFloatSetRat),
		I.Func("(*Float).SetInf", (*big.Float).SetInf, execmFloatSetInf),
		I.Func("(*Float).Set", (*big.Float).Set, execmFloatSet),
		I.Func("(*Float).Copy", (*big.Float).Copy, execmFloatCopy),
		I.Func("(*Float).Uint64", (*big.Float).Uint64, execmFloatUint64),
		I.Func("(*Float).Int64", (*big.Float).Int64, execmFloatInt64),
		I.Func("(*Float).Float32", (*big.Float).Float32, execmFloatFloat32),
		I.Func("(*Float).Float64", (*big.Float).Float64, execmFloatFloat64),
		I.Func("(*Float).Int", (*big.Float).Int, execmFloatInt),
		I.Func("(*Float).Rat", (*big.Float).Rat, execmFloatRat),
		I.Func("(*Float).Abs", (*big.Float).Abs, execmFloatAbs),
		I.Func("(*Float).Neg", (*big.Float).Neg, execmFloatNeg),
		I.Func("(*Float).Add", (*big.Float).Add, execmFloatAdd),
		I.Func("(*Float).Sub", (*big.Float).Sub, execmFloatSub),
		I.Func("(*Float).Mul", (*big.Float).Mul, execmFloatMul),
		I.Func("(*Float).Quo", (*big.Float).Quo, execmFloatQuo),
		I.Func("(*Float).Cmp", (*big.Float).Cmp, execmFloatCmp),
		I.Func("(*Float).SetString", (*big.Float).SetString, execmFloatSetString),
		I.Func("(*Float).Parse", (*big.Float).Parse, execmFloatParse),
		I.Func("(*Float).Scan", (*big.Float).Scan, execmFloatScan),
		I.Func("(*Float).GobEncode", (*big.Float).GobEncode, execmFloatGobEncode),
		I.Func("(*Float).GobDecode", (*big.Float).GobDecode, execmFloatGobDecode),
		I.Func("(*Float).MarshalText", (*big.Float).MarshalText, execmFloatMarshalText),
		I.Func("(*Float).UnmarshalText", (*big.Float).UnmarshalText, execmFloatUnmarshalText),
		I.Func("(*Float).Text", (*big.Float).Text, execmFloatText),
		I.Func("(*Float).String", (*big.Float).String, execmFloatString),
		I.Func("(*Float).Append", (*big.Float).Append, execmFloatAppend),
		I.Func("(*Float).Format", (*big.Float).Format, execmFloatFormat),
		I.Func("(*Float).Sqrt", (*big.Float).Sqrt, execmFloatSqrt),
		I.Func("(*Int).Sign", (*big.Int).Sign, execmIntSign),
		I.Func("(*Int).SetInt64", (*big.Int).SetInt64, execmIntSetInt64),
		I.Func("(*Int).SetUint64", (*big.Int).SetUint64, execmIntSetUint64),
		I.Func("(*Int).Set", (*big.Int).Set, execmIntSet),
		I.Func("(*Int).Bits", (*big.Int).Bits, execmIntBits),
		I.Func("(*Int).SetBits", (*big.Int).SetBits, execmIntSetBits),
		I.Func("(*Int).Abs", (*big.Int).Abs, execmIntAbs),
		I.Func("(*Int).Neg", (*big.Int).Neg, execmIntNeg),
		I.Func("(*Int).Add", (*big.Int).Add, execmIntAdd),
		I.Func("(*Int).Sub", (*big.Int).Sub, execmIntSub),
		I.Func("(*Int).Mul", (*big.Int).Mul, execmIntMul),
		I.Func("(*Int).MulRange", (*big.Int).MulRange, execmIntMulRange),
		I.Func("(*Int).Binomial", (*big.Int).Binomial, execmIntBinomial),
		I.Func("(*Int).Quo", (*big.Int).Quo, execmIntQuo),
		I.Func("(*Int).Rem", (*big.Int).Rem, execmIntRem),
		I.Func("(*Int).QuoRem", (*big.Int).QuoRem, execmIntQuoRem),
		I.Func("(*Int).Div", (*big.Int).Div, execmIntDiv),
		I.Func("(*Int).Mod", (*big.Int).Mod, execmIntMod),
		I.Func("(*Int).DivMod", (*big.Int).DivMod, execmIntDivMod),
		I.Func("(*Int).Cmp", (*big.Int).Cmp, execmIntCmp),
		I.Func("(*Int).CmpAbs", (*big.Int).CmpAbs, execmIntCmpAbs),
		I.Func("(*Int).Int64", (*big.Int).Int64, execmIntInt64),
		I.Func("(*Int).Uint64", (*big.Int).Uint64, execmIntUint64),
		I.Func("(*Int).IsInt64", (*big.Int).IsInt64, execmIntIsInt64),
		I.Func("(*Int).IsUint64", (*big.Int).IsUint64, execmIntIsUint64),
		I.Func("(*Int).SetString", (*big.Int).SetString, execmIntSetString),
		I.Func("(*Int).SetBytes", (*big.Int).SetBytes, execmIntSetBytes),
		I.Func("(*Int).Bytes", (*big.Int).Bytes, execmIntBytes),
		I.Func("(*Int).BitLen", (*big.Int).BitLen, execmIntBitLen),
		I.Func("(*Int).TrailingZeroBits", (*big.Int).TrailingZeroBits, execmIntTrailingZeroBits),
		I.Func("(*Int).Exp", (*big.Int).Exp, execmIntExp),
		I.Func("(*Int).GCD", (*big.Int).GCD, execmIntGCD),
		I.Func("(*Int).Rand", (*big.Int).Rand, execmIntRand),
		I.Func("(*Int).ModInverse", (*big.Int).ModInverse, execmIntModInverse),
		I.Func("(*Int).ModSqrt", (*big.Int).ModSqrt, execmIntModSqrt),
		I.Func("(*Int).Lsh", (*big.Int).Lsh, execmIntLsh),
		I.Func("(*Int).Rsh", (*big.Int).Rsh, execmIntRsh),
		I.Func("(*Int).Bit", (*big.Int).Bit, execmIntBit),
		I.Func("(*Int).SetBit", (*big.Int).SetBit, execmIntSetBit),
		I.Func("(*Int).And", (*big.Int).And, execmIntAnd),
		I.Func("(*Int).AndNot", (*big.Int).AndNot, execmIntAndNot),
		I.Func("(*Int).Or", (*big.Int).Or, execmIntOr),
		I.Func("(*Int).Xor", (*big.Int).Xor, execmIntXor),
		I.Func("(*Int).Not", (*big.Int).Not, execmIntNot),
		I.Func("(*Int).Sqrt", (*big.Int).Sqrt, execmIntSqrt),
		I.Func("(*Int).Text", (*big.Int).Text, execmIntText),
		I.Func("(*Int).Append", (*big.Int).Append, execmIntAppend),
		I.Func("(*Int).String", (*big.Int).String, execmIntString),
		I.Func("(*Int).Format", (*big.Int).Format, execmIntFormat),
		I.Func("(*Int).Scan", (*big.Int).Scan, execmIntScan),
		I.Func("(*Int).GobEncode", (*big.Int).GobEncode, execmIntGobEncode),
		I.Func("(*Int).GobDecode", (*big.Int).GobDecode, execmIntGobDecode),
		I.Func("(*Int).MarshalText", (*big.Int).MarshalText, execmIntMarshalText),
		I.Func("(*Int).UnmarshalText", (*big.Int).UnmarshalText, execmIntUnmarshalText),
		I.Func("(*Int).MarshalJSON", (*big.Int).MarshalJSON, execmIntMarshalJSON),
		I.Func("(*Int).UnmarshalJSON", (*big.Int).UnmarshalJSON, execmIntUnmarshalJSON),
		I.Func("(*Int).ProbablyPrime", (*big.Int).ProbablyPrime, execmIntProbablyPrime),
		I.Func("Jacobi", big.Jacobi, execJacobi),
		I.Func("NewFloat", big.NewFloat, execNewFloat),
		I.Func("NewInt", big.NewInt, execNewInt),
		I.Func("NewRat", big.NewRat, execNewRat),
		I.Func("ParseFloat", big.ParseFloat, execParseFloat),
		I.Func("(*Rat).SetFloat64", (*big.Rat).SetFloat64, execmRatSetFloat64),
		I.Func("(*Rat).Float32", (*big.Rat).Float32, execmRatFloat32),
		I.Func("(*Rat).Float64", (*big.Rat).Float64, execmRatFloat64),
		I.Func("(*Rat).SetFrac", (*big.Rat).SetFrac, execmRatSetFrac),
		I.Func("(*Rat).SetFrac64", (*big.Rat).SetFrac64, execmRatSetFrac64),
		I.Func("(*Rat).SetInt", (*big.Rat).SetInt, execmRatSetInt),
		I.Func("(*Rat).SetInt64", (*big.Rat).SetInt64, execmRatSetInt64),
		I.Func("(*Rat).SetUint64", (*big.Rat).SetUint64, execmRatSetUint64),
		I.Func("(*Rat).Set", (*big.Rat).Set, execmRatSet),
		I.Func("(*Rat).Abs", (*big.Rat).Abs, execmRatAbs),
		I.Func("(*Rat).Neg", (*big.Rat).Neg, execmRatNeg),
		I.Func("(*Rat).Inv", (*big.Rat).Inv, execmRatInv),
		I.Func("(*Rat).Sign", (*big.Rat).Sign, execmRatSign),
		I.Func("(*Rat).IsInt", (*big.Rat).IsInt, execmRatIsInt),
		I.Func("(*Rat).Num", (*big.Rat).Num, execmRatNum),
		I.Func("(*Rat).Denom", (*big.Rat).Denom, execmRatDenom),
		I.Func("(*Rat).Cmp", (*big.Rat).Cmp, execmRatCmp),
		I.Func("(*Rat).Add", (*big.Rat).Add, execmRatAdd),
		I.Func("(*Rat).Sub", (*big.Rat).Sub, execmRatSub),
		I.Func("(*Rat).Mul", (*big.Rat).Mul, execmRatMul),
		I.Func("(*Rat).Quo", (*big.Rat).Quo, execmRatQuo),
		I.Func("(*Rat).Scan", (*big.Rat).Scan, execmRatScan),
		I.Func("(*Rat).SetString", (*big.Rat).SetString, execmRatSetString),
		I.Func("(*Rat).String", (*big.Rat).String, execmRatString),
		I.Func("(*Rat).RatString", (*big.Rat).RatString, execmRatRatString),
		I.Func("(*Rat).FloatString", (*big.Rat).FloatString, execmRatFloatString),
		I.Func("(*Rat).GobEncode", (*big.Rat).GobEncode, execmRatGobEncode),
		I.Func("(*Rat).GobDecode", (*big.Rat).GobDecode, execmRatGobDecode),
		I.Func("(*Rat).MarshalText", (*big.Rat).MarshalText, execmRatMarshalText),
		I.Func("(*Rat).UnmarshalText", (*big.Rat).UnmarshalText, execmRatUnmarshalText),
		I.Func("(RoundingMode).String", (big.RoundingMode).String, execmRoundingModeString),
	)
	I.RegisterTypes(
		I.Type("Accuracy", reflect.TypeOf((*big.Accuracy)(nil)).Elem()),
		I.Type("ErrNaN", reflect.TypeOf((*big.ErrNaN)(nil)).Elem()),
		I.Type("Float", reflect.TypeOf((*big.Float)(nil)).Elem()),
		I.Type("Int", reflect.TypeOf((*big.Int)(nil)).Elem()),
		I.Type("Rat", reflect.TypeOf((*big.Rat)(nil)).Elem()),
		I.Type("RoundingMode", reflect.TypeOf((*big.RoundingMode)(nil)).Elem()),
		I.Type("Word", reflect.TypeOf((*big.Word)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("Above", qspec.Int8, big.Above),
		I.Const("AwayFromZero", qspec.Uint8, big.AwayFromZero),
		I.Const("Below", qspec.Int8, big.Below),
		I.Const("Exact", qspec.Int8, big.Exact),
		I.Const("MaxBase", qspec.ConstBoundRune, big.MaxBase),
		I.Const("MaxExp", qspec.ConstUnboundInt, big.MaxExp),
		I.Const("MaxPrec", qspec.Uint64, uint64(big.MaxPrec)),
		I.Const("MinExp", qspec.ConstUnboundInt, big.MinExp),
		I.Const("ToNearestAway", qspec.Uint8, big.ToNearestAway),
		I.Const("ToNearestEven", qspec.Uint8, big.ToNearestEven),
		I.Const("ToNegativeInf", qspec.Uint8, big.ToNegativeInf),
		I.Const("ToPositiveInf", qspec.Uint8, big.ToPositiveInf),
		I.Const("ToZero", qspec.Uint8, big.ToZero),
	)
}
