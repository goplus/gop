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
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// -----------------------------------------------------------------------------

func GOPMOD(dir string) (file string, err error) {
	if dir == "" {
		dir = "."
	}
	if dir, err = filepath.Abs(dir); err != nil {
		return
	}
	for dir != "" {
		file = filepath.Join(dir, "gop.mod")
		if fi, e := os.Lstat(file); e == nil && !fi.IsDir() {
			return
		}
		file = filepath.Join(dir, "go.mod")
		if fi, e := os.Lstat(file); e == nil && !fi.IsDir() {
			return
		}
		if dir, file = filepath.Split(strings.TrimRight(dir, "/\\")); file == "" {
			break
		}
	}
	return "", syscall.ENOENT
}

// -----------------------------------------------------------------------------
