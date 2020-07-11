/*
 Copyright 2020 The GoPlus Authors (goplus.org)

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package constant

import (
	"math/big"
	"testing"

	"github.com/goplus/gop/token"
	"github.com/qiniu/x/ts"
)

// ----------------------------------------------------------------------------

func TestOp(test *testing.T) {
	t := ts.New(test)
	x := MakeRatFromString("1.3e-2")
	y := Make(int64(3))
	t.Case("-Rat 13/1000", UnaryOp(token.SUB, x, 0)).PropEqual("kind", Rat).PropEqual("string", "-13/1000")
	t.Case("-Int 3", UnaryOp(token.SUB, y, 0)).PropEqual("kind", Int).PropEqual("string", "-3")
	t.Case("Rat 13/1000 + Int 3", BinaryOp(x, token.ADD, y)).PropEqual("kind", Rat).PropEqual("string", "3013/1000")
	t.Case("Rat 1/2 + Float 1.5", BinaryOp(MakeRatFromString("1p-1"), token.ADD, MakeFloat64(1.5))).PropEqual("kind", Rat).PropEqual("string", "2")
	t.Case("Rat 1.2 + Rat 1/2", BinaryOp(MakeRatFromString("1.2"), token.ADD, MakeRatFromString("1p-1"))).PropEqual("kind", Rat).PropEqual("string", "17/10")
	t.Case("Int 1 + Int 2", BinaryOp(MakeInt64(1), token.ADD, MakeInt64(2))).PropEqual("kind", Int).PropEqual("string", "3")
	t.Case("Cmp(Int 1 != Int 2)", Compare(MakeInt64(1), token.NEQ, MakeInt64(2))).Equal(true)
	t.Case("Cmp(Rat 1p1 != Int 2)", Compare(MakeRatFromString("1p1"), token.EQL, MakeInt64(2))).Equal(true)
	t.Case("Cmp(Rat 1p1 != Rat 2)", Compare(MakeRatFromString("1p1"), token.EQL, MakeRatFromString("2"))).Equal(true)
	t.Case("Shift(Int 1, Int 2)", Shift(MakeInt64(1), token.SHL, 2)).Equal(MakeInt64(4))
	t.Case("Shift(Rat 1, Int 2)", Shift(MakeRatFromString("1"), token.SHL, 2)).PropEqual("kind", Rat).PropEqual("string", "4")
}

func TestMakeVal(test *testing.T) {
	t := ts.New(test)
	t.Call(MakeUnknown).PropEqual("kind", Unknown).PropEqual("string", "unknown")
	t.Call(MakeBool, true).PropEqual("kind", Bool).PropEqual("string", "true")
	t.Call(MakeString, "hello").PropEqual("kind", String).PropEqual("string", `"hello"`)
	t.Call(MakeInt64, int64(-10000)).PropEqual("kind", Int).PropEqual("string", "-10000")
	t.Call(MakeUint64, uint64(10000)).PropEqual("kind", Int).PropEqual("string", "10000")
	t.Call(MakeFloat64, 1e4).PropEqual("kind", Float).PropEqual("string", "10000")
	t.Call(Make, false).PropEqual("kind", Bool).PropEqual("string", "false")
	t.Call(MakeImag, MakeInt64(1000)).PropEqual("kind", Complex).PropEqual("string", "(0 + 1000i)")
	t.Call(MakeImag, MakeRat(int64(1000))).PropEqual("kind", Complex).PropEqual("string", "(0 + 1000i)")
	t.Call(MakeFromLiteral, "1x1r", token.RAT, uint(0)).PropEqual("kind", Unknown)
	t.Call(MakeFromLiteral, "1r", token.RAT, uint(0)).PropEqual("kind", Rat).PropEqual("string", "1")
	t.Call(Val, MakeFromLiteral("1.2r", token.RAT, uint(0))).Equal(big.NewRat(6, 5))
	t.Call(Val, MakeFromLiteral("1.2", token.FLOAT, uint(0))).Equal(big.NewRat(6, 5))
	t.Call(Val, MakeFromLiteral("3.14159265358979323846264338327950288419716939937510582097494459", token.FLOAT, uint(0))).
		PropEqual("string", "314159265358979323846264338327950288419716939937510582097494459/100000000000000000000000000000000000000000000000000000000000000")
}

func TestMakeRat(test *testing.T) {
	t := ts.New(test)
	t.Call(MakeRat, int64(1)).PropEqual("kind", Rat).PropEqual("string", "1")
	t.Call(MakeRat, big.NewInt(1)).PropEqual("kind", Rat).PropEqual("string", "1")
	t.Call(MakeRat, big.NewRat(1, 3)).PropEqual("kind", Rat).PropEqual("string", "1/3")
}

func TestMakeRatFromString(test *testing.T) {
	t := ts.New(test)
	t.Call(MakeRatFromString, "1x").PropEqual("kind", Unknown)
	t.Call(MakeRatFromString, "1e3x").PropEqual("kind", Unknown)
	t.Call(MakeRatFromString, "1").PropEqual("kind", Rat).PropEqual("string", "1")
	t.Call(MakeRatFromString, "1.234567").PropEqual("kind", Rat).PropEqual("string", "1234567/1000000")
	t.Call(MakeRatFromString, ".234567").PropEqual("kind", Rat).PropEqual("string", "234567/1000000")
	t.Call(MakeRatFromString, "1_111e2").PropEqual("kind", Rat).PropEqual("string", "111100")
	t.Call(MakeRatFromString, "1_111E-2").PropEqual("kind", Rat).PropEqual("string", "1111/100")
	t.Call(MakeRatFromString, "1111p-2").PropEqual("kind", Rat).PropEqual("string", "1111/4")
	t.Call(MakeRatFromString, "1111P2").PropEqual("kind", Rat).PropEqual("string", "4444")
	t.Call(MakeRatFromString, "1.3e-2").PropEqual("kind", Rat).PropEqual("string", "13/1000")
}

func TestVal(test *testing.T) {
	t := ts.New(test)
	t.Call(BoolVal, MakeBool(true)).Equal(true)
	t.Call(StringVal, MakeString("hello")).Equal("hello")
	t.Call(Int64Val, MakeInt64(1000)).Equal(int64(1000)).Next().Equal(true)
	t.Call(Int64Val, MakeRat(int64(1000))).Equal(int64(1000)).Next().Equal(true)
	t.Call(Uint64Val, MakeInt64(1000)).Equal(uint64(1000)).Next().Equal(true)
	t.Call(Uint64Val, MakeRat(int64(1000))).Equal(uint64(1000)).Next().Equal(true)
	t.Call(Uint64Val, MakeInt64(-1000)).Equal(uint64(18446744073709550616)).Next().Equal(false)
	t.Call(Float32Val, MakeInt64(-1000)).Equal(float32(-1000)).Next().Equal(true)
	t.Call(Float32Val, MakeRat(int64(-1000))).Equal(float32(-1000)).Next().Equal(true)
	t.Call(Float64Val, MakeInt64(-1000)).Equal(float64(-1000)).Next().Equal(true)
	t.Call(Float64Val, MakeRat(int64(-1000))).Equal(float64(-1000)).Next().Equal(true)
	t.Call(Val, MakeInt64(-1000)).Equal(int64(-1000))
	t.Call(BitLen, MakeInt64(4)).Equal(3)
	t.Call(Sign, MakeInt64(4)).Equal(1)
	t.Call(Sign, MakeRat(int64(-1000))).Equal(-1)
	t.Call(Bytes, MakeInt64(0x408)).Equal([]byte{8, 4})
	t.Call(MakeFromBytes, []byte{8, 4}).Equal(MakeInt64(0x408))
	t.Call(Num, MakeRat(int64(4))).Equal(MakeInt64(4))
	t.Call(Num, MakeInt64(4)).Equal(MakeInt64(4))
	t.Call(Denom, MakeRat(int64(4))).Equal(MakeInt64(1))
	t.Call(Denom, MakeInt64(4)).Equal(MakeInt64(1))
	t.Call(Real, MakeRat(int64(4))).Equal(MakeRat(int64(4)))
	t.Call(Real, MakeImag(MakeInt64(4))).Equal(MakeInt64(0))
	t.Call(Imag, MakeRat(int64(4))).Equal(MakeInt64(0))
	t.Call(Imag, MakeInt64(4)).Equal(MakeInt64(0))
}

func TestToVal(test *testing.T) {
	t := ts.New(test)
	t.Call(IsInt, MakeFloat64(10)).Equal(true)
	t.Call(IsInt, MakeString("10")).Equal(false)
	t.Call(IsInt, MakeInt64(10)).Equal(true)
	t.Call(ToInt, MakeRat(int64(10))).Equal(MakeInt64(10))
	t.Call(ToInt, MakeFloat64(10)).Equal(MakeInt64(10))
	t.Call(ToInt, MakeFloat64(10.8)).Equal(MakeUnknown())
	t.Call(ToInt, MakeFloat64(10)).Equal(MakeInt64(10))
	t.Call(ToFloat, MakeRat(int64(10))).PropEqual("kind", Float).PropEqual("string", "10")
	t.Call(ToFloat, MakeInt64(10)).PropEqual("kind", Float).PropEqual("string", "10")
	t.Call(ToComplex, MakeRat(int64(10))).PropEqual("kind", Complex).PropEqual("string", "(10 + 0i)")
	t.Call(ToComplex, MakeInt64(10)).PropEqual("kind", Complex).PropEqual("string", "(10 + 0i)")
}

// ----------------------------------------------------------------------------
