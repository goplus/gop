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
// type_Op: type operator
// xxxx__N: the Nth overload function
//

type Gopb_untyped_uint uint

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

// Gopo_bigint_Add: func (a bigint) + (b bigint) bigint
func Gopo_bigint_Add(a, b Gope_bigint) Gope_bigint {
	return tmpint(a, b).Add(a, b)
}

// Gopo_bigint_Sub: func (a bigint) - (b bigint) bigint
func Gopo_bigint_Sub(a, b Gope_bigint) Gope_bigint {
	return tmpint(a, b).Sub(a, b)
}

// Gopo_bigint_Mul: func (a bigint) * (b bigint) bigint
func Gopo_bigint_Mul(a, b Gope_bigint) Gope_bigint {
	return tmpint(a, b).Mul(a, b)
}

// Gopo_bigint_Quo: func (a bigint) / (b bigint) bigint {
func Gopo_bigint_Quo(a, b Gope_bigint) Gope_bigint {
	return tmpint(a, b).Quo(a, b)
}

// Gopo_bigint_Rem: func (a bigint) % (b bigint) bigint
func Gopo_bigint_Rem(a, b Gope_bigint) Gope_bigint {
	return tmpint(a, b).Rem(a, b)
}

// Gopo_bigint_Or: func (a bigint) | (b bigint) bigint
func Gopo_bigint_Or(a, b Gope_bigint) Gope_bigint {
	return tmpint(a, b).Or(a, b)
}

// Gopo_bigint_Xor: func (a bigint) ^ (b bigint) bigint
func Gopo_bigint_Xor(a, b Gope_bigint) Gope_bigint {
	return tmpint(a, b).Xor(a, b)
}

// Gopo_bigint_And: func (a bigint) & (b bigint) bigint
func Gopo_bigint_And(a, b Gope_bigint) Gope_bigint {
	return tmpint(a, b).And(a, b)
}

// Gopo_bigint_AndNot: func (a bigint) &^ (b bigint) bigint
func Gopo_bigint_AndNot(a, b Gope_bigint) Gope_bigint {
	return tmpint(a, b).AndNot(a, b)
}

// Gopo_bigint_Lsh: func (a bigint) << (n untyped_uint) bigint
func Gopo_bigint_Lsh(a Gope_bigint, n Gopb_untyped_uint) Gope_bigint {
	return tmpint1(a).Lsh(a, uint(n))
}

// Gopo_bigint_Rsh: func (a bigint) >> (n untyped_uint) bigint
func Gopo_bigint_Rsh(a Gope_bigint, n Gopb_untyped_uint) Gope_bigint {
	return tmpint1(a).Rsh(a, uint(n))
}

// Gopo_bigint_LT: func (a bigint) < (b bigint) bool
func Gopo_bigint_LT(a, b Gope_bigint) bool {
	return a.Cmp(b) < 0
}

// Gopo_bigint_LE: func (a bigint) <= (b bigint) bool
func Gopo_bigint_LE(a, b Gope_bigint) bool {
	return a.Cmp(b) <= 0
}

// Gopo_bigint_GT: func (a bigint) > (b bigint) bool
func Gopo_bigint_GT(a, b Gope_bigint) bool {
	return a.Cmp(b) > 0
}

// Gopo_bigint_GE: func (a bigint) >= (b bigint) bool
func Gopo_bigint_GE(a, b Gope_bigint) bool {
	return a.Cmp(b) >= 0
}

// Gopo_bigint_EQ: func (a bigint) == (b bigint) bool
func Gopo_bigint_EQ(a, b Gope_bigint) bool {
	return a.Cmp(b) == 0
}

// Gopo_bigint_NE: func (a bigint) != (b bigint) bool
func Gopo_bigint_NE(a, b Gope_bigint) bool {
	return a.Cmp(b) != 0
}

// Gopo_bigint_Neg: func -(a bigint) bigint
func Gopo_bigint_Neg(a Gope_bigint) Gope_bigint {
	return tmpint1(a).Neg(a)
}

// Gopo_bigint_Not: func ^(a bigint) bigint
func Gopo_bigint_Not(a Gope_bigint) Gope_bigint {
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

// Gopo_bigrat_Add: func (a bigrat) + (b bigrat) bigrat
func Gopo_bigrat_Add(a, b Gope_bigrat) Gope_bigrat {
	return tmprat(a, b).Add(a, b)
}

// Gopo_bigrat_Sub: func (a bigrat) - (b bigrat) bigrat
func Gopo_bigrat_Sub(a, b Gope_bigrat) Gope_bigrat {
	return tmprat(a, b).Sub(a, b)
}

// Gopo_bigrat_Mul: func (a bigrat) * (b bigrat) bigrat
func Gopo_bigrat_Mul(a, b Gope_bigrat) Gope_bigrat {
	return tmprat(a, b).Mul(a, b)
}

// Gopo_bigrat_Quo: func (a bigrat) / (b bigrat) bigrat
func Gopo_bigrat_Quo(a, b Gope_bigrat) Gope_bigrat {
	return tmprat(a, b).Quo(a, b)
}

// Gopo_bigrat_LT: func (a bigrat) < (b bigrat) bool
func Gopo_bigrat_LT(a, b Gope_bigrat) bool {
	return a.Cmp(b) < 0
}

// Gopo_bigrat_LE: func (a bigrat) <= (b bigrat) bool
func Gopo_bigrat_LE(a, b Gope_bigrat) bool {
	return a.Cmp(b) <= 0
}

// Gopo_bigrat_GT: func (a bigrat) > (b bigrat) bool
func Gopo_bigrat_GT(a, b Gope_bigrat) bool {
	return a.Cmp(b) > 0
}

// Gopo_bigrat_GE: func (a bigrat) >= (b bigrat) bool
func Gopo_bigrat_GE(a, b Gope_bigrat) bool {
	return a.Cmp(b) >= 0
}

// Gopo_bigrat_EQ: func (a bigrat) == (b bigrat) bool
func Gopo_bigrat_EQ(a, b Gope_bigrat) bool {
	return a.Cmp(b) == 0
}

// Gopo_bigrat_NE: func (a bigrat) != (b bigrat) bool
func Gopo_bigrat_NE(a, b Gope_bigrat) bool {
	return a.Cmp(b) != 0
}

// Gopo_bigrat_Neg: func -(a bigrat) bigrat
func Gopo_bigrat_Neg(a Gope_bigrat) Gope_bigrat {
	return tmprat1(a).Neg(a)
}

// Gopo_bigrat_Inv: func /(a bigrat) bigrat
func Gopo_bigrat_Inv(a Gope_bigrat) Gope_bigrat {
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
