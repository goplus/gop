package exec

import (
	"fmt"
	"reflect"

	qlang "qlang.io/spec"
)

// -----------------------------------------------------------------------------
// ForRange

type iForRange struct {
	args  []int
	start int
	end   int
}

const (
	BreakForRange    = -1
	ContinueForRange = -2
)

func (p *iForRange) execBody(stk *Stack, ctx *Context) {

	data := ctx.Code.data
	ipEnd := p.end

	ctx.ip = p.start
	for ctx.ip != ipEnd {
		instr := data[ctx.ip]
		ctx.ip++
		instr.Exec(stk, ctx)
		if ctx.ip < 0 { // continue or break?
			return
		}
	}
}

func (p *iForRange) Exec(stk *Stack, ctx *Context) {

	val, ok := stk.Pop()
	if !ok {
		panic("unexpected")
	}

	done := ctx.ip
	args := p.args
	narg := len(args)

	if ch, ok := val.(*qlang.Chan); ok {
		if narg > 1 {
			panic("too many variables in range")
		}
		v := ch.Data
		for {
			x, ok := v.Recv()
			if !ok {
				break
			}
			if narg > 0 {
				item := args[0]
				varRef := ctx.getVarRef(item)
				*varRef = x.Interface()
			}
			p.execBody(stk, ctx)
			if ctx.ip == BreakForRange {
				break
			}
			if ctx.ip == ContinueForRange {
				continue
			}
		}
		ctx.ip = done
		return
	}

	v := reflect.ValueOf(val)
	kind := v.Kind()
	switch kind {
	case reflect.Slice:
		n := v.Len()
		for i := 0; i < n; i++ {
			if narg > 0 {
				index := args[0]
				varRef := ctx.getVarRef(index)
				*varRef = i
			}
			if narg > 1 {
				item := args[1]
				varRef := ctx.getVarRef(item)
				*varRef = v.Index(i).Interface()
			}
			p.execBody(stk, ctx)
			if ctx.ip == BreakForRange {
				break
			}
			if ctx.ip == ContinueForRange {
				continue
			}
		}
	case reflect.Map:
		keys := v.MapKeys()
		for _, key := range keys {
			if narg > 0 {
				index := args[0]
				varRef := ctx.getVarRef(index)
				*varRef = key.Interface()
			}
			if narg > 1 {
				item := args[1]
				varRef := ctx.getVarRef(item)
				*varRef = v.MapIndex(key).Interface()
			}
			p.execBody(stk, ctx)
			if ctx.ip == BreakForRange {
				break
			}
			if ctx.ip == ContinueForRange {
				continue
			}
		}
	default:
		panic(fmt.Errorf("type `%v` doesn't support `range`", v.Type()))
	}
	ctx.ip = done
}

// ForRange returns an instruction that represents for..range
//
func ForRange(args []int, start, end int) Instr {

	return &iForRange{args, start, end}
}

// -----------------------------------------------------------------------------
