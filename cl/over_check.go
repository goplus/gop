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
	"reflect"

	"github.com/goplus/gop/constant"
	"github.com/goplus/gop/token"
)

// 0 = 32bit, 1 = 64bit
const ptrSizeIndex = (^uintptr(0) >> 63)

func doesOverflow(v constant.Value, kind reflect.Kind) bool {
	if kind < reflect.Int || kind > reflect.Complex128 {
		return false
	}
	return constant.Compare(v, token.LSS, mincval[kind]) || constant.Compare(v, token.GTR, maxcval[kind])
}

var (
	maxcval = make(map[reflect.Kind]constant.Value)
	mincval = make(map[reflect.Kind]constant.Value)
)

func init() {
	maxcval[reflect.Int8] = constant.MakeInt64(0x7f)
	mincval[reflect.Int8] = constant.MakeInt64(-0x80)
	maxcval[reflect.Int16] = constant.MakeInt64(0x7fff)
	mincval[reflect.Int16] = constant.MakeInt64(-0x8000)
	maxcval[reflect.Int32] = constant.MakeInt64(0x7fffffff)
	mincval[reflect.Int32] = constant.MakeInt64(-0x80000000)
	maxcval[reflect.Int64] = constant.MakeInt64(0x7fffffffffffffff)
	mincval[reflect.Int64] = constant.MakeInt64(-0x8000000000000000)

	maxcval[reflect.Int] = []constant.Value{maxcval[reflect.Int32], maxcval[reflect.Int64]}[ptrSizeIndex]
	mincval[reflect.Int] = []constant.Value{mincval[reflect.Int32], mincval[reflect.Int64]}[ptrSizeIndex]

	maxcval[reflect.Uint8] = constant.MakeUint64(0xff)
	mincval[reflect.Uint8] = constant.MakeUint64(0)
	maxcval[reflect.Uint16] = constant.MakeUint64(0xffff)
	mincval[reflect.Uint16] = constant.MakeUint64(0)
	maxcval[reflect.Uint32] = constant.MakeUint64(0xffffffff)
	mincval[reflect.Uint32] = constant.MakeUint64(0)
	maxcval[reflect.Uint64] = constant.MakeUint64(0xffffffffffffffff)
	mincval[reflect.Uint64] = constant.MakeUint64(0)

	maxcval[reflect.Uint] = []constant.Value{maxcval[reflect.Uint32], maxcval[reflect.Uint64]}[ptrSizeIndex]
	mincval[reflect.Uint] = []constant.Value{mincval[reflect.Uint32], mincval[reflect.Uint64]}[ptrSizeIndex]
	maxcval[reflect.Uintptr] = maxcval[reflect.Uint]
	mincval[reflect.Uintptr] = mincval[reflect.Uint]

	maxcval[reflect.Float32] = constant.MakeFromLiteral("33554431p103", token.FLOAT, 0)
	mincval[reflect.Float32] = constant.MakeFromLiteral("-33554431p103", token.FLOAT, 0)
	maxcval[reflect.Float64] = constant.MakeFromLiteral("18014398509481983p970", token.FLOAT, 0)
	mincval[reflect.Float64] = constant.MakeFromLiteral("-18014398509481983p970", token.FLOAT, 0)

	maxcval[reflect.Complex64] = maxcval[reflect.Float32]
	mincval[reflect.Complex64] = mincval[reflect.Float32]
	maxcval[reflect.Complex128] = maxcval[reflect.Float64]
	mincval[reflect.Complex128] = mincval[reflect.Float64]
}
