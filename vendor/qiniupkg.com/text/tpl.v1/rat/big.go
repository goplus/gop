package rat

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"

	"qiniupkg.com/text/tpl.v1/interpreter.util"
)

// -----------------------------------------------------------------------------

const Grammar = `

term = factor *('*' factor/mul | '/' factor/quo)

doc = term *('+' term/add | '-' term/sub)

factor =
	INT/push |
	'-' factor/neg |
	'(' doc ')' |
	(IDENT '(' doc %= ','/ARITY ')')/call |
	'+' factor
`

var (
	ErrUnsupportedRetType = errors.New("unsupported return type of function")
	ErrFncallWithoutArity = errors.New("function call without arity")
)

var (
	zero       = new(big.Rat)
	typeRatPtr = reflect.TypeOf(zero)
)

// -----------------------------------------------------------------------------

func Denom(a *big.Rat) *big.Rat {
	return new(big.Rat).SetInt(a.Denom())
}

func Num(a *big.Rat) *big.Rat {
	return new(big.Rat).SetInt(a.Num())
}

func Abs(a *big.Rat) *big.Rat {
	return new(big.Rat).Abs(a)
}

func Inv(a *big.Rat) *big.Rat {
	return new(big.Rat).Inv(a)
}

func Neg(a *big.Rat) *big.Rat {
	return new(big.Rat).Neg(a)
}

func Mul(a, b *big.Rat) *big.Rat {
	return new(big.Rat).Mul(a, b)
}

func Quo(a, b *big.Rat) *big.Rat {
	return new(big.Rat).Quo(a, b)
}

func Add(a, b *big.Rat) *big.Rat {
	return new(big.Rat).Add(a, b)
}

func Sub(a, b *big.Rat) *big.Rat {
	return new(big.Rat).Sub(a, b)
}

func Max(args ...*big.Rat) (max *big.Rat) {

	if len(args) == 0 {
		return zero
	}

	max = args[0]
	for i := 1; i < len(args); i++ {
		if args[i].Cmp(max) > 0 {
			max = args[i]
		}
	}
	return
}

func Min(args ...*big.Rat) (min *big.Rat) {

	if len(args) == 0 {
		return zero
	}

	min = args[0]
	for i := 1; i < len(args); i++ {
		if args[i].Cmp(min) < 0 {
			min = args[i]
		}
	}
	return
}

// -----------------------------------------------------------------------------

type Stack []*big.Rat

func NewStack() *Stack {

	stk := make(Stack, 0, 16)
	return &stk
}

func (stk *Stack) Clear() {

	*stk = (*stk)[:0]
}

func (stk *Stack) Push(v *big.Rat) {

	*stk = append(*stk, v)
}

func (stk *Stack) PushInt(v int64) {

	val := big.NewRat(v, 1)
	*stk = append(*stk, val)
}

func (stk *Stack) PushString(v string) { // 这个v可以很大，超出int64范围

	val, _ := new(big.Rat).SetString(v)
	*stk = append(*stk, val)
}

func (stk *Stack) Pop() (v *big.Rat, ok bool) {

	n := len(*stk)
	if n > 0 {
		v, ok = (*stk)[n-1], true
		*stk = (*stk)[:n-1]
	}
	return
}

func (stk *Stack) PushRet(ret []reflect.Value) error {

	for _, v := range ret {
		if v.Type() != typeRatPtr {
			return ErrUnsupportedRetType
		}
		val := v.Interface().(*big.Rat)
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

	stk.Push(big.NewRat(int64(arity), 1))
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

func (p *Calculator) Ret() (v *big.Rat, ok bool) {

	v, ok = p.stk.Pop()
	p.stk.Clear()
	return
}

func (p *Calculator) Call(name string) error {

	if fn, ok := fntable[name]; ok {
		if arity, ok := p.stk.Pop(); ok && arity.IsInt() {
			err := interpreter.Call(p.stk, fn, int(arity.Num().Int64()))
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
	"inv":   Inv,
	"max":   Max,
	"min":   Min,
	"denom": Denom,
	"num":   Num,

	"$neg":   Neg,
	"$mul":   Mul,
	"$quo":   Quo,
	"$add":   Add,
	"$sub":   Sub,
	"$ARITY": Arity,
	"$push":  (*Stack).PushString,
	"$call":  (*Calculator).Call,
}

// -----------------------------------------------------------------------------
