package format

import (
	"bytes"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/format"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
)

var printFuncs = map[string]string{
	"Print":   "print",
	"Printf":  "printf",
	"Println": "println",
}

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

	replaceFmtPrintCalls(file)
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

func processAssignment(stmt *ast.AssignStmt) bool {
	var keep bool

	for _, expr := range stmt.Rhs {
		switch val := expr.(type) {
		case *ast.CallExpr:
			if processCallExpr(val) {
				keep = true
			}
		case *ast.SelectorExpr:
			if processSelectorExpr(val) {
				keep = true
			}
		}
	}

	return keep
}

func processCallExpr(expr *ast.CallExpr) bool {
	if selExpr, ok := expr.Fun.(*ast.SelectorExpr); ok {
		if ident, ok := selExpr.X.(*ast.Ident); ok {
			if ident.Name == "fmt" {
				p, ok := printFuncs[selExpr.Sel.Name]
				if !ok {
					return true
				}

				selExpr.Sel.Name = p
				expr.Fun = selExpr.Sel
			}
		}
	}

	return false
}

func processDeclStmt(stmt *ast.DeclStmt) bool {
	if decl, ok := stmt.Decl.(*ast.GenDecl); ok {
		for _, spec := range decl.Specs {
			if valSpec, ok := spec.(*ast.ValueSpec); ok {
				if selExpr, ok := valSpec.Type.(*ast.SelectorExpr); ok {
					if processSelectorExpr(selExpr) {
						return true
					}
				}
			}
		}
	}

	return false
}

func processExprStmt(stmt *ast.ExprStmt) bool {
	if callExpr, ok := stmt.X.(*ast.CallExpr); ok {
		return processCallExpr(callExpr)
	}

	return false
}

func processSelectorExpr(expr *ast.SelectorExpr) bool {
	if ident, ok := expr.X.(*ast.Ident); ok {
		if ident.Name == "fmt" {
			if _, ok := printFuncs[expr.Sel.Name]; !ok {
				return true
			}
		}
	}

	return false
}

// replaceFmtPrintCalls replaces fmt print funcs to builtin print funcs.
// also removes fmt import if necessary.
func replaceFmtPrintCalls(file *ast.File) {
	var keepFmtImport bool

	for _, decl := range file.Decls {
		if replacePrintCall(decl) {
			keepFmtImport = true
		}
	}

	if keepFmtImport {
		return
	}

	trimFmtImport(file)
}

func replacePrintCall(decl ast.Decl) bool {
	var keepFmtImport bool

	if fn, ok := decl.(*ast.FuncDecl); ok {
		for _, stmt := range fn.Body.List {
			switch val := stmt.(type) {
			case *ast.AssignStmt:
				if processAssignment(val) {
					keepFmtImport = true
				}
			case *ast.DeclStmt:
				if processDeclStmt(val) {
					keepFmtImport = true
				}
			case *ast.ExprStmt:
				if processExprStmt(val) {
					keepFmtImport = true
				}
			}
		}
	}

	return keepFmtImport
}

func trimFmtImport(file *ast.File) {
	var indice []int
	for i, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			if genDecl.Tok != token.IMPORT {
				continue
			}
			for _, spec := range genDecl.Specs {
				if imp, ok := spec.(*ast.ImportSpec); ok {
					if imp.Path.Value == `"fmt"` {
						indice = append(indice, i)
					}
				}
			}
		}
	}

	for _, index := range indice {
		file.Decls = append(file.Decls[:index], file.Decls[index+1:]...)
	}
}

// -----------------------------------------------------------------------------
