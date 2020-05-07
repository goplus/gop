package exec

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
	bitsFuncvArityMax      = (1 << bitsFuncvArity) - 1

	bitsOpVar        = bitsOp + bitsVarScope
	bitsOpVarShift   = bitsInstr - bitsOpVar
	bitsOpVarOperand = (1 << bitsOpVarShift) - 1
)

// A Instr represents a instruction of the executor.
type Instr = uint32

const (
	opInvalid     = (1 << bitsOp) - 1
	opPushInt     = 0  // intKind(3) intVal(23)
	opPushUint    = 1  // intKind(3) intVal(23)
	opPushValSpec = 2  // valSpec(26) - false=0, true=1
	opPushFloatR  = 3  // floatKind(2) floatIdx(24)
	opPushStringR = 4  // stringIdx(26)
	opPushIntR    = 5  // intKind(3) intIdx(23)
	opPushUintR   = 6  // intKind(3) intIdx(23)
	opBuiltinOp   = 7  // reserved(16) kind(5) builtinOp(5)
	opReserved    = 8  // reserved for future
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
	opReturn      = 24 // reserved(26)
)

const (
	iInvalid        = (opInvalid << bitsOpShift)
	iPushFalse      = (opPushValSpec << bitsOpShift)
	iPushTrue       = (opPushValSpec << bitsOpShift) | 1
	iPushNil        = (opPushValSpec << bitsOpShift) | 2
	iPushUnresolved = (opInvalid << bitsOpShift)
)

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
}

// NewCode returns a new Code object.
func NewCode() *Code {
	return &Code{data: make([]Instr, 0, 64)}
}

// Len returns code length.
func (p *Code) Len() int {
	return len(p.data)
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
	NestDepth uint32
}

// NewBuilder creates a new Code Builder instance.
func NewBuilder(code *Code) *Builder {
	if code == nil {
		code = NewCode()
	}
	return &Builder{
		code:      code,
		valConsts: make(map[interface{}]*valUnresolved),
		labels:    make(map[*Label]int),
		funcs:     make(map[*FuncInfo]int),
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
