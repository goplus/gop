package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/qiniu/qlang"
	qipt "github.com/qiniu/qlang/cl/interpreter"
	qall "github.com/qiniu/qlang/lib/qlang.all"
)

// -----------------------------------------------------------------------------

func main() {

	qall.InitSafe(false)
	qlang.Import("", qipt.Exports)

	if len(os.Args) > 1 {
		lang, err := qlang.New(qlang.InsertSemis)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fname := os.Args[1]
		b, err := ioutil.ReadFile(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		err = lang.SafeExec(b, fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(3)
		}
		return
	}

	qall.Copyright()
	lang, err := qlang.New(nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \t\r\n")
		if line == "" {
			continue
		}
		err := lang.SafeEval(line)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		v, _ := lang.Ret()
		fmt.Printf("> %v\n\n", v)
	}
}

// -----------------------------------------------------------------------------
