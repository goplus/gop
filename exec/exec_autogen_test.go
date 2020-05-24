package exec

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

// -----------------------------------------------------------------------------

var opAutogenOps = [...]string{
	OpAdd:       "Add",
	OpSub:       "Sub",
	OpMul:       "Mul",
	OpDiv:       "Div",
	OpMod:       "Mod",
	OpBitAnd:    "BitAnd",
	OpBitOr:     "BitOr",
	OpBitXor:    "BitXor",
	OpBitAndNot: "BitAndNot",
	OpBitSHL:    "BitSHL",
	OpBitSHR:    "BitSHR",
	OpLT:        "LT",
	OpLE:        "LE",
	OpGT:        "GT",
	OpGE:        "GE",
	OpEQ:        "EQ",
	OpNE:        "NE",
	OpLAnd:      "LAnd",
	OpLOr:       "LOr",
	OpNeg:       "Neg",
	OpNot:       "Not",
	OpBitNot:    "BitNot",
}

const autogenOpHeader = `package exec
`

const autogenBinaryOpUintTempl = `
func exec$Op$Type(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].($type) $op toUint(p.data[n-1])
	p.data = p.data[:n-1]
}
`

const autogenBinaryOpTempl = `
func exec$Op$Type(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].($type) $op p.data[n-1].($type)
	p.data = p.data[:n-1]
}
`

const autogenUnaryOpTempl = `
func exec$Op$Type(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = $opp.data[n-1].($type)
}
`

const autogenBuiltinOpHeader = `
var builtinOps = [...]func(i Instr, p *Context){`

const autogenBuiltinOpItemTempl = `
	(int($Type) << bitsOperator) | int(Op$Op): exec$Op$Type,`

const autogenBuiltinOpFooter = `
}
`

func autogenOpWithTempl(f *os.File, op Operator, Op string, templ string) {
	i := op.GetInfo()
	if templ == "" {
		templ = autogenBinaryOpTempl
		if i.InSecond == bitNone {
			templ = autogenUnaryOpTempl
		} else if i.InSecond == bitsAllIntUint {
			templ = autogenBinaryOpUintTempl
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

func _TestOpAutogen(t *testing.T) {
	f, err := os.Create("exec_op_autogen.go")
	if err != nil {
		t.Fatal("TestAutogen failed:", err)
	}
	defer f.Close()
	fmt.Fprint(f, autogenOpHeader)
	fmt.Fprint(f, autogenBuiltinOpHeader)
	for i, Op := range opAutogenOps {
		if Op != "" {
			autogenOpWithTempl(f, Operator(i), Op, autogenBuiltinOpItemTempl)
		}
	}
	fmt.Fprint(f, autogenBuiltinOpFooter)
	for i, Op := range opAutogenOps {
		if Op != "" {
			autogenOpWithTempl(f, Operator(i), Op, "")
		}
	}
}

func newKindValue(kind Kind) reflect.Value {
	t := TypeFromKind(kind)
	o := reflect.New(t).Elem()
	if kind >= Int && kind <= Int64 {
		o.SetInt(1)
	}
	if kind >= Uint && kind <= Uintptr {
		o.SetUint(1)
	}
	if kind >= Float32 && kind <= Float64 {
		o.SetFloat(1.0)
	}
	if kind >= Complex64 && kind <= Complex128 {
		o.SetComplex(1.0)
	}
	return o
}

func TestExecAutogenOp(t *testing.T) {
	add := OpAdd.GetInfo()
	fmt.Println("+", add, OpAdd.String())
	for i, execOp := range builtinOps {
		if execOp == nil {
			continue
		}
		kind := Kind(i >> bitsOperator)
		op := Operator(i & ((1 << bitsOperator) - 1))
		vars := []interface{}{
			newKindValue(kind).Interface(),
			newKindValue(kind).Interface(),
		}
		CallBuiltinOp(kind, op, vars...)
	}
}

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

func _TestAddrOpAutogen(t *testing.T) {
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

const autogenCastOpHeader = `package exec

import (
	"unsafe"
)
`

const autogenCastOpTempl = `
func exec$OpFrom$Type(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = $op(p.data[n-1].($type))
}
`

const autogenBuiltinCastOpHeader = `
var builtinCastOps = [...]func(i Instr, p *Context){`

const autogenBuiltinCastOpItemTempl = `
	(int($Op) << bitsKind) | int($Type): exec$OpFrom$Type,`

const autogenBuiltinCastOpFooter = `
}
`

func autogenCastOpWithTempl(f *os.File, op Kind, templ string) {
	i := builtinTypes[op]
	if templ == "" {
		templ = autogenCastOpTempl
	}
	for kind := Bool; kind <= UnsafePointer; kind++ {
		if (i.castFrom&(1<<kind)) == 0 || kind == op {
			continue
		}
		if kind == reflect.Ptr {
			continue
		}
		typ := TypeFromKind(kind).String()
		Typ := strings.Title(typ)
		opstr := op.String()
		Op := strings.Title(opstr)
		if op == UnsafePointer {
			Op = "UnsafePointer"
		}
		repl := strings.NewReplacer("$Op", Op, "$op", opstr, "$Type", Typ, "$type", typ)
		text := repl.Replace(templ)
		fmt.Fprint(f, text)
	}
}

func _TestCastOpAutogen(t *testing.T) {
	f, err := os.Create("exec_cast_autogen.go")
	if err != nil {
		t.Fatal("TestAutogen failed:", err)
	}
	defer f.Close()
	fmt.Fprint(f, autogenCastOpHeader)
	fmt.Fprint(f, autogenBuiltinCastOpHeader)
	for kind := Bool; kind <= UnsafePointer; kind++ {
		i := builtinTypes[kind]
		if i.castFrom != 0 {
			autogenCastOpWithTempl(f, kind, autogenBuiltinCastOpItemTempl)
		}
	}
	fmt.Fprint(f, autogenBuiltinCastOpFooter)
	for kind := Bool; kind <= UnsafePointer; kind++ {
		i := builtinTypes[kind]
		if i.castFrom != 0 {
			autogenCastOpWithTempl(f, kind, "")
		}
	}
}

func TestExecAutogenCastOp(t *testing.T) {
	for i, execOp := range builtinCastOps {
		if execOp == nil {
			continue
		}
		op := Kind(i >> bitsKind)
		kind := Kind(i & ((1 << bitsKind) - 1))
		vars := []interface{}{
			newKindValue(kind).Interface(),
		}
		callExecCastOp(kind, op, vars...)
	}
}

func callExecCastOp(kind Kind, op Kind, data ...interface{}) {
	i := (int(op) << bitsKind) | int(kind)
	ctx := newSimpleContext(data)
	execTypeCast((opAddrOp<<bitsOpShift)|uint32(i), ctx)
}

// -----------------------------------------------------------------------------
