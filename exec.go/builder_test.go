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

package exec

import (
	"fmt"
	"testing"

	"github.com/qiniu/qlang/v6/exec.spec"
	"github.com/qiniu/qlang/v6/token"

	qexec "github.com/qiniu/qlang/v6/exec"
	_ "github.com/qiniu/qlang/v6/lib/builtin"
)

// I is a Go package instance.
var I = qexec.FindGoPackage("")

// -----------------------------------------------------------------------------

func TestBasic(t *testing.T) {
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

	fmt.Println(code.String())
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
	fset := token.NewFileSet()
	foo := fset.AddFile("foo.ql", fset.Base(), 100)
	bar := fset.AddFile("bar.ql", fset.Base(), 100)
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

	fmt.Println(code.String())
}

// -----------------------------------------------------------------------------
