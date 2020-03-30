package qall

import (
	"github.com/qiniu/qlang/lib/bufio"
	"github.com/qiniu/qlang/lib/bytes"
	"github.com/qiniu/qlang/lib/crypto/md5"
	"github.com/qiniu/qlang/lib/encoding/hex"
	"github.com/qiniu/qlang/lib/encoding/json"
	"github.com/qiniu/qlang/lib/eqlang.v1"
	"github.com/qiniu/qlang/lib/errors"
	"github.com/qiniu/qlang/lib/io"
	"github.com/qiniu/qlang/lib/io/ioutil"
	"github.com/qiniu/qlang/lib/math"
	"github.com/qiniu/qlang/lib/meta"
	"github.com/qiniu/qlang/lib/net/http"
	"github.com/qiniu/qlang/lib/os"
	"github.com/qiniu/qlang/lib/path"
	"github.com/qiniu/qlang/lib/reflect"
	"github.com/qiniu/qlang/lib/runtime"
	"github.com/qiniu/qlang/lib/strconv"
	"github.com/qiniu/qlang/lib/strings"
	"github.com/qiniu/qlang/lib/sync"
	"github.com/qiniu/qlang/lib/terminal"
	"github.com/qiniu/qlang/lib/tpl.v1/extractor"
	"github.com/qiniu/qlang/lib/version"
	"github.com/qiniu/qlang/spec"

	// qlang builtin modules
	_ "github.com/qiniu/qlang/lib/builtin"
	_ "github.com/qiniu/qlang/lib/chan"
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
