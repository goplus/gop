// Package smtp provide Go+ "net/smtp" package, as "net/smtp" package in Go.
package smtp

import (
	tls "crypto/tls"
	net "net"
	smtp "net/smtp"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func execiAuthNext(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(smtp.Auth).Next(args[1].([]byte), args[2].(bool))
	p.Ret(3, ret0, ret1)
}

func execiAuthStart(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1, ret2 := args[0].(smtp.Auth).Start(args[1].(*smtp.ServerInfo))
	p.Ret(2, ret0, ret1, ret2)
}

func execCRAMMD5Auth(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := smtp.CRAMMD5Auth(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execmClientClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*smtp.Client).Close()
	p.Ret(1, ret0)
}

func execmClientHello(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*smtp.Client).Hello(args[1].(string))
	p.Ret(2, ret0)
}

func execmClientStartTLS(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*smtp.Client).StartTLS(args[1].(*tls.Config))
	p.Ret(2, ret0)
}

func execmClientTLSConnectionState(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*smtp.Client).TLSConnectionState()
	p.Ret(1, ret0, ret1)
}

func execmClientVerify(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*smtp.Client).Verify(args[1].(string))
	p.Ret(2, ret0)
}

func toType0(v interface{}) smtp.Auth {
	if v == nil {
		return nil
	}
	return v.(smtp.Auth)
}

func execmClientAuth(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*smtp.Client).Auth(toType0(args[1]))
	p.Ret(2, ret0)
}

func execmClientMail(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*smtp.Client).Mail(args[1].(string))
	p.Ret(2, ret0)
}

func execmClientRcpt(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*smtp.Client).Rcpt(args[1].(string))
	p.Ret(2, ret0)
}

func execmClientData(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*smtp.Client).Data()
	p.Ret(1, ret0, ret1)
}

func execmClientExtension(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*smtp.Client).Extension(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmClientReset(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*smtp.Client).Reset()
	p.Ret(1, ret0)
}

func execmClientNoop(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*smtp.Client).Noop()
	p.Ret(1, ret0)
}

func execmClientQuit(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*smtp.Client).Quit()
	p.Ret(1, ret0)
}

func execDial(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := smtp.Dial(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func toType1(v interface{}) net.Conn {
	if v == nil {
		return nil
	}
	return v.(net.Conn)
}

func execNewClient(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := smtp.NewClient(toType1(args[0]), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execPlainAuth(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := smtp.PlainAuth(args[0].(string), args[1].(string), args[2].(string), args[3].(string))
	p.Ret(4, ret0)
}

func execSendMail(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	ret0 := smtp.SendMail(args[0].(string), toType0(args[1]), args[2].(string), args[3].([]string), args[4].([]byte))
	p.Ret(5, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("net/smtp")

func init() {
	I.RegisterFuncs(
		I.Func("(Auth).Next", (smtp.Auth).Next, execiAuthNext),
		I.Func("(Auth).Start", (smtp.Auth).Start, execiAuthStart),
		I.Func("CRAMMD5Auth", smtp.CRAMMD5Auth, execCRAMMD5Auth),
		I.Func("(*Client).Close", (*smtp.Client).Close, execmClientClose),
		I.Func("(*Client).Hello", (*smtp.Client).Hello, execmClientHello),
		I.Func("(*Client).StartTLS", (*smtp.Client).StartTLS, execmClientStartTLS),
		I.Func("(*Client).TLSConnectionState", (*smtp.Client).TLSConnectionState, execmClientTLSConnectionState),
		I.Func("(*Client).Verify", (*smtp.Client).Verify, execmClientVerify),
		I.Func("(*Client).Auth", (*smtp.Client).Auth, execmClientAuth),
		I.Func("(*Client).Mail", (*smtp.Client).Mail, execmClientMail),
		I.Func("(*Client).Rcpt", (*smtp.Client).Rcpt, execmClientRcpt),
		I.Func("(*Client).Data", (*smtp.Client).Data, execmClientData),
		I.Func("(*Client).Extension", (*smtp.Client).Extension, execmClientExtension),
		I.Func("(*Client).Reset", (*smtp.Client).Reset, execmClientReset),
		I.Func("(*Client).Noop", (*smtp.Client).Noop, execmClientNoop),
		I.Func("(*Client).Quit", (*smtp.Client).Quit, execmClientQuit),
		I.Func("Dial", smtp.Dial, execDial),
		I.Func("NewClient", smtp.NewClient, execNewClient),
		I.Func("PlainAuth", smtp.PlainAuth, execPlainAuth),
		I.Func("SendMail", smtp.SendMail, execSendMail),
	)
	I.RegisterTypes(
		I.Type("Auth", reflect.TypeOf((*smtp.Auth)(nil)).Elem()),
		I.Type("Client", reflect.TypeOf((*smtp.Client)(nil)).Elem()),
		I.Type("ServerInfo", reflect.TypeOf((*smtp.ServerInfo)(nil)).Elem()),
	)
}
