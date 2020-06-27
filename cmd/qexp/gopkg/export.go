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
	"go/importer"
	"go/types"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
)

// -----------------------------------------------------------------------------

func exportFunc(e *Exporter, o *types.Func, prefix string) (err error) {
	e.ExportFunc(o)
	return nil
}

func exportTypeName(e *Exporter, o *types.TypeName) (err error) {
	t := o.Type().(*types.Named)
	n := t.NumMethods()
	for i := 0; i < n; i++ {
		m := t.Method(i)
		if !m.Exported() {
			continue
		}
		exportFunc(e, m, "  ")
	}
	return nil
}

func exportVar(o *types.Var) (err error) {
	return nil
}

func exportConst(o *types.Const) (err error) {
	return nil
}

// Export utility
func Export(pkgPath string, createFile func(string) (io.WriteCloser, error)) (err error) {
	imp := importer.Default()
	pkg, err := imp.Import(pkgPath)
	if err != nil {
		return
	}
	goplus, err := goplusRoot()
	if err != nil {
		return
	}
	pkgDir := filepath.Join(goplus, "lib", pkg.Path())
	f, err := createFile(pkgDir)
	if err != nil {
		return
	}
	defer f.Close()
	gbl := pkg.Scope()
	names := gbl.Names()
	e := NewExporter(f, pkg)
	for _, name := range names {
		obj := gbl.Lookup(name)
		if !obj.Exported() {
			continue
		}
		switch o := obj.(type) {
		case *types.Var:
			err = exportVar(o)
		case *types.Const:
			err = exportConst(o)
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
	return e.Close()
}

// -----------------------------------------------------------------------------

func goplusRoot() (root string, err error) {
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
	goplusPrefix = []byte("module github.com/qiniu/goplus")
)

func hasFile(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.Mode().IsRegular()
}

// -----------------------------------------------------------------------------
