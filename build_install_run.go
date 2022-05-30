/*
 * Copyright (c) 2022 The GoPlus Authors (goplus.org). All rights reserved.
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
	"errors"
	"log"
	"os"

	"github.com/goplus/gop/x/gocmd"
)

// -----------------------------------------------------------------------------

func InstallDir(dir string, conf *Config, install *gocmd.InstallConfig) (err error) {
	_, _, err = GenGo(dir, conf, false)
	if err != nil {
		return
	}
	return gocmd.Install(dir, install)
}

func InstallPkgPath(workDir, pkgPath string, conf *Config, install *gocmd.InstallConfig) (err error) {
	localDir, recursively, err := GenGoPkgPath(workDir, pkgPath, conf, true)
	if err != nil {
		return
	}
	old := chdir(localDir)
	defer os.Chdir(old)
	return gocmd.Install(cwdParam(recursively), install)
}

func cwdParam(recursively bool) string {
	if recursively {
		return "./..."
	}
	return "."
}

func InstallFiles(files []string, conf *Config, install *gocmd.InstallConfig) (err error) {
	files, err = GenGoFiles("", files, conf)
	if err != nil {
		return
	}
	return gocmd.InstallFiles(files, install)
}

func chdir(dir string) string {
	old, err := os.Getwd()
	if err != nil {
		log.Panicln(err)
	}
	err = os.Chdir(dir)
	if err != nil {
		log.Panicln(err)
	}
	return old
}

// -----------------------------------------------------------------------------

func BuildDir(dir string, conf *Config, build *gocmd.BuildConfig) (err error) {
	_, _, err = GenGo(dir, conf, false)
	if err != nil {
		return
	}
	return gocmd.Build(dir, build)
}

func BuildPkgPath(workDir, pkgPath string, conf *Config, build *gocmd.BuildConfig) (err error) {
	localDir, recursively, err := GenGoPkgPath(workDir, pkgPath, conf, false)
	if err != nil {
		return
	}
	old := chdirAndMod(localDir)
	defer restoreDirAndMod(old)
	return gocmd.Build(cwdParam(recursively), build)
}

func BuildFiles(files []string, conf *Config, build *gocmd.BuildConfig) (err error) {
	files, err = GenGoFiles("", files, conf)
	if err != nil {
		return
	}
	return gocmd.BuildFiles(files, build)
}

func chdirAndMod(dir string) string {
	os.Chmod(dir, 0777)
	return chdir(dir)
}

func restoreDirAndMod(old string) {
	os.Chmod(".", 0555)
	os.Chdir(old)
}

// -----------------------------------------------------------------------------

func RunDir(dir string, args []string, conf *Config, run *gocmd.RunConfig) (err error) {
	_, _, err = GenGo(dir, conf, false)
	if err != nil {
		return
	}
	return gocmd.RunDir(dir, args, run)
}

func RunPkgPath(pkgPath string, args []string, chDir bool, conf *Config, run *gocmd.RunConfig) (err error) {
	localDir, recursively, err := GenGoPkgPath("", pkgPath, conf, true)
	if err != nil {
		return
	}
	if recursively {
		return errors.New("can't use ... pattern for `gop run` command")
	}
	if chDir {
		old := chdir(localDir)
		defer os.Chdir(old)
		localDir = "."
	}
	return gocmd.RunDir(localDir, args, run)
}

func RunFiles(autogen string, files []string, args []string, conf *Config, run *gocmd.RunConfig) (err error) {
	files, err = GenGoFiles(autogen, files, conf)
	if err != nil {
		return
	}
	return gocmd.RunFiles(files, args, run)
}

// -----------------------------------------------------------------------------

func TestDir(dir string, conf *Config, test *gocmd.TestConfig) (err error) {
	_, _, err = GenGo(dir, conf, true)
	if err != nil {
		return
	}
	return gocmd.Test(dir, test)
}

func TestPkgPath(workDir, pkgPath string, conf *Config, test *gocmd.TestConfig) (err error) {
	localDir, recursively, err := GenGoPkgPath(workDir, pkgPath, conf, false)
	if err != nil {
		return
	}
	old := chdirAndMod(localDir)
	defer restoreDirAndMod(old)
	return gocmd.Test(cwdParam(recursively), test)
}

func TestFiles(files []string, conf *Config, test *gocmd.TestConfig) (err error) {
	files, err = GenGoFiles("", files, conf)
	if err != nil {
		return
	}
	return gocmd.TestFiles(files, test)
}

// -----------------------------------------------------------------------------
