/*
 * Copyright (c) 2023 The GoPlus Authors (goplus.org). All rights reserved.
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

package fsx

import (
	"io/fs"
	"os"
	"path/filepath"
)

// -----------------------------------------------------------------------------

// FileSystem represents a file system.
type FileSystem interface {
	ReadDir(dirname string) ([]fs.DirEntry, error)
	ReadFile(filename string) ([]byte, error)
	Join(elem ...string) string

	// Base returns the last element of path.
	// Trailing path separators are removed before extracting the last element.
	// If the path is empty, Base returns ".".
	// If the path consists entirely of separators, Base returns a single separator.
	Base(filename string) string

	// Abs returns an absolute representation of path.
	Abs(path string) (string, error)
}

// -----------------------------------------------------------------------------

type localFS struct{}

func (p localFS) ReadDir(dirname string) ([]fs.DirEntry, error) {
	return os.ReadDir(dirname)
}

func (p localFS) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func (p localFS) Join(elem ...string) string {
	return filepath.Join(elem...)
}

// Base returns the last element of path.
// Trailing path separators are removed before extracting the last element.
// If the path is empty, Base returns ".".
// If the path consists entirely of separators, Base returns a single separator.
func (p localFS) Base(filename string) string {
	return filepath.Base(filename)
}

// Abs returns an absolute representation of path.
func (p localFS) Abs(path string) (string, error) {
	return filepath.Abs(path)
}

var Local FileSystem = localFS{}

// -----------------------------------------------------------------------------
