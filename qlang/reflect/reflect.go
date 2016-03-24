package reflect

import (
	"reflect"
)

// -----------------------------------------------------------------------------

var Exports = map[string]interface{}{
	"valueOf":   reflect.ValueOf,
	"typeOf":    reflect.TypeOf,
	"indirect":  reflect.Indirect,
	"makeSlice": reflect.MakeSlice,
	"makeMap":   reflect.MakeMap,
	"zero":      reflect.Zero,

	"Map":       int(reflect.Map),
	"Slice":     int(reflect.Slice),
	"Interface": int(reflect.Interface),
	"Int":       int(reflect.Int),
}

// -----------------------------------------------------------------------------

