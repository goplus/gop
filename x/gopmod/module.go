/*
 * Copyright (c) 2021 The GoPlus Authors (goplus.org). All rights reserved.
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

package gopmod

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	"github.com/goplus/gop/token"
	"github.com/goplus/gop/x/mod/modfetch"
	"github.com/goplus/gop/x/mod/modload"
	"golang.org/x/mod/module"
)

func Imports(dir string) (imps []string, err error) {
	var recursive bool
	if strings.HasSuffix(dir, "/...") {
		dir, recursive = dir[:len(dir)-4], true
	}
	mod, err := Load(dir)
	if err != nil {
		return
	}
	err = mod.RegisterClasses()
	if err != nil {
		return
	}
	return mod.Imports(dir, recursive)
}

// -----------------------------------------------------------------------------

type Module struct {
	modload.Module
	classes map[string]*Class
	depmods []module.Version
	fset    *token.FileSet
}

type PkgType int

const (
	PkgtStandard PkgType = iota
	PkgtModule
	PkgtLocal
	PkgtExtern
	PkgtInvalid = -1
)

func (p *Module) PkgType(pkgPath string) PkgType {
	if pkgPath == "" {
		return PkgtInvalid
	}
	if isPkgInMod(pkgPath, p.Path()) {
		return PkgtModule
	}
	c := pkgPath[0]
	if c == '/' || c == '.' {
		return PkgtLocal
	}
	pos := strings.Index(pkgPath, "/")
	if pos > 0 {
		pkgPath = pkgPath[:pos]
	}
	if strings.Contains(pkgPath, ".") {
		return PkgtExtern
	}
	return PkgtStandard
}

func isPkgInMod(pkgPath, modPath string) bool {
	if strings.HasPrefix(pkgPath, modPath) {
		suffix := pkgPath[len(modPath):]
		return suffix == "" || suffix[0] == '/'
	}
	return false
}

func (p *Module) LookupExternPkg(pkgPath string) (m module.Version, relDir string, ok bool) {
	for _, m = range p.depmods {
		if isPkgInMod(pkgPath, m.Path) {
			relDir, ok = "."+pkgPath[len(m.Path):], true
			break
		}
	}
	return
}

func (p *Module) LookupMod(modPath string) (m module.Version, ok bool) {
	for _, m = range p.depmods {
		if m.Path == modPath {
			ok = true
			break
		}
	}
	return
}

func getDepMods(mod modload.Module) []module.Version {
	depmods := mod.DepMods()
	ret := make([]module.Version, 0, len(depmods))
	for _, m := range depmods {
		ret = append(ret, m)
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Path > ret[j].Path
	})
	return ret
}

func New(mod modload.Module) *Module {
	classes := make(map[string]*Class)
	depmods := getDepMods(mod)
	fset := token.NewFileSet()
	return &Module{classes: classes, depmods: depmods, Module: mod, fset: fset}
}

func Load(dir string) (*Module, error) {
	mod, err := modload.Load(dir)
	if err != nil {
		return nil, err
	}
	return New(mod), nil
}

func LoadMod(mod module.Version) (p *Module, err error) {
	p, err = loadModFrom(mod)
	if err != syscall.ENOENT {
		return
	}
	mod, err = modfetch.Get(mod.String())
	if err != nil {
		return
	}
	return loadModFrom(mod)
}

func loadModFrom(mod module.Version) (p *Module, err error) {
	dir, err := modfetch.ModCachePath(mod)
	if err != nil {
		return
	}
	return Load(dir)
}

func (p *Module) Imports(dir string, recursive bool) (imps []string, err error) {
	if !isLocal(dir) {
		return nil, fmt.Errorf("`%s` is not a local directory", dir)
	}
	imports := make(map[string]struct{})
	err = p.parseImports(imports, dir, recursive)
	imps = getKeys(imports)
	return
}

func isLocal(modPath string) bool {
	return strings.HasPrefix(modPath, ".") || strings.HasPrefix(modPath, "/")
}

func getKeys(v map[string]struct{}) []string {
	keys := make([]string, 0, len(v))
	for key := range v {
		keys = append(keys, key)
	}
	return keys
}

func (p *Module) parseImports(imports map[string]struct{}, dir string, recursive bool) (err error) {
	list, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	var errs ErrorList
	for _, d := range list {
		fname := d.Name()
		if strings.HasPrefix(fname, "_") { // skip this file/directory
			continue
		}
		if d.IsDir() {
			if recursive {
				p.parseImports(imports, filepath.Join(dir, fname), true)
			}
			continue
		}
		ext := filepath.Ext(fname)
		switch ext {
		case ".gop":
			p.parseGopImport(&errs, imports, filepath.Join(dir, fname))
		case ".go":
			if !strings.HasPrefix(fname, "gop_autogen") {
				p.parseGoImport(&errs, imports, filepath.Join(dir, fname))
			}
		default:
			if c, ok := p.classes[ext]; ok {
				for _, pkgPath := range c.PkgPaths {
					imports[pkgPath] = struct{}{}
				}
				p.parseGopImport(&errs, imports, filepath.Join(dir, fname))
			}
		}
	}
	if len(errs) > 0 {
		err = errs
	}
	return
}

// -----------------------------------------------------------------------------

type ErrorList []error

func (e ErrorList) Error() string {
	errStrs := make([]string, len(e))
	for i, err := range e {
		errStrs[i] = err.Error()
	}
	return strings.Join(errStrs, "\n")
}

// -----------------------------------------------------------------------------
