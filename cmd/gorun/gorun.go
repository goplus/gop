package main

import (
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/goplus/gop/x/gopmod"
)

type goFile struct {
	file string
	code []byte
}

func newGoProj(file string) *gopmod.Project {
	code, err := os.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}
	src := &goFile{file: file, code: code}
	return &gopmod.Project{
		Source:        src,
		AutoGenFile:   file,
		FriendlyFname: filepath.Base(file),
	}
}

func (p *goFile) Fingerp() [20]byte { // source code fingerprint
	return sha1.Sum(p.code)
}

func (p *goFile) IsDirty(outFile string, temp bool) bool {
	return true
}

func (p *goFile) GenGo(outFile string) error {
	if p.file != outFile {
		return os.WriteFile(outFile, p.code, 0666)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "Usage: gorun <file.go>\n\n")
		return
	}
	goProj := newGoProj(os.Args[1])
	ctx := gopmod.New("")
	cmd := ctx.GoCommand("run", goProj, os.Args[2:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	err := cmd.Run()
	if err != nil {
		switch e := err.(type) {
		case *exec.ExitError:
			os.Exit(e.ExitCode())
		default:
			log.Fatalln(err)
		}
	}
}
