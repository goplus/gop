/*
 Copyright 2020 The GoPlus Authors (goplus.org)

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package bytecode

import (
	"reflect"
	"time"

	"github.com/goplus/gop/exec.spec"
	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

type varScope struct {
	vars   varsContext
	addrs  []reflect.Value
	parent *varScope
}

// A Context represents the context of an executor.
type Context struct {
	Stack
	varScope
	code    *Code
	defers  *theDefer
	updates []func(*Context)
	ip      int
	base    int
}

func (p *Context) regUpdate(fn func(*Context)) {
	p.updates = append(p.updates, fn)
}

func newSimpleContext(data []interface{}) *Context {
	return &Context{Stack: Stack{data: data}}
}

// NewContext returns a new context of an executor.
func NewContext(in exec.Code) *Context {
	code := in.(*Code)
	p := &Context{
		code: code,
	}
	p.Init()
	if len(code.vlist) > 0 {
		p.vars = code.makeVarsContext(p)
	}
	return p
}

// Go starts a new goroutine to run.
func (ctx *Context) Go(arity int, f func(goctx *Context)) {
	goctx := &Context{
		code: ctx.code,
	}
	goctx.Init()
	base := len(ctx.data) - arity
	parent := ctx.varScope
	goctx.parent = &parent
	goctx.data = append(goctx.data, ctx.data[base:]...)
	ctx.data = ctx.data[:base]
	go f(goctx)
}

// CloneSetVarScope clone already set varScope to new context
func (ctx *Context) CloneSetVarScope(new *Context) {
	if ctx.vars.IsValid() {
		for i := 0; i < ctx.vars.NumField(); i++ {
			new.varScope.setVar(uint32(i), ctx.varScope.getVar(uint32(i)))
		}
	}
	new.defers = ctx.defers
	new.updates = ctx.updates
	for _, fn := range new.updates {
		fn(new)
	}
}

type savedScopeCtx struct {
	base int
	ip   int
	varScope
}

func (ctx *Context) switchScope(parent *varScope, vmgr *varManager) (old savedScopeCtx) {
	old.base = ctx.base
	old.ip = ctx.ip
	old.varScope = ctx.varScope
	ctx.base = len(ctx.data)
	ctx.parent = parent
	ctx.vars = vmgr.makeVarsContext(ctx)

	size := ctx.Len()
	ctx.addrs = make([]reflect.Value, size)
	for i := size; i > 0; i-- {
		v := reflect.ValueOf(ctx.Get(-i))
		if v.Kind() == reflect.Struct {
			nv := reflect.New(v.Type()).Elem()
			nv.Set(v)
			ctx.addrs[size-i] = nv
		} else {
			ctx.addrs[size-i] = v
		}
	}
	return
}

func (ctx *Context) restoreScope(old savedScopeCtx) {
	ctx.ip = old.ip
	ctx.base = old.base
	ctx.varScope = old.varScope
	ctx.addrs = old.addrs
}

func (ctx *Context) getScope(local bool) *varScope {
	scope := ctx.parent
	if scope == nil || local {
		vs := ctx.varScope
		return &vs
	}
	for scope.parent != nil {
		scope = scope.parent
	}
	return scope
}

// Run executes the code.
func (ctx *Context) Run() {
	defer ctx.execDefers()
	ctx.Exec(0, ctx.code.Len())
}

// UpdateCode update increasable code for current context
func (ctx *Context) UpdateCode(in exec.Code) {
	code := in.(*Code)
	ctx.Init()
	if len(code.vlist) > 0 {
		vars := code.makeVarsContext(ctx)
		if ctx.vars.IsValid() {
			for i := 0; i < ctx.vars.NumField(); i++ {
				vars.Field(i).Set(ctx.vars.Field(i))
			}
		}
		ctx.vars = vars
	}
	for _, fn := range ctx.updates {
		fn(ctx)
	}
	ctx.code = code
}

// Exec executes a code block from ip to ipEnd.
func (ctx *Context) Exec(ip, ipEnd int) (currentIP int) {
	const allowProfile = true
	var lastInstr Instr
	var start time.Time
	var data = ctx.code.data
	ctx.ip = ip
	for ctx.ip < ipEnd {
		i := data[ctx.ip]
		ctx.ip++
		if allowProfile && doProfile {
			if lastInstr != 0 {
				instrProfile(lastInstr, time.Since(start))
			}
			lastInstr, start = i, time.Now()
		}
		switch i >> bitsOpShift {
		case opPushInt:
			const mask = uint32(bitsOpIntOperand >> 1)
			switch i & ^mask {
			case opPushInt << bitsOpShift: // push kind=int
				ctx.Push(int(i & mask))
			default:
				execPushInt(i, ctx)
			}
		case opBuiltinOp:
			execBuiltinOp(i, ctx)
		case opCallFunc:
			fun := ctx.code.funs[i&bitsOperand]
			fun.exec(ctx, ctx.getScope(fun.nestDepth > 1))
		case opJmp:
			execJmp(i, ctx)
		case opJmpIf:
			execJmpIf(i, ctx)
		case opPushConstR:
			execPushConstR(i, ctx)
		case opLoadVar:
			execLoadVar(i, ctx)
		case opStoreVar:
			execStoreVar(i, ctx)
		case opCallFuncv:
			execFuncv(i, ctx)
		case opCallGoFunc:
			execGoFunc(i, ctx)
		case opCallGoFuncv:
			execGoFuncv(i, ctx)
		case opReturn:
			currentIP = ctx.ip
			if i == iReturn {
				ctx.ip = int(i)
			} else {
				ctx.ip = ipReturnN
			}
			goto finished
		case opPushUint:
			execPushUint(i, ctx)
		default:
			if fn := execTable[i>>bitsOpShift]; fn != nil {
				fn(i, ctx)
			} else {
				log.Panicln("Exec: unknown instr -", i>>bitsOpShift, "ip:", ctx.ip-1)
			}
		}
	}
finished:
	if allowProfile && doProfile {
		if lastInstr != 0 {
			instrProfile(lastInstr, time.Since(start))
		}
	}
	return
}

var _execTable = [...]func(i Instr, p *Context){
	opCallGoFunc:    execGoFunc,
	opCallGoFuncv:   execGoFuncv,
	opCallFunc:      execFunc,
	opCallFuncv:     execFuncv,
	opPushInt:       execPushInt,
	opPushUint:      execPushUint,
	opPushValSpec:   execPushValSpec,
	opPushConstR:    execPushConstR,
	opIndex:         execIndex,
	opMake:          execMake,
	opAppend:        execAppend,
	opBuiltinOp:     execBuiltinOp,
	opJmp:           execJmp,
	opJmpIf:         execJmpIf,
	opCaseNE:        execCaseNE,
	opPop:           execPop,
	opLoadVar:       execLoadVar,
	opStoreVar:      execStoreVar,
	opAddrVar:       execAddrVar,
	opAddrOp:        execAddrOp,
	opLoadGoVar:     execLoadGoVar,
	opStoreGoVar:    execStoreGoVar,
	opAddrGoVar:     execAddrGoVar,
	opLoad:          execLoad,
	opAddr:          execAddr,
	opStore:         execStore,
	opClosure:       execClosure,
	opCallClosure:   execCallClosure,
	opGoClosure:     execGoClosure,
	opCallGoClosure: execCallGoClosure,
	opMakeArray:     execMakeArray,
	opMakeMap:       execMakeMap,
	opZero:          execZero,
	opForPhrase:     execForPhrase,
	opLstComprehens: execListComprehension,
	opMapComprehens: execMapComprehension,
	opTypeCast:      execTypeCast,
	opSlice:         execSlice,
	opSlice3:        execSlice3,
	opMapIndex:      execMapIndex,
	opGoBuiltin:     execGoBuiltin,
	opErrWrap:       execErrWrap,
	opWrapIfErr:     execWrapIfErr,
	opDefer:         execDefer,
	opGo:            execGo,
	opLoadField:     execLoadField,
	opStoreField:    execStoreField,
	opAddrField:     execAddrField,
	opStruct:        execStruct,
	opSend:          execSend,
	opRecv:          execRecv,
}

var execTable []func(i Instr, p *Context)

func init() {
	execTable = _execTable[:]
}

// -----------------------------------------------------------------------------
