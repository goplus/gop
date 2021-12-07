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

	replacePrintCalls(file)
}

func identEqual(v *ast.Ident, name string) bool {
	return v != nil && v.Name == name
}

func findFmtImports(decls []ast.Decl) []string {
	var aliases []string

	for _, decl := range decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			if genDecl.Tok != token.IMPORT {
				continue
			}

			for _, spec := range genDecl.Specs {
				if imp, ok := spec.(*ast.ImportSpec); ok {
					if imp.Path.Value == `"fmt"` {
						if imp.Name != nil {
							aliases = append(aliases, imp.Name.Name)
						} else {
							aliases = append(aliases, "fmt")
						}
					}
				}
			}
		}
	}

	return aliases
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

func processAssignment(stmt *ast.AssignStmt, fmts []string) []string {
	var keeps []string

	for _, expr := range stmt.Rhs {
		switch val := expr.(type) {
		case *ast.CallExpr:
			keeps = append(keeps, processCallExpr(val, fmts)...)
		case *ast.SelectorExpr:
			keeps = append(keeps, processSelectorExpr(val, fmts)...)
		}
	}

	return keeps
}

func processCallExpr(expr *ast.CallExpr, fmts []string) (keeps []string) {
	if selExpr, ok := expr.Fun.(*ast.SelectorExpr); ok {
		if ident, ok := selExpr.X.(*ast.Ident); ok {
			if contains(fmts, ident.Name) {
				p, ok := printFuncs[selExpr.Sel.Name]
				if !ok {
					return []string{ident.Name}
				}

				selExpr.Sel.Name = p
				expr.Fun = selExpr.Sel
			}
		}
	}

	return
}

func processDeclStmt(stmt *ast.DeclStmt, fmts []string) []string {
	var keeps []string
	if decl, ok := stmt.Decl.(*ast.GenDecl); ok {
		for _, spec := range decl.Specs {
			if valSpec, ok := spec.(*ast.ValueSpec); ok {
				if selExpr, ok := valSpec.Type.(*ast.SelectorExpr); ok {
					keeps = append(keeps, processSelectorExpr(selExpr, fmts)...)
				}
			}
		}
	}

	return keeps
}

func processExprStmt(stmt *ast.ExprStmt, fmts []string) (keeps []string) {
	if callExpr, ok := stmt.X.(*ast.CallExpr); ok {
		return processCallExpr(callExpr, fmts)
	}

	return
}

func processSelectorExpr(expr *ast.SelectorExpr, fmts []string) (keeps []string) {
	if ident, ok := expr.X.(*ast.Ident); ok {
		if contains(fmts, ident.Name) {
			if _, ok := printFuncs[expr.Sel.Name]; !ok {
				return []string{ident.Name}
			}
		}
	}

	return
}

// replacePrintCalls replaces fmt print funcs to builtin print funcs.
// also removes fmt import if necessary.
func replacePrintCalls(file *ast.File) {
	var keeps []string
	fmts := findFmtImports(file.Decls)
	for _, decl := range file.Decls {
		keeps = append(keeps, replacePrintCall(decl, fmts)...)
	}

	trimFmtImport(file, keeps)
}

func replacePrintCall(decl ast.Decl, fmts []string) []string {
	var keeps []string

	if fn, ok := decl.(*ast.FuncDecl); ok {
		for _, stmt := range fn.Body.List {
			switch val := stmt.(type) {
			case *ast.AssignStmt:
				keeps = append(keeps, processAssignment(val, fmts)...)
			case *ast.DeclStmt:
				keeps = append(keeps, processDeclStmt(val, fmts)...)
			case *ast.ExprStmt:
				keeps = append(keeps, processExprStmt(val, fmts)...)
			}
		}
	}

	return keeps
}

func trimFmtImport(file *ast.File, keeps []string) {
	var declIndice []int

	for i, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok != token.IMPORT {
			continue
		}

		var indice []int
		for j, spec := range genDecl.Specs {
			if imp, ok := spec.(*ast.ImportSpec); ok {
				if imp.Path.Value == `"fmt"` {
					// has alias
					if imp.Name != nil && !contains(keeps, imp.Name.Name) {
						indice = append(indice, j)
					} else if !contains(keeps, "fmt") {
						indice = append(indice, j)
					}
				}
			}
		}

		genDecl.Specs = removeSpecsIf(genDecl.Specs, func(i int) bool {
			for _, val := range indice {
				if i == val {
					return true
				}
			}
			return false
		})
		if len(genDecl.Specs) == 0 {
			declIndice = append(declIndice, i)
		}
	}

	file.Decls = removeDeclsIf(file.Decls, func(i int) bool {
		for _, val := range declIndice {
			if i == val {
				return true
			}
		}
		return false
	})
}

// -----------------------------------------------------------------------------

func contains(items []string, val string) bool {
	for _, item := range items {
		if val == item {
			return true
		}
	}
	return false
}

func removeDeclsIf(decls []ast.Decl, fn func(i int) bool) []ast.Decl {
	ret := decls[:0]
	for i := range decls {
		if fn(i) {
			continue
		}
		ret = append(ret, decls[i])
	}
	return ret
}

func removeSpecsIf(specs []ast.Spec, fn func(i int) bool) []ast.Spec {
	ret := specs[:0]
	for i := range specs {
		if fn(i) {
			continue
		}
		ret = append(ret, specs[i])
	}
	return ret
}
