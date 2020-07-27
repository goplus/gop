package cltest

import (
	"fmt"
	"os"
	"testing"

	"github.com/goplus/gop/ast/asttest"
	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
	"github.com/qiniu/x/ts"

	exec "github.com/goplus/gop/exec/bytecode"
	_ "github.com/goplus/gop/lib" // libraries
)

// -----------------------------------------------------------------------------

// Expect runs a script and check expected output and if panic or not.
func Expect(t *testing.T, script string, expected string, panicMsg ...interface{}) {
	fset := token.NewFileSet()
	fs := asttest.NewSingleFileFS("/foo", "bar.gop", script)
	pkgs, err := parser.ParseFSDir(fset, fs, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	e := ts.StartExpecting(t, ts.CapStdout)
	defer e.Close()

	e.Call(func() {
		bar := pkgs["main"]
		b := exec.NewBuilder(nil)
		pkg, err := cl.NewPackage(b.Interface(), bar, fset, cl.PkgActClMain)
		if err != nil {
			t.Fatal("Compile failed:", err)
		}
		cl.Debug(pkg)
		code := b.Resolve()
		code.Dump(os.Stderr)
		fmt.Fprintln(os.Stderr)
		exec.NewContext(code).Run()
	}).Panic(panicMsg...).Expect(expected)
}

// -----------------------------------------------------------------------------

// Call runs a script and gets the last expression value to check
func Call(t *testing.T, script string, idx ...int) *ts.TestCase {
	fset := token.NewFileSet()
	fs := asttest.NewSingleFileFS("/foo", "bar.gop", script)
	pkgs, err := parser.ParseFSDir(fset, fs, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	return ts.New(t).Call(func() interface{} {
		bar := pkgs["main"]
		b := exec.NewBuilder(nil)
		pkg, err := cl.NewPackage(b.Interface(), bar, fset, cl.PkgActClMain)
		if err != nil {
			t.Fatal("Compile failed:", err)
		}
		cl.Debug(pkg)
		code := b.Resolve()
		ctx := exec.NewContext(code)
		ctx.Run()
		return ctx.Get(append(idx, -1)[0])
	})
}

// -----------------------------------------------------------------------------
