/*
 * Copyright (c) 2023 The GoPlus Authors (goplus.org). All rights reserved.
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
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/goplus/gop/cl/outline"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/x/c2go"
	"github.com/goplus/gop/x/gopenv"
	"github.com/goplus/mod/gopmod"
)

// -----------------------------------------------------------------------------

func Outline(dir string, conf *Config) (out outline.Package, err error) {
	if dir, err = filepath.Abs(dir); err != nil {
		return
	}
	mod, err := LoadMod(dir)
	if err != nil {
		return
	}

	if conf == nil {
		conf = new(Config)
	}
	filterConf := conf.Filter
	filter := func(fi fs.FileInfo) bool {
		if filterConf != nil && !filterConf(fi) {
			return false
		}
		fname := fi.Name()
		if pos := strings.Index(fname, "."); pos > 0 {
			fname = fname[:pos]
		}
		return !strings.HasSuffix(fname, "_test")
	}
	fset := conf.Fset
	if fset == nil {
		fset = token.NewFileSet()
	}
	pkgs, err := parser.ParseDirEx(fset, dir, parser.Config{
		ClassKind: mod.ClassKind,
		Filter:    filter,
		Mode:      parser.ParseComments,
	})
	if err != nil {
		return
	}
	if len(pkgs) == 0 {
		err = ErrNotFound
		return
	}

	imp := conf.Importer
	if imp == nil {
		gop := conf.Gop
		if gop == nil {
			gop = gopenv.Get()
		}
		imp = NewImporter(mod, gop, fset)
	}

	for name, pkg := range pkgs {
		if out.Valid() {
			err = fmt.Errorf("%w: %s, %s", ErrMultiPackges, name, out.Pkg().Name())
			return
		}
		if len(pkg.Files)+len(pkg.GoFiles) == 0 { // no Go/Go+ source files
			break
		}
		relPart, _ := filepath.Rel(mod.Root(), dir)
		pkgPath := path.Join(mod.Path(), filepath.ToSlash(relPart))
		out, err = outline.NewPackage(pkgPath, pkg, &outline.Config{
			Fset:        fset,
			Importer:    imp,
			LookupClass: mod.LookupClass,
			LookupPub:   c2go.LookupPub(mod),
		})
		if err != nil {
			return
		}
	}
	if !out.Valid() {
		err = ErrNotFound
	}
	return
}

// -----------------------------------------------------------------------------

func OutlinePkgPath(workDir, pkgPath string, conf *Config, allowExtern bool) (out outline.Package, err error) {
	mod, err := gopmod.Load(workDir)
	if NotFound(err) && allowExtern {
		remotePkgPathDo(pkgPath, func(pkgDir, modDir string) {
			modFile := chmodModfile(modDir)
			defer os.Chmod(modFile, modReadonly)
			out, err = Outline(pkgDir, conf)
		}, func(e error) {
			err = e
		})
		return
	} else if err != nil {
		return
	}

	pkg, err := mod.Lookup(pkgPath)
	if err != nil {
		return
	}
	if pkg.Type == gopmod.PkgtExtern {
		modFile := chmodModfile(pkg.ModDir)
		defer os.Chmod(modFile, modReadonly)
	}
	return Outline(pkg.Dir, conf)
}

func chmodModfile(modDir string) string {
	modFile := modDir + "/go.mod"
	os.Chmod(modFile, modWritable)
	return modFile
}

// -----------------------------------------------------------------------------
