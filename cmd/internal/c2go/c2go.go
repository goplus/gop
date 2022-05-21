package c2go

import (
	"os"

	"github.com/goplus/gop/cmd/internal/base"

	c2go "github.com/goplus/c2go/cmd/c2go/impl"
)

// gop c2go
var Cmd = &base.Command{
	UsageLine: "gop " + c2go.ShortUsage,
	Short:     "Run c2go (convert C to Go) tools",
}

func init() {
	Cmd.Flag.Usage = func() {
		Cmd.Usage(os.Stderr)
	}
	Cmd.Run = runCmd
}

func runCmd(cmd *base.Command, args []string) {
	c2go.Main(&Cmd.Flag, args)
}
