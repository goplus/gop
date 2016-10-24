package interpreter

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"qiniupkg.com/text/tpl.v1"
	"qiniupkg.com/text/tpl.v1/interpreter.util"
)

// -----------------------------------------------------------------------------

var (
	typeIntf    = reflect.TypeOf((*interface{})(nil)).Elem()
	typeEng     = reflect.TypeOf((*interpreter.Engine)(nil)).Elem()
	zeroIntfVal = reflect.Zero(typeIntf)
)

func ParseInt(lit string) (v int64, err error) {

	switch lit[0] {
	case '0':
		if len(lit) == 1 {
			return
		}
		switch lit[1] {
		case 'x', 'X':
			return strconv.ParseInt(lit[2:], 16, 64)
		default:
			return strconv.ParseInt(lit[1:], 8, 64)
		}
	default:
		return strconv.ParseInt(lit, 10, 64)
	}
}

func ParseFloat(lit string) (v float64, err error) {

	switch lit[0] {
	case '0':
		if len(lit) == 1 {
			return
		}
		switch lit[1] {
		case 'x', 'X':
			v1, err1 := strconv.ParseInt(lit[2:], 16, 64)
			return float64(v1), err1
		default:
			v1, err1 := strconv.ParseInt(lit[1:], 8, 64)
			if err1 == nil {
				return float64(v1), nil
			}
		}
	}
	return strconv.ParseFloat(lit, 64)
}

func GetValue(t reflect.Type, lit string) (ret reflect.Value, err error) {

	switch t.Kind() {
	case reflect.String:
		return reflect.ValueOf(lit), nil
	case reflect.Float64:
		v, err1 := ParseFloat(lit)
		if err1 != nil {
			err = err1
			return
		}
		return reflect.ValueOf(v), nil
	case reflect.Int64:
		v, err1 := ParseInt(lit)
		if err1 != nil {
			err = err1
			return
		}
		return reflect.ValueOf(v), nil
	case reflect.Int:
		v, err1 := ParseInt(lit)
		if err1 != nil {
			err = err1
			return
		}
		return reflect.ValueOf(int(v)), nil
	}
	err = errors.New("unsupported type: " + t.String())
	return
}

// -----------------------------------------------------------------------------

type RuntimeError struct {
	Grammar tpl.Grammar
	Ctx     tpl.Context
	Err     error
	Src     []tpl.Token
}

func (p *RuntimeError) Error() string {
	if p.Ctx == nil {
		return fmt.Sprintf("runtime error: %v", p.Err)
	}
	s := p.Ctx.Tokener()
	file, line := tpl.FileLine(p.Src, s)
	if file == "" {
		return fmt.Sprintf("line %d: runtime error: %v", line, p.Err)
	}
	return fmt.Sprintf("%s:%d: runtime error: %v", file, line, p.Err)
}

// -----------------------------------------------------------------------------

type Options struct {
	Scanner  tpl.Tokener
	ScanMode tpl.ScanMode
}

type Engine struct {
	*tpl.CompileRet
	Interpreter interpreter.Interface
	mute        int
}

type fnCaller interface {
	CallFn(fn interface{})
}

var (
	optionsDef Options
	ptrError   *error
	typeError  = reflect.TypeOf(ptrError).Elem()
	typeInt    = reflect.TypeOf(0)
)

var (
	InsertSemis = &Options{ScanMode: tpl.InsertSemis}
)

func New(ipt interpreter.Interface, options *Options) (p *Engine, err error) {

	if options == nil {
		options = &optionsDef
	}

	p = new(Engine)

	stk := ipt.Stack()
	tstk := reflect.TypeOf(stk)
	tipt := reflect.TypeOf(ipt)
	compiler := new(tpl.Compiler)
	vEngine := reflect.ValueOf(interpreter.Engine(p))
	fnCaller, hasFnCaller := ipt.(fnCaller)

	fntable := ipt.Fntable()
	initer := func() {
		p.mute = 0 // bugfix: 语法匹配过程有可能有异常，需要通过 init 完成 mute 状态的重置
	}
	marker := func(g tpl.Grammar, mark string) tpl.Grammar {
		fn, ok := fntable["$"+mark]
		if !ok {
			switch mark {
			case "_mute":
				return tpl.Action(g, func(tokens []tpl.Token, g tpl.Grammar) {
					p.mute++
				})
			case "_unmute":
				return tpl.Action(g, func(tokens []tpl.Token, g tpl.Grammar) {
					p.mute--
				})
			case "_tr":
				return tpl.Action(g, func(tokens []tpl.Token, g tpl.Grammar) {
					println(string(tpl.Source(tokens, compiler.Scanner)))
				})
			default:
				panic("invalid grammar mark: " + mark)
			}
		}
		tfn := reflect.TypeOf(fn)
		if tfn.Kind() != reflect.Func {
			panic("function table item not a function: " + mark)
		}
		if tfn.IsVariadic() {
			panic("invalid mark - doesn't support variadic function: " + mark)
		}
		n := tfn.NumIn()
		if n < 1 {
			panic("invalid function (no param): " + mark)
		}
		trecvr := tfn.In(0)
		kindRecvr := 0
		if trecvr == tstk {
			kindRecvr = 1
		} else if trecvr == tipt {
			kindRecvr = 2
		}
		if kindRecvr > 0 {
			if n > 2 {
				panic("invalid function (too many params): " + mark)
			}
			if nout := tfn.NumOut(); nout > 0 {
				if nout != 1 {
					panic("invalid function (too many return values): " + mark)
				}
				if tfn.Out(0) != typeError {
					panic("invalid function (return type isn't error): " + mark)
				}
			}
		}
		idxLen := 0
		if mark[0] == '_' {
			idxLen++
		}
		isLen := (mark[idxLen] >= 'A' && mark[idxLen] <= 'Z') // 用大写的mark表示取Grammar.Len()
		action := func(tokens []tpl.Token, g tpl.Grammar) {
			if p.mute > 0 { // 禁止所有非 '_' 开头的 action
				if mark[0] != '_' || p.mute > 1 {
					return
				}
			}
			defer func() {
				if e := recover(); e != nil {
					if err, ok := e.(*RuntimeError); ok {
						panic(err)
					}
					if err, ok := e.(*tpl.MatchError); ok {
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
					err = &RuntimeError{
						Grammar: g,
						Ctx:     p.Matcher,
						Err:     err,
						Src:     tokens,
					}
					panic(err)
				}
			}()
			if kindRecvr > 0 {
				args := make([]reflect.Value, n)
				if kindRecvr == 1 {
					args[0] = reflect.ValueOf(stk)
				} else {
					args[0] = reflect.ValueOf(p.Interpreter)
				}
				if n == 2 {
					targ1 := tfn.In(1)
					if isLen {
						if targ1 != typeInt {
							panic("invalid function (second argument type isn't int): " + mark)
						}
						args[1] = reflect.ValueOf(g.Len())
					} else {
						switch targ1 {
						case typeIntf:
							if len(tokens) == 0 { // use nil as empty source
								args[1] = zeroIntfVal
							} else {
								args[1] = reflect.ValueOf(tokens)
							}
						case typeEng:
							args[1] = vEngine
						default:
							v, err := GetValue(targ1, tokens[0].Literal)
							if err != nil {
								panic("invalid function `" + mark + "` argument, " + err.Error())
							}
							args[1] = v
						}
					}
				}
				ret := reflect.ValueOf(fn).Call(args)
				if len(ret) > 0 {
					if retv := ret[0].Interface(); retv != nil {
						panic(retv.(error))
					}
				}
			} else if hasFnCaller {
				fnCaller.CallFn(fn)
			} else {
				if err := interpreter.DoCall(stk, fn, n, n, false); err != nil {
					panic("call function `$" + mark + "`: " + err.Error())
				}
			}
		}
		return tpl.Action(g, action)
	}

	grammar := ipt.Grammar()
	compiler.Grammar = []byte(grammar)
	compiler.Scanner = options.Scanner
	compiler.ScanMode = options.ScanMode
	compiler.Marker = marker
	compiler.Init = initer
	ret, err := compiler.Cl()
	if err != nil {
		return
	}

	p.CompileRet = &ret
	p.Interpreter = ipt
	return
}

func (p *Engine) Source(src interface{}) (text []byte) {

	if src == nil {
		return
	}
	return tpl.Source(src.([]tpl.Token), p.Scanner)
}

func (p *Engine) FileLine(src interface{}) (f interpreter.FileLine) {

	if src == nil {
		return
	}
	f.File, f.Line = tpl.FileLine(src.([]tpl.Token), p.Scanner)
	return
}

func (p *Engine) EvalCode(fn interpreter.Interface, name string, src interface{}) error {

	old := p.Interpreter
	mute := p.mute
	p.Interpreter = fn
	defer func() { // 正常的 return 使用了 panic，所以需要保护
		p.Interpreter = old
		p.mute = mute
	}()
	return p.EvalSub(name, src)
}

// -----------------------------------------------------------------------------
