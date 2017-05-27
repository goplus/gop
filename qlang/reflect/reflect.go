package reflect

import (
	"reflect"
)

// -----------------------------------------------------------------------------

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"_name":     "reflect",
	"valueOf":   reflect.ValueOf,
	"typeOf":    reflect.TypeOf,
	"indirect":  reflect.Indirect,
	"makeSlice": reflect.MakeSlice,
	"makeMap":   reflect.MakeMap,
	"zero":      reflect.Zero,

	"ValueOf":   reflect.ValueOf,
	"TypeOf":    reflect.TypeOf,
	"Indirect":  reflect.Indirect,
	"MakeSlice": reflect.MakeSlice,
	"MakeMap":   reflect.MakeMap,
	"Zero":      reflect.Zero,

	"Map":       reflect.Map,
	"Slice":     reflect.Slice,
	"Interface": reflect.Interface,
	"Int":       reflect.Int,
}

// -----------------------------------------------------------------------------
