/*
 * Copyright (c) 2023 The GoPlus Authors (goplus.org). All rights reserved.
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

package c2go

import (
	"path/filepath"

	"github.com/goplus/mod/gopmod"
)

// -----------------------------------------------------------------------------

// LookupPub returns a anonymous function required by cl.NewPackage.
func LookupPub(mod *gopmod.Module) func(pkgPath string) (pubfile string, err error) {
	return func(pkgPath string) (pubfile string, err error) {
		pkg, err := mod.Lookup(pkgPath)
		if err == nil {
			pubfile = filepath.Join(pkg.Dir, "c2go.a.pub")
		}
		return
	}
}

// -----------------------------------------------------------------------------
