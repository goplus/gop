package format

import (
	"bytes"
	"testing"

	"github.com/goplus/gop/format"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
)

// -----------------------------------------------------------------------------

func testFormat(t *testing.T, name string, src, expect string) {
	t.Run(name, func(t *testing.T) {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "foo.gop", src, parser.ParseComments)
		if err != nil {
			t.Fatal("parser.ParseFile failed:", err)
		}
		Gopstyle(f)
		var buf bytes.Buffer
		err = format.Node(&buf, fset, f)
		if err != nil {
			t.Fatal("format.Node failed:", err)
		}
		if ret := buf.String(); ret != expect {
			t.Fatalf("%s => Expect:\n%s\n=> Got:\n%s\n", name, expect, ret)
		}
	})
}

func TestBasic1(t *testing.T) {
	testFormat(t, "hello world", `package main

import "fmt"

func main() {
	fmt.Println("Hello world")
}
`, `import "fmt"

fmt.Println("Hello world")
`)
}

func TestBasic2(t *testing.T) {
	testFormat(t, "hello world", `package main

import "fmt"

func main() {
	fmt.Println("Hello world")
}

func f() {
}
`, `import "fmt"

func f() {
}

fmt.Println("Hello world")
`)
}

func TestBasic3(t *testing.T) {
	testFormat(t, "hello world", `package main

import "fmt"

func f() {
}
`, `import "fmt"

func f() {
}
`)
}

// -----------------------------------------------------------------------------
