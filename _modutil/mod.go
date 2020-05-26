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

package modutil

import (
	"go/build"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/visualfc/fastmod"
)

var (
	// BuildContext is the default build context.
	BuildContext = &build.Default
)

// -----------------------------------------------------------------------------

// gPkgModPath - base path that `go get` places all versioned packages.
var gPkgModPath string
var gGorootSrcPath string

// LookupModFile finds go.mod file for a package directory.
func LookupModFile(dir string) (file string, err error) {
	file, err = fastmod.LookupModFile(dir)
	if err != nil {
		return
	}
	if file == "" {
		err = syscall.ENOENT
	}
	return
}

// -----------------------------------------------------------------------------

type modules = *fastmod.ModuleList

var gMods modules
var onceMods sync.Once

func getModules() modules {
	if gMods == nil {
		onceMods.Do(func() {
			gPkgModPath = fastmod.GetPkgModPath(BuildContext) + "/"
			gGorootSrcPath = BuildContext.GOROOT + "/src/"
			gMods = fastmod.NewModuleList(BuildContext)
		})
	}
	return gMods
}

// LoadModule loads a module from specified dir.
func LoadModule(dir string) (mod *Module, err error) {
	impl, err := getModules().LoadModule(dir)
	if err != nil {
		return
	}
	isMod := strings.HasPrefix(dir, gPkgModPath)
	modKind := 0
	if isMod {
		modKind = modKindMod
	} else if impl.Path() == "std" {
		modKind = modKindStd
	}
	return &Module{impl, dir, modKind}, nil
}

// -----------------------------------------------------------------------------

const (
	modKindMod = 0x01
	modKindStd = 0x02
)

// Module represents a loaded module.
type Module struct {
	impl    *fastmod.Module
	dir     string
	modKind int
}

// Lookup returns package info, if found.
func (p Module) Lookup(pkg string) (pi PackageInfo, err error) {
	dir := filepath.Join(BuildContext.GOROOT, "src", pkg)
	if isDirExist(dir) {
		return PackageInfo{VersionPkgPath: pkg, Location: dir, Type: PkgTypeStd}, nil
	}
	_, dir, typ := p.impl.Lookup(pkg)
	if typ == fastmod.PkgTypeNil {
		err = syscall.ENOENT
		return
	}
	if typ == PkgTypeDepMod {
		pkg = strings.TrimPrefix(dir, gPkgModPath)
	}
	return PackageInfo{VersionPkgPath: pkg, Location: dir, Type: typ}, nil
}

// ModFile returns `go.mod` file path of this module. eg. `$HOME/work/qiniu/qlang/go.mod`
func (p Module) ModFile() string {
	return p.impl.ModFile()
}

// RootPath returns root path of this module. eg. `$HOME/work/qiniu/qlang`
func (p Module) RootPath() string {
	return p.impl.ModDir()
}

// PkgPath returns PkgPath of this module. eg. `github.com/qiniu/qlang/v6`
func (p Module) PkgPath() string {
	switch p.modKind {
	case modKindStd:
		return strings.TrimPrefix(p.dir, gGorootSrcPath)
	default:
		return p.impl.Path()
	}
}

// VersionPkgPath returns VersionPkgPath of this module. eg. `github.com/qiniu/qlang/v6@v1.0.0`
func (p Module) VersionPkgPath() string {
	switch p.modKind {
	case modKindMod:
		return strings.TrimPrefix(p.impl.ModDir(), gPkgModPath)
	case modKindStd:
		return strings.TrimPrefix(p.dir, gGorootSrcPath)
	default:
		return p.impl.Path()
	}
}

func isDirExist(dir string) bool {
	fi, err := os.Stat(dir)
	return err == nil && fi.IsDir()
}

// -----------------------------------------------------------------------------

// PkgType represents a package type.
type PkgType = fastmod.PkgType

const (
	// PkgTypeStd - a std module
	PkgTypeStd = fastmod.PkgTypeGoroot

	// PkgTypeGopath - a module found at $GOPATH/src
	PkgTypeGopath = fastmod.PkgTypeGoroot

	// PkgTypeThis - this module itself
	PkgTypeThis = fastmod.PkgTypeMod

	// PkgTypeChild - child module of this module
	PkgTypeChild = fastmod.PkgTypeLocal

	// PkgTypeDepMod - a depended module found at $GOPATH/pkg/mod
	PkgTypeDepMod = fastmod.PkgTypeDepMod

	// PkgTypeLocalDep - a module that rewrites to local
	PkgTypeLocalDep = fastmod.PkgTypeLocalMod
)

// PackageInfo represents a package info.
type PackageInfo struct {
	Location       string
	VersionPkgPath string
	Type           PkgType
}

// -----------------------------------------------------------------------------
