The Q Language Export Tool
========

Export go package to qlang module


The Q Language : https://github.com/qiniu/qlang


Usages:
```
qexport [-contexts=""] [-defctx=false] [-outpath="./qlang"] packages
	 
The packages for go package list or std for golang all standard packages.

  -contexts string
    	optional comma-separated list of <goos>-<goarch>[-cgo] to override default contexts.
  -defctx
    	optional use default context for build, default use all contexts.
  -outpath string
    	optional set export root path (default "./qlang")
```

Examples:

```
export sync
> qexport sync

export html and html/template package
> qexport html html/template

export all package
> qexport std
```


Export pacakge runtime:

``` go
package runtime

import (
	"runtime"
)

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"_name": "runtime",

	"compiler":       runtime.Compiler,
	"GOARCH":         runtime.GOARCH,
	"GOOS":           runtime.GOOS,
	"memProfileRate": runtime.MemProfileRate,

	"blockProfile":        runtime.BlockProfile,
	"breakpoint":          runtime.Breakpoint,
	"CPUProfile":          runtime.CPUProfile,
	"caller":              runtime.Caller,
	"callers":             runtime.Callers,
	"GC":                  runtime.GC,
	"GOMAXPROCS":          runtime.GOMAXPROCS,
	"GOROOT":              runtime.GOROOT,
	"goexit":              runtime.Goexit,
	"goroutineProfile":    runtime.GoroutineProfile,
	"gosched":             runtime.Gosched,
	"lockOSThread":        runtime.LockOSThread,
	"memProfile":          runtime.MemProfile,
	"numCPU":              runtime.NumCPU,
	"numCgoCall":          runtime.NumCgoCall,
	"numGoroutine":        runtime.NumGoroutine,
	"readMemStats":        runtime.ReadMemStats,
	"readTrace":           runtime.ReadTrace,
	"setBlockProfileRate": runtime.SetBlockProfileRate,
	"setCPUProfileRate":   runtime.SetCPUProfileRate,
	"setFinalizer":        runtime.SetFinalizer,
	"stack":               runtime.Stack,
	"startTrace":          runtime.StartTrace,
	"stopTrace":           runtime.StopTrace,
	"threadCreateProfile": runtime.ThreadCreateProfile,
	"unlockOSThread":      runtime.UnlockOSThread,
	"version":             runtime.Version,

	"blockProfileRecord":      varBlockProfileRecord,
	"blockProfileRecordArray": newBlockProfileRecordArray,
	"funcForPC":               runtime.FuncForPC,
	"memProfileRecord":        varMemProfileRecord,
	"memProfileRecordArray":   newMemProfileRecordArray,
	"memStats":                varMemStats,
	"memStatsArray":           newMemStatsArray,
	"stackRecord":             varStackRecord,
	"stackRecordArray":        newStackRecordArray,
	"typeAssertionError":      varTypeAssertionError,
	"typeAssertionErrorArray": newTypeAssertionErrorArray,
}

func varBlockProfileRecord() runtime.BlockProfileRecord {
	var v runtime.BlockProfileRecord
	return v
}

func newBlockProfileRecordArray(n int) []runtime.BlockProfileRecord {
	return make([]runtime.BlockProfileRecord, n)
}

func varMemProfileRecord() runtime.MemProfileRecord {
	var v runtime.MemProfileRecord
	return v
}

func newMemProfileRecordArray(n int) []runtime.MemProfileRecord {
	return make([]runtime.MemProfileRecord, n)
}

func varMemStats() runtime.MemStats {
	var v runtime.MemStats
	return v
}

func newMemStatsArray(n int) []runtime.MemStats {
	return make([]runtime.MemStats, n)
}

func varStackRecord() runtime.StackRecord {
	var v runtime.StackRecord
	return v
}

func newStackRecordArray(n int) []runtime.StackRecord {
	return make([]runtime.StackRecord, n)
}

func varTypeAssertionError() runtime.TypeAssertionError {
	var v runtime.TypeAssertionError
	return v
}

func newTypeAssertionErrorArray(n int) []runtime.TypeAssertionError {
	return make([]runtime.TypeAssertionError, n)
}
```

		