package qlang_all

import (
	"github.com/qiniu/qlang/lib/bufio"
	"github.com/qiniu/qlang/lib/crypto/md5"
	"github.com/qiniu/qlang/lib/encoding/json"
	"github.com/qiniu/qlang/lib/io/ioutil"
	"github.com/qiniu/qlang/lib/math"
	"github.com/qiniu/qlang/lib/os"
	"github.com/qiniu/qlang/lib/path"
	"github.com/qiniu/qlang/lib/reflect"
	"github.com/qiniu/qlang/lib/runtime"
	"github.com/qiniu/qlang/lib/strconv"
	"github.com/qiniu/qlang/lib/strings"
	"github.com/qiniu/qlang/lib/sync"
	"github.com/qiniu/qlang/lib/tpl.v1/extractor"
	"github.com/qiniu/qlang/lib/version"
	qlang "github.com/qiniu/qlang/spec"

	_ "github.com/qiniu/qlang/lib/builtin"
	_ "github.com/qiniu/qlang/lib/chan"
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
	qlang.Import("runtime", runtime.Exports)
	qlang.Import("strconv", strconv.Exports)
	qlang.Import("strings", strings.Exports)
	qlang.Import("sync", sync.Exports)
	qlang.Import("extractor", extractor.Exports)
}

// -----------------------------------------------------------------------------
