package exec

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// -----------------------------------------------------------------------------

var opAutogenAddrOps = [...]string{
	OpAddAssign:       "AddAssign",
	OpSubAssign:       "SubAssign",
	OpMulAssign:       "MulAssign",
	OpDivAssign:       "DivAssign",
	OpModAssign:       "ModAssign",
	OpBitAndAssign:    "BitAndAssign",
	OpBitOrAssign:     "BitOrAssign",
	OpBitXorAssign:    "BitXorAssign",
	OpBitAndNotAssign: "BitAndNotAssign",
	OpBitSHLAssign:    "BitSHLAssign",
	OpBitSHRAssign:    "BitSHRAssign",
	OpInc:             "Inc",
	OpDec:             "Dec",
}

const autogenAddrOpHeader = `package exec
`

const autogenBinaryAddrOpUintTempl = `
func exec$Op$Type(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*$type) $op toUint(p.data[n-2])
	p.data = p.data[:n-2]
}
`

const autogenBinaryAddrOpTempl = `
func exec$Op$Type(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*$type) $op p.data[n-2].($type)
	p.data = p.data[:n-2]
}
`

const autogenUnaryAddrOpTempl = `
func exec$Op$Type(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*$type))$op
	p.data = p.data[:n-1]
}
`

const autogenBuiltinAddrOpHeader = `
var builtinAddrOps = [...]func(i Instr, p *Context){`

const autogenBuiltinAddrOpItemTempl = `
	(int(Op$Op) << bitsKind) | int($Type): exec$Op$Type,`

const autogenBuiltinAddrOpFooter = `
}
`

func autogenAddrOpWithTempl(f *os.File, op AddrOperator, Op string, templ string) {
	i := op.GetInfo()
	if templ == "" {
		templ = autogenBinaryAddrOpTempl
		if i.InSecond == bitNone {
			templ = autogenUnaryAddrOpTempl
		} else if i.InSecond == bitsAllIntUint {
			templ = autogenBinaryAddrOpUintTempl
		}
	}
	for kind := Bool; kind <= UnsafePointer; kind++ {
		if (i.InFirst & (1 << kind)) == 0 {
			continue
		}
		typ := TypeFromKind(kind).String()
		Typ := strings.Title(typ)
		repl := strings.NewReplacer("$Op", Op, "$op", i.Lit, "$Type", Typ, "$type", typ)
		text := repl.Replace(templ)
		fmt.Fprint(f, text)
	}
}

func TestAddrOpAutogen(t *testing.T) {
	f, err := os.Create("exec_addrop_autogen.go")
	if err != nil {
		t.Fatal("TestAutogen failed:", err)
	}
	defer f.Close()
	fmt.Fprint(f, autogenAddrOpHeader)
	fmt.Fprint(f, autogenBuiltinAddrOpHeader)
	for i, Op := range opAutogenAddrOps {
		if Op != "" {
			autogenAddrOpWithTempl(f, AddrOperator(i), Op, autogenBuiltinAddrOpItemTempl)
		}
	}
	fmt.Fprint(f, autogenBuiltinAddrOpFooter)
	for i, Op := range opAutogenAddrOps {
		if Op != "" {
			autogenAddrOpWithTempl(f, AddrOperator(i), Op, "")
		}
	}
}

func TestExecAutogenAddrOp(t *testing.T) {
	add := OpAddAssign.GetInfo()
	fmt.Println("+", add, OpAddAssign.String())
	fmt.Println("=", OpAssign.String())
	fmt.Println("*", OpAddrVal.String())
	for i, execOp := range builtinAddrOps {
		if execOp == nil {
			continue
		}
		op := AddrOperator(i >> bitsKind)
		kind := Kind(i & ((1 << bitsKind) - 1))
		vars := []interface{}{
			newKindValue(kind).Interface(),
			newKindValue(kind).Addr().Interface(),
		}
		CallAddrOp(kind, op, vars...)
		callExecAddrOp(kind, op, vars...)
	}
}

func callExecAddrOp(kind Kind, op AddrOperator, data ...interface{}) {
	i := (int(op) << bitsKind) | int(kind)
	ctx := newSimpleContext(data)
	execAddrOp((opAddrOp<<bitsOpShift)|uint32(i), ctx)
}

// -----------------------------------------------------------------------------

func TestVar(t *testing.T) {
	sprint, ok := I.FindFuncv("Sprint")
	strcat, ok2 := I.FindFunc("strcat")
	if !ok || !ok2 {
		t.Fatal("FindFunc failed: Sprintf/strcat")
	}

	x := NewVar(TyString, "x")
	y := NewVar(TyString, "y")
	code := NewBuilder(nil).
		DefineVar(x, y).
		Push(5).
		Push("32").
		CallGoFuncv(sprint, 2).
		StoreVar(x). // x = sprint(5, "32")
		Push("78").
		LoadVar(x).
		CallGoFunc(strcat).
		StoreVar(y). // y = strcat("78", x)
		Resolve()

	ctx := NewContext(code)
	ctx.Exec(0, code.Len())
	if v := ctx.getVar(1); v != "78532" {
		t.Fatal("y != 78532, ret =", v)
	}
}

func TestParentCtx(t *testing.T) {
	sprint, ok := I.FindFuncv("Sprint")
	strcat, ok2 := I.FindFunc("strcat")
	if !ok || !ok2 {
		t.Fatal("FindFunc failed: Sprintf/strcat")
	}

	x := NewVar(TyString, "x")
	y := NewVar(TyString, "y")
	z := NewVar(TyString, "z")
	code := NewBuilder(nil).
		DefineVar(z).
		DefineFunc(NewFunc("", 2).Args()).
		DefineVar(x, y).
		Push(5).
		Push("32").
		CallGoFuncv(sprint, 2).
		StoreVar(x). // x = sprint(5, "32")
		LoadVar(z).
		LoadVar(x).
		CallGoFunc(strcat).
		StoreVar(z). // z = strcat(z, x)
		Resolve()

	p1 := NewContext(code)
	p1.SetVar(z, "78")

	p2 := newContextEx(p1, p1.Stack, p1.code, nil)
	ctx := newContextEx(p2, p2.Stack, p2.code, newVarManager(x, y))
	ctx.Exec(0, code.Len())
	if v := p1.GetVar(z); v != "78532" {
		t.Fatal("z != 78532, ret =", v)
	}
}

func TestAddrVar(t *testing.T) {
	sprint, ok := I.FindFuncv("Sprint")
	strcat, ok2 := I.FindFunc("strcat")
	if !ok || !ok2 {
		t.Fatal("FindFunc failed: Sprintf/strcat")
	}

	x := NewVar(TyString, "x")
	y := NewVar(TyString, "y")
	z := NewVar(TyString, "z")
	code := NewBuilder(nil).
		DefineVar(z).
		DefineFunc(NewFunc("", 2).Args()).
		DefineVar(x, y).
		Push(5).
		Push("32").
		CallGoFuncv(sprint, 2).
		StoreVar(x). // x = sprint(5, "32")
		LoadVar(z).
		AddrVar(x). // &x
		AddrOp(String, OpAddrVal).
		CallGoFunc(strcat).
		AddrVar(z).
		AddrOp(String, OpAssign). // z = strcat(z, *&x)
		Resolve()

	p1 := NewContext(code)
	p1.SetVar(z, "78")

	p2 := newContextEx(p1, p1.Stack, p1.code, nil)
	ctx := newContextEx(p2, p2.Stack, p2.code, newVarManager(x, y))
	ctx.Exec(0, code.Len())
	if v := p1.GetVar(z); v != "78532" {
		t.Fatal("y != 78532, ret =", v)
	}
}

// -----------------------------------------------------------------------------
