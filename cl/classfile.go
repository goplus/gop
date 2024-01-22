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
	goast "go/ast"
	"go/constant"
	"go/types"
	"log"
	"path/filepath"
	"strings"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/token"
	"github.com/goplus/gox"
	"github.com/goplus/mod/modfile"
)

// -----------------------------------------------------------------------------

type gmxClass struct {
	tname string // class type
	ext   string
	proj  *gmxProject
}

type gmxProject struct {
	gameClass  string             // <gmtype>.gmx
	game       gox.Ref            // Game
	sprite     map[string]gox.Ref // .spx => Sprite
	sptypes    []string           // <sptype>.spx
	scheds     []string
	schedStmts []goast.Stmt // nil or len(scheds) == 2 (delayload)
	pkgImps    []*gox.PkgRef
	pkgPaths   []string
	hasScheds  bool
	gameIsPtr  bool
}

func (p *gmxProject) getScheds(cb *gox.CodeBuilder) []goast.Stmt {
	if p == nil || !p.hasScheds {
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

func ClassNameAndExt(file string) (name, ext string) {
	fname := filepath.Base(file)
	name, ext = modfile.SplitFname(fname)
	if idx := strings.Index(name, "."); idx > 0 {
		name = name[:idx]
	}
	return
}

func loadClass(ctx *pkgCtx, pkg *gox.Package, file string, f *ast.File, conf *Config) *gmxProject {
	tname, ext := ClassNameAndExt(file)
	gt, ok := conf.LookupClass(ext)
	if !ok {
		panic("TODO: class not found")
	}
	p, ok := ctx.projs[gt.Ext]
	if !ok {
		pkgPaths := gt.PkgPaths
		p = &gmxProject{pkgPaths: pkgPaths}
		ctx.projs[gt.Ext] = p

		p.pkgImps = make([]*gox.PkgRef, len(pkgPaths))
		for i, pkgPath := range pkgPaths {
			p.pkgImps[i] = pkg.Import(pkgPath)
		}
		spx := p.pkgImps[0]
		if gt.Class != "" {
			p.game, p.gameIsPtr = spxRef(spx, gt.Class)
		}
		p.sprite = make(map[string]types.Object)
		for _, v := range gt.Works {
			obj, _ := spxRef(spx, v.Class)
			p.sprite[v.Ext] = obj
		}
		if x := getStringConst(spx, "Gop_sched"); x != "" {
			p.scheds, p.hasScheds = strings.SplitN(x, ",", 2), true
		}
	}
	if f.IsProj {
		if p.gameClass != "" {
			panic("TODO: multiple project files found")
		}
		if tname == "main" {
			tname = gt.Class
		}
		p.gameClass = tname
	} else {
		p.sptypes = append(p.sptypes, tname)
	}
	ctx.classes[f] = gmxClass{tname, ext, p}
	if debugLoad {
		log.Println("==> InitClass", tname, "isProj:", f.IsProj)
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

func spxTryRef(spx *gox.PkgRef, typ string) (obj types.Object, isPtr bool) {
	if strings.HasPrefix(typ, "*") {
		typ, isPtr = typ[1:], true
	}
	obj = spx.TryRef(typ)
	return
}

func spxRef(spx *gox.PkgRef, typ string) (obj gox.Ref, isPtr bool) {
	obj, isPtr = spxTryRef(spx, typ)
	if obj == nil {
		panic(spx.Path() + "." + typ + " not found")
	}
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

func getFields(f *ast.File) []ast.Spec {
	for _, decl := range f.Decls {
		if g, ok := decl.(*ast.GenDecl); ok {
			if g.Tok == token.VAR {
				return g.Specs
			}
			continue
		}
		break
	}
	return nil
}

func setBodyHandler(ctx *blockCtx) {
	if proj := ctx.proj; proj != nil { // in a Go+ class file
		if scheds := proj.getScheds(ctx.cb); scheds != nil {
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

func gmxMainFunc(p *gox.Package, ctx *pkgCtx, noAutoGenMain bool) func() {
	if len(ctx.projs) == 1 { // only one project file
		var proj *gmxProject
		for _, v := range ctx.projs {
			proj = v
			break
		}
		scope := p.Types.Scope()
		class := proj.gameClass
		var o types.Object
		if class != "" {
			o = scope.Lookup(proj.gameClass)
			if noAutoGenMain && o != nil && hasMethod(o, "MainEntry") {
				noAutoGenMain = false
			}
		} else {
			o = proj.game
		}
		if !noAutoGenMain && o != nil {
			// new(Game).Main()
			// new(Game).Main(workers...)
			fn := p.NewFunc(nil, "main", nil, nil, false)
			return func() {
				new := p.Builtin().Ref("new")
				cb := fn.BodyStart(p).Val(new).Val(o).Call(1).MemberVal("Main")

				sig := cb.Get(-1).Type.(*types.Signature)
				narg := gmxMainNarg(sig)
				if narg > 0 {
					narg = len(proj.sptypes)
					for _, spt := range proj.sptypes {
						sp := scope.Lookup(spt)
						cb.Val(new).Val(sp).Call(1)
					}
				}

				cb.Call(narg).EndStmt().End()
			}
		}
	}
	return nil
}

func gmxMainNarg(sig *types.Signature) int {
	if fex, ok := gox.CheckFuncEx(sig); ok {
		if trm, ok := fex.(*gox.TyTemplateRecvMethod); ok {
			sig = trm.Func.Type().(*types.Signature)
			return sig.Params().Len() - 1
		}
	}
	return sig.Params().Len()
}

func hasMethod(o types.Object, name string) bool {
	if obj, ok := o.(*types.TypeName); ok {
		if t, ok := obj.Type().(*types.Named); ok {
			for i, n := 0, t.NumMethods(); i < n; i++ {
				if t.Method(i).Name() == name {
					return true
				}
			}
		}
	}
	return false
}

// -----------------------------------------------------------------------------
