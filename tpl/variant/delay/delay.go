/*
 * Copyright (c) 2025 The GoPlus Authors (goplus.org). All rights reserved.
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

package delay

import (
	"github.com/goplus/gop/tpl/token"
	"github.com/goplus/gop/tpl/variant"
)

// -----------------------------------------------------------------------------

// Value represents a delayed value.
type Value = func() any

// Eval evaluates a value.
func Eval(v any) any {
	if d, ok := v.(Value); ok {
		return d()
	}
	return v
}

// -----------------------------------------------------------------------------

// Compare delays a compare operation.
func Compare(op token.Token, x, y any) any {
	return func() any {
		return variant.Compare(op, x, y)
	}
}

// MathOp delays a math operation.
func MathOp(op token.Token, x, y any) any {
	return func() any {
		return variant.MathOp(op, x, y)
	}
}

// LogicOp delays a logic operation.
func LogicOp(op token.Token, x, y any) any {
	return func() any {
		return variant.LogicOp(op, x, y)
	}
}

// UnaryOp delays a unary operation.
func UnaryOp(op token.Token, x any) any {
	return func() any {
		return variant.UnaryOp(op, x)
	}
}

// Call delays to call a function.
func Call(name string, args ...any) any {
	return func() any {
		return variant.Call(name, args...)
	}
}

// -----------------------------------------------------------------------------

// InitUniverse initializes the universe module with the specified modules.
func InitUniverse(names ...string) {
	variant.InitUniverse(names...)
}

// -----------------------------------------------------------------------------
