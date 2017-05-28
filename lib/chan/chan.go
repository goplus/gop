package channel

import (
	"fmt"
	"reflect"

	qlang "qlang.io/spec"
	"qlang.io/spec/types"
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
	return &qlang.Chan{Data: reflect.MakeChan(t, n)}
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

func makeChan(typ reflect.Type, buffer ...int) interface{} {

	n := 0
	if len(buffer) > 0 {
		n = buffer[0]
	}
	return &qlang.Chan{Data: reflect.MakeChan(typ, n)}
}

func chanOf(elem interface{}) interface{} {

	if t, ok := elem.(qlang.GoTyper); ok {
		tchan := reflect.ChanOf(reflect.BothDir, t.GoType())
		return qlang.NewType(tchan)
	}
	panic(fmt.Sprintf("invalid chan T: `%v` isn't a qlang type", elem))
}

func init() {
	qlang.ChanIn = ChanIn
	qlang.ChanOut = ChanOut
	qlang.MakeChan = makeChan
	qlang.ChanOf = chanOf
	qlang.Import("", exports)
}

// -----------------------------------------------------------------------------
