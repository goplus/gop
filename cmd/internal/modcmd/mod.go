package modcmd

import (
	"github.com/goplus/gop/cmd/internal/base"
)

var Cmd = &base.Command{
	UsageLine: "gop mod",
	Short:     "module maintenance",

	Commands: []*base.Command{
		cmdInit,
		cmdDownload,
		cmdTidy,
	},
}
