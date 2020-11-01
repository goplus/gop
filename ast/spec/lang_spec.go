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

package spec

import (
	"reflect"
)

// -----------------------------------------------------------------------------

// A ConstKind represents the specific kind of type that a Type represents.
// The zero Kind is not a valid kind.
type ConstKind = reflect.Kind

const (
	// BigInt - bound type - bigint
	BigInt = ConstKind(reflect.UnsafePointer + 1)
	// BigRat - bound type - bigrat
	BigRat = ConstKind(reflect.UnsafePointer + 2)
	// BigFloat - bound type - bigfloat
	BigFloat = ConstKind(reflect.UnsafePointer + 3)
	// ConstBoundRune - bound type: rune
	ConstBoundRune = reflect.Int32
	// ConstBoundString - bound type: string
	ConstBoundString = reflect.String
	// ConstUnboundInt - unbound int type
	ConstUnboundInt = ConstKind(reflect.UnsafePointer + 4)
	// ConstUnboundFloat - unbound float type
	ConstUnboundFloat = ConstKind(reflect.UnsafePointer + 5)
	// ConstUnboundComplex - unbound complex type
	ConstUnboundComplex = ConstKind(reflect.UnsafePointer + 6)
	// ConstUnboundPtr - nil: unbound ptr
	ConstUnboundPtr = ConstKind(reflect.UnsafePointer + 7)
	// Slice - bound type: slice
	Slice = reflect.Slice
	// Map - bound type: map
	Map = reflect.Map
	// Chan - bound type: chan
	Chan = reflect.Chan
	// Ptr - bound type: ptr
	Ptr = reflect.Ptr
)

// IsConstBound checks a const is bound or not.
func IsConstBound(kind ConstKind) bool {
	return kind <= BigFloat
}

func KindName(kind ConstKind) string {
	switch kind {
	case BigInt:
		return "bigint"
	case BigRat:
		return "bigrat"
	case BigFloat:
		return "bigfloat"
	case ConstUnboundInt:
		return "unbound int"
	case ConstUnboundFloat:
		return "unbound float"
	case ConstUnboundComplex:
		return "unbound complex"
	case ConstUnboundPtr:
		return "unbound ptr"
	}
	return kind.String()
}

// -----------------------------------------------------------------------------
