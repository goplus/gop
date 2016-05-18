package os

import (
	"os"
	"strconv"
)

// -----------------------------------------------------------------------------

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"args":      os.Args[1:],
	"stdin":     os.Stdin,
	"stderr":    os.Stderr,
	"stdout":    os.Stdout,
	"open":      os.Open,
	"exit":      os.Exit,
	"_initSafe": _initSafe,
}

func _initSafe(table map[string]interface{}, dummy func(...interface{}) interface{}) {

	table["open"] = dummy
	table["exit"] = SafeExit
}

// SafeExit is a safe way to quit qlang application.
//
func SafeExit(code int) {

	panic(strconv.Itoa(code))
}

// -----------------------------------------------------------------------------
