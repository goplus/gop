// Package tls provide Go+ "crypto/tls" package, as "crypto/tls" package in Go.
package tls

import (
	tls "crypto/tls"
	net "net"
	reflect "reflect"
	time "time"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execmCertificateRequestInfoSupportsCertificate(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*tls.CertificateRequestInfo).SupportsCertificate(args[1].(*tls.Certificate))
	p.Ret(2, ret0)
}

func execCipherSuiteName(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := tls.CipherSuiteName(args[0].(uint16))
	p.Ret(1, ret0)
}

func execCipherSuites(_ int, p *gop.Context) {
	ret0 := tls.CipherSuites()
	p.Ret(0, ret0)
}

func toType0(v interface{}) net.Conn {
	if v == nil {
		return nil
	}
	return v.(net.Conn)
}

func execClient(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := tls.Client(toType0(args[0]), args[1].(*tls.Config))
	p.Ret(2, ret0)
}

func execmClientHelloInfoSupportsCertificate(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*tls.ClientHelloInfo).SupportsCertificate(args[1].(*tls.Certificate))
	p.Ret(2, ret0)
}

func execiClientSessionCacheGet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(tls.ClientSessionCache).Get(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execiClientSessionCachePut(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(tls.ClientSessionCache).Put(args[1].(string), args[2].(*tls.ClientSessionState))
	p.PopN(3)
}

func execmConfigClone(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*tls.Config).Clone()
	p.Ret(1, ret0)
}

func execmConfigSetSessionTicketKeys(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*tls.Config).SetSessionTicketKeys(args[1].([][32]byte))
	p.PopN(2)
}

func execmConfigBuildNameToCertificate(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*tls.Config).BuildNameToCertificate()
	p.PopN(1)
}

func execmConnLocalAddr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*tls.Conn).LocalAddr()
	p.Ret(1, ret0)
}

func execmConnRemoteAddr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*tls.Conn).RemoteAddr()
	p.Ret(1, ret0)
}

func execmConnSetDeadline(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*tls.Conn).SetDeadline(args[1].(time.Time))
	p.Ret(2, ret0)
}

func execmConnSetReadDeadline(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*tls.Conn).SetReadDeadline(args[1].(time.Time))
	p.Ret(2, ret0)
}

func execmConnSetWriteDeadline(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*tls.Conn).SetWriteDeadline(args[1].(time.Time))
	p.Ret(2, ret0)
}

func execmConnWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*tls.Conn).Write(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmConnRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*tls.Conn).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmConnClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*tls.Conn).Close()
	p.Ret(1, ret0)
}

func execmConnCloseWrite(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*tls.Conn).CloseWrite()
	p.Ret(1, ret0)
}

func execmConnHandshake(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*tls.Conn).Handshake()
	p.Ret(1, ret0)
}

func execmConnConnectionState(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*tls.Conn).ConnectionState()
	p.Ret(1, ret0)
}

func execmConnOCSPResponse(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*tls.Conn).OCSPResponse()
	p.Ret(1, ret0)
}

func execmConnVerifyHostname(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*tls.Conn).VerifyHostname(args[1].(string))
	p.Ret(2, ret0)
}

func execmConnectionStateExportKeyingMaterial(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := args[0].(*tls.ConnectionState).ExportKeyingMaterial(args[1].(string), args[2].([]byte), args[3].(int))
	p.Ret(4, ret0, ret1)
}

func execDial(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := tls.Dial(args[0].(string), args[1].(string), args[2].(*tls.Config))
	p.Ret(3, ret0, ret1)
}

func execDialWithDialer(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := tls.DialWithDialer(args[0].(*net.Dialer), args[1].(string), args[2].(string), args[3].(*tls.Config))
	p.Ret(4, ret0, ret1)
}

func execInsecureCipherSuites(_ int, p *gop.Context) {
	ret0 := tls.InsecureCipherSuites()
	p.Ret(0, ret0)
}

func execListen(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := tls.Listen(args[0].(string), args[1].(string), args[2].(*tls.Config))
	p.Ret(3, ret0, ret1)
}

func execLoadX509KeyPair(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := tls.LoadX509KeyPair(args[0].(string), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execNewLRUClientSessionCache(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := tls.NewLRUClientSessionCache(args[0].(int))
	p.Ret(1, ret0)
}

func toType1(v interface{}) net.Listener {
	if v == nil {
		return nil
	}
	return v.(net.Listener)
}

func execNewListener(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := tls.NewListener(toType1(args[0]), args[1].(*tls.Config))
	p.Ret(2, ret0)
}

func execmRecordHeaderErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(tls.RecordHeaderError).Error()
	p.Ret(1, ret0)
}

func execServer(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := tls.Server(toType0(args[0]), args[1].(*tls.Config))
	p.Ret(2, ret0)
}

func execX509KeyPair(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := tls.X509KeyPair(args[0].([]byte), args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("crypto/tls")

func init() {
	I.RegisterFuncs(
		I.Func("(*CertificateRequestInfo).SupportsCertificate", (*tls.CertificateRequestInfo).SupportsCertificate, execmCertificateRequestInfoSupportsCertificate),
		I.Func("CipherSuiteName", tls.CipherSuiteName, execCipherSuiteName),
		I.Func("CipherSuites", tls.CipherSuites, execCipherSuites),
		I.Func("Client", tls.Client, execClient),
		I.Func("(*ClientHelloInfo).SupportsCertificate", (*tls.ClientHelloInfo).SupportsCertificate, execmClientHelloInfoSupportsCertificate),
		I.Func("(ClientSessionCache).Get", (tls.ClientSessionCache).Get, execiClientSessionCacheGet),
		I.Func("(ClientSessionCache).Put", (tls.ClientSessionCache).Put, execiClientSessionCachePut),
		I.Func("(*Config).Clone", (*tls.Config).Clone, execmConfigClone),
		I.Func("(*Config).SetSessionTicketKeys", (*tls.Config).SetSessionTicketKeys, execmConfigSetSessionTicketKeys),
		I.Func("(*Config).BuildNameToCertificate", (*tls.Config).BuildNameToCertificate, execmConfigBuildNameToCertificate),
		I.Func("(*Conn).LocalAddr", (*tls.Conn).LocalAddr, execmConnLocalAddr),
		I.Func("(*Conn).RemoteAddr", (*tls.Conn).RemoteAddr, execmConnRemoteAddr),
		I.Func("(*Conn).SetDeadline", (*tls.Conn).SetDeadline, execmConnSetDeadline),
		I.Func("(*Conn).SetReadDeadline", (*tls.Conn).SetReadDeadline, execmConnSetReadDeadline),
		I.Func("(*Conn).SetWriteDeadline", (*tls.Conn).SetWriteDeadline, execmConnSetWriteDeadline),
		I.Func("(*Conn).Write", (*tls.Conn).Write, execmConnWrite),
		I.Func("(*Conn).Read", (*tls.Conn).Read, execmConnRead),
		I.Func("(*Conn).Close", (*tls.Conn).Close, execmConnClose),
		I.Func("(*Conn).CloseWrite", (*tls.Conn).CloseWrite, execmConnCloseWrite),
		I.Func("(*Conn).Handshake", (*tls.Conn).Handshake, execmConnHandshake),
		I.Func("(*Conn).ConnectionState", (*tls.Conn).ConnectionState, execmConnConnectionState),
		I.Func("(*Conn).OCSPResponse", (*tls.Conn).OCSPResponse, execmConnOCSPResponse),
		I.Func("(*Conn).VerifyHostname", (*tls.Conn).VerifyHostname, execmConnVerifyHostname),
		I.Func("(*ConnectionState).ExportKeyingMaterial", (*tls.ConnectionState).ExportKeyingMaterial, execmConnectionStateExportKeyingMaterial),
		I.Func("Dial", tls.Dial, execDial),
		I.Func("DialWithDialer", tls.DialWithDialer, execDialWithDialer),
		I.Func("InsecureCipherSuites", tls.InsecureCipherSuites, execInsecureCipherSuites),
		I.Func("Listen", tls.Listen, execListen),
		I.Func("LoadX509KeyPair", tls.LoadX509KeyPair, execLoadX509KeyPair),
		I.Func("NewLRUClientSessionCache", tls.NewLRUClientSessionCache, execNewLRUClientSessionCache),
		I.Func("NewListener", tls.NewListener, execNewListener),
		I.Func("(RecordHeaderError).Error", (tls.RecordHeaderError).Error, execmRecordHeaderErrorError),
		I.Func("Server", tls.Server, execServer),
		I.Func("X509KeyPair", tls.X509KeyPair, execX509KeyPair),
	)
	I.RegisterTypes(
		I.Type("Certificate", reflect.TypeOf((*tls.Certificate)(nil)).Elem()),
		I.Type("CertificateRequestInfo", reflect.TypeOf((*tls.CertificateRequestInfo)(nil)).Elem()),
		I.Type("CipherSuite", reflect.TypeOf((*tls.CipherSuite)(nil)).Elem()),
		I.Type("ClientAuthType", reflect.TypeOf((*tls.ClientAuthType)(nil)).Elem()),
		I.Type("ClientHelloInfo", reflect.TypeOf((*tls.ClientHelloInfo)(nil)).Elem()),
		I.Type("ClientSessionCache", reflect.TypeOf((*tls.ClientSessionCache)(nil)).Elem()),
		I.Type("ClientSessionState", reflect.TypeOf((*tls.ClientSessionState)(nil)).Elem()),
		I.Type("Config", reflect.TypeOf((*tls.Config)(nil)).Elem()),
		I.Type("Conn", reflect.TypeOf((*tls.Conn)(nil)).Elem()),
		I.Type("ConnectionState", reflect.TypeOf((*tls.ConnectionState)(nil)).Elem()),
		I.Type("CurveID", reflect.TypeOf((*tls.CurveID)(nil)).Elem()),
		I.Type("RecordHeaderError", reflect.TypeOf((*tls.RecordHeaderError)(nil)).Elem()),
		I.Type("RenegotiationSupport", reflect.TypeOf((*tls.RenegotiationSupport)(nil)).Elem()),
		I.Type("SignatureScheme", reflect.TypeOf((*tls.SignatureScheme)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("CurveP256", qspec.Uint16, tls.CurveP256),
		I.Const("CurveP384", qspec.Uint16, tls.CurveP384),
		I.Const("CurveP521", qspec.Uint16, tls.CurveP521),
		I.Const("ECDSAWithP256AndSHA256", qspec.Uint16, tls.ECDSAWithP256AndSHA256),
		I.Const("ECDSAWithP384AndSHA384", qspec.Uint16, tls.ECDSAWithP384AndSHA384),
		I.Const("ECDSAWithP521AndSHA512", qspec.Uint16, tls.ECDSAWithP521AndSHA512),
		I.Const("ECDSAWithSHA1", qspec.Uint16, tls.ECDSAWithSHA1),
		I.Const("Ed25519", qspec.Uint16, tls.Ed25519),
		I.Const("NoClientCert", qspec.Int, tls.NoClientCert),
		I.Const("PKCS1WithSHA1", qspec.Uint16, tls.PKCS1WithSHA1),
		I.Const("PKCS1WithSHA256", qspec.Uint16, tls.PKCS1WithSHA256),
		I.Const("PKCS1WithSHA384", qspec.Uint16, tls.PKCS1WithSHA384),
		I.Const("PKCS1WithSHA512", qspec.Uint16, tls.PKCS1WithSHA512),
		I.Const("PSSWithSHA256", qspec.Uint16, tls.PSSWithSHA256),
		I.Const("PSSWithSHA384", qspec.Uint16, tls.PSSWithSHA384),
		I.Const("PSSWithSHA512", qspec.Uint16, tls.PSSWithSHA512),
		I.Const("RenegotiateFreelyAsClient", qspec.Int, tls.RenegotiateFreelyAsClient),
		I.Const("RenegotiateNever", qspec.Int, tls.RenegotiateNever),
		I.Const("RenegotiateOnceAsClient", qspec.Int, tls.RenegotiateOnceAsClient),
		I.Const("RequestClientCert", qspec.Int, tls.RequestClientCert),
		I.Const("RequireAndVerifyClientCert", qspec.Int, tls.RequireAndVerifyClientCert),
		I.Const("RequireAnyClientCert", qspec.Int, tls.RequireAnyClientCert),
		I.Const("TLS_AES_128_GCM_SHA256", qspec.Uint16, tls.TLS_AES_128_GCM_SHA256),
		I.Const("TLS_AES_256_GCM_SHA384", qspec.Uint16, tls.TLS_AES_256_GCM_SHA384),
		I.Const("TLS_CHACHA20_POLY1305_SHA256", qspec.Uint16, tls.TLS_CHACHA20_POLY1305_SHA256),
		I.Const("TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA", qspec.Uint16, tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA),
		I.Const("TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256", qspec.Uint16, tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256),
		I.Const("TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256", qspec.Uint16, tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256),
		I.Const("TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA", qspec.Uint16, tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA),
		I.Const("TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384", qspec.Uint16, tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384),
		I.Const("TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305", qspec.Uint16, tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305),
		I.Const("TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256", qspec.Uint16, tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256),
		I.Const("TLS_ECDHE_ECDSA_WITH_RC4_128_SHA", qspec.Uint16, tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA),
		I.Const("TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA", qspec.Uint16, tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA),
		I.Const("TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA", qspec.Uint16, tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA),
		I.Const("TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256", qspec.Uint16, tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256),
		I.Const("TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256", qspec.Uint16, tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256),
		I.Const("TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA", qspec.Uint16, tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA),
		I.Const("TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384", qspec.Uint16, tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384),
		I.Const("TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305", qspec.Uint16, tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305),
		I.Const("TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256", qspec.Uint16, tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256),
		I.Const("TLS_ECDHE_RSA_WITH_RC4_128_SHA", qspec.Uint16, tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA),
		I.Const("TLS_FALLBACK_SCSV", qspec.Uint16, tls.TLS_FALLBACK_SCSV),
		I.Const("TLS_RSA_WITH_3DES_EDE_CBC_SHA", qspec.Uint16, tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA),
		I.Const("TLS_RSA_WITH_AES_128_CBC_SHA", qspec.Uint16, tls.TLS_RSA_WITH_AES_128_CBC_SHA),
		I.Const("TLS_RSA_WITH_AES_128_CBC_SHA256", qspec.Uint16, tls.TLS_RSA_WITH_AES_128_CBC_SHA256),
		I.Const("TLS_RSA_WITH_AES_128_GCM_SHA256", qspec.Uint16, tls.TLS_RSA_WITH_AES_128_GCM_SHA256),
		I.Const("TLS_RSA_WITH_AES_256_CBC_SHA", qspec.Uint16, tls.TLS_RSA_WITH_AES_256_CBC_SHA),
		I.Const("TLS_RSA_WITH_AES_256_GCM_SHA384", qspec.Uint16, tls.TLS_RSA_WITH_AES_256_GCM_SHA384),
		I.Const("TLS_RSA_WITH_RC4_128_SHA", qspec.Uint16, tls.TLS_RSA_WITH_RC4_128_SHA),
		I.Const("VerifyClientCertIfGiven", qspec.Int, tls.VerifyClientCertIfGiven),
		I.Const("VersionSSL30", qspec.ConstUnboundInt, tls.VersionSSL30),
		I.Const("VersionTLS10", qspec.ConstUnboundInt, tls.VersionTLS10),
		I.Const("VersionTLS11", qspec.ConstUnboundInt, tls.VersionTLS11),
		I.Const("VersionTLS12", qspec.ConstUnboundInt, tls.VersionTLS12),
		I.Const("VersionTLS13", qspec.ConstUnboundInt, tls.VersionTLS13),
		I.Const("X25519", qspec.Uint16, tls.X25519),
	)
}
