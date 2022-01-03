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
	"errors"
	"syscall"

	"github.com/goplus/gop/x/mod/modfetch"
	"github.com/goplus/gop/x/mod/modfile"
	"github.com/goplus/gop/x/mod/modload"
	"golang.org/x/mod/module"
)

type Class = modfile.Classfile

var (
	ClassSpx = &Class{
		ProjExt:  ".gmx",
		WorkExt:  ".spx",
		PkgPaths: []string{"github.com/goplus/spx", "math"},
	}
)

var (
	ErrNotClassFileMod = errors.New("not a classfile module")
)

// -----------------------------------------------------------------------------

func (p *Module) IsClass(ext string) bool {
	_, ok := p.classes[ext]
	return ok
}

func (p *Module) LookupClass(ext string) (c *Class, ok bool) {
	c, ok = p.classes[ext]
	return
}

func (p *Module) RegisterClasses(registerClass ...func(c *Class)) (err error) {
	var regcls func(c *Class)
	if registerClass != nil {
		regcls = registerClass[0]
	}
	p.registerClass(ClassSpx, regcls)
	if c := p.Classfile; c != nil {
		p.registerClass(c, regcls)
	}
	for _, r := range p.Register {
		if err = p.registerMod(r.ClassfileMod, regcls); err != nil {
			return
		}
	}
	return
}

func (p *Module) registerMod(modPath string, regcls func(c *Class)) (err error) {
	mod, ok := p.LookupMod(modPath)
	if !ok {
		return syscall.ENOENT
	}
	err = p.registerClassFrom(mod, regcls)
	if err != syscall.ENOENT {
		return
	}
	mod, _, err = modfetch.Get(mod.String())
	if err != nil && err != syscall.EEXIST {
		return
	}
	return p.registerClassFrom(mod, regcls)
}

func (p *Module) registerClassFrom(modVer module.Version, regcls func(c *Class)) (err error) {
	dir, err := modfetch.ModCachePath(modVer)
	if err != nil {
		return
	}
	mod, err := modload.Load(dir)
	if err != nil {
		return
	}
	c := mod.Classfile
	if c == nil {
		return ErrNotClassFileMod
	}
	p.registerClass(c, regcls)
	return
}

func (p *Module) registerClass(c *Class, regcls func(c *Class)) {
	p.classes[c.ProjExt] = c
	if c.WorkExt != "" {
		p.classes[c.WorkExt] = c
	}
	if regcls != nil {
		regcls(c)
	}
}

// -----------------------------------------------------------------------------
