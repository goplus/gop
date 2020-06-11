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

package bytecode

import (
	"fmt"
	"testing"

	"github.com/qiniu/goplus/exec.spec"
	"github.com/qiniu/x/errors"
)

// -----------------------------------------------------------------------------

func TestErrWrap(t *testing.T) {
	errorf, ok := I.FindFuncv("Errorf")
	if !ok {
		t.Fatal("FindFuncv failed: Errorf")
	}

	defer func() {
		if e := recover(); e != nil {
			frame, ok := e.(*errors.Frame)
			if !ok {
				t.Fatal("TestErrWrap failed:", e)
			}
			fmt.Println(frame.Args...)
		}
	}()
	frame := &errors.Frame{
		Pkg:  "main",
		Func: "TestErrWrap",
		Code: `errorf("not found")?`,
		File: `./flow_test.go`,
		Line: 49,
	}
	code := newBuilder().
		Push("arg1").
		Push("arg2").
		Push("arg3").
		Push(123).
		Push("not found").
		CallGoFuncv(errorf, 1, 1).
		ErrWrap(2, nil, frame, 3).
		Resolve()

	ctx := NewContext(code)
	ctx.base = 3
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != 8 {
		t.Fatal("v != 8, ret =", v)
	}
}

func TestErrWrap1(t *testing.T) {
	errorf, ok := I.FindFuncv("Errorf")
	if !ok {
		t.Fatal("FindFuncv failed: Errorf")
	}

	frame := &errors.Frame{
		Pkg:  "main",
		Func: "TestErrWrap",
		Code: `errorf("not found")?`,
		File: `./flow_test.go`,
		Line: 45,
	}
	retErr := NewVar(TyError, "err")
	code := newBuilder().
		DefineVar(retErr).
		Push("arg1").
		Push("arg2").
		Push("arg3").
		Push(123).
		Push("not found").
		CallGoFuncv(errorf, 1, 1).
		ErrWrap(2, retErr, frame, 3).
		Resolve()

	ctx := NewContext(code)
	ctx.base = 3
	ctx.Exec(0, code.Len())

	if e := ctx.GetVar(retErr); e != nil {
		frame, ok := e.(*errors.Frame)
		if !ok {
			t.Fatal("TestErrWrap1 failed:", e)
		}
		fmt.Println(frame.Args...)
	} else {
		t.Fatal("TestErrWrap1 failed: retErr not set")
	}
}

func TestErrWrap2(t *testing.T) {
	code := newBuilder().
		Push(123).
		Push(nil).
		ErrWrap(2, nil, nil, 0).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != 123 {
		t.Fatal("v != 123, ret =", v)
	}
}

func TestWrapIfErr(t *testing.T) {
	l := NewLabel("")
	code := newBuilder().
		Push(123).
		Push(nil).
		WrapIfErr(2, l).
		Push(0).
		Label(l).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != 123 {
		t.Fatal("v != 123, ret =", v)
	}
}

func TestWrapIfErr2(t *testing.T) {
	l := NewLabel("")
	code := newBuilder().
		Push(123).
		Push(true).
		WrapIfErr(2, l).
		Push(10).
		Label(l).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != 10 {
		t.Fatal("v != 10, ret =", v)
	}
}

func TestIf1(t *testing.T) {
	label1 := defaultImpl.NewLabel("a").(*Label)
	label2 := NewLabel("b")
	code := newBuilder().
		Push(true).
		JmpIf(exec.JcNil, label1).
		Push(50).
		Push(6).
		BuiltinOp(Int, OpQuo).
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
	_ = label1.Name()
}

func TestIf11(t *testing.T) {
	label1 := defaultImpl.NewLabel("a").(*Label)
	label2 := NewLabel("b")
	code := newBuilder().
		Push(true).
		JmpIf(0, label1).
		Push(50).
		Push(6).
		BuiltinOp(Int, OpQuo).
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
	_ = label1.Name()
}

func TestIf2(t *testing.T) {
	label1 := NewLabel("a")
	label2 := NewLabel("b")
	code := newBuilder().
		Push(nil).
		JmpIf(exec.JcNil, label1).
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

func TestIf22(t *testing.T) {
	label1 := NewLabel("a")
	label2 := NewLabel("b")
	code := newBuilder().
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
	code := newBuilder().
		Push(true).
		JmpIf(exec.JcNotNil, label1).
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

func TestIf33(t *testing.T) {
	label1 := NewLabel("a")
	label2 := NewLabel("b")
	code := newBuilder().
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
	code := newBuilder().
		Push(nil).
		JmpIf(exec.JcNotNil, label1).
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

func TestIf44(t *testing.T) {
	label1 := NewLabel("a")
	label2 := NewLabel("b")
	code := newBuilder().
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
	code := newBuilder().
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
	code := newBuilder().
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
	code := newBuilder().
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
