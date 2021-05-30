// Package runtime provide Go+ "runtime" package, as "runtime" package in Go.
package runtime

import (
	reflect "reflect"
	runtime "runtime"
	unsafe "unsafe"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execBlockProfile(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := runtime.BlockProfile(args[0].([]runtime.BlockProfileRecord))
	p.Ret(1, ret0, ret1)
}

func execBreakpoint(_ int, p *gop.Context) {
	runtime.Breakpoint()
	p.PopN(0)
}

func execCPUProfile(_ int, p *gop.Context) {
	ret0 := runtime.CPUProfile()
	p.Ret(0, ret0)
}

func execCaller(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2, ret3 := runtime.Caller(args[0].(int))
	p.Ret(1, ret0, ret1, ret2, ret3)
}

func execCallers(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := runtime.Callers(args[0].(int), args[1].([]uintptr))
	p.Ret(2, ret0)
}

func execCallersFrames(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := runtime.CallersFrames(args[0].([]uintptr))
	p.Ret(1, ret0)
}

func execiErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(runtime.Error).Error()
	p.Ret(1, ret0)
}

func execiErrorRuntimeError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(runtime.Error).RuntimeError()
	p.PopN(1)
}

func execmFramesNext(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*runtime.Frames).Next()
	p.Ret(1, ret0, ret1)
}

func execmFuncName(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*runtime.Func).Name()
	p.Ret(1, ret0)
}

func execmFuncEntry(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*runtime.Func).Entry()
	p.Ret(1, ret0)
}

func execmFuncFileLine(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*runtime.Func).FileLine(args[1].(uintptr))
	p.Ret(2, ret0, ret1)
}

func execFuncForPC(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := runtime.FuncForPC(args[0].(uintptr))
	p.Ret(1, ret0)
}

func execGC(_ int, p *gop.Context) {
	runtime.GC()
	p.PopN(0)
}

func execGOMAXPROCS(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := runtime.GOMAXPROCS(args[0].(int))
	p.Ret(1, ret0)
}

func execGOROOT(_ int, p *gop.Context) {
	ret0 := runtime.GOROOT()
	p.Ret(0, ret0)
}

func execGoexit(_ int, p *gop.Context) {
	runtime.Goexit()
	p.PopN(0)
}

func execGoroutineProfile(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := runtime.GoroutineProfile(args[0].([]runtime.StackRecord))
	p.Ret(1, ret0, ret1)
}

func execGosched(_ int, p *gop.Context) {
	runtime.Gosched()
	p.PopN(0)
}

func execKeepAlive(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	runtime.KeepAlive(args[0])
	p.PopN(1)
}

func execLockOSThread(_ int, p *gop.Context) {
	runtime.LockOSThread()
	p.PopN(0)
}

func execMemProfile(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := runtime.MemProfile(args[0].([]runtime.MemProfileRecord), args[1].(bool))
	p.Ret(2, ret0, ret1)
}

func execmMemProfileRecordInUseBytes(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*runtime.MemProfileRecord).InUseBytes()
	p.Ret(1, ret0)
}

func execmMemProfileRecordInUseObjects(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*runtime.MemProfileRecord).InUseObjects()
	p.Ret(1, ret0)
}

func execmMemProfileRecordStack(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*runtime.MemProfileRecord).Stack()
	p.Ret(1, ret0)
}

func execMutexProfile(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := runtime.MutexProfile(args[0].([]runtime.BlockProfileRecord))
	p.Ret(1, ret0, ret1)
}

func execNumCPU(_ int, p *gop.Context) {
	ret0 := runtime.NumCPU()
	p.Ret(0, ret0)
}

func execNumCgoCall(_ int, p *gop.Context) {
	ret0 := runtime.NumCgoCall()
	p.Ret(0, ret0)
}

func execNumGoroutine(_ int, p *gop.Context) {
	ret0 := runtime.NumGoroutine()
	p.Ret(0, ret0)
}

func execReadMemStats(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	runtime.ReadMemStats(args[0].(*runtime.MemStats))
	p.PopN(1)
}

func execReadTrace(_ int, p *gop.Context) {
	ret0 := runtime.ReadTrace()
	p.Ret(0, ret0)
}

func execSetBlockProfileRate(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	runtime.SetBlockProfileRate(args[0].(int))
	p.PopN(1)
}

func execSetCPUProfileRate(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	runtime.SetCPUProfileRate(args[0].(int))
	p.PopN(1)
}

func execSetCgoTraceback(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	runtime.SetCgoTraceback(args[0].(int), args[1].(unsafe.Pointer), args[2].(unsafe.Pointer), args[3].(unsafe.Pointer))
	p.PopN(4)
}

func execSetFinalizer(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	runtime.SetFinalizer(args[0], args[1])
	p.PopN(2)
}

func execSetMutexProfileFraction(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := runtime.SetMutexProfileFraction(args[0].(int))
	p.Ret(1, ret0)
}

func execStack(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := runtime.Stack(args[0].([]byte), args[1].(bool))
	p.Ret(2, ret0)
}

func execmStackRecordStack(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*runtime.StackRecord).Stack()
	p.Ret(1, ret0)
}

func execStartTrace(_ int, p *gop.Context) {
	ret0 := runtime.StartTrace()
	p.Ret(0, ret0)
}

func execStopTrace(_ int, p *gop.Context) {
	runtime.StopTrace()
	p.PopN(0)
}

func execThreadCreateProfile(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := runtime.ThreadCreateProfile(args[0].([]runtime.StackRecord))
	p.Ret(1, ret0, ret1)
}

func execmTypeAssertionErrorRuntimeError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*runtime.TypeAssertionError).RuntimeError()
	p.PopN(1)
}

func execmTypeAssertionErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*runtime.TypeAssertionError).Error()
	p.Ret(1, ret0)
}

func execUnlockOSThread(_ int, p *gop.Context) {
	runtime.UnlockOSThread()
	p.PopN(0)
}

func execVersion(_ int, p *gop.Context) {
	ret0 := runtime.Version()
	p.Ret(0, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("runtime")

func init() {
	I.RegisterFuncs(
		I.Func("BlockProfile", runtime.BlockProfile, execBlockProfile),
		I.Func("Breakpoint", runtime.Breakpoint, execBreakpoint),
		I.Func("CPUProfile", runtime.CPUProfile, execCPUProfile),
		I.Func("Caller", runtime.Caller, execCaller),
		I.Func("Callers", runtime.Callers, execCallers),
		I.Func("CallersFrames", runtime.CallersFrames, execCallersFrames),
		I.Func("(Error).Error", (runtime.Error).Error, execiErrorError),
		I.Func("(Error).RuntimeError", (runtime.Error).RuntimeError, execiErrorRuntimeError),
		I.Func("(*Frames).Next", (*runtime.Frames).Next, execmFramesNext),
		I.Func("(*Func).Name", (*runtime.Func).Name, execmFuncName),
		I.Func("(*Func).Entry", (*runtime.Func).Entry, execmFuncEntry),
		I.Func("(*Func).FileLine", (*runtime.Func).FileLine, execmFuncFileLine),
		I.Func("FuncForPC", runtime.FuncForPC, execFuncForPC),
		I.Func("GC", runtime.GC, execGC),
		I.Func("GOMAXPROCS", runtime.GOMAXPROCS, execGOMAXPROCS),
		I.Func("GOROOT", runtime.GOROOT, execGOROOT),
		I.Func("Goexit", runtime.Goexit, execGoexit),
		I.Func("GoroutineProfile", runtime.GoroutineProfile, execGoroutineProfile),
		I.Func("Gosched", runtime.Gosched, execGosched),
		I.Func("KeepAlive", runtime.KeepAlive, execKeepAlive),
		I.Func("LockOSThread", runtime.LockOSThread, execLockOSThread),
		I.Func("MemProfile", runtime.MemProfile, execMemProfile),
		I.Func("(*MemProfileRecord).InUseBytes", (*runtime.MemProfileRecord).InUseBytes, execmMemProfileRecordInUseBytes),
		I.Func("(*MemProfileRecord).InUseObjects", (*runtime.MemProfileRecord).InUseObjects, execmMemProfileRecordInUseObjects),
		I.Func("(*MemProfileRecord).Stack", (*runtime.MemProfileRecord).Stack, execmMemProfileRecordStack),
		I.Func("MutexProfile", runtime.MutexProfile, execMutexProfile),
		I.Func("NumCPU", runtime.NumCPU, execNumCPU),
		I.Func("NumCgoCall", runtime.NumCgoCall, execNumCgoCall),
		I.Func("NumGoroutine", runtime.NumGoroutine, execNumGoroutine),
		I.Func("ReadMemStats", runtime.ReadMemStats, execReadMemStats),
		I.Func("ReadTrace", runtime.ReadTrace, execReadTrace),
		I.Func("SetBlockProfileRate", runtime.SetBlockProfileRate, execSetBlockProfileRate),
		I.Func("SetCPUProfileRate", runtime.SetCPUProfileRate, execSetCPUProfileRate),
		I.Func("SetCgoTraceback", runtime.SetCgoTraceback, execSetCgoTraceback),
		I.Func("SetFinalizer", runtime.SetFinalizer, execSetFinalizer),
		I.Func("SetMutexProfileFraction", runtime.SetMutexProfileFraction, execSetMutexProfileFraction),
		I.Func("Stack", runtime.Stack, execStack),
		I.Func("(*StackRecord).Stack", (*runtime.StackRecord).Stack, execmStackRecordStack),
		I.Func("StartTrace", runtime.StartTrace, execStartTrace),
		I.Func("StopTrace", runtime.StopTrace, execStopTrace),
		I.Func("ThreadCreateProfile", runtime.ThreadCreateProfile, execThreadCreateProfile),
		I.Func("(*TypeAssertionError).RuntimeError", (*runtime.TypeAssertionError).RuntimeError, execmTypeAssertionErrorRuntimeError),
		I.Func("(*TypeAssertionError).Error", (*runtime.TypeAssertionError).Error, execmTypeAssertionErrorError),
		I.Func("UnlockOSThread", runtime.UnlockOSThread, execUnlockOSThread),
		I.Func("Version", runtime.Version, execVersion),
	)
	I.RegisterVars(
		I.Var("MemProfileRate", &runtime.MemProfileRate),
	)
	I.RegisterTypes(
		I.Type("BlockProfileRecord", reflect.TypeOf((*runtime.BlockProfileRecord)(nil)).Elem()),
		I.Type("Error", reflect.TypeOf((*runtime.Error)(nil)).Elem()),
		I.Type("Frame", reflect.TypeOf((*runtime.Frame)(nil)).Elem()),
		I.Type("Frames", reflect.TypeOf((*runtime.Frames)(nil)).Elem()),
		I.Type("Func", reflect.TypeOf((*runtime.Func)(nil)).Elem()),
		I.Type("MemProfileRecord", reflect.TypeOf((*runtime.MemProfileRecord)(nil)).Elem()),
		I.Type("MemStats", reflect.TypeOf((*runtime.MemStats)(nil)).Elem()),
		I.Type("StackRecord", reflect.TypeOf((*runtime.StackRecord)(nil)).Elem()),
		I.Type("TypeAssertionError", reflect.TypeOf((*runtime.TypeAssertionError)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("Compiler", qspec.ConstBoundString, runtime.Compiler),
		I.Const("GOARCH", qspec.String, runtime.GOARCH),
		I.Const("GOOS", qspec.String, runtime.GOOS),
	)
}
