// Package rand provide Go+ "crypto/rand" package, as "crypto/rand" package in Go.
package rand

import (
	rand "crypto/rand"
	io "io"
	big "math/big"

	gop "github.com/goplus/gop"
)

func toType0(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execInt(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := rand.Int(toType0(args[0]), args[1].(*big.Int))
	p.Ret(2, ret0, ret1)
}

func execPrime(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := rand.Prime(toType0(args[0]), args[1].(int))
	p.Ret(2, ret0, ret1)
}

func execRead(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := rand.Read(args[0].([]byte))
	p.Ret(1, ret0, ret1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("crypto/rand")

func init() {
	I.RegisterFuncs(
		I.Func("Int", rand.Int, execInt),
		I.Func("Prime", rand.Prime, execPrime),
		I.Func("Read", rand.Read, execRead),
	)
	I.RegisterVars(
		I.Var("Reader", &rand.Reader),
	)
}
