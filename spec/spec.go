package spec

import (
	"github.com/qiniu/qlang/v6/exec"
)

// -----------------------------------------------------------------------------

// A Context represents the context of an executor.
type Context = exec.Context

// -----------------------------------------------------------------------------

// GoPackage represents a Go package.
type GoPackage = exec.GoPackage

// NewGoPackage creates a new builtin Go Package.
func NewGoPackage(pkgPath string) *GoPackage {
	return exec.NewGoPackage(pkgPath)
}

// ToStrings converts []interface{} into []string.
func ToStrings(args []interface{}) []string {
	ret := make([]string, len(args))
	for i, arg := range args {
		ret[i] = arg.(string)
	}
	return ret
}

// -----------------------------------------------------------------------------
