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

package gop

import (
	"os"
	"os/exec"

	"github.com/goplus/mod"
	"github.com/goplus/mod/env"
	"github.com/goplus/mod/gopmod"
	"github.com/goplus/mod/modfetch"
)

func Tidy(dir string, gop *env.Gop) (err error) {
	mod, err := gopmod.Load(dir, mod.GopModOnly)
	if err != nil {
		return
	}

	depMods, err := GenDepMods(mod, mod.Root(), true)
	if err != nil {
		return
	}

	old := mod.DepMods()
	for modPath := range old {
		if _, ok := depMods[modPath]; !ok { // removed
			mod.DropRequire(modPath)
		}
	}
	for modPath := range depMods {
		if _, ok := old[modPath]; !ok { // added
			if newMod, e := modfetch.Get(modPath); e != nil {
				return e
			} else {
				mod.AddRequire(newMod.Path, newMod.Version)
			}
		}
	}

	mod.Cleanup()
	err = mod.Save()
	if err != nil {
		return
	}

	err = genGoDir(mod.Root(), &Config{DontUpdateGoMod: true, Gop: gop}, true, true)
	if err != nil {
		return
	}

	err = mod.UpdateGoMod(gop, true)
	if err != nil {
		return
	}

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = mod.Root()
	return cmd.Run()
}
