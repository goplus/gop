/*
 Copyright 2020 The GoPlus Authors (goplus.org)

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package parser

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/qiniu/goplus/ast"
	"github.com/qiniu/goplus/scanner"
	"github.com/qiniu/goplus/token"
)

// -----------------------------------------------------------------------------

// FileSystem represents a file system.
type FileSystem interface {
	ReadDir(dirname string) ([]os.FileInfo, error)
	ReadFile(filename string) ([]byte, error)
	Join(elem ...string) string
}

type localFS struct{}

func (p localFS) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

func (p localFS) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func (p localFS) Join(elem ...string) string {
	return filepath.Join(elem...)
}

var local FileSystem = localFS{}

// ParseGopFiles parses the Go+ source files under directory or single Go+ source file.
// The target specifies the directory or single Go+ source file.
//
// The ParseGopFiles should return the map of packages to run Go+ script, even the target is single file.
//
// If the file or directory couldn't be read, a nil map and the respective error are
// returned.
// If the target is directory and a parse error occurred, a non-nil but incomplete map and the first error encountered are returned.
func ParseGopFiles(fset *token.FileSet, target string, isDir bool, mode Mode) (pkgs map[string]*ast.Package, err error) {
	if !isDir {
		file, err := ParseFile(fset, target, nil, mode)
		if err != nil {
			return pkgs, err
		}
		pkgs = make(map[string]*ast.Package)
		pkg := astFileToPkg(file, target)
		pkgs[file.Name.Name] = pkg
		return pkgs, nil
	}
	return ParseDir(fset, target, nil, mode)
}

// -----------------------------------------------------------------------------
// ParseDir calls ParseFSDir by passing a local filesystem.
//
func ParseDir(fset *token.FileSet, path string, filter func(os.FileInfo) bool, mode Mode) (pkgs map[string]*ast.Package, first error) {
	return ParseFSDir(fset, local, path, filter, mode)
}

// ParseFSDir calls ParseFile for all files with names ending in ".gop" in the
// directory specified by path and returns a map of package name -> package
// AST with all the packages found.
//
// If filter != nil, only the files with os.FileInfo entries passing through
// the filter (and ending in ".gop") are considered. The mode bits are passed
// to ParseFile unchanged. Position information is recorded in fset, which
// must not be nil.
//
// If the directory couldn't be read, a nil map and the respective error are
// returned. If a parse error occurred, a non-nil but incomplete map and the
// first error encountered are returned.
//
func ParseFSDir(fset *token.FileSet, fs FileSystem, path string, filter func(os.FileInfo) bool, mode Mode) (pkgs map[string]*ast.Package, first error) {
	list, err := fs.ReadDir(path)
	if err != nil {
		return nil, err
	}
	pkgs = make(map[string]*ast.Package)
	for _, d := range list {
		if !strings.HasPrefix(d.Name(), "_") && strings.HasSuffix(d.Name(), ".gop") && (filter == nil || filter(d)) {
			filename := fs.Join(path, d.Name())
			if filedata, err := fs.ReadFile(filename); err == nil {
				if src, err := ParseFSFile(fset, fs, filename, filedata, mode); err == nil {
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
	return parseFileEx(fset, filename, code, mode)
}

// TODO: should not add package info and init|main function.
// If do this, parsing will display error line number when error occur
func parseFileEx(fset *token.FileSet, filename string, code []byte, mode Mode) (f *ast.File, err error) {
	var b []byte
	var isMod, hasUnnamed bool
	var fsetTmp = token.NewFileSet()
	f, err = parseFile(fsetTmp, filename, code, PackageClauseOnly)
	if err != nil {
		n := len(code)
		b = make([]byte, n+28)
		copy(b, "package main;")
		copy(b[13:], code)
		code = b[:n+13]
	} else {
		isMod = f.Name.Name != "main"
	}
	_, err = parseFile(fsetTmp, filename, code, mode)
	if err != nil {
		if errlist, ok := err.(scanner.ErrorList); ok {
			if e := errlist[0]; strings.HasPrefix(e.Msg, "expected declaration") {
				n := len(code)
				idx := e.Pos.Offset
				if b == nil {
					b = make([]byte, n+14)
					copy(b, code[:idx])
				}
				copy(b[idx+12:], code[idx:n])
				if isMod {
					copy(b[idx:], "func init(){")
				} else {
					copy(b[idx:], "func main(){")
				}
				b[n+12] = '\n'
				b[n+13] = '}'
				code = b[:n+14]
				hasUnnamed = true
				err = nil
			}
		}
	}
	if err == nil {
		f, err = parseFile(fset, filename, code, mode)
		if err == nil {
			f.HasUnnamed = hasUnnamed
		}
	}
	return
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
