/*
 Copyright 2021 The GoPlus Authors (goplus.org)
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

package builtin

import (
	"math/big"
)

//
// Gopb_: builtin type (of Go+) prefix
// Gope_: type extend prefix
// Gopo_: operator prefix
// Gopc_: type convert prefix
//
// xxxx__N: the Nth overload function
//

type Gopb_ninteger uint

// -----------------------------------------------------------------------------
// type bigint

// A Gope_bigint represents a signed multi-precision integer.
// The zero value for a Gope_bigint represents nil.
type Gope_bigint struct {
	*big.Int
}

func tmpint(a, b Gope_bigint) Gope_bigint {
	return Gope_bigint{new(big.Int)}
}

func tmpint1(a Gope_bigint) Gope_bigint {
	return Gope_bigint{new(big.Int)}
}

// IsNil returns a bigint object is nil or not
func (a Gope_bigint) IsNil() bool {
	return a.Int == nil
}

// Gopo_Add: func (a bigint) + (b bigint) bigint
func (a Gope_bigint) Gopo_Add(b Gope_bigint) Gope_bigint {
	return Gope_bigint{tmpint(a, b).Add(a.Int, b.Int)}
}

// Gopo_Sub: func (a bigint) - (b bigint) bigint
func (a Gope_bigint) Gopo_Sub(b Gope_bigint) Gope_bigint {
	return Gope_bigint{tmpint(a, b).Sub(a.Int, b.Int)}
}

// Gopo_Mul: func (a bigint) * (b bigint) bigint
func (a Gope_bigint) Gopo_Mul(b Gope_bigint) Gope_bigint {
	return Gope_bigint{tmpint(a, b).Mul(a.Int, b.Int)}
}

// Gopo_Quo: func (a bigint) / (b bigint) bigint {
func (a Gope_bigint) Gopo_Quo(b Gope_bigint) Gope_bigint {
	return Gope_bigint{tmpint(a, b).Quo(a.Int, b.Int)}
}

// Gopo_Rem: func (a bigint) % (b bigint) bigint
func (a Gope_bigint) Gopo_Rem(b Gope_bigint) Gope_bigint {
	return Gope_bigint{tmpint(a, b).Rem(a.Int, b.Int)}
}

// Gopo_Or: func (a bigint) | (b bigint) bigint
func (a Gope_bigint) Gopo_Or(b Gope_bigint) Gope_bigint {
	return Gope_bigint{tmpint(a, b).Or(a.Int, b.Int)}
}

// Gopo_Xor: func (a bigint) ^ (b bigint) bigint
func (a Gope_bigint) Gopo_Xor(b Gope_bigint) Gope_bigint {
	return Gope_bigint{tmpint(a, b).Xor(a.Int, b.Int)}
}

// Gopo_And: func (a bigint) & (b bigint) bigint
func (a Gope_bigint) Gopo_And(b Gope_bigint) Gope_bigint {
	return Gope_bigint{tmpint(a, b).And(a.Int, b.Int)}
}

// Gopo_AndNot: func (a bigint) &^ (b bigint) bigint
func (a Gope_bigint) Gopo_AndNot(b Gope_bigint) Gope_bigint {
	return Gope_bigint{tmpint(a, b).AndNot(a.Int, b.Int)}
}

// Gopo_Lsh: func (a bigint) << (n untyped_uint) bigint
func (a Gope_bigint) Gopo_Lsh(n Gopb_ninteger) Gope_bigint {
	return Gope_bigint{tmpint1(a).Lsh(a.Int, uint(n))}
}

// Gopo_Rsh: func (a bigint) >> (n untyped_uint) bigint
func (a Gope_bigint) Gopo_Rsh(n Gopb_ninteger) Gope_bigint {
	return Gope_bigint{tmpint1(a).Rsh(a.Int, uint(n))}
}

// Gopo_LT: func (a bigint) < (b bigint) bool
func (a Gope_bigint) Gopo_LT(b Gope_bigint) bool {
	return a.Cmp(b.Int) < 0
}

// (a Gope_bigint) Gopo_LE: func (a bigint) <= (b bigint) bool
func (a Gope_bigint) Gopo_LE(b Gope_bigint) bool {
	return a.Cmp(b.Int) <= 0
}

// Gopo_GT: func (a bigint) > (b bigint) bool
func (a Gope_bigint) Gopo_GT(b Gope_bigint) bool {
	return a.Cmp(b.Int) > 0
}

// Gopo_GE: func (a bigint) >= (b bigint) bool
func (a Gope_bigint) Gopo_GE(b Gope_bigint) bool {
	return a.Cmp(b.Int) >= 0
}

// Gopo_EQ: func (a bigint) == (b bigint) bool
func (a Gope_bigint) Gopo_EQ(b Gope_bigint) bool {
	return a.Cmp(b.Int) == 0
}

// Gopo_NE: func (a bigint) != (b bigint) bool
func (a Gope_bigint) Gopo_NE(b Gope_bigint) bool {
	return a.Cmp(b.Int) != 0
}

// Gopo_Neg: func -(a bigint) bigint
func (a Gope_bigint) Gopo_Neg() Gope_bigint {
	return Gope_bigint{tmpint1(a).Neg(a.Int)}
}

// Gopo_Not: func ^(a bigint) bigint
func (a Gope_bigint) Gopo_Not() Gope_bigint {
	return Gope_bigint{tmpint1(a).Not(a.Int)}
}

// Gopc_bigint__0: func bigint() bigint
func Gopc_bigint__0() Gope_bigint {
	return Gope_bigint{new(big.Int)}
}

// Gopc_bigint__1: func bigint(x int64) bigint
func Gopc_bigint__1(x int64) Gope_bigint {
	return Gope_bigint{big.NewInt(x)}
}

// -----------------------------------------------------------------------------
// type bigrat

// A Gope_bigrat represents a quotient a/b of arbitrary precision.
// The zero value for a Gope_bigrat represents nil.
type Gope_bigrat struct {
	*big.Rat
}

func tmprat(a, b Gope_bigrat) Gope_bigrat {
	return Gope_bigrat{new(big.Rat)}
}

func tmprat1(a Gope_bigrat) Gope_bigrat {
	return Gope_bigrat{new(big.Rat)}
}

// IsNil returns a bigrat object is nil or not
func (a Gope_bigrat) IsNil() bool {
	return a.Rat == nil
}

// Gopo_Add: func (a bigrat) + (b bigrat) bigrat
func (a Gope_bigrat) Gopo_Add(b Gope_bigrat) Gope_bigrat {
	return Gope_bigrat{tmprat(a, b).Add(a.Rat, b.Rat)}
}

// Gopo_Sub: func (a bigrat) - (b bigrat) bigrat
func (a Gope_bigrat) Gopo_Sub(b Gope_bigrat) Gope_bigrat {
	return Gope_bigrat{tmprat(a, b).Sub(a.Rat, b.Rat)}
}

// Gopo_Mul: func (a bigrat) * (b bigrat) bigrat
func (a Gope_bigrat) Gopo_Mul(b Gope_bigrat) Gope_bigrat {
	return Gope_bigrat{tmprat(a, b).Mul(a.Rat, b.Rat)}
}

// Gopo_Quo: func (a bigrat) / (b bigrat) bigrat
func (a Gope_bigrat) Gopo_Quo(b Gope_bigrat) Gope_bigrat {
	return Gope_bigrat{tmprat(a, b).Quo(a.Rat, b.Rat)}
}

// Gopo_LT: func (a bigrat) < (b bigrat) bool
func (a Gope_bigrat) Gopo_LT(b Gope_bigrat) bool {
	return a.Cmp(b.Rat) < 0
}

// Gopo_LE: func (a bigrat) <= (b bigrat) bool
func (a Gope_bigrat) Gopo_LE(b Gope_bigrat) bool {
	return a.Cmp(b.Rat) <= 0
}

// Gopo_GT: func (a bigrat) > (b bigrat) bool
func (a Gope_bigrat) Gopo_GT(b Gope_bigrat) bool {
	return a.Cmp(b.Rat) > 0
}

// Gopo_GE: func (a bigrat) >= (b bigrat) bool
func (a Gope_bigrat) Gopo_GE(b Gope_bigrat) bool {
	return a.Cmp(b.Rat) >= 0
}

// Gopo_EQ: func (a bigrat) == (b bigrat) bool
func (a Gope_bigrat) Gopo_EQ(b Gope_bigrat) bool {
	return a.Cmp(b.Rat) == 0
}

// Gopo_NE: func (a bigrat) != (b bigrat) bool
func (a Gope_bigrat) Gopo_NE(b Gope_bigrat) bool {
	return a.Cmp(b.Rat) != 0
}

// Gopo_Neg: func -(a bigrat) bigrat
func (a Gope_bigrat) Gopo_Neg() Gope_bigrat {
	return Gope_bigrat{tmprat1(a).Neg(a.Rat)}
}

// Gopo_Inv: func /(a bigrat) bigrat
func (a Gope_bigrat) Gopo_Inv() Gope_bigrat {
	return Gope_bigrat{tmprat1(a).Inv(a.Rat)}
}

// Gopc_bigrat__0: func bigrat() bigrat
func Gopc_bigrat__0() Gope_bigrat {
	return Gope_bigrat{new(big.Rat)}
}

// Gopc_bigrat__1: func bigrat(a bigint) bigrat
func Gopc_bigrat__1(a Gope_bigint) Gope_bigrat {
	return Gope_bigrat{new(big.Rat).SetInt(a.Int)}
}

// Gopc_bigrat__2: func bigrat(a, b int64) bigrat
func Gopc_bigrat__2(a, b int64) Gope_bigrat {
	return Gope_bigrat{big.NewRat(a, b)}
}

// -----------------------------------------------------------------------------
// type bigfloat

// A Gope_bigfloat represents a multi-precision floating point number.
// The zero value for a Gope_bigfloat represents nil.
type Gope_bigfloat struct {
	*big.Float
}

// -----------------------------------------------------------------------------
