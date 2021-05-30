// Package exec provide Go+ "os/exec" package, as "os/exec" package in Go.
package exec

import (
	context "context"
	exec "os/exec"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func execmCmdString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*exec.Cmd).String()
	p.Ret(1, ret0)
}

func execmCmdRun(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*exec.Cmd).Run()
	p.Ret(1, ret0)
}

func execmCmdStart(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*exec.Cmd).Start()
	p.Ret(1, ret0)
}

func execmCmdWait(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*exec.Cmd).Wait()
	p.Ret(1, ret0)
}

func execmCmdOutput(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*exec.Cmd).Output()
	p.Ret(1, ret0, ret1)
}

func execmCmdCombinedOutput(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*exec.Cmd).CombinedOutput()
	p.Ret(1, ret0, ret1)
}

func execmCmdStdinPipe(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*exec.Cmd).StdinPipe()
	p.Ret(1, ret0, ret1)
}

func execmCmdStdoutPipe(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*exec.Cmd).StdoutPipe()
	p.Ret(1, ret0, ret1)
}

func execmCmdStderrPipe(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*exec.Cmd).StderrPipe()
	p.Ret(1, ret0, ret1)
}

func execCommand(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := exec.Command(args[0].(string), gop.ToStrings(args[1:])...)
	p.Ret(arity, ret0)
}

func toType0(v interface{}) context.Context {
	if v == nil {
		return nil
	}
	return v.(context.Context)
}

func execCommandContext(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := exec.CommandContext(toType0(args[0]), args[1].(string), gop.ToStrings(args[2:])...)
	p.Ret(arity, ret0)
}

func execmErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*exec.Error).Error()
	p.Ret(1, ret0)
}

func execmErrorUnwrap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*exec.Error).Unwrap()
	p.Ret(1, ret0)
}

func execmExitErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*exec.ExitError).Error()
	p.Ret(1, ret0)
}

func execLookPath(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := exec.LookPath(args[0].(string))
	p.Ret(1, ret0, ret1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("os/exec")

func init() {
	I.RegisterFuncs(
		I.Func("(*Cmd).String", (*exec.Cmd).String, execmCmdString),
		I.Func("(*Cmd).Run", (*exec.Cmd).Run, execmCmdRun),
		I.Func("(*Cmd).Start", (*exec.Cmd).Start, execmCmdStart),
		I.Func("(*Cmd).Wait", (*exec.Cmd).Wait, execmCmdWait),
		I.Func("(*Cmd).Output", (*exec.Cmd).Output, execmCmdOutput),
		I.Func("(*Cmd).CombinedOutput", (*exec.Cmd).CombinedOutput, execmCmdCombinedOutput),
		I.Func("(*Cmd).StdinPipe", (*exec.Cmd).StdinPipe, execmCmdStdinPipe),
		I.Func("(*Cmd).StdoutPipe", (*exec.Cmd).StdoutPipe, execmCmdStdoutPipe),
		I.Func("(*Cmd).StderrPipe", (*exec.Cmd).StderrPipe, execmCmdStderrPipe),
		I.Func("(*Error).Error", (*exec.Error).Error, execmErrorError),
		I.Func("(*Error).Unwrap", (*exec.Error).Unwrap, execmErrorUnwrap),
		I.Func("(*ExitError).Error", (*exec.ExitError).Error, execmExitErrorError),
		I.Func("LookPath", exec.LookPath, execLookPath),
	)
	I.RegisterFuncvs(
		I.Funcv("Command", exec.Command, execCommand),
		I.Funcv("CommandContext", exec.CommandContext, execCommandContext),
	)
	I.RegisterVars(
		I.Var("ErrNotFound", &exec.ErrNotFound),
	)
	I.RegisterTypes(
		I.Type("Cmd", reflect.TypeOf((*exec.Cmd)(nil)).Elem()),
		I.Type("Error", reflect.TypeOf((*exec.Error)(nil)).Elem()),
		I.Type("ExitError", reflect.TypeOf((*exec.ExitError)(nil)).Elem()),
	)
}
