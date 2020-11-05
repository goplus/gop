package work

import (
	"os"
	"path/filepath"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
)

type Package struct {
	Name   string       // package name
	Dir    string       // package source dir
	Pkg    *ast.Package // ast.Package
	Target string       // target name
	Input  string       // package input argument
	Error  error        // parser Error
}

func LoadPackages(fset *token.FileSet, args []string) []*Package {
	var all []*Package
	for _, arg := range args {
		path, err := filepath.Abs(arg)
		if err != nil {
			all = append(all, &Package{Input: arg, Error: err})
			continue
		}
		fi, err := os.Stat(path)
		if err != nil {
			all = append(all, &Package{Input: arg, Error: err})
			continue
		}
		var pkgs map[string]*ast.Package
		var dir string
		_, target := filepath.Split(path)
		if fi.IsDir() {
			dir = path
			pkgs, err = parser.ParseDir(fset, path, nil, 0)
		} else {
			dir, _ = filepath.Split(path)
			dir = filepath.Clean(dir)
			ext := filepath.Ext(target)
			target = target[:len(target)-len(ext)]
			pkgs, err = parser.Parse(fset, path, nil, 0)
		}
		for name, pkg := range pkgs {
			all = append(all, &Package{Input: arg, Name: name, Dir: dir, Pkg: pkg, Target: target, Error: err})
		}
	}
	return all
}

// func ParserPackage(fset *token.FileSet, path string) (*ast.Package, error) {
// 	pkgs, err := ParserPackages(fset, path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(pkgs) != 1 {
// 		return nil, fmt.Errorf("too many packages (%d) in the same directory", len(pkgs))
// 	}
// 	return getPkg(pkgs), nil
// }

// func ParserPackages(fset *token.FileSet, path string) (pkgs map[string]*ast.Package, err error) {
// 	isDir, err := IsDir(path)
// 	if err != nil {
// 		return nil, os.ErrInvalid
// 	}
// 	if isDir {
// 		pkgs, err = parser.ParseDir(fset, path, nil, 0)
// 	} else {
// 		pkgs, err = parser.Parse(fset, path, nil, 0)
// 	}
// 	return
// }

// // IsDir checks a target path is dir or not.
// func IsDir(target string) (bool, error) {
// 	fi, err := os.Stat(target)
// 	if err != nil {
// 		return false, err
// 	}
// 	return fi.IsDir(), nil
// }

// func ParserTarget(target string) (isDir bool, targeName string, targetDir string, err error) {
// 	path, err := filepath.Abs(target)
// 	if err != nil {
// 		return
// 	}
// 	var fi os.FileInfo
// 	fi, err = os.Stat(target)
// 	if err != nil {
// 		return
// 	}
// 	targeName = fi.Name()
// 	if fi.IsDir() {
// 		isDir = true
// 		targetDir = target
// 	}
// 	return
// }

// func getPkg(pkgs map[string]*ast.Package) *ast.Package {
// 	for _, pkg := range pkgs {
// 		return pkg
// 	}
// 	return nil
// }
