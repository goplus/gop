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
	prj := goprj.NewProject()
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
