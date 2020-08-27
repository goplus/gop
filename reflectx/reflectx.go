package reflectx

import (
	"reflect"
)

type UserType struct {
	reflect.Type
	elem  reflect.Type
	key   reflect.Type
	field interface{}
}

func (t *UserType) Elem() reflect.Type {
	if t.elem != nil {
		return t.elem
	}
	return t.Type.Elem()
}

func (t *UserType) Key() reflect.Type {
	if t.key != nil {
		return t.key
	}
	return t.Type.Key()
}

func (t *UserType) FieldByName(name string) (sf reflect.StructField, ok bool) {
	if st, ok := t.field.([]reflect.StructField); ok {
		for _, t := range st {
			if t.Name == name {
				return t, true
			}
		}
	}
	return t.Type.FieldByName(name)
}

func (t *UserType) Field(i int) reflect.StructField {
	if st, ok := t.field.([]reflect.StructField); ok {
		return st[i]
	}
	return t.Type.Field(i)
}

func NewUserType(t reflect.Type, value interface{}) reflect.Type {
	return &UserType{Type: t, field: value}
}

func IsUserType(t reflect.Type) bool {
	_, ok := t.(*UserType)
	return ok
}

func ToType(t reflect.Type) reflect.Type {
	if ut, ok := t.(*UserType); ok {
		return ut.Type
	}
	return t
}

func ToTypes(typs []reflect.Type) []reflect.Type {
	ret := make([]reflect.Type, len(typs))
	for i := 0; i < len(typs); i++ {
		if ut, ok := typs[i].(*UserType); ok {
			ret[i] = ut.Type
		} else {
			ret[i] = typs[i]
		}
	}
	return ret
}

func StructOf(fields []reflect.StructField) reflect.Type {
	t := reflect.StructOf(fields)
	return &UserType{Type: t, field: fields}
}

func PtrTo(t reflect.Type) reflect.Type {
	if ut, ok := t.(*UserType); ok {
		return &UserType{Type: reflect.PtrTo(ut.Type), elem: ut}
	}
	return reflect.PtrTo(t)
}

func FuncOf(in, out []reflect.Type, variadic bool) reflect.Type {
	return reflect.FuncOf(ToTypes(in), ToTypes(out), variadic)
}

func SliceOf(t reflect.Type) reflect.Type {
	if ut, ok := t.(*UserType); ok {
		return &UserType{Type: reflect.SliceOf(ut.Type), elem: ut}
	}
	return reflect.SliceOf(t)
}

func ArrayOf(count int, elem reflect.Type) reflect.Type {
	if ut, ok := elem.(*UserType); ok {
		return &UserType{Type: reflect.ArrayOf(count, ut.Type), elem: ut}
	}
	return reflect.ArrayOf(count, elem)
}

func ChanOf(dir reflect.ChanDir, typ reflect.Type) reflect.Type {
	if ut, ok := typ.(*UserType); ok {
		return &UserType{Type: reflect.ChanOf(dir, ut.Type), elem: ut}
	}
	return reflect.ChanOf(dir, typ)
}

func MapOf(key reflect.Type, value reflect.Type) reflect.Type {
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

func MakeSlice(typ reflect.Type, len, cap int) reflect.Value {
	return reflect.MakeSlice(ToType(typ), len, cap)
}

func MakeMap(typ reflect.Type) reflect.Value {
	return reflect.MakeMap(ToType(typ))
}

func MakeMapWithSize(typ reflect.Type, n int) reflect.Value {
	return reflect.MakeMapWithSize(ToType(typ), n)
}

func MakeChan(typ reflect.Type, buffer int) reflect.Value {
	return reflect.MakeChan(ToType(typ), buffer)
}

func New(t reflect.Type) reflect.Value {
	return reflect.New(ToType(t))
}
