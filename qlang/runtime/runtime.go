package runtime

import (
	"runtime"
)

// -----------------------------------------------------------------------------

var Exports = map[string]interface{}{
	"GOMAXPROCS": runtime.GOMAXPROCS,
	"GOROOT":     runtime.GOROOT,
}

// -----------------------------------------------------------------------------

