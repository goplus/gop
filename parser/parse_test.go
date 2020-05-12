package parser

import (
	"fmt"
	"go/ast"
	"reflect"
	"testing"

	"github.com/qiniu/qlang/ast/asttest"
	"github.com/qiniu/qlang/token"
)

// -----------------------------------------------------------------------------

var fsTest1 = asttest.NewSingleFileFS("/foo", "bar.ql", `package bar; import "io"

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
	file := bar.Files["/foo/bar.ql"]
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

var fsTest2 = asttest.NewSingleFileFS("/foo", "bar.ql", `import "io"

func New() (*Bar, error) {
	return nil, io.EOF
}

bar, err := New()
if err != nil {
	log.Println(err)
}`)

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
	file := bar.Files["/foo/bar.ql"]
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
