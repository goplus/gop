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

package parser

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/goplus/gop/parser/parsertest"
	"github.com/goplus/gop/token"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

func init() {
	log.SetOutputLevel(log.Ldebug)
	log.SetFlags(log.Llevel)
}

func testFrom(t *testing.T, pkgDir, sel string) {
	if sel != "" && !strings.Contains(pkgDir, sel) {
		return
	}
	log.Debug("Parsing", pkgDir)
	fset := token.NewFileSet()
	pkgs, err := ParseDir(fset, pkgDir, nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseDir failed:", err, len(pkgs))
	}
	for _, pkg := range pkgs {
		b, err := ioutil.ReadFile(pkgDir + "/parser.expect")
		if err != nil {
			t.Fatal("Parsing", pkgDir, "-", err)
		}
		parsertest.Expect(t, pkg, string(b))
		return
	}
}

func TestFromTestdata(t *testing.T) {
	sel := "exists"
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Getwd failed:", err)
	}
	dir = path.Join(dir, "./_testdata")
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal("ReadDir failed:", err)
	}
	for _, fi := range fis {
		testFrom(t, dir+"/"+fi.Name(), sel)
	}
}

// -----------------------------------------------------------------------------
