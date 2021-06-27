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

// Package builtin provide Go+ builtin stuffs, including builtin constants,
// types and functions.
package builtin

import (
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/goplus/gop"
	"github.com/goplus/gop/ast/gopiter"

	qspec "github.com/goplus/gop/exec.spec"
	exec "github.com/goplus/gop/exec/bytecode"
)

// -----------------------------------------------------------------------------

// QexecPrint instr
func QexecPrint(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	n, err := fmt.Print(args...)
	p.Ret(arity, n, err)
}

// QexecPrintf instr
func QexecPrintf(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	n, err := fmt.Printf(args[0].(string), args[1:]...)
	p.Ret(arity, n, err)
}

// QexecErrorf instr
func QexecErrorf(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	err := fmt.Errorf(args[0].(string), args[1:]...)
	p.Ret(arity, err)
}

// QexecPrintln instr
func QexecPrintln(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	n, err := fmt.Println(args...)
	p.Ret(arity, n, err)
}

// QexecFprintln instr
func QexecFprintln(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	n, err := fmt.Fprintln(args[0].(io.Writer), args[1:]...)
	p.Ret(arity, n, err)
}

// QexecIs instr
func QexecIs(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	is := errors.Is(gop.ToError(args[0]), gop.ToError(args[1]))
	p.Ret(2, is)
}

// QNewIter instr
func QNewIter(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	p.Ret(1, gopiter.NewIter(args[0]))
}

// QNext instr
func QNext(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	p.Ret(1, gopiter.Next(args[0].(gopiter.Iterator)))
}

// QKey instr
func QKey(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	gopiter.Key(args[0].(gopiter.Iterator), args[1])
	p.PopN(2)
}

// QValue instr
func QValue(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	gopiter.Value(args[0].(gopiter.Iterator), args[1])
	p.PopN(2)
}

func _gop_Select(v ...interface{}) (int, interface{}) {
	var cases []reflect.SelectCase
	for i := 0; i < len(v); i += 3 {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectDir(v[i].(int)),
			Chan: reflect.ValueOf(v[i+1]),
			Send: reflect.ValueOf(v[i+2]),
		})
	}
	chosen, recv, recvOK := reflect.Select(cases)
	if recvOK {
		return chosen, recv.Interface()
	}
	return chosen, nil
}

func QexecSelect(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	chosen, recv := _gop_Select(args...)
	p.Ret(2, chosen, recv)
}

// FuncGoInfo returns Go package and function name of a Go+ builtin function.
func FuncGoInfo(f string) ([2]string, bool) {
	fi, ok := builtinFnvs[f]
	return fi, ok
}

// -----------------------------------------------------------------------------

func execPanic(_ int, p *gop.Context) {
	panic(p.Pop())
}

// -----------------------------------------------------------------------------

// I is a Go package instance.
var I = gop.NewGoPackage("")

var builtinFnvs = map[string][2]string{
	"errorf":       {"fmt", "Errorf"},
	"print":        {"fmt", "Print"},
	"printf":       {"fmt", "Printf"},
	"println":      {"fmt", "Println"},
	"fprintf":      {"fmt", "Fprintf"},
	"_gop_NewIter": {"github.com/goplus/gop/ast/gopiter", "NewIter"},
	"_gop_Next":    {"github.com/goplus/gop/ast/gopiter", "Next"},
	"_gop_Key":     {"github.com/goplus/gop/ast/gopiter", "Key"},
	"_gop_Value":   {"github.com/goplus/gop/ast/gopiter", "Value"},
}

func init() {
	I.RegisterFuncs(
		I.Func("panic", qlPanic, execPanic),
		I.Func("is", errors.Is, QexecIs),
		I.Func("_gop_NewIter", gopiter.NewIter, QNewIter),
		I.Func("_gop_Next", gopiter.Next, QNext),
		I.Func("_gop_Key", gopiter.Key, QKey),
		I.Func("_gop_Value", gopiter.Value, QValue),
	)
	I.RegisterFuncvs(
		I.Funcv("errorf", fmt.Errorf, QexecErrorf),
		I.Funcv("print", fmt.Print, QexecPrint),
		I.Funcv("printf", fmt.Printf, QexecPrintf),
		I.Funcv("println", fmt.Println, QexecPrintln),
		I.Funcv("fprintln", fmt.Fprintln, QexecFprintln),
		I.Funcv("$select", _gop_Select, QexecSelect),
	)
	I.RegisterConsts(
		I.Const("true", reflect.Bool, true),
		I.Const("false", reflect.Bool, false),
		I.Const("nil", exec.ConstUnboundPtr, nil),
	)
	I.RegisterTypes(
		I.Type("bool", qspec.TyBool),
		I.Type("int", qspec.TyInt),
		I.Type("int8", qspec.TyInt8),
		I.Type("int16", qspec.TyInt16),
		I.Type("int32", qspec.TyInt32),
		I.Type("int64", qspec.TyInt64),
		I.Type("uint", qspec.TyUint),
		I.Type("uint8", qspec.TyUint8),
		I.Type("uint16", qspec.TyUint16),
		I.Type("uint32", qspec.TyUint32),
		I.Type("uint64", qspec.TyUint64),
		I.Type("uintptr", qspec.TyUintptr),
		I.Type("float32", qspec.TyFloat32),
		I.Type("float64", qspec.TyFloat64),
		I.Type("complex64", qspec.TyComplex64),
		I.Type("complex128", qspec.TyComplex128),
		I.Type("string", qspec.TyString),
		I.Type("error", qspec.TyError),
		I.Type("byte", qspec.TyByte),
		I.Type("rune", qspec.TyRune),
	)
}

// -----------------------------------------------------------------------------
