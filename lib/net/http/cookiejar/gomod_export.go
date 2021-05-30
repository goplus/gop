// Package cookiejar provide Go+ "net/http/cookiejar" package, as "net/http/cookiejar" package in Go.
package cookiejar

import (
	http "net/http"
	cookiejar "net/http/cookiejar"
	url "net/url"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func execmJarCookies(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*cookiejar.Jar).Cookies(args[1].(*url.URL))
	p.Ret(2, ret0)
}

func execmJarSetCookies(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*cookiejar.Jar).SetCookies(args[1].(*url.URL), args[2].([]*http.Cookie))
	p.PopN(3)
}

func execNew(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := cookiejar.New(args[0].(*cookiejar.Options))
	p.Ret(1, ret0, ret1)
}

func execiPublicSuffixListPublicSuffix(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(cookiejar.PublicSuffixList).PublicSuffix(args[1].(string))
	p.Ret(2, ret0)
}

func execiPublicSuffixListString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(cookiejar.PublicSuffixList).String()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("net/http/cookiejar")

func init() {
	I.RegisterFuncs(
		I.Func("(*Jar).Cookies", (*cookiejar.Jar).Cookies, execmJarCookies),
		I.Func("(*Jar).SetCookies", (*cookiejar.Jar).SetCookies, execmJarSetCookies),
		I.Func("New", cookiejar.New, execNew),
		I.Func("(PublicSuffixList).PublicSuffix", (cookiejar.PublicSuffixList).PublicSuffix, execiPublicSuffixListPublicSuffix),
		I.Func("(PublicSuffixList).String", (cookiejar.PublicSuffixList).String, execiPublicSuffixListString),
	)
	I.RegisterTypes(
		I.Type("Jar", reflect.TypeOf((*cookiejar.Jar)(nil)).Elem()),
		I.Type("Options", reflect.TypeOf((*cookiejar.Options)(nil)).Elem()),
		I.Type("PublicSuffixList", reflect.TypeOf((*cookiejar.PublicSuffixList)(nil)).Elem()),
	)
}
