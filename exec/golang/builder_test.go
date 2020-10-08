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

func TestIndex(t *testing.T) {
	codeExp := `package main

import fmt "fmt"

var a []float64

func main() {
	a = []float64{3.2, 1.2, 2.4}
	a[1] = 1.6
	fmt.Println(a[0], &a[1])
}
`
	println, _ := I.FindFuncv("println")
	a := NewVar(reflect.SliceOf(exec.TyFloat64), "a")
	code := NewBuilder("main", nil, nil).Interface().
		DefineVar(a).
		EndStmt(nil, &stmtState{rhsBase: 0}).
		Push(3.2).
		Push(1.2).
		Push(2.4).
		MakeArray(reflect.SliceOf(exec.TyFloat64), 3).
		StoreVar(a).
		EndStmt(nil, &stmtState{rhsBase: 0}).
		Push(1.6).
		LoadVar(a).
		Push(1).
		SetIndex(-1).
		EndStmt(nil, &stmtState{rhsBase: 0}).
		LoadVar(a).
		Index(0).
		LoadVar(a).
		AddrIndex(1).
		CallGoFuncv(println, 2, 2).
		EndStmt(nil, &stmtState{rhsBase: 0}).
		Resolve()

	codeGen := code.(*Code).String()
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

	code := NewBuilder("main", nil, nil).Interface().
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

	codeGen := code.(*Code).String()
	if codeGen != codeExp {
		fmt.Println(codeGen)
		t.Fatal("TestGoVar failed: codeGen != codeExp")
	}
}

type testInfo struct {
	Info string
}

type testPoint struct {
	X int
	Y int
}

type testRect struct {
	testInfo
	Pt1 testPoint
	Pt2 *testPoint
}

var gRect *testRect

func init() {
	gRect = &testRect{}
	gRect.Pt2 = &testPoint{}
}

func getTestRect() *testRect {
	return gRect
}

func execTestRect(_ int, p *qexec.Context) {
	p.Ret(0, getTestRect())
}

func TestGoField(t *testing.T) {
	codeExp := `package main

import (
	fmt "fmt"
	golang "github.com/goplus/gop/exec/golang"
	pkg_field "pkg_field"
)

var y golang.testRect

func main() {
	y = pkg_field.Rect2
	pkg_field.Rect.Info = "hello"
	pkg_field.Rect.Pt1.X = -1
	pkg_field.Rect.Pt2.Y = -2
	y.Info = "world"
	y.Pt1.X = -10
	y.Pt2.Y = -20
	pkg_field.GetRect().Info = "next"
	pkg_field.GetRect().Pt1.X = 101
	pkg_field.GetRect().Pt2.Y = 102
	pkg_field.Rect2 = y
	fmt.Println(pkg_field.Rect.Info, pkg_field.Rect.Pt2.Y, y.Info, y.Pt2.Y, pkg_field.GetRect().Info, pkg_field.GetRect().Pt2.Y, &pkg_field.Rect.Pt1, &y.Pt1, &pkg_field.GetRect().Pt1)
}
`
	println, _ := I.FindFuncv("println")

	pkg := qexec.NewGoPackage("pkg_field")

	rc := testRect{}
	rc.Info = "Info"
	rc.Pt1 = testPoint{10, 20}
	rc.Pt2 = &testPoint{30, 40}

	rc2 := testRect{}
	rc2.Pt2 = &testPoint{}

	pkg.RegisterVars(
		pkg.Var("Rect", &rc),
		pkg.Var("Rect2", &rc2),
	)
	pkg.RegisterFuncs(
		pkg.Func("GetRect", getTestRect, execTestRect),
	)

	x, ok := pkg.FindVar("Rect")
	if !ok {
		t.Fatal("FindVar failed:", x)
	}
	xtyp := reflect.TypeOf(rc)

	x2, ok := pkg.FindVar("Rect2")
	if !ok {
		t.Fatal("FindVar failed:", x2)
	}
	fnTestRect, ok := pkg.FindFunc("GetRect")
	if !ok {
		t.Fatal("FindFunc failed:", fnTestRect)
	}

	gtyp := reflect.TypeOf(gRect)
	ytyp := reflect.TypeOf(rc)
	y := NewVar(ytyp, "y")
	u := NewVar(exec.TyInt, "_")
	b := NewBuilder("main", nil, nil)

	code := b.Interface().
		DefineVar(y). // var y testRect
		DefineVar(u). // var _ int
		EndStmt(nil, &stmtState{rhsBase: 0}).
		LoadGoVar(x2).
		StoreVar(y). // y = pkg_field.Rect2
		EndStmt(nil, &stmtState{rhsBase: 0}).
		Push("hello").
		AddrGoVar(x).
		StoreField(xtyp, []int{0, 0}). // pkg_field.Rect.Info = "hello"
		EndStmt(nil, &stmtState{rhsBase: 0}).
		Push(-1).
		AddrGoVar(x).
		StoreField(xtyp, []int{1, 0}). // pkg_field.Rect.Pt1.X = -1
		EndStmt(nil, &stmtState{rhsBase: 0}).
		Push(-2).
		AddrGoVar(x).
		StoreField(xtyp, []int{2, 1}). // pkg_field.Rect.Pt2.Y = -2
		EndStmt(nil, &stmtState{rhsBase: 0}).
		Push("world").
		AddrVar(y).
		StoreField(ytyp, []int{0, 0}). // y.Info = "world"
		EndStmt(nil, &stmtState{rhsBase: 0}).
		Push(-10).
		AddrVar(y).
		StoreField(ytyp, []int{1, 0}). // y.Pt1.X = -10
		EndStmt(nil, &stmtState{rhsBase: 0}).
		Push(-20).
		AddrVar(y).
		StoreField(ytyp, []int{2, 1}). // y.Pt2.Y = -20
		EndStmt(nil, &stmtState{rhsBase: 0}).
		Push("next").
		CallGoFunc(fnTestRect, 0).
		StoreField(gtyp, []int{0, 0}). // pkg_field.GetRect().Info = "next"
		EndStmt(nil, &stmtState{rhsBase: 0}).
		Push(101).
		CallGoFunc(fnTestRect, 0).
		StoreField(gtyp, []int{1, 0}). // pkg_field.GetRect().Pt1.X = 101
		EndStmt(nil, &stmtState{rhsBase: 0}).
		Push(102).
		CallGoFunc(fnTestRect, 0).
		StoreField(gtyp, []int{2, 1}). // pkg_field.GetRect().Pt2.Y = 102
		EndStmt(nil, &stmtState{rhsBase: 0}).
		LoadVar(y).
		StoreGoVar(x2). // pkg_field.Rect2 = y
		EndStmt(nil, &stmtState{rhsBase: 0}).
		LoadGoVar(x).
		LoadField(xtyp, []int{0, 0}).
		LoadGoVar(x).
		LoadField(xtyp, []int{2, 1}).
		LoadVar(y).
		LoadField(ytyp, []int{0, 0}).
		LoadVar(y).
		LoadField(ytyp, []int{2, 1}).
		CallGoFunc(fnTestRect, 0).
		LoadField(gtyp, []int{0, 0}).
		CallGoFunc(fnTestRect, 0).
		LoadField(gtyp, []int{2, 1}).
		AddrGoVar(x).
		AddrField(xtyp, []int{1}).
		AddrVar(y).
		AddrField(ytyp, []int{1}).
		CallGoFunc(fnTestRect, 0).
		AddrField(gtyp, []int{1}).
		CallGoFuncv(println, 9, 9). // println(pkg_field.Rect.Info, pkg_field.Rect.Pt2.Y, y.Info, y.Pt2.Y, pkg_field.GetRect().Info, pkg_field.GetRect().Pt2.Y, &pkg_field.Rect.Pt1, &y.Pt1, &pkg_field.GetRect().Pt1)
		EndStmt(nil, &stmtState{rhsBase: 0}).
		Resolve()

	if v := code.(*Code); v.String() != codeExp {
		fmt.Println(v.String())
		t.Fatal("TestGoVar failed: codeGen != codeExp")
	}
}

// -----------------------------------------------------------------------------

func TestStruct(t *testing.T) {
	println, _ := I.FindFuncv("println")

	codeExp := `package main

import fmt "fmt"

func main() {
	fmt.Println(struct {
		Name string
		Age  int
	}{Name: "bar", Age: 30})
}
`
	fields := []reflect.StructField{
		{
			Name: "Name",
			Type: reflect.TypeOf(""),
		},
		{
			Name: "Age",
			Type: reflect.TypeOf(0),
		},
	}
	typ := reflect.StructOf(fields)

	code := NewBuilder("main", nil, nil).Interface().
		Push(0).
		Push("bar").
		Push(1).
		Push(30).
		Struct(typ, 2).
		CallGoFuncv(println, 1, 1).
		EndStmt(nil, &stmtState{rhsBase: 0}).
		Resolve()

	codeGen := code.(*Code).String()
	if codeGen != codeExp {
		fmt.Println(codeGen)
		fmt.Println(codeExp)
		t.Fatal("TestStruct failed: codeGen != codeExp")
	}
}

func TestStruct2(t *testing.T) {
	println, _ := I.FindFuncv("println")

	codeExp := `package main

import fmt "fmt"

func main() {
	fmt.Println(&struct {
		Name string
		Age  int
	}{Name: "bar", Age: 30})
}
`
	fields := []reflect.StructField{
		{
			Name: "Name",
			Type: reflect.TypeOf(""),
		},
		{
			Name: "Age",
			Type: reflect.TypeOf(0),
		},
	}
	typ := reflect.StructOf(fields)

	code := NewBuilder("main", nil, nil).Interface().
		Push(0).
		Push("bar").
		Push(1).
		Push(30).
		Struct(reflect.PtrTo(typ), 2).
		CallGoFuncv(println, 1, 1).
		EndStmt(nil, &stmtState{rhsBase: 0}).
		Resolve()

	codeGen := code.(*Code).String()
	if codeGen != codeExp {
		fmt.Println(codeGen)
		fmt.Println(codeExp)
		t.Fatal("TestStruct failed: codeGen != codeExp")
	}
}

func TestType(t *testing.T) {
	println, _ := I.FindFuncv("println")

	codeExp := `package main

import fmt "fmt"

type Person struct {
	Name string
	Age  int
}

var p Person

func main() {
	p = Person{Name: "bar", Age: 30}
	fmt.Println(p)
}
`
	typ := reflect.StructOf([]reflect.StructField{
		{
			Name: "Name",
			Type: exec.TyString,
		},
		{
			Name: "Age",
			Type: exec.TyInt,
		},
	})

	p := NewVar(typ, "p")

	code := NewBuilder("main", nil, nil).Interface().
		DefineVar(p).
		DefineType(typ, "Person").
		EndStmt(nil, &stmtState{rhsBase: 0}).
		Push(0).
		Push("bar").
		Push(1).
		Push(30).
		Struct(typ, 2).
		StoreVar(p).
		EndStmt(nil, &stmtState{rhsBase: 0}).
		LoadVar(p).
		CallGoFuncv(println, 1, 1).
		EndStmt(nil, &stmtState{rhsBase: 0}).
		Resolve()

	codeGen := code.(*Code).String()
	if codeGen != codeExp {
		fmt.Println(codeGen)
		fmt.Println(codeExp)
		t.Fatal("TestStruct failed: codeGen != codeExp")
	}
}

func TestReserved(t *testing.T) {
	code := NewBuilder("main", nil, nil)
	off := code.Reserve()
	code.ReservedAsPush(off, 123)
	if !reflect.DeepEqual(code.reserveds[off].Expr.(*ast.BasicLit), IntConst(123)) {
		t.Fatal("TestReserved failed: reserveds is not set to", 123)
	}
}

// -----------------------------------------------------------------------------
