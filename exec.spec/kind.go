/*
 Copyright 2020 Qiniu Cloud (qiniu.com)

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

package exec

import (
	"math/big"
	"reflect"
	"unsafe"

	"github.com/qiniu/goplus/ast/spec"
)

// -----------------------------------------------------------------------------

// A Kind represents the specific kind of type that a Type represents.
type Kind = reflect.Kind

const (
	// Bool type
	Bool = reflect.Bool
	// Int type
	Int = reflect.Int
	// Int8 type
	Int8 = reflect.Int8
	// Int16 type
	Int16 = reflect.Int16
	// Int32 type
	Int32 = reflect.Int32
	// Int64 type
	Int64 = reflect.Int64
	// Uint type
	Uint = reflect.Uint
	// Uint8 type
	Uint8 = reflect.Uint8
	// Uint16 type
	Uint16 = reflect.Uint16
	// Uint32 type
	Uint32 = reflect.Uint32
	// Uint64 type
	Uint64 = reflect.Uint64
	// Uintptr type
	Uintptr = reflect.Uintptr
	// Float32 type
	Float32 = reflect.Float32
	// Float64 type
	Float64 = reflect.Float64
	// Complex64 type
	Complex64 = reflect.Complex64
	// Complex128 type
	Complex128 = reflect.Complex128
	// String type
	String = reflect.String
	// UnsafePointer type
	UnsafePointer = reflect.UnsafePointer
	// BigInt type
	BigInt = spec.BigInt
	// BigRat type
	BigRat = spec.BigRat
	// BigFloat type
	BigFloat = spec.BigFloat
)

var (
	// TyBool type
	TyBool = reflect.TypeOf(true)
	// TyInt type
	TyInt = reflect.TypeOf(int(0))
	// TyInt8 type
	TyInt8 = reflect.TypeOf(int8(0))
	// TyInt16 type
	TyInt16 = reflect.TypeOf(int16(0))
	// TyInt32 type
	TyInt32 = reflect.TypeOf(int32(0))
	// TyInt64 type
	TyInt64 = reflect.TypeOf(int64(0))
	// TyUint type
	TyUint = reflect.TypeOf(uint(0))
	// TyUint8 type
	TyUint8 = reflect.TypeOf(uint8(0))
	// TyUint16 type
	TyUint16 = reflect.TypeOf(uint16(0))
	// TyUint32 type
	TyUint32 = reflect.TypeOf(uint32(0))
	// TyUint64 type
	TyUint64 = reflect.TypeOf(uint64(0))
	// TyUintptr type
	TyUintptr = reflect.TypeOf(uintptr(0))
	// TyFloat32 type
	TyFloat32 = reflect.TypeOf(float32(0))
	// TyFloat64 type
	TyFloat64 = reflect.TypeOf(float64(0))
	// TyComplex64 type
	TyComplex64 = reflect.TypeOf(complex64(0))
	// TyComplex128 type
	TyComplex128 = reflect.TypeOf(complex128(0))
	// TyString type
	TyString = reflect.TypeOf("")
	// TyUnsafePointer type
	TyUnsafePointer = reflect.TypeOf(unsafe.Pointer(nil))
	// TyEmptyInterface type
	TyEmptyInterface = reflect.TypeOf((*interface{})(nil)).Elem()
	// TyError type
	TyError = reflect.TypeOf((*error)(nil)).Elem()
	// TyBigInt type
	TyBigInt = reflect.TypeOf((*big.Int)(nil))
	// TyBigRat type
	TyBigRat = reflect.TypeOf((*big.Rat)(nil))
	// TyBigFloat type
	TyBigFloat = reflect.TypeOf((*big.Float)(nil))
)

var (
	// TyByte type
	TyByte = reflect.TypeOf(byte(0))
	// TyRune type
	TyRune = reflect.TypeOf(rune(0))

	// TyEmptyInterfaceSlice type
	TyEmptyInterfaceSlice = reflect.SliceOf(TyEmptyInterface)
)

type bTI struct { // builtin type info
	typ  reflect.Type
	size uintptr
	//castFrom uint64
}

var builtinTypes = [...]bTI{
	Bool:          {TyBool, 1},                                  // 0},
	Int:           {TyInt, unsafe.Sizeof(int(0))},               // bitsAllReal},
	Int8:          {TyInt8, 1},                                  // bitsAllReal | bitBigInt},
	Int16:         {TyInt16, 2},                                 // bitsAllReal | bitBigInt},
	Int32:         {TyInt32, 4},                                 // bitsAllReal | bitBigInt},
	Int64:         {TyInt64, 8},                                 // bitsAllReal | bitBigInt},
	Uint:          {TyUint, unsafe.Sizeof(uint(0))},             // bitsAllReal | bitBigInt},
	Uint8:         {TyUint8, 1},                                 // bitsAllReal | bitBigInt},
	Uint16:        {TyUint16, 2},                                // bitsAllReal | bitBigInt},
	Uint32:        {TyUint32, 4},                                // bitsAllReal | bitBigInt},
	Uint64:        {TyUint64, 8},                                // bitsAllReal | bitBigInt},
	Uintptr:       {TyUintptr, unsafe.Sizeof(uintptr(0))},       // bitsAllReal | bitBigInt},
	Float32:       {TyFloat32, 4},                               // bitsAllReal | bitBigInt},
	Float64:       {TyFloat64, 8},                               // bitsAllReal | bitBigInt},
	Complex64:     {TyComplex64, 8},                             // bitsAllComplex},
	Complex128:    {TyComplex128, 16},                           // bitsAllComplex},
	String:        {TyString, unsafe.Sizeof(string('0'))},       // bitsAllIntUint},
	UnsafePointer: {TyUnsafePointer, unsafe.Sizeof(uintptr(0))}, // 0},
	BigInt:        {TyBigInt, unsafe.Sizeof(uintptr(0))},        // bitsAllReal | bitBigInt},
	BigRat:        {TyBigRat, unsafe.Sizeof(uintptr(0))},        // bitsAllReal | bitBigInt},
	BigFloat:      {TyBigFloat, unsafe.Sizeof(uintptr(0))},      // bitsAllReal | bitBigInt},
}

// TypeFromKind returns the type who has this kind.
func TypeFromKind(kind Kind) reflect.Type {
	return builtinTypes[kind].typ
}

// SizeofKind returns sizeof type who has this kind.
func SizeofKind(kind Kind) uintptr {
	return builtinTypes[kind].size
}

const (
	// BuiltinTypesLen - len(builtinTypes)
	BuiltinTypesLen = len(builtinTypes)
)

// -----------------------------------------------------------------------------

// A ConstKind represents the specific kind of type that a Type represents.
// The zero Kind is not a valid kind.
type ConstKind = spec.ConstKind

const (
	// ConstBoundRune - bound type: rune
	ConstBoundRune = spec.ConstBoundRune
	// ConstBoundString - bound type: string
	ConstBoundString = spec.ConstBoundString
	// ConstUnboundInt - unbound int type
	ConstUnboundInt = spec.ConstUnboundInt
	// ConstUnboundFloat - unbound float type
	ConstUnboundFloat = spec.ConstUnboundFloat
	// ConstUnboundComplex - unbound complex type
	ConstUnboundComplex = spec.ConstUnboundComplex
	// ConstUnboundPtr - nil: unbound ptr
	ConstUnboundPtr = spec.ConstUnboundPtr
)

// GoConstInfo represents a Go constant information.
type GoConstInfo struct {
	Pkg   GoPackage
	Name  string
	Kind  ConstKind
	Value interface{}
}

// -----------------------------------------------------------------------------
