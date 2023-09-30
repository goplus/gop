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
	"go/token"
	"strings"
	"syscall"

	"github.com/goplus/gop/cl/outline"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/x/gopenv"
	"github.com/goplus/mod/gopmod"
)

// -----------------------------------------------------------------------------

func Outline(dir string, conf *Config) (out outline.Package, err error) {
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
		ClassKind: mod.ClassKind,
		Filter:    conf.Filter,
		Mode:      parser.ParseComments,
	})
	if err != nil {
		return
	}
	if len(pkgs) == 0 {
		err = syscall.ENOENT
		return
	}

	imp := conf.Importer
	if imp == nil {
		imp = NewImporter(mod, gop, fset)
	}

	for name, pkg := range pkgs {
		if strings.HasSuffix(name, "_test") {
			continue
		}
		if out.Valid() {
			err = errMultiPackges
			return
		}
		if len(pkg.Files)+len(pkg.GoFiles) == 0 { // no Go/Go+ source files
			break
		}
		out, err = outline.NewPackage("", pkg, &outline.Config{
			Fset:        fset,
			WorkingDir:  dir,
			Importer:    imp,
			LookupClass: mod.LookupClass,
			LookupPub:   lookupPub(mod),
		})
		if err != nil {
			return
		}
	}
	if !out.Valid() {
		err = syscall.ENOENT
	}
	return
}

// -----------------------------------------------------------------------------

func OutlinePkgPath(workDir, pkgPath string, conf *Config, allowExtern bool) (out outline.Package, err error) {
	mod, err := gopmod.Load(workDir, 0)
	if NotFound(err) && allowExtern {
		remotePkgPathDo(pkgPath, func(dir string) {
			out, err = Outline(dir, conf)
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
	return Outline(pkg.Dir, conf)
}

// -----------------------------------------------------------------------------
