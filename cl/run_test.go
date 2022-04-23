package cl_test

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/parser/parsertest"
	"github.com/goplus/gop/scanner"
	"github.com/goplus/gox"
)

var (
	tmpDir     string
	tmpFileIdx int64
)

func init() {
	home, err := os.Getwd()
	check(err)

	tmpDir = home + "/.gop/tmp/"
	err = os.MkdirAll(tmpDir, 0755)
	check(err)
}

func check(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func checkWith(err error, stdout, stderr io.Writer) {
	if err != nil {
		fatalWith(err, stdout, stderr)
	}
}

func fatalWith(err error, stdout, stderr io.Writer) {
	if o, ok := getBytes(stdout, stderr); ok {
		os.Stderr.Write(o.Bytes())
	}
	log.Panicln(err)
}

type iBytes interface {
	Bytes() []byte
}

func getBytes(stdout, stderr io.Writer) (o iBytes, ok bool) {
	if o, ok = stderr.(iBytes); ok {
		return
	}
	o, ok = stdout.(iBytes)
	return
}

func goRun(t *testing.T, code []byte) string {
	idx := atomic.AddInt64(&tmpFileIdx, 1)
	infile := tmpDir + strconv.FormatInt(idx, 10) + ".go"
	err := os.WriteFile(infile, []byte(code), 0666)
	check(err)

	var stdout, stderr bytes.Buffer
	cmd := exec.Command("go", "run", infile)
	cmd.Dir = tmpDir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	os.Remove(infile)
	checkWith(err, &stdout, &stderr)
	return stdout.String()
}

func genGo(t *testing.T, conf *cl.Config, gopcode string) []byte {
	cl.SetDisableRecover(true)
	defer cl.SetDisableRecover(false)

	fs := parsertest.NewSingleFileFS("/foo", "bar.gop", gopcode)
	pkgs, err := parser.ParseFSDir(gblFset, fs, "/foo", parser.Config{Mode: parser.ParseComments})
	if err != nil {
		scanner.PrintError(os.Stderr, err)
		t.Fatal("ParseFSDir:", err)
	}
	pkg, err := cl.NewPackage("", pkgs["main"], conf)
	if err != nil {
		t.Fatal("NewPackage:", err)
	}
	var b bytes.Buffer
	err = gox.WriteTo(&b, pkg, false)
	if err != nil {
		t.Fatal("gox.WriteTo failed:", err)
	}
	return b.Bytes()
}

func testRun(t *testing.T, gopcode, expected string) {
	code := genGo(t, gblConf, gopcode)
	result := goRun(t, code)
	if result != expected {
		t.Fatalf("=> Result:\n%s\n=> Expected:\n%s\n", result, expected)
	}
}

// -----------------------------------------------------------------------------

func TestUint128_println(t *testing.T) {
	testRun(t, `
var x uint128 = 1
println x
`, `1
`)
}

func TestUint128_init(t *testing.T) {
	testRun(t, `
var x uint128 = 1 << 65
var y = x + 1
println x
println y
`, `36893488147419103232
36893488147419103233
`)
}

func TestUint128_cast(t *testing.T) {
	testRun(t, `
println uint128(1 << 65), uint128()
`, `36893488147419103232 0
`)
}

func TestUint128_printf(t *testing.T) {
	testRun(t, `
var x uint128 = 1
printf "%4d\n", x
`, `   1
`)
}

func TestUint128_scanf(t *testing.T) {
	testRun(t, `
import "fmt"

var name string
var age uint128
fmt.Sscanf("Kim is 22 years old", "%s is %d years old", &name, &age)
println name, age
`, `Kim 22
`)
}

// -----------------------------------------------------------------------------

func TestBigint_println(t *testing.T) {
	testRun(t, `
var x bigint = 1
println x
`, `1
`)
}

func TestBigint_init(t *testing.T) {
	testRun(t, `
var x bigint = 1 << 65
var y = x + 1
println x
println y
`, `36893488147419103232
36893488147419103233
`)
}

func TestBigint_cast(t *testing.T) {
	testRun(t, `
println bigint(1 << 65), bigint()
`, `36893488147419103232 0
`)
}

func TestBigint_printf(t *testing.T) {
	testRun(t, `
var x bigint = 1
printf "%4d\n", x
`, `   1
`)
}

// -----------------------------------------------------------------------------
