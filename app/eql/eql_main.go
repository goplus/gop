package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"qlang.io/exec.v2"
	"qlang.io/qlang.v2/qlang"

	qipt "qlang.io/qlang.v2/interpreter"
	"qlang.io/qlang/eql.v1"
	qall "qlang.io/qlang/qlang.all"
)

// -----------------------------------------------------------------------------

const usage = `
Usage:
    eql example.eql -o <output.txt> --key1=value1 --key2=value2 ...
`

// eql example.eql -o <output.txt> --key1=value1 --key2=value2 ...
//
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
		os.Exit(-1)
	}
	lang.SetLibs(libs)

	vars = lang.Context
	eql.DefaultVars = vars

	paseFlags()
	if source == "" {
		fmt.Fprintln(os.Stderr, usage)
		return
	}

	b, err := ioutil.ReadFile(source)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	code, err := eql.Parse(string(b))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	if output != "" {
		f, err := os.Create(output)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(3)
		}
		defer f.Close()
		os.Stdout = f
	}

	err = lang.SafeExec(code, source)
	if err != nil {
		os.Remove(output)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(4)
	}
}

// -----------------------------------------------------------------------------

var (
	source string
	output string
	vars   *exec.Context
)

func paseFlags() {

	vars.SetVar("imports", "")
	for i := 0; i < len(os.Args); i++ {
		switch arg := os.Args[i]; arg {
		case "-o":
			if i+1 >= len(os.Args) {
				fmt.Fprintln(os.Stderr, "ERROR: switch -o doesn't have parameters, please use -o <output>")
				os.Exit(10)
			}
			output = os.Args[i+1]
			i++
		default:
			if strings.HasPrefix(arg, "--") {
				kv := arg[2:]
				pos := strings.Index(kv, "=")
				if pos < 0 {
					fmt.Fprintf(os.Stderr, "ERROR: invalid switch `%s`\n", arg)
					os.Exit(11)
				}
				vars.SetVar(kv[:pos], kv[pos+1:])
			} else if arg[0] == '-' {
				fmt.Fprintf(os.Stderr, "ERROR: unknown switch `%s`\n", arg)
				os.Exit(12)
			} else {
				source = arg
			}
		}
	}
}

// -----------------------------------------------------------------------------
