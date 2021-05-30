// Package bufio provide Go+ "bufio" package, as "bufio" package in Go.
package bufio

import (
	bufio "bufio"
	io "io"
	reflect "reflect"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execNewReadWriter(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bufio.NewReadWriter(args[0].(*bufio.Reader), args[1].(*bufio.Writer))
	p.Ret(2, ret0)
}

func toType0(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execNewReader(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bufio.NewReader(toType0(args[0]))
	p.Ret(1, ret0)
}

func execNewReaderSize(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bufio.NewReaderSize(toType0(args[0]), args[1].(int))
	p.Ret(2, ret0)
}

func execNewScanner(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bufio.NewScanner(toType0(args[0]))
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
	ret0 := bufio.NewWriter(toType1(args[0]))
	p.Ret(1, ret0)
}

func execNewWriterSize(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bufio.NewWriterSize(toType1(args[0]), args[1].(int))
	p.Ret(2, ret0)
}

func execmReaderSize(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*bufio.Reader).Size()
	p.Ret(1, ret0)
}

func execmReaderReset(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*bufio.Reader).Reset(toType0(args[1]))
	p.PopN(2)
}

func execmReaderPeek(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*bufio.Reader).Peek(args[1].(int))
	p.Ret(2, ret0, ret1)
}

func execmReaderDiscard(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*bufio.Reader).Discard(args[1].(int))
	p.Ret(2, ret0, ret1)
}

func execmReaderRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*bufio.Reader).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmReaderReadByte(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*bufio.Reader).ReadByte()
	p.Ret(1, ret0, ret1)
}

func execmReaderUnreadByte(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*bufio.Reader).UnreadByte()
	p.Ret(1, ret0)
}

func execmReaderReadRune(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := args[0].(*bufio.Reader).ReadRune()
	p.Ret(1, ret0, ret1, ret2)
}

func execmReaderUnreadRune(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*bufio.Reader).UnreadRune()
	p.Ret(1, ret0)
}

func execmReaderBuffered(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*bufio.Reader).Buffered()
	p.Ret(1, ret0)
}

func execmReaderReadSlice(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*bufio.Reader).ReadSlice(args[1].(byte))
	p.Ret(2, ret0, ret1)
}

func execmReaderReadLine(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := args[0].(*bufio.Reader).ReadLine()
	p.Ret(1, ret0, ret1, ret2)
}

func execmReaderReadBytes(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*bufio.Reader).ReadBytes(args[1].(byte))
	p.Ret(2, ret0, ret1)
}

func execmReaderReadString(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*bufio.Reader).ReadString(args[1].(byte))
	p.Ret(2, ret0, ret1)
}

func execmReaderWriteTo(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*bufio.Reader).WriteTo(toType1(args[1]))
	p.Ret(2, ret0, ret1)
}

func execScanBytes(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1, ret2 := bufio.ScanBytes(args[0].([]byte), args[1].(bool))
	p.Ret(2, ret0, ret1, ret2)
}

func execScanLines(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1, ret2 := bufio.ScanLines(args[0].([]byte), args[1].(bool))
	p.Ret(2, ret0, ret1, ret2)
}

func execScanRunes(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1, ret2 := bufio.ScanRunes(args[0].([]byte), args[1].(bool))
	p.Ret(2, ret0, ret1, ret2)
}

func execScanWords(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1, ret2 := bufio.ScanWords(args[0].([]byte), args[1].(bool))
	p.Ret(2, ret0, ret1, ret2)
}

func execmScannerErr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*bufio.Scanner).Err()
	p.Ret(1, ret0)
}

func execmScannerBytes(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*bufio.Scanner).Bytes()
	p.Ret(1, ret0)
}

func execmScannerText(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*bufio.Scanner).Text()
	p.Ret(1, ret0)
}

func execmScannerScan(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*bufio.Scanner).Scan()
	p.Ret(1, ret0)
}

func execmScannerBuffer(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*bufio.Scanner).Buffer(args[1].([]byte), args[2].(int))
	p.PopN(3)
}

func execmScannerSplit(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*bufio.Scanner).Split(args[1].(bufio.SplitFunc))
	p.PopN(2)
}

func execmWriterSize(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*bufio.Writer).Size()
	p.Ret(1, ret0)
}

func execmWriterReset(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*bufio.Writer).Reset(toType1(args[1]))
	p.PopN(2)
}

func execmWriterFlush(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*bufio.Writer).Flush()
	p.Ret(1, ret0)
}

func execmWriterAvailable(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*bufio.Writer).Available()
	p.Ret(1, ret0)
}

func execmWriterBuffered(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*bufio.Writer).Buffered()
	p.Ret(1, ret0)
}

func execmWriterWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*bufio.Writer).Write(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmWriterWriteByte(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*bufio.Writer).WriteByte(args[1].(byte))
	p.Ret(2, ret0)
}

func execmWriterWriteRune(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*bufio.Writer).WriteRune(args[1].(rune))
	p.Ret(2, ret0, ret1)
}

func execmWriterWriteString(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*bufio.Writer).WriteString(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmWriterReadFrom(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*bufio.Writer).ReadFrom(toType0(args[1]))
	p.Ret(2, ret0, ret1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("bufio")

func init() {
	I.RegisterFuncs(
		I.Func("NewReadWriter", bufio.NewReadWriter, execNewReadWriter),
		I.Func("NewReader", bufio.NewReader, execNewReader),
		I.Func("NewReaderSize", bufio.NewReaderSize, execNewReaderSize),
		I.Func("NewScanner", bufio.NewScanner, execNewScanner),
		I.Func("NewWriter", bufio.NewWriter, execNewWriter),
		I.Func("NewWriterSize", bufio.NewWriterSize, execNewWriterSize),
		I.Func("(*Reader).Size", (*bufio.Reader).Size, execmReaderSize),
		I.Func("(*Reader).Reset", (*bufio.Reader).Reset, execmReaderReset),
		I.Func("(*Reader).Peek", (*bufio.Reader).Peek, execmReaderPeek),
		I.Func("(*Reader).Discard", (*bufio.Reader).Discard, execmReaderDiscard),
		I.Func("(*Reader).Read", (*bufio.Reader).Read, execmReaderRead),
		I.Func("(*Reader).ReadByte", (*bufio.Reader).ReadByte, execmReaderReadByte),
		I.Func("(*Reader).UnreadByte", (*bufio.Reader).UnreadByte, execmReaderUnreadByte),
		I.Func("(*Reader).ReadRune", (*bufio.Reader).ReadRune, execmReaderReadRune),
		I.Func("(*Reader).UnreadRune", (*bufio.Reader).UnreadRune, execmReaderUnreadRune),
		I.Func("(*Reader).Buffered", (*bufio.Reader).Buffered, execmReaderBuffered),
		I.Func("(*Reader).ReadSlice", (*bufio.Reader).ReadSlice, execmReaderReadSlice),
		I.Func("(*Reader).ReadLine", (*bufio.Reader).ReadLine, execmReaderReadLine),
		I.Func("(*Reader).ReadBytes", (*bufio.Reader).ReadBytes, execmReaderReadBytes),
		I.Func("(*Reader).ReadString", (*bufio.Reader).ReadString, execmReaderReadString),
		I.Func("(*Reader).WriteTo", (*bufio.Reader).WriteTo, execmReaderWriteTo),
		I.Func("ScanBytes", bufio.ScanBytes, execScanBytes),
		I.Func("ScanLines", bufio.ScanLines, execScanLines),
		I.Func("ScanRunes", bufio.ScanRunes, execScanRunes),
		I.Func("ScanWords", bufio.ScanWords, execScanWords),
		I.Func("(*Scanner).Err", (*bufio.Scanner).Err, execmScannerErr),
		I.Func("(*Scanner).Bytes", (*bufio.Scanner).Bytes, execmScannerBytes),
		I.Func("(*Scanner).Text", (*bufio.Scanner).Text, execmScannerText),
		I.Func("(*Scanner).Scan", (*bufio.Scanner).Scan, execmScannerScan),
		I.Func("(*Scanner).Buffer", (*bufio.Scanner).Buffer, execmScannerBuffer),
		I.Func("(*Scanner).Split", (*bufio.Scanner).Split, execmScannerSplit),
		I.Func("(*Writer).Size", (*bufio.Writer).Size, execmWriterSize),
		I.Func("(*Writer).Reset", (*bufio.Writer).Reset, execmWriterReset),
		I.Func("(*Writer).Flush", (*bufio.Writer).Flush, execmWriterFlush),
		I.Func("(*Writer).Available", (*bufio.Writer).Available, execmWriterAvailable),
		I.Func("(*Writer).Buffered", (*bufio.Writer).Buffered, execmWriterBuffered),
		I.Func("(*Writer).Write", (*bufio.Writer).Write, execmWriterWrite),
		I.Func("(*Writer).WriteByte", (*bufio.Writer).WriteByte, execmWriterWriteByte),
		I.Func("(*Writer).WriteRune", (*bufio.Writer).WriteRune, execmWriterWriteRune),
		I.Func("(*Writer).WriteString", (*bufio.Writer).WriteString, execmWriterWriteString),
		I.Func("(*Writer).ReadFrom", (*bufio.Writer).ReadFrom, execmWriterReadFrom),
	)
	I.RegisterVars(
		I.Var("ErrAdvanceTooFar", &bufio.ErrAdvanceTooFar),
		I.Var("ErrBufferFull", &bufio.ErrBufferFull),
		I.Var("ErrFinalToken", &bufio.ErrFinalToken),
		I.Var("ErrInvalidUnreadByte", &bufio.ErrInvalidUnreadByte),
		I.Var("ErrInvalidUnreadRune", &bufio.ErrInvalidUnreadRune),
		I.Var("ErrNegativeAdvance", &bufio.ErrNegativeAdvance),
		I.Var("ErrNegativeCount", &bufio.ErrNegativeCount),
		I.Var("ErrTooLong", &bufio.ErrTooLong),
	)
	I.RegisterTypes(
		I.Type("ReadWriter", reflect.TypeOf((*bufio.ReadWriter)(nil)).Elem()),
		I.Type("Reader", reflect.TypeOf((*bufio.Reader)(nil)).Elem()),
		I.Type("Scanner", reflect.TypeOf((*bufio.Scanner)(nil)).Elem()),
		I.Type("SplitFunc", reflect.TypeOf((*bufio.SplitFunc)(nil)).Elem()),
		I.Type("Writer", reflect.TypeOf((*bufio.Writer)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("MaxScanTokenSize", qspec.ConstUnboundInt, bufio.MaxScanTokenSize),
	)
}
