// Package crc64 provide Go+ "hash/crc64" package, as "hash/crc64" package in Go.
package crc64

import (
	crc64 "hash/crc64"
	reflect "reflect"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execChecksum(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := crc64.Checksum(args[0].([]byte), args[1].(*crc64.Table))
	p.Ret(2, ret0)
}

func execMakeTable(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := crc64.MakeTable(args[0].(uint64))
	p.Ret(1, ret0)
}

func execNew(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := crc64.New(args[0].(*crc64.Table))
	p.Ret(1, ret0)
}

func execUpdate(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := crc64.Update(args[0].(uint64), args[1].(*crc64.Table), args[2].([]byte))
	p.Ret(3, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("hash/crc64")

func init() {
	I.RegisterFuncs(
		I.Func("Checksum", crc64.Checksum, execChecksum),
		I.Func("MakeTable", crc64.MakeTable, execMakeTable),
		I.Func("New", crc64.New, execNew),
		I.Func("Update", crc64.Update, execUpdate),
	)
	I.RegisterTypes(
		I.Type("Table", reflect.TypeOf((*crc64.Table)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("ECMA", qspec.Uint64, uint64(crc64.ECMA)),
		I.Const("ISO", qspec.Uint64, uint64(crc64.ISO)),
		I.Const("Size", qspec.ConstUnboundInt, crc64.Size),
	)
}
