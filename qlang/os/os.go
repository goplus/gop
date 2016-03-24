package os

import (
	"os"
	"strconv"
)

// -----------------------------------------------------------------------------

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

func SafeExit(code int) {

	panic(strconv.Itoa(code))
}

// -----------------------------------------------------------------------------
