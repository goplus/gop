/*
 * Copyright (c) 2024 The GoPlus Authors (goplus.org). All rights reserved.
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
	"bytes"
	"os"
	"testing"

	"github.com/goplus/gogen"
	"github.com/goplus/gop"
	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/parser/fsx/memfs"
	"github.com/goplus/gop/scanner"
	"github.com/goplus/gop/token"
	"github.com/goplus/mod/env"
	"github.com/goplus/mod/modfile"
)

var (
	Gop  *env.Gop
	Conf *cl.Config
)

func init() {
	Gop = &env.Gop{Version: "1.0"}
	gogen.SetDebug(gogen.DbgFlagAll)
	cl.SetDebug(cl.DbgFlagAll | cl.FlagNoMarkAutogen)
	fset := token.NewFileSet()
	imp := gop.NewImporter(nil, Gop, fset)
	Conf = &cl.Config{
		Fset:          fset,
		Importer:      imp,
		Recorder:      gopRecorder{},
		LookupClass:   LookupClass,
		NoFileLine:    true,
		NoAutoGenMain: true,
	}
}

func LookupClass(ext string) (c *modfile.Project, ok bool) {
	switch ext {
	case ".tgmx", ".tspx":
		return &modfile.Project{
			Ext: ".tgmx", Class: "*MyGame",
			Works:    []*modfile.Class{{Ext: ".tspx", Class: "Sprite"}},
			PkgPaths: []string{"github.com/goplus/gop/cl/internal/spx", "math"}}, true
	case ".t2gmx", ".t2spx", ".t2spx2":
		return &modfile.Project{
			Ext: ".t2gmx", Class: "Game",
			Works: []*modfile.Class{{Ext: ".t2spx", Class: "Sprite"},
				{Ext: ".t2spx2", Class: "Sprite2"}},
			PkgPaths: []string{"github.com/goplus/gop/cl/internal/spx2"}}, true
	case "_t3spx.gox", ".t3spx2":
		return &modfile.Project{
			Works: []*modfile.Class{{Ext: "_t3spx.gox", Class: "Sprite"},
				{Ext: ".t3spx2", Class: "Sprite2"}},
			PkgPaths: []string{"github.com/goplus/gop/cl/internal/spx2"}}, true
	case "_spx.gox":
		return &modfile.Project{
			Ext: "_spx.gox", Class: "Game",
			Works:    []*modfile.Class{{Ext: "_spx.gox", Class: "Sprite"}},
			PkgPaths: []string{"github.com/goplus/gop/cl/internal/spx3", "math"},
			Import:   []*modfile.Import{{Path: "github.com/goplus/gop/cl/internal/spx3/jwt"}}}, true
	case "_xtest.gox":
		return &modfile.Project{
			Ext: "_xtest.gox", Class: "App",
			Works:    []*modfile.Class{{Ext: "_xtest.gox", Class: "Case"}},
			PkgPaths: []string{"github.com/goplus/gop/test", "testing"}}, true
	}
	return
}

func Named(t *testing.T, name string, gopcode, expected string) {
	t.Run(name, func(t *testing.T) {
		Do(t, gopcode, expected)
	})
}

func Do(t *testing.T, gopcode, expected string) {
	DoExt(t, Conf, "main", gopcode, expected)
}

func DoWithFname(t *testing.T, gopcode, expected string, fname string) {
	fs := memfs.SingleFile("/foo", fname, gopcode)
	DoFS(t, Conf, fs, "main", expected)
}

func DoExt(t *testing.T, conf *cl.Config, pkgname, gopcode, expected string) {
	fs := memfs.SingleFile("/foo", "bar.gop", gopcode)
	DoFS(t, conf, fs, pkgname, expected)
}

func Mixed(t *testing.T, pkgname, gocode, gopcode, expected string, outline ...bool) {
	conf := *Conf
	conf.Outline = (outline != nil && outline[0])
	fs := memfs.TwoFiles("/foo", "a.go", gocode, "b.gop", gopcode)
	DoFS(t, &conf, fs, pkgname, expected)
}

func DoFS(t *testing.T, conf *cl.Config, fs parser.FileSystem, pkgname, expected string) {
	cl.SetDisableRecover(true)
	defer cl.SetDisableRecover(false)

	fset := conf.Fset
	pkgs, err := parser.ParseFSDir(fset, fs, "/foo", parser.Config{Mode: parser.ParseComments})
	if err != nil {
		scanner.PrintError(os.Stderr, err)
		t.Fatal("ParseFSDir:", err)
	}
	bar := pkgs[pkgname]
	pkg, err := cl.NewPackage("github.com/goplus/gop/cl", bar, conf)
	if err != nil {
		t.Fatal("NewPackage:", err)
	}
	var b bytes.Buffer
	err = pkg.WriteTo(&b)
	if err != nil {
		t.Fatal("gogen.WriteTo failed:", err)
	}
	result := b.String()
	if result != expected {
		t.Fatalf("\nResult:\n%s\nExpected:\n%s\n", result, expected)
	}
}
