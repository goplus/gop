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

package parser

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/token"
)

const (
	DbgFlagParseOutput = 1 << iota
	DbgFlagParseError
	DbgFlagAll = DbgFlagParseOutput | DbgFlagParseError
)

var (
	debugParseOutput bool
	debugParseError  bool
)

func SetDebug(dbgFlags int) {
	debugParseOutput = (dbgFlags & DbgFlagParseOutput) != 0
	debugParseError = (dbgFlags & DbgFlagParseError) != 0
}

// -----------------------------------------------------------------------------

// FileSystem represents a file system.
type FileSystem interface {
	ReadDir(dirname string) ([]fs.FileInfo, error)
	ReadFile(filename string) ([]byte, error)
	Join(elem ...string) string
}

type localFS struct{}

func (p localFS) ReadDir(dirname string) ([]fs.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

func (p localFS) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func (p localFS) Join(elem ...string) string {
	return filepath.Join(elem...)
}

var local FileSystem = localFS{}

// Parse parses a single Go+ source file. The target specifies the Go+ source file.
// If the file couldn't be read, a nil map and the respective error are returned.
func Parse(fset *token.FileSet, target string, src interface{}, mode Mode) (pkgs map[string]*ast.Package, err error) {
	file, err := ParseFile(fset, target, src, mode)
	if err != nil {
		return
	}
	pkgs = make(map[string]*ast.Package)
	pkgs[file.Name.Name] = astFileToPkg(file, target)
	return
}

// astFileToPkg translate ast.File to ast.Package
func astFileToPkg(file *ast.File, fileName string) (pkg *ast.Package) {
	pkg = &ast.Package{
		Name:  file.Name.Name,
		Files: make(map[string]*ast.File),
	}
	pkg.Files[fileName] = file
	return
}

// -----------------------------------------------------------------------------

// ParseDir calls ParseFSDir by passing a local filesystem.
//
func ParseDir(fset *token.FileSet, path string, filter func(fs.FileInfo) bool, mode Mode) (pkgs map[string]*ast.Package, first error) {
	return ParseFSDir(fset, local, path, Config{Filter: filter, Mode: mode})
}

type Config struct {
	IsClass func(ext string) (isProj bool, ok bool)
	Filter  func(fs.FileInfo) bool
	Mode    Mode
}

// ParseDirEx calls ParseFSDir by passing a local filesystem.
//
func ParseDirEx(fset *token.FileSet, path string, conf Config) (pkgs map[string]*ast.Package, first error) {
	return ParseFSDir(fset, local, path, conf)
}

// ParseFSDir calls ParseFile for all files with names ending in ".gop" in the
// directory specified by path and returns a map of package name -> package
// AST with all the packages found.
//
// If filter != nil, only the files with fs.FileInfo entries passing through
// the filter (and ending in ".gop") are considered. The mode bits are passed
// to ParseFile unchanged. Position information is recorded in fset, which
// must not be nil.
//
// If the directory couldn't be read, a nil map and the respective error are
// returned. If a parse error occurred, a non-nil but incomplete map and the
// first error encountered are returned.
//
func ParseFSDir(fset *token.FileSet, fs FileSystem, path string, conf Config) (pkgs map[string]*ast.Package, first error) {
	list, err := fs.ReadDir(path)
	if err != nil {
		return nil, err
	}
	if conf.IsClass == nil {
		conf.IsClass = defaultIsClass
	}
	pkgs = make(map[string]*ast.Package)
	for _, d := range list {
		if d.IsDir() {
			continue
		}
		fname := d.Name()
		ext := filepath.Ext(fname)
		var isProj, isClass bool
		switch ext {
		case ".gop":
		case ".go":
			if strings.HasPrefix(fname, "gop_autogen") {
				continue
			}
		default:
			if isProj, isClass = conf.IsClass(ext); !isClass {
				continue
			}
		}
		if !strings.HasPrefix(fname, "_") && (conf.Filter == nil || conf.Filter(d)) {
			filename := fs.Join(path, fname)
			if filedata, err := fs.ReadFile(filename); err == nil {
				if src, err := ParseFSFile(fset, fs, filename, filedata, conf.Mode); err == nil {
					src.IsProj, src.IsClass = isProj, isClass
					name := src.Name.Name
					pkg, found := pkgs[name]
					if !found {
						pkg = &ast.Package{
							Name:  name,
							Files: make(map[string]*ast.File),
						}
						pkgs[name] = pkg
					}
					pkg.Files[filename] = src
				} else if first == nil {
					first = err
				}
			} else if first == nil {
				first = err
			}
		}
	}
	return
}

func defaultIsClass(ext string) (isProj bool, ok bool) {
	switch ext {
	case ".gmx":
		isProj = true
		fallthrough
	case ".spx":
		ok = true
	}
	return
}

// -----------------------------------------------------------------------------

func ParseFiles(fset *token.FileSet, files []string, mode Mode) (map[string]*ast.Package, error) {
	return ParseFSFiles(fset, local, files, mode)
}

func ParseFSFiles(fset *token.FileSet, fs FileSystem, files []string, mode Mode) (map[string]*ast.Package, error) {
	ret := map[string]*ast.Package{}
	for _, file := range files {
		f, err := ParseFSFile(fset, fs, file, nil, mode)
		if err != nil {
			return nil, err
		}
		pkgName := f.Name.Name
		pkg, ok := ret[pkgName]
		if !ok {
			pkg = &ast.Package{Name: pkgName, Files: make(map[string]*ast.File)}
			ret[pkgName] = pkg
		}
		pkg.Files[file] = f
	}
	return ret, nil
}

// -----------------------------------------------------------------------------

// ParseFile parses the source code of a single Go+ source file and returns the corresponding ast.File node.
func ParseFile(fset *token.FileSet, filename string, src interface{}, mode Mode) (f *ast.File, err error) {
	return ParseFSFile(fset, local, filename, src, mode)
}

// ParseFSFile parses the source code of a single Go+ source file and returns the corresponding ast.File node.
func ParseFSFile(fset *token.FileSet, fs FileSystem, filename string, src interface{}, mode Mode) (f *ast.File, err error) {
	var code []byte
	if src == nil {
		code, err = fs.ReadFile(filename)
	} else {
		code, err = readSource(src)
	}
	if err != nil {
		return
	}
	return parseFile(fset, filename, code, mode)
}

var (
	errInvalidSource = errors.New("invalid source")
)

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
	return nil, errInvalidSource
}

// -----------------------------------------------------------------------------
