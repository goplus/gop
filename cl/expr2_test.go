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

package cl_test

import (
	"testing"

	"github.com/goplus/gop/cl/cltest"
)

// -----------------------------------------------------------------------------

func TestUnbound(t *testing.T) {
	cltest.Expect(t,
		`println("Hello " + "qiniu:", 123, 4.5, 7i)`,
		"Hello qiniu: 123 4.5 (0+7i)\n",
	)
}

func TestPanic(t *testing.T) {
	cltest.Expect(t,
		`panic("Helo")`, "", "Helo",
	)
}

func TestTypeCast(t *testing.T) {
	cltest.Call(t, `
	x := []byte("hello")
	x
	`).Equal([]byte("hello"))
}

// -----------------------------------------------------------------------------
