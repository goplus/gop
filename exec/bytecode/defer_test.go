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

func TestDefer1(t *testing.T) {
	b := newBuilder()
	off := b.Reserve()
	b.ReservedAsInstr(off, &iDefer{start: 1, end: 3})
	code := b.Resolve()

	ctx := NewContext(code)

	if i := ctx.code.data[0] >> bitsOpShift; i != opDeferOp {
		t.Fatal("opDeferOp != opDeferOp, ret =", i)
	}
}

func TestDefer2(t *testing.T) {
	b := newBuilder()
	l1 := NewLabel("")
	l2 := NewLabel("")
	off := b.Reserve()
	b.Label(l1)
	b.Push(1)
	b.Push(2)
	b.Push(3)
	b.Label(l2)
	b.ReservedAsInstr(off, b.Defer(l1, l2))
	code := b.Resolve()
	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	var i = 3
	for {
		if ctx.Len() < 1 {
			break
		}
		v := ctx.Pop()
		if v != i {
			t.Fatalf("v != %d, ret = %v", i, v)
		}
		i--
	}
}
