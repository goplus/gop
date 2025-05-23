Overload Func/Method/Operator/Types
=====

### Overload Funcs

Define `overload func` in `inline func literal` style (see [overloadfunc1/add.xgo](demo/overloadfunc1/add.xgo)):

```go
func add = (
	func(a, b int) int {
		return a + b
	}
	func(a, b string) string {
		return a + b
	}
)

println add(100, 7)
println add("Hello", "World")
```

Define `overload func` in `ident` style (see [overloadfunc2/mul.xgo](demo/overloadfunc2/mul.xgo)):

```go
func mulInt(a, b int) int {
	return a * b
}

func mulFloat(a, b float64) float64 {
	return a * b
}

func mul = (
	mulInt
	mulFloat
)

println mul(100, 7)
println mul(1.2, 3.14)
```

### Overload Methods

Define `overload method` (see [overloadmethod/method.xgo](demo/overloadmethod/method.xgo)):

```go
type foo struct {
}

func (a *foo) mulInt(b int) *foo {
	println "mulInt"
	return a
}

func (a *foo) mulFoo(b *foo) *foo {
	println "mulFoo"
	return a
}

func (foo).mul = (
	(foo).mulInt
	(foo).mulFoo
)

var a, b *foo
var c = a.mul(100)
var d = a.mul(c)
```

### Overload Unary Operators

Define `overload unary operator` (see [overloadop1/overloadop.xgo](demo/overloadop1/overloadop.xgo)):

```go
type foo struct {
}

func -(a foo) (ret foo) {
	println "-a"
	return
}

func ++(a foo) {
	println "a++"
}

var a foo
var b = -a
a++
```

### Overload Binary Operators

Define `overload binary operator` (see [overloadop1/overloadop.xgo](demo/overloadop1/overloadop.xgo)):

```go
type foo struct {
}

func (a foo) * (b foo) (ret foo) {
	println "a * b"
	return
}

func (a foo) != (b foo) bool {
	println "a != b"
	return true
}

var a, b foo
var c = a * b
var d = a != b
```

However, `binary operator` usually need to support interoperability between multiple types. In this case it becomes more complex (see [overloadop2/overloadop.xgo](demo/overloadop2/overloadop.xgo)):

```go
type foo struct {
}

func (a foo) mulInt(b int) (ret foo) {
	println "a * int"
	return
}

func (a foo) mulFoo(b foo) (ret foo) {
	println "a * b"
	return
}

func intMulFoo(a int, b foo) (ret foo) {
	println "int * b"
	return
}

func (foo).* = (
	(foo).mulInt
	(foo).mulFoo
	intMulFoo
)

var a, b foo
var c = a * 10
var d = a * b
var e = 10 * a
```

### Overload Types

TODO

### Overload Typecast

TODO

