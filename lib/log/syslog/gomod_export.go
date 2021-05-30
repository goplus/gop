// Package syslog provide Go+ "log/syslog" package, as "log/syslog" package in Go.
package syslog

import (
	syslog "log/syslog"
	reflect "reflect"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execDial(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := syslog.Dial(args[0].(string), args[1].(string), args[2].(syslog.Priority), args[3].(string))
	p.Ret(4, ret0, ret1)
}

func execNew(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := syslog.New(args[0].(syslog.Priority), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execNewLogger(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := syslog.NewLogger(args[0].(syslog.Priority), args[1].(int))
	p.Ret(2, ret0, ret1)
}

func execmWriterWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*syslog.Writer).Write(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmWriterClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*syslog.Writer).Close()
	p.Ret(1, ret0)
}

func execmWriterEmerg(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*syslog.Writer).Emerg(args[1].(string))
	p.Ret(2, ret0)
}

func execmWriterAlert(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*syslog.Writer).Alert(args[1].(string))
	p.Ret(2, ret0)
}

func execmWriterCrit(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*syslog.Writer).Crit(args[1].(string))
	p.Ret(2, ret0)
}

func execmWriterErr(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*syslog.Writer).Err(args[1].(string))
	p.Ret(2, ret0)
}

func execmWriterWarning(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*syslog.Writer).Warning(args[1].(string))
	p.Ret(2, ret0)
}

func execmWriterNotice(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*syslog.Writer).Notice(args[1].(string))
	p.Ret(2, ret0)
}

func execmWriterInfo(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*syslog.Writer).Info(args[1].(string))
	p.Ret(2, ret0)
}

func execmWriterDebug(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*syslog.Writer).Debug(args[1].(string))
	p.Ret(2, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("log/syslog")

func init() {
	I.RegisterFuncs(
		I.Func("Dial", syslog.Dial, execDial),
		I.Func("New", syslog.New, execNew),
		I.Func("NewLogger", syslog.NewLogger, execNewLogger),
		I.Func("(*Writer).Write", (*syslog.Writer).Write, execmWriterWrite),
		I.Func("(*Writer).Close", (*syslog.Writer).Close, execmWriterClose),
		I.Func("(*Writer).Emerg", (*syslog.Writer).Emerg, execmWriterEmerg),
		I.Func("(*Writer).Alert", (*syslog.Writer).Alert, execmWriterAlert),
		I.Func("(*Writer).Crit", (*syslog.Writer).Crit, execmWriterCrit),
		I.Func("(*Writer).Err", (*syslog.Writer).Err, execmWriterErr),
		I.Func("(*Writer).Warning", (*syslog.Writer).Warning, execmWriterWarning),
		I.Func("(*Writer).Notice", (*syslog.Writer).Notice, execmWriterNotice),
		I.Func("(*Writer).Info", (*syslog.Writer).Info, execmWriterInfo),
		I.Func("(*Writer).Debug", (*syslog.Writer).Debug, execmWriterDebug),
	)
	I.RegisterTypes(
		I.Type("Priority", reflect.TypeOf((*syslog.Priority)(nil)).Elem()),
		I.Type("Writer", reflect.TypeOf((*syslog.Writer)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("LOG_ALERT", qspec.Int, syslog.LOG_ALERT),
		I.Const("LOG_AUTH", qspec.Int, syslog.LOG_AUTH),
		I.Const("LOG_AUTHPRIV", qspec.Int, syslog.LOG_AUTHPRIV),
		I.Const("LOG_CRIT", qspec.Int, syslog.LOG_CRIT),
		I.Const("LOG_CRON", qspec.Int, syslog.LOG_CRON),
		I.Const("LOG_DAEMON", qspec.Int, syslog.LOG_DAEMON),
		I.Const("LOG_DEBUG", qspec.Int, syslog.LOG_DEBUG),
		I.Const("LOG_EMERG", qspec.Int, syslog.LOG_EMERG),
		I.Const("LOG_ERR", qspec.Int, syslog.LOG_ERR),
		I.Const("LOG_FTP", qspec.Int, syslog.LOG_FTP),
		I.Const("LOG_INFO", qspec.Int, syslog.LOG_INFO),
		I.Const("LOG_KERN", qspec.Int, syslog.LOG_KERN),
		I.Const("LOG_LOCAL0", qspec.Int, syslog.LOG_LOCAL0),
		I.Const("LOG_LOCAL1", qspec.Int, syslog.LOG_LOCAL1),
		I.Const("LOG_LOCAL2", qspec.Int, syslog.LOG_LOCAL2),
		I.Const("LOG_LOCAL3", qspec.Int, syslog.LOG_LOCAL3),
		I.Const("LOG_LOCAL4", qspec.Int, syslog.LOG_LOCAL4),
		I.Const("LOG_LOCAL5", qspec.Int, syslog.LOG_LOCAL5),
		I.Const("LOG_LOCAL6", qspec.Int, syslog.LOG_LOCAL6),
		I.Const("LOG_LOCAL7", qspec.Int, syslog.LOG_LOCAL7),
		I.Const("LOG_LPR", qspec.Int, syslog.LOG_LPR),
		I.Const("LOG_MAIL", qspec.Int, syslog.LOG_MAIL),
		I.Const("LOG_NEWS", qspec.Int, syslog.LOG_NEWS),
		I.Const("LOG_NOTICE", qspec.Int, syslog.LOG_NOTICE),
		I.Const("LOG_SYSLOG", qspec.Int, syslog.LOG_SYSLOG),
		I.Const("LOG_USER", qspec.Int, syslog.LOG_USER),
		I.Const("LOG_UUCP", qspec.Int, syslog.LOG_UUCP),
		I.Const("LOG_WARNING", qspec.Int, syslog.LOG_WARNING),
	)
}
