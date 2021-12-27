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
	"github.com/goplus/gop/x/mod/modload"
)

// -----------------------------------------------------------------------------

func UpdateGoMod(dir string) {
	p, err := modload.Load(dir)
	if err != nil {
		return
	}
	p.UpdateGoMod(true)
	/*
		if p.classModFile != nil && p.classModFile.Classfile != nil {
			gengo.RegisterPkgFlags(p.classModFile.Classfile.ProjExt, gengo.PkgFlagGmx)
			gengo.RegisterPkgFlags(p.classModFile.Classfile.WorkExt, gengo.PkgFlagSpx)
			cl.RegisterClassFileType(p.classModFile.Classfile.ProjExt,
				p.classModFile.Classfile.WorkExt, p.classModFile.Classfile.PkgPaths...)
		}
	*/
}

// -----------------------------------------------------------------------------
