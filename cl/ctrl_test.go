package qlang_test

import (
	"testing"

	"qlang.io/cl/qlang"
	_ "qlang.io/lib/builtin"
)

// -----------------------------------------------------------------------------

func TestIf(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(`fn { 0 }; if true { x = 3 } else { x = 5 }`), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 3 {
		t.Fatal("x =", v)
	}
}

func TestNormalFor(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(`for i = 3; i < 10; i++ {}`), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("i"); !ok || v != 10 {
		t.Fatal("i =", 10)
	}
}

func TestBreak(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(`for i = 3; i < 10; i++ { x = 1; y = 2; break }`), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("i"); !ok || v != 3 {
		t.Fatal("i =", v)
	}
}

func TestIfBreak(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(`for i = 3; i < 10; i++ { x = 1; y = 2; if i == 8 { z = 3; break } }`), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("i"); !ok || v != 8 {
		t.Fatal("i =", v)
	}
}

func TestContinue(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(`for i = 3; i < 10; i++ { continue; i=100 }`), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("i"); !ok || v != 10 {
		t.Fatal("i =", v)
	}
}

func TestSwitchBreak(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(`
		for i = 3; i < 10; i++ {
			x = 1
			y = 2
			switch {
			case true:
				z = 3
				break
			}
			i = 100
		}`), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("i"); !ok || v != 101 {
		t.Fatal("i =", v)
	}
}

func TestIfContinue(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(`for i = 3; i < 10; i++ { x = 1; y = 2; if true { z = 3; continue }; i=100 }`), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("i"); !ok || v != 10 {
		t.Fatal("i =", v)
	}
}

func TestSwitchContinue(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(`
		for i = 3; i < 10; i++ {
			x = 1
			y = 2
			switch {
			case true:
				z = 3
				continue
			}
			i = 100
		}`), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("i"); !ok || v != 10 {
		t.Fatal("i =", v)
	}
}

func TestSwitchContinue2(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(`
		for i = 3; i < 10; i++ {
			x = 1
			y = 2
			switch {
			case true:
				z = 4
			default:
				z = 3
				continue
			}
			i = 100
		}`), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("i"); !ok || v != 101 {
		t.Fatal("i =", v)
	}
}

func TestSwitchDefaultContinue(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(`
		for i = 3; i < 10; i++ {
			x = 1
			y = 2
			switch {
			case false:
				z = 4
			default:
				z = 3
				continue
			}
			i = 100
		}`), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("i"); !ok || v != 10 {
		t.Fatal("i =", v)
	}
}

// -----------------------------------------------------------------------------

func TestForRange1(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(`
		x = 0
		for range [1, 2] {
			x++
		}`), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 2 {
		t.Fatal("x =", v)
	}
}

func TestForRange2(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(`
		x = 0
		for i = range [10, 20] {
			x += i+1
		}`), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 3 {
		t.Fatal("x =", v)
	}
}

func TestForRange3(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(`
		x = 0
		for _, v = range [10, 20] {
			x += v
		}`), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 30 {
		t.Fatal("x =", v)
	}
}

func TestForRange4(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(`
		x = 0
		for range {10: 3, 20: 4} {
			x++
		}`), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 2 {
		t.Fatal("x =", v)
	}
}

func TestForRange5(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(`
		x = 0
		for k = range {10: 3, 20: 4} {
			x += k
		}`), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 30 {
		t.Fatal("x =", v)
	}
}

func TestForRange6(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(`
		x = 0
		for _, v = range {10: 3, 20: 4} {
			x += v
		}`), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 7 {
		t.Fatal("x =", v)
	}
}

func TestForRange7(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(`
		x = 0
		for k, v = range {10: 3, 20: 4} {
			x += k + v
		}`), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 37 {
		t.Fatal("x =", v)
	}
}

func TestForRange8(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(`

		ch = make(chan int, 2)
		ch <- 10
		ch <- 20
		close(ch)

		x = 0
		for range ch {
			x++
		}`), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 2 {
		t.Fatal("x =", v)
	}
}

func TestForRange9(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(`

		ch = make(chan int, 2)
		ch <- 10
		ch <- 20
		close(ch)

		x = 0
		for v = range ch {
			x += v
		}`), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 30 {
		t.Fatal("x =", v)
	}
}

func TestForRange10(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(`
		x = 0
		for i, v = range [10, 20, 30, 40] {
			if true {
				if i == 2 {
					break
				}
			}
			x += v
		}`), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 30 {
		t.Fatal("x =", v)
	}
}

func TestForRange11(t *testing.T) {

	lang := qlang.New()
	err := lang.SafeExec([]byte(`
		x = 0
		for i, v = range [10, 20, 30, 40] {
			if true {
				if i == 2 {
					continue
				}
			}
			x += v
		}`), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.GetVar("x"); !ok || v != 70 {
		t.Fatal("x =", v)
	}
}

// -----------------------------------------------------------------------------
