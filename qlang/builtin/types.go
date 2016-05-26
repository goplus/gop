package builtin

import (
	"fmt"
	"reflect"
)

// -----------------------------------------------------------------------------

var (
	gotyInt       = reflect.TypeOf(int(0))
	gotyInt8      = reflect.TypeOf(int8(0))
	gotyInt16     = reflect.TypeOf(int16(0))
	gotyInt32     = reflect.TypeOf(int32(0))
	gotyInt64     = reflect.TypeOf(int64(0))
	gotyUint      = reflect.TypeOf(uint(0))
	gotyUint8     = reflect.TypeOf(uint8(0))
	gotyUint16    = reflect.TypeOf(uint16(0))
	gotyUint32    = reflect.TypeOf(uint32(0))
	gotyUint64    = reflect.TypeOf(uint64(0))
	gotyFloat32   = reflect.TypeOf(float32(0))
	gotyFloat64   = reflect.TypeOf(float64(0))
	gotyString    = reflect.TypeOf(string(""))
	gotyInterface = reflect.TypeOf((*interface{})(nil)).Elem()
)

// -----------------------------------------------------------------------------

// A Type represents a qlang builtin type.
//
type Type struct {
	t reflect.Type
}

// NewType returns a qlang builtin type object.
//
func NewType(t reflect.Type) *Type {

	return &Type{t: t}
}

// GoType returns the underlying go type. required by `qlang type` spec.
//
func (p *Type) GoType() reflect.Type {

	return p.t
}

// NewInstance creates a new instance of a qlang type. required by `qlang type` spec.
//
func (p *Type) NewInstance(args ...interface{}) interface{} {

	ret := reflect.New(p.t)
	if len(args) > 0 {
		v := reflect.ValueOf(args[0]).Convert(p.t)
		ret.Set(v)
	}
	return ret.Interface()
}

// Call returns `T(a)`.
//
func (p *Type) Call(a interface{}) interface{} {

	return reflect.ValueOf(a).Convert(p.t).Interface()
}

func (p *Type) String() string {

	return p.t.String()
}

// TySliceOf represents the `[]T` type.
//
func TySliceOf(elem reflect.Type) *Type {

	return &Type{t: reflect.SliceOf(elem)}
}

// TyMapOf represents the `map[key]elem` type.
//
func TyMapOf(key, elem reflect.Type) *Type {

	return &Type{t: reflect.MapOf(key, elem)}
}

// TyPtrTo represents the `*T` type.
//
func TyPtrTo(elem reflect.Type) *Type {

	return &Type{t: reflect.PtrTo(elem)}
}

// TyByte represents the `byte` type.
//
var TyByte = TyUint8

// TyFloat represents the `float` type.
//
var TyFloat = TyFloat64

// -----------------------------------------------------------------------------

type goTyper interface {
	GoType() reflect.Type
}

// Elem returns *a
//
func Elem(a interface{}) interface{} {

	if t, ok := a.(goTyper); ok {
		return TyPtrTo(t.GoType())
	}
	return reflect.ValueOf(a).Elem().Interface()
}

// Slice returns []T
//
func Slice(elem interface{}) interface{} {

	if t, ok := elem.(goTyper); ok {
		return TySliceOf(t.GoType())
	}
	panic(fmt.Sprintf("invalid []T: `%v` isn't a qlang type", elem))
}

// Map returns map[key]elem
//
func Map(key, elem interface{}) interface{} {

	tkey, ok := key.(goTyper)
	if !ok {
		panic(fmt.Sprintf("invalid map[key]elem: key `%v` isn't a qlang type", key))
	}
	telem, ok := elem.(goTyper)
	if !ok {
		panic(fmt.Sprintf("invalid map[key]elem: elem `%v` isn't a qlang type", elem))
	}
	return TyMapOf(tkey.GoType(), telem.GoType())
}

// -----------------------------------------------------------------------------
