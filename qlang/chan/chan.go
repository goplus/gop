package channel

import (
	"reflect"

	"qlang.io/qlang/builtin/types"
	"qlang.io/qlang.spec.v1"
)

// -----------------------------------------------------------------------------

func TrySend(p *types.Chan, v interface{}) interface{} {

	if v == qlang.Undefined {
		panic("can't send `undefined` value to a channel")
	}
	if ok := p.Data.TrySend(reflect.ValueOf(v)); ok {
		return nil
	}
	return qlang.Undefined
}

func TryRecv(p *types.Chan) interface{} {

	v, _ := p.Data.TryRecv()
	if v.IsValid() {
		return v.Interface()
	}
	return qlang.Undefined
}

func Send(p *types.Chan, v interface{}) {

	if v == qlang.Undefined {
		panic("can't send `undefined` value to a channel")
	}
	p.Data.Send(reflect.ValueOf(v))
}

func Recv(p *types.Chan) interface{} {

	v, ok := p.Data.Recv()
	if ok {
		return v.Interface()
	}
	return qlang.Undefined
}

// -----------------------------------------------------------------------------

func ChanIn(ch1, val interface{}, try bool) interface{} {

	ch := ch1.(*types.Chan)
	if try {
		return TrySend(ch, val)
	}
	Send(ch, val)
	return nil
}

func ChanOut(ch1 interface{}, try bool) interface{} {

	ch := ch1.(*types.Chan)
	if try {
		return TryRecv(ch)
	}
	return Recv(ch)
}

func ChanOf(typ interface{}) interface{} {

	return reflect.ChanOf(reflect.BothDir, types.Reflect(typ))
}

func Mkchan(typ interface{}, buffer ...int) *types.Chan {

	n := 0
	if len(buffer) > 0 {
		n = buffer[0]
	}
	t := reflect.ChanOf(reflect.BothDir, types.Reflect(typ))
	return &types.Chan{reflect.MakeChan(t, n)}
}

// -----------------------------------------------------------------------------

var exports = map[string]interface{}{
	"chanOf": ChanOf,
	"mkchan": Mkchan,
}

func init() {
	qlang.ChanIn = ChanIn
	qlang.ChanOut = ChanOut
	qlang.Import("", exports)
}

// -----------------------------------------------------------------------------

