package exec

import (
	"errors"
	"fmt"
	"reflect"

	qlang "qlang.io/spec"
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
	typeFloat64  = reflect.TypeOf(float64(0))
	typeIntf     = reflect.TypeOf((*interface{})(nil)).Elem()
	typeFunction = reflect.TypeOf((*Function)(nil))
)

// -----------------------------------------------------------------------------
// Call

type iCall struct {
	vfn   reflect.Value
	n     int // 期望的参数个数
	arity int // 实际传入的参数个数
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

func function2Func(in *reflect.Value, t reflect.Type) {

	fn := in.MethodByName("Call")
	wrap := func(args []reflect.Value) (results []reflect.Value) {
		argsWithStk := make([]reflect.Value, len(args)+1)
		stk := NewStack() // issue #127: 使用新的stack，因为回调过程有可能是从新的goroutine发起的
		argsWithStk[0] = reflect.ValueOf(stk)
		copy(argsWithStk[1:], args)
		out := fn.Call(argsWithStk)
		n := t.NumOut()
		if n == 0 {
			return
		}
		ret := out[0]
		if n == 1 {
			return []reflect.Value{toType(ret, t.Out(0))}
		}
		if ret.Kind() != reflect.Slice || ret.Len() != n {
			panic(fmt.Sprintf("unexpected return value count, we need `%d` values", n))
		}
		results = make([]reflect.Value, n)
		for i := 0; i < n; i++ {
			results[i] = toType(ret.Index(i), t.Out(i))
		}
		return
	}
	*in = reflect.MakeFunc(t, wrap)
}

func toType(v reflect.Value, t reflect.Type) reflect.Value {

	if t.Kind() != reflect.Interface {
		return v.Elem()
	}
	if v.IsNil() {
		return reflect.New(t).Elem()
	}
	return v.Elem().Convert(t)
}

func validateType(in *reflect.Value, t, tfn reflect.Type) {

	kind := in.Kind()
	if kind == reflect.Invalid {
		*in = reflect.Zero(t) // work around `reflect: Call using zero Value argument`
		return
	}

	tkind := t.Kind()
	if tkind == reflect.Interface {
		if tfn != nil && qlang.DontTyNormalize[tfn] { // don't normalize input type
			return
		}
		switch {
		case kind > reflect.Int && kind <= reflect.Int64:
			*in = reflect.ValueOf(int(in.Int()))
		case kind >= reflect.Uint && kind <= reflect.Uintptr:
			*in = reflect.ValueOf(int(in.Uint()))
		case kind == reflect.Float32:
			*in = reflect.ValueOf(in.Float())
		}
		return
	}

	tin := in.Type()
	if tin == t {
		return
	}

	switch tkind {
	case reflect.Struct:
		if kind == reflect.Ptr {
			tin = tin.Elem()
			if tin == t {
				*in = in.Elem()
				return
			}
		}
	case reflect.Func:
		if tin == typeFunction {
			function2Func(in, t)
			return
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
type iCaller interface {
	Call(stk *Stack, args ...interface{}) interface{}
}

func (arity iCallFn) Exec(stk *Stack, ctx *Context) {

	in, ok := stk.PopArgs(int(arity) + 1)
	if !ok {
		panic(ErrStackDamaged)
	}

	vfn := in[0]
	tfn := vfn.Type()
	var tfn0 reflect.Type
	if vfn.Kind() != reflect.Func { // 这不是func，而是Function/其他可调用对象
		tfn0 = tfn
		if caller, ok := vfn.Interface().(iCaller); ok {
			vfn = reflect.ValueOf(func(args ...interface{}) interface{} {
				return caller.Call(stk, args...)
			})
		} else {
			vfnt := vfn.MethodByName("Call")
			if vfnt.IsValid() {
				vfn = vfnt
			} else {
				vfn = reflect.Indirect(vfn).FieldByName("Call")
				if vfn.Kind() == reflect.Interface {
					vfn = vfn.Elem()
				}
			}
		}
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
