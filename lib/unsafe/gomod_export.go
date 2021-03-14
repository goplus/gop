// Package unsafe provide Go+ "unsafe" package, as "unsafe" package in Go.
package unsafe

import (
	"reflect"
	"unsafe"

	gop "github.com/goplus/gop"
)

func execSizeof(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := reflect.TypeOf(args[0]).Size()
	p.Ret(1, ret0)
}

func execAlignof(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := uintptr(reflect.TypeOf(args[0]).Align())
	p.Ret(1, ret0)
}

func execOffsetof(_ int, p *gop.Context) {
	p.Ret(1, 0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("unsafe")

func sizeof(any interface{}) uintptr {
	return reflect.TypeOf(any).Size()
}

func alignof(any interface{}) uintptr {
	return uintptr(reflect.TypeOf(any).Align())
}

func offsetof(any interface{}) uintptr {
	return 0
}

func init() {
	I.RegisterFuncs(
		I.Func("Sizeof", sizeof, execSizeof),
		I.Func("Alignof", alignof, execAlignof),
		I.Func("Offsetof", offsetof, execOffsetof),
	)
	I.RegisterTypes(
		I.Type("Pointer", reflect.TypeOf((*unsafe.Pointer)(nil)).Elem()),
	)
}
