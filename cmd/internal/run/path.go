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

package run

import (
	"os"
	"path/filepath"
	"syscall"

	"github.com/goplus/gop/build"
	"github.com/goplus/gop/cl"
	"github.com/qiniu/x/log"
)

const ErrNotFound = syscall.ENOENT

const (
	ENV_GOPROOT = "GOPROOT"
	ENV_HOME    = "HOME"
)

func findGoModFile(dir string) (modfile string, noCacheFile bool, err error) {
	modfile, err = cl.FindGoModFile(dir)

	if err != nil {
		gopRoot, err := findGopRoot()
		if err == nil {
			modfile = filepath.Join(gopRoot, "go.mod")
			return modfile, true, nil
		}
	}
	return
}

func findGopRoot() (string, error) {
	envGopRoot := os.Getenv(ENV_GOPROOT)
	if envGopRoot != "" {
		// GOPROOT must valid
		if isValidGopRoot(envGopRoot) {
			return envGopRoot, nil
		}
		log.Panicf("\n%s (%s) is not valid\n", ENV_GOPROOT, envGopRoot)
	}

	// if parent directory is a valid gop root, use it
	exePath, err := executableRealPath()
	if err == nil {
		dir := filepath.Dir(exePath)
		parentDir := filepath.Dir(dir)
		if parentDir != dir && isValidGopRoot(parentDir) {
			return parentDir, nil
		}
	}

	// check build.GopRoot, if it is valid, use it
	if build.GopRoot != "" && isValidGopRoot(build.GopRoot) {
		return build.GopRoot, nil
	}

	// Compatible with old GOPROOT
	if home := os.Getenv(ENV_HOME); home != "" {
		gopRoot := filepath.Join(home, "gop")
		if isValidGopRoot(gopRoot) {
			return gopRoot, nil
		}
		goplusRoot := filepath.Join(home, "goplus")
		if isValidGopRoot(goplusRoot) {
			return goplusRoot, nil
		}
	}

	return "", ErrNotFound
}

// Mockable for testing.
var executable = func() (string, error) {
	return os.Executable()
}

func executableRealPath() (path string, err error) {
	path, err = executable()
	if err != nil {
		return
	}

	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		return
	}

	path, err = filepath.Abs(path)
	if err != nil {
		return
	}

	return
}

func isFileExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}

func isDirExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && st.IsDir()
}

// A valid gop root has go.mod, go.sum, builtin/
func isValidGopRoot(path string) bool {
	if !isFileExists(filepath.Join(path, "go.mod")) ||
		!isFileExists(filepath.Join(path, "go.sum")) ||
		!isDirExists(filepath.Join(path, "builtin")) {
		return false
	}

	return true
}
