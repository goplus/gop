The Q Language Export Tool
========

Export go package to qlang module


The Q Language : https://github.com/qiniu/qlang


Usages:
```
qexport [-contexts=""] [-defctx=false] [-convnew=true] [-skiperrimpl=true] [-outpath="./qlang"] packages

The packages for go package list or std for golang all standard packages.

  -contexts string
    	optional comma-separated list of <goos>-<goarch>[-cgo] to override default contexts.
  -convnew
    	optional convert NewType func to type func (default true)
  -defctx
    	optional use default context for build, default use all contexts.
  -outpath string
    	optional set export root path (default "./qlang")
  -skiperrimpl
    	optional skip error interface implement struct. (default true)
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

* runtime.go
``` go
package runtime

import (
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
	"setBlockProfileRate": runtime.SetBlockProfileRate,
	"setCPUProfileRate":   runtime.SetCPUProfileRate,
	"setFinalizer":        runtime.SetFinalizer,
	"stack":               runtime.Stack,
	"threadCreateProfile": runtime.ThreadCreateProfile,
	"unlockOSThread":      runtime.UnlockOSThread,
	"version":             runtime.Version,

	"BlockProfileRecord": qlang.StructOf((*runtime.BlockProfileRecord)(nil)),
	"Func":               qlang.StructOf((*runtime.Func)(nil)),
	"funcForPC":          runtime.FuncForPC,
	"MemProfileRecord":   qlang.StructOf((*runtime.MemProfileRecord)(nil)),
	"MemStats":           qlang.StructOf((*runtime.MemStats)(nil)),
	"StackRecord":        qlang.StructOf((*runtime.StackRecord)(nil)),
}
```
* runtime-go15.go
``` go
// +build go1.5

package runtime

import "runtime"

func init() {
	Exports["readTrace"] = runtime.ReadTrace
	Exports["startTrace"] = runtime.StartTrace
	Exports["stopTrace"] = runtime.StopTrace
}
```
		