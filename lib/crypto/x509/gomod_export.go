// Package x509 provide Go+ "crypto/x509" package, as "crypto/x509" package in Go.
package x509

import (
	ecdsa "crypto/ecdsa"
	rsa "crypto/rsa"
	x509 "crypto/x509"
	pkix "crypto/x509/pkix"
	pem "encoding/pem"
	io "io"
	reflect "reflect"
	time "time"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execmCertPoolAddCert(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*x509.CertPool).AddCert(args[1].(*x509.Certificate))
	p.PopN(2)
}

func execmCertPoolAppendCertsFromPEM(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*x509.CertPool).AppendCertsFromPEM(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmCertPoolSubjects(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*x509.CertPool).Subjects()
	p.Ret(1, ret0)
}

func execmCertificateVerify(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*x509.Certificate).Verify(args[1].(x509.VerifyOptions))
	p.Ret(2, ret0, ret1)
}

func execmCertificateVerifyHostname(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*x509.Certificate).VerifyHostname(args[1].(string))
	p.Ret(2, ret0)
}

func execmCertificateEqual(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*x509.Certificate).Equal(args[1].(*x509.Certificate))
	p.Ret(2, ret0)
}

func execmCertificateCheckSignatureFrom(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*x509.Certificate).CheckSignatureFrom(args[1].(*x509.Certificate))
	p.Ret(2, ret0)
}

func execmCertificateCheckSignature(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := args[0].(*x509.Certificate).CheckSignature(args[1].(x509.SignatureAlgorithm), args[2].([]byte), args[3].([]byte))
	p.Ret(4, ret0)
}

func execmCertificateCheckCRLSignature(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*x509.Certificate).CheckCRLSignature(args[1].(*pkix.CertificateList))
	p.Ret(2, ret0)
}

func toType0(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execmCertificateCreateCRL(_ int, p *gop.Context) {
	args := p.GetArgs(6)
	ret0, ret1 := args[0].(*x509.Certificate).CreateCRL(toType0(args[1]), args[2], args[3].([]pkix.RevokedCertificate), args[4].(time.Time), args[5].(time.Time))
	p.Ret(6, ret0, ret1)
}

func execmCertificateInvalidErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(x509.CertificateInvalidError).Error()
	p.Ret(1, ret0)
}

func execmCertificateRequestCheckSignature(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*x509.CertificateRequest).CheckSignature()
	p.Ret(1, ret0)
}

func execmConstraintViolationErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(x509.ConstraintViolationError).Error()
	p.Ret(1, ret0)
}

func execCreateCertificate(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	ret0, ret1 := x509.CreateCertificate(toType0(args[0]), args[1].(*x509.Certificate), args[2].(*x509.Certificate), args[3], args[4])
	p.Ret(5, ret0, ret1)
}

func execCreateCertificateRequest(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := x509.CreateCertificateRequest(toType0(args[0]), args[1].(*x509.CertificateRequest), args[2])
	p.Ret(3, ret0, ret1)
}

func execDecryptPEMBlock(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := x509.DecryptPEMBlock(args[0].(*pem.Block), args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execEncryptPEMBlock(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	ret0, ret1 := x509.EncryptPEMBlock(toType0(args[0]), args[1].(string), args[2].([]byte), args[3].([]byte), args[4].(x509.PEMCipher))
	p.Ret(5, ret0, ret1)
}

func execmHostnameErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(x509.HostnameError).Error()
	p.Ret(1, ret0)
}

func execmInsecureAlgorithmErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(x509.InsecureAlgorithmError).Error()
	p.Ret(1, ret0)
}

func execIsEncryptedPEMBlock(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := x509.IsEncryptedPEMBlock(args[0].(*pem.Block))
	p.Ret(1, ret0)
}

func execMarshalECPrivateKey(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := x509.MarshalECPrivateKey(args[0].(*ecdsa.PrivateKey))
	p.Ret(1, ret0, ret1)
}

func execMarshalPKCS1PrivateKey(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := x509.MarshalPKCS1PrivateKey(args[0].(*rsa.PrivateKey))
	p.Ret(1, ret0)
}

func execMarshalPKCS1PublicKey(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := x509.MarshalPKCS1PublicKey(args[0].(*rsa.PublicKey))
	p.Ret(1, ret0)
}

func execMarshalPKCS8PrivateKey(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := x509.MarshalPKCS8PrivateKey(args[0])
	p.Ret(1, ret0, ret1)
}

func execMarshalPKIXPublicKey(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := x509.MarshalPKIXPublicKey(args[0])
	p.Ret(1, ret0, ret1)
}

func execNewCertPool(_ int, p *gop.Context) {
	ret0 := x509.NewCertPool()
	p.Ret(0, ret0)
}

func execParseCRL(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := x509.ParseCRL(args[0].([]byte))
	p.Ret(1, ret0, ret1)
}

func execParseCertificate(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := x509.ParseCertificate(args[0].([]byte))
	p.Ret(1, ret0, ret1)
}

func execParseCertificateRequest(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := x509.ParseCertificateRequest(args[0].([]byte))
	p.Ret(1, ret0, ret1)
}

func execParseCertificates(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := x509.ParseCertificates(args[0].([]byte))
	p.Ret(1, ret0, ret1)
}

func execParseDERCRL(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := x509.ParseDERCRL(args[0].([]byte))
	p.Ret(1, ret0, ret1)
}

func execParseECPrivateKey(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := x509.ParseECPrivateKey(args[0].([]byte))
	p.Ret(1, ret0, ret1)
}

func execParsePKCS1PrivateKey(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := x509.ParsePKCS1PrivateKey(args[0].([]byte))
	p.Ret(1, ret0, ret1)
}

func execParsePKCS1PublicKey(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := x509.ParsePKCS1PublicKey(args[0].([]byte))
	p.Ret(1, ret0, ret1)
}

func execParsePKCS8PrivateKey(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := x509.ParsePKCS8PrivateKey(args[0].([]byte))
	p.Ret(1, ret0, ret1)
}

func execParsePKIXPublicKey(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := x509.ParsePKIXPublicKey(args[0].([]byte))
	p.Ret(1, ret0, ret1)
}

func execmPublicKeyAlgorithmString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(x509.PublicKeyAlgorithm).String()
	p.Ret(1, ret0)
}

func execmSignatureAlgorithmString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(x509.SignatureAlgorithm).String()
	p.Ret(1, ret0)
}

func execSystemCertPool(_ int, p *gop.Context) {
	ret0, ret1 := x509.SystemCertPool()
	p.Ret(0, ret0, ret1)
}

func execmSystemRootsErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(x509.SystemRootsError).Error()
	p.Ret(1, ret0)
}

func execmUnhandledCriticalExtensionError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(x509.UnhandledCriticalExtension).Error()
	p.Ret(1, ret0)
}

func execmUnknownAuthorityErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(x509.UnknownAuthorityError).Error()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("crypto/x509")

func init() {
	I.RegisterFuncs(
		I.Func("(*CertPool).AddCert", (*x509.CertPool).AddCert, execmCertPoolAddCert),
		I.Func("(*CertPool).AppendCertsFromPEM", (*x509.CertPool).AppendCertsFromPEM, execmCertPoolAppendCertsFromPEM),
		I.Func("(*CertPool).Subjects", (*x509.CertPool).Subjects, execmCertPoolSubjects),
		I.Func("(*Certificate).Verify", (*x509.Certificate).Verify, execmCertificateVerify),
		I.Func("(*Certificate).VerifyHostname", (*x509.Certificate).VerifyHostname, execmCertificateVerifyHostname),
		I.Func("(*Certificate).Equal", (*x509.Certificate).Equal, execmCertificateEqual),
		I.Func("(*Certificate).CheckSignatureFrom", (*x509.Certificate).CheckSignatureFrom, execmCertificateCheckSignatureFrom),
		I.Func("(*Certificate).CheckSignature", (*x509.Certificate).CheckSignature, execmCertificateCheckSignature),
		I.Func("(*Certificate).CheckCRLSignature", (*x509.Certificate).CheckCRLSignature, execmCertificateCheckCRLSignature),
		I.Func("(*Certificate).CreateCRL", (*x509.Certificate).CreateCRL, execmCertificateCreateCRL),
		I.Func("(CertificateInvalidError).Error", (x509.CertificateInvalidError).Error, execmCertificateInvalidErrorError),
		I.Func("(*CertificateRequest).CheckSignature", (*x509.CertificateRequest).CheckSignature, execmCertificateRequestCheckSignature),
		I.Func("(ConstraintViolationError).Error", (x509.ConstraintViolationError).Error, execmConstraintViolationErrorError),
		I.Func("CreateCertificate", x509.CreateCertificate, execCreateCertificate),
		I.Func("CreateCertificateRequest", x509.CreateCertificateRequest, execCreateCertificateRequest),
		I.Func("DecryptPEMBlock", x509.DecryptPEMBlock, execDecryptPEMBlock),
		I.Func("EncryptPEMBlock", x509.EncryptPEMBlock, execEncryptPEMBlock),
		I.Func("(HostnameError).Error", (x509.HostnameError).Error, execmHostnameErrorError),
		I.Func("(InsecureAlgorithmError).Error", (x509.InsecureAlgorithmError).Error, execmInsecureAlgorithmErrorError),
		I.Func("IsEncryptedPEMBlock", x509.IsEncryptedPEMBlock, execIsEncryptedPEMBlock),
		I.Func("MarshalECPrivateKey", x509.MarshalECPrivateKey, execMarshalECPrivateKey),
		I.Func("MarshalPKCS1PrivateKey", x509.MarshalPKCS1PrivateKey, execMarshalPKCS1PrivateKey),
		I.Func("MarshalPKCS1PublicKey", x509.MarshalPKCS1PublicKey, execMarshalPKCS1PublicKey),
		I.Func("MarshalPKCS8PrivateKey", x509.MarshalPKCS8PrivateKey, execMarshalPKCS8PrivateKey),
		I.Func("MarshalPKIXPublicKey", x509.MarshalPKIXPublicKey, execMarshalPKIXPublicKey),
		I.Func("NewCertPool", x509.NewCertPool, execNewCertPool),
		I.Func("ParseCRL", x509.ParseCRL, execParseCRL),
		I.Func("ParseCertificate", x509.ParseCertificate, execParseCertificate),
		I.Func("ParseCertificateRequest", x509.ParseCertificateRequest, execParseCertificateRequest),
		I.Func("ParseCertificates", x509.ParseCertificates, execParseCertificates),
		I.Func("ParseDERCRL", x509.ParseDERCRL, execParseDERCRL),
		I.Func("ParseECPrivateKey", x509.ParseECPrivateKey, execParseECPrivateKey),
		I.Func("ParsePKCS1PrivateKey", x509.ParsePKCS1PrivateKey, execParsePKCS1PrivateKey),
		I.Func("ParsePKCS1PublicKey", x509.ParsePKCS1PublicKey, execParsePKCS1PublicKey),
		I.Func("ParsePKCS8PrivateKey", x509.ParsePKCS8PrivateKey, execParsePKCS8PrivateKey),
		I.Func("ParsePKIXPublicKey", x509.ParsePKIXPublicKey, execParsePKIXPublicKey),
		I.Func("(PublicKeyAlgorithm).String", (x509.PublicKeyAlgorithm).String, execmPublicKeyAlgorithmString),
		I.Func("(SignatureAlgorithm).String", (x509.SignatureAlgorithm).String, execmSignatureAlgorithmString),
		I.Func("SystemCertPool", x509.SystemCertPool, execSystemCertPool),
		I.Func("(SystemRootsError).Error", (x509.SystemRootsError).Error, execmSystemRootsErrorError),
		I.Func("(UnhandledCriticalExtension).Error", (x509.UnhandledCriticalExtension).Error, execmUnhandledCriticalExtensionError),
		I.Func("(UnknownAuthorityError).Error", (x509.UnknownAuthorityError).Error, execmUnknownAuthorityErrorError),
	)
	I.RegisterVars(
		I.Var("ErrUnsupportedAlgorithm", &x509.ErrUnsupportedAlgorithm),
		I.Var("IncorrectPasswordError", &x509.IncorrectPasswordError),
	)
	I.RegisterTypes(
		I.Type("CertPool", reflect.TypeOf((*x509.CertPool)(nil)).Elem()),
		I.Type("Certificate", reflect.TypeOf((*x509.Certificate)(nil)).Elem()),
		I.Type("CertificateInvalidError", reflect.TypeOf((*x509.CertificateInvalidError)(nil)).Elem()),
		I.Type("CertificateRequest", reflect.TypeOf((*x509.CertificateRequest)(nil)).Elem()),
		I.Type("ConstraintViolationError", reflect.TypeOf((*x509.ConstraintViolationError)(nil)).Elem()),
		I.Type("ExtKeyUsage", reflect.TypeOf((*x509.ExtKeyUsage)(nil)).Elem()),
		I.Type("HostnameError", reflect.TypeOf((*x509.HostnameError)(nil)).Elem()),
		I.Type("InsecureAlgorithmError", reflect.TypeOf((*x509.InsecureAlgorithmError)(nil)).Elem()),
		I.Type("InvalidReason", reflect.TypeOf((*x509.InvalidReason)(nil)).Elem()),
		I.Type("KeyUsage", reflect.TypeOf((*x509.KeyUsage)(nil)).Elem()),
		I.Type("PEMCipher", reflect.TypeOf((*x509.PEMCipher)(nil)).Elem()),
		I.Type("PublicKeyAlgorithm", reflect.TypeOf((*x509.PublicKeyAlgorithm)(nil)).Elem()),
		I.Type("SignatureAlgorithm", reflect.TypeOf((*x509.SignatureAlgorithm)(nil)).Elem()),
		I.Type("SystemRootsError", reflect.TypeOf((*x509.SystemRootsError)(nil)).Elem()),
		I.Type("UnhandledCriticalExtension", reflect.TypeOf((*x509.UnhandledCriticalExtension)(nil)).Elem()),
		I.Type("UnknownAuthorityError", reflect.TypeOf((*x509.UnknownAuthorityError)(nil)).Elem()),
		I.Type("VerifyOptions", reflect.TypeOf((*x509.VerifyOptions)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("CANotAuthorizedForExtKeyUsage", qspec.Int, x509.CANotAuthorizedForExtKeyUsage),
		I.Const("CANotAuthorizedForThisName", qspec.Int, x509.CANotAuthorizedForThisName),
		I.Const("DSA", qspec.Int, x509.DSA),
		I.Const("DSAWithSHA1", qspec.Int, x509.DSAWithSHA1),
		I.Const("DSAWithSHA256", qspec.Int, x509.DSAWithSHA256),
		I.Const("ECDSA", qspec.Int, x509.ECDSA),
		I.Const("ECDSAWithSHA1", qspec.Int, x509.ECDSAWithSHA1),
		I.Const("ECDSAWithSHA256", qspec.Int, x509.ECDSAWithSHA256),
		I.Const("ECDSAWithSHA384", qspec.Int, x509.ECDSAWithSHA384),
		I.Const("ECDSAWithSHA512", qspec.Int, x509.ECDSAWithSHA512),
		I.Const("Ed25519", qspec.Int, x509.Ed25519),
		I.Const("Expired", qspec.Int, x509.Expired),
		I.Const("ExtKeyUsageAny", qspec.Int, x509.ExtKeyUsageAny),
		I.Const("ExtKeyUsageClientAuth", qspec.Int, x509.ExtKeyUsageClientAuth),
		I.Const("ExtKeyUsageCodeSigning", qspec.Int, x509.ExtKeyUsageCodeSigning),
		I.Const("ExtKeyUsageEmailProtection", qspec.Int, x509.ExtKeyUsageEmailProtection),
		I.Const("ExtKeyUsageIPSECEndSystem", qspec.Int, x509.ExtKeyUsageIPSECEndSystem),
		I.Const("ExtKeyUsageIPSECTunnel", qspec.Int, x509.ExtKeyUsageIPSECTunnel),
		I.Const("ExtKeyUsageIPSECUser", qspec.Int, x509.ExtKeyUsageIPSECUser),
		I.Const("ExtKeyUsageMicrosoftCommercialCodeSigning", qspec.Int, x509.ExtKeyUsageMicrosoftCommercialCodeSigning),
		I.Const("ExtKeyUsageMicrosoftKernelCodeSigning", qspec.Int, x509.ExtKeyUsageMicrosoftKernelCodeSigning),
		I.Const("ExtKeyUsageMicrosoftServerGatedCrypto", qspec.Int, x509.ExtKeyUsageMicrosoftServerGatedCrypto),
		I.Const("ExtKeyUsageNetscapeServerGatedCrypto", qspec.Int, x509.ExtKeyUsageNetscapeServerGatedCrypto),
		I.Const("ExtKeyUsageOCSPSigning", qspec.Int, x509.ExtKeyUsageOCSPSigning),
		I.Const("ExtKeyUsageServerAuth", qspec.Int, x509.ExtKeyUsageServerAuth),
		I.Const("ExtKeyUsageTimeStamping", qspec.Int, x509.ExtKeyUsageTimeStamping),
		I.Const("IncompatibleUsage", qspec.Int, x509.IncompatibleUsage),
		I.Const("KeyUsageCRLSign", qspec.Int, x509.KeyUsageCRLSign),
		I.Const("KeyUsageCertSign", qspec.Int, x509.KeyUsageCertSign),
		I.Const("KeyUsageContentCommitment", qspec.Int, x509.KeyUsageContentCommitment),
		I.Const("KeyUsageDataEncipherment", qspec.Int, x509.KeyUsageDataEncipherment),
		I.Const("KeyUsageDecipherOnly", qspec.Int, x509.KeyUsageDecipherOnly),
		I.Const("KeyUsageDigitalSignature", qspec.Int, x509.KeyUsageDigitalSignature),
		I.Const("KeyUsageEncipherOnly", qspec.Int, x509.KeyUsageEncipherOnly),
		I.Const("KeyUsageKeyAgreement", qspec.Int, x509.KeyUsageKeyAgreement),
		I.Const("KeyUsageKeyEncipherment", qspec.Int, x509.KeyUsageKeyEncipherment),
		I.Const("MD2WithRSA", qspec.Int, x509.MD2WithRSA),
		I.Const("MD5WithRSA", qspec.Int, x509.MD5WithRSA),
		I.Const("NameConstraintsWithoutSANs", qspec.Int, x509.NameConstraintsWithoutSANs),
		I.Const("NameMismatch", qspec.Int, x509.NameMismatch),
		I.Const("NotAuthorizedToSign", qspec.Int, x509.NotAuthorizedToSign),
		I.Const("PEMCipher3DES", qspec.Int, x509.PEMCipher3DES),
		I.Const("PEMCipherAES128", qspec.Int, x509.PEMCipherAES128),
		I.Const("PEMCipherAES192", qspec.Int, x509.PEMCipherAES192),
		I.Const("PEMCipherAES256", qspec.Int, x509.PEMCipherAES256),
		I.Const("PEMCipherDES", qspec.Int, x509.PEMCipherDES),
		I.Const("PureEd25519", qspec.Int, x509.PureEd25519),
		I.Const("RSA", qspec.Int, x509.RSA),
		I.Const("SHA1WithRSA", qspec.Int, x509.SHA1WithRSA),
		I.Const("SHA256WithRSA", qspec.Int, x509.SHA256WithRSA),
		I.Const("SHA256WithRSAPSS", qspec.Int, x509.SHA256WithRSAPSS),
		I.Const("SHA384WithRSA", qspec.Int, x509.SHA384WithRSA),
		I.Const("SHA384WithRSAPSS", qspec.Int, x509.SHA384WithRSAPSS),
		I.Const("SHA512WithRSA", qspec.Int, x509.SHA512WithRSA),
		I.Const("SHA512WithRSAPSS", qspec.Int, x509.SHA512WithRSAPSS),
		I.Const("TooManyConstraints", qspec.Int, x509.TooManyConstraints),
		I.Const("TooManyIntermediates", qspec.Int, x509.TooManyIntermediates),
		I.Const("UnconstrainedName", qspec.Int, x509.UnconstrainedName),
		I.Const("UnknownPublicKeyAlgorithm", qspec.Int, x509.UnknownPublicKeyAlgorithm),
		I.Const("UnknownSignatureAlgorithm", qspec.Int, x509.UnknownSignatureAlgorithm),
	)
}
