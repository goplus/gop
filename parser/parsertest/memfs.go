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

package parsertest

import (
	"os"
	"path"
	"syscall"
	"time"
)

// -----------------------------------------------------------------------------

type memFileInfo struct {
	name string
}

func (p *memFileInfo) Name() string {
	return p.name
}

func (p *memFileInfo) Size() int64 {
	return 0
}

func (p *memFileInfo) Mode() os.FileMode {
	return 0
}

func (p *memFileInfo) ModTime() (t time.Time) {
	return
}

func (p *memFileInfo) IsDir() bool {
	return false
}

func (p *memFileInfo) Sys() interface{} {
	return nil
}

// MemFS represents a file system in memory.
type MemFS struct {
	dirs  map[string][]string
	files map[string]string
}

// NewMemFS creates a MemFS instance.
func NewMemFS(dirs map[string][]string, files map[string]string) *MemFS {
	return &MemFS{dirs: dirs, files: files}
}

// ReadDir reads the directory named by dirname and returns
// a list of directory entries sorted by filename.
func (p *MemFS) ReadDir(dirname string) ([]os.FileInfo, error) {
	if items, ok := p.dirs[dirname]; ok {
		fis := make([]os.FileInfo, len(items))
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
func (p *MemFS) ReadFile(filename string) ([]byte, error) {
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
func (p *MemFS) Join(elem ...string) string {
	return path.Join(elem...)
}

// -----------------------------------------------------------------------------

// NewSingleFileFS creates a file system that only contains a single file.
func NewSingleFileFS(dir string, fname string, data string) *MemFS {
	return NewMemFS(map[string][]string{
		dir: {fname},
	}, map[string]string{
		path.Join(dir, fname): data,
	})
}

// -----------------------------------------------------------------------------
