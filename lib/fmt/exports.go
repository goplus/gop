package fmt

import (
	"fmt"

	"github.com/qiniu/qlang/lib/builtin"

	qlang "github.com/qiniu/qlang/spec"
)

// -----------------------------------------------------------------------------

// I is a Go package instance.
var I = qlang.NewGoPackage("fmt")

func init() {
	I.RegisterVariadicFuncs(
		I.Func("Print", fmt.Print, builtin.QexecPrint),
		I.Func("Printf", fmt.Printf, builtin.QexecPrintf),
		I.Func("Println", fmt.Println, builtin.QexecPrintln),
		I.Func("Fprintln", fmt.Fprintln, builtin.QexecFprintln),
	)
}

// -----------------------------------------------------------------------------
