/*
 Copyright 2020 Qiniu Cloud (qiniu.com)

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
	"testing"

	"github.com/qiniu/qlang/v6/cl"
	"github.com/qiniu/qlang/v6/exec.spec"
	"github.com/qiniu/qlang/v6/token"

	qexec "github.com/qiniu/qlang/v6/exec/bytecode"
	_ "github.com/qiniu/qlang/v6/lib/builtin"
	_ "github.com/qiniu/qlang/v6/lib/fmt"
	_ "github.com/qiniu/qlang/v6/lib/reflect"
	_ "github.com/qiniu/qlang/v6/lib/strings"
)

// I is a Go package instance.
var I = qexec.FindGoPackage("")

func init() {
	cl.CallBuiltinOp = qexec.CallBuiltinOp
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
	code := NewBuilder(nil, nil).
		Push(1).
		Push(2).
		BuiltinOp(exec.Int, exec.OpAdd).
		CallGoFuncv(println, 1).
		EndStmt(nil).
		Push(complex64(3+2i)).
		CallGoFuncv(println, 1).
		EndStmt(nil).
		Resolve()

	codeGen := code.String()
	if codeGen != codeExp {
		fmt.Println(codeGen)
		t.Fatal("TestBasic failed: codeGen != codeExp")
	}
}

// -----------------------------------------------------------------------------

type node struct {
	pos token.Pos
}

func (p *node) Pos() token.Pos {
	return p.pos
}

func (p *node) End() token.Pos {
	return p.pos + 1
}

func TestFileLine(t *testing.T) {
	codeExp := `package main

import fmt "fmt"

func main() { 
//line ./foo.ql:1
	fmt.Println(1 + 2)
//line ./bar.ql:1
	fmt.Println(complex64((3 + 2i)))
}
`
	fset := token.NewFileSet()
	foo := fset.AddFile("./foo.ql", fset.Base(), 100)
	bar := fset.AddFile("./bar.ql", fset.Base(), 100)
	foo.SetLines([]int{0, 10, 20, 80, 100})
	bar.SetLines([]int{0, 10, 20, 80, 100})
	node1 := &node{23}
	node2 := &node{123}
	println, _ := I.FindFuncv("println")
	code := NewBuilder(nil, fset).
		Push(1).
		Push(2).
		BuiltinOp(exec.Int, exec.OpAdd).
		CallGoFuncv(println, 1).
		EndStmt(node1).
		Push(complex64(3+2i)).
		CallGoFuncv(println, 1).
		EndStmt(node2).
		Resolve()

	codeGen := code.String()
	if codeGen != codeExp {
		fmt.Println(codeGen)
		t.Fatal("TestFileLine failed: codeGen != codeExp")
	}
}

// -----------------------------------------------------------------------------
