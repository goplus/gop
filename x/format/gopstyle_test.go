package format

import (
	"testing"
)

// -----------------------------------------------------------------------------

func testFormat(t *testing.T, name string, src, expect string) {
	t.Run(name, func(t *testing.T) {
		result, err := GopstyleSource([]byte(src), name)
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

// this is main
func main() {
	// say hello
	fmt.Println("Hello world")
}
`, `import "fmt"

// this is main

// say hello
fmt.Println("Hello world")
`)
}

func TestBasic2(t *testing.T) {
	testFormat(t, "hello world", `package main

import "fmt"

// this is main
func main() {
	// say hello
	fmt.Println("Hello world")
}

func f() {
}
`, `import "fmt"

// this is main
func main() {
	// say hello
	fmt.Println("Hello world")
}

func f() {
}
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
