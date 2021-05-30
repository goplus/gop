// Package binary provide Go+ "encoding/binary" package, as "encoding/binary" package in Go.
package binary

import (
	binary "encoding/binary"
	io "io"
	reflect "reflect"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execiByteOrderPutUint16(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(binary.ByteOrder).PutUint16(args[1].([]byte), args[2].(uint16))
	p.PopN(3)
}

func execiByteOrderPutUint32(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(binary.ByteOrder).PutUint32(args[1].([]byte), args[2].(uint32))
	p.PopN(3)
}

func execiByteOrderPutUint64(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(binary.ByteOrder).PutUint64(args[1].([]byte), args[2].(uint64))
	p.PopN(3)
}

func execiByteOrderString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(binary.ByteOrder).String()
	p.Ret(1, ret0)
}

func execiByteOrderUint16(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(binary.ByteOrder).Uint16(args[1].([]byte))
	p.Ret(2, ret0)
}

func execiByteOrderUint32(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(binary.ByteOrder).Uint32(args[1].([]byte))
	p.Ret(2, ret0)
}

func execiByteOrderUint64(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(binary.ByteOrder).Uint64(args[1].([]byte))
	p.Ret(2, ret0)
}

func execPutUvarint(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := binary.PutUvarint(args[0].([]byte), args[1].(uint64))
	p.Ret(2, ret0)
}

func execPutVarint(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := binary.PutVarint(args[0].([]byte), args[1].(int64))
	p.Ret(2, ret0)
}

func toType0(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func toType1(v interface{}) binary.ByteOrder {
	if v == nil {
		return nil
	}
	return v.(binary.ByteOrder)
}

func execRead(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := binary.Read(toType0(args[0]), toType1(args[1]), args[2])
	p.Ret(3, ret0)
}

func toType2(v interface{}) io.ByteReader {
	if v == nil {
		return nil
	}
	return v.(io.ByteReader)
}

func execReadUvarint(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := binary.ReadUvarint(toType2(args[0]))
	p.Ret(1, ret0, ret1)
}

func execReadVarint(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := binary.ReadVarint(toType2(args[0]))
	p.Ret(1, ret0, ret1)
}

func execSize(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := binary.Size(args[0])
	p.Ret(1, ret0)
}

func execUvarint(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := binary.Uvarint(args[0].([]byte))
	p.Ret(1, ret0, ret1)
}

func execVarint(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := binary.Varint(args[0].([]byte))
	p.Ret(1, ret0, ret1)
}

func toType3(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execWrite(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := binary.Write(toType3(args[0]), toType1(args[1]), args[2])
	p.Ret(3, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("encoding/binary")

func init() {
	I.RegisterFuncs(
		I.Func("(ByteOrder).PutUint16", (binary.ByteOrder).PutUint16, execiByteOrderPutUint16),
		I.Func("(ByteOrder).PutUint32", (binary.ByteOrder).PutUint32, execiByteOrderPutUint32),
		I.Func("(ByteOrder).PutUint64", (binary.ByteOrder).PutUint64, execiByteOrderPutUint64),
		I.Func("(ByteOrder).String", (binary.ByteOrder).String, execiByteOrderString),
		I.Func("(ByteOrder).Uint16", (binary.ByteOrder).Uint16, execiByteOrderUint16),
		I.Func("(ByteOrder).Uint32", (binary.ByteOrder).Uint32, execiByteOrderUint32),
		I.Func("(ByteOrder).Uint64", (binary.ByteOrder).Uint64, execiByteOrderUint64),
		I.Func("PutUvarint", binary.PutUvarint, execPutUvarint),
		I.Func("PutVarint", binary.PutVarint, execPutVarint),
		I.Func("Read", binary.Read, execRead),
		I.Func("ReadUvarint", binary.ReadUvarint, execReadUvarint),
		I.Func("ReadVarint", binary.ReadVarint, execReadVarint),
		I.Func("Size", binary.Size, execSize),
		I.Func("Uvarint", binary.Uvarint, execUvarint),
		I.Func("Varint", binary.Varint, execVarint),
		I.Func("Write", binary.Write, execWrite),
	)
	I.RegisterVars(
		I.Var("BigEndian", &binary.BigEndian),
		I.Var("LittleEndian", &binary.LittleEndian),
	)
	I.RegisterTypes(
		I.Type("ByteOrder", reflect.TypeOf((*binary.ByteOrder)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("MaxVarintLen16", qspec.ConstUnboundInt, binary.MaxVarintLen16),
		I.Const("MaxVarintLen32", qspec.ConstUnboundInt, binary.MaxVarintLen32),
		I.Const("MaxVarintLen64", qspec.ConstUnboundInt, binary.MaxVarintLen64),
	)
}
