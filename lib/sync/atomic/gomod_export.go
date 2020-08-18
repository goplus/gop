// Package atomic provide Go+ "sync/atomic" package, as "sync/atomic" package in Go.
package atomic

import (
	reflect "reflect"
	atomic "sync/atomic"
	unsafe "unsafe"

	gop "github.com/goplus/gop"
)

func execAddInt32(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := atomic.AddInt32(args[0].(*int32), args[1].(int32))
	p.Ret(2, ret0)
}

func execAddInt64(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := atomic.AddInt64(args[0].(*int64), args[1].(int64))
	p.Ret(2, ret0)
}

func execAddUint32(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := atomic.AddUint32(args[0].(*uint32), args[1].(uint32))
	p.Ret(2, ret0)
}

func execAddUint64(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := atomic.AddUint64(args[0].(*uint64), args[1].(uint64))
	p.Ret(2, ret0)
}

func execAddUintptr(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := atomic.AddUintptr(args[0].(*uintptr), args[1].(uintptr))
	p.Ret(2, ret0)
}

func execCompareAndSwapInt32(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := atomic.CompareAndSwapInt32(args[0].(*int32), args[1].(int32), args[2].(int32))
	p.Ret(3, ret0)
}

func execCompareAndSwapInt64(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := atomic.CompareAndSwapInt64(args[0].(*int64), args[1].(int64), args[2].(int64))
	p.Ret(3, ret0)
}

func execCompareAndSwapPointer(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := atomic.CompareAndSwapPointer(args[0].(*unsafe.Pointer), args[1].(unsafe.Pointer), args[2].(unsafe.Pointer))
	p.Ret(3, ret0)
}

func execCompareAndSwapUint32(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := atomic.CompareAndSwapUint32(args[0].(*uint32), args[1].(uint32), args[2].(uint32))
	p.Ret(3, ret0)
}

func execCompareAndSwapUint64(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := atomic.CompareAndSwapUint64(args[0].(*uint64), args[1].(uint64), args[2].(uint64))
	p.Ret(3, ret0)
}

func execCompareAndSwapUintptr(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := atomic.CompareAndSwapUintptr(args[0].(*uintptr), args[1].(uintptr), args[2].(uintptr))
	p.Ret(3, ret0)
}

func execLoadInt32(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := atomic.LoadInt32(args[0].(*int32))
	p.Ret(1, ret0)
}

func execLoadInt64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := atomic.LoadInt64(args[0].(*int64))
	p.Ret(1, ret0)
}

func execLoadPointer(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := atomic.LoadPointer(args[0].(*unsafe.Pointer))
	p.Ret(1, ret0)
}

func execLoadUint32(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := atomic.LoadUint32(args[0].(*uint32))
	p.Ret(1, ret0)
}

func execLoadUint64(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := atomic.LoadUint64(args[0].(*uint64))
	p.Ret(1, ret0)
}

func execLoadUintptr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := atomic.LoadUintptr(args[0].(*uintptr))
	p.Ret(1, ret0)
}

func execStoreInt32(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	atomic.StoreInt32(args[0].(*int32), args[1].(int32))
	p.PopN(2)
}

func execStoreInt64(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	atomic.StoreInt64(args[0].(*int64), args[1].(int64))
	p.PopN(2)
}

func execStorePointer(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	atomic.StorePointer(args[0].(*unsafe.Pointer), args[1].(unsafe.Pointer))
	p.PopN(2)
}

func execStoreUint32(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	atomic.StoreUint32(args[0].(*uint32), args[1].(uint32))
	p.PopN(2)
}

func execStoreUint64(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	atomic.StoreUint64(args[0].(*uint64), args[1].(uint64))
	p.PopN(2)
}

func execStoreUintptr(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	atomic.StoreUintptr(args[0].(*uintptr), args[1].(uintptr))
	p.PopN(2)
}

func execSwapInt32(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := atomic.SwapInt32(args[0].(*int32), args[1].(int32))
	p.Ret(2, ret0)
}

func execSwapInt64(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := atomic.SwapInt64(args[0].(*int64), args[1].(int64))
	p.Ret(2, ret0)
}

func execSwapPointer(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := atomic.SwapPointer(args[0].(*unsafe.Pointer), args[1].(unsafe.Pointer))
	p.Ret(2, ret0)
}

func execSwapUint32(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := atomic.SwapUint32(args[0].(*uint32), args[1].(uint32))
	p.Ret(2, ret0)
}

func execSwapUint64(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := atomic.SwapUint64(args[0].(*uint64), args[1].(uint64))
	p.Ret(2, ret0)
}

func execSwapUintptr(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := atomic.SwapUintptr(args[0].(*uintptr), args[1].(uintptr))
	p.Ret(2, ret0)
}

func execmValueLoad(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*atomic.Value).Load()
	p.Ret(1, ret0)
}

func execmValueStore(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*atomic.Value).Store(args[1])
	p.PopN(2)
}

// I is a Go package instance.
var I = gop.NewGoPackage("sync/atomic")

func init() {
	I.RegisterFuncs(
		I.Func("AddInt32", atomic.AddInt32, execAddInt32),
		I.Func("AddInt64", atomic.AddInt64, execAddInt64),
		I.Func("AddUint32", atomic.AddUint32, execAddUint32),
		I.Func("AddUint64", atomic.AddUint64, execAddUint64),
		I.Func("AddUintptr", atomic.AddUintptr, execAddUintptr),
		I.Func("CompareAndSwapInt32", atomic.CompareAndSwapInt32, execCompareAndSwapInt32),
		I.Func("CompareAndSwapInt64", atomic.CompareAndSwapInt64, execCompareAndSwapInt64),
		I.Func("CompareAndSwapPointer", atomic.CompareAndSwapPointer, execCompareAndSwapPointer),
		I.Func("CompareAndSwapUint32", atomic.CompareAndSwapUint32, execCompareAndSwapUint32),
		I.Func("CompareAndSwapUint64", atomic.CompareAndSwapUint64, execCompareAndSwapUint64),
		I.Func("CompareAndSwapUintptr", atomic.CompareAndSwapUintptr, execCompareAndSwapUintptr),
		I.Func("LoadInt32", atomic.LoadInt32, execLoadInt32),
		I.Func("LoadInt64", atomic.LoadInt64, execLoadInt64),
		I.Func("LoadPointer", atomic.LoadPointer, execLoadPointer),
		I.Func("LoadUint32", atomic.LoadUint32, execLoadUint32),
		I.Func("LoadUint64", atomic.LoadUint64, execLoadUint64),
		I.Func("LoadUintptr", atomic.LoadUintptr, execLoadUintptr),
		I.Func("StoreInt32", atomic.StoreInt32, execStoreInt32),
		I.Func("StoreInt64", atomic.StoreInt64, execStoreInt64),
		I.Func("StorePointer", atomic.StorePointer, execStorePointer),
		I.Func("StoreUint32", atomic.StoreUint32, execStoreUint32),
		I.Func("StoreUint64", atomic.StoreUint64, execStoreUint64),
		I.Func("StoreUintptr", atomic.StoreUintptr, execStoreUintptr),
		I.Func("SwapInt32", atomic.SwapInt32, execSwapInt32),
		I.Func("SwapInt64", atomic.SwapInt64, execSwapInt64),
		I.Func("SwapPointer", atomic.SwapPointer, execSwapPointer),
		I.Func("SwapUint32", atomic.SwapUint32, execSwapUint32),
		I.Func("SwapUint64", atomic.SwapUint64, execSwapUint64),
		I.Func("SwapUintptr", atomic.SwapUintptr, execSwapUintptr),
		I.Func("(*Value).Load", (*atomic.Value).Load, execmValueLoad),
		I.Func("(*Value).Store", (*atomic.Value).Store, execmValueStore),
	)
	I.RegisterTypes(
		I.Type("Value", reflect.TypeOf((*atomic.Value)(nil)).Elem()),
	)
}
