package qlang_test

import (
	"testing"

	_ "qlang.io/qlang/builtin"
	_ "qlang.io/qlang/chan"
	"qlang.io/qlang/sync"
	"qlang.io/qlang.v2/qlang"
)

func init() {
	qlang.Import("sync", sync.Exports)
}

// -----------------------------------------------------------------------------

const testChanCode = `

mutex = sync.mutex()
ch = mkchan("bool", 2)

x = 1

go fn {
	defer ch <- true

	mutex.lock()
	defer mutex.unlock()
	x; x++
}

go fn {
	defer ch <- true

	mutex.lock()
	defer mutex.unlock()
	x; x++
}

<-ch
<-ch
`

func TestChan(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeExec([]byte(testChanCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.Var("x"); !ok || v != 3 {
		t.Fatal("x != 3, x =", v)
	}
}

// -----------------------------------------------------------------------------

const testGoroutineCode = `

mutex = sync.mutex()

wg = sync.waitGroup()
wg.add(2)

x = 1

go fn {
	defer wg.done()

	mutex.lock()
	defer mutex.unlock()
	x; x++
}

go fn {
	defer wg.done()

	mutex.lock()
	defer mutex.unlock()
	x; x++
}

wg.wait()
`

func TestGoroutine(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeExec([]byte(testGoroutineCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.Var("x"); !ok || v != 3 {
		t.Fatal("x != 3, x =", v)
	}
}

// -----------------------------------------------------------------------------


