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

package variant

import (
	"github.com/goplus/gop/tpl/token"
)

// -----------------------------------------------------------------------------

// DelayValue represents a delayed value.
type DelayValue = func() any

// Delay delays a value.
func Delay(getter DelayValue) any {
	return getter
}

// Eval evaluates a value.
func Eval(v any) any {
	if d, ok := v.(DelayValue); ok {
		return d()
	}
	return v
}

// -----------------------------------------------------------------------------

func cmpInt(op token.Token, x, y int) bool {
	switch op {
	case token.EQ, token.ASSIGN:
		return x == y
	case token.NE, token.BIDIARROW:
		return x != y
	case token.LT:
		return x < y
	case token.LE:
		return x <= y
	case token.GT:
		return x > y
	case token.GE:
		return x >= y
	}
	panic("unreachable")
}

func cmpFloat(op token.Token, x, y float64) bool {
	switch op {
	case token.EQ, token.ASSIGN:
		return x == y
	case token.NE, token.BIDIARROW:
		return x != y
	case token.LT:
		return x < y
	case token.LE:
		return x <= y
	case token.GT:
		return x > y
	case token.GE:
		return x >= y
	}
	panic("unreachable")
}

func cmpString(op token.Token, x, y string) bool {
	switch op {
	case token.EQ, token.ASSIGN:
		return x == y
	case token.NE, token.BIDIARROW:
		return x != y
	case token.LT:
		return x < y
	case token.LE:
		return x <= y
	case token.GT:
		return x > y
	case token.GE:
		return x >= y
	}
	panic("unreachable")
}

func cmpBool(op token.Token, x, y bool) bool {
	switch op {
	case token.EQ, token.ASSIGN:
		return x == y
	case token.NE, token.BIDIARROW:
		return x != y
	}
	panic("unreachable")
}

// Compare compares two values.
func Compare(op token.Token, x, y any) bool {
	x, y = Eval(x), Eval(y)
	switch x := x.(type) {
	case int:
		switch y := y.(type) {
		case int:
			return cmpInt(op, x, y)
		case float64:
			return cmpFloat(op, float64(x), y)
		}
	case float64:
		switch y := y.(type) {
		case int:
			return cmpFloat(op, x, float64(y))
		case float64:
			return cmpFloat(op, x, y)
		}
	case string:
		if y, ok := y.(string); ok {
			return cmpString(op, x, y)
		}
	case bool:
		if y, ok := y.(bool); ok {
			return cmpBool(op, x, y)
		}
	}
	panic("compare: invalid operation")
}

// DelayCompare delays a compare operation.
func DelayCompare(op token.Token, x, y any) any {
	return func() any {
		return Compare(op, x, y)
	}
}

// -----------------------------------------------------------------------------

func mopInt(op token.Token, x, y int) int {
	switch op {
	case token.ADD:
		return x + y
	case token.SUB:
		return x - y
	case token.MUL:
		return x * y
	case token.QUO:
		return x / y
	case token.REM:
		return x % y
	}
	panic("unreachable")
}

func mopFloat(op token.Token, x, y float64) float64 {
	switch op {
	case token.ADD:
		return x + y
	case token.SUB:
		return x - y
	case token.MUL:
		return x * y
	case token.QUO:
		return x / y
	}
	panic("unreachable")
}

func mopString(op token.Token, x, y string) string {
	switch op {
	case token.ADD:
		return x + y
	}
	panic("unreachable")
}

// MathOp does a math operation.
func MathOp(op token.Token, x, y any) any {
	x, y = Eval(x), Eval(y)
	switch x := x.(type) {
	case int:
		switch y := y.(type) {
		case int:
			return mopInt(op, x, y)
		case float64:
			return mopFloat(op, float64(x), y)
		}
	case float64:
		switch y := y.(type) {
		case int:
			return mopFloat(op, x, float64(y))
		case float64:
			return mopFloat(op, x, y)
		}
	case string:
		if y, ok := y.(string); ok {
			return mopString(op, x, y)
		}
	}
	panic("mathOp: invalid operation")
}

// DelayMathOp delays a math operation.
func DelayMathOp(op token.Token, x, y any) any {
	return func() any {
		return MathOp(op, x, y)
	}
}

// -----------------------------------------------------------------------------

// LogicOp does a logic operation.
func LogicOp(op token.Token, x, y any) bool {
	x, y = Eval(x), Eval(y)
	switch x := x.(type) {
	case bool:
		if y, ok := y.(bool); ok {
			switch op {
			case token.LAND:
				return x && y
			case token.LOR:
				return x || y
			}
		}
	}
	panic("logicOp: invalid operation")
}

// DelayLogicOp delays a logic operation.
func DelayLogicOp(op token.Token, x, y any) any {
	return func() any {
		return LogicOp(op, x, y)
	}
}

// -----------------------------------------------------------------------------

// UnaryOp does a unary operation.
func UnaryOp(op token.Token, x any) any {
	x = Eval(x)
	switch x := x.(type) {
	case int:
		switch op {
		case token.SUB:
			return -x
		case token.ADD:
			return x
		}
	case float64:
		switch op {
		case token.SUB:
			return -x
		case token.ADD:
			return x
		}
	case bool:
		if op == token.NOT {
			return !x
		}
	}
	panic("unaryOp: invalid operation")
}

// DelayUnaryOp delays a unary operation.
func DelayUnaryOp(op token.Token, x any) any {
	return func() any {
		return UnaryOp(op, x)
	}
}

// -----------------------------------------------------------------------------
