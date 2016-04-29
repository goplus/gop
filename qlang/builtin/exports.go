package builtin

import (
	"fmt"
	"reflect"

	"qlang.io/qlang.spec.v1"
)

// -----------------------------------------------------------------------------

var exports = map[string]interface{}{
	"append":    Append,
	"copy":      Copy,
	"delete":    Delete,
	"get":       Get,
	"len":       Len,
	"cap":       Cap,
	"mkmap":     Mkmap,
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

	"float":  Float,
	"int8":   Int8,
	"int16":  Int16,
	"int32":  Int32,
	"int64":  Int64,
	"int":    Int,
	"uint":   Uint,
	"byte":   Byte,
	"uint8":  Byte,
	"uint16": Uint16,
	"uint32": Uint32,
	"uint64": Uint64,
	"max":    Max,
	"min":    Min,

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
	qlang.GetVar = GetVar
	qlang.Get = Get
	qlang.SetIndex = SetIndex
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
