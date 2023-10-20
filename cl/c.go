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

package cl

import (
	"strings"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gox"
)

const (
	c2goInvalid = iota
	c2goStandard
	c2goUserDef
)

func checkC2go(pkgPath string) (realPkgPath string, kind int) {
	if strings.HasPrefix(pkgPath, "C") {
		if len(pkgPath) == 1 {
			return "libc", c2goStandard
		}
		if pkgPath[1] == '/' {
			realPkgPath = pkgPath[2:]
			if strings.IndexByte(realPkgPath, '/') < 0 {
				kind = c2goStandard
			} else {
				kind = c2goUserDef
			}
		}
	}
	return
}

func c2goBase(base string) string {
	if base == "" {
		base = "github.com/goplus/"
	} else if !strings.HasSuffix(base, "/") {
		base += "/"
	}
	return base
}

// -----------------------------------------------------------------------------

func loadC2goPkg(ctx *blockCtx, realPath string, src *ast.BasicLit) *gox.PkgRef {
	cpkg, err := ctx.cpkgs.Import(realPath)
	if err != nil {
		ctx.handleErrorf(src.Pos(),
			"%v not found or not a valid C package (c2go.a.pub file not found).\n", realPath)
		return nil
	}
	ctx.clookups = append(ctx.clookups, cpkg)
	return cpkg.Pkg()
}

// -----------------------------------------------------------------------------
