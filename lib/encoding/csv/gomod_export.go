// Package csv provide Go+ "encoding/csv" package, as "encoding/csv" package in Go.
package csv

import (
	csv "encoding/csv"
	io "io"
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
	ret0 := csv.NewReader(toType0(args[0]))
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
	ret0 := csv.NewWriter(toType1(args[0]))
	p.Ret(1, ret0)
}

func execmParseErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*csv.ParseError).Error()
	p.Ret(1, ret0)
}

func execmParseErrorUnwrap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*csv.ParseError).Unwrap()
	p.Ret(1, ret0)
}

func execmReaderRead(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*csv.Reader).Read()
	p.Ret(1, ret0, ret1)
}

func execmReaderReadAll(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*csv.Reader).ReadAll()
	p.Ret(1, ret0, ret1)
}

func execmWriterWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*csv.Writer).Write(args[1].([]string))
	p.Ret(2, ret0)
}

func execmWriterFlush(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*csv.Writer).Flush()
	p.PopN(1)
}

func execmWriterError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*csv.Writer).Error()
	p.Ret(1, ret0)
}

func execmWriterWriteAll(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*csv.Writer).WriteAll(args[1].([][]string))
	p.Ret(2, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("encoding/csv")

func init() {
	I.RegisterFuncs(
		I.Func("NewReader", csv.NewReader, execNewReader),
		I.Func("NewWriter", csv.NewWriter, execNewWriter),
		I.Func("(*ParseError).Error", (*csv.ParseError).Error, execmParseErrorError),
		I.Func("(*ParseError).Unwrap", (*csv.ParseError).Unwrap, execmParseErrorUnwrap),
		I.Func("(*Reader).Read", (*csv.Reader).Read, execmReaderRead),
		I.Func("(*Reader).ReadAll", (*csv.Reader).ReadAll, execmReaderReadAll),
		I.Func("(*Writer).Write", (*csv.Writer).Write, execmWriterWrite),
		I.Func("(*Writer).Flush", (*csv.Writer).Flush, execmWriterFlush),
		I.Func("(*Writer).Error", (*csv.Writer).Error, execmWriterError),
		I.Func("(*Writer).WriteAll", (*csv.Writer).WriteAll, execmWriterWriteAll),
	)
	I.RegisterVars(
		I.Var("ErrBareQuote", &csv.ErrBareQuote),
		I.Var("ErrFieldCount", &csv.ErrFieldCount),
		I.Var("ErrQuote", &csv.ErrQuote),
		I.Var("ErrTrailingComma", &csv.ErrTrailingComma),
	)
	I.RegisterTypes(
		I.Type("ParseError", reflect.TypeOf((*csv.ParseError)(nil)).Elem()),
		I.Type("Reader", reflect.TypeOf((*csv.Reader)(nil)).Elem()),
		I.Type("Writer", reflect.TypeOf((*csv.Writer)(nil)).Elem()),
	)
}
