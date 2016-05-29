The Q Language Export Tool
========

Export go package to qlang module


The Q Language : https://github.com/qiniu/qlang


Usages:
```
qexport [-contexts=""] [-defctx=false] [skiperrimpl=true] [lower=false] [-outpath="./qlang"] packages

The packages for go package list or std for golang all standard packages.

  -contexts string
    	optional comma-separated list of <goos>-<goarch>[-cgo] to override default contexts.
  -defctx
    	optional use default context for build, default use all contexts.
  -lower
    	optional export name is first lower
  -outpath string
    	optional set export root path (default "./qlang")
  -skiperrimpl
    	optional skip error interface implement struct. (default true)```

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
	"reflect"
	"runtime"

	"qlang.io/qlang.spec.v1"
)

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"_name": "runtime",

	"Compiler": runtime.Compiler,
	"GOARCH":   runtime.GOARCH,
	"GOOS":     runtime.GOOS,

	"MemProfileRate": runtime.MemProfileRate,

	"BlockProfile":        runtime.BlockProfile,
	"Breakpoint":          runtime.Breakpoint,
	"CPUProfile":          runtime.CPUProfile,
	"Caller":              runtime.Caller,
	"Callers":             runtime.Callers,
	"GC":                  runtime.GC,
	"GOMAXPROCS":          runtime.GOMAXPROCS,
	"GOROOT":              runtime.GOROOT,
	"Goexit":              runtime.Goexit,
	"GoroutineProfile":    runtime.GoroutineProfile,
	"Gosched":             runtime.Gosched,
	"LockOSThread":        runtime.LockOSThread,
	"MemProfile":          runtime.MemProfile,
	"NumCPU":              runtime.NumCPU,
	"NumCgoCall":          runtime.NumCgoCall,
	"NumGoroutine":        runtime.NumGoroutine,
	"ReadMemStats":        runtime.ReadMemStats,
	"ReadTrace":           runtime.ReadTrace,
	"SetBlockProfileRate": runtime.SetBlockProfileRate,
	"SetCPUProfileRate":   runtime.SetCPUProfileRate,
	"SetFinalizer":        runtime.SetFinalizer,
	"Stack":               runtime.Stack,
	"StartTrace":          runtime.StartTrace,
	"StopTrace":           runtime.StopTrace,
	"ThreadCreateProfile": runtime.ThreadCreateProfile,
	"UnlockOSThread":      runtime.UnlockOSThread,
	"Version":             runtime.Version,

	"BlockProfileRecord": qlang.NewType(reflect.TypeOf((*runtime.BlockProfileRecord)(nil)).Elem()),
	"Func":               qlang.NewType(reflect.TypeOf((*runtime.Func)(nil)).Elem()),
	"FuncForPC":          runtime.FuncForPC,
	"MemProfileRecord":   qlang.NewType(reflect.TypeOf((*runtime.MemProfileRecord)(nil)).Elem()),
	"MemStats":           qlang.NewType(reflect.TypeOf((*runtime.MemStats)(nil)).Elem()),
	"StackRecord":        qlang.NewType(reflect.TypeOf((*runtime.StackRecord)(nil)).Elem()),
}

```

		