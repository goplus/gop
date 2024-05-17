/*
 * Copyright (c) 2022 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cl_test

import (
	"testing"
)

func TestHelloLLGo(t *testing.T) {
	gopClTest(t, `
import "c"

c.printf C"Hello, world!\n"
`, `package main

import "github.com/goplus/llgo/c"

func main() {
	c.Printf(c.Str("Hello, world!\n"))
}
`)
}

func TestPyCall(t *testing.T) {
	gopClTest(t, `
import (
	"c"
	"py"
	"py/math"
)

x := math.sqrt(py.float(2))
c.printf C"sqrt(2) = %f\n", x.float64
`, `package main

import (
	"github.com/goplus/llgo/c"
	"github.com/goplus/llgo/py"
	"github.com/goplus/llgo/py/math"
)

func main() {
	x := math.Sqrt(py.Float(2))
	c.Printf(c.Str("sqrt(2) = %f\n"), x.Float64())
}
`)
}
