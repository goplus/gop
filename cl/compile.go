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

func loadFile(ctx *blockCtx, f *ast.File, imports map[string]string) {
	file := newFileCtx(ctx)
	last := len(f.Decls) - 1
	ctx.file = file
	for name, pkg := range imports {
		ctx.file.imports[name] = pkg
	}
	for i, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			loadFunc(ctx, d, f.NoEntrypoint && i == last)
		case *ast.GenDecl:
			switch d.Tok {
			case token.IMPORT:
				loadImports(file, d)
			case token.TYPE:
				loadTypes(ctx, d)
			case token.CONST:
				loadConsts(ctx, d)
			case token.VAR:
				compileStmt(ctx, &ast.DeclStmt{decl})
			default:
				log.Panicln("tok:", d.Tok, "spec:", reflect.TypeOf(d.Specs).Elem())
			}
		default:
			log.Panicln("gopkg.Package.load: unknown decl -", reflect.TypeOf(decl))
		}
	}
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

func loadType(ctx *blockCtx, spec *ast.TypeSpec) {
	if ctx.exists(spec.Name.Name) {
		log.Panicln("loadType failed: symbol exists -", spec.Name.Name)
	}
	t := toType(ctx, spec.Type).(reflect.Type)
	typ := reflectx.NamedTypeOf(ctx.pkg.Name, spec.Name.Name, t)
	if t.Kind() == reflect.Interface {
		pkg := bytecode.FindGoPackage(ctx.pkg.Name)
		if pkg == nil {
			pkg = bytecode.NewGoPackage(ctx.pkg.Name)
		}
		registerInterface(pkg.(*bytecode.GoPackage), typ)
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
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		fnName := "(" + name + ")." + method.Name
		registerInterfaceMethod(pkg, fnName, typ, method.Name, method.Func, method.Type)
	}
}

func registerInterfaceMethod(p *bytecode.GoPackage, fnname string, t reflect.Type, name string, fun reflect.Value, typ reflect.Type) (addr uint32, kind exec.SymbolKind) {
	var tin []reflect.Type
	var fnInterface interface{}
	isVariadic := typ.IsVariadic()
	numIn := typ.NumIn()
	numOut := typ.NumOut()
	var arity int
	if fun.IsValid() { // type method
		tin = make([]reflect.Type, numIn)
		for i := 0; i < numIn; i++ {
			tin[i] = typ.In(i)
		}
		fnInterface = fun.Interface()
		arity = numIn
	} else { // interface method
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
		arity = numIn + 1
	}
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
		if fun.IsValid() {
			out = fun.Call(in)
		} else {
			out = in[0].MethodByName(name).Call(in[1:])
		}
		if numOut > 0 {
			iout := make([]interface{}, numOut)
			for i := 0; i < numOut; i++ {
				iout[i] = out[i].Interface()
			}
			if isVariadic {
				p.Ret(arity, iout...)
			} else {
				p.Ret(arity+1, iout...)
			}
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

// -----------------------------------------------------------------------------
