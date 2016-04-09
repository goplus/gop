package qlang_test

import (
	"testing"

	_ "qlang.io/qlang/builtin"
	"qlang.io/qlang.v2/qlang"
)

// -----------------------------------------------------------------------------

func TestIf(t *testing.T) {

	lang, err := qlang.New(nil)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeExec([]byte(`fn { 0 }; if true { x = 3 } else { x = 5 }`), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.Var("x"); !ok || v != 3 {
		t.Fatal("x =", v)
	}
}

func TestNormalFor(t *testing.T) {

	lang, err := qlang.New(nil)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeExec([]byte(`for i = 3; i < 10; i++ {}`), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.Var("i"); !ok || v != 10 {
		t.Fatal("i =", 10)
	}
}

func TestBreak(t *testing.T) {

	lang, err := qlang.New(nil)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeExec([]byte(`for i = 3; i < 10; i++ { x = 1; y = 2; break }`), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.Var("i"); !ok || v != 3 {
		t.Fatal("i =", v)
	}
}

func TestIfBreak(t *testing.T) {

	lang, err := qlang.New(nil)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeExec([]byte(`for i = 3; i < 10; i++ { x = 1; y = 2; if i == 8 { z = 3; break } }`), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.Var("i"); !ok || v != 8 {
		t.Fatal("i =", v)
	}
}

func TestContinue(t *testing.T) {

	lang, err := qlang.New(nil)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeExec([]byte(`for i = 3; i < 10; i++ { continue; i=100 }`), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.Var("i"); !ok || v != 10 {
		t.Fatal("i =", v)
	}
}

func TestSwitchBreak(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeExec([]byte(`
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
	if v, ok := lang.Var("i"); !ok || v != 101 {
		t.Fatal("i =", v)
	}
}

func TestIfContinue(t *testing.T) {

	lang, err := qlang.New(nil)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeExec([]byte(`for i = 3; i < 10; i++ { x = 1; y = 2; if true { z = 3; continue }; i=100 }`), "")
	if err != nil {
		t.Fatal("qlang.SafeExec:", err)
	}
	if v, ok := lang.Var("i"); !ok || v != 10 {
		t.Fatal("i =", v)
	}
}

func TestSwitchContinue(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeExec([]byte(`
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
	if v, ok := lang.Var("i"); !ok || v != 10 {
		t.Fatal("i =", v)
	}
}

func TestSwitchContinue2(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeExec([]byte(`
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
	if v, ok := lang.Var("i"); !ok || v != 101 {
		t.Fatal("i =", v)
	}
}

func TestSwitchDefaultContinue(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New:", err)
	}

	err = lang.SafeExec([]byte(`
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
	if v, ok := lang.Var("i"); !ok || v != 10 {
		t.Fatal("i =", v)
	}
}

// -----------------------------------------------------------------------------

