//go:build ignore
// +build ignore

/*
 * Copyright (c) 2021 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func checkPathExist(path string, isDir bool) bool {
	stat, err := os.Stat(path)
	isExists := !os.IsNotExist(err)
	if isDir {
		return isExists && stat.IsDir()
	}
	return isExists
}

// Path returns single path to check
type Path struct {
	path  string
	isDir bool
}

func (p *Path) checkExists(rootDir string) bool {
	absPath := filepath.Join(rootDir, p.path)
	return checkPathExist(absPath, p.isDir)
}

func getGopRoot() string {
	pwd, _ := os.Getwd()

	pathsToCheck := []Path{
		{path: "cmd/gop", isDir: true},
		{path: "builtin", isDir: true},
		{path: "go.mod", isDir: false},
		{path: "go.sum", isDir: false},
	}

	for _, path := range pathsToCheck {
		if !path.checkExists(pwd) {
			println("Error: This script should be run at the root directory of gop repository.")
			os.Exit(1)
		}
	}
	return pwd
}

var gopRoot = getGopRoot()
var initCommandExecuteEnv = os.Environ()
var commandExecuteEnv = initCommandExecuteEnv

func execCommand(command string, arg ...string) (string, string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(command, arg...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Env = commandExecuteEnv
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func getRevCommit(tag string) string {
	commit, stderr, err := execCommand("git", "rev-parse", "--verify", tag)
	if err != nil || stderr != "" {
		return ""
	}
	return strings.TrimRight(commit, "\n")
}

func getGitInfo() (string, bool) {
	gitDir := filepath.Join(gopRoot, ".git")
	if checkPathExist(gitDir, true) {
		return getRevCommit("HEAD"), true
	}
	return "", false
}

func getBuildDateTime() string {
	now := time.Now()
	return now.Format("2006-01-02_15-04-05")
}

func getBuildVer() string {
	tagRet, tagErr, err := execCommand("git", "describe", "--tags")
	if err != nil || tagErr != "" {
		return ""
	}
	return strings.TrimRight(tagRet, "\n")
}

func getGopBuildFlags() string {
	defaultGopRoot := gopRoot
	if gopRootFinal := os.Getenv("GOPROOT_FINAL"); gopRootFinal != "" {
		defaultGopRoot = gopRootFinal
	}
	buildFlags := fmt.Sprintf("-X \"github.com/goplus/gop/env.defaultGopRoot=%s\"", defaultGopRoot)
	buildFlags += fmt.Sprintf(" -X \"github.com/goplus/gop/env.buildDate=%s\"", getBuildDateTime())
	if commit, ok := getGitInfo(); ok {
		buildFlags += fmt.Sprintf(" -X github.com/goplus/gop/env.buildCommit=%s", commit)
		if buildVer := getBuildVer(); buildVer != "" {
			buildFlags += fmt.Sprintf(" -X github.com/goplus/gop/env.buildVersion=%s", buildVer)
		}
	}
	return buildFlags
}

func detectGopBinPath() string {
	return filepath.Join(gopRoot, "bin")
}

func buildGoplusTools(useGoProxy bool) {
	commandsDir := filepath.Join(gopRoot, "cmd")
	buildFlags := getGopBuildFlags()

	if useGoProxy {
		println("Info: we will use goproxy.cn as a Go proxy to accelerate installing process.")
		commandExecuteEnv = append(commandExecuteEnv,
			"GOPROXY=https://goproxy.cn,direct",
		)
	}

	// Install Go+ binary files under current ./bin directory.
	gopBinPath := detectGopBinPath()
	clean()
	if err := os.Mkdir(gopBinPath, 0755); err != nil {
		println(err.Error())
		println("Error: Go+ can't create ./bin directory to put build assets.")
		os.Exit(1)
	}

	println("Installing Go+ tools...")
	os.Chdir(commandsDir)
	buildOutput, buildErr, err := execCommand("go", "build", "-o", gopBinPath, "-v", "-ldflags", buildFlags, "./...")
	println(buildErr)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	println(buildOutput)

	println("Go+ tools installed successfully!")
	showHelpPostInstall()
}

func showHelpPostInstall() {
	println("\nNEXT STEP:")
	println("\nWe just installed Go+ into the directory: ", detectGopBinPath())
	message := `
To setup a better Go+ development environment,
we recommend you add the above install directory into your PATH environment variable.
	`
	println(message)
}

func runTestcases() {
	println("Start running testcases.")
	os.Chdir(gopRoot)

	coverage := "-coverprofile=coverage.txt"
	gopCommand := filepath.Join(detectGopBinPath(), "gop")
	if !checkPathExist(gopCommand, false) {
		println("Error: Go+ must be installed before running testcases.")
		os.Exit(1)
	}

	testOutput, testErr, err := execCommand(gopCommand, "test", coverage, "-covermode=atomic", "./...")
	println(testOutput)
	println(testErr)
	if err != nil {
		println(err.Error())
	}

	println("End running testcases.")
}

func clean() {
	gopBinPath := detectGopBinPath()
	if checkPathExist(gopBinPath, true) {
		if err := os.RemoveAll(gopBinPath); err != nil {
			println(err.Error())
		}
	}
}

func uninstall() {
	println("Uninstalling Go+ and related tools.")
	clean()
	println("Go+ and related tools uninstalled successfully.")
}

func isInChina() bool {
	const prefix = "LANG=\""
	out, errMsg, err := execCommand("locale")
	if err != nil || errMsg != "" {
		return false
	}
	if strings.HasPrefix(out, prefix) {
		out = out[len(prefix):]
		return strings.HasPrefix(out, "zh_CN") || strings.HasPrefix(out, "zh_HK")
	}
	return false
}

func main() {
	isInstall := flag.Bool("install", false, "Install Go+")
	isTest := flag.Bool("test", false, "Run testcases")
	isUninstall := flag.Bool("uninstall", false, "Uninstall Go+")
	isGoProxy := flag.Bool("proxy", false, "Set GOPROXY for people in China")
	isAutoProxy := flag.Bool("autoproxy", false, "Check to set GOPROXY automatically")

	flag.Parse()

	useGoProxy := *isGoProxy
	if !useGoProxy && *isAutoProxy {
		useGoProxy = isInChina()
	}
	flagActionMap := map[*bool]func(){
		isInstall:   func() { buildGoplusTools(useGoProxy) },
		isUninstall: uninstall,
		isTest:      runTestcases,
	}

	// Sort flags, for example: install flag should be checked earlier than test flag.
	flags := []*bool{isInstall, isTest, isUninstall}
	hasActionDone := false

	for _, flag := range flags {
		if *flag {
			flagActionMap[flag]()
			hasActionDone = true
		}
	}

	if !hasActionDone {
		println("Usage:\n")
		flag.PrintDefaults()
	}
}
