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
	"log"
	"os"

	"github.com/goplus/gop/x/gocmd"
)

// -----------------------------------------------------------------------------

func InstallDir(dir string, conf *Config, install *gocmd.InstallConfig) (err error) {
	err = GenGo(dir, conf)
	if err != nil {
		return
	}
	return gocmd.Install(dir, install)
}

func InstallPkgPath(workDir, pkgPath string, conf *Config, install *gocmd.InstallConfig) (err error) {
	localDir, err := GenGoPkgPath(workDir, pkgPath, conf, true)
	if err != nil {
		return
	}
	old := chdir(localDir)
	defer os.Chdir(old)
	return gocmd.Install(localDir, install)
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
	err = GenGo(dir, conf)
	if err != nil {
		return
	}
	return gocmd.Build(dir, build)
}

func BuildPkgPath(workDir, pkgPath string, conf *Config, build *gocmd.BuildConfig) (err error) {
	localDir, err := GenGoPkgPath(workDir, pkgPath, conf, false)
	if err != nil {
		return
	}
	old := chdirAndMod(localDir)
	defer restoreDirAndMod(old)
	return gocmd.Build(localDir, build)
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
