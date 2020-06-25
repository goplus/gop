package main

import (
	"go/build"
	"os"
	"path/filepath"
	"syscall"
)

var (
	// BuildContext is the default build context.
	BuildContext = &build.Default
)

// PkgKind - package kind
type PkgKind int

const (
	// PkgKindMod - go module
	PkgKindMod PkgKind = 0x01
	// PkgKindStd - standard go package
	PkgKindStd PkgKind = 0x02
)

// Import a package
func Import(path string) (kind PkgKind, dir string, err error) {
	goroot := BuildContext.GOROOT
	dir = filepath.Join(goroot, "src", path)
	if isDir(dir) {
		kind = PkgKindStd
		return
	}
	return 0, "", syscall.ENOENT
}

func isDir(dir string) bool {
	fi, err := os.Stat(dir)
	return err == nil && fi.IsDir()
}
