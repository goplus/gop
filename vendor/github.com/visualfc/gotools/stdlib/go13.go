// +build !go1.4

package stdlib

import (
	"go/build"
	"os"
	"path/filepath"
)

func ImportStdPkg(context *build.Context, path string, mode build.ImportMode) (*build.Package, error) {
	realpath := filepath.Join(context.GOROOT, "src", "pkg", path)
	if _, err := os.Stat(realpath); err != nil {
		realpath = filepath.Join(context.GOROOT, "src", path)
	}
	pkg, err := context.ImportDir(realpath, 0)
	pkg.ImportPath = path
	return pkg, err
}
