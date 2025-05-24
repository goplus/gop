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

package math

import (
	"math"

	"github.com/goplus/xgo/tpl/token"
	"github.com/goplus/xgo/tpl/variant"
)

// -----------------------------------------------------------------------------

type namedFn struct {
	name string
	fn   func(...any) any
}

type math1 struct {
	name string
	fn   func(float64) float64
}

type math2 struct {
	name string
	fn   func(x, y float64) float64
}

type mathIntFloat struct {
	name string
	fn   func(x int, y float64) float64
}

func float1(args []any) float64 {
	if len(args) != 1 {
		panic("function call: arity mismatch")
	}
	return variant.Float(args[0])
}

func float2(args []any) (float64, float64) {
	if len(args) != 2 {
		panic("function call: arity mismatch")
	}
	return variant.Float(args[0]), variant.Float(args[1])
}

func intFloat(args []any) (int, float64) {
	if len(args) != 2 {
		panic("function call: arity mismatch")
	}
	return variant.Int(args[0]), variant.Float(args[1])
}

func regMath1(mod *variant.Module, name string, fn func(float64) float64) {
	mod.Insert(name, func(args ...any) any {
		return fn(float1(args))
	})
}

func regMath2(mod *variant.Module, name string, fn func(x, y float64) float64) {
	mod.Insert(name, func(args ...any) any {
		return fn(float2(args))
	})
}

func regMathIntFloat(mod *variant.Module, name string, fn func(x int, y float64) float64) {
	mod.Insert(name, func(args ...any) any {
		return fn(intFloat(args))
	})
}

// -----------------------------------------------------------------------------

func topValue(op token.Token, args []any) any {
	ret := variant.Eval(args[0])
	for _, v := range args[1:] {
		if variant.Compare(op, v, ret) {
			ret = v
		}
	}
	return ret
}

// Max returns the maximum value.
func Max(args ...any) any {
	if len(args) == 0 {
		panic("max: arity mismatch")
	}
	return topValue(token.GT, args)
}

// Min returns the minimum value.
func Min(args ...any) any {
	if len(args) == 0 {
		panic("min: arity mismatch")
	}
	return topValue(token.LT, args)
}

// -----------------------------------------------------------------------------

var fnsMath1 = [...]math1{
	{"abs", math.Abs},
	{"acos", math.Acos},
	{"acosh", math.Acosh},
	{"asin", math.Asin},
	{"asinh", math.Asinh},
	{"atan", math.Atan},
	{"atanh", math.Atanh},
	{"ceil", math.Ceil},
	{"cos", math.Cos},
	{"cosh", math.Cosh},
	{"erf", math.Erf},
	{"erfc", math.Erfc},
	{"exp", math.Exp},
	{"expm1", math.Expm1},
	{"floor", math.Floor},
	{"gamma", math.Gamma},
	{"j0", math.J0},
	{"j1", math.J1},
	{"log", math.Log},
	{"log10", math.Log10},
	{"log1p", math.Log1p},
	{"log2", math.Log2},
	{"round", math.Round},
	{"roundToEven", math.RoundToEven},
	{"sin", math.Sin},
	{"sinh", math.Sinh},
	{"sqrt", math.Sqrt},
	{"tan", math.Tan},
	{"tanh", math.Tanh},
	{"trunc", math.Trunc},
	{"y0", math.Y0},
	{"y1", math.Y1},
}

var fnsMath2 = [...]math2{
	{"atan2", math.Atan2},
	{"copysign", math.Copysign},
	{"dim", math.Dim},
	{"hypot", math.Hypot},
	{"mod", math.Mod},
	{"nextafter", math.Nextafter},
	{"pow", math.Pow},
	{"remainder", math.Remainder},
}

var fnsMathIntFloat = [...]mathIntFloat{
	{"jn", math.Jn},
	{"yn", math.Yn},
}

var fnsNamed = [...]namedFn{
	{"max", Max},
	{"min", Min},
}

func init() {
	mod := variant.NewModule("math")
	for _, m := range fnsMath1 {
		regMath1(mod, m.name, m.fn)
	}
	for _, m := range fnsMath2 {
		regMath2(mod, m.name, m.fn)
	}
	for _, m := range fnsMathIntFloat {
		regMathIntFloat(mod, m.name, m.fn)
	}
	for _, m := range fnsNamed {
		mod.Insert(m.name, m.fn)
	}
}

// -----------------------------------------------------------------------------
