package qlang

import (
	"errors"
	"reflect"
	"strconv"

	"qlang.io/qlang.spec.v1"
	"qiniupkg.com/text/tpl.v1/interpreter.util"
)

// -----------------------------------------------------------------------------

const Grammar = `

term1 = factor *('*' factor/mul | '/' factor/quo | '%' factor/mod)

term2 = term1 *('+' term1/add | '-' term1/sub)

term3 = term2 *('<' term2/lt | '>' term2/gt | "==" term2/eq | "<=" term2/le | ">=" term2/ge | "!=" term2/ne)

term4 = term3 *("&&"/_mute term3/_code/_unmute/and)

expr = term4 *("||"/_mute term4/_code/_unmute/or)

s =
	(IDENT '='! expr)/assign |
	(IDENT ','!)/name IDENT/name % ','/ARITY '=' expr /massign |
	(IDENT "++")/inc |
	(IDENT "--")/dec |
	(IDENT "+="! expr)/adda |
	(IDENT "-="! expr)/suba |
	(IDENT "*="! expr)/mula |
	(IDENT "/="! expr)/quoa |
	(IDENT "%="! expr)/moda |
	"return" ?expr/ARITY /return |
	"defer"/_mute! expr/_code/_unmute/defer |
	expr

doc = s *(';'/clear s | ';'/pushn)

ifbody = '{' ?doc/_code '}'

swbody = *("case"! expr/_code ':' ?doc/_code)/_ARITY ?("default"! ':' ?doc/_code)/_ARITY

fnbody = '(' IDENT/name %= ','/ARITY ?"..."/ARITY ')' '{'/_mute ?doc/_code '}'/_unmute

clsname = '(' IDENT/ref ')' | IDENT/ref

newargs = ?('(' expr %= ','/ARITY ')')/ARITY

classb = "fn"! IDENT/name fnbody ?';'/mfn

atom =
	'(' expr %= ','/ARITY ?"..."/ARITY ?',' ')'/call |
	'.' (IDENT|"class"|"new"|"recover"|"main")/mref |
	'[' ?expr/ARITY ?':'/ARITY ?expr/ARITY ']'/index

factor =
	INT/pushi |
	FLOAT/pushf |
	STRING/pushs |
	CHAR/pushc |
	(IDENT/ref | '('! expr ')' | "fn"! fnbody/fn | '[' expr %= ','/ARITY ?',' ']'/slice) *atom |
	"if"/_mute! expr/_code ifbody *("elif" expr/_code ifbody)/_ARITY ?("else" ifbody)/_ARITY/_unmute/if |
	"switch"/_mute! ?(~'{' expr)/_code '{' swbody '}'/_unmute/switch |
	"for"/_mute! (~'{' s)/_code %= ';'/_ARITY '{' ?doc/_code '}'/_unmute/for |
	"new"! clsname newargs /new |
	"class"! '{' *classb/ARITY '}'/class |
	"recover" '(' ')'/recover |
	"main" '{'/_mute ?doc/_code '}'/_unmute/main |
	'{'! (expr ':' expr) %= ','/ARITY ?',' '}'/map |
	'!' factor/not |
	'-' factor/neg |
	'+' factor
`

var (
	ErrRefWithoutObject = errors.New("reference without object")
	ErrAssignWithoutVal = errors.New("variable assign without value")
)

// -----------------------------------------------------------------------------

type Stack []interface{}

func NewStack() *Stack {

	stk := make(Stack, 0, 16)
	return &stk
}

func (stk *Stack) Push(v interface{}) {

	*stk = append(*stk, v)
}

func (stk *Stack) PushInt(v int) {

	*stk = append(*stk, v)
}

func (stk *Stack) PushFloat(v float64) {

	*stk = append(*stk, v)
}

func (stk *Stack) PushByte(lit string) {

	v, multibyte, tail, err := strconv.UnquoteChar(lit[1:len(lit)-1], '\'')
	if err != nil {
		panic("invalid char `" + lit + "`: " + err.Error())
	}
	if tail != "" || multibyte {
		panic("invalid char: " + lit)
	}
	*stk = append(*stk, byte(v))
}

func (stk *Stack) PushString(lit string) {

	v, err := strconv.Unquote(lit)
	if err != nil {
		panic("invalid string `" + lit + "`: " + err.Error())
	}
	*stk = append(*stk, v)
}

func (stk *Stack) Pop() (v interface{}, ok bool) {

	n := len(*stk)
	if n > 0 {
		v, ok = (*stk)[n-1], true
		*stk = (*stk)[:n-1]
	}
	return
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

func (stk *Stack) PushRet(ret []reflect.Value) error {

	switch len(ret) {
	case 0:
		stk.Push(nil)
	case 1:
		stk.Push(castVal(ret[0]))
	default:
		slice := make([]interface{}, len(ret))
		for i, v := range ret {
			slice[i] = castVal(v)
		}
		stk.Push(slice)
	}
	return nil
}

func (stk *Stack) PopArgs(arity int) (args []reflect.Value, ok bool) {

	pstk := *stk
	n := len(pstk)
	if n >= arity {
		args, ok = make([]reflect.Value, arity), true
		n -= arity
		for i := 0; i < arity; i++ {
			args[i] = reflect.ValueOf(pstk[n+i])
		}
		*stk = pstk[:n]
	}
	return
}

func (stk *Stack) popNArgs(arity int) []interface{} {

	pstk := *stk
	n := len(pstk)
	if n >= arity {
		args := make([]interface{}, arity)
		n -= arity
		for i := 0; i < arity; i++ {
			args[i] = pstk[n+i]
		}
		*stk = pstk[:n]
		return args
	}
	panic("unexpected argument count")
}

func (stk *Stack) popFnArgs(arity int) []string {

	ok := false
	pstk := *stk
	n := len(pstk)
	if n >= arity {
		args := make([]string, arity)
		n -= arity
		for i := 0; i < arity; i++ {
			if args[i], ok = pstk[n+i].(string); !ok {
				panic("function argument isn't a symbol?")
			}
		}
		*stk = pstk[:n]
		return args
	}
	panic("unexpected argument count")
}

func (stk *Stack) popArity() int {

	if v, ok := stk.Pop(); ok {
		if arity, ok := v.(int); ok {
			return arity
		}
	}
	panic("no arity")
}

func PushNil(stk *Stack) {

	stk.Push(nil)
}

func Arity(stk *Stack, arity int) {

	stk.Push(arity)
}

func Code(stk *Stack, code interface{}) {

	stk.Push(code)
}

func Index(stk *Stack) {

	var arg1, arg2 interface{}

	arity2 := stk.popArity()
	if arity2 != 0 {
		arg2, _ = stk.Pop()
	}
	arityMid := stk.popArity()
	arity1 := stk.popArity()
	args := stk.popNArgs(arity1 + 1)

	if arityMid == 0 {
		if arity1 == 0 {
			panic("call operator[] without index")
		}
		stk.Push(qlang.Get(args[0], args[1]))
		return
	}
	if arity1 != 0 {
		arg1 = args[1]
	}
	stk.Push(qlang.SubSlice(args[0], arg1, arg2))
}

// -----------------------------------------------------------------------------

type Defer struct {
	next *Defer
	src  interface{}
}

type Interpreter struct {
	vars   map[string]interface{}
	ret    interface{}
	Recov  interface{}
	parent *Interpreter
	stk    *Stack
	defers *Defer
	base   int
}

func New() *Interpreter {

	stk := NewStack()
	vars := make(map[string]interface{})
	return &Interpreter{vars: vars, stk: stk}
}

func (p *Interpreter) Grammar() string {

	return Grammar
}

func (p *Interpreter) Fntable() map[string]interface{} {

	return qlang.Fntable
}

func (p *Interpreter) Vars() map[string]interface{} {

	return p.vars
}

func (p *Interpreter) Var(name string) (v interface{}, ok bool) {

	v, ok = p.vars[name]
	return
}

func (p *Interpreter) Stack() interpreter.Stack {

	return p.stk
}

func (p *Interpreter) Ret() (v interface{}, ok bool) {

	v, ok = p.stk.Pop()
	p.Clear()
	return
}

func (p *Interpreter) VMap() error {

	arity := p.stk.popArity()
	return interpreter.Call(p.stk, qlang.MapFrom, arity*2)
}

func (p *Interpreter) VSlice() error {

	arity := p.stk.popArity()
	return interpreter.Call(p.stk, qlang.SliceFrom, arity)
}

func (p *Interpreter) VCall() error {

	stk := p.stk
	variadic := stk.popArity()
	arity := stk.popArity()
	if variadic != 0 {
		if arity == 0 {
			panic("what do you mean of `...`?")
		}
		if val, ok := stk.Pop(); ok {
			v := reflect.ValueOf(val)
			if v.Kind() != reflect.Slice {
				panic("apply `...` on non-slice object")
			}
			n := v.Len()
			for i := 0; i < n; i++ {
				stk.Push(v.Index(i).Interface())
			}
			arity += (n - 1)
		} else {
			panic("unexpected")
		}
	}
	return interpreter.VCall(stk, arity)
}

// -----------------------------------------------------------------------------

var exports = map[string]interface{}{
	"$ARITY":   Arity,
	"$_ARITY":  Arity,
	"$_code":   Code,
	"$index":   Index,
	"$name":    PushName,
	"$pushn":   PushNil,
	"$pushi":   (*Stack).PushInt,
	"$pushf":   (*Stack).PushFloat,
	"$pushs":   (*Stack).PushString,
	"$pushc":   (*Stack).PushByte,
	"$mref":    (*Interpreter).MemberRef,
	"$ref":     (*Interpreter).Ref,
	"$slice":   (*Interpreter).VSlice,
	"$map":     (*Interpreter).VMap,
	"$call":    (*Interpreter).VCall,
	"$assign":  (*Interpreter).Assign,
	"$massign": (*Interpreter).MultiAssign,
	"$inc":     (*Interpreter).Inc,
	"$dec":     (*Interpreter).Dec,
	"$adda":    (*Interpreter).AddAssign,
	"$suba":    (*Interpreter).SubAssign,
	"$mula":    (*Interpreter).MulAssign,
	"$quoa":    (*Interpreter).QuoAssign,
	"$moda":    (*Interpreter).ModAssign,
	"$defer":   (*Interpreter).Defer,
	"$recover": (*Interpreter).Recover,
	"$return":  (*Interpreter).Return,
	"$fn":      (*Interpreter).Function,
	"$main":    (*Interpreter).Main,
	"$mfn":     (*Interpreter).MemberFuncDecl,
	"$class":   (*Interpreter).Class,
	"$new":     (*Interpreter).New,
	"$clear":   (*Interpreter).Clear,
	"$if":      (*Interpreter).If,
	"$switch":  (*Interpreter).Switch,
	"$for":     (*Interpreter).For,
	"$and":     (*Interpreter).And,
	"$or":      (*Interpreter).Or,
}

func init() {
	qlang.Import("", exports)
}

// -----------------------------------------------------------------------------
