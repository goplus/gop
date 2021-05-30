// Package expvar provide Go+ "expvar" package, as "expvar" package in Go.
package expvar

import (
	expvar "expvar"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func execDo(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	expvar.Do(args[0].(func(expvar.KeyValue)))
	p.PopN(1)
}

func execmFloatValue(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*expvar.Float).Value()
	p.Ret(1, ret0)
}

func execmFloatString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*expvar.Float).String()
	p.Ret(1, ret0)
}

func execmFloatAdd(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*expvar.Float).Add(args[1].(float64))
	p.PopN(2)
}

func execmFloatSet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*expvar.Float).Set(args[1].(float64))
	p.PopN(2)
}

func execmFuncValue(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(expvar.Func).Value()
	p.Ret(1, ret0)
}

func execmFuncString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(expvar.Func).String()
	p.Ret(1, ret0)
}

func execGet(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := expvar.Get(args[0].(string))
	p.Ret(1, ret0)
}

func execHandler(_ int, p *gop.Context) {
	ret0 := expvar.Handler()
	p.Ret(0, ret0)
}

func execmIntValue(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*expvar.Int).Value()
	p.Ret(1, ret0)
}

func execmIntString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*expvar.Int).String()
	p.Ret(1, ret0)
}

func execmIntAdd(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*expvar.Int).Add(args[1].(int64))
	p.PopN(2)
}

func execmIntSet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*expvar.Int).Set(args[1].(int64))
	p.PopN(2)
}

func execmMapString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*expvar.Map).String()
	p.Ret(1, ret0)
}

func execmMapInit(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*expvar.Map).Init()
	p.Ret(1, ret0)
}

func execmMapGet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*expvar.Map).Get(args[1].(string))
	p.Ret(2, ret0)
}

func toType0(v interface{}) expvar.Var {
	if v == nil {
		return nil
	}
	return v.(expvar.Var)
}

func execmMapSet(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*expvar.Map).Set(args[1].(string), toType0(args[2]))
	p.PopN(3)
}

func execmMapAdd(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*expvar.Map).Add(args[1].(string), args[2].(int64))
	p.PopN(3)
}

func execmMapAddFloat(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*expvar.Map).AddFloat(args[1].(string), args[2].(float64))
	p.PopN(3)
}

func execmMapDelete(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*expvar.Map).Delete(args[1].(string))
	p.PopN(2)
}

func execmMapDo(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*expvar.Map).Do(args[1].(func(expvar.KeyValue)))
	p.PopN(2)
}

func execNewFloat(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := expvar.NewFloat(args[0].(string))
	p.Ret(1, ret0)
}

func execNewInt(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := expvar.NewInt(args[0].(string))
	p.Ret(1, ret0)
}

func execNewMap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := expvar.NewMap(args[0].(string))
	p.Ret(1, ret0)
}

func execNewString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := expvar.NewString(args[0].(string))
	p.Ret(1, ret0)
}

func execPublish(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	expvar.Publish(args[0].(string), toType0(args[1]))
	p.PopN(2)
}

func execmStringValue(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*expvar.String).Value()
	p.Ret(1, ret0)
}

func execmStringString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*expvar.String).String()
	p.Ret(1, ret0)
}

func execmStringSet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*expvar.String).Set(args[1].(string))
	p.PopN(2)
}

func execiVarString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(expvar.Var).String()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("expvar")

func init() {
	I.RegisterFuncs(
		I.Func("Do", expvar.Do, execDo),
		I.Func("(*Float).Value", (*expvar.Float).Value, execmFloatValue),
		I.Func("(*Float).String", (*expvar.Float).String, execmFloatString),
		I.Func("(*Float).Add", (*expvar.Float).Add, execmFloatAdd),
		I.Func("(*Float).Set", (*expvar.Float).Set, execmFloatSet),
		I.Func("(Func).Value", (expvar.Func).Value, execmFuncValue),
		I.Func("(Func).String", (expvar.Func).String, execmFuncString),
		I.Func("Get", expvar.Get, execGet),
		I.Func("Handler", expvar.Handler, execHandler),
		I.Func("(*Int).Value", (*expvar.Int).Value, execmIntValue),
		I.Func("(*Int).String", (*expvar.Int).String, execmIntString),
		I.Func("(*Int).Add", (*expvar.Int).Add, execmIntAdd),
		I.Func("(*Int).Set", (*expvar.Int).Set, execmIntSet),
		I.Func("(*Map).String", (*expvar.Map).String, execmMapString),
		I.Func("(*Map).Init", (*expvar.Map).Init, execmMapInit),
		I.Func("(*Map).Get", (*expvar.Map).Get, execmMapGet),
		I.Func("(*Map).Set", (*expvar.Map).Set, execmMapSet),
		I.Func("(*Map).Add", (*expvar.Map).Add, execmMapAdd),
		I.Func("(*Map).AddFloat", (*expvar.Map).AddFloat, execmMapAddFloat),
		I.Func("(*Map).Delete", (*expvar.Map).Delete, execmMapDelete),
		I.Func("(*Map).Do", (*expvar.Map).Do, execmMapDo),
		I.Func("NewFloat", expvar.NewFloat, execNewFloat),
		I.Func("NewInt", expvar.NewInt, execNewInt),
		I.Func("NewMap", expvar.NewMap, execNewMap),
		I.Func("NewString", expvar.NewString, execNewString),
		I.Func("Publish", expvar.Publish, execPublish),
		I.Func("(*String).Value", (*expvar.String).Value, execmStringValue),
		I.Func("(*String).String", (*expvar.String).String, execmStringString),
		I.Func("(*String).Set", (*expvar.String).Set, execmStringSet),
		I.Func("(Var).String", (expvar.Var).String, execiVarString),
	)
	I.RegisterTypes(
		I.Type("Float", reflect.TypeOf((*expvar.Float)(nil)).Elem()),
		I.Type("Func", reflect.TypeOf((*expvar.Func)(nil)).Elem()),
		I.Type("Int", reflect.TypeOf((*expvar.Int)(nil)).Elem()),
		I.Type("KeyValue", reflect.TypeOf((*expvar.KeyValue)(nil)).Elem()),
		I.Type("Map", reflect.TypeOf((*expvar.Map)(nil)).Elem()),
		I.Type("String", reflect.TypeOf((*expvar.String)(nil)).Elem()),
		I.Type("Var", reflect.TypeOf((*expvar.Var)(nil)).Elem()),
	)
}
