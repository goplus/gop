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
	"errors"
	"fmt"
	"path"
	"reflect"
	"syscall"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/ast/astutil"
	"github.com/goplus/gop/exec.spec"
	"github.com/goplus/gop/exec/bytecode"
	"github.com/goplus/gop/token"
	"github.com/goplus/reflectx"
	"github.com/qiniu/x/log"
)

var (
	// ErrNotFound error
	ErrNotFound = syscall.ENOENT

	// ErrNotAMainPackage error
	ErrNotAMainPackage = errors.New("not a main package")

	// ErrMainFuncNotFound error
	ErrMainFuncNotFound = errors.New("main function not found")

	// ErrSymbolNotVariable error
	ErrSymbolNotVariable = errors.New("symbol exists but not a variable")

	// ErrSymbolNotFunc error
	ErrSymbolNotFunc = errors.New("symbol exists but not a func")

	// ErrSymbolNotType error
	ErrSymbolNotType = errors.New("symbol exists but not a type")
)

var (
	// CallBuiltinOp calls BuiltinOp
	CallBuiltinOp func(kind exec.Kind, op exec.Operator, data ...interface{}) interface{}
)

// -----------------------------------------------------------------------------

// CompileError represents a compiling time error.
type CompileError struct {
	At  ast.Node
	Err error
}

func (p *CompileError) Error() string {
	return p.Err.Error()
}

func newError(at ast.Node, format string, params ...interface{}) *CompileError {
	err := fmt.Errorf(format, params...)
	return &CompileError{at, err}
}

func logError(ctx *blockCtx, at ast.Node, format string, params ...interface{}) {
	err := newError(at, format, params...)
	log.Error(err)
}

func logPanic(ctx *blockCtx, at ast.Node, format string, params ...interface{}) {
	err := newError(at, format, params...)
	log.Panicln(err)
}

func logNonIntegerIdxPanic(ctx *blockCtx, v ast.Node, kind reflect.Kind) {
	logPanic(ctx, v, `non-integer %v index %v`, kind, ctx.code(v))
}

func logIllTypeMapIndexPanic(ctx *blockCtx, v ast.Node, t, typIdx reflect.Type) {
	logPanic(ctx, v, `cannot use %v (type %v) as type %v in map index`, ctx.code(v), t, typIdx)
}

// -----------------------------------------------------------------------------

// - varName => *exec.Var
// - stkVarName => *stackVar
// - pkgName => pkgPath
// - funcName => *funcDecl
// - typeName => *typeDecl
//
type iSymbol = interface{}

type iVar interface {
	inCurrentCtx(ctx *blockCtx) bool
	getType() reflect.Type
}

type execVar struct {
	v exec.Var
}

func (p *execVar) inCurrentCtx(ctx *blockCtx) bool {
	return ctx.out.InCurrentCtx(p.v)
}

func (p *execVar) getType() reflect.Type {
	return p.v.Type()
}

type stackVar struct {
	typ   reflect.Type
	index int32
}

func (p *stackVar) inCurrentCtx(ctx *blockCtx) bool {
	return true
}

func (p *stackVar) getType() reflect.Type {
	return p.typ
}

// -----------------------------------------------------------------------------

// A Package represents a Go+ package.
type Package struct {
	syms map[string]iSymbol
}

// PkgAct represents a package compiling action.
type PkgAct int

const (
	// PkgActNone - do nothing
	PkgActNone PkgAct = iota
	// PkgActClMain - compile main function
	PkgActClMain
	// PkgActClAll - compile all things
	PkgActClAll
)

// NewPackage creates a Go+ package instance.
func NewPackage(out exec.Builder, pkg *ast.Package, fset *token.FileSet, act PkgAct) (p *Package, err error) {
	return NewPackageEx(out, pkg, fset, act, nil)
}

func NewPackageEx(out exec.Builder, pkg *ast.Package, fset *token.FileSet, act PkgAct, imports map[string]string) (p *Package, err error) {
	if pkg == nil {
		return nil, ErrNotFound
	}
	if CallBuiltinOp == nil {
		log.Panicln("NewPackage failed: variable CallBuiltinOp is uninitialized")
	}
	// reset reflectx
	reflectx.Reset()
	//
	p = &Package{}
	ctxPkg := newPkgCtx(out, pkg, fset)
	ctx := newGblBlockCtx(ctxPkg)
	for _, f := range pkg.Files {
		loadFile(ctx, f, imports)
	}
	switch act {
	case PkgActClAll:
		for _, sym := range ctx.syms {
			if f, ok := sym.(*funcDecl); ok && f.fi != nil {
				ctxPkg.use(f)
			}
		}
		if pkg.Name != "main" {
			break
		}
		fallthrough
	case PkgActClMain:
		if pkg.Name != "main" {
			return nil, ErrNotAMainPackage
		}
		entry, err := ctx.findFunc("main")
		if err != nil {
			if err == ErrNotFound {
				err = ErrMainFuncNotFound
			}
			return p, err
		}
		if entry.ctx.noExecCtx {
			ctx.file = entry.ctx.file
			compileBlockStmtWithout(ctx, entry.body)
			ctx.checkLabels()
		} else {
			out.CallFunc(entry.Get(), 0)
			ctxPkg.use(entry)
		}
		out.Return(-1)
	}
	ctxPkg.resolveFuncs()

	p.syms = ctx.syms
	return
}

// SymKind represents a symbol kind.
type SymKind uint

const (
	// SymInvalid - invalid symbol kind
	SymInvalid SymKind = iota
	// SymConst - symbol is a const
	SymConst
	// SymVar - symbol is a variable
	SymVar
	// SymFunc - symbol is a function
	SymFunc
	// SymType - symbol is a type
	SymType
)

// Find lookups a symbol and returns it's kind and the object instance.
func (p *Package) Find(name string) (kind SymKind, v interface{}, ok bool) {
	if v, ok = p.syms[name]; !ok {
		return
	}
	switch v.(type) {
	case *constVal:
		kind = SymConst
	case *exec.Var:
		kind = SymVar
	case *funcDecl:
		kind = SymFunc
	case *typeDecl:
		kind = SymType
	default:
		log.Panicln("Package.Find: unknown symbol type -", reflect.TypeOf(v))
	}
	return
}

type namedType struct {
	name     string
	methods  []string
	pmethods []string
}

func loadNamedType(ctx *blockCtx, d *ast.FuncDecl) {
	typ := d.Recv.List[0].Type
	var ptr bool
	if expr, ok := typ.(*ast.StarExpr); ok {
		typ = expr.X
		ptr = true
	}
	name := typ.(*ast.Ident).Name
	nt, ok := ctx.named[name]
	if !ok {
		nt = &namedType{name: name}
		ctx.named[name] = nt
	}
	if ptr {
		nt.pmethods = append(nt.pmethods, d.Name.Name)
	} else {
		nt.methods = append(nt.methods, d.Name.Name)
	}
}

func loadFile(ctx *blockCtx, f *ast.File, imports map[string]string) {
	file := newFileCtx(ctx)
	last := len(f.Decls) - 1
	ctx.file = file
	for name, pkg := range imports {
		ctx.file.imports[name] = pkg
	}
	// load import
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			switch d.Tok {
			case token.IMPORT:
				loadImports(file, d)
			}
		}
	}
	//	load named type methods
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Recv != nil {
				loadNamedType(ctx, d)
			}
		}
	}
	// load type & const
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			switch d.Tok {
			case token.TYPE:
				loadTypes(ctx, d)
			case token.CONST:
				loadConsts(ctx, d)
			}
		}
	}
	// load func & method
	for i, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			loadFunc(ctx, d, f.NoEntrypoint && i == last)
		}
	}
	// register methods
	pkg := bytecode.FindGoPackage(ctx.pkg.Name)
	if pkg == nil {
		pkg = bytecode.NewGoPackage(ctx.pkg.Name)
	}
	var decls []*typeDecl
	for typ, decl := range ctx.types {
		if typ.Kind() == reflect.Interface {
			registerInterface(pkg.(*bytecode.GoPackage), typ)
			continue
		}
		if decl.Type.Kind() == reflect.Struct {
			typ := decl.Type
			for i := 0; i < typ.NumField(); i++ {
				_, ftyp := countPtr(typ.Field(i).Type)
				if d, ok := ctx.types[ftyp]; ok {
					decl.Depends = append(decl.Depends, d)
				}
			}
		}
		decls = append(decls, decl)
	}
	registerTypeDecls(ctx, pkg.(*bytecode.GoPackage), decls)
	for _, mt := range ctx.mtypeList {
		mt.Update(ctx.mtype)
		mt.RegisterMethod(pkg.(*bytecode.GoPackage))
		ctx.out.MethodOf(mt.typ, mt.infos)
	}
	// load vars
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			switch d.Tok {
			case token.VAR:
				compileStmt(ctx, &ast.DeclStmt{decl})
			}
		}
	}
}

func registerTypeDecls(ctx *blockCtx, pkg *bytecode.GoPackage, decls []*typeDecl) {
	for _, decl := range decls {
		if len(decl.Depends) > 0 {
			registerTypeDecls(ctx, pkg, decl.Depends)
		}
		if decl.Register {
			continue
		}
		decl.Register = true
		typ := decl.Type
		var infos []exec.FuncInfo
		for _, mfun := range decl.Methods {
			ctx.use(mfun)
			infos = append(infos, mfun.Get())
		}
		mt := NewMethodType(typ, infos)
		ctx.mtypeList = append(ctx.mtypeList, mt)
		decl.Type = mt.typ
		ctx.mtype[typ] = mt.typ
	}
}

func extractRealType(ctx *blockCtx, t reflect.Type) reflect.Type {
	n, vt := countPtr(t)
	if r, ok := ctx.mtype[vt]; ok {
		for i := 0; i < n; i++ {
			r = reflect.PtrTo(r)
		}
		return r
	}
	return t
}

func extractStructType(ctx *blockCtx, typ reflect.Type) reflect.Type {
	n := typ.NumField()
	fs := make([]reflect.StructField, n, n)
	for i := 0; i < n; i++ {
		field := typ.Field(i)
		if !ast.IsExported(field.Name) {
			field.PkgPath = ctx.pkg.Name
		}
		field.Type = extractRealType(ctx, field.Type)
		fs[i] = field
	}
	return reflectx.NamedStructOf(typ.PkgPath(), typ.Name(), fs)
}

func loadImports(ctx *fileCtx, d *ast.GenDecl) {
	for _, item := range d.Specs {
		loadImport(ctx, item.(*ast.ImportSpec))
	}
}

func loadImport(ctx *fileCtx, spec *ast.ImportSpec) {
	var pkgPath = astutil.ToString(spec.Path)
	var name string
	if spec.Name != nil {
		name = spec.Name.Name
		switch name {
		case "_", ".":
			panic("not impl")
		}
	} else {
		name = path.Base(pkgPath)
	}
	ctx.imports[name] = pkgPath
}

func loadTypes(ctx *blockCtx, d *ast.GenDecl) {
	for _, item := range d.Specs {
		loadType(ctx, item.(*ast.TypeSpec))
	}
}

func toMethods(nt *namedType) []reflectx.Method {
	ftyp := reflect.FuncOf(nil, nil, false)
	fn := func([]reflect.Value) []reflect.Value { return nil }
	count := len(nt.methods) + len(nt.pmethods)
	ms := make([]reflectx.Method, count, count)
	for i, v := range nt.methods {
		ms[i] = reflectx.MakeMethod(v, false, ftyp, fn)
	}
	sz := len(nt.methods)
	for i, v := range nt.pmethods {
		ms[i+sz] = reflectx.MakeMethod(v, true, ftyp, fn)
	}
	return ms
}

func loadType(ctx *blockCtx, spec *ast.TypeSpec) {
	if ctx.exists(spec.Name.Name) {
		log.Panicln("loadType failed: symbol exists -", spec.Name.Name)
	}
	t := toType(ctx, spec.Type).(reflect.Type)
	typ := reflectx.NamedTypeOf(ctx.pkg.Name, spec.Name.Name, t)
	if nt, ok := ctx.named[spec.Name.Name]; ok {
		typ = reflectx.MethodOf(typ, toMethods(nt))
	} else if typ.Kind() == reflect.Struct {
		typ = reflectx.MethodOf(typ, nil)
	}
	ctx.out.DefineType(typ, spec.Name.Name)

	tDecl := &typeDecl{
		Type: typ,
	}
	ctx.syms[spec.Name.Name] = tDecl
	ctx.types[typ] = tDecl
}

var (
	iotaIndex int
	iotaUsed  bool
)

func newIotaValue() *constVal {
	iotaUsed = true
	return newConstVal(iotaIndex, exec.ConstUnboundInt)
}

func loadConsts(ctx *blockCtx, d *ast.GenDecl) {
	iotaIndex = 0
	iotaUsed = false
	var last *ast.ValueSpec
	for _, item := range d.Specs {
		spec := item.(*ast.ValueSpec)
		if spec.Type == nil && spec.Values == nil {
			spec.Type = last.Type
			spec.Values = last.Values
		}
		nnames := len(spec.Names)
		nvalue := len(spec.Values)
		if nvalue < nnames {
			log.Panicf("missing value in const declaration")
		} else if nvalue > nnames {
			log.Panicf("extra expression in const declaration")
		} else {
			for i := 0; i < nnames; i++ {
				loadConst(ctx, spec.Names[i].Name, spec.Type, spec.Values[i])
			}
			if iotaUsed {
				iotaIndex++
			}
		}
		last = spec
	}
}

func loadVars(ctx *blockCtx, d *ast.GenDecl, stmt ast.Stmt) {
	for _, item := range d.Specs {
		spec := item.(*ast.ValueSpec)
		for i := 0; i < len(spec.Names); i++ {
			start := ctx.out.StartStmt(stmt)
			name := spec.Names[i].Name
			if len(spec.Values) > i {
				loadVar(ctx, name, spec.Type, spec.Values[i])
			} else {
				loadVar(ctx, name, spec.Type, nil)
			}
			ctx.out.EndStmt(stmt, start)
		}
	}
}

func loadConst(ctx *blockCtx, name string, typ ast.Expr, value ast.Expr) {
	if name != "_" && ctx.exists(name) {
		log.Panicln("loadConst failed: symbol exists -", name)
	}
	compileExpr(ctx, value)
	in := ctx.infer.Pop()
	c := in.(*constVal)
	if typ != nil {
		t := toType(ctx, typ).(reflect.Type)
		v := boundConst(c.v, t)
		c.v = v
		c.kind = t.Kind()
	}
	if name != "_" {
		ctx.syms[name] = c
	}
}

func loadVar(ctx *blockCtx, name string, typ ast.Expr, value ast.Expr) {
	var t reflect.Type
	if typ != nil {
		t = toType(ctx, typ).(reflect.Type)
	}
	if value != nil {
		expr := compileExpr(ctx, value)
		in := ctx.infer.Get(-1)
		if t == nil {
			t = boundType(in.(iValue))
		}
		expr()
		addr := ctx.insertVar(name, t)
		checkType(addr.getType(), in, ctx.out)
		ctx.infer.PopN(1)
		ctx.out.StoreVar(addr.v)
	} else {
		ctx.insertVar(name, t)
	}
}

func loadFunc(ctx *blockCtx, d *ast.FuncDecl, isUnnamed bool) {
	var name = d.Name.Name
	if d.Recv != nil {
		recv := astutil.ToRecv(d.Recv)
		funCtx := newExecBlockCtx(ctx)
		funCtx.noExecCtx = isUnnamed
		funCtx.funcCtx = newFuncCtx(nil)
		ctx.insertMethod(recv, name, d, funCtx)
	} else if name == "init" {
		log.Panicln("loadFunc TODO: init")
	} else {
		funCtx := newExecBlockCtx(ctx)
		funCtx.noExecCtx = isUnnamed
		funCtx.funcCtx = newFuncCtx(nil)
		ctx.insertFunc(name, newFuncDecl(name, nil, d.Type, d.Body, funCtx))
	}
}

func registerInterface(pkg *bytecode.GoPackage, typ reflect.Type) {
	name := typ.Name()
	if name == "" {
		name = typ.String()
	}
	pkg.RegisterTypes(pkg.Type(name, typ))

	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		fnName := "(" + name + ")." + method.Name
		registerInterfaceMethod(pkg, fnName, typ, method.Name, method.Type)
	}
}

func registerInterfaceMethod(p *bytecode.GoPackage, fnname string, t reflect.Type, name string, typ reflect.Type) (addr uint32, kind exec.SymbolKind) {
	var tin []reflect.Type
	var fnInterface interface{}
	isVariadic := typ.IsVariadic()
	numIn := typ.NumIn()
	numOut := typ.NumOut()
	tin = make([]reflect.Type, numIn+1)
	tin[0] = t
	for i := 1; i < numIn+1; i++ {
		tin[i] = typ.In(i - 1)
	}
	tout := make([]reflect.Type, numOut)
	for i := 0; i < numOut; i++ {
		tout[i] = typ.Out(i)
	}
	typFunc := reflect.FuncOf(tin, tout, isVariadic)
	fnInterface = reflect.Zero(typFunc).Interface()
	arity := numIn + 1
	fnExec := func(i int, p *bytecode.Context) {
		if isVariadic {
			arity = i
		}
		args := p.GetArgs(arity)
		in := make([]reflect.Value, arity)
		for i, arg := range args {
			if arg != nil {
				in[i] = reflect.ValueOf(arg)
			} else if i >= numIn {
				in[i] = reflect.Zero(tin[numIn-1])
			} else {
				in[i] = reflect.Zero(tin[i])
			}
		}
		var out []reflect.Value
		out = in[0].MethodByName(name).Call(in[1:])
		if numOut > 0 {
			iout := make([]interface{}, numOut)
			for i := 0; i < numOut; i++ {
				iout[i] = out[i].Interface()
			}
			p.Ret(arity, iout...)
		}
	}
	if isVariadic {
		info := p.Funcv(fnname, fnInterface, fnExec)
		base := p.RegisterFuncvs(info)
		addr = uint32(base)
		kind = exec.SymbolFuncv
	} else {
		info := p.Func(fnname, fnInterface, fnExec)
		base := p.RegisterFuncs(info)
		addr = uint32(base)
		kind = exec.SymbolFunc
	}
	return
}

func extractFuncType(typ reflect.Type) reflect.Type {
	numIn := typ.NumIn()
	numOut := typ.NumOut()
	in := make([]reflect.Type, numIn-1)
	out := make([]reflect.Type, numOut)
	for i := 1; i < numIn; i++ {
		in[i-1] = typ.In(i)
	}
	for i := 0; i < numOut; i++ {
		out[i] = typ.Out(i)
	}
	return reflect.FuncOf(in, out, typ.IsVariadic())
}

type MethodType struct {
	typ     reflect.Type
	infos   []*exec.MethodInfo
	methods []reflectx.Method
	imap    map[string]*exec.MethodInfo
}

func NewMethodType(typ reflect.Type, infos []exec.FuncInfo) *MethodType {
	imap := make(map[string]*exec.MethodInfo)
	var methods []reflectx.Method
	var minfos []*exec.MethodInfo
	for _, fi := range infos {
		ftyp := fi.Type()
		mtyp := extractFuncType(ftyp)
		ptr := ftyp.In(0).Kind() == reflect.Ptr
		mi := &exec.MethodInfo{Info: fi}
		m := reflectx.MakeMethod(fi.Name(), ptr, mtyp, func(args []reflect.Value) []reflect.Value {
			return mi.Func(args)
		})
		imap[fi.Name()] = mi
		minfos = append(minfos, mi)
		methods = append(methods, m)
	}
	reflectx.UpdateMethod(typ, methods, make(map[reflect.Type]reflect.Type))
	return &MethodType{
		typ:     typ,
		methods: methods,
		infos:   minfos,
		imap:    imap,
	}
}

func (p *MethodType) Update(rmap map[reflect.Type]reflect.Type) {
	reflectx.UpdateMethod(p.typ, p.methods, rmap)
}

func (p *MethodType) RegisterMethod(pkg *bytecode.GoPackage) {
	name := p.typ.Name()
	pkg.RegisterTypes(pkg.Type(name, p.typ))
	registerTypeMethods(pkg, p.typ, p.imap)
}

func registerTypeMethods(pkg *bytecode.GoPackage, typ reflect.Type, imap map[string]*exec.MethodInfo) {
	name := typ.Name()
	skip := make(map[string]bool)
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		fnName := "(" + name + ")." + method.Name
		m, ok := imap[method.Name]
		if !ok {
			continue
		}
		skip[method.Name] = true
		registerTypeMethod(pkg, fnName, method.Func, m.Info)
	}
	typ = reflect.PtrTo(typ)
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		if skip[method.Name] {
			continue
		}
		m, ok := imap[method.Name]
		if !ok {
			continue
		}
		fnName := "(*" + name + ")." + method.Name
		registerTypeMethod(pkg, fnName, method.Func, m.Info)
	}
}

func registerTypeMethod(p *bytecode.GoPackage, fnname string, fun reflect.Value, fi exec.FuncInfo) (addr uint32, kind exec.SymbolKind) {
	if fi.IsVariadic() {
		fnExec := func(i int, p *bytecode.Context) {
			p.Callv(fi, uint32(i))
		}
		info := p.Funcv(fnname, fun.Interface(), fnExec)
		base := p.RegisterFuncvs(info)
		addr = uint32(base)
		kind = exec.SymbolFuncv
	} else {
		fnExec := func(i int, p *bytecode.Context) {
			p.Call(fi)
		}
		info := p.Func(fnname, fun.Interface(), fnExec)
		base := p.RegisterFuncs(info)
		addr = uint32(base)
		kind = exec.SymbolFunc
	}
	return
}

// -----------------------------------------------------------------------------
