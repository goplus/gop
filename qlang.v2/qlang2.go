package qlang

import (
	"strings"

	"qiniupkg.com/text/tpl.v1/interpreter.util"
	"qlang.io/exec.v2"
	"qlang.io/qlang.spec.v1"

	ipt "qiniupkg.com/text/tpl.v1/interpreter"
)

// -----------------------------------------------------------------------------

const grammar = `

term1 = factor *(
	'*' factor/mul | '/' factor/quo | '%' factor/mod |
	"<<" factor/lshr | ">>" factor/rshr | '&' factor/bitand | "&^" factor/andnot)

term2 = term1 *('+' term1/add | '-' term1/sub | '|' term1/bitor | '^' term1/xor)

term31 = term2 *('<' term2/lt | '>' term2/gt | "==" term2/eq | "<=" term2/le | ">=" term2/ge | "!=" term2/ne)

term3 = term31 *("<-" term31/chin)

term4 = term3 *("&&"/_mute term3/_code/_unmute/and)

expr = term4 *("||"/_mute term4/_code/_unmute/or)

sexpr = expr (
	'='/tovar! expr/assign |
	','/tovar! expr/tovar % ','/ARITY '=' expr % ','/ARITY /massign |
	"++"/tovar/inc | "--"/tovar/dec |
	"+="/tovar! expr/adda | "-="/tovar! expr/suba |
	"*="/tovar! expr/mula | "/="/tovar! expr/quoa | "%="/tovar! expr/moda |
	"^="/tovar! expr/xora | "<<="/tovar! expr/lshra | ">>="/tovar! expr/rshra |
	"&="/tovar! expr/bitanda | "|="/tovar! expr/bitora | "&^="/tovar! expr/andnota | 1/pop)

s = (
	"if"/_mute! expr/_code body *("elif" expr/_code body)/_ARITY ?("else" body)/_ARITY/_unmute/if |
	"switch"/_mute! ?(~'{' expr)/_code '{' swbody '}'/_unmute/switch |
	"for"/_mute/_urange! fhead body/_unmute/for |
	"return"! expr %= ','/ARITY /return |
	"break" /brk |
	"continue" /cont |
	"include"! STRING/include |
	"import"! (STRING ?("as" IDENT/name)/ARITY)/import |
	"export"! IDENT/name % ','/ARITY /export |
	"defer"/_mute! expr/_code/_unmute/defer |
	"go"/_mute! expr/_code/_unmute/go |
	sexpr)/xline

doc = ?s *(';' ?s)

body = '{' doc/_code '}'

fhead = (~'{' s)/_code %= ';'/_ARITY

frange = ?(IDENT/name % ','/ARITY '=')/ARITY "range" expr

swbody = *("case"! expr/_code ':' doc/_code)/_ARITY ?("default"! ':' doc/_code)/_ARITY

fnbody = '(' IDENT/name %= ','/ARITY ?"..."/ARITY ')' '{'/_mute doc/_code '}'/_unmute

afn = '{'/_mute doc/_code '}'/_unmute/afn

member = IDENT | "class" | "new" | "recover" | "main" | "import" | "as" | "export" | "include" | "type"

newargs = ?('(' expr %= ','/ARITY ')')/ARITY

classb = "fn"! IDENT/name fnbody ?';'/mfn

methods = *classb/ARITY

clsname = IDENT/ref *('.' member/mref)

atom =
	'('! expr %= ','/ARITY ?"..."/ARITY ?',' ')'/call |
	'.'! member/mref |
	'['! ?expr/ARITY ?':'/ARITY ?expr/ARITY ']'/index

factor =
	INT/pushi |
	FLOAT/pushf |
	STRING/pushs |
	CHAR/pushc |
	(IDENT/ref | '('! expr ')' |
	"fn"! (~'{' fnbody/fn | afn) | '[' expr %= ','/ARITY ?',' ']'/slice) *atom |
	"new"! ('('! clsname ')' | clsname) newargs /new |
	"range"! expr/_range |
	"class"! '{' *classb/ARITY '}'/class |
	"recover"! '(' ')'/recover |
	"type"! '(' expr ')'/type |
	"main"! afn |
	'{'! (expr ':' expr) %= ','/ARITY ?',' '}'/map |
	'!' factor/not |
	'^' factor/bitnot |
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

// A Compiler represents a qlang compiler.
//
type Compiler struct {
	Opts  *ipt.Options
	code  *exec.Code
	ipt   interpreter.Engine
	libs  []string
	exits []func()
	mods  map[string]module
	gvars map[string]interface{}
	gstk  exec.Stack
	bctx  blockCtx
	forRg bool
	inFor bool
}

// New returns a qlang compiler instance.
//
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

// SetLibs sets searching paths when qlang searchs a library (ie. import a module).
//
func (p *Compiler) SetLibs(libs string) {

	p.libs = strings.Split(libs, ":")
}

// Vars returns compiling time variables, eg. __file__, __dir__, etc.
//
func (p *Compiler) Vars() map[string]interface{} {

	return p.gvars
}

// Code returns the generated code.
//
func (p *Compiler) Code() *exec.Code {

	return p.code
}

// Grammar returns the qlang compiler's grammar. It is required by tpl.Interpreter engine.
//
func (p *Compiler) Grammar() string {

	return grammar
}

// Fntable returns the qlang compiler's function table. It is required by tpl.Interpreter engine.
//
func (p *Compiler) Fntable() map[string]interface{} {

	return qlang.Fntable
}

// Stack returns nil (no stack). It is required by tpl.Interpreter engine.
//
func (p *Compiler) Stack() interpreter.Stack {

	return nil
}

func (p *Compiler) vMap() {

	arity := p.popArity()
	p.code.Block(exec.Call(qlang.MapFrom, arity*2))
}

func (p *Compiler) vSlice() {

	arity := p.popArity()
	p.code.Block(exec.SliceFrom(arity))
}

func (p *Compiler) vCall() {

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

func (p *Compiler) pop() {

	p.code.Block(exec.PopEx())
}

// CallFn generates a function call instruction. It is required by tpl.Interpreter engine.
//
func (p *Compiler) CallFn(fn interface{}) {

	p.code.Block(exec.Call(fn))
}

// -----------------------------------------------------------------------------

// DumpCode is mode how to dump code.
// 1 means to dump code with `rem` instruction; 2 means to dump clean code; 0 means don't dump code.
//
var DumpCode int

func (p *Compiler) codeLine(src interface{}) {

	ipt := p.ipt
	if ipt == nil {
		return
	}

	f := ipt.FileLine(src)
	p.code.CodeLine(f.File, f.Line)
	if DumpCode == 1 {
		text := string(ipt.Source(src))
		p.code.Block(exec.Rem(f.File, f.Line, text))
	}
}

// -----------------------------------------------------------------------------

var exports = map[string]interface{}{
	"$ARITY":   (*Compiler).arity,
	"$_ARITY":  (*Compiler).arity,
	"$_code":   (*Compiler).pushCode,
	"$name":    (*Compiler).pushName,
	"$pushi":   (*Compiler).pushInt,
	"$pushf":   (*Compiler).pushFloat,
	"$pushs":   (*Compiler).pushString,
	"$pushc":   (*Compiler).pushByte,
	"$index":   (*Compiler).index,
	"$mref":    (*Compiler).memberRef,
	"$ref":     (*Compiler).ref,
	"$tovar":   (*Compiler).toVar,
	"$slice":   (*Compiler).vSlice,
	"$map":     (*Compiler).vMap,
	"$call":    (*Compiler).vCall,
	"$assign":  (*Compiler).assign,
	"$massign": (*Compiler).multiAssign,
	"$inc":     (*Compiler).inc,
	"$dec":     (*Compiler).dec,
	"$adda":    (*Compiler).addAssign,
	"$suba":    (*Compiler).subAssign,
	"$mula":    (*Compiler).mulAssign,
	"$quoa":    (*Compiler).quoAssign,
	"$moda":    (*Compiler).modAssign,
	"$xora":    (*Compiler).xorAssign,
	"$lshra":   (*Compiler).lshrAssign,
	"$rshra":   (*Compiler).rshrAssign,
	"$bitanda": (*Compiler).bitandAssign,
	"$bitora":  (*Compiler).bitorAssign,
	"$andnota": (*Compiler).andnotAssign,
	"$defer":   (*Compiler).fnDefer,
	"$go":      (*Compiler).fnGo,
	"$chin":    (*Compiler).chanIn,
	"$chout":   (*Compiler).chanOut,
	"$recover": (*Compiler).fnRecover,
	"$return":  (*Compiler).fnReturn,
	"$fn":      (*Compiler).function,
	"$afn":     (*Compiler).anonymFn,
	"$include": (*Compiler).include,
	"$import":  (*Compiler).fnImport,
	"$export":  (*Compiler).export,
	"$mfn":     (*Compiler).memberFuncDecl,
	"$class":   (*Compiler).fnClass,
	"$new":     (*Compiler).fnNew,
	"$type":    (*Compiler).fnType,
	"$if":      (*Compiler).fnIf,
	"$switch":  (*Compiler).fnSwitch,
	"$for":     (*Compiler).fnFor,
	"$_urange": (*Compiler).unsetRange,
	"$_range":  (*Compiler).setRange,
	"$brk":     (*Compiler).fnBreak,
	"$cont":    (*Compiler).fnContinue,
	"$and":     (*Compiler).and,
	"$or":      (*Compiler).or,
	"$pop":     (*Compiler).pop,
	"$xline":   (*Compiler).codeLine,
}

func init() {
	qlang.Import("", exports)
}

// -----------------------------------------------------------------------------
