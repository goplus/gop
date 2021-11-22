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

package builtin

import (
	"math/big"
)

const (
	GopPackage = true // to indicate this is a Go+ package
)

//
// Gop_: Go+ object prefix
// Gop_xxx_Cast: type Gop_xxx typecast
// xxxx__N: the Nth overload function
//

type Gop_ninteger = uint

func Gop_istmp(a interface{}) bool {
	return false
}

// -----------------------------------------------------------------------------

type Gop_untyped_bigint *big.Int
type Gop_untyped_bigrat *big.Rat
type Gop_untyped_bigfloat *big.Float

type Gop_untyped_bigint_Default = Gop_bigint
type Gop_untyped_bigrat_Default = Gop_bigrat
type Gop_untyped_bigfloat_Default = Gop_bigfloat

func Gop_untyped_bigint_Init__0(x int) Gop_untyped_bigint {
	panic("make compiler happy")
}

func Gop_untyped_bigrat_Init__0(x int) Gop_untyped_bigrat {
	panic("make compiler happy")
}

func Gop_untyped_bigrat_Init__1(x Gop_untyped_bigint) Gop_untyped_bigrat {
	panic("make compiler happy")
}

// -----------------------------------------------------------------------------
// type bigint

// A Gop_bigint represents a signed multi-precision integer.
// The zero value for a Gop_bigint represents nil.
type Gop_bigint struct {
	*big.Int
}

func tmpint(a, b Gop_bigint) Gop_bigint {
	if Gop_istmp(a) {
		return a
	} else if Gop_istmp(b) {
		return b
	}
	return Gop_bigint{new(big.Int)}
}

func tmpint1(a Gop_bigint) Gop_bigint {
	if Gop_istmp(a) {
		return a
	}
	return Gop_bigint{new(big.Int)}
}

// IsNil returns a bigint object is nil or not
func (a Gop_bigint) IsNil() bool {
	return a.Int == nil
}

// Gop_Assign: func (a bigint) = (b bigint)
func (a Gop_bigint) Gop_Assign(b Gop_bigint) {
	if Gop_istmp(b) {
		*a.Int = *b.Int
	} else {
		a.Int.Set(b.Int)
	}
}

// Gop_Add: func (a bigint) + (b bigint) bigint
func (a Gop_bigint) Gop_Add(b Gop_bigint) Gop_bigint {
	return Gop_bigint{tmpint(a, b).Add(a.Int, b.Int)}
}

// Gop_Sub: func (a bigint) - (b bigint) bigint
func (a Gop_bigint) Gop_Sub(b Gop_bigint) Gop_bigint {
	return Gop_bigint{tmpint(a, b).Sub(a.Int, b.Int)}
}

// Gop_Mul: func (a bigint) * (b bigint) bigint
func (a Gop_bigint) Gop_Mul(b Gop_bigint) Gop_bigint {
	return Gop_bigint{tmpint(a, b).Mul(a.Int, b.Int)}
}

// Gop_Quo: func (a bigint) / (b bigint) bigint {
func (a Gop_bigint) Gop_Quo(b Gop_bigint) Gop_bigint {
	return Gop_bigint{tmpint(a, b).Quo(a.Int, b.Int)}
}

// Gop_Rem: func (a bigint) % (b bigint) bigint
func (a Gop_bigint) Gop_Rem(b Gop_bigint) Gop_bigint {
	return Gop_bigint{tmpint(a, b).Rem(a.Int, b.Int)}
}

// Gop_Or: func (a bigint) | (b bigint) bigint
func (a Gop_bigint) Gop_Or(b Gop_bigint) Gop_bigint {
	return Gop_bigint{tmpint(a, b).Or(a.Int, b.Int)}
}

// Gop_Xor: func (a bigint) ^ (b bigint) bigint
func (a Gop_bigint) Gop_Xor(b Gop_bigint) Gop_bigint {
	return Gop_bigint{tmpint(a, b).Xor(a.Int, b.Int)}
}

// Gop_And: func (a bigint) & (b bigint) bigint
func (a Gop_bigint) Gop_And(b Gop_bigint) Gop_bigint {
	return Gop_bigint{tmpint(a, b).And(a.Int, b.Int)}
}

// Gop_AndNot: func (a bigint) &^ (b bigint) bigint
func (a Gop_bigint) Gop_AndNot(b Gop_bigint) Gop_bigint {
	return Gop_bigint{tmpint(a, b).AndNot(a.Int, b.Int)}
}

// Gop_Lsh: func (a bigint) << (n untyped_uint) bigint
func (a Gop_bigint) Gop_Lsh(n Gop_ninteger) Gop_bigint {
	return Gop_bigint{tmpint1(a).Lsh(a.Int, uint(n))}
}

// Gop_Rsh: func (a bigint) >> (n untyped_uint) bigint
func (a Gop_bigint) Gop_Rsh(n Gop_ninteger) Gop_bigint {
	return Gop_bigint{tmpint1(a).Rsh(a.Int, uint(n))}
}

// Gop_LT: func (a bigint) < (b bigint) bool
func (a Gop_bigint) Gop_LT(b Gop_bigint) bool {
	return a.Cmp(b.Int) < 0
}

// Gop_LE: func (a bigint) <= (b bigint) bool
func (a Gop_bigint) Gop_LE(b Gop_bigint) bool {
	return a.Cmp(b.Int) <= 0
}

// Gop_GT: func (a bigint) > (b bigint) bool
func (a Gop_bigint) Gop_GT(b Gop_bigint) bool {
	return a.Cmp(b.Int) > 0
}

// Gop_GE: func (a bigint) >= (b bigint) bool
func (a Gop_bigint) Gop_GE(b Gop_bigint) bool {
	return a.Cmp(b.Int) >= 0
}

// Gop_EQ: func (a bigint) == (b bigint) bool
func (a Gop_bigint) Gop_EQ(b Gop_bigint) bool {
	return a.Cmp(b.Int) == 0
}

// Gop_NE: func (a bigint) != (b bigint) bool
func (a Gop_bigint) Gop_NE(b Gop_bigint) bool {
	return a.Cmp(b.Int) != 0
}

// Gop_Neg: func -(a bigint) bigint
func (a Gop_bigint) Gop_Neg() Gop_bigint {
	return Gop_bigint{tmpint1(a).Neg(a.Int)}
}

// Gop_Pos: func +(a bigint) bigint
func (a Gop_bigint) Gop_Pos() Gop_bigint {
	return a
}

// Gop_Not: func ^(a bigint) bigint
func (a Gop_bigint) Gop_Not() Gop_bigint {
	return Gop_bigint{tmpint1(a).Not(a.Int)}
}

// Gop_Add: func (a bigint) += (b bigint)
func (a Gop_bigint) Gop_AddAssign(b Gop_bigint) {
	a.Int.Add(a.Int, b.Int)
}

// Gop_Sub: func (a bigint) -= (b bigint)
func (a Gop_bigint) Gop_SubAssign(b Gop_bigint) {
	a.Int.Sub(a.Int, b.Int)
}

// Gop_Mul: func (a bigint) *= (b bigint)
func (a Gop_bigint) Gop_MulAssign(b Gop_bigint) {
	a.Int.Mul(a.Int, b.Int)
}

// Gop_Quo: func (a bigint) /= (b bigint) {
func (a Gop_bigint) Gop_QuoAssign(b Gop_bigint) {
	a.Int.Quo(a.Int, b.Int)
}

// Gop_Rem: func (a bigint) %= (b bigint)
func (a Gop_bigint) Gop_RemAssign(b Gop_bigint) {
	a.Int.Rem(a.Int, b.Int)
}

// Gop_Or: func (a bigint) |= (b bigint)
func (a Gop_bigint) Gop_OrAssign(b Gop_bigint) {
	a.Int.Or(a.Int, b.Int)
}

// Gop_Xor: func (a bigint) ^= (b bigint)
func (a Gop_bigint) Gop_XorAssign(b Gop_bigint) {
	a.Int.Xor(a.Int, b.Int)
}

// Gop_And: func (a bigint) &= (b bigint)
func (a Gop_bigint) Gop_AndAssign(b Gop_bigint) {
	a.Int.And(a.Int, b.Int)
}

// Gop_AndNot: func (a bigint) &^= (b bigint)
func (a Gop_bigint) Gop_AndNotAssign(b Gop_bigint) {
	a.Int.AndNot(a.Int, b.Int)
}

// Gop_Lsh: func (a bigint) <<= (n untyped_uint)
func (a Gop_bigint) Gop_LshAssign(n Gop_ninteger) {
	a.Int.Lsh(a.Int, uint(n))
}

// Gop_Rsh: func (a bigint) >>= (n untyped_uint)
func (a Gop_bigint) Gop_RshAssign(n Gop_ninteger) {
	a.Int.Rsh(a.Int, uint(n))
}

// Gop_bigint_Cast: func bigint() bigint
func Gop_bigint_Cast__0() Gop_bigint {
	return Gop_bigint{new(big.Int)}
}

// Gop_bigint_Cast: func bigint(x int64) bigint
func Gop_bigint_Cast__1(x int64) Gop_bigint {
	return Gop_bigint{big.NewInt(x)}
}

// Gop_bigint_Cast: func bigint(x int) bigint
func Gop_bigint_Cast__2(x int) Gop_bigint {
	return Gop_bigint{big.NewInt(int64(x))}
}

// Gop_bigint_Cast: func bigint(x uint64) bigint
func Gop_bigint_Cast__3(x uint64) Gop_bigint {
	return Gop_bigint{new(big.Int).SetUint64(x)}
}

// Gop_bigint_Cast: func bigint(x uint) bigint
func Gop_bigint_Cast__4(x uint) Gop_bigint {
	return Gop_bigint{new(big.Int).SetUint64(uint64(x))}
}

// Gop_bigint_Cast: func bigint(x *big.Int) bigint
func Gop_bigint_Cast__5(x *big.Int) Gop_bigint {
	return Gop_bigint{x}
}

// Gop_bigint_Cast: func bigint(x *big.Rat) bigint
func Gop_bigint_Cast__6(x *big.Rat) Gop_bigint {
	if x.IsInt() {
		return Gop_bigint{x.Num()}
	}
	ret, _ := new(big.Float).SetRat(x).Int(nil)
	return Gop_bigint{ret}
}

// Gop_bigint_Init: func bigint.init(x int) bigint
func Gop_bigint_Init__0(x int) Gop_bigint {
	return Gop_bigint{big.NewInt(int64(x))}
}

// Gop_bigint_Init: func bigint.init(x *big.Int) bigint
func Gop_bigint_Init__1(x *big.Int) Gop_bigint {
	return Gop_bigint{x}
}

// Gop_bigint_Init: func bigint.init(x *big.Rat) bigint
func Gop_bigint_Init__2(x *big.Rat) Gop_bigint {
	if x.IsInt() {
		return Gop_bigint{x.Num()}
	}
	panic("TODO: can't init bigint from bigrat")
}

// -----------------------------------------------------------------------------
// type bigrat

// A Gop_bigrat represents a quotient a/b of arbitrary precision.
// The zero value for a Gop_bigrat represents nil.
type Gop_bigrat struct {
	*big.Rat
}

func tmprat(a, b Gop_bigrat) Gop_bigrat {
	if Gop_istmp(a) {
		return a
	} else if Gop_istmp(b) {
		return b
	}
	return Gop_bigrat{new(big.Rat)}
}

func tmprat1(a Gop_bigrat) Gop_bigrat {
	if Gop_istmp(a) {
		return a
	}
	return Gop_bigrat{new(big.Rat)}
}

// IsNil returns a bigrat object is nil or not
func (a Gop_bigrat) IsNil() bool {
	return a.Rat == nil
}

// Gop_Assign: func (a bigrat) = (b bigrat)
func (a Gop_bigrat) Gop_Assign(b Gop_bigrat) {
	if Gop_istmp(b) {
		*a.Rat = *b.Rat
	} else {
		a.Rat.Set(b.Rat)
	}
}

// Gop_Add: func (a bigrat) + (b bigrat) bigrat
func (a Gop_bigrat) Gop_Add(b Gop_bigrat) Gop_bigrat {
	return Gop_bigrat{tmprat(a, b).Add(a.Rat, b.Rat)}
}

// Gop_Sub: func (a bigrat) - (b bigrat) bigrat
func (a Gop_bigrat) Gop_Sub(b Gop_bigrat) Gop_bigrat {
	return Gop_bigrat{tmprat(a, b).Sub(a.Rat, b.Rat)}
}

// Gop_Mul: func (a bigrat) * (b bigrat) bigrat
func (a Gop_bigrat) Gop_Mul(b Gop_bigrat) Gop_bigrat {
	return Gop_bigrat{tmprat(a, b).Mul(a.Rat, b.Rat)}
}

// Gop_Quo: func (a bigrat) / (b bigrat) bigrat
func (a Gop_bigrat) Gop_Quo(b Gop_bigrat) Gop_bigrat {
	return Gop_bigrat{tmprat(a, b).Quo(a.Rat, b.Rat)}
}

// Gop_LT: func (a bigrat) < (b bigrat) bool
func (a Gop_bigrat) Gop_LT(b Gop_bigrat) bool {
	return a.Cmp(b.Rat) < 0
}

// Gop_LE: func (a bigrat) <= (b bigrat) bool
func (a Gop_bigrat) Gop_LE(b Gop_bigrat) bool {
	return a.Cmp(b.Rat) <= 0
}

// Gop_GT: func (a bigrat) > (b bigrat) bool
func (a Gop_bigrat) Gop_GT(b Gop_bigrat) bool {
	return a.Cmp(b.Rat) > 0
}

// Gop_GE: func (a bigrat) >= (b bigrat) bool
func (a Gop_bigrat) Gop_GE(b Gop_bigrat) bool {
	return a.Cmp(b.Rat) >= 0
}

// Gop_EQ: func (a bigrat) == (b bigrat) bool
func (a Gop_bigrat) Gop_EQ(b Gop_bigrat) bool {
	return a.Cmp(b.Rat) == 0
}

// Gop_NE: func (a bigrat) != (b bigrat) bool
func (a Gop_bigrat) Gop_NE(b Gop_bigrat) bool {
	return a.Cmp(b.Rat) != 0
}

// Gop_Neg: func -(a bigrat) bigrat
func (a Gop_bigrat) Gop_Neg() Gop_bigrat {
	return Gop_bigrat{tmprat1(a).Neg(a.Rat)}
}

// Gop_Pos: func +(a bigrat) bigrat
func (a Gop_bigrat) Gop_Pos() Gop_bigrat {
	return a
}

// Gop_Inv: func /(a bigrat) bigrat
func (a Gop_bigrat) Gop_Inv() Gop_bigrat {
	return Gop_bigrat{tmprat1(a).Inv(a.Rat)}
}

// Gop_Add: func (a bigrat) += (b bigrat)
func (a Gop_bigrat) Gop_AddAssign(b Gop_bigrat) {
	a.Rat.Add(a.Rat, b.Rat)
}

// Gop_Sub: func (a bigrat) -= (b bigrat)
func (a Gop_bigrat) Gop_SubAssign(b Gop_bigrat) {
	a.Rat.Sub(a.Rat, b.Rat)
}

// Gop_Mul: func (a bigrat) *= (b bigrat)
func (a Gop_bigrat) Gop_MulAssign(b Gop_bigrat) {
	a.Rat.Mul(a.Rat, b.Rat)
}

// Gop_Quo: func (a bigrat) /= (b bigrat)
func (a Gop_bigrat) Gop_QuoAssign(b Gop_bigrat) {
	a.Rat.Quo(a.Rat, b.Rat)
}

// Gop_bigrat_Cast: func bigrat() bigrat
func Gop_bigrat_Cast__0() Gop_bigrat {
	return Gop_bigrat{new(big.Rat)}
}

// Gop_bigrat_Cast: func bigrat(a bigint) bigrat
func Gop_bigrat_Cast__1(a Gop_bigint) Gop_bigrat {
	return Gop_bigrat{new(big.Rat).SetInt(a.Int)}
}

// Gop_bigrat_Cast: func bigrat(a *big.Int) bigrat
func Gop_bigrat_Cast__2(a *big.Int) Gop_bigrat {
	return Gop_bigrat{new(big.Rat).SetInt(a)}
}

// Gop_bigrat_Cast: func bigrat(a, b int64) bigrat
func Gop_bigrat_Cast__3(a, b int64) Gop_bigrat {
	return Gop_bigrat{big.NewRat(a, b)}
}

// Gop_bigrat_Cast: func bigrat(a *big.Rat) bigrat
func Gop_bigrat_Cast__4(a *big.Rat) Gop_bigrat {
	return Gop_bigrat{a}
}

// Gop_bigrat_Init: func bigrat.init(x untyped_int) bigrat
func Gop_bigrat_Init__0(x int) Gop_bigrat {
	return Gop_bigrat{big.NewRat(int64(x), 1)}
}

// Gop_bigrat_Init: func bigrat.init(x untyped_bigint) bigrat
func Gop_bigrat_Init__1(x Gop_untyped_bigint) Gop_bigrat {
	return Gop_bigrat{new(big.Rat).SetInt(x)}
}

// Gop_bigrat_Init: func bigrat.init(x *big.Rat) bigrat
func Gop_bigrat_Init__2(x *big.Rat) Gop_bigrat {
	return Gop_bigrat{x}
}

// -----------------------------------------------------------------------------
// type bigfloat

// A Gop_bigfloat represents a multi-precision floating point number.
// The zero value for a Gop_bigfloat represents nil.
type Gop_bigfloat struct {
	*big.Float
}

// -----------------------------------------------------------------------------
