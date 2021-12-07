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
`, `// this is main

// say hello
println("Hello world")
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
`, `// this is main
func main() {
	// say hello
	println("Hello world")
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
`, `func f() {
}
`)
}

// -----------------------------------------------------------------------------

func TestPrint(t *testing.T) {
	testFormat(t, "print", `package main

import "fmt"

func f() {
	fmt.Print("hello")
}
`, `func f() {
	print("hello")
}
`)
}

func TestPrintf(t *testing.T) {
	testFormat(t, "print", `package main

import "fmt"

func f() {
	fmt.Printf("hello")
}
`, `func f() {
	printf("hello")
}
`)
}

func TestPrintln(t *testing.T) {
	testFormat(t, "print", `package main

import "fmt"

func f() {
	fmt.Println("hello")
}
`, `func f() {
	println("hello")
}
`)
}

func TestPrintlnGroup(t *testing.T) {
	testFormat(t, "print", `package main

import (
	"fmt"
)

func f() {
	fmt.Println("hello")
}
`, `func f() {
	println("hello")
}
`)
}

func TestPrintlnWithOtherFmtCalls(t *testing.T) {
	testFormat(t, "print", `package main

import "errors"
import "fmt"

func f() {
	fmt.Errorf("%w", errors.New("hello"))
	fmt.Println("hello")
}
`, `import "errors"

func f() {
	errorf("%w", errors.New("hello"))
	println("hello")
}
`)
}

func TestPrintlnWithOtherFmtCallsGroup(t *testing.T) {
	testFormat(t, "print", `package main

import (
	"errors"
	"fmt"
)

func f() {
	fmt.Errorf("%w", errors.New("hello"))
	fmt.Println("hello")
}
`, `import (
	"errors"
)

func f() {
	errorf("%w", errors.New("hello"))
	println("hello")
}
`)
}

func TestPrintlnWithOtherFmtCallsWithAssign(t *testing.T) {
	testFormat(t, "print", `package main

import "errors"
import "fmt"

func f() {
	_ = fmt.Errorf("%w", errors.New("hello"))
	fmt.Println("hello")
}
`, `import "errors"

func f() {
	_ = errorf("%w", errors.New("hello"))
	println("hello")
}
`)
}

func TestPrintlnWithOtherFmtCallsWithGroupWithAssign(t *testing.T) {
	testFormat(t, "print", `package main

import (
	"errors"
	"fmt"
)

func f() {
	_ = fmt.Errorf("%w", errors.New("hello"))
	fmt.Println("hello")
}
`, `import (
	"errors"
)

func f() {
	_ = errorf("%w", errors.New("hello"))
	println("hello")
}
`)
}

func TestPrintlnWithOtherFmtDecls(t *testing.T) {
	testFormat(t, "print", `package main

import "fmt"

func f() {
	_ = fmt.Stringer
	fmt.Println("hello")
}
`, `import "fmt"

func f() {
	_ = fmt.Stringer
	println("hello")
}
`)
}

func TestPrintlnWithOtherFmtVars(t *testing.T) {
	testFormat(t, "print", `package main

import "fmt"

func f() {
	var _ fmt.Stringer
	fmt.Println("hello")
}
`, `import "fmt"

func f() {
	var _ fmt.Stringer
	println("hello")
}
`)
}

func TestPrintlnImportAlias(t *testing.T) {
	testFormat(t, "print", `package main

import fmt1 "fmt"

func f() {
	fmt1.Println("hello")
}
`, `func f() {
	println("hello")
}
`)
}

func TestPrintlnImportMultiAliases(t *testing.T) {
	testFormat(t, "print", `package main

import (
	fmt1 "fmt"
	fmt2 "fmt"
)

func f() {
	fmt1.Println(1)
	fmt2.Println(2)
}
`, `func f() {
	println(1)
	println(2)
}
`)
}

func TestPrintlnImportMultiAliasesDifferentGroups(t *testing.T) {
	testFormat(t, "print", `package main

import "fmt"

import (
	fmt1 "fmt"
	fmt2 "fmt"
)

func f() {
	var _ fmt.Stringer
	fmt1.Println(1)
	fmt2.Println(2)
}
`, `import "fmt"

func f() {
	var _ fmt.Stringer
	println(1)
	println(2)
}
`)
}
