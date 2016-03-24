package builtin

import (
	"fmt"
	"reflect"

	"qlang.io/qlang.spec.v1"
)

// -----------------------------------------------------------------------------

var exports = map[string]interface{}{
	"append":    Append,
	"delete":    Delete,
	"get":       Get,
	"len":       Len,
	"map":       Map,
	"mapFrom":   MapFrom,
	"mapOf":     MapOf,
	"panic":     Panic,
	"panicf":    Panicf,
	"printf":    fmt.Printf,
	"println":   fmt.Println,
	"fprintln":  fmt.Fprintln,
	"set":       Set,
	"slice":     Slice,
	"sliceFrom": SliceFrom,
	"sliceOf":   SliceOf,
	"string":    String,
	"sub":       SubSlice,
	"type":      reflect.TypeOf,

	"float": Float,
	"int":   Int,
	"byte":  Byte,
	"max":   Max,
	"min":   Min,

	"undefined": qlang.Undefined,
	"nil":       nil,
	"true":      true,
	"false":     false,

	"$neg": Neg,
	"$mul": Mul,
	"$quo": Quo,
	"$mod": Mod,
	"$add": Add,
	"$sub": Sub,

	"$lt":  LT,
	"$gt":  GT,
	"$le":  LE,
	"$ge":  GE,
	"$eq":  EQ,
	"$ne":  NE,
	"$not": Not,
}

func init() {
	qlang.SubSlice = SubSlice
	qlang.SliceFrom = SliceFrom
	qlang.MapFrom = MapFrom
	qlang.EQ = EQ
	qlang.Get = Get
	qlang.Add = Add
	qlang.Sub = Sub
	qlang.Mul = Mul
	qlang.Quo = Quo
	qlang.Mod = Mod
	qlang.Inc = Inc
	qlang.Dec = Dec
	qlang.Import("", exports)
}

// -----------------------------------------------------------------------------
