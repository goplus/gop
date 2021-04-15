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

// Package bytecode implements a bytecode backend for the Go+ language.
package bytecode

import (
	"bufio"
	"io"
	"reflect"
	"strconv"

	"github.com/goplus/gop/exec.spec"
)

// -----------------------------------------------------------------------------

const (
	bitsInstr = 32
	bitsOp    = 6

	bitsIntKind    = 3
	bitsFuncKind   = 2
	bitsFuncvArity = 10
	bitsVarScope   = 6
	bitsAssignOp   = 4
	bitsIndexOp    = 2
	bitsIsPtr      = 2
	bitsTwoValue   = 2

	bitsOpShift = bitsInstr - bitsOp
	bitsOperand = (1 << bitsOpShift) - 1

	bitsOpIndexShift   = bitsInstr - (bitsOp + bitsIndexOp)
	bitsOpIndexOperand = (1 << bitsOpIndexShift) - 1

	bitsOpZeroShift   = bitsInstr - (bitsOp + bitsIsPtr)
	bitsOpZeroOperand = (1 << bitsOpZeroShift) - 1

	bitsOpTypeAssertShift        = bitsInstr - (bitsOp + bitsTwoValue)
	bitsOpTypeAssertShiftOperand = (1 << bitsOpTypeAssertShift) - 1

	bitsOpInt        = bitsOp + bitsIntKind
	bitsOpIntShift   = bitsInstr - bitsOpInt
	bitsOpIntOperand = (1 << bitsOpIntShift) - 1

	bitsOpCallFuncv        = bitsOp + bitsFuncvArity
	bitsOpCallFuncvShift   = bitsInstr - bitsOpCallFuncv
	bitsOpCallFuncvOperand = (1 << bitsOpCallFuncvShift) - 1

	bitsFuncvArityOperand = (1 << bitsFuncvArity) - 1
	bitsFuncvArityVar     = bitsFuncvArityOperand
	bitsFuncvArityMax     = bitsFuncvArityOperand - 1

	bitsOpClosure        = bitsOp + bitsFuncKind
	bitsOpClosureShift   = bitsInstr - bitsOpClosure
	bitsOpClosureOperand = (1 << bitsOpClosureShift) - 1

	bitsOpVar        = bitsOp + bitsVarScope
	bitsOpVarShift   = bitsInstr - bitsOpVar
	bitsOpVarOperand = (1 << bitsOpVarShift) - 1

	bitsOpCaseNE        = bitsOp + 10
	bitsOpCaseNEShift   = bitsInstr - bitsOpCaseNE
	bitsOpCaseNEOperand = (1 << bitsOpCaseNEShift) - 1
)

// A Instr represents a instruction of the executor.
type Instr = uint32

const (
	opInvalid       = 0
	opCallGoFunc    = 1  // addr(26) - call a Go function
	opCallGoFuncv   = 2  // funvArity(10) addr(16) - call a Go function with variadic args
	opPushInt       = 3  // intKind(3) intVal(23)
	opPushUint      = 4  // intKind(3) intVal(23)
	opPushValSpec   = 5  // valSpec(26) - false=0, true=1
	opPushConstR    = 6  // idx(26)
	opIndex         = 7  // indexOp(2) idx(24)
	opMake          = 8  // funvArity(10) type(16)
	opAppend        = 9  // arity(26)
	opBuiltinOp     = 10 // reserved(16) kind(5) builtinOp(5)
	opJmp           = 11 // reserved(2) offset(24)
	opJmpIf         = 12 // flag(4) offset(22)
	opCaseNE        = 13 // n(10) offset(16)
	opPop           = 14 // n(26)
	opLoadVar       = 15 // varScope(6) addr(20)
	opStoreVar      = 16 // varScope(6) addr(20)
	opAddrVar       = 17 // varScope(6) addr(20) - load a variable's address
	opLoadGoVar     = 18 // addr(26)
	opStoreGoVar    = 19 // addr(26)
	opAddrGoVar     = 20 // addr(26)
	opAddrOp        = 21 // reserved(17) addressOp(4) kind(5)
	opCallFunc      = 22 // addr(26)
	opCallFuncv     = 23 // funvArity(10) addr(16)
	opReturn        = 24 // n(26)
	opLoad          = 25 // index(26)
	opAddr          = 26 // index(26)
	opStore         = 27 // index(26)
	opClosure       = 28 // funcKind(2) addr(24)
	opCallClosure   = 29 // arity(26)
	opGoClosure     = 30 // funcKind(2) addr(24)
	opCallGoClosure = 31 // arity(26)
	opMakeArray     = 32 // funvArity(10) type(16)
	opMakeMap       = 33 // funvArity(10) type(16)
	opZero          = 34 // isPtr(2) type(24)
	opForPhrase     = 35 // addr(26)
	opLstComprehens = 36 // addr(26)
	opMapComprehens = 37 // addr(26)
	opTypeCast      = 38 // type(26)
	opSlice         = 39 // i(13) j(13)
	opSlice3        = 40 // i(13) j(13)
	opMapIndex      = 41 // reserved(25) set(1)
	opGoBuiltin     = 42 // op(26)
	opErrWrap       = 43 // idx(26)
	opWrapIfErr     = 44 // reserved(2) offset(24)
	opDefer         = 45 // reserved(26)
	opGo            = 46 // arity(26)
	opLoadField     = 47 // op(26)
	opStoreField    = 48 // op(26)
	opAddrField     = 49 // op(26)
	opStruct        = 50 // funvArity(10) type(16)
	opSend          = 51 // reserved(26)
	opRecv          = 52 // reserved(26)
	opTypeAssert    = 53 // twoValue(2) type(24)
	opTypeMethod    = 54
)

const (
	iInvalid        = (opInvalid << bitsOpShift)
	iPushFalse      = (opPushValSpec << bitsOpShift)
	iPushTrue       = (opPushValSpec << bitsOpShift) | 1
	iPushNil        = (opPushValSpec << bitsOpShift) | 2
	iPushUnresolved = (opInvalid << bitsOpShift)

	iReturn = (opReturn << bitsOpShift) | (0xffffffff & bitsOperand)
)

const (
	ipInvalid = 0x7fffffff    // return
	ipReturnN = ipInvalid - 1 // return val1, val2, ...
)

// DecodeInstr returns
func DecodeInstr(i Instr) (InstrInfo, int32, int32) {
	op := i >> bitsOpShift
	v := instrInfos[op]
	p1, p2 := getParam(int32(i<<bitsOp), v.Params>>8)
	p2, _ = getParam(p2, v.Params&0xff)
	return v, p1, p2
}

func getParam(v int32, bits uint16) (int32, int32) {
	return v >> (32 - bits), v << bits
}

// InstrInfo represents the information of an instr.
type InstrInfo struct {
	Name   string
	Arg1   string
	Arg2   string
	Params uint16
}

var instrInfos = []InstrInfo{
	opInvalid:       {"invalid", "", "", 0},
	opCallGoFunc:    {"callGoFunc", "", "addr", 26},                       // addr(26) - call a Go function
	opCallGoFuncv:   {"callGoFuncv", "funvArity", "addr", (10 << 8) | 16}, // funvArity(10) addr(16) - call a Go function with variadic args
	opPushInt:       {"pushInt", "intKind", "intVal", (3 << 8) | 23},      // intKind(3) intVal(23)
	opPushUint:      {"pushUint", "intKind", "intVal", (3 << 8) | 23},     // intKind(3) intVal(23)
	opPushValSpec:   {"pushValSpec", "", "valSpec", 26},                   // valSpec(26) - false=0, true=1
	opPushConstR:    {"pushConstR", "", "idx", 26},                        // idx(26)
	opIndex:         {"index", "indexOp", "idx", (2 << 8) | 24},           // indexOp(2) idx(24)
	opMake:          {"make", "funvArity", "type", (10 << 8) | 16},        // funvArity(10) type(16)
	opAppend:        {"append", "", "arity", 26},                          // arity(26)
	opBuiltinOp:     {"builtinOp", "kind", "op", (21 << 8) | 5},           // reserved(16) kind(5) builtinOp(5)
	opJmp:           {"jmp", "", "offset", 26},                            // offset(26)
	opJmpIf:         {"jmpIf", "flag", "offset", (4 << 8) | 22},           // flag(4) offset(22)
	opCaseNE:        {"caseNE", "n", "offset", (10 << 8) | 16},            // n(10) offset(16)
	opPop:           {"pop", "", "n", 26},                                 // n(26)
	opLoadVar:       {"loadVar", "varScope", "addr", (6 << 8) | 20},       // varScope(6) addr(20)
	opStoreVar:      {"storeVar", "varScope", "addr", (6 << 8) | 20},      // varScope(6) addr(20)
	opAddrVar:       {"addrVar", "varScope", "addr", (6 << 8) | 20},       // varScope(6) addr(20) - load a variable's address
	opLoadGoVar:     {"loadGoVar", "", "addr", 26},                        // addr(26)
	opStoreGoVar:    {"storeGoVar", "", "addr", 26},                       // addr(26)
	opAddrGoVar:     {"addrGoVar", "", "addr", 26},                        // addr(26)
	opAddrOp:        {"addrOp", "op", "kind", (21 << 8) | 5},              // reserved(17) addressOp(4) kind(5)
	opCallFunc:      {"callFunc", "", "addr", 26},                         // addr(26)
	opCallFuncv:     {"callFuncv", "funvArity", "addr", (10 << 8) | 16},   // funvArity(10) addr(16)
	opReturn:        {"return", "", "n", 26},                              // n(26)
	opLoad:          {"load", "", "index", 26},                            // index(26)
	opAddr:          {"addr", "", "index", 26},                            // index(26)
	opStore:         {"store", "", "index", 26},                           // index(26)
	opClosure:       {"closure", "funcKind", "addr", (2 << 8) | 24},       // funcKind(2) addr(24)
	opCallClosure:   {"callClosure", "", "arity", 26},                     // arity(26)
	opGoClosure:     {"closureGo", "funcKind", "addr", (2 << 8) | 24},     // funcKind(2) addr(24)
	opCallGoClosure: {"callGoClosure", "", "arity", 26},                   // arity(26)
	opMakeArray:     {"makeArray", "funvArity", "type", (10 << 8) | 16},   // funvArity(10) type(16)
	opMakeMap:       {"makeMap", "funvArity", "type", (10 << 8) | 16},     // funvArity(10) type(16)
	opZero:          {"zero", "isPtr", "type", (2 << 8) | 24},             // isPtr(2) type(24)
	opForPhrase:     {"forPhrase", "", "addr", 26},                        // addr(26)
	opLstComprehens: {"listComprehension", "", "addr", 26},                // addr(26)
	opMapComprehens: {"mapComprehension", "", "addr", 26},                 // addr(26)
	opTypeCast:      {"typeCast", "", "type", 26},                         // type(26)
	opSlice:         {"slice", "i", "j", (13 << 8) | 13},                  // i(13) j(13)
	opSlice3:        {"slice3", "i", "j", (13 << 8) | 13},                 // i(13) j(13)
	opMapIndex:      {"mapIndex", "", "set", 26},                          // reserved(25) set(1)
	opGoBuiltin:     {"goBuiltin", "", "op", 26},                          // op(26)
	opErrWrap:       {"errWrap", "", "idx", 26},                           // idx(26)
	opWrapIfErr:     {"wrapIfErr", "", "offset", 26},                      // reserved(2) offset(24)
	opDefer:         {"defer", "", "", 0},                                 // reserved(26)
	opGo:            {"go", "", "arity", 26},                              // arity(26)
	opLoadField:     {"loadField", "", "", 26},                            // addr(26)
	opStoreField:    {"storeField", "", "", 26},                           // addr(26)
	opAddrField:     {"addrField", "", "", 26},                            // addr(26)
	opStruct:        {"struct", "funvArity", "type", (10 << 8) | 16},      // funvArity(10) type(16)
	opSend:          {"send", "", "", 0},                                  // reserved(26)
	opRecv:          {"recv", "", "", 0},                                  // reserved(26)
	opTypeAssert:    {"typeAssert", "twoValue", "type", (2 << 8) | 24},    // type(26)
	opTypeMethod:    {"typeMethod", "", "type", 26},                       // type(26)
}

// -----------------------------------------------------------------------------

// A Code represents generated instructions to execute.
type Code struct {
	data        []Instr
	valConsts   []interface{}
	funs        []*FuncInfo
	funvs       []*FuncInfo
	comprehens  []*Comprehension
	fors        []*ForPhrase
	types       []reflect.Type
	structs     []StructInfo
	errWraps    []errWrap
	typeMethods map[reflect.Type][]*exec.MethodInfo
	varManager
}

// NewCode returns a new Code object.
func NewCode() *Code {
	return &Code{data: make([]Instr, 0, 64), typeMethods: make(map[reflect.Type][]*exec.MethodInfo)}
}

// Len returns code length.
func (p *Code) Len() int {
	return len(p.data)
}

// Dump dumps code.
func (p *Code) Dump(w io.Writer) {
	DumpCodeBlock(w, p.data...)
}

// DumpCodeBlock dumps a code block.
func DumpCodeBlock(w io.Writer, data ...Instr) {
	b := bufio.NewWriter(w)
	for _, i := range data {
		v, p1, p2 := DecodeInstr(i)
		b.WriteString(v.Name)
		b.WriteByte(' ')
		if (v.Params & 0xff00) != 0 {
			b.WriteString(v.Arg1)
			b.WriteByte('=')
			b.WriteString(strconv.Itoa(int(p1)))
			b.WriteByte(' ')
		}
		if (v.Params & 0xff) != 0 {
			b.WriteString(v.Arg2)
			b.WriteByte('=')
			b.WriteString(strconv.Itoa(int(p2)))
		}
		b.WriteByte('\n')
	}
	b.Flush()
}

// -----------------------------------------------------------------------------

type anyUnresolved struct {
	offs []int
}

// Builder is a class that generates executing byte code.
type Builder struct {
	code   *Code
	labels map[*Label]int
	funcs  map[*FuncInfo]int
	types  map[reflect.Type]uint32
	*varManager
}

// NewBuilder creates a new Code Builder instance.
func NewBuilder(code *Code) *Builder {
	if code == nil {
		code = NewCode()
	}
	return &Builder{
		code:       code,
		labels:     make(map[*Label]int),
		funcs:      make(map[*FuncInfo]int),
		types:      make(map[reflect.Type]uint32),
		varManager: &code.varManager,
	}
}

// Resolve resolves all unresolved labels/functions/consts/etc.
func (p *Builder) Resolve() *Code {
	p.resolveLabels()
	p.resolveFuncs()
	return p.code
}

// -----------------------------------------------------------------------------

// Reserved represents a reserved instruction position.
type Reserved = exec.Reserved

// InvalidReserved is an invalid reserved position.
const InvalidReserved = exec.InvalidReserved

// Reserve reserves an instruction.
func (p *Builder) Reserve() Reserved {
	code := p.code
	idx := len(code.data)
	code.data = append(code.data, iInvalid)
	return Reserved(idx)
}

func (p *Builder) MethodOf(typ reflect.Type, infos []*exec.MethodInfo) {
	p.code.typeMethods[typ] = infos
	i := (opTypeMethod << bitsOpShift) | p.requireType(typ)
	p.code.data = append(p.code.data, i)
}

// -----------------------------------------------------------------------------
