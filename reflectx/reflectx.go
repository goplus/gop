package reflectx

import (
	"reflect"
)

type UserType struct {
	reflect.Type
	elem *UserType
}

func (t *UserType) Elem() reflect.Type {
	if t.elem != nil {
		return t.elem
	}
	return t.Type.Elem()
}

func NewUserType(t reflect.Type) reflect.Type {
	return &UserType{Type: t}
}

func PtrTo(t reflect.Type) reflect.Type {
	if ut, ok := t.(*UserType); ok {
		return &UserType{reflect.PtrTo(ut.Type), ut}
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
