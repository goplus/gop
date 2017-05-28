package qlang_test

import (
	"testing"

	"qlang.io/cl/qlang"
	_ "qlang.io/lib/builtin"
	_ "qlang.io/lib/chan"
	"qlang.io/lib/sync"
)

func init() {
	qlang.Import("sync", sync.Exports)
}

// -----------------------------------------------------------------------------

const testChanCode = `

mutex = sync.mutex()
ch = make(chan bool, 2)

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

	lang := qlang.New()
	err := lang.SafeExec([]byte(testChanCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 3 {
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

	lang := qlang.New()
	err := lang.SafeExec([]byte(testGoroutineCode), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 3 {
		t.Fatal("x != 3, x =", v)
	}
}

// -----------------------------------------------------------------------------
