// Package mail provide Go+ "net/mail" package, as "net/mail" package in Go.
package mail

import (
	io "io"
	mail "net/mail"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func execmAddressString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*mail.Address).String()
	p.Ret(1, ret0)
}

func execmAddressParserParse(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*mail.AddressParser).Parse(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmAddressParserParseList(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*mail.AddressParser).ParseList(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmHeaderGet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(mail.Header).Get(args[1].(string))
	p.Ret(2, ret0)
}

func execmHeaderDate(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(mail.Header).Date()
	p.Ret(1, ret0, ret1)
}

func execmHeaderAddressList(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(mail.Header).AddressList(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execParseAddress(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := mail.ParseAddress(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execParseAddressList(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := mail.ParseAddressList(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execParseDate(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := mail.ParseDate(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func toType0(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execReadMessage(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := mail.ReadMessage(toType0(args[0]))
	p.Ret(1, ret0, ret1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("net/mail")

func init() {
	I.RegisterFuncs(
		I.Func("(*Address).String", (*mail.Address).String, execmAddressString),
		I.Func("(*AddressParser).Parse", (*mail.AddressParser).Parse, execmAddressParserParse),
		I.Func("(*AddressParser).ParseList", (*mail.AddressParser).ParseList, execmAddressParserParseList),
		I.Func("(Header).Get", (mail.Header).Get, execmHeaderGet),
		I.Func("(Header).Date", (mail.Header).Date, execmHeaderDate),
		I.Func("(Header).AddressList", (mail.Header).AddressList, execmHeaderAddressList),
		I.Func("ParseAddress", mail.ParseAddress, execParseAddress),
		I.Func("ParseAddressList", mail.ParseAddressList, execParseAddressList),
		I.Func("ParseDate", mail.ParseDate, execParseDate),
		I.Func("ReadMessage", mail.ReadMessage, execReadMessage),
	)
	I.RegisterVars(
		I.Var("ErrHeaderNotPresent", &mail.ErrHeaderNotPresent),
	)
	I.RegisterTypes(
		I.Type("Address", reflect.TypeOf((*mail.Address)(nil)).Elem()),
		I.Type("AddressParser", reflect.TypeOf((*mail.AddressParser)(nil)).Elem()),
		I.Type("Header", reflect.TypeOf((*mail.Header)(nil)).Elem()),
		I.Type("Message", reflect.TypeOf((*mail.Message)(nil)).Elem()),
	)
}
