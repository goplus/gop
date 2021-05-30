// Package xml provide Go+ "encoding/xml" package, as "encoding/xml" package in Go.
package xml

import (
	xml "encoding/xml"
	io "io"
	reflect "reflect"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execmCharDataCopy(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(xml.CharData).Copy()
	p.Ret(1, ret0)
}

func execmCommentCopy(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(xml.Comment).Copy()
	p.Ret(1, ret0)
}

func execCopyToken(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := xml.CopyToken(args[0])
	p.Ret(1, ret0)
}

func execmDecoderDecode(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*xml.Decoder).Decode(args[1])
	p.Ret(2, ret0)
}

func execmDecoderDecodeElement(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*xml.Decoder).DecodeElement(args[1], args[2].(*xml.StartElement))
	p.Ret(3, ret0)
}

func execmDecoderSkip(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*xml.Decoder).Skip()
	p.Ret(1, ret0)
}

func execmDecoderToken(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*xml.Decoder).Token()
	p.Ret(1, ret0, ret1)
}

func execmDecoderRawToken(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*xml.Decoder).RawToken()
	p.Ret(1, ret0, ret1)
}

func execmDecoderInputOffset(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*xml.Decoder).InputOffset()
	p.Ret(1, ret0)
}

func execmDirectiveCopy(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(xml.Directive).Copy()
	p.Ret(1, ret0)
}

func execmEncoderIndent(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*xml.Encoder).Indent(args[1].(string), args[2].(string))
	p.PopN(3)
}

func execmEncoderEncode(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*xml.Encoder).Encode(args[1])
	p.Ret(2, ret0)
}

func execmEncoderEncodeElement(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*xml.Encoder).EncodeElement(args[1], args[2].(xml.StartElement))
	p.Ret(3, ret0)
}

func execmEncoderEncodeToken(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*xml.Encoder).EncodeToken(args[1])
	p.Ret(2, ret0)
}

func execmEncoderFlush(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*xml.Encoder).Flush()
	p.Ret(1, ret0)
}

func toType0(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execEscape(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	xml.Escape(toType0(args[0]), args[1].([]byte))
	p.PopN(2)
}

func execEscapeText(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := xml.EscapeText(toType0(args[0]), args[1].([]byte))
	p.Ret(2, ret0)
}

func execMarshal(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := xml.Marshal(args[0])
	p.Ret(1, ret0, ret1)
}

func execMarshalIndent(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := xml.MarshalIndent(args[0], args[1].(string), args[2].(string))
	p.Ret(3, ret0, ret1)
}

func execiMarshalerMarshalXML(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(xml.Marshaler).MarshalXML(args[1].(*xml.Encoder), args[2].(xml.StartElement))
	p.Ret(3, ret0)
}

func execiMarshalerAttrMarshalXMLAttr(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(xml.MarshalerAttr).MarshalXMLAttr(args[1].(xml.Name))
	p.Ret(2, ret0, ret1)
}

func toType1(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execNewDecoder(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := xml.NewDecoder(toType1(args[0]))
	p.Ret(1, ret0)
}

func execNewEncoder(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := xml.NewEncoder(toType0(args[0]))
	p.Ret(1, ret0)
}

func toType2(v interface{}) xml.TokenReader {
	if v == nil {
		return nil
	}
	return v.(xml.TokenReader)
}

func execNewTokenDecoder(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := xml.NewTokenDecoder(toType2(args[0]))
	p.Ret(1, ret0)
}

func execmProcInstCopy(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(xml.ProcInst).Copy()
	p.Ret(1, ret0)
}

func execmStartElementCopy(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(xml.StartElement).Copy()
	p.Ret(1, ret0)
}

func execmStartElementEnd(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(xml.StartElement).End()
	p.Ret(1, ret0)
}

func execmSyntaxErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*xml.SyntaxError).Error()
	p.Ret(1, ret0)
}

func execmTagPathErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*xml.TagPathError).Error()
	p.Ret(1, ret0)
}

func execiTokenReaderToken(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(xml.TokenReader).Token()
	p.Ret(1, ret0, ret1)
}

func execUnmarshal(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := xml.Unmarshal(args[0].([]byte), args[1])
	p.Ret(2, ret0)
}

func execmUnmarshalErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(xml.UnmarshalError).Error()
	p.Ret(1, ret0)
}

func execiUnmarshalerUnmarshalXML(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(xml.Unmarshaler).UnmarshalXML(args[1].(*xml.Decoder), args[2].(xml.StartElement))
	p.Ret(3, ret0)
}

func execiUnmarshalerAttrUnmarshalXMLAttr(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(xml.UnmarshalerAttr).UnmarshalXMLAttr(args[1].(xml.Attr))
	p.Ret(2, ret0)
}

func execmUnsupportedTypeErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*xml.UnsupportedTypeError).Error()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("encoding/xml")

func init() {
	I.RegisterFuncs(
		I.Func("(CharData).Copy", (xml.CharData).Copy, execmCharDataCopy),
		I.Func("(Comment).Copy", (xml.Comment).Copy, execmCommentCopy),
		I.Func("CopyToken", xml.CopyToken, execCopyToken),
		I.Func("(*Decoder).Decode", (*xml.Decoder).Decode, execmDecoderDecode),
		I.Func("(*Decoder).DecodeElement", (*xml.Decoder).DecodeElement, execmDecoderDecodeElement),
		I.Func("(*Decoder).Skip", (*xml.Decoder).Skip, execmDecoderSkip),
		I.Func("(*Decoder).Token", (*xml.Decoder).Token, execmDecoderToken),
		I.Func("(*Decoder).RawToken", (*xml.Decoder).RawToken, execmDecoderRawToken),
		I.Func("(*Decoder).InputOffset", (*xml.Decoder).InputOffset, execmDecoderInputOffset),
		I.Func("(Directive).Copy", (xml.Directive).Copy, execmDirectiveCopy),
		I.Func("(*Encoder).Indent", (*xml.Encoder).Indent, execmEncoderIndent),
		I.Func("(*Encoder).Encode", (*xml.Encoder).Encode, execmEncoderEncode),
		I.Func("(*Encoder).EncodeElement", (*xml.Encoder).EncodeElement, execmEncoderEncodeElement),
		I.Func("(*Encoder).EncodeToken", (*xml.Encoder).EncodeToken, execmEncoderEncodeToken),
		I.Func("(*Encoder).Flush", (*xml.Encoder).Flush, execmEncoderFlush),
		I.Func("Escape", xml.Escape, execEscape),
		I.Func("EscapeText", xml.EscapeText, execEscapeText),
		I.Func("Marshal", xml.Marshal, execMarshal),
		I.Func("MarshalIndent", xml.MarshalIndent, execMarshalIndent),
		I.Func("(Marshaler).MarshalXML", (xml.Marshaler).MarshalXML, execiMarshalerMarshalXML),
		I.Func("(MarshalerAttr).MarshalXMLAttr", (xml.MarshalerAttr).MarshalXMLAttr, execiMarshalerAttrMarshalXMLAttr),
		I.Func("NewDecoder", xml.NewDecoder, execNewDecoder),
		I.Func("NewEncoder", xml.NewEncoder, execNewEncoder),
		I.Func("NewTokenDecoder", xml.NewTokenDecoder, execNewTokenDecoder),
		I.Func("(ProcInst).Copy", (xml.ProcInst).Copy, execmProcInstCopy),
		I.Func("(StartElement).Copy", (xml.StartElement).Copy, execmStartElementCopy),
		I.Func("(StartElement).End", (xml.StartElement).End, execmStartElementEnd),
		I.Func("(*SyntaxError).Error", (*xml.SyntaxError).Error, execmSyntaxErrorError),
		I.Func("(*TagPathError).Error", (*xml.TagPathError).Error, execmTagPathErrorError),
		I.Func("(TokenReader).Token", (xml.TokenReader).Token, execiTokenReaderToken),
		I.Func("Unmarshal", xml.Unmarshal, execUnmarshal),
		I.Func("(UnmarshalError).Error", (xml.UnmarshalError).Error, execmUnmarshalErrorError),
		I.Func("(Unmarshaler).UnmarshalXML", (xml.Unmarshaler).UnmarshalXML, execiUnmarshalerUnmarshalXML),
		I.Func("(UnmarshalerAttr).UnmarshalXMLAttr", (xml.UnmarshalerAttr).UnmarshalXMLAttr, execiUnmarshalerAttrUnmarshalXMLAttr),
		I.Func("(*UnsupportedTypeError).Error", (*xml.UnsupportedTypeError).Error, execmUnsupportedTypeErrorError),
	)
	I.RegisterVars(
		I.Var("HTMLAutoClose", &xml.HTMLAutoClose),
		I.Var("HTMLEntity", &xml.HTMLEntity),
	)
	I.RegisterTypes(
		I.Type("Attr", reflect.TypeOf((*xml.Attr)(nil)).Elem()),
		I.Type("CharData", reflect.TypeOf((*xml.CharData)(nil)).Elem()),
		I.Type("Comment", reflect.TypeOf((*xml.Comment)(nil)).Elem()),
		I.Type("Decoder", reflect.TypeOf((*xml.Decoder)(nil)).Elem()),
		I.Type("Directive", reflect.TypeOf((*xml.Directive)(nil)).Elem()),
		I.Type("Encoder", reflect.TypeOf((*xml.Encoder)(nil)).Elem()),
		I.Type("EndElement", reflect.TypeOf((*xml.EndElement)(nil)).Elem()),
		I.Type("Marshaler", reflect.TypeOf((*xml.Marshaler)(nil)).Elem()),
		I.Type("MarshalerAttr", reflect.TypeOf((*xml.MarshalerAttr)(nil)).Elem()),
		I.Type("Name", reflect.TypeOf((*xml.Name)(nil)).Elem()),
		I.Type("ProcInst", reflect.TypeOf((*xml.ProcInst)(nil)).Elem()),
		I.Type("StartElement", reflect.TypeOf((*xml.StartElement)(nil)).Elem()),
		I.Type("SyntaxError", reflect.TypeOf((*xml.SyntaxError)(nil)).Elem()),
		I.Type("TagPathError", reflect.TypeOf((*xml.TagPathError)(nil)).Elem()),
		I.Type("Token", reflect.TypeOf((*xml.Token)(nil)).Elem()),
		I.Type("TokenReader", reflect.TypeOf((*xml.TokenReader)(nil)).Elem()),
		I.Type("UnmarshalError", reflect.TypeOf((*xml.UnmarshalError)(nil)).Elem()),
		I.Type("Unmarshaler", reflect.TypeOf((*xml.Unmarshaler)(nil)).Elem()),
		I.Type("UnmarshalerAttr", reflect.TypeOf((*xml.UnmarshalerAttr)(nil)).Elem()),
		I.Type("UnsupportedTypeError", reflect.TypeOf((*xml.UnsupportedTypeError)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("Header", qspec.ConstBoundString, xml.Header),
	)
}
