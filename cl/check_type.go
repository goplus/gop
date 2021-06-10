/*
 Copyright 2020 The GoPlus Authors (goplus.org)

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package cl

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/exec.spec"
	"github.com/goplus/gop/token"
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

func checkFuncArgs(tfn iFuncType, args []interface{}, b exec.Builder) {
	for i, arg := range args {
		texp := tfn.In(i)
		checkType(texp, arg, b)
	}
}

func checkUnaryOp(kind exec.Kind, op exec.Operator, x interface{}, b exec.Builder) {
	if xcons, xok := x.(*constVal); xok {
		if xcons.reserve != -1 {
			xv := boundConst(xcons.v, exec.TypeFromKind(kind))
			xcons.reserve.Push(b, xv)
		}
	}
}

func checkBinaryOp(kind exec.Kind, op exec.Operator, x, y interface{}, b exec.Builder) {
	if xcons, xok := x.(*constVal); xok {
		if xcons.reserve != -1 {
			if xcons.kind == exec.ConstUnboundPtr {
				xcons.reserve.Push(b, nil)
			} else {
				if kind == reflect.Interface {
					xv := boundConst(xcons.v, xcons.boundType())
					xcons.reserve.Push(b, xv)
				} else {
					xv := boundConst(xcons.v, exec.TypeFromKind(kind))
					xcons.reserve.Push(b, xv)
				}
			}
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
			if ycons.kind == exec.ConstUnboundPtr {
				ycons.reserve.Push(b, nil)
			} else {
				if kind == reflect.Interface {
					yv := boundConst(ycons.v, ycons.boundType())
					ycons.reserve.Push(b, yv)
				} else {
					yv := boundConst(ycons.v, exec.TypeFromKind(kind))
					ycons.reserve.Push(b, yv)
				}
			}
		}
	}
}

func checkOpMatchType(op exec.Operator, x, y interface{}) error {
	switch op {
	case exec.OpEQ, exec.OpNE:
		ix := x.(iValue)
		iy := y.(iValue)
		xkind, ykind := ix.Kind(), iy.Kind()
		if xkind == reflect.Interface || ykind == reflect.Interface {
			return nil
		}
		if xkind == exec.ConstUnboundPtr && ykind == exec.ConstUnboundPtr {
			return nil
		} else if xkind == exec.ConstUnboundPtr {
			if !IsPtrKind(ykind) {
				return fmt.Errorf("mismatched types nil and %v", boundType(iy))
			}
		} else if ykind == exec.ConstUnboundPtr {
			if !IsPtrKind(xkind) {
				return fmt.Errorf("mismatched types %v and nil", boundType(ix))
			}
		} else if tx, ty := boundType(ix), boundType(iy); tx != ty {
			return fmt.Errorf("mismatched types %v and %v", tx, ty)
		}
	}
	return nil
}

func IsPtrKind(kind reflect.Kind) bool {
	return kind >= reflect.Chan && kind <= reflect.Slice || kind == reflect.UnsafePointer
}

func checkType(t reflect.Type, v interface{}, b exec.Builder) {
	if cons, ok := v.(*constVal); ok {
		cons.bound(t, b)
	} else if lsh, ok := v.(*lshValue); ok {
		lsh.checkType(t)
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
			switch typVal.Kind() {
			case reflect.Chan,
				reflect.Func,
				reflect.Slice,
				reflect.Array,
				reflect.Struct,
				reflect.Map:
				if !typVal.ConvertibleTo(t) {
					log.Panicf("checkType: cannot use `%v` as type `%v` in argument to produce", typVal, t)
				}
			default:
				log.Panicf("checkType: unexptected value type, require `%v(%p)`, but got `%v(%p)`\n", t, t, typVal, typVal)
			}
		}
	}
}

func checkIntType(v interface{}, b exec.Builder) {
	if cons, ok := v.(*constVal); ok {
		cons.bound(exec.TyInt, b)
	} else if lsh, ok := v.(*lshValue); ok {
		lsh.checkType(exec.TyInt)
	} else {
		iv := v.(iValue)
		n := iv.NumValues()
		if n != 1 {
			panicExprNotValue(n)
		}
		t := iv.Type()
		if kind := t.Kind(); kind >= reflect.Int && kind <= reflect.Uintptr {
			if kind != reflect.Int {
				b.TypeCast(t, exec.TyInt)
			}
		} else {
			log.Panicf("checkIntType: unexptected value type, require integer type, but got `%v`\n", t)
		}
	}
}

func checkElementType(t reflect.Type, elts []interface{}, i, n, step int, b exec.Builder) {
	for ; i < n; i += step {
		checkType(t, elts[i], b)
	}
}

func checkCaseCompare(x, y interface{}, b exec.Builder) {
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
