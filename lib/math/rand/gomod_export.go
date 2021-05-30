// Package rand provide Go+ "math/rand" package, as "math/rand" package in Go.
package rand

import (
	rand "math/rand"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func execExpFloat64(_ int, p *gop.Context) {
	ret0 := rand.ExpFloat64()
	p.Ret(0, ret0)
}

func execFloat32(_ int, p *gop.Context) {
	ret0 := rand.Float32()
	p.Ret(0, ret0)
}

func execFloat64(_ int, p *gop.Context) {
	ret0 := rand.Float64()
	p.Ret(0, ret0)
}

func execInt(_ int, p *gop.Context) {
	ret0 := rand.Int()
	p.Ret(0, ret0)
}

func execInt31(_ int, p *gop.Context) {
	ret0 := rand.Int31()
	p.Ret(0, ret0)
}

func execInt31n(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := rand.Int31n(args[0].(int32))
	p.Ret(1, ret0)
}

func execInt63(_ int, p *gop.Context) {
	ret0 := rand.Int63()
	p.Ret(0, ret0)
}

func execInt63n(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := rand.Int63n(args[0].(int64))
	p.Ret(1, ret0)
}

func execIntn(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := rand.Intn(args[0].(int))
	p.Ret(1, ret0)
}

func toType0(v interface{}) rand.Source {
	if v == nil {
		return nil
	}
	return v.(rand.Source)
}

func execNew(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := rand.New(toType0(args[0]))
	p.Ret(1, ret0)
}

func execNewSource(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := rand.NewSource(args[0].(int64))
	p.Ret(1, ret0)
}

func execNewZipf(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := rand.NewZipf(args[0].(*rand.Rand), args[1].(float64), args[2].(float64), args[3].(uint64))
	p.Ret(4, ret0)
}

func execNormFloat64(_ int, p *gop.Context) {
	ret0 := rand.NormFloat64()
	p.Ret(0, ret0)
}

func execPerm(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := rand.Perm(args[0].(int))
	p.Ret(1, ret0)
}

func execmRandExpFloat64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*rand.Rand).ExpFloat64()
	p.Ret(1, ret0)
}

func execmRandNormFloat64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*rand.Rand).NormFloat64()
	p.Ret(1, ret0)
}

func execmRandSeed(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*rand.Rand).Seed(args[1].(int64))
	p.PopN(2)
}

func execmRandInt63(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*rand.Rand).Int63()
	p.Ret(1, ret0)
}

func execmRandUint32(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*rand.Rand).Uint32()
	p.Ret(1, ret0)
}

func execmRandUint64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*rand.Rand).Uint64()
	p.Ret(1, ret0)
}

func execmRandInt31(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*rand.Rand).Int31()
	p.Ret(1, ret0)
}

func execmRandInt(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*rand.Rand).Int()
	p.Ret(1, ret0)
}

func execmRandInt63n(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*rand.Rand).Int63n(args[1].(int64))
	p.Ret(2, ret0)
}

func execmRandInt31n(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*rand.Rand).Int31n(args[1].(int32))
	p.Ret(2, ret0)
}

func execmRandIntn(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*rand.Rand).Intn(args[1].(int))
	p.Ret(2, ret0)
}

func execmRandFloat64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*rand.Rand).Float64()
	p.Ret(1, ret0)
}

func execmRandFloat32(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*rand.Rand).Float32()
	p.Ret(1, ret0)
}

func execmRandPerm(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*rand.Rand).Perm(args[1].(int))
	p.Ret(2, ret0)
}

func execmRandShuffle(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*rand.Rand).Shuffle(args[1].(int), args[2].(func(i int, j int)))
	p.PopN(3)
}

func execmRandRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*rand.Rand).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execRead(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := rand.Read(args[0].([]byte))
	p.Ret(1, ret0, ret1)
}

func execSeed(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	rand.Seed(args[0].(int64))
	p.PopN(1)
}

func execShuffle(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	rand.Shuffle(args[0].(int), args[1].(func(i int, j int)))
	p.PopN(2)
}

func execiSourceInt63(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(rand.Source).Int63()
	p.Ret(1, ret0)
}

func execiSourceSeed(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(rand.Source).Seed(args[1].(int64))
	p.PopN(2)
}

func execiSource64Uint64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(rand.Source64).Uint64()
	p.Ret(1, ret0)
}

func execUint32(_ int, p *gop.Context) {
	ret0 := rand.Uint32()
	p.Ret(0, ret0)
}

func execUint64(_ int, p *gop.Context) {
	ret0 := rand.Uint64()
	p.Ret(0, ret0)
}

func execmZipfUint64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*rand.Zipf).Uint64()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("math/rand")

func init() {
	I.RegisterFuncs(
		I.Func("ExpFloat64", rand.ExpFloat64, execExpFloat64),
		I.Func("Float32", rand.Float32, execFloat32),
		I.Func("Float64", rand.Float64, execFloat64),
		I.Func("Int", rand.Int, execInt),
		I.Func("Int31", rand.Int31, execInt31),
		I.Func("Int31n", rand.Int31n, execInt31n),
		I.Func("Int63", rand.Int63, execInt63),
		I.Func("Int63n", rand.Int63n, execInt63n),
		I.Func("Intn", rand.Intn, execIntn),
		I.Func("New", rand.New, execNew),
		I.Func("NewSource", rand.NewSource, execNewSource),
		I.Func("NewZipf", rand.NewZipf, execNewZipf),
		I.Func("NormFloat64", rand.NormFloat64, execNormFloat64),
		I.Func("Perm", rand.Perm, execPerm),
		I.Func("(*Rand).ExpFloat64", (*rand.Rand).ExpFloat64, execmRandExpFloat64),
		I.Func("(*Rand).NormFloat64", (*rand.Rand).NormFloat64, execmRandNormFloat64),
		I.Func("(*Rand).Seed", (*rand.Rand).Seed, execmRandSeed),
		I.Func("(*Rand).Int63", (*rand.Rand).Int63, execmRandInt63),
		I.Func("(*Rand).Uint32", (*rand.Rand).Uint32, execmRandUint32),
		I.Func("(*Rand).Uint64", (*rand.Rand).Uint64, execmRandUint64),
		I.Func("(*Rand).Int31", (*rand.Rand).Int31, execmRandInt31),
		I.Func("(*Rand).Int", (*rand.Rand).Int, execmRandInt),
		I.Func("(*Rand).Int63n", (*rand.Rand).Int63n, execmRandInt63n),
		I.Func("(*Rand).Int31n", (*rand.Rand).Int31n, execmRandInt31n),
		I.Func("(*Rand).Intn", (*rand.Rand).Intn, execmRandIntn),
		I.Func("(*Rand).Float64", (*rand.Rand).Float64, execmRandFloat64),
		I.Func("(*Rand).Float32", (*rand.Rand).Float32, execmRandFloat32),
		I.Func("(*Rand).Perm", (*rand.Rand).Perm, execmRandPerm),
		I.Func("(*Rand).Shuffle", (*rand.Rand).Shuffle, execmRandShuffle),
		I.Func("(*Rand).Read", (*rand.Rand).Read, execmRandRead),
		I.Func("Read", rand.Read, execRead),
		I.Func("Seed", rand.Seed, execSeed),
		I.Func("Shuffle", rand.Shuffle, execShuffle),
		I.Func("(Source).Int63", (rand.Source).Int63, execiSourceInt63),
		I.Func("(Source).Seed", (rand.Source).Seed, execiSourceSeed),
		I.Func("(Source64).Int63", (rand.Source).Int63, execiSourceInt63),
		I.Func("(Source64).Seed", (rand.Source).Seed, execiSourceSeed),
		I.Func("(Source64).Uint64", (rand.Source64).Uint64, execiSource64Uint64),
		I.Func("Uint32", rand.Uint32, execUint32),
		I.Func("Uint64", rand.Uint64, execUint64),
		I.Func("(*Zipf).Uint64", (*rand.Zipf).Uint64, execmZipfUint64),
	)
	I.RegisterTypes(
		I.Type("Rand", reflect.TypeOf((*rand.Rand)(nil)).Elem()),
		I.Type("Source", reflect.TypeOf((*rand.Source)(nil)).Elem()),
		I.Type("Source64", reflect.TypeOf((*rand.Source64)(nil)).Elem()),
		I.Type("Zipf", reflect.TypeOf((*rand.Zipf)(nil)).Elem()),
	)
}
