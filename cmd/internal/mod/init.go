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

// gop mod init

package mod

import (
	"log"

	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/cmd/internal/modload"
)

var cmdInit = &base.Command{
	UsageLine: "gop mod init [module]",
	Short:     "initialize new module in current directory",
}

func init() {
	cmdInit.Run = runInit
}

func runInit(cmd *base.Command, args []string) {
	if len(args) > 1 {
		log.Fatalf("gop mod init: too many arguments")
	}
	var modPath string
	if len(args) == 1 {
		modPath = args[0]
	}

	// modfetch.InitArgs(".", args...)
	modload.CreateModFile(modPath) // does all the hard work
	modload.LoadModFile()
	modload.SyncGoMod()
	modload.SyncGopMod()
}
