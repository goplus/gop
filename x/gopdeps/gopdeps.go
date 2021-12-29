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

package gopdeps

import (
	"path"
	"strconv"
	"strings"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
)

// -----------------------------------------------------------------------------

func (p *ImportsParser) ParseGopImport(gopfile string) (err error) {
	f, err := parser.ParseFile(p.fset, gopfile, nil, parser.ImportsOnly)
	if err != nil {
		return
	}
	for _, imp := range f.Imports {
		p.importGop(imp)
	}
	return
}

func (p *ImportsParser) importGop(spec *ast.ImportSpec) {
	pkgPath := p.canonicalGop(gopToString(spec.Path))
	p.imports[pkgPath] = struct{}{}
}

func (p *ImportsParser) canonicalGop(pkgPath string) string {
	if strings.HasPrefix(pkgPath, "gop/") {
		return "github.com/goplus/" + pkgPath
	} else if strings.HasPrefix(pkgPath, ".") {
		return path.Join(p.mod.Path(), pkgPath)
	}
	return pkgPath
}

func gopToString(l *ast.BasicLit) string {
	if l.Kind == token.STRING {
		s, err := strconv.Unquote(l.Value)
		if err == nil {
			return s
		}
	}
	panic("TODO: toString - convert ast.BasicLit to string failed")
}

// -----------------------------------------------------------------------------
