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

package golang

import (
	"fmt"
	"go/ast"
	"reflect"
	"testing"

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/exec.spec"
	"github.com/qiniu/x/log"

	qexec "github.com/goplus/gop/exec/bytecode"
	_ "github.com/goplus/gop/lib"
)

// I is a Go package instance.
var I = qexec.FindGoPackage("")

func init() {
	cl.CallBuiltinOp = qexec.CallBuiltinOp
	log.SetFlags(log.Ldefault &^ log.LstdFlags)
	log.SetOutputLevel(0)
}

// -----------------------------------------------------------------------------

func TestBuild(t *testing.T) {
	codeExp := `package main

import fmt "fmt"

func main() {
	fmt.Println(1 + 2)
	fmt.Println(complex64((3 + 2i)))
}
`
	println, _ := I.FindFuncv("println")
	code := NewBuilder("main", nil, nil).
		Push(1).
		Push(2).
		BuiltinOp(exec.Int, exec.OpAdd).
		CallGoFuncv(println, 1, 1).
		EndStmt(nil, &stmtState{rhsBase: 0}).
		Push(complex64(3+2i)).
		CallGoFuncv(println, 1, 1).
		EndStmt(nil, &stmtState{rhsBase: 0}).
		Resolve()

	codeGen := code.String()
	if codeGen != codeExp {
		fmt.Println(codeGen)
		t.Fatal("TestBasic failed: codeGen != codeExp")
	}
}

func TestGoVar(t *testing.T) {
	codeExp := `package main

import (
	fmt "fmt"
	pkg "pkg"
)

func main() {
	pkg.X, pkg.Y = 5, 6
	fmt.Println(pkg.X, pkg.Y)
	fmt.Println(&pkg.X, *&pkg.Y)
}
`
	println, _ := I.FindFuncv("println")

	pkg := qexec.NewGoPackage("pkg")
	pkg.RegisterVars(
		pkg.Var("X", new(int)),
		pkg.Var("Y", new(int)),
	)

	x, ok := pkg.FindVar("X")
	if !ok {
		t.Fatal("FindVar failed:", x)
	}
	y, ok := pkg.FindVar("Y")
	if !ok {
		t.Fatal("FindVar failed:", y)
	}

	code := NewBuilder("main", nil, nil).
		Push(5).
		Push(6).
		StoreGoVar(y).
		StoreGoVar(x).
		EndStmt(nil, &stmtState{rhsBase: 0}).
		LoadGoVar(x).
		LoadGoVar(y).
		CallGoFuncv(println, 2, 2).
		EndStmt(nil, &stmtState{rhsBase: 0}).
		AddrGoVar(x).
		AddrGoVar(y).
		AddrOp(exec.Int, exec.OpAddrVal).
		CallGoFuncv(println, 2, 2).
		EndStmt(nil, &stmtState{rhsBase: 0}).
		Resolve()

	codeGen := code.String()
	if codeGen != codeExp {
		fmt.Println(codeGen)
		t.Fatal("TestGoVar failed: codeGen != codeExp")
	}
}

// -----------------------------------------------------------------------------

func TestReserved(t *testing.T) {
	code := NewBuilder("main", nil, nil)
	off := code.Reserve()
	code.ReservedAsPush(off, 123)
	if !reflect.DeepEqual(code.reserveds[off].Expr.(*ast.BasicLit), IntConst(123)) {
		t.Fatal("TestReserved failed: reserveds is not set to", 123)
	}
}

func TestReserved2(t *testing.T) {
	code := NewBuilder("main", nil, nil)
	defer func() {
		err := recover()
		if !reflect.DeepEqual(err, "The method defer under the builder of golang is not yet supported") {
			t.Fatal("TestReserved failed: Defer is not yet supported now")
		}
	}()
	var start = NewLabel("")
	var end = NewLabel("")
	code.Interface().Defer(start, end)
}

func TestReserved3(t *testing.T) {
	code := NewBuilder("main", nil, nil)
	off := code.Reserve()
	defer func() {
		err := recover()
		if !reflect.DeepEqual(err, "todo\n") {
			t.Fatal("TestReserved failed: ReservedAsInstr is todo now")
		}
	}()
	code.Interface().ReservedAsInstr(off, nil)
}
