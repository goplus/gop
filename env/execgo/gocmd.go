package execgo

import (
	"fmt"
	"os/exec"

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

func Command(op string, args ...string) *exec.Cmd {
	exargs := make([]string, 1, len(args)+3)
	exargs[0] = op
	exargs = appendLdflags(exargs, op)
	exargs = append(exargs, args...)
	return exec.Command("go", exargs...)
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
