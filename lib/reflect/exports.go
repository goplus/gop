package fmt

import (
	"reflect"

	qlang "github.com/qiniu/qlang/v6/spec"
)

func execTypeOf(zero int, p *qlang.Context) {
	args := p.GetArgs(1)
	args[0] = reflect.TypeOf(args[0])
}

func execValueOf(zero int, p *qlang.Context) {
	args := p.GetArgs(1)
	args[0] = reflect.ValueOf(args[0])
}

// -----------------------------------------------------------------------------

// I is a Go package instance.
var I = qlang.NewGoPackage("reflect")

func init() {
	I.RegisterFuncs(
		I.Func("TypeOf", reflect.TypeOf, execTypeOf),
		I.Func("ValueOf", reflect.ValueOf, execValueOf),
	)
}

// -----------------------------------------------------------------------------
