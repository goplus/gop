package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/cmd/internal/build"
	"github.com/goplus/gop/cmd/internal/export"
	"github.com/goplus/gop/cmd/internal/gengo"
	"github.com/goplus/gop/cmd/internal/gopfmt"
	"github.com/goplus/gop/cmd/internal/help"
	"github.com/goplus/gop/cmd/internal/install"
	"github.com/goplus/gop/cmd/internal/repl"
	"github.com/goplus/gop/cmd/internal/run"

	_ "github.com/goplus/gop/lib"
)

func mainUsage() {
	help.PrintUsage(os.Stderr, base.Gop)
	os.Exit(2)
}

func init() {
	base.Usage = mainUsage
	base.Gop.Commands = []*base.Command{
		run.Cmd,
		gengo.Cmd,
		gopfmt.Cmd,
		export.Cmd,
		repl.Cmd,
		install.Cmd,
		build.Cmd,
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		base.Usage()
	}

	base.CmdName = args[0] // for error messages
	if args[0] == "help" {
		help.Help(os.Stderr, args[1:])
		return
	}

BigCmdLoop:
	for bigCmd := base.Gop; ; {
		for _, cmd := range bigCmd.Commands {
			if cmd.Name() != args[0] {
				continue
			}
			args = args[1:]
			if len(cmd.Commands) > 0 {
				bigCmd = cmd
				if len(args) == 0 {
					help.PrintUsage(os.Stderr, bigCmd)
					os.Exit(2)
				}
				if args[0] == "help" {
					help.Help(os.Stderr, append(strings.Split(base.CmdName, " "), args[1:]...))
					return
				}
				base.CmdName += " " + args[0]
				continue BigCmdLoop
			}
			if !cmd.Runnable() {
				continue
			}
			cmd.Run(cmd, args)
			return
		}
		helpArg := ""
		if i := strings.LastIndex(base.CmdName, " "); i >= 0 {
			helpArg = " " + base.CmdName[:i]
		}
		fmt.Fprintf(os.Stderr, "gop %s: unknown command\nRun 'gop help%s' for usage.\n", base.CmdName, helpArg)
		os.Exit(2)
	}
}
