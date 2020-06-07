package bytecode_test

import (
	"os"
	"path"
	"testing"

	"github.com/qiniu/goplus/cl/cltest"
)

// -----------------------------------------------------------------------------

func TestFromTestdata(t *testing.T) {
	sel, exclude := "", ""
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Getwd failed:", err)
	}
	dir = path.Join(dir, "../golang/testdata")
	cltest.FromTestdata(t, dir, sel, exclude)
}

// -----------------------------------------------------------------------------
