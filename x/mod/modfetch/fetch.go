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
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
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
	cmd := exec.Command("go", "get", "-d", "-x", modPath)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	comments, infos := parserGetResults(stderr.String())
	if err != nil {
		if mod, isClass, ok := tryFoundModule(modPath, infos); ok {
			return mod, isClass, nil
		}
		err = &ExecCmdError{Err: err, Stderr: []byte(infos)}
		return
	}
	if len(infos) > 0 {
		return getResult(infos)
	}
	if len(comments) > 0 {
		if m, ok := foundModuleInComment(modPath, comments[0]); ok {
			modPath = m
		}
	}
	return getFromCache(modPath)
}

func getResult(data string) (mod module.Version, isClass bool, err error) {
	// go: downloading github.com/xushiwei/foogop v0.1.0
	const downloading = "go: downloading "
	if strings.HasPrefix(data, downloading) {
		if pos := strings.IndexByte(data, '\n'); pos > 0 {
			fmt.Fprintln(os.Stderr, "gop:", data[4:pos])
		}
		return tryConvGoMod(data[len(downloading):], &data)
	}
	// go1.17  go get: added github.com/xushiwei/foogop v0.1.0
	// go1.18  go: added github.com/xushiwei/foogop v0.1.0
	const added = ": added "
	if pos := strings.Index(data, added); pos > 1 {
		return tryConvGoMod(data[pos+len(added):], &data)
	}
	err = fmt.Errorf("unknown go get result: %v", data)
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
	modRoot, modVer, err := lookupFromCache(modPath)
	if err != nil {
		return
	}
	isClass, err = convGoMod(modRoot)
	return
}

func lookupFromCache(modPath string) (modRoot string, mod module.Version, err error) {
	mod.Path = modPath
	pos := strings.IndexByte(modPath, '@')
	if pos > 0 {
		mod.Path, mod.Version = modPath[:pos], modPath[pos+1:]
	}
	encPath, err := module.EscapePath(mod.Path)
	if err != nil {
		return
	}
	modRoot = filepath.Join(GOMODCACHE, encPath+"@"+mod.Version)
	if pos > 0 { // has version
		modRoot, mod.Path, err = foundModRoot(encPath, mod.Version)
		return
	}
	dir, fname := filepath.Split(modRoot)
	fis, err := os.ReadDir(dir)
	if err != nil {
		err = errors.Unwrap(err)
		return
	}
	err = syscall.ENOENT
	for _, fi := range fis {
		if fi.IsDir() {
			if name := fi.Name(); strings.HasPrefix(name, fname) {
				ver := name[len(fname):]
				if semver.Compare(mod.Version, ver) < 0 {
					modRoot, mod.Version, err = dir+name, ver, nil
				}
			}
		}
	}
	return
}

// parser from `go get -x pkgpath`
func parserGetResults(data string) (comment []string, result string) {
	for _, line := range strings.Split(data, "\n") {
		if strings.HasPrefix(line, "# get ") {
			if strings.Contains(line, "@v/list: 200 OK") {
				comment = append(comment, line)
			}
		} else if strings.HasPrefix(line, "go: ") {
			result = line + "\n"
		} else if strings.HasPrefix(line, "go get: ") {
			result = line + "\n"
		}
	}
	return
}

// found module from go get pkgpath/@v/list comment
// # get https://goproxy.cn/github.com/goplus/spx/@v/list: 200 OK (0.782s)
func foundModuleInComment(modPath string, data string) (mod string, ok bool) {
	if strings.Contains(data, modPath+"/@v/list") {
		return modPath, true
	}
	dir, _ := path.Split(modPath)
	if dir == "" {
		return "", false
	}
	return foundModuleInComment(dir[:len(dir)-1], data)
}

// found pkgpath@ver to modpath@ver
// github.com/goplus/spx/tutorial/04-Bullet@v1.0.0-rc5 => github.com/goplus/spx@v1.0.0-rc5
func foundModRoot(modPath string, ver string) (modRoot string, mod string, err error) {
	modRoot = filepath.Join(GOMODCACHE, modPath+"@"+ver)
	if fi, e := os.Stat(modRoot); e == nil {
		if fi.IsDir() {
			return modRoot, modPath, nil
		}
		return "", "", syscall.ENOENT
	}
	dir, _ := path.Split(modPath)
	if dir == "" {
		return "", "", syscall.ENOENT
	}
	return foundModRoot(dir[:len(dir)-1], ver)
}

// found module msg: found, but does not contain package
func tryFoundModule(modPath string, data string) (mod module.Version, isClass bool, ok bool) {
	// go get: module github.com/goplus/spx@v1.0.0-rc5 found, but does not contain package github.com/goplus/spx/tutorial/04-Bullet
	// go get: module github.com/goplus/spx@main found (v1.0.0-rc3.0.20220110030840-d39f5107c481), but does not contain package github.com/goplus/spx/tutorial/00-Hello
	// go1.17 syntax go get: module
	// go1.18 syntax go: module
	const getmodule = ": module "
	if pos := strings.Index(data, getmodule); pos > 0 {
		if pos := strings.IndexByte(data, '\n'); pos > 0 {
			data = data[:pos]
		}
		data = data[pos+len(getmodule):]
		const notpkg = ", but does not contain package"
		if pos := strings.Index(data, notpkg); pos > 0 {
			list := strings.Split(data[:pos], " ")
			switch len(list) {
			case 2:
				dir := list[0]
				var err error
				mod, isClass, err = getFromCache(dir)
				ok = (err == nil)
			case 3:
				ver := list[2]
				if len(ver) > 3 && ver[0] == '(' && ver[len(ver)-1] == ')' {
					mpath, _ := splitVerson(list[0])
					var err error
					mod, isClass, err = getFromCache(mpath + "@" + ver[1:len(ver)-1])
					ok = (err == nil)
				}
			}
		}
	}
	return
}

func splitVerson(modPath string) (path, version string) {
	pos := strings.IndexByte(modPath, '@')
	if pos > 0 {
		return modPath[:pos], modPath[pos+1:]
	}
	return modPath, ""
}

// -----------------------------------------------------------------------------
