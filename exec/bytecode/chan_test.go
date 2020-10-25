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

func TestChan(t *testing.T) {
	println, ok := I.FindFuncv("Println")
	if !ok {
		t.Fatal("FindFuncv failed: Println")
	}

	b := newBuilder()

	v := NewVar(reflect.ChanOf(3, reflect.TypeOf(0)), "p")
	i := NewVar(reflect.TypeOf(0), "i")
	expect(t,
		func() {
			b.DefineVar(v)
			b.DefineVar(i)
			b.Push(10)
			b.Make(reflect.ChanOf(3, reflect.TypeOf(0)), 1)
			b.StoreVar(v)
			b.LoadVar(v)
			b.Push(3)
			b.Send()
			b.LoadVar(v)
			b.Recv()
			b.StoreVar(i)
			b.LoadVar(i)
			b.CallGoFuncv(println, 1, 1)
			code := b.Resolve()
			code.(*Code).Dump(os.Stderr)
			ctx := NewContext(code)
			ctx.Exec(0, code.Len())
		},
		"3\n",
	)
}
