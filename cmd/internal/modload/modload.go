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

package modload

import (
	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/cmd/gengo"
	"github.com/goplus/gop/x/gopmod"
)

// -----------------------------------------------------------------------------

func UpdateGoMod(dir string) {
	p, err := gopmod.Load(dir)
	if err != nil {
		return
	}
	p.UpdateGoMod(true)
	p.RegisterClasses(func(c *gopmod.Class) {
		gengo.RegisterPkgFlags(c.ProjExt, gengo.PkgFlagGmx)
		gengo.RegisterPkgFlags(c.WorkExt, gengo.PkgFlagSpx)
		if c != gopmod.ClassSpx {
			cl.RegisterClassFileType(c.ProjExt, c.WorkExt, c.PkgPaths...)
		}
	})
}

// -----------------------------------------------------------------------------
