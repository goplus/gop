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

package cl

import (
	"fmt"
	"testing"

	"github.com/qiniu/qlang/v6/ast"
	"github.com/qiniu/qlang/v6/ast/asttest"
	"github.com/qiniu/qlang/v6/exec"
	"github.com/qiniu/qlang/v6/parser"
	"github.com/qiniu/qlang/v6/token"
	"github.com/qiniu/x/log"

	_ "github.com/qiniu/qlang/v6/lib/builtin"
	_ "github.com/qiniu/qlang/v6/lib/fmt"
	_ "github.com/qiniu/qlang/v6/lib/reflect"
	_ "github.com/qiniu/qlang/v6/lib/strings"
)

func init() {
	log.SetFlags(log.Ldefault &^ log.LstdFlags)
	log.SetOutputLevel(log.Ldebug)
}

func newPackage(out *exec.Builder, pkg *ast.Package) (p *Package, noExecCtx bool, err error) {
	if p, err = NewPackage(out, pkg); err != nil {
		return
	}
	ctxPkg := newPkgCtx(out)
	ctx := newGblBlockCtx(ctxPkg, nil)
	ctx.syms = p.syms
	entry, _ := ctx.findFunc("main")
	noExecCtx = isNoExecCtx(ctx, entry.body)
	return
}

// -----------------------------------------------------------------------------

var fsTestBasic = asttest.NewSingleFileFS("/foo", "bar.ql", `
	println("Hello", "xsw", "- nice to meet you!")
	println("Hello, world!")
	return
`)

func TestBasic(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseFSDir(fset, fsTestBasic, "/foo", nil, 0)
	if err != nil || len(pkgs) != 1 {
		t.Fatal("ParseFSDir failed:", err, len(pkgs))
	}

	bar := pkgs["main"]
	b := exec.NewBuilder(nil)
	_, noExecCtx, err := newPackage(b, bar)
	if err != nil || !noExecCtx {
		t.Fatal("Compile failed:", err)
	}
	code := b.Resolve()

	ctx := exec.NewContext(code)
	ctx.Exec(0, code.Len())
	fmt.Println("results:", ctx.Get(-2), ctx.Get(-1))
	if v := ctx.Get(-1); v != nil {
		t.Fatal("error:", v)
	}
	if v := ctx.Get(-2); v != int(14) {
		t.Fatal("n:", v)
	}
}

// -----------------------------------------------------------------------------
