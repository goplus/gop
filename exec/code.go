package exec

import (
	"bufio"
	"io"
	"strconv"
)

// -----------------------------------------------------------------------------

const (
	bitsInstr = 32
	bitsOp    = 6

	bitsIntKind    = 3
	bitsFloatKind  = 2
	bitsFuncvArity = 10
	bitsVarScope   = 6
	bitsAssignOp   = 4

	bitsOpShift = bitsInstr - bitsOp
	bitsOperand = (1 << bitsOpShift) - 1

	bitsOpInt        = bitsOp + bitsIntKind
	bitsOpIntShift   = bitsInstr - bitsOpInt
	bitsOpIntOperand = (1 << bitsOpIntShift) - 1

	bitsOpFloat        = bitsOp + bitsFloatKind
	bitsOpFloatShift   = bitsInstr - bitsOpFloat
	bitsOpFloatOperand = (1 << bitsOpFloatShift) - 1

	bitsOpCallFuncv        = bitsOp + bitsFuncvArity
	bitsOpCallFuncvShift   = bitsInstr - bitsOpCallFuncv
	bitsOpCallFuncvOperand = (1 << bitsOpCallFuncvShift) - 1
	bitsFuncvArityOperand  = (1 << bitsFuncvArity) - 1
	bitsFuncvArityVar      = bitsFuncvArityOperand
	bitsFuncvArityMax      = bitsFuncvArityOperand - 1

	bitsOpVar        = bitsOp + bitsVarScope
	bitsOpVarShift   = bitsInstr - bitsOpVar
	bitsOpVarOperand = (1 << bitsOpVarShift) - 1
)

// A Instr represents a instruction of the executor.
type Instr = uint32

const (
	opInvalid     = 0
	opPushInt     = 1  // intKind(3) intVal(23)
	opPushUint    = 2  // intKind(3) intVal(23)
	opPushFloatR  = 3  // floatKind(2) floatIdx(24)
	opPushStringR = 4  // stringIdx(26)
	opPushIntR    = 5  // intKind(3) intIdx(23)
	opPushUintR   = 6  // intKind(3) intIdx(23)
	opPushValSpec = 7  // valSpec(26) - false=0, true=1
	opBuiltinOp   = 8  // reserved(16) kind(5) builtinOp(5)
	opJmp         = 9  // offset(26)
	opJmpIfFalse  = 10 // offset(26)
	opCaseNE      = 11 // offset(26)
	opPop         = 12 // n(26)
	opCallGoFunc  = 13 // addr(26) - call a Go function
	opCallGoFuncv = 14 // funvArity(10) addr(16) - call a Go function with variadic args
	opLoadVar     = 15 // varScope(6) addr(20)
	opStoreVar    = 16 // varScope(6) addr(20)
	opAddrVar     = 17 // varScope(6) addr(20) - load a variable's address
	opLoadGoVar   = 18 // addr(26)
	opStoreGoVar  = 19 // addr(26)
	opAddrGoVar   = 20 // addr(26)
	opAddrOp      = 21 // reserved(17) addressOp(4) kind(5)
	opCallFunc    = 22 // addr(26)
	opCallFuncv   = 23 // funvArity(10) addr(16)
	opReturn      = 24 // n(26)
	opLoad        = 25 // index(26)
	opStore       = 26 // index(26)
)

const (
	iInvalid        = (opInvalid << bitsOpShift)
	iPushFalse      = (opPushValSpec << bitsOpShift)
	iPushTrue       = (opPushValSpec << bitsOpShift) | 1
	iPushNil        = (opPushValSpec << bitsOpShift) | 2
	iPushUnresolved = (opInvalid << bitsOpShift)
	iReturn         = (opReturn << bitsOpShift) | (0xffffffff & bitsOperand)
)

const (
	ipInvalid = 0x7fffffff
	ipReturnN = ipInvalid - 1
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
	opInvalid:     {"invalid", "", "", 0},
	opPushInt:     {"pushInt", "intKind", "intVal", (3 << 8) | 23},        // intKind(3) intVal(23)
	opPushUint:    {"pushUint", "intKind", "intVal", (3 << 8) | 23},       // intKind(3) intVal(23)
	opPushFloatR:  {"pushFloatR", "floatKind", "floatIdx", (2 << 8) | 24}, // floatKind(2) floatIdx(24)
	opPushStringR: {"pushStringR", "", "stringIdx", 26},                   // stringIdx(26)
	opPushIntR:    {"pushIntR", "intKind", "intIdx", (3 << 8) | 23},       // intKind(3) intIdx(23)
	opPushUintR:   {"pushUintR", "intKind", "intIdx", (3 << 8) | 23},      // intKind(3) intIdx(23)
	opPushValSpec: {"pushValSpec", "", "valSpec", 26},                     // valSpec(26) - false=0, true=1
	opBuiltinOp:   {"builtinOp", "kind", "op", (21 << 8) | 5},             // reserved(16) kind(5) builtinOp(5)
	opJmp:         {"jmp", "", "offset", 26},                              // offset(26)
	opJmpIfFalse:  {"jmpIfFalse", "", "offset", 26},                       // offset(26)
	opCaseNE:      {"caseNE", "", "offset", 26},                           // offset(26)
	opPop:         {"pop", "", "n", 26},                                   // n(26)
	opCallGoFunc:  {"callGoFunc", "", "addr", 26},                         // addr(26) - call a Go function
	opCallGoFuncv: {"callGoFuncv", "funvArity", "addr", (10 << 8) | 16},   // funvArity(10) addr(16) - call a Go function with variadic args
	opLoadVar:     {"loadVar", "varScope", "addr", (6 << 8) | 20},         // varScope(6) addr(20)
	opStoreVar:    {"storeVar", "varScope", "addr", (6 << 8) | 20},        // varScope(6) addr(20)
	opAddrVar:     {"addrVar", "varScope", "addr", (6 << 8) | 20},         // varScope(6) addr(20) - load a variable's address
	opLoadGoVar:   {"loadGoVar", "", "addr", 26},                          // addr(26)
	opStoreGoVar:  {"storeGoVar", "", "addr", 26},                         // addr(26)
	opAddrGoVar:   {"addrGoVar", "", "addr", 26},                          // addr(26)
	opAddrOp:      {"addrOp", "op", "kind", (21 << 8) | 5},                // reserved(17) addressOp(4) kind(5)
	opCallFunc:    {"callFunc", "", "addr", 26},                           // addr(26)
	opCallFuncv:   {"callFuncv", "funvArity", "addr", (10 << 8) | 16},     // funvArity(10) addr(16)
	opReturn:      {"return", "", "n", 26},                                // n(26)
	opLoad:        {"load", "", "index", 26},                              // index(26)
	opStore:       {"store", "", "index", 26},                             // index(26)
}

// -----------------------------------------------------------------------------

// A Code represents generated instructions to execute.
type Code struct {
	data         []Instr
	stringConsts []string
	intConsts    []int64
	uintConsts   []uint64
	valConsts    []interface{}
	funs         []*FuncInfo
	funvs        []*FuncInfo
	structs      []StructInfo
	varManager
}

// NewCode returns a new Code object.
func NewCode() *Code {
	return &Code{data: make([]Instr, 0, 64)}
}

// Len returns code length.
func (p *Code) Len() int {
	return len(p.data)
}

// Dump dumps code.
func (p *Code) Dump(w io.Writer) {
	b := bufio.NewWriter(w)
	for _, i := range p.data {
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

type valUnresolved struct {
	op   Instr
	offs []int
}

// Builder class.
type Builder struct {
	code      *Code
	valConsts map[interface{}]*valUnresolved
	labels    map[*Label]int
	funcs     map[*FuncInfo]int
	*varManager
}

// NewBuilder creates a new Code Builder instance.
func NewBuilder(code *Code) *Builder {
	if code == nil {
		code = NewCode()
	}
	return &Builder{
		code:       code,
		valConsts:  make(map[interface{}]*valUnresolved),
		labels:     make(map[*Label]int),
		funcs:      make(map[*FuncInfo]int),
		varManager: &code.varManager,
	}
}

// Resolve resolves all unresolved labels/functions/consts/etc.
func (p *Builder) Resolve() *Code {
	p.resolveLabels()
	p.resolveConsts()
	p.resolveFuncs()
	return p.code
}

// -----------------------------------------------------------------------------

// Reserved represents a reserved instruction position.
type Reserved int

// InvalidReserved is an invalid reserved position.
const InvalidReserved Reserved = -1

// Reserve reserves an instruction.
func (p *Builder) Reserve() Reserved {
	code := p.code
	idx := len(code.data)
	code.data = append(code.data, iInvalid)
	return Reserved(idx)
}

// -----------------------------------------------------------------------------
