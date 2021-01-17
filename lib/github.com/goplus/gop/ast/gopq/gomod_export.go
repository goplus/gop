// Package gopq provide Go+ "github.com/goplus/gop/ast/gopq" package, as "github.com/goplus/gop/ast/gopq" package in Go.
package gopq

import (
	token "go/token"
	os "os"
	reflect "reflect"

	gop "github.com/goplus/gop"
	gopq "github.com/goplus/gop/ast/gopq"
	parser "github.com/goplus/gop/parser"
	token1 "github.com/goplus/gop/token"
)

func toType0(v interface{}) gopq.Node {
	if v == nil {
		return nil
	}
	return v.(gopq.Node)
}

func execNameOf(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := gopq.NameOf(toType0(args[0]))
	p.Ret(1, ret0)
}

func execNewSource(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := gopq.NewSource(args[0].(*token.FileSet), args[1].(string), args[2].(func(os.FileInfo) bool), args[3].(parser.Mode))
	p.Ret(4, ret0, ret1)
}

func toType1(v interface{}) parser.FileSystem {
	if v == nil {
		return nil
	}
	return v.(parser.FileSystem)
}

func execNewSourceFrom(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	ret0, ret1 := gopq.NewSourceFrom(args[0].(*token.FileSet), toType1(args[1]), args[2].(string), args[3].(func(os.FileInfo) bool), args[4].(parser.Mode))
	p.Ret(5, ret0, ret1)
}

func execiNodeEnd(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(gopq.Node).End()
	p.Ret(1, ret0)
}

func execiNodeForEach(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(gopq.Node).ForEach(args[1].(func(node gopq.Node) error))
	p.Ret(2, ret0)
}

func execiNodeObj(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(gopq.Node).Obj()
	p.Ret(1, ret0)
}

func execiNodePos(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(gopq.Node).Pos()
	p.Ret(1, ret0)
}

func execiNodeEnumForEach(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(gopq.NodeEnum).ForEach(args[1].(func(node gopq.Node) error))
	p.Ret(2, ret0)
}

func execmNodeSetOk(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(gopq.NodeSet).Ok()
	p.Ret(1, ret0)
}

func execmNodeSetFuncDecl(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(gopq.NodeSet).FuncDecl()
	p.Ret(1, ret0)
}

func execmNodeSetGenDecl(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(gopq.NodeSet).GenDecl(args[1].(token1.Token))
	p.Ret(2, ret0)
}

func execmNodeSetTypeSpec(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(gopq.NodeSet).TypeSpec()
	p.Ret(1, ret0)
}

func execmNodeSetVarSpec(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(gopq.NodeSet).VarSpec()
	p.Ret(1, ret0)
}

func execmNodeSetConstSpec(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(gopq.NodeSet).ConstSpec()
	p.Ret(1, ret0)
}

func execmNodeSetImportSpec(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(gopq.NodeSet).ImportSpec()
	p.Ret(1, ret0)
}

func execmNodeSetOne(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(gopq.NodeSet).One()
	p.Ret(1, ret0)
}

func execmNodeSetCache(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(gopq.NodeSet).Cache()
	p.Ret(1, ret0)
}

func execmNodeSetAny(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(gopq.NodeSet).Any()
	p.Ret(1, ret0)
}

func execmNodeSetChild(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(gopq.NodeSet).Child()
	p.Ret(1, ret0)
}

func execmNodeSetMatch(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(gopq.NodeSet).Match(args[1].(func(node gopq.Node) bool))
	p.Ret(2, ret0)
}

func execmNodeSetName(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(gopq.NodeSet).Name()
	p.Ret(1, ret0)
}

func execmNodeSetToString(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(gopq.NodeSet).ToString(args[1].(func(node gopq.Node) string))
	p.Ret(2, ret0)
}

func execmNodeSetCollect(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(gopq.NodeSet).Collect()
	p.Ret(1, ret0, ret1)
}

func execmNodeSetCollectOne(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := args[0].(gopq.NodeSet).CollectOne(gop.ToBools(args[1:])...)
	p.Ret(arity, ret0, ret1)
}

func toSlice0(args []interface{}) []gopq.Node {
	ret := make([]gopq.Node, len(args))
	for i, arg := range args {
		ret[i] = toType0(arg)
	}
	return ret
}

func execNodes(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := gopq.Nodes(toSlice0(args)...)
	p.Ret(arity, ret0)
}

func execOne(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := gopq.One(toType0(args[0]))
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("github.com/goplus/gop/ast/gopq")

func init() {
	I.RegisterFuncs(
		I.Func("NameOf", gopq.NameOf, execNameOf),
		I.Func("NewSource", gopq.NewSource, execNewSource),
		I.Func("NewSourceFrom", gopq.NewSourceFrom, execNewSourceFrom),
		I.Func("(Node).End", (gopq.Node).End, execiNodeEnd),
		I.Func("(Node).ForEach", (gopq.Node).ForEach, execiNodeForEach),
		I.Func("(Node).Obj", (gopq.Node).Obj, execiNodeObj),
		I.Func("(Node).Pos", (gopq.Node).Pos, execiNodePos),
		I.Func("(NodeEnum).ForEach", (gopq.NodeEnum).ForEach, execiNodeEnumForEach),
		I.Func("(NodeSet).Ok", (gopq.NodeSet).Ok, execmNodeSetOk),
		I.Func("(NodeSet).FuncDecl", (gopq.NodeSet).FuncDecl, execmNodeSetFuncDecl),
		I.Func("(NodeSet).GenDecl", (gopq.NodeSet).GenDecl, execmNodeSetGenDecl),
		I.Func("(NodeSet).TypeSpec", (gopq.NodeSet).TypeSpec, execmNodeSetTypeSpec),
		I.Func("(NodeSet).VarSpec", (gopq.NodeSet).VarSpec, execmNodeSetVarSpec),
		I.Func("(NodeSet).ConstSpec", (gopq.NodeSet).ConstSpec, execmNodeSetConstSpec),
		I.Func("(NodeSet).ImportSpec", (gopq.NodeSet).ImportSpec, execmNodeSetImportSpec),
		I.Func("(NodeSet).One", (gopq.NodeSet).One, execmNodeSetOne),
		I.Func("(NodeSet).Cache", (gopq.NodeSet).Cache, execmNodeSetCache),
		I.Func("(NodeSet).Any", (gopq.NodeSet).Any, execmNodeSetAny),
		I.Func("(NodeSet).Child", (gopq.NodeSet).Child, execmNodeSetChild),
		I.Func("(NodeSet).Match", (gopq.NodeSet).Match, execmNodeSetMatch),
		I.Func("(NodeSet).Name", (gopq.NodeSet).Name, execmNodeSetName),
		I.Func("(NodeSet).ToString", (gopq.NodeSet).ToString, execmNodeSetToString),
		I.Func("(NodeSet).Collect", (gopq.NodeSet).Collect, execmNodeSetCollect),
		I.Func("One", gopq.One, execOne),
	)
	I.RegisterFuncvs(
		I.Funcv("(NodeSet).CollectOne", (gopq.NodeSet).CollectOne, execmNodeSetCollectOne),
		I.Funcv("Nodes", gopq.Nodes, execNodes),
	)
	I.RegisterVars(
		I.Var("ErrBreak", &gopq.ErrBreak),
		I.Var("ErrNotFound", &gopq.ErrNotFound),
		I.Var("ErrTooManyNodes", &gopq.ErrTooManyNodes),
	)
	I.RegisterTypes(
		I.Type("Node", reflect.TypeOf((*gopq.Node)(nil)).Elem()),
		I.Type("NodeEnum", reflect.TypeOf((*gopq.NodeEnum)(nil)).Elem()),
		I.Type("NodeSet", reflect.TypeOf((*gopq.NodeSet)(nil)).Elem()),
	)
}
