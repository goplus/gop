package interpreter

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrStackDamaged       = errors.New("unexpected: stack damaged")
	ErrArgumentsNotEnough = errors.New("argument not enough")
)

var (
	typeFloat64 = reflect.TypeOf(float64(0))
)

// -----------------------------------------------------------------------------

type FileLine struct {
	File string
	Line int
}

type Stack interface {
	PopArgs(arity int) (args []reflect.Value, ok bool)
	PushRet(ret []reflect.Value) error
}

type Interface interface {
	Grammar() string
	Fntable() map[string]interface{}
	Stack() Stack
}

type Engine interface {
	EvalCode(ip Interface, name string, src interface{}) error
	FileLine(src interface{}) FileLine
	Source(src interface{}) []byte
}

func Call(stk Stack, fn interface{}, arity int) error {

	tfn := reflect.TypeOf(fn)
	n := tfn.NumIn()

	isVariadic := tfn.IsVariadic() // 可变参数函数
	if isVariadic {
		if arity < n-1 {
			return ErrArgumentsNotEnough
		}
	} else if arity != n {
		return fmt.Errorf("invalid argument count: require %d, but we got %d", n, arity)
	}

	return DoCall(stk, fn, arity, n, isVariadic)
}

func DoCall(stk Stack, fn interface{}, arity, n int, isVariadic bool) error {

	tfn := reflect.TypeOf(fn)

	var ok bool
	var in []reflect.Value
	if arity > 0 {
		if in, ok = stk.PopArgs(arity); !ok {
			return ErrStackDamaged
		}
		if isVariadic {
			n--
			t := tfn.In(n).Elem()
			for i := n; i < arity; i++ {
				if err := validateType(&in[i], t); err != nil {
					return err
				}
			}
		}
		for i := 0; i < n; i++ {
			if err := validateType(&in[i], tfn.In(i)); err != nil {
				return err
			}
		}
	}
	vfn := reflect.ValueOf(fn)
	out := vfn.Call(in)
	return stk.PushRet(out)
}

func VCall(stk Stack, arity int) error {

	in, ok := stk.PopArgs(arity + 1)
	if !ok {
		return ErrStackDamaged
	}

	vfn := in[0]
	if vfn.Kind() != reflect.Func { // 这不是func，而是Function对象
		vfn = vfn.MethodByName("Call")
	}
	tfn := vfn.Type()
	n := tfn.NumIn()

	isVariadic := tfn.IsVariadic() // 可变参数函数
	if isVariadic {
		if arity < n-1 {
			return ErrArgumentsNotEnough
		}
	} else if arity != n {
		return fmt.Errorf("invalid argument count: require %d, but we got %d", n, arity)
	}

	in = in[1:]
	if isVariadic {
		n--
		t := tfn.In(n).Elem()
		for i := n; i < arity; i++ {
			if err := validateType(&in[i], t); err != nil {
				return err
			}
		}
	}
	for i := 0; i < n; i++ {
		if err := validateType(&in[i], tfn.In(i)); err != nil {
			return err
		}
	}

	out := vfn.Call(in)
	return stk.PushRet(out)
}

func validateType(in *reflect.Value, t reflect.Type) error {

	switch t.Kind() {
	case reflect.Interface, reflect.Ptr:
		if !in.IsValid() {
			*in = reflect.Zero(t) // work around `reflect: Call using zero Value argument`
		}
		return nil
	}

	tin := in.Type()
	if tin == t {
		return nil
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
		return nil
	}

lzErr:
	return fmt.Errorf("invalid argument type: require `%v`, but we got `%v`", t, tin)
}

// -----------------------------------------------------------------------------
