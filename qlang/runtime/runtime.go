package runtime

import (
	"runtime"
)

// -----------------------------------------------------------------------------

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"GOMAXPROCS": runtime.GOMAXPROCS,
	"GOROOT":     runtime.GOROOT,
}

// -----------------------------------------------------------------------------
