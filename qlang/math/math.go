package math

import (
	"math"
	"reflect"
	"strings"

	qlang "qlang.io/spec"
)

// -----------------------------------------------------------------------------

func init() {

	fnt := qlang.Fntable

	fnt["e"] = math.E
	fnt["pi"] = math.Pi
	fnt["phi"] = math.Phi
	fnt["Inf"] = math.Inf(1)
	fnt["NaN"] = math.NaN()
}

// -----------------------------------------------------------------------------

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"_name":     "math",
	"abs":       math.Abs,
	"acos":      math.Acos,
	"acosh":     math.Acosh,
	"asin":      math.Asin,
	"asinh":     math.Asinh,
	"atan":      math.Atan,
	"atan2":     math.Atan2,
	"atanh":     math.Atanh,
	"cbrt":      math.Cbrt,
	"ceil":      math.Ceil,
	"copysign":  math.Copysign,
	"cos":       math.Cos,
	"cosh":      math.Cosh,
	"dim":       math.Dim,
	"erf":       math.Erf,
	"erfc":      math.Erfc,
	"exp":       math.Exp,
	"exp2":      math.Exp2,
	"expm1":     math.Expm1,
	"floor":     math.Floor,
	"gamma":     math.Gamma,
	"hypot":     math.Hypot,
	"inf":       math.Inf,
	"j0":        math.J0,
	"j1":        math.J1,
	"jn":        math.Jn,
	"ldexp":     math.Ldexp,
	"ln":        math.Log,
	"log":       math.Log,
	"log10":     math.Log10,
	"log1p":     math.Log1p,
	"log2":      math.Log2,
	"logb":      math.Logb,
	"mod":       mod,
	"nextafter": math.Nextafter,
	"pow":       math.Pow,
	"pow10":     math.Pow10,
	"remainder": math.Remainder,
	"sin":       math.Sin,
	"sinh":      math.Sinh,
	"sqrt":      math.Sqrt,
	"tan":       math.Tan,
	"tanh":      math.Tanh,
	"trunc":     math.Trunc,
	"y0":        math.Y0,
	"y1":        math.Y1,
	"yn":        math.Yn,
}

// -----------------------------------------------------------------------------

func mod(a, b interface{}) interface{} {

	return math.Mod(castFloat(a), castFloat(b))
}

func castFloat(a interface{}) float64 {

	switch a1 := a.(type) {
	case int:
		return float64(a1)
	case float64:
		return a1
	}
	panicUnsupportedFn("float", a)
	return 0
}

func panicUnsupportedFn(fn string, args ...interface{}) interface{} {

	targs := make([]string, len(args))
	for i, a := range args {
		targs[i] = reflect.TypeOf(a).String()
	}
	panic("unsupported function: " + fn + "(" + strings.Join(targs, ",") + ")")
}

// -----------------------------------------------------------------------------
