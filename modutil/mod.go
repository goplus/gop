package modutil

import (
	"go/build"
	"os"
	"path/filepath"
	"syscall"

	"github.com/visualfc/fastmod"
)

var (
	// BuildContext is the default build context.
	BuildContext = &build.Default
)

// -----------------------------------------------------------------------------

// GetPkgModPath returns base path that `go get` places all versioned packages.
func GetPkgModPath() string {
	return fastmod.GetPkgModPath(BuildContext)
}

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

// Modules represents a fast module loader.
type Modules struct {
	impl *fastmod.ModuleList
}

// GetModules returns a fast module loader.
func GetModules() Modules {
	impl := fastmod.NewModuleList(BuildContext)
	return Modules{impl}
}

// Load loads a module from specified dir.
func (p Modules) Load(dir string) (mod Module, err error) {
	impl, err := p.impl.LoadModule(dir)
	if err != nil {
		return
	}
	return Module{impl}, nil
}

// -----------------------------------------------------------------------------

// Module represents a loaded module.
type Module struct {
	impl *fastmod.Module
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
	return PackageInfo{Location: dir, Type: typ}, nil
}

// ModFile returns `go.mod` file path of this module.
func (p Module) ModFile() string {
	return p.impl.ModFile()
}

// Path returns root path of this module.
func (p Module) Path() string {
	return p.impl.ModDir()
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
	Type     PkgType
}

// -----------------------------------------------------------------------------
