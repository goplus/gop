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

// Package strings provide Go+ "strings" package, as "strings" package in Go.
package strings

import (
	"strings"

	"github.com/qiniu/goplus/gop"
)

// -----------------------------------------------------------------------------

func execNewReplacer(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	repl := strings.NewReplacer(gop.ToStrings(args)...)
	p.Ret(arity, repl)
}

func execReplacerReplace(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret := args[0].(*strings.Replacer).Replace(args[1].(string))
	p.Ret(2, ret)
}

// -----------------------------------------------------------------------------

// I is a Go package instance.
var I = gop.NewGoPackage("strings")

func init() {
	I.RegisterFuncvs(
		I.Funcv("NewReplacer", strings.NewReplacer, execNewReplacer),
	)
	I.RegisterFuncs(
		I.Func("(*Replacer).Replace", (*strings.Replacer).Replace, execReplacerReplace),
	)
}

// -----------------------------------------------------------------------------
