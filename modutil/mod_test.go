package modutil_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/qiniu/qlang/modutil"
)

const (
	a  = 30
	a2 = 90
	a3
	b = iota
	c
	d, e = 2, 3
	f, g
	h = iota
)

func TestMod(t *testing.T) {
	fmt.Println(a3, b, c, f, g, h)
	if b != 3 || c != 4 {
		t.Fatal("const")
	}
	mod, err := modutil.LoadModule(".")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(mod.ModFile(), mod.RootPath(), mod.VersionPkgPath())
	pkg, err := mod.Lookup("github.com/qiniu/x/log")
	if err != nil || pkg.Type != modutil.PkgTypeDepMod {
		t.Fatal(err, pkg)
	}
	fmt.Println(pkg)
	pkg, err = mod.Lookup("fmt")
	if err != nil || pkg.Type != modutil.PkgTypeStd {
		t.Fatal(err, pkg, "goroot:", modutil.BuildContext.GOROOT)
	}
	fmt.Println(pkg)
}

func TestModFile(t *testing.T) {
	mod, err := modutil.LookupModFile(".")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(mod, "go.mod") {
		t.Fatal("not go.mod:", mod)
	}
}
