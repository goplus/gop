package lib

//go:generate go run ../cmd/qexp bytes flag fmt io/... log os reflect strconv strings sync/... time
//go:generate go run ../cmd/qexp github.com/goplus/gop/ast/gopq github.com/goplus/gop/ast/goptest

import (
	_ "github.com/goplus/gop/lib/builtin" // builtin
	_ "github.com/goplus/gop/lib/bytes"
	_ "github.com/goplus/gop/lib/errors"
	_ "github.com/goplus/gop/lib/flag"
	_ "github.com/goplus/gop/lib/fmt"
	_ "github.com/goplus/gop/lib/github.com/goplus/gop/ast/gopq"
	_ "github.com/goplus/gop/lib/github.com/goplus/gop/ast/goptest"
	_ "github.com/goplus/gop/lib/io"
	_ "github.com/goplus/gop/lib/io/ioutil"
	_ "github.com/goplus/gop/lib/log"
	_ "github.com/goplus/gop/lib/os"
	_ "github.com/goplus/gop/lib/reflect"
	_ "github.com/goplus/gop/lib/strconv"
	_ "github.com/goplus/gop/lib/strings"
	_ "github.com/goplus/gop/lib/sync"
	_ "github.com/goplus/gop/lib/sync/atomic"
	_ "github.com/goplus/gop/lib/time"
)
