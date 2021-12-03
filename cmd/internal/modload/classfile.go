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
package modload

import (
	"log"
	"path/filepath"

	"golang.org/x/mod/module"

	"github.com/goplus/gop/x/mod/modfetch"
	"github.com/goplus/gop/x/mod/modfile"
)

func LoadClassFile() {
	if modFile.Register == nil {
		return
	}

	var dir string
	var err error
	var claassMod module.Version

	for _, require := range modFile.Require {
		if require.Mod.Path == modFile.Register.ClassfileMod {
			claassMod = require.Mod
		}
	}
	for _, replace := range modFile.Replace {
		if replace.Old.Path == modFile.Register.ClassfileMod {
			claassMod = replace.New
		}
	}

	if claassMod.Version != "" {
		dir, err = modfetch.Download(claassMod)
		if err != nil {
			log.Fatalf("gop: download classsfile module error %v", err)
		}
	} else {
		dir = claassMod.Path
	}

	if dir == "" {
		log.Fatalf("gop: can't find classfile path in require statment")
	}

	gopmod := filepath.Join(modRoot, dir, "gop.mod")
	data, err := modfetch.Read(gopmod)
	if err != nil {
		log.Fatalf("gop: %v", err)
	}

	var fixed bool
	f, err := modfile.Parse(gopmod, data, fixVersion(&fixed))
	if err != nil {
		// Errors returned by modfile.Parse begin with file:line.
		log.Fatalf("go: errors parsing go.mod:\n%s\n", err)
	}
	classModFile = f
}
