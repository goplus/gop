package spec

import (
	"reflect"
)

// -----------------------------------------------------------------------------

// A ConstKind represents the specific kind of type that a Type represents.
// The zero Kind is not a valid kind.
type ConstKind = reflect.Kind

const (
	// ConstBoundRune - bound type: rune
	ConstBoundRune = reflect.Int32
	// ConstBoundString - bound type: string
	ConstBoundString = reflect.String
	// ConstUnboundInt - unbound int type
	ConstUnboundInt = ConstKind(reflect.UnsafePointer + 3)
	// ConstUnboundFloat - unbound float type
	ConstUnboundFloat = ConstKind(reflect.UnsafePointer + 4)
	// ConstUnboundComplex - unbound complex type
	ConstUnboundComplex = ConstKind(reflect.UnsafePointer + 5)
	// ConstUnboundPtr - nil: unbound ptr
	ConstUnboundPtr = ConstKind(reflect.UnsafePointer + 6)
)

// IsConstBound checks a const is bound or not.
func IsConstBound(kind ConstKind) bool {
	return kind <= reflect.UnsafePointer
}

// -----------------------------------------------------------------------------
