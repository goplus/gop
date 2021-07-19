/*
 Copyright 2021 The GoPlus Authors (goplus.org)

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

package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
	"github.com/goplus/gox"
)

// -----------------------------------------------------------------------------

func init() {
	cl.GenGoPkg = genGopkg
}

func getPkg(pkgs map[string]*ast.Package) *ast.Package {
	for _, pkg := range pkgs {
		return pkg
	}
	return nil
}

func genGopkg(pkgDir string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			switch v := e.(type) {
			case string:
				err = errors.New(v)
			case error:
				err = v
			default:
				panic(e)
			}
		}
	}()

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkgDir, nil, 0)
	if err != nil {
		return
	}
	if len(pkgs) != 1 {
		return fmt.Errorf("too many packages (%d) in the same directory", len(pkgs))
	}

	pkg := getPkg(pkgs)
	conf := &cl.Config{}
	out, err := cl.NewPackage("", pkg, fset, conf)
	if err != nil {
		return
	}
	err = os.MkdirAll(pkgDir, 0777)
	if err != nil {
		return
	}
	return gox.WriteFile(pkgDir+"/gop_autogen.go", out)
}

// -----------------------------------------------------------------------------
