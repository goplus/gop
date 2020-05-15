package fmt

import (
	"fmt"

	"github.com/qiniu/qlang/v6/lib/builtin"

	qlang "github.com/qiniu/qlang/v6/spec"
)

// -----------------------------------------------------------------------------

// I is a Go package instance.
var I = qlang.NewGoPackage("fmt")

func init() {
	I.RegisterFuncvs(
		I.Funcv("Print", fmt.Print, builtin.QexecPrint),
		I.Funcv("Printf", fmt.Printf, builtin.QexecPrintf),
		I.Funcv("Println", fmt.Println, builtin.QexecPrintln),
		I.Funcv("Fprintln", fmt.Fprintln, builtin.QexecFprintln),
	)
}

// -----------------------------------------------------------------------------
