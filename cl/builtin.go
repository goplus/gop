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

	"github.com/goplus/gogen"
)

// -----------------------------------------------------------------------------

func initMathBig(_ *gogen.Package, conf *gogen.Config, big gogen.PkgRef) {
	conf.UntypedBigInt = big.Ref("UntypedBigint").Type().(*types.Named)
	conf.UntypedBigRat = big.Ref("UntypedBigrat").Type().(*types.Named)
	conf.UntypedBigFloat = big.Ref("UntypedBigfloat").Type().(*types.Named)
}

func initBuiltinFns(builtin *types.Package, scope *types.Scope, pkg gogen.PkgRef, fns []string) {
	for _, fn := range fns {
		fnTitle := string(fn[0]-'a'+'A') + fn[1:]
		scope.Insert(gogen.NewOverloadFunc(token.NoPos, builtin, fn, pkg.Ref(fnTitle)))
	}
}

func initBuiltin(_ *gogen.Package, builtin *types.Package, os, fmt, ng, iox, buil, reflect gogen.PkgRef) {
	scope := builtin.Scope()
	if ng.Types != nil {
		typs := []string{"bigint", "bigrat", "bigfloat"}
		for _, typ := range typs {
			name := string(typ[0]-('a'-'A')) + typ[1:]
			scope.Insert(types.NewTypeName(token.NoPos, builtin, typ, ng.Ref(name).Type()))
		}
		scope.Insert(types.NewTypeName(token.NoPos, builtin, "uint128", ng.Ref("Uint128").Type()))
		scope.Insert(types.NewTypeName(token.NoPos, builtin, "int128", ng.Ref("Int128").Type()))
	}
	if fmt.Types != nil {
		scope.Insert(gogen.NewOverloadFunc(token.NoPos, builtin, "echo", fmt.Ref("Println")))
		initBuiltinFns(builtin, scope, fmt, []string{
			"print", "println", "printf", "errorf",
			"fprint", "fprintln", "fprintf",
			"sprint", "sprintln", "sprintf",
		})
	}
	if os.Types != nil {
		initBuiltinFns(builtin, scope, os, []string{
			"open", "create",
		})
	}
	if iox.Types != nil {
		initBuiltinFns(builtin, scope, iox, []string{
			"lines",
		})
		scope.Insert(gogen.NewOverloadFunc(token.NoPos, builtin, "blines", iox.Ref("BLines")))
	}
	if reflect.Types != nil {
		scope.Insert(gogen.NewOverloadFunc(token.NoPos, builtin, "type", reflect.Ref("TypeOf")))
	}
	if buil.Types != nil {
		scope.Insert(gogen.NewOverloadFunc(token.NoPos, builtin, "newRange", buil.Ref("NewRange__0")))
	}
	scope.Insert(types.NewTypeName(token.NoPos, builtin, "any", gogen.TyEmptyInterface))
}

func newBuiltinDefault(pkg *gogen.Package, conf *gogen.Config) *types.Package {
	builtin := types.NewPackage("", "")
	fmt := pkg.TryImport("fmt")
	os := pkg.TryImport("os")
	reflect := pkg.TryImport("reflect")
	buil := pkg.TryImport("github.com/goplus/gop/builtin")
	ng := pkg.TryImport("github.com/goplus/gop/builtin/ng")
	iox := pkg.TryImport("github.com/goplus/gop/builtin/iox")
	strx := pkg.TryImport("github.com/qiniu/x/stringutil")
	stringslice := pkg.TryImport("github.com/goplus/gop/builtin/stringslice")
	pkg.TryImport("strconv")
	pkg.TryImport("strings")
	if ng.Types != nil {
		initMathBig(pkg, conf, ng)
	}
	initBuiltin(pkg, builtin, os, fmt, ng, iox, buil, reflect)
	gogen.InitBuiltin(pkg, builtin, conf)
	if strx.Types != nil {
		ti := pkg.BuiltinTI(types.Typ[types.String])
		ti.AddMethods(
			&gogen.BuiltinMethod{Name: "Capitalize", Fn: strx.Ref("Capitalize")},
		)
	}
	if stringslice.Types != nil {
		ti := pkg.BuiltinTI(types.NewSlice(types.Typ[types.String]))
		ti.AddMethods(
			&gogen.BuiltinMethod{Name: "Capitalize", Fn: stringslice.Ref("Capitalize")},
			&gogen.BuiltinMethod{Name: "ToTitle", Fn: stringslice.Ref("ToTitle")},
			&gogen.BuiltinMethod{Name: "ToUpper", Fn: stringslice.Ref("ToUpper")},
			&gogen.BuiltinMethod{Name: "ToLower", Fn: stringslice.Ref("ToLower")},
			&gogen.BuiltinMethod{Name: "Repeat", Fn: stringslice.Ref("Repeat")},
			&gogen.BuiltinMethod{Name: "Replace", Fn: stringslice.Ref("Replace")},
			&gogen.BuiltinMethod{Name: "ReplaceAll", Fn: stringslice.Ref("ReplaceAll")},
			&gogen.BuiltinMethod{Name: "Trim", Fn: stringslice.Ref("Trim")},
			&gogen.BuiltinMethod{Name: "TrimSpace", Fn: stringslice.Ref("TrimSpace")},
			&gogen.BuiltinMethod{Name: "TrimLeft", Fn: stringslice.Ref("TrimLeft")},
			&gogen.BuiltinMethod{Name: "TrimRight", Fn: stringslice.Ref("TrimRight")},
			&gogen.BuiltinMethod{Name: "TrimPrefix", Fn: stringslice.Ref("TrimPrefix")},
			&gogen.BuiltinMethod{Name: "TrimSuffix", Fn: stringslice.Ref("TrimSuffix")},
		)
	}
	return builtin
}

// -----------------------------------------------------------------------------
