// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package modcmd

import (
	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/cmd/internal/modfetch"
	"github.com/goplus/gop/cmd/internal/modload"
)

var cmdTidy = &base.Command{
	UsageLine: "gop mod tidy [-e] [-v]",
	Short:     "add missing and remove unused modules",
}

func init() {
	cmdTidy.Run = runTidy
}

func runTidy(cmd *base.Command, args []string) {
	modload.LoadModFile()
	modload.SyncGoMod()
	modfetch.TidyArgs(".", args...)
	modload.SyncGopMod()
}
