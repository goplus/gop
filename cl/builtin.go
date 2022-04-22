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
	conf.UntypedBigInt = big.Ref("UntypedInt").Type().(*types.Named)
	conf.UntypedBigRat = big.Ref("UntypedRat").Type().(*types.Named)
	conf.UntypedBigFloat = big.Ref("UntypedFloat").Type().(*types.Named)
}

func initBuiltin(pkg gox.PkgImporter, builtin *types.Package, fmt, big, ng, buil *gox.PkgRef) {
	scope := builtin.Scope()
	typs := []string{"int", "rat", "float"}
	for _, typ := range typs {
		name := string(typ[0]-('a'-'A')) + typ[1:]
		scope.Insert(types.NewTypeName(token.NoPos, builtin, "big"+typ, big.Ref(name).Type()))
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
	scope.Insert(gox.NewOverloadFunc(token.NoPos, builtin, "newRange", buil.Ref("NewRange__0")))
	scope.Insert(types.NewTypeName(token.NoPos, builtin, "uint128", ng.Ref("Uint128").Type()))
	scope.Insert(types.NewTypeName(token.NoPos, builtin, "int128", ng.Ref("Int128").Type()))
	scope.Insert(types.NewTypeName(token.NoPos, builtin, "any", gox.TyEmptyInterface))
}

func newBuiltinDefault(pkg gox.PkgImporter, conf *gox.Config) *types.Package {
	builtin := types.NewPackage("", "")
	fmt := pkg.Import("fmt")
	buil := pkg.Import("github.com/goplus/gop/builtin")
	big := pkg.Import("github.com/goplus/gop/builtin/big")
	ng := pkg.Import("github.com/goplus/gop/builtin/ng")
	pkg.Import("strconv")
	pkg.Import("strings")
	initMathBig(pkg, conf, big)
	initBuiltin(pkg, builtin, fmt, big, ng, buil)
	gox.InitBuiltin(pkg, builtin, conf)
	return builtin
}

// -----------------------------------------------------------------------------
