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
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/goplus/mod/gopmod"
	"github.com/goplus/mod/modcache"
	"github.com/goplus/mod/modfetch"
	"github.com/qiniu/x/errors"
)

const (
	testingGoFile    = "_test"
	autoGenFile      = "gop_autogen.go"
	autoGenTestFile  = "gop_autogen_test.go"
	autoGen2TestFile = "gop_autogen2_test.go"
)

// -----------------------------------------------------------------------------

func GenGo(dir string, conf *Config, genTestPkg bool) (string, bool, error) {
	recursively := strings.HasSuffix(dir, "/...")
	if recursively {
		dir = dir[:len(dir)-4]
	}
	return dir, recursively, genGoDir(dir, conf, genTestPkg, recursively)
}

func genGoDir(dir string, conf *Config, genTestPkg, recursively bool) (err error) {
	if recursively {
		var list errors.List
		fn := func(path string, d fs.DirEntry, err error) error {
			if err == nil && d.IsDir() {
				if strings.HasPrefix(d.Name(), "_") { // skip _
					return filepath.SkipDir
				}
				if e := genGoIn(path, conf, genTestPkg, true); e != nil {
					list.Add(e)
				}
			}
			return err
		}
		err = filepath.WalkDir(dir, fn)
		if err != nil {
			return errors.NewWith(err, `filepath.WalkDir(dir, fn)`, -2, "filepath.WalkDir", dir, fn)
		}
		return list.ToError()
	}
	return genGoIn(dir, conf, genTestPkg, false)
}

func genGoIn(dir string, conf *Config, genTestPkg, prompt bool) (err error) {
	out, test, err := LoadDir(dir, conf, genTestPkg, prompt)
	if err != nil {
		if err == syscall.ENOENT { // no Go+ source files
			return nil
		}
		return errors.NewWith(err, `LoadDir(dir, conf, genTestPkg, prompt)`, -5, "gop.LoadDir", dir, conf, genTestPkg, prompt)
	}

	os.MkdirAll(dir, 0755)
	file := filepath.Join(dir, autoGenFile)
	err = out.WriteFile(file)
	if err != nil {
		return errors.NewWith(err, `out.WriteFile(file)`, -2, "(*gox.Package).WriteFile", out, file)
	}

	testFile := filepath.Join(dir, autoGenTestFile)
	err = out.WriteFile(testFile, testingGoFile)
	if err != nil && err != syscall.ENOENT {
		return errors.NewWith(err, `out.WriteFile(testFile, testingGoFile)`, -2, "(*gox.Package).WriteFile", out, testFile, testingGoFile)
	}

	if test != nil {
		testFile = filepath.Join(dir, autoGen2TestFile)
		err = test.WriteFile(testFile, testingGoFile)
		if err != nil {
			return errors.NewWith(err, `test.WriteFile(testFile, testingGoFile)`, -2, "(*gox.Package).WriteFile", test, testFile, testingGoFile)
		}
	} else {
		err = nil
	}
	return
}

// -----------------------------------------------------------------------------

const (
	modWritable = 0755
	modReadonly = 0555
)

func GenGoPkgPath(workDir, pkgPath string, conf *Config, allowExtern bool) (localDir string, recursively bool, err error) {
	recursively = strings.HasSuffix(pkgPath, "/...")
	if recursively {
		pkgPath = pkgPath[:len(pkgPath)-4]
	}

	mod, err := gopmod.Load(workDir, 0)
	if NotFound(err) && allowExtern {
		remotePkgPathDo(pkgPath, func(dir string) {
			os.Chmod(dir, modWritable)
			defer os.Chmod(dir, modReadonly)
			localDir = dir
			err = genGoDir(dir, conf, false, recursively)
		}, func(e error) {
			err = e
		})
		return
	} else if err != nil {
		return
	}

	pkg, err := mod.Lookup(pkgPath)
	if err != nil {
		return
	}
	localDir = pkg.Dir
	if pkg.Type == gopmod.PkgtExtern {
		os.Chmod(localDir, modWritable)
		defer os.Chmod(localDir, modReadonly)
	}
	err = genGoDir(localDir, conf, false, recursively)
	return
}

func remotePkgPathDo(pkgPath string, doSth func(dir string), onErr func(e error)) {
	modVer, leftPart, err := modfetch.GetPkg(pkgPath, "")
	if err != nil {
		onErr(err)
	} else if dir, err := modcache.Path(modVer); err != nil {
		onErr(err)
	} else {
		doSth(filepath.Join(dir, leftPart))
	}
}

// -----------------------------------------------------------------------------

func GenGoFiles(autogen string, files []string, conf *Config) (result []string, err error) {
	if autogen == "" {
		autogen = "gop_autogen.go"
		if len(files) == 1 {
			file := files[0]
			srcDir, fname := filepath.Split(file)
			if hasMultiFiles(srcDir, ".gop") {
				autogen = filepath.Join(srcDir, "gop_autogen_"+fname+".go")
			}
		}
	}
	out, err := LoadFiles(files, conf)
	if err != nil {
		err = errors.NewWith(err, `LoadFiles(files, conf)`, -2, "gop.LoadFiles", files, conf)
		return
	}
	result = append(result, autogen)
	err = out.WriteFile(autogen)
	if err != nil {
		err = errors.NewWith(err, `out.WriteFile(autogen)`, -2, "(*gox.Package).WriteFile", out, autogen)
	}
	return
}

func hasMultiFiles(srcDir string, ext string) bool {
	var has bool
	if f, err := os.Open(srcDir); err == nil {
		defer f.Close()
		fis, _ := f.ReadDir(-1)
		for _, fi := range fis {
			if !fi.IsDir() && filepath.Ext(fi.Name()) == ext {
				if has {
					return true
				}
				has = true
			}
		}
	}
	return false
}

// -----------------------------------------------------------------------------
