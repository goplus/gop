package qshell

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"qlang.io/cl/qlang"
	"qlang.io/lib/terminal"

	qipt "qlang.io/cl/interpreter"
	qall "qlang.io/lib/qlang.all"
)

var (
	historyFile = os.Getenv("HOME") + "/.qlang.history"
	notFound    = interface{}(errors.New("not found"))
)

// Main is the entry of qlang/qlang.safe shell
//
func Main(safeMode bool) {

	qall.InitSafe(safeMode)
	qlang.Import("", qipt.Exports)
	qlang.Import("qlang", qlang.Exports)
	qlang.SetDumpCode(os.Getenv("QLANG_DUMPCODE"))

	libs := os.Getenv("QLANG_PATH")
	if libs == "" {
		libs = os.Getenv("HOME") + "/qlang"
	}

	lang := qlang.New()
	lang.SetLibs(libs)

	// exec source
	//
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
	if safeMode {
		fmt.Printf("Use Ctrl-D (i.e. EOF) to exit.\n\n")
	} else {
		fmt.Printf("Use exit() or Ctrl-D (i.e. EOF) to exit.\n\n")
	}

	var ret interface{}
	qlang.SetOnPop(func(v interface{}) {
		ret = v
	})

	var tokener tokener
	term := terminal.New(">>> ", "... ", tokener.ReadMore)
	term.SetWordCompleter(func(line string, pos int) (head string, completions []string, tail string) {
		return line[:pos], []string{"  "}, line[pos:]
	})

	term.LoadHistroy(historyFile) // load/save histroy
	defer term.SaveHistroy(historyFile)

	for {
		expr, err := term.Scan()
		if err != nil {
			if err == terminal.ErrPromptAborted {
				continue
			} else if err == io.EOF {
				fmt.Println("^D")
				break
			}
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		expr = strings.TrimSpace(expr)
		if expr == "" {
			continue
		}
		ret = notFound
		err = lang.SafeEval(expr)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		if ret != notFound {
			fmt.Println(ret)
		}
	}
}

// -----------------------------------------------------------------------------

type tokener struct {
	level int
	instr bool
}

var dontReadMoreChars = "+-})];"
var puncts = "([=,*/%|&<>^.:"

func readMore(line string) bool {

	n := len(line)
	if n == 0 {
		return false
	}

	pos := strings.IndexByte(dontReadMoreChars, line[n-1])
	if pos == 0 || pos == 1 {
		return n >= 2 && line[n-2] != dontReadMoreChars[pos]
	}
	return pos < 0 && strings.IndexByte(puncts, line[n-1]) >= 0
}

func findEnd(line string, c byte) int {

	for i := 0; i < len(line); i++ {
		switch line[i] {
		case c:
			return i
		case '\\':
			i++
		}
	}
	return -1
}

func (p *tokener) ReadMore(expr string, line string) (string, bool) { // read more line check

	ret := expr + line + "\n"
	for {
		if p.instr {
			pos := strings.IndexByte(line, '`')
			if pos < 0 {
				return ret, true
			}
			line = line[pos+1:]
			p.instr = false
		}

		pos := strings.IndexAny(line, "{}`'\"")
		if pos < 0 {
			if p.level != 0 {
				return ret, true
			}
			line = strings.TrimRight(line, " \t")
			return ret, readMore(line)
		}
		switch c := line[pos]; c {
		case '{':
			p.level++
		case '}':
			p.level--
		case '`':
			p.instr = true
		default:
			line = line[pos+1:]
			pos = findEnd(line, c)
			if pos < 0 {
				return ret, p.level != 0
			}
		}
		line = line[pos+1:]
	}
}

// -----------------------------------------------------------------------------
