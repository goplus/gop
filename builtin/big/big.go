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

package big

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

type UntypedInt *big.Int
type UntypedRat *big.Rat
type UntypedFloat *big.Float

type UntypedInt_Default = Int
type UntypedRat_Default = Rat
type UntypedFloat_Default = Float

func UntypedInt_Init__0(x int) UntypedInt {
	panic("make compiler happy")
}

func UntypedRat_Init__0(x int) UntypedRat {
	panic("make compiler happy")
}

func UntypedRat_Init__1(x UntypedInt) UntypedRat {
	panic("make compiler happy")
}

// -----------------------------------------------------------------------------
// type bigint

// A Int represents a signed multi-precision integer.
// The zero value for a Int represents nil.
type Int struct {
	*big.Int
}

func tmpint(a, b Int) Int {
	if Gop_istmp(a) {
		return a
	} else if Gop_istmp(b) {
		return b
	}
	return Int{new(big.Int)}
}

func tmpint1(a Int) Int {
	if Gop_istmp(a) {
		return a
	}
	return Int{new(big.Int)}
}

// IsNil returns a bigint object is nil or not
func (a Int) IsNil() bool {
	return a.Int == nil
}

// Gop_Assign: func (a bigint) = (b bigint)
func (a Int) Gop_Assign(b Int) {
	if Gop_istmp(b) {
		*a.Int = *b.Int
	} else {
		a.Int.Set(b.Int)
	}
}

// Gop_Add: func (a bigint) + (b bigint) bigint
func (a Int) Gop_Add(b Int) Int {
	return Int{tmpint(a, b).Add(a.Int, b.Int)}
}

// Gop_Sub: func (a bigint) - (b bigint) bigint
func (a Int) Gop_Sub(b Int) Int {
	return Int{tmpint(a, b).Sub(a.Int, b.Int)}
}

// Gop_Mul: func (a bigint) * (b bigint) bigint
func (a Int) Gop_Mul(b Int) Int {
	return Int{tmpint(a, b).Mul(a.Int, b.Int)}
}

// Gop_Quo: func (a bigint) / (b bigint) bigint {
func (a Int) Gop_Quo(b Int) Int {
	return Int{tmpint(a, b).Quo(a.Int, b.Int)}
}

// Gop_Rem: func (a bigint) % (b bigint) bigint
func (a Int) Gop_Rem(b Int) Int {
	return Int{tmpint(a, b).Rem(a.Int, b.Int)}
}

// Gop_Or: func (a bigint) | (b bigint) bigint
func (a Int) Gop_Or(b Int) Int {
	return Int{tmpint(a, b).Or(a.Int, b.Int)}
}

// Gop_Xor: func (a bigint) ^ (b bigint) bigint
func (a Int) Gop_Xor(b Int) Int {
	return Int{tmpint(a, b).Xor(a.Int, b.Int)}
}

// Gop_And: func (a bigint) & (b bigint) bigint
func (a Int) Gop_And(b Int) Int {
	return Int{tmpint(a, b).And(a.Int, b.Int)}
}

// Gop_AndNot: func (a bigint) &^ (b bigint) bigint
func (a Int) Gop_AndNot(b Int) Int {
	return Int{tmpint(a, b).AndNot(a.Int, b.Int)}
}

// Gop_Lsh: func (a bigint) << (n untyped_uint) bigint
func (a Int) Gop_Lsh(n Gop_ninteger) Int {
	return Int{tmpint1(a).Lsh(a.Int, uint(n))}
}

// Gop_Rsh: func (a bigint) >> (n untyped_uint) bigint
func (a Int) Gop_Rsh(n Gop_ninteger) Int {
	return Int{tmpint1(a).Rsh(a.Int, uint(n))}
}

// Gop_LT: func (a bigint) < (b bigint) bool
func (a Int) Gop_LT(b Int) bool {
	return a.Cmp(b.Int) < 0
}

// Gop_LE: func (a bigint) <= (b bigint) bool
func (a Int) Gop_LE(b Int) bool {
	return a.Cmp(b.Int) <= 0
}

// Gop_GT: func (a bigint) > (b bigint) bool
func (a Int) Gop_GT(b Int) bool {
	return a.Cmp(b.Int) > 0
}

// Gop_GE: func (a bigint) >= (b bigint) bool
func (a Int) Gop_GE(b Int) bool {
	return a.Cmp(b.Int) >= 0
}

// Gop_EQ: func (a bigint) == (b bigint) bool
func (a Int) Gop_EQ(b Int) bool {
	return a.Cmp(b.Int) == 0
}

// Gop_NE: func (a bigint) != (b bigint) bool
func (a Int) Gop_NE(b Int) bool {
	return a.Cmp(b.Int) != 0
}

// Gop_Neg: func -(a bigint) bigint
func (a Int) Gop_Neg() Int {
	return Int{tmpint1(a).Neg(a.Int)}
}

// Gop_Pos: func +(a bigint) bigint
func (a Int) Gop_Pos() Int {
	return a
}

// Gop_Not: func ^(a bigint) bigint
func (a Int) Gop_Not() Int {
	return Int{tmpint1(a).Not(a.Int)}
}

// Gop_Add: func (a bigint) += (b bigint)
func (a Int) Gop_AddAssign(b Int) {
	a.Int.Add(a.Int, b.Int)
}

// Gop_Sub: func (a bigint) -= (b bigint)
func (a Int) Gop_SubAssign(b Int) {
	a.Int.Sub(a.Int, b.Int)
}

// Gop_Mul: func (a bigint) *= (b bigint)
func (a Int) Gop_MulAssign(b Int) {
	a.Int.Mul(a.Int, b.Int)
}

// Gop_Quo: func (a bigint) /= (b bigint) {
func (a Int) Gop_QuoAssign(b Int) {
	a.Int.Quo(a.Int, b.Int)
}

// Gop_Rem: func (a bigint) %= (b bigint)
func (a Int) Gop_RemAssign(b Int) {
	a.Int.Rem(a.Int, b.Int)
}

// Gop_Or: func (a bigint) |= (b bigint)
func (a Int) Gop_OrAssign(b Int) {
	a.Int.Or(a.Int, b.Int)
}

// Gop_Xor: func (a bigint) ^= (b bigint)
func (a Int) Gop_XorAssign(b Int) {
	a.Int.Xor(a.Int, b.Int)
}

// Gop_And: func (a bigint) &= (b bigint)
func (a Int) Gop_AndAssign(b Int) {
	a.Int.And(a.Int, b.Int)
}

// Gop_AndNot: func (a bigint) &^= (b bigint)
func (a Int) Gop_AndNotAssign(b Int) {
	a.Int.AndNot(a.Int, b.Int)
}

// Gop_Lsh: func (a bigint) <<= (n untyped_uint)
func (a Int) Gop_LshAssign(n Gop_ninteger) {
	a.Int.Lsh(a.Int, uint(n))
}

// Gop_Rsh: func (a bigint) >>= (n untyped_uint)
func (a Int) Gop_RshAssign(n Gop_ninteger) {
	a.Int.Rsh(a.Int, uint(n))
}

// Int_Cast: func bigint() bigint
func Int_Cast__0() Int {
	return Int{new(big.Int)}
}

// Int_Cast: func bigint(x int64) bigint
func Int_Cast__1(x int64) Int {
	return Int{big.NewInt(x)}
}

// Int_Cast: func bigint(x int) bigint
func Int_Cast__2(x int) Int {
	return Int{big.NewInt(int64(x))}
}

// Int_Cast: func bigint(x uint64) bigint
func Int_Cast__3(x uint64) Int {
	return Int{new(big.Int).SetUint64(x)}
}

// Int_Cast: func bigint(x uint) bigint
func Int_Cast__4(x uint) Int {
	return Int{new(big.Int).SetUint64(uint64(x))}
}

// Int_Cast: func bigint(x *big.Int) bigint
func Int_Cast__5(x *big.Int) Int {
	return Int{x}
}

// Int_Cast: func bigint(x *big.Rat) bigint
func Int_Cast__6(x *big.Rat) Int {
	if x.IsInt() {
		return Int{x.Num()}
	}
	ret, _ := new(big.Float).SetRat(x).Int(nil)
	return Int{ret}
}

// Int_Init: func bigint.init(x int) bigint
func Int_Init__0(x int) Int {
	return Int{big.NewInt(int64(x))}
}

// Int_Init: func bigint.init(x *big.Int) bigint
func Int_Init__1(x *big.Int) Int {
	return Int{x}
}

// Int_Init: func bigint.init(x *big.Rat) bigint
func Int_Init__2(x *big.Rat) Int {
	if x.IsInt() {
		return Int{x.Num()}
	}
	panic("TODO: can't init bigint from bigrat")
}

// -----------------------------------------------------------------------------
// type bigrat

// A Rat represents a quotient a/b of arbitrary precision.
// The zero value for a Rat represents nil.
type Rat struct {
	*big.Rat
}

func tmprat(a, b Rat) Rat {
	if Gop_istmp(a) {
		return a
	} else if Gop_istmp(b) {
		return b
	}
	return Rat{new(big.Rat)}
}

func tmprat1(a Rat) Rat {
	if Gop_istmp(a) {
		return a
	}
	return Rat{new(big.Rat)}
}

// IsNil returns a bigrat object is nil or not
func (a Rat) IsNil() bool {
	return a.Rat == nil
}

// Gop_Assign: func (a bigrat) = (b bigrat)
func (a Rat) Gop_Assign(b Rat) {
	if Gop_istmp(b) {
		*a.Rat = *b.Rat
	} else {
		a.Rat.Set(b.Rat)
	}
}

// Gop_Add: func (a bigrat) + (b bigrat) bigrat
func (a Rat) Gop_Add(b Rat) Rat {
	return Rat{tmprat(a, b).Add(a.Rat, b.Rat)}
}

// Gop_Sub: func (a bigrat) - (b bigrat) bigrat
func (a Rat) Gop_Sub(b Rat) Rat {
	return Rat{tmprat(a, b).Sub(a.Rat, b.Rat)}
}

// Gop_Mul: func (a bigrat) * (b bigrat) bigrat
func (a Rat) Gop_Mul(b Rat) Rat {
	return Rat{tmprat(a, b).Mul(a.Rat, b.Rat)}
}

// Gop_Quo: func (a bigrat) / (b bigrat) bigrat
func (a Rat) Gop_Quo(b Rat) Rat {
	return Rat{tmprat(a, b).Quo(a.Rat, b.Rat)}
}

// Gop_LT: func (a bigrat) < (b bigrat) bool
func (a Rat) Gop_LT(b Rat) bool {
	return a.Cmp(b.Rat) < 0
}

// Gop_LE: func (a bigrat) <= (b bigrat) bool
func (a Rat) Gop_LE(b Rat) bool {
	return a.Cmp(b.Rat) <= 0
}

// Gop_GT: func (a bigrat) > (b bigrat) bool
func (a Rat) Gop_GT(b Rat) bool {
	return a.Cmp(b.Rat) > 0
}

// Gop_GE: func (a bigrat) >= (b bigrat) bool
func (a Rat) Gop_GE(b Rat) bool {
	return a.Cmp(b.Rat) >= 0
}

// Gop_EQ: func (a bigrat) == (b bigrat) bool
func (a Rat) Gop_EQ(b Rat) bool {
	return a.Cmp(b.Rat) == 0
}

// Gop_NE: func (a bigrat) != (b bigrat) bool
func (a Rat) Gop_NE(b Rat) bool {
	return a.Cmp(b.Rat) != 0
}

// Gop_Neg: func -(a bigrat) bigrat
func (a Rat) Gop_Neg() Rat {
	return Rat{tmprat1(a).Neg(a.Rat)}
}

// Gop_Pos: func +(a bigrat) bigrat
func (a Rat) Gop_Pos() Rat {
	return a
}

// Gop_Inv: func /(a bigrat) bigrat
func (a Rat) Gop_Inv() Rat {
	return Rat{tmprat1(a).Inv(a.Rat)}
}

// Gop_Add: func (a bigrat) += (b bigrat)
func (a Rat) Gop_AddAssign(b Rat) {
	a.Rat.Add(a.Rat, b.Rat)
}

// Gop_Sub: func (a bigrat) -= (b bigrat)
func (a Rat) Gop_SubAssign(b Rat) {
	a.Rat.Sub(a.Rat, b.Rat)
}

// Gop_Mul: func (a bigrat) *= (b bigrat)
func (a Rat) Gop_MulAssign(b Rat) {
	a.Rat.Mul(a.Rat, b.Rat)
}

// Gop_Quo: func (a bigrat) /= (b bigrat)
func (a Rat) Gop_QuoAssign(b Rat) {
	a.Rat.Quo(a.Rat, b.Rat)
}

// Rat_Cast: func bigrat() bigrat
func Rat_Cast__0() Rat {
	return Rat{new(big.Rat)}
}

// Rat_Cast: func bigrat(a bigint) bigrat
func Rat_Cast__1(a Int) Rat {
	return Rat{new(big.Rat).SetInt(a.Int)}
}

// Rat_Cast: func bigrat(a *big.Int) bigrat
func Rat_Cast__2(a *big.Int) Rat {
	return Rat{new(big.Rat).SetInt(a)}
}

// Rat_Cast: func bigrat(a, b int64) bigrat
func Rat_Cast__3(a, b int64) Rat {
	return Rat{big.NewRat(a, b)}
}

// Rat_Cast: func bigrat(a *big.Rat) bigrat
func Rat_Cast__4(a *big.Rat) Rat {
	return Rat{a}
}

// Rat_Init: func bigrat.init(x untyped_int) bigrat
func Rat_Init__0(x int) Rat {
	return Rat{big.NewRat(int64(x), 1)}
}

// Rat_Init: func bigrat.init(x untyped_bigint) bigrat
func Rat_Init__1(x UntypedInt) Rat {
	return Rat{new(big.Rat).SetInt(x)}
}

// Rat_Init: func bigrat.init(x *big.Rat) bigrat
func Rat_Init__2(x *big.Rat) Rat {
	return Rat{x}
}

// -----------------------------------------------------------------------------
// type bigfloat

// A Float represents a multi-precision floating point number.
// The zero value for a Float represents nil.
type Float struct {
	*big.Float
}

// -----------------------------------------------------------------------------
