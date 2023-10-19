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
	goast "go/ast"
	"go/types"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/token"
)

// A Checker maintains the state of the type checker.
// It must be created with NewChecker.
type Checker struct {
	c *types.Checker
}

// NewChecker returns a new Checker instance for a given package.
// Package files may be added incrementally via checker.Files.
func NewChecker(
	conf *types.Config, fset *token.FileSet, pkg *types.Package,
	goInfo *types.Info, gopInfo *Info) *Checker {
	check := types.NewChecker(conf, fset, pkg, goInfo)
	return &Checker{check}
}

// Files checks the provided files as part of the checker's package.
func (p *Checker) Files(goFiles []*goast.File, gopFiles []*ast.File) error {
	if len(gopFiles) == 0 {
		return p.c.Files(goFiles)
	}
	// goxls: todo
	return nil
}
