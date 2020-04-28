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
			gMods = fastmod.NewModuleList(BuildContext)
		})
	}
	return gMods
}

// LoadModule loads a module from specified dir.
func LoadModule(dir string) (mod Module, err error) {
	impl, err := getModules().LoadModule(dir)
	if err != nil {
		return
	}
	isMod := strings.HasPrefix(dir, gPkgModPath)
	return Module{impl, isMod}, nil
}

// -----------------------------------------------------------------------------

// Module represents a loaded module.
type Module struct {
	impl  *fastmod.Module
	isMod bool
}

// Lookup returns package info, if found.
func (p Module) Lookup(pkg string) (pi PackageInfo, err error) {
	dir := filepath.Join(BuildContext.GOROOT, "src", pkg)
	if isDirExist(dir) {
		return PackageInfo{Location: dir, Type: PkgTypeStd}, nil
	}
	_, dir, typ := p.impl.Lookup(pkg)
	if typ == fastmod.PkgTypeNil {
		err = syscall.ENOENT
		return
	}
	if typ == PkgTypeDepMod {
		pkg = strings.TrimPrefix(dir, gPkgModPath)
	}
	return PackageInfo{PkgPath: pkg, Location: dir, Type: typ}, nil
}

// ModFile returns `go.mod` file path of this module. eg. `$HOME/work/qiniu/qlang/go.mod`
func (p Module) ModFile() string {
	return p.impl.ModFile()
}

// RootPath returns root path of this module. eg. `$HOME/work/qiniu/qlang`
func (p Module) RootPath() string {
	return p.impl.ModDir()
}

// PkgPath returns PkgPath (with version) of this module. eg. `github.com/qiniu/qlang@v1.0.0`
func (p Module) PkgPath() string {
	if p.isMod {
		return strings.TrimPrefix(p.impl.ModDir(), gPkgModPath)
	}
	return p.impl.Path()
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
	Location string
	PkgPath  string
	Type     PkgType
}

// -----------------------------------------------------------------------------
