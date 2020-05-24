/*
 Copyright 2020 Qiniu Cloud (七牛云)

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

package fmt

import (
	"reflect"

	qlang "github.com/qiniu/qlang/v6/spec"
)

func execTypeOf(zero int, p *qlang.Context) {
	args := p.GetArgs(1)
	args[0] = reflect.TypeOf(args[0])
}

func execValueOf(zero int, p *qlang.Context) {
	args := p.GetArgs(1)
	args[0] = reflect.ValueOf(args[0])
}

// -----------------------------------------------------------------------------

// I is a Go package instance.
var I = qlang.NewGoPackage("reflect")

func init() {
	I.RegisterFuncs(
		I.Func("TypeOf", reflect.TypeOf, execTypeOf),
		I.Func("ValueOf", reflect.ValueOf, execValueOf),
	)
}

// -----------------------------------------------------------------------------
