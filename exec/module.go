package exec

import (
	"sync"
)

// -----------------------------------------------------------------------------
// AnonymFn(匿名函数)

type iAnonymFn struct {
	start  int
	end    int
	symtbl map[string]int
}

func (p *iAnonymFn) Exec(stk *Stack, ctx *Context) {

	fn := NewFunction(nil, p.start, p.end, p.symtbl, nil, false)
	fn.Parent = ctx
	stk.Push(fn.Call(stk))
}

// AnonymFn returns an instruction that creates an anonymous function.
//
func AnonymFn(start, end int, symtbl map[string]int) Instr {

	return &iAnonymFn{start, end, symtbl}
}

// -----------------------------------------------------------------------------
// Module

type importMod struct {
	exports map[string]interface{}
	sync.Mutex
}

func (p *importMod) Lock() (exports map[string]interface{}, uninited bool) {

	p.Mutex.Lock()
	if p.exports == nil {
		return nil, true
	}
	exports = p.exports
	p.Mutex.Unlock()
	return
}

type moduleMgr struct {
	mods  map[string]*importMod
	mutex sync.Mutex
}

func (p *moduleMgr) get(id string) *importMod {

	p.mutex.Lock()
	mod, ok := p.mods[id]
	if !ok {
		mod = new(importMod)
		p.mods[id] = mod
	}
	p.mutex.Unlock()
	return mod
}

type iModule struct {
	start  int
	end    int
	symtbl map[string]int
	id     string
}

func (p *iModule) Exec(stk *Stack, ctx *Context) {

	mod := ctx.modmgr.get(p.id)
	exports, uninited := mod.Lock()
	if uninited {
		defer mod.Unlock()
		modCtx := &Context{
			Code:   ctx.Code,
			Stack:  ctx.Stack,
			modmgr: ctx.modmgr,
		}
		modCtx.initVars(p.symtbl)
		modFn := NewFunction(nil, p.start, p.end, p.symtbl, nil, false)
		modFn.ExtCall(modCtx)
		exports = modCtx.Exports()
		mod.exports = exports
	}
	stk.Push(exports)
}

// Module returns an instruction that creates a new module.
//
func Module(id string, start, end int, symtbl map[string]int) Instr {

	return &iModule{start, end, symtbl, id}
}

// -----------------------------------------------------------------------------

type iAs struct {
	name int
}

func (p *iAs) Exec(stk *Stack, ctx *Context) {

	v, ok := stk.Pop()
	if !ok {
		panic(ErrStackDamaged)
	}
	ctx.FastSetVar(p.name, v)
}

// As returns an instruction that specifies an alias name of a module.
//
func As(name int) Instr {

	return &iAs{name}
}

// -----------------------------------------------------------------------------

type iExport struct {
	names []string
}

func (p *iExport) Exec(stk *Stack, ctx *Context) {

	ctx.export = append(ctx.export, p.names...)
}

// Export returns an instruction that exports some module symbols.
//
func Export(names ...string) Instr {

	return &iExport{names}
}

// -----------------------------------------------------------------------------
