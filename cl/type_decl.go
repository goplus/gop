package cl

import (
	"reflect"
	"strconv"

	"github.com/qiniu/qlang/ast"
	"github.com/qiniu/qlang/exec"
	"github.com/qiniu/x/log"
)

type iType = reflect.Type

func toTypeEx(ctx *blockCtx, typ ast.Expr) (t iType, variadic bool) {
	if t = toType(ctx, typ); t != nil {
		return
	}
	if v, ok := typ.(*ast.Ellipsis); ok {
		elem := toType(ctx, v.Elt)
		return reflect.SliceOf(elem), true
	}
	log.Panicln("toType: unknown -", reflect.TypeOf(typ))
	return nil, false
}

func toType(ctx *blockCtx, typ ast.Expr) iType {
	switch v := typ.(type) {
	case *ast.Ident:
		return toIdentType(ctx, v.Name)
	case *ast.SelectorExpr:
		return toExternalType(ctx, v)
	case *ast.StarExpr:
		elem := toType(ctx, v.X)
		return reflect.PtrTo(elem)
	case *ast.ArrayType:
		return toArrayType(ctx, v)
	case *ast.FuncType:
		return toFuncType(ctx, v)
	case *ast.InterfaceType:
		return toInterfaceType(ctx, v)
	case *ast.StructType:
		return toStructType(ctx, v)
	case *ast.MapType:
		key := toType(ctx, v.Key)
		elem := toType(ctx, v.Value)
		return reflect.MapOf(key, elem)
	case *ast.ChanType:
		val := toType(ctx, v.Value)
		return reflect.ChanOf(toChanDir(v.Dir), val)
	}
	return nil
}

func toChanDir(dir ast.ChanDir) (ret reflect.ChanDir) {
	if (dir & ast.SEND) != 0 {
		ret |= reflect.SendDir
	}
	if (dir & ast.RECV) != 0 {
		ret |= reflect.RecvDir
	}
	return
}

func toFuncType(ctx *blockCtx, t *ast.FuncType) iType {
	in, variadic := toTypesEx(ctx, t.Params)
	out := toTypes(ctx, t.Results)
	return reflect.FuncOf(in, out, variadic)
}

func buildFuncType(fi *exec.FuncInfo, ctx *blockCtx, t *ast.FuncType) {
	in, args, variadic := toArgTypes(ctx, t.Params)
	rets := toReturnTypes(ctx, t.Results)
	if variadic {
		fi.Vargs(in...)
	} else {
		fi.Args(in...)
	}
	fi.Return(rets...)
	ctx.insertFuncVars(in, args, rets)
}

func toTypes(ctx *blockCtx, fields *ast.FieldList) (types []iType) {
	if fields == nil {
		return
	}
	for _, field := range fields.List {
		n := len(field.Names)
		if n == 0 {
			n = 1
		}
		typ := toType(ctx, field.Type)
		if typ == nil {
			log.Panicln("toType: unknown -", reflect.TypeOf(field.Type))
		}
		for i := 0; i < n; i++ {
			types = append(types, typ)
		}
	}
	return
}

func toTypesEx(ctx *blockCtx, fields *ast.FieldList) ([]iType, bool) {
	var types []iType
	last := len(fields.List) - 1
	for i := 0; i <= last; i++ {
		field := fields.List[i]
		n := len(field.Names)
		if n == 0 {
			n = 1
		}
		typ, variadic := toTypeEx(ctx, field.Type)
		for i := 0; i < n; i++ {
			types = append(types, typ)
		}
		if variadic {
			if i != last {
				log.Panicln("toTypes failed: the variadic type isn't last argument?")
			}
			return types, true
		}
	}
	return types, false
}

func toReturnTypes(ctx *blockCtx, fields *ast.FieldList) (vars []*exec.Var) {
	if fields == nil {
		return
	}
	index := 0
	for _, field := range fields.List {
		n := len(field.Names)
		typ := toType(ctx, field.Type)
		if typ == nil {
			log.Panicln("toType: unknown -", reflect.TypeOf(field.Type))
		}
		if n == 0 {
			index++
			vars = append(vars, exec.NewVar(typ, strconv.Itoa(index)))
		} else {
			for i := 0; i < n; i++ {
				vars = append(vars, exec.NewVar(typ, field.Names[i].Name))
			}
			index += n
		}
	}
	return
}

func toArgTypes(ctx *blockCtx, fields *ast.FieldList) ([]iType, []string, bool) {
	var types []iType
	var names []string
	last := len(fields.List) - 1
	for i := 0; i <= last; i++ {
		field := fields.List[i]
		n := len(field.Names)
		if n == 0 {
			names = append(names, "")
			n = 1
		} else {
			for _, fld := range field.Names {
				names = append(names, fld.Name)
			}
		}
		typ, variadic := toTypeEx(ctx, field.Type)
		for i := 0; i < n; i++ {
			types = append(types, typ)
		}
		if variadic {
			if i != last {
				log.Panicln("toTypes failed: the variadic type isn't last argument?")
			}
			return types, names, true
		}
	}
	return types, names, false
}

func toStructType(ctx *blockCtx, v *ast.StructType) iType {
	panic("toStructType: todo")
}

func toInterfaceType(ctx *blockCtx, v *ast.InterfaceType) iType {
	panic("toInterfaceType: todo")
}

func toExternalType(ctx *blockCtx, v *ast.SelectorExpr) iType {
	panic("toExternalType: todo")
}

func toIdentType(ctx *blockCtx, ident string) iType {
	if typ, ok := ctx.builtin.FindType(ident); ok {
		return typ
	}
	log.Panicln("toIdentType failed: unknown ident -", ident)
	return nil
}

func toArrayType(ctx *blockCtx, v *ast.ArrayType) iType {
	panic("toArrayType: todo")
}

// -----------------------------------------------------------------------------

type typeDecl struct {
	Methods map[string]*methodDecl
	Alias   bool
}

type methodDecl struct {
	recv    string // recv object name
	pointer int
	typ     *ast.FuncType
	body    *ast.BlockStmt
	file    *fileCtx
}

type funcDecl struct {
	typ  *ast.FuncType
	body *ast.BlockStmt
	ctx  *blockCtx
	fi   *exec.FuncInfo
	t    reflect.Type
	used bool
}

func newFuncDecl(name string, typ *ast.FuncType, body *ast.BlockStmt, ctx *blockCtx) *funcDecl {
	fi := exec.NewFunc(name, 1)
	return &funcDecl{typ: typ, body: body, ctx: ctx, fi: fi}
}

func (p *funcDecl) getFuncInfo() *exec.FuncInfo {
	if !p.fi.IsTypeValid() {
		buildFuncType(p.fi, p.ctx, p.typ)
	}
	return p.fi
}

func (p *funcDecl) typeOf() iType {
	if p.t == nil {
		p.t = p.getFuncInfo().Type()
	}
	return p.t
}

func (p *funcDecl) compile() {
	fun := p.getFuncInfo()
	ctx := p.ctx
	out := ctx.out
	out.DefineFunc(fun)
	compileBlockStmt(ctx, p.body)
	out.EndFunc(fun)
}

// -----------------------------------------------------------------------------
