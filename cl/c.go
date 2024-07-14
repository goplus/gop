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
)

// -----------------------------------------------------------------------------

const (
	pathLibc   = "github.com/goplus/llgo/c"
	pathLibpy  = "github.com/goplus/llgo/py"
	pathLibcpp = "github.com/goplus/llgo/cpp"
)

func simplifyGopPackage(pkgPath string) string {
	if strings.HasPrefix(pkgPath, "gop/") {
		return "github.com/goplus/" + pkgPath
	}
	return pkgPath
}

func simplifyPkgPath(pkgPath string) string {
	switch pkgPath {
	case "c":
		return pathLibc
	case "py":
		return pathLibpy
	default:
		if strings.HasPrefix(pkgPath, "c/") {
			return pathLibc + pkgPath[1:]
		}
		if strings.HasPrefix(pkgPath, "py/") {
			return pathLibpy + pkgPath[2:]
		}
		if strings.HasPrefix(pkgPath, "cpp/") {
			return pathLibcpp + pkgPath[3:]
		}
		return simplifyGopPackage(pkgPath)
	}
}

// -----------------------------------------------------------------------------
