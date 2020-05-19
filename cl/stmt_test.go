package cl

import (
	"fmt"
	"os"
	"testing"

	"github.com/qiniu/qlang/v6/ast/asttest"
	"github.com/qiniu/qlang/v6/exec"
	"github.com/qiniu/qlang/v6/parser"
	"github.com/qiniu/qlang/v6/token"
)

// -----------------------------------------------------------------------------

var fsTestSwif = asttest.NewSingleFileFS("/foo", "bar.ql", `
	x := 0
	t := "Hello"
	switch {
	case t == "xsw":
		x = 3
	case t == "Hello", t == "world":
		x = 5
	default:
		x = 7
	}
	x
`)

func TestSwitchIf(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestSwif, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, err = NewPackage(b, bar)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	fmt.Println("x:", ctx.Get(-1))
	if v := ctx.Get(-1); v != 5 {
		t.Fatal("x:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestSwif2 = asttest.NewSingleFileFS("/foo", "bar.ql", `
	x := 0
	t := "Hello"
	switch {
	case t == "xsw":
		x = 3
	case t == "hi", t == "world":
		x = 11
	default:
		x = 7
	}
	x
`)

func TestSwitchIfDefault(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestSwif2, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, err = NewPackage(b, bar)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	code.Dump(os.Stdout)

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	fmt.Println("x:", ctx.Get(-1))
	if v := ctx.Get(-1); v != 7 {
		t.Fatal("x:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestSw = asttest.NewSingleFileFS("/foo", "bar.ql", `
	x := 0
	switch t := "Hello"; t {
	case "xsw":
		x = 3
	case "Hello", "world":
		x = 5
	default:
		x= 7
	}
	x
`)

func TestSwitch(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestSw, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, err = NewPackage(b, bar)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	fmt.Println("x:", ctx.Get(-1))
	if v := ctx.Get(-1); v != 5 {
		t.Fatal("x:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestSw2 = asttest.NewSingleFileFS("/foo", "bar.ql", `
	x := 0
	t := "Hello"
	switch t {
	}
	switch t {
	case "world", "Hello":
		x = 5
	case "xsw":
		x = 3
	}
	x
`)

func TestSwitch2(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestSw2, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, err = NewPackage(b, bar)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	fmt.Println("x:", ctx.Get(-1))
	if v := ctx.Get(-1); v != 5 {
		t.Fatal("x:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestSw3 = asttest.NewSingleFileFS("/foo", "bar.ql", `
	x := 0
	t := "Hello"
	switch t {
	default:
		x = 7
	case "world", "hi":
		x = 5
	case "xsw":
		x = 3
	}
	x
`)

func TestDefault(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestSw3, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, err = NewPackage(b, bar)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	fmt.Println("x:", ctx.Get(-1))
	if v := ctx.Get(-1); v != 7 {
		t.Fatal("x:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestIf = asttest.NewSingleFileFS("/foo", "bar.ql", `
	x := 0
	if t := false; t {
		x = 3
	} else {
		x = 5
	}
	x
`)

func TestIf(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestIf, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, err = NewPackage(b, bar)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	fmt.Println("x:", ctx.Get(-1))
	if v := ctx.Get(-1); v != 5 {
		t.Fatal("x:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestIf2 = asttest.NewSingleFileFS("/foo", "bar.ql", `
	x := 5
	if true {
		x = 3
	}
	x
`)

func TestIf2(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestIf2, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, err = NewPackage(b, bar)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	fmt.Println("x:", ctx.Get(-1))
	if v := ctx.Get(-1); v != 3 {
		t.Fatal("x:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestReturn = asttest.NewSingleFileFS("/foo", "bar.ql", `
	import (
		"fmt"
		"strings"
	)

	func foo(x string) string {
		return strings.NewReplacer("?", "!").Replace(x)
	}

	fmt.Println(foo("Hello, world???"))
`)

func TestReturn(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestReturn, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, err = NewPackage(b, bar)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	fmt.Println("results:", ctx.Get(-2), ctx.Get(-1))
	if v := ctx.Get(-1); v != nil {
		t.Fatal("error:", v)
	}
	if v := ctx.Get(-2); v != int(16) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestFunc = asttest.NewSingleFileFS("/foo", "bar.ql", `
	import "fmt"

	func foo(x string) (n int, err error) {
		n, err = fmt.Println("x: " + x)
		return
	}

	foo("Hello, world!")
`)

func TestFunc(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestFunc, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, err = NewPackage(b, bar)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	fmt.Println("results:", ctx.Get(-2), ctx.Get(-1))
	if v := ctx.Get(-1); v != nil {
		t.Fatal("error:", v)
	}
	if v := ctx.Get(-2); v != int(17) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------

var fsTestClosure = asttest.NewSingleFileFS("/foo", "bar.ql", `
	import "fmt"

	foo := func(prompt string) (n int, err error) {
		n, err = fmt.Println(prompt + x)
		return
	}

	x := "Hello, world!"
	foo("x: ")
`)

func TestClosure(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestClosure, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, err = NewPackage(b, bar)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	fmt.Println("results:", ctx.Get(-2), ctx.Get(-1))
	if v := ctx.Get(-1); v != nil {
		t.Fatal("error:", v)
	}
	if v := ctx.Get(-2); v != int(17) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------
