/*
 * Copyright (c) 2021 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ng

import (
	"math/big"
)

func Gop_istmp(a any) bool {
	return false
}

// -----------------------------------------------------------------------------

type UntypedBigint *big.Int
type UntypedBigrat *big.Rat
type UntypedBigfloat *big.Float

type UntypedBigint_Default = Bigint
type UntypedBigrat_Default = Bigrat
type UntypedBigfloat_Default = Bigfloat

func UntypedBigint_Init__0(x int) UntypedBigint {
	panic(makeCompilerHappy)
}

func UntypedBigrat_Init__0(x int) UntypedBigrat {
	panic(makeCompilerHappy)
}

func UntypedBigrat_Init__1(x UntypedBigint) UntypedBigrat {
	panic(makeCompilerHappy)
}

const (
	makeCompilerHappy = "make compiler happy"
)

// -----------------------------------------------------------------------------
// type bigint

// A Bigint represents a signed multi-precision integer.
// The zero value for a Bigint represents nil.
type Bigint struct {
	*big.Int
}

func tmpint(a, b Bigint) Bigint {
	if Gop_istmp(a) {
		return a
	} else if Gop_istmp(b) {
		return b
	}
	return Bigint{new(big.Int)}
}

func tmpint1(a Bigint) Bigint {
	if Gop_istmp(a) {
		return a
	}
	return Bigint{new(big.Int)}
}

// IsNil returns a bigint object is nil or not
func (a Bigint) IsNil() bool {
	return a.Int == nil
}

// Gop_Add: func (a bigint) + (b bigint) bigint
func (a Bigint) Gop_Add(b Bigint) Bigint {
	return Bigint{tmpint(a, b).Add(a.Int, b.Int)}
}

// Gop_Sub: func (a bigint) - (b bigint) bigint
func (a Bigint) Gop_Sub(b Bigint) Bigint {
	return Bigint{tmpint(a, b).Sub(a.Int, b.Int)}
}

// Gop_Mul: func (a bigint) * (b bigint) bigint
func (a Bigint) Gop_Mul(b Bigint) Bigint {
	return Bigint{tmpint(a, b).Mul(a.Int, b.Int)}
}

// Gop_Quo: func (a bigint) / (b bigint) bigint {
func (a Bigint) Gop_Quo(b Bigint) Bigint {
	return Bigint{tmpint(a, b).Quo(a.Int, b.Int)}
}

// Gop_Rem: func (a bigint) % (b bigint) bigint
func (a Bigint) Gop_Rem(b Bigint) Bigint {
	return Bigint{tmpint(a, b).Rem(a.Int, b.Int)}
}

// Gop_Or: func (a bigint) | (b bigint) bigint
func (a Bigint) Gop_Or(b Bigint) Bigint {
	return Bigint{tmpint(a, b).Or(a.Int, b.Int)}
}

// Gop_Xor: func (a bigint) ^ (b bigint) bigint
func (a Bigint) Gop_Xor(b Bigint) Bigint {
	return Bigint{tmpint(a, b).Xor(a.Int, b.Int)}
}

// Gop_And: func (a bigint) & (b bigint) bigint
func (a Bigint) Gop_And(b Bigint) Bigint {
	return Bigint{tmpint(a, b).And(a.Int, b.Int)}
}

// Gop_AndNot: func (a bigint) &^ (b bigint) bigint
func (a Bigint) Gop_AndNot(b Bigint) Bigint {
	return Bigint{tmpint(a, b).AndNot(a.Int, b.Int)}
}

// Gop_Lsh: func (a bigint) << (n untyped_uint) bigint
func (a Bigint) Gop_Lsh(n Gop_ninteger) Bigint {
	return Bigint{tmpint1(a).Lsh(a.Int, uint(n))}
}

// Gop_Rsh: func (a bigint) >> (n untyped_uint) bigint
func (a Bigint) Gop_Rsh(n Gop_ninteger) Bigint {
	return Bigint{tmpint1(a).Rsh(a.Int, uint(n))}
}

// Gop_LT: func (a bigint) < (b bigint) bool
func (a Bigint) Gop_LT(b Bigint) bool {
	return a.Cmp(b.Int) < 0
}

// Gop_LE: func (a bigint) <= (b bigint) bool
func (a Bigint) Gop_LE(b Bigint) bool {
	return a.Cmp(b.Int) <= 0
}

// Gop_GT: func (a bigint) > (b bigint) bool
func (a Bigint) Gop_GT(b Bigint) bool {
	return a.Cmp(b.Int) > 0
}

// Gop_GE: func (a bigint) >= (b bigint) bool
func (a Bigint) Gop_GE(b Bigint) bool {
	return a.Cmp(b.Int) >= 0
}

// Gop_EQ: func (a bigint) == (b bigint) bool
func (a Bigint) Gop_EQ(b Bigint) bool {
	return a.Cmp(b.Int) == 0
}

// Gop_NE: func (a bigint) != (b bigint) bool
func (a Bigint) Gop_NE(b Bigint) bool {
	return a.Cmp(b.Int) != 0
}

// Gop_Neg: func -(a bigint) bigint
func (a Bigint) Gop_Neg() Bigint {
	return Bigint{tmpint1(a).Neg(a.Int)}
}

// Gop_Dup: func +(a bigint) bigint
func (a Bigint) Gop_Dup() Bigint {
	return Bigint{new(big.Int).Set(a.Int)}
}

// Gop_Not: func ^(a bigint) bigint
func (a Bigint) Gop_Not() Bigint {
	return Bigint{tmpint1(a).Not(a.Int)}
}

// Gop_Inc: func ++(b bigint)
func (a Bigint) Gop_Inc() {
	a.Int.Add(a.Int, big1)
}

// Gop_Dec: func --(b bigint)
func (a Bigint) Gop_Dec() {
	a.Int.Sub(a.Int, big1)
}

// Gop_AddAssign: func (a bigint) += (b bigint)
func (a Bigint) Gop_AddAssign(b Bigint) {
	a.Int.Add(a.Int, b.Int)
}

// Gop_SubAssign: func (a bigint) -= (b bigint)
func (a Bigint) Gop_SubAssign(b Bigint) {
	a.Int.Sub(a.Int, b.Int)
}

// Gop_MulAssign: func (a bigint) *= (b bigint)
func (a Bigint) Gop_MulAssign(b Bigint) {
	a.Int.Mul(a.Int, b.Int)
}

// Gop_QuoAssign: func (a bigint) /= (b bigint) {
func (a Bigint) Gop_QuoAssign(b Bigint) {
	a.Int.Quo(a.Int, b.Int)
}

// Gop_RemAssign: func (a bigint) %= (b bigint)
func (a Bigint) Gop_RemAssign(b Bigint) {
	a.Int.Rem(a.Int, b.Int)
}

// Gop_OrAssign: func (a bigint) |= (b bigint)
func (a Bigint) Gop_OrAssign(b Bigint) {
	a.Int.Or(a.Int, b.Int)
}

// Gop_XorAssign: func (a bigint) ^= (b bigint)
func (a Bigint) Gop_XorAssign(b Bigint) {
	a.Int.Xor(a.Int, b.Int)
}

// Gop_AndAssign: func (a bigint) &= (b bigint)
func (a Bigint) Gop_AndAssign(b Bigint) {
	a.Int.And(a.Int, b.Int)
}

// Gop_AndNotAssign: func (a bigint) &^= (b bigint)
func (a Bigint) Gop_AndNotAssign(b Bigint) {
	a.Int.AndNot(a.Int, b.Int)
}

// Gop_LshAssign: func (a bigint) <<= (n untyped_uint)
func (a Bigint) Gop_LshAssign(n Gop_ninteger) {
	a.Int.Lsh(a.Int, uint(n))
}

// Gop_RshAssign: func (a bigint) >>= (n untyped_uint)
func (a Bigint) Gop_RshAssign(n Gop_ninteger) {
	a.Int.Rsh(a.Int, uint(n))
}

// -----------------------------------------------------------------------------

// Gop_Rcast: func int64(x bigint) int64
func (a Bigint) Gop_Rcast__0() int64 {
	return a.Int64()
}

// Gop_Rcast: func int64(x bigint) (int64, bool)
func (a Bigint) Gop_Rcast__1() (int64, bool) {
	return a.Int64(), a.IsInt64()
}

// Gop_Rcast: func uint64(x bigint) uint64
func (a Bigint) Gop_Rcast__2() uint64 {
	return a.Uint64()
}

// Gop_Rcast: func uint64(x bigint) (uint64, bool)
func (a Bigint) Gop_Rcast__3() (uint64, bool) {
	return a.Uint64(), a.IsUint64()
}

// Bigint_Cast: func bigint(x untyped_int) bigint
func Bigint_Cast__0(x int) Bigint {
	return Bigint{big.NewInt(int64(x))}
}

// Bigint_Cast: func bigint(x untyped_bigint) bigint
func Bigint_Cast__1(x UntypedBigint) Bigint {
	return Bigint{x}
}

// Bigint_Cast: func bigint(x int64) bigint
func Bigint_Cast__2(x int64) Bigint {
	return Bigint{big.NewInt(x)}
}

// Bigint_Cast: func bigint(x uint64) bigint
func Bigint_Cast__3(x uint64) Bigint {
	return Bigint{new(big.Int).SetUint64(x)}
}

// Bigint_Cast: func bigint(x uint) bigint
func Bigint_Cast__4(x uint) Bigint {
	return Bigint{new(big.Int).SetUint64(uint64(x))}
}

// Bigint_Cast: func bigint(x *big.Int) bigint
func Bigint_Cast__5(x *big.Int) Bigint {
	return Bigint{x}
}

// Bigint_Cast: func bigint(x *big.Rat) bigint
func Bigint_Cast__6(x *big.Rat) Bigint {
	if x.IsInt() {
		return Bigint{x.Num()}
	}
	ret, _ := new(big.Float).SetRat(x).Int(nil)
	return Bigint{ret}
}

// Bigint_Cast: func bigint() bigint
func Bigint_Cast__7() Bigint {
	return Bigint{new(big.Int)}
}

// Bigint_Init: func bigint.init(x int) bigint
func Bigint_Init__0(x int) Bigint {
	return Bigint{big.NewInt(int64(x))}
}

// Bigint_Init: func bigint.init(x untyped_bigint) bigint
func Bigint_Init__1(x UntypedBigint) Bigint {
	return Bigint{x}
}

// Bigint_Init: func bigint.init(x *big.Int) bigint
func Bigint_Init__2(x *big.Int) Bigint {
	return Bigint{x}
}

// -----------------------------------------------------------------------------
// type bigrat

// A Bigrat represents a quotient a/b of arbitrary precision.
// The zero value for a Bigrat represents nil.
type Bigrat struct {
	*big.Rat
}

func tmprat(a, b Bigrat) Bigrat {
	if Gop_istmp(a) {
		return a
	} else if Gop_istmp(b) {
		return b
	}
	return Bigrat{new(big.Rat)}
}

func tmprat1(a Bigrat) Bigrat {
	if Gop_istmp(a) {
		return a
	}
	return Bigrat{new(big.Rat)}
}

// IsNil returns a bigrat object is nil or not
func (a Bigrat) IsNil() bool {
	return a.Rat == nil
}

// Gop_Assign: func (a bigrat) = (b bigrat)
func (a Bigrat) Gop_Assign(b Bigrat) {
	if Gop_istmp(b) {
		*a.Rat = *b.Rat
	} else {
		a.Rat.Set(b.Rat)
	}
}

// Gop_Add: func (a bigrat) + (b bigrat) bigrat
func (a Bigrat) Gop_Add(b Bigrat) Bigrat {
	return Bigrat{tmprat(a, b).Add(a.Rat, b.Rat)}
}

// Gop_Sub: func (a bigrat) - (b bigrat) bigrat
func (a Bigrat) Gop_Sub(b Bigrat) Bigrat {
	return Bigrat{tmprat(a, b).Sub(a.Rat, b.Rat)}
}

// Gop_Mul: func (a bigrat) * (b bigrat) bigrat
func (a Bigrat) Gop_Mul(b Bigrat) Bigrat {
	return Bigrat{tmprat(a, b).Mul(a.Rat, b.Rat)}
}

// Gop_Quo: func (a bigrat) / (b bigrat) bigrat
func (a Bigrat) Gop_Quo(b Bigrat) Bigrat {
	return Bigrat{tmprat(a, b).Quo(a.Rat, b.Rat)}
}

// Gop_LT: func (a bigrat) < (b bigrat) bool
func (a Bigrat) Gop_LT(b Bigrat) bool {
	return a.Cmp(b.Rat) < 0
}

// Gop_LE: func (a bigrat) <= (b bigrat) bool
func (a Bigrat) Gop_LE(b Bigrat) bool {
	return a.Cmp(b.Rat) <= 0
}

// Gop_GT: func (a bigrat) > (b bigrat) bool
func (a Bigrat) Gop_GT(b Bigrat) bool {
	return a.Cmp(b.Rat) > 0
}

// Gop_GE: func (a bigrat) >= (b bigrat) bool
func (a Bigrat) Gop_GE(b Bigrat) bool {
	return a.Cmp(b.Rat) >= 0
}

// Gop_EQ: func (a bigrat) == (b bigrat) bool
func (a Bigrat) Gop_EQ(b Bigrat) bool {
	return a.Cmp(b.Rat) == 0
}

// Gop_NE: func (a bigrat) != (b bigrat) bool
func (a Bigrat) Gop_NE(b Bigrat) bool {
	return a.Cmp(b.Rat) != 0
}

// Gop_Neg: func -(a bigrat) bigrat
func (a Bigrat) Gop_Neg() Bigrat {
	return Bigrat{tmprat1(a).Neg(a.Rat)}
}

// Gop_Dup: func +(a bigrat) bigrat
func (a Bigrat) Gop_Dup() Bigrat {
	return Bigrat{new(big.Rat).Set(a.Rat)}
}

// Gop_Inv: func /(a bigrat) bigrat
func (a Bigrat) Gop_Inv() Bigrat {
	return Bigrat{tmprat1(a).Inv(a.Rat)}
}

// Gop_Add: func (a bigrat) += (b bigrat)
func (a Bigrat) Gop_AddAssign(b Bigrat) {
	a.Rat.Add(a.Rat, b.Rat)
}

// Gop_Sub: func (a bigrat) -= (b bigrat)
func (a Bigrat) Gop_SubAssign(b Bigrat) {
	a.Rat.Sub(a.Rat, b.Rat)
}

// Gop_Mul: func (a bigrat) *= (b bigrat)
func (a Bigrat) Gop_MulAssign(b Bigrat) {
	a.Rat.Mul(a.Rat, b.Rat)
}

// Gop_Quo: func (a bigrat) /= (b bigrat)
func (a Bigrat) Gop_QuoAssign(b Bigrat) {
	a.Rat.Quo(a.Rat, b.Rat)
}

// -----------------------------------------------------------------------------

// Bigrat_Cast: func bigrat(x untyped_int) bigrat
func Bigrat_Cast__0(x int) Bigrat {
	return Bigrat{big.NewRat(int64(x), 1)}
}

// Bigrat_Cast: func bigrat(a untyped_bigint) bigrat
func Bigrat_Cast__1(a UntypedBigint) Bigrat {
	return Bigrat{new(big.Rat).SetInt(a)}
}

// Bigrat_Cast: func bigrat(a *big.Int) bigrat
func Bigrat_Cast__2(a *big.Int) Bigrat {
	return Bigrat{new(big.Rat).SetInt(a)}
}

// Bigrat_Cast: func bigrat(a bigint) bigrat
func Bigrat_Cast__3(a Bigint) Bigrat {
	return Bigrat{new(big.Rat).SetInt(a.Int)}
}

// Bigrat_Cast: func bigrat(a *big.Rat) bigrat
func Bigrat_Cast__4(a *big.Rat) Bigrat {
	return Bigrat{a}
}

// Bigrat_Cast: func bigrat() bigrat
func Bigrat_Cast__5() Bigrat {
	return Bigrat{new(big.Rat)}
}

// Bigrat_Cast: func bigrat(a, b int64) bigrat
func Bigrat_Cast__6(a, b int64) Bigrat {
	return Bigrat{big.NewRat(a, b)}
}

// Bigrat_Init: func bigrat.init(x untyped_int) bigrat
func Bigrat_Init__0(x int) Bigrat {
	return Bigrat{big.NewRat(int64(x), 1)}
}

// Bigrat_Init: func bigrat.init(x untyped_bigint) bigrat
func Bigrat_Init__1(x UntypedBigint) Bigrat {
	return Bigrat{new(big.Rat).SetInt(x)}
}

// Bigrat_Init: func bigrat.init(x *big.Rat) bigrat
func Bigrat_Init__2(x *big.Rat) Bigrat {
	return Bigrat{x}
}

// -----------------------------------------------------------------------------
// type bigfloat

// A Bigfloat represents a multi-precision floating point number.
// The zero value for a Float represents nil.
type Bigfloat struct {
	*big.Float
}

// -----------------------------------------------------------------------------
