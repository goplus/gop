package main

import (
	"os"

	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/cmd/internal/repl"

	_ "github.com/goplus/gop/lib"
)

func main() {
	base.Main(repl.Cmd, "igop", os.Args[1:])
}
