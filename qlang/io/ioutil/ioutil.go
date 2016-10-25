package ioutil

import (
	"io/ioutil"

	qlang "qlang.io/spec"
)

// -----------------------------------------------------------------------------

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"_name":     "io/ioutil",
	"_initSafe": _initSafe,
	"discard":   ioutil.Discard,

	"nopCloser": ioutil.NopCloser,
	"readAll":   ioutil.ReadAll,
	"readDir":   ioutil.ReadDir,
	"readFile":  ioutil.ReadFile,
	"tempDir":   ioutil.TempDir,
	"tempFile":  ioutil.TempFile,
	"writeFile": ioutil.WriteFile,
}

func _initSafe(mod qlang.Module) {

	mod.Disable("readDir", "readFile", "writeFile")
}

// -----------------------------------------------------------------------------
