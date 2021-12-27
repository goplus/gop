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
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/env"
	"github.com/goplus/gop/x/mod/modload"
)

var cmdInit = &base.Command{
	UsageLine: "gop mod init [module]",
	Short:     "initialize new module in current directory",
}

func init() {
	cmdInit.Run = runInit
}

func runInit(cmd *base.Command, args []string) {
	switch len(args) {
	case 0:
		fatal(`Example usage:
	'gop mod init example.com/m' to initialize a v0 or v1 module
	'gop mod init example.com/m/v2' to initialize a v2 module

Run 'gop help mod init' for more information.`)
	case 1:
	default:
		fatal("gop mod init: too many arguments")
	}
	modPath := args[0]
	mod, err := modload.Create(".", modPath, goMainVer(), env.MainVersion)
	if err != nil {
		fatal(err)
	}
	err = mod.Save()
	if err != nil {
		fatal(err)
	}
}

func goMainVer() string {
	ver := strings.TrimPrefix(runtime.Version(), "go")
	if pos := strings.Index(ver, "."); pos > 0 {
		pos++
		if pos2 := strings.Index(ver[pos:], "."); pos2 > 0 {
			ver = ver[:pos+pos2]
		}
	}
	return ver
}

func fatal(msg interface{}) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
