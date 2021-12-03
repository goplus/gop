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

package cl

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/goplus/gop/env"
	"github.com/goplus/gop/x/mod/modfile"
	"github.com/goplus/gox"
	"golang.org/x/tools/go/packages"
)

// -----------------------------------------------------------------------------

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
	cached     *gox.LoadPkgsCached
	genGoPkg   func(pkgDir string, base *Config) error
	LoadPkgs   gox.LoadPkgsFunc
	BaseConfig *Config
}

func initPkgsLoader(base *Config) {
	if base.ModRootDir == "" {
		modfile, err := env.GOPMOD(base.Dir)
		if err != nil {
			log.Panicln("gop.mod not found:", err)
		}
		base.ModRootDir, _ = filepath.Split(modfile)
	}
	p := &PkgsLoader{genGoPkg: base.GenGoPkg, BaseConfig: base}
	if base.PersistLoadPkgs {
		if base.CacheFile == "" && base.ModRootDir != "" {
			dir := base.ModRootDir + "/.gop"
			os.MkdirAll(dir, 0755)
			base.CacheFile = dir + "/gop.cache"
		}
	}
	if base.CacheFile != "" {
		p.cached = gox.OpenLoadPkgsCached(base.CacheFile, p.Load)
		p.LoadPkgs = p.cached.Load
	} else if base.CacheLoadPkgs {
		p.LoadPkgs = gox.NewLoadPkgsCached(p.Load)
	} else {
		p.LoadPkgs = p.loadPkgsNoCache
	}
	base.PkgsLoader = p
}

func (p *PkgsLoader) Save() error {
	if p.cached == nil {
		return nil
	}
	return p.cached.Save()
}

func (p *PkgsLoader) GenGoPkgs(cfg *packages.Config, notFounds []string) (err error) {
	if p.genGoPkg == nil {
		return syscall.ENOENT
	}
	modfile, err := env.GOPMOD(cfg.Dir)
	if err != nil {
		return
	}
	pkgPath, err := GetModulePath(modfile)
	if err != nil {
		return
	}
	root, _ := filepath.Split(modfile)
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
	var nretry int
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
			if nretry > 1 {
				log.Println("Load packages too many times:", notFounds)
				return loadPkgs, err
			}
			if e := p.GenGoPkgs(cfg, notFounds); e == nil {
				nretry++
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
