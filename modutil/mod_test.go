package modutil_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/qiniu/qlang/modutil"
)

func TestMod(t *testing.T) {
	mods := modutil.GetModules()
	mod, err := mods.Load(".")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(mod.ModFile(), mod.Path())
	pkg, err := mod.Lookup("github.com/qiniu/x/log")
	if err != nil || pkg.Type != modutil.PkgTypeDepMod {
		t.Fatal(err, pkg)
	}
	fmt.Println(pkg)
	if !strings.HasPrefix(pkg.Location, modutil.GetPkgModPath()) {
		t.Fatal("not at mod path:", pkg.Location)
	}
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
