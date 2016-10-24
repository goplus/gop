package types

import (
	"fmt"
	"reflect"

	qlang "qlang.io/spec"
)

// -----------------------------------------------------------------------------

var (
	// Var is reflect.Type of interface{}
	Var = reflect.TypeOf((*interface{})(nil)).Elem()
)

// Reflect returns reflect.Type of typ.
//
func Reflect(typ interface{}) reflect.Type {

	if t, ok := typ.(qlang.GoTyper); ok {
		return t.GoType()
	}
	if t, ok := typ.(string); ok {
		if v, ok := builtinTypes[t]; ok {
			return v
		}
	} else if t, ok := typ.(reflect.Type); ok {
		return t
	}
	panic(fmt.Errorf("unknown type: `%v`", typ))
}

var builtinTypes = map[string]reflect.Type{
	"int":           reflect.TypeOf(0),
	"bool":          reflect.TypeOf(false),
	"float":         reflect.TypeOf(float64(0)),
	"string":        reflect.TypeOf(""),
	"byte":          reflect.TypeOf(byte(0)),
	"var":           Var,
	"int:int":       reflect.TypeOf(map[int]int(nil)),
	"int:float":     reflect.TypeOf(map[int]float64(nil)),
	"int:string":    reflect.TypeOf(map[int]string(nil)),
	"int:var":       reflect.TypeOf(map[int]interface{}(nil)),
	"string:int":    reflect.TypeOf(map[string]int(nil)),
	"string:float":  reflect.TypeOf(map[string]float64(nil)),
	"string:string": reflect.TypeOf(map[string]string(nil)),
	"string:var":    reflect.TypeOf(map[string]interface{}(nil)),
}

// -----------------------------------------------------------------------------
