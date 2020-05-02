package cl

import (
	"testing"

	"github.com/qiniu/qlang/ast/asttest"
	"github.com/qiniu/qlang/exec"
	"github.com/qiniu/qlang/parser"
	"github.com/qiniu/qlang/token"
	"github.com/qiniu/x/log"
)

func init() {
	log.SetFlags(log.Ldefault &^ log.LstdFlags)
	log.SetOutputLevel(log.Ldebug)
}

var fsTestBasic = asttest.NewSingleFileFS("/foo", "bar.ql", `
	println("Hello, world!")
`)

func TestBasic(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestBasic, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, err = NewPackage(b, bar)
	if err != nil {
		t.Fatal("Compile failed:", err)
	}
}
