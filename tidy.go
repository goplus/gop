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
	"github.com/qiniu/x/errors"
)

func Tidy(dir string, gop *env.Gop) (err error) {
	modObj, err := gopmod.Load(dir, mod.GopModOnly)
	if err != nil {
		return errors.NewWith(err, `gopmod.Load(dir, mod.GopModOnly)`, -2, "gopmod.Load", dir, mod.GopModOnly)
	}

	modRoot := modObj.Root()
	depMods, err := GenDepMods(modObj, modRoot, true)
	if err != nil {
		return errors.NewWith(err, `GenDepMods(modObj, modRoot, true)`, -2, "gop.GenDepMods", modObj, modRoot, true)
	}

	old := modObj.DepMods()
	for modPath := range old {
		if _, ok := depMods[modPath]; !ok { // removed
			modObj.DropRequire(modPath)
		}
	}
	for modPath := range depMods {
		if _, ok := old[modPath]; !ok { // added
			if newMod, e := modfetch.Get(modPath); e != nil {
				return errors.NewWith(e, `modfetch.Get(modPath)`, -1, "modfetch.Get", modPath)
			} else {
				modObj.AddRequire(newMod.Path, newMod.Version)
			}
		}
	}

	modObj.Cleanup()
	err = modObj.Save()
	if err != nil {
		return errors.NewWith(err, `modObj.Save()`, -2, "(*gopmod.Module).Save")
	}

	conf := &Config{DontUpdateGoMod: true, Gop: gop}
	err = genGoDir(modRoot, conf, true, true)
	if err != nil {
		return errors.NewWith(err, `genGoDir(modRoot, conf, true, true)`, -2, "gop.genGoDir", modRoot, conf, true, true)
	}

	err = modObj.UpdateGoMod(gop, true)
	if err != nil {
		return errors.NewWith(err, `modObj.UpdateGoMod(gop, true)`, -2, "(*gopmod.Module).UpdateGoMod", gop, true)
	}

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = modRoot
	err = cmd.Run()
	if err != nil {
		err = errors.NewWith(err, `cmd.Run()`, -2, "(*exec.Cmd).Run")
	}
	return
}
