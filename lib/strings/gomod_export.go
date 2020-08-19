// Package strings provide Go+ "strings" package, as "strings" package in Go.
package strings

import (
	io "io"
	reflect "reflect"
	strings "strings"
	unicode "unicode"

	gop "github.com/goplus/gop"
)

func execmBuilderString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*strings.Builder).String()
	p.Ret(1, ret0)
}

func execmBuilderLen(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*strings.Builder).Len()
	p.Ret(1, ret0)
}

func execmBuilderCap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*strings.Builder).Cap()
	p.Ret(1, ret0)
}

func execmBuilderReset(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*strings.Builder).Reset()
	p.PopN(1)
}

func execmBuilderGrow(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*strings.Builder).Grow(args[1].(int))
	p.PopN(2)
}

func execmBuilderWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*strings.Builder).Write(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmBuilderWriteByte(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*strings.Builder).WriteByte(args[1].(byte))
	p.Ret(2, ret0)
}

func execmBuilderWriteRune(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*strings.Builder).WriteRune(args[1].(rune))
	p.Ret(2, ret0, ret1)
}

func execmBuilderWriteString(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*strings.Builder).WriteString(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execCompare(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.Compare(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execContains(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.Contains(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execContainsAny(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.ContainsAny(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execContainsRune(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.ContainsRune(args[0].(string), args[1].(rune))
	p.Ret(2, ret0)
}

func execCount(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.Count(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execEqualFold(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.EqualFold(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execFields(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := strings.Fields(args[0].(string))
	p.Ret(1, ret0)
}

func execFieldsFunc(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.FieldsFunc(args[0].(string), args[1].(func(rune) bool))
	p.Ret(2, ret0)
}

func execHasPrefix(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.HasPrefix(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execHasSuffix(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.HasSuffix(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execIndex(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.Index(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execIndexAny(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.IndexAny(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execIndexByte(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.IndexByte(args[0].(string), args[1].(byte))
	p.Ret(2, ret0)
}

func execIndexFunc(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.IndexFunc(args[0].(string), args[1].(func(rune) bool))
	p.Ret(2, ret0)
}

func execIndexRune(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.IndexRune(args[0].(string), args[1].(rune))
	p.Ret(2, ret0)
}

func execJoin(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.Join(args[0].([]string), args[1].(string))
	p.Ret(2, ret0)
}

func execLastIndex(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.LastIndex(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execLastIndexAny(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.LastIndexAny(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execLastIndexByte(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.LastIndexByte(args[0].(string), args[1].(byte))
	p.Ret(2, ret0)
}

func execLastIndexFunc(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.LastIndexFunc(args[0].(string), args[1].(func(rune) bool))
	p.Ret(2, ret0)
}

func execMap(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.Map(args[0].(func(rune) rune), args[1].(string))
	p.Ret(2, ret0)
}

func execNewReader(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := strings.NewReader(args[0].(string))
	p.Ret(1, ret0)
}

func execNewReplacer(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := strings.NewReplacer(gop.ToStrings(args)...)
	p.Ret(arity, ret0)
}

func execmReaderLen(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*strings.Reader).Len()
	p.Ret(1, ret0)
}

func execmReaderSize(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*strings.Reader).Size()
	p.Ret(1, ret0)
}

func execmReaderRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*strings.Reader).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmReaderReadAt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*strings.Reader).ReadAt(args[1].([]byte), args[2].(int64))
	p.Ret(3, ret0, ret1)
}

func execmReaderReadByte(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*strings.Reader).ReadByte()
	p.Ret(1, ret0, ret1)
}

func execmReaderUnreadByte(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*strings.Reader).UnreadByte()
	p.Ret(1, ret0)
}

func execmReaderReadRune(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := args[0].(*strings.Reader).ReadRune()
	p.Ret(1, ret0, ret1, ret2)
}

func execmReaderUnreadRune(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*strings.Reader).UnreadRune()
	p.Ret(1, ret0)
}

func execmReaderSeek(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*strings.Reader).Seek(args[1].(int64), args[2].(int))
	p.Ret(3, ret0, ret1)
}

func toType0(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execmReaderWriteTo(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*strings.Reader).WriteTo(toType0(args[1]))
	p.Ret(2, ret0, ret1)
}

func execmReaderReset(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*strings.Reader).Reset(args[1].(string))
	p.PopN(2)
}

func execRepeat(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.Repeat(args[0].(string), args[1].(int))
	p.Ret(2, ret0)
}

func execReplace(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := strings.Replace(args[0].(string), args[1].(string), args[2].(string), args[3].(int))
	p.Ret(4, ret0)
}

func execReplaceAll(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := strings.ReplaceAll(args[0].(string), args[1].(string), args[2].(string))
	p.Ret(3, ret0)
}

func execmReplacerReplace(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*strings.Replacer).Replace(args[1].(string))
	p.Ret(2, ret0)
}

func execmReplacerWriteString(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*strings.Replacer).WriteString(toType0(args[1]), args[2].(string))
	p.Ret(3, ret0, ret1)
}

func execSplit(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.Split(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execSplitAfter(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.SplitAfter(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execSplitAfterN(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := strings.SplitAfterN(args[0].(string), args[1].(string), args[2].(int))
	p.Ret(3, ret0)
}

func execSplitN(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := strings.SplitN(args[0].(string), args[1].(string), args[2].(int))
	p.Ret(3, ret0)
}

func execTitle(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := strings.Title(args[0].(string))
	p.Ret(1, ret0)
}

func execToLower(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := strings.ToLower(args[0].(string))
	p.Ret(1, ret0)
}

func execToLowerSpecial(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.ToLowerSpecial(args[0].(unicode.SpecialCase), args[1].(string))
	p.Ret(2, ret0)
}

func execToTitle(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := strings.ToTitle(args[0].(string))
	p.Ret(1, ret0)
}

func execToTitleSpecial(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.ToTitleSpecial(args[0].(unicode.SpecialCase), args[1].(string))
	p.Ret(2, ret0)
}

func execToUpper(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := strings.ToUpper(args[0].(string))
	p.Ret(1, ret0)
}

func execToUpperSpecial(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.ToUpperSpecial(args[0].(unicode.SpecialCase), args[1].(string))
	p.Ret(2, ret0)
}

func execToValidUTF8(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.ToValidUTF8(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execTrim(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.Trim(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execTrimFunc(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.TrimFunc(args[0].(string), args[1].(func(rune) bool))
	p.Ret(2, ret0)
}

func execTrimLeft(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.TrimLeft(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execTrimLeftFunc(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.TrimLeftFunc(args[0].(string), args[1].(func(rune) bool))
	p.Ret(2, ret0)
}

func execTrimPrefix(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.TrimPrefix(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execTrimRight(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.TrimRight(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execTrimRightFunc(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.TrimRightFunc(args[0].(string), args[1].(func(rune) bool))
	p.Ret(2, ret0)
}

func execTrimSpace(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := strings.TrimSpace(args[0].(string))
	p.Ret(1, ret0)
}

func execTrimSuffix(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strings.TrimSuffix(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("strings")

func init() {
	I.RegisterFuncs(
		I.Func("(*Builder).String", (*strings.Builder).String, execmBuilderString),
		I.Func("(*Builder).Len", (*strings.Builder).Len, execmBuilderLen),
		I.Func("(*Builder).Cap", (*strings.Builder).Cap, execmBuilderCap),
		I.Func("(*Builder).Reset", (*strings.Builder).Reset, execmBuilderReset),
		I.Func("(*Builder).Grow", (*strings.Builder).Grow, execmBuilderGrow),
		I.Func("(*Builder).Write", (*strings.Builder).Write, execmBuilderWrite),
		I.Func("(*Builder).WriteByte", (*strings.Builder).WriteByte, execmBuilderWriteByte),
		I.Func("(*Builder).WriteRune", (*strings.Builder).WriteRune, execmBuilderWriteRune),
		I.Func("(*Builder).WriteString", (*strings.Builder).WriteString, execmBuilderWriteString),
		I.Func("Compare", strings.Compare, execCompare),
		I.Func("Contains", strings.Contains, execContains),
		I.Func("ContainsAny", strings.ContainsAny, execContainsAny),
		I.Func("ContainsRune", strings.ContainsRune, execContainsRune),
		I.Func("Count", strings.Count, execCount),
		I.Func("EqualFold", strings.EqualFold, execEqualFold),
		I.Func("Fields", strings.Fields, execFields),
		I.Func("FieldsFunc", strings.FieldsFunc, execFieldsFunc),
		I.Func("HasPrefix", strings.HasPrefix, execHasPrefix),
		I.Func("HasSuffix", strings.HasSuffix, execHasSuffix),
		I.Func("Index", strings.Index, execIndex),
		I.Func("IndexAny", strings.IndexAny, execIndexAny),
		I.Func("IndexByte", strings.IndexByte, execIndexByte),
		I.Func("IndexFunc", strings.IndexFunc, execIndexFunc),
		I.Func("IndexRune", strings.IndexRune, execIndexRune),
		I.Func("Join", strings.Join, execJoin),
		I.Func("LastIndex", strings.LastIndex, execLastIndex),
		I.Func("LastIndexAny", strings.LastIndexAny, execLastIndexAny),
		I.Func("LastIndexByte", strings.LastIndexByte, execLastIndexByte),
		I.Func("LastIndexFunc", strings.LastIndexFunc, execLastIndexFunc),
		I.Func("Map", strings.Map, execMap),
		I.Func("NewReader", strings.NewReader, execNewReader),
		I.Func("(*Reader).Len", (*strings.Reader).Len, execmReaderLen),
		I.Func("(*Reader).Size", (*strings.Reader).Size, execmReaderSize),
		I.Func("(*Reader).Read", (*strings.Reader).Read, execmReaderRead),
		I.Func("(*Reader).ReadAt", (*strings.Reader).ReadAt, execmReaderReadAt),
		I.Func("(*Reader).ReadByte", (*strings.Reader).ReadByte, execmReaderReadByte),
		I.Func("(*Reader).UnreadByte", (*strings.Reader).UnreadByte, execmReaderUnreadByte),
		I.Func("(*Reader).ReadRune", (*strings.Reader).ReadRune, execmReaderReadRune),
		I.Func("(*Reader).UnreadRune", (*strings.Reader).UnreadRune, execmReaderUnreadRune),
		I.Func("(*Reader).Seek", (*strings.Reader).Seek, execmReaderSeek),
		I.Func("(*Reader).WriteTo", (*strings.Reader).WriteTo, execmReaderWriteTo),
		I.Func("(*Reader).Reset", (*strings.Reader).Reset, execmReaderReset),
		I.Func("Repeat", strings.Repeat, execRepeat),
		I.Func("Replace", strings.Replace, execReplace),
		I.Func("ReplaceAll", strings.ReplaceAll, execReplaceAll),
		I.Func("(*Replacer).Replace", (*strings.Replacer).Replace, execmReplacerReplace),
		I.Func("(*Replacer).WriteString", (*strings.Replacer).WriteString, execmReplacerWriteString),
		I.Func("Split", strings.Split, execSplit),
		I.Func("SplitAfter", strings.SplitAfter, execSplitAfter),
		I.Func("SplitAfterN", strings.SplitAfterN, execSplitAfterN),
		I.Func("SplitN", strings.SplitN, execSplitN),
		I.Func("Title", strings.Title, execTitle),
		I.Func("ToLower", strings.ToLower, execToLower),
		I.Func("ToLowerSpecial", strings.ToLowerSpecial, execToLowerSpecial),
		I.Func("ToTitle", strings.ToTitle, execToTitle),
		I.Func("ToTitleSpecial", strings.ToTitleSpecial, execToTitleSpecial),
		I.Func("ToUpper", strings.ToUpper, execToUpper),
		I.Func("ToUpperSpecial", strings.ToUpperSpecial, execToUpperSpecial),
		I.Func("ToValidUTF8", strings.ToValidUTF8, execToValidUTF8),
		I.Func("Trim", strings.Trim, execTrim),
		I.Func("TrimFunc", strings.TrimFunc, execTrimFunc),
		I.Func("TrimLeft", strings.TrimLeft, execTrimLeft),
		I.Func("TrimLeftFunc", strings.TrimLeftFunc, execTrimLeftFunc),
		I.Func("TrimPrefix", strings.TrimPrefix, execTrimPrefix),
		I.Func("TrimRight", strings.TrimRight, execTrimRight),
		I.Func("TrimRightFunc", strings.TrimRightFunc, execTrimRightFunc),
		I.Func("TrimSpace", strings.TrimSpace, execTrimSpace),
		I.Func("TrimSuffix", strings.TrimSuffix, execTrimSuffix),
	)
	I.RegisterFuncvs(
		I.Funcv("NewReplacer", strings.NewReplacer, execNewReplacer),
	)
	I.RegisterTypes(
		I.Type("Builder", reflect.TypeOf((*strings.Builder)(nil)).Elem()),
		I.Type("Reader", reflect.TypeOf((*strings.Reader)(nil)).Elem()),
		I.Type("Replacer", reflect.TypeOf((*strings.Replacer)(nil)).Elem()),
	)
}
