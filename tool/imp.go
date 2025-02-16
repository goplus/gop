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

package tool

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"go/token"
	"go/types"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/goplus/gogen/packages"
	"github.com/goplus/gogen/packages/cache"
	"github.com/goplus/mod/env"
	"github.com/goplus/mod/gopmod"
	"github.com/goplus/mod/modfetch"
	"github.com/goplus/mod/modfile"
)

// -----------------------------------------------------------------------------

// Importer represents a Go+ importer.
type Importer struct {
	impFrom *packages.Importer
	mod     *gopmod.Module
	gop     *env.Gop
	fset    *token.FileSet

	Flags GenFlags // can change this for loading Go+ modules
}

// NewImporter creates a Go+ Importer.
func NewImporter(mod *gopmod.Module, gop *env.Gop, fset *token.FileSet) *Importer {
	const (
		defaultFlags = GenFlagPrompt | GenFlagPrintError
	)
	if mod == nil || !mod.HasModfile() {
		if modGop, e := gopmod.LoadFrom(filepath.Join(gop.Root, "go.mod"), ""); e == nil {
			modGop.ImportClasses()
			mod = modGop
		} else {
			mod = gopmod.Default
		}
	}
	dir := mod.Root()
	impFrom := packages.NewImporter(fset, dir)
	ret := &Importer{mod: mod, gop: gop, impFrom: impFrom, fset: fset, Flags: defaultFlags}
	impFrom.SetCache(cache.New(ret.PkgHash))
	return ret
}

func (p *Importer) SetTags(tags string) {
	p.impFrom.SetTags(tags)
	if c, ok := p.impFrom.Cache().(*cache.Impl); ok {
		c.SetTags(tags)
	}
}

// CacheFile returns file path of the cache.
func (p *Importer) CacheFile() string {
	cacheDir, _ := os.UserCacheDir()
	cacheDir += "/gop-build/"
	os.MkdirAll(cacheDir, 0755)

	fname := ""
	h := sha256.New()
	if root := p.mod.Root(); root != "" {
		io.WriteString(h, root)
		fname = filepath.Base(root)
	}
	if tags := p.impFrom.Tags(); tags != "" {
		io.WriteString(h, tags)
	}
	hash := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	return cacheDir + hash + fname
}

// Cache returns the cache object.
func (p *Importer) Cache() *cache.Impl {
	return p.impFrom.Cache().(*cache.Impl)
}

// PkgHash calculates hash value for a package.
// It is required by cache.New func.
func (p *Importer) PkgHash(pkgPath string, self bool) string {
	if pkg, e := p.mod.Lookup(pkgPath); e == nil {
		switch pkg.Type {
		case gopmod.PkgtStandard:
			return cache.HashSkip
		case gopmod.PkgtExtern:
			if pkg.Real.Version != "" {
				return pkg.Real.String()
			}
			fallthrough
		case gopmod.PkgtModule:
			return dirHash(p.mod, p.gop, pkg.Dir, self)
		}
	}
	if isPkgInMod(pkgPath, gopMod) {
		return cache.HashSkip
	}
	log.Println("PkgHash: unexpected package -", pkgPath)
	return cache.HashInvalid
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

func dirHash(mod *gopmod.Module, gop *env.Gop, dir string, self bool) string {
	h := sha256.New()
	if self {
		fmt.Fprintf(h, "go\t%s\n", runtime.Version())
		fmt.Fprintf(h, "gop\t%s\n", gop.Version)
	}
	if fis, err := os.ReadDir(dir); err == nil {
		for _, fi := range fis {
			if fi.IsDir() {
				continue
			}
			fname := fi.Name()
			if strings.HasPrefix(fname, "_") || !canCl(mod, fname) {
				continue
			}
			if v, e := fi.Info(); e == nil {
				fmt.Fprintf(h, "file\t%s\t%x\t%x\n", fname, v.Size(), v.ModTime().UnixNano())
			}
		}
	}
	return base64.RawStdEncoding.EncodeToString(h.Sum(nil))
}

func canCl(mod *gopmod.Module, fname string) bool {
	switch path.Ext(fname) {
	case ".go", ".gop", ".gox":
		return true
	default:
		ext := modfile.ClassExt(fname)
		return mod.IsClass(ext)
	}
}

// -----------------------------------------------------------------------------
