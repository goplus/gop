package exec

import (
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"
	"sort"
	"strings"

	qlang "qlang.io/spec"
)

// -----------------------------------------------------------------------------
// type Stack

// A Stack represents a FILO container.
//
type Stack struct {
	data []interface{}
}

// NewStack returns a new Stack.
//
func NewStack() *Stack {

	data := make([]interface{}, 0, 16)
	return &Stack{data}
}

// Push pushs a value into this stack.
//
func (p *Stack) Push(v interface{}) {

	p.data = append(p.data, v)
}

// Top returns the last pushed value, if it exists.
//
func (p *Stack) Top() (v interface{}, ok bool) {

	n := len(p.data)
	if n > 0 {
		v, ok = p.data[n-1], true
	}
	return
}

// Pop pops a value from this stack.
//
func (p *Stack) Pop() (v interface{}, ok bool) {

	n := len(p.data)
	if n > 0 {
		v, ok = p.data[n-1], true
		p.data = p.data[:n-1]
	}
	return
}

// PushRet pushs a function call result.
//
func (p *Stack) PushRet(ret []reflect.Value) error {

	switch len(ret) {
	case 0:
		p.Push(nil)
	case 1:
		p.Push(ret[0].Interface())
	default:
		slice := make([]interface{}, len(ret))
		for i, v := range ret {
			slice[i] = v.Interface()
		}
		p.Push(slice)
	}
	return nil
}

// PopArgs pops arguments of a function call.
//
func (p *Stack) PopArgs(arity int) (args []reflect.Value, ok bool) {

	pstk := p.data
	n := len(pstk)
	if n >= arity {
		args, ok = make([]reflect.Value, arity), true
		n -= arity
		for i := 0; i < arity; i++ {
			args[i] = reflect.ValueOf(pstk[n+i])
		}
		p.data = pstk[:n]
	}
	return
}

// PopNArgs pops arguments of a function call.
//
func (p *Stack) PopNArgs(arity int) []interface{} {

	pstk := p.data
	n := len(pstk)
	if n >= arity {
		args := make([]interface{}, arity)
		n -= arity
		for i := 0; i < arity; i++ {
			args[i] = pstk[n+i]
		}
		p.data = pstk[:n]
		return args
	}
	panic("unexpected argument count")
}

// PopFnArgs pops argument names of a function call.
//
func (p *Stack) PopFnArgs(arity int) []string {

	ok := false
	pstk := p.data
	n := len(pstk)
	if n >= arity {
		args := make([]string, arity)
		n -= arity
		for i := 0; i < arity; i++ {
			if args[i], ok = pstk[n+i].(string); !ok {
				panic("function argument isn't a symbol")
			}
		}
		p.data = pstk[:n]
		return args
	}
	panic("unexpected argument count")
}

// BaseFrame returns current stack size.
//
func (p *Stack) BaseFrame() int {

	return len(p.data)
}

// SetFrame sets stack to new size.
//
func (p *Stack) SetFrame(n int) {

	p.data = p.data[:n]
}

// -----------------------------------------------------------------------------

type variable struct {
	Name int
}

type variables struct {
	vars   []interface{}
	symtbl map[string]int // symbol table
}

// initVars initializes the variable table.
//
func (p *variables) initVars(symtbl map[string]int) {

	n := len(symtbl)
	if n > 0 {
		p.vars = make([]interface{}, n)
	}
	p.symtbl = symtbl
}

// ResizeVars is reserved for internal use.
//
func (p *variables) ResizeVars() {

	vars := p.vars
	n := len(p.symtbl)
	i := len(vars)
	if i >= n {
		return
	}
	if i == 0 {
		vars = make([]interface{}, 0, n)
	}
	for ; i < n; i++ {
		vars = append(vars, qlang.Undefined)
	}
	p.vars = vars
}

// Vars is deprecated. please use `CopyVars` method.
//
func (p *variables) Vars() []interface{} {

	return p.vars
}

// CopyVars copies and returns all variables.
//
func (p *variables) CopyVars() map[string]interface{} {

	vars := make(map[string]interface{})
	for k, name := range p.symtbl {
		vars[k] = p.vars[name]
	}
	return vars
}

// ResetVars resets all variables.
//
func (p *variables) ResetVars(vars map[string]interface{}) {

	for k, name := range p.symtbl {
		p.vars[name] = vars[k]
	}
}

// Var returns a variable value.
//
func (p *variables) Var(name string) interface{} {

	if k, ok := p.symtbl[name]; ok {
		return p.vars[k]
	}
	return qlang.Undefined
}

// GetVar returns a variable value.
//
func (p *variables) GetVar(name string) (v interface{}, ok bool) {

	if k, ok := p.symtbl[name]; ok {
		return p.vars[k], true
	}
	return qlang.Undefined, false
}

// SetVar sets a variable value.
//
func (p *variables) SetVar(name string, v interface{}) {

	k, ok := p.symtbl[name]
	if !ok {
		k = len(p.symtbl)
		if k != len(p.vars) {
			panic("ERROR: Variables need to resize (call `ResizeVars` method please)")
		}
		p.symtbl[name] = k
		p.vars = append(p.vars, qlang.Undefined)
	}
	p.vars[k] = v
}

// UnsetVar deletes a variable.
//
func (p *variables) UnsetVar(name string) {

	if k, ok := p.symtbl[name]; ok {
		p.vars[k] = qlang.Undefined
	}
}

// FastGetVar returns a variable value.
//
func (p *variables) FastGetVar(name int) interface{} {

	return p.vars[name]
}

// FastRefVar returns a variable address.
//
func (p *variables) FastRefVar(name int) *interface{} {

	return &p.vars[name]
}

// FastSetVar sets a variable value.
//
func (p *variables) FastSetVar(name int, v interface{}) {

	p.vars[name] = v
}

// -----------------------------------------------------------------------------
// type Context

type theDefer struct {
	next  *theDefer
	start int
	end   int
}

// A Context represents the context of an executor.
//
type Context struct {
	variables
	Stack  *Stack
	Code   *Code
	parent *Context
	defers *theDefer
	modmgr *moduleMgr
	Recov  interface{}
	ret    interface{}
	export []string
	ip     int
	base   int
	onsel  bool // on select
	noextv bool // don't cache extern var
}

// NewContextEx returns a new context of an executor.
//
func NewContextEx(symtbl map[string]int) *Context {

	mods := make(map[string]*importMod)
	modmgr := &moduleMgr{
		mods: mods,
	}
	p := &Context{modmgr: modmgr}
	p.initVars(symtbl)
	return p
}

// Exports returns a module exports.
//
func (p *Context) Exports() map[string]interface{} {

	export := make(map[string]interface{}, len(p.export))
	for _, name := range p.export {
		export[name], _ = p.GetVar(name)
	}
	return export
}

// ExecBlock executes an anonym function.
//
func (p *Context) ExecBlock(ip, ipEnd int, symtbl map[string]int) {

	mod := NewFunction(nil, ip, ipEnd, symtbl, nil, false)
	mod.ExtCall(p)
}

// ExecDefers executes defer blocks.
//
func (p *Context) ExecDefers() {

	d := p.defers
	if d == nil {
		return
	}

	p.defers = nil
	code := p.Code
	stk := p.Stack
	for {
		code.Exec(d.start, d.end, stk, p)
		d = d.next
		if d == nil {
			break
		}
	}
}

// -----------------------------------------------------------------------------

// A Error represents a qlang runtime error.
//
type Error struct {
	Err   error
	File  string
	Line  int
	Stack []byte
}

func (p *Error) Error() string {

	var sep string
	var stk []byte
	if qlang.DumpStack {
		stk = p.Stack
		sep = "\n\n"
	}

	if p.Line == 0 {
		return fmt.Sprintf("%v%s%s", p.Err, sep, stk)
	}
	if p.File == "" {
		return fmt.Sprintf("line %d: %v%s%s", p.Line, p.Err, sep, stk)
	}
	return fmt.Sprintf("%s:%d: %v%s%s", p.File, p.Line, p.Err, sep, stk)
}

// -----------------------------------------------------------------------------
// type Code

// A Instr represents a instruction of the executor.
//
type Instr interface {
	Exec(stk *Stack, ctx *Context)
}

// RefToVar converts a value reference instruction into a assignable variable instruction.
//
type RefToVar interface {
	ToVar() Instr
}

type optimizableArityGetter interface {
	OptimizableGetArity() int
}

type ipFileLine struct {
	ip   int
	line int
	file string
}

// A Code represents generated instructions to execute.
//
type Code struct {
	data  []Instr
	lines []*ipFileLine // ip => (file,line)
}

// A ReservedInstr represents a reserved instruction to be assigned.
//
type ReservedInstr struct {
	code *Code
	idx  int
}

// New returns a new Code object.
//
func New(data ...Instr) *Code {

	return &Code{data, nil}
}

// CodeLine informs current file and line.
//
func (p *Code) CodeLine(file string, line int) {

	p.lines = append(p.lines, &ipFileLine{ip: len(p.data), file: file, line: line})
}

// Line returns file line of a instruction position.
//
func (p *Code) Line(ip int) (file string, line int) {

	idx := sort.Search(len(p.lines), func(i int) bool {
		return ip < p.lines[i].ip
	})
	if idx < len(p.lines) {
		t := p.lines[idx]
		return t.file, t.line
	}
	return "", 0
}

// Len returns code length.
//
func (p *Code) Len() int {

	return len(p.data)
}

// Reserve reserves an instruction and returns it.
//
func (p *Code) Reserve() ReservedInstr {

	idx := len(p.data)
	p.data = append(p.data, nil)
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

// CheckConst returns the value, if code[ip] is a const instruction.
//
func (p *Code) CheckConst(ip int) (v interface{}, ok bool) {

	if instr, ok := p.data[ip].(*iPush); ok {
		return instr.v, true
	}
	return
}

func appendInstrOptimized(data []Instr, instr Instr, arity int) []Instr {

	n := len(data)
	base := n - arity
	for i := base; i < n; i++ {
		if _, ok := data[i].(*iPush); !ok {
			return append(data, instr)
		}
	}
	args := make([]interface{}, arity)
	for i := base; i < n; i++ {
		args[i-base] = data[i].(*iPush).v
	}
	stk := &Stack{data: args}
	instr.Exec(stk, nil)
	return append(data[:base], Push(stk.data[0]))
}

// Block appends some instructions to code.
//
func (p *Code) Block(code ...Instr) int {

	for _, instr := range code {
		if g, ok := instr.(optimizableArityGetter); ok {
			arity := g.OptimizableGetArity()
			p.data = appendInstrOptimized(p.data, instr, arity)
		} else {
			p.data = append(p.data, instr)
		}
	}
	return len(p.data)
}

// ToVar converts the last instruction from ref to var.
//
func (p *Code) ToVar() {

	data := p.data
	idx := len(data) - 1
	if cvt, ok := data[idx].(RefToVar); ok {
		data[idx] = cvt.ToVar()
	} else {
		panic("expr is not assignable")
	}
}

// Exec executes a code block from ip to ipEnd.
//
func (p *Code) Exec(ip, ipEnd int, stk *Stack, ctx *Context) {

	defer func() {
		if e := recover(); e != nil {
			if e == ErrReturn {
				panic(e)
			}
			if err, ok := e.(*Error); ok {
				panic(err)
			}
			err, ok := e.(error)
			if !ok {
				if s, ok := e.(string); ok {
					err = errors.New(s)
				} else {
					panic(e)
				}
			}
			file, line := p.Line(ctx.ip - 1)
			err = &Error{
				Err:   err,
				File:  file,
				Line:  line,
				Stack: debug.Stack(),
			}
			panic(err)
		}
	}()

	ctx.ip = ip
	data := p.data
	for ctx.ip != ipEnd {
		instr := data[ctx.ip]
		ctx.ip++
		instr.Exec(stk, ctx)
	}
}

// Dump dumps code instructions within a range.
//
func (p *Code) Dump(ranges ...int) {

	start := 0
	end := len(p.data)
	if len(ranges) > 0 {
		start = ranges[0]
		if len(ranges) > 1 {
			end = ranges[1]
		}
	}
	for i, instr := range p.data[start:end] {
		fmt.Printf("==> %04d: %s %v\n", i+start, instrName(instr), instr)
	}
}

func instrName(instr interface{}) string {

	if instr == nil {
		return "<nil>"
	}
	t := reflect.TypeOf(instr).String()
	if strings.HasPrefix(t, "*") {
		t = t[1:]
	}
	if strings.HasPrefix(t, "exec.") {
		t = t[5:]
	}
	if strings.HasPrefix(t, "i") {
		t = t[1:]
	}
	return t
}

// -----------------------------------------------------------------------------
