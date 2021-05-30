// Package context provide Go+ "context" package, as "context" package in Go.
package context

import (
	context "context"
	reflect "reflect"
	time "time"

	gop "github.com/goplus/gop"
)

func execBackground(_ int, p *gop.Context) {
	ret0 := context.Background()
	p.Ret(0, ret0)
}

func execiContextDeadline(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(context.Context).Deadline()
	p.Ret(1, ret0, ret1)
}

func execiContextDone(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(context.Context).Done()
	p.Ret(1, ret0)
}

func execiContextErr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(context.Context).Err()
	p.Ret(1, ret0)
}

func execiContextValue(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(context.Context).Value(args[1])
	p.Ret(2, ret0)
}

func execTODO(_ int, p *gop.Context) {
	ret0 := context.TODO()
	p.Ret(0, ret0)
}

func toType0(v interface{}) context.Context {
	if v == nil {
		return nil
	}
	return v.(context.Context)
}

func execWithCancel(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := context.WithCancel(toType0(args[0]))
	p.Ret(1, ret0, ret1)
}

func execWithDeadline(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := context.WithDeadline(toType0(args[0]), args[1].(time.Time))
	p.Ret(2, ret0, ret1)
}

func execWithTimeout(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := context.WithTimeout(toType0(args[0]), args[1].(time.Duration))
	p.Ret(2, ret0, ret1)
}

func execWithValue(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := context.WithValue(toType0(args[0]), args[1], args[2])
	p.Ret(3, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("context")

func init() {
	I.RegisterFuncs(
		I.Func("Background", context.Background, execBackground),
		I.Func("(Context).Deadline", (context.Context).Deadline, execiContextDeadline),
		I.Func("(Context).Done", (context.Context).Done, execiContextDone),
		I.Func("(Context).Err", (context.Context).Err, execiContextErr),
		I.Func("(Context).Value", (context.Context).Value, execiContextValue),
		I.Func("TODO", context.TODO, execTODO),
		I.Func("WithCancel", context.WithCancel, execWithCancel),
		I.Func("WithDeadline", context.WithDeadline, execWithDeadline),
		I.Func("WithTimeout", context.WithTimeout, execWithTimeout),
		I.Func("WithValue", context.WithValue, execWithValue),
	)
	I.RegisterVars(
		I.Var("Canceled", &context.Canceled),
		I.Var("DeadlineExceeded", &context.DeadlineExceeded),
	)
	I.RegisterTypes(
		I.Type("CancelFunc", reflect.TypeOf((*context.CancelFunc)(nil)).Elem()),
		I.Type("Context", reflect.TypeOf((*context.Context)(nil)).Elem()),
	)
}
