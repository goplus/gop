package exec

// -----------------------------------------------------------------------------
/*
func TestFunc(t *testing.T) {
	strcat, ok := I.FindFunc("strcat")
	if !ok {
		t.Fatal("FindFunc failed: strcat")
	}
	fmt.Println("strcat:", strcat.GetInfo())

	foo := NewFunc("foo")
	code := NewBuilder(nil).
		Push(nil).
		Push("x").
		Push("sw").
		CallFunc(foo).
		Return().
		DefineFunc(
			foo.Return(TyString).
			Args(TyString, TyString)).
		LoadArg(-2).
		LoadArg(-1).
		CallGoFunc(strcat).
		StoreArg(-3).
		EndFunc().
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := checkPop(ctx); v != "Hello, 1.3, 1, xsw" {
		t.Fatal("format 1.3 1 `x` `sw` strcat sprintf != `Hello, 1.3, 1, xsw`, ret =", v)
	}
}

/*
func TestFuncv(t *testing.T) {
	sprintf, ok := I.FindVariadicFunc("Sprintf")
	if !ok {
		t.Fatal("FindFunc failed: Sprintf")
	}
	fmt.Println("sprintf:", sprintf.GetInfo())

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
*/

// -----------------------------------------------------------------------------
