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

package gopmod

import (
	"os"
	"path/filepath"
)

// -----------------------------------------------------------------------------

const (
	FlagGoAsGoPlus = 1 << iota
)

func (p *Context) OpenProject(flags int, args ...string) (proj *Project, err error) {
	if len(args) != 1 {
		return p.openFromGopFiles(args)
	}
	src := args[0]
	fi, err := os.Stat(src)
	if err != nil {
		return
	}
	if fi.IsDir() {
		return p.openFromDir(src)
	}
	if (flags&FlagGoAsGoPlus) == 0 && filepath.Ext(src) == ".go" {
		return openFromGoFile(src)
	}
	return p.openFromGopFiles(args)
}

// -----------------------------------------------------------------------------

func (p *Context) openFromDir(dir string) (proj *Project, err error) {
	panic("todo")
}

// -----------------------------------------------------------------------------
