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
	"io/fs"
	"testing"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/token"
	fsx "github.com/qiniu/x/http/fs"
)

// -----------------------------------------------------------------------------

type fileInfo struct {
	*fsx.FileInfo
}

func (p fileInfo) Info() (fs.FileInfo, error) {
	return nil, fs.ErrNotExist
}

func TestFilter(t *testing.T) {
	d := fsx.NewFileInfo("foo.go", 10)
	if filter(d, func(fi fs.FileInfo) bool {
		return false
	}) {
		t.Fatal("TestFilter: true?")
	}
	if !filter(d, func(fi fs.FileInfo) bool {
		return true
	}) {
		t.Fatal("TestFilter: false?")
	}
	d2 := fileInfo{d}
	if filter(d2, func(fi fs.FileInfo) bool {
		return true
	}) {
		t.Fatal("TestFilter: true?")
	}
}

func TestAssert(t *testing.T) {
	defer func() {
		if e := recover(); e != "go/parser internal error: panic msg" {
			t.Fatal("TestAssert:", e)
		}
	}()
	assert(false, "panic msg")
}

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
	t.Helper()
	fset := token.NewFileSet()
	_, err := Parse(fset, "/foo/bar.gop", code, 0)
	if err == nil || err.Error() != errExp {
		t.Fatal("testErrCode error:", err)
	}
}

func testErrCodeParseExpr(t *testing.T, code string, errExp, panicExp string) {
	defer func() {
		if e := recover(); e != nil {
			if panicMsg(e) != panicExp {
				t.Fatal("testErrCodeParseExpr panic:", e)
			}
		}
	}()
	t.Helper()
	_, err := ParseExpr(code)
	if err == nil || err.Error() != errExp {
		t.Fatal("testErrCodeParseExpr error:", err)
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
	t.Helper()
	fset := token.NewFileSet()
	_, err := Parse(fset, "/foo/bar.gox", code, ParseGoPlusClass)
	if err == nil || err.Error() != errExp {
		t.Fatal("testErrCode error:", err)
	}
}

func TestErrLabel(t *testing.T) {
	testErrCode(t, `a.x:`, `/foo/bar.gop:1:4: illegal label declaration`, ``)
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

func TestErrTooManyParseExpr(t *testing.T) {
	testErrCodeParseExpr(t, `func() int {
  var
  var
  var
  var
  var
  var
  var
  var
  var
  var
  var
  var
}()
`, `3:3: expected 'IDENT', found 'var' (and 10 more errors)`, ``)
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
func h() { var }
`, `/foo/bar.gop:2:16: expected 'IDENT', found '}' (and 10 more errors)`, ``)
}

func TestErrInFunc(t *testing.T) {
	testErrCode(t, `func test() {
	a,
}`, `/foo/bar.gop:2:2: expected 1 expression (and 2 more errors)`, ``)
	testErrCode(t, `func test() {
	a.test, => {
	}
}`, `/foo/bar.gop:2:10: expected operand, found '=>' (and 1 more errors)`, ``)
	testErrCode(t, `func test() {
		,
	}
}`, `/foo/bar.gop:2:3: expected statement, found ',' (and 1 more errors)`, ``)
}

// -----------------------------------------------------------------------------

func TestClassErrCode(t *testing.T) {
	testClassErrCode(t, `var (
	A,B
	v int
)
`, `/foo/bar.gox:2:2: missing variable type or initialization`, ``)
	testClassErrCode(t, `var (
	A.*B
	v int
)
`, `/foo/bar.gox:2:4: expected 'IDENT', found '*' (and 2 more errors)`, ``)
	testClassErrCode(t, `var (
	[]A
	v int
)
`, `/foo/bar.gox:2:2: expected 'IDENT', found '['`, ``)
	testClassErrCode(t, `var (
	*[]A
	v int
)
`, `/foo/bar.gox:2:3: expected 'IDENT', found '['`, ``)
	testClassErrCode(t, `var (
	v int = 10
)
`, `/foo/bar.gox:2:8: syntax error: cannot assign value to field in class file`, ``)
	testClassErrCode(t, `var (
	v = 10
)
`, `/foo/bar.gox:2:4: syntax error: cannot assign value to field in class file`, ``)
	testClassErrCode(t, `
var (
)
const c = 100
const d
`, `/foo/bar.gox:5:7: missing constant value`, ``)
}

func TestErrGlobal(t *testing.T) {
	testErrCode(t, `func test() {}
}`, `/foo/bar.gop:2:1: expected statement, found '}'`, ``)
}

func TestErrCompositeLiteral(t *testing.T) {
	testErrCode(t, `println (T[int]){a: 1, b: 2}
`, `/foo/bar.gop:1:10: cannot parenthesize type in composite literal`, ``)
}

func TestErrSelectorExpr(t *testing.T) {
	testErrCode(t, `
x.
*p
`, `/foo/bar.gop:3:1: expected selector or type assertion, found '*'`, ``)
}

func TestErrStringLitEx(t *testing.T) {
	testErrCode(t, `
println "${ ... }"
`, "/foo/bar.gop:2:13: expected operand, found '...'", ``)
	testErrCode(t, `
println "${b"
`, "/foo/bar.gop:2:11: invalid $ expression: ${ doesn't end with }", ``)
	testErrCode(t, `
println "$a${b}"
`, "/foo/bar.gop:2:10: invalid $ expression: neither `${ ... }` nor `$$`", ``)
}

func TestErrStringLiteral(t *testing.T) {
	testErrCode(t, `run "
`, `/foo/bar.gop:1:5: string literal not terminated`, ``)
}

func TestErrFieldDecl(t *testing.T) {
	testErrCode(t, `
type T struct {
	*(Foo)
}
`, `/foo/bar.gop:3:3: cannot parenthesize embedded type`, ``)
	testErrCode(t, `
type T struct {
	(Foo)
}
`, `/foo/bar.gop:3:2: cannot parenthesize embedded type`, ``)
	testErrCode(t, `
type T struct {
	(*Foo)
}
`, `/foo/bar.gop:3:2: cannot parenthesize embedded type`, ``)
}

func TestParseFieldDecl(t *testing.T) {
	var p parser
	p.init(token.NewFileSet(), "/foo/bar.gop", []byte(`type T struct {
}
`), 0)
	p.parseFieldDecl(nil)
}

func TestCheckExpr(t *testing.T) {
	var p parser
	p.init(token.NewFileSet(), "/foo/bar.gop", []byte(``), 0)
	p.checkExpr(&ast.Ellipsis{})
	p.checkExpr(&ast.ElemEllipsis{})
	p.checkExpr(&ast.StarExpr{})
	p.checkExpr(&ast.IndexListExpr{})
	p.checkExpr(&ast.FuncType{})
	p.checkExpr(&ast.FuncLit{})
}

func TestErrFuncDecl(t *testing.T) {
	testErrCode(t, `func test()
{
}
`, `/foo/bar.gop:2:1: unexpected semicolon or newline before {`, ``)
	testErrCode(t, `func test() +1
`, `/foo/bar.gop:1:13: expected ';', found '+'`, ``)
	testErrCode(t, `
func (a T) +{}
`, `/foo/bar.gop:2:12: expected type, found '+'`, ``)
	testErrCode(t, `func +(a T, b T) {}
`, `/foo/bar.gop:1:6: overload operator can only have one parameter`, ``)
}

func TestNumberUnitLit(t *testing.T) {
	var p parser
	p.checkExpr(&ast.NumberUnitLit{})
	p.toIdent(&ast.NumberUnitLit{})
}

// -----------------------------------------------------------------------------
