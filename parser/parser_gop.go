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
	"os"
	"path"
	"strings"

	goast "go/ast"
	goparser "go/parser"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/parser/fsx"
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

var (
	ErrUnknownFileKind = errors.New("unknown file kind")
)

type FileSystem = fsx.FileSystem

// -----------------------------------------------------------------------------

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
func ParseDir(fset *token.FileSet, path string, filter func(fs.FileInfo) bool, mode Mode) (pkgs map[string]*ast.Package, first error) {
	return ParseFSDir(fset, fsx.Local, path, Config{Filter: filter, Mode: mode})
}

type Config struct {
	ClassKind func(fname string) (isProj, ok bool)
	Filter    func(fs.FileInfo) bool
	Mode      Mode
}

// ParseDirEx calls ParseFSDir by passing a local filesystem.
func ParseDirEx(fset *token.FileSet, path string, conf Config) (pkgs map[string]*ast.Package, first error) {
	return ParseFSDir(fset, fsx.Local, path, conf)
}

// ParseEntry calls ParseFSEntry by passing a local filesystem.
func ParseEntry(fset *token.FileSet, filename string, src interface{}, conf Config) (f *ast.File, err error) {
	return ParseFSEntry(fset, fsx.Local, filename, src, conf)
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
func ParseFSDir(fset *token.FileSet, fs FileSystem, dir string, conf Config) (pkgs map[string]*ast.Package, first error) {
	if conf.Mode&SaveAbsFile != 0 {
		dir, _ = fs.Abs(dir)
		conf.Mode &^= SaveAbsFile
	}
	list, err := fs.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	if conf.ClassKind == nil {
		conf.ClassKind = defaultClassKind
	}
	pkgs = make(map[string]*ast.Package)
	for _, d := range list {
		if d.IsDir() {
			continue
		}
		fname := d.Name()
		ext := path.Ext(fname)
		var isProj, isClass, isNormalGox, useGoParser bool
		switch ext {
		case ".gop":
		case ".go":
			if strings.HasPrefix(fname, "gop_autogen") {
				continue
			}
			useGoParser = (conf.Mode & ParseGoAsGoPlus) == 0
		case ".gox":
			isNormalGox = true
			fallthrough
		default:
			if isProj, isClass = conf.ClassKind(fname); isClass {
				isNormalGox = false
			} else if isNormalGox { // not found Go+ class by ext, but is a .gox file
				isClass = true
			} else { // unknown fileKind
				continue
			}
		}
		mode := conf.Mode
		if isClass {
			mode |= ParseGoPlusClass
		}
		if !strings.HasPrefix(fname, "_") && (conf.Filter == nil || filter(d, conf.Filter)) {
			filename := fs.Join(dir, fname)
			if useGoParser {
				if filedata, err := fs.ReadFile(filename); err == nil {
					if src, err := goparser.ParseFile(fset, filename, filedata, goparser.Mode(conf.Mode&goReservedFlags)); err == nil {
						pkg := reqPkg(pkgs, src.Name.Name)
						if pkg.GoFiles == nil {
							pkg.GoFiles = make(map[string]*goast.File)
						}
						pkg.GoFiles[filename] = src
					} else {
						first = err
					}
				} else {
					first = err
				}
			} else {
				f, err := ParseFSFile(fset, fs, filename, nil, mode)
				if f != nil {
					f.IsProj, f.IsClass = isProj, isClass
					f.IsNormalGox = isNormalGox
					if f.Name != nil {
						pkg := reqPkg(pkgs, f.Name.Name)
						pkg.Files[filename] = f
					}
				}
				if err != nil && first == nil {
					first = err
				}
			}
		}
	}
	return
}

// ParseFSEntry parses the source code of a single Go+ source file and returns the corresponding ast.File node.
// Compared to ParseFSFile, ParseFSEntry detects fileKind by its filename.
func ParseFSEntry(fset *token.FileSet, fs FileSystem, filename string, src interface{}, conf Config) (f *ast.File, err error) {
	fname := fs.Base(filename)
	ext := path.Ext(fname)
	var isProj, isClass, isNormalGox bool
	switch ext {
	case ".gop", ".go":
	case ".gox":
		isNormalGox = true
		fallthrough
	default:
		if conf.ClassKind == nil {
			conf.ClassKind = defaultClassKind
		}
		if isProj, isClass = conf.ClassKind(fname); isClass {
			isNormalGox = false
		} else if isNormalGox { // not found Go+ class by ext, but is a .gox file
			isClass = true
		} else {
			return nil, ErrUnknownFileKind
		}
	}
	mode := conf.Mode
	if isClass {
		mode |= ParseGoPlusClass
	}
	f, err = ParseFSFile(fset, fs, filename, src, mode)
	if f != nil {
		f.IsProj, f.IsClass = isProj, isClass
		f.IsNormalGox = isNormalGox
	}
	return
}

func filter(d fs.DirEntry, fn func(fs.FileInfo) bool) bool {
	fi, err := d.Info()
	return err != nil || fn(fi)
}

func reqPkg(pkgs map[string]*ast.Package, name string) *ast.Package {
	pkg, found := pkgs[name]
	if !found {
		pkg = &ast.Package{
			Name:  name,
			Files: make(map[string]*ast.File),
		}
		pkgs[name] = pkg
	}
	return pkg
}

func defaultClassKind(fname string) (isProj bool, ok bool) {
	ext := path.Ext(fname)
	switch ext {
	case ".gmx":
		return true, true
	case ".spx":
		return fname == "main.spx", true
	}
	return
}

// -----------------------------------------------------------------------------

func ParseFiles(fset *token.FileSet, files []string, mode Mode) (map[string]*ast.Package, error) {
	return ParseFSFiles(fset, fsx.Local, files, mode)
}

func ParseFSFiles(fset *token.FileSet, fs FileSystem, files []string, mode Mode) (map[string]*ast.Package, error) {
	ret := map[string]*ast.Package{}
	fabs := (mode & SaveAbsFile) != 0
	if fabs {
		mode &^= SaveAbsFile
	}
	for _, file := range files {
		if fabs {
			file, _ = fs.Abs(file)
		}
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
	return ParseFSFile(fset, fsx.Local, filename, src, mode)
}

// ParseFSFile parses the source code of a single Go+ source file and returns the corresponding ast.File node.
func ParseFSFile(fset *token.FileSet, fs FileSystem, filename string, src interface{}, mode Mode) (f *ast.File, err error) {
	code, err := readSourceFS(fs, filename, src)
	if err != nil {
		return
	}
	if mode&SaveAbsFile != 0 {
		filename, _ = fs.Abs(filename)
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
		return io.ReadAll(s)
	}
	return nil, errInvalidSource
}

// If src != nil, readSource converts src to a []byte if possible;
// otherwise it returns an error. If src == nil, readSource returns
// the result of reading the file specified by filename.
func readSourceLocal(filename string, src interface{}) ([]byte, error) {
	if src != nil {
		return readSource(src)
	}
	return os.ReadFile(filename)
}

func readSourceFS(fs FileSystem, filename string, src interface{}) ([]byte, error) {
	if src != nil {
		return readSource(src)
	}
	return fs.ReadFile(filename)
}

// -----------------------------------------------------------------------------
