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
	"testing"
)

// -----------------------------------------------------------------------------

func TestIf1(t *testing.T) {
	label1 := NewLabel("a")
	label2 := NewLabel("b")
	code := NewBuilder(nil).
		Push(true).
		JmpIf(0, label1).
		Push(50).
		Push(6).
		BuiltinOp(Int, OpDiv).
		Jmp(label2).
		Label(label1).
		Push(5).
		Push(2).
		BuiltinOp(Int, OpMod).
		Label(label2).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != 8 {
		t.Fatal("50 6 div != 8, ret =", v)
	}
}

func TestIf2(t *testing.T) {
	label1 := NewLabel("a")
	label2 := NewLabel("b")
	code := NewBuilder(nil).
		Push(false).
		JmpIf(0, label1).
		Push(5).
		Push(6).
		BuiltinOp(Int, OpMul).
		Jmp(label2).
		Label(label1).
		Push(5.0).
		Push(2.0).
		BuiltinOp(Float64, OpMul).
		Label(label2).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != 10.0 {
		t.Fatal("5.0 2.0 mul != 10.0, ret =", v)
	}
}

func TestIf3(t *testing.T) {
	label1 := NewLabel("a")
	label2 := NewLabel("b")
	code := NewBuilder(nil).
		Push(true).
		JmpIf(1, label1).
		Push(5).
		Push(6).
		BuiltinOp(Int, OpMul).
		Jmp(label2).
		Label(label1).
		Push(5.0).
		Push(2.0).
		BuiltinOp(Float64, OpMul).
		Label(label2).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != 10.0 {
		t.Fatal("5.0 2.0 mul != 10.0, ret =", v)
	}
}

func TestIf4(t *testing.T) {
	label1 := NewLabel("a")
	label2 := NewLabel("b")
	code := NewBuilder(nil).
		Push(false).
		JmpIf(1, label1).
		Push(5).
		Push(6).
		BuiltinOp(Int, OpMul).
		Jmp(label2).
		Label(label1).
		Push(5.0).
		Push(2.0).
		BuiltinOp(Float64, OpMul).
		Label(label2).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != 30 {
		t.Fatal("5.0 2.0 mul != 10.0, ret =", v)
	}
}

// -----------------------------------------------------------------------------

func TestCase1(t *testing.T) {
	done := NewLabel("done")
	label1 := NewLabel("a")
	label2 := NewLabel("b")
	code := NewBuilder(nil).
		Push(1).
		Push(1).
		CaseNE(label1, 1).
		Push(5).
		Push(6).
		BuiltinOp(Int, OpAdd).
		Jmp(done).
		Label(label1).
		Push(2).
		CaseNE(label2, 1).
		Push(5).
		Push(2).
		BuiltinOp(Int, OpMod).
		Jmp(done).
		Label(label2).
		Default().
		Push(100).
		Label(done).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != 11 {
		t.Fatal("5 6 add != 11, ret =", v)
	}
}

func TestCase2(t *testing.T) {
	done := NewLabel("done")
	label1 := NewLabel("a")
	label2 := NewLabel("b")
	code := NewBuilder(nil).
		Push(2).
		Push(1).
		Push(0).
		CaseNE(label1, 2).
		Push(5).
		Push(6).
		BuiltinOp(Int, OpMul).
		Jmp(done).
		Label(label1).
		Push(3).
		Push(2).
		Push(100).
		CaseNE(label2, 3).
		Push(5).
		Push(2).
		BuiltinOp(Int, OpSub).
		Jmp(done).
		Label(label2).
		Default().
		Push(100).
		Label(done).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != 3 {
		t.Fatal("5 2 sub != 3, ret =", v)
	}
}

func TestDefault(t *testing.T) {
	done := NewLabel("done")
	label1 := NewLabel("a")
	label2 := NewLabel("b")
	code := NewBuilder(nil).
		Push(3).
		Push(1).
		CaseNE(label1, 1).
		Push(5).
		Push(6).
		BuiltinOp(Int, OpMul).
		Jmp(done).
		Label(label1).
		Push(2).
		CaseNE(label2, 1).
		Push(5).
		Push(2).
		BuiltinOp(Int, OpMod).
		Jmp(done).
		Label(label2).
		Default().
		Push(100).
		Label(done).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != 100 {
		t.Fatal("100 != 100, ret =", v)
	}
}

// -----------------------------------------------------------------------------
