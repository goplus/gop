package exec

import (
	"errors"
	"fmt"
	"reflect"

	"qlang.io/qlang.spec.v1"
)

var (
	// ErrStackDamaged is returned when stack is damaged.
	ErrStackDamaged = errors.New("unexpected: stack damaged")

	// ErrArityRequired is returned when calling `Call` without providing `arity`.
	ErrArityRequired = errors.New("arity required")

	// ErrArgumentsNotEnough is returned when calling a function without enough arguments.
	ErrArgumentsNotEnough = errors.New("arguments not enough")
)

var (
	typeFloat64 = reflect.TypeOf(float64(0))
	typeIntf    = reflect.TypeOf((*interface{})(nil)).Elem()
)

// -----------------------------------------------------------------------------
// Call

type iCall struct {
	vfn   reflect.Value
	n     int // 期望的参数个数
	arity int // 实际传入的参数个数
}

func (p *iCall) OptimizableGetArity() int {

	return p.arity
}

func (p *iCall) Exec(stk *Stack, ctx *Context) {

	tfn := p.vfn.Type()
	n, arity := p.n, p.arity

	var ok bool
	var in []reflect.Value
	if arity > 0 {
		if in, ok = stk.PopArgs(arity); !ok {
			panic(ErrStackDamaged)
		}
		if tfn.IsVariadic() {
			n--
			t := tfn.In(n).Elem()
			for i := n; i < arity; i++ {
				validateType(&in[i], t, nil)
			}
		}
		for i := 0; i < n; i++ {
			validateType(&in[i], tfn.In(i), nil)
		}
	}
	out := p.vfn.Call(in)
	err := stk.PushRet(out)
	if err != nil {
		panic(err)
	}
}

func validateType(in *reflect.Value, t, tfn reflect.Type) {

	tkind := t.Kind()
	switch tkind {
	case reflect.Interface:
		if tfn != nil && qlang.DontTyNormalize[tfn] { // don't normalize input type
			return
		}
		switch kind := in.Kind(); {
		case kind == reflect.Invalid:
			*in = reflect.Zero(t) // work around `reflect: Call using zero Value argument`
		case kind > reflect.Int && kind <= reflect.Int64:
			*in = reflect.ValueOf(int(in.Int()))
		case kind >= reflect.Uint && kind <= reflect.Uintptr:
			*in = reflect.ValueOf(int(in.Uint()))
		case kind == reflect.Float32:
			*in = reflect.ValueOf(in.Float())
		}
		return
	case reflect.Ptr, reflect.Slice, reflect.Map:
		if !in.IsValid() {
			*in = reflect.Zero(t) // work around `reflect: Call using zero Value argument`
			return
		}
	}

	tin := in.Type()
	if tin == t {
		return
	}

	kind := in.Kind()
	switch tkind {
	case reflect.Struct:
		if kind == reflect.Ptr {
			tin = tin.Elem()
			if tin == t {
				*in = in.Elem()
				return
			}
		}
	default:
		if tkind == kind || convertible(kind, tkind) {
			*in = in.Convert(t)
			return
		}
	}
	panic(fmt.Errorf("invalid argument type: require `%v`, but we got `%v`", t, tin))
}

func convertible(kind, tkind reflect.Kind) bool {

	if tkind >= reflect.Int && tkind <= reflect.Uintptr {
		return kind >= reflect.Int && kind <= reflect.Uintptr
	}
	if tkind == reflect.Float64 || tkind == reflect.Float32 {
		return kind >= reflect.Int && kind <= reflect.Float64
	}
	return false
}

// Call returns a function call instruction.
//
func Call(fn interface{}, varity ...int) Instr {

	tfn := reflect.TypeOf(fn)
	n := tfn.NumIn()
	arity := 0
	if len(varity) == 0 {
		arity = n
	} else {
		arity = varity[0]
	}

	isVariadic := tfn.IsVariadic() // 可变参数函数
	if isVariadic {
		if len(varity) == 0 {
			panic(ErrArityRequired)
		}
		if arity < n-1 {
			panic(ErrArgumentsNotEnough)
		}
	} else if arity != n {
		panic(fmt.Errorf("invalid argument count: require %d, but we got %d", n, arity))
	}

	return &iCall{reflect.ValueOf(fn), n, arity}
}

// -----------------------------------------------------------------------------
// CallFn

type iCallFn int

func (arity iCallFn) OptimizableGetArity() int {

	return int(arity) + 1
}

func (arity iCallFn) Exec(stk *Stack, ctx *Context) {

	in, ok := stk.PopArgs(int(arity) + 1)
	if !ok {
		panic(ErrStackDamaged)
	}

	vfn := in[0]
	tfn := vfn.Type()
	var tfn0 reflect.Type
	if vfn.Kind() != reflect.Func { // 这不是func，而是Function对象
		tfn0 = tfn
		vfn = vfn.MethodByName("Call")
		tfn = vfn.Type()
	}
	n := tfn.NumIn()

	isVariadic := tfn.IsVariadic() // 可变参数函数
	if isVariadic {
		if int(arity) < n-1 {
			panic(ErrArgumentsNotEnough)
		}
	} else if int(arity) != n {
		panic(fmt.Errorf("invalid argument count: require %d, but we got %d", n, arity))
	}

	in = in[1:]
	if isVariadic {
		n--
		t := tfn.In(n).Elem()
		for i := n; i < int(arity); i++ {
			validateType(&in[i], t, tfn0)
		}
	}
	for i := 0; i < n; i++ {
		validateType(&in[i], tfn.In(i), tfn0)
	}

	out := vfn.Call(in)
	err := stk.PushRet(out)
	if err != nil {
		panic(err)
	}
}

// CallFn returns a function call instruction.
//
func CallFn(arity int) Instr {
	return iCallFn(arity)
}

// -----------------------------------------------------------------------------
// CallFnv

type iCallFnv int

func (arity iCallFnv) OptimizableGetArity() int {

	return int(arity) + 1
}

func (arity iCallFnv) Exec(stk *Stack, ctx *Context) {

	instr := iCallFn(arity)
	if val, ok := stk.Pop(); ok {
		v := reflect.ValueOf(val)
		if v.Kind() != reflect.Slice {
			panic("apply `...` on non-slice object")
		}
		n := v.Len()
		for i := 0; i < n; i++ {
			stk.Push(v.Index(i).Interface())
		}
		instr += iCallFn(n - 1)
	} else {
		panic("unexpected")
	}

	instr.Exec(stk, ctx)
}

// CallFnv returns a function call instruction.
//
func CallFnv(arity int) Instr {
	return iCallFnv(arity)
}

// -----------------------------------------------------------------------------

type sliceFrom int

var nilVarSlice = make([]interface{}, 0)

func (p sliceFrom) OptimizableGetArity() int {

	return int(p)
}

func (p sliceFrom) Exec(stk *Stack, ctx *Context) {

	if p == 0 {
		stk.data = append(stk.data, nilVarSlice)
		return
	}
	n := len(stk.data) - int(p)
	stk.data[n] = qlang.SliceFrom(stk.data[n:]...)
	stk.data = stk.data[:n+1]
}

// SliceFrom is an instruction that creates slice in [a1, a2, ...] form.
//
func SliceFrom(arity int) Instr {

	return sliceFrom(arity)
}

// -----------------------------------------------------------------------------

type slice int

func (p slice) OptimizableGetArity() int {

	return 1
}

func (p slice) Exec(stk *Stack, ctx *Context) {

	n := len(stk.data) - 1
	stk.data[n] = qlang.Slice(stk.data[n])
}

// Slice is an instruction that returns []T.
//
var Slice Instr = slice(0)

// -----------------------------------------------------------------------------
