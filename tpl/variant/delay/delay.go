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

package delay

import (
	"github.com/goplus/xgo/tpl/token"
	"github.com/goplus/xgo/tpl/variant"
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

// EvalOp delays to evaluate a value.
func EvalOp(expr any, fn func(v any)) any {
	return func() any {
		fn(Eval(expr))
		return nil
	}
}

// List delays to convert the matching result of (R % ",") to a flat list.
// R % "," means R *("," R)
func List(in []any, fn func(flat []any)) any {
	return func() any {
		fn(variant.List(in))
		return nil
	}
}

// ListOp delays to convert the matching result of (R % ",") to a flat list.
// R % "," means R *("," R)
func ListOp[T any](in []any, op func(v any) T, fn func(flat []T)) any {
	return func() any {
		fn(variant.ListOp(in, op))
		return nil
	}
}

// RangeOp delays to travel the matching result of (R % ",") and call fn(result of R).
// R % "," means R *("," R)
func RangeOp(in []any, fn func(v any)) any {
	return func() any {
		variant.RangeOp(in, fn)
		return nil
	}
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
func Call(needList bool, name string, arglist any) any {
	return func() any {
		return variant.Call(needList, name, arglist)
	}
}

// CallObject delays to call a function object.
func CallObject(needList bool, fn any, arglist any) any {
	return func() any {
		return variant.CallObject(needList, fn, arglist)
	}
}

// -----------------------------------------------------------------------------

// ValueOf delays to get a value.
func ValueOf(name string, getval func(name string) (any, bool)) any {
	return func() any {
		v, ok := getval(name)
		if !ok {
			panic(name + " is undefined")
		}
		return v
	}
}

// SetValue delays to set a value.
func SetValue(name string, setval func(name string, val any), expr any) any {
	return func() any {
		setval(name, Eval(expr))
		return nil
	}
}

// ChgValue delays to change a value.
func ChgValue(name string, chgval func(string, func(oldv any) any), chg func(oldv any) any) any {
	return func() any {
		chgval(name, chg)
		return nil
	}
}

// -----------------------------------------------------------------------------

// StmtList delays a statement list.
func StmtList(stmts []any) any {
	return func() any {
		for _, stmt := range stmts {
			Eval(stmt)
		}
		return nil
	}
}

// IfElse delays an if-else statement.
func IfElse(cond, ifBody, elseStmt any, elseBodyAt int) any {
	return func() any {
		if Eval(cond).(bool) {
			Eval(ifBody)
		} else if elseStmt != nil {
			Eval(elseStmt.([]any)[elseBodyAt])
		}
		return nil
	}
}

// While delays a while loop.
func While(cond, body any) any {
	return func() any {
		for Eval(cond).(bool) {
			Eval(body)
		}
		return nil
	}
}

// RepeatUntil delays a repeat-until loop.
func RepeatUntil(body, cond any) any {
	return func() any {
		for {
			Eval(body)
			if Eval(cond).(bool) {
				break
			}
		}
		return nil
	}
}

// -----------------------------------------------------------------------------

// InitUniverse initializes the universe module with the specified modules.
func InitUniverse(names ...string) {
	variant.InitUniverse(names...)
}

// -----------------------------------------------------------------------------
