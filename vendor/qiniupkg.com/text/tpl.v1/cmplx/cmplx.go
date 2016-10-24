package cmplx

import (
	"errors"
	"fmt"
	"math/cmplx"
	"reflect"
	"strconv"

	"qiniupkg.com/text/tpl.v1/interpreter.util"
)

// -----------------------------------------------------------------------------

const Grammar = `

term = factor *('*' factor/mul | '/' factor/quo)

doc = term *('+' term/add | '-' term/sub)

factor =
	FLOAT/push |
	IMAG/pushi |
	'-' factor/neg |
	'(' doc ')' |
	(IDENT '(' doc %= ','/ARITY ')')/call |
	'+' factor
`

var (
	ErrUnsupportedRetType = errors.New("unsupported return type of function")
	ErrFncallWithoutArity = errors.New("function call without arity")
)

// -----------------------------------------------------------------------------

func Abs(a complex128) complex128 {
	return complex(cmplx.Abs(a), 0)
}

func Neg(a complex128) complex128 {
	return -a
}

func Mul(a, b complex128) complex128 {
	return a * b
}

func Quo(a, b complex128) complex128 {
	return a / b
}

func Add(a, b complex128) complex128 {
	return a + b
}

func Sub(a, b complex128) complex128 {
	return a - b
}

func Phase(a complex128) complex128 {
	return complex(cmplx.Phase(a), 0)
}

func Rect(a, b complex128) complex128 {
	if imag(a) != 0 || imag(b) != 0 {
		panic("function `rect` argument type must be `float64`")
	}
	return cmplx.Rect(real(a), real(b))
}

// -----------------------------------------------------------------------------

func Real(a complex128) complex128 {
	return complex(real(a), 0)
}

func Imag(a complex128) complex128 {
	return complex(imag(a), 0)
}

func Complex(a, b complex128) complex128 {
	if imag(a) != 0 || imag(b) != 0 {
		panic("function `complex` argument type must be `float64`")
	}
	return complex(real(a), real(b))
}

// -----------------------------------------------------------------------------

type Stack []complex128

func NewStack() *Stack {

	stk := make(Stack, 0, 16)
	return &stk
}

func (stk *Stack) Clear() {

	*stk = (*stk)[:0]
}

func (stk *Stack) Push(v complex128) {

	*stk = append(*stk, v)
}

func (stk *Stack) PushReal(v float64) {

	*stk = append(*stk, complex(v, 0))
}

func (stk *Stack) PushImag(v string) {

	val, _ := strconv.ParseFloat(v[:len(v)-1], 64)
	*stk = append(*stk, complex(0, val))
}

func (stk *Stack) Pop() (v complex128, ok bool) {

	n := len(*stk)
	if n > 0 {
		v, ok = (*stk)[n-1], true
		*stk = (*stk)[:n-1]
	}
	return
}

func (stk *Stack) PushRet(ret []reflect.Value) error {

	for _, v := range ret {
		var val complex128
		switch kind := v.Kind(); {
		case kind == reflect.Complex128 || kind == reflect.Complex64:
			val = v.Complex()
		case kind == reflect.Float64 || kind == reflect.Float32:
			val = complex(v.Float(), 0)
		case kind >= reflect.Int && kind <= reflect.Int64:
			val = complex(float64(v.Int()), 0)
		case kind >= reflect.Uint && kind <= reflect.Uintptr:
			val = complex(float64(v.Uint()), 0)
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

	stk.Push(complex(float64(arity), 0))
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

func (p *Calculator) Ret() (v complex128, ok bool) {

	v, ok = p.stk.Pop()
	p.stk.Clear()
	return
}

func (p *Calculator) Call(name string) error {

	if fn, ok := fntable[name]; ok {
		if arity, ok := p.stk.Pop(); ok && imag(arity) == 0 {
			err := interpreter.Call(p.stk, fn, int(real(arity)))
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
	"abs":   Abs,
	"acos":  cmplx.Acos,
	"acosh": cmplx.Acosh,
	"asin":  cmplx.Asin,
	"asinh": cmplx.Asinh,
	"atan":  cmplx.Atan,
	"atanh": cmplx.Atanh,
	"conj":  cmplx.Conj,
	"cos":   cmplx.Cos,
	"cosh":  cmplx.Cosh,
	"cot":   cmplx.Cot,
	"exp":   cmplx.Exp,
	"inf":   cmplx.Inf,
	"log":   cmplx.Log,
	"log10": cmplx.Log10,
	"NaN":   cmplx.NaN,
	"phase": Phase,
	"pow":   cmplx.Pow,
	"rect":  Rect,
	"sin":   cmplx.Sin,
	"sinh":  cmplx.Sinh,
	"sqrt":  cmplx.Sqrt,
	"tan":   cmplx.Tan,
	"tanh":  cmplx.Tanh,

	"real":    Real,
	"imag":    Imag,
	"complex": Complex,

	"$neg":   Neg,
	"$mul":   Mul,
	"$quo":   Quo,
	"$add":   Add,
	"$sub":   Sub,
	"$ARITY": Arity,
	"$push":  (*Stack).PushReal,
	"$pushi": (*Stack).PushImag,
	"$call":  (*Calculator).Call,
}

// -----------------------------------------------------------------------------
