package run

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/goplus/gop/env"
	"github.com/goplus/gop/x/gopmod"
	"github.com/qiniu/x/log"
)

func findGoModFile(dir string) (modfile string, noCacheFile bool, err error) {
	modfile, err = env.GOPMOD(dir)
	if err != nil {
		modfile = filepath.Join(env.GOPROOT(), "go.mod")
		return modfile, true, nil
	}
	return
}

func findGoModDir(dir string) (string, bool) {
	modfile, nocachefile, err := findGoModFile(dir)
	if err != nil {
		log.Fatalln("findGoModFile:", err)
	}
	return filepath.Dir(modfile), nocachefile
}

func gopRun(source string, args ...string) {
	ctx := gopmod.New("")
	flags := 0
	if *flagGop {
		flags = gopmod.FlagGoAsGoPlus
	}
	goProj, err := ctx.OpenProject(flags, source)
	if err != nil {
		log.Fatalln("OpenProject failed:", err)
	}
	goProj.ExecArgs = args
	goProj.FlagNRINC = *flagNorun
	goProj.FlagRTOE = *flagRTOE
	if goProj.FlagRTOE {
		goProj.UseDefaultCtx = true
	}
	cmd := ctx.GoCommand("run", goProj)
	if cmd.IsValid() {
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
}
