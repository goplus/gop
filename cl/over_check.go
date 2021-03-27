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

package cl

import (
	"math/big"
	"reflect"
)

const ptrSize = 4 << (^uintptr(0) >> 63)

func doesOverflowInt(v *big.Int, kind reflect.Kind) bool {
	return v.Cmp(minintval[kind]) < 0 || v.Cmp(maxintval[kind]) > 0
}

func doesOverflowFloat(v *big.Float, kind reflect.Kind) bool {
	return v.Cmp(minfltval[kind]) < 0 || v.Cmp(maxfltval[kind]) > 0
}

var (
	maxintval = make(map[reflect.Kind]*big.Int)
	minintval = make(map[reflect.Kind]*big.Int)
	maxfltval = make(map[reflect.Kind]*big.Float)
	minfltval = make(map[reflect.Kind]*big.Float)
)

func init() {
	maxintval[reflect.Int8] = new(big.Int).SetInt64(0x7f)
	minintval[reflect.Int8] = new(big.Int).SetInt64(-0x80)
	maxintval[reflect.Int16] = new(big.Int).SetInt64(0x7fff)
	minintval[reflect.Int16] = new(big.Int).SetInt64(-0x8000)
	maxintval[reflect.Int32] = new(big.Int).SetInt64(0x7fffffff)
	minintval[reflect.Int32] = new(big.Int).SetInt64(-0x80000000)
	maxintval[reflect.Int64], _ = new(big.Int).SetString("0x7fffffffffffffff", 0)
	minintval[reflect.Int64], _ = new(big.Int).SetString("-0x8000000000000000", 0)

	zero := big.NewInt(0)
	maxintval[reflect.Uint8] = new(big.Int).SetUint64(0xff)
	minintval[reflect.Uint8] = zero
	maxintval[reflect.Uint16] = new(big.Int).SetUint64(0xffff)
	minintval[reflect.Uint16] = zero
	maxintval[reflect.Uint32] = new(big.Int).SetUint64(0xffffffff)
	minintval[reflect.Uint32] = zero
	maxintval[reflect.Uint64], _ = new(big.Int).SetString("0xffffffffffffffff", 0)
	minintval[reflect.Uint64] = zero

	minintval[reflect.Uint] = zero
	if ptrSize == 4 {
		maxintval[reflect.Int] = maxintval[reflect.Int32]
		minintval[reflect.Int] = minintval[reflect.Int32]
		maxintval[reflect.Uint] = maxintval[reflect.Uint32]
	} else {
		maxintval[reflect.Int] = maxintval[reflect.Int64]
		minintval[reflect.Int] = minintval[reflect.Int64]
		maxintval[reflect.Uint] = maxintval[reflect.Uint64]
	}
	maxintval[reflect.Interface] = maxintval[reflect.Int]
	minintval[reflect.Interface] = minintval[reflect.Int]

	maxfltval[reflect.Float32], _ = new(big.Float).SetString("33554431p103")
	minfltval[reflect.Float32], _ = new(big.Float).SetString("-33554431p103")
	maxfltval[reflect.Float64], _ = new(big.Float).SetString("18014398509481983p970")
	minfltval[reflect.Float64], _ = new(big.Float).SetString("-18014398509481983p970")
	maxfltval[reflect.Interface] = maxfltval[reflect.Float64]
	minfltval[reflect.Interface] = minfltval[reflect.Float64]
}
