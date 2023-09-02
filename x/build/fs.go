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

package build

import (
	"bytes"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

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

type fileFS struct {
	dir      string
	filename string
	data     []byte
	info     fs.FileInfo
}

func newFileFS(filename string, src interface{}) (f *fileFS, err error) {
	dir, name := filepath.Split(filename)
	var data []byte
	var info fs.FileInfo
	if src != nil {
		data, err = readSource(src)
		if err != nil {
			return
		}
		info = &memFileInfo{name: name}
	} else {
		info, err = os.Stat(filename)
		if err != nil {
			return
		}
		data, err = os.ReadFile(filename)
		if err != nil {
			return
		}
	}
	return &fileFS{dir: dir, filename: filename, data: data, info: info}, nil
}

func (p fileFS) ReadDir(dirname string) ([]fs.FileInfo, error) {
	return []fs.FileInfo{p.info}, nil
}

func (p fileFS) ReadFile(filename string) ([]byte, error) {
	return p.data, nil
}

func (p fileFS) Join(elem ...string) string {
	return filepath.Join(elem...)
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
		return ioutil.ReadAll(s)
	}
	return nil, os.ErrInvalid
}
