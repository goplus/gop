package format

import (
	"github.com/goplus/gop/ast"
)

// -----------------------------------------------------------------------------

func Gopstyle(file *ast.File) {
	if identEqual(file.Name, "main") {
		file.NoPkgDecl = true
	}
	if idx := findFuncDecl(file.Decls, "main"); idx >= 0 {
		last := len(file.Decls) - 1
		if idx != last {
			fn := file.Decls[idx]
			copy(file.Decls[idx:], file.Decls[idx+1:])
			file.Decls[last] = fn
		}
		file.NoEntrypoint = true
	}
}

func identEqual(v *ast.Ident, name string) bool {
	return v != nil && v.Name == name
}

func findFuncDecl(decls []ast.Decl, name string) int {
	for i, decl := range decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if identEqual(fn.Name, name) {
				return i
			}
		}
	}
	return -1
}

// -----------------------------------------------------------------------------
