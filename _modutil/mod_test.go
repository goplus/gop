/*
 Copyright 2020 The GoPlus Authors (goplus.org)

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package modutil_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/qiniu/goplus/modutil"
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
