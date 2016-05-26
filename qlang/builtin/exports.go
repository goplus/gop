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
	"mkslice":   Mkslice,
	"slice":     Mkslice,
	"sliceFrom": SliceFrom,
	"sliceOf":   SliceOf,
	"sub":       SubSlice,
	"type":      reflect.TypeOf,

	"float":   TyFloat64,
	"float64": TyFloat64,
	"float32": TyFloat32,
	"int8":    TyInt8,
	"int16":   TyInt16,
	"int32":   TyInt32,
	"int64":   TyInt64,
	"int":     TyInt,
	"uint":    TyUint,
	"byte":    TyUint8,
	"uint8":   TyUint8,
	"uint16":  TyUint16,
	"uint32":  TyUint32,
	"uint64":  TyUint64,
	"string":  TyString,

	"max": Max,
	"min": Min,

	"undefined": qlang.Undefined,
	"nil":       nil,
	"true":      true,
	"false":     false,

	"$elem": Elem,
	"$neg":  Neg,
	"$mul":  Mul,
	"$quo":  Quo,
	"$mod":  Mod,
	"$add":  Add,
	"$sub":  Sub,

	"$xor":    Xor,
	"$lshr":   Lshr,
	"$rshr":   Rshr,
	"$bitand": BitAnd,
	"$bitor":  BitOr,
	"$bitnot": BitNot,
	"$andnot": AndNot,

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
	qlang.Slice = Slice
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
	qlang.Xor = Xor
	qlang.Lshr = Lshr
	qlang.Rshr = Rshr
	qlang.BitAnd = BitAnd
	qlang.BitOr = BitOr
	qlang.AndNot = AndNot
	qlang.Inc = Inc
	qlang.Dec = Dec
	qlang.Import("", exports)
}

// -----------------------------------------------------------------------------
