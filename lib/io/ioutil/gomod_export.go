// Package ioutil provide Go+ "io/ioutil" package, as "io/ioutil" package in Go.
package ioutil

import (
	io "io"
	ioutil "io/ioutil"
	os "os"

	gop "github.com/qiniu/goplus/gop"
)

func toType0(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execNopCloser(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := ioutil.NopCloser(toType0(args[0]))
	p.Ret(1, ret0)
}

func execReadAll(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := ioutil.ReadAll(toType0(args[0]))
	p.Ret(1, ret0, ret1)
}

func execReadDir(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := ioutil.ReadDir(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execReadFile(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := ioutil.ReadFile(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execTempDir(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := ioutil.TempDir(args[0].(string), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execTempFile(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := ioutil.TempFile(args[0].(string), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execWriteFile(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := ioutil.WriteFile(args[0].(string), args[1].([]byte), args[2].(os.FileMode))
	p.Ret(3, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("io/ioutil")

func init() {
	I.RegisterFuncs(
		I.Func("NopCloser", ioutil.NopCloser, execNopCloser),
		I.Func("ReadAll", ioutil.ReadAll, execReadAll),
		I.Func("ReadDir", ioutil.ReadDir, execReadDir),
		I.Func("ReadFile", ioutil.ReadFile, execReadFile),
		I.Func("TempDir", ioutil.TempDir, execTempDir),
		I.Func("TempFile", ioutil.TempFile, execTempFile),
		I.Func("WriteFile", ioutil.WriteFile, execWriteFile),
	)
}
