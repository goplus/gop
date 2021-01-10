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
	"testing"

	"github.com/goplus/gop/parser/parsertest"
	"github.com/goplus/gop/token"
)

// -----------------------------------------------------------------------------

var fsTestStd = parsertest.NewSingleFileFS("/foo", "bar.gop", `package bar; import "io"
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
	parsertest.Expect(t, bar, `package bar

file bar.gop
noEntrypoint
ast.GenDecl:
  Tok: import
  ast.ImportSpec:
    ast.BasicLit:
      Kind: STRING
      Value: "io"
ast.FuncDecl:
  ast.Ident:
    Name: init
  ast.FuncType:
    ast.FieldList:
  ast.BlockStmt:
    ast.AssignStmt:
      ast.Ident:
        Name: x
      Tok: :=
      ast.BasicLit:
        Kind: INT
        Value: 0
    ast.IfStmt:
      ast.AssignStmt:
        ast.Ident:
          Name: t
        Tok: :=
        ast.Ident:
          Name: false
      ast.Ident:
        Name: t
      ast.BlockStmt:
        ast.AssignStmt:
          ast.Ident:
            Name: x
          Tok: =
          ast.BasicLit:
            Kind: INT
            Value: 3
      ast.BlockStmt:
        ast.AssignStmt:
          ast.Ident:
            Name: x
          Tok: =
          ast.BasicLit:
            Kind: INT
            Value: 5
    ast.ExprStmt:
      ast.CallExpr:
        ast.Ident:
          Name: println
        ast.BasicLit:
          Kind: STRING
          Value: "x:"
        ast.Ident:
          Name: x
    ast.AssignStmt:
      ast.Ident:
        Name: x
      Tok: =
      ast.BasicLit:
        Kind: INT
        Value: 0
    ast.SwitchStmt:
      ast.AssignStmt:
        ast.Ident:
          Name: s
        Tok: :=
        ast.BasicLit:
          Kind: STRING
          Value: "Hello"
      ast.Ident:
        Name: s
      ast.BlockStmt:
        ast.CaseClause:
          ast.AssignStmt:
            ast.Ident:
              Name: x
            Tok: =
            ast.BasicLit:
              Kind: INT
              Value: 7
        ast.CaseClause:
          ast.BasicLit:
            Kind: STRING
            Value: "world"
          ast.BasicLit:
            Kind: STRING
            Value: "hi"
          ast.AssignStmt:
            ast.Ident:
              Name: x
            Tok: =
            ast.BasicLit:
              Kind: INT
              Value: 5
        ast.CaseClause:
          ast.BasicLit:
            Kind: STRING
            Value: "xsw"
          ast.AssignStmt:
            ast.Ident:
              Name: x
            Tok: =
            ast.BasicLit:
              Kind: INT
              Value: 3
    ast.ExprStmt:
      ast.CallExpr:
        ast.Ident:
          Name: println
        ast.BasicLit:
          Kind: STRING
          Value: "x:"
        ast.Ident:
          Name: x
    ast.AssignStmt:
      ast.Ident:
        Name: c
      Tok: :=
      ast.CallExpr:
        ast.Ident:
          Name: make
        ast.ChanType:
          ast.Ident:
            Name: bool
        ast.BasicLit:
          Kind: INT
          Value: 100
    ast.SelectStmt:
      ast.BlockStmt:
        ast.CommClause:
          ast.SendStmt:
            ast.Ident:
              Name: c
            ast.Ident:
              Name: true
        ast.CommClause:
          ast.AssignStmt:
            ast.Ident:
              Name: v
            Tok: :=
            ast.UnaryExpr:
              Op: <-
              ast.Ident:
                Name: c
        ast.CommClause:
          ast.ExprStmt:
            ast.CallExpr:
              ast.Ident:
                Name: panic
              ast.BasicLit:
                Kind: STRING
                Value: "error"
`)
}

// -----------------------------------------------------------------------------
