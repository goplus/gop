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

package repl

import (
	"errors"
	"fmt"
	"strings"

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"

	exec "github.com/goplus/gop/exec/bytecode"
)

// UI represents the UserInterface interacting abstraction.
type UI interface {
	SetPrompt(prompt string)
	Printf(format string, a ...interface{})
}

// REPL type
type REPL struct {
	src          string       // the whole source code from repl
	preContext   exec.Context // store the context after exec
	ip           int          // store the ip after exec
	continueMode bool         // switch to control the promot type
	term         UI           // liner instance
}

const (
	// ContinuePrompt - the current code statement is not completed.
	ContinuePrompt string = "... "
	// NormalPrompt - start of a code statement.
	NormalPrompt string = ">>> "
)

// New creates a REPL object.
func New() *REPL {
	return &REPL{}
}

// SetUI initializes UI.
func (r *REPL) SetUI(term UI) {
	r.term = term
	term.SetPrompt(NormalPrompt)
}

// Run handles one line.
func (r *REPL) Run(line string) {
	if r.continueMode {
		r.continueModeByLine(line)
	}
	if !r.continueMode {
		r.run(line)
	}
}

// continueModeByLine check if continue-mode should continue :)
func (r *REPL) continueModeByLine(line string) {
	if line != "" {
		r.src += line + "\n"
		return
	}
	// input nothing means jump out continue mode
	r.continueMode = false
	r.term.SetPrompt(NormalPrompt)
}

// run execute the input line
func (r *REPL) run(newLine string) (err error) {
	src := r.src + newLine + "\n"
	defer func() {
		if errR := recover(); errR != nil {
			r.dumpErr(newLine, errR)
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
			r.term.SetPrompt(ContinuePrompt)
			r.continueMode = true
			err = nil
			return
		}
		r.term.Printf("ParseGopFiles err: %v\n", err)
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
		r.term.Printf("NewPackage err %+v\n", err)
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
	size := ctx.Len()
	var dump []string
	for i := 0; i < size; i++ {
		dump = append(dump, fmt.Sprintf("%v", ctx.Get(i-size)))
	}
	if len(dump) > 0 {
		r.term.Printf("%v\n", strings.Join(dump, ","))
	}
	return
}

func (r *REPL) dumpErr(line string, err interface{}) {
	r.term.Printf("code run fail : %v\n", line)
	r.term.Printf("repl err: %v\n", err)
}
