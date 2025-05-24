//go:build go1.18
// +build go1.18

/*
 * Copyright (c) 2022 The XGo Authors (xgo.dev). All rights reserved.
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
	"testing"

	"github.com/goplus/xgo/parser/parsertest"
	"github.com/goplus/xgo/token"
)

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
	pkgs, err := Parse(fset, "/foo/bar.xgo", testStdCode, ParseComments)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("Parse failed:", err, len(pkgs))
	}
	bar := pkgs["bar"]
	parsertest.Expect(t, bar, `package bar

file bar.xgo
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

func TestFromInstance(t *testing.T) {
	testFromDir(t, "", "./_instance")
}

// -----------------------------------------------------------------------------
