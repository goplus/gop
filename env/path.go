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

package env

import (
	"log"
	"os"
	"path/filepath"
	"syscall"
)

var (
	// This is set by the linker.
	defaultGopRoot string
)

// GOPROOT returns the root of the Go+ tree. It uses the GOPROOT environment variable,
// if set at process start, or else the root used during the Go+ build.
func GOPROOT() string {
	gopRoot, err := findGopRoot()
	if err != nil {
		log.Panicln("GOPROOT not found:", err)
	}
	return gopRoot
}

const (
	envGOPROOT = "GOPROOT"
)

func findGopRoot() (string, error) {
	envGopRoot := os.Getenv(envGOPROOT)
	if envGopRoot != "" {
		// GOPROOT must valid
		if isValidGopRoot(envGopRoot) {
			return envGopRoot, nil
		}
		log.Panicf("\n%s (%s) is not valid\n", envGOPROOT, envGopRoot)
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

	// check defaultGopRoot, if it is valid, use it
	if defaultGopRoot != "" && isValidGopRoot(defaultGopRoot) {
		return defaultGopRoot, nil
	}

	// Compatible with old GOPROOT
	if home := HOME(); home != "" {
		gopRoot := filepath.Join(home, "gop")
		if isValidGopRoot(gopRoot) {
			return gopRoot, nil
		}
		goplusRoot := filepath.Join(home, "goplus")
		if isValidGopRoot(goplusRoot) {
			return goplusRoot, nil
		}
	}
	return "", syscall.ENOENT
}

// Mockable for testing.
var executable = func() (string, error) {
	return os.Executable()
}

func executableRealPath() (path string, err error) {
	path, err = executable()
	if err == nil {
		path, err = filepath.EvalSymlinks(path)
		if err == nil {
			path, err = filepath.Abs(path)
		}
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

func isValidGopRoot(path string) bool {
	return isDirExists(filepath.Join(path, "cmd/gop")) && isFileExists(filepath.Join(path, "go.mod"))
}
