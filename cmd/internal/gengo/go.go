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

// Package gengo implements the ``gop go'' command.
package gengo

import (
	"fmt"
	"os"
	"strings"

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/cmd/gengo"
	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gox"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------
/*
var (
	errTestFailed = errors.New("test failed")
)

func testPkg(p *gengo.Runner, dir string, flags int) error {
	if flags == gengo.PkgFlagGo { // don't test Go packages
		return nil
	}
	cmd1 := exec.Command("go", "run", path.Join(dir, "gop_autogen.go"))
	gorun, err := cmd1.CombinedOutput()
	if err != nil {
		os.Stderr.Write(gorun)
		fmt.Fprintf(os.Stderr, "[ERROR] `%v` failed: %v\n", cmd1, err)
		return err
	}
	cmd2 := exec.Command("gop", "run", "-quiet", dir) // -quiet: don't generate any log
	qrun, err := cmd2.CombinedOutput()
	if err != nil {
		os.Stderr.Write(qrun)
		fmt.Fprintf(os.Stderr, "[ERROR] `%v` failed: %v\n", cmd2, err)
		return err
	}
	if !bytes.Equal(gorun, qrun) {
		fmt.Fprintf(os.Stderr, "[ERROR] Output has differences!\n")
		fmt.Fprintf(os.Stderr, ">>> Output of `%v`:\n", cmd1)
		os.Stderr.Write(gorun)
		fmt.Fprintf(os.Stderr, "\n>>> Output of `%v`:\n", cmd2)
		os.Stderr.Write(qrun)
		return errTestFailed
	}
	return nil
}
*/
// -----------------------------------------------------------------------------

// Cmd - gop go
var Cmd = &base.Command{
	UsageLine: "gop go [-debug -test -slow] <gopSrcDir>",
	Short:     "Convert Go+ packages into Go packages",
}

var (
	flag      = &Cmd.Flag
	flagDebug = flag.Bool("debug", false, "set log level to debug")
	flagTest  = flag.Bool("test", false, "test Go+ package")
	flagSlow  = flag.Bool("slow", false, "don't cache imported packages")
)

func init() {
	Cmd.Run = runCmd
}

func runCmd(cmd *base.Command, args []string) {
	err := flag.Parse(args)
	if err != nil {
		log.Fatalln("parse input arguments failed:", err)
	}
	if flag.NArg() < 1 {
		cmd.Usage(os.Stderr)
		return
	}
	if *flagDebug {
		log.SetOutputLevel(log.Ldebug)
		gox.SetDebug(gox.DbgFlagAll)
	}
	dir := flag.Arg(0)
	dir = strings.TrimSuffix(dir, "/...")
	runner := new(gengo.Runner)
	runner.SetAfter(func(p *gengo.Runner, dir string, flags int) error {
		errs := p.ResetErrors()
		if errs != nil {
			for _, err := range errs {
				fmt.Fprintln(os.Stderr, err)
			}
			fmt.Fprintln(os.Stderr)
		} else if *flagTest {
			panic("gop go -test: not impl")
		}
		return nil
	})
	runner.GenGo(dir, true, &cl.Config{CacheLoadPkgs: !*flagSlow})
	errs := runner.Errors()
	if errs != nil {
		for _, err := range errs {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(-1)
	}
}

// -----------------------------------------------------------------------------
