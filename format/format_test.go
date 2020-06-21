// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package format

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/qiniu/goplus/parser"
	"github.com/qiniu/goplus/token"
)

const testfile = "format_test.go"

func diff(t *testing.T, dst, src []byte) {
	line := 1
	offs := 0 // line offset
	for i := 0; i < len(dst) && i < len(src); i++ {
		d := dst[i]
		s := src[i]
		if d != s {
			t.Errorf("dst:%d: %s\n", line, dst[offs:i+1])
			t.Errorf("src:%d: %s\n", line, src[offs:i+1])
			return
		}
		if s == '\n' {
			line++
			offs = i + 1
		}
	}
	if len(dst) != len(src) {
		t.Errorf("len(dst) = %d, len(src) = %d\nsrc = %q", len(dst), len(src), src)
	}
}

func TestNode(t *testing.T) {
	src, err := ioutil.ReadFile(testfile)
	if err != nil {
		t.Fatal(err)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, testfile, src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer

	if err = Node(&buf, fset, file); err != nil {
		t.Fatal("Node failed:", err)
	}

	diff(t, buf.Bytes(), src)
}

func TestSource(t *testing.T) {
	src, err := ioutil.ReadFile(testfile)
	if err != nil {
		t.Fatal(err)
	}

	res, err := Source(src)
	if err != nil {
		t.Fatal("Source failed:", err)
	}

	diff(t, res, src)
}

var (
	gop_src1 = `package main

import (
	"strconv"
)

func add(x, y string) (int, error) {
	return strconv.Atoi(x)? + strconv.Atoi(y)?, nil
}

func addSafe(x, y string) int {
	return strconv.Atoi(x)?:0 + strconv.Atoi(y)?:0
}

func test() {
	println(add("100", "23")!)

	sum, err := add("10", "abc")
	println(sum, err)

	println(addSafe("10", "abc"))
}

func main() {
	println("Hello, Go+")
	println(1r << 129)
	println(1/3r + 2/7r*2)

	arr := [1, 3, 5, 7, 11, 13, 17, 19]
	println(arr)

	e := [[a, b] for a <- arr, a < b for b <- arr, b > 2]

	x := {x: i for i, x <- [1, 3, 5, 7, 11]}
	y := {x: i for i, x <- [1, 3, 5, 7, 11], i%2 == 1}
	z := {v: k for k, v <- {1: "Hello", 3: "Hi", 5: "xsw", 7: "Go+"}, k > 3}

	m := {"Hi": 1, "Go+": 2}
	println(m)
	println({v: k for k, v <- m})
	println([k for k, _ <- m])
	println([v for v <- m])

	sum := 0
	for x <- [1, 3, 5, 7, 11], x > 3 {
		sum += x
	}
	println("sum(5,7,11):", sum)
}
`
	gop_src2 = `println("Hello, Go+")
println(1r << 129)
println(1/3r + 2/7r*2)

arr := [1, 3, 5, 7, 11, 13, 17, 19]
println(arr)
println([x * x for x <- arr, x > 3])

m := {"Hi": 1, "Go+": 2}
println(m)
println({v: k for k, v <- m})
println([k for k, _ <- m])
println([v for v <- m])`
)

func TestGopSource(t *testing.T) {
	src := []byte(gop_src2)
	res, err := Source(src)
	if err != nil {
		t.Fatal("Source failed:", err)
	}
	fmt.Println(string(res))
	diff(t, res, src)
}

// Test cases that are expected to fail are marked by the prefix "ERROR".
// The formatted result must look the same as the input for successful tests.
var tests = []string{
	// declaration lists
	`import "go/format"`,
	"var x int",
	"var x int\n\ntype T struct{}",

	// statement lists
	"x := 0",
	"f(a, b, c)\nvar x int = f(1, 2, 3)",

	// indentation, leading and trailing space
	"\tx := 0\n\tgo f()",
	"\tx := 0\n\tgo f()\n\n\n",
	"\n\t\t\n\n\tx := 0\n\tgo f()\n\n\n",
	"\n\t\t\n\n\t\t\tx := 0\n\t\t\tgo f()\n\n\n",
	"\n\t\t\n\n\t\t\tx := 0\n\t\t\tconst s = `\nfoo\n`\n\n\n",     // no indentation added inside raw strings
	"\n\t\t\n\n\t\t\tx := 0\n\t\t\tconst s = `\n\t\tfoo\n`\n\n\n", // no indentation removed inside raw strings

	// comments
	"/* Comment */",
	"\t/* Comment */ ",
	"\n/* Comment */ ",
	"i := 5 /* Comment */",         // issue #5551
	"\ta()\n//line :1",             // issue #11276
	"\t//xxx\n\ta()\n//line :2",    // issue #11276
	"\ta() //line :1\n\tb()\n",     // issue #11276
	"x := 0\n//line :1\n//line :2", // issue #11276

	// whitespace
	"",     // issue #11275
	" ",    // issue #11275
	"\t",   // issue #11275
	"\t\t", // issue #11275
	"\n",   // issue #11275
	"\n\n", // issue #11275
	"\t\n", // issue #11275

	// erroneous programs
	"ERROR1 + 2 +",
	"ERRORx :=  0",
}

func String(s string) (string, error) {
	res, err := Source([]byte(s))
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func TestPartial(t *testing.T) {
	for _, src := range tests {
		if strings.HasPrefix(src, "ERROR") {
			// test expected to fail
			src = src[5:] // remove ERROR prefix
			res, err := String(src)
			if err == nil && res == src {
				t.Errorf("formatting succeeded but was expected to fail:\n%q", src)
			}
		} else {
			// test expected to succeed
			res, err := String(src)
			if err != nil {
				t.Errorf("formatting failed (%s):\n%q", err, src)
			} else if res != src {
				t.Errorf("formatting incorrect:\nsource: %q\nresult: %q", src, res)
			}
		}
		break
	}
}
