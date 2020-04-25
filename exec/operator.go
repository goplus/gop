package exec

import (
	"fmt"
	"reflect"
	"strconv"
)

// -----------------------------------------------------------------------------

// Operator type.
type Operator uint

const (
	// OpAdd '+' String/Int/Uint/Float/Complex
	OpAdd Operator = iota
	// OpSub '-' Int/Uint/Float/Complex
	OpSub
	// OpMul '*' Int/Uint/Float/Complex
	OpMul
	// OpDiv '/' Int/Uint/Float/Complex
	OpDiv
	// OpMod '%' Int/Uint
	OpMod
	// OpBitAnd '&' Int/Uint
	OpBitAnd
	// OpBitOr '|' Int/Uint
	OpBitOr
	// OpBitXor '^' Int/Uint
	OpBitXor
	// OpBitAndNot '&^' Int/Uint
	OpBitAndNot
	// OpBitSHL '<<' Int/Uint, Uint
	OpBitSHL
	// OpBitSHR '>>' Int/Uint, Uint
	OpBitSHR
	// OpLT '<' String/Int/Uint/Float
	OpLT
	// OpLE '<=' String/Int/Uint/Float
	OpLE
	// OpGT '>' String/Int/Uint/Float
	OpGT
	// OpGE '>=' String/Int/Uint/Float
	OpGE
	// OpEQ '==' ComparableType
	// Slice, map, and function values are not comparable. However, as a special case, a slice, map,
	// or function value may be compared to the predeclared identifier nil.
	OpEQ
	// OpEQNil '==' nil
	OpEQNil
	// OpNE '!=' ComparableType
	OpNE
	// OpNENil '!=' nil
	OpNENil
	// OpLAnd '&&' Bool
	OpLAnd
	// OpLOr '||' Bool
	OpLOr
	// OpNeg '-'
	OpNeg
	// OpNot '!'
	OpNot
	// OpBitNot '^'
	OpBitNot
)

var opNames = [...]string{
	OpAdd:       "+",
	OpSub:       "-",
	OpMul:       "*",
	OpDiv:       "/",
	OpMod:       "%",
	OpBitAnd:    "&",
	OpBitOr:     "|",
	OpBitXor:    "^",
	OpBitAndNot: "&^",
	OpBitSHL:    "<<",
	OpBitSHR:    ">>",
	OpLT:        "<",
	OpLE:        "<=",
	OpGT:        ">",
	OpGE:        ">=",
	OpEQ:        "==",
	OpEQNil:     "== nil",
	OpNE:        "!=",
	OpNENil:     "!= nil",
	OpLAnd:      "&&",
	OpLOr:       "||",
	OpNeg:       "-",
	OpNot:       "!",
	OpBitNot:    "^",
}

func (op Operator) String() string {
	if int(op) < len(opNames) {
		return opNames[op]
	}
	return "op" + strconv.Itoa(int(op))
}

// -----------------------------------------------------------------------------

// A Kind represents the specific kind of type that a Type represents.
type Kind = reflect.Kind

const (
	// Bool type
	Bool = reflect.Bool
	// Int type
	Int = reflect.Int
	// Int8 type
	Int8 = reflect.Int8
	// Int16 type
	Int16 = reflect.Int16
	// Int32 type
	Int32 = reflect.Int32
	// Int64 type
	Int64 = reflect.Int64
	// Uint type
	Uint = reflect.Uint
	// Uint8 type
	Uint8 = reflect.Uint8
	// Uint16 type
	Uint16 = reflect.Uint16
	// Uint32 type
	Uint32 = reflect.Uint32
	// Uint64 type
	Uint64 = reflect.Uint64
	// Uintptr type
	Uintptr = reflect.Uintptr
	// Float32 type
	Float32 = reflect.Float32
	// Float64 type
	Float64 = reflect.Float64
	// Complex64 type
	Complex64 = reflect.Complex64
	// Complex128 type
	Complex128 = reflect.Complex128

	// String type
	String = reflect.String

	// Chan type
	Chan = reflect.Chan
	// Func type
	Func = reflect.Func
	// Interface type
	Interface = reflect.Interface
	// Map type
	Map = reflect.Map
	// Ptr type
	Ptr = reflect.Ptr
	// Slice type
	Slice = reflect.Slice
	// UnsafePointer type
	UnsafePointer = reflect.UnsafePointer
)

// -----------------------------------------------------------------------------

func execAddInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) + p.data[n-1].(int)
	p.data = p.data[:n-1]
}

func execAddInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) + p.data[n-1].(int8)
	p.data = p.data[:n-1]
}

func execAddInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) + p.data[n-1].(int16)
	p.data = p.data[:n-1]
}

func execAddInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) + p.data[n-1].(int32)
	p.data = p.data[:n-1]
}

func execAddInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) + p.data[n-1].(int64)
	p.data = p.data[:n-1]
}

func execAddUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) + p.data[n-1].(uint)
	p.data = p.data[:n-1]
}

func execAddUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) + p.data[n-1].(uint8)
	p.data = p.data[:n-1]
}

func execAddUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) + p.data[n-1].(uint16)
	p.data = p.data[:n-1]
}

func execAddUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) + p.data[n-1].(uint32)
	p.data = p.data[:n-1]
}

func execAddUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) + p.data[n-1].(uint64)
	p.data = p.data[:n-1]
}

func execAddUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) + p.data[n-1].(uintptr)
	p.data = p.data[:n-1]
}

func execAddFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float32) + p.data[n-1].(float32)
	p.data = p.data[:n-1]
}

func execAddFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float64) + p.data[n-1].(float64)
	p.data = p.data[:n-1]
}

func execAddComplex64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(complex64) + p.data[n-1].(complex64)
	p.data = p.data[:n-1]
}

func execAddComplex128(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(complex128) + p.data[n-1].(complex128)
	p.data = p.data[:n-1]
}

func execAddString(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(string) + p.data[n-1].(string)
	p.data = p.data[:n-1]
}

// -----------------------------------------------------------------------------

func execSubInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) - p.data[n-1].(int)
	p.data = p.data[:n-1]
}

func execSubInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) - p.data[n-1].(int8)
	p.data = p.data[:n-1]
}

func execSubInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) - p.data[n-1].(int16)
	p.data = p.data[:n-1]
}

func execSubInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) - p.data[n-1].(int32)
	p.data = p.data[:n-1]
}

func execSubInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) - p.data[n-1].(int64)
	p.data = p.data[:n-1]
}

func execSubUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) - p.data[n-1].(uint)
	p.data = p.data[:n-1]
}

func execSubUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) - p.data[n-1].(uint8)
	p.data = p.data[:n-1]
}

func execSubUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) - p.data[n-1].(uint16)
	p.data = p.data[:n-1]
}

func execSubUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) - p.data[n-1].(uint32)
	p.data = p.data[:n-1]
}

func execSubUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) - p.data[n-1].(uint64)
	p.data = p.data[:n-1]
}

func execSubUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) - p.data[n-1].(uintptr)
	p.data = p.data[:n-1]
}

func execSubFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float32) - p.data[n-1].(float32)
	p.data = p.data[:n-1]
}

func execSubFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float64) - p.data[n-1].(float64)
	p.data = p.data[:n-1]
}

func execSubComplex64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(complex64) - p.data[n-1].(complex64)
	p.data = p.data[:n-1]
}

func execSubComplex128(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(complex128) - p.data[n-1].(complex128)
	p.data = p.data[:n-1]
}

// -----------------------------------------------------------------------------

func execMulInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) * p.data[n-1].(int)
	p.data = p.data[:n-1]
}

func execMulInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) * p.data[n-1].(int8)
	p.data = p.data[:n-1]
}

func execMulInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) * p.data[n-1].(int16)
	p.data = p.data[:n-1]
}

func execMulInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) * p.data[n-1].(int32)
	p.data = p.data[:n-1]
}

func execMulInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) * p.data[n-1].(int64)
	p.data = p.data[:n-1]
}

func execMulUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) * p.data[n-1].(uint)
	p.data = p.data[:n-1]
}

func execMulUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) * p.data[n-1].(uint8)
	p.data = p.data[:n-1]
}

func execMulUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) * p.data[n-1].(uint16)
	p.data = p.data[:n-1]
}

func execMulUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) * p.data[n-1].(uint32)
	p.data = p.data[:n-1]
}

func execMulUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) * p.data[n-1].(uint64)
	p.data = p.data[:n-1]
}

func execMulUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) * p.data[n-1].(uintptr)
	p.data = p.data[:n-1]
}

func execMulFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float32) * p.data[n-1].(float32)
	p.data = p.data[:n-1]
}

func execMulFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float64) * p.data[n-1].(float64)
	p.data = p.data[:n-1]
}

func execMulComplex64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(complex64) * p.data[n-1].(complex64)
	p.data = p.data[:n-1]
}

func execMulComplex128(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(complex128) * p.data[n-1].(complex128)
	p.data = p.data[:n-1]
}

// -----------------------------------------------------------------------------

func execDivInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) / p.data[n-1].(int)
	p.data = p.data[:n-1]
}

func execDivInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) / p.data[n-1].(int8)
	p.data = p.data[:n-1]
}

func execDivInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) / p.data[n-1].(int16)
	p.data = p.data[:n-1]
}

func execDivInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) / p.data[n-1].(int32)
	p.data = p.data[:n-1]
}

func execDivInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) / p.data[n-1].(int64)
	p.data = p.data[:n-1]
}

func execDivUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) / p.data[n-1].(uint)
	p.data = p.data[:n-1]
}

func execDivUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) / p.data[n-1].(uint8)
	p.data = p.data[:n-1]
}

func execDivUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) / p.data[n-1].(uint16)
	p.data = p.data[:n-1]
}

func execDivUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) / p.data[n-1].(uint32)
	p.data = p.data[:n-1]
}

func execDivUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) / p.data[n-1].(uint64)
	p.data = p.data[:n-1]
}

func execDivUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) / p.data[n-1].(uintptr)
	p.data = p.data[:n-1]
}

func execDivFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float32) / p.data[n-1].(float32)
	p.data = p.data[:n-1]
}

func execDivFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float64) / p.data[n-1].(float64)
	p.data = p.data[:n-1]
}

func execDivComplex64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(complex64) / p.data[n-1].(complex64)
	p.data = p.data[:n-1]
}

func execDivComplex128(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(complex128) / p.data[n-1].(complex128)
	p.data = p.data[:n-1]
}

// -----------------------------------------------------------------------------

func execModInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) % p.data[n-1].(int)
	p.data = p.data[:n-1]
}

func execModInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) % p.data[n-1].(int8)
	p.data = p.data[:n-1]
}

func execModInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) % p.data[n-1].(int16)
	p.data = p.data[:n-1]
}

func execModInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) % p.data[n-1].(int32)
	p.data = p.data[:n-1]
}

func execModInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) % p.data[n-1].(int64)
	p.data = p.data[:n-1]
}

func execModUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) % p.data[n-1].(uint)
	p.data = p.data[:n-1]
}

func execModUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) % p.data[n-1].(uint8)
	p.data = p.data[:n-1]
}

func execModUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) % p.data[n-1].(uint16)
	p.data = p.data[:n-1]
}

func execModUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) % p.data[n-1].(uint32)
	p.data = p.data[:n-1]
}

func execModUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) % p.data[n-1].(uint64)
	p.data = p.data[:n-1]
}

func execModUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) % p.data[n-1].(uintptr)
	p.data = p.data[:n-1]
}

func execBuiltinOp(i Instr, ctx *Context) {
	if fn, ok := builtinOps[int(i&bitsOperand)]; ok {
		fn(i, ctx)
	} else {
		panic("execBuiltinOp: invalid builtinOp")
	}
}

// -----------------------------------------------------------------------------

const (
	bitsKind     = 5 // Kind count = 26
	bitsOperator = 5 // Operator count = 24
)

var builtinOps = map[int]func(i Instr, p *Context){
	(int(Int) << bitsOperator) | int(OpAdd):        execAddInt,
	(int(Int8) << bitsOperator) | int(OpAdd):       execAddInt8,
	(int(Int16) << bitsOperator) | int(OpAdd):      execAddInt16,
	(int(Int32) << bitsOperator) | int(OpAdd):      execAddInt32,
	(int(Int64) << bitsOperator) | int(OpAdd):      execAddInt64,
	(int(Uint) << bitsOperator) | int(OpAdd):       execAddUint,
	(int(Uint8) << bitsOperator) | int(OpAdd):      execAddUint8,
	(int(Uint16) << bitsOperator) | int(OpAdd):     execAddUint16,
	(int(Uint32) << bitsOperator) | int(OpAdd):     execAddUint32,
	(int(Uint64) << bitsOperator) | int(OpAdd):     execAddUint64,
	(int(Uintptr) << bitsOperator) | int(OpAdd):    execAddUintptr,
	(int(String) << bitsOperator) | int(OpAdd):     execAddString,
	(int(Float32) << bitsOperator) | int(OpAdd):    execAddFloat32,
	(int(Float64) << bitsOperator) | int(OpAdd):    execAddFloat64,
	(int(Complex64) << bitsOperator) | int(OpAdd):  execAddComplex64,
	(int(Complex128) << bitsOperator) | int(OpAdd): execAddComplex128,
	(int(Int) << bitsOperator) | int(OpSub):        execSubInt,
	(int(Int8) << bitsOperator) | int(OpSub):       execSubInt8,
	(int(Int16) << bitsOperator) | int(OpSub):      execSubInt16,
	(int(Int32) << bitsOperator) | int(OpSub):      execSubInt32,
	(int(Int64) << bitsOperator) | int(OpSub):      execSubInt64,
	(int(Uint) << bitsOperator) | int(OpSub):       execSubUint,
	(int(Uint8) << bitsOperator) | int(OpSub):      execSubUint8,
	(int(Uint16) << bitsOperator) | int(OpSub):     execSubUint16,
	(int(Uint32) << bitsOperator) | int(OpSub):     execSubUint32,
	(int(Uint64) << bitsOperator) | int(OpSub):     execSubUint64,
	(int(Uintptr) << bitsOperator) | int(OpSub):    execSubUintptr,
	(int(Float32) << bitsOperator) | int(OpSub):    execSubFloat32,
	(int(Float64) << bitsOperator) | int(OpSub):    execSubFloat64,
	(int(Complex64) << bitsOperator) | int(OpSub):  execSubComplex64,
	(int(Complex128) << bitsOperator) | int(OpSub): execSubComplex128,
	(int(Int) << bitsOperator) | int(OpMul):        execMulInt,
	(int(Int8) << bitsOperator) | int(OpMul):       execMulInt8,
	(int(Int16) << bitsOperator) | int(OpMul):      execMulInt16,
	(int(Int32) << bitsOperator) | int(OpMul):      execMulInt32,
	(int(Int64) << bitsOperator) | int(OpMul):      execMulInt64,
	(int(Uint) << bitsOperator) | int(OpMul):       execMulUint,
	(int(Uint8) << bitsOperator) | int(OpMul):      execMulUint8,
	(int(Uint16) << bitsOperator) | int(OpMul):     execMulUint16,
	(int(Uint32) << bitsOperator) | int(OpMul):     execMulUint32,
	(int(Uint64) << bitsOperator) | int(OpMul):     execMulUint64,
	(int(Uintptr) << bitsOperator) | int(OpMul):    execMulUintptr,
	(int(Float32) << bitsOperator) | int(OpMul):    execMulFloat32,
	(int(Float64) << bitsOperator) | int(OpMul):    execMulFloat64,
	(int(Complex64) << bitsOperator) | int(OpMul):  execMulComplex64,
	(int(Complex128) << bitsOperator) | int(OpMul): execMulComplex128,
	(int(Int) << bitsOperator) | int(OpDiv):        execDivInt,
	(int(Int8) << bitsOperator) | int(OpDiv):       execDivInt8,
	(int(Int16) << bitsOperator) | int(OpDiv):      execDivInt16,
	(int(Int32) << bitsOperator) | int(OpDiv):      execDivInt32,
	(int(Int64) << bitsOperator) | int(OpDiv):      execDivInt64,
	(int(Uint) << bitsOperator) | int(OpDiv):       execDivUint,
	(int(Uint8) << bitsOperator) | int(OpDiv):      execDivUint8,
	(int(Uint16) << bitsOperator) | int(OpDiv):     execDivUint16,
	(int(Uint32) << bitsOperator) | int(OpDiv):     execDivUint32,
	(int(Uint64) << bitsOperator) | int(OpDiv):     execDivUint64,
	(int(Uintptr) << bitsOperator) | int(OpDiv):    execDivUintptr,
	(int(Float32) << bitsOperator) | int(OpDiv):    execDivFloat32,
	(int(Float64) << bitsOperator) | int(OpDiv):    execDivFloat64,
	(int(Complex64) << bitsOperator) | int(OpDiv):  execDivComplex64,
	(int(Complex128) << bitsOperator) | int(OpDiv): execDivComplex128,
	(int(Int) << bitsOperator) | int(OpMod):        execModInt,
	(int(Int8) << bitsOperator) | int(OpMod):       execModInt8,
	(int(Int16) << bitsOperator) | int(OpMod):      execModInt16,
	(int(Int32) << bitsOperator) | int(OpMod):      execModInt32,
	(int(Int64) << bitsOperator) | int(OpMod):      execModInt64,
	(int(Uint) << bitsOperator) | int(OpMod):       execModUint,
	(int(Uint8) << bitsOperator) | int(OpMod):      execModUint8,
	(int(Uint16) << bitsOperator) | int(OpMod):     execModUint16,
	(int(Uint32) << bitsOperator) | int(OpMod):     execModUint32,
	(int(Uint64) << bitsOperator) | int(OpMod):     execModUint64,
	(int(Uintptr) << bitsOperator) | int(OpMod):    execModUintptr,
}

func (p *Code) builtinOp(kind Kind, op Operator) error {
	i := (int(kind) << bitsOperator) | int(op)
	if _, ok := builtinOps[i]; ok {
		p.data = append(p.data, (opBuiltinOp<<bitsOpShift)|uint32(i))
		return nil
	}
	return fmt.Errorf("builtinOp: type %v doesn't support operator %v", kind, op)
}

// BuiltinOp instr
func (ctx *Builder) BuiltinOp(kind Kind, op Operator) *Builder {
	err := ctx.code.builtinOp(kind, op)
	if err != nil {
		panic(err)
	}
	return ctx
}

// -----------------------------------------------------------------------------
