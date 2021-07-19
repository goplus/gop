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

package cl

import (
	"fmt"
	"go/token"
	"go/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/goplus/gox"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/packages"
)

// -----------------------------------------------------------------------------

func FindGoModFile(dir string) (file string, err error) {
	if dir == "" {
		dir = "."
	}
	if dir, err = filepath.Abs(dir); err != nil {
		return
	}
	for dir != "" {
		file = filepath.Join(dir, "go.mod")
		if fi, e := os.Lstat(file); e == nil && !fi.IsDir() {
			return
		}
		dir, _ = filepath.Split(dir)
	}
	return "", syscall.ENOENT
}

func GetModulePath(file string) (pkgPath string, err error) {
	src, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	f, err := modfile.ParseLax(file, src, nil)
	if err != nil {
		return
	}
	return f.Module.Mod.Path, nil
}

// -----------------------------------------------------------------------------

var (
	GenGoPkg func(pkgDir string) error
	GopGo    func(cfg *packages.Config, notFounds []string) error
)

func GenGoPkgs(cfg *packages.Config, notFounds []string) (err error) {
	if GenGoPkg == nil {
		return syscall.ENOENT
	}
	file, err := FindGoModFile(cfg.Dir)
	if err != nil {
		return
	}
	root, _ := filepath.Split(file)
	pkgPath, err := GetModulePath(file)
	if err != nil {
		return
	}
	pkgPathSlash := pkgPath + "/"
	for _, notFound := range notFounds {
		if strings.HasPrefix(notFound, pkgPathSlash) || notFound == pkgPath {
			if err = GenGoPkg(root + notFound[len(pkgPath):]); err != nil {
				return
			}
		} else {
			return syscall.ENOENT
		}
	}
	return nil
}

func LoadGop(cfg *packages.Config, patterns ...string) ([]*packages.Package, error) {
retry:
	loadPkgs, err := packages.Load(cfg, patterns...)
	if err == nil && GenGoPkg != nil {
		var notFounds []string
		packages.Visit(loadPkgs, nil, func(pkg *packages.Package) {
			const goGetCmd = "go get "
			for _, err := range pkg.Errors {
				if pos := strings.LastIndex(err.Msg, goGetCmd); pos > 0 {
					notFounds = append(notFounds, err.Msg[pos+len(goGetCmd):])
				}
			}
		})
		if notFounds != nil {
			gopGenGo := GopGo
			if gopGenGo == nil {
				gopGenGo = GenGoPkgs
			}
			if e := gopGenGo(cfg, notFounds); e == nil {
				goto retry
			}
		}
	}
	return loadPkgs, err
}

// LoadGopPkgs loads and returns the Go+ packages named by the given pkgPaths.
func LoadGopPkgs(at *gox.Package, importPkgs map[string]*gox.PkgRef, pkgPaths ...string) int {
	conf := at.InternalGetLoadConfig()
	loadPkgs, err := LoadGop(conf, pkgPaths...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if n := packages.PrintErrors(loadPkgs); n > 0 {
		return n
	}
	for _, loadPkg := range loadPkgs {
		if pkg, ok := importPkgs[loadPkg.PkgPath]; ok && pkg.ID == "" {
			pkg.ID = loadPkg.ID
			pkg.Errors = loadPkg.Errors
			pkg.Types = loadPkg.Types
			pkg.Fset = loadPkg.Fset
			pkg.Module = loadPkg.Module
			pkg.IllTyped = loadPkg.IllTyped
		}
	}
	return 0
}

// -----------------------------------------------------------------------------

func initGopBuiltin(pkg gox.PkgImporter, builtin *types.Package) {
	big := pkg.Import("github.com/goplus/gop/builtin")
	big.Ref("Gop_bigint")
	scope := big.Types.Scope()
	for i, n := 0, scope.Len(); i < n; i++ {
		names := scope.Names()
		for _, name := range names {
			builtin.Scope().Insert(scope.Lookup(name))
		}
	}
}

func initBuiltin(pkg gox.PkgImporter, builtin *types.Package) {
	fmt := pkg.Import("fmt")
	fns := []string{"print", "println", "printf", "errorf", "fprint", "fprintln", "fprintf"}
	for _, fn := range fns {
		fnTitle := string(fn[0]-'a'+'A') + fn[1:]
		builtin.Scope().Insert(gox.NewOverloadFunc(token.NoPos, builtin, fn, fmt.Ref(fnTitle)))
	}
}

func newBuiltinDefault(pkg gox.PkgImporter, prefix string, contracts *gox.BuiltinContracts) *types.Package {
	builtin := types.NewPackage("", "")
	initBuiltin(pkg, builtin)
	initGopBuiltin(pkg, builtin)
	gox.InitBuiltinOps(builtin, prefix, contracts)
	gox.InitBuiltinFuncs(builtin)
	return builtin
}

// -----------------------------------------------------------------------------
