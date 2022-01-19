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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goplus/gop/x/gengo"
	"github.com/goplus/gop/x/gopprojs"
	"github.com/goplus/gop/x/mod/modfetch"
	"golang.org/x/mod/module"
)

// -----------------------------------------------------------------------------

type handleEvent struct {
	lastErr error
}

func (p *handleEvent) OnStart(pkgPath string) {
	fmt.Fprintln(os.Stderr, pkgPath)
}

func (p *handleEvent) OnInfo(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func (p *handleEvent) OnErr(stage string, err error) {
	p.lastErr = fmt.Errorf("%s: %v", stage, err)
	fmt.Fprintln(os.Stderr, p.lastErr)
}

func (p *handleEvent) OnEnd() {
}

// -----------------------------------------------------------------------------

const (
	FlagGoAsGoPlus = 1 << iota
)

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

func OpenProject(flags int, src gopprojs.Proj) (ctx *Context, proj *Project, err error) {
	switch v := src.(type) {
	case *gopprojs.FilesProj:
		ctx = New(".")
		proj, err = ctx.OpenFiles(flags, v.Files...)
		return
	case *gopprojs.DirProj:
		return OpenDir(flags, v.Dir)
	case *gopprojs.PkgPathProj:
		return OpenPkgPath(flags, v.Path)
	}
	panic("OpenProject: unexpected source")
}

func OpenDir(flags int, dir string) (ctx *Context, proj *Project, err error) {
	err = os.Chdir(dir)
	if err != nil {
		return
	}
	ev := new(handleEvent)
	if !gengo.GenGo(gengo.Config{Event: ev}, false, dir, dir) {
		return nil, nil, ev.lastErr
	}
	ctx = New(dir)
	proj, err = ctx.OpenFiles(0, filepath.Join(dir, "gop_autogen.go"))
	return
}

func OpenPkgPath(flags int, pkgPath string) (ctx *Context, proj *Project, err error) {
	dir := gopGetDefaultPath()
	err = os.Chdir(dir)
	if err != nil {
		return
	}
	modPath, leftPart := splitPkgPath(pkgPath)
	var modVer module.Version
	modVer, _, err = modfetch.Get(modPath)
	if err != nil {
		return
	}
	dir, err = modfetch.ModCachePath(modVer)
	if err != nil {
		return
	}
	pkg, _ := splitVerson(modPath)
	if strings.HasPrefix(pkg, modVer.Path+"/") {
		subdir := filepath.Join(dir, pkg[len(modVer.Path)+1:])
		return OpenDir(flags, subdir)
	}
	return OpenDir(flags, dir+leftPart)
}

func splitVerson(modPath string) (path, version string) {
	pos := strings.IndexByte(modPath, '@')
	if pos > 0 {
		return modPath[:pos], modPath[pos+1:]
	}
	return modPath, ""
}

func splitPkgPath(pkgPath string) (modPath, leftPart string) {
	return pkgPath, ""
}

// -----------------------------------------------------------------------------
