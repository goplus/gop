package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strings"
	"os"

	"qlang.io/qlang.v2/qlang"
	qall "qlang.io/qlang/qlang.all"
)

// -----------------------------------------------------------------------------

func main() {

	qall.InitSafe(false)

	libs := os.Getenv("QLANG_PATH")
	if libs == "" {
		libs = os.Getenv("HOME") + "/qlang"
	}
	if os.Getenv("QLANG_DUMPCODE") != "" {
		qlang.DumpCode = true
	}

	if len(os.Args) > 1 {
		lang, err := qlang.New(qlang.InsertSemis)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		lang.SetLibs(libs)
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

	var ret interface{}
	qlang.SetOnPop(func(v interface{}) {
		ret = v
	})

	lang, err := qlang.New(nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	lang.SetLibs(libs)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \t\r\n")
		if line == "" {
			continue
		}
		ret = nil
		err := lang.SafeEval(line)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		fmt.Printf("> %v\n\n", ret)
	}
}

// -----------------------------------------------------------------------------

