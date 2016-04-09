package qlang

import (
	"strings"

	"qlang.io/exec.v2"
	"qlang.io/qlang.spec.v1"
	"qiniupkg.com/text/tpl.v1/interpreter.util"

	ipt "qiniupkg.com/text/tpl.v1/interpreter"
)

// -----------------------------------------------------------------------------

const Grammar = `

term1 = factor *('*' factor/mul | '/' factor/quo | '%' factor/mod)

term2 = term1 *('+' term1/add | '-' term1/sub)

term31 = term2 *('<' term2/lt | '>' term2/gt | "==" term2/eq | "<=" term2/le | ">=" term2/ge | "!=" term2/ne)

term3 = term31 *("<-" term31/chin)

term4 = term3 *("&&"/_mute term3/_code/_unmute/and)

expr = term4 *("||"/_mute term4/_code/_unmute/or)

s = (
	(IDENT '='! expr)/assign |
	(IDENT ','!)/name IDENT/name % ','/ARITY '=' expr /massign |
	(IDENT "++")/inc |
	(IDENT "--")/dec |
	(IDENT "+="! expr)/adda |
	(IDENT "-="! expr)/suba |
	(IDENT "*="! expr)/mula |
	(IDENT "/="! expr)/quoa |
	(IDENT "%="! expr)/moda |
	"if"/_mute! expr/_code body *("elif" expr/_code body)/_ARITY ?("else" body)/_ARITY/_unmute/if |
	"switch"/_mute! ?(~'{' expr)/_code '{' swbody '}'/_unmute/switch |
	"for"/_mute! fhead body/_unmute/for |
	"return"! ?expr/ARITY /return |
	"break" /brk |
	"continue" /cont |
	"include"! STRING/include |
	"import"! (STRING ?("as" IDENT/name)/ARITY)/import |
	"export"! IDENT/name % ','/ARITY /export |
	"defer"/_mute! expr/_code/_unmute/defer |
	"go"/_mute! expr/_code/_unmute/go |
	(expr/pop))/xline

doc = ?(s *(';' ?s))

body = '{' doc/_code '}'

fhead = (~'{' s)/_code %= ';'/_ARITY

swbody = *("case"! expr/_code ':' doc/_code)/_ARITY ?("default"! ':' doc/_code)/_ARITY

fnbody = '(' IDENT/name %= ','/ARITY ?"..."/ARITY ')' '{'/_mute doc/_code '}'/_unmute

afn = '{'/_mute doc/_code '}'/_unmute/afn

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
	(IDENT/ref | '('! expr ')' |
	"fn"! (~'{' fnbody/fn | afn) | '[' expr %= ','/ARITY ?',' ']'/slice) *atom |
	"new"! clsname newargs /new |
	"class"! '{' *classb/ARITY '}'/class |
	"recover"! '(' ')'/recover |
	"main"! afn |
	'{'! (expr ':' expr) %= ','/ARITY ?',' '}'/map |
	'!' factor/not |
	'-' factor/neg |
	"<-" factor/chout |
	'+' factor
`

// -----------------------------------------------------------------------------

type module struct {
	start, end int
}

type instrNode struct {
	prev  *instrNode
	instr exec.ReservedInstr
}

func (p *instrNode) JmpTo(where int) {
	for p != nil {
		instr := exec.Jmp(where - p.instr.Next())
		p.instr.Set(instr)
		p = p.prev
	}
}

func (p *instrNode) MergeTo(parent *instrNode) *instrNode {
	for p != nil {
		p, p.prev, parent = p.prev, parent, p
	}
	return parent
}

type blockCtx struct {
	brks  *instrNode
	conts *instrNode
}

func (p *blockCtx) MergeTo(parent *blockCtx) {
	parent.brks = p.brks.MergeTo(parent.brks)
	parent.conts = p.conts.MergeTo(parent.conts)
}

func (p *blockCtx) MergeSw(old *blockCtx, done int) {
	old.conts = p.conts.MergeTo(old.conts)
	p.brks.JmpTo(done)
	*p = *old
}

type Compiler struct {
	Opts  *ipt.Options
	code  *exec.Code
	libs  []string
	exits []func()
	mods  map[string]module
	gvars map[string]interface{}
	gstk  exec.Stack
	bctx  blockCtx
}

func New() *Compiler {

	gvars := make(map[string]interface{})
	mods := make(map[string]module)
	return &Compiler{
		code:  exec.New(),
		mods:  mods,
		gvars: gvars,
		Opts:  ipt.InsertSemis,
	}
}

func (p *Compiler) SetLibs(libs string) {

	p.libs = strings.Split(libs, ":")
}

func (p *Compiler) Vars() map[string]interface{} {

	return p.gvars
}

func (p *Compiler) Code() *exec.Code {

	return p.code
}

func (p *Compiler) Grammar() string {

	return Grammar
}

func (p *Compiler) Fntable() map[string]interface{} {

	return qlang.Fntable
}

func (p *Compiler) Stack() interpreter.Stack {

	return nil
}

func (p *Compiler) VMap() {

	arity := p.popArity()
	p.code.Block(exec.Call(qlang.MapFrom, arity*2))
}

func (p *Compiler) VSlice() {

	arity := p.popArity()
	p.code.Block(exec.Call(qlang.SliceFrom, arity))
}

func (p *Compiler) VCall() {

	variadic := p.popArity()
	arity := p.popArity()
	if variadic != 0 {
		if arity == 0 {
			panic("what do you mean of `...`?")
		}
		p.code.Block(exec.CallFnv(arity))
	} else {
		p.code.Block(exec.CallFn(arity))
	}
}

func (p *Compiler) Index() {

	arity2 := p.popArity()
	arityMid := p.popArity()
	arity1 := p.popArity()

	if arityMid == 0 {
		if arity1 == 0 {
			panic("call operator[] without index")
		}
		p.code.Block(exec.Call(qlang.Get, 2))
	} else {
		p.code.Block(exec.Op3(qlang.SubSlice, arity1 != 0, arity2 != 0))
	}
}

func (p *Compiler) Pop() {

	p.code.Block(exec.PopEx())
}

func (p *Compiler) CodeLine(f *interpreter.FileLine) {

	p.code.CodeLine(f.File, f.Line)
}

func (p *Compiler) CallFn(fn interface{}) {

	p.code.Block(exec.Call(fn))
}

// -----------------------------------------------------------------------------

var exports = map[string]interface{}{
	"$ARITY":   (*Compiler).Arity,
	"$_ARITY":  (*Compiler).Arity,
	"$_code":   (*Compiler).PushCode,
	"$name":    (*Compiler).PushName,
	"$pushi":   (*Compiler).PushInt,
	"$pushf":   (*Compiler).PushFloat,
	"$pushs":   (*Compiler).PushString,
	"$pushc":   (*Compiler).PushByte,
	"$index":   (*Compiler).Index,
	"$mref":    (*Compiler).MemberRef,
	"$ref":     (*Compiler).Ref,
	"$slice":   (*Compiler).VSlice,
	"$map":     (*Compiler).VMap,
	"$call":    (*Compiler).VCall,
	"$assign":  (*Compiler).Assign,
	"$massign": (*Compiler).MultiAssign,
	"$inc":     (*Compiler).Inc,
	"$dec":     (*Compiler).Dec,
	"$adda":    (*Compiler).AddAssign,
	"$suba":    (*Compiler).SubAssign,
	"$mula":    (*Compiler).MulAssign,
	"$quoa":    (*Compiler).QuoAssign,
	"$moda":    (*Compiler).ModAssign,
	"$defer":   (*Compiler).Defer,
	"$go":      (*Compiler).Go,
	"$chin":    (*Compiler).ChanIn,
	"$chout":   (*Compiler).ChanOut,
	"$recover": (*Compiler).Recover,
	"$return":  (*Compiler).Return,
	"$fn":      (*Compiler).Function,
	"$afn":     (*Compiler).AnonymFn,
	"$include": (*Compiler).Include,
	"$import":  (*Compiler).Import,
	"$export":  (*Compiler).Export,
	"$mfn":     (*Compiler).MemberFuncDecl,
	"$class":   (*Compiler).Class,
	"$new":     (*Compiler).New,
	"$clear":   (*Compiler).Clear,
	"$if":      (*Compiler).If,
	"$switch":  (*Compiler).Switch,
	"$for":     (*Compiler).For,
	"$brk":     (*Compiler).Break,
	"$cont":    (*Compiler).Continue,
	"$and":     (*Compiler).And,
	"$or":      (*Compiler).Or,
	"$pop":     (*Compiler).Pop,
	"$xline":   (*Compiler).CodeLine,
}

func init() {
	qlang.Import("", exports)
}

// -----------------------------------------------------------------------------
