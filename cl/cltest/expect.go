package cltest

import (
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

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	pkg, err := cl.NewPackage(b.Interface(), bar, fset, cl.PkgActClMain)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	cl.Debug(pkg)
	code := b.Resolve()

	ctx := exec.NewContext(code)
	e := ts.StartExpecting(t, ts.CapStdout)
	defer e.Close()
	e.Call(ctx.Exec, 0, code.Len()).Expect(expected).Panic(panicMsg...)
}

// -----------------------------------------------------------------------------

// Call runs a script and gets the last expression value to check
func Call(t *testing.T, script string) *ts.TestCase {

	fset := token.NewFileSet()
	fs := asttest.NewSingleFileFS("/foo", "bar.gop", script)
	pkgs, err := parser.ParseFSDir(fset, fs, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	pkg, err := cl.NewPackage(b.Interface(), bar, fset, cl.PkgActClMain)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
	cl.Debug(pkg)
	code := b.Resolve()

	ctx := exec.NewContext(code)
	return ts.New(t).Call(func() interface{} {
		ctx.Exec(0, code.Len())
		return ctx.Get(-1)
	})
}

// -----------------------------------------------------------------------------
