/*
 Copyright 2020 The GoPlus Authors (goplus.org)

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package parser

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/ast/asttest"
	"github.com/goplus/gop/token"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

var fsTestStd = asttest.NewSingleFileFS("/foo", "bar.gop", `package bar; import "io"
	// comment
	x := 0
	if t := false; t {
		x = 3
	} else {
		x = 5
	}
	println("x:", x)

	// comment 1
	// comment 2
	x = 0
	switch s := "Hello"; s {
	default:
		x = 7
	case "world", "hi":
		x = 5
	case "xsw":
		x = 3
	}
	println("x:", x)

	c := make(chan bool, 100)
	select {
	case c <- true:
	case v := <-c:
	default:
		panic("error")
	}
`)

func TestStd(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := ParseFSDir(fset, fsTestStd, "/foo", nil, ParseComments)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}
	bar := pkgs["bar"]
	file := bar.Files["/foo/bar.gop"]
	fmt.Println("Pkg:", file.Name)
	for _, decl := range file.Decls {
		fmt.Println("decl:", reflect.TypeOf(decl))
		switch d := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch v := spec.(type) {
				case *ast.ImportSpec:
					fmt.Println(" - import:", v.Path.Value)
				}
			}
		case *ast.FuncDecl:
			fmt.Println(" - func:", d.Name.Name)
		}
	}
}

// -----------------------------------------------------------------------------

var fsTestStd2 = asttest.NewSingleFileFS("/foo", "bar.gop", `package bar; import "io"
	x := []float64{1, 3.4, 5}
	y := map[string]float64{"Hello": 1, "xsw": 3.4}
	println("x:", x, "y:", y)

	a := [...]float64{1, 3.4, 5}
	b := [...]float64{1, 3: 3.4, 5}
	c := []float64{2: 1.2, 3, 6: 4.5}
	println("a:", a, "b:", b, "c:", c)
`)

func TestStd2(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := ParseFSDir(fset, fsTestStd2, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}
	bar := pkgs["bar"]
	file := bar.Files["/foo/bar.gop"]
	fmt.Println("Pkg:", file.Name)
	for _, decl := range file.Decls {
		fmt.Println("decl:", reflect.TypeOf(decl))
		switch d := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch v := spec.(type) {
				case *ast.ImportSpec:
					fmt.Println(" - import:", v.Path.Value)
				}
			}
		case *ast.FuncDecl:
			fmt.Println(" - func:", d.Name.Name)
		}
	}
}

// -----------------------------------------------------------------------------

var fsTestStdFor = asttest.NewSingleFileFS("/foo", "bar.gop", `
	n := 0
	for range [1, 3, 5, 7, 11] {
		n++
	}
	println("n:", n)

	sum := 0
	for _, x := range [1, 3, 5, 7, 11] {
		if x > 3 {
			sum += x
		}
	}
	println("sum(1,3,5,7,11):", sum)

	sum = 0
	for i := 1; i < 100; i++ {
		sum += i
	}
	println("sum(1-100):", sum)
`)

func TestStdFor(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := ParseFSDir(fset, fsTestStdFor, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}
	bar := pkgs["main"]
	file := bar.Files["/foo/bar.gop"]
	fmt.Println("Pkg:", file.Name)
	for _, decl := range file.Decls {
		fmt.Println("decl:", reflect.TypeOf(decl))
		switch d := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch v := spec.(type) {
				case *ast.ImportSpec:
					fmt.Println(" - import:", v.Path.Value)
				}
			}
		case *ast.FuncDecl:
			fmt.Println(" - func:", d.Name.Name)
		}
	}
}

// -----------------------------------------------------------------------------

var fsTestBuild = asttest.NewSingleFileFS("/foo", "bar.gop", `
	type cstring string

	title := "Hello,world!2020-05-27"
	s := (*cstring)(&title)
	println(title[0:len(title)-len("2006-01-02")])
`)

func TestBuild(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := ParseFSDir(fset, fsTestBuild, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}
	bar := pkgs["main"]
	file := bar.Files["/foo/bar.gop"]
	fmt.Println("Pkg:", file.Name)
	for _, decl := range file.Decls {
		fmt.Println("decl:", reflect.TypeOf(decl))
		switch d := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch v := spec.(type) {
				case *ast.ImportSpec:
					fmt.Println(" - import:", v.Path.Value)
				}
			}
		case *ast.FuncDecl:
			fmt.Println(" - func:", d.Name.Name)
		}
	}
}

// -----------------------------------------------------------------------------

func testFrom(t *testing.T, pkgDir, sel string) {
	if sel != "" && !strings.Contains(pkgDir, sel) {
		return
	}
	log.Debug("Parsing", pkgDir)
	fset := token.NewFileSet()
	pkgs, err := ParseDir(fset, pkgDir, nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseDir failed:", err, len(pkgs))
	}
}

func TestFromTestdata(t *testing.T) {
	sel := ""
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Getwd failed:", err)
	}
	dir = path.Join(dir, "../exec/golang/testdata")
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal("ReadDir failed:", err)
	}
	for _, fi := range fis {
		testFrom(t, dir+"/"+fi.Name(), sel)
	}
}

func TestFromTestdata2(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Getwd failed:", err)
	}
	dir = path.Join(dir, "./testdata")
	testFrom(t, dir, "")
}

// -----------------------------------------------------------------------------

var fsTestArray = asttest.NewSingleFileFS("/foo", "bar.gop", `
println([1][0])
println([1,2][0])
println([][]int{})
println([2][2]int{})
`)

func TestArray(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := ParseFSDir(fset, fsTestArray, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}
	bar := pkgs["main"]
	file := bar.Files["/foo/bar.gop"]
	fmt.Println("Pkg:", file.Name)
	for _, decl := range file.Decls {
		fmt.Println("decl:", reflect.TypeOf(decl))
		switch d := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch v := spec.(type) {
				case *ast.ImportSpec:
					fmt.Println(" - import:", v.Path.Value)
				}
			}
		case *ast.FuncDecl:
			fmt.Println(" - func:", d.Name.Name)
		}
	}
}
