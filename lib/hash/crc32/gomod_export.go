// Package crc32 provide Go+ "hash/crc32" package, as "hash/crc32" package in Go.
package crc32

import (
	crc32 "hash/crc32"
	reflect "reflect"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execChecksum(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := crc32.Checksum(args[0].([]byte), args[1].(*crc32.Table))
	p.Ret(2, ret0)
}

func execChecksumIEEE(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := crc32.ChecksumIEEE(args[0].([]byte))
	p.Ret(1, ret0)
}

func execMakeTable(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := crc32.MakeTable(args[0].(uint32))
	p.Ret(1, ret0)
}

func execNew(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := crc32.New(args[0].(*crc32.Table))
	p.Ret(1, ret0)
}

func execNewIEEE(_ int, p *gop.Context) {
	ret0 := crc32.NewIEEE()
	p.Ret(0, ret0)
}

func execUpdate(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := crc32.Update(args[0].(uint32), args[1].(*crc32.Table), args[2].([]byte))
	p.Ret(3, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("hash/crc32")

func init() {
	I.RegisterFuncs(
		I.Func("Checksum", crc32.Checksum, execChecksum),
		I.Func("ChecksumIEEE", crc32.ChecksumIEEE, execChecksumIEEE),
		I.Func("MakeTable", crc32.MakeTable, execMakeTable),
		I.Func("New", crc32.New, execNew),
		I.Func("NewIEEE", crc32.NewIEEE, execNewIEEE),
		I.Func("Update", crc32.Update, execUpdate),
	)
	I.RegisterVars(
		I.Var("IEEETable", &crc32.IEEETable),
	)
	I.RegisterTypes(
		I.Type("Table", reflect.TypeOf((*crc32.Table)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("Castagnoli", qspec.Uint64, uint64(crc32.Castagnoli)),
		I.Const("IEEE", qspec.Uint64, uint64(crc32.IEEE)),
		I.Const("Koopman", qspec.Uint64, uint64(crc32.Koopman)),
		I.Const("Size", qspec.ConstUnboundInt, crc32.Size),
	)
}
