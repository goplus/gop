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

package bytecode

import (
	"testing"
)

// -----------------------------------------------------------------------------

func TestConst1(t *testing.T) {
	code := newBuilder().
		Push(int64(1 << 32)).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != int64(1<<32) {
		t.Fatal("1<<32 != 1<<32, ret =", v)
	}
}

func TestConst2(t *testing.T) {
	code := newBuilder().
		Push(uint64(1 << 32)).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != uint64(1<<32) {
		t.Fatal("1<<32 != 1<<32, ret =", v)
	}
}

func TestConst3(t *testing.T) {
	code := newBuilder().
		Push(uint32(1 << 30)).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != uint32(1<<30) {
		t.Fatal("1<<30 != 1<<30, ret =", v)
	}
}

func TestConst4(t *testing.T) {
	code := newBuilder().
		Push(int32(1 << 30)).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != int32(1<<30) {
		t.Fatal("1<<30 != 1<<30, ret =", v)
	}
}

func TestConst5(t *testing.T) {
	code := newBuilder().
		Push(uint(1 << 12)).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != uint(1<<12) {
		t.Fatal("1<<12 != 1<<12, ret =", v)
	}
}

func TestConst6(t *testing.T) {
	code := newBuilder().
		Push(1.12).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := ctx.Pop(); v != 1.12 {
		t.Fatal("1.12 != 1.12, ret =", v)
	}
}

func TestNil(t *testing.T) {
	code := newBuilder().
		Push(nil).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := ctx.Pop(); v != nil {
		t.Fatal("nil != nil, ret =", v)
	}
}

func TestReserve(t *testing.T) {
	b := NewBuilder(nil)
	off := b.Reserve()
	b.ReservedAsPush(off, 1.12)
	code := b.Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != 1.12 {
		t.Fatal("1.12 != 1.12, ret =", v)
	}
	i := b.Interface()
	_ = i.EndStmt(nil, nil).StartStmt(nil)
	g := i.Pop(0).(*iBuilder).GetPackage()
	_ = g.GetGoFuncInfo(0)
	_ = g.GetGoFuncvInfo(0)
}

func TestReserve2(t *testing.T) {
	b := newBuilder()
	off := b.Reserve()
	b.ReservedAsPush(off, 1.12)
	code := b.Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != 1.12 {
		t.Fatal("1.12 != 1.12, ret =", v)
	}
}

// -----------------------------------------------------------------------------
