package format

import (
	"testing"
)

// -----------------------------------------------------------------------------

func testFormat(t *testing.T, name string, src, expect string) {
	t.Run(name, func(t *testing.T) {
		result, err := Source([]byte(src))
		if err != nil {
			t.Fatal("format.Source failed:", err)
		}
		if ret := string(result); ret != expect {
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
