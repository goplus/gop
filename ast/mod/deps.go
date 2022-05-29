/*
 * Copyright (c) 2022 The GoPlus Authors (goplus.org). All rights reserved.
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

package mod

import (
	"strconv"
	"strings"

	goast "go/ast"
	gotoken "go/token"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/token"
)

// ----------------------------------------------------------------------------

type Deps struct {
	HandlePkg func(pkgPath string)
}

func (p Deps) Load(pkg *ast.Package, withGopStd bool) {
	for _, f := range pkg.Files {
		p.LoadFile(f, withGopStd)
	}
	for _, f := range pkg.GoFiles {
		p.LoadGoFile(f)
	}
}

func (p Deps) LoadGoFile(f *goast.File) {
	for _, imp := range f.Imports {
		path := imp.Path
		if path.Kind == gotoken.STRING {
			if s, err := strconv.Unquote(path.Value); err == nil {
				if s == "C" {
					continue
				}
				p.HandlePkg(s)
			}
		}
	}
}

func (p Deps) LoadFile(f *ast.File, withGopStd bool) {
	for _, imp := range f.Imports {
		path := imp.Path
		if path.Kind == token.STRING {
			if s, err := strconv.Unquote(path.Value); err == nil {
				p.gopPkgPath(s, withGopStd)
			}
		}
	}
}

func (p Deps) gopPkgPath(s string, withGopStd bool) {
	if strings.HasPrefix(s, "gop/") {
		if !withGopStd {
			return
		}
		s = "github.com/goplus/gop/" + s[4:]
	} else if strings.HasPrefix(s, "C") {
		if len(s) == 1 {
			s = "github.com/goplus/libc"
		} else if s[1] == '/' {
			s = s[2:]
			if strings.IndexByte(s, '/') < 0 {
				s = "github.com/goplus/" + s
			}
		}
	}
	p.HandlePkg(s)
}

// ----------------------------------------------------------------------------
