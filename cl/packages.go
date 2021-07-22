/*
 Copyright 2021 The GoPlus Authors (goplus.org)

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package cl

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/goplus/gox"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/packages"
)

// -----------------------------------------------------------------------------

func FindGoModFile(dir string) (file string, err error) {
	if dir == "" {
		dir = "."
	}
	if dir, err = filepath.Abs(dir); err != nil {
		return
	}
	for dir != "" {
		file = filepath.Join(dir, "go.mod")
		if fi, e := os.Lstat(file); e == nil && !fi.IsDir() {
			return
		}
		if dir, file = filepath.Split(strings.TrimRight(dir, "/\\")); file == "" {
			break
		}
	}
	return "", syscall.ENOENT
}

func GetModulePath(file string) (pkgPath string, err error) {
	src, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	f, err := modfile.ParseLax(file, src, nil)
	if err != nil {
		return
	}
	return f.Module.Mod.Path, nil
}

// -----------------------------------------------------------------------------

type PkgsLoader struct {
	genGoPkg   func(pkgDir string, base *Config) error
	LoadPkgs   gox.LoadPkgsFunc
	BaseConfig *Config
}

func initPkgsLoader(base *Config) {
	p := &PkgsLoader{genGoPkg: base.GenGoPkg, BaseConfig: base}
	if base.CacheLoadPkgs {
		p.LoadPkgs = gox.NewLoadPkgsCached(p.Load)
	} else {
		p.LoadPkgs = p.loadPkgsNoCache
	}
	base.PkgsLoader = p
}

func (p *PkgsLoader) GenGoPkgs(cfg *packages.Config, notFounds []string) (err error) {
	if p.genGoPkg == nil {
		return syscall.ENOENT
	}
	file, err := FindGoModFile(cfg.Dir)
	if err != nil {
		return
	}
	root, _ := filepath.Split(file)
	pkgPath, err := GetModulePath(file)
	if err != nil {
		return
	}
	pkgPathSlash := pkgPath + "/"
	for _, notFound := range notFounds {
		if strings.HasPrefix(notFound, pkgPathSlash) || notFound == pkgPath {
			if err = p.genGoPkg(root+notFound[len(pkgPath):], p.BaseConfig); err != nil {
				return
			}
		} else {
			return syscall.ENOENT
		}
	}
	return nil
}

func (p *PkgsLoader) Load(cfg *packages.Config, patterns ...string) ([]*packages.Package, error) {
retry:
	loadPkgs, err := packages.Load(cfg, patterns...)
	if err == nil && p.genGoPkg != nil {
		var notFounds []string
		packages.Visit(loadPkgs, nil, func(pkg *packages.Package) {
			const goGetCmd = "go get "
			for _, err := range pkg.Errors {
				if pos := strings.LastIndex(err.Msg, goGetCmd); pos > 0 {
					notFounds = append(notFounds, err.Msg[pos+len(goGetCmd):])
				}
			}
		})
		if notFounds != nil {
			if e := p.GenGoPkgs(cfg, notFounds); e == nil {
				goto retry
			}
		}
	}
	return loadPkgs, err
}

func (p *PkgsLoader) loadPkgsNoCache(at *gox.Package, imports map[string]*gox.PkgRef, pkgPaths ...string) int {
	conf := at.InternalGetLoadConfig()
	loadPkgs, err := p.Load(conf, pkgPaths...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if n := packages.PrintErrors(loadPkgs); n > 0 {
		return n
	}
	for _, loadPkg := range loadPkgs {
		gox.LoadGoPkg(at, imports, loadPkg)
	}
	return 0
}

// -----------------------------------------------------------------------------
