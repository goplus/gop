package qlang_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/qiniu/qlang"

	qlangv1 "github.com/qiniu/qlang/cl"
	qlangspec "github.com/qiniu/qlang/spec" // 导入 qlang spec
	qip "github.com/qiniu/text/tpl/interpreter"

	_ "github.com/qiniu/qlang/lib/builtin" // 导入 builtin 包
)

// -----------------------------------------------------------------------------

func TestInit(t *testing.T) {

	lang, err := qlang.New(nil)
	if err != nil {
		t.Fatal("qlang.New failed:", err)
	}

	err = lang.SafeEval(`a = switch { true: 1 false: 2 }`)
	if err == nil {
		t.Fatal("Eval ok?")
	}

	err = lang.SafeEval(`a = switch { case true: 1 case false: 2 }`)
	if err != nil {
		t.Fatal("Eval failed:", err)
	}

	v, ok := lang.Ret()
	if !ok {
		t.Fatal("Ret failed:", err)
	}

	if v != 1 {
		t.Fatal("invalid value of v:", v)
	}
}

func TestDefer1(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New failed:", err)
	}

	err = lang.SafeEval(`
		a = 1
		defer fn() {
			a; a = 2
		}()

		if true {
			b = 1
		}
		println(a)
		a = 3
	`)
	if err != nil {
		t.Fatal("Eval failed:", err)
	}
	if v, ok := lang.Var("a"); !ok || v != 2 {
		t.Fatal("invalid value of v:", v)
	}
}

func TestDefer2(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New failed:", err)
	}

	err = lang.SafeEval(`
		a = 1
		y = fn() {
			defer fn() {
				a; a = 2
			}()

			if true {
				b = 1
			}
			println(a)
			a = 3
		}
		y()
	`)
	if err != nil {
		t.Fatal("Eval failed:", err)
	}
	if v, ok := lang.Var("a"); !ok || v != 2 {
		t.Fatal("invalid value of v:", v)
	}
}

func TestReturn(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New failed:", err)
	}

	errReturn2 := &qip.RuntimeError{Err: errors.New("return")}
	if qlangv1.ErrReturn == errReturn2 {
		t.Fatal("ErrReturn == errReturn2")
	}

	err = lang.SafeEval(`
		x = fn() { return 2 }
		a = x()
	`)
	if err != nil {
		t.Fatal("Eval failed:", err)
	}
	if v, ok := lang.Var("a"); !ok || v != 2 {
		t.Fatal("invalid value of v:", v)
	}
}

func TestPanic(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New failed:", err)
	}

	err = lang.SafeEval(`
		panic("error message")
		println("code unreachable")
	`)
	if err == nil || strings.Index(err.Error(), "error message") < 0 {
		t.Fatal("panic failed:", err)
	}
}

func TestRecover(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New failed:", err)
	}

	err = lang.SafeEval(`
		a = 1
		defer fn() {
			a; a = recover()
		}()
		panic(2)
	`)
	if err != nil {
		t.Fatal("recover failed:", err)
	}
	if v, ok := lang.Var("a"); !ok || v != 2 {
		t.Fatal("invalid value of v:", v)
	}
}

func TestMap(t *testing.T) {

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New failed:", err)
	}

	err = lang.SafeEval(`
		x = {"a": 1, "b": 2}
		a = x["a"]
		c = x["c"]
		if c == undefined {
			println("ok!", a)
		}
	`)
	if err != nil {
		t.Fatal("Eval failed:", err)
	}
	if a, ok := lang.Var("a"); !ok || a != 1 {
		t.Fatal("invalid value of a:", a)
	}
	if c, ok := lang.Var("c"); !ok || c != qlangspec.Undefined {
		t.Fatal("invalid value of c:", c)
	}
}

// -----------------------------------------------------------------------------

type AutoCaller struct {
}

func NewAutoCaller() *AutoCaller {
	return new(AutoCaller)
}

func (p *AutoCaller) Text() string {
	return "abc"
}

func (p *AutoCaller) Get(a int) int {
	return a + 1
}

func TestAutoCall(t *testing.T) {

	qlang.Import("", map[string]interface{}{
		"newAutoCaller": NewAutoCaller,
	})

	qlang.SetAutoCall(reflect.TypeOf((*AutoCaller)(nil)))

	lang, err := qlang.New(qlang.InsertSemis)
	if err != nil {
		t.Fatal("qlang.New failed:", err)
	}

	err = lang.SafeEval(`
		p = newAutoCaller()
		println(p.text)
		a = p.text
		b = p.get(1)
	`)
	if err != nil {
		t.Fatal("Eval failed:", err)
	}
	if a, ok := lang.Var("a"); !ok || a != "abc" {
		t.Fatal("invalid value of a:", a)
	}
	if b, ok := lang.Var("b"); !ok || b != 2 {
		t.Fatal("invalid value of b:", b)
	}
}

// -----------------------------------------------------------------------------
