// Package url provide Go+ "net/url" package, as "net/url" package in Go.
package url

import (
	url "net/url"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func execmErrorUnwrap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*url.Error).Unwrap()
	p.Ret(1, ret0)
}

func execmErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*url.Error).Error()
	p.Ret(1, ret0)
}

func execmErrorTimeout(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*url.Error).Timeout()
	p.Ret(1, ret0)
}

func execmErrorTemporary(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*url.Error).Temporary()
	p.Ret(1, ret0)
}

func execmEscapeErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(url.EscapeError).Error()
	p.Ret(1, ret0)
}

func execmInvalidHostErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(url.InvalidHostError).Error()
	p.Ret(1, ret0)
}

func execParse(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := url.Parse(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execParseQuery(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := url.ParseQuery(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execParseRequestURI(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := url.ParseRequestURI(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execPathEscape(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := url.PathEscape(args[0].(string))
	p.Ret(1, ret0)
}

func execPathUnescape(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := url.PathUnescape(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execQueryEscape(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := url.QueryEscape(args[0].(string))
	p.Ret(1, ret0)
}

func execQueryUnescape(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := url.QueryUnescape(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execmURLEscapedPath(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*url.URL).EscapedPath()
	p.Ret(1, ret0)
}

func execmURLString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*url.URL).String()
	p.Ret(1, ret0)
}

func execmURLIsAbs(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*url.URL).IsAbs()
	p.Ret(1, ret0)
}

func execmURLParse(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*url.URL).Parse(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmURLResolveReference(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*url.URL).ResolveReference(args[1].(*url.URL))
	p.Ret(2, ret0)
}

func execmURLQuery(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*url.URL).Query()
	p.Ret(1, ret0)
}

func execmURLRequestURI(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*url.URL).RequestURI()
	p.Ret(1, ret0)
}

func execmURLHostname(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*url.URL).Hostname()
	p.Ret(1, ret0)
}

func execmURLPort(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*url.URL).Port()
	p.Ret(1, ret0)
}

func execmURLMarshalBinary(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*url.URL).MarshalBinary()
	p.Ret(1, ret0, ret1)
}

func execmURLUnmarshalBinary(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*url.URL).UnmarshalBinary(args[1].([]byte))
	p.Ret(2, ret0)
}

func execUser(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := url.User(args[0].(string))
	p.Ret(1, ret0)
}

func execUserPassword(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := url.UserPassword(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execmUserinfoUsername(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*url.Userinfo).Username()
	p.Ret(1, ret0)
}

func execmUserinfoPassword(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*url.Userinfo).Password()
	p.Ret(1, ret0, ret1)
}

func execmUserinfoString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*url.Userinfo).String()
	p.Ret(1, ret0)
}

func execmValuesGet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(url.Values).Get(args[1].(string))
	p.Ret(2, ret0)
}

func execmValuesSet(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(url.Values).Set(args[1].(string), args[2].(string))
	p.PopN(3)
}

func execmValuesAdd(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(url.Values).Add(args[1].(string), args[2].(string))
	p.PopN(3)
}

func execmValuesDel(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(url.Values).Del(args[1].(string))
	p.PopN(2)
}

func execmValuesEncode(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(url.Values).Encode()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("net/url")

func init() {
	I.RegisterFuncs(
		I.Func("(*Error).Unwrap", (*url.Error).Unwrap, execmErrorUnwrap),
		I.Func("(*Error).Error", (*url.Error).Error, execmErrorError),
		I.Func("(*Error).Timeout", (*url.Error).Timeout, execmErrorTimeout),
		I.Func("(*Error).Temporary", (*url.Error).Temporary, execmErrorTemporary),
		I.Func("(EscapeError).Error", (url.EscapeError).Error, execmEscapeErrorError),
		I.Func("(InvalidHostError).Error", (url.InvalidHostError).Error, execmInvalidHostErrorError),
		I.Func("Parse", url.Parse, execParse),
		I.Func("ParseQuery", url.ParseQuery, execParseQuery),
		I.Func("ParseRequestURI", url.ParseRequestURI, execParseRequestURI),
		I.Func("PathEscape", url.PathEscape, execPathEscape),
		I.Func("PathUnescape", url.PathUnescape, execPathUnescape),
		I.Func("QueryEscape", url.QueryEscape, execQueryEscape),
		I.Func("QueryUnescape", url.QueryUnescape, execQueryUnescape),
		I.Func("(*URL).EscapedPath", (*url.URL).EscapedPath, execmURLEscapedPath),
		I.Func("(*URL).String", (*url.URL).String, execmURLString),
		I.Func("(*URL).IsAbs", (*url.URL).IsAbs, execmURLIsAbs),
		I.Func("(*URL).Parse", (*url.URL).Parse, execmURLParse),
		I.Func("(*URL).ResolveReference", (*url.URL).ResolveReference, execmURLResolveReference),
		I.Func("(*URL).Query", (*url.URL).Query, execmURLQuery),
		I.Func("(*URL).RequestURI", (*url.URL).RequestURI, execmURLRequestURI),
		I.Func("(*URL).Hostname", (*url.URL).Hostname, execmURLHostname),
		I.Func("(*URL).Port", (*url.URL).Port, execmURLPort),
		I.Func("(*URL).MarshalBinary", (*url.URL).MarshalBinary, execmURLMarshalBinary),
		I.Func("(*URL).UnmarshalBinary", (*url.URL).UnmarshalBinary, execmURLUnmarshalBinary),
		I.Func("User", url.User, execUser),
		I.Func("UserPassword", url.UserPassword, execUserPassword),
		I.Func("(*Userinfo).Username", (*url.Userinfo).Username, execmUserinfoUsername),
		I.Func("(*Userinfo).Password", (*url.Userinfo).Password, execmUserinfoPassword),
		I.Func("(*Userinfo).String", (*url.Userinfo).String, execmUserinfoString),
		I.Func("(Values).Get", (url.Values).Get, execmValuesGet),
		I.Func("(Values).Set", (url.Values).Set, execmValuesSet),
		I.Func("(Values).Add", (url.Values).Add, execmValuesAdd),
		I.Func("(Values).Del", (url.Values).Del, execmValuesDel),
		I.Func("(Values).Encode", (url.Values).Encode, execmValuesEncode),
	)
	I.RegisterTypes(
		I.Type("Error", reflect.TypeOf((*url.Error)(nil)).Elem()),
		I.Type("EscapeError", reflect.TypeOf((*url.EscapeError)(nil)).Elem()),
		I.Type("InvalidHostError", reflect.TypeOf((*url.InvalidHostError)(nil)).Elem()),
		I.Type("URL", reflect.TypeOf((*url.URL)(nil)).Elem()),
		I.Type("Userinfo", reflect.TypeOf((*url.Userinfo)(nil)).Elem()),
		I.Type("Values", reflect.TypeOf((*url.Values)(nil)).Elem()),
	)
}
