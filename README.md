Q Language - A script language for Go
========

[![LICENSE](https://img.shields.io/github/license/qiniu/qlang.svg)](https://github.com/qiniu/qlang/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/qiniu/qlang.svg?branch=master)](https://travis-ci.org/qiniu/qlang)
[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/qlang)](https://goreportcard.com/report/github.com/qiniu/qlang)
[![GitHub release](https://img.shields.io/github/v/tag/qiniu/qlang.svg?label=release)](https://github.com/qiniu/qlang/releases)
[![Coverage Status](https://codecov.io/gh/qiniu/qlang/branch/master/graph/badge.svg)](https://codecov.io/gh/qiniu/qlang)
[![GoDoc](https://img.shields.io/badge/Godoc-reference-blue.svg)](https://godoc.org/github.com/qiniu/qlang)

[![Qiniu Logo](http://open.qiniudn.com/logo.png)](http://www.qiniu.com/)

最新的 qlang v6 版本在能力上较以前的版本有极大的调整。其核心变化为：

- 完全推翻重来，从动态类型转向静态类型！
- 完全兼容 Go 语言文法。
- 在 Go 语言兼容基础上，保留当初 qlang 动态类型版本的重要特性。比如：

```
a := [1, 2, 3.4]
// 等价于 Go 中的  a := []float64{1, 2, 3.4}

b := {"a": 1, "b": 3.0}
// 等价于  b := map[string]float64{"a": 1, "b": 3.0}

c := {"a": 1, "b": "Hello"}
// 等价于 c := map[string]interface{}{"a": 1, "b": "Hello"}
```

当然也会放弃一些特性，比如：

```
a = 1   // 需要改为 a := 1，放弃该特性是为了让编译器更好地发现变量名冲突。
```

关于新版本的详细规划，参考：

* https://github.com/qiniu/qlang/issues/176

代码样例：

* https://github.com/qiniu/qlang/tree/v6.x/tutorial


## 老版本

当前 qlang v6 还在快速迭代中。在正式场合建议使用正式 release 的版本：

* https://github.com/qiniu/qlang/releases

最近的老版本代码可以从 qlang v1.5 分支获得：

* https://github.com/qiniu/qlang/tree/v1.5
