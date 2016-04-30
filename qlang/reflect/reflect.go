package reflect

import (
	"reflect"
)

// -----------------------------------------------------------------------------

// Exports is export table of this module.
//
var Exports = map[string]interface{}{
	"valueOf":   reflect.ValueOf,
	"typeOf":    reflect.TypeOf,
	"indirect":  reflect.Indirect,
	"makeSlice": reflect.MakeSlice,
	"makeMap":   reflect.MakeMap,
	"zero":      reflect.Zero,

	"Map":       reflect.Map,
	"Slice":     reflect.Slice,
	"Interface": reflect.Interface,
	"Int":       reflect.Int,
}

// -----------------------------------------------------------------------------
