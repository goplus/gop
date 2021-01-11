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

package gopkg

import (
	"bytes"
	"errors"
	"fmt"
	"go/importer"
	"go/token"
	"go/types"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/goplus/gop/mod/semver"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

func exportFunc(e *Exporter, o *types.Func, prefix string) (err error) {
	e.ExportFunc(o)
	return nil
}

func exportTypeName(e *Exporter, o *types.TypeName) (err error) {
	if t, ok := o.Type().(*types.Named); ok {
		e.ExportType(o)
		n := t.NumMethods()
		for i := 0; i < n; i++ {
			m := t.Method(i)
			if !m.Exported() {
				continue
			}
			exportFunc(e, m, "  ")
		}
	}
	return nil
}

func exportVar(e *Exporter, o *types.Var) (err error) {
	e.ExportVar(o)
	return nil
}

func exportConst(e *Exporter, o *types.Const) (err error) {
	return e.ExportConst(o)
}

// ParsePkgVer
// github.com/qiniu/x@v1.11.5/log => github.com/qiniu/x/log github.com/qiniu/x/@v1.11.5 log
func ParsePkgVer(pkgPath string) (pkg string, path string, sub string) {
	i := strings.Index(pkgPath, "@")
	if i == -1 {
		return pkgPath, "", ""
	}
	j := strings.Index(pkgPath[i:], "/")
	if j == -1 {
		return pkgPath[:i], pkgPath, ""
	}
	return pkgPath[:i] + pkgPath[i:][j:], pkgPath[:i+j], pkgPath[i+j+1:]
}

func findLastVerPkg(pkgDirBase string, name string) (verName string) {
	verName = name
	fis, err := ioutil.ReadDir(pkgDirBase)
	if err != nil {
		return
	}
	name += "@"
	var verMax string
	for _, fi := range fis {
		if fi.IsDir() {
			if v := fi.Name(); strings.HasPrefix(v, name) {
				ver := v[len(name):]
				if verMax == "" || semver.Compare(verMax, ver) < 0 {
					verName, verMax = v, ver
				}
			}
		}
	}
	return
}

const (
	pkgStandard = 0
	pkgExternal = 1
	pkgGoplus   = 2
)

const (
	github    = "github.com/"
	goplus    = "goplus/gop"
	githubLen = len(github)
	pkgGopLen = len(goplus) + githubLen
)

func getPkgKind(pkgPath string) int {
	if strings.HasPrefix(pkgPath, github) {
		if strings.HasPrefix(pkgPath[githubLen:], goplus) {
			if len(pkgPath) == pkgGopLen || pkgPath[pkgGopLen] == '/' {
				return pkgGoplus
			}
		}
		return pkgExternal
	}
	return pkgStandard
}

// LookupMod looup pkgPath srcDir from GOMOD root
// github.com/qiniu/x/log
// github.com/qiniu/x@v1.11.5/log
func LookupMod(pkgPath string) (srcDir string, err error) {
	parts := strings.Split(pkgPath, "/")
	n := len(parts)
	if n < 3 {
		return "", ErrInvalidPkgPath
	}
	srcDir = filepath.Join(getModRoot(), parts[0], parts[1])
	noVer := strings.Index(parts[2], "@") == -1
	if noVer {
		parts[2] = findLastVerPkg(srcDir, parts[2])
	}
	srcDir = filepath.Join(srcDir, parts[2])
	if n > 3 {
		srcDir = filepath.Join(srcDir, filepath.Join(parts[3:]...))
	}
	s, err := os.Stat(srcDir)
	if err != nil {
		return "", err
	}
	if !s.IsDir() {
		return "", os.ErrInvalid
	}
	return
}

// ImportSource import a Go package from pkgPath and srcDir
func ImportSource(pkgPath string, srcDir string) (*types.Package, error) {
	imp := importer.ForCompiler(token.NewFileSet(), "source", nil)
	return imp.(types.ImporterFrom).ImportFrom(pkgPath, srcDir, 0)
}

// Import imports a Go package.
func Import(pkgPath string) (*types.Package, error) {
	var imp types.Importer
	var srcDir string
	switch getPkgKind(pkgPath) {
	case pkgExternal:
		parts := strings.Split(pkgPath, "/")
		n := len(parts)
		if n < 3 {
			return nil, ErrInvalidPkgPath
		}
		srcDir = filepath.Join(getModRoot(), parts[0], parts[1])
		noVer := strings.Index(parts[2], "@") == -1
		if noVer {
			parts[2] = findLastVerPkg(srcDir, parts[2])
		} else {
			pkgPath, _, _ = ParsePkgVer(pkgPath)
		}
		srcDir = filepath.Join(srcDir, parts[2])
		if n > 3 {
			srcDir += "/" + strings.Join(parts[3:], "/")
		}
		fmt.Fprintln(os.Stderr, "export", srcDir)
		imp = importer.ForCompiler(token.NewFileSet(), "source", nil)
	case pkgGoplus:
		goplusRoot, err := GoPlusRoot()
		if err != nil {
			return nil, err
		}
		srcDir = goplusRoot + pkgPath[pkgGopLen:]
		fmt.Fprintln(os.Stderr, "export", srcDir)
		imp = importer.ForCompiler(token.NewFileSet(), "source", nil)
	default:
		imp = importer.Default()
	}
	return imp.(types.ImporterFrom).ImportFrom(pkgPath, srcDir, 0)
}

var (
	gModRoot string
	gOnce    sync.Once
)

func getModRoot() string {
	gOnce.Do(func() {
		gModRoot = os.Getenv("GOMODCACHE")
		if gModRoot == "" {
			paths := strings.Split(os.Getenv("GOPATH"), string(filepath.ListSeparator))
			gModRoot = paths[0] + "/pkg/mod"
		}
	})
	return gModRoot
}

var (
	// ErrIgnore error
	ErrIgnore = errors.New("ignore empty exported package")
	// ErrInvalidPkgPath error
	ErrInvalidPkgPath = errors.New("invalid package path")
)

// ExportPackage export types.Package to io.Writer
func ExportPackage(pkg *types.Package, w io.Writer) (err error) {
	gbl := pkg.Scope()
	names := gbl.Names()
	e := NewExporter(w, pkg)
	for _, name := range names {
		obj := gbl.Lookup(name)
		if !obj.Exported() {
			continue
		}
		switch o := obj.(type) {
		case *types.Var:
			err = exportVar(e, o)
		case *types.Const:
			err = exportConst(e, o)
		case *types.TypeName:
			err = exportTypeName(e, o)
		case *types.Func:
			err = exportFunc(e, o, "")
		default:
			log.Panicln("export failed: unknown type -", reflect.TypeOf(o))
		}
		if err != nil {
			return
		}
	}
	if e.IsEmpty() {
		return ErrIgnore
	}
	return e.Close()
}

// Export utility
func Export(pkgPath string, w io.Writer) (err error) {
	pkg, err := Import(pkgPath)
	if err != nil {
		return err
	}
	return ExportPackage(pkg, w)
}

// -----------------------------------------------------------------------------

// GoPlusRoot returns Go+ root path.
func GoPlusRoot() (root string, err error) {
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	for {
		modfile := filepath.Join(dir, "go.mod")
		if hasFile(modfile) {
			if isGoplus(modfile) {
				return dir, nil
			}
			return "", errors.New("current directory is not under goplus root")
		}
		next := filepath.Dir(dir)
		if dir == next {
			return "", errors.New("go.mod not found, please run under goplus root")
		}
		dir = next
	}
}

func isGoplus(modfile string) bool {
	b, err := ioutil.ReadFile(modfile)
	return err == nil && bytes.HasPrefix(b, goplusPrefix)
}

var (
	goplusPrefix = []byte("module github.com/goplus/gop")
)

func hasFile(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.Mode().IsRegular()
}

// -----------------------------------------------------------------------------
