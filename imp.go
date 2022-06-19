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
	"os/exec"
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
	fset    *token.FileSet
}

func NewImporter(mod *gopmod.Module, gop *env.Gop, fset *token.FileSet) *Importer {
	dir := ""
	if mod.IsValid() {
		dir = mod.Root()
	}
	impFrom := packages.NewImporter(fset, dir)
	return &Importer{mod: mod, gop: gop, impFrom: impFrom, fset: fset}
}

func (p *Importer) Import(pkgPath string) (pkg *types.Package, err error) {
	const (
		gop = "github.com/goplus/gop"
	)
	if strings.HasPrefix(pkgPath, gop) {
		if suffix := pkgPath[len(gop):]; suffix == "" || suffix[0] == '/' {
			gopRoot := p.gop.Root
			if suffix == "/cl/internal/gop-in-go/foo" {
				if err = p.genGoExtern(gopRoot+suffix, false); err != nil {
					return
				}
			}
			return p.impFrom.ImportFrom(pkgPath, gopRoot, 0)
		}
	}
	if mod := p.mod; mod.IsValid() {
		ret, e := mod.Lookup(pkgPath)
		if e != nil {
			return nil, e
		}
		switch ret.Type {
		case gopmod.PkgtExtern:
			isExtern := ret.Real.Version != ""
			if isExtern {
				if _, err = modfetch.Get(ret.Real.String()); err != nil {
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
		case gopmod.PkgtModule, gopmod.PkgtLocal:
			if err = p.genGoExtern(ret.Dir, false); err != nil {
				return
			}
		case gopmod.PkgtStandard:
			return p.impFrom.ImportFrom(pkgPath, p.gop.Root, 0)
		}
	}
	return p.impFrom.Import(pkgPath)
}

func (p *Importer) genGoExtern(dir string, isExtern bool) (err error) {
	gosum := filepath.Join(dir, "go.sum")
	if _, err = os.Lstat(gosum); err == nil { // has go.sum
		return
	}

	if isExtern {
		os.Chmod(dir, modWritable)
		defer os.Chmod(dir, modReadonly)
	}

	genfile := filepath.Join(dir, autoGenFile)
	if _, err = os.Lstat(genfile); err != nil { // has gop_autogen.go
		err = genGoIn(dir, &Config{Gop: p.gop, Importer: p, Fset: p.fset}, false, false)
		if err != nil {
			return
		}
	}

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	err = cmd.Run()
	return
}

// -----------------------------------------------------------------------------
