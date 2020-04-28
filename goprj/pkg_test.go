package goprj_test

import (
	"testing"

	"github.com/qiniu/qlang/goprj"
	"github.com/qiniu/x/log"
)

func init() {
	log.SetFlags(log.Llevel)
	log.SetOutputLevel(log.Ldebug)
}

func Test(t *testing.T) {
	prj, err := goprj.Open(".")
	if err != nil {
		t.Fatal(err)
	}
	pkg, err := prj.LoadPackage("github.com/qiniu/qlang/goprj")
	if err != nil {
		t.Fatal(err)
	}
	if pkg.Source().Name != "goprj" {
		t.Fatal("please run test in this package directory")
	}
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
