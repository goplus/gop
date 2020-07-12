package cl_test

import (
	"testing"

	"github.com/goplus/gop/cl/cltest"
)

// -----------------------------------------------------------------------------

func TestUnbound(t *testing.T) {
	cltest.Expect(t,
		`println("Hello " + "qiniu:", 123, 4.5, 7i)`,
		"Hello qiniu: 123 4.5 (0+7i)\n",
	)
}

func TestPanic(t *testing.T) {
	cltest.Expect(t,
		`panic("Helo")`, "", "Helo",
	)
}

// -----------------------------------------------------------------------------
