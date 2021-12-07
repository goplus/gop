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

package format

import (
	"github.com/goplus/gop/ast"
)

// -----------------------------------------------------------------------------

var printFuncs = map[string]string{
	"Errorf":   "errorf",
	"Print":    "print",
	"Printf":   "printf",
	"Println":  "println",
	"Fprint":   "fprint",
	"Fprintf":  "fprintf",
	"Fprintln": "fprintln",
	"Sprint":   "sprint",
	"Sprintf":  "sprintf",
	"Sprintln": "sprintln",
}

func fmtToBuiltin(ctx *importCtx, sel *ast.Ident, ref *ast.Expr) bool {
	if ctx.pkgPath == "fmt" {
		if v, ok := printFuncs[sel.Name]; ok {
			*ref = &ast.Ident{NamePos: sel.NamePos, Name: v}
			return true
		}
	}
	return false
}

func commandStyleFirst(v *ast.CallExpr) {
	v.NoParenEnd = v.Rparen
}

// -----------------------------------------------------------------------------
