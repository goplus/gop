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
	"github.com/qiniu/x/errors"
)

// -----------------------------------------------------------------------------

// InstallDir installs a Go+ package directory.
// if conf != nil && conf.Context == nil, it will be set with `gox.NewContext()`.
func InstallDir(dir string, conf *Config, install *gocmd.InstallConfig) (err error) {
	_, _, err = GenGo(dir, conf, false)
	if err != nil {
		return errors.NewWith(err, `GenGo(dir, conf, false)`, -2, "gop.GenGo", dir, conf, false)
	}
	return gocmd.Install(dir, install)
}

// InstallPkgPath installs a Go+ package.
// if conf != nil && conf.Context == nil, it will be set with `gox.NewContext()`.
func InstallPkgPath(workDir, pkgPath string, conf *Config, install *gocmd.InstallConfig) (err error) {
	localDir, recursively, err := GenGoPkgPath(workDir, pkgPath, conf, true)
	if err != nil {
		return errors.NewWith(err, `GenGoPkgPath(workDir, pkgPath, conf, true)`, -2, "gop.GenGoPkgPath", workDir, pkgPath, conf, true)
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

// InstallFiles installs specified Go+ files.
// if conf != nil && conf.Context == nil, it will be set with `gox.NewContext()`.
func InstallFiles(files []string, conf *Config, install *gocmd.InstallConfig) (err error) {
	files, err = GenGoFiles("", files, conf)
	if err != nil {
		return errors.NewWith(err, `GenGoFiles("", files, conf)`, -2, "gop.GenGoFiles", "", files, conf)
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

// BuildDir builds a Go+ package directory.
// if conf != nil && conf.Context == nil, it will be set with `gox.NewContext()`.
func BuildDir(dir string, conf *Config, build *gocmd.BuildConfig) (err error) {
	_, _, err = GenGo(dir, conf, false)
	if err != nil {
		return errors.NewWith(err, `GenGo(dir, conf, false)`, -2, "gop.GenGo", dir, conf, false)
	}
	return gocmd.Build(dir, build)
}

// BuildPkgPath builds a Go+ package.
// if conf != nil && conf.Context == nil, it will be set with `gox.NewContext()`.
func BuildPkgPath(workDir, pkgPath string, conf *Config, build *gocmd.BuildConfig) (err error) {
	localDir, recursively, err := GenGoPkgPath(workDir, pkgPath, conf, false)
	if err != nil {
		return errors.NewWith(err, `GenGoPkgPath(workDir, pkgPath, conf, false)`, -2, "gop.GenGoPkgPath", workDir, pkgPath, conf, false)
	}
	old, mod := chdirAndMod(localDir)
	defer restoreDirAndMod(old, mod)
	return gocmd.Build(cwdParam(recursively), build)
}

// BuildFiles builds specified Go+ files.
// if conf != nil && conf.Context == nil, it will be set with `gox.NewContext()`.
func BuildFiles(files []string, conf *Config, build *gocmd.BuildConfig) (err error) {
	files, err = GenGoFiles("", files, conf)
	if err != nil {
		return errors.NewWith(err, `GenGoFiles("", files, conf)`, -2, "gop.GenGoFiles", "", files, conf)
	}
	return gocmd.BuildFiles(files, build)
}

func chdirAndMod(dir string) (old string, mod os.FileMode) {
	mod = 0755
	if info, err := os.Stat(dir); err == nil {
		mod = info.Mode().Perm()
	}
	os.Chmod(dir, 0777)
	old = chdir(dir)
	return
}

func restoreDirAndMod(old string, mod os.FileMode) {
	os.Chmod(".", mod)
	os.Chdir(old)
}

// -----------------------------------------------------------------------------

// RunDir runs an application from a Go+ package directory.
// if conf != nil && conf.Context == nil, it will be set with `gox.NewContext()`.
func RunDir(dir string, args []string, conf *Config, run *gocmd.RunConfig) (err error) {
	_, _, err = GenGo(dir, conf, false)
	if err != nil {
		return errors.NewWith(err, `GenGo(dir, conf, false)`, -2, "gop.GenGo", dir, conf, false)
	}
	return gocmd.RunDir(dir, args, run)
}

// RunPkgPath runs an application from a Go+ package.
// if conf != nil && conf.Context == nil, it will be set with `gox.NewContext()`.
func RunPkgPath(pkgPath string, args []string, chDir bool, conf *Config, run *gocmd.RunConfig) (err error) {
	localDir, recursively, err := GenGoPkgPath("", pkgPath, conf, true)
	if err != nil {
		return errors.NewWith(err, `GenGoPkgPath("", pkgPath, conf, true)`, -2, "gop.GenGoPkgPath", "", pkgPath, conf, true)
	}
	if recursively {
		return errors.NewWith(errors.New("can't use ... pattern for `gop run` command"), `recursively`, -1, "", recursively)
	}
	if chDir {
		old := chdir(localDir)
		defer os.Chdir(old)
		localDir = "."
	}
	return gocmd.RunDir(localDir, args, run)
}

// RunFiles runs an application from specified Go+ files.
// if conf != nil && conf.Context == nil, it will be set with `gox.NewContext()`.
func RunFiles(autogen string, files []string, args []string, conf *Config, run *gocmd.RunConfig) (err error) {
	files, err = GenGoFiles(autogen, files, conf)
	if err != nil {
		return errors.NewWith(err, `GenGoFiles(autogen, files, conf)`, -2, "gop.GenGoFiles", autogen, files, conf)
	}
	return gocmd.RunFiles(files, args, run)
}

// -----------------------------------------------------------------------------

// TestDir tests a Go+ package directory.
// if conf != nil && conf.Context == nil, it will be set with `gox.NewContext()`.
func TestDir(dir string, conf *Config, test *gocmd.TestConfig) (err error) {
	_, _, err = GenGo(dir, conf, true)
	if err != nil {
		return errors.NewWith(err, `GenGo(dir, conf, true)`, -2, "gop.GenGo", dir, conf, true)
	}
	return gocmd.Test(dir, test)
}

// TestPkgPath tests a Go+ package.
// if conf != nil && conf.Context == nil, it will be set with `gox.NewContext()`.
func TestPkgPath(workDir, pkgPath string, conf *Config, test *gocmd.TestConfig) (err error) {
	localDir, recursively, err := GenGoPkgPath(workDir, pkgPath, conf, false)
	if err != nil {
		return errors.NewWith(err, `GenGoPkgPath(workDir, pkgPath, conf, false)`, -2, "gop.GenGoPkgPath", workDir, pkgPath, conf, false)
	}
	old, mod := chdirAndMod(localDir)
	defer restoreDirAndMod(old, mod)
	return gocmd.Test(cwdParam(recursively), test)
}

// TestFiles tests specified Go+ files.
// if conf != nil && conf.Context == nil, it will be set with `gox.NewContext()`.
func TestFiles(files []string, conf *Config, test *gocmd.TestConfig) (err error) {
	files, err = GenGoFiles("", files, conf)
	if err != nil {
		return errors.NewWith(err, `GenGoFiles("", files, conf)`, -2, "gop.GenGoFiles", "", files, conf)
	}
	return gocmd.TestFiles(files, test)
}

// -----------------------------------------------------------------------------
