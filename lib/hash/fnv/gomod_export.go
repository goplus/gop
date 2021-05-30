// Package fnv provide Go+ "hash/fnv" package, as "hash/fnv" package in Go.
package fnv

import (
	fnv "hash/fnv"

	gop "github.com/goplus/gop"
)

func execNew128(_ int, p *gop.Context) {
	ret0 := fnv.New128()
	p.Ret(0, ret0)
}

func execNew128a(_ int, p *gop.Context) {
	ret0 := fnv.New128a()
	p.Ret(0, ret0)
}

func execNew32(_ int, p *gop.Context) {
	ret0 := fnv.New32()
	p.Ret(0, ret0)
}

func execNew32a(_ int, p *gop.Context) {
	ret0 := fnv.New32a()
	p.Ret(0, ret0)
}

func execNew64(_ int, p *gop.Context) {
	ret0 := fnv.New64()
	p.Ret(0, ret0)
}

func execNew64a(_ int, p *gop.Context) {
	ret0 := fnv.New64a()
	p.Ret(0, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("hash/fnv")

func init() {
	I.RegisterFuncs(
		I.Func("New128", fnv.New128, execNew128),
		I.Func("New128a", fnv.New128a, execNew128a),
		I.Func("New32", fnv.New32, execNew32),
		I.Func("New32a", fnv.New32a, execNew32a),
		I.Func("New64", fnv.New64, execNew64),
		I.Func("New64a", fnv.New64a, execNew64a),
	)
}
