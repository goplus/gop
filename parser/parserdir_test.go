/*
 * Copyright (c) 2021 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package parser

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	"github.com/qiniu/x/log"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/parser/parsertest"
	"github.com/goplus/gop/scanner"
	"github.com/goplus/gop/token"
)

// -----------------------------------------------------------------------------

func init() {
	log.SetFlags(log.Llongfile)
	SetDebug(DbgFlagAll)
}

func TestReadSource(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	if _, err := readSource(buf); err != nil {
		t.Fatal("readSource failed:", err)
	}
	sr := strings.NewReader("")
	if _, err := readSource(sr); err != nil {
		t.Fatal("readSource strings.Reader failed:", err)
	}
	if _, err := readSource(0); err == nil {
		t.Fatal("readSource int failed: no error?")
	}
}

func TestParseFiles(t *testing.T) {
	fset := token.NewFileSet()
	if _, err := ParseFiles(fset, []string{"/foo/bar/not-exists"}, PackageClauseOnly); err == nil {
		t.Fatal("ParseFiles failed: no error?")
	}
}

func TestIparseFileInvalidSrc(t *testing.T) {
	fset := token.NewFileSet()
	if _, err := parseFile(fset, "/foo/bar/not-exists", 1, PackageClauseOnly); err != errInvalidSource {
		t.Fatal("ParseFile failed: not errInvalidSource?")
	}
}

func TestIparseFileNoFset(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Fatal("ParseFile failed: no error?")
		}
	}()
	parseFile(nil, "/foo/bar/not-exists", nil, PackageClauseOnly)
}

func TestParseDir(t *testing.T) {
	fset := token.NewFileSet()
	if _, err := ParseDir(fset, "/foo/bar/not-exists", nil, PackageClauseOnly); err == nil {
		t.Fatal("ParseDir failed: no error?")
	}
}

func testFrom(t *testing.T, pkgDir, sel string, exclude Mode) {
	if sel != "" && !strings.Contains(pkgDir, sel) {
		return
	}
	log.Println("Parsing", pkgDir)
	fset := token.NewFileSet()
	pkgs, err := ParseDir(fset, pkgDir, nil, (Trace|ParseGoFiles|ParseComments)&^exclude)
	if err != nil || len(pkgs) != 1 {
		if errs, ok := err.(scanner.ErrorList); ok {
			for _, e := range errs {
				t.Log(e)
			}
		}
		t.Fatal("ParseDir failed:", err, reflect.TypeOf(err), len(pkgs))
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

func TestParseGo(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := ParseDir(fset, "./_testdata/functype", nil, Trace)
	if err != nil {
		t.Fatal("TestParseGo: ParseDir failed -", err)
	}
	if len(pkgs) != 0 {
		t.Fatal("TestParseGo failed: len(pkgs) =", len(pkgs))
	}
}

func TestParseGoFiles(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := ParseFiles(fset, []string{"./_testdata/functype/functype.go"}, Trace)
	if err != nil {
		t.Fatal("TestParseGoFiles: ParseFiles failed -", err)
	}
	if len(pkgs) != 1 {
		t.Fatal("TestParseGoFiles failed: len(pkgs) =", len(pkgs))
	}
}

func TestFromTestdata(t *testing.T) {
	sel := ""
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
		name := fi.Name()
		if strings.HasPrefix(name, "_") {
			continue
		}
		t.Run(name, func(t *testing.T) {
			testFrom(t, dir+"/"+name, sel, 0)
		})
	}
}

func TestRegisterFileType(t *testing.T) {
	RegisterFileType(".gsh", ast.FileTypeSpx)
	func() {
		defer func() {
			if e := recover(); e == nil {
				t.Fatal("TestRegisterFileType failed: no error?")
			}
		}()
		RegisterFileType(".gshx", ast.FileTypeGop)
	}()
	func() {
		defer func() {
			if e := recover(); e == nil {
				t.Fatal("TestRegisterFileType failed: no error?")
			}
		}()
		RegisterFileType(".gsh", ast.FileTypeGmx)
	}()
}

// -----------------------------------------------------------------------------
