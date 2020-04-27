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
	pkgDir := "."
	names, err := goprj.OpenPkgNames(pkgDir)
	if err != nil {
		t.Fatal(err)
	}
	types := goprj.NewUniqueTypes()
	pkg, err := goprj.LoadPackage(pkgDir, names, types)
	if err != nil {
		t.Fatal(err)
	}
	if pkg.Name() != "goprj" {
		t.Fatal("please run test in this package directory")
	}
}
