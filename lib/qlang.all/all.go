package qall

import (
	"lib/bufio"
	"lib/bytes"
	"lib/crypto/md5"
	"lib/encoding/hex"
	"lib/encoding/json"
	"lib/eqlang"
	"lib/errors"
	"lib/io"
	"lib/io/ioutil"
	"lib/math"
	"lib/meta"
	"lib/net/http"
	"lib/os"
	"lib/path"
	"lib/reflect"
	"lib/runtime"
	"lib/strconv"
	"lib/strings"
	"lib/sync"
	"lib/terminal"
	"lib/tpl/extractor"
	"lib/version"
	qlang "spec"

	// qlang builtin modules
	_ "lib/builtin"
	_ "lib/chan"
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
