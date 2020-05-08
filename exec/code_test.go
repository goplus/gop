package exec

import "github.com/qiniu/x/log"

// -----------------------------------------------------------------------------

func checkPop(ctx *Context) interface{} {
	if ctx.Len() < 1 {
		log.Panicln("checkPop failed: no data.")
	}
	v := ctx.Get(-1)
	ctx.PopN(1)
	if ctx.Len() > 0 {
		log.Panicln("checkPop failed: too many data:", ctx.Len())
	}
	return v
}

// -----------------------------------------------------------------------------
