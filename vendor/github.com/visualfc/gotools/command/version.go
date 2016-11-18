// Copyright 2011-2015 visualfc <visualfc@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package command

import (
	"os"
	"runtime"
)

func init() {
	Register(cmdVersion)
}

var AppVersion string = "1.0"

var cmdVersion = &Command{
	Run:       runVersion,
	UsageLine: "version",
	Short:     "print tool version",
	Long:      `Version prints the version.`,
}

func runVersion(cmd *Command, args []string) error {
	if len(args) != 0 {
		cmd.PrintUsage()
		return os.ErrInvalid
	}

	cmd.Printf("%s version %s [%s %s/%s]\n", AppName, AppVersion, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	return nil
}
