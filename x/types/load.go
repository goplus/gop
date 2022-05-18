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

// This package is experimental. Don't use it in production.
package types

import (
	"go/token"
	"go/types"
	"strings"

	goast "go/ast"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/ast/togo"
)

// A Config specifies the configuration for type checking.
// The zero value for Config is a ready-to-use default configuration.
type Config = types.Config

// A Package describes a Go/Go+ package.
type Package = types.Package

// Load loads a package and returns the resulting package object and
// the first error if any.
//
// The package is marked as complete if no errors occurred, otherwise it is
// incomplete. See Config.Error for controlling behavior in the presence of
// errors.
func Load(fset *token.FileSet, in *ast.Package, conf *Config) (pkg *Package, err error) {
	n := len(in.Files) + len(in.GoFiles)
	files := make([]*goast.File, 0, n)
	for filename, f := range in.Files {
		if !strings.HasSuffix(filename, "_test.gop") {
			files = append(files, togo.ASTFile(f, 0))
		}
	}
	for filename, f := range in.GoFiles {
		if !strings.HasSuffix(filename, "_test.go") {
			files = append(files, f)
		}
	}
	return conf.Check(in.Name, fset, files, nil)
}
