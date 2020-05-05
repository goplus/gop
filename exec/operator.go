package exec

import (
	"fmt"
	"reflect"
	"strconv"
	"unsafe"
)

// -----------------------------------------------------------------------------

// Operator type.
type Operator uint

const (
	// OpInvalid - invalid operator
	OpInvalid Operator = iota
	// OpAdd '+' String/Int/Uint/Float/Complex
	OpAdd
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

const (
	// SameAsFirst means the second argument is same as first argument type.
	SameAsFirst = reflect.Invalid
)

const (
	bitNone          = 0
	bitSameAsFirst   = 1 << SameAsFirst
	bitBool          = 1 << Bool
	bitInt           = 1 << Int
	bitInt8          = 1 << Int8
	bitInt16         = 1 << Int16
	bitInt32         = 1 << Int32
	bitInt64         = 1 << Int64
	bitUint          = 1 << Uint
	bitUint8         = 1 << Uint8
	bitUint16        = 1 << Uint16
	bitUint32        = 1 << Uint32
	bitUint64        = 1 << Uint64
	bitUintptr       = 1 << Uintptr
	bitFloat32       = 1 << Float32
	bitFloat64       = 1 << Float64
	bitComplex64     = 1 << Complex64
	bitComplex128    = 1 << Complex128
	bitString        = 1 << String
	bitUnsafePointer = 1 << UnsafePointer

	bitsAllInt     = bitInt | bitInt8 | bitInt16 | bitInt32 | bitInt64
	bitsAllUint    = bitUint | bitUint8 | bitUint16 | bitUint32 | bitUint64 | bitUintptr
	bitsAllIntUint = bitsAllInt | bitsAllUint
	bitsAllFloat   = bitFloat32 | bitFloat64
	bitsAllReal    = bitsAllIntUint | bitsAllFloat
	bitsAllComplex = bitComplex64 | bitComplex128
	bitsAllNumber  = bitsAllReal | bitsAllComplex
)

// OperatorInfo represents an operator information.
type OperatorInfo struct {
	Lit      string
	InFirst  uint64       // first argument supported types.
	InSecond uint64       // second argument supported types. It may have SameAsFirst flag.
	Out      reflect.Kind // result type. It may be SameAsFirst.
}

var opInfos = [...]OperatorInfo{
	OpAdd:       {"+", bitsAllNumber | bitString, bitSameAsFirst, SameAsFirst},
	OpSub:       {"-", bitsAllNumber, bitSameAsFirst, SameAsFirst},
	OpMul:       {"*", bitsAllNumber, bitSameAsFirst, SameAsFirst},
	OpDiv:       {"/", bitsAllNumber, bitSameAsFirst, SameAsFirst},
	OpMod:       {"%", bitsAllIntUint, bitSameAsFirst, SameAsFirst},
	OpBitAnd:    {"&", bitsAllIntUint, bitSameAsFirst, SameAsFirst},
	OpBitOr:     {"|", bitsAllIntUint, bitSameAsFirst, SameAsFirst},
	OpBitXor:    {"^", bitsAllIntUint, bitSameAsFirst, SameAsFirst},
	OpBitAndNot: {"&^", bitsAllIntUint, bitSameAsFirst, SameAsFirst},
	OpBitSHL:    {"<<", bitsAllIntUint, bitsAllIntUint, SameAsFirst},
	OpBitSHR:    {">>", bitsAllIntUint, bitsAllIntUint, SameAsFirst},
	OpLT:        {"<", bitsAllReal | bitString, bitSameAsFirst, Bool},
	OpLE:        {"<=", bitsAllReal | bitString, bitSameAsFirst, Bool},
	OpGT:        {">", bitsAllReal | bitString, bitSameAsFirst, Bool},
	OpGE:        {">=", bitsAllReal | bitString, bitSameAsFirst, Bool},
	OpEQ:        {"==", bitsAllNumber | bitString, bitSameAsFirst, Bool},
	OpEQNil:     {"== nil", bitUnsafePointer, bitNone, Bool},
	OpNE:        {"!=", bitsAllNumber | bitString, bitSameAsFirst, Bool},
	OpNENil:     {"!= nil", bitUnsafePointer, bitNone, Bool},
	OpLAnd:      {"&&", bitBool, bitBool, Bool},
	OpLOr:       {"||", bitBool, bitBool, Bool},
	OpNeg:       {"-", bitsAllNumber, bitNone, SameAsFirst},
	OpNot:       {"!", bitBool, bitNone, Bool},
	OpBitNot:    {"^", bitsAllIntUint, bitNone, SameAsFirst},
}

// GetInfo returns the information of this operator.
func (op Operator) GetInfo() *OperatorInfo {
	return &opInfos[op]
}

func (op Operator) String() string {
	if int(op) < len(opInfos) {
		return opInfos[op].Lit
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

var (
	// TyBool type
	TyBool = reflect.TypeOf(true)
	// TyInt type
	TyInt = reflect.TypeOf(int(0))
	// TyInt8 type
	TyInt8 = reflect.TypeOf(int8(0))
	// TyInt16 type
	TyInt16 = reflect.TypeOf(int16(0))
	// TyInt32 type
	TyInt32 = reflect.TypeOf(int32(0))
	// TyInt64 type
	TyInt64 = reflect.TypeOf(int64(0))
	// TyUint type
	TyUint = reflect.TypeOf(uint(0))
	// TyUint8 type
	TyUint8 = reflect.TypeOf(uint8(0))
	// TyUint16 type
	TyUint16 = reflect.TypeOf(uint16(0))
	// TyUint32 type
	TyUint32 = reflect.TypeOf(uint32(0))
	// TyUint64 type
	TyUint64 = reflect.TypeOf(uint64(0))
	// TyUintptr type
	TyUintptr = reflect.TypeOf(uintptr(0))
	// TyFloat32 type
	TyFloat32 = reflect.TypeOf(float32(0))
	// TyFloat64 type
	TyFloat64 = reflect.TypeOf(float64(0))
	// TyComplex64 type
	TyComplex64 = reflect.TypeOf(complex64(0))
	// TyComplex128 type
	TyComplex128 = reflect.TypeOf(complex128(0))
	// TyString type
	TyString = reflect.TypeOf("")
	// TyUnsafePointer type
	TyUnsafePointer = reflect.TypeOf(unsafe.Pointer(nil))
	// TyEmptyInterface type
	TyEmptyInterface = reflect.TypeOf((*interface{})(nil)).Elem()
)

type bTI struct { // builtin type info
	typ  reflect.Type
	size uintptr
}

var builtinTypes = [...]bTI{
	Bool:          {TyBool, 1},
	Int:           {TyInt, unsafe.Sizeof(int(0))},
	Int8:          {TyInt8, 1},
	Int16:         {TyInt16, 2},
	Int32:         {TyInt32, 4},
	Int64:         {TyInt64, 8},
	Uint:          {TyUint, unsafe.Sizeof(uint(0))},
	Uint8:         {TyUint8, 1},
	Uint16:        {TyUint16, 2},
	Uint32:        {TyUint32, 4},
	Uint64:        {TyUint64, 8},
	Uintptr:       {TyUintptr, unsafe.Sizeof(uintptr(0))},
	Float32:       {TyFloat32, 4},
	Float64:       {TyFloat64, 8},
	Complex64:     {TyComplex64, 8},
	Complex128:    {TyComplex128, 16},
	String:        {TyString, unsafe.Sizeof(string('0'))},
	UnsafePointer: {TyUnsafePointer, unsafe.Sizeof(uintptr(0))},
}

// TypeFromKind returns the type who has this kind.
func TypeFromKind(kind Kind) reflect.Type {
	return builtinTypes[kind].typ
}

// SizeofKind returns sizeof type who has this kind.
func SizeofKind(kind Kind) uintptr {
	return builtinTypes[kind].size
}

// -----------------------------------------------------------------------------

func execBuiltinOp(i Instr, p *Context) {
	if fn := builtinOps[int(i&bitsOperand)]; fn != nil {
		fn(0, p)
	} else {
		panic("execBuiltinOp: invalid builtinOp")
	}
}

// -----------------------------------------------------------------------------

const (
	bitsKind     = 5 // Kind count = 26
	bitsOperator = 5 // Operator count = 24
)

func (p *Code) builtinOp(kind Kind, op Operator) error {
	i := (int(kind) << bitsOperator) | int(op)
	if fn := builtinOps[i]; fn != nil {
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

// CallBuiltinOp calls BuiltinOp
func CallBuiltinOp(kind Kind, op Operator, data ...interface{}) interface{} {
	if fn := builtinOps[(int(kind)<<bitsOperator)|int(op)]; fn != nil {
		ctx := newSimpleContext(data)
		fn(0, ctx)
		return ctx.Get(-1)
	}
	panic("CallBuiltinOp: invalid builtinOp")
}

// -----------------------------------------------------------------------------
