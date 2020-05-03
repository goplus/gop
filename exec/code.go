package exec

// -----------------------------------------------------------------------------

const (
	bitsInstr       = 32
	bitsOp          = 6
	bitsIntKind     = 3
	bitsFloatKind   = 2
	bitsGoFunvArity = 10

	bitsOpShift = bitsInstr - bitsOp
	bitsOperand = (1 << bitsOpShift) - 1

	bitsOpInt        = bitsOp + bitsIntKind
	bitsOpIntShift   = bitsInstr - bitsOpInt
	bitsOpIntOperand = (1 << bitsOpIntShift) - 1

	bitsOpFloat        = bitsOp + bitsFloatKind
	bitsOpFloatShift   = bitsInstr - bitsOpFloat
	bitsOpFloatOperand = (1 << bitsOpFloatShift) - 1

	bitsOpCallGoFunv        = bitsOp + bitsGoFunvArity
	bitsOpCallGoFunvShift   = bitsInstr - bitsOpCallGoFunv
	bitsOpCallGoFunvOperand = (1 << bitsOpCallGoFunvShift) - 1
)

// A Instr represents a instruction of the executor.
type Instr = uint32

const (
	opInvalid     = (1 << bitsOp) - 1
	opPushInt     = 0
	opPushUint    = 1
	opPushFloat   = 2
	opPushValSpec = 3
	opPushStringR = 4
	opPushIntR    = 5
	opPushUintR   = 6
	opPushFloatR  = 7
	opBuiltinOp   = 8
	opJmp         = 9
	opJmpIfFalse  = 10
	opCaseNE      = 11
	opPop         = 12
	opCallGoFun   = 13 // call a Go function
	opCallGoFunv  = 14 // call a Go function with variadic args
)

const (
	iInvalid        = (opInvalid << bitsOpShift)
	iPushFalse      = (opPushValSpec << bitsOpShift)
	iPushTrue       = (opPushValSpec << bitsOpShift) | 1
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
	}
}

// Resolve resolves all unresolved labels/functions/consts/etc.
func (p *Builder) Resolve() *Code {
	p.resolveLabels()
	p.resolveConsts()
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
