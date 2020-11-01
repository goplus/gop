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
package work

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/goplus/gop/ast"
)

// Build is to describe the building/installing process of a gop build/install
type Build struct {
	Pkgs          ast.Package // Pkg list parsed from "go list -json ./..." command
	NewGOPATH     string      // the new GOPATH
	OriGOPATH     string      // the original GOPATH
	WorkingDir    string      // the working directory
	TmpDir        string      // the temporary directory to build the project
	TmpWorkingDir string      // the working directory in the temporary directory, which is corresponding to the current directory in the project directory
	Target        string      // the binary name that go build generate
	// keep compatible with go commands:
	// gop build [-o output] [-i] [build flags] [packages]
	// gop install [-i] [build flags] [packages]
	BuildFlags string // Build flags
	Packages   string // Packages that needs to build
}

// NewInstall creates a Build struct which can install from goc temporary directory
func NewInstall(buildflags string, args []string, workingDir string) (*Build, error) {
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
		packages = append(packages, filepath.Join(dir, "gop_autogen.go"))
		err := GenGo(dir, filepath.Join(b.TmpWorkingDir, arg))
		if err != nil {
			return nil, err
		}
	}
	b.Packages = strings.Join(packages, " ")

	return b, nil
}

// CreateTmpWorkingDir moves the projects into a temporary directory
func (b *Build) CreateTmpWorkingDir() error {
	var err error
	b.TmpDir = filepath.Join(os.TempDir(), tmpFolderName(b.WorkingDir))

	// Create a new tmp folder
	err = os.MkdirAll(filepath.Join(b.TmpDir, "src"), os.ModePerm)
	if err != nil {
		return fmt.Errorf("Fail to create the temporary build directory. The err is: %v", err)
	}
	log.Printf("Tmp project generated in: %v", b.TmpDir)
	if err != nil {
		log.Fatalln("Fail to move the project to temporary directory")
		return err
	}
	_, lastdir := filepath.Split(b.WorkingDir)
	b.TmpWorkingDir = filepath.Join(b.TmpDir, "src", lastdir)
	return nil
}

// Install use the 'go install' tool to install packages
func (b *Build) Install() error {
	log.Println("Go building in temp...")
	cmd := exec.Command("/bin/bash", "-c", "go install "+b.BuildFlags+" "+b.Packages)
	cmd.Dir = b.TmpWorkingDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	whereToInstall, err := b.findWhereToInstall()
	if err != nil {
		// ignore the err
		log.Fatalf("No place to install: %v", err)
	}
	// Change the temp GOBIN, to force binary install to original place
	cmd.Env = append(os.Environ(), fmt.Sprintf("GOBIN=%v", whereToInstall))
	if b.NewGOPATH != "" {
		// Change to temp GOPATH for go install command
		cmd.Env = append(cmd.Env, fmt.Sprintf("GOPATH=%v", b.NewGOPATH))
	}

	log.Printf("go install cmd is: %v", cmd.Args)
	err = cmd.Start()
	if err != nil {
		log.Fatalf("Fail to execute: %v. The error is: %v", cmd.Args, err)
		return err
	}
	if err = cmd.Wait(); err != nil {
		log.Fatalf("go install failed. The error is: %v", err)
		return err
	}
	log.Printf("Go install successful. Binary installed in: %v", whereToInstall)
	return nil
}

func (b *Build) findWhereToInstall() (string, error) {
	if GOPBIN := os.Getenv("GOPBIN"); GOPBIN != "" {
		return GOPBIN, nil
	}
	return filepath.Join(os.Getenv("HOME"), "gop", "bin"), nil
}

// tmpFolderName uses the first six characters of the input path's SHA256 checksum
// as the suffix.
func tmpFolderName(path string) string {
	sum := sha256.Sum256([]byte(path))
	h := fmt.Sprintf("%x", sum[:6])

	return "goplus-build-" + h
}
