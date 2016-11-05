package interpreter

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"qiniupkg.com/text/tpl.v1"
	"qiniupkg.com/text/tpl.v1/interpreter"
	"qlang.io/exec"
)

var (
	ErrNotObject = errors.New("not object")
)

// -----------------------------------------------------------------------------

type Options struct {
	Scanner  tpl.Tokener
	ScanMode tpl.ScanMode
}

type Engine struct {
	*tpl.CompileRet
	Interpreter *exec.Object
	mute        int
}

var (
	optionsDef Options
	typeFunc   = reflect.TypeOf((*exec.Function)(nil))
)

var (
	InsertSemis = &Options{ScanMode: tpl.InsertSemis}
)

func lastWord(s string) string {

	idx := strings.LastIndexFunc(s, func(c rune) bool {
		return c >= 'A' && c <= 'Z'
	})
	if idx >= 0 {
		return s[idx:]
	}
	return ""
}

func call(estk *exec.Stack, ip *exec.Object, name string) interface{} {

	if fn, ok := ip.Cls.Fns[name]; ok {
		return fn.Call(estk, ip)
	}
	panic("method not found: " + name)
}

func callN(estk *exec.Stack, ip *exec.Object, name string, args ...interface{}) interface{} {

	if fn, ok := ip.Cls.Fns[name]; ok {
		args1 := make([]interface{}, 1, len(args)+1)
		args1[0] = ip
		return fn.Call(estk, append(args1, args...)...)
	}
	panic("method not found: " + name)
}

func callFn(estk *exec.Stack, stk *exec.Object, fn interface{}, arity int) {

	var in []reflect.Value
	if arity > 0 {
		ret := reflect.ValueOf(callN(estk, stk, "popArgs", arity))
		if ret.Kind() != reflect.Slice {
			panic("method stack.popArgs doesn't return a slice")
		}
		n := ret.Len()
		in = make([]reflect.Value, n)
		for i := 0; i < n; i++ {
			in[i] = ret.Index(i)
		}
	}
	vfn := reflect.ValueOf(fn)
	out := vfn.Call(in)
	if len(out) == 0 {
		callN(estk, stk, "push", nil)
	} else {
		callN(estk, stk, "push", out[0].Interface())
	}
}

func callVfn(estk *exec.Stack, stk *exec.Object, fn *exec.Function, arity int) {

	v := callN(estk, stk, "popArgs", arity)
	if args, ok := v.([]interface{}); ok {
		callN(estk, stk, "push", fn.Call(estk, args...))
		return
	}

	ret := reflect.ValueOf(v)
	if ret.Kind() != reflect.Slice {
		panic("method stack.popArgs doesn't return a slice")
	}
	n := ret.Len()
	in := make([]interface{}, n)
	for i := 0; i < n; i++ {
		in[i] = ret.Index(i).Interface()
	}
	callN(estk, stk, "push", fn.Call(estk, in...))
}

func New(ip interface{}, options *Options) (p *Engine, err error) {

	defer func() {
		if e := recover(); e != nil {
			switch v := e.(type) {
			case string:
				err = errors.New(v)
			case error:
				err = v
			default:
				panic(e)
			}
		}
	}()

	estk := exec.NewStack()
	p = newEngine(estk, ip, options)
	return
}

func (p *Engine) Exec(code []byte, fname string) (err error) {

	defer func() {
		if e := recover(); e != nil {
			switch v := e.(type) {
			case string:
				err = errors.New(v)
			case error:
				err = v
			default:
				panic(e)
			}
		}
	}()

	err = p.MatchExactly(code, fname)
	return
}

func (p *Engine) Eval(expr string) (err error) {

	return p.Exec([]byte(expr), "")
}

func newEngine(estk *exec.Stack, ip1 interface{}, options *Options) (p *Engine) {

	if options == nil {
		options = &optionsDef
	}

	ip, ok := ip1.(*exec.Object)
	if !ok {
		panic(ErrNotObject)
	}

	p = new(Engine)
	vstk := call(estk, ip, "stack")
	stk, ok := vstk.(*exec.Object)
	if !ok {
		panic("stack isn't an object")
	}

	vfntable := call(estk, ip, "fntable")
	fntable, ok := vfntable.(map[string]interface{})
	if !ok {
		panic("fntable isn't a map[string]interface{} object")
	}

	compiler := new(tpl.Compiler)

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
		if tfn != typeFunc {
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
			action := func(tokens []tpl.Token, g tpl.Grammar) {
				if p.mute > 0 { // 禁止所有非 '_' 开头的 action
					if mark[0] != '_' || p.mute > 1 {
						return
					}
				}
				callFn(estk, stk, fn, n)
			}
			return tpl.Action(g, action)
		}

		vfn := fn.(*exec.Function)
		if vfn.Variadic {
			panic("invalid mark - doesn't support variadic function: " + mark)
		}
		n := len(vfn.Args)
		if n < 1 {
			panic("invalid function (no param): " + mark)
		}
		kindRecvr := 0
		if vfn.Cls == stk.Cls {
			kindRecvr = 1
		} else if vfn.Cls == ip.Cls {
			kindRecvr = 2
		}
		if kindRecvr > 0 && n > 2 {
			panic("invalid function (too many params): " + mark)
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
			if kindRecvr > 0 {
				args := make([]interface{}, n)
				if kindRecvr == 1 {
					args[0] = stk
				} else {
					args[0] = p.Interpreter
				}
				if n == 2 {
					if isLen {
						args[1] = g.Len()
					} else {
						switch lastWord(mark) {
						case "Src":
							if len(tokens) != 0 {
								args[1] = tokens
							}
						case "Engine":
							args[1] = p
						case "Float":
							v, err := strconv.ParseFloat(tokens[0].Literal, 64)
							if err != nil {
								panic(fmt.Sprintf("parse mark `%s` failed: %v", mark, err))
							}
							args[1] = v
						case "Int":
							v, err := strconv.Atoi(tokens[0].Literal)
							if err != nil {
								panic(fmt.Sprintf("parse mark `%s` failed: %v", mark, err))
							}
							args[1] = v
						default:
							args[1] = tokens[0].Literal
						}
					}
				}
				vfn.Call(estk, args...)
			} else {
				callVfn(estk, stk, vfn, n)
			}
		}
		return tpl.Action(g, action)
	}

	vgrammar := call(estk, ip, "grammar")
	grammar, ok := vgrammar.(string)
	if !ok {
		panic("grammar isn't a string object")
	}

	compiler.Grammar = []byte(grammar)
	compiler.Scanner = options.Scanner
	compiler.ScanMode = options.ScanMode
	compiler.Marker = marker
	compiler.Init = initer
	ret, err := compiler.Cl()
	if err != nil {
		panic(err)
	}

	p.CompileRet = &ret
	p.Interpreter = ip
	return
}

func (p *Engine) EvalCode(fn *exec.Object, name string, src interface{}) error {

	old := p.Interpreter
	mute := p.mute
	p.Interpreter = fn
	defer func() { // 正常的 return 使用了 panic，所以需要保护
		p.Interpreter = old
		p.mute = mute
	}()
	return p.EvalSub(name, src)
}

func ClassOf(v interface{}) interface{} {

	if o, ok := v.(*exec.Object); ok {
		return o.Cls
	}
	return nil
}

func NewRuntimeError(msg string) error {

	return &interpreter.RuntimeError{Err: errors.New(msg)}
}

// -----------------------------------------------------------------------------

var Exports = map[string]interface{}{
	"classof":         ClassOf,
	"newRuntimeError": NewRuntimeError,
	"interpreter":     New,
	"insertSemis":     InsertSemis,
}

// -----------------------------------------------------------------------------
