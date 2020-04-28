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
}
