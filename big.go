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

package gop

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

// Gopo_Bigint__Neg: func -(a bigint) bigint {
func Gopo_Bigint__Neg(a Bigint) Bigint {
	return tmpint1(a).Neg(a)
}

// Gopo_Bigint__Not: func ^(a bigint) bigint {
func Gopo_Bigint__Not(a Bigint) Bigint {
	return tmpint1(a).Not(a)
}

// -----------------------------------------------------------------------------
// type bigrat

type Bigrat = *big.Rat

// -----------------------------------------------------------------------------
// type bigfloat

type Bigfloat = *big.Float

// -----------------------------------------------------------------------------
