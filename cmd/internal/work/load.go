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

package work

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
)

type Package struct {
	Name      string       // package name
	Dir       string       // package source dir
	Pkg       *ast.Package // ast.Package
	Target    string       // build target
	IsDir     bool         // build input is dir
	GenGoFile string       // generate go file path
}

func LoadPackages(fset *token.FileSet, args []string) (pkgs []*Package, errs []error) {
	for _, arg := range args {
		path, err := filepath.Abs(arg)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		fi, err := os.Stat(path)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		var apkgs map[string]*ast.Package
		var dir string
		var target string
		isDir := fi.IsDir()
		if isDir {
			dir = path
			_, target = filepath.Split(path)
			apkgs, err = parser.ParseDir(fset, path, nil, 0)
		} else {
			dir, target = filepath.Split(path)
			dir = filepath.Clean(dir)
			ext := filepath.Ext(target)
			target = target[:len(target)-len(ext)]
			apkgs, err = parser.Parse(fset, path, nil, 0)
		}
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if runtime.GOOS == "windows" {
			target += ".exe"
		}
		for name, pkg := range apkgs {
			var mainTraget string
			if pkg.Name == "main" {
				mainTraget = target
			}
			pkgs = append(pkgs, &Package{IsDir: isDir, Name: name, Dir: dir, Pkg: pkg, Target: mainTraget})
		}
	}
	return
}
