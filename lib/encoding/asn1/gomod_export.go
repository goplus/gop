// Package asn1 provide Go+ "encoding/asn1" package, as "encoding/asn1" package in Go.
package asn1

import (
	asn1 "encoding/asn1"
	reflect "reflect"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execmBitStringAt(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(asn1.BitString).At(args[1].(int))
	p.Ret(2, ret0)
}

func execmBitStringRightAlign(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(asn1.BitString).RightAlign()
	p.Ret(1, ret0)
}

func execMarshal(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := asn1.Marshal(args[0])
	p.Ret(1, ret0, ret1)
}

func execMarshalWithParams(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := asn1.MarshalWithParams(args[0], args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmObjectIdentifierEqual(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(asn1.ObjectIdentifier).Equal(args[1].(asn1.ObjectIdentifier))
	p.Ret(2, ret0)
}

func execmObjectIdentifierString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(asn1.ObjectIdentifier).String()
	p.Ret(1, ret0)
}

func execmStructuralErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(asn1.StructuralError).Error()
	p.Ret(1, ret0)
}

func execmSyntaxErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(asn1.SyntaxError).Error()
	p.Ret(1, ret0)
}

func execUnmarshal(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := asn1.Unmarshal(args[0].([]byte), args[1])
	p.Ret(2, ret0, ret1)
}

func execUnmarshalWithParams(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := asn1.UnmarshalWithParams(args[0].([]byte), args[1], args[2].(string))
	p.Ret(3, ret0, ret1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("encoding/asn1")

func init() {
	I.RegisterFuncs(
		I.Func("(BitString).At", (asn1.BitString).At, execmBitStringAt),
		I.Func("(BitString).RightAlign", (asn1.BitString).RightAlign, execmBitStringRightAlign),
		I.Func("Marshal", asn1.Marshal, execMarshal),
		I.Func("MarshalWithParams", asn1.MarshalWithParams, execMarshalWithParams),
		I.Func("(ObjectIdentifier).Equal", (asn1.ObjectIdentifier).Equal, execmObjectIdentifierEqual),
		I.Func("(ObjectIdentifier).String", (asn1.ObjectIdentifier).String, execmObjectIdentifierString),
		I.Func("(StructuralError).Error", (asn1.StructuralError).Error, execmStructuralErrorError),
		I.Func("(SyntaxError).Error", (asn1.SyntaxError).Error, execmSyntaxErrorError),
		I.Func("Unmarshal", asn1.Unmarshal, execUnmarshal),
		I.Func("UnmarshalWithParams", asn1.UnmarshalWithParams, execUnmarshalWithParams),
	)
	I.RegisterVars(
		I.Var("NullBytes", &asn1.NullBytes),
		I.Var("NullRawValue", &asn1.NullRawValue),
	)
	I.RegisterTypes(
		I.Type("BitString", reflect.TypeOf((*asn1.BitString)(nil)).Elem()),
		I.Type("Enumerated", reflect.TypeOf((*asn1.Enumerated)(nil)).Elem()),
		I.Type("Flag", reflect.TypeOf((*asn1.Flag)(nil)).Elem()),
		I.Type("ObjectIdentifier", reflect.TypeOf((*asn1.ObjectIdentifier)(nil)).Elem()),
		I.Type("RawContent", reflect.TypeOf((*asn1.RawContent)(nil)).Elem()),
		I.Type("RawValue", reflect.TypeOf((*asn1.RawValue)(nil)).Elem()),
		I.Type("StructuralError", reflect.TypeOf((*asn1.StructuralError)(nil)).Elem()),
		I.Type("SyntaxError", reflect.TypeOf((*asn1.SyntaxError)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("ClassApplication", qspec.ConstUnboundInt, asn1.ClassApplication),
		I.Const("ClassContextSpecific", qspec.ConstUnboundInt, asn1.ClassContextSpecific),
		I.Const("ClassPrivate", qspec.ConstUnboundInt, asn1.ClassPrivate),
		I.Const("ClassUniversal", qspec.ConstUnboundInt, asn1.ClassUniversal),
		I.Const("TagBMPString", qspec.ConstUnboundInt, asn1.TagBMPString),
		I.Const("TagBitString", qspec.ConstUnboundInt, asn1.TagBitString),
		I.Const("TagBoolean", qspec.ConstUnboundInt, asn1.TagBoolean),
		I.Const("TagEnum", qspec.ConstUnboundInt, asn1.TagEnum),
		I.Const("TagGeneralString", qspec.ConstUnboundInt, asn1.TagGeneralString),
		I.Const("TagGeneralizedTime", qspec.ConstUnboundInt, asn1.TagGeneralizedTime),
		I.Const("TagIA5String", qspec.ConstUnboundInt, asn1.TagIA5String),
		I.Const("TagInteger", qspec.ConstUnboundInt, asn1.TagInteger),
		I.Const("TagNull", qspec.ConstUnboundInt, asn1.TagNull),
		I.Const("TagNumericString", qspec.ConstUnboundInt, asn1.TagNumericString),
		I.Const("TagOID", qspec.ConstUnboundInt, asn1.TagOID),
		I.Const("TagOctetString", qspec.ConstUnboundInt, asn1.TagOctetString),
		I.Const("TagPrintableString", qspec.ConstUnboundInt, asn1.TagPrintableString),
		I.Const("TagSequence", qspec.ConstUnboundInt, asn1.TagSequence),
		I.Const("TagSet", qspec.ConstUnboundInt, asn1.TagSet),
		I.Const("TagT61String", qspec.ConstUnboundInt, asn1.TagT61String),
		I.Const("TagUTCTime", qspec.ConstUnboundInt, asn1.TagUTCTime),
		I.Const("TagUTF8String", qspec.ConstUnboundInt, asn1.TagUTF8String),
	)
}
