package exec

import (
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
	"sort"

	"qlang.io/qlang.spec.v1"
)

// -----------------------------------------------------------------------------
// type Stack

type Stack struct {
	data []interface{}
}

func NewStack() *Stack {

	data := make([]interface{}, 0, 16)
	return &Stack{data}
}

func (p *Stack) Push(v interface{}) {

	p.data = append(p.data, v)
}

func (p *Stack) Top() (v interface{}, ok bool) {

	n := len(p.data)
	if n > 0 {
		v, ok = p.data[n-1], true
	}
	return
}

func (p *Stack) Pop() (v interface{}, ok bool) {

	n := len(p.data)
	if n > 0 {
		v, ok = p.data[n-1], true
		p.data = p.data[:n-1]
	}
	return
}

func (p *Stack) PushRet(ret []reflect.Value) error {

	switch len(ret) {
	case 0:
		p.Push(nil)
	case 1:
		p.Push(castVal(ret[0]))
	default:
		slice := make([]interface{}, len(ret))
		for i, v := range ret {
			slice[i] = castVal(v)
		}
		p.Push(slice)
	}
	return nil
}

func castVal(v reflect.Value) interface{} {

	switch kind := v.Kind(); {
	case kind == reflect.Float64 || kind == reflect.Int:
		return v.Interface()
	case kind > reflect.Int && kind <= reflect.Int64:
		if kind != reflect.Uint8 {
			return int(v.Int())
		} else {
			return v.Interface()
		}
	case kind >= reflect.Uint && kind <= reflect.Uintptr:
		return int(v.Uint())
	case kind == reflect.Float32:
		return v.Float()
	default:
		return v.Interface()
	}
}

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

func (p *Stack) PopFnArgs(arity int) []string {

	ok := false
	pstk := p.data
	n := len(pstk)
	if n >= arity {
		args := make([]string, arity)
		n -= arity
		for i := 0; i < arity; i++ {
			if args[i], ok = pstk[n+i].(string); !ok {
				panic("function argument isn't a symbol?")
			}
		}
		p.data = pstk[:n]
		return args
	}
	panic("unexpected argument count")
}

func (p *Stack) BaseFrame() int {

	return len(p.data)
}

func (p *Stack) SetFrame(n int) {

	p.data = p.data[:n]
}

// -----------------------------------------------------------------------------
// type Context

type theDefer struct {
	next  *theDefer
	start int
	end   int
}

type Context struct {
	parent *Context
	defers *theDefer
	modmgr *moduleMgr
	Stack  *Stack
	Code   *Code
	Recov  interface{}
	ret    interface{}
	vars   map[string]interface{}
	export []string
	ip     int
	base   int
	onsel  bool // on select
}

func NewContext() *Context {

	vars := make(map[string]interface{})
	mods := make(map[string]*importMod)
	modmgr := &moduleMgr{
		mods: mods,
	}
	return &Context{vars: vars, modmgr: modmgr}
}

func (p *Context) Exports() map[string]interface{} {

	export := make(map[string]interface{}, len(p.export))
	vars := p.vars
	for _, name := range p.export {
		export[name] = vars[name]
	}
	return export
}

func (p *Context) Vars() map[string]interface{} {

	return p.vars
}

func (p *Context) Var(name string) (v interface{}, ok bool) {

	v, ok = p.vars[name]
	return
}

func (p *Context) SetVar(name string, v interface{}) {

	p.vars[name] = v
}

func (p *Context) Unset(name string) {

	delete(p.vars, name)
}

func (p *Context) ExecBlock(ip, ipEnd int) {

	mod := NewFunction(nil, ip, ipEnd, nil, false)
	mod.ExtCall(p)
}

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

type Error struct {
	Err   error
	File  string
	Line  int
	Stack []byte
}

func (p *Error) Error() string {

	var stk []byte
	if qlang.DumpStack {
		stk = p.Stack
	}

	if p.Line == 0 {
		return fmt.Sprintf("%v\n\n%s", p.Err, stk)
	}
	if p.File == "" {
		return fmt.Sprintf("line %d: %v\n\n%s", p.Line, p.Err, stk)
	}
	return fmt.Sprintf("%s:%d: %v\n\n%s", p.File, p.Line, p.Err, stk)
}

// -----------------------------------------------------------------------------
// type Code

type Instr interface {
	Exec(stk *Stack, ctx *Context)
}

type ipFileLine struct {
	ip   int
	line int
	file string
}

type Code struct {
	data  []Instr
	lines []*ipFileLine // ip => (file,line)
}

type ReservedInstr struct {
	code *Code
	idx  int
}

func New(data ...Instr) *Code {

	return &Code{data, nil}
}

func (p *Code) CodeLine(file string, line int) {

	p.lines = append(p.lines, &ipFileLine{ip: len(p.data), file: file, line: line})
}

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

func (p *Code) Len() int {

	return len(p.data)
}

func (p *Code) Reserve() ReservedInstr {

	idx := len(p.data)
	p.data = append(p.data, nil)
	return ReservedInstr{p, idx}
}

func (p ReservedInstr) Set(code Instr) {

	p.code.data[p.idx] = code
}

func (p ReservedInstr) Next() int {

	return p.idx + 1
}

func (a ReservedInstr) Delta(b ReservedInstr) int {

	return a.idx - b.idx
}

func (p *Code) Block(code ...Instr) int {

	p.data = append(p.data, code...)
	return len(p.data)
}

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
				Err:  err,
				File: file,
				Line: line,
				Stack: debug.Stack(),
			}
			panic(err)
		}
	}()

	ctx.ip = ip
	for ctx.ip != ipEnd {
		instr := p.data[ctx.ip]
		ctx.ip++
		instr.Exec(stk, ctx)
	}
}

func (p *Code) Dump() {

	for i, instr := range p.data {
		fmt.Printf("==> %04d: %s %v\n", i, instrName(instr), instr)
	}
}

func instrName(instr interface{}) string {

	if instr == nil {
		return "<nil>"
	}
	t := reflect.TypeOf(instr).String()
	if strings.HasPrefix(t, "*")  {
		t = t[1:]
	}
	if strings.HasPrefix(t, "exec.i") {
		t = t[6:]
	}
	return t
}

// -----------------------------------------------------------------------------
