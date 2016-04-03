package exec

// -----------------------------------------------------------------------------

type iBlock struct {
	start int
	end   int
}

func (p *iBlock) Exec(stk *Stack, ctx *Context) {

	ip := ctx.ip
	ctx.Code.Exec(p.start, p.end, stk, ctx)
	ctx.ip = ip
}

func Block(start, end int) Instr {

	return &iBlock{start, end}
}

// -----------------------------------------------------------------------------

type iModule struct {
	start int
	end   int
	id    string
}

func (p *iModule) Exec(stk *Stack, ctx *Context) {

	exports, ok := ctx.mods[p.id]
	if !ok {
		parent := &Context{
			Code:  ctx.Code,
			Stack: ctx.Stack,
			mods:  ctx.mods,
			vars:  make(map[string]interface{}),
		}
		fn := NewFunction(nil, p.start, p.end, nil, false)		
		fn.parent = parent
		_, mod := fn.ExtCall()
		exports = mod.Exports()
		ctx.mods[p.id] = exports
	}
	stk.Push(exports)
}

func Module(id string, start, end int) Instr {

	return &iModule{start, end, id}
}

// -----------------------------------------------------------------------------

type iAs struct {
	name string
}

func (p *iAs) Exec(stk *Stack, ctx *Context) {

	name := p.name
	if _, ok := ctx.vars[name]; ok { // 符号已经存在
		panic("import `" + name + "` error: ident exists")
	}

	v, ok := stk.Pop()
	if !ok {
		panic(ErrStackDamaged)
	}

	ctx.vars[name] = v
}

func As(name string) Instr {

	return &iAs{name}
}

// -----------------------------------------------------------------------------

type iExport struct {
	names []string
}

func (p *iExport) Exec(stk *Stack, ctx *Context) {

	ctx.export = append(ctx.export, p.names...)
}

func Export(names ...string) Instr {

	return &iExport{names}
}

// -----------------------------------------------------------------------------

