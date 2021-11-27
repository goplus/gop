// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package modcmd

import (
	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/cmd/internal/modfetch"
	"github.com/goplus/gop/cmd/internal/modload"
)

var cmdDownload = &base.Command{
	UsageLine: "gop mod download [-x] [-json] [modules]",
	Short:     "download modules to local cache",
}

func init() {
	cmdDownload.Run = runDownload // break init cycle
}

func runDownload(cmd *base.Command, args []string) {
	modfetch.DownloadArgs(".", args...)
	modload.SyncGopMod()
}
