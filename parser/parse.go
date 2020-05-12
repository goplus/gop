package parser

import (
	"bytes"
	"errors"
	"go/parser"
	"go/scanner"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/qiniu/qlang/ast"
	"github.com/qiniu/qlang/token"
)

// A Mode value is a set of flags (or 0). They control the amount of source code parsed
// and other optional parser functionality.
type Mode = parser.Mode

const (
	// PackageClauseOnly - stop parsing after package clause
	PackageClauseOnly = parser.PackageClauseOnly

	// ImportsOnly - stop parsing after import declarations
	ImportsOnly = parser.ImportsOnly

	// ParseComments - parse comments and add them to AST
	ParseComments = parser.ParseComments

	// Trace - print a trace of parsed productions
	Trace = parser.Trace

	// DeclarationErrors - report declaration errors
	DeclarationErrors = parser.DeclarationErrors

	// AllErrors - report all errors (not just the first 10 on different lines)
	AllErrors = parser.AllErrors
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
				if src, err := ParseFile(fset, filename, filedata, mode); err == nil {
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
	return parseFile(fset, filename, code, mode)
}

func parseFile(fset *token.FileSet, filename string, code []byte, mode Mode) (f *ast.File, err error) {
	var b []byte
	var isMod bool
	var fsetTmp = token.NewFileSet()
	f, err = parser.ParseFile(fsetTmp, filename, code, PackageClauseOnly)
	if err != nil {
		n := len(code)
		b = make([]byte, n+28)
		copy(b, "package main;")
		copy(b[13:], code)
		code = b[:n+13]
	} else {
		isMod = f.Name.Name != "main"
	}
	_, err = parser.ParseFile(fsetTmp, filename, code, mode)
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
		f, err = parser.ParseFile(fset, filename, code, mode)
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
