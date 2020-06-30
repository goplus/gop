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
	"go/importer"
	"go/types"
	"io"
	"log"
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

func exportConst(e *Exporter, o *types.Const) (err error) {
	return e.ExportConst(o)
}

func Import(pkgPath string) (*types.Package, error) {
	pkg, err := importer.Default().Import(pkgPath)
	if err != nil {
		return nil, err
	}
	return pkg, nil
}

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
			err = exportVar(o)
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
