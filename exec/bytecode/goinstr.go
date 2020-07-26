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

	"github.com/goplus/gop/exec.spec"
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
	data := reflect.ValueOf(p.Pop())
	switch data.Kind() {
	case reflect.Map:
		c.execMapRange(data, p)
	default:
		c.execListRange(data, p)
	}
}

func (c *ForPhrase) execListRange(data reflect.Value, ctx *Context) {
	data = reflect.Indirect(data)
	var n = data.Len()
	var ip, ipCond, ipEnd = ctx.ip, c.Cond, c.End
	var key, val = c.Key, c.Value
	var blockScope = c.block != nil
	var old savedScopeCtx
Loop:
	for i := 0; i < n; i++ {
		if key != nil {
			ctx.setVar(key.idx, i)
		}
		if val != nil {
			ctx.setVar(val.idx, data.Index(i).Interface())
		}
		if blockScope { // TODO: move out of `for` statement
			parent := ctx.varScope
			old = ctx.switchScope(&parent, &c.block.varManager)
		}
		if ipCond > 0 {
			ctx.Exec(ip, ipCond)
			if ok := ctx.Pop().(bool); ok {
				ctx.Exec(ipCond, ipEnd)
			}
		} else {
			ctx.Exec(ip, ipEnd)
		}
		if blockScope {
			ctx.restoreScope(old)
		}
		switch ctx.ip {
		case iReturn, ipReturnN:
			return
		case iBreak:
			break Loop
		case iContinue:
			continue
		}
	}
	ctx.ip = ipEnd
}

func (c *ForPhrase) execMapRange(data reflect.Value, ctx *Context) {
	var iter = data.MapRange()
	var ip, ipCond, ipEnd = ctx.ip, c.Cond, c.End
	var key, val = c.Key, c.Value
	var blockScope = c.block != nil
	var old savedScopeCtx
Loop:
	for iter.Next() {
		if key != nil {
			ctx.setVar(key.idx, iter.Key().Interface())
		}
		if val != nil {
			ctx.setVar(val.idx, iter.Value().Interface())
		}
		if blockScope {
			parent := ctx.varScope
			old = ctx.switchScope(&parent, &c.block.varManager)
		}
		if ipCond > 0 {
			ctx.Exec(ip, ipCond)
			if ok := ctx.Pop().(bool); ok {
				ctx.Exec(ipCond, ipEnd)
			}
		} else {
			ctx.Exec(ip, ipEnd)
		}
		if blockScope {
			ctx.restoreScope(old)
		}
		switch ctx.ip {
		case iReturn, ipReturnN:
			return
		case iBreak:
			break Loop
		case iContinue:
			continue
		}
	}
	ctx.ip = ipEnd
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
	var ret, set reflect.Value
	kind := typSlice.Kind()
	if kind == reflect.Slice {
		ret = reflect.MakeSlice(typSlice, arity, arity)
		set = ret
	} else if kind == reflect.Array {
		ret = reflect.New(typSlice)
		set = ret.Elem()
	} else {
		log.Panic("makeArray bad type:", typSlice)
	}
	for i, arg := range args {
		set.Index(i).Set(getElementOf(arg, typSlice))
	}
	if kind == reflect.Slice {
		p.Ret(arity, ret.Interface())
	} else {
		p.Ret(arity, ret.Elem().Interface())
	}
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

func execIndex(i Instr, p *Context) {
	idx := int(i & setIndexOperand)
	if idx == setIndexOperand {
		idx = p.Pop().(int)
	}
	n := len(p.data)
	v := reflect.Indirect(reflect.ValueOf(p.data[n-1])).Index(idx)
	if (i & setIndexFlag) != 0 { // value sliceData $idx $setIndex
		v.Set(reflect.ValueOf(p.data[n-2]))
		p.PopN(2)
	} else { // sliceData $idx $setIndex
		p.data[n-1] = v.Interface()
	}
}

func execMapIndex(i Instr, p *Context) {
	n := len(p.data)
	key := reflect.ValueOf(p.data[n-1])
	v := reflect.ValueOf(p.data[n-2])
	switch i & bitsOperand {
	case 1: // value mapData $key $setMapIndex
		v.SetMapIndex(key, reflect.ValueOf(p.data[n-3]))
		p.PopN(3)
	default: // mapData $key $mapIndex
		value := v.MapIndex(key)
		if !value.IsValid() {
			value = reflect.Zero(v.Type().Elem())
		}
		p.Ret(2, value.Interface())
	}
}

func popSliceIndexs(instr Instr, p *Context) (i, j int) {
	instr &= bitsOperand
	i = int(instr >> 13)
	j = int(instr & sliceIndexMask)
	if j == sliceIndexMask {
		j = p.Pop().(int)
	} else if j == sliceIndexMask-1 {
		j = -2
	}
	if i == sliceIndexMask {
		i = p.Pop().(int)
	} else if i == sliceIndexMask-1 {
		i = 0
	}
	return
}

func execSlice(instr Instr, p *Context) {
	i, j := popSliceIndexs(instr, p)
	n := len(p.data)
	v := reflect.Indirect(reflect.ValueOf(p.data[n-1]))
	if j == -2 {
		j = v.Len()
	}
	p.data[n-1] = v.Slice(i, j).Interface()
}

func execSlice3(instr Instr, p *Context) {
	k := p.Pop().(int)
	i, j := popSliceIndexs(instr, p)
	n := len(p.data)
	v := reflect.ValueOf(p.data[n-1])
	p.data[n-1] = v.Slice3(i, j, k).Interface()
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
func (p *Builder) ForPhrase(f *ForPhrase, key, val *Var, hasExecCtx ...bool) *Builder {
	f.Key, f.Value = key, val
	if key != nil {
		p.DefineVar(key)
	}
	if val != nil {
		p.DefineVar(val)
	}
	if hasExecCtx != nil && hasExecCtx[0] {
		f.block = newBlockCtx(p.nestDepth+1, p.varManager)
		p.varManager = &f.block.varManager
		log.Debug("ForPhrase:", f.block.nestDepth)
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
	if arity > 2 {
		panic("make arity > 2")
	}
	i := (opMake << bitsOpShift) | (uint32(arity) << bitsOpCallFuncvShift) | p.newType(typ)
	p.code.data = append(p.code.data, i)
	return p
}

// MapIndex instr
func (p *Builder) MapIndex() *Builder {
	p.code.data = append(p.code.data, opMapIndex<<bitsOpShift)
	return p
}

// SetMapIndex instr
func (p *Builder) SetMapIndex() *Builder {
	p.code.data = append(p.code.data, (opMapIndex<<bitsOpShift)|1)
	return p
}

// Index instr
func (p *Builder) Index(idx int) *Builder {
	if idx >= setIndexOperand {
		p.Push(idx)
		idx = -1
	}
	i := (opIndex << bitsOpShift) | uint32(idx&setIndexOperand)
	p.code.data = append(p.code.data, i)
	return p
}

// SetIndex instr
func (p *Builder) SetIndex(idx int) *Builder {
	if idx >= setIndexOperand {
		p.Push(idx)
		idx = -1
	}
	i := (opIndex<<bitsOpShift | setIndexFlag) | uint32(idx&setIndexOperand)
	p.code.data = append(p.code.data, i)
	return p
}

const (
	setIndexFlag    = (1 << 25)
	setIndexOperand = setIndexFlag - 1
	sliceIndexMask  = (1 << 13) - 1
	// SliceConstIndexLast - slice const index max
	SliceConstIndexLast = exec.SliceConstIndexLast
	// SliceDefaultIndex - unspecified index
	SliceDefaultIndex = exec.SliceDefaultIndex
)

// Slice instr
func (p *Builder) Slice(i, j int) *Builder { // i = -1, -2
	if i > SliceConstIndexLast {
		panic("i > SliceConstIndexLast")
	}
	if j > SliceConstIndexLast {
		p.Push(j)
		j = -1
	}
	instr := (opSlice << bitsOpShift) | uint32(i&sliceIndexMask)<<13 | uint32(j&sliceIndexMask)
	p.code.data = append(p.code.data, instr)
	return p
}

// Slice3 instr
func (p *Builder) Slice3(i, j, k int) *Builder {
	if i > SliceConstIndexLast {
		panic("i > SliceConstIndexLast")
	}
	if k == -2 || j == -2 || j > SliceConstIndexLast {
		panic("k == SliceDefaultIndex || j == SliceDefaultIndex || j > SliceConstIndexLast")
	}
	if k >= 0 {
		p.Push(k)
	}
	instr := (opSlice3 << bitsOpShift) | uint32(i&sliceIndexMask)<<13 | uint32(j&sliceIndexMask)
	p.code.data = append(p.code.data, instr)
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
	if exec.SizeofKind(kind) > 0 {
		return uint32(kind)
	}
	return p.newType(typ)
}

func (p *Builder) newType(typ reflect.Type) uint32 {
	if ityp, ok := p.types[typ]; ok {
		return ityp
	}
	code := p.code
	ityp := uint32(len(code.types) + exec.BuiltinTypesLen)
	code.types = append(code.types, typ)
	p.types[typ] = ityp
	return ityp
}

func getType(ityp uint32, ctx *Context) reflect.Type {
	if ityp < uint32(exec.BuiltinTypesLen) {
		return exec.TypeFromKind(exec.Kind(ityp))
	}
	return ctx.code.types[ityp-uint32(exec.BuiltinTypesLen)]
}

// -----------------------------------------------------------------------------
