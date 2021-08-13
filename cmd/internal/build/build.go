/*
 Copyright 2021 The GoPlus Authors (goplus.org)

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

// Package build implements the ``gop build'' command.
package build

import (
	"github.com/goplus/gop/cmd/internal/base"
)

// -----------------------------------------------------------------------------

// Cmd - gop build
var Cmd = &base.Command{
	UsageLine: "gop build [-v] [-o output] <gopSrcDir|gopSrcFile>",
	Short:     "Build Go+ files",
}

var (
	flagBuildOutput string
	flagVerbose     bool
	flag            = &Cmd.Flag
)

func init() {
	flag.StringVar(&flagBuildOutput, "o", "", "go build output file")
	flag.BoolVar(&flagVerbose, "v", false, "print the names of packages as they are compiled.")
	Cmd.Run = runCmd
}

func runCmd(cmd *base.Command, args []string) {
	panic("TODO: gop build not impl")
}

// -----------------------------------------------------------------------------
