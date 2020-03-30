package runtime

import (
	"runtime"
)

func newBlockProfileRecords(n int) []runtime.BlockProfileRecord {
	return make([]runtime.BlockProfileRecord, n)
}

func newMemProfileRecords(n int) []runtime.MemProfileRecord {
	return make([]runtime.MemProfileRecord, n)
}

func newStackRecords(n int) []runtime.StackRecord {
	return make([]runtime.StackRecord, n)
}

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"_name":          "runtime",
	"memProfileRate": runtime.MemProfileRate,

	"compiler": runtime.Compiler,
	"GOARCH":   runtime.GOARCH,
	"GOOS":     runtime.GOOS,

	"blockProfileRecords": newBlockProfileRecords,
	"memProfileRecords":   newMemProfileRecords,
	"stackRecords":        newStackRecords,
	"blockProfile":        runtime.BlockProfile,
	"memProfile":          runtime.MemProfile,
	"CPUProfile":          runtime.CPUProfile,
	"threadCreateProfile": runtime.ThreadCreateProfile,
	"goroutineProfile":    runtime.GoroutineProfile,
	"setBlockProfileRate": runtime.SetBlockProfileRate,
	"setCPUProfileRate":   runtime.SetCPUProfileRate,

	"breakpoint":     runtime.Breakpoint,
	"caller":         runtime.Caller,
	"callers":        runtime.Callers,
	"funcForPC":      runtime.FuncForPC,
	"GC":             runtime.GC,
	"GOMAXPROCS":     runtime.GOMAXPROCS,
	"GOROOT":         runtime.GOROOT,
	"goexit":         runtime.Goexit,
	"gosched":        runtime.Gosched,
	"lockOSThread":   runtime.LockOSThread,
	"numCPU":         runtime.NumCPU,
	"numCgoCall":     runtime.NumCgoCall,
	"numGoroutine":   runtime.NumGoroutine,
	"readMemStats":   runtime.ReadMemStats,
	"setFinalizer":   runtime.SetFinalizer,
	"stack":          runtime.Stack,
	"unlockOSThread": runtime.UnlockOSThread,
	"version":        runtime.Version,
}
