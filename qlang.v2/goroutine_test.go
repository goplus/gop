package qlang_test

import (
	"testing"

	_ "qlang.io/qlang/builtin"
	"qlang.io/qlang/sync"
	"qlang.io/qlang.v2/qlang"
)

// -----------------------------------------------------------------------------

const testGoroutineCode = `

wg = sync.waitGroup()
wg.add(2)

x = 1

go fn {
	defer wg.done()
	x; x++
}

go fn {
	defer wg.done()
	x; x++
}

wg.wait()
`

func TestGoroutine(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}
	qlang.Import("sync", sync.Exports)

	err = lang.SafeExec([]byte(testGoroutineCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.Var("x"); !ok || v != 3 {
		t.Fatal("x != 3, x =", v)
	}
}

// -----------------------------------------------------------------------------


