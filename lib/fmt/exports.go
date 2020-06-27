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

// Package fmt provide Go+ "fmt" package, as "fmt" package in Go.
package fmt

import (
	"fmt"

	"github.com/qiniu/goplus/gop"
	"github.com/qiniu/goplus/lib/builtin"
)

// -----------------------------------------------------------------------------

// I is a Go package instance.
var I = gop.NewGoPackage("fmt")

func init() {
	I.RegisterFuncvs(
		I.Funcv("errorf", fmt.Errorf, builtin.QexecErrorf),
		I.Funcv("Print", fmt.Print, builtin.QexecPrint),
		I.Funcv("Printf", fmt.Printf, builtin.QexecPrintf),
		I.Funcv("Println", fmt.Println, builtin.QexecPrintln),
		I.Funcv("Fprintln", fmt.Fprintln, builtin.QexecFprintln),
	)
}

// -----------------------------------------------------------------------------
