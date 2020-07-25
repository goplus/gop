package repl

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/cmd/internal/base"
	exec "github.com/goplus/gop/exec/bytecode"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
	"github.com/peterh/liner"
)

func init() {
	Cmd.Run = runCmd
}

var Cmd = &base.Command{
	UsageLine: "repl",
	Short:     "Play Go+ in console",
}

type repl struct {
	src          string
	preContext   exec.Context
	ip           int
	prompt       string
	continueMode bool
}

const (
	continuePrompt string = "..."
	standardPrompt string = ">>>"
	welcome        string = "welcome to Go+ console!"
)

func runCmd(cmd *base.Command, args []string) {
	fmt.Println(welcome)
	liner := liner.NewLiner()
	defer liner.Close()
	rInstance := repl{prompt: standardPrompt}
	for {
		l, err := liner.Prompt(string(rInstance.prompt))
		if err != nil {
			if err == io.EOF {
				fmt.Printf("\n")
				break
			}
			fmt.Printf("Problem reading line: %v\n", err)
			continue
		}
		if l != "" {
			// mainly for scroll up
			liner.AppendHistory(l)
		}
		if rInstance.continueMode {
			if l != "" {
				rInstance.src += string(l) + "\n"
				continue
			}
			// input nothing means jump out continue mode
			rInstance.continueMode = false
			rInstance.prompt = standardPrompt
		}

		if nil == rInstance.run(string(l)) {
			rInstance.src += string(l) + "\n"
		}
	}
}

// run execute the input line
func (r *repl) run(newLine string) (err error) {
	src := r.src + newLine + "\n"
	defer func() {
		if errR := recover(); errR != nil {
			replErr(newLine, errR)
			err = errors.New("panic err")
			// TODO: Need a better way to log and show the stack when crash
			// It is too long to print stack on terminal even only print part of the them:
		}
	}()
	fset := token.NewFileSet()
	pkgs, err := parser.ParseGopFiles(fset, "", false, src, 0)
	if err != nil {
		if strings.Contains(err.Error(), `expected ')', found 'EOF'`) ||
			strings.Contains(err.Error(), "expected '}', found 'EOF'") {
			r.prompt = continuePrompt
			r.continueMode = true
			err = nil
			return
		}
		fmt.Println("ParseGopFiles err", err)
		return
	}
	cl.CallBuiltinOp = exec.CallBuiltinOp

	b := exec.NewBuilder(nil)

	_, err = cl.NewPackage(b.Interface(), pkgs["main"], fset, cl.PkgActClMain)
	if err != nil {
		if err == cl.ErrMainFuncNotFound {
			err = nil
			return
		}
		fmt.Printf("NewPackage err %+v\n", err)
		return
	}
	code := b.Resolve()
	ctx := exec.NewContext(code)
	if r.ip != 0 {
		// not first time, restore pre var
		r.preContext.CloneSetterScope(ctx)
	}
	currentIp := ctx.Exec(r.ip, code.Len())
	r.preContext = *ctx
	r.ip = currentIp - 1
	return
}
func replErr(line string, err interface{}) {
	fmt.Println("code run fail : ", line)
	fmt.Println("repl err: ", err)
}
