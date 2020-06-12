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
	"strconv"

	"github.com/qiniu/goplus/gop"
)

func execAtoi(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	v, err := strconv.Atoi(args[0].(string))
	p.Ret(1, v, err)
}

// -----------------------------------------------------------------------------

// I is a Go package instance.
var I = gop.NewGoPackage("strconv")

func init() {
	I.RegisterFuncs(
		I.Func("Atoi", strconv.Atoi, execAtoi),
	)
}

// -----------------------------------------------------------------------------
