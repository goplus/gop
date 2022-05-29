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

package mod

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/goplus/gop"
	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/x/gopenv"
	"github.com/goplus/mod"
	"github.com/goplus/mod/gopmod"
	"github.com/goplus/mod/modfetch"
)

// gop mod tidy
var cmdTidy = &base.Command{
	UsageLine: "gop mod tidy [-e -v]",
	Short:     "add missing and remove unused modules",
}

func init() {
	cmdTidy.Run = runTidy
}

func runTidy(cmd *base.Command, args []string) {
	mod, err := gopmod.Load(".", mod.GopModOnly)
	if err != nil {
		if err == syscall.ENOENT {
			fmt.Fprintln(os.Stderr, "gop.mod not found")
			os.Exit(1)
		}
		log.Panicln(err)
	}

	depMods, err := gop.GenDepMods(mod, mod.Root(), true)
	check(err)

	tidy(mod, depMods)
}

func tidy(mod *gopmod.Module, depMods map[string]struct{}) {
	old := mod.DepMods()
	for modPath := range old {
		if _, ok := depMods[modPath]; !ok { // removed
			mod.DropRequire(modPath)
		}
	}
	for modPath := range depMods {
		if _, ok := old[modPath]; !ok { // added
			if newMod, err := modfetch.Get(modPath); err != nil {
				log.Fatalln(err)
			} else {
				mod.AddRequire(newMod.Path, newMod.Version)
			}
		}
	}

	err := mod.Save()
	check(err)

	_, _, err = gop.GenGo(mod.Root()+"/...", &gop.Config{DontUpdateGoMod: true})
	check(err)

	err = mod.UpdateGoMod(gopenv.Get(), true)
	check(err)

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = mod.Root()
	err = cmd.Run()
	check(err)
}
