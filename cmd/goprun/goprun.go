package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/goplus/gop/x/gopmod"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "Usage: gorun <file.go> [switch ...]\n\n")
		return
	}
	ctx := gopmod.New("")
	goProj, err := ctx.OpenProject(0, os.Args[1])
	if err != nil {
		fmt.Fprint(os.Stderr, "OpenProject failed:", err)
		return
	}
	goProj.ExecArgs = os.Args[2:]
	cmd := ctx.GoCommand("run", goProj)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	err = cmd.Run()
	if err != nil {
		switch e := err.(type) {
		case *exec.ExitError:
			os.Exit(e.ExitCode())
		default:
			log.Fatalln(err)
		}
	}
}
