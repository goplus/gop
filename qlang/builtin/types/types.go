package types

import (
	"fmt"
	"reflect"
)

// -----------------------------------------------------------------------------

type Chan struct {
	Data reflect.Value
}

// -----------------------------------------------------------------------------

var (
	Var = reflect.TypeOf((*interface{})(nil)).Elem()
)

func Reflect(typ interface{}) reflect.Type {

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
	"string:int":    reflect.TypeOf(map[string]int(nil)),
	"string:float":  reflect.TypeOf(map[string]float64(nil)),
	"string:string": reflect.TypeOf(map[string]string(nil)),
	"string:var":    reflect.TypeOf(map[string]interface{}(nil)),
}

// -----------------------------------------------------------------------------

