package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/goplus/gop/ast/togo"
	gopp "github.com/goplus/gop/parser"
	"golang.org/x/tools/go/gcexportdata"
)

var (
	compiler = flag.String("c", "source", `compiler to use: source, gc or gccgo`)
	internal = flag.Bool("i", false, "print internal declarations")
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: gopdecl [-i] [source.gop ...]\n")
	flag.PrintDefaults()
}

func isDir(name string) bool {
	if fi, err := os.Lstat(name); err == nil {
		return fi.IsDir()
	}
	return false
}

func isPublic(name string) bool {
	for _, c := range name {
		return unicode.IsUpper(c)
	}
	return false
}

var (
	goModCache string
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		usage()
		return
	}
	goModCache = initGoModCache()

	var fset = token.NewFileSet()
	var files []*ast.File
	var imp types.Importer
	switch *compiler {
	case "gc":
		packages := make(map[string]*types.Package)
		imp = gcexportdata.NewImporter(fset, packages)
	case "source":
		imp = importer.ForCompiler(fset, "source", nil)
	default:
		log.Panicln("Unsupport compiler:", *compiler)
	}

	if infile := flag.Arg(0); isDir(infile) {
		pkgs, first := gopp.ParseDir(fset, infile, func(fi fs.FileInfo) bool {
			name := fi.Name()
			if strings.HasPrefix(name, "_") {
				return false
			}
			return !strings.HasSuffix(strings.TrimSuffix(name, filepath.Ext(name)), "_test")
		}, 0)
		check(first)
		for name, pkg := range pkgs {
			if !strings.HasSuffix(name, "_test") {
				for _, f := range pkg.Files {
					files = append(files, togo.ASTFile(f, 0))
				}
				break
			}
		}
	} else {
		for i, n := 0, flag.NArg(); i < n; i++ {
			infile = flag.Arg(i)
			switch filepath.Ext(infile) {
			case ".gop":
				f, err := gopp.ParseFile(fset, infile, nil, 0)
				check(err)
				files = append(files, togo.ASTFile(f, 0))
			case ".go":
				f, err := parser.ParseFile(fset, infile, nil, 0)
				check(err)
				files = append(files, f)
			default:
				log.Panicln("Unknown support file:", infile)
			}
		}
	}

	conf := &types.Config{
		Importer:         imp,
		IgnoreFuncBodies: true,
		Error: func(err error) {
			log.Println(err)
		},
		DisableUnusedImportCheck: true,
	}

	pkg, err := conf.Check("", fset, files, nil)
	check(err)

	scope := pkg.Scope()
	names := scope.Names()
	for _, name := range names {
		if *internal || isPublic(name) {
			fmt.Println(scope.Lookup(name))
		}
	}
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
