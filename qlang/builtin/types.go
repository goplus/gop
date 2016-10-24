package builtin

import (
	"fmt"
	"reflect"

	qlang "qlang.io/spec"
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
	gotyString    = reflect.TypeOf("")
	gotyBool      = reflect.TypeOf(false)
	gotyInterface = reflect.TypeOf((*interface{})(nil)).Elem()
)

// TyByte represents the `byte` type.
//
var TyByte = TyUint8

// TyFloat represents the `float` type.
//
var TyFloat = TyFloat64

// -----------------------------------------------------------------------------

type tyVar int

func (p tyVar) GoType() reflect.Type {

	return gotyInterface
}

// NewInstance creates a new instance of a qlang type. required by `qlang type` spec.
//
func (p tyVar) NewInstance(args ...interface{}) interface{} {

	ret := new(interface{})
	if len(args) > 0 {
		*ret = args[0]
	}
	return ret
}

func (p tyVar) Call(a interface{}) interface{} {

	return a
}

func (p tyVar) String() string {

	return "var"
}

// TyVar represents the `var` type.
//
var TyVar = tyVar(0)

// -----------------------------------------------------------------------------

type goSliceFrom int
type goTypeOf int

func (p goSliceFrom) Call(a ...interface{}) interface{} {
	return SliceFrom(a...)
}

func (p goTypeOf) Call(a interface{}) reflect.Type {
	return reflect.TypeOf(a)
}

var sliceFrom = goSliceFrom(0)
var typeOf = goTypeOf(0)

func init() {
	t1 := reflect.TypeOf(TyVar)
	t2 := reflect.TypeOf(typeOf)
	t3 := reflect.TypeOf(sliceFrom)
	qlang.SetDontTyNormalize(t1)
	qlang.SetDontTyNormalize(t2)
	qlang.SetDontTyNormalize(t3)
}

// -----------------------------------------------------------------------------

// Elem returns *a
//
func Elem(a interface{}) interface{} {

	if t, ok := a.(qlang.GoTyper); ok {
		return qlang.TyPtrTo(t.GoType())
	}
	return reflect.ValueOf(a).Elem().Interface()
}

// Slice returns []T
//
func Slice(elem interface{}) interface{} {

	if t, ok := elem.(qlang.GoTyper); ok {
		return qlang.TySliceOf(t.GoType())
	}
	panic(fmt.Sprintf("invalid []T: `%v` isn't a qlang type", elem))
}

// Map returns map[key]elem
//
func Map(key, elem interface{}) interface{} {

	tkey, ok := key.(qlang.GoTyper)
	if !ok {
		panic(fmt.Sprintf("invalid map[key]elem: key `%v` isn't a qlang type", key))
	}
	telem, ok := elem.(qlang.GoTyper)
	if !ok {
		panic(fmt.Sprintf("invalid map[key]elem: elem `%v` isn't a qlang type", elem))
	}
	return qlang.TyMapOf(tkey.GoType(), telem.GoType())
}

// -----------------------------------------------------------------------------

// Make creates a instance of qlang builtin type (slice, map and chan)
//
func Make(typ qlang.GoTyper, args ...int) interface{} {

	t := typ.GoType()
	switch t.Kind() {
	case reflect.Slice:
		n, cap := 0, 0
		if len(args) == 1 {
			n = args[0]
			cap = n
		} else if len(args) > 1 {
			n, cap = args[0], args[1]
		}
		return reflect.MakeSlice(t, n, cap).Interface()
	case reflect.Map:
		return reflect.MakeMap(t).Interface()
	case reflect.Chan:
		return qlang.MakeChan(t, args...)
	}
	panic(fmt.Sprintf("cannot make type `%v`", typ))
}

// -----------------------------------------------------------------------------
