package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strings"
	"os"

	"qlang.io/qlang.v1/qlang"
	qall "qlang.io/qlang/qlang.all"
	qipt "qlang.io/qlang.v1/interpreter"
)

// -----------------------------------------------------------------------------

func main() {

	qall.InitSafe(true)
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
			return
		}
		err = lang.SafeExec(b, fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
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

