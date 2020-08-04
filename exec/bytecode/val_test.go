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

func TestVals(t *testing.T) {
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
			b.StoreVal(reflect.ValueOf(p).Interface())
			b.StoreVar(v)
			b.LoadVar(v)
			b.CallGoFuncv(println, 1, 1)
			b.Push("foo")
			b.LoadVar(v)
			b.StoreVal("Name")
			b.SetField()
			b.StoreVal(v)
			b.LoadVar(v)
			b.CallGoFuncv(println, 1, 1)
			b.LoadVar(v)
			b.StoreVal("Name")
			b.CallField()
			b.CallGoFuncv(println, 1, 1)
			code := b.Resolve()
			code.(*Code).Dump(os.Stderr)
			ctx := NewContext(code)
			ctx.Exec(0, code.Len())
		},
		"&{bar 30}\n&{foo 30}\nfoo\n",
	)
}

func TestVals2(t *testing.T) {
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
			b.StoreVal(reflect.ValueOf(p).Interface())
			b.StoreVar(v)
			b.LoadVar(v)
			b.CallGoFuncv(println, 1, 1)
			b.Push("foo")
			b.LoadVar(v)
			b.StoreVal("Name")
			b.SetField()
			b.StoreVal(v)
			b.LoadVar(v)
			b.CallGoFuncv(println, 1, 1)
			b.LoadVar(v)
			b.StoreVal("Name")
			b.CallField()
			b.CallGoFuncv(println, 1, 1)
			code := b.Resolve()
			code.(*Code).Dump(os.Stderr)
			ctx := NewContext(code)
			ctx.Exec(0, code.Len())
		},
		"{bar 30}\n{bar 30}\nbar\n",
	)
}
