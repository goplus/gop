package exec

import (
	"reflect"

	"github.com/qiniu/x/log"
)

func execListComprehension(i Instr, p *Context) {
	addr := i & bitsOperand
	c := p.code.comprehens[addr]
	base := len(p.data)
	p.Exec(p.ip, c.End)
	makeArray(c.TypeOut, len(p.data)-base, p)
}

func execMapComprehension(i Instr, p *Context) {
	addr := i & bitsOperand
	c := p.code.comprehens[addr]
	base := len(p.data)
	p.Exec(p.ip, c.End)
	makeMap(c.TypeOut, (len(p.data)-base)>>1, p)
}

func execForPhrase(i Instr, p *Context) {
	addr := i & bitsOperand
	p.code.fors[addr].exec(p)
}

func (c *ForPhrase) exec(p *Context) {
	var ctxBody *Context
	if c.block != nil {
		ctxBody = newContextEx(p, p.Stack, p.code, &c.block.varManager)
	} else {
		ctxBody = p
	}
	data := reflect.ValueOf(p.Pop())
	switch data.Kind() {
	case reflect.Map:
		c.execMapRange(data, p, ctxBody)
	default:
		c.execListRange(data, p, ctxBody)
	}
}

func (c *ForPhrase) execListRange(data reflect.Value, ctxFor, ctxBody *Context) {
	n := data.Len()
	ip, ipCond, ipEnd := ctxFor.ip, c.Cond, c.End
	key, val := c.Key, c.Value
	for i := 0; i < n; i++ {
		if key != nil {
			ctxFor.setVar(key.idx, i)
		}
		if val != nil {
			ctxFor.setVar(val.idx, data.Index(i).Interface())
		}
		if ipCond > 0 {
			ctxBody.Exec(ip, ipCond)
			if ok := ctxBody.Pop().(bool); ok {
				ctxBody.Exec(ipCond, ipEnd)
			}
		} else {
			ctxBody.Exec(ip, ipEnd)
		}
	}
	ctxFor.ip = ipEnd
}

func (c *ForPhrase) execMapRange(data reflect.Value, ctxFor, ctxBody *Context) {
	iter := data.MapRange()
	ip, ipCond, ipEnd := ctxFor.ip, c.Cond, c.End
	key, val := c.Key, c.Value
	for iter.Next() {
		if key != nil {
			ctxFor.setVar(key.idx, iter.Key().Interface())
		}
		if val != nil {
			ctxFor.setVar(val.idx, iter.Value().Interface())
		}
		if ipCond > 0 {
			ctxBody.Exec(ip, ipCond)
			if ok := ctxBody.Pop().(bool); ok {
				ctxBody.Exec(ipCond, ipEnd)
			}
		} else {
			ctxBody.Exec(ip, ipEnd)
		}
	}
	ctxFor.ip = ipEnd
}

func execMakeArray(i Instr, p *Context) {
	typSlice := getType(i&bitsOpCallFuncvOperand, p)
	arity := int((i >> bitsOpCallFuncvShift) & bitsFuncvArityOperand)
	if arity == bitsFuncvArityVar { // args...
		v := reflect.ValueOf(p.Get(-1))
		n := v.Len()
		ret := reflect.MakeSlice(typSlice, n, n)
		reflect.Copy(ret, v)
		p.Ret(1, ret.Interface())
	} else {
		if arity == bitsFuncvArityMax {
			arity = p.Pop().(int) + bitsFuncvArityMax
		}
		makeArray(typSlice, arity, p)
	}
}

func makeArray(typSlice reflect.Type, arity int, p *Context) {
	args := p.GetArgs(arity)
	var ret reflect.Value
	if typSlice.Kind() == reflect.Slice {
		ret = reflect.MakeSlice(typSlice, arity, arity)
	} else {
		ret = reflect.New(typSlice).Elem()
	}
	for i, arg := range args {
		ret.Index(i).Set(getElementOf(arg, typSlice))
	}
	p.Ret(arity, ret.Interface())
}

func execAppend(i Instr, p *Context) {
	arity := int(i & bitsOperand)
	if arity == bitsFuncvArityVar { // args...
		args := p.GetArgs(2)
		ret := reflect.AppendSlice(reflect.ValueOf(args[0]), reflect.ValueOf(args[1]))
		p.Ret(2, ret.Interface())
	} else {
		args := p.GetArgs(arity)
		ret := reflect.Append(reflect.ValueOf(args[0]), ToValues(args[1:])...)
		p.Ret(arity, ret.Interface())
	}
}

func execMake(i Instr, p *Context) {
	typ := getType(i&bitsOpCallFuncvOperand, p)
	arity := int((i >> bitsOpCallFuncvShift) & bitsFuncvArityOperand)
	switch typ.Kind() {
	case reflect.Slice:
		var cap = p.Get(-1).(int)
		var len int
		if arity > 1 {
			len = p.Get(-2).(int)
		} else {
			len = cap
		}
		p.Ret(arity, reflect.MakeSlice(typ, len, cap).Interface())
	case reflect.Map:
		if arity == 0 {
			p.Push(reflect.MakeMap(typ).Interface())
			return
		}
		n := p.Get(-1).(int)
		p.Ret(arity, reflect.MakeMapWithSize(typ, n).Interface())
	case reflect.Chan:
		var buffer int
		if arity > 0 {
			buffer = p.Get(-1).(int)
		}
		p.Ret(arity, reflect.MakeChan(typ, buffer).Interface())
	default:
		panic("make: unexpected type")
	}
}

func execMakeMap(i Instr, p *Context) {
	typMap := getType(i&bitsOpCallFuncvOperand, p)
	arity := int((i >> bitsOpCallFuncvShift) & bitsFuncvArityOperand)
	if arity == bitsFuncvArityMax {
		arity = p.Pop().(int) + bitsFuncvArityMax
	}
	makeMap(typMap, arity, p)
}

func makeMap(typMap reflect.Type, arity int, p *Context) {
	n := arity << 1
	args := p.GetArgs(n)
	ret := reflect.MakeMapWithSize(typMap, arity)
	for i := 0; i < n; i += 2 {
		key := getKeyOf(args[i], typMap)
		val := getElementOf(args[i+1], typMap)
		ret.SetMapIndex(key, val)
	}
	p.Ret(n, ret.Interface())
}

func execTypeCast(i Instr, p *Context) {
	args := p.GetArgs(1)
	typ := getType(i&bitsOperand, p)
	args[0] = reflect.ValueOf(args[0]).Convert(typ).Interface()
}

func execZero(i Instr, p *Context) {
	typ := getType(i&bitsOperand, p)
	p.Push(reflect.Zero(typ).Interface())
}

// ToValues converts []interface{} into []reflect.Value.
func ToValues(args []interface{}) []reflect.Value {
	ret := make([]reflect.Value, len(args))
	for i, arg := range args {
		ret[i] = reflect.ValueOf(arg)
	}
	return ret
}

// -----------------------------------------------------------------------------

// ForPhrase represents a for range phrase.
type ForPhrase struct {
	Key, Value *Var // Key, Value may be nil
	Cond, End  int
	TypeIn     reflect.Type
	block      *blockCtx
}

// NewForPhrase creates a new ForPhrase instance.
func NewForPhrase(in reflect.Type) *ForPhrase {
	return &ForPhrase{TypeIn: in}
}

// NewForPhraseWith creates a new ForPhrase instance with executing context.
func NewForPhraseWith(in reflect.Type, nestDepth uint32) *ForPhrase {
	return &ForPhrase{TypeIn: in, block: newBlockCtx(nestDepth)}
}

// Comprehension represents a list/map comprehension.
type Comprehension struct {
	TypeOut reflect.Type
	End     int
}

// NewComprehension creates a new Comprehension instance.
func NewComprehension(out reflect.Type) *Comprehension {
	return &Comprehension{TypeOut: out}
}

// ForPhrase instr
func (p *Builder) ForPhrase(f *ForPhrase, key, val *Var) *Builder {
	f.Key, f.Value = key, val
	if key != nil {
		p.DefineVar(key)
	}
	if val != nil {
		p.DefineVar(val)
	}
	if f.block != nil {
		f.block.parent = p.varManager
		p.varManager = &f.block.varManager
	}
	code := p.code
	addr := uint32(len(code.fors))
	code.fors = append(code.fors, f)
	code.data = append(code.data, (opForPhrase<<bitsOpShift)|addr)
	return p
}

// FilterForPhrase instr
func (p *Builder) FilterForPhrase(f *ForPhrase) *Builder {
	f.Cond = len(p.code.data)
	return p
}

// EndForPhrase instr
func (p *Builder) EndForPhrase(f *ForPhrase) *Builder {
	f.End = len(p.code.data)
	if f.block != nil {
		p.varManager = f.block.parent
	}
	return p
}

// ListComprehension instr
func (p *Builder) ListComprehension(c *Comprehension) *Builder {
	code := p.code
	addr := uint32(len(code.comprehens))
	code.comprehens = append(code.comprehens, c)
	code.data = append(code.data, (opLstComprehens<<bitsOpShift)|addr)
	return p
}

// MapComprehension instr
func (p *Builder) MapComprehension(c *Comprehension) *Builder {
	code := p.code
	addr := uint32(len(code.comprehens))
	code.comprehens = append(code.comprehens, c)
	code.data = append(code.data, (opMapComprehens<<bitsOpShift)|addr)
	return p
}

// EndComprehension instr
func (p *Builder) EndComprehension(c *Comprehension) *Builder {
	c.End = len(p.code.data)
	return p
}

// -----------------------------------------------------------------------------

// Append instr
func (p *Builder) Append(typ reflect.Type, arity int) *Builder {
	if arity < 0 {
		arity = bitsFuncvArityVar
	}
	i := (opAppend << bitsOpShift) | uint32(arity)
	p.code.data = append(p.code.data, i)
	return p
}

// MakeArray instr
func (p *Builder) MakeArray(typ reflect.Type, arity int) *Builder {
	if arity < 0 {
		if typ.Kind() == reflect.Array {
			log.Panicln("MakeArray failed: can't be variadic.")
		}
		arity = bitsFuncvArityVar
	} else if arity >= bitsFuncvArityMax {
		p.Push(arity - bitsFuncvArityMax)
		arity = bitsFuncvArityMax
	}
	i := (opMakeArray << bitsOpShift) | (uint32(arity) << bitsOpCallFuncvShift) | p.newType(typ)
	p.code.data = append(p.code.data, i)
	return p
}

// MakeMap instr
func (p *Builder) MakeMap(typ reflect.Type, arity int) *Builder {
	if arity < 0 {
		log.Panicln("MakeMap failed: can't be variadic.")
	} else if arity >= bitsFuncvArityMax {
		p.Push(arity - bitsFuncvArityMax)
		arity = bitsFuncvArityMax
	}
	i := (opMakeMap << bitsOpShift) | (uint32(arity) << bitsOpCallFuncvShift) | p.newType(typ)
	p.code.data = append(p.code.data, i)
	return p
}

// Make instr
func (p *Builder) Make(typ reflect.Type, arity int) *Builder {
	i := (opMake << bitsOpShift) | (uint32(arity) << bitsOpCallFuncvShift) | p.newType(typ)
	p.code.data = append(p.code.data, i)
	return p
}

// TypeCast instr
func (p *Builder) TypeCast(from, to reflect.Type) *Builder {
	i := (opTypeCast << bitsOpShift) | p.requireType(to)
	p.code.data = append(p.code.data, i)
	return p
}

// Zero instr
func (p *Builder) Zero(typ reflect.Type) *Builder {
	i := (opZero << bitsOpShift) | p.requireType(typ)
	p.code.data = append(p.code.data, i)
	return p
}

func (p *Builder) requireType(typ reflect.Type) uint32 {
	kind := typ.Kind()
	bt := builtinTypes[kind]
	if bt.size > 0 {
		return uint32(kind)
	}
	return p.newType(typ)
}

func (p *Builder) newType(typ reflect.Type) uint32 {
	if ityp, ok := p.types[typ]; ok {
		return ityp
	}
	code := p.code
	ityp := uint32(len(code.types) + len(builtinTypes))
	code.types = append(code.types, typ)
	p.types[typ] = ityp
	return ityp
}

func getType(ityp uint32, ctx *Context) reflect.Type {
	if ityp < uint32(len(builtinTypes)) {
		return builtinTypes[ityp].typ
	}
	return ctx.code.types[ityp-uint32(len(builtinTypes))]
}

// -----------------------------------------------------------------------------
