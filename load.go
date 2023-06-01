/*
 * Copyright (c) 2022 The GoPlus Authors (goplus.org). All rights reserved.
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

package gop

import (
	"fmt"
	"go/token"
	"go/types"
	"io/fs"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/x/gopenv"
	"github.com/goplus/gox"
	"github.com/goplus/mod/env"
	"github.com/goplus/mod/gopmod"
	"github.com/qiniu/x/errors"
)

type Config struct {
	Gop      *env.Gop
	Fset     *token.FileSet
	Filter   func(fs.FileInfo) bool
	Importer types.Importer

	DontUpdateGoMod     bool
	DontCheckModChanged bool
}

// -----------------------------------------------------------------------------

func NotFound(err error) bool {
	return errors.Err(err) == syscall.ENOENT
}

func LoadMod(dir string, gop *env.Gop, conf *Config) (mod *gopmod.Module, err error) {
	mod, err = gopmod.Load(dir, 0)
	if err != nil && !NotFound(err) {
		err = errors.NewWith(err, `gopmod.Load(dir, 0`, -2, "gopmod.Load", dir, 0)
		return
	}
	if mod != nil {
		err = mod.RegisterClasses()
		if err != nil {
			err = errors.NewWith(err, `mod.RegisterClasses()`, -2, "(*gopmod.Module).RegisterClasses", mod)
			return
		}
		if !conf.DontUpdateGoMod {
			err = mod.UpdateGoMod(gop, !conf.DontCheckModChanged)
			if err != nil {
				err = errors.NewWith(err, `mod.UpdateGoMod(gop, !conf.DontCheckModChanged)`, -2, "(*gopmod.Module).UpdateGoMod", mod, gop, !conf.DontCheckModChanged)
			}
		}
		return
	}
	return new(gopmod.Module), nil
}

func lookupPub(mod *gopmod.Module) func(pkgPath string) (pubfile string, err error) {
	return func(pkgPath string) (pubfile string, err error) {
		if mod.File == nil { // no go.mod/gop.mod file
			return "", syscall.ENOENT
		}
		pkg, err := mod.Lookup(pkgPath)
		if err == nil {
			pubfile = filepath.Join(pkg.Dir, "c2go.a.pub")
		}
		return
	}
}

// -----------------------------------------------------------------------------

func LoadDir(dir string, conf *Config, genTestPkg bool, promptGenGo ...bool) (out, test *gox.Package, err error) {
	if conf == nil {
		conf = new(Config)
	}
	fset := conf.Fset
	if fset == nil {
		fset = token.NewFileSet()
	}
	gop := conf.Gop
	if gop == nil {
		gop = gopenv.Get()
	}
	mod, err := LoadMod(dir, gop, conf)
	if err != nil {
		return
	}

	pkgs, err := parser.ParseDirEx(fset, dir, parser.Config{
		IsClass: mod.IsClass,
		Filter:  conf.Filter,
		Mode:    parser.ParseComments,
	})
	if err != nil {
		return
	}
	if len(pkgs) == 0 {
		return nil, nil, syscall.ENOENT
	}

	if promptGenGo != nil && promptGenGo[0] {
		fmt.Printf("GenGo %v ...\n", dir)
	}

	imp := conf.Importer
	if imp == nil {
		imp = NewImporter(mod, gop, fset)
	}

	var pkgTest *ast.Package
	var clConf = &cl.Config{
		WorkingDir:  dir,
		Fset:        fset,
		Importer:    imp,
		LookupClass: mod.LookupClass,
		LookupPub:   lookupPub(mod),
	}
	for name, pkg := range pkgs {
		if strings.HasSuffix(name, "_test") {
			if pkgTest != nil {
				return nil, nil, errMultiTestPackges
			}
			pkgTest = pkg
			continue
		}
		if out != nil {
			return nil, nil, errMultiPackges
		}
		if len(pkg.Files) == 0 { // no Go+ source files
			break
		}
		out, err = cl.NewPackage("", pkg, clConf)
		if err != nil {
			return
		}
	}
	if out == nil {
		return nil, nil, syscall.ENOENT
	}
	if pkgTest != nil && genTestPkg {
		test, err = cl.NewPackage("", pkgTest, clConf)
	}
	return
}

// -----------------------------------------------------------------------------

func LoadFiles(files []string, conf *Config) (out *gox.Package, err error) {
	if conf == nil {
		conf = new(Config)
	}
	fset := conf.Fset
	if fset == nil {
		fset = token.NewFileSet()
	}
	gop := conf.Gop
	if gop == nil {
		gop = gopenv.Get()
	}
	mod, err := LoadMod("", gop, conf)
	if err != nil {
		err = errors.NewWith(err, `loadMod("", gop, conf)`, -2, "gop.loadMod", "", gop, conf)
		return
	}

	pkgs, err := parser.ParseFiles(fset, files, parser.ParseComments)
	if err != nil {
		err = errors.NewWith(err, `parser.ParseFiles(fset, files, parser.ParseComments)`, -2, "parser.ParseFiles", fset, files, parser.ParseComments)
		return
	}
	if len(pkgs) != 1 {
		err = errors.NewWith(errMultiPackges, `len(pkgs) != 1`, -1, "!=", len(pkgs), 1)
		return
	}
	for _, pkg := range pkgs {
		imp := conf.Importer
		if imp == nil {
			imp = NewImporter(mod, gop, fset)
		}
		clConf := &cl.Config{
			Fset:        fset,
			Importer:    imp,
			LookupClass: mod.LookupClass,
			LookupPub:   lookupPub(mod),
		}
		out, err = cl.NewPackage("", pkg, clConf)
		if err != nil {
			err = errors.NewWith(err, `cl.NewPackage("", pkg, clConf)`, -2, "cl.NewPackage", "", pkg, clConf)
		}
		break
	}
	return
}

// -----------------------------------------------------------------------------

var (
	errMultiPackges     = errors.New("multiple packages")
	errMultiTestPackges = errors.New("multiple test packages")
)

// -----------------------------------------------------------------------------
