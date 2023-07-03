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
	"testing"

	"github.com/goplus/gop/parser/parsertest"
	"github.com/goplus/gop/token"
)

// -----------------------------------------------------------------------------

func panicMsg(e interface{}) string {
	switch v := e.(type) {
	case string:
		return v
	case error:
		return v.Error()
	}
	return ""
}

func testErrCode(t *testing.T, code string, errExp, panicExp string) {
	defer func() {
		if e := recover(); e != nil {
			if panicMsg(e) != panicExp {
				t.Fatal("testErrCode panic:", e)
			}
		}
	}()
	fset := token.NewFileSet()
	_, err := Parse(fset, "/foo/bar.gop", code, 0)
	if err == nil || err.Error() != errExp {
		t.Fatal("testErrCode error:", err)
	}
}

func testClassErrCode(t *testing.T, code string, errExp, panicExp string) {
	defer func() {
		if e := recover(); e != nil {
			if panicMsg(e) != panicExp {
				t.Fatal("testErrCode panic:", e)
			}
		}
	}()
	fset := token.NewFileSet()
	_, err := Parse(fset, "/foo/bar.gopx", code, ParseGoPlusClass)
	if err == nil || err.Error() != errExp {
		t.Fatal("testErrCode error:", err)
	}
}

func TestErrTuple(t *testing.T) {
	testErrCode(t, `println (1,2)*2`, `/foo/bar.gop:1:9: tuple is not supported`, ``)
	testErrCode(t, `println 2*(1,2)`, `/foo/bar.gop:1:13: expected ')', found ','`, ``)
	testErrCode(t, `func test() (int,int) { return (100,100)`, `/foo/bar.gop:1:32: tuple is not supported`, ``)
}

func TestErrOperand(t *testing.T) {
	testErrCode(t, `a :=`, `/foo/bar.gop:1:5: expected operand, found 'EOF'`, ``)
}

func TestErrMissingComma(t *testing.T) {
	testErrCode(t, `func a(b int c)`, `/foo/bar.gop:1:14: missing ',' in parameter list`, ``)
}

func TestErrLambda(t *testing.T) {
	testErrCode(t, `func test(v string, f func( int)) {
}
test "hello" => {
	println "lambda",x
}
`, `/foo/bar.gop:3:6: expected 'IDENT', found "hello"`, ``)
	testErrCode(t, `func test(v string, f func( int)) {
}
test "hello", "x" => {
	println "lambda",x
}
`, `/foo/bar.gop:3:15: expected 'IDENT', found "x"`, ``)
	testErrCode(t, `func test(v string, f func( int)) {
}
test "hello", ("x") => {
	println "lambda",x
}
`, `/foo/bar.gop:3:16: expected 'IDENT', found "x"`, ``)
	testErrCode(t, `func test(v string, f func(int,int)) {
}
test "hello", (x, "y") => {
	println "lambda",x,y
}
`, `/foo/bar.gop:3:19: expected 'IDENT', found "y"`, ``)
}

func TestErrTooMany(t *testing.T) {
	testErrCode(t, `
func f() { var }
func g() { var }
func h() { var }
func h() { var }
func h() { var }
func h() { var }
func h() { var }
func h() { var }
func h() { var }
func h() { var }
func h() { var }
`, `/foo/bar.gop:2:16: expected 'IDENT', found '}' (and 10 more errors)`, ``)
}

// -----------------------------------------------------------------------------

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
                      Value:
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

func TestClassErrCode(t *testing.T) {
	testClassErrCode(t, `var (
	A,B
	v int
)
`, `/foo/bar.gopx:2:2: missing variable type or initialization`, ``)
	testClassErrCode(t, `var (
	A.*B
	v int
)
`, `/foo/bar.gopx:2:4: expected 'IDENT', found '*' (and 2 more errors)`, ``)
	testClassErrCode(t, `var (
	[]A
	v int
)
`, `/foo/bar.gopx:2:2: expected 'IDENT', found '['`, ``)
	testClassErrCode(t, `var (
	*[]A
	v int
)
`, `/foo/bar.gopx:2:3: expected 'IDENT', found '['`, ``)
	testClassErrCode(t, `var (
	v int = 10
)
`, `/foo/bar.gopx:2:8: syntax error: cannot assign value to field in class file`, ``)
	testClassErrCode(t, `var (
	v = 10
)
`, `/foo/bar.gopx:2:4: syntax error: cannot assign value to field in class file`, ``)
}

// -----------------------------------------------------------------------------
