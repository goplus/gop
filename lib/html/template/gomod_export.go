// Package template provide Go+ "html/template" package, as "html/template" package in Go.
package template

import (
	template "html/template"
	io "io"
	reflect "reflect"
	parse "text/template/parse"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execmErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*template.Error).Error()
	p.Ret(1, ret0)
}

func toType0(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execHTMLEscape(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	template.HTMLEscape(toType0(args[0]), args[1].([]byte))
	p.PopN(2)
}

func execHTMLEscapeString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := template.HTMLEscapeString(args[0].(string))
	p.Ret(1, ret0)
}

func execHTMLEscaper(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := template.HTMLEscaper(args...)
	p.Ret(arity, ret0)
}

func execIsTrue(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := template.IsTrue(args[0])
	p.Ret(1, ret0, ret1)
}

func execJSEscape(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	template.JSEscape(toType0(args[0]), args[1].([]byte))
	p.PopN(2)
}

func execJSEscapeString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := template.JSEscapeString(args[0].(string))
	p.Ret(1, ret0)
}

func execJSEscaper(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := template.JSEscaper(args...)
	p.Ret(arity, ret0)
}

func toType1(v interface{}) error {
	if v == nil {
		return nil
	}
	return v.(error)
}

func execMust(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := template.Must(args[0].(*template.Template), toType1(args[1]))
	p.Ret(2, ret0)
}

func execNew(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := template.New(args[0].(string))
	p.Ret(1, ret0)
}

func execParseFiles(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := template.ParseFiles(gop.ToStrings(args)...)
	p.Ret(arity, ret0, ret1)
}

func execParseGlob(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := template.ParseGlob(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execmTemplateTemplates(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*template.Template).Templates()
	p.Ret(1, ret0)
}

func execmTemplateOption(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*template.Template).Option(gop.ToStrings(args[1:])...)
	p.Ret(arity, ret0)
}

func execmTemplateExecute(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*template.Template).Execute(toType0(args[1]), args[2])
	p.Ret(3, ret0)
}

func execmTemplateExecuteTemplate(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := args[0].(*template.Template).ExecuteTemplate(toType0(args[1]), args[2].(string), args[3])
	p.Ret(4, ret0)
}

func execmTemplateDefinedTemplates(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*template.Template).DefinedTemplates()
	p.Ret(1, ret0)
}

func execmTemplateParse(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*template.Template).Parse(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmTemplateAddParseTree(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*template.Template).AddParseTree(args[1].(string), args[2].(*parse.Tree))
	p.Ret(3, ret0, ret1)
}

func execmTemplateClone(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*template.Template).Clone()
	p.Ret(1, ret0, ret1)
}

func execmTemplateNew(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*template.Template).New(args[1].(string))
	p.Ret(2, ret0)
}

func execmTemplateName(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*template.Template).Name()
	p.Ret(1, ret0)
}

func execmTemplateFuncs(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*template.Template).Funcs(args[1].(template.FuncMap))
	p.Ret(2, ret0)
}

func execmTemplateDelims(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*template.Template).Delims(args[1].(string), args[2].(string))
	p.Ret(3, ret0)
}

func execmTemplateLookup(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*template.Template).Lookup(args[1].(string))
	p.Ret(2, ret0)
}

func execmTemplateParseFiles(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := args[0].(*template.Template).ParseFiles(gop.ToStrings(args[1:])...)
	p.Ret(arity, ret0, ret1)
}

func execmTemplateParseGlob(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*template.Template).ParseGlob(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execURLQueryEscaper(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := template.URLQueryEscaper(args...)
	p.Ret(arity, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("html/template")

func init() {
	I.RegisterFuncs(
		I.Func("(*Error).Error", (*template.Error).Error, execmErrorError),
		I.Func("HTMLEscape", template.HTMLEscape, execHTMLEscape),
		I.Func("HTMLEscapeString", template.HTMLEscapeString, execHTMLEscapeString),
		I.Func("IsTrue", template.IsTrue, execIsTrue),
		I.Func("JSEscape", template.JSEscape, execJSEscape),
		I.Func("JSEscapeString", template.JSEscapeString, execJSEscapeString),
		I.Func("Must", template.Must, execMust),
		I.Func("New", template.New, execNew),
		I.Func("ParseGlob", template.ParseGlob, execParseGlob),
		I.Func("(*Template).Templates", (*template.Template).Templates, execmTemplateTemplates),
		I.Func("(*Template).Execute", (*template.Template).Execute, execmTemplateExecute),
		I.Func("(*Template).ExecuteTemplate", (*template.Template).ExecuteTemplate, execmTemplateExecuteTemplate),
		I.Func("(*Template).DefinedTemplates", (*template.Template).DefinedTemplates, execmTemplateDefinedTemplates),
		I.Func("(*Template).Parse", (*template.Template).Parse, execmTemplateParse),
		I.Func("(*Template).AddParseTree", (*template.Template).AddParseTree, execmTemplateAddParseTree),
		I.Func("(*Template).Clone", (*template.Template).Clone, execmTemplateClone),
		I.Func("(*Template).New", (*template.Template).New, execmTemplateNew),
		I.Func("(*Template).Name", (*template.Template).Name, execmTemplateName),
		I.Func("(*Template).Funcs", (*template.Template).Funcs, execmTemplateFuncs),
		I.Func("(*Template).Delims", (*template.Template).Delims, execmTemplateDelims),
		I.Func("(*Template).Lookup", (*template.Template).Lookup, execmTemplateLookup),
		I.Func("(*Template).ParseGlob", (*template.Template).ParseGlob, execmTemplateParseGlob),
	)
	I.RegisterFuncvs(
		I.Funcv("HTMLEscaper", template.HTMLEscaper, execHTMLEscaper),
		I.Funcv("JSEscaper", template.JSEscaper, execJSEscaper),
		I.Funcv("ParseFiles", template.ParseFiles, execParseFiles),
		I.Funcv("(*Template).Option", (*template.Template).Option, execmTemplateOption),
		I.Funcv("(*Template).ParseFiles", (*template.Template).ParseFiles, execmTemplateParseFiles),
		I.Funcv("URLQueryEscaper", template.URLQueryEscaper, execURLQueryEscaper),
	)
	I.RegisterTypes(
		I.Type("CSS", reflect.TypeOf((*template.CSS)(nil)).Elem()),
		I.Type("Error", reflect.TypeOf((*template.Error)(nil)).Elem()),
		I.Type("ErrorCode", reflect.TypeOf((*template.ErrorCode)(nil)).Elem()),
		I.Type("FuncMap", reflect.TypeOf((*template.FuncMap)(nil)).Elem()),
		I.Type("HTML", reflect.TypeOf((*template.HTML)(nil)).Elem()),
		I.Type("HTMLAttr", reflect.TypeOf((*template.HTMLAttr)(nil)).Elem()),
		I.Type("JS", reflect.TypeOf((*template.JS)(nil)).Elem()),
		I.Type("JSStr", reflect.TypeOf((*template.JSStr)(nil)).Elem()),
		I.Type("Srcset", reflect.TypeOf((*template.Srcset)(nil)).Elem()),
		I.Type("Template", reflect.TypeOf((*template.Template)(nil)).Elem()),
		I.Type("URL", reflect.TypeOf((*template.URL)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("ErrAmbigContext", qspec.Int, template.ErrAmbigContext),
		I.Const("ErrBadHTML", qspec.Int, template.ErrBadHTML),
		I.Const("ErrBranchEnd", qspec.Int, template.ErrBranchEnd),
		I.Const("ErrEndContext", qspec.Int, template.ErrEndContext),
		I.Const("ErrNoSuchTemplate", qspec.Int, template.ErrNoSuchTemplate),
		I.Const("ErrOutputContext", qspec.Int, template.ErrOutputContext),
		I.Const("ErrPartialCharset", qspec.Int, template.ErrPartialCharset),
		I.Const("ErrPartialEscape", qspec.Int, template.ErrPartialEscape),
		I.Const("ErrPredefinedEscaper", qspec.Int, template.ErrPredefinedEscaper),
		I.Const("ErrRangeLoopReentry", qspec.Int, template.ErrRangeLoopReentry),
		I.Const("ErrSlashAmbig", qspec.Int, template.ErrSlashAmbig),
		I.Const("OK", qspec.Int, template.OK),
	)
}
