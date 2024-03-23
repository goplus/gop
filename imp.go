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
	"bytes"
	"fmt"
	"go/token"
	"go/types"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/goplus/gogen/packages"
	"github.com/goplus/mod/env"
	"github.com/goplus/mod/gopmod"
	"github.com/goplus/mod/modfetch"
	"github.com/qiniu/x/errors"
)

// -----------------------------------------------------------------------------

type Importer struct {
	impFrom *packages.Importer
	mod     *gopmod.Module
	gop     *env.Gop
	fset    *token.FileSet

	Flags GenFlags // can change this for loading Go+ modules
}

func NewImporter(mod *gopmod.Module, gop *env.Gop, fset *token.FileSet) *Importer {
	const (
		defaultFlags = GenFlagPrompt | GenFlagPrintError
	)
	if mod == nil {
		mod = gopmod.Default
	}
	dir := ""
	if mod.HasModfile() {
		dir = mod.Root()
	}
	impFrom := packages.NewImporter(fset, dir)
	// is gop root? load gop pkg cache
	if len(gop.Root) > 0 {
		impFrom.GetPkgCache().GoListExportCacheSync(gop.Root,
			"github.com/goplus/gop/builtin",
			"github.com/goplus/gop/builtin/ng",
			"github.com/goplus/gop/builtin/iox",
		)
	}
	return &Importer{mod: mod, gop: gop, impFrom: impFrom, fset: fset, Flags: defaultFlags}
}

const (
	gopMod = "github.com/goplus/gop"
)

// Import imports a Go/Go+ package.
func (p *Importer) Import(pkgPath string) (pkg *types.Package, err error) {
	if strings.HasPrefix(pkgPath, gopMod) {
		if suffix := pkgPath[len(gopMod):]; suffix == "" || suffix[0] == '/' {
			gopRoot := p.gop.Root
			if suffix == "/cl/internal/gop-in-go/foo" { // for test github.com/goplus/gop/cl
				if err = p.genGoExtern(gopRoot+suffix, false); err != nil {
					return
				}
			}
			return p.impFrom.ImportFrom(pkgPath, gopRoot, 0)
		}
	}
	if isPkgInMod(pkgPath, "github.com/qiniu/x") {
		return p.impFrom.ImportFrom(pkgPath, p.gop.Root, 0)
	}
	if mod := p.mod; mod.HasModfile() {
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
			modDir := ret.ModDir
			goModfile := filepath.Join(modDir, "go.mod")
			if _, e := os.Lstat(goModfile); e != nil { // no go.mod
				os.Chmod(modDir, modWritable)
				defer os.Chmod(modDir, modReadonly)
				os.WriteFile(goModfile, defaultGoMod(ret.ModPath), 0644)
			}
			p.checkGopPackage(pkgPath)
			return p.impFrom.ImportFrom(pkgPath, ret.ModDir, 0)
		case gopmod.PkgtModule, gopmod.PkgtLocal:
			if err = p.genGoExtern(ret.Dir, false); err != nil {
				return
			}
		case gopmod.PkgtStandard:
			return p.impFrom.ImportFrom(pkgPath, p.gop.Root, 0)
		}
	}
	p.checkGopPackage(pkgPath)
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
		err = genGoIn(dir, &Config{Gop: p.gop, Importer: p, Fset: p.fset}, false, p.Flags, &gen)
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

// check import pkg is a go+ package
func (p *Importer) checkGopPackage(pkgPath string) {
	cacheInfo, ok := p.impFrom.GetPkgCache().GetPkgCache(pkgPath)
	if !ok {
		return
	}
	// is go+ package with package code change?
	if checkPkgGopSourcesStatus(cacheInfo.PkgDir, cacheInfo.PkgExport) {
		return
	}
	// rebuild package
	reBuildGopPackage(cacheInfo.PkgDir)
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

// check whether source files of the go+ package has changed
func checkPkgGopSourcesStatus(dir, pkgFile string) bool {
	// no need to check dir because dir use go list to get
	pkgInfo, err := os.Stat(pkgFile)
	if err != nil {
		return false
	}
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		// filter not .gxx files,because gop file ext is .gxx
		fileExt := filepath.Ext(info.Name())
		if fileExt == "" || !(fileExt != ".go" && fileExt[1] == 'g') {
			return nil
		}
		if info.ModTime().After(pkgInfo.ModTime()) {
			return fmt.Errorf("package need compile, dir:%s, file:%s", dir, info.Name())
		}
		return nil
	})
	if err != nil {
		log.Println(err.Error())
	}
	return err == nil
}

// https://github.com/goplus/gop/issues/1821
// exec gop build to generate go+ package
func reBuildGopPackage(dir string) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("gop", "build", dir)
	cmd.Dir = dir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		// exit compile
		panic(errors.New(stderr.String()))
	}
}

// -----------------------------------------------------------------------------
