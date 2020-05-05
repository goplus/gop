package exec

import (
	"reflect"

	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

// A StructField describes a single field in a struct.
type StructField = reflect.StructField

// A StructInfo represents a struct information.
type StructInfo struct {
	Fields []StructField
	typ    reflect.Type
}

// Struct creates a new struct.
func Struct(fields []StructField) *StructInfo {
	return &StructInfo{
		Fields: fields,
		typ:    reflect.StructOf(fields),
	}
}

// Type returns the struct type.
func (p *StructInfo) Type() reflect.Type {
	return p.typ
}

// -----------------------------------------------------------------------------

type varsContext = reflect.Value

func makeVarList(vars []*Var) []StructField {
	items := make([]StructField, len(vars))
	for i, v := range vars {
		items[i].Type = v.Type
		items[i].Name = v.Name
	}
	return items
}

func makeVarsContext(vars []*Var, ctx *Context) varsContext {
	t := Struct(makeVarList(vars)).Type()
	if log.CanOutput(log.Ldebug) {
		nestDepth := ctx.getNestDepth()
		for i, v := range vars {
			if v.NestDepth != nestDepth || v.idx != uint32(i) {
				log.Panicln("makeVarsContext failed: unexpected variable -", v.Name)
			}
		}
	}
	return reflect.New(t).Elem()
}

func (ctx *Context) addrVar(idx uint32) interface{} {
	return ctx.vars.Field(int(idx)).Addr().Interface()
}

func (ctx *Context) getVar(idx uint32) interface{} {
	return ctx.vars.Field(int(idx)).Interface()
}

func (ctx *Context) setVar(idx uint32, v interface{}) {
	ctx.vars.Field(int(idx)).Set(reflect.ValueOf(v))
}

// -----------------------------------------------------------------------------
