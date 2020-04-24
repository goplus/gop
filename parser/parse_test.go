package parser

import (
	"fmt"
	"go/ast"
	"os"
	"path"
	"reflect"
	"syscall"
	"testing"
	"time"

	"github.com/qiniu/qlang/token"
)

// -----------------------------------------------------------------------------

type memFileInfo struct {
	name string
}

func (p *memFileInfo) Name() string {
	return p.name
}

func (p *memFileInfo) Size() int64 {
	return 0
}

func (p *memFileInfo) Mode() os.FileMode {
	return 0
}

func (p *memFileInfo) ModTime() (t time.Time) {
	return
}

func (p *memFileInfo) IsDir() bool {
	return false
}

func (p *memFileInfo) Sys() interface{} {
	return nil
}

type memFS struct {
	dirs  map[string][]string
	files map[string]string
}

func (p *memFS) ReadDir(dirname string) ([]os.FileInfo, error) {
	if items, ok := p.dirs[dirname]; ok {
		fis := make([]os.FileInfo, len(items))
		for i, item := range items {
			fis[i] = &memFileInfo{name: item}
		}
		return fis, nil
	}
	return nil, syscall.ENOENT
}

func (p *memFS) ReadFile(filename string) ([]byte, error) {
	if data, ok := p.files[filename]; ok {
		return []byte(data), nil
	}
	return nil, syscall.ENOENT
}

func (p *memFS) Join(elem ...string) string {
	return path.Join(elem...)
}

// -----------------------------------------------------------------------------

var fsTest1 = &memFS{
	dirs: map[string][]string{
		"/foo": {
			"bar.ql",
		},
	},
	files: map[string]string{
		"/foo/bar.ql": `package bar; import "io"

func New() (*Bar, error) {
	return nil, io.EOF
}

bar, err := New()
if err != nil {
	log.Println(err)
}`},
}

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

var fsTest2 = &memFS{
	dirs: map[string][]string{
		"/foo": {
			"bar.ql",
		},
	},
	files: map[string]string{
		"/foo/bar.ql": `import "io"

func New() (*Bar, error) {
	return nil, io.EOF
}

bar, err := New()
if err != nil {
	log.Println(err)
}`},
}

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
