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
	"reflect"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/ast/spec"
	"github.com/goplus/gop/exec.spec"
)

// -----------------------------------------------------------------------------

type goInstr = func(ctx *blockCtx, v *ast.CallExpr, ct callType) func()

type goInstrInfo struct {
	instr goInstr
}

var goinstrs map[string]goInstrInfo

func init() {
	goinstrs = map[string]goInstrInfo{
		"append":          {igoAppend},
		"copy":            {igoCopy},
		"delete":          {igoDelete},
		"len":             {igoLen},
		"cap":             {igoCap},
		"make":            {igoMake},
		"new":             {igoNew},
		"complex":         {igoComplex},
		"real":            {igoReal},
		"imag":            {igoImag},
		"close":           {igoClose},
		"recover":         {igoRecover},
		"unsafe.Sizeof":   {igoSizeof},
		"unsafe.Alignof":  {igoAlignof},
		"unsafe.Offsetof": {igoOffsetof},
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
		ctx.infer.Ret(1, &goValue{t: exec.TyInt})
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
	if n := len(v.Args); n != 2 {
		if n < 2 {
			log.Panicln("not enough arguments in call to", ctx.code(v))
		} else {
			log.Panicln("too many arguments in call to", ctx.code(v))
		}
	}
	xExpr := compileExpr(ctx, v.Args[0])
	yExpr := compileExpr(ctx, v.Args[1])
	xin := ctx.infer.Get(-1).(iValue)
	yin := ctx.infer.Get(-2).(iValue)
	xkind := xin.Kind()
	ykind := yin.Kind()
	switch xkind {
	case spec.ConstUnboundInt, spec.ConstUnboundFloat, reflect.Float32, reflect.Float64:
	default:
		log.Panicf("invalid operation: %v (arguments have type %v, expected floating-point)\n", ctx.code(v), xkind)
	}
	switch ykind {
	case spec.ConstUnboundInt, spec.ConstUnboundFloat, reflect.Float32, reflect.Float64:
	default:
		log.Panicf("invalid operation: %v (arguments have type %v, expected floating-point)\n", ctx.code(v), ykind)
	}
	if (xkind == reflect.Float32 && ykind == reflect.Float64) ||
		(xkind == reflect.Float64 && ykind == reflect.Float32) {
		log.Panicf("invalid operation: %v) (mismatched types %v and %v)", ctx.code(v), ykind, xkind)
	}
	var elem reflect.Type
	var typ reflect.Type
	if xkind == reflect.Float32 || ykind == reflect.Float32 {
		elem = exec.TyFloat32
		typ = exec.TyComplex64
	} else {
		elem = exec.TyFloat64
		typ = exec.TyComplex128
	}
	ctx.infer.Ret(2, &goValue{t: typ})
	return func() {
		xExpr()
		yExpr()
		if c, ok := xin.(*constVal); ok {
			c.bound(elem, ctx.out)
		}
		if c, ok := yin.(*constVal); ok {
			c.bound(elem, ctx.out)
		}
		ctx.out.GoBuiltin(typ, exec.GobComplex)
	}
}

// func real/imag(ComplexType) func
func igoRealOrImag(ctx *blockCtx, v *ast.CallExpr, op exec.GoBuiltin) func() {
	if n := len(v.Args); n != 1 {
		if n == 0 {
			log.Panicf("missing argument to %v: %v", op, ctx.code(v))
		} else {
			log.Panicf("too many arguments to %v: %v", op, ctx.code(v))
		}
	}
	expr := compileExpr(ctx, v.Args[0])
	in := ctx.infer.Get(-1).(iValue)
	kind := in.Kind()
	switch kind {
	case spec.ConstUnboundInt, spec.ConstUnboundFloat, spec.ConstUnboundComplex, reflect.Complex64, reflect.Complex128:
	default:
		log.Panicf("invalid argument %v (type %v) for %v\n", ctx.code(v.Args[0]), kind, op)
	}
	var ctyp reflect.Type
	var typ reflect.Type
	if kind == reflect.Complex64 {
		ctyp = exec.TyComplex64
		typ = exec.TyFloat32
	} else {
		ctyp = exec.TyComplex128
		typ = exec.TyFloat64
	}
	ctx.infer.Ret(1, &goValue{t: typ})
	return func() {
		expr()
		if c, ok := in.(*constVal); ok {
			c.bound(ctyp, ctx.out)
		}
		ctx.out.GoBuiltin(typ, op)
	}
}

// func real(c ComplexType) FloatType
func igoReal(ctx *blockCtx, v *ast.CallExpr, ct callType) func() {
	if ct != callExpr {
		log.Panicf("%s discards result of %s\n", gCallTypes[ct], ctx.code(v))
	}
	return igoRealOrImag(ctx, v, exec.GobReal)
}

// func imag(c ComplexType) FloatType
func igoImag(ctx *blockCtx, v *ast.CallExpr, ct callType) func() {
	if ct != callExpr {
		log.Panicf("%s discards result of %s\n", gCallTypes[ct], ctx.code(v))
	}
	return igoRealOrImag(ctx, v, exec.GobImag)
}

// func close(c chan<- Type)
func igoClose(ctx *blockCtx, v *ast.CallExpr, ct callType) func() {
	if n := len(v.Args); n != 1 {
		if n == 0 {
			log.Panicf("missing argument to close: %v", ctx.code(v))
		} else {
			log.Panicf("too many arguments to close: %v", ctx.code(v))
		}
	}
	arg := v.Args[0]
	expr := compileExpr(ctx, arg)
	in := ctx.infer.Pop().(iValue)
	kind := in.Kind()
	if kind != reflect.Chan {
		log.Panicf("invalid operation: close(%v) (non-chan type %v)", ctx.code(arg), spec.KindName(kind))
	}
	typ := in.Type()
	if typ.ChanDir() == reflect.RecvDir {
		log.Panicf("invalid operation: close(%v) (cannot close receive-only channel)", ctx.code(arg))
	}
	if ct == callExpr {
		ctx.infer.Push(&nonValue{})
	}
	return func() {
		expr()
		builder(ctx, ct).GoBuiltin(typ, exec.GobClose)
	}
}

// func recover() interface{}
func igoRecover(ctx *blockCtx, v *ast.CallExpr, ct callType) func() {
	if n := len(v.Args); n != 0 {
		log.Panicf("too many arguments to recover: %v", ctx.code(v))
	}
	typ := exec.TyEmptyInterface
	ctx.infer.Push(&goValue{t: typ})
	return func() {
		builder(ctx, ct).GoBuiltin(typ, exec.GobRecover)
	}
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
			if typ.PkgPath() != "" {
				c := newConstVal(cons.v, cons.kind)
				c.typed = typ
				ctx.infer.Ret(1, &goValue{t: typ, c: c})
				return func() {
					pushConstVal(ctx.out, c)
				}
			}
			return func() {
				pushConstVal(ctx.out, cons)
			}
		}
		if lsh, ok := in.(*lshValue); ok {
			lsh.bound(typ)
		}
	} else if kind == reflect.Interface {
		if cons, ok := in.(*constVal); ok {
			if cons.kind == exec.ConstUnboundPtr {
				cons.kind = reflect.UnsafePointer
			} else {
				kind := cons.boundKind()
				cons.v = boundConst(cons, exec.TypeFromKind(kind))
				cons.kind = kind
			}
		}
	} else if kind == reflect.Ptr {
		if cons, ok := in.(*constVal); ok && cons.kind == exec.ConstUnboundPtr {
			cons.kind = reflect.UnsafePointer
		}
	}

	ctx.infer.Ret(1, &goValue{t: typ})
	return func() {
		xExpr()
		iv := in.(iValue)
		n := iv.NumValues()
		if n != 1 {
			panicExprNotValue(n)
		}
		tIn := iv.Type()
		tkind := tIn.Kind()
		kind := typ.Kind()
		var conv bool
		switch kind {
		case reflect.UnsafePointer:
			if tkind == reflect.Ptr || tkind == reflect.Uintptr {
				conv = true
			}
		case reflect.Ptr, reflect.Uintptr:
			if tkind == reflect.UnsafePointer {
				conv = true
			}
		}
		if !conv {
			conv = tIn.ConvertibleTo(typ)
		}
		if !conv {
			log.Panicf("compileTypeCast: can't convert type `%v` to `%v`\n", tIn, typ)
		}
		ctx.out.TypeCast(tIn, typ)
	}
}

// func unsafe.Sizeof(ArbitraryType) uintptr
func igoSizeof(ctx *blockCtx, v *ast.CallExpr, ct callType) func() {
	return igoUnsafe(ctx, v, ct, "Sizeof")
}

// func unsafe.Alignof(ArbitraryType) uintptr
func igoAlignof(ctx *blockCtx, v *ast.CallExpr, ct callType) func() {
	return igoUnsafe(ctx, v, ct, "Alignof")
}

// func unsafe.Offsetof(ArbitraryType) uintptr
func igoOffsetof(ctx *blockCtx, v *ast.CallExpr, ct callType) func() {
	return igoUnsafe(ctx, v, ct, "Offsetof")
}

func igoUnsafe(ctx *blockCtx, v *ast.CallExpr, ct callType, name string) func() {
	if len(v.Args) == 0 {
		log.Panicf("missing argument to unsafe.%v: %v", name, ctx.code(v))
	} else if len(v.Args) > 1 {
		log.Panicf("too many arguments to unsafe.%v: %v", name, ctx.code(v))
	}
	ctx.infer.Pop()
	for _, arg := range v.Args {
		compileExpr(ctx, arg)
	}
	switch name {
	case "Sizeof":
		c := ctx.infer.Get(-1).(iValue)
		t := boundType(c)
		ctx.infer.Pop()
		ret := newConstVal(int64(t.Size()), exec.Uintptr)
		ctx.infer.Push(ret)
		return func() {
			pushConstVal(ctx.out, ret)
		}
	case "Alignof":
		c := ctx.infer.Get(-1).(iValue)
		t := boundType(c)
		ctx.infer.Pop()
		ret := newConstVal(int64(t.Align()), exec.Uintptr)
		ctx.infer.Push(ret)
		return func() {
			pushConstVal(ctx.out, ret)
		}
	case "Offsetof":
		if ctx.fieldStructType == nil || ctx.fieldIndex == nil {
			log.Panicf("invalid expression %v", ctx.code(v))
		}
		ctx.infer.Pop()
		typ := ctx.fieldStructType
		for typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		offset := typ.FieldByIndex(ctx.fieldIndex).Offset
		ret := newConstVal(int64(offset), exec.Uintptr)
		ctx.infer.Push(ret)
		return func() {
			pushConstVal(ctx.out, ret)
		}
	}
	return nil
}

// -----------------------------------------------------------------------------
