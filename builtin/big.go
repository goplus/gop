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

// A Bigint represents a signed multi-precision integer.
// The zero value for a Bigint represents the value 0.
type Bigint = *big.Int

func tmpint(a, b Bigint) Bigint {
	return new(big.Int)
}

func tmpint1(a Bigint) Bigint {
	return new(big.Int)
}

// Gopo_Bigint__Add: func (a bigint) + (b bigint) bigint
func Gopo_Bigint__Add(a, b Bigint) Bigint {
	return tmpint(a, b).Add(a, b)
}

// Gopo_Bigint__Sub: func (a bigint) - (b bigint) bigint
func Gopo_Bigint__Sub(a, b Bigint) Bigint {
	return tmpint(a, b).Sub(a, b)
}

// Gopo_Bigint__Mul: func (a bigint) * (b bigint) bigint
func Gopo_Bigint__Mul(a, b Bigint) Bigint {
	return tmpint(a, b).Mul(a, b)
}

// Gopo_Bigint__Quo: func (a bigint) / (b bigint) bigint {
func Gopo_Bigint__Quo(a, b Bigint) Bigint {
	return tmpint(a, b).Quo(a, b)
}

// Gopo_Bigint__Rem: func (a bigint) % (b bigint) bigint
func Gopo_Bigint__Rem(a, b Bigint) Bigint {
	return tmpint(a, b).Rem(a, b)
}

// Gopo_Bigint__Or: func (a bigint) | (b bigint) bigint
func Gopo_Bigint__Or(a, b Bigint) Bigint {
	return tmpint(a, b).Or(a, b)
}

// Gopo_Bigint__Xor: func (a bigint) ^ (b bigint) bigint
func Gopo_Bigint__Xor(a, b Bigint) Bigint {
	return tmpint(a, b).Xor(a, b)
}

// Gopo_Bigint__And: func (a bigint) & (b bigint) bigint
func Gopo_Bigint__And(a, b Bigint) Bigint {
	return tmpint(a, b).And(a, b)
}

// Gopo_Bigint__AndNot: func (a bigint) &^ (b bigint) bigint
func Gopo_Bigint__AndNot(a, b Bigint) Bigint {
	return tmpint(a, b).AndNot(a, b)
}

// Gopo_Bigint__Lsh: func (a bigint) << (n untyped_uint) bigint
func Gopo_Bigint__Lsh(a Bigint, n untyped_uint) Bigint {
	return tmpint1(a).Lsh(a, n)
}

// Gopo_Bigint__Rsh: func (a bigint) >> (n untyped_uint) bigint
func Gopo_Bigint__Rsh(a Bigint, n untyped_uint) Bigint {
	return tmpint1(a).Rsh(a, n)
}

// Gopo_Bigint__LT: func (a bigint) < (b bigint) bool
func Gopo_Bigint__LT(a, b Bigint) bool {
	return a.Cmp(b) < 0
}

// Gopo_Bigint__LE: func (a bigint) <= (b bigint) bool
func Gopo_Bigint__LE(a, b Bigint) bool {
	return a.Cmp(b) <= 0
}

// Gopo_Bigint__GT: func (a bigint) > (b bigint) bool
func Gopo_Bigint__GT(a, b Bigint) bool {
	return a.Cmp(b) > 0
}

// Gopo_Bigint__GE: func (a bigint) >= (b bigint) bool
func Gopo_Bigint__GE(a, b Bigint) bool {
	return a.Cmp(b) >= 0
}

// Gopo_Bigint__EQ: func (a bigint) == (b bigint) bool
func Gopo_Bigint__EQ(a, b Bigint) bool {
	return a.Cmp(b) == 0
}

// Gopo_Bigint__NE: func (a bigint) != (b bigint) bool
func Gopo_Bigint__NE(a, b Bigint) bool {
	return a.Cmp(b) != 0
}

// Gopo_Bigint__Neg: func -(a bigint) bigint
func Gopo_Bigint__Neg(a Bigint) Bigint {
	return tmpint1(a).Neg(a)
}

// Gopo_Bigint__Not: func ^(a bigint) bigint
func Gopo_Bigint__Not(a Bigint) Bigint {
	return tmpint1(a).Not(a)
}

// Gopc_Bigint__0: func bigint() bigint
func Gopc_Bigint__0() Bigint {
	return new(big.Int)
}

// Gopc_Bigint__1: func bigint(x int64) bigint
func Gopc_Bigint__1(x int64) Bigint {
	return big.NewInt(x)
}

// -----------------------------------------------------------------------------
// type bigrat

type Bigrat = *big.Rat

func tmprat(a, b Bigrat) Bigrat {
	return new(big.Rat)
}

func tmprat1(a Bigrat) Bigrat {
	return new(big.Rat)
}

// Gopo_Bigrat__Add: func (a bigrat) + (b bigrat) bigrat
func Gopo_Bigrat__Add(a, b Bigrat) Bigrat {
	return tmprat(a, b).Add(a, b)
}

// Gopo_Bigrat__Sub: func (a bigrat) - (b bigrat) bigrat
func Gopo_Bigrat__Sub(a, b Bigrat) Bigrat {
	return tmprat(a, b).Sub(a, b)
}

// Gopo_Bigrat__Mul: func (a bigrat) * (b bigrat) bigrat
func Gopo_Bigrat__Mul(a, b Bigrat) Bigrat {
	return tmprat(a, b).Mul(a, b)
}

// Gopo_Bigrat__Quo: func (a bigrat) / (b bigrat) bigrat
func Gopo_Bigrat__Quo(a, b Bigrat) Bigrat {
	return tmprat(a, b).Quo(a, b)
}

// Gopo_Bigrat__LT: func (a bigrat) < (b bigrat) bool
func Gopo_Bigrat__LT(a, b Bigrat) bool {
	return a.Cmp(b) < 0
}

// Gopo_Bigrat__LE: func (a bigrat) <= (b bigrat) bool
func Gopo_Bigrat__LE(a, b Bigrat) bool {
	return a.Cmp(b) <= 0
}

// Gopo_Bigrat__GT: func (a bigrat) > (b bigrat) bool
func Gopo_Bigrat__GT(a, b Bigrat) bool {
	return a.Cmp(b) > 0
}

// Gopo_Bigrat__GE: func (a bigrat) >= (b bigrat) bool
func Gopo_Bigrat__GE(a, b Bigrat) bool {
	return a.Cmp(b) >= 0
}

// Gopo_Bigrat__EQ: func (a bigrat) == (b bigrat) bool
func Gopo_Bigrat__EQ(a, b Bigrat) bool {
	return a.Cmp(b) == 0
}

// Gopo_Bigrat__NE: func (a bigrat) != (b bigrat) bool
func Gopo_Bigrat__NE(a, b Bigrat) bool {
	return a.Cmp(b) != 0
}

// Gopo_Bigrat__Neg: func -(a bigrat) bigrat
func Gopo_Bigrat__Neg(a Bigrat) Bigrat {
	return tmprat1(a).Neg(a)
}

// Gopo_Bigrat__Inv: func /(a bigrat) bigrat
func Gopo_Bigrat__Inv(a Bigrat) Bigrat {
	return tmprat1(a).Inv(a)
}

// Gopc_Bigrat__0: func bigrat() bigrat
func Gopc_Bigrat__0() Bigrat {
	return new(big.Rat)
}

// Gopc_Bigrat__1: func bigrat(a bigint) bigrat
func Gopc_Bigrat__1(a Bigint) Bigrat {
	return new(big.Rat).SetInt(a)
}

// Gopc_Bigrat__2: func bigrat(a, b int64) bigrat
func Gopc_Bigrat__2(a, b int64) Bigrat {
	return big.NewRat(a, b)
}

// -----------------------------------------------------------------------------
// type bigfloat

type Bigfloat = *big.Float

// -----------------------------------------------------------------------------
