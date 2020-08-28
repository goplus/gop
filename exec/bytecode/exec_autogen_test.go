/*
 Copyright 2020 The GoPlus Authors (goplus.org)

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package bytecode

import (
	"fmt"
	"os"
	"strings"
	"testing"

	exec "github.com/goplus/gop/exec.spec"
	"github.com/goplus/gop/reflect"
)

const (
	OpNot       = OpLNot
	OpDiv       = OpQuo
	OpBitAnd    = OpAnd
	OpBitOr     = OpOr
	OpBitXor    = OpXor
	OpBitAndNot = OpAndNot
	OpBitSHL    = OpLsh
	OpBitSHR    = OpRsh

	OpDivAssign       = OpQuoAssign
	OpBitAndAssign    = OpAndAssign
	OpBitOrAssign     = OpOrAssign
	OpBitXorAssign    = OpXorAssign
	OpBitAndNotAssign = OpAndNotAssign
	OpBitSHLAssign    = OpLshAssign
	OpBitSHRAssign    = OpRshAssign
)

// -----------------------------------------------------------------------------

var opAutogenOps = [...]string{
	OpAdd:    "Add",
	OpSub:    "Sub",
	OpMul:    "Mul",
	OpQuo:    "Quo",
	OpMod:    "Mod",
	OpAnd:    "And",
	OpOr:     "Or",
	OpXor:    "Xor",
	OpAndNot: "AndNot",
	OpLsh:    "Lsh",
	OpRsh:    "Rsh",
	OpLT:     "LT",
	OpLE:     "LE",
	OpGT:     "GT",
	OpGE:     "GE",
	OpEQ:     "EQ",
	OpNE:     "NE",
	OpLAnd:   "LAnd",
	OpLOr:    "LOr",
	OpLNot:   "LNot",
	OpNeg:    "Neg",
	OpBitNot: "BitNot",
}

const autogenOpHeader = `package bytecode

import (
	"math/big"
)
`

const autogenStdBinaryOpUintTempl = `
func exec$Op$Type(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].($type) $op toUint(p.data[n-1])
	p.data = p.data[:n-1]
}
`

const autogenBigBinaryOpUintTempl = `
func exec$Op$Type(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = new($type).$Op(p.data[n-2].(*$type), toUint(p.data[n-1]))
	p.data = p.data[:n-1]
}
`

var autogenBinaryOpUintTempl = [2]string{
	autogenStdBinaryOpUintTempl,
	autogenBigBinaryOpUintTempl,
}

const autogenStdBinaryOpTempl = `
func exec$Op$Type(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].($type) $op p.data[n-1].($type)
	p.data = p.data[:n-1]
}
`

const autogenBigBinaryOpTempl = `
func exec$Op$Type(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = new($type).$Op(p.data[n-2].(*$type), p.data[n-1].(*$type))
	p.data = p.data[:n-1]
}
`

const autogenBigBoolBinaryOpTempl = `
func exec$Op$Type(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(*$type).Cmp(p.data[n-1].(*$type)) $op 0
	p.data = p.data[:n-1]
}
`

var autogenBinaryOpTempl = [2]string{
	autogenStdBinaryOpTempl,
	autogenBigBinaryOpTempl,
}

var autogenBoolOpTempl = [2]string{
	autogenStdBinaryOpTempl,
	autogenBigBoolBinaryOpTempl,
}

const autogenStdUnaryOpTempl = `
func exec$Op$Type(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = $opp.data[n-1].($type)
}
`

const autogenBigUnaryOpTempl = `
func exec$Op$Type(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = new($type).$Op(p.data[n-1].(*$type))
}
`

const autogenBigNotTempl = `
func exec$Op$Type(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = new($type).Not(p.data[n-1].(*$type))
}
`

var autogenUnaryOpTempl = [2]string{
	autogenStdUnaryOpTempl,
	autogenBigUnaryOpTempl,
}

const autogenBuiltinOpHeader = `
var builtinOps = [...]func(i Instr, p *Context){
	(int(BigInt) << bitsOperator) | int(OpQuo): execQuoBigInt,`

const autogenStdBuiltinOpItemTempl = `
	(int($Type) << bitsOperator) | int(Op$Op): exec$Op$Type,`

var autogenBuiltinOpItemTempl = [2]string{
	autogenStdBuiltinOpItemTempl,
	autogenStdBuiltinOpItemTempl,
}

const autogenBuiltinOpFooter = `
}
`

func autogenOpWithTempl(f *os.File, op Operator, Op string, templ *[2]string) {
	i := op.GetInfo()
	if templ == nil {
		templ = &autogenBinaryOpTempl
		if i.InSecond == exec.BitNone {
			templ = &autogenUnaryOpTempl
		} else if i.InSecond == exec.BitsAllIntUint {
			templ = &autogenBinaryOpUintTempl
		} else if i.Out == Bool && i.InFirst != (1<<Bool) {
			templ = &autogenBoolOpTempl
		}
	}
	for kind := Bool; kind <= BigFloat; kind++ {
		if (i.InFirst & (1 << kind)) == 0 {
			continue
		}
		typ := exec.TypeFromKind(kind).String()
		Typ, tpl := typ, templ[0]
		if strings.HasPrefix(typ, "*") { // *Big.Int, *Big.Rat, *Big.Float
			Typ, tpl = typ[1:4]+typ[5:], templ[1]
			typ = typ[1:]
			if templ == &autogenUnaryOpTempl && Op == "BitNot" {
				tpl = autogenBigNotTempl
			}
		}
		Typ = strings.Title(Typ)
		repl := strings.NewReplacer("$Op", Op, "$op", i.Lit, "$Type", Typ, "$type", typ)
		text := repl.Replace(tpl)
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
			autogenOpWithTempl(f, Operator(i), Op, &autogenBuiltinOpItemTempl)
		}
	}
	fmt.Fprint(f, autogenBuiltinOpFooter)
	for i, Op := range opAutogenOps {
		if Op != "" {
			autogenOpWithTempl(f, Operator(i), Op, nil)
		}
	}
}

func newKindValue(kind Kind) reflect.Value {
	t := exec.TypeFromKind(kind)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	o := reflect.New(t).Elem()
	if kind >= Int && kind <= Int64 {
		o.SetInt(1)
	} else if kind >= Uint && kind <= Uintptr {
		o.SetUint(1)
	} else if kind >= Float32 && kind <= Float64 {
		o.SetFloat(1.0)
	} else if kind >= Complex64 && kind <= Complex128 {
		o.SetComplex(1.0)
	} else if kind > UnsafePointer {
		x := reflect.ValueOf(int64(1))
		o = o.Addr()
		o.MethodByName("SetInt64").Call([]reflect.Value{x})
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
	OpAddAssign:    "AddAssign",
	OpSubAssign:    "SubAssign",
	OpMulAssign:    "MulAssign",
	OpQuoAssign:    "QuoAssign",
	OpModAssign:    "ModAssign",
	OpAndAssign:    "AndAssign",
	OpOrAssign:     "OrAssign",
	OpXorAssign:    "XorAssign",
	OpAndNotAssign: "AndNotAssign",
	OpLshAssign:    "LshAssign",
	OpRshAssign:    "RshAssign",
	OpInc:          "Inc",
	OpDec:          "Dec",
}

const autogenAddrOpHeader = `package bytecode

import (
	"math/big"
)
`

const autogenStdBinaryAddrOpUintTempl = `
func exec$Op$Type(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*$type) $op toUint(p.data[n-2])
	p.data = p.data[:n-2]
}
`

const autogenBigBinaryAddrOpUintTempl = `
func exec$Op$Type(i Instr, p *Context) {
	n := len(p.data)
	x := p.data[n-1].(*$type)
	x.$OP(x, toUint(p.data[n-2]))
	p.data = p.data[:n-2]
}
`

var autogenBinaryAddrOpUintTempl = [2]string{
	autogenStdBinaryAddrOpUintTempl,
	autogenBigBinaryAddrOpUintTempl,
}

const autogenStdBinaryAddrOpTempl = `
func exec$Op$Type(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*$type) $op p.data[n-2].($type)
	p.data = p.data[:n-2]
}
`

const autogenBigBinaryAddrOpTempl = `
func exec$Op$Type(i Instr, p *Context) {
	n := len(p.data)
	x := p.data[n-1].(*$type)
	x.$OP(x, p.data[n-2].(*$type))
	p.data = p.data[:n-2]
}
`

var autogenBinaryAddrOpTempl = [2]string{
	autogenStdBinaryAddrOpTempl,
	autogenBigBinaryAddrOpTempl,
}

const autogenStdUnaryAddrOpTempl = `
func exec$Op$Type(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*$type))$op
	p.data = p.data[:n-1]
}
`

const autogenBigUnaryAddrOpTempl = `
func exec$Op$Type(i Instr, p *Context) {
	n := len(p.data)
	x := p.data[n-1].(*$type)
	x.$OP(x, $one)
	p.data = p.data[:n-1]
}
`

var autogenUnaryAddrOpTempl = [2]string{
	autogenStdUnaryAddrOpTempl,
	autogenBigUnaryAddrOpTempl,
}

const autogenBuiltinAddrOpHeader = `
var builtinAddrOps = [...]func(i Instr, p *Context){`

const autogenStdBuiltinAddrOpItemTempl = `
	(int(Op$Op) << bitsKind) | int($Type): exec$Op$Type,`

var autogenBuiltinAddrOpItemTempl = [2]string{
	autogenStdBuiltinAddrOpItemTempl,
	autogenStdBuiltinAddrOpItemTempl,
}

const autogenBuiltinAddrOpFooter = `
}
`

func autogenAddrOpWithTempl(f *os.File, op AddrOperator, Op string, templ *[2]string) {
	i := op.GetInfo()
	if templ == nil {
		templ = &autogenBinaryAddrOpTempl
		if i.InSecond == exec.BitNone {
			templ = &autogenUnaryAddrOpTempl
		} else if i.InSecond == exec.BitsAllIntUint {
			templ = &autogenBinaryAddrOpUintTempl
		}
	}
	for kind := Bool; kind <= BigFloat; kind++ {
		if (i.InFirst & (1 << kind)) == 0 {
			continue
		}
		typ := exec.TypeFromKind(kind).String()
		Typ, tpl, one := typ, templ[0], "1"
		if strings.HasPrefix(typ, "*") { // *Big.Int, *Big.Rat, *Big.Float
			Typ, tpl = typ[1:4]+typ[5:], templ[1]
			typ = typ[1:]
			one = ones[kind-BigInt]
		}
		Typ = strings.Title(Typ)
		Op2 := "Add"
		if i.InSecond == exec.BitNone {
			if i.Lit == "--" {
				Op2 = "Sub"
			}
		} else {
			Op2 = strings.TrimSuffix(Op, "Assign")
		}
		repl := strings.NewReplacer(
			"$Op", Op, "$OP", Op2, "$one", one, "$op", i.Lit, "$Type", Typ, "$type", typ)
		text := repl.Replace(tpl)
		fmt.Fprint(f, text)
	}
}

var (
	ones = [3]string{"bigIntOne", "bigRatOne", "bigFloatOne"}
)

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
			autogenAddrOpWithTempl(f, AddrOperator(i), Op, &autogenBuiltinAddrOpItemTempl)
		}
	}
	fmt.Fprint(f, autogenBuiltinAddrOpFooter)
	for i, Op := range opAutogenAddrOps {
		if Op != "" {
			autogenAddrOpWithTempl(f, AddrOperator(i), Op, nil)
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
		addr := newKindValue(kind)
		if kind < BigInt {
			addr = addr.Addr()
		}
		vars := []interface{}{
			newKindValue(kind).Interface(),
			addr.Interface(),
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

/*
const autogenCastOpHeader = `package bytecode
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

func execTypeCast(i Instr, p *Context) {
	if fn := builtinCastOps[int(i&bitsOperand)]; fn != nil {
		fn(0, p)
		return
	}
	log.Panicln("execTypeCast: invalid instr -", i)
}
*/

// -----------------------------------------------------------------------------
