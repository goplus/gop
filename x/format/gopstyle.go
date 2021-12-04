package format

import (
	"bytes"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/format"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
)

// -----------------------------------------------------------------------------

func GopstyleSource(src []byte, filename ...string) (ret []byte, err error) {
	var fname string
	if filename != nil {
		fname = filename[0]
	}
	fset := token.NewFileSet()
	var f *ast.File
	if f, err = parser.ParseFile(fset, fname, src, parser.ParseComments); err == nil {
		Gopstyle(f)
		var buf bytes.Buffer
		if err = format.Node(&buf, fset, f); err == nil {
			ret = buf.Bytes()
		}
	}
	return
}

// -----------------------------------------------------------------------------

func Gopstyle(file *ast.File) {
	if identEqual(file.Name, "main") {
		file.NoPkgDecl = true
	}
	if idx := findFuncDecl(file.Decls, "main"); idx >= 0 {
		last := len(file.Decls) - 1
		if idx == last {
			file.NoEntrypoint = true
			// TODO: idx != last: swap main func to last
			// TODO: should also swap file.Comments
			/*
				fn := file.Decls[idx]
				copy(file.Decls[idx:], file.Decls[idx+1:])
				file.Decls[last] = fn
			*/
		}
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
