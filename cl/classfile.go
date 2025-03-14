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
	gotoken "go/token"
	"go/types"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/goplus/gogen"
	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/token"
	"github.com/goplus/mod/modfile"
)

// -----------------------------------------------------------------------------

type gmxClass struct {
	tname   string // class type
	clsfile string
	ext     string
	proj    *gmxProject
}

type spxObj struct {
	obj  gogen.Ref
	proj string
}

type gmxProject struct {
	gameClass  string    // <gmtype>.gmx
	game       gogen.Ref // Game (project base class)
	spfeats    spriteFeat
	sprite     map[string]spxObj // .spx => Sprite
	sptypes    []string          // <sptype>.spx
	scheds     []string
	schedStmts []goast.Stmt // nil or len(scheds) == 2 (delayload)
	pkgImps    []gogen.PkgRef
	pkgPaths   []string
	autoimps   map[string]pkgImp // auto-import statement in gop.mod
	hasScheds  bool
	gameIsPtr  bool
	isTest     bool
	hasMain_   bool
}

type spriteFeat uint

const (
	spriteClassfname spriteFeat = 1 << iota
	spriteClassclone
)

func spriteFeatures(game gogen.Ref) (feats spriteFeat) {
	if mainFn := findMethod(game, "Main"); mainFn != nil {
		sig := mainFn.Type().(*types.Signature)
		if t, ok := gogen.CheckSigFuncEx(sig); ok {
			if t, ok := t.(*gogen.TyTemplateRecvMethod); ok {
				sig = t.Func.Type().(*types.Signature)
			}
		}
		if sig.Variadic() {
			in := sig.Params()
			last := in.At(in.Len() - 1)
			elt := last.Type().(*types.Slice).Elem()
			if tn, ok := elt.(*types.Named); ok {
				elt = tn.Underlying()
			}
			if intf, ok := elt.(*types.Interface); ok {
				for i, n := 0, intf.NumMethods(); i < n; i++ {
					switch intf.Method(i).Name() {
					case "Classfname":
						feats |= spriteClassfname
					case "Classclone":
						feats |= spriteClassclone
					}
				}
			}
		}
	}
	return
}

func (p *gmxProject) hasMain() bool {
	if !p.hasMain_ {
		imps := p.pkgImps
		p.hasMain_ = len(imps) > 0 && imps[0].TryRef("GopTestClass") != nil
	}
	return p.hasMain_
}

func (p *gmxProject) getScheds(cb *gogen.CodeBuilder) []goast.Stmt {
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

var (
	repl = strings.NewReplacer(":", "", "#", "", "-", "_", ".", "_")
)

func ClassNameAndExt(file string) (name, clsfile, ext string) {
	fname := filepath.Base(file)
	clsfile, ext = modfile.SplitFname(fname)
	name = clsfile
	if strings.ContainsAny(name, ":#-.") {
		name = repl.Replace(name)
	}
	return
}

// GetFileClassType get ast.File classType
func GetFileClassType(file *ast.File, filename string, lookupClass func(ext string) (c *Project, ok bool)) (classType string, isTest bool) {
	if file.IsClass {
		var ext string
		classType, _, ext = ClassNameAndExt(filename)
		if file.IsNormalGox {
			isTest = strings.HasSuffix(ext, "_test.gox")
			if !isTest && classType == "main" {
				classType = "_main"
			}
		} else {
			isTest = strings.HasSuffix(ext, "test.gox")
		}
		if file.IsProj && classType == "main" {
			if gt, ok := lookupClass(ext); ok {
				classType = gt.Class
			}
		} else if isTest {
			classType = casePrefix + testNameSuffix(classType)
		}
	} else if strings.HasSuffix(filename, "_test.gop") {
		isTest = true
	}
	return
}

func isGoxTestFile(ext string) bool {
	return strings.HasSuffix(ext, "test.gox")
}

func loadClass(ctx *pkgCtx, pkg *gogen.Package, file string, f *ast.File, conf *Config) *gmxProject {
	tname, clsfile, ext := ClassNameAndExt(file)
	gt, ok := conf.LookupClass(ext)
	if !ok {
		panic("TODO: class not found")
	}
	p, ok := ctx.projs[gt.Ext]
	if !ok {
		pkgPaths := gt.PkgPaths
		p = &gmxProject{pkgPaths: pkgPaths, isTest: isGoxTestFile(ext)}
		ctx.projs[gt.Ext] = p

		p.pkgImps = make([]gogen.PkgRef, len(pkgPaths))
		for i, pkgPath := range pkgPaths {
			p.pkgImps[i] = pkg.Import(pkgPath)
		}

		if len(gt.Import) > 0 {
			autoimps := make(map[string]pkgImp)
			for _, imp := range gt.Import {
				pkgi := pkg.Import(imp.Path)
				name := imp.Name
				if name == "" {
					name = pkgi.Types.Name()
				}
				pkgName := types.NewPkgName(token.NoPos, pkg.Types, name, pkgi.Types)
				autoimps[name] = pkgImp{pkgi, pkgName}
			}
			p.autoimps = autoimps
		}

		spx := p.pkgImps[0]
		if gt.Class != "" {
			p.game, p.gameIsPtr = spxRef(spx, gt.Class)
			p.spfeats = spriteFeatures(p.game)
		}
		p.sprite = make(map[string]spxObj)
		for _, v := range gt.Works {
			obj, _ := spxRef(spx, v.Class)
			p.sprite[v.Ext] = spxObj{obj, v.Project}
		}
		if x := getStringConst(spx, "Gop_sched"); x != "" {
			p.scheds, p.hasScheds = strings.SplitN(x, ",", 2), true
		}
	}
	if f.IsProj {
		if tname == "main" {
			tname = gt.Class
		}
		if p.gameClass != "" {
			log.Panicln("multiple project files found:", tname, p.gameClass)
		}
		p.gameClass = tname
		p.hasMain_ = f.HasShadowEntry()
	} else {
		p.sptypes = append(p.sptypes, tname)
	}
	ctx.classes[f] = &gmxClass{tname, clsfile, ext, p}
	if debugLoad {
		log.Println("==> InitClass", tname, "isProj:", f.IsProj)
	}
	return p
}

func spxLookup(pkgImps []gogen.PkgRef, name string) gogen.Ref {
	for _, pkg := range pkgImps {
		if o := pkg.TryRef(name); o != nil {
			return o
		}
	}
	panic("spxLookup: symbol not found - " + name)
}

func spxTryRef(spx gogen.PkgRef, typ string) (obj types.Object, isPtr bool) {
	if strings.HasPrefix(typ, "*") {
		typ, isPtr = typ[1:], true
	}
	obj = spx.TryRef(typ)
	return
}

func spxRef(spx gogen.PkgRef, typ string) (obj gogen.Ref, isPtr bool) {
	obj, isPtr = spxTryRef(spx, typ)
	if obj == nil {
		panic(spx.Types.Name() + "." + typ + " not found")
	}
	return
}

func getStringConst(spx gogen.PkgRef, name string) string {
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
				gogen.InsertStmtFront(body, scheds[idx])
			})
		}
	}
}

const (
	casePrefix = "case"
)

func testNameSuffix(testType string) string {
	if c := testType[0]; c >= 'A' && c <= 'Z' {
		return testType
	}
	return "_" + testType
}

func gmxTestFunc(pkg *gogen.Package, testType string, isProj bool) {
	if isProj {
		genTestFunc(pkg, "TestMain", testType, "m", "M")
	} else {
		name := testNameSuffix(testType)
		genTestFunc(pkg, "Test"+name, casePrefix+name, "t", "T")
	}
}

func genTestFunc(pkg *gogen.Package, name, testType, param, paramType string) {
	testing := pkg.Import("testing")
	objT := testing.Ref(paramType)
	paramT := types.NewParam(token.NoPos, pkg.Types, param, types.NewPointer(objT.Type()))
	params := types.NewTuple(paramT)

	pkg.NewFunc(nil, name, params, nil, false).BodyStart(pkg).
		Val(pkg.Builtin().Ref("new")).Val(pkg.Ref(testType)).Call(1).
		MemberVal("TestMain").Val(paramT).Call(1).EndStmt().
		End()
}

func gmxCheckProjs(pkg *gogen.Package, ctx *pkgCtx) (*gmxProject, bool) {
	var projMain, projNoMain *gmxProject
	var multiMain, multiNoMain bool
	for _, v := range ctx.projs {
		if v.isTest {
			continue
		}
		if v.hasMain() {
			if projMain != nil {
				multiMain = true
			} else {
				projMain = v
			}
		} else {
			if projNoMain != nil {
				multiNoMain = true
			} else {
				projNoMain = v
			}
		}
		if v.game != nil { // just to make testcase happy
			gmxProjMain(pkg, ctx, v)
		}
	}
	if projMain != nil {
		return projMain, multiMain
	}
	return projNoMain, multiNoMain
}

func gmxProjMain(pkg *gogen.Package, parent *pkgCtx, proj *gmxProject) {
	base := proj.game
	classType := proj.gameClass
	if classType == "" {
		classType = base.Name()
		proj.gameClass = classType
	}

	ld := getTypeLoader(parent, parent.syms, token.NoPos, proj.gameClass)
	if ld.typ == nil {
		ld.typ = func() {
			if debugLoad {
				log.Println("==> Load > NewType", classType)
			}
			old, _ := pkg.SetCurFile(defaultGoFile, true)
			defer pkg.RestoreCurFile(old)

			baseType := base.Type()
			if proj.gameIsPtr {
				baseType = types.NewPointer(baseType)
			}
			name := base.Name()
			flds := []*types.Var{
				types.NewField(token.NoPos, pkg.Types, name, baseType, true),
			}
			decl := pkg.NewTypeDefs().NewType(classType)
			ld.typInit = func() { // decycle
				if debugLoad {
					log.Println("==> Load > InitType", classType)
				}
				old, _ := pkg.SetCurFile(defaultGoFile, true)
				defer pkg.RestoreCurFile(old)

				decl.InitType(pkg, types.NewStruct(flds, nil))
			}
			parent.tylds = append(parent.tylds, ld)
		}
	}
	ld.methods = append(ld.methods, func() {
		old, _ := pkg.SetCurFile(defaultGoFile, true)
		defer pkg.RestoreCurFile(old)
		doInitType(ld)

		t := pkg.Ref(classType).Type()
		recv := types.NewParam(token.NoPos, pkg.Types, "this", types.NewPointer(t))
		fn := pkg.NewFunc(recv, "Main", nil, nil, false)

		parent.inits = append(parent.inits, func() {
			old, _ := pkg.SetCurFile(defaultGoFile, true)
			defer pkg.RestoreCurFile(old)

			cb := fn.BodyStart(pkg).Typ(base.Type()).MemberVal("Main")

			// force remove //line comments for main func
			cb.SetComments(nil, false)

			sigParams := cb.Get(-1).Type.(*types.Signature).Params()
			if _, ok := sigParams.At(0).Type().(*types.Pointer); !ok {
				cb.Val(recv) // template recv method
			} else {
				cb.Val(recv).MemberRef(base.Name()).UnaryOp(gotoken.AND)
			}
			narg := sigParams.Len()
			if narg > 1 {
				narg = 1 + len(proj.sptypes)
				new := pkg.Builtin().Ref("new")
				for _, spt := range proj.sptypes {
					sp := pkg.Ref(spt)
					cb.Val(new).Val(sp).Call(1)
				}
			}

			cb.Call(narg).EndStmt().End()
		})
	})
}

func genMainFunc(pkg *gogen.Package, gameClass string) {
	if o := pkg.TryRef(gameClass); o != nil {
		// new(gameClass).Main()
		new := pkg.Builtin().Ref("new")
		pkg.NewFunc(nil, "main", nil, nil, false).
			BodyStart(pkg).Val(new).Val(o).Call(1).MemberVal("Main").Call(0).EndStmt().End()
	}
}

func findMethod(o types.Object, name string) *types.Func {
	if obj, ok := o.(*types.TypeName); ok {
		if t, ok := obj.Type().(*types.Named); ok {
			for i, n := 0, t.NumMethods(); i < n; i++ {
				f := t.Method(i)
				if f.Name() == name {
					return f
				}
			}
		}
	}
	return nil
}

func makeMainSig(recv *types.Var, f *types.Func) *types.Signature {
	const (
		paramNameTempl = "_gop_arg0"
	)
	sig := f.Type().(*types.Signature)
	in := sig.Params()
	nin := in.Len()
	pkg := recv.Pkg()
	params := make([]*types.Var, nin)
	paramName := []byte(paramNameTempl)
	for i := 0; i < nin; i++ {
		paramName[len(paramNameTempl)-1] = byte('0' + i)
		params[i] = types.NewParam(token.NoPos, pkg, string(paramName), in.At(i).Type())
	}
	return types.NewSignatureType(recv, nil, nil, types.NewTuple(params...), nil, false)
}

func astFnClassfname(c *gmxClass) *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: "Classfname",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{Type: &ast.Ident{Name: "string"}},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(c.clsfile)},
					},
				},
			},
		},
		Shadow: true,
	}
}

func astFnClassclone() *ast.FuncDecl {
	ret := &ast.Ident{Name: "_gop_ret"}
	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: "Classclone",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{Type: &ast.Ident{Name: "any"}},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{ret},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.StarExpr{X: &ast.Ident{Name: "this"}},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{Op: token.AND, X: ret},
					},
				},
			},
		},
		Shadow: true,
	}
}

func astEmptyEntrypoint(f *ast.File) {
	var entry = getEntrypoint(f)
	var hasEntry bool
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Name.Name == entry {
				hasEntry = true
			}
		}
	}
	if !hasEntry {
		f.Decls = append(f.Decls, &ast.FuncDecl{
			Name: &ast.Ident{
				Name: entry,
			},
			Type: &ast.FuncType{
				Params: &ast.FieldList{},
			},
			Body:   &ast.BlockStmt{},
			Shadow: true,
		})
	}
}

func getEntrypoint(f *ast.File) string {
	switch {
	case f.IsProj:
		return "MainEntry"
	case f.IsClass:
		return "Main"
	case inMainPkg(f):
		return "main"
	default:
		return "init"
	}
}

func inMainPkg(f *ast.File) bool {
	return f.Name.Name == "main"
}

// -----------------------------------------------------------------------------
