package os

import (
	"os"
	"strconv"

	"qlang.io/qlang.spec.v1"
)

// -----------------------------------------------------------------------------

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"_name":     "os",
	"args":      os.Args[1:],
	"stdin":     os.Stdin,
	"stderr":    os.Stderr,
	"stdout":    os.Stdout,
	"open":      os.Open,
	"exit":      os.Exit,
	"_initSafe": _initSafe,
}

func _initSafe(mod qlang.Module) {

	mod.Disable("open")
	mod.Exports["exit"] = SafeExit
}

// SafeExit is a safe way to quit qlang application.
//
func SafeExit(code int) {

	panic(strconv.Itoa(code))
}

// -----------------------------------------------------------------------------
