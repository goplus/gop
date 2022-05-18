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

func GenGo(dir string, conf *Config) (err error) {
	recursively := strings.HasSuffix(dir, "/...")
	if recursively {
		dir = dir[:len(dir)-4]
		return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err == nil && d.IsDir() {
				err = genGoDir(path, conf, true)
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
	return genGoDir(dir, conf, false)
}

func genGoDir(dir string, conf *Config, prompt bool) (err error) {
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
