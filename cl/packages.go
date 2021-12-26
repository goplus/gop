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
	"io/ioutil"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/goplus/gop/env"
	"github.com/goplus/gop/x/mod/modfile"
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
	genGoPkg   func(pkgDir string, base *Config) error
	BaseConfig *Config
}

func initPkgsLoader(base *Config) {
	p := &PkgsLoader{genGoPkg: base.GenGoPkg, BaseConfig: base}
	base.PkgsLoader = p
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

// -----------------------------------------------------------------------------
