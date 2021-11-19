/*
 * Copyright (c) 2021 The GoPlus Authors (goplus.org). All rights reserved.
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

package cl

import (
	"go/token"
	"go/types"

	"github.com/goplus/gox"
)

// -----------------------------------------------------------------------------

func initMathBig(pkg gox.PkgImporter, conf *gox.Config, big *gox.PkgRef) {
	big.EnsureImported()
	conf.UntypedBigInt = big.Ref("Gop_untyped_bigint").Type().(*types.Named)
	conf.UntypedBigRat = big.Ref("Gop_untyped_bigrat").Type().(*types.Named)
	conf.UntypedBigFloat = big.Ref("Gop_untyped_bigfloat").Type().(*types.Named)
}

func initBuiltin(pkg gox.PkgImporter, builtin *types.Package, fmt, big *gox.PkgRef) {
	scope := builtin.Scope()
	typs := []string{"bigint", "bigrat", "bigfloat"}
	for _, typ := range typs {
		scope.Insert(types.NewTypeName(token.NoPos, builtin, typ, big.Ref("Gop_"+typ).Type()))
	}
	fns := []string{
		"print", "println", "printf", "errorf",
		"fprint", "fprintln", "fprintf",
		"sprint", "sprintln", "sprintf",
	}
	for _, fn := range fns {
		fnTitle := string(fn[0]-'a'+'A') + fn[1:]
		scope.Insert(gox.NewOverloadFunc(token.NoPos, builtin, fn, fmt.Ref(fnTitle)))
	}
	scope.Insert(gox.NewOverloadFunc(token.NoPos, builtin, "newRange", big.Ref("NewRange__0")))
}

func newBuiltinDefault(pkg gox.PkgImporter, conf *gox.Config) *types.Package {
	builtin := types.NewPackage("", "")
	fmt := pkg.Import("fmt")
	big := pkg.Import("github.com/goplus/gop/builtin")
	pkg.Import("strconv")
	pkg.Import("strings")
	initMathBig(pkg, conf, big)
	initBuiltin(pkg, builtin, fmt, big)
	gox.InitBuiltin(pkg, builtin, conf)
	return builtin
}

// -----------------------------------------------------------------------------
