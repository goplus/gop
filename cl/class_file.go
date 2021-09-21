/*
 Copyright 2021 The GoPlus Authors (goplus.org)

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package cl

import (
	"go/constant"
	"go/types"
	"path/filepath"
	"strings"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
	"github.com/goplus/gox"
)

// -----------------------------------------------------------------------------

var (
	gmxPkgPaths = map[string]string{
		".gmx": "github.com/goplus/spx",
	}
)

// RegisterClassFileType registers Go+ class file types.
func RegisterClassFileType(extGmx, extSpx string, pkgPath string) {
	parser.RegisterFileType(extGmx, ast.FileTypeGmx)
	parser.RegisterFileType(extSpx, ast.FileTypeSpx)
	if _, ok := gmxPkgPaths[extGmx]; !ok {
		gmxPkgPaths[extGmx] = pkgPath
	}
}

// -----------------------------------------------------------------------------

type gmxSettings struct {
	Pkg    string
	Class  string
	spx    *gox.PkgRef
	game   gox.Ref
	sprite gox.Ref
}

func newGmx(pkg *gox.Package, file string, gmx *ast.File, conf *Config) *gmxSettings {
	_, name := filepath.Split(file)
	ext := filepath.Ext(name)
	if idx := strings.Index(name, "."); idx > 0 {
		name = name[:idx]
	}
	p := &gmxSettings{Pkg: gmxPkgPaths[ext], Class: name}
	p.spx = pkg.Import(p.Pkg)
	p.game = spxRef(p.spx, "Gop_game", "Game")
	p.sprite = spxRef(p.spx, "Gop_sprite", "Sprite")
	return p
}

func getDefaultClass(file string) string {
	_, name := filepath.Split(file)
	if idx := strings.Index(name, "."); idx > 0 {
		name = name[:idx]
	}
	return name
}

func spxRef(spx *gox.PkgRef, name, typ string) gox.Ref {
	if v := getStringConst(spx, name); v != "" {
		typ = v
	}
	return spx.Ref(typ)
}

func getStringConst(spx *gox.PkgRef, name string) string {
	if o := spx.Ref(name); o != nil {
		if c, ok := o.(*types.Const); ok {
			return constant.StringVal(c.Val())
		}
	}
	return ""
}

func getFields(ctx *blockCtx, f *ast.File) (specs []ast.Spec) {
	decls := f.Decls
	i, n := 0, len(decls)
	for i < n {
		g, ok := decls[i].(*ast.GenDecl)
		if !ok {
			break
		}
		if g.Tok != token.IMPORT && g.Tok != token.CONST {
			if g.Tok == token.VAR {
				specs, g.Specs = g.Specs, nil
			}
			break
		}
		i++ // skip import, if any
	}
	return
}

// -----------------------------------------------------------------------------
