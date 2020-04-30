package exec

import (
	"go/ast"
	"go/token"
	"log"
	"reflect"
	"strconv"

	"github.com/qiniu/qlang/goprj"
)

// -----------------------------------------------------------------------------

// InferType infers type from a ast.Expr.
func InferType(pkg *goprj.Package, expr ast.Expr, reserved int) goprj.Type {
	switch v := expr.(type) {
	case *ast.CallExpr:
		if t, ok := inferTypeFromFun(pkg, v.Fun); ok {
			return t
		}
	case *ast.UnaryExpr:
		switch v.Op {
		case token.AND: // &X
			t := InferType(pkg, v.X, reserved)
			return pkg.Project().UniqueType(&goprj.PointerType{Elem: t})
		default:
			log.Fatalln("InferType: unknown UnaryExpr -", v.Op)
		}
	case *ast.SelectorExpr:
		_ = v
	}
	log.Fatalln("InferType: unknown -", reflect.TypeOf(expr))
	return nil
}

func inferTypeFromFun(pkg *goprj.Package, fun ast.Expr) (t goprj.Type, ok bool) {
	switch v := fun.(type) {
	case *ast.SelectorExpr:
		switch recv := v.X.(type) {
		case *ast.Ident:
			fnt, err := pkg.FindPackageType(recv.Name, v.Sel.Name)
			if err == nil {
				if fn, ok := fnt.(*goprj.FuncType); ok {
					return &goprj.RetType{Results: fn.Results}, true
				}
			}
			log.Fatalln("inferTypeFromFun: FindPackageType error -", err)
		default:
			log.Fatalln("inferTypeFromFun:", reflect.TypeOf(v.X))
		}
	case *ast.Ident:
		fnt, err := pkg.FindType(v.Name)
		if err == nil {
			if fn, ok := fnt.(*goprj.FuncType); ok {
				return &goprj.RetType{Results: fn.Results}, true
			}
		}
		log.Fatalln("inferTypeFromFun: FindType error -", err)
	default:
		log.Fatalln("inferTypeFromFun: unknown -", reflect.TypeOf(fun))
	}
	return nil, false
}

// InferConst infers constant value from a ast.Expr.
func InferConst(pkg *goprj.Package, expr ast.Expr, i int) (typ goprj.Type, val interface{}) {
	switch v := expr.(type) {
	case *ast.BasicLit:
		switch v.Kind {
		case token.INT:
			n, err := strconv.ParseInt(v.Value, 0, 0)
			if err != nil {
				n2, err2 := strconv.ParseUint(v.Value, 0, 0)
				if err2 != nil {
					log.Fatalln("InferConst: strconv.ParseInt failed:", err2)
				}
				return goprj.AtomType(goprj.Uint), uint(n2)
			}
			return goprj.AtomType(goprj.Int), int(n)
		case token.FLOAT:
			n, err := strconv.ParseFloat(v.Value, 64)
			if err != nil {
				log.Fatalln("InferConst: strconv.ParseFloat failed:", err)
			}
			return goprj.AtomType(goprj.Float64), n
		case token.CHAR, token.STRING:
			n, err := strconv.Unquote(v.Value)
			if err != nil {
				log.Fatalln("InferConst: strconv.Unquote failed:", err)
			}
			if v.Kind == token.CHAR {
				for _, c := range n {
					return goprj.AtomType(goprj.Rune), c
				}
				panic("not here")
			}
			return goprj.AtomType(goprj.String), n
		default:
			log.Fatalln("InferConst: unknown -", expr)
		}
	case *ast.Ident:
	case *ast.SelectorExpr:
	case *ast.BinaryExpr:
	default:
		if i < 0 {
			log.Fatalln("InferConst: unknown -", reflect.TypeOf(expr), "-", expr)
		}
	}
	return &goprj.UninferedType{Expr: expr}, expr
}

// -----------------------------------------------------------------------------
