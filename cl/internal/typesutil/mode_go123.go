//go:build go1.23
// +build go1.23

/*
 * Copyright (c) 2025 The XGo Authors (xgo.dev). All rights reserved.
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

const (
	Invalid  OperandMode = iota // operand is invalid
	NoValue                     // operand represents no value (result of a function call w/o result)
	Builtin                     // operand is a built-in function
	TypExpr                     // operand is a type
	Constant                    // operand is a constant; the operand's typ is a Basic type
	Variable                    // operand is an addressable variable
	MapIndex                    // operand is a map index expression (acts like a variable on lhs, commaok on rhs of an assignment)
	Value                       // operand is a computed value
	NilValue                    // operand is the nil value - only used by types2
	CommaOK                     // like value, but operand may be used in a comma,ok expression
	CommaErr                    // like commaok, but second value is error, not boolean
	CgoFunc                     // operand is a cgo function
)
