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

type untyped_uint = uint

// -----------------------------------------------------------------------------
// type bigint

// A Gope_bigint represents a signed multi-precision integer.
// The zero value for a Gope_bigint represents nil.
type Gope_bigint = *big.Int

func tmpint(a, b Gope_bigint) Gope_bigint {
	return new(big.Int)
}

func tmpint1(a Gope_bigint) Gope_bigint {
	return new(big.Int)
}

// Gopo_bigint__Add: func (a bigint) + (b bigint) bigint
func Gopo_bigint__Add(a, b Gope_bigint) Gope_bigint {
	return tmpint(a, b).Add(a, b)
}

// Gopo_bigint__Sub: func (a bigint) - (b bigint) bigint
func Gopo_bigint__Sub(a, b Gope_bigint) Gope_bigint {
	return tmpint(a, b).Sub(a, b)
}

// Gopo_bigint__Mul: func (a bigint) * (b bigint) bigint
func Gopo_bigint__Mul(a, b Gope_bigint) Gope_bigint {
	return tmpint(a, b).Mul(a, b)
}

// Gopo_bigint__Quo: func (a bigint) / (b bigint) bigint {
func Gopo_bigint__Quo(a, b Gope_bigint) Gope_bigint {
	return tmpint(a, b).Quo(a, b)
}

// Gopo_bigint__Rem: func (a bigint) % (b bigint) bigint
func Gopo_bigint__Rem(a, b Gope_bigint) Gope_bigint {
	return tmpint(a, b).Rem(a, b)
}

// Gopo_bigint__Or: func (a bigint) | (b bigint) bigint
func Gopo_bigint__Or(a, b Gope_bigint) Gope_bigint {
	return tmpint(a, b).Or(a, b)
}

// Gopo_bigint__Xor: func (a bigint) ^ (b bigint) bigint
func Gopo_bigint__Xor(a, b Gope_bigint) Gope_bigint {
	return tmpint(a, b).Xor(a, b)
}

// Gopo_bigint__And: func (a bigint) & (b bigint) bigint
func Gopo_bigint__And(a, b Gope_bigint) Gope_bigint {
	return tmpint(a, b).And(a, b)
}

// Gopo_bigint__AndNot: func (a bigint) &^ (b bigint) bigint
func Gopo_bigint__AndNot(a, b Gope_bigint) Gope_bigint {
	return tmpint(a, b).AndNot(a, b)
}

// Gopo_bigint__Lsh: func (a bigint) << (n untyped_uint) bigint
func Gopo_bigint__Lsh(a Gope_bigint, n untyped_uint) Gope_bigint {
	return tmpint1(a).Lsh(a, n)
}

// Gopo_bigint__Rsh: func (a bigint) >> (n untyped_uint) bigint
func Gopo_bigint__Rsh(a Gope_bigint, n untyped_uint) Gope_bigint {
	return tmpint1(a).Rsh(a, n)
}

// Gopo_bigint__LT: func (a bigint) < (b bigint) bool
func Gopo_bigint__LT(a, b Gope_bigint) bool {
	return a.Cmp(b) < 0
}

// Gopo_bigint__LE: func (a bigint) <= (b bigint) bool
func Gopo_bigint__LE(a, b Gope_bigint) bool {
	return a.Cmp(b) <= 0
}

// Gopo_bigint__GT: func (a bigint) > (b bigint) bool
func Gopo_bigint__GT(a, b Gope_bigint) bool {
	return a.Cmp(b) > 0
}

// Gopo_bigint__GE: func (a bigint) >= (b bigint) bool
func Gopo_bigint__GE(a, b Gope_bigint) bool {
	return a.Cmp(b) >= 0
}

// Gopo_bigint__EQ: func (a bigint) == (b bigint) bool
func Gopo_bigint__EQ(a, b Gope_bigint) bool {
	return a.Cmp(b) == 0
}

// Gopo_bigint__NE: func (a bigint) != (b bigint) bool
func Gopo_bigint__NE(a, b Gope_bigint) bool {
	return a.Cmp(b) != 0
}

// Gopo_bigint__Neg: func -(a bigint) bigint
func Gopo_bigint__Neg(a Gope_bigint) Gope_bigint {
	return tmpint1(a).Neg(a)
}

// Gopo_bigint__Not: func ^(a bigint) bigint
func Gopo_bigint__Not(a Gope_bigint) Gope_bigint {
	return tmpint1(a).Not(a)
}

// Gopc_bigint__0: func bigint() bigint
func Gopc_bigint__0() Gope_bigint {
	return new(big.Int)
}

// Gopc_bigint__1: func bigint(x int64) bigint
func Gopc_bigint__1(x int64) Gope_bigint {
	return big.NewInt(x)
}

// -----------------------------------------------------------------------------
// type bigrat

// A Gope_bigrat represents a quotient a/b of arbitrary precision.
// The zero value for a Gope_bigrat represents nil.
type Gope_bigrat = *big.Rat

func tmprat(a, b Gope_bigrat) Gope_bigrat {
	return new(big.Rat)
}

func tmprat1(a Gope_bigrat) Gope_bigrat {
	return new(big.Rat)
}

// Gopo_bigrat__Add: func (a bigrat) + (b bigrat) bigrat
func Gopo_bigrat__Add(a, b Gope_bigrat) Gope_bigrat {
	return tmprat(a, b).Add(a, b)
}

// Gopo_bigrat__Sub: func (a bigrat) - (b bigrat) bigrat
func Gopo_bigrat__Sub(a, b Gope_bigrat) Gope_bigrat {
	return tmprat(a, b).Sub(a, b)
}

// Gopo_bigrat__Mul: func (a bigrat) * (b bigrat) bigrat
func Gopo_bigrat__Mul(a, b Gope_bigrat) Gope_bigrat {
	return tmprat(a, b).Mul(a, b)
}

// Gopo_bigrat__Quo: func (a bigrat) / (b bigrat) bigrat
func Gopo_bigrat__Quo(a, b Gope_bigrat) Gope_bigrat {
	return tmprat(a, b).Quo(a, b)
}

// Gopo_bigrat__LT: func (a bigrat) < (b bigrat) bool
func Gopo_bigrat__LT(a, b Gope_bigrat) bool {
	return a.Cmp(b) < 0
}

// Gopo_bigrat__LE: func (a bigrat) <= (b bigrat) bool
func Gopo_bigrat__LE(a, b Gope_bigrat) bool {
	return a.Cmp(b) <= 0
}

// Gopo_bigrat__GT: func (a bigrat) > (b bigrat) bool
func Gopo_bigrat__GT(a, b Gope_bigrat) bool {
	return a.Cmp(b) > 0
}

// Gopo_bigrat__GE: func (a bigrat) >= (b bigrat) bool
func Gopo_bigrat__GE(a, b Gope_bigrat) bool {
	return a.Cmp(b) >= 0
}

// Gopo_bigrat__EQ: func (a bigrat) == (b bigrat) bool
func Gopo_bigrat__EQ(a, b Gope_bigrat) bool {
	return a.Cmp(b) == 0
}

// Gopo_bigrat__NE: func (a bigrat) != (b bigrat) bool
func Gopo_bigrat__NE(a, b Gope_bigrat) bool {
	return a.Cmp(b) != 0
}

// Gopo_bigrat__Neg: func -(a bigrat) bigrat
func Gopo_bigrat__Neg(a Gope_bigrat) Gope_bigrat {
	return tmprat1(a).Neg(a)
}

// Gopo_bigrat__Inv: func /(a bigrat) bigrat
func Gopo_bigrat__Inv(a Gope_bigrat) Gope_bigrat {
	return tmprat1(a).Inv(a)
}

// Gopc_bigrat__0: func bigrat() bigrat
func Gopc_bigrat__0() Gope_bigrat {
	return new(big.Rat)
}

// Gopc_bigrat__1: func bigrat(a bigint) bigrat
func Gopc_bigrat__1(a Gope_bigint) Gope_bigrat {
	return new(big.Rat).SetInt(a)
}

// Gopc_bigrat__2: func bigrat(a, b int64) bigrat
func Gopc_bigrat__2(a, b int64) Gope_bigrat {
	return big.NewRat(a, b)
}

// -----------------------------------------------------------------------------
// type bigfloat

// A Gope_bigfloat represents a multi-precision floating point number.
// The zero value for a Gope_bigfloat represents nil.
type Gope_bigfloat = *big.Float

// -----------------------------------------------------------------------------
