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


Export pacakge go/build:

``` go
package build

import (
	"go/build"
)

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"_name": "go/build",

	"allowBinary":   build.AllowBinary,
	"findOnly":      build.FindOnly,
	"ignoreVendor":  build.IgnoreVendor,
	"importComment": build.ImportComment,

	"default": build.Default,
	"toolDir": build.ToolDir,

	"archChar":      build.ArchChar,
	"isLocalImport": build.IsLocalImport,

	"context":                   newContext,
	"contextArray":              newContextArray,
	"multiplePackageError":      newMultiplePackageError,
	"multiplePackageErrorArray": newMultiplePackageErrorArray,
	"noGoError":                 newNoGoError,
	"noGoErrorArray":            newNoGoErrorArray,
	"import":                    build.Import,
	"importDir":                 build.ImportDir,
}

func newContext() *build.Context {
	return new(build.Context)
}

func newContextArray(n int) []build.Context {
	return make([]build.Context, n)
}

func newMultiplePackageError() *build.MultiplePackageError {
	return new(build.MultiplePackageError)
}

func newMultiplePackageErrorArray(n int) []build.MultiplePackageError {
	return make([]build.MultiplePackageError, n)
}

func newNoGoError() *build.NoGoError {
	return new(build.NoGoError)
}

func newNoGoErrorArray(n int) []build.NoGoError {
	return make([]build.NoGoError, n)
}
```

		