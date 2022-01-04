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
	"path"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/goplus/gop/token"
	"github.com/goplus/gop/x/mod/modfetch"
	"github.com/goplus/gop/x/mod/modload"
	"golang.org/x/mod/module"
)

func Imports(dir string) (imps []string, err error) {
	dir, recursive := parsePattern(dir)
	mod, err := Load(dir)
	if err != nil {
		return
	}
	err = mod.RegisterClasses()
	if err != nil {
		return
	}
	imports := make(map[string]none)
	err = mod.Imports(imports, dir, recursive)
	if err != nil {
		return
	}
	return getSortedKeys(imports), nil
}

func List(pattern ...string) (pkgPaths []string, err error) {
	mod, err := Load(".")
	if err != nil {
		return
	}
	err = mod.RegisterClasses()
	if err != nil {
		return
	}
	ret := make(map[string]none)
	err = mod.List(ret, pattern...)
	if err != nil {
		return
	}
	return getSortedKeys(ret), nil
}

func getSortedKeys(v map[string]none) []string {
	keys := make([]string, 0, len(v))
	for key := range v {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// -----------------------------------------------------------------------------

type depmodInfo struct {
	path string
	real module.Version
}

type Module struct {
	modload.Module
	classes map[string]*Class
	depmods []depmodInfo
	fset    *token.FileSet
}

// PkgType specifies a package type.
type PkgType int

const (
	PkgtStandard PkgType = iota // a standard Go/Go+ package
	PkgtModule                  // a package in this module (in standard form)
	PkgtLocal                   // a package in this module (in relative path form)
	PkgtExtern                  // an extarnal package
	PkgtInvalid  = -1           // an invalid package
)

// PkgType returns the package type of specified package.
func (p *Module) PkgType(pkgPath string) PkgType {
	if pkgPath == "" {
		return PkgtInvalid
	}
	if isPkgInMod(pkgPath, p.Path()) {
		return PkgtModule
	}
	if pkgPath[0] == '.' {
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

// LookupExternPkg lookups a external package from depended modules.
// If modVer.Path is replace to be a local path, it will be canonical to an absolute path.
func (p *Module) LookupExternPkg(pkgPath string) (modPath string, modVer module.Version, ok bool) {
	for _, m := range p.depmods {
		if isPkgInMod(pkgPath, m.path) {
			modPath, modVer, ok = m.path, m.real, true
			break
		}
	}
	return
}

// LookupMod lookups a depended module.
// If modVer.Path is replace to be a local path, it will be canonical to an absolute path.
func (p *Module) LookupMod(modPath string) (modVer module.Version, ok bool) {
	for _, m := range p.depmods {
		if m.path == modPath {
			modVer, ok = m.real, true
			break
		}
	}
	return
}

func getDepMods(mod modload.Module) []depmodInfo {
	depmods := mod.DepMods()
	ret := make([]depmodInfo, 0, len(depmods))
	for path, m := range depmods {
		ret = append(ret, depmodInfo{path: path, real: m})
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].path > ret[j].path
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
	mod, _, err = modfetch.Get(mod.String())
	if err != nil && err != syscall.EEXIST {
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

func (p *Module) List(ret map[string]struct{}, pattern ...string) (err error) {
	modRoot, modPath := p.Root(), p.Path()
	for _, pat := range pattern {
		if err = p.listPkgs(ret, pat, modRoot, modPath); err != nil {
			return
		}
	}
	return
}

func (p *Module) listPkgs(ret map[string]none, pat, modRoot, modPath string) error {
	pat, recursive := parsePattern(pat)
	if isLocal(pat) {
		patAbs, err1 := filepath.Abs(pat)
		patRel, err2 := filepath.Rel(modRoot, patAbs)
		if err1 != nil || err2 != nil || strings.HasPrefix(patRel, "..") {
			return fmt.Errorf("directory `%s` outside available modules", pat)
		}
		pkgPathBase := path.Join(modPath, filepath.ToSlash(patRel))
		return p.doListLocalPkgs(ret, pkgPathBase, pat, recursive)
	}
	return p.doListPkgs(ret, pat, recursive)
}

func parsePattern(pat string) (string, bool) {
	const multi = "/..."
	recursive := strings.HasSuffix(pat, multi)
	if recursive {
		pat = pat[:len(pat)-len(multi)]
		if pat == "" {
			pat = "/"
		}
	}
	return pat, recursive
}

func isLocal(ns string) bool {
	if len(ns) > 0 {
		switch c := ns[0]; c {
		case '/', '\\', '.':
			return true
		default:
			return len(ns) >= 2 && ns[1] == ':' && ('A' <= c && c <= 'Z' || 'a' <= c && c <= 'z')
		}
	}
	return false
}

func (p *Module) doListLocalPkgs(ret map[string]none, pkgPathBase, pat string, recursive bool) (err error) {
	fis, err := os.ReadDir(pat)
	if err != nil {
		return
	}
	noSouceFile := true
	for _, fi := range fis {
		name := fi.Name()
		if strings.HasPrefix(name, "_") {
			continue
		}
		if fi.IsDir() {
			if recursive {
				err = p.doListLocalPkgs(ret, pkgPathBase+"/"+name, pat+"/"+name, true)
				if err != nil {
					return
				}
			}
		} else if noSouceFile {
			ext := path.Ext(name)
			switch ext {
			case ".gop", ".go":
				noSouceFile = false
			default:
				if _, ok := p.classes[ext]; ok {
					noSouceFile = false
				}
			}
		}
	}
	if !noSouceFile {
		ret[pkgPathBase] = none{}
	}
	return
}

func (p *Module) doListPkgs(ret map[string]none, pkgPath string, recursive bool) error {
	switch p.PkgType(pkgPath) {
	case PkgtStandard:
		if !recursive {
			ret[pkgPath] = none{}
			return nil
		}
	case PkgtModule:
		relDir := p.Root() + pkgPath[len(p.Path()):]
		return p.doListLocalPkgs(ret, pkgPath, relDir, recursive)
	}
	modPath, modVer, ok := p.LookupExternPkg(pkgPath)
	if !ok {
		panic("TODO: external package not found")
	}
	mod, err := LoadMod(modVer)
	if err != nil {
		return err
	}
	relDir := mod.Root() + pkgPath[len(modPath):]
	return mod.doListLocalPkgs(ret, pkgPath, relDir, recursive)
}

func (p *Module) Imports(imports map[string]struct{}, dir string, recursive bool) (err error) {
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
				if err = p.Imports(imports, filepath.Join(dir, fname), true); err != nil {
					return
				}
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
					imports[pkgPath] = none{}
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

type ChangeInfo struct {
	ModTime    time.Time
	GopFileNum int
	SourceNum  int
}

func (p *Module) ChangeInfo(dir string) (ci ChangeInfo, err error) {
	list, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	dir += "/"
	var fi os.FileInfo
	for _, d := range list {
		fname := d.Name()
		if strings.HasPrefix(fname, "_") || d.IsDir() {
			continue
		}
		ext := filepath.Ext(fname)
		switch ext {
		case ".gop":
			ci.GopFileNum++
		case ".go":
			if strings.HasPrefix(fname, "gop_autogen") {
				continue
			}
		default:
			if _, ok := p.classes[ext]; !ok {
				continue
			}
			ci.GopFileNum++
		}
		ci.SourceNum++
		fi, err = os.Stat(dir + fname)
		if err != nil {
			return
		}
		t := fi.ModTime()
		if t.After(ci.ModTime) {
			ci.ModTime = t
		}
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
