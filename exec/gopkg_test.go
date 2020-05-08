package exec

import (
	"fmt"
	"testing"

	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

func init() {
	log.SetOutputLevel(log.Ldebug)
}

func Strcat(a, b string) string {
	return a + b
}

func execStrcat(arity uint32, p *Context) {
	args := p.GetArgs(2)
	ret := Strcat(args[0].(string), args[1].(string))
	p.Ret(2, ret)
}

func execSprint(arity uint32, p *Context) {
	args := p.GetArgs(arity)
	s := fmt.Sprint(args...)
	p.Ret(arity, s)
}

func execSprintf(arity uint32, p *Context) {
	args := p.GetArgs(arity)
	s := fmt.Sprintf(args[0].(string), args[1:]...)
	p.Ret(arity, s)
}

// I is a Go package instance.
var I = NewGoPackage("")

func init() {
	I.RegisterFuncvs(
		I.Funcv("Sprint", fmt.Sprint, execSprint),
		I.Funcv("Sprintf", fmt.Sprintf, execSprintf),
	)
	I.RegisterFuncs(
		I.Func("strcat", Strcat, execStrcat),
	)
	I.RegisterVars(
		I.Var("x", new(int)),
	)
}

func TestSprint(t *testing.T) {
	sprint, ok := I.FindFuncv("Sprint")
	if !ok {
		t.Fatal("FindFuncv failed: Sprint")
	}

	code := NewBuilder(nil).
		Push(5).
		Push("32").
		CallGoFuncv(sprint, 2).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != "532" {
		t.Fatal("5 `32` sprint != 532, ret =", v)
	}
}

func TestSprintf(t *testing.T) {
	sprintf, ok := I.FindFuncv("Sprintf")
	strcat, ok2 := I.FindFunc("strcat")
	if !ok || !ok2 {
		t.Fatal("FindFunc failed: Sprintf/strcat")
	}
	fmt.Println("sprintf:", sprintf.GetInfo())
	fmt.Println("strcat:", strcat.GetInfo())

	code := NewBuilder(nil).
		Push("Hello, %v, %d, %s").
		Push(1.3).
		Push(1).
		Push("x").
		Push("sw").
		CallGoFunc(strcat).
		CallGoFuncv(sprintf, 4).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != "Hello, 1.3, 1, xsw" {
		t.Fatal("format 1.3 1 `x` `sw` strcat sprintf != `Hello, 1.3, 1, xsw`, ret =", v)
	}
}

func TestLargeArity(t *testing.T) {
	sprint, kind, ok := FindGoPackage("").Find("Sprint")
	if !ok || kind != SymbolFuncv {
		t.Fatal("Find failed: Sprint")
	}

	b := NewBuilder(nil)
	ret := ""
	for i := 0; i < bitsFuncvArityMax+1; i++ {
		b.Push("32")
		ret += "32"
	}
	code := b.
		CallGoFuncv(GoFuncvAddr(sprint), bitsFuncvArityMax+1).
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != ret {
		t.Fatal("32 times(1024) sprint != `32` times(1024), ret =", v)
	}
}
