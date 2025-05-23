/*
 * Copyright (c) 2024 The XGo Authors (xgo.dev). All rights reserved.
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

package cltest

import (
	"os"
	"testing"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/parser/fsx/memfs"
	"github.com/goplus/gop/scanner"
)

func Error(t *testing.T, msg, src string) {
	ErrorEx(t, "main", "bar.xgo", msg, src)
}

func ErrorEx(t *testing.T, pkgname, filename, msg, src string) {
	fs := memfs.SingleFile("/foo", filename, src)
	pkgs, err := parser.ParseFSDir(Conf.Fset, fs, "/foo", parser.Config{})
	if err != nil {
		scanner.PrintError(os.Stderr, err)
		t.Fatal("parser.ParseFSDir failed")
	}
	conf := *Conf
	conf.NoFileLine = false
	conf.RelativeBase = "/foo"
	bar := pkgs[pkgname]
	_, err = cl.NewPackage("", bar, &conf)
	if err == nil {
		t.Fatal("no error?")
	}
	if ret := err.Error(); ret != msg {
		t.Fatalf("\nError: \"%s\"\nExpected: \"%s\"\n", ret, msg)
	}
}

func ErrorAst(t *testing.T, pkgname, filename, msg, src string) {
	f, _ := parser.ParseFile(Conf.Fset, filename, src, parser.AllErrors)
	pkg := &ast.Package{
		Name:  pkgname,
		Files: map[string]*ast.File{filename: f},
	}
	conf := *Conf
	conf.NoFileLine = false
	conf.RelativeBase = "/foo"
	_, err := cl.NewPackage("", pkg, &conf)
	if err == nil {
		t.Fatal("no error?")
	}
	if ret := err.Error(); ret != msg {
		t.Fatalf("\nError: \"%s\"\nExpected: \"%s\"\n", ret, msg)
	}
}
