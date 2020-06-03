/*
 Copyright 2020 Qiniu Cloud (qiniu.com)

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

package bytecode

import (
	"fmt"
	"reflect"

	exec "github.com/qiniu/qlang/v6/exec.spec"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

// Operator type.
type Operator = exec.Operator

const (
	// OpInvalid - invalid operator
	OpInvalid = exec.OpInvalid
	// OpAdd '+' String/Int/Uint/Float/Complex
	OpAdd = exec.OpAdd
	// OpSub '-' Int/Uint/Float/Complex
	OpSub = exec.OpSub
	// OpMul '*' Int/Uint/Float/Complex
	OpMul = exec.OpMul
	// OpDiv '/' Int/Uint/Float/Complex
	OpDiv = exec.OpDiv
	// OpMod '%' Int/Uint
	OpMod = exec.OpMod
	// OpBitAnd '&' Int/Uint
	OpBitAnd = exec.OpBitAnd
	// OpBitOr '|' Int/Uint
	OpBitOr = exec.OpBitOr
	// OpBitXor '^' Int/Uint
	OpBitXor = exec.OpBitXor
	// OpBitAndNot '&^' Int/Uint
	OpBitAndNot = exec.OpBitAndNot
	// OpBitSHL '<<' Int/Uint, Uint
	OpBitSHL = exec.OpBitSHL
	// OpBitSHR '>>' Int/Uint, Uint
	OpBitSHR = exec.OpBitSHR
	// OpLT '<' String/Int/Uint/Float
	OpLT = exec.OpLT
	// OpLE '<=' String/Int/Uint/Float
	OpLE = exec.OpLE
	// OpGT '>' String/Int/Uint/Float
	OpGT = exec.OpGT
	// OpGE '>=' String/Int/Uint/Float
	OpGE = exec.OpGE
	// OpEQ '==' ComparableType
	// Slice, map, and function values are not comparable. However, as a special case, a slice, map,
	// or function value may be compared to the predeclared identifier nil.
	OpEQ = exec.OpEQ
	// OpEQNil '==' nil
	OpEQNil = exec.OpEQNil
	// OpNE '!=' ComparableType
	OpNE = exec.OpNE
	// OpNENil '!=' nil
	OpNENil = exec.OpNENil
	// OpLAnd '&&' Bool
	OpLAnd = exec.OpLAnd
	// OpLOr '||' Bool
	OpLOr = exec.OpLOr
	// OpNeg '-'
	OpNeg = exec.OpNeg
	// OpNot '!'
	OpNot = exec.OpNot
	// OpBitNot '^'
	OpBitNot = exec.OpBitNot
)

const (
	// SameAsFirst means the second argument is same as first argument type.
	SameAsFirst = exec.SameAsFirst
)

// OperatorInfo represents an operator information.
type OperatorInfo = exec.OperatorInfo

// -----------------------------------------------------------------------------

// A Kind represents the specific kind of type that a Type represents.
type Kind = exec.Kind

const (
	// Bool type
	Bool = exec.Bool
	// Int type
	Int = exec.Int
	// Int8 type
	Int8 = exec.Int8
	// Int16 type
	Int16 = exec.Int16
	// Int32 type
	Int32 = exec.Int32
	// Int64 type
	Int64 = exec.Int64
	// Uint type
	Uint = exec.Uint
	// Uint8 type
	Uint8 = exec.Uint8
	// Uint16 type
	Uint16 = exec.Uint16
	// Uint32 type
	Uint32 = exec.Uint32
	// Uint64 type
	Uint64 = exec.Uint64
	// Uintptr type
	Uintptr = exec.Uintptr
	// Float32 type
	Float32 = exec.Float32
	// Float64 type
	Float64 = exec.Float64
	// Complex64 type
	Complex64 = exec.Complex64
	// Complex128 type
	Complex128 = exec.Complex128
	// String type
	String = exec.String
	// UnsafePointer type
	UnsafePointer = exec.UnsafePointer
)

// -----------------------------------------------------------------------------

func toUint(v interface{}) uint {
	switch n := v.(type) {
	case int:
		return uint(n)
	case uint:
		return n
	case uint32:
		return uint(n)
	case int32:
		return uint(n)
	case uint64:
		return uint(n)
	case int64:
		return uint(n)
	case uintptr:
		return uint(n)
	case uint16:
		return uint(n)
	case int16:
		return uint(n)
	case uint8:
		return uint(n)
	case int8:
		return uint(n)
	default:
		log.Panicln("toUint failed: unsupport type -", reflect.TypeOf(v))
		return 0
	}
}

func execBuiltinOp(i Instr, p *Context) {
	if fn := builtinOps[int(i&bitsOperand)]; fn != nil {
		fn(0, p)
		return
	}
	log.Panicln("execBuiltinOp: invalid instr -", i)
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
