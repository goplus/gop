/*
 Copyright 2020 Qiniu Cloud (qiniu.com)

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

package builtin

import (
	"fmt"
	"io"
	"reflect"

	"github.com/qiniu/qlang/v6/exec"
	qlang "github.com/qiniu/qlang/v6/spec"
)

// -----------------------------------------------------------------------------

// QexecPrint instr
func QexecPrint(arity int, p *qlang.Context) {
	args := p.GetArgs(arity)
	n, err := fmt.Print(args...)
	p.Ret(arity, n, err)
}

// QexecPrintf instr
func QexecPrintf(arity int, p *qlang.Context) {
	args := p.GetArgs(arity)
	n, err := fmt.Printf(args[0].(string), args[1:]...)
	p.Ret(arity, n, err)
}

// QexecPrintln instr
func QexecPrintln(arity int, p *qlang.Context) {
	args := p.GetArgs(arity)
	n, err := fmt.Println(args...)
	p.Ret(arity, n, err)
}

// QexecFprintln instr
func QexecFprintln(arity int, p *qlang.Context) {
	args := p.GetArgs(arity)
	n, err := fmt.Fprintln(args[0].(io.Writer), args[1:]...)
	p.Ret(arity, n, err)
}

// -----------------------------------------------------------------------------

func execPanic(zero int, p *qlang.Context) {
	panic(p.Pop())
}

// -----------------------------------------------------------------------------

// I is a Go package instance.
var I = qlang.NewGoPackage("")

func init() {
	I.RegisterFuncs(
		I.Func("panic", qlPanic, execPanic),
	)
	I.RegisterFuncvs(
		I.Funcv("print", fmt.Print, QexecPrint),
		I.Funcv("printf", fmt.Printf, QexecPrintf),
		I.Funcv("println", fmt.Println, QexecPrintln),
		I.Funcv("fprintln", fmt.Fprintln, QexecFprintln),
	)
	I.RegisterConsts(
		I.Const("true", reflect.Bool, true),
		I.Const("false", reflect.Bool, false),
		I.Const("nil", exec.ConstUnboundPtr, nil),
	)
	I.RegisterTypes(
		I.Type("bool", exec.TyBool),
		I.Type("int", exec.TyInt),
		I.Type("int8", exec.TyInt8),
		I.Type("int16", exec.TyInt16),
		I.Type("int32", exec.TyInt32),
		I.Type("int64", exec.TyInt64),
		I.Type("uint", exec.TyUint),
		I.Type("uint8", exec.TyUint8),
		I.Type("uint16", exec.TyUint16),
		I.Type("uint32", exec.TyUint32),
		I.Type("uint64", exec.TyUint64),
		I.Type("uintptr", exec.TyUintptr),
		I.Type("float32", exec.TyFloat32),
		I.Type("float64", exec.TyFloat64),
		I.Type("complex64", exec.TyComplex64),
		I.Type("complex128", exec.TyComplex128),
		I.Type("string", exec.TyString),
		I.Type("error", exec.TyError),
		I.Type("byte", exec.TyByte),
		I.Type("rune", exec.TyRune),
	)
}

// -----------------------------------------------------------------------------
