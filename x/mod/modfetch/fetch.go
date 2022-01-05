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
	"path/filepath"
	"strings"
	"syscall"

	"github.com/goplus/gop/x/mod/modload"
	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
)

// -----------------------------------------------------------------------------

type (
	ExecCmdError = modload.ExecCmdError
)

func Get(modPath string, noCache ...bool) (mod module.Version, isClass bool, err error) {
	if noCache == nil || !noCache[0] {
		mod, isClass, err = getFromCache(modPath)
		if err != syscall.ENOENT {
			return
		}
	}
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

func getFromCache(modPath string) (modVer module.Version, isClass bool, err error) {
	modRoot, ver, err := lookupFromCache(modPath)
	if err != nil {
		return
	}
	modVer = module.Version{Path: modPath, Version: ver}
	isClass, err = convGoMod(modRoot)
	return
}

func lookupFromCache(modPath string) (modRoot string, modVer string, err error) {
	encPath, err := module.EscapePath(modPath)
	if err != nil {
		return
	}
	prefix := filepath.Join(GOMODCACHE, encPath)
	if pos := strings.IndexByte(modPath, '@'); pos >= 0 {
		fi, err := os.Stat(prefix)
		if err != nil || !fi.IsDir() {
			return "", "", syscall.ENOENT
		}
		return prefix, modPath[pos+1:], nil
	}
	dir, fname := filepath.Split(prefix)
	fis, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	fname += "@"
	err = syscall.ENOENT
	for _, fi := range fis {
		if fi.IsDir() {
			if name := fi.Name(); strings.HasPrefix(name, fname) {
				ver := name[len(fname):]
				if semver.Compare(modVer, ver) < 0 {
					modRoot, modVer, err = dir+name, ver, nil
				}
			}
		}
	}
	return
}

// -----------------------------------------------------------------------------
