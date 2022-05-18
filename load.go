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
	"errors"
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
	"github.com/goplus/gox/packages"
	"github.com/goplus/mod/env"
	"github.com/goplus/mod/gopmod"
	"github.com/goplus/mod/modcache"
	"github.com/goplus/mod/modfetch"
)

type Config struct {
	Gop      *env.Gop
	Fset     *token.FileSet
	Filter   func(fs.FileInfo) bool
	Importer types.Importer

	UpdateGoMod     bool
	CheckModChanged bool
}

// -----------------------------------------------------------------------------

func LoadDir(dir string, conf *Config) (out, test *gox.Package, err error) {
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
	imp := conf.Importer
	if imp == nil {
		imp = packages.NewImporter(fset, dir)
	}

	mod, err := gopmod.Load(dir, gop)
	if err != nil && err != syscall.ENOENT {
		return
	}
	if mod != nil {
		err = mod.RegisterClasses()
		if err != nil {
			return
		}
		if conf.UpdateGoMod {
			err = mod.UpdateGoMod(gop, conf.CheckModChanged)
			if err != nil {
				return
			}
		}
	} else {
		mod = new(gopmod.Module)
	}

	pkgs, err := parser.ParseDirEx(fset, dir, parser.Config{
		IsClass: mod.IsClass,
		Filter:  conf.Filter,
		Mode:    parser.ParseComments,
	})
	if err != nil {
		return
	}

	var pkgTest *ast.Package
	var clConf = &cl.Config{
		WorkingDir:  dir,
		GopRoot:     gop.Root,
		Fset:        fset,
		Importer:    imp,
		LookupClass: mod.LookupClass,
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
	if pkgTest != nil {
		test, err = cl.NewPackage("", pkgTest, clConf)
	}
	return
}

// -----------------------------------------------------------------------------

func LoadPkgPath(pkgPath string, conf *Config) (out, test *gox.Package, err error) {
	getPkgPathDo(pkgPath, gopEnv(conf), func(dir string) {
		out, test, err = LoadDir(dir, conf)
	}, func(e error) {
		err = e
	})
	return
}

func getPkgPathDo(pkgPath string, gop *env.Gop, doSth func(dir string), onErr func(e error)) {
	modPath, leftPart := splitPkgPath(pkgPath)
	modVer, _, err := modfetch.Get(gop, modPath)
	if err != nil {
		onErr(err)
	} else if dir, err := modcache.Path(modVer); err != nil {
		onErr(err)
	} else {
		doSth(filepath.Join(dir, leftPart))
	}
}

func splitPkgPath(pkgPath string) (modPath, leftPart string) {
	if strings.HasPrefix(pkgPath, "github.com/") {
		parts := strings.SplitN(pkgPath, "/", 4)
		if len(parts) > 3 {
			leftPart = parts[3]
			if pos := strings.IndexByte(leftPart, '@'); pos > 0 {
				parts[2] += leftPart[pos:]
				leftPart = leftPart[:pos]
			}
			modPath = strings.Join(parts[:3], "/")
			return
		}
	}
	return pkgPath, "" // TODO:
}

func gopEnv(conf *Config) *env.Gop {
	if conf != nil {
		if gop := conf.Gop; gop != nil {
			return gop
		}
	}
	return gopenv.Get()
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

	pkgs, err := parser.ParseFiles(fset, files, parser.ParseComments)
	if err != nil {
		return
	}
	if len(pkgs) != 1 {
		return nil, errMultiPackges
	}
	for _, pkg := range pkgs {
		out, err = cl.NewPackage("", pkg, &cl.Config{
			Fset:     fset,
			Importer: conf.Importer,
		})
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
