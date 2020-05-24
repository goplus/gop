/*
 Copyright 2020 Qiniu Cloud (七牛云)

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
	"go/scanner"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/qiniu/qlang/v6/ast"
	"github.com/qiniu/qlang/v6/token"
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

// -----------------------------------------------------------------------------

// ParseDir calls ParseFSDir by passing a local filesystem.
//
func ParseDir(fset *token.FileSet, path string, filter func(os.FileInfo) bool, mode Mode) (pkgs map[string]*ast.Package, first error) {
	return ParseFSDir(fset, local, path, filter, mode)
}

// ParseFSDir calls ParseFile for all files with names ending in ".ql" in the
// directory specified by path and returns a map of package name -> package
// AST with all the packages found.
//
// If filter != nil, only the files with os.FileInfo entries passing through
// the filter (and ending in ".ql") are considered. The mode bits are passed
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
		if strings.HasSuffix(d.Name(), ".ql") && (filter == nil || filter(d)) {
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

// ParseFile parses the source code of a single qlang source file and returns the corresponding ast.File node.
func ParseFile(fset *token.FileSet, filename string, src interface{}, mode Mode) (f *ast.File, err error) {
	return ParseFSFile(fset, local, filename, src, mode)
}

// ParseFSFile parses the source code of a single qlang source file and returns the corresponding ast.File node.
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

func parseFileEx(fset *token.FileSet, filename string, code []byte, mode Mode) (f *ast.File, err error) {
	var b []byte
	var isMod bool
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
					b = make([]byte, n+13)
					copy(b, code[:idx])
				}
				copy(b[idx+12:], code[idx:n])
				if isMod {
					copy(b[idx:], "func init(){")
				} else {
					copy(b[idx:], "func main(){")
				}
				b[n+12] = '}'
				code = b[:n+13]
				err = nil
			}
		}
	}
	if err == nil {
		f, err = parseFile(fset, filename, code, mode)
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

func (p *parser) parseSliceLit(lbrack token.Pos, len ast.Expr) ast.Expr {
	elts := make([]ast.Expr, 1, 8)
	elts[0] = len
	for p.tok == token.COMMA {
		p.next()
		elt := p.parseRHS()
		elts = append(elts, elt)
	}
	rbrack := p.expect(token.RBRACK)
	return &ast.SliceLit{Lbrack: lbrack, Elts: elts, Rbrack: rbrack}
}

func newSliceLit(lbrack, rbrack token.Pos, len ast.Expr) ast.Expr {
	var elts []ast.Expr
	if len != nil {
		elts = []ast.Expr{len}
	}
	return &ast.SliceLit{Lbrack: lbrack, Elts: elts, Rbrack: rbrack}
}

// -----------------------------------------------------------------------------
