package channel

import (
	"reflect"

	"qlang.io/qlang.spec.v1"
	"qlang.io/qlang.spec.v1/types"
)

// -----------------------------------------------------------------------------

func TrySend(p *qlang.Chan, v interface{}) interface{} {

	if v == qlang.Undefined {
		panic("can't send `undefined` value to a channel")
	}
	if ok := p.Data.TrySend(reflect.ValueOf(v)); ok {
		return nil
	}
	return qlang.Undefined
}

func TryRecv(p *qlang.Chan) interface{} {

	v, _ := p.Data.TryRecv()
	if v.IsValid() {
		return v.Interface()
	}
	return qlang.Undefined
}

func Send(p *qlang.Chan, v interface{}) {

	if v == qlang.Undefined {
		panic("can't send `undefined` value to a channel")
	}
	p.Data.Send(reflect.ValueOf(v))
}

func Recv(p *qlang.Chan) interface{} {

	v, ok := p.Data.Recv()
	if ok {
		return v.Interface()
	}
	return qlang.Undefined
}

// -----------------------------------------------------------------------------

func ChanIn(ch1, val interface{}, try bool) interface{} {

	ch := ch1.(*qlang.Chan)
	if try {
		return TrySend(ch, val)
	}
	Send(ch, val)
	return nil
}

func ChanOut(ch1 interface{}, try bool) interface{} {

	ch := ch1.(*qlang.Chan)
	if try {
		return TryRecv(ch)
	}
	return Recv(ch)
}

func ChanOf(typ interface{}) interface{} {

	return reflect.ChanOf(reflect.BothDir, types.Reflect(typ))
}

func Mkchan(typ interface{}, buffer ...int) *qlang.Chan {

	n := 0
	if len(buffer) > 0 {
		n = buffer[0]
	}
	t := reflect.ChanOf(reflect.BothDir, types.Reflect(typ))
	return &qlang.Chan{reflect.MakeChan(t, n)}
}

func Close(ch1 interface{}) {

	ch := ch1.(*qlang.Chan)
	ch.Data.Close()
}

// -----------------------------------------------------------------------------

var exports = map[string]interface{}{
	"chanOf": ChanOf,
	"mkchan": Mkchan,
	"close":  Close,
}

func init() {
	qlang.ChanIn = ChanIn
	qlang.ChanOut = ChanOut
	qlang.Import("", exports)
}

// -----------------------------------------------------------------------------

