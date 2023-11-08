/*
 * Copyright (c) 2023 The GoPlus Authors (goplus.org). All rights reserved.
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

package typesutil

import (
	"go/types"
	"syscall"

	"github.com/goplus/gop"
	"github.com/goplus/gop/token"
	"github.com/goplus/gop/x/gopenv"
	"github.com/goplus/mod/env"
	"github.com/goplus/mod/gopmod"
)

// -----------------------------------------------------------------------------

type nilImporter struct{}

func (p nilImporter) Import(path string) (ret *types.Package, err error) {
	return nil, syscall.ENOENT
}

// -----------------------------------------------------------------------------

type importer struct {
	imp types.Importer
	alt *gop.Importer

	mod  *gopmod.Module
	gop  *env.Gop
	fset *token.FileSet
}

func newImporter(imp types.Importer, mod *gopmod.Module, gop *env.Gop, fset *token.FileSet) types.Importer {
	if imp == nil {
		imp = nilImporter{}
	}
	if gop == nil {
		gop = gopenv.Get()
	}
	return &importer{imp, nil, mod, gop, fset}
}

func (p *importer) Import(path string) (ret *types.Package, err error) {
	ret, err = p.imp.Import(path)
	if err == nil {
		return
	}
	if p.alt == nil {
		p.alt = gop.NewImporter(p.mod, p.gop, p.fset)
	}
	return p.alt.Import(path)
}

// -----------------------------------------------------------------------------
