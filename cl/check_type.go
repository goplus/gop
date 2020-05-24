package cl

import (
	"errors"
	"reflect"

	"github.com/qiniu/qlang/v6/ast"
	"github.com/qiniu/qlang/v6/exec"
	"github.com/qiniu/qlang/v6/token"
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

func isEllipsis(v *ast.CallExpr) bool {
	return v.Ellipsis != token.NoPos
}

func checkFuncCall(tfn iFuncType, isMethod int, v *ast.CallExpr, ctx *blockCtx) (arity int) {
	nargIn := len(v.Args) + isMethod
	nargExp := tfn.NumIn()
	variadic := tfn.IsVariadic()
	args := ctx.infer.GetArgs(nargIn)
	if isEllipsis(v) {
		if !variadic {
			log.Panicln("checkFuncCall: call a non variadic function with ...")
		}
		if nargIn != nargExp {
			log.Panicln("checkFuncCall: call with unexpected argument count")
		}
		checkFuncArgs(tfn, args, ctx.out)
		arity = -1
	} else {
		if len(v.Args) == 1 {
			recv, in := args[0], args[isMethod].(iValue)
			n := in.NumValues()
			if n != 1 {
				if n == 0 {
					log.Panicln("checkFuncCall:", ErrFuncArgNoReturnValue)
				}
				args = make([]interface{}, 0, n+isMethod)
				if isMethod > 0 {
					args = append(args, recv)
				}
				for i := 0; i < n; i++ {
					args = append(args, in.Value(i))
				}
			}
		}
		if !variadic {
			if len(args) != nargExp {
				log.Panicln("checkFuncCall: call with unexpected argument count")
			}
			checkFuncArgs(tfn, args, ctx.out)
		} else {
			nargExp--
			if len(args) < nargExp {
				log.Panicln("checkFuncCall: argument count is not enough")
			}
			checkFuncArgs(tfn, args[:nargExp], ctx.out)
			nVariadic := len(args) - nargExp
			if nVariadic > 0 {
				checkElementType(tfn.In(nargExp).Elem(), args[nargExp:], 0, nVariadic, 1, ctx.out)
			}
		}
		arity = len(args)
	}
	ctx.infer.PopN(nargIn)
	return
}

func checkFuncArgs(tfn iFuncType, args []interface{}, b *exec.Builder) {
	for i, arg := range args {
		texp := tfn.In(i)
		checkType(texp, arg, b)
	}
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
		iv := v.(iValue)
		n := iv.NumValues()
		if n != 1 {
			panicExprNotValue(n)
		}
		typVal := iv.Type()
		if kind := t.Kind(); kind == reflect.Interface {
			if !typVal.Implements(t) {
				log.Panicf("checkType: type `%v` doesn't implments interface `%v`", typVal, t)
			}
		} else if t != typVal {
			log.Panicf("checkType: unexptected value type, require `%v`, but got `%v`\n", t, typVal)
		}
	}
}

func checkIntType(v interface{}, b *exec.Builder) {
	if cons, ok := v.(*constVal); ok {
		cons.bound(exec.TyInt, b)
	} else {
		iv := v.(iValue)
		n := iv.NumValues()
		if n != 1 {
			panicExprNotValue(n)
		}
		t := iv.Type()
		if kind := t.Kind(); kind >= reflect.Int && kind <= reflect.Uintptr {
			if kind != reflect.Int {
				b.TypeCast(kind, reflect.Int)
			}
		} else {
			log.Panicf("checkIntType: unexptected value type, require integer type, but got `%v`\n", t)
		}
	}
}

func checkElementType(t reflect.Type, elts []interface{}, i, n, step int, b *exec.Builder) {
	for ; i < n; i += step {
		checkType(t, elts[i], b)
	}
}

func checkCaseCompare(x, y interface{}, b *exec.Builder) {
	kind, _ := binaryOpResult(exec.OpEQ, x, y)
	checkBinaryOp(kind, exec.OpEQ, x, y, b)
}

func checkBool(v interface{}) {
	if !isBool(v.(iValue)) {
		log.Panicln("checkBool failed: bool expression required.")
	}
}

func panicExprNotValue(n int) {
	if n == 0 {
		log.Panicln("checkFuncCall:", ErrFuncArgNoReturnValue)
	} else {
		log.Panicln("checkFuncCall:", ErrFuncArgCantBeMultiValue)
	}
}

// -----------------------------------------------------------------------------
