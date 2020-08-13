package bytecode

import (
	"reflect"

	"github.com/goplus/gop/exec.spec"
)

type Type struct {
	name string
	typ  reflect.Type

	off uint32
}

// NewType creates a type instance.
func NewType(typ reflect.Type, name string) *Type {
	return &Type{typ: typ, name: name}
}

// DefineVar defines types.
func (p *Builder) DefineType(typ exec.Type) *Builder {
	t := typ.(*Type)
	t.off = p.newType(typ.Type())
	return p
}

func (p *Type) Type() reflect.Type {
	return p.typ
}

// IsUnnamedOut returns if variable unnamed or not.
func (p *Type) IsUnnamedOut() bool {
	c := p.name[0]
	return c >= '0' && c <= '9'
}

func (p *Type) Name() string {
	return p.name
}
