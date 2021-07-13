/*
 Copyright 2021 The GoPlus Authors (goplus.org)

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
	"go/token"
	"go/types"

	"github.com/goplus/gox"
)

// -----------------------------------------------------------------------------

func initBuiltin(pkg gox.PkgImporter, builtin *types.Package) {
	fmt := pkg.Import("fmt")
	fns := []string{"print", "println", "printf", "errorf", "fprint", "fprintln", "fprintf"}
	for _, fn := range fns {
		fnTitle := string(fn[0]-'a'+'A') + fn[1:]
		builtin.Scope().Insert(gox.NewOverloadFunc(token.NoPos, builtin, fn, fmt.Ref(fnTitle)))
	}
}

func newBuiltinDefault(pkg gox.PkgImporter, prefix *gox.NamePrefix, contracts *gox.BuiltinContracts) *types.Package {
	builtin := types.NewPackage("", "")
	initBuiltin(pkg, builtin)
	gox.InitBuiltinOps(builtin, prefix, contracts)
	gox.InitBuiltinFuncs(builtin)
	return builtin
}

// -----------------------------------------------------------------------------
