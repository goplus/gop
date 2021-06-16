/*
 Copyright 2020 The GoPlus Authors (goplus.org)

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

// Package cl compiles Go+ syntax trees (ast) into a backend code.
// For now the supported backends are `bytecode` and `golang`.
package cl

import (
	"reflect"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/ast/astutil"
	"github.com/goplus/gop/exec.spec"
	"github.com/goplus/gop/token"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

type pkgCtx struct {
	exec.Package
	infer   exec.Stack
	builtin exec.GoPackage
	out     exec.Builder
	usedfns []*funcDecl
	types   map[reflect.Type]*typeDecl
	pkg     *ast.Package
	fset    *token.FileSet
}

func newPkgCtx(out exec.Builder, pkg *ast.Package, fset *token.FileSet) *pkgCtx {
	pkgOut := out.GetPackage()
	builtin := pkgOut.FindGoPackage("")
	p := &pkgCtx{Package: pkgOut, builtin: builtin, out: out, pkg: pkg, fset: fset}
	p.types = make(map[reflect.Type]*typeDecl)
	p.infer.Init()
	return p
}

func (p *pkgCtx) code(v ast.Node) string {
	if expr, ok := v.(*ast.ParenExpr); ok {
		v = expr.X
	}
	_, code := p.getCodeInfo(v)
	return code
}

func (p *pkgCtx) getCodeInfo(v ast.Node) (token.Position, string) {
	start, end := v.Pos(), v.End()
	pos := p.fset.Position(start)
	if f, ok := p.pkg.Files[pos.Filename]; ok {
		return pos, string(f.Code[pos.Offset : pos.Offset+int(end-start)])
	}
	log.Panicln("pkgCtx.getCodeInfo failed: file not found -", pos.Filename)
	return pos, ""
}

func (p *pkgCtx) use(f *funcDecl) {
	if f.used {
		return
	}
	p.usedfns = append(p.usedfns, f)
	f.used = true
}

func (p *pkgCtx) resolveFuncs() {
	for {
		n := len(p.usedfns)
		if n == 0 {
			break
		}
		f := p.usedfns[n-1]
		p.usedfns = p.usedfns[:n-1]
		f.Compile()
	}
}

type fileCtx struct {
	*blockCtx // it's global blockCtx
	imports   map[string]string
}

func newFileCtx(block *blockCtx) *fileCtx {
	return &fileCtx{blockCtx: block, imports: make(map[string]string)}
}

// -----------------------------------------------------------------------------

type funcCtx struct {
	fun          exec.FuncInfo
	labels       map[string]*flowLabel
	currentFlow  *flowCtx
	currentLabel *ast.LabeledStmt
}

func newFuncCtx(fun exec.FuncInfo) *funcCtx {
	return &funcCtx{labels: map[string]*flowLabel{}, fun: fun}
}

type flowLabel struct {
	ctx *blockCtx
	exec.Label
	jumps []*blockCtx
}

type flowCtx struct {
	parent    *flowCtx
	name      string
	postLabel exec.Label
	doneLabel exec.Label
}

func (fc *funcCtx) nextFlow(post, done exec.Label, name string) {
	fc.currentFlow = &flowCtx{
		parent:    fc.currentFlow,
		name:      name,
		postLabel: post,
		doneLabel: done,
	}
}

func (fc *funcCtx) getBreakLabel(labelName string) (label exec.Label) {
	for i := fc.currentFlow; i != nil; i = i.parent {
		if i.doneLabel != nil {
			if labelName == "" || i.name == labelName {
				return i.doneLabel
			}
		}
	}
	return nil
}

func (fc *funcCtx) getContinueLabel(labelName string) (label exec.Label) {
	for i := fc.currentFlow; i != nil; i = i.parent {
		if i.postLabel != nil {
			if labelName == "" || i.name == labelName {
				return i.postLabel
			}
		}
	}
	return nil
}

func (fc *funcCtx) checkLabels() {
	for name, fl := range fc.labels {
		if fl.ctx == nil {
			log.Panicf("label %s not defined\n", name)
		}
		if !checkLabel(fl) {
			log.Panicf("goto %s jumps into illegal block\n", name)
		}
	}
	fc.labels = map[string]*flowLabel{}
}

func checkLabel(fl *flowLabel) bool {
	to := fl.ctx
	for _, j := range fl.jumps {
		if !j.canJmpTo(to) {
			return false
		}
	}
	return true
}

type blockCtx struct {
	*pkgCtx
	*funcCtx
	file            *fileCtx
	parent          *blockCtx
	syms            map[string]iSymbol
	fieldStructType reflect.Type
	fieldIndex      []int
	fieldExprX      func()
	indirect        int
	underscore      int
	noExecCtx       bool
	takeAddr        bool
	checkFlag       bool
	checkLoadAddr   bool
}

// function block ctx
func newExecBlockCtx(parent *blockCtx) *blockCtx {
	return &blockCtx{
		pkgCtx:    parent.pkgCtx,
		file:      parent.file,
		parent:    parent,
		syms:      make(map[string]iSymbol),
		noExecCtx: false,
	}
}

// normal block ctx, eg. if/switch/for/etc.
func newNormBlockCtx(parent *blockCtx) *blockCtx {
	return newNormBlockCtxEx(parent, true)
}

func newNormBlockCtxEx(parent *blockCtx, noExecCtx bool) *blockCtx {
	return &blockCtx{
		pkgCtx:    parent.pkgCtx,
		file:      parent.file,
		parent:    parent,
		funcCtx:   parent.funcCtx,
		syms:      make(map[string]iSymbol),
		noExecCtx: noExecCtx,
	}
}

// global block ctx
func newGblBlockCtx(pkg *pkgCtx) *blockCtx {
	return &blockCtx{
		pkgCtx:    pkg,
		parent:    nil,
		syms:      make(map[string]iSymbol),
		noExecCtx: true,
		funcCtx:   newFuncCtx(nil),
	}
}

func (ctx *blockCtx) resetFieldIndex() {
	ctx.fieldStructType = nil
	ctx.fieldIndex = nil
	ctx.fieldExprX = nil
}

func (p *blockCtx) requireLabel(name string) exec.Label {
	fl := p.labels[name]
	if fl == nil {
		fl = &flowLabel{
			Label: p.NewLabel(name),
		}
		p.labels[name] = fl
	}
	fl.jumps = append(fl.jumps, p)
	return fl.Label
}

func (p *blockCtx) defineLabel(name string) exec.Label {
	fl, ok := p.labels[name]
	if ok {
		if fl.ctx != nil {
			log.Panicf("label %s already defined at other position \n", name)
		}
		fl.ctx = p
	} else {
		fl = &flowLabel{
			ctx:   p,
			Label: p.NewLabel(name),
		}
		p.labels[name] = fl
	}
	return fl.Label
}

func (p *blockCtx) canJmpTo(to *blockCtx) bool {
	for from := p; from != nil; from = from.parent {
		if from == to {
			return true
		}
	}
	return false
}

func (p *blockCtx) getNestDepth() (nestDepth uint32) {
	for {
		if !p.noExecCtx {
			nestDepth++
		}
		if p = p.parent; p == nil {
			return
		}
	}
}

func (p *blockCtx) exists(name string) (ok bool) {
	if _, ok = p.syms[name]; ok {
		return
	}
	if p.parent == nil { // it's global blockCtx
		_, ok = p.file.imports[name]
	}
	return
}

func (p *blockCtx) find(name string) (sym interface{}, ok bool) {
	ctx := p
	for ; p != nil; p = p.parent {
		if sym, ok = p.syms[name]; ok {
			return
		}
	}
	sym, ok = ctx.file.imports[name]
	return
}

func (p *blockCtx) findType(name string) (decl *typeDecl, err error) {
	v, ok := p.find(name)
	if !ok {
		return nil, ErrNotFound
	}
	if decl, ok = v.(*typeDecl); ok {
		return
	}
	return nil, ErrSymbolNotType
}

func (p *blockCtx) findFunc(name string) (addr *funcDecl, err error) {
	v, ok := p.find(name)
	if !ok {
		return nil, ErrNotFound
	}
	if addr, ok = v.(*funcDecl); ok {
		return
	}
	return nil, ErrSymbolNotFunc
}

// getCtxVar finds a var in currentCtx only
func (p *blockCtx) getCtxVar(name string) (addr iVar, err error) {
	v, ok := p.syms[name]
	if !ok {
		return nil, ErrNotFound
	}
	if addr, ok = v.(iVar); ok {
		return
	}
	return nil, ErrSymbolNotVariable
}

func (p *blockCtx) findVar(name string) (addr iVar, err error) {
	v, ok := p.find(name)
	if !ok {
		return nil, ErrNotFound
	}
	if addr, ok = v.(iVar); ok {
		return
	}
	return nil, ErrSymbolNotVariable
}

func (p *blockCtx) insertFuncVars(in []reflect.Type, args []string, rets []exec.Var) {
	n := len(args)
	if n > 0 {
		for i := n - 1; i >= 0; i-- {
			name := args[i]
			if name == "" { // unnamed argument
				continue
			}
			if p.exists(name) {
				log.Panicln("insertStkVars failed: symbol exists -", name)
			}
			p.syms[name] = &stackVar{index: int32(i - n), typ: in[i]}
		}
	}
	for _, ret := range rets {
		if ret.IsUnnamedOut() {
			continue
		}
		p.syms[ret.Name()] = &execVar{ret}
	}
}

func (p *blockCtx) insertVar(name string, typ reflect.Type, inferOnly ...bool) *execVar {
	if name != "_" && p.exists(name) {
		log.Panicln("insertVar failed: symbol exists -", name)
	}
	v := p.NewVar(typ, name)
	if inferOnly == nil {
		p.out.DefineVar(v)
	}
	ev := &execVar{v}
	p.syms[name] = ev
	return ev
}

func (p *blockCtx) insertFunc(name string, fun *funcDecl) {
	if p.exists(name) {
		log.Panicln("insertFunc failed: symbol exists -", name)
	}
	p.syms[name] = fun
}

func (p *blockCtx) insertMethod(recv astutil.RecvInfo, methodName string, decl *ast.FuncDecl, ctx *blockCtx) {
	if p.parent != nil {
		log.Panicln("insertMethod failed: unexpected - non global method declaration?")
	}
	typ, err := p.findType(recv.Type)
	if err == ErrNotFound {
		typ = new(typeDecl)
		p.syms[recv.Type] = typ
	} else if err != nil {
		log.Panicln("insertMethod failed:", err)
	} else if typ.Alias {
		log.Panicln("insertMethod failed: alias?")
	}

	var t reflect.Type = typ.Type
	pointer := recv.Pointer
	for pointer > 0 {
		t = reflect.PtrTo(t)
		pointer--
	}

	method := newFuncDecl(methodName, decl.Recv, decl.Type, decl.Body, ctx)

	if typ.Methods == nil {
		typ.Methods = map[string]*funcDecl{methodName: method}
	} else {
		if _, ok := typ.Methods[methodName]; ok {
			log.Panicln("insertMethod failed: method exists -", recv.Type, methodName)
		}
		typ.Methods[methodName] = method
	}
}

// -----------------------------------------------------------------------------

func (c *blockCtx) findMethod(typ reflect.Type, methodName string) (*funcDecl, bool) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if tDecl, ok := c.types[typ]; ok {
		fDecl, ok := tDecl.Methods[methodName]
		return fDecl, ok
	}
	return nil, false
}
