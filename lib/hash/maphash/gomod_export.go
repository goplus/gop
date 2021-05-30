// Package maphash provide Go+ "hash/maphash" package, as "hash/maphash" package in Go.
package maphash

import (
	maphash "hash/maphash"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func execmHashWriteByte(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*maphash.Hash).WriteByte(args[1].(byte))
	p.Ret(2, ret0)
}

func execmHashWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*maphash.Hash).Write(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmHashWriteString(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*maphash.Hash).WriteString(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmHashSeed(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*maphash.Hash).Seed()
	p.Ret(1, ret0)
}

func execmHashSetSeed(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*maphash.Hash).SetSeed(args[1].(maphash.Seed))
	p.PopN(2)
}

func execmHashReset(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*maphash.Hash).Reset()
	p.PopN(1)
}

func execmHashSum64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*maphash.Hash).Sum64()
	p.Ret(1, ret0)
}

func execmHashSum(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*maphash.Hash).Sum(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmHashSize(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*maphash.Hash).Size()
	p.Ret(1, ret0)
}

func execmHashBlockSize(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*maphash.Hash).BlockSize()
	p.Ret(1, ret0)
}

func execMakeSeed(_ int, p *gop.Context) {
	ret0 := maphash.MakeSeed()
	p.Ret(0, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("hash/maphash")

func init() {
	I.RegisterFuncs(
		I.Func("(*Hash).WriteByte", (*maphash.Hash).WriteByte, execmHashWriteByte),
		I.Func("(*Hash).Write", (*maphash.Hash).Write, execmHashWrite),
		I.Func("(*Hash).WriteString", (*maphash.Hash).WriteString, execmHashWriteString),
		I.Func("(*Hash).Seed", (*maphash.Hash).Seed, execmHashSeed),
		I.Func("(*Hash).SetSeed", (*maphash.Hash).SetSeed, execmHashSetSeed),
		I.Func("(*Hash).Reset", (*maphash.Hash).Reset, execmHashReset),
		I.Func("(*Hash).Sum64", (*maphash.Hash).Sum64, execmHashSum64),
		I.Func("(*Hash).Sum", (*maphash.Hash).Sum, execmHashSum),
		I.Func("(*Hash).Size", (*maphash.Hash).Size, execmHashSize),
		I.Func("(*Hash).BlockSize", (*maphash.Hash).BlockSize, execmHashBlockSize),
		I.Func("MakeSeed", maphash.MakeSeed, execMakeSeed),
	)
	I.RegisterTypes(
		I.Type("Hash", reflect.TypeOf((*maphash.Hash)(nil)).Elem()),
		I.Type("Seed", reflect.TypeOf((*maphash.Seed)(nil)).Elem()),
	)
}
