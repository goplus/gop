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

	"github.com/qiniu/goplus/cl"
	"github.com/qiniu/goplus/exec.spec"
	"github.com/qiniu/x/log"

	qexec "github.com/qiniu/goplus/exec/bytecode"
	_ "github.com/qiniu/goplus/lib"
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

// -----------------------------------------------------------------------------
