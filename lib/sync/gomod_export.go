// Package sync provide Go+ "sync" package, as "sync" package in Go.
package sync

import (
	sync "sync"

	gop "github.com/qiniu/goplus/gop"
)

func execCondWait(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*sync.Cond).Wait()
	p.PopN(1)
}

func execCondSignal(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*sync.Cond).Signal()
	p.PopN(1)
}

func execCondBroadcast(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*sync.Cond).Broadcast()
	p.PopN(1)
}

func execMapLoad(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*sync.Map).Load(args[1])
	p.Ret(2, ret0, ret1)
}

func execMapStore(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*sync.Map).Store(args[1], args[2])
	p.PopN(3)
}

func execMapLoadOrStore(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*sync.Map).LoadOrStore(args[1], args[2])
	p.Ret(3, ret0, ret1)
}

func execMapDelete(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*sync.Map).Delete(args[1])
	p.PopN(2)
}

func execMapRange(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*sync.Map).Range(args[1].(func(key interface{}, value interface{}) bool))
	p.PopN(2)
}

func execMutexLock(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*sync.Mutex).Lock()
	p.PopN(1)
}

func execMutexUnlock(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*sync.Mutex).Unlock()
	p.PopN(1)
}

func toType0(v interface{}) sync.Locker {
	if v == nil {
		return nil
	}
	return v.(sync.Locker)
}

func execNewCond(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := sync.NewCond(toType0(args[0]))
	p.Ret(1, ret0)
}

func execOnceDo(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*sync.Once).Do(args[1].(func()))
	p.PopN(2)
}

func execPoolPut(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*sync.Pool).Put(args[1])
	p.PopN(2)
}

func execPoolGet(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*sync.Pool).Get()
	p.Ret(1, ret0)
}

func execRWMutexRLock(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*sync.RWMutex).RLock()
	p.PopN(1)
}

func execRWMutexRUnlock(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*sync.RWMutex).RUnlock()
	p.PopN(1)
}

func execRWMutexLock(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*sync.RWMutex).Lock()
	p.PopN(1)
}

func execRWMutexUnlock(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*sync.RWMutex).Unlock()
	p.PopN(1)
}

func execRWMutexRLocker(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*sync.RWMutex).RLocker()
	p.Ret(1, ret0)
}

func execWaitGroupAdd(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*sync.WaitGroup).Add(args[1].(int))
	p.PopN(2)
}

func execWaitGroupDone(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*sync.WaitGroup).Done()
	p.PopN(1)
}

func execWaitGroupWait(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*sync.WaitGroup).Wait()
	p.PopN(1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("sync")

func init() {
	I.RegisterFuncs(
		I.Func("(*Cond).Wait", (*sync.Cond).Wait, execCondWait),
		I.Func("(*Cond).Signal", (*sync.Cond).Signal, execCondSignal),
		I.Func("(*Cond).Broadcast", (*sync.Cond).Broadcast, execCondBroadcast),
		I.Func("(*Map).Load", (*sync.Map).Load, execMapLoad),
		I.Func("(*Map).Store", (*sync.Map).Store, execMapStore),
		I.Func("(*Map).LoadOrStore", (*sync.Map).LoadOrStore, execMapLoadOrStore),
		I.Func("(*Map).Delete", (*sync.Map).Delete, execMapDelete),
		I.Func("(*Map).Range", (*sync.Map).Range, execMapRange),
		I.Func("(*Mutex).Lock", (*sync.Mutex).Lock, execMutexLock),
		I.Func("(*Mutex).Unlock", (*sync.Mutex).Unlock, execMutexUnlock),
		I.Func("NewCond", sync.NewCond, execNewCond),
		I.Func("(*Once).Do", (*sync.Once).Do, execOnceDo),
		I.Func("(*Pool).Put", (*sync.Pool).Put, execPoolPut),
		I.Func("(*Pool).Get", (*sync.Pool).Get, execPoolGet),
		I.Func("(*RWMutex).RLock", (*sync.RWMutex).RLock, execRWMutexRLock),
		I.Func("(*RWMutex).RUnlock", (*sync.RWMutex).RUnlock, execRWMutexRUnlock),
		I.Func("(*RWMutex).Lock", (*sync.RWMutex).Lock, execRWMutexLock),
		I.Func("(*RWMutex).Unlock", (*sync.RWMutex).Unlock, execRWMutexUnlock),
		I.Func("(*RWMutex).RLocker", (*sync.RWMutex).RLocker, execRWMutexRLocker),
		I.Func("(*WaitGroup).Add", (*sync.WaitGroup).Add, execWaitGroupAdd),
		I.Func("(*WaitGroup).Done", (*sync.WaitGroup).Done, execWaitGroupDone),
		I.Func("(*WaitGroup).Wait", (*sync.WaitGroup).Wait, execWaitGroupWait),
	)
}
