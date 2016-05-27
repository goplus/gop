package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"qlang.io/qlang.v2/qlang"
	"qlang.io/qlang/terminal"

	qspec "qlang.io/qlang.spec.v1"
	qipt "qlang.io/qlang.v2/interpreter"
	qall "qlang.io/qlang/qlang.all"
)

var (
	historyFile = os.Getenv("HOME") + "/.qlang.history"
)

func main() {
	qall.InitSafe(false)
	qlang.Import("", qipt.Exports)
	qlang.SetDumpCode(os.Getenv("QLANG_DUMPCODE"))

	libs := os.Getenv("QLANG_PATH")
	if libs == "" {
		libs = os.Getenv("HOME") + "/qlang"
	}

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	lang.SetLibs(libs)

	// exec source
	if len(os.Args) > 1 {
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

	// interpreter
	qall.Copyright()

	var ret interface{}
	qlang.SetOnPop(func(v interface{}) {
		ret = v
	})

	term := terminal.New()
	term.LoadHistroy(historyFile) // load/save histroy
	defer term.SaveHistroy(historyFile)

	fnReadMore := func(expr string, line string) (string, bool) { // read more line check
		if strings.HasSuffix(line, "\\") {
			return expr + line[:len(line)-1], true
		}
		return expr + line + "\n", false
	}

	for {
		expr, err := term.Scan(">>> ", fnReadMore)
		if err != nil {
			if err == terminal.ErrPromptAborted {
				fmt.Println("Aborted")
				continue
			} else if err == io.EOF {
				break
			}
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		expr = strings.TrimSpace(expr)
		if expr == "" {
			continue
		}
		ret = qspec.Undefined
		err = lang.SafeEval(expr)
		if err != nil {
			fmt.Println(strings.TrimSpace(err.Error()))
			continue
		}
		if ret != qspec.Undefined {
			fmt.Println(ret)
		}
	}
}

// -----------------------------------------------------------------------------
