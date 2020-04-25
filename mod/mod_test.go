package mod

import (
	"fmt"
	"strings"
	"testing"
)

func TestMod(t *testing.T) {
	mods := GetModules()
	mod, err := mods.Load(".")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(mod.ModFile(), mod.Path())
	pkg, err := mod.Lookup("github.com/qiniu/x")
	if err != nil || pkg.Type != PkgTypeExternal {
		t.Fatal(err, pkg)
	}
	fmt.Println(pkg)
	if !strings.HasPrefix(pkg.Location, GetPkgModPath()) {
		t.Fatal("not at mod path:", pkg.Location)
	}
	pkg, err = mod.Lookup("fmt")
	if err != nil || pkg.Type != PkgTypeStd {
		t.Fatal(err, pkg, "goroot:", BuildContext.GOROOT)
	}
	fmt.Println(pkg)
}

func TestModFile(t *testing.T) {
	mod, err := LookupModFile(".")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(mod, "go.mod") {
		t.Fatal("not go.mod:", mod)
	}
}
