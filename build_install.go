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
	"github.com/goplus/gop/x/gocmd"
)

// -----------------------------------------------------------------------------

func InstallDir(dir string, conf *Config, install *gocmd.InstallConfig) (err error) {
	err = GenGo(dir, conf)
	if err != nil {
		return
	}
	return gocmd.InstallDir(dir, install)
}

func InstallPkgPath(pkgPath string, conf *Config, install *gocmd.InstallConfig) (err error) {
	localDir, err := GenGoPkgPath(pkgPath, conf)
	if err != nil {
		return
	}
	return gocmd.InstallDir(localDir, install)
}

func InstallFiles(files []string, conf *Config, install *gocmd.InstallConfig) (err error) {
	files, err = GenGoFiles("", files, conf)
	if err != nil {
		return
	}
	return gocmd.InstallFiles(files, install)
}

// -----------------------------------------------------------------------------

func BuildDir(dir string, conf *Config, build *gocmd.InstallConfig) (err error) {
	err = GenGo(dir, conf)
	if err != nil {
		return
	}
	return gocmd.BuildDir(dir, build)
}

func BuildPkgPath(pkgPath string, conf *Config, build *gocmd.InstallConfig) (err error) {
	localDir, err := GenGoPkgPath(pkgPath, conf)
	if err != nil {
		return
	}
	return gocmd.BuildDir(localDir, build)
}

func BuildFiles(files []string, conf *Config, build *gocmd.InstallConfig) (err error) {
	files, err = GenGoFiles("", files, conf)
	if err != nil {
		return
	}
	return gocmd.BuildFiles(files, build)
}

// -----------------------------------------------------------------------------
