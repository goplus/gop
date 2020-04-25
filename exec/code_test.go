package exec

// -----------------------------------------------------------------------------

func checkPop(ctx *Context) interface{} {
	v, ok := ctx.Pop()
	if !ok {
		panic("checkPop failed: no data")
	}
	_, ok = ctx.Pop()
	if ok {
		panic("checkPop failed: too many data")
	}
	return v
}

// -----------------------------------------------------------------------------
