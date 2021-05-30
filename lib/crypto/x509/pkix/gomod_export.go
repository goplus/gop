// Package pkix provide Go+ "crypto/x509/pkix" package, as "crypto/x509/pkix" package in Go.
package pkix

import (
	pkix "crypto/x509/pkix"
	reflect "reflect"
	time "time"

	gop "github.com/goplus/gop"
)

func execmCertificateListHasExpired(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*pkix.CertificateList).HasExpired(args[1].(time.Time))
	p.Ret(2, ret0)
}

func execmNameFillFromRDNSequence(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*pkix.Name).FillFromRDNSequence(args[1].(*pkix.RDNSequence))
	p.PopN(2)
}

func execmNameToRDNSequence(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(pkix.Name).ToRDNSequence()
	p.Ret(1, ret0)
}

func execmNameString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(pkix.Name).String()
	p.Ret(1, ret0)
}

func execmRDNSequenceString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(pkix.RDNSequence).String()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("crypto/x509/pkix")

func init() {
	I.RegisterFuncs(
		I.Func("(*CertificateList).HasExpired", (*pkix.CertificateList).HasExpired, execmCertificateListHasExpired),
		I.Func("(*Name).FillFromRDNSequence", (*pkix.Name).FillFromRDNSequence, execmNameFillFromRDNSequence),
		I.Func("(Name).ToRDNSequence", (pkix.Name).ToRDNSequence, execmNameToRDNSequence),
		I.Func("(Name).String", (pkix.Name).String, execmNameString),
		I.Func("(RDNSequence).String", (pkix.RDNSequence).String, execmRDNSequenceString),
	)
	I.RegisterTypes(
		I.Type("AlgorithmIdentifier", reflect.TypeOf((*pkix.AlgorithmIdentifier)(nil)).Elem()),
		I.Type("AttributeTypeAndValue", reflect.TypeOf((*pkix.AttributeTypeAndValue)(nil)).Elem()),
		I.Type("AttributeTypeAndValueSET", reflect.TypeOf((*pkix.AttributeTypeAndValueSET)(nil)).Elem()),
		I.Type("CertificateList", reflect.TypeOf((*pkix.CertificateList)(nil)).Elem()),
		I.Type("Extension", reflect.TypeOf((*pkix.Extension)(nil)).Elem()),
		I.Type("Name", reflect.TypeOf((*pkix.Name)(nil)).Elem()),
		I.Type("RDNSequence", reflect.TypeOf((*pkix.RDNSequence)(nil)).Elem()),
		I.Type("RelativeDistinguishedNameSET", reflect.TypeOf((*pkix.RelativeDistinguishedNameSET)(nil)).Elem()),
		I.Type("RevokedCertificate", reflect.TypeOf((*pkix.RevokedCertificate)(nil)).Elem()),
		I.Type("TBSCertificateList", reflect.TypeOf((*pkix.TBSCertificateList)(nil)).Elem()),
	)
}
