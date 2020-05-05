package cl

import (
	"errors"
	"reflect"

	"github.com/qiniu/qlang/exec"
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

func checkFuncCall(tfn iFuncType, args []interface{}, b *exec.Builder) (arity int) {
	narg := tfn.NumIn()
	variadic := tfn.IsVariadic()
	if variadic {
		narg--
	}
	if len(args) == 1 {
		n := args[0].(iValue).NumValues()
		if n != 1 { // TODO
			return n
		}
	}
	for idx, arg := range args {
		var treq reflect.Type
		if variadic && idx >= narg {
			treq = tfn.In(narg).Elem()
		} else {
			treq = tfn.In(idx)
		}
		if v, ok := arg.(*constVal); ok {
			v.bound(treq, b)
		}
		n := arg.(iValue).NumValues()
		if n != 1 {
			if n == 0 {
				log.Fatalln("checkFuncCall:", ErrFuncArgNoReturnValue)
			} else {
				log.Fatalln("checkFuncCall:", ErrFuncArgCantBeMultiValue)
			}
		}
	}
	return len(args)
}

func checkBinaryOp(kind exec.Kind, op exec.Operator, x, y interface{}, b *exec.Builder) {
	if xcons, xok := x.(*constVal); xok {
		if xcons.reserve != -1 {
			xv, ok := boundConst(xcons.v, exec.TypeFromKind(kind))
			if !ok {
				log.Fatalln("checkBinaryOp: invalid operator", kind, "argument type.")
			}
			xcons.reserve.Push(b, xv)
		}
	}
	if ycons, yok := y.(*constVal); yok {
		i := op.GetInfo()
		if i.InSecond != (1 << exec.SameAsFirst) {
			if (uint64(ycons.kind) & i.InSecond) == 0 {
				log.Fatalln("checkBinaryOp: invalid operator", kind, "argument type.")
			}
			kind = ycons.kind
		}
		if ycons.reserve != -1 {
			yv, ok := boundConst(ycons.v, exec.TypeFromKind(kind))
			if !ok {
				log.Fatalln("checkBinaryOp: invalid operator", kind, "argument type.")
			}
			ycons.reserve.Push(b, yv)
		}
	}
}

func checkType(t reflect.Type, v interface{}, b *exec.Builder) {
	if cons, ok := v.(*constVal); ok {
		cons.bound(t, b)
	}
}

// -----------------------------------------------------------------------------
