// Package quotedprintable provide Go+ "mime/quotedprintable" package, as "mime/quotedprintable" package in Go.
package quotedprintable

import (
	io "io"
	quotedprintable "mime/quotedprintable"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func toType0(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execNewReader(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := quotedprintable.NewReader(toType0(args[0]))
	p.Ret(1, ret0)
}

func toType1(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execNewWriter(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := quotedprintable.NewWriter(toType1(args[0]))
	p.Ret(1, ret0)
}

func execmReaderRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*quotedprintable.Reader).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmWriterWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*quotedprintable.Writer).Write(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmWriterClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*quotedprintable.Writer).Close()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("mime/quotedprintable")

func init() {
	I.RegisterFuncs(
		I.Func("NewReader", quotedprintable.NewReader, execNewReader),
		I.Func("NewWriter", quotedprintable.NewWriter, execNewWriter),
		I.Func("(*Reader).Read", (*quotedprintable.Reader).Read, execmReaderRead),
		I.Func("(*Writer).Write", (*quotedprintable.Writer).Write, execmWriterWrite),
		I.Func("(*Writer).Close", (*quotedprintable.Writer).Close, execmWriterClose),
	)
	I.RegisterTypes(
		I.Type("Reader", reflect.TypeOf((*quotedprintable.Reader)(nil)).Elem()),
		I.Type("Writer", reflect.TypeOf((*quotedprintable.Writer)(nil)).Elem()),
	)
}
