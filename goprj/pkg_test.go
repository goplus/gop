package goprj_test

import (
	"go/ast"
	"testing"

	"github.com/qiniu/qlang/exec"
	"github.com/qiniu/qlang/goprj"
	"github.com/qiniu/x/log"
)

func init() {
	log.SetFlags(log.Llevel)
	//log.SetOutputLevel(log.Ldebug)
}

type typeInferrer struct {
}

func (p *typeInferrer) InferType(pkg *goprj.Package, expr ast.Expr, reserved int) goprj.Type {
	if reserved < 0 { // *RetType
		exec.InferType(pkg, expr, reserved)
	}
	return &goprj.UninferedType{expr}
}

func (p *typeInferrer) InferConst(pkg *goprj.Package, expr ast.Expr, i int) (typ goprj.Type, val interface{}) {
	typ, val = exec.InferConst(pkg, expr, i)
	if _, ok := typ.(*goprj.UninferedType); ok {
		if i < 0 {
			return goprj.AtomType(goprj.Int), 0
		}
	}
	return
}

func Test(t *testing.T) {
	prj := goprj.NewProject()
	prj.TypeInferrer = &typeInferrer{}
	pkg, err := prj.OpenPackage(".")
	if err != nil {
		t.Fatal(err)
	}
	if pkg.Source().Name != "goprj" {
		t.Fatal("please run test in this package directory")
	}
	if pkg.ThisModule().PkgPath() != "github.com/qiniu/qlang/goprj" {
		t.Fatal("PkgPath:", pkg.ThisModule().PkgPath())
	}
	return
	log.Debug("------------------------------------------------------")
	pkg2, err := pkg.LoadPackage("github.com/qiniu/qlang/modutil")
	if err != nil {
		t.Fatal(err)
	}
	if pkg2.Source().Name != "modutil" {
		t.Fatal("please run test in this package directory")
	}
	log.Debug("------------------------------------------------------")
	prjPath := "github.com/visualfc/fastmod"
	pkg3, err := pkg2.LoadPackage(prjPath)
	if err != nil {
		t.Fatal(err)
	}
	if pkg3.ThisModule().VersionPkgPath() != "github.com/visualfc/fastmod@v1.3.3" {
		t.Fatal("PkgPath:", pkg3.ThisModule().VersionPkgPath())
	}
}
