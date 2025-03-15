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

package memfs

import (
	"io/fs"
	"path"
	"path/filepath"
	"syscall"
	"time"
)

// -----------------------------------------------------------------------------

type memFileInfo struct {
	name string
	size int
}

func (p *memFileInfo) Name() string {
	return filepath.Base(p.name)
}

func (p *memFileInfo) Size() int64 {
	return int64(p.size)
}

func (p *memFileInfo) Mode() fs.FileMode {
	return 0
}

func (p *memFileInfo) Type() fs.FileMode {
	return 0
}

func (p *memFileInfo) ModTime() (t time.Time) {
	return time.Now()
}

func (p *memFileInfo) IsDir() bool {
	return false
}

func (p *memFileInfo) Sys() any {
	return nil
}

func (p *memFileInfo) Info() (fs.FileInfo, error) {
	return p, nil
}

// -----------------------------------------------------------------------------

// FS represents a file system in memory.
type FS struct {
	dirs  map[string][]string
	files map[string]string
}

// New creates a file system instance.
func New(dirs map[string][]string, files map[string]string) *FS {
	return &FS{dirs: dirs, files: files}
}

// ReadDir reads the directory named by dirname and returns
// a list of directory entries sorted by filename.
func (p *FS) ReadDir(dirname string) ([]fs.DirEntry, error) {
	if items, ok := p.dirs[dirname]; ok {
		fis := make([]fs.DirEntry, len(items))
		for i, item := range items {
			fis[i] = &memFileInfo{name: item}
		}
		return fis, nil
	}
	return nil, syscall.ENOENT
}

// ReadFile reads the file named by filename and returns the contents.
// A successful call returns err == nil, not err == EOF. Because ReadFile
// reads the whole file, it does not treat an EOF from Read as an error
// to be reported.
func (p *FS) ReadFile(filename string) ([]byte, error) {
	if data, ok := p.files[filename]; ok {
		return []byte(data), nil
	}
	return nil, syscall.ENOENT
}

// Join joins any number of path elements into a single path,
// separating them with slashes. Empty elements are ignored.
// The result is Cleaned. However, if the argument list is
// empty or all its elements are empty, Join returns
// an empty string.
func (p *FS) Join(elem ...string) string {
	return path.Join(elem...)
}

// Base returns the last element of path.
// Trailing slashes are removed before extracting the last element.
// If the path is empty, Base returns ".".
// If the path consists entirely of slashes, Base returns "/".
func (p *FS) Base(filename string) string {
	return path.Base(filename)
}

func (p *FS) Abs(path string) (string, error) {
	return path, nil
}

// -----------------------------------------------------------------------------

// SingleFile creates a file system that only contains a single file.
func SingleFile(dir string, fname string, data string) *FS {
	return New(map[string][]string{
		dir: {fname},
	}, map[string]string{
		path.Join(dir, fname): data,
	})
}

// TwoFiles creates a file system that contains two files.
func TwoFiles(dir string, fname1, data1, fname2, data2 string) *FS {
	return New(map[string][]string{
		dir: {fname1, fname2},
	}, map[string]string{
		path.Join(dir, fname1): data1,
		path.Join(dir, fname2): data2,
	})
}

// -----------------------------------------------------------------------------
