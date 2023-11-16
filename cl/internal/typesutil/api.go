/*
 * Copyright (c) 2023 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package typesutil

import (
	"go/constant"
	"go/types"
	"unsafe"

	"github.com/goplus/gox"
)

// An operandMode specifies the (addressing) mode of an operand.
type operandMode byte

const (
	invalid   operandMode = iota // operand is invalid
	novalue                      // operand represents no value (result of a function call w/o result)
	builtin                      // operand is a built-in function
	typexpr                      // operand is a type
	constant_                    // operand is a constant; the operand's typ is a Basic type
	variable                     // operand is an addressable variable
	mapindex                     // operand is a map index expression (acts like a variable on lhs, commaok on rhs of an assignment)
	value                        // operand is a computed value
	commaok                      // like value, but operand may be used in a comma,ok expression
	commaerr                     // like commaok, but second value is error, not boolean
	cgofunc                      // operand is a cgo function
)

// TypeAndValue reports the type and value (for constants)
// of the corresponding expression.
type TypeAndValue struct {
	mode  operandMode
	Type  types.Type
	Value constant.Value
}

func NewTypeAndValueForType(typ types.Type) (ret types.TypeAndValue) {
	switch t := typ.(type) {
	case *gox.TypeType:
		typ = t.Type()
	}
	ret.Type = types.Default(typ)
	(*TypeAndValue)(unsafe.Pointer(&ret)).mode = typexpr
	return
}

func NewTypeAndValueForValue(typ types.Type, val constant.Value) (ret types.TypeAndValue) {
	var mode operandMode
	if val != nil {
		mode = constant_
	} else {
		mode = value
	}
	ret.Type = types.Default(typ)
	ret.Value = val
	(*TypeAndValue)(unsafe.Pointer(&ret)).mode = mode
	return
}

func NewTypeAndValueForVariable(typ types.Type) (ret types.TypeAndValue) {
	ret.Type = typ
	(*TypeAndValue)(unsafe.Pointer(&ret)).mode = variable
	return
}

func NewTypeAndValueForCallResult(typ types.Type) (ret types.TypeAndValue) {
	if typ == nil {
		ret.Type = &types.Tuple{}
		(*TypeAndValue)(unsafe.Pointer(&ret)).mode = novalue
		return
	}
	ret.Type = typ
	(*TypeAndValue)(unsafe.Pointer(&ret)).mode = value
	return
}

func NewTypeAndValueForObject(obj types.Object) (ret types.TypeAndValue) {
	var mode operandMode
	var val constant.Value
	switch v := obj.(type) {
	case *types.Const:
		mode = constant_
		val = v.Val()
	case *types.TypeName:
		mode = typexpr
	case *types.Var:
		mode = variable
	case *types.Func:
		mode = value
	case *types.Builtin:
		mode = builtin
	case *types.Nil:
		mode = value
	}
	ret.Type = obj.Type()
	ret.Value = val
	(*TypeAndValue)(unsafe.Pointer(&ret)).mode = mode
	return
}
