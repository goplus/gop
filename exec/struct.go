package exec

import (
	"reflect"
)

// -----------------------------------------------------------------------------

// A StructField describes a single field in a struct.
type StructField = reflect.StructField

// A StructInfo represents a struct information.
type StructInfo struct {
	Fields []StructField
	typ    reflect.Type
}

// Struct creates a new struct.
func Struct(fields []StructField) *StructInfo {
	return &StructInfo{
		Fields: fields,
		typ:    reflect.StructOf(fields),
	}
}

// Type returns the struct type.
func (p *StructInfo) Type() reflect.Type {
	return p.typ
}

// -----------------------------------------------------------------------------
