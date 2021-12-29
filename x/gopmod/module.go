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
	"golang.org/x/mod/module"

	"github.com/goplus/gop/token"
	"github.com/goplus/gop/x/mod/modfile"
	"github.com/goplus/gop/x/mod/modload"
)

// -----------------------------------------------------------------------------

type Module struct {
	modload.Module
	imports map[string]struct{}
	classes map[string]*modfile.Classfile
	gengo   func(act string, mod module.Version)
	fset    *token.FileSet
}

func New(mod modload.Module) *Module {
	imps := make(map[string]struct{})
	classes := make(map[string]*modfile.Classfile)
	return &Module{imports: imps, classes: classes, Module: mod, fset: token.NewFileSet()}
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

func (p *Module) Imports() []string {
	return getKeys(p.imports)
}

func getKeys(v map[string]struct{}) []string {
	keys := make([]string, 0, len(v))
	for key := range v {
		keys = append(keys, key)
	}
	return keys
}

// -----------------------------------------------------------------------------
