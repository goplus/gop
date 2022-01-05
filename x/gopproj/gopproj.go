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

	"github.com/goplus/gop/x/gengo"
	"github.com/goplus/gop/x/gopprojs"
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
	ev := new(handleEvent)
	if !gengo.GenGo(gengo.Config{Event: ev}, false, dir) {
		return nil, ev.lastErr
	}
	return p.OpenFiles(0, filepath.Join(dir, "gop_autogen.go"))
}

func (p *Context) OpenPkgPath(flags int, pkgPath string) (proj *Project, err error) {
	panic("todo")
}

// -----------------------------------------------------------------------------
