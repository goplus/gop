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
		items[i].Name = v.name
	}
	return items
}

func makeVarsContext(vars []*Var, ctx *Context) varsContext {
	t := Struct(makeVarList(vars)).Type()
	if log.CanOutput(log.Ldebug) {
		nestDepth := ctx.getNestDepth()
		for i, v := range vars {
			if v.nestDepth != nestDepth || v.idx != uint32(i) {
				log.Panicln("makeVarsContext failed: unexpected var -", v.name[1:], v.nestDepth, nestDepth)
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
	x := ctx.vars.Field(int(idx))
	setValue(x, v)
}

func setValue(x reflect.Value, v interface{}) {
	if v != nil {
		x.Set(reflect.ValueOf(v))
		return
	}
	x.Set(reflect.Zero(x.Type()))
}

func getValueOf(v interface{}, t reflect.Type) reflect.Value {
	if v != nil {
		return reflect.ValueOf(v)
	}
	return reflect.Zero(t)
}

func getKeyOf(v interface{}, t reflect.Type) reflect.Value {
	if v != nil {
		return reflect.ValueOf(v)
	}
	return reflect.Zero(t.Key())
}

func getElementOf(v interface{}, t reflect.Type) reflect.Value {
	if v != nil {
		return reflect.ValueOf(v)
	}
	return reflect.Zero(t.Elem())
}

func getArgOf(v interface{}, t reflect.Type, i int) reflect.Value {
	if v != nil {
		return reflect.ValueOf(v)
	}
	return reflect.Zero(t.In(i))
}

func getRetOf(v interface{}, fun *FuncInfo, i int) reflect.Value {
	if v != nil {
		return reflect.ValueOf(v)
	}
	return reflect.Zero(fun.vlist[i].Type)
}

// -----------------------------------------------------------------------------
