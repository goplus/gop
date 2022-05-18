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
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	testingGoFile    = "_test"
	autoGenFile      = "gop_autogen.go"
	autoGenTestFile  = "gop_autogen_test.go"
	autoGen2TestFile = "gop_autogen2_test.go"
)

// -----------------------------------------------------------------------------

func GenGo(dir string, conf *Config) (err error) {
	recursively := strings.HasSuffix(dir, "/...")
	if recursively {
		dir = dir[:len(dir)-4]
	}
	return genGoDir(dir, conf, recursively)
}

func genGoDir(dir string, conf *Config, recursively bool) (err error) {
	if recursively {
		return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err == nil && d.IsDir() {
				err = genGoIn(path, conf, true)
				if err != nil {
					if err == syscall.ENOENT {
						err = nil
					} else {
						fmt.Fprintln(os.Stderr, err)
					}
				}
			}
			return err
		})
	}
	return genGoIn(dir, conf, false)
}

func genGoIn(dir string, conf *Config, prompt bool) (err error) {
	out, test, err := LoadDir(dir, conf)
	if err != nil {
		return
	}

	if prompt {
		fmt.Printf("GenGo %v ...\n", dir)
	}

	os.MkdirAll(dir, 0755)
	file := filepath.Join(dir, autoGenFile)
	err = out.WriteFile(file)
	if err != nil {
		return
	}

	err = out.WriteFile(filepath.Join(dir, autoGenTestFile), testingGoFile)
	if err != nil && err != syscall.ENOENT {
		return
	}

	if test != nil {
		err = test.WriteFile(filepath.Join(dir, autoGen2TestFile), testingGoFile)
	} else {
		err = nil
	}
	return
}

// -----------------------------------------------------------------------------

func GenGoPkgPath(pkgPath string, conf *Config) (err error) {
	recursively := strings.HasSuffix(pkgPath, "/...")
	if recursively {
		pkgPath = pkgPath[:len(pkgPath)-4]
	}

	getPkgPathDo(pkgPath, gopEnv(conf), func(dir string) {
		os.Chmod(dir, 0755)
		defer os.Chmod(dir, 0555)
		err = genGoDir(dir, conf, recursively)
	}, func(e error) {
		err = e
	})
	return
}

// -----------------------------------------------------------------------------

func GenGoFiles(autogen string, files []string, conf *Config) (err error) {
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
		return
	}
	return out.WriteFile(autogen)
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
