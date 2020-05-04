package exec

// -----------------------------------------------------------------------------

func checkPop(ctx *Context) interface{} {
	if ctx.Len() < 1 {
		panic("checkPop failed: no data")
	}
	v := ctx.Get(-1)
	ctx.PopN(1)
	if ctx.Len() > 0 {
		panic("checkPop failed: too many data")
	}
	return v
}

// -----------------------------------------------------------------------------
