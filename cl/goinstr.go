package cl

import (
	"log"
	"reflect"

	"github.com/qiniu/qlang/v6/ast"
)

// -----------------------------------------------------------------------------

type goInstr = func(ctx *blockCtx, v *ast.CallExpr) func()

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
func igoAppend(ctx *blockCtx, v *ast.CallExpr) func() {
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
func igoCopy(ctx *blockCtx, v *ast.CallExpr) func() {
	panic("todo")
}

// func delete(m map[Type]Type1, key Type)
func igoDelete(ctx *blockCtx, v *ast.CallExpr) func() {
	panic("todo")
}

// func len(v Type) int
func igoLen(ctx *blockCtx, v *ast.CallExpr) func() {
	panic("todo")
}

// func cap(v Type) int
func igoCap(ctx *blockCtx, v *ast.CallExpr) func() {
	panic("todo")
}

// func make(t Type, size ...IntegerType) Type
func igoMake(ctx *blockCtx, v *ast.CallExpr) func() {
	if len(v.Args) < 1 {
		log.Panicln("make: argument count not enough")
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
func igoNew(ctx *blockCtx, v *ast.CallExpr) func() {
	panic("todo")
}

// func complex(r, i FloatType) ComplexType
func igoComplex(ctx *blockCtx, v *ast.CallExpr) func() {
	panic("todo")
}

// func real(c ComplexType) FloatType
func igoReal(ctx *blockCtx, v *ast.CallExpr) func() {
	panic("todo")
}

// func imag(c ComplexType) FloatType
func igoImag(ctx *blockCtx, v *ast.CallExpr) func() {
	panic("todo")
}

// func close(c chan<- Type)
func igoClose(ctx *blockCtx, v *ast.CallExpr) func() {
	panic("todo")
}

// func recover() interface{}
func igoRecover(ctx *blockCtx, v *ast.CallExpr) func() {
	panic("todo")
}

func compileTypeCast(typ reflect.Type, ctx *blockCtx, v *ast.CallExpr) func() {
	if len(v.Args) != 1 {
		log.Panicln("compileTypeCast: invalid argument count, please use `type(expr)`")
	}
	ctx.infer.Push(&goValue{typ})
	return func() {
		compileExpr(ctx, v.Args[0])()
		in := ctx.infer.Pop()
		if cons, ok := in.(*constVal); ok {
			cons.bound(typ, ctx.out)
			return
		}
		iv := in.(iValue)
		n := iv.NumValues()
		if n != 1 {
			panicExprNotValue(n)
		}
		tIn := iv.Type()
		if !tIn.ConvertibleTo(typ) {
			log.Panicf("compileTypeCast: can't convert type `%v` to `%v`\n", tIn, typ)
		}
		ctx.out.TypeCast(tIn, typ)
	}
}

// -----------------------------------------------------------------------------
