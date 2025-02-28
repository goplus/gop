/*
 * Copyright (c) 2025 The GoPlus Authors (goplus.org). All rights reserved.
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
	"go/scanner"
	"log"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	"github.com/goplus/gop/tpl/parser/parsertest"
	"github.com/goplus/gop/tpl/token"
)

func testFrom(t *testing.T, pkgDir, sel string) {
	if sel != "" && !strings.Contains(pkgDir, sel) {
		return
	}
	t.Helper()
	log.Println("Parsing", pkgDir)
	fset := token.NewFileSet()
	f, err := ParseFile(fset, pkgDir+"/in.gop", nil, 0)
	if err != nil {
		if errs, ok := err.(scanner.ErrorList); ok {
			for _, e := range errs {
				t.Log(e)
			}
		}
		t.Fatal("ParseFile failed:", err, reflect.TypeOf(err))
	}
	b, err := os.ReadFile(pkgDir + "/out.expect")
	if err != nil {
		t.Fatal("Parsing", pkgDir, "-", err)
	}
	parsertest.Expect(t, f, string(b))
}

func testFromDir(t *testing.T, sel, relDir string) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Getwd failed:", err)
	}
	dir = path.Join(dir, relDir)
	fis, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal("ReadDir failed:", err)
	}
	for _, fi := range fis {
		name := fi.Name()
		if strings.HasPrefix(name, "_") {
			continue
		}
		t.Run(name, func(t *testing.T) {
			testFrom(t, dir+"/"+name, sel)
		})
	}
}

func TestFromTestdata(t *testing.T) {
	testFromDir(t, "", "./_testdata")
}
