package number

import (
	"errors"
	"fmt"
	"math"
	"reflect"

	"qiniupkg.com/text/tpl.v1/interpreter.util"
)

// -----------------------------------------------------------------------------
//go:generate tplgen -i -g Grammar -f static_interpreter.go
const Grammar = `

term = factor *('*' factor/mul | '/' factor/quo | '%' factor/mod)

doc = term *('+' term/add | '-' term/sub)

factor =
	FLOAT/push |
	'-' factor/neg |
	'(' doc ')' |
	(IDENT '(' doc %= ','/ARITY ')')/call |
	IDENT/ident |
	'+' factor
`

var (
	ErrUnsupportedRetType = errors.New("unsupported return type of function")
	ErrFncallWithoutArity = errors.New("function call without arity")
)

// -----------------------------------------------------------------------------

func Neg(a float64) float64 {
	return -a
}

func Mul(a, b float64) float64 {
	return a * b
}

func Quo(a, b float64) float64 {
	return a / b
}

func Add(a, b float64) float64 {
	return a + b
}

func Sub(a, b float64) float64 {
	return a - b
}

func Inf(a float64) float64 {
	v := 0
	if a < 0 {
		v = -1
	}
	return math.Inf(v)
}

func Jn(a, b float64) float64 {
	return math.Jn(int(a), b)
}

func Yn(a, b float64) float64 {
	return math.Yn(int(a), b)
}

func Ldexp(a, b float64) float64 {
	return math.Ldexp(a, int(b))
}

func Pow10(a float64) float64 {
	return math.Pow10(int(a))
}

func Max(args ...float64) (max float64) {

	if len(args) == 0 {
		return
	}

	max = args[0]
	for i := 1; i < len(args); i++ {
		if args[i] > max {
			max = args[i]
		}
	}
	return
}

func Min(args ...float64) (min float64) {

	if len(args) == 0 {
		return
	}

	min = args[0]
	for i := 1; i < len(args); i++ {
		if args[i] < min {
			min = args[i]
		}
	}
	return
}

// -----------------------------------------------------------------------------

type Stack []float64

func NewStack() *Stack {

	stk := make(Stack, 0, 16)
	return &stk
}

func (stk *Stack) Clear() {

	*stk = (*stk)[:0]
}

func (stk *Stack) Push(v float64) {

	*stk = append(*stk, v)
}

func (stk *Stack) PushIdent(v string) {

	var val float64
	switch v {
	case "e":
		val = math.E
	case "pi":
		val = math.Pi
	case "phi":
		val = math.Phi
	case "Inf":
		val = math.Inf(1)
	case "NaN":
		val = math.NaN()
	default:
		panic("undefined symbol: " + v)
	}
	*stk = append(*stk, val)
}

func (stk *Stack) Pop() (v float64, ok bool) {

	n := len(*stk)
	if n > 0 {
		v, ok = (*stk)[n-1], true
		*stk = (*stk)[:n-1]
	}
	return
}

func (stk *Stack) PushRet(ret []reflect.Value) error {

	for _, v := range ret {
		var val float64
		switch kind := v.Kind(); {
		case kind == reflect.Float64 || kind == reflect.Float32:
			val = v.Float()
		case kind >= reflect.Int && kind <= reflect.Int64:
			val = float64(v.Int())
		case kind >= reflect.Uint && kind <= reflect.Uintptr:
			val = float64(v.Uint())
		default:
			return ErrUnsupportedRetType
		}
		stk.Push(val)
	}
	return nil
}

func (stk *Stack) PopArgs(arity int) (args []reflect.Value, ok bool) {

	pstk := *stk
	n := len(pstk)
	if n >= arity {
		args, ok = make([]reflect.Value, arity), true
		n -= arity
		for i := 0; i < arity; i++ {
			args[i] = reflect.ValueOf(pstk[n+i])
		}
		*stk = pstk[:n]
	}
	return
}

func Arity(stk *Stack, arity int) {

	stk.Push(float64(arity))
}

// -----------------------------------------------------------------------------

type Calculator struct {
	stk *Stack
}

func New() *Calculator {

	stk := NewStack()
	return &Calculator{stk}
}

func (p *Calculator) Grammar() string {

	return Grammar
}

func (p *Calculator) Fntable() map[string]interface{} {

	return fntable
}

func (p *Calculator) Stack() interpreter.Stack {

	return p.stk
}

func (p *Calculator) Ret() (v float64, ok bool) {

	v, ok = p.stk.Pop()
	p.stk.Clear()
	return
}

func (p *Calculator) Call(name string) error {

	if fn, ok := fntable[name]; ok {
		if arity, ok := p.stk.Pop(); ok {
			err := interpreter.Call(p.stk, fn, int(arity))
			if err != nil {
				return fmt.Errorf("call function `%s` failed: %v", name, err)
			}
			return nil
		}
		return ErrFncallWithoutArity
	}
	return fmt.Errorf("function `%s` not found", name)
}

func init() {
	fntable = Fntable
}

var fntable map[string]interface{}

// -----------------------------------------------------------------------------

var Fntable = map[string]interface{}{
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
	"inf":       Inf,
	"j0":        math.J0,
	"j1":        math.J1,
	"jn":        Jn,
	"ldexp":     Ldexp,
	"ln":        math.Log,
	"log":       math.Log,
	"log10":     math.Log10,
	"log1p":     math.Log1p,
	"log2":      math.Log2,
	"logb":      math.Logb,
	"max":       Max,
	"min":       Min,
	"mod":       math.Mod,
	"NaN":       math.NaN,
	"nextafter": math.Nextafter,
	"pow":       math.Pow,
	"pow10":     Pow10,
	"remainder": math.Remainder,
	"sin":       math.Sin,
	"sinh":      math.Sinh,
	"sqrt":      math.Sqrt,
	"tan":       math.Tan,
	"tanh":      math.Tanh,
	"trunc":     math.Trunc,
	"y0":        math.Y0,
	"y1":        math.Y1,
	"yn":        Yn,

	"$neg":   Neg,
	"$mul":   Mul,
	"$quo":   Quo,
	"$mod":   math.Mod,
	"$add":   Add,
	"$sub":   Sub,
	"$ARITY": Arity,
	"$push":  (*Stack).Push,
	"$ident": (*Stack).PushIdent,
	"$call":  (*Calculator).Call,
}

// -----------------------------------------------------------------------------
