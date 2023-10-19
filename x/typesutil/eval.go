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

package typesutil

import (
	"go/types"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/token"
)

// CheckExpr type checks the expression expr as if it had appeared at position
// pos of package pkg. Type information about the expression is recorded in
// info. The expression may be an identifier denoting an uninstantiated generic
// function or type.
//
// If pkg == nil, the Universe scope is used and the provided
// position pos is ignored. If pkg != nil, and pos is invalid,
// the package scope is used. Otherwise, pos must belong to the
// package.
//
// An error is returned if pos is not within the package or
// if the node cannot be type-checked.
//
// Note: Eval and CheckExpr should not be used instead of running Check
// to compute types and values, but in addition to Check, as these
// functions ignore the context in which an expression is used (e.g., an
// assignment). Thus, top-level untyped constants will return an
// untyped type rather than the respective context-specific type.
func CheckExpr(fset *token.FileSet, pkg *types.Package, pos token.Pos, expr ast.Expr, info *Info) (err error) {
	panic("todo")
}
