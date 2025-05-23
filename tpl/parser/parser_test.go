/*
 * Copyright (c) 2025 The XGo Authors (xgo.dev). All rights reserved.
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

package parser_test

import (
	"log"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	gopp "github.com/goplus/gop/parser"
	"github.com/goplus/gop/tpl/ast"
	"github.com/goplus/gop/tpl/parser"
	"github.com/goplus/gop/tpl/parser/parsertest"
	"github.com/goplus/gop/tpl/scanner"
	"github.com/goplus/gop/tpl/token"
)

func testFrom(t *testing.T, pkgDir, sel string) {
	if sel != "" && !strings.Contains(pkgDir, sel) {
		return
	}
	t.Helper()
	log.Println("Parsing", pkgDir)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, pkgDir+"/in.gop", nil, &parser.Config{
		ParseRetProc: func(file *token.File, src []byte, offset int) (ast.Node, scanner.ErrorList) {
			return gopp.ParseExprEx(file, src, offset, 0)
		},
	})
	if err != nil {
		if errs, ok := err.(scanner.ErrorList); ok {
			for _, e := range errs {
				t.Log(e)
			}
		}
		t.Fatal("ParseFile failed:", err, reflect.TypeOf(err))
	}
	b, _ := os.ReadFile(pkgDir + "/out.expect")
	parsertest.Expect(t, pkgDir+"/result.txt", f, b)
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
