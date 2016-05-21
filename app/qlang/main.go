package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	qipt "qlang.io/qlang.v2/interpreter"
	"qlang.io/qlang.v2/qlang"
	qall "qlang.io/qlang/qlang.all"

	"gopkg.in/readline.v1"
)

// -----------------------------------------------------------------------------

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

	qall.Copyright()

	var ret interface{}
	var poped bool
	qlang.SetOnPop(func(v interface{}) {
		ret, poped = v, true
	})

	historyFile := os.Getenv("HOME") + "/.qlang.history"
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          ">>> ",
		HistoryFile:     historyFile,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				if len(line) == 0 {
					break
				} else {
					continue
				}
			} else if err == io.EOF {
				break
			}

			fmt.Fprintln(os.Stderr, err)
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		poped = false
		err = lang.SafeEval(line)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		if poped {
			fmt.Println(ret)
		}
	}
}

// -----------------------------------------------------------------------------
