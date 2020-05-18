package cl

import (
	"errors"
	"reflect"

	"github.com/qiniu/qlang/v6/exec"
	"github.com/qiniu/x/log"
)

type iFuncType interface {
	In(i int) reflect.Type
	Out(i int) reflect.Type
	NumIn() int
	NumOut() int
	IsVariadic() bool
}

var (
	// ErrFuncArgNoReturnValue error.
	ErrFuncArgNoReturnValue = errors.New("function argument expression doesn't have return value")
	// ErrFuncArgCantBeMultiValue error.
	ErrFuncArgCantBeMultiValue = errors.New("function argument expression can't be multi values")
)

func checkFuncCall(tfn iFuncType, isMethod int, args []interface{}, b *exec.Builder) (arity int) {
	narg := tfn.NumIn() - isMethod
	variadic := tfn.IsVariadic()
	if variadic {
		narg--
	}
	if len(args) == 1 {
		n := args[0].(iValue).NumValues()
		if n != 1 { // TODO
			return n + isMethod
		}
	}
	for idx, arg := range args {
		var treq reflect.Type
		if variadic && idx >= narg {
			treq = tfn.In(narg + isMethod).Elem()
		} else {
			treq = tfn.In(idx + isMethod)
		}
		if v, ok := arg.(*constVal); ok {
			v.bound(treq, b)
		}
		n := arg.(iValue).NumValues()
		if n != 1 {
			if n == 0 {
				log.Panicln("checkFuncCall:", ErrFuncArgNoReturnValue)
			} else {
				log.Panicln("checkFuncCall:", ErrFuncArgCantBeMultiValue)
			}
		}
	}
	return len(args) + isMethod
}

func checkBinaryOp(kind exec.Kind, op exec.Operator, x, y interface{}, b *exec.Builder) {
	if xcons, xok := x.(*constVal); xok {
		if xcons.reserve != -1 {
			xv, ok := boundConst(xcons.v, exec.TypeFromKind(kind))
			if !ok {
				log.Panicln("checkBinaryOp: invalid operator", kind, "argument type.")
			}
			xcons.reserve.Push(b, xv)
		}
	}
	if ycons, yok := y.(*constVal); yok {
		i := op.GetInfo()
		if i.InSecond != (1 << exec.SameAsFirst) {
			if (uint64(ycons.kind) & i.InSecond) == 0 {
				log.Panicln("checkBinaryOp: invalid operator", kind, "argument type.")
			}
			kind = ycons.boundKind()
		}
		if ycons.reserve != -1 {
			yv, ok := boundConst(ycons.v, exec.TypeFromKind(kind))
			if !ok {
				log.Panicln("checkBinaryOp: invalid operator", kind, "argument type.")
			}
			ycons.reserve.Push(b, yv)
		}
	}
}

func checkType(t reflect.Type, v interface{}, b *exec.Builder) {
	if cons, ok := v.(*constVal); ok {
		cons.bound(t, b)
	} else {
		typVal := v.(iValue).Type()
		if kind := t.Kind(); kind == reflect.Interface {
			if !typVal.Implements(t) {
				log.Panicln("checkType: type doesn't implments interface -", typVal, t)
			}
		} else if t != typVal {
			log.Panicln("checkType: mismatched value type -", t, typVal)
		}
	}
}

func checkElementType(t reflect.Type, elts []interface{}, i, n, step int, b *exec.Builder) {
	for ; i < n; i += step {
		checkType(t, elts[i], b)
	}
}

func checkBool(v interface{}) {
	if !isBool(v.(iValue)) {
		log.Panicln("checkBool failed: bool expression required.")
	}
}

// -----------------------------------------------------------------------------
