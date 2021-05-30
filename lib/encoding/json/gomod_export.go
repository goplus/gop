// Package json provide Go+ "encoding/json" package, as "encoding/json" package in Go.
package json

import (
	bytes "bytes"
	json "encoding/json"
	io "io"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func execCompact(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := json.Compact(args[0].(*bytes.Buffer), args[1].([]byte))
	p.Ret(2, ret0)
}

func execmDecoderUseNumber(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*json.Decoder).UseNumber()
	p.PopN(1)
}

func execmDecoderDisallowUnknownFields(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*json.Decoder).DisallowUnknownFields()
	p.PopN(1)
}

func execmDecoderDecode(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*json.Decoder).Decode(args[1])
	p.Ret(2, ret0)
}

func execmDecoderBuffered(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*json.Decoder).Buffered()
	p.Ret(1, ret0)
}

func execmDecoderToken(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*json.Decoder).Token()
	p.Ret(1, ret0, ret1)
}

func execmDecoderMore(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*json.Decoder).More()
	p.Ret(1, ret0)
}

func execmDecoderInputOffset(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*json.Decoder).InputOffset()
	p.Ret(1, ret0)
}

func execmDelimString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(json.Delim).String()
	p.Ret(1, ret0)
}

func execmEncoderEncode(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*json.Encoder).Encode(args[1])
	p.Ret(2, ret0)
}

func execmEncoderSetIndent(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*json.Encoder).SetIndent(args[1].(string), args[2].(string))
	p.PopN(3)
}

func execmEncoderSetEscapeHTML(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*json.Encoder).SetEscapeHTML(args[1].(bool))
	p.PopN(2)
}

func execHTMLEscape(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	json.HTMLEscape(args[0].(*bytes.Buffer), args[1].([]byte))
	p.PopN(2)
}

func execIndent(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := json.Indent(args[0].(*bytes.Buffer), args[1].([]byte), args[2].(string), args[3].(string))
	p.Ret(4, ret0)
}

func execmInvalidUTF8ErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*json.InvalidUTF8Error).Error()
	p.Ret(1, ret0)
}

func execmInvalidUnmarshalErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*json.InvalidUnmarshalError).Error()
	p.Ret(1, ret0)
}

func execMarshal(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := json.Marshal(args[0])
	p.Ret(1, ret0, ret1)
}

func execMarshalIndent(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := json.MarshalIndent(args[0], args[1].(string), args[2].(string))
	p.Ret(3, ret0, ret1)
}

func execiMarshalerMarshalJSON(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(json.Marshaler).MarshalJSON()
	p.Ret(1, ret0, ret1)
}

func execmMarshalerErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*json.MarshalerError).Error()
	p.Ret(1, ret0)
}

func execmMarshalerErrorUnwrap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*json.MarshalerError).Unwrap()
	p.Ret(1, ret0)
}

func toType0(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execNewDecoder(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := json.NewDecoder(toType0(args[0]))
	p.Ret(1, ret0)
}

func toType1(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execNewEncoder(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := json.NewEncoder(toType1(args[0]))
	p.Ret(1, ret0)
}

func execmNumberString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(json.Number).String()
	p.Ret(1, ret0)
}

func execmNumberFloat64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(json.Number).Float64()
	p.Ret(1, ret0, ret1)
}

func execmNumberInt64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(json.Number).Int64()
	p.Ret(1, ret0, ret1)
}

func execmRawMessageMarshalJSON(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(json.RawMessage).MarshalJSON()
	p.Ret(1, ret0, ret1)
}

func execmRawMessageUnmarshalJSON(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*json.RawMessage).UnmarshalJSON(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmSyntaxErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*json.SyntaxError).Error()
	p.Ret(1, ret0)
}

func execUnmarshal(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := json.Unmarshal(args[0].([]byte), args[1])
	p.Ret(2, ret0)
}

func execmUnmarshalFieldErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*json.UnmarshalFieldError).Error()
	p.Ret(1, ret0)
}

func execmUnmarshalTypeErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*json.UnmarshalTypeError).Error()
	p.Ret(1, ret0)
}

func execiUnmarshalerUnmarshalJSON(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(json.Unmarshaler).UnmarshalJSON(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmUnsupportedTypeErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*json.UnsupportedTypeError).Error()
	p.Ret(1, ret0)
}

func execmUnsupportedValueErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*json.UnsupportedValueError).Error()
	p.Ret(1, ret0)
}

func execValid(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := json.Valid(args[0].([]byte))
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("encoding/json")

func init() {
	I.RegisterFuncs(
		I.Func("Compact", json.Compact, execCompact),
		I.Func("(*Decoder).UseNumber", (*json.Decoder).UseNumber, execmDecoderUseNumber),
		I.Func("(*Decoder).DisallowUnknownFields", (*json.Decoder).DisallowUnknownFields, execmDecoderDisallowUnknownFields),
		I.Func("(*Decoder).Decode", (*json.Decoder).Decode, execmDecoderDecode),
		I.Func("(*Decoder).Buffered", (*json.Decoder).Buffered, execmDecoderBuffered),
		I.Func("(*Decoder).Token", (*json.Decoder).Token, execmDecoderToken),
		I.Func("(*Decoder).More", (*json.Decoder).More, execmDecoderMore),
		I.Func("(*Decoder).InputOffset", (*json.Decoder).InputOffset, execmDecoderInputOffset),
		I.Func("(Delim).String", (json.Delim).String, execmDelimString),
		I.Func("(*Encoder).Encode", (*json.Encoder).Encode, execmEncoderEncode),
		I.Func("(*Encoder).SetIndent", (*json.Encoder).SetIndent, execmEncoderSetIndent),
		I.Func("(*Encoder).SetEscapeHTML", (*json.Encoder).SetEscapeHTML, execmEncoderSetEscapeHTML),
		I.Func("HTMLEscape", json.HTMLEscape, execHTMLEscape),
		I.Func("Indent", json.Indent, execIndent),
		I.Func("(*InvalidUTF8Error).Error", (*json.InvalidUTF8Error).Error, execmInvalidUTF8ErrorError),
		I.Func("(*InvalidUnmarshalError).Error", (*json.InvalidUnmarshalError).Error, execmInvalidUnmarshalErrorError),
		I.Func("Marshal", json.Marshal, execMarshal),
		I.Func("MarshalIndent", json.MarshalIndent, execMarshalIndent),
		I.Func("(Marshaler).MarshalJSON", (json.Marshaler).MarshalJSON, execiMarshalerMarshalJSON),
		I.Func("(*MarshalerError).Error", (*json.MarshalerError).Error, execmMarshalerErrorError),
		I.Func("(*MarshalerError).Unwrap", (*json.MarshalerError).Unwrap, execmMarshalerErrorUnwrap),
		I.Func("NewDecoder", json.NewDecoder, execNewDecoder),
		I.Func("NewEncoder", json.NewEncoder, execNewEncoder),
		I.Func("(Number).String", (json.Number).String, execmNumberString),
		I.Func("(Number).Float64", (json.Number).Float64, execmNumberFloat64),
		I.Func("(Number).Int64", (json.Number).Int64, execmNumberInt64),
		I.Func("(RawMessage).MarshalJSON", (json.RawMessage).MarshalJSON, execmRawMessageMarshalJSON),
		I.Func("(*RawMessage).UnmarshalJSON", (*json.RawMessage).UnmarshalJSON, execmRawMessageUnmarshalJSON),
		I.Func("(*SyntaxError).Error", (*json.SyntaxError).Error, execmSyntaxErrorError),
		I.Func("Unmarshal", json.Unmarshal, execUnmarshal),
		I.Func("(*UnmarshalFieldError).Error", (*json.UnmarshalFieldError).Error, execmUnmarshalFieldErrorError),
		I.Func("(*UnmarshalTypeError).Error", (*json.UnmarshalTypeError).Error, execmUnmarshalTypeErrorError),
		I.Func("(Unmarshaler).UnmarshalJSON", (json.Unmarshaler).UnmarshalJSON, execiUnmarshalerUnmarshalJSON),
		I.Func("(*UnsupportedTypeError).Error", (*json.UnsupportedTypeError).Error, execmUnsupportedTypeErrorError),
		I.Func("(*UnsupportedValueError).Error", (*json.UnsupportedValueError).Error, execmUnsupportedValueErrorError),
		I.Func("Valid", json.Valid, execValid),
	)
	I.RegisterTypes(
		I.Type("Decoder", reflect.TypeOf((*json.Decoder)(nil)).Elem()),
		I.Type("Delim", reflect.TypeOf((*json.Delim)(nil)).Elem()),
		I.Type("Encoder", reflect.TypeOf((*json.Encoder)(nil)).Elem()),
		I.Type("InvalidUTF8Error", reflect.TypeOf((*json.InvalidUTF8Error)(nil)).Elem()),
		I.Type("InvalidUnmarshalError", reflect.TypeOf((*json.InvalidUnmarshalError)(nil)).Elem()),
		I.Type("Marshaler", reflect.TypeOf((*json.Marshaler)(nil)).Elem()),
		I.Type("MarshalerError", reflect.TypeOf((*json.MarshalerError)(nil)).Elem()),
		I.Type("Number", reflect.TypeOf((*json.Number)(nil)).Elem()),
		I.Type("RawMessage", reflect.TypeOf((*json.RawMessage)(nil)).Elem()),
		I.Type("SyntaxError", reflect.TypeOf((*json.SyntaxError)(nil)).Elem()),
		I.Type("Token", reflect.TypeOf((*json.Token)(nil)).Elem()),
		I.Type("UnmarshalFieldError", reflect.TypeOf((*json.UnmarshalFieldError)(nil)).Elem()),
		I.Type("UnmarshalTypeError", reflect.TypeOf((*json.UnmarshalTypeError)(nil)).Elem()),
		I.Type("Unmarshaler", reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()),
		I.Type("UnsupportedTypeError", reflect.TypeOf((*json.UnsupportedTypeError)(nil)).Elem()),
		I.Type("UnsupportedValueError", reflect.TypeOf((*json.UnsupportedValueError)(nil)).Elem()),
	)
}
