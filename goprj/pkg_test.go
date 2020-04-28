package goprj_test

import (
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
	"testing"

	"github.com/qiniu/qlang/goprj"
	"github.com/qiniu/x/log"
)

func init() {
	log.SetFlags(log.Llevel)
	log.SetOutputLevel(log.Ldebug)
}

type simpleTypeInferer struct {
}

func (p *simpleTypeInferer) InferType(expr ast.Expr) goprj.Type {
	return &goprj.UninferedType{expr}
}

func (p *simpleTypeInferer) InferConst(expr ast.Expr, i int) (typ goprj.Type, val interface{}) {
	switch v := expr.(type) {
	case *ast.BasicLit:
		switch v.Kind {
		case token.INT:
			n, err := strconv.ParseInt(v.Value, 0, 0)
			if err != nil {
				log.Fatalln("InferConst: strconv.ParseInt failed:", err)
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
	case *ast.SelectorExpr:
		if i < 0 {
			return goprj.AtomType(goprj.Int), 0
		}
	case *ast.BinaryExpr:
		if i < 0 {
			return goprj.AtomType(goprj.Int), 0
		}
	default:
		if i < 0 {
			log.Fatalln("InferConst: unknown -", reflect.TypeOf(expr), "-", expr)
		}
	}
	return &goprj.UninferedType{expr}, expr
}

func Test(t *testing.T) {
	prj, err := goprj.Open(".")
	if err != nil {
		t.Fatal(err)
	}
	prj.TypeInferrer = &simpleTypeInferer{}
	pkg, err := prj.LoadPackage("github.com/qiniu/qlang/goprj")
	if err != nil {
		t.Fatal(err)
	}
	if pkg.Source().Name != "goprj" {
		t.Fatal("please run test in this package directory")
	}
	if pkg.PkgPath() != "github.com/qiniu/qlang/goprj" {
		t.Fatal("PkgPath:", pkg.PkgPath())
	}
	return
	log.Debug("------------------------------------------------------")
	pkg2, err := prj.LoadPackage("github.com/qiniu/qlang/modutil")
	if err != nil {
		t.Fatal(err)
	}
	if pkg2.Source().Name != "modutil" {
		t.Fatal("please run test in this package directory")
	}
	log.Debug("------------------------------------------------------")
	prjPath := "github.com/visualfc/fastmod"
	pi, err := prj.ThisModule().Lookup(prjPath)
	if err != nil {
		t.Fatal(err)
	}
	prj2, err := goprj.Open(pi.Location)
	if err != nil {
		t.Fatal(err)
	}
	pkg3, err := prj2.Load()
	if err != nil {
		t.Fatal(err)
	}
	if pkg3.PkgPath() != "github.com/visualfc/fastmod@v1.3.3" {
		t.Fatal("PkgPath:", pkg3.PkgPath())
	}
}
