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
	"testing"

	"github.com/goplus/gop/exec.spec"
	"github.com/qiniu/x/log"
	"github.com/qiniu/x/ts"
)

var (
	// TyBool type
	TyBool = exec.TyBool
	// TyInt type
	TyInt = exec.TyInt
	// TyInt8 type
	TyInt8 = exec.TyInt8
	// TyInt16 type
	TyInt16 = exec.TyInt16
	// TyInt32 type
	TyInt32 = exec.TyInt32
	// TyInt64 type
	TyInt64 = exec.TyInt64
	// TyUint type
	TyUint = exec.TyUint
	// TyUint8 type
	TyUint8 = exec.TyUint8
	// TyUint16 type
	TyUint16 = exec.TyUint16
	// TyUint32 type
	TyUint32 = exec.TyUint32
	// TyUint64 type
	TyUint64 = exec.TyUint64
	// TyUintptr type
	TyUintptr = exec.TyUintptr
	// TyFloat32 type
	TyFloat32 = exec.TyFloat32
	// TyFloat64 type
	TyFloat64 = exec.TyFloat64
	// TyComplex64 type
	TyComplex64 = exec.TyComplex64
	// TyComplex128 type
	TyComplex128 = exec.TyComplex128
	// TyString type
	TyString = exec.TyString
	// TyUnsafePointer type
	TyUnsafePointer = exec.TyUnsafePointer
	// TyEmptyInterface type
	TyEmptyInterface = exec.TyEmptyInterface
	// TyError type
	TyError = exec.TyError
)

var (
	// TyByte type
	TyByte = exec.TyByte
	// TyRune type
	TyRune = exec.TyRune
	// TyEmptyInterfaceSlice type
	TyEmptyInterfaceSlice = exec.TyEmptyInterfaceSlice
)

// -----------------------------------------------------------------------------

var (
	defaultImpl = NewPackage(nil)
)

func init() {
	SetProfile(true)
}

func newBuilder() exec.Builder {
	return NewBuilder(nil).Interface()
}

func newFunc(name string, nestDepth uint32) exec.FuncInfo {
	return defaultImpl.NewFunc(name, nestDepth)
}

func setPackage(f exec.FuncInfo, code exec.Code) {
	f.(*iFuncInfo).Pkg = NewPackage(code.(*Code))
}

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

func expect(t *testing.T, f func(), expected string, panicMsg ...interface{}) {
	e := ts.StartExpecting(t, ts.CapStdout)
	defer e.Close()
	e.Call(f).Panic(panicMsg...).Expect(expected)
}

// -----------------------------------------------------------------------------
