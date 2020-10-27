package work

import (
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// NewBuild creates a Build struct which can build from gop temporary directory,
// and generate binary in current working directory
func NewBuild(buildflags string, args []string, workingDir string, outputDir string) (*Build, error) {
	b := &Build{
		BuildFlags: buildflags,
		WorkingDir: workingDir,
	}
	if err := b.CreateTmpWorkingDir(); err != nil {
		return nil, err
	}
	b.OriGOPATH = os.Getenv("GOPATH")
	if b.OriGOPATH == "" {
		b.NewGOPATH = b.TmpDir
	} else {
		b.NewGOPATH = fmt.Sprintf("%v:%v", b.TmpDir, b.OriGOPATH)
	}
	var packages []string
	for _, arg := range args {
		dir, _ := filepath.Split(arg)
		fset := token.NewFileSet()
		pkgs, _ := parser.ParseDir(fset, dir, nil, 0)
		if _, ok := pkgs["main"]; ok {
			b.Target = filepath.Join(b.WorkingDir, dir)
		}
		packages = append(packages, filepath.Join(dir, "gop_autogen.go"))
		err := GenGo(dir, filepath.Join(b.TmpWorkingDir, arg))
		if err != nil {
			return nil, err
		}
	}
	b.Packages = strings.Join(packages, " ")
	return b, nil
}

// Build calls 'go build' tool to do building
func (b *Build) Build() error {
	log.Println("Go building in temp...")
	// new -o will overwrite  previous ones
	b.BuildFlags = b.BuildFlags + " -o " + b.Target
	cmd := exec.Command("/bin/bash", "-c", "go build "+b.BuildFlags+" "+b.Packages)
	cmd.Dir = b.TmpWorkingDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if b.NewGOPATH != "" {
		// Change to temp GOPATH for go install command
		cmd.Env = append(os.Environ(), fmt.Sprintf("GOPATH=%v", b.NewGOPATH))
	}

	log.Printf("go build cmd is: %v", cmd.Args)
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("fail to execute: %v, err: %w", cmd.Args, err)
	}
	if err = cmd.Wait(); err != nil {
		return fmt.Errorf("fail to execute: %v, err: %w", cmd.Args, err)
	}
	log.Println("Go build exit successful.")
	return nil
}
