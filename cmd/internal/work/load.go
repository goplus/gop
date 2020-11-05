package work

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
)

type Package struct {
	Name      string       // package name
	Dir       string       // package source dir
	Pkg       *ast.Package // ast.Package
	Target    string       // build target
	Input     string       // input argument
	GenGoFile string       // generate go file path
}

func LoadPackages(fset *token.FileSet, args []string) (pkgs []*Package, errs []error) {
	for _, arg := range args {
		path, err := filepath.Abs(arg)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		fi, err := os.Stat(path)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		var apkgs map[string]*ast.Package
		var dir string
		var target string
		if fi.IsDir() {
			dir = path
			_, target = filepath.Split(path)
			apkgs, err = parser.ParseDir(fset, path, nil, 0)
		} else {
			dir, target = filepath.Split(path)
			dir = filepath.Clean(dir)
			ext := filepath.Ext(target)
			target = target[:len(target)-len(ext)]
			apkgs, err = parser.Parse(fset, path, nil, 0)
		}
		if runtime.GOOS == "window" {
			target += ".exe"
		}
		if err != nil {
			errs = append(errs, err)
			continue
		}
		for name, pkg := range apkgs {
			var mainTraget string
			if pkg.Name == "main" {
				mainTraget = target
			}
			pkgs = append(pkgs, &Package{Input: arg, Name: name, Dir: dir, Pkg: pkg, Target: mainTraget})
		}
	}
	return
}
