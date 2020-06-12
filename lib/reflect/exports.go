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

// Package fmt provide Go+ "reflect" package, as "reflect" package in Go.
package fmt

import (
	"reflect"

	"github.com/qiniu/goplus/gop"
)

func execTypeOf(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0] = reflect.TypeOf(args[0])
}

func execValueOf(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0] = reflect.ValueOf(args[0])
}

// -----------------------------------------------------------------------------

// I is a Go package instance.
var I = gop.NewGoPackage("reflect")

func init() {
	I.RegisterFuncs(
		I.Func("TypeOf", reflect.TypeOf, execTypeOf),
		I.Func("ValueOf", reflect.ValueOf, execValueOf),
	)
}

// -----------------------------------------------------------------------------
