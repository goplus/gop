package qlang_all

import (
	"qlang.io/qlang/bufio"
	"qlang.io/qlang/crypto/md5"
	"qlang.io/qlang/encoding/json"
	"qlang.io/qlang/io/ioutil"
	"qlang.io/qlang/math"
	"qlang.io/qlang/os"
	"qlang.io/qlang/path"
	"qlang.io/qlang/reflect"
	"qlang.io/qlang/strconv"
	"qlang.io/qlang/strings"
	"qlang.io/qlang/tpl.v1/extractor"
	"qlang.io/qlang/version"
	"qlang.io/qlang.spec.v1"

	_ "qlang.io/qlang/builtin"
)

// -----------------------------------------------------------------------------

func Copyright() {
	version.Copyright()
}

func InitSafe(safeMode bool) {

	qlang.SafeMode = safeMode

	qlang.Import("", math.Exports) // import math as builtin package
	qlang.Import("bufio", bufio.Exports)
	qlang.Import("md5", md5.Exports)
	qlang.Import("ioutil", ioutil.Exports)
	qlang.Import("json", json.Exports)
	qlang.Import("math", math.Exports)
	qlang.Import("os", os.Exports)
	qlang.Import("path", path.Exports)
	qlang.Import("reflect", reflect.Exports)
	qlang.Import("strconv", strconv.Exports)
	qlang.Import("strings", strings.Exports)
	qlang.Import("extractor", extractor.Exports)
}

// -----------------------------------------------------------------------------
