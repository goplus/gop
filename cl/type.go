package cl

import (
	"reflect"

	"github.com/qiniu/qlang/ast/astutil"
	"github.com/qiniu/qlang/exec"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

// Type represents a qlang type.
type Type interface {
}

type value interface {
	TypeOf() Type
}

// -----------------------------------------------------------------------------

type goFunc struct {
	v    *exec.GoFuncInfo
	addr uint32
	kind exec.SymbolKind
}

func (p *goFunc) TypeOf() Type {
	return reflect.TypeOf(p.v.This)
}

func getGoFunc(addr uint32, kind exec.SymbolKind) *goFunc {
	var fi *exec.GoFuncInfo
	switch kind {
	case exec.SymbolFunc:
		fi = exec.GoFuncAddr(addr).GetInfo()
	case exec.SymbolVariadicFunc:
		fi = exec.GoVariadicFuncAddr(addr).GetInfo()
	default:
		log.Fatalln("getGoFunc: unknown -", kind, addr)
	}
	return &goFunc{v: fi, addr: addr, kind: kind}
}

// -----------------------------------------------------------------------------

type constVal struct {
	v    interface{}
	kind astutil.ConstKind
}

func (p *constVal) TypeOf() Type {
	if astutil.IsConstBound(p.kind) {
		return reflect.TypeOf(p.v)
	}
	return p.kind
}

// -----------------------------------------------------------------------------
