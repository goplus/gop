package build

import (
	"fmt"
)

// The value of variables come form `go build -ldflags '-X "build.Date=xxxxx" -X "build.CommitID=xxxx"' `
var (
	// Date build time
	Date string
	// Branch current git branch
	Branch string
	// Commit git commit id
	Commit string
)

// Build information
func Build() string {
	//if Date == "" {
	//	return ""
	//}
	return fmt.Sprintf("%s(%s) %s", Branch, Commit, Date)
}
