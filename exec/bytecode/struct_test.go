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
	"os"
	"reflect"
	"testing"
)

// -----------------------------------------------------------------------------

type Person struct {
	Name string
	Age  int
}

func TestStruct(t *testing.T) {
	println, ok := I.FindFuncv("Println")
	if !ok {
		t.Fatal("FindFuncv failed: Println")
	}

	b := newBuilder()
	p := Person{
		Name: "bar",
		Age:  30,
	}

	v := NewVar(reflect.TypeOf(p), "p")
	expect(t,
		func() {
			b.DefineVar(v)
			b.Push(0)
			b.Push("bar")
			b.Push(1)
			b.Push(30)
			b.Struct(reflect.TypeOf(p), 2)
			b.StoreVar(v)
			b.LoadVar(v)
			b.CallGoFuncv(println, 1, 1)
			code := b.Resolve()
			code.(*Code).Dump(os.Stderr)
			ctx := NewContext(code)
			ctx.Exec(0, code.Len())
		},
		"{bar 30}\n",
	)
}

func TestStruct2(t *testing.T) {
	println, ok := I.FindFuncv("Println")
	if !ok {
		t.Fatal("FindFuncv failed: Println")
	}

	b := newBuilder()
	p := &Person{
		Name: "bar",
		Age:  30,
	}

	v := NewVar(reflect.TypeOf(p), "p")
	expect(t,
		func() {
			b.DefineVar(v)
			b.Push(0)
			b.Push("bar")
			b.Push(1)
			b.Push(30)
			b.Struct(reflect.TypeOf(p), 2)
			b.StoreVar(v)
			b.LoadVar(v)
			b.CallGoFuncv(println, 1, 1)
			code := b.Resolve()
			code.(*Code).Dump(os.Stderr)
			ctx := NewContext(code)
			ctx.Exec(0, code.Len())
		},
		"&{bar 30}\n",
	)
}

func TestCloneStruct(t *testing.T) {
	type obj struct {
		a int
	}
	v1 := reflect.ValueOf(obj{a: 10})
	v2 := cloneStruct(v1)

	t.Logf("v1: %+v , v2: %+v", v1.Interface(), v2.Interface())

	if !reflect.DeepEqual(v1.Interface(), v2.Interface()) {
		t.Fatalf("v1: %+v , v2: %+v", v1.Interface(), v2.Interface())
	}

	if !reflect.DeepEqual(v1.Interface(), InterfaceOf(v1)) {
		t.Fatalf("v1 interface:  %+v  |  %+v", v1.Interface(), InterfaceOf(v1))
	}

	if !reflect.DeepEqual(v1.Interface(), InterfaceOf(v1.Interface())) {
		t.Fatalf("v1 interface:  %+v  |  %+v", v1.Interface(), InterfaceOf(v1.Interface()))
	}

	if !reflect.DeepEqual(v1, valueOf(v1)) {
		t.Fatalf("valueOf(v1) :  %+v  |  %+v", v1, valueOf(v1))
	}

	if !reflect.DeepEqual(v1, valueOf(InterfaceOf(v1))) {
		t.Fatalf("valueOf(v1) :  %+v  |  %+v", v1, valueOf(InterfaceOf(v1)))
	}
}
