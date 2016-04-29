package exec

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrStackDamaged       = errors.New("unexpected: stack damaged")
	ErrArityRequired      = errors.New("arity required")
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
				validateType(&in[i], t)
			}
		}
		for i := 0; i < n; i++ {
			validateType(&in[i], tfn.In(i))
		}
	}
	out := p.vfn.Call(in)
	err := stk.PushRet(out)
	if err != nil {
		panic(err)
	}
}

func validateType(in *reflect.Value, t reflect.Type) {

	switch t.Kind() {
	case reflect.Interface:
		switch kind := in.Kind(); {
		case kind == reflect.Invalid:
			*in = reflect.Zero(t) // work around `reflect: Call using zero Value argument`
		case kind > reflect.Int && kind <= reflect.Int64:
			*in = reflect.ValueOf(int(in.Int()))
		case kind > reflect.Uint && kind <= reflect.Uintptr:
			*in = reflect.ValueOf(int(in.Uint()))
		case kind == reflect.Float32:
			*in = reflect.ValueOf(in.Float())
		}
		return
	case reflect.Ptr:
		if !in.IsValid() {
			*in = reflect.Zero(t) // work around `reflect: Call using zero Value argument`
		}
		return
	}

	tin := in.Type()
	if tin == t {
		return
	}

	if t == typeFloat64 {
		var val float64
		switch kind := tin.Kind(); {
		case kind >= reflect.Int && kind <= reflect.Int64:
			val = float64(in.Int())
		case kind >= reflect.Uint && kind <= reflect.Uintptr:
			val = float64(in.Uint())
		case kind == reflect.Float32:
			val = in.Float()
		default:
			goto lzErr
		}
		*in = reflect.ValueOf(val)
		return
	}

lzErr:
	panic(fmt.Errorf("invalid argument type: require `%v`, but we got `%v`", t, tin))
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

func (arity iCallFn) Exec(stk *Stack, ctx *Context) {

	in, ok := stk.PopArgs(int(arity) + 1)
	if !ok {
		panic(ErrStackDamaged)
	}

	vfn := in[0]
	if vfn.Kind() != reflect.Func { // 这不是func，而是Function对象
		vfn = vfn.MethodByName("Call")
	}
	tfn := vfn.Type()
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
			validateType(&in[i], t)
		}
	}
	for i := 0; i < n; i++ {
		validateType(&in[i], tfn.In(i))
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
