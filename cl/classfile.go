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

package cl

import (
	"go/constant"
	"go/types"
	"path/filepath"
	"strings"

	goast "go/ast"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
	"github.com/goplus/gox"
)

// -----------------------------------------------------------------------------

type gmxInfo struct {
	extSpx   string
	pkgPaths []string
}

var (
	gmxTypes = map[string]gmxInfo{
		".gmx": {".spx", []string{"github.com/goplus/spx", "math"}},
		".spc": {"", []string{"github.com/qbox/gopla/spc", "github.com/goplus/spx", "math"}},
		// TODO: dynamic register
	}
)

// RegisterClassFileType registers Go+ class file types.
func RegisterClassFileType(extGmx, extSpx string, pkgPaths ...string) {
	if pkgPaths == nil {
		panic("RegisterClassFileType: no pkgPath specified")
	}
	parser.RegisterFileType(extGmx, ast.FileTypeGmx)
	if extSpx != "" {
		parser.RegisterFileType(extSpx, ast.FileTypeSpx)
	}
	if _, ok := gmxTypes[extGmx]; !ok {
		gmxTypes[extGmx] = gmxInfo{extSpx, pkgPaths}
	}
}

// -----------------------------------------------------------------------------

type gmxSettings struct {
	gameClass  string
	extSpx     string
	game       gox.Ref
	sprite     gox.Ref
	scheds     []string
	schedStmts []goast.Stmt // nil or len(scheds) == 2 (delayload)
	pkgImps    []*gox.PkgRef
	pkgPaths   []string
	hasScheds  bool
	gameIsPtr  bool
}

func (p *gmxSettings) getScheds(cb *gox.CodeBuilder) []goast.Stmt {
	if !p.hasScheds {
		return nil
	}
	if p.schedStmts == nil {
		p.schedStmts = make([]goast.Stmt, 2)
		for i, v := range p.scheds {
			fn := cb.Val(spxLookup(p.pkgImps, v)).Call(0).InternalStack().Pop().Val
			p.schedStmts[i] = &goast.ExprStmt{X: fn}
		}
		if len(p.scheds) < 2 {
			p.schedStmts[1] = p.schedStmts[0]
		}
	}
	return p.schedStmts
}

func newGmx(pkg *gox.Package, file string) *gmxSettings {
	_, name := filepath.Split(file)
	ext := filepath.Ext(name)
	if idx := strings.Index(name, "."); idx > 0 {
		name = name[:idx]
		if name == "main" {
			name = "_main"
		}
	}
	gt := gmxTypes[ext]
	pkgPaths := gt.pkgPaths
	p := &gmxSettings{extSpx: gt.extSpx, gameClass: name, pkgPaths: pkgPaths}
	p.pkgImps = make([]*gox.PkgRef, len(pkgPaths))
	for i, pkgPath := range pkgPaths {
		p.pkgImps[i] = pkg.Import(pkgPath)
	}
	spx := p.pkgImps[0]
	p.game, p.gameIsPtr = spxRef(spx, "Gop_game", "Game")
	if gt.extSpx != "" {
		p.sprite, _ = spxRef(spx, "Gop_sprite", "Sprite")
	}
	if x := getStringConst(spx, "Gop_sched"); x != "" {
		p.scheds, p.hasScheds = strings.SplitN(x, ",", 2), true
	}
	return p
}

func spxLookup(pkgImps []*gox.PkgRef, name string) gox.Ref {
	for _, pkg := range pkgImps {
		if o := pkg.TryRef(name); o != nil {
			return o
		}
	}
	panic("spxLookup: symbol not found - " + name)
}

func getDefaultClass(file string) string {
	_, name := filepath.Split(file)
	if idx := strings.Index(name, "."); idx > 0 {
		name = name[:idx]
	}
	return name
}

func spxRef(spx *gox.PkgRef, name, typ string) (obj gox.Ref, isPtr bool) {
	if v := getStringConst(spx, name); v != "" {
		typ = v
		if strings.HasPrefix(typ, "*") {
			typ, isPtr = typ[1:], true
		}
	}
	obj = spx.Ref(typ)
	return
}

func getStringConst(spx *gox.PkgRef, name string) string {
	if o := spx.TryRef(name); o != nil {
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

func setBodyHandler(ctx *blockCtx) {
	if ctx.fileType > 0 { // in a Go+ class file
		if scheds := ctx.getScheds(ctx.cb); scheds != nil {
			ctx.cb.SetBodyHandler(func(body *goast.BlockStmt, kind int) {
				idx := 0
				if len(body.List) == 0 {
					idx = 1
				}
				gox.InsertStmtFront(body, scheds[idx])
			})
		}
	}
}

func gmxMainFunc(p *gox.Package, ctx *pkgCtx) {
	if o := p.Types.Scope().Lookup(ctx.gameClass); o != nil && hasMethod(o, "MainEntry") {
		// new(Game).Main()
		p.NewFunc(nil, "main", nil, nil, false).BodyStart(p).
			Val(p.Builtin().Ref("new")).Val(o).Call(1).
			MemberVal("Main").Call(0).EndStmt().
			End()
	}
}

// -----------------------------------------------------------------------------
