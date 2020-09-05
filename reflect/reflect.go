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

// Package reflect implements golang reflect type wrapper
package reflect

import (
	"reflect"
)

// A Kind represents the specific kind of type that a Type represents.
// The zero Kind is not a valid kind.
type Kind = reflect.Kind

const (
	Invalid reflect.Kind = iota
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	Array
	Chan
	Func
	Interface
	Map
	Ptr
	Slice
	String
	Struct
	UnsafePointer
)

// ChanDir represents a channel type's direction.
type ChanDir = reflect.ChanDir

const (
	RecvDir ChanDir             = 1 << iota // <-chan
	SendDir                                 // chan<-
	BothDir = RecvDir | SendDir             // chan
)

// Type is the representation of a Go type.
type Type = reflect.Type

// Value is the reflection interface to a Go value.
type Value = reflect.Value

// A StructField describes a single field in a struct.
type StructField = reflect.StructField

// A StructTag is the tag string in a struct field.
type StructTag = reflect.StructTag

// Method represents a single method.
type Method = reflect.Method

// UserType is user-defined struct and type
type UserType struct {
	Type
	elem   Type
	key    Type
	field  []StructField
	method []Method
	in     []Type
	out    []Type
}

func (t *UserType) NumIn() int {
	return len(t.in)
}

func (t *UserType) In(i int) Type {
	return t.in[i]
}

func (t *UserType) NumOut() int {
	return len(t.out)
}

func (t *UserType) Out(i int) Type {
	return t.out[i]
}

func (t *UserType) Elem() Type {
	if t.elem != nil {
		return t.elem
	}
	return t.Type.Elem()
}

func (t *UserType) Key() Type {
	if t.key != nil {
		return t.key
	}
	return t.Type.Key()
}

func (t *UserType) NumMethod() int {
	if t.method != nil {
		return len(t.method)
	}
	return t.Type.NumMethod()
}

func (t *UserType) Method(i int) reflect.Method {
	if t.method != nil {
		return t.method[i]
	}
	return t.Type.Method(i)
}

func (t *UserType) MethodByName(name string) (reflect.Method, bool) {
	for _, m := range t.method {
		if m.Name == name {
			return m, true
		}
	}
	return t.Type.MethodByName(name)
}

func (t *UserType) FieldByName(name string) (sf StructField, ok bool) {
	for _, t := range t.field {
		if t.Name == name {
			return t, true
		}
	}
	return t.Type.FieldByName(name)
}

func (t *UserType) Field(i int) StructField {
	if t.field != nil {
		return t.field[i]
	}
	return t.Type.Field(i)
}

func (t *UserType) FieldByIndex(index []int) (f StructField) {
	f.Type = t
	for i, x := range index {
		if i > 0 {
			ft := f.Type
			if ft.Kind() == Ptr && ft.Elem().Kind() == Struct {
				ft = ft.Elem()
			}
			f.Type = ft
		}
		f = f.Type.Field(x)
	}
	return
}

func (t *UserType) ConvertibleTo(u Type) bool {
	return t.Type.ConvertibleTo(toType(u))
}

func (t *UserType) Implements(u Type) bool {
	return Implements(t, u)
}

func NewUserType(t Type) Type {
	return &UserType{Type: t}
}

func IsUserType(t Type) bool {
	_, ok := t.(*UserType)
	return ok
}

func toType(t Type) Type {
	if ut, ok := t.(*UserType); ok {
		return ut.Type
	}
	return t
}

func toTypes(typs []Type) []Type {
	ret := make([]Type, len(typs))
	for i := 0; i < len(typs); i++ {
		if ut, ok := typs[i].(*UserType); ok {
			ret[i] = ut.Type
		} else {
			ret[i] = typs[i]
		}
	}
	return ret
}

func StructOf(fields []StructField) Type {
	t := reflect.StructOf(fields)
	return &UserType{Type: t, field: fields}
}

var (
	emptyInterface = reflect.TypeOf((*interface{})(nil)).Elem()
)

func InterfaceOf(methods []Method) Type {
	return &UserType{Type: emptyInterface, method: methods}
}

func PtrTo(t Type) Type {
	if ut, ok := t.(*UserType); ok {
		return &UserType{Type: reflect.PtrTo(ut.Type), elem: ut}
	}
	return reflect.PtrTo(t)
}

func FuncOf(in, out []Type, variadic bool) Type {
	return &UserType{Type: reflect.FuncOf(toTypes(in), toTypes(out), variadic), in: in, out: out}
}

func SliceOf(t Type) Type {
	if ut, ok := t.(*UserType); ok {
		return &UserType{Type: reflect.SliceOf(ut.Type), elem: ut}
	}
	return reflect.SliceOf(t)
}

func ArrayOf(count int, elem Type) Type {
	if ut, ok := elem.(*UserType); ok {
		return &UserType{Type: reflect.ArrayOf(count, ut.Type), elem: ut}
	}
	return reflect.ArrayOf(count, elem)
}

func ChanOf(dir reflect.ChanDir, typ Type) Type {
	if ut, ok := typ.(*UserType); ok {
		return &UserType{Type: reflect.ChanOf(dir, ut.Type), elem: ut}
	}
	return reflect.ChanOf(dir, typ)
}

func MapOf(key Type, value Type) Type {
	var user bool
	uKey := key
	uValue := value
	if ut, ok := key.(*UserType); ok {
		key = ut.Type
		uKey = ut
		user = true
	}
	if ut, ok := value.(*UserType); ok {
		value = ut.Type
		uValue = ut
		user = true
	}
	if user {
		return &UserType{Type: reflect.MapOf(key, value), key: uKey, elem: uValue}
	}
	return reflect.MapOf(key, value)
}

func MakeSlice(typ Type, len, cap int) Value {
	return reflect.MakeSlice(toType(typ), len, cap)
}

func MakeMap(typ Type) Value {
	return reflect.MakeMap(toType(typ))
}

func MakeMapWithSize(typ Type, n int) Value {
	return reflect.MakeMapWithSize(toType(typ), n)
}

func MakeChan(typ Type, buffer int) Value {
	return reflect.MakeChan(toType(typ), buffer)
}

func MakeFunc(typ Type, fn func(args []Value) (results []Value)) Value {
	return reflect.MakeFunc(toType(typ), fn)
}

func New(t Type) Value {
	return reflect.New(toType(t))
}

func Indirect(v Value) Value {
	return reflect.Indirect(v)
}

func Copy(dst, src Value) int {
	return reflect.Copy(dst, src)
}

func AppendSlice(s, t Value) Value {
	return reflect.AppendSlice(s, t)
}

func Append(s Value, x ...Value) Value {
	return reflect.Append(s, x...)
}

func EqualType(t1, t2 Type) bool {
	return toType(t1) == toType(t2)
}

func DeepEqual(x interface{}, y interface{}) bool {
	return reflect.DeepEqual(x, y)
}

func ConvertibleTo(from Type, to Type) bool {
	return toType(from).ConvertibleTo(toType(to))
}

func Convert(v Value, t Type) Value {
	return v.Convert(toType(t))
}

func equalTypeMethod(t Type, u Type) bool {
	if t.NumIn() != u.NumIn()+1 {
		return false
	}
	if t.NumOut() != u.NumOut() {
		return false
	}
	for i := 0; i < u.NumIn(); i++ {
		if t.In(i+1) != u.In(i) {
			return false
		}
	}
	for i := 0; i < u.NumOut(); i++ {
		if t.Out(i) != u.Out(i) {
			return false
		}
	}
	return true
}

func Implements(t Type, u Type) bool {
	if u == nil {
		panic("reflect: nil type passed to Type.Implements")
	}
	if u.Kind() != Interface {
		panic("reflect: non-interface type passed to Type.Implements")
	}
	if t.Kind() != Interface && !IsUserType(u) {
		return toType(t).Implements(u)
	}
	ucount := u.NumMethod()
	if ucount == 0 {
		return true
	}
	tcount := t.NumMethod()
	if t.Kind() == reflect.Interface {
		i := 0
		for j := 0; j < tcount; j++ {
			tm := t.Method(j)
			um := u.Method(i)
			if tm.Name == um.Name && EqualType(tm.Type, um.Type) {
				if i++; i >= ucount {
					return true
				}
			}
		}
	} else {
		i := 0
		for j := 0; j < tcount; j++ {
			tm := t.Method(j)
			um := u.Method(i)
			if tm.Name == um.Name && equalTypeMethod(tm.Type, um.Type) {
				if i++; i >= ucount {
					return true
				}
			}
		}
	}
	return false
}

func Zero(t Type) Value {
	return reflect.Zero(toType(t))
}

func TypeOf(i interface{}) Type {
	return reflect.TypeOf(i)
}

func ValueOf(i interface{}) Value {
	return reflect.ValueOf(i)
}
