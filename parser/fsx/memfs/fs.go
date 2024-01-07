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

package memfs

import (
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type FileFS struct {
	data []byte
	info fs.DirEntry
}

type dirEntry struct {
	fs.FileInfo
}

func (p *dirEntry) Type() fs.FileMode {
	return p.FileInfo.Mode().Type()
}

func (p *dirEntry) Info() (fs.FileInfo, error) {
	return p.FileInfo, nil
}

func File(filename string, src interface{}) (f *FileFS, err error) {
	var data []byte
	var info fs.DirEntry
	if src != nil {
		data, err = readSource(src)
		if err != nil {
			return
		}
		info = &memFileInfo{name: filename, size: len(data)}
	} else {
		fi, e := os.Stat(filename)
		if e != nil {
			return nil, e
		}
		data, err = os.ReadFile(filename)
		if err != nil {
			return
		}
		info = &dirEntry{fi}
	}
	return &FileFS{data: data, info: info}, nil
}

func (p *FileFS) ReadDir(dirname string) ([]fs.DirEntry, error) {
	return []fs.DirEntry{p.info}, nil
}

func (p *FileFS) ReadFile(filename string) ([]byte, error) {
	return p.data, nil
}

func (p *FileFS) Join(elem ...string) string {
	return filepath.Join(elem...)
}

func (p *FileFS) Base(filename string) string {
	return filepath.Base(filename)
}

func (p *FileFS) Abs(path string) (string, error) {
	return path, nil
}

func readSource(src interface{}) ([]byte, error) {
	switch s := src.(type) {
	case string:
		return []byte(s), nil
	case []byte:
		return s, nil
	case *bytes.Buffer:
		// is io.Reader, but src is already available in []byte form
		if s != nil {
			return s.Bytes(), nil
		}
	case io.Reader:
		return io.ReadAll(s)
	}
	return nil, os.ErrInvalid
}
