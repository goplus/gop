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

package modfetch

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/goplus/gop/x/mod/modload"
	"golang.org/x/mod/module"
)

// -----------------------------------------------------------------------------

type (
	ExecCmdError = modload.ExecCmdError
)

func Get(modPath string) (mod module.Version, isClass bool, err error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("go", "get", modPath)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		err = &ExecCmdError{Err: err, Stderr: stderr.Bytes()}
		return
	}
	return getResult(stderr.String())
}

func getResult(data string) (mod module.Version, isClass bool, err error) {
	if data == "" {
		err = syscall.EEXIST
		return
	}
	// go: downloading github.com/xushiwei/foogop v0.1.0
	const downloading = "go: downloading "
	if strings.HasPrefix(data, downloading) {
		return tryConvGoMod(data[len(downloading):], &data)
	}
	// go get: added github.com/xushiwei/foogop v0.1.0
	const added = "go get: added "
	if strings.HasPrefix(data, added) {
		return tryConvGoMod(data[len(added):], &data)
	}
	return
}

func tryConvGoMod(data string, next *string) (mod module.Version, isClass bool, err error) {
	err = syscall.ENOENT
	if pos := strings.IndexByte(data, '\n'); pos > 0 {
		line := data[:pos]
		*next = data[pos+1:]
		if pos = strings.IndexByte(line, ' '); pos > 0 {
			mod.Path, mod.Version = line[:pos], line[pos+1:]
			if dir, e := ModCachePath(mod); e == nil {
				isClass, err = convGoMod(dir)
			}
		}
	}
	return
}

func convGoMod(dir string) (isClass bool, err error) {
	mod, err := modload.Load(dir)
	if err != nil {
		return
	}
	os.Chmod(dir, 0755)
	defer os.Chmod(dir, 0555)
	return mod.Classfile != nil, mod.UpdateGoMod(true)
}

// -----------------------------------------------------------------------------
