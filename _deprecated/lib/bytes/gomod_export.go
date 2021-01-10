// Package bytes provide Go+ "bytes" package, as "bytes" package in Go.
package bytes

import (
	bytes "bytes"
	io "io"
	reflect "reflect"
	unicode "unicode"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execmBufferBytes(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*bytes.Buffer).Bytes()
	p.Ret(1, ret0)
}

func execmBufferString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*bytes.Buffer).String()
	p.Ret(1, ret0)
}

func execmBufferLen(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*bytes.Buffer).Len()
	p.Ret(1, ret0)
}

func execmBufferCap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*bytes.Buffer).Cap()
	p.Ret(1, ret0)
}

func execmBufferTruncate(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*bytes.Buffer).Truncate(args[1].(int))
	p.PopN(2)
}

func execmBufferReset(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*bytes.Buffer).Reset()
	p.PopN(1)
}

func execmBufferGrow(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*bytes.Buffer).Grow(args[1].(int))
	p.PopN(2)
}

func execmBufferWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*bytes.Buffer).Write(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmBufferWriteString(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*bytes.Buffer).WriteString(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func toType0(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execmBufferReadFrom(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*bytes.Buffer).ReadFrom(toType0(args[1]))
	p.Ret(2, ret0, ret1)
}

func toType1(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execmBufferWriteTo(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*bytes.Buffer).WriteTo(toType1(args[1]))
	p.Ret(2, ret0, ret1)
}

func execmBufferWriteByte(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*bytes.Buffer).WriteByte(args[1].(byte))
	p.Ret(2, ret0)
}

func execmBufferWriteRune(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*bytes.Buffer).WriteRune(args[1].(rune))
	p.Ret(2, ret0, ret1)
}

func execmBufferRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*bytes.Buffer).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmBufferNext(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*bytes.Buffer).Next(args[1].(int))
	p.Ret(2, ret0)
}

func execmBufferReadByte(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*bytes.Buffer).ReadByte()
	p.Ret(1, ret0, ret1)
}

func execmBufferReadRune(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := args[0].(*bytes.Buffer).ReadRune()
	p.Ret(1, ret0, ret1, ret2)
}

func execmBufferUnreadRune(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*bytes.Buffer).UnreadRune()
	p.Ret(1, ret0)
}

func execmBufferUnreadByte(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*bytes.Buffer).UnreadByte()
	p.Ret(1, ret0)
}

func execmBufferReadBytes(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*bytes.Buffer).ReadBytes(args[1].(byte))
	p.Ret(2, ret0, ret1)
}

func execmBufferReadString(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*bytes.Buffer).ReadString(args[1].(byte))
	p.Ret(2, ret0, ret1)
}

func execCompare(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.Compare(args[0].([]byte), args[1].([]byte))
	p.Ret(2, ret0)
}

func execContains(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.Contains(args[0].([]byte), args[1].([]byte))
	p.Ret(2, ret0)
}

func execContainsAny(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.ContainsAny(args[0].([]byte), args[1].(string))
	p.Ret(2, ret0)
}

func execContainsRune(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.ContainsRune(args[0].([]byte), args[1].(rune))
	p.Ret(2, ret0)
}

func execCount(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.Count(args[0].([]byte), args[1].([]byte))
	p.Ret(2, ret0)
}

func execEqual(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.Equal(args[0].([]byte), args[1].([]byte))
	p.Ret(2, ret0)
}

func execEqualFold(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.EqualFold(args[0].([]byte), args[1].([]byte))
	p.Ret(2, ret0)
}

func execFields(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bytes.Fields(args[0].([]byte))
	p.Ret(1, ret0)
}

func execFieldsFunc(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.FieldsFunc(args[0].([]byte), args[1].(func(rune) bool))
	p.Ret(2, ret0)
}

func execHasPrefix(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.HasPrefix(args[0].([]byte), args[1].([]byte))
	p.Ret(2, ret0)
}

func execHasSuffix(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.HasSuffix(args[0].([]byte), args[1].([]byte))
	p.Ret(2, ret0)
}

func execIndex(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.Index(args[0].([]byte), args[1].([]byte))
	p.Ret(2, ret0)
}

func execIndexAny(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.IndexAny(args[0].([]byte), args[1].(string))
	p.Ret(2, ret0)
}

func execIndexByte(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.IndexByte(args[0].([]byte), args[1].(byte))
	p.Ret(2, ret0)
}

func execIndexFunc(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.IndexFunc(args[0].([]byte), args[1].(func(r rune) bool))
	p.Ret(2, ret0)
}

func execIndexRune(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.IndexRune(args[0].([]byte), args[1].(rune))
	p.Ret(2, ret0)
}

func execJoin(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.Join(args[0].([][]byte), args[1].([]byte))
	p.Ret(2, ret0)
}

func execLastIndex(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.LastIndex(args[0].([]byte), args[1].([]byte))
	p.Ret(2, ret0)
}

func execLastIndexAny(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.LastIndexAny(args[0].([]byte), args[1].(string))
	p.Ret(2, ret0)
}

func execLastIndexByte(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.LastIndexByte(args[0].([]byte), args[1].(byte))
	p.Ret(2, ret0)
}

func execLastIndexFunc(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.LastIndexFunc(args[0].([]byte), args[1].(func(r rune) bool))
	p.Ret(2, ret0)
}

func execMap(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.Map(args[0].(func(r rune) rune), args[1].([]byte))
	p.Ret(2, ret0)
}

func execNewBuffer(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bytes.NewBuffer(args[0].([]byte))
	p.Ret(1, ret0)
}

func execNewBufferString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bytes.NewBufferString(args[0].(string))
	p.Ret(1, ret0)
}

func execNewReader(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bytes.NewReader(args[0].([]byte))
	p.Ret(1, ret0)
}

func execmReaderLen(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*bytes.Reader).Len()
	p.Ret(1, ret0)
}

func execmReaderSize(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*bytes.Reader).Size()
	p.Ret(1, ret0)
}

func execmReaderRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*bytes.Reader).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmReaderReadAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*bytes.Reader).ReadAt(args[1].([]byte), args[2].(int64))
	p.Ret(3, ret0, ret1)
}

func execmReaderReadByte(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*bytes.Reader).ReadByte()
	p.Ret(1, ret0, ret1)
}

func execmReaderUnreadByte(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*bytes.Reader).UnreadByte()
	p.Ret(1, ret0)
}

func execmReaderReadRune(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := args[0].(*bytes.Reader).ReadRune()
	p.Ret(1, ret0, ret1, ret2)
}

func execmReaderUnreadRune(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*bytes.Reader).UnreadRune()
	p.Ret(1, ret0)
}

func execmReaderSeek(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*bytes.Reader).Seek(args[1].(int64), args[2].(int))
	p.Ret(3, ret0, ret1)
}

func execmReaderWriteTo(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*bytes.Reader).WriteTo(toType1(args[1]))
	p.Ret(2, ret0, ret1)
}

func execmReaderReset(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*bytes.Reader).Reset(args[1].([]byte))
	p.PopN(2)
}

func execRepeat(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.Repeat(args[0].([]byte), args[1].(int))
	p.Ret(2, ret0)
}

func execReplace(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := bytes.Replace(args[0].([]byte), args[1].([]byte), args[2].([]byte), args[3].(int))
	p.Ret(4, ret0)
}

func execReplaceAll(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := bytes.ReplaceAll(args[0].([]byte), args[1].([]byte), args[2].([]byte))
	p.Ret(3, ret0)
}

func execRunes(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bytes.Runes(args[0].([]byte))
	p.Ret(1, ret0)
}

func execSplit(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.Split(args[0].([]byte), args[1].([]byte))
	p.Ret(2, ret0)
}

func execSplitAfter(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.SplitAfter(args[0].([]byte), args[1].([]byte))
	p.Ret(2, ret0)
}

func execSplitAfterN(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := bytes.SplitAfterN(args[0].([]byte), args[1].([]byte), args[2].(int))
	p.Ret(3, ret0)
}

func execSplitN(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := bytes.SplitN(args[0].([]byte), args[1].([]byte), args[2].(int))
	p.Ret(3, ret0)
}

func execTitle(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bytes.Title(args[0].([]byte))
	p.Ret(1, ret0)
}

func execToLower(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bytes.ToLower(args[0].([]byte))
	p.Ret(1, ret0)
}

func execToLowerSpecial(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.ToLowerSpecial(args[0].(unicode.SpecialCase), args[1].([]byte))
	p.Ret(2, ret0)
}

func execToTitle(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bytes.ToTitle(args[0].([]byte))
	p.Ret(1, ret0)
}

func execToTitleSpecial(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.ToTitleSpecial(args[0].(unicode.SpecialCase), args[1].([]byte))
	p.Ret(2, ret0)
}

func execToUpper(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bytes.ToUpper(args[0].([]byte))
	p.Ret(1, ret0)
}

func execToUpperSpecial(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.ToUpperSpecial(args[0].(unicode.SpecialCase), args[1].([]byte))
	p.Ret(2, ret0)
}

func execToValidUTF8(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.ToValidUTF8(args[0].([]byte), args[1].([]byte))
	p.Ret(2, ret0)
}

func execTrim(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.Trim(args[0].([]byte), args[1].(string))
	p.Ret(2, ret0)
}

func execTrimFunc(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.TrimFunc(args[0].([]byte), args[1].(func(r rune) bool))
	p.Ret(2, ret0)
}

func execTrimLeft(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.TrimLeft(args[0].([]byte), args[1].(string))
	p.Ret(2, ret0)
}

func execTrimLeftFunc(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.TrimLeftFunc(args[0].([]byte), args[1].(func(r rune) bool))
	p.Ret(2, ret0)
}

func execTrimPrefix(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.TrimPrefix(args[0].([]byte), args[1].([]byte))
	p.Ret(2, ret0)
}

func execTrimRight(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.TrimRight(args[0].([]byte), args[1].(string))
	p.Ret(2, ret0)
}

func execTrimRightFunc(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.TrimRightFunc(args[0].([]byte), args[1].(func(r rune) bool))
	p.Ret(2, ret0)
}

func execTrimSpace(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := bytes.TrimSpace(args[0].([]byte))
	p.Ret(1, ret0)
}

func execTrimSuffix(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := bytes.TrimSuffix(args[0].([]byte), args[1].([]byte))
	p.Ret(2, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("bytes")

func init() {
	I.RegisterFuncs(
		I.Func("(*Buffer).Bytes", (*bytes.Buffer).Bytes, execmBufferBytes),
		I.Func("(*Buffer).String", (*bytes.Buffer).String, execmBufferString),
		I.Func("(*Buffer).Len", (*bytes.Buffer).Len, execmBufferLen),
		I.Func("(*Buffer).Cap", (*bytes.Buffer).Cap, execmBufferCap),
		I.Func("(*Buffer).Truncate", (*bytes.Buffer).Truncate, execmBufferTruncate),
		I.Func("(*Buffer).Reset", (*bytes.Buffer).Reset, execmBufferReset),
		I.Func("(*Buffer).Grow", (*bytes.Buffer).Grow, execmBufferGrow),
		I.Func("(*Buffer).Write", (*bytes.Buffer).Write, execmBufferWrite),
		I.Func("(*Buffer).WriteString", (*bytes.Buffer).WriteString, execmBufferWriteString),
		I.Func("(*Buffer).ReadFrom", (*bytes.Buffer).ReadFrom, execmBufferReadFrom),
		I.Func("(*Buffer).WriteTo", (*bytes.Buffer).WriteTo, execmBufferWriteTo),
		I.Func("(*Buffer).WriteByte", (*bytes.Buffer).WriteByte, execmBufferWriteByte),
		I.Func("(*Buffer).WriteRune", (*bytes.Buffer).WriteRune, execmBufferWriteRune),
		I.Func("(*Buffer).Read", (*bytes.Buffer).Read, execmBufferRead),
		I.Func("(*Buffer).Next", (*bytes.Buffer).Next, execmBufferNext),
		I.Func("(*Buffer).ReadByte", (*bytes.Buffer).ReadByte, execmBufferReadByte),
		I.Func("(*Buffer).ReadRune", (*bytes.Buffer).ReadRune, execmBufferReadRune),
		I.Func("(*Buffer).UnreadRune", (*bytes.Buffer).UnreadRune, execmBufferUnreadRune),
		I.Func("(*Buffer).UnreadByte", (*bytes.Buffer).UnreadByte, execmBufferUnreadByte),
		I.Func("(*Buffer).ReadBytes", (*bytes.Buffer).ReadBytes, execmBufferReadBytes),
		I.Func("(*Buffer).ReadString", (*bytes.Buffer).ReadString, execmBufferReadString),
		I.Func("Compare", bytes.Compare, execCompare),
		I.Func("Contains", bytes.Contains, execContains),
		I.Func("ContainsAny", bytes.ContainsAny, execContainsAny),
		I.Func("ContainsRune", bytes.ContainsRune, execContainsRune),
		I.Func("Count", bytes.Count, execCount),
		I.Func("Equal", bytes.Equal, execEqual),
		I.Func("EqualFold", bytes.EqualFold, execEqualFold),
		I.Func("Fields", bytes.Fields, execFields),
		I.Func("FieldsFunc", bytes.FieldsFunc, execFieldsFunc),
		I.Func("HasPrefix", bytes.HasPrefix, execHasPrefix),
		I.Func("HasSuffix", bytes.HasSuffix, execHasSuffix),
		I.Func("Index", bytes.Index, execIndex),
		I.Func("IndexAny", bytes.IndexAny, execIndexAny),
		I.Func("IndexByte", bytes.IndexByte, execIndexByte),
		I.Func("IndexFunc", bytes.IndexFunc, execIndexFunc),
		I.Func("IndexRune", bytes.IndexRune, execIndexRune),
		I.Func("Join", bytes.Join, execJoin),
		I.Func("LastIndex", bytes.LastIndex, execLastIndex),
		I.Func("LastIndexAny", bytes.LastIndexAny, execLastIndexAny),
		I.Func("LastIndexByte", bytes.LastIndexByte, execLastIndexByte),
		I.Func("LastIndexFunc", bytes.LastIndexFunc, execLastIndexFunc),
		I.Func("Map", bytes.Map, execMap),
		I.Func("NewBuffer", bytes.NewBuffer, execNewBuffer),
		I.Func("NewBufferString", bytes.NewBufferString, execNewBufferString),
		I.Func("NewReader", bytes.NewReader, execNewReader),
		I.Func("(*Reader).Len", (*bytes.Reader).Len, execmReaderLen),
		I.Func("(*Reader).Size", (*bytes.Reader).Size, execmReaderSize),
		I.Func("(*Reader).Read", (*bytes.Reader).Read, execmReaderRead),
		I.Func("(*Reader).ReadAt", (*bytes.Reader).ReadAt, execmReaderReadAt),
		I.Func("(*Reader).ReadByte", (*bytes.Reader).ReadByte, execmReaderReadByte),
		I.Func("(*Reader).UnreadByte", (*bytes.Reader).UnreadByte, execmReaderUnreadByte),
		I.Func("(*Reader).ReadRune", (*bytes.Reader).ReadRune, execmReaderReadRune),
		I.Func("(*Reader).UnreadRune", (*bytes.Reader).UnreadRune, execmReaderUnreadRune),
		I.Func("(*Reader).Seek", (*bytes.Reader).Seek, execmReaderSeek),
		I.Func("(*Reader).WriteTo", (*bytes.Reader).WriteTo, execmReaderWriteTo),
		I.Func("(*Reader).Reset", (*bytes.Reader).Reset, execmReaderReset),
		I.Func("Repeat", bytes.Repeat, execRepeat),
		I.Func("Replace", bytes.Replace, execReplace),
		I.Func("ReplaceAll", bytes.ReplaceAll, execReplaceAll),
		I.Func("Runes", bytes.Runes, execRunes),
		I.Func("Split", bytes.Split, execSplit),
		I.Func("SplitAfter", bytes.SplitAfter, execSplitAfter),
		I.Func("SplitAfterN", bytes.SplitAfterN, execSplitAfterN),
		I.Func("SplitN", bytes.SplitN, execSplitN),
		I.Func("Title", bytes.Title, execTitle),
		I.Func("ToLower", bytes.ToLower, execToLower),
		I.Func("ToLowerSpecial", bytes.ToLowerSpecial, execToLowerSpecial),
		I.Func("ToTitle", bytes.ToTitle, execToTitle),
		I.Func("ToTitleSpecial", bytes.ToTitleSpecial, execToTitleSpecial),
		I.Func("ToUpper", bytes.ToUpper, execToUpper),
		I.Func("ToUpperSpecial", bytes.ToUpperSpecial, execToUpperSpecial),
		I.Func("ToValidUTF8", bytes.ToValidUTF8, execToValidUTF8),
		I.Func("Trim", bytes.Trim, execTrim),
		I.Func("TrimFunc", bytes.TrimFunc, execTrimFunc),
		I.Func("TrimLeft", bytes.TrimLeft, execTrimLeft),
		I.Func("TrimLeftFunc", bytes.TrimLeftFunc, execTrimLeftFunc),
		I.Func("TrimPrefix", bytes.TrimPrefix, execTrimPrefix),
		I.Func("TrimRight", bytes.TrimRight, execTrimRight),
		I.Func("TrimRightFunc", bytes.TrimRightFunc, execTrimRightFunc),
		I.Func("TrimSpace", bytes.TrimSpace, execTrimSpace),
		I.Func("TrimSuffix", bytes.TrimSuffix, execTrimSuffix),
	)
	I.RegisterVars(
		I.Var("ErrTooLarge", &bytes.ErrTooLarge),
	)
	I.RegisterTypes(
		I.Type("Buffer", reflect.TypeOf((*bytes.Buffer)(nil)).Elem()),
		I.Type("Reader", reflect.TypeOf((*bytes.Reader)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("MinRead", qspec.ConstUnboundInt, bytes.MinRead),
	)
}
