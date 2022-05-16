/*
 * Copyright (c) 2022 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"flag"
	"fmt"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	goast "go/ast"
	goparser "go/parser"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/x/types"
	"github.com/goplus/gox/packages"
)

var (
	verboseFlag = flag.Bool("v", false, "print verbose information")
	internal    = flag.Bool("i", false, "print internal declarations")
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

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		usage()
		return
	}
	verbose := *verboseFlag
	if verbose {
		log.Println("==> Parsing ...")
	}

	var fset = token.NewFileSet()
	var in *ast.Package
	if infile := flag.Arg(0); isDir(infile) {
		pkgs, first := parser.ParseDir(fset, infile, nil, parser.ParseGoFilesByGo)
		check(first)
		for name, pkg := range pkgs {
			if !strings.HasSuffix(name, "_test") {
				in = pkg
				break
			}
		}
	} else {
		in = &ast.Package{
			Files:   make(map[string]*ast.File),
			GoFiles: make(map[string]*goast.File),
		}
		for i, n := 0, flag.NArg(); i < n; i++ {
			infile = flag.Arg(i)
			switch filepath.Ext(infile) {
			case ".gop":
				f, err := parser.ParseFile(fset, infile, nil, 0)
				check(err)
				in.Files[infile] = f
			case ".go":
				f, err := goparser.ParseFile(fset, infile, nil, 0)
				check(err)
				in.GoFiles[infile] = f
			default:
				log.Panicln("Unknown support file:", infile)
			}
		}
	}

	if verbose {
		log.Println("==> Loading ...")
	}
	conf := &types.Config{
		Importer:         packages.NewImporter(fset),
		IgnoreFuncBodies: true,
		Error: func(err error) {
			log.Println(err)
		},
		DisableUnusedImportCheck: true,
	}
	pkg, err := types.Load(fset, in, conf)
	check(err)

	if verbose {
		log.Println("==> Printing ...")
	}
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
