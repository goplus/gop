/*
 Copyright 2020 The GoPlus Authors (goplus.org)

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

// Package repl implements the ``gop repl'' command.
package repl

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
	"github.com/peterh/liner"

	exec "github.com/goplus/gop/exec/bytecode"
)

func init() {
	Cmd.Run = runCmd
}

// Cmd - gop repl
var Cmd = &base.Command{
	UsageLine: "repl",
	Short:     "Play Go+ in console",
}

type repl struct {
	src          string       // the whole source code from repl
	preContext   exec.Context // store the context after exec
	ip           int          // store the ip after exec
	prompt       string       // the prompt type in console
	continueMode bool         // switch to control the promot type
	liner        *liner.State // liner instance
}

const (
	continuePrompt string = "... "
	standardPrompt string = ">>> "
	welcome        string = "welcome to Go+ console!"
)

func runCmd(cmd *base.Command, args []string) {
	fmt.Println(welcome)
	liner := liner.NewLiner()
	rInstance := repl{prompt: standardPrompt, liner: liner}
	defer rInstance.liner.Close()

	for {
		l, err := rInstance.liner.Prompt(string(rInstance.prompt))
		if err != nil {
			if err == io.EOF {
				fmt.Printf("\n")
				break
			}
			fmt.Printf("Problem reading line: %v\n", err)
			continue
		}
		rInstance.replOneLine(l)
	}
}

// replOneLine handle one line
func (r *repl) replOneLine(line string) {
	if line != "" {
		// mainly for scroll up
		r.liner.AppendHistory(line)
	}
	if r.continueMode {
		r.continueModeByLine(line)
	}
	if !r.continueMode {
		r.run(line)
	}
}

// continueModeByLine check if continue-mode should continue :)
func (r *repl) continueModeByLine(line string) {
	if line != "" {
		r.src += line + "\n"
		return
	}
	// input nothing means jump out continue mode
	r.continueMode = false
	r.prompt = standardPrompt
}

// run execute the input line
func (r *repl) run(newLine string) (err error) {
	src := r.src + newLine + "\n"
	defer func() {
		if errR := recover(); errR != nil {
			replErr(newLine, errR)
			err = errors.New("panic err")
			// TODO: Need a better way to log and show the stack when crash
			// It is too long to print stack on terminal even only print part of the them; not friendly to user
		}
		if err == nil {
			r.src = src
		}
	}()
	fset := token.NewFileSet()
	pkgs, err := parser.Parse(fset, "", src, 0)
	if err != nil {
		// check if into continue mode
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
		// if it is not the first time, restore pre var
		r.preContext.CloneSetVarScope(ctx)
	}
	currentIP := ctx.Exec(r.ip, code.Len())
	r.preContext = *ctx
	// "currentip - 1" is the index of `return`
	// next time it will replace by new code from newLine
	r.ip = currentIP - 1
	return
}

func replErr(line string, err interface{}) {
	fmt.Println("code run fail : ", line)
	fmt.Println("repl err: ", err)
}
