The Q Language Export Tool
========

Export go package to qlang module


The Q Language : https://github.com/qiniu/qlang


### Usages:

```
Export go packages to qlang modules.

Usage:
  qexport [option] packages

The packages for go package list or std for golang all standard packages.

  -contexts string
    	optional comma-separated list of <goos>-<goarch>[-cgo] to override default contexts.
  -convnew
    	optional convert NewType func to type func (default true)
  -defctx
    	optional use default context for build, default use all contexts.
  -lowercase
    	optional use qlang lower case style. (default true)
  -outpath string
    	optional set export root path (default "./qlang")
  -skiperrimpl
    	optional skip error interface implement struct. (default true)
  -updatepath string
    	option set update qlang package root
```

### Examples

```
export sync
> qexport sync

export html and html/template package
> qexport html html/template

export all package
> qexport std

export io and update from qlang.io/qlang/io
> qexport -updatepath qlang.io/qlang io

export all package and update from qlang.io
> qexport -updatepath qlang.io/qlang std
```

### 导出包

```
导出一个包
> qexport bufio

导出多个包
> qexport bufio io io/ioutil

导出标准包
> qexport std 
```

### 更新包
```
导出bufio包，复制qlang.io/qlang中bufio包到输出目录并更新。
> qexport -updatepath qlang.io/qlang bufio

导出多个包，复制qlang.io/qlang中对应包到输出目录并作更新。
> qexport -updatepath qlang.io/qlang bufio io io/ioutil

导出标准包，复制qlang.io/qlang中对应包到输出目录并作更新。
> qexport -updatepath qlang.io/qlang std

```