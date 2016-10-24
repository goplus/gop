package qall

import (
	"qlang.io/qlang/bufio"
	"qlang.io/qlang/bytes"
	"qlang.io/qlang/crypto/md5"
	"qlang.io/qlang/encoding/hex"
	"qlang.io/qlang/encoding/json"
	"qlang.io/qlang/eqlang"
	"qlang.io/qlang/errors"
	"qlang.io/qlang/io"
	"qlang.io/qlang/io/ioutil"
	"qlang.io/qlang/math"
	"qlang.io/qlang/meta"
	"qlang.io/qlang/net/http"
	"qlang.io/qlang/os"
	"qlang.io/qlang/path"
	"qlang.io/qlang/reflect"
	"qlang.io/qlang/runtime"
	"qlang.io/qlang/strconv"
	"qlang.io/qlang/strings"
	"qlang.io/qlang/sync"
	"qlang.io/qlang/terminal"
	"qlang.io/qlang/tpl/extractor"
	"qlang.io/qlang/version"
	qlang "qlang.io/spec"

	// qlang builtin modules
	_ "qlang.io/qlang/builtin"
	_ "qlang.io/qlang/chan"
)

// -----------------------------------------------------------------------------

// Copyright prints qlang copyright information.
//
func Copyright() {
	version.Copyright()
}

// InitSafe inits qlang and imports modules.
//
func InitSafe(safeMode bool) {

	qlang.SafeMode = safeMode

	qlang.Import("", math.Exports) // import math as builtin package
	qlang.Import("", meta.Exports) // import meta package
	qlang.Import("bufio", bufio.Exports)
	qlang.Import("bytes", bytes.Exports)
	qlang.Import("md5", md5.Exports)
	qlang.Import("io", io.Exports)
	qlang.Import("ioutil", ioutil.Exports)
	qlang.Import("hex", hex.Exports)
	qlang.Import("json", json.Exports)
	qlang.Import("errors", errors.Exports)
	qlang.Import("eqlang", eqlang.Exports)
	qlang.Import("math", math.Exports)
	qlang.Import("os", os.Exports)
	qlang.Import("", os.InlineExports)
	qlang.Import("path", path.Exports)
	qlang.Import("http", http.Exports)
	qlang.Import("reflect", reflect.Exports)
	qlang.Import("runtime", runtime.Exports)
	qlang.Import("strconv", strconv.Exports)
	qlang.Import("strings", strings.Exports)
	qlang.Import("sync", sync.Exports)
	qlang.Import("terminal", terminal.Exports)
	qlang.Import("extractor", extractor.Exports)
}

// -----------------------------------------------------------------------------
