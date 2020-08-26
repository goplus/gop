package reflectx

import (
	"reflect"
)

type UserType struct {
	reflect.Type
	User interface{}
}

func PtrTo(t reflect.Type) reflect.Type {
	if ut, ok := t.(*UserType); ok {
		return &UserType{reflect.PtrTo(ut.Type), ut.User}
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

func ToTypes(in []reflect.Type) []reflect.Type {
	out := make([]reflect.Type, len(in))
	for i := 0; i < len(in); i++ {
		if ut, ok := in[i].(*UserType); ok {
			out[i] = ut.Type
		} else {
			out[i] = in[i]
		}
	}
	return out
}

func FuncOf(in, out []reflect.Type, variadic bool) reflect.Type {
	return reflect.FuncOf(ToTypes(in), ToTypes(out), variadic)
}

func New(t reflect.Type) reflect.Value {
	return reflect.New(ToType(t))
}
