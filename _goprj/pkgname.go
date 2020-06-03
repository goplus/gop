/*
 Copyright 2020 Qiniu Cloud (qiniu.com)

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

package goprj

import (
	"errors"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"

	"github.com/qiniu/goplus/modutil"
	"github.com/qiniu/x/log"
)

var (
	// ErrMultiPackages error.
	ErrMultiPackages = errors.New("multi packages in the same dir")
)

func filterTest(fi os.FileInfo) bool {
	return !strings.HasSuffix(fi.Name(), "_test.go")
}

// GetPkgName returns package name of a module by specified directory.
func GetPkgName(dir string) (string, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, filterTest, parser.PackageClauseOnly)
	if err != nil {
		return "", err
	}
	if len(pkgs) != 1 {
		delete(pkgs, "main")
		if len(pkgs) != 1 {
			return "", ErrMultiPackages
		}
	}
	for name := range pkgs {
		return name, nil
	}
	panic("not here")
}

// LookupPkgName lookups a package name by specified PkgPath.
func LookupPkgName(prjMod *modutil.Module, pkg string) (name string, err error) {
	pi, err := prjMod.Lookup(pkg)
	if err != nil {
		return
	}
	if pi.Type == modutil.PkgTypeStd {
		return path.Base(pkg), nil
	}
	return GetPkgName(pi.Location)
}

// -----------------------------------------------------------------------------

// PkgNames keeps a package name cache: pkg => pkgName
type PkgNames struct {
	names map[string]string // pkg => pkgName
}

// NewPkgNames creates a new PkgNames object.
func NewPkgNames() PkgNames {
	return PkgNames{
		names: make(map[string]string),
	}
}

// LookupPkgName lookups a package name by specified PkgPath.
func (p PkgNames) LookupPkgName(prjMod *modutil.Module, pkg string) string {
	if pkg == "C" { // import "C"
		return pkg
	}
	if name, ok := p.names[pkg]; ok {
		return name
	}
	name, err := LookupPkgName(prjMod, pkg)
	if err != nil {
		log.Fatalln("LookupPkgName failed:", err, "-", pkg)
	}
	p.names[pkg] = name
	return name
}

// -----------------------------------------------------------------------------
