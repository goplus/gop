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

package cl

import (
	"testing"

	exec "github.com/goplus/gop/exec.spec"
)

// -----------------------------------------------------------------------------

func TestValue(t *testing.T) {
	g := &goValue{t: exec.TyInt}
	_ = g.Value(0)

	nv := &nonValue{0}
	_ = nv.Kind()
	_ = nv.NumValues()
	_ = nv.Type()
	_ = nv.Value(0)

	f := new(qlFunc)
	_ = f.Kind()
	_ = f.Value(0)

	gf := new(goFunc)
	_ = gf.Kind()
	_ = gf.NumValues()
	_ = gf.Type()
	_ = gf.Value(0)

	c := new(constVal)
	_ = c.Value(0)
}

func TestOverflowsIntBy(t *testing.T) {
	if isOverflowsIntByInt64(1<<63-1, 64) {
		t.Fatal("not overflows", 1<<63-1)
	}
	if !isOverflowsIntByInt64(1<<31, 32) {
		t.Fatal("overflows", 1<<31)
	}
	if isOverflowsIntByUint64(1<<63-1, 64) {
		t.Fatal("not overflows", uint64(1<<63-1))
	}
	if !isOverflowsIntByUint64(1<<64-1, 64) {
		t.Fatal("not overflows", uint64(1<<64-1))
	}
	if !isOverflowsIntByUint64(1<<32, 32) {
		t.Fatal("overflows", 1<<32)
	}
}

// -----------------------------------------------------------------------------
