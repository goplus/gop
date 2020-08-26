package reflectx

import (
	"reflect"
)

type UserType struct {
	reflect.Type
	elem  *UserType
	value interface{}
}

func (t *UserType) Elem() reflect.Type {
	if t.elem != nil {
		return t.elem
	}
	return t.Type.Elem()
}

func (t *UserType) FieldByName(name string) (sf reflect.StructField, ok bool) {
	if st, ok := t.value.([]reflect.StructField); ok {
		for _, t := range st {
			if t.Name == name {
				return t, true
			}
		}
	}
	return t.Type.FieldByName(name)
}

func (t *UserType) Field(i int) reflect.StructField {
	if st, ok := t.value.([]reflect.StructField); ok {
		return st[i]
	}
	return t.Type.Field(i)
}

func NewUserType(t reflect.Type, value interface{}) reflect.Type {
	return &UserType{Type: t, value: value}
}

func StructOf(fields []reflect.StructField) reflect.Type {
	t := reflect.StructOf(fields)
	return &UserType{Type: t, value: fields}
}

func PtrTo(t reflect.Type) reflect.Type {
	if ut, ok := t.(*UserType); ok {
		return &UserType{Type: reflect.PtrTo(ut.Type), elem: ut}
	}
	return reflect.PtrTo(t)
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

func FuncOf(in, out []reflect.Type, variadic bool) reflect.Type {
	return reflect.FuncOf(ToTypes(in), ToTypes(out), variadic)
}

func New(t reflect.Type) reflect.Value {
	return reflect.New(ToType(t))
}
