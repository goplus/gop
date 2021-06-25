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
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/ast/asttest"
	"github.com/goplus/gop/token"
	"github.com/qiniu/x/log"
)

func init() {
	log.SetFlags(log.Ldefault &^ log.LstdFlags)
	log.SetOutputLevel(log.Ldebug)
}

// -----------------------------------------------------------------------------

var fsTest1 = asttest.NewSingleFileFS("/foo", "bar.gop", `package bar; import "io"
func New() (*Bar, error) {
	return nil, io.EOF
}

bar, err := New()
if err != nil {
	log.Println(err)
}`)

func TestParseBarPackage(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := ParseFSDir(fset, fsTest1, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}
	bar := pkgs["bar"]
	file := bar.Files["/foo/bar.gop"]
	fmt.Println("Pkg:", file.Name)
	if !checkFSTest(file, []string{`"io"`}, []string{"init", "New"}) {
		t.Error("file check file")
	}
}

// -----------------------------------------------------------------------------

var fsTest2 = asttest.NewSingleFileFS("/foo", "bar.gop", `import "io"
func New() (*Bar, error) {
	return nil, io.EOF
}

bar, err := New()
if err != nil {
	log.Println(err)
}
`)

func TestParseNoPackageAndGlobalCode(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := ParseFSDir(fset, fsTest2, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}
	bar, isMain := pkgs["main"]
	if !isMain {
		t.Fatal("TestParseNoPackageAndGlobalCode failed: not main")
	}
	file := bar.Files["/foo/bar.gop"]
	fmt.Println("Pkg:", file.Name)
	if !checkFSTest(file, []string{`"io"`}, []string{"New", "main"}) {
		t.Error("file check file")
	}
}

// -----------------------------------------------------------------------------

const script3 = `package bar; import "io"
func New() (*Bar, error) {
	return nil, io.EOF
}

bar, err := New()
if err != nil {
	log.Println(err)
}`

func TestParse(t *testing.T) {
	fset := token.NewFileSet()
	local = asttest.NewSingleFileFS("/foo", "bar.gop", script3)
	defer func() {
		local = localFS{}
	}()

	// test parse file
	pkgs, err := Parse(fset, "/foo/bar.gop", nil, 0)
	if err != nil {
		t.Error("parse err!", err)
		return
	}
	pkg, ok := pkgs["bar"]
	if ok != true {
		t.Error("pkg not found")
		return
	}
	file, ok := pkg.Files["/foo/bar.gop"]
	if ok != true {
		t.Error("file not found")
		return
	}
	if !checkFSTest(file, []string{`"io"`}, []string{"init", "New"}) {
		t.Error("file check file")
	}
}

func doTestParse(t *testing.T, src interface{}) {
	fset := token.NewFileSet()
	pkgs, err := Parse(fset, "/foo/bar.gop", src, 0)
	if err != nil {
		t.Error("parse err!", err)
		return
	}
	pkg, ok := pkgs["bar"]
	if ok != true {
		t.Error("pkg not found")
		return
	}
	file, ok := pkg.Files["/foo/bar.gop"]
	if ok != true {
		t.Error("file not found")
		return
	}
	if !checkFSTest(file, []string{`"io"`}, []string{"init", "New"}) {
		t.Error("file check file")
	}
}

func TestParse2(t *testing.T) {
	doTestParse(t, script3)
	doTestParse(t, bytes.NewBufferString(script3))
	doTestParse(t, strings.NewReader(script3))
}

func checkFSTest(file *ast.File, targetImport []string, targetFuncDecl []string) (pass bool) {
	var foundImport, foundFunc int
	for _, decl := range file.Decls {
		fmt.Println("decl:", reflect.TypeOf(decl))
		switch d := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch v := spec.(type) {
				case *ast.ImportSpec:
					for _, ti := range targetImport {
						if ti == v.Path.Value {
							foundImport++
						}
					}
				}
			}
		case *ast.FuncDecl:
			fmt.Println(d.Name.Name)
			for _, tfd := range targetFuncDecl {
				if tfd == d.Name.Name {
					foundFunc++
				}
			}
		}
	}
	if foundFunc != len(targetFuncDecl) || foundImport != len(targetImport) {
		fmt.Println("invalid check", foundImport, foundFunc)
		return false
	}
	return true
}

// -----------------------------------------------------------------------------

var fsTestMapLit = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := {"Hello": 1, "xsw": 3.4}
	println("x:", x)
`)

func TestMapLit(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := ParseFSDir(fset, fsTestMapLit, "/foo", nil, Trace)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}
	bar, isMain := pkgs["main"]
	if !isMain {
		t.Fatal("TestMap failed: not main")
	}
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

var fsTestSliceLit = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := [1, 3.4, 5.6]
	println("x:", x)
`)

func TestSliceLit(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := ParseFSDir(fset, fsTestSliceLit, "/foo", nil, Trace)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}
	bar, isMain := pkgs["main"]
	if !isMain {
		t.Fatal("TestMap failed: not main")
	}
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

var fsTestSliceLit2 = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := [5.6]
	println("x:", x)
`)

func TestSliceLit2(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := ParseFSDir(fset, fsTestSliceLit2, "/foo", nil, Trace)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}
	bar, isMain := pkgs["main"]
	if !isMain {
		t.Fatal("TestMap failed: not main")
	}
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

var fsTestSliceLit3 = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := []
	println("x:", x)
`)

func TestSliceLit3(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := ParseFSDir(fset, fsTestSliceLit3, "/foo", nil, Trace)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}
	bar, isMain := pkgs["main"]
	if !isMain {
		t.Fatal("TestMap failed: not main")
	}
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

var fsTestListComprehension = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := [x*x for x <- [1, 2, 3, 4], x > 2]
	println("x:", x)
`)

func TestListComprehension(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := ParseFSDir(fset, fsTestListComprehension, "/foo", nil, Trace)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}
	bar, isMain := pkgs["main"]
	if !isMain {
		t.Fatal("TestMap failed: not main")
	}
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

var fsTestMapComprehension = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := {v: k*k for k, v <- [3, 5, 7, 11]}
	println("x:", x)
`)

func TestMapComprehension(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := ParseFSDir(fset, fsTestMapComprehension, "/foo", nil, Trace)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}
	bar, isMain := pkgs["main"]
	if !isMain {
		t.Fatal("TestMap failed: not main")
	}
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

var fsTestMapComprehensionCond = asttest.NewSingleFileFS("/foo", "bar.gop", `
	x := {v: k*k for k, v <- [3, 5, 7, 11], k % 2 == 1}
	println("x:", x)

	x = map[int]int{}
	for k, v := range []int{3, 5, 7, 11} {
		if k % 2 == 1 {
			x[v] = k*k
		}
	}
	println("x:", x)
`)

func TestMapComprehensionCond(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := ParseFSDir(fset, fsTestMapComprehensionCond, "/foo", nil, Trace)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}
	bar, isMain := pkgs["main"]
	if !isMain {
		t.Fatal("TestMap failed: not main")
	}
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

var fsTestListComprehensionEx = asttest.NewSingleFileFS("/foo", "bar.gop", `
	arr := [1, 2, 3, 4, 5, 6]
	x := [[a, b] for a <- arr, a < b for b <- arr, b > 2]
	println("x:", x)
`)

func TestListComprehensionEx(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := ParseFSDir(fset, fsTestListComprehensionEx, "/foo", nil, Trace)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}
	bar, isMain := pkgs["main"]
	if !isMain {
		t.Fatal("TestMap failed: not main")
	}
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

var fsTestFor = asttest.NewSingleFileFS("/foo", "bar.gop", `
	sum := 0
	for x <- [1, 3, 5, 7, 11], x > 3 {
		sum += x
	}
	println("sum(5,7,11):", sum)

	sum = 0
	for i, x <- [1, 3, 5, 7, 11] {
		sum += x
	}
	println("sum(1,3,5,7,11):", sum)
`)

func TestFor(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := ParseFSDir(fset, fsTestFor, "/foo", nil, Trace)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}
	bar, isMain := pkgs["main"]
	if !isMain {
		t.Fatal("TestMap failed: not main")
	}
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

var fsTestMake = asttest.NewSingleFileFS("/foo", "bar.gop", `
	a := make([]int, 0, 4)
	a = append(a, [1, 2, 3]...)

	b := make(map[string]interface{ Error() string })
	b := make(chan bool)
`)

func TestMake(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := ParseFSDir(fset, fsTestMake, "/foo", nil, Trace)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}
	bar, isMain := pkgs["main"]
	if !isMain {
		t.Fatal("TestMap failed: not main")
	}
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

var fsTestXOR = asttest.NewSingleFileFS("/foo", "bar.gop", `
	println(^uint32(1))
`)

func TestXOR(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := ParseFSDir(fset, fsTestXOR, "/foo", nil, Trace)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}
	bar, isMain := pkgs["main"]
	if !isMain {
		t.Fatal("TestTILDE failed: not main")
	}
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

var fsTestRational = asttest.NewSingleFileFS("/foo", "bar.gop", `
	println(5/7r + 3.4r)
`)

func TestRational(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := ParseFSDir(fset, fsTestRational, "/foo", nil, Trace)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}
	bar, isMain := pkgs["main"]
	if !isMain {
		t.Fatal("TestTILDE failed: not main")
	}
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

var valids = []string{
	"package p\n",
	`package p;`,
	`package p; import "fmt"; func f() { fmt.Println("Hello, World!") };`,
	`package p; func f() { if f(T{}) {} };`,
	`package p; func f() { _ = <-chan int(nil) };`,
	`package p; func f() { _ = (<-chan int)(nil) };`,
	`package p; func f() { _ = (<-chan <-chan int)(nil) };`,
	`package p; func f() { _ = <-chan <-chan <-chan <-chan <-int(nil) };`,
	`package p; func f(func() func() func());`,
	`package p; func f(...T);`,
	`package p; func f(float, ...int);`,
	`package p; func f(x int, a ...int) { f(0, a...); f(1, a...,) };`,
	`package p; func f(int,) {};`,
	`package p; func f(...int,) {};`,
	`package p; func f(x ...int,) {};`,
	`package p; type T []int; var a []bool; func f() { if a[T{42}[0]] {} };`,
	`package p; type T []int; func g(int) bool { return true }; func f() { if g(T{42}[0]) {} };`,
	`package p; type T []int; func f() { for _ = range []int{T{42}[0]} {} };`,
	`package p; var a = T{{1, 2}, {3, 4}}`,
	`package p; func f() { select { case <- c: case c <- d: case c <- <- d: case <-c <- d: } };`,
	`package p; func f() { select { case x := (<-c): } };`,
	`package p; func f() { if ; true {} };`,
	`package p; func f() { switch ; {} };`,
	`package p; func f() { for _ = range "foo" + "bar" {} };`,
	`package p; func f() { var s []int; g(s[:], s[i:], s[:j], s[i:j], s[i:j:k], s[:j:k]) };`,
	`package p; var ( _ = (struct {*T}).m; _ = (interface {T}).m )`,
	`package p; func ((T),) m() {}`,
	`package p; func ((*T),) m() {}`,
	`package p; func (*(T),) m() {}`,
	`package p; func _(x []int) { for range x {} }`,
	`package p; func _() { if [T{}.n]int{} {} }`,
	`package p; func _() { map[int]int{}[0]++; map[int]int{}[0] += 1 }`,
	`package p; func _(x interface{f()}) { interface{f()}(x).f() }`,
	`package p; func _(x chan int) { chan int(x) <- 0 }`,
	`package p; const (x = 0; y; z)`, // issue 9639
	`package p; var _ = map[P]int{P{}:0, {}:1}`,
	`package p; var _ = map[*P]int{&P{}:0, {}:1}`,
	`package p; type T = int`,
	`package p; type (T = p.T; _ = struct{}; x = *T)`,
}

func TestParseExpr(t *testing.T) {
	// just kicking the tires:
	// a valid arithmetic expression
	src := "a + b"
	x, err := ParseExpr(src)
	if err != nil {
		t.Errorf("ParseExpr(%q): %v", src, err)
	}
	// sanity check
	if _, ok := x.(*ast.BinaryExpr); !ok {
		t.Errorf("ParseExpr(%q): got %T, want *ast.BinaryExpr", src, x)
	}

	// a valid type expression
	src = "struct{x *int}"
	x, err = ParseExpr(src)
	if err != nil {
		t.Errorf("ParseExpr(%q): %v", src, err)
	}
	// sanity check
	if _, ok := x.(*ast.StructType); !ok {
		t.Errorf("ParseExpr(%q): got %T, want *ast.StructType", src, x)
	}

	// an invalid expression
	src = "a + *"
	x, err = ParseExpr(src)
	if err == nil {
		t.Errorf("ParseExpr(%q): got no error", src)
	}
	if x == nil {
		t.Errorf("ParseExpr(%q): got no (partial) result", src)
	}
	if _, ok := x.(*ast.BinaryExpr); !ok {
		t.Errorf("ParseExpr(%q): got %T, want *ast.BinaryExpr", src, x)
	}

	// a valid expression followed by extra tokens is invalid
	src = "a[i] := x"
	if _, err := ParseExpr(src); err == nil {
		t.Errorf("ParseExpr(%q): got no error", src)
	}

	// a semicolon is not permitted unless automatically inserted
	src = "a + b\n"
	if _, err := ParseExpr(src); err != nil {
		t.Errorf("ParseExpr(%q): got error %s", src, err)
	}
	src = "a + b;"
	if _, err := ParseExpr(src); err == nil {
		t.Errorf("ParseExpr(%q): got no error", src)
	}

	// various other stuff following a valid expression
	const validExpr = "a + b"
	const anything = "dh3*#D)#_"
	for _, c := range "!)]};," {
		src := validExpr + string(c) + anything
		if _, err := ParseExpr(src); err == nil {
			t.Errorf("ParseExpr(%q): got no error", src)
		}
	}

	// ParseExpr must not crash
	for _, src := range valids {
		ParseExpr(src)
	}
}

// -----------------------------------------------------------------------------
