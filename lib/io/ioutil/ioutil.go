package ioutil

import (
	"io/ioutil"
)

// -----------------------------------------------------------------------------

var Exports = map[string]interface{}{
	"nopCloser": ioutil.NopCloser,
	"readAll":   ioutil.ReadAll,
	"readDir":   ioutil.ReadDir,
	"readFile":  ioutil.ReadFile,
	"writeFile": ioutil.WriteFile,
	"discard":   ioutil.Discard,
	"_initSafe": _initSafe,
}

func _initSafe(table map[string]interface{}, dummy func(...interface{}) interface{}) {

	table["readDir"] = dummy
	table["readFile"] = dummy
	table["writeFile"] = dummy
}

// -----------------------------------------------------------------------------
