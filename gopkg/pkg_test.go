package gopkg_test

import (
	"testing"

	"github.com/qiniu/qlang/gopkg"
	"github.com/qiniu/x/log"
)

func init() {
	log.SetFlags(log.Llevel)
	log.SetOutputLevel(log.Ldebug)
}

func Test(t *testing.T) {
	pkg, err := gopkg.Load(".")
	if err != nil {
		t.Fatal(err)
	}
	if pkg.Name() != "gopkg" {
		t.Fatal("please run test in this package directory")
	}
	pkg.Export()
}
