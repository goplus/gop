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
	"fmt"
	"io"

	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/repl"
	"github.com/peterh/liner"
)

// -----------------------------------------------------------------------------

// LinerUI implements repl.UI interface.
type LinerUI struct {
	state  *liner.State
	prompt string
}

// SetPrompt is required by repl.UI interface.
func (u *LinerUI) SetPrompt(prompt string) {
	u.prompt = prompt
}

// Printf is required by repl.UI interface.
func (u *LinerUI) Printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

// -----------------------------------------------------------------------------

func init() {
	Cmd.Run = runCmd
}

// Cmd - gop repl
var Cmd = &base.Command{
	UsageLine: "repl",
	Short:     "Play Go+ in console",
}

const (
	welcome string = "Welcome to Go+ console!"
)

func runCmd(cmd *base.Command, args []string) {
	fmt.Println(welcome)
	state := liner.NewLiner()
	defer state.Close()

	state.SetCtrlCAborts(true)
	state.SetMultiLineMode(true)

	ui := &LinerUI{state: state}
	repl := repl.New()
	repl.SetUI(ui)

	for {
		line, err := ui.state.Prompt(ui.prompt)
		if err != nil {
			if err == liner.ErrPromptAborted || err == io.EOF {
				fmt.Printf("\n")
				break
			}
			fmt.Printf("Problem reading line: %v\n", err)
			continue
		}
		if line != "" {
			state.AppendHistory(line)
		}
		repl.Run(line)
	}
}

// -----------------------------------------------------------------------------
