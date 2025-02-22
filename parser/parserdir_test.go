/*
 * Copyright (c) 2021 The GoPlus Authors (goplus.org). All rights reserved.
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

package parser

import (
	"bytes"
	"os"
	"path"
	"reflect"
	"strings"
	"syscall"
	"testing"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/parser/fsx"
	"github.com/goplus/gop/parser/fsx/memfs"
	"github.com/goplus/gop/parser/parsertest"
	"github.com/goplus/gop/scanner"
	"github.com/goplus/gop/token"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

func init() {
	log.SetFlags(log.Llongfile)
	SetDebug(DbgFlagAll)
}

var testStdCode = `package bar; import "io"
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
`

func TestStd(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := Parse(fset, "/foo/bar.gop", testStdCode, ParseComments)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("Parse failed:", err, len(pkgs))
	}
	bar := pkgs["bar"]
	parsertest.Expect(t, bar, `package bar

file bar.gop
noEntrypoint
ast.GenDecl:
  Tok: import
  Specs:
    ast.ImportSpec:
      Path:
        ast.BasicLit:
          Kind: STRING
          Value: "io"
ast.FuncDecl:
  Doc:
    ast.CommentGroup:
      List:
        ast.Comment:
          Text: // comment
  Name:
    ast.Ident:
      Name: main
  Type:
    ast.FuncType:
      Params:
        ast.FieldList:
  Body:
    ast.BlockStmt:
      List:
        ast.AssignStmt:
          Lhs:
            ast.Ident:
              Name: x
          Tok: :=
          Rhs:
            ast.BasicLit:
              Kind: INT
              Value: 0
        ast.IfStmt:
          Init:
            ast.AssignStmt:
              Lhs:
                ast.Ident:
                  Name: t
              Tok: :=
              Rhs:
                ast.Ident:
                  Name: false
          Cond:
            ast.Ident:
              Name: t
          Body:
            ast.BlockStmt:
              List:
                ast.AssignStmt:
                  Lhs:
                    ast.Ident:
                      Name: x
                  Tok: =
                  Rhs:
                    ast.BasicLit:
                      Kind: INT
                      Value: 3
          Else:
            ast.BlockStmt:
              List:
                ast.AssignStmt:
                  Lhs:
                    ast.Ident:
                      Name: x
                  Tok: =
                  Rhs:
                    ast.BasicLit:
                      Kind: INT
                      Value: 5
        ast.ExprStmt:
          X:
            ast.CallExpr:
              Fun:
                ast.Ident:
                  Name: println
              Args:
                ast.BasicLit:
                  Kind: STRING
                  Value: "x:"
                ast.Ident:
                  Name: x
        ast.AssignStmt:
          Lhs:
            ast.Ident:
              Name: x
          Tok: =
          Rhs:
            ast.BasicLit:
              Kind: INT
              Value: 0
        ast.SwitchStmt:
          Init:
            ast.AssignStmt:
              Lhs:
                ast.Ident:
                  Name: s
              Tok: :=
              Rhs:
                ast.BasicLit:
                  Kind: STRING
                  Value: "Hello"
          Tag:
            ast.Ident:
              Name: s
          Body:
            ast.BlockStmt:
              List:
                ast.CaseClause:
                  Body:
                    ast.AssignStmt:
                      Lhs:
                        ast.Ident:
                          Name: x
                      Tok: =
                      Rhs:
                        ast.BasicLit:
                          Kind: INT
                          Value: 7
                ast.CaseClause:
                  List:
                    ast.BasicLit:
                      Kind: STRING
                      Value: "world"
                    ast.BasicLit:
                      Kind: STRING
                      Value: "hi"
                  Body:
                    ast.AssignStmt:
                      Lhs:
                        ast.Ident:
                          Name: x
                      Tok: =
                      Rhs:
                        ast.BasicLit:
                          Kind: INT
                          Value: 5
                ast.CaseClause:
                  List:
                    ast.BasicLit:
                      Kind: STRING
                      Value: "xsw"
                  Body:
                    ast.AssignStmt:
                      Lhs:
                        ast.Ident:
                          Name: x
                      Tok: =
                      Rhs:
                        ast.BasicLit:
                          Kind: INT
                          Value: 3
        ast.ExprStmt:
          X:
            ast.CallExpr:
              Fun:
                ast.Ident:
                  Name: println
              Args:
                ast.BasicLit:
                  Kind: STRING
                  Value: "x:"
                ast.Ident:
                  Name: x
        ast.AssignStmt:
          Lhs:
            ast.Ident:
              Name: c
          Tok: :=
          Rhs:
            ast.CallExpr:
              Fun:
                ast.Ident:
                  Name: make
              Args:
                ast.ChanType:
                  Value:
                    ast.Ident:
                      Name: bool
                ast.BasicLit:
                  Kind: INT
                  Value: 100
        ast.SelectStmt:
          Body:
            ast.BlockStmt:
              List:
                ast.CommClause:
                  Comm:
                    ast.SendStmt:
                      Chan:
                        ast.Ident:
                          Name: c
                      Values:
                        ast.Ident:
                          Name: true
                ast.CommClause:
                  Comm:
                    ast.AssignStmt:
                      Lhs:
                        ast.Ident:
                          Name: v
                      Tok: :=
                      Rhs:
                        ast.UnaryExpr:
                          Op: <-
                          X:
                            ast.Ident:
                              Name: c
                ast.CommClause:
                  Body:
                    ast.ExprStmt:
                      X:
                        ast.CallExpr:
                          Fun:
                            ast.Ident:
                              Name: panic
                          Args:
                            ast.BasicLit:
                              Kind: STRING
                              Value: "error"
`)
}

func TestParseExprFrom(t *testing.T) {
	fset := token.NewFileSet()
	if _, err := ParseExprFrom(fset, "/foo/bar/not-exists", nil, 0); err == nil {
		t.Fatal("ParseExprFrom: no error?")
	}
	if _, err := ParseExpr("1+1\n;"); err == nil || err.Error() != "2:1: expected 'EOF', found ';'" {
		t.Fatal("ParseExpr:", err)
	}
}

func TestReadSource(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	if _, err := readSource(buf); err != nil {
		t.Fatal("readSource failed:", err)
	}
	sr := strings.NewReader("")
	if _, err := readSource(sr); err != nil {
		t.Fatal("readSource strings.Reader failed:", err)
	}
	if _, err := readSource(0); err == nil {
		t.Fatal("readSource int failed: no error?")
	}
	if _, err := readSourceLocal("/foo/bar/not-exists", nil); err == nil {
		t.Fatal("readSourceLocal int failed: no error?")
	}
}

func TestParseFiles(t *testing.T) {
	fset := token.NewFileSet()
	if _, err := ParseFSFiles(fset, fsx.Local, []string{"/foo/bar/not-exists"}, PackageClauseOnly|SaveAbsFile); err == nil {
		t.Fatal("ParseFiles failed: no error?")
	}
}

func TestIparseFileInvalidSrc(t *testing.T) {
	fset := token.NewFileSet()
	if _, err := parseFile(fset, "/foo/bar/not-exists", 1, PackageClauseOnly); err != errInvalidSource {
		t.Fatal("ParseFile failed: not errInvalidSource?")
	}
}

func TestIparseFileNoFset(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Fatal("ParseFile failed: no error?")
		}
	}()
	parseFile(nil, "/foo/bar/not-exists", nil, PackageClauseOnly)
}

func TestParseDir(t *testing.T) {
	fset := token.NewFileSet()
	if _, err := ParseDir(fset, "/foo/bar/not-exists", nil, PackageClauseOnly); err == nil {
		t.Fatal("ParseDir failed: no error?")
	}
}

func testFrom(t *testing.T, pkgDir, sel string, exclude Mode) {
	if sel != "" && !strings.Contains(pkgDir, sel) {
		return
	}
	t.Helper()
	log.Println("Parsing", pkgDir)
	fset := token.NewFileSet()
	pkgs, err := ParseDir(fset, pkgDir, nil, (Trace|ParseComments|ParseGoAsGoPlus)&^exclude)
	if err != nil || len(pkgs) != 1 {
		if errs, ok := err.(scanner.ErrorList); ok {
			for _, e := range errs {
				t.Log(e)
			}
		}
		t.Fatal("ParseDir failed:", err, reflect.TypeOf(err), len(pkgs))
	}
	for _, pkg := range pkgs {
		b, err := os.ReadFile(pkgDir + "/parser.expect")
		if err != nil {
			t.Fatal("Parsing", pkgDir, "-", err)
		}
		parsertest.Expect(t, pkg, string(b))
		return
	}
}

func TestParseGo(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := ParseDirEx(fset, "./_testdata/functype", Config{Mode: Trace | ParseGoAsGoPlus})
	if err != nil {
		t.Fatal("TestParseGo: ParseDir failed -", err)
	}
	if len(pkgs) != 1 {
		t.Fatal("TestParseGo failed: len(pkgs) =", len(pkgs))
	}
}

func TestParseGoFiles(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := ParseFiles(fset, []string{"./_testdata/functype/functype.go"}, Trace|ParseGoAsGoPlus)
	if err != nil {
		t.Fatal("TestParseGoFiles: ParseFiles failed -", err)
	}
	if len(pkgs) != 1 {
		t.Fatal("TestParseGoFiles failed: len(pkgs) =", len(pkgs))
	}
}

func TestParseEntries(t *testing.T) {
	doTestParseEntries(t, Config{})
	_, err := ParseEntries(nil, []string{"/not-found/gopfile.gox"}, Config{})
	if err == nil {
		t.Fatal("ParseEntries: no error?")
	}
}

func TestParseEntries_SaveAbsFile(t *testing.T) {
	doTestParseEntries(t, Config{Mode: SaveAbsFile})
}

func doTestParseEntries(t *testing.T, confReal Config) {
	doTestParseEntry(t, func(fset *token.FileSet, filename string, src interface{}, conf Config) (f *ast.File, err error) {
		fs, _ := memfs.File(filename, src)
		pkgs, err := ParseFSEntries(fset, fs, []string{filename}, confReal)
		if err != nil {
			return
		}
		for _, pkg := range pkgs {
			for _, file := range pkg.Files {
				f = file
				return
			}
		}
		t.Fatal("TestParseEntries: no source?")
		return nil, syscall.ENOENT
	})
}

func TestParseEntry(t *testing.T) {
	doTestParseEntry(t, ParseEntry)
}

func doTestParseEntry(t *testing.T, parseEntry func(fset *token.FileSet, filename string, src interface{}, conf Config) (f *ast.File, err error)) {
	fset := token.NewFileSet()
	src, err := os.ReadFile("./_testdata/functype/functype.go")
	if err != nil {
		t.Fatal("os.ReadFile:", err)
	}
	conf := Config{}
	t.Run(".gop file", func(t *testing.T) {
		f, err := parseEntry(fset, "./functype.gop", src, conf)
		if err != nil {
			t.Fatal("ParseEntry failed:", err)
		}
		if f.IsClass || f.IsProj || f.IsNormalGox {
			t.Fatal("ParseEntry functype.gop:", f.IsClass, f.IsProj, f.IsNormalGox)
		}
	})
	t.Run(".gox file", func(t *testing.T) {
		f, err := parseEntry(fset, "./functype.gox", src, conf)
		if err != nil {
			t.Fatal("ParseEntry failed:", err)
		}
		if !f.IsClass || f.IsProj || !f.IsNormalGox {
			t.Fatal("ParseEntry functype.gox:", f.IsClass, f.IsProj, f.IsNormalGox)
		}
	})
	t.Run(".foo.gox file", func(t *testing.T) {
		f, err := parseEntry(fset, "./functype.foo.gox", src, conf)
		if err != nil {
			t.Fatal("ParseEntry failed:", err)
		}
		if !f.IsClass || f.IsProj || !f.IsNormalGox {
			t.Fatal("ParseEntry functype.foo.gox:", f.IsClass, f.IsProj, f.IsNormalGox)
		}
	})
	t.Run(".foo file", func(t *testing.T) {
		_, err := parseEntry(fset, "./functype.foo", src, conf)
		if err != ErrUnknownFileKind {
			t.Fatal("ParseEntry failed:", err)
		}
	})
	t.Run(".spx file", func(t *testing.T) {
		f, err := parseEntry(fset, "./main.spx", src, conf)
		if err != nil {
			t.Fatal("ParseEntry failed:", err)
		}
		if !f.IsClass || !f.IsProj || f.IsNormalGox {
			t.Fatal("ParseEntry main.spx:", f.IsClass, f.IsProj, f.IsNormalGox)
		}
	})
}

func TestParseEntry2(t *testing.T) {
	fset := token.NewFileSet()
	src, err := os.ReadFile("./_testdata/functype/functype.go")
	if err != nil {
		t.Fatal("os.ReadFile:", err)
	}
	conf := Config{}
	conf.ClassKind = func(fname string) (isProj bool, ok bool) {
		if strings.HasSuffix(fname, "_yap.gox") {
			return true, true
		}
		return defaultClassKind(fname)
	}
	t.Run("_yap.gox file", func(t *testing.T) {
		f, err := ParseEntry(fset, "./functype_yap.gox", src, conf)
		if err != nil {
			t.Fatal("ParseEntry failed:", err)
		}
		if !f.IsClass || !f.IsProj || f.IsNormalGox {
			t.Fatal("ParseEntry functype_yap.gox:", f.IsClass, f.IsProj, f.IsNormalGox)
		}
	})
}

func TestSaveAbsFile(t *testing.T) {
	fset := token.NewFileSet()
	src, err := os.ReadFile("./_testdata/functype/functype.go")
	if err != nil {
		t.Fatal("os.ReadFile:", err)
	}
	conf := Config{}
	conf.Mode = SaveAbsFile
	t.Run(".gop file", func(t *testing.T) {
		f, err := ParseEntry(fset, "./functype.gop", src, conf)
		if err != nil {
			t.Fatal("ParseEntry failed:", err)
		}
		if f.IsClass || f.IsProj || f.IsNormalGox {
			t.Fatal("ParseEntry functype.gop:", f.IsClass, f.IsProj, f.IsNormalGox)
		}
	})
	t.Run(".gop file", func(t *testing.T) {
		f, err := ParseFile(fset, "./functype.gop", src, conf.Mode)
		if err != nil {
			t.Fatal("ParseEntry failed:", err)
		}
		if f.IsClass || f.IsProj || f.IsNormalGox {
			t.Fatal("ParseEntry functype.gop:", f.IsClass, f.IsProj, f.IsNormalGox)
		}
	})
	t.Run("dir", func(t *testing.T) {
		_, err := ParseDirEx(fset, "./_nofmt/cmdlinestyle1", conf)
		if err != nil {
			t.Fatal("ParseDirEx failed:", err)
		}
	})
}

func TestGopAutoGen(t *testing.T) {
	fset := token.NewFileSet()
	fs := memfs.SingleFile("/foo", "gop_autogen.go", ``)
	pkgs, err := ParseFSDir(fset, fs, "/foo", Config{})
	if err != nil {
		t.Fatal("ParseFSDir:", err)
	}
	if len(pkgs) != 0 {
		t.Fatal("TestGopAutoGen:", len(pkgs))
	}
}

func TestGoFile(t *testing.T) {
	fset := token.NewFileSet()
	fs := memfs.SingleFile("/foo", "test.go", `package foo`)
	pkgs, err := ParseFSDir(fset, fs, "/foo", Config{})
	if err != nil {
		t.Fatal("ParseFSDir:", err)
	}
	if len(pkgs) != 1 || pkgs["foo"].GoFiles["/foo/test.go"] == nil {
		t.Fatal("TestGoFile:", len(pkgs))
	}
}

func TestErrParse(t *testing.T) {
	fset := token.NewFileSet()
	fs := memfs.SingleFile("/foo", "test.go", `package foo bar`)
	_, err := ParseFSDir(fset, fs, "/foo", Config{})
	if err == nil {
		t.Fatal("ParseFSDir test.go: no error?")
	}

	fs = memfs.SingleFile("/foo", "test.gop", `package foo bar`)
	_, err = ParseFSDir(fset, fs, "/foo", Config{})
	if err == nil {
		t.Fatal("ParseFSDir test.gop: no error?")
	}

	fs = memfs.New(map[string][]string{"/foo": {"test.go"}}, map[string]string{})
	_, err = ParseFSDir(fset, fs, "/foo", Config{})
	if err != syscall.ENOENT {
		t.Fatal("ParseFSDir:", err)
	}

	fs = memfs.SingleFile("/foo", "test.abc.gox", `package foo`)
	_, err = ParseFSDir(fset, fs, "/foo", Config{})
	if err != nil {
		t.Fatal("ParseFSDir failed:", err)
	}
}

func testFromDir(t *testing.T, sel, relDir string) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Getwd failed:", err)
	}
	dir = path.Join(dir, relDir)
	fis, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal("ReadDir failed:", err)
	}
	for _, fi := range fis {
		name := fi.Name()
		if strings.HasPrefix(name, "_") {
			continue
		}
		t.Run(name, func(t *testing.T) {
			testFrom(t, dir+"/"+name, sel, 0)
		})
	}
}

func TestFromTestdata(t *testing.T) {
	testFromDir(t, "", "./_testdata")
}

func TestFromNofmt(t *testing.T) {
	testFromDir(t, "", "./_nofmt")
}

// -----------------------------------------------------------------------------
