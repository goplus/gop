package exec

import "reflect"

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

const (
	bitsKind     = 5 // Kind count = 26
	bitsOperator = 5 // Operator count = 24
)

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

// -----------------------------------------------------------------------------

// BuiltinOp instr
func (p *Code) BuiltinOp(kind Kind, op Operator) {
}

// -----------------------------------------------------------------------------
