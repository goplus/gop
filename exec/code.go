package exec

// -----------------------------------------------------------------------------

const (
	bitsOp        = 6
	bitsIntKind   = 3
	bitsFloatKind = 2

	bitsOpShift = 32 - bitsOp
	bitsOperand = (1 << bitsOpShift) - 1

	bitsOpIntShift   = 32 - (bitsOp + bitsIntKind)
	bitsOpIntOperand = (1 << bitsOpIntShift) - 1

	bitsOpFloatShift   = 32 - (bitsOp + bitsFloatKind)
	bitsOpFloatOperand = (1 << bitsOpFloatShift) - 1
)

// A Instr represents a instruction of the executor.
//
type Instr = uint32

const (
	opInvalid     = (1 << bitsOp) - 1
	opPushInt     = 0x00
	opPushUint    = 0x01
	opPushFloat   = 0x02
	opPushValSpec = 0x03
	opPushStringR = 0x04
	opPushIntR    = 0x05
	opPushUintR   = 0x06
	opPushFloatR  = 0x07
	opBuiltinOp   = 0x08
)

const (
	iPushFalse      = (opPushValSpec << bitsOpShift)
	iPushTrue       = (opPushValSpec << bitsOpShift) | 1
	iPushUnresolved = (opInvalid << bitsOpShift)
)

// -----------------------------------------------------------------------------

// A Code represents generated instructions to execute.
//
type Code struct {
	data         []Instr
	stringConsts []string
	intConsts    []int64
	uintConsts   []uint64
	valConsts    []interface{}
}

// A ReservedInstr represents a reserved instruction to be assigned.
//
type ReservedInstr struct {
	code *Code
	idx  int
}

// Reserve reserves an instruction and returns it.
//
func (p *Code) Reserve() ReservedInstr {
	idx := len(p.data)
	p.data = append(p.data, 0)
	return ReservedInstr{p, idx}
}

// Set sets a reserved instruction.
//
func (p ReservedInstr) Set(code Instr) {
	p.code.data[p.idx] = code
}

// Next returns next instruction position.
//
func (p ReservedInstr) Next() int {
	return p.idx + 1
}

// Delta returns distance from b to p.
//
func (p ReservedInstr) Delta(b ReservedInstr) int {
	return p.idx - b.idx
}

// NewCode returns a new Code object.
//
func NewCode() *Code {

	return &Code{data: make([]Instr, 0, 64)}
}

// Len returns code length.
//
func (p *Code) Len() int {
	return len(p.data)
}

// -----------------------------------------------------------------------------

type valUnresolved struct {
	op   Instr
	offs []int
}

// Builder class.
type Builder struct {
	code      *Code
	valConsts map[interface{}]*valUnresolved
}

// NewBuilder creates a new Code Builder instance.
func NewBuilder(code *Code) *Builder {
	if code == nil {
		code = NewCode()
	}
	return &Builder{
		code:      code,
		valConsts: make(map[interface{}]*valUnresolved),
	}
}

// Resolve resolves all unresolved consts/functions/etc.
func (ctx *Builder) Resolve() (code *Code) {
	code = ctx.code
	resolveConsts(ctx, code)
	return
}

// -----------------------------------------------------------------------------
