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
	"strings"
	"time"
)

func getcwd() string {
	path, _ := os.Getwd()
	return path
}

var gopRoot = getcwd()

func checkPathExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func getGopHomeDir() string {
	path, _ := os.UserHomeDir()
	return path + "/gop"
}

func execCommand(command string, arg ...string) (string, string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(command, arg...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func getGitInfo() (string, string) {
	gitDir := gopRoot + "/.git"
	noBranch := "nobranch"
	noCommit := "nocommit"

	if checkPathExist(gitDir) {
		branch, stderr, _ := execCommand("git", "branch", "--show-current")
		if len(stderr) > 0 {
			println(stderr)
			branch = noBranch
		} else {
			branch = strings.TrimRight(branch, "\n")
		}

		commit, stderr, _ := execCommand("git", "rev-parse", "--verify", "HEAD")
		if len(stderr) > 0 {
			println(stderr)
			commit = noCommit
		} else {
			commit = strings.TrimRight(commit, "\n")
		}

		return branch, commit
	}

	return noBranch, noCommit
}

func getBuildDateTime() string {
	now := time.Now()
	return now.Format("2006-01-02_15-04-05")
}

func getGopBuildFlags() string {
	branch, commit := getGitInfo()
	buildDateTime := getBuildDateTime()

	buildFlags := fmt.Sprintf("-X github.com/goplus/gop/build.Date=%s ", buildDateTime)
	buildFlags += fmt.Sprintf("-X github.com/goplus/gop/build.Commit=%s ", commit)
	buildFlags += fmt.Sprintf("-X github.com/goplus/gop/build.Branch=%s", branch)

	return buildFlags
}

func buildGoplusTools() {
	commandsDir := gopRoot + "/cmd"
	if !checkPathExist(commandsDir) {
		println("Error: This script should be run at the root directory of gop repository.")
		os.Exit(1)
	}

	buildFlags := getGopBuildFlags()

	println("Installing Go+ tools...")
	os.Chdir(commandsDir)
	buildOutput, buildErr, err := execCommand("go", "install", "-v", "-ldflags", buildFlags, "./...")
	if err != nil {
		println(err)
		os.Exit(1)
	}
	println(buildErr)
	println(buildOutput)
	println("Go+ tools installed successfully!")
}

func linkGoplusToHome() {
	gopHomeDir := getGopHomeDir()

	fmt.Printf("Linking %s to %s\n", gopRoot, gopHomeDir)

	os.Chdir(gopRoot)
	if gopHomeDir != gopRoot && !checkPathExist(gopHomeDir) {
		err := os.Symlink(gopRoot, gopHomeDir)
		if err != nil {
			println(err)
		}
	}

	fmt.Printf("%s linked to %s successfully!\n", gopRoot, gopHomeDir)
}

func runTestcases() {
	println("Start running testcases.")
	os.Chdir(gopRoot)

	path, _ := os.LookupEnv("PATH")
	homeDir, _ := os.UserHomeDir()

	goPath, ok := os.LookupEnv("GOPATH")
	if ok {
		path += fmt.Sprintf(":%s/bin", goPath)
	} else {
		path += fmt.Sprintf(":%s/go/bin", homeDir)
	}

	cmd := exec.Command("gop", "test", "-v", "-coverprofile=coverage.txt", "-covermode=atomic", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "PATH="+path)
	err := cmd.Run()
	if err != nil {
		println(err)
	}
	println("End running testcases.")
}

func localInstall() {
	buildGoplusTools()
	linkGoplusToHome()
	println("Go+ is now installed.")
}

func main() {
	isBuild := flag.Bool("build", false, "Build the Go+")
	isTest := flag.Bool("test", false, "Run testcases")

	flag.Parse()

	if !*isBuild && !*isTest {
		localInstall()
		return
	}

	if *isBuild {
		buildGoplusTools()
	}

	if *isTest {
		runTestcases()
	}
}
