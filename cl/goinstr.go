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
	"log"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/exec.spec"
	"github.com/goplus/gop/reflect"
)

// -----------------------------------------------------------------------------

type goInstr = func(ctx *blockCtx, v *ast.CallExpr, ct callType) func()

type goInstrInfo struct {
	instr goInstr
}

var goinstrs map[string]goInstrInfo

func init() {
	goinstrs = map[string]goInstrInfo{
		"append":  {igoAppend},
		"copy":    {igoCopy},
		"delete":  {igoDelete},
		"len":     {igoLen},
		"cap":     {igoCap},
		"make":    {igoMake},
		"new":     {igoNew},
		"complex": {igoComplex},
		"real":    {igoReal},
		"imag":    {igoImag},
		"close":   {igoClose},
		"recover": {igoRecover},
	}
}

// -----------------------------------------------------------------------------

/*
func checkSliceType(args []interface{}, i int) (typ reflect.Type, ok bool) {
	if i >= len(args) {
		return
	}
	typ = args[i].(iValue).Type()
	ok = typ.Kind() == reflect.Slice
	return
}
*/

// func append(slice []Type, elems ...Type) []Type
func igoAppend(ctx *blockCtx, v *ast.CallExpr, ct callType) func() {
	if ct != callExpr {
		log.Panicf("%s discards result of %s\n", gCallTypes[ct], ctx.code(v))
	}
	if len(v.Args) < 2 {
		log.Panicln("append: argument count not enough")
	}
	sliceExpr := compileExpr(ctx, v.Args[0])
	sliceTy := ctx.infer.Get(-1).(iValue).Type()
	if sliceTy.Kind() != reflect.Slice {
		log.Panicln("append: first argument not a slice")
	}
	return func() {
		sliceExpr()
		n1 := len(v.Args) - 1
		for i := 1; i <= n1; i++ {
			compileExpr(ctx, v.Args[i])()
		}
		args := ctx.infer.GetArgs(n1)
		elem := sliceTy.Elem()
		if isEllipsis(v) {
			if n1 != 1 {
				log.Panicln("append: please use `append(slice, elems...)`")
			}
			checkType(sliceTy, args[0], ctx.out)
			ctx.infer.PopN(1)
			ctx.out.Append(elem, -1)
		} else {
			checkElementType(elem, args, 0, n1, 1, ctx.out)
			ctx.infer.PopN(n1)
			ctx.out.Append(elem, n1+1)
		}
	}
}

// func copy(dst, src []Type) int
func igoCopy(ctx *blockCtx, v *ast.CallExpr, ct callType) func() {
	if len(v.Args) < 2 {
		log.Panicln("not enough arguments in call to copy")
	}
	if len(v.Args) > 2 {
		log.Panicln("too many arguments in call to copy")
	}
	dstExpr := compileExpr(ctx, v.Args[0])
	dstTy := ctx.infer.Get(-1).(iValue).Type()
	if dstTy.Kind() != reflect.Slice {
		log.Panicln("arguments to copy must be slices; have ", dstTy.Kind())
	}
	if ct == callExpr {
		ctx.infer.Ret(1, &goValue{exec.TyInt})
	}
	return func() {
		dstExpr()
		compileExpr(ctx, v.Args[1])()
		srcTy := ctx.infer.Get(-1).(iValue).Type()
		switch srcTy.Kind() {
		case reflect.Slice:
			if srcTy.Elem().Kind() != dstTy.Elem().Kind() {
				log.Panicf("arguments to copy have different element types: %s(%s) %s(%s)", dstTy.Kind(), dstTy.Elem().Kind(), srcTy.Kind(), srcTy.Elem().Kind())
			}
		case reflect.String:
			if dstTy.Elem().Kind() != reflect.Uint8 {
				log.Panicln("arguments to copy have different element types:", dstTy.Kind(), srcTy.Kind())
			}
		default:
			log.Panicln("second argument to copy should be slice or string; have", srcTy.Kind())
		}
		ctx.infer.Pop()
		builder(ctx, ct).GoBuiltin(dstTy, exec.GobCopy)
	}
}

// func delete(m map[Type]Type1, key Type)
func igoDelete(ctx *blockCtx, v *ast.CallExpr, ct callType) func() {
	narg := len(v.Args)
	if narg != 2 {
		if narg < 2 {
			log.Panicln("missing second (key) argument to delete")
		}
		log.Panicln("too many arguments to delete")
	}
	mapExpr := compileExpr(ctx, v.Args[0])
	mapType := ctx.infer.Pop().(iValue).Type()
	if mapType.Kind() != reflect.Map {
		log.Panicln("first argument to delete must be map; have", mapType)
	}
	ctx.infer.Push(&nonValue{})
	return func() {
		mapExpr()
		compileExpr(ctx, v.Args[1])()
		key := ctx.infer.Pop().(iValue)
		checkType(mapType.Key(), key, ctx.out)
		builder(ctx, ct).GoBuiltin(mapType, exec.GobDelete)
	}
}

// func len/cap(v Type) int
func igoLenOrCap(ctx *blockCtx, v *ast.CallExpr, op exec.GoBuiltin) func() {
	narg := len(v.Args)
	if narg != 1 {
		if narg < 1 {
			logPanic(ctx, v, `missing argument to %v: %v`, op, ctx.code(v))
		}
		logPanic(ctx, v, `too many arguments to %v: %v`, op, ctx.code(v))
	}
	expr := compileExpr(ctx, v.Args[0])
	x := ctx.infer.Get(-1)
	typ := x.(iValue).Type()
	kind := typ.Kind()
	switch kind {
	case reflect.Array:
		n := typ.Len()
		ctx.infer.Ret(1, newConstVal(n, exec.Int))
		return func() {
			ctx.out.Push(n)
		}
	default:
		if kind == reflect.String {
			if op == exec.GobCap {
				logPanic(ctx, v, `invalid argument a (type %v) for cap`, typ)
			}
			if cons, ok := x.(*constVal); ok {
				n := len(cons.v.(string))
				ctx.infer.Ret(1, newConstVal(n, exec.Int))
				return func() {
					ctx.out.Push(n)
				}
			}
		}
		ctx.infer.Ret(1, &goValue{t: exec.TyInt})
	}
	return func() {
		switch kind {
		case reflect.Slice, reflect.String, reflect.Map, reflect.Chan:
			if op == exec.GobCap && kind == reflect.Chan {
				logPanic(ctx, v, `invalid argument a (type %v) for cap`, typ)
			}
			expr()
			ctx.out.GoBuiltin(typ, op)
		default:
			logPanic(ctx, v, `invalid argument a (type %v) for %v`, typ, op)
		}
	}
}

// func len(v Type) int
func igoLen(ctx *blockCtx, v *ast.CallExpr, ct callType) func() {
	if ct != callExpr {
		log.Panicf("%s discards result of %s\n", gCallTypes[ct], ctx.code(v))
	}
	return igoLenOrCap(ctx, v, exec.GobLen)
}

// func cap(v Type) int
func igoCap(ctx *blockCtx, v *ast.CallExpr, ct callType) func() {
	if ct != callExpr {
		log.Panicf("%s discards result of %s\n", gCallTypes[ct], ctx.code(v))
	}
	return igoLenOrCap(ctx, v, exec.GobCap)
}

// func make(t Type, size ...IntegerType) Type
func igoMake(ctx *blockCtx, v *ast.CallExpr, ct callType) func() {
	if ct != callExpr {
		log.Panicf("%s discards result of %s\n", gCallTypes[ct], ctx.code(v))
	}
	if len(v.Args) < 1 {
		logPanic(ctx, v, `missing argument to make: %v`, ctx.code(v))
	}
	t := toType(ctx, v.Args[0])
	kind := t.Kind()
	switch kind {
	case reflect.Slice, reflect.Map, reflect.Chan:
	default:
		log.Panicln("make: invalid type, must be a slice, map or chan -", t)
	}
	typ := t.(reflect.Type)
	ctx.infer.Push(&goValue{t: typ})
	return func() {
		n1 := len(v.Args) - 1
		for i := 1; i <= n1; i++ {
			compileExpr(ctx, v.Args[i])()
			checkIntType(ctx.infer.Pop(), ctx.out)
		}
		if kind == reflect.Slice && n1 == 0 {
			log.Panicln("make slice: must specified length of slice")
		}
		ctx.out.Make(typ, n1)
	}
}

// func new(Type) *Type
func igoNew(ctx *blockCtx, v *ast.CallExpr, ct callType) func() {
	if ct != callExpr {
		log.Panicf("%s discards result of %s\n", gCallTypes[ct], ctx.code(v))
	}
	if n := len(v.Args); n != 1 {
		if n == 0 {
			log.Panicln("missing argument to new")
		}
		log.Panicf("too many arguments to new(%v)\n", ctx.code(v.Args[0]))
	}
	t := toType(ctx, v.Args[0]).(reflect.Type)
	ptrT := reflect.PtrTo(t)
	ctx.infer.Push(&goValue{t: ptrT})
	return func() {
		ctx.out.New(t)
	}
}

// func complex(r, i FloatType) ComplexType
func igoComplex(ctx *blockCtx, v *ast.CallExpr, ct callType) func() {
	if ct != callExpr {
		log.Panicf("%s discards result of %s\n", gCallTypes[ct], ctx.code(v))
	}
	panic("todo")
}

// func real(c ComplexType) FloatType
func igoReal(ctx *blockCtx, v *ast.CallExpr, ct callType) func() {
	if ct != callExpr {
		log.Panicf("%s discards result of %s\n", gCallTypes[ct], ctx.code(v))
	}
	panic("todo")
}

// func imag(c ComplexType) FloatType
func igoImag(ctx *blockCtx, v *ast.CallExpr, ct callType) func() {
	if ct != callExpr {
		log.Panicf("%s discards result of %s\n", gCallTypes[ct], ctx.code(v))
	}
	panic("todo")
}

// func close(c chan<- Type)
func igoClose(ctx *blockCtx, v *ast.CallExpr, ct callType) func() {
	panic("todo")
}

// func recover() interface{}
func igoRecover(ctx *blockCtx, v *ast.CallExpr, ct callType) func() {
	panic("todo")
}

func compileTypeCast(typ reflect.Type, ctx *blockCtx, v *ast.CallExpr) func() {
	if len(v.Args) != 1 {
		log.Panicln("compileTypeCast: invalid argument count, please use `type(expr)`")
	}
	xExpr := compileExpr(ctx, v.Args[0])
	in := ctx.infer.Get(-1)
	kind := typ.Kind()
	if kind <= reflect.Complex128 || kind == reflect.String { // can be constant
		if cons, ok := in.(*constVal); ok {
			cons.kind = typ.Kind()

			if reflect.IsUserType(typ) {
				ctx.infer.Ret(1, &goValue{typ})
				return func() {
					pushConstVal(ctx.out, cons)
					ctx.out.TypeCast(cons.Type(), typ)
				}
			}
			return func() {
				pushConstVal(ctx.out, cons)
			}
		}
	}
	ctx.infer.Ret(1, &goValue{typ})
	return func() {
		xExpr()
		iv := in.(iValue)
		n := iv.NumValues()
		if n != 1 {
			panicExprNotValue(n)
		}
		tIn := iv.Type()
		if !reflect.ConvertibleTo(tIn, typ) {
			log.Panicf("compileTypeCast: can't convert type `%v` to `%v`\n", tIn, typ)
		}
		ctx.out.TypeCast(tIn, typ)
	}
}

// -----------------------------------------------------------------------------
