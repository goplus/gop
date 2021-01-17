// Package sync provide Go+ "sync" package, as "sync" package in Go.
package sync

import (
	reflect "reflect"
	sync "sync"

	gop "github.com/goplus/gop"
)

func execmCondWait(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*sync.Cond).Wait()
	p.PopN(1)
}

func execmCondSignal(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*sync.Cond).Signal()
	p.PopN(1)
}

func execmCondBroadcast(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*sync.Cond).Broadcast()
	p.PopN(1)
}

func execiLockerLock(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(sync.Locker).Lock()
	p.PopN(1)
}

func execiLockerUnlock(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(sync.Locker).Unlock()
	p.PopN(1)
}

func execmMapLoad(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*sync.Map).Load(args[1])
	p.Ret(2, ret0, ret1)
}

func execmMapStore(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*sync.Map).Store(args[1], args[2])
	p.PopN(3)
}

func execmMapLoadOrStore(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*sync.Map).LoadOrStore(args[1], args[2])
	p.Ret(3, ret0, ret1)
}

func execmMapDelete(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*sync.Map).Delete(args[1])
	p.PopN(2)
}

func execmMapRange(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*sync.Map).Range(args[1].(func(key interface{}, value interface{}) bool))
	p.PopN(2)
}

func execmMutexLock(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*sync.Mutex).Lock()
	p.PopN(1)
}

func execmMutexUnlock(_ int, p *gop.Context) {
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

func execmOnceDo(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*sync.Once).Do(args[1].(func()))
	p.PopN(2)
}

func execmPoolPut(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*sync.Pool).Put(args[1])
	p.PopN(2)
}

func execmPoolGet(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*sync.Pool).Get()
	p.Ret(1, ret0)
}

func execmRWMutexRLock(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*sync.RWMutex).RLock()
	p.PopN(1)
}

func execmRWMutexRUnlock(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*sync.RWMutex).RUnlock()
	p.PopN(1)
}

func execmRWMutexLock(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*sync.RWMutex).Lock()
	p.PopN(1)
}

func execmRWMutexUnlock(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*sync.RWMutex).Unlock()
	p.PopN(1)
}

func execmRWMutexRLocker(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*sync.RWMutex).RLocker()
	p.Ret(1, ret0)
}

func execmWaitGroupAdd(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*sync.WaitGroup).Add(args[1].(int))
	p.PopN(2)
}

func execmWaitGroupDone(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*sync.WaitGroup).Done()
	p.PopN(1)
}

func execmWaitGroupWait(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*sync.WaitGroup).Wait()
	p.PopN(1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("sync")

func init() {
	I.RegisterFuncs(
		I.Func("(*Cond).Wait", (*sync.Cond).Wait, execmCondWait),
		I.Func("(*Cond).Signal", (*sync.Cond).Signal, execmCondSignal),
		I.Func("(*Cond).Broadcast", (*sync.Cond).Broadcast, execmCondBroadcast),
		I.Func("(Locker).Lock", (sync.Locker).Lock, execiLockerLock),
		I.Func("(Locker).Unlock", (sync.Locker).Unlock, execiLockerUnlock),
		I.Func("(*Map).Load", (*sync.Map).Load, execmMapLoad),
		I.Func("(*Map).Store", (*sync.Map).Store, execmMapStore),
		I.Func("(*Map).LoadOrStore", (*sync.Map).LoadOrStore, execmMapLoadOrStore),
		I.Func("(*Map).Delete", (*sync.Map).Delete, execmMapDelete),
		I.Func("(*Map).Range", (*sync.Map).Range, execmMapRange),
		I.Func("(*Mutex).Lock", (*sync.Mutex).Lock, execmMutexLock),
		I.Func("(*Mutex).Unlock", (*sync.Mutex).Unlock, execmMutexUnlock),
		I.Func("NewCond", sync.NewCond, execNewCond),
		I.Func("(*Once).Do", (*sync.Once).Do, execmOnceDo),
		I.Func("(*Pool).Put", (*sync.Pool).Put, execmPoolPut),
		I.Func("(*Pool).Get", (*sync.Pool).Get, execmPoolGet),
		I.Func("(*RWMutex).RLock", (*sync.RWMutex).RLock, execmRWMutexRLock),
		I.Func("(*RWMutex).RUnlock", (*sync.RWMutex).RUnlock, execmRWMutexRUnlock),
		I.Func("(*RWMutex).Lock", (*sync.RWMutex).Lock, execmRWMutexLock),
		I.Func("(*RWMutex).Unlock", (*sync.RWMutex).Unlock, execmRWMutexUnlock),
		I.Func("(*RWMutex).RLocker", (*sync.RWMutex).RLocker, execmRWMutexRLocker),
		I.Func("(*WaitGroup).Add", (*sync.WaitGroup).Add, execmWaitGroupAdd),
		I.Func("(*WaitGroup).Done", (*sync.WaitGroup).Done, execmWaitGroupDone),
		I.Func("(*WaitGroup).Wait", (*sync.WaitGroup).Wait, execmWaitGroupWait),
	)
	I.RegisterTypes(
		I.Type("Cond", reflect.TypeOf((*sync.Cond)(nil)).Elem()),
		I.Type("Locker", reflect.TypeOf((*sync.Locker)(nil)).Elem()),
		I.Type("Map", reflect.TypeOf((*sync.Map)(nil)).Elem()),
		I.Type("Mutex", reflect.TypeOf((*sync.Mutex)(nil)).Elem()),
		I.Type("Once", reflect.TypeOf((*sync.Once)(nil)).Elem()),
		I.Type("Pool", reflect.TypeOf((*sync.Pool)(nil)).Elem()),
		I.Type("RWMutex", reflect.TypeOf((*sync.RWMutex)(nil)).Elem()),
		I.Type("WaitGroup", reflect.TypeOf((*sync.WaitGroup)(nil)).Elem()),
	)
}
