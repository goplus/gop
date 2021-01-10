// Package strconv provide Go+ "strconv" package, as "strconv" package in Go.
package strconv

import (
	reflect "reflect"
	strconv "strconv"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execAppendBool(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strconv.AppendBool(args[0].([]byte), args[1].(bool))
	p.Ret(2, ret0)
}

func execAppendFloat(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	ret0 := strconv.AppendFloat(args[0].([]byte), args[1].(float64), args[2].(byte), args[3].(int), args[4].(int))
	p.Ret(5, ret0)
}

func execAppendInt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := strconv.AppendInt(args[0].([]byte), args[1].(int64), args[2].(int))
	p.Ret(3, ret0)
}

func execAppendQuote(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strconv.AppendQuote(args[0].([]byte), args[1].(string))
	p.Ret(2, ret0)
}

func execAppendQuoteRune(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strconv.AppendQuoteRune(args[0].([]byte), args[1].(rune))
	p.Ret(2, ret0)
}

func execAppendQuoteRuneToASCII(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strconv.AppendQuoteRuneToASCII(args[0].([]byte), args[1].(rune))
	p.Ret(2, ret0)
}

func execAppendQuoteRuneToGraphic(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strconv.AppendQuoteRuneToGraphic(args[0].([]byte), args[1].(rune))
	p.Ret(2, ret0)
}

func execAppendQuoteToASCII(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strconv.AppendQuoteToASCII(args[0].([]byte), args[1].(string))
	p.Ret(2, ret0)
}

func execAppendQuoteToGraphic(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strconv.AppendQuoteToGraphic(args[0].([]byte), args[1].(string))
	p.Ret(2, ret0)
}

func execAppendUint(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := strconv.AppendUint(args[0].([]byte), args[1].(uint64), args[2].(int))
	p.Ret(3, ret0)
}

func execAtoi(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := strconv.Atoi(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execCanBackquote(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := strconv.CanBackquote(args[0].(string))
	p.Ret(1, ret0)
}

func execFormatBool(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := strconv.FormatBool(args[0].(bool))
	p.Ret(1, ret0)
}

func execFormatFloat(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := strconv.FormatFloat(args[0].(float64), args[1].(byte), args[2].(int), args[3].(int))
	p.Ret(4, ret0)
}

func execFormatInt(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strconv.FormatInt(args[0].(int64), args[1].(int))
	p.Ret(2, ret0)
}

func execFormatUint(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := strconv.FormatUint(args[0].(uint64), args[1].(int))
	p.Ret(2, ret0)
}

func execIsGraphic(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := strconv.IsGraphic(args[0].(rune))
	p.Ret(1, ret0)
}

func execIsPrint(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := strconv.IsPrint(args[0].(rune))
	p.Ret(1, ret0)
}

func execItoa(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := strconv.Itoa(args[0].(int))
	p.Ret(1, ret0)
}

func execmNumErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*strconv.NumError).Error()
	p.Ret(1, ret0)
}

func execmNumErrorUnwrap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*strconv.NumError).Unwrap()
	p.Ret(1, ret0)
}

func execParseBool(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := strconv.ParseBool(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execParseFloat(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := strconv.ParseFloat(args[0].(string), args[1].(int))
	p.Ret(2, ret0, ret1)
}

func execParseInt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := strconv.ParseInt(args[0].(string), args[1].(int), args[2].(int))
	p.Ret(3, ret0, ret1)
}

func execParseUint(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := strconv.ParseUint(args[0].(string), args[1].(int), args[2].(int))
	p.Ret(3, ret0, ret1)
}

func execQuote(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := strconv.Quote(args[0].(string))
	p.Ret(1, ret0)
}

func execQuoteRune(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := strconv.QuoteRune(args[0].(rune))
	p.Ret(1, ret0)
}

func execQuoteRuneToASCII(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := strconv.QuoteRuneToASCII(args[0].(rune))
	p.Ret(1, ret0)
}

func execQuoteRuneToGraphic(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := strconv.QuoteRuneToGraphic(args[0].(rune))
	p.Ret(1, ret0)
}

func execQuoteToASCII(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := strconv.QuoteToASCII(args[0].(string))
	p.Ret(1, ret0)
}

func execQuoteToGraphic(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := strconv.QuoteToGraphic(args[0].(string))
	p.Ret(1, ret0)
}

func execUnquote(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := strconv.Unquote(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execUnquoteChar(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1, ret2, ret3 := strconv.UnquoteChar(args[0].(string), args[1].(byte))
	p.Ret(2, ret0, ret1, ret2, ret3)
}

// I is a Go package instance.
var I = gop.NewGoPackage("strconv")

func init() {
	I.RegisterFuncs(
		I.Func("AppendBool", strconv.AppendBool, execAppendBool),
		I.Func("AppendFloat", strconv.AppendFloat, execAppendFloat),
		I.Func("AppendInt", strconv.AppendInt, execAppendInt),
		I.Func("AppendQuote", strconv.AppendQuote, execAppendQuote),
		I.Func("AppendQuoteRune", strconv.AppendQuoteRune, execAppendQuoteRune),
		I.Func("AppendQuoteRuneToASCII", strconv.AppendQuoteRuneToASCII, execAppendQuoteRuneToASCII),
		I.Func("AppendQuoteRuneToGraphic", strconv.AppendQuoteRuneToGraphic, execAppendQuoteRuneToGraphic),
		I.Func("AppendQuoteToASCII", strconv.AppendQuoteToASCII, execAppendQuoteToASCII),
		I.Func("AppendQuoteToGraphic", strconv.AppendQuoteToGraphic, execAppendQuoteToGraphic),
		I.Func("AppendUint", strconv.AppendUint, execAppendUint),
		I.Func("Atoi", strconv.Atoi, execAtoi),
		I.Func("CanBackquote", strconv.CanBackquote, execCanBackquote),
		I.Func("FormatBool", strconv.FormatBool, execFormatBool),
		I.Func("FormatFloat", strconv.FormatFloat, execFormatFloat),
		I.Func("FormatInt", strconv.FormatInt, execFormatInt),
		I.Func("FormatUint", strconv.FormatUint, execFormatUint),
		I.Func("IsGraphic", strconv.IsGraphic, execIsGraphic),
		I.Func("IsPrint", strconv.IsPrint, execIsPrint),
		I.Func("Itoa", strconv.Itoa, execItoa),
		I.Func("(*NumError).Error", (*strconv.NumError).Error, execmNumErrorError),
		I.Func("(*NumError).Unwrap", (*strconv.NumError).Unwrap, execmNumErrorUnwrap),
		I.Func("ParseBool", strconv.ParseBool, execParseBool),
		I.Func("ParseFloat", strconv.ParseFloat, execParseFloat),
		I.Func("ParseInt", strconv.ParseInt, execParseInt),
		I.Func("ParseUint", strconv.ParseUint, execParseUint),
		I.Func("Quote", strconv.Quote, execQuote),
		I.Func("QuoteRune", strconv.QuoteRune, execQuoteRune),
		I.Func("QuoteRuneToASCII", strconv.QuoteRuneToASCII, execQuoteRuneToASCII),
		I.Func("QuoteRuneToGraphic", strconv.QuoteRuneToGraphic, execQuoteRuneToGraphic),
		I.Func("QuoteToASCII", strconv.QuoteToASCII, execQuoteToASCII),
		I.Func("QuoteToGraphic", strconv.QuoteToGraphic, execQuoteToGraphic),
		I.Func("Unquote", strconv.Unquote, execUnquote),
		I.Func("UnquoteChar", strconv.UnquoteChar, execUnquoteChar),
	)
	I.RegisterVars(
		I.Var("ErrRange", &strconv.ErrRange),
		I.Var("ErrSyntax", &strconv.ErrSyntax),
	)
	I.RegisterTypes(
		I.Type("NumError", reflect.TypeOf((*strconv.NumError)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("IntSize", qspec.ConstUnboundInt, strconv.IntSize),
	)
}
