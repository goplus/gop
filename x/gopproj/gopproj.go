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

package gopproj

import (
	"path/filepath"

	"github.com/goplus/gop/x/gopprojs"
)

// -----------------------------------------------------------------------------

const (
	FlagGoAsGoPlus = 1 << iota
)

func (p *Context) OpenProject(flags int, src gopprojs.Proj) (proj *Project, err error) {
	switch v := src.(type) {
	case *gopprojs.FilesProj:
		return p.OpenFiles(flags, v.Files...)
	case *gopprojs.DirProj:
		return p.OpenDir(flags, v.Dir)
	case *gopprojs.PkgPathProj:
		return p.OpenPkgPath(flags, v.Path)
	}
	panic("OpenProject: unexpected source")
}

func (p *Context) OpenFiles(flags int, args ...string) (proj *Project, err error) {
	if len(args) != 1 {
		return p.openFromGopFiles(args)
	}
	src := args[0]
	if (flags&FlagGoAsGoPlus) == 0 && filepath.Ext(src) == ".go" {
		return openFromGoFile(src)
	}
	return p.openFromGopFiles(args)
}

func (p *Context) OpenDir(flags int, dir string) (proj *Project, err error) {
	panic("todo")
}

func (p *Context) OpenPkgPath(flags int, pkgPath string) (proj *Project, err error) {
	panic("todo")
}

// -----------------------------------------------------------------------------
