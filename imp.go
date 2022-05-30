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
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"github.com/goplus/gox/packages"
	"github.com/goplus/mod/env"
	"github.com/goplus/mod/gopmod"
	"github.com/goplus/mod/modfetch"
)

// -----------------------------------------------------------------------------

type Importer struct {
	impFrom *packages.Importer
	mod     *gopmod.Module
	gop     *env.Gop
	gopRoot string
}

func NewImporter(mod *gopmod.Module, gop *env.Gop, fset *token.FileSet) *Importer {
	dir := ""
	if mod.IsValid() {
		dir = mod.Root()
	}
	impFrom := packages.NewImporter(fset, dir)
	return &Importer{mod: mod, gopRoot: gop.Root, gop: gop, impFrom: impFrom}
}

func (p *Importer) Import(pkgPath string) (pkg *types.Package, err error) {
	const (
		gop = "github.com/goplus/gop"
	)
	if strings.HasPrefix(pkgPath, gop) {
		if suffix := pkgPath[len(gop):]; suffix == "" || suffix[0] == '/' {
			if suffix == "/cl/internal/gop-in-go/foo" {
				if err = p.genGoExtern(p.gopRoot+suffix, false); err != nil {
					return
				}
			}
			return p.impFrom.ImportFrom(pkgPath, p.gopRoot, 0)
		}
	}
	if mod := p.mod; mod.IsValid() {
		if mod.PkgType(pkgPath) == gopmod.PkgtExtern {
			ret, modVer, e := mod.LookupExternPkg(pkgPath)
			if e != nil {
				return nil, e
			}
			isExtern := modVer.Version != ""
			if isExtern {
				if _, err = modfetch.Get(modVer.String()); err != nil {
					return
				}
			}
			modfile := filepath.Join(ret.ModDir, "gop.mod")
			if _, e := os.Lstat(modfile); e == nil { // has gop.mod
				if err = p.genGoExtern(ret.Dir, isExtern); err != nil {
					return
				}
			}
			return p.impFrom.ImportFrom(pkgPath, ret.ModDir, 0)
		}
	}
	return p.impFrom.Import(pkgPath)
}

func (p *Importer) genGoExtern(dir string, isExtern bool) (err error) {
	genfile := filepath.Join(dir, autoGenFile)
	if _, err = os.Lstat(genfile); err == nil { // has gop_autogen.go
		return
	}
	if isExtern {
		os.Chmod(dir, modWritable)
		defer os.Chmod(dir, modReadonly)
	}
	return genGoIn(dir, &Config{Gop: p.gop, Importer: p}, false, false)
}

// -----------------------------------------------------------------------------
