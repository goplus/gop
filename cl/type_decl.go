package cl

import (
	"reflect"

	"github.com/qiniu/qlang/ast"
)

// -----------------------------------------------------------------------------

type funcTypeDecl struct {
	X *ast.FuncType
}

type methodDecl struct {
	Recv    string // recv object name
	Pointer int
	Type    *funcTypeDecl
	Body    *ast.BlockStmt
	file    *fileCtx
}

type typeDecl struct {
	Methods map[string]*methodDecl
	Alias   bool
}

type funcDecl struct {
	typ  *funcTypeDecl
	body *ast.BlockStmt
	file *fileCtx
	t    reflect.Type
}

func (p *funcDecl) typeOf() reflect.Type {
	if p.t == nil {
		p.t = p.buildType()
	}
	return p.t
}

func (p *funcDecl) buildType() reflect.Type {
	panic("todo")
}

// -----------------------------------------------------------------------------
