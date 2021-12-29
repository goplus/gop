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
	"os"
	"path/filepath"
	"strings"

	"github.com/goplus/gop/token"
	"github.com/goplus/gop/x/mod/modfile"
	"github.com/goplus/gop/x/mod/modload"
	"golang.org/x/mod/module"
)

// -----------------------------------------------------------------------------

type Module struct {
	modload.Module
	classes map[string]*modfile.Classfile
	vers    map[string]module.Version
	gengo   func(act string, mod module.Version)
	fset    *token.FileSet
}

func New(mod modload.Module) *Module {
	classes := make(map[string]*modfile.Classfile)
	vers := mod.Deps()
	fset := token.NewFileSet()
	return &Module{classes: classes, vers: vers, Module: mod, fset: fset}
}

func Load(dir string) (*Module, error) {
	mod, err := modload.Load(dir)
	if err != nil {
		return nil, err
	}
	return New(mod), nil
}

func (p *Module) SetGenGo(gengo func(act string, mod module.Version)) {
	p.gengo = gengo
}

func (p *Module) Imports(dir string) (imps []string, err error) {
	imports := make(map[string]struct{})
	err = p.parseImports(imports, dir)
	imps = getKeys(imports)
	return
}

func getKeys(v map[string]struct{}) []string {
	keys := make([]string, 0, len(v))
	for key := range v {
		keys = append(keys, key)
	}
	return keys
}

func (p *Module) parseImports(imports map[string]struct{}, dir string) (err error) {
	list, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	var errs ErrorList
	for _, d := range list {
		if d.IsDir() {
			continue
		}
		fname := d.Name()
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
