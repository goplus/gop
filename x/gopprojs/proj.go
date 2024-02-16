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

package gopprojs

import (
	"errors"
	"path/filepath"
	"syscall"
)

// -----------------------------------------------------------------------------

type Proj = interface {
	projObj()
}

type FilesProj struct {
	Files []string
}

type PkgPathProj struct {
	Path string
}

type DirProj struct {
	Dir string
}

func (p *FilesProj) projObj()   {}
func (p *PkgPathProj) projObj() {}
func (p *DirProj) projObj()     {}

// -----------------------------------------------------------------------------

func ParseOne(args ...string) (proj Proj, next []string, err error) {
	if len(args) == 0 {
		return nil, nil, syscall.ENOENT
	}
	arg := args[0]
	if isFile(arg) {
		n := 1
		for n < len(args) && isFile(args[n]) {
			n++
		}
		return &FilesProj{Files: args[:n]}, args[n:], nil
	}
	if isLocal(arg) {
		return &DirProj{Dir: arg}, args[1:], nil
	}
	return &PkgPathProj{Path: arg}, args[1:], nil
}

func isFile(fname string) bool {
	n := len(filepath.Ext(fname))
	return n > 1
}

func isLocal(ns string) bool {
	if len(ns) > 0 {
		switch c := ns[0]; c {
		case '/', '\\', '.':
			return true
		default:
			return len(ns) >= 2 && ns[1] == ':' && ('A' <= c && c <= 'Z' || 'a' <= c && c <= 'z')
		}
	}
	return false
}

// -----------------------------------------------------------------------------

func ParseAll(args ...string) (projs []Proj, err error) {
	var hasFiles, hasNotFiles bool
	for {
		proj, next, e := ParseOne(args...)
		if e != nil {
			if hasFiles && hasNotFiles {
				return nil, ErrMixedFilesProj
			}
			return
		}
		if _, ok := proj.(*FilesProj); ok {
			hasFiles = true
		} else {
			hasNotFiles = true
		}
		projs = append(projs, proj)
		args = next
	}
}

var (
	ErrMixedFilesProj = errors.New("mixed files project")
)

// -----------------------------------------------------------------------------
