/*
 Copyright 2020 The GoPlus Authors (goplus.org)

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package exec

import (
	"reflect"
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
	// OpQuo '/' Int/Uint/Float/Complex
	OpQuo
	// OpMod '%' Int/Uint
	OpMod
	// OpAnd '&' Int/Uint
	OpAnd
	// OpOr '|' Int/Uint
	OpOr
	// OpXor '^' Int/Uint
	OpXor
	// OpAndNot '&^' Int/Uint
	OpAndNot
	// OpLsh '<<' Int/Uint, Uint
	OpLsh
	// OpRsh '>>' Int/Uint, Uint
	OpRsh
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
	// OpLNot '!'
	OpLNot
	// OpNeg '-'
	OpNeg
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
	bitPtr           = 1 << reflect.Ptr
	bitBigInt        = 1 << BigInt
	bitBigRat        = 1 << BigRat
	bitBigFloat      = 1 << BigFloat

	bitsAllInt     = bitInt | bitInt8 | bitInt16 | bitInt32 | bitInt64
	bitsAllUint    = bitUint | bitUint8 | bitUint16 | bitUint32 | bitUint64 | bitUintptr
	bitsAllIntUint = bitsAllInt | bitsAllUint
	bitsAllFloat   = bitFloat32 | bitFloat64
	bitsAllReal    = bitsAllIntUint | bitsAllFloat | bitBigRat | bitBigFloat
	bitsAllComplex = bitComplex64 | bitComplex128
	bitsAllNumber  = bitsAllReal | bitsAllComplex
	bitsAllPtr     = bitPtr | bitUintptr | bitUnsafePointer
)

// OperatorInfo represents an operator information.
type OperatorInfo struct {
	Lit      string
	InFirst  uint64       // first argument supported types.
	InSecond uint64       // second argument supported types. It may have SameAsFirst flag.
	Out      reflect.Kind // result type. It may be SameAsFirst.
}

var opInfos = [...]OperatorInfo{
	OpAdd:    {"+", bitsAllNumber | bitString | bitBigInt, bitSameAsFirst, SameAsFirst},
	OpSub:    {"-", bitsAllNumber | bitBigInt, bitSameAsFirst, SameAsFirst},
	OpMul:    {"*", bitsAllNumber | bitBigInt, bitSameAsFirst, SameAsFirst},
	OpQuo:    {"/", bitsAllNumber, bitSameAsFirst, SameAsFirst}, // NOTE: not include `bigint/bigint`
	OpMod:    {"%", bitsAllIntUint | bitBigInt, bitSameAsFirst, SameAsFirst},
	OpAnd:    {"&", bitsAllIntUint | bitBigInt, bitSameAsFirst, SameAsFirst},
	OpOr:     {"|", bitsAllIntUint | bitBigInt, bitSameAsFirst, SameAsFirst},
	OpXor:    {"^", bitsAllIntUint | bitBigInt, bitSameAsFirst, SameAsFirst},
	OpAndNot: {"&^", bitsAllIntUint | bitBigInt, bitSameAsFirst, SameAsFirst},
	OpLsh:    {"<<", bitsAllIntUint | bitBigInt, bitsAllIntUint, SameAsFirst},
	OpRsh:    {">>", bitsAllIntUint | bitBigInt, bitsAllIntUint, SameAsFirst},
	OpLT:     {"<", bitsAllReal | bitString | bitBigInt, bitSameAsFirst, Bool},
	OpLE:     {"<=", bitsAllReal | bitString | bitBigInt, bitSameAsFirst, Bool},
	OpGT:     {">", bitsAllReal | bitString | bitBigInt, bitSameAsFirst, Bool},
	OpGE:     {">=", bitsAllReal | bitString | bitBigInt, bitSameAsFirst, Bool},
	OpEQ:     {"==", bitsAllNumber | bitBool | bitString | bitBigInt | bitUnsafePointer, bitSameAsFirst, Bool},
	OpEQNil:  {"== nil", bitUnsafePointer, bitNone, Bool},
	OpNE:     {"!=", bitsAllNumber | bitBool | bitString | bitBigInt | bitUnsafePointer, bitSameAsFirst, Bool},
	OpNENil:  {"!= nil", bitUnsafePointer, bitNone, Bool},
	OpLAnd:   {"&&", bitBool, bitBool, Bool},
	OpLOr:    {"||", bitBool, bitBool, Bool},
	OpNeg:    {"-", bitsAllNumber | bitBigInt, bitNone, SameAsFirst},
	OpLNot:   {"!", bitBool, bitNone, Bool},
	OpBitNot: {"^", bitsAllIntUint | bitBigInt, bitNone, SameAsFirst},
}

// GetInfo returns the information of this operator.
func (op Operator) GetInfo() *OperatorInfo {
	return &opInfos[op]
}

func (op Operator) String() string {
	return opInfos[op].Lit
}

// -----------------------------------------------------------------------------

// AddrOperator type.
type AddrOperator Operator

const (
	// OpAddrVal `*addr`
	OpAddrVal = AddrOperator(0)
	// OpAddAssign `+=`
	OpAddAssign = AddrOperator(OpAdd)
	// OpSubAssign `-=`
	OpSubAssign = AddrOperator(OpSub)
	// OpMulAssign `*=`
	OpMulAssign = AddrOperator(OpMul)
	// OpQuoAssign `/=`
	OpQuoAssign = AddrOperator(OpQuo)
	// OpModAssign `%=`
	OpModAssign = AddrOperator(OpMod)

	// OpAndAssign '&='
	OpAndAssign = AddrOperator(OpAnd)
	// OpOrAssign '|='
	OpOrAssign = AddrOperator(OpOr)
	// OpXorAssign '^='
	OpXorAssign = AddrOperator(OpXor)
	// OpAndNotAssign '&^='
	OpAndNotAssign = AddrOperator(OpAndNot)
	// OpLshAssign '<<='
	OpLshAssign = AddrOperator(OpLsh)
	// OpRshAssign '>>='
	OpRshAssign = AddrOperator(OpRsh)
	// OpAssign `=`
	OpAssign AddrOperator = iota
	// OpInc '++'
	OpInc
	// OpDec '--'
	OpDec
)

// AddrOperatorInfo represents an addr-operator information.
type AddrOperatorInfo struct {
	Lit      string
	InFirst  uint64 // first argument supported types.
	InSecond uint64 // second argument supported types. It may have SameAsFirst flag.
}

var addropInfos = [...]AddrOperatorInfo{
	OpAddAssign:    {"+=", bitsAllNumber | bitString | bitBigInt, bitSameAsFirst},
	OpSubAssign:    {"-=", bitsAllNumber | bitBigInt, bitSameAsFirst},
	OpMulAssign:    {"*=", bitsAllNumber | bitBigInt, bitSameAsFirst},
	OpQuoAssign:    {"/=", bitsAllNumber | bitBigInt, bitSameAsFirst},
	OpModAssign:    {"%=", bitsAllIntUint, bitSameAsFirst},
	OpAndAssign:    {"&=", bitsAllIntUint, bitSameAsFirst},
	OpOrAssign:     {"|=", bitsAllIntUint, bitSameAsFirst},
	OpXorAssign:    {"^=", bitsAllIntUint, bitSameAsFirst},
	OpAndNotAssign: {"&^=", bitsAllIntUint, bitSameAsFirst},
	OpLshAssign:    {"<<=", bitsAllIntUint, bitsAllIntUint},
	OpRshAssign:    {">>=", bitsAllIntUint, bitsAllIntUint},
	OpInc:          {"++", bitsAllNumber | bitBigInt, bitNone},
	OpDec:          {"--", bitsAllNumber | bitBigInt, bitNone},
}

// GetInfo returns the information of this operator.
func (op AddrOperator) GetInfo() *AddrOperatorInfo {
	return &addropInfos[op]
}

func (op AddrOperator) String() string {
	switch op {
	case OpAddrVal:
		return "*"
	case OpAssign:
		return "="
	default:
		return addropInfos[op].Lit
	}
}

// -----------------------------------------------------------------------------

// GoBuiltin represents go builtin func.
type GoBuiltin uint

func (p GoBuiltin) String() string {
	return goBuiltinNames[p]
}

const (
	gobInvalid GoBuiltin = iota
	// GobLen - len: 1
	GobLen
	// GobCap - cap: 2
	GobCap
	// GobCopy - copy: 3
	GobCopy
	// GobDelete - delete: 4
	GobDelete
	// GobComplex - complex: 5
	GobComplex
	// GobReal - real: 6
	GobReal
	// GobImag - imag: 7
	GobImag
	// GobClose - close: 8
	GobClose
	// GobRecover - recover:9
	GobRecover
)

var goBuiltinNames = [...]string{
	GobLen:     "len",
	GobCap:     "cap",
	GobCopy:    "copy",
	GobDelete:  "delete",
	GobComplex: "complex",
	GobReal:    "real",
	GobImag:    "imag",
	GobClose:   "close",
	GobRecover: "recover",
}

// -----------------------------------------------------------------------------

// Operator related constants, for Operator/AddrOperator instr.
const (
	// BitNone - bitNone
	BitNone = bitNone
	// BitsAllIntUint - bitsAllIntUint
	BitsAllIntUint = bitsAllIntUint
)

// -----------------------------------------------------------------------------

// Slice related constants, for Slice/Slice3 instr.
const (
	// SliceConstIndexLast - slice const index max
	SliceConstIndexLast = (1 << 13) - 3
	// SliceDefaultIndex - unspecified index
	SliceDefaultIndex = -2
)

// -----------------------------------------------------------------------------
