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
	flags   GenFlags
}

func NewImporter(mod *gopmod.Module, gop *env.Gop, fset *token.FileSet) *Importer {
	const (
		defaultFlags = GenFlagPrompt | GenFlagPrintError
	)
	if mod == nil {
		mod = gopmod.Default
	}
	dir := ""
	if hasModfile(mod) {
		dir = mod.Root()
	}
	impFrom := packages.NewImporter(fset, dir)
	return &Importer{mod: mod, gop: gop, impFrom: impFrom, fset: fset, flags: defaultFlags}
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
	if mod := p.mod; hasModfile(mod) {
		ret, e := mod.Lookup(pkgPath)
		if e != nil {
			if isPkgInMod(pkgPath, "github.com/qiniu/x") {
				return p.impFrom.ImportFrom(pkgPath, p.gop.Root, 0)
			}
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
			modDir := ret.ModDir
			goModfile := filepath.Join(modDir, "go.mod")
			if _, e := os.Lstat(goModfile); e != nil { // no go.mod
				os.Chmod(modDir, modWritable)
				defer os.Chmod(modDir, modReadonly)
				os.WriteFile(goModfile, defaultGoMod(ret.ModPath), 0644)
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
	genfile := filepath.Join(dir, autoGenFile)
	if _, err = os.Lstat(genfile); err != nil { // no gop_autogen.go
		if isExtern {
			os.Chmod(dir, modWritable)
			defer os.Chmod(dir, modReadonly)
		}
		gen := false
		err = genGoIn(dir, &Config{Gop: p.gop, Importer: p, Fset: p.fset}, false, p.flags, &gen)
		if err != nil {
			return
		}
		if gen {
			cmd := exec.Command("go", "mod", "tidy")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Dir = dir
			err = cmd.Run()
		}
	}
	return
}

func isPkgInMod(pkgPath, modPath string) bool {
	if strings.HasPrefix(pkgPath, modPath) {
		suffix := pkgPath[len(modPath):]
		return suffix == "" || suffix[0] == '/'
	}
	return false
}

func defaultGoMod(modPath string) []byte {
	return []byte(`module ` + modPath + `

go 1.16
`)
}

// -----------------------------------------------------------------------------
