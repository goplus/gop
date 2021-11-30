package gopmod

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/goplus/gop/env"
)

// -----------------------------------------------------------------------------

const (
	ldFlagVersion   = "-X \"github.com/goplus/gop/env.buildVersion=%s\""
	ldFlagBuildDate = "-X \"github.com/goplus/gop/env.buildDate=%s\""
	ldFlagBuildRev  = "-X \"github.com/goplus/gop/env.buildCommit=%s\""
	ldFlagGopRoot   = "-X \"github.com/goplus/gop/env.defaultGopRoot=%s\""
)

const (
	ldFlagAll = ldFlagVersion + " " + ldFlagBuildDate + " " + ldFlagBuildRev + " " + ldFlagGopRoot
)

var (
	GOPVERSION   = env.Version()
	GOPBUILDDATE = env.BuildDate()
	GOPBUILDREV  = env.BuildRevision()
	GOPROOT      = env.GOPROOT()
)

func LoadFlags() string {
	return fmt.Sprintf(ldFlagAll, GOPVERSION, GOPBUILDDATE, GOPBUILDREV, GOPROOT)
}

// -----------------------------------------------------------------------------

type GoCmd struct {
	*exec.Cmd
	after func() error
}

func (p GoCmd) Run() error {
	if err := p.Cmd.Run(); err != nil {
		return err
	}
	return p.after()
}

func goCommand(dir, op string, t *goTarget, args ...string) (ret GoCmd) {
	exargs := make([]string, 1, len(args)+6)
	exargs[0] = op                     // 1
	exargs = appendLdflags(exargs, op) // 2
	if op == "run" && t.defctx {       // 2
		afterDir := dir
		dir, _ = filepath.Split(t.goFile)
		exargs[0] = "build"
		exargs = append(exargs, "-o", t.outFile)
		ret.after = func() error {
			return runCommand(afterDir, t.outFile, args...)
		}
	} else {
		exargs = append(exargs, args...) // len(args)
	}
	exargs = append(exargs, t.goFile) // 1
	ret.Cmd = exec.Command("go", exargs...)
	ret.Cmd.Dir = dir
	return
}

func appendLdflags(exargs []string, op string) []string {
	for _, v := range opsWithLdflags {
		if op == v {
			return append(exargs, "-ldflags", LoadFlags())
		}
	}
	return exargs
}

var (
	opsWithLdflags = []string{"run", "install", "build", "test"}
)

// -----------------------------------------------------------------------------
