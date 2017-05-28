package main

import (
	"fmt"
	"os"
	"strings"

	"qlang.io/cl/interpreter"
	"qlang.io/cl/qlang"
	"qlang.io/exec"
	"qlang.io/lib/eqlang"
	"qlang.io/lib/qlang.all"
)

// -----------------------------------------------------------------------------

const usage = `
Usage:
    eql <templatefile> [-i -j <jsoninput> -o <outputfile> --key1=value1 --key2=value2 ...]
    eql <templatedir> [-i -j <jsoninput> -o <outputdir> --key1=value1 --key2=value2 ...]

<templatefile>: the template file to format.
<templatedir>:  the template dir to format.

-i:              use stdin as json input.
-j <jsoninput>:  json input string.
-o <outputfile>: output result into this file, default is stdout.
-o <outputdir>:  output result into this directory, default is <templatedir> format result.
--key=value:     (key, value) pair that overrides json input specified by -i or -j <jsoninput>.
`

func main() {

	qall.InitSafe(false)
	qlang.Import("", interpreter.Exports)
	qlang.SetDumpCode(os.Getenv("QLANG_DUMPCODE"))

	libs := os.Getenv("QLANG_PATH")
	if libs == "" {
		libs = os.Getenv("HOME") + "/qlang"
	}

	lang := qlang.New()
	lang.SetLibs(libs)

	vars = lang.Context
	eql := eqlang.New(lang)

	parseFlags()
	if source == "" {
		fmt.Fprintln(os.Stderr, usage)
		return
	}

	fi, err := os.Stat(source)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-2)
	}

	if fi.IsDir() {
		global := lang.CopyVars()
		err = eql.ExecuteDir(global, source, output)
	} else {
		err = eql.ExecuteFile(source, output)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-3)
	}
}

// -----------------------------------------------------------------------------

var (
	source string
	output string
	vars   *exec.Context
)

func parseFlags() {

	vars.SetVar("imports", "")
	for i := 1; i < len(os.Args); i++ {
		switch arg := os.Args[i]; arg {
		case "-o":
			if i+1 >= len(os.Args) {
				fmt.Fprintln(os.Stderr, "ERROR: switch -o doesn't have parameters, please use -o <output>")
				os.Exit(10)
			}
			output = os.Args[i+1]
			i++
		case "-i":
			ret, err := eqlang.InputFile("-")
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: stdin isn't valid json - %v\n", err)
				os.Exit(12)
			}
			vars.ResetVars(ret)
		case "-j":
			if i+1 >= len(os.Args) {
				fmt.Fprintln(os.Stderr, "ERROR: switch -j doesn't have parameters, please use -j <jsoninput>")
				os.Exit(10)
			}
			ret, err := eqlang.Input(os.Args[i+1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: invalid parameter `-j <jsoninput>` - %v\n", err)
				os.Exit(12)
			}
			vars.ResetVars(ret)
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
