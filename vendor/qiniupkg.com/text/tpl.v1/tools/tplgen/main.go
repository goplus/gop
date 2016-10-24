package main

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"runtime"
	"text/template"

	"qiniupkg.com/text/tpl.v1/generator"
)

func exitIfNotNil(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// find value of ident 'grammar' in GenDecl.
func findInGenDecl(genDecl *ast.GenDecl, grammarName string) string {
	for _, spec := range genDecl.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if ok {
			index := -1
			for i, ident := range valueSpec.Names {
				if ident.Name == grammarName {
					index = i
				}
			}
			if index != -1 {
				basicLit, ok := valueSpec.Values[index].(*ast.BasicLit)
				if ok {
					return basicLit.Value
				}
			}
		}
	}
	return ""
}

func findInAssignment(assignStmt *ast.AssignStmt, grammarName string) string {
	index := -1
	for i, expr := range assignStmt.Lhs {
		ident, ok := expr.(*ast.Ident)
		if ok {
			if ident.Name == grammarName {
				index = i
				break
			}
		}
	}
	if index != -1 {
		basicLit, ok := assignStmt.Rhs[index].(*ast.BasicLit)
		if ok {
			return basicLit.Value
		}
	}
	return ""
}

func findInDecl(decl ast.Decl, grammarName string) string {
	genDecl, ok := decl.(*ast.GenDecl)
	if ok {
		g := findInGenDecl(genDecl, grammarName)
		if g != "" {
			return g
		}
	}
	funcDecl, ok := decl.(*ast.FuncDecl)
	if ok {
		for _, stmt := range funcDecl.Body.List {
			declStmt, ok := stmt.(*ast.DeclStmt)
			if ok {
				g := findInDecl(declStmt.Decl, grammarName)
				if g != "" {
					return g
				}
			}
			assignStmt, ok := stmt.(*ast.AssignStmt)
			if ok {
				g := findInAssignment(assignStmt, grammarName)
				if g != "" {
					return g
				}
			}
		}
	}
	return ""
}

func main() {
	var (
		grammarContent string
		packageName    string
	)
	grammarName := flag.String("g", "grammar", "grammar to parse")
	outputFileName := flag.String("f", "static_compiler.go", "generated static compiler file")
	genInterpreter := flag.Bool("i", false, "if generate static interpreter")
	flag.Parse()
	os.Remove(*outputFileName)
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, ".", nil, parser.ParseComments)
	exitIfNotNil(err)
	// 只做一次搜索,找到语法就结束.
	for name, pkg := range pkgs {
		for _, f := range pkg.Files {
			for _, decl := range f.Decls {
				grammarContent = findInDecl(decl, *grammarName)
				if grammarContent != "" {
					packageName = name
					goto OUT
				}
			}
		}
	}
OUT:
	if grammarContent == "" {
		exitIfNotNil(errors.New(*grammarName + " not found."))
	}
	code, err := generator.GenStaticCode(grammarContent[1 : len(grammarContent)-1])
	exitIfNotNil(err)
	fd, err := os.OpenFile(*outputFileName, os.O_RDWR|os.O_CREATE, 0664)
	exitIfNotNil(err)
	defer fd.Close()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	dir := path.Dir(filename)
	t, err := template.ParseFiles(path.Join(dir, "output.tpl"))
	exitIfNotNil(err)
	t.Execute(fd, map[string]interface{}{
		"packageName":    packageName,
		"genInterpreter": *genInterpreter,
		"code":           code,
	})
}
