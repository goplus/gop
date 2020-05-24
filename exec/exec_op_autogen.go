/*
 Copyright 2020 Qiniu Cloud (七牛云)

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

package exec

var builtinOps = [...]func(i Instr, p *Context){
	(int(Int) << bitsOperator) | int(OpAdd):           execAddInt,
	(int(Int8) << bitsOperator) | int(OpAdd):          execAddInt8,
	(int(Int16) << bitsOperator) | int(OpAdd):         execAddInt16,
	(int(Int32) << bitsOperator) | int(OpAdd):         execAddInt32,
	(int(Int64) << bitsOperator) | int(OpAdd):         execAddInt64,
	(int(Uint) << bitsOperator) | int(OpAdd):          execAddUint,
	(int(Uint8) << bitsOperator) | int(OpAdd):         execAddUint8,
	(int(Uint16) << bitsOperator) | int(OpAdd):        execAddUint16,
	(int(Uint32) << bitsOperator) | int(OpAdd):        execAddUint32,
	(int(Uint64) << bitsOperator) | int(OpAdd):        execAddUint64,
	(int(Uintptr) << bitsOperator) | int(OpAdd):       execAddUintptr,
	(int(Float32) << bitsOperator) | int(OpAdd):       execAddFloat32,
	(int(Float64) << bitsOperator) | int(OpAdd):       execAddFloat64,
	(int(Complex64) << bitsOperator) | int(OpAdd):     execAddComplex64,
	(int(Complex128) << bitsOperator) | int(OpAdd):    execAddComplex128,
	(int(String) << bitsOperator) | int(OpAdd):        execAddString,
	(int(Int) << bitsOperator) | int(OpSub):           execSubInt,
	(int(Int8) << bitsOperator) | int(OpSub):          execSubInt8,
	(int(Int16) << bitsOperator) | int(OpSub):         execSubInt16,
	(int(Int32) << bitsOperator) | int(OpSub):         execSubInt32,
	(int(Int64) << bitsOperator) | int(OpSub):         execSubInt64,
	(int(Uint) << bitsOperator) | int(OpSub):          execSubUint,
	(int(Uint8) << bitsOperator) | int(OpSub):         execSubUint8,
	(int(Uint16) << bitsOperator) | int(OpSub):        execSubUint16,
	(int(Uint32) << bitsOperator) | int(OpSub):        execSubUint32,
	(int(Uint64) << bitsOperator) | int(OpSub):        execSubUint64,
	(int(Uintptr) << bitsOperator) | int(OpSub):       execSubUintptr,
	(int(Float32) << bitsOperator) | int(OpSub):       execSubFloat32,
	(int(Float64) << bitsOperator) | int(OpSub):       execSubFloat64,
	(int(Complex64) << bitsOperator) | int(OpSub):     execSubComplex64,
	(int(Complex128) << bitsOperator) | int(OpSub):    execSubComplex128,
	(int(Int) << bitsOperator) | int(OpMul):           execMulInt,
	(int(Int8) << bitsOperator) | int(OpMul):          execMulInt8,
	(int(Int16) << bitsOperator) | int(OpMul):         execMulInt16,
	(int(Int32) << bitsOperator) | int(OpMul):         execMulInt32,
	(int(Int64) << bitsOperator) | int(OpMul):         execMulInt64,
	(int(Uint) << bitsOperator) | int(OpMul):          execMulUint,
	(int(Uint8) << bitsOperator) | int(OpMul):         execMulUint8,
	(int(Uint16) << bitsOperator) | int(OpMul):        execMulUint16,
	(int(Uint32) << bitsOperator) | int(OpMul):        execMulUint32,
	(int(Uint64) << bitsOperator) | int(OpMul):        execMulUint64,
	(int(Uintptr) << bitsOperator) | int(OpMul):       execMulUintptr,
	(int(Float32) << bitsOperator) | int(OpMul):       execMulFloat32,
	(int(Float64) << bitsOperator) | int(OpMul):       execMulFloat64,
	(int(Complex64) << bitsOperator) | int(OpMul):     execMulComplex64,
	(int(Complex128) << bitsOperator) | int(OpMul):    execMulComplex128,
	(int(Int) << bitsOperator) | int(OpDiv):           execDivInt,
	(int(Int8) << bitsOperator) | int(OpDiv):          execDivInt8,
	(int(Int16) << bitsOperator) | int(OpDiv):         execDivInt16,
	(int(Int32) << bitsOperator) | int(OpDiv):         execDivInt32,
	(int(Int64) << bitsOperator) | int(OpDiv):         execDivInt64,
	(int(Uint) << bitsOperator) | int(OpDiv):          execDivUint,
	(int(Uint8) << bitsOperator) | int(OpDiv):         execDivUint8,
	(int(Uint16) << bitsOperator) | int(OpDiv):        execDivUint16,
	(int(Uint32) << bitsOperator) | int(OpDiv):        execDivUint32,
	(int(Uint64) << bitsOperator) | int(OpDiv):        execDivUint64,
	(int(Uintptr) << bitsOperator) | int(OpDiv):       execDivUintptr,
	(int(Float32) << bitsOperator) | int(OpDiv):       execDivFloat32,
	(int(Float64) << bitsOperator) | int(OpDiv):       execDivFloat64,
	(int(Complex64) << bitsOperator) | int(OpDiv):     execDivComplex64,
	(int(Complex128) << bitsOperator) | int(OpDiv):    execDivComplex128,
	(int(Int) << bitsOperator) | int(OpMod):           execModInt,
	(int(Int8) << bitsOperator) | int(OpMod):          execModInt8,
	(int(Int16) << bitsOperator) | int(OpMod):         execModInt16,
	(int(Int32) << bitsOperator) | int(OpMod):         execModInt32,
	(int(Int64) << bitsOperator) | int(OpMod):         execModInt64,
	(int(Uint) << bitsOperator) | int(OpMod):          execModUint,
	(int(Uint8) << bitsOperator) | int(OpMod):         execModUint8,
	(int(Uint16) << bitsOperator) | int(OpMod):        execModUint16,
	(int(Uint32) << bitsOperator) | int(OpMod):        execModUint32,
	(int(Uint64) << bitsOperator) | int(OpMod):        execModUint64,
	(int(Uintptr) << bitsOperator) | int(OpMod):       execModUintptr,
	(int(Int) << bitsOperator) | int(OpBitAnd):        execBitAndInt,
	(int(Int8) << bitsOperator) | int(OpBitAnd):       execBitAndInt8,
	(int(Int16) << bitsOperator) | int(OpBitAnd):      execBitAndInt16,
	(int(Int32) << bitsOperator) | int(OpBitAnd):      execBitAndInt32,
	(int(Int64) << bitsOperator) | int(OpBitAnd):      execBitAndInt64,
	(int(Uint) << bitsOperator) | int(OpBitAnd):       execBitAndUint,
	(int(Uint8) << bitsOperator) | int(OpBitAnd):      execBitAndUint8,
	(int(Uint16) << bitsOperator) | int(OpBitAnd):     execBitAndUint16,
	(int(Uint32) << bitsOperator) | int(OpBitAnd):     execBitAndUint32,
	(int(Uint64) << bitsOperator) | int(OpBitAnd):     execBitAndUint64,
	(int(Uintptr) << bitsOperator) | int(OpBitAnd):    execBitAndUintptr,
	(int(Int) << bitsOperator) | int(OpBitOr):         execBitOrInt,
	(int(Int8) << bitsOperator) | int(OpBitOr):        execBitOrInt8,
	(int(Int16) << bitsOperator) | int(OpBitOr):       execBitOrInt16,
	(int(Int32) << bitsOperator) | int(OpBitOr):       execBitOrInt32,
	(int(Int64) << bitsOperator) | int(OpBitOr):       execBitOrInt64,
	(int(Uint) << bitsOperator) | int(OpBitOr):        execBitOrUint,
	(int(Uint8) << bitsOperator) | int(OpBitOr):       execBitOrUint8,
	(int(Uint16) << bitsOperator) | int(OpBitOr):      execBitOrUint16,
	(int(Uint32) << bitsOperator) | int(OpBitOr):      execBitOrUint32,
	(int(Uint64) << bitsOperator) | int(OpBitOr):      execBitOrUint64,
	(int(Uintptr) << bitsOperator) | int(OpBitOr):     execBitOrUintptr,
	(int(Int) << bitsOperator) | int(OpBitXor):        execBitXorInt,
	(int(Int8) << bitsOperator) | int(OpBitXor):       execBitXorInt8,
	(int(Int16) << bitsOperator) | int(OpBitXor):      execBitXorInt16,
	(int(Int32) << bitsOperator) | int(OpBitXor):      execBitXorInt32,
	(int(Int64) << bitsOperator) | int(OpBitXor):      execBitXorInt64,
	(int(Uint) << bitsOperator) | int(OpBitXor):       execBitXorUint,
	(int(Uint8) << bitsOperator) | int(OpBitXor):      execBitXorUint8,
	(int(Uint16) << bitsOperator) | int(OpBitXor):     execBitXorUint16,
	(int(Uint32) << bitsOperator) | int(OpBitXor):     execBitXorUint32,
	(int(Uint64) << bitsOperator) | int(OpBitXor):     execBitXorUint64,
	(int(Uintptr) << bitsOperator) | int(OpBitXor):    execBitXorUintptr,
	(int(Int) << bitsOperator) | int(OpBitAndNot):     execBitAndNotInt,
	(int(Int8) << bitsOperator) | int(OpBitAndNot):    execBitAndNotInt8,
	(int(Int16) << bitsOperator) | int(OpBitAndNot):   execBitAndNotInt16,
	(int(Int32) << bitsOperator) | int(OpBitAndNot):   execBitAndNotInt32,
	(int(Int64) << bitsOperator) | int(OpBitAndNot):   execBitAndNotInt64,
	(int(Uint) << bitsOperator) | int(OpBitAndNot):    execBitAndNotUint,
	(int(Uint8) << bitsOperator) | int(OpBitAndNot):   execBitAndNotUint8,
	(int(Uint16) << bitsOperator) | int(OpBitAndNot):  execBitAndNotUint16,
	(int(Uint32) << bitsOperator) | int(OpBitAndNot):  execBitAndNotUint32,
	(int(Uint64) << bitsOperator) | int(OpBitAndNot):  execBitAndNotUint64,
	(int(Uintptr) << bitsOperator) | int(OpBitAndNot): execBitAndNotUintptr,
	(int(Int) << bitsOperator) | int(OpBitSHL):        execBitSHLInt,
	(int(Int8) << bitsOperator) | int(OpBitSHL):       execBitSHLInt8,
	(int(Int16) << bitsOperator) | int(OpBitSHL):      execBitSHLInt16,
	(int(Int32) << bitsOperator) | int(OpBitSHL):      execBitSHLInt32,
	(int(Int64) << bitsOperator) | int(OpBitSHL):      execBitSHLInt64,
	(int(Uint) << bitsOperator) | int(OpBitSHL):       execBitSHLUint,
	(int(Uint8) << bitsOperator) | int(OpBitSHL):      execBitSHLUint8,
	(int(Uint16) << bitsOperator) | int(OpBitSHL):     execBitSHLUint16,
	(int(Uint32) << bitsOperator) | int(OpBitSHL):     execBitSHLUint32,
	(int(Uint64) << bitsOperator) | int(OpBitSHL):     execBitSHLUint64,
	(int(Uintptr) << bitsOperator) | int(OpBitSHL):    execBitSHLUintptr,
	(int(Int) << bitsOperator) | int(OpBitSHR):        execBitSHRInt,
	(int(Int8) << bitsOperator) | int(OpBitSHR):       execBitSHRInt8,
	(int(Int16) << bitsOperator) | int(OpBitSHR):      execBitSHRInt16,
	(int(Int32) << bitsOperator) | int(OpBitSHR):      execBitSHRInt32,
	(int(Int64) << bitsOperator) | int(OpBitSHR):      execBitSHRInt64,
	(int(Uint) << bitsOperator) | int(OpBitSHR):       execBitSHRUint,
	(int(Uint8) << bitsOperator) | int(OpBitSHR):      execBitSHRUint8,
	(int(Uint16) << bitsOperator) | int(OpBitSHR):     execBitSHRUint16,
	(int(Uint32) << bitsOperator) | int(OpBitSHR):     execBitSHRUint32,
	(int(Uint64) << bitsOperator) | int(OpBitSHR):     execBitSHRUint64,
	(int(Uintptr) << bitsOperator) | int(OpBitSHR):    execBitSHRUintptr,
	(int(Int) << bitsOperator) | int(OpLT):            execLTInt,
	(int(Int8) << bitsOperator) | int(OpLT):           execLTInt8,
	(int(Int16) << bitsOperator) | int(OpLT):          execLTInt16,
	(int(Int32) << bitsOperator) | int(OpLT):          execLTInt32,
	(int(Int64) << bitsOperator) | int(OpLT):          execLTInt64,
	(int(Uint) << bitsOperator) | int(OpLT):           execLTUint,
	(int(Uint8) << bitsOperator) | int(OpLT):          execLTUint8,
	(int(Uint16) << bitsOperator) | int(OpLT):         execLTUint16,
	(int(Uint32) << bitsOperator) | int(OpLT):         execLTUint32,
	(int(Uint64) << bitsOperator) | int(OpLT):         execLTUint64,
	(int(Uintptr) << bitsOperator) | int(OpLT):        execLTUintptr,
	(int(Float32) << bitsOperator) | int(OpLT):        execLTFloat32,
	(int(Float64) << bitsOperator) | int(OpLT):        execLTFloat64,
	(int(String) << bitsOperator) | int(OpLT):         execLTString,
	(int(Int) << bitsOperator) | int(OpLE):            execLEInt,
	(int(Int8) << bitsOperator) | int(OpLE):           execLEInt8,
	(int(Int16) << bitsOperator) | int(OpLE):          execLEInt16,
	(int(Int32) << bitsOperator) | int(OpLE):          execLEInt32,
	(int(Int64) << bitsOperator) | int(OpLE):          execLEInt64,
	(int(Uint) << bitsOperator) | int(OpLE):           execLEUint,
	(int(Uint8) << bitsOperator) | int(OpLE):          execLEUint8,
	(int(Uint16) << bitsOperator) | int(OpLE):         execLEUint16,
	(int(Uint32) << bitsOperator) | int(OpLE):         execLEUint32,
	(int(Uint64) << bitsOperator) | int(OpLE):         execLEUint64,
	(int(Uintptr) << bitsOperator) | int(OpLE):        execLEUintptr,
	(int(Float32) << bitsOperator) | int(OpLE):        execLEFloat32,
	(int(Float64) << bitsOperator) | int(OpLE):        execLEFloat64,
	(int(String) << bitsOperator) | int(OpLE):         execLEString,
	(int(Int) << bitsOperator) | int(OpGT):            execGTInt,
	(int(Int8) << bitsOperator) | int(OpGT):           execGTInt8,
	(int(Int16) << bitsOperator) | int(OpGT):          execGTInt16,
	(int(Int32) << bitsOperator) | int(OpGT):          execGTInt32,
	(int(Int64) << bitsOperator) | int(OpGT):          execGTInt64,
	(int(Uint) << bitsOperator) | int(OpGT):           execGTUint,
	(int(Uint8) << bitsOperator) | int(OpGT):          execGTUint8,
	(int(Uint16) << bitsOperator) | int(OpGT):         execGTUint16,
	(int(Uint32) << bitsOperator) | int(OpGT):         execGTUint32,
	(int(Uint64) << bitsOperator) | int(OpGT):         execGTUint64,
	(int(Uintptr) << bitsOperator) | int(OpGT):        execGTUintptr,
	(int(Float32) << bitsOperator) | int(OpGT):        execGTFloat32,
	(int(Float64) << bitsOperator) | int(OpGT):        execGTFloat64,
	(int(String) << bitsOperator) | int(OpGT):         execGTString,
	(int(Int) << bitsOperator) | int(OpGE):            execGEInt,
	(int(Int8) << bitsOperator) | int(OpGE):           execGEInt8,
	(int(Int16) << bitsOperator) | int(OpGE):          execGEInt16,
	(int(Int32) << bitsOperator) | int(OpGE):          execGEInt32,
	(int(Int64) << bitsOperator) | int(OpGE):          execGEInt64,
	(int(Uint) << bitsOperator) | int(OpGE):           execGEUint,
	(int(Uint8) << bitsOperator) | int(OpGE):          execGEUint8,
	(int(Uint16) << bitsOperator) | int(OpGE):         execGEUint16,
	(int(Uint32) << bitsOperator) | int(OpGE):         execGEUint32,
	(int(Uint64) << bitsOperator) | int(OpGE):         execGEUint64,
	(int(Uintptr) << bitsOperator) | int(OpGE):        execGEUintptr,
	(int(Float32) << bitsOperator) | int(OpGE):        execGEFloat32,
	(int(Float64) << bitsOperator) | int(OpGE):        execGEFloat64,
	(int(String) << bitsOperator) | int(OpGE):         execGEString,
	(int(Int) << bitsOperator) | int(OpEQ):            execEQInt,
	(int(Int8) << bitsOperator) | int(OpEQ):           execEQInt8,
	(int(Int16) << bitsOperator) | int(OpEQ):          execEQInt16,
	(int(Int32) << bitsOperator) | int(OpEQ):          execEQInt32,
	(int(Int64) << bitsOperator) | int(OpEQ):          execEQInt64,
	(int(Uint) << bitsOperator) | int(OpEQ):           execEQUint,
	(int(Uint8) << bitsOperator) | int(OpEQ):          execEQUint8,
	(int(Uint16) << bitsOperator) | int(OpEQ):         execEQUint16,
	(int(Uint32) << bitsOperator) | int(OpEQ):         execEQUint32,
	(int(Uint64) << bitsOperator) | int(OpEQ):         execEQUint64,
	(int(Uintptr) << bitsOperator) | int(OpEQ):        execEQUintptr,
	(int(Float32) << bitsOperator) | int(OpEQ):        execEQFloat32,
	(int(Float64) << bitsOperator) | int(OpEQ):        execEQFloat64,
	(int(Complex64) << bitsOperator) | int(OpEQ):      execEQComplex64,
	(int(Complex128) << bitsOperator) | int(OpEQ):     execEQComplex128,
	(int(String) << bitsOperator) | int(OpEQ):         execEQString,
	(int(Int) << bitsOperator) | int(OpNE):            execNEInt,
	(int(Int8) << bitsOperator) | int(OpNE):           execNEInt8,
	(int(Int16) << bitsOperator) | int(OpNE):          execNEInt16,
	(int(Int32) << bitsOperator) | int(OpNE):          execNEInt32,
	(int(Int64) << bitsOperator) | int(OpNE):          execNEInt64,
	(int(Uint) << bitsOperator) | int(OpNE):           execNEUint,
	(int(Uint8) << bitsOperator) | int(OpNE):          execNEUint8,
	(int(Uint16) << bitsOperator) | int(OpNE):         execNEUint16,
	(int(Uint32) << bitsOperator) | int(OpNE):         execNEUint32,
	(int(Uint64) << bitsOperator) | int(OpNE):         execNEUint64,
	(int(Uintptr) << bitsOperator) | int(OpNE):        execNEUintptr,
	(int(Float32) << bitsOperator) | int(OpNE):        execNEFloat32,
	(int(Float64) << bitsOperator) | int(OpNE):        execNEFloat64,
	(int(Complex64) << bitsOperator) | int(OpNE):      execNEComplex64,
	(int(Complex128) << bitsOperator) | int(OpNE):     execNEComplex128,
	(int(String) << bitsOperator) | int(OpNE):         execNEString,
	(int(Bool) << bitsOperator) | int(OpLAnd):         execLAndBool,
	(int(Bool) << bitsOperator) | int(OpLOr):          execLOrBool,
	(int(Int) << bitsOperator) | int(OpNeg):           execNegInt,
	(int(Int8) << bitsOperator) | int(OpNeg):          execNegInt8,
	(int(Int16) << bitsOperator) | int(OpNeg):         execNegInt16,
	(int(Int32) << bitsOperator) | int(OpNeg):         execNegInt32,
	(int(Int64) << bitsOperator) | int(OpNeg):         execNegInt64,
	(int(Uint) << bitsOperator) | int(OpNeg):          execNegUint,
	(int(Uint8) << bitsOperator) | int(OpNeg):         execNegUint8,
	(int(Uint16) << bitsOperator) | int(OpNeg):        execNegUint16,
	(int(Uint32) << bitsOperator) | int(OpNeg):        execNegUint32,
	(int(Uint64) << bitsOperator) | int(OpNeg):        execNegUint64,
	(int(Uintptr) << bitsOperator) | int(OpNeg):       execNegUintptr,
	(int(Float32) << bitsOperator) | int(OpNeg):       execNegFloat32,
	(int(Float64) << bitsOperator) | int(OpNeg):       execNegFloat64,
	(int(Complex64) << bitsOperator) | int(OpNeg):     execNegComplex64,
	(int(Complex128) << bitsOperator) | int(OpNeg):    execNegComplex128,
	(int(Bool) << bitsOperator) | int(OpNot):          execNotBool,
	(int(Int) << bitsOperator) | int(OpBitNot):        execBitNotInt,
	(int(Int8) << bitsOperator) | int(OpBitNot):       execBitNotInt8,
	(int(Int16) << bitsOperator) | int(OpBitNot):      execBitNotInt16,
	(int(Int32) << bitsOperator) | int(OpBitNot):      execBitNotInt32,
	(int(Int64) << bitsOperator) | int(OpBitNot):      execBitNotInt64,
	(int(Uint) << bitsOperator) | int(OpBitNot):       execBitNotUint,
	(int(Uint8) << bitsOperator) | int(OpBitNot):      execBitNotUint8,
	(int(Uint16) << bitsOperator) | int(OpBitNot):     execBitNotUint16,
	(int(Uint32) << bitsOperator) | int(OpBitNot):     execBitNotUint32,
	(int(Uint64) << bitsOperator) | int(OpBitNot):     execBitNotUint64,
	(int(Uintptr) << bitsOperator) | int(OpBitNot):    execBitNotUintptr,
}

func execAddInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) + p.data[n-1].(int)
	p.data = p.data[:n-1]
}

func execAddInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) + p.data[n-1].(int8)
	p.data = p.data[:n-1]
}

func execAddInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) + p.data[n-1].(int16)
	p.data = p.data[:n-1]
}

func execAddInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) + p.data[n-1].(int32)
	p.data = p.data[:n-1]
}

func execAddInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) + p.data[n-1].(int64)
	p.data = p.data[:n-1]
}

func execAddUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) + p.data[n-1].(uint)
	p.data = p.data[:n-1]
}

func execAddUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) + p.data[n-1].(uint8)
	p.data = p.data[:n-1]
}

func execAddUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) + p.data[n-1].(uint16)
	p.data = p.data[:n-1]
}

func execAddUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) + p.data[n-1].(uint32)
	p.data = p.data[:n-1]
}

func execAddUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) + p.data[n-1].(uint64)
	p.data = p.data[:n-1]
}

func execAddUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) + p.data[n-1].(uintptr)
	p.data = p.data[:n-1]
}

func execAddFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float32) + p.data[n-1].(float32)
	p.data = p.data[:n-1]
}

func execAddFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float64) + p.data[n-1].(float64)
	p.data = p.data[:n-1]
}

func execAddComplex64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(complex64) + p.data[n-1].(complex64)
	p.data = p.data[:n-1]
}

func execAddComplex128(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(complex128) + p.data[n-1].(complex128)
	p.data = p.data[:n-1]
}

func execAddString(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(string) + p.data[n-1].(string)
	p.data = p.data[:n-1]
}

func execSubInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) - p.data[n-1].(int)
	p.data = p.data[:n-1]
}

func execSubInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) - p.data[n-1].(int8)
	p.data = p.data[:n-1]
}

func execSubInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) - p.data[n-1].(int16)
	p.data = p.data[:n-1]
}

func execSubInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) - p.data[n-1].(int32)
	p.data = p.data[:n-1]
}

func execSubInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) - p.data[n-1].(int64)
	p.data = p.data[:n-1]
}

func execSubUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) - p.data[n-1].(uint)
	p.data = p.data[:n-1]
}

func execSubUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) - p.data[n-1].(uint8)
	p.data = p.data[:n-1]
}

func execSubUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) - p.data[n-1].(uint16)
	p.data = p.data[:n-1]
}

func execSubUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) - p.data[n-1].(uint32)
	p.data = p.data[:n-1]
}

func execSubUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) - p.data[n-1].(uint64)
	p.data = p.data[:n-1]
}

func execSubUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) - p.data[n-1].(uintptr)
	p.data = p.data[:n-1]
}

func execSubFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float32) - p.data[n-1].(float32)
	p.data = p.data[:n-1]
}

func execSubFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float64) - p.data[n-1].(float64)
	p.data = p.data[:n-1]
}

func execSubComplex64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(complex64) - p.data[n-1].(complex64)
	p.data = p.data[:n-1]
}

func execSubComplex128(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(complex128) - p.data[n-1].(complex128)
	p.data = p.data[:n-1]
}

func execMulInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) * p.data[n-1].(int)
	p.data = p.data[:n-1]
}

func execMulInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) * p.data[n-1].(int8)
	p.data = p.data[:n-1]
}

func execMulInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) * p.data[n-1].(int16)
	p.data = p.data[:n-1]
}

func execMulInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) * p.data[n-1].(int32)
	p.data = p.data[:n-1]
}

func execMulInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) * p.data[n-1].(int64)
	p.data = p.data[:n-1]
}

func execMulUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) * p.data[n-1].(uint)
	p.data = p.data[:n-1]
}

func execMulUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) * p.data[n-1].(uint8)
	p.data = p.data[:n-1]
}

func execMulUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) * p.data[n-1].(uint16)
	p.data = p.data[:n-1]
}

func execMulUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) * p.data[n-1].(uint32)
	p.data = p.data[:n-1]
}

func execMulUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) * p.data[n-1].(uint64)
	p.data = p.data[:n-1]
}

func execMulUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) * p.data[n-1].(uintptr)
	p.data = p.data[:n-1]
}

func execMulFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float32) * p.data[n-1].(float32)
	p.data = p.data[:n-1]
}

func execMulFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float64) * p.data[n-1].(float64)
	p.data = p.data[:n-1]
}

func execMulComplex64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(complex64) * p.data[n-1].(complex64)
	p.data = p.data[:n-1]
}

func execMulComplex128(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(complex128) * p.data[n-1].(complex128)
	p.data = p.data[:n-1]
}

func execDivInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) / p.data[n-1].(int)
	p.data = p.data[:n-1]
}

func execDivInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) / p.data[n-1].(int8)
	p.data = p.data[:n-1]
}

func execDivInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) / p.data[n-1].(int16)
	p.data = p.data[:n-1]
}

func execDivInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) / p.data[n-1].(int32)
	p.data = p.data[:n-1]
}

func execDivInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) / p.data[n-1].(int64)
	p.data = p.data[:n-1]
}

func execDivUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) / p.data[n-1].(uint)
	p.data = p.data[:n-1]
}

func execDivUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) / p.data[n-1].(uint8)
	p.data = p.data[:n-1]
}

func execDivUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) / p.data[n-1].(uint16)
	p.data = p.data[:n-1]
}

func execDivUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) / p.data[n-1].(uint32)
	p.data = p.data[:n-1]
}

func execDivUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) / p.data[n-1].(uint64)
	p.data = p.data[:n-1]
}

func execDivUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) / p.data[n-1].(uintptr)
	p.data = p.data[:n-1]
}

func execDivFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float32) / p.data[n-1].(float32)
	p.data = p.data[:n-1]
}

func execDivFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float64) / p.data[n-1].(float64)
	p.data = p.data[:n-1]
}

func execDivComplex64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(complex64) / p.data[n-1].(complex64)
	p.data = p.data[:n-1]
}

func execDivComplex128(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(complex128) / p.data[n-1].(complex128)
	p.data = p.data[:n-1]
}

func execModInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) % p.data[n-1].(int)
	p.data = p.data[:n-1]
}

func execModInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) % p.data[n-1].(int8)
	p.data = p.data[:n-1]
}

func execModInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) % p.data[n-1].(int16)
	p.data = p.data[:n-1]
}

func execModInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) % p.data[n-1].(int32)
	p.data = p.data[:n-1]
}

func execModInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) % p.data[n-1].(int64)
	p.data = p.data[:n-1]
}

func execModUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) % p.data[n-1].(uint)
	p.data = p.data[:n-1]
}

func execModUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) % p.data[n-1].(uint8)
	p.data = p.data[:n-1]
}

func execModUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) % p.data[n-1].(uint16)
	p.data = p.data[:n-1]
}

func execModUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) % p.data[n-1].(uint32)
	p.data = p.data[:n-1]
}

func execModUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) % p.data[n-1].(uint64)
	p.data = p.data[:n-1]
}

func execModUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) % p.data[n-1].(uintptr)
	p.data = p.data[:n-1]
}

func execBitAndInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) & p.data[n-1].(int)
	p.data = p.data[:n-1]
}

func execBitAndInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) & p.data[n-1].(int8)
	p.data = p.data[:n-1]
}

func execBitAndInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) & p.data[n-1].(int16)
	p.data = p.data[:n-1]
}

func execBitAndInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) & p.data[n-1].(int32)
	p.data = p.data[:n-1]
}

func execBitAndInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) & p.data[n-1].(int64)
	p.data = p.data[:n-1]
}

func execBitAndUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) & p.data[n-1].(uint)
	p.data = p.data[:n-1]
}

func execBitAndUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) & p.data[n-1].(uint8)
	p.data = p.data[:n-1]
}

func execBitAndUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) & p.data[n-1].(uint16)
	p.data = p.data[:n-1]
}

func execBitAndUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) & p.data[n-1].(uint32)
	p.data = p.data[:n-1]
}

func execBitAndUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) & p.data[n-1].(uint64)
	p.data = p.data[:n-1]
}

func execBitAndUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) & p.data[n-1].(uintptr)
	p.data = p.data[:n-1]
}

func execBitOrInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) | p.data[n-1].(int)
	p.data = p.data[:n-1]
}

func execBitOrInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) | p.data[n-1].(int8)
	p.data = p.data[:n-1]
}

func execBitOrInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) | p.data[n-1].(int16)
	p.data = p.data[:n-1]
}

func execBitOrInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) | p.data[n-1].(int32)
	p.data = p.data[:n-1]
}

func execBitOrInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) | p.data[n-1].(int64)
	p.data = p.data[:n-1]
}

func execBitOrUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) | p.data[n-1].(uint)
	p.data = p.data[:n-1]
}

func execBitOrUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) | p.data[n-1].(uint8)
	p.data = p.data[:n-1]
}

func execBitOrUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) | p.data[n-1].(uint16)
	p.data = p.data[:n-1]
}

func execBitOrUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) | p.data[n-1].(uint32)
	p.data = p.data[:n-1]
}

func execBitOrUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) | p.data[n-1].(uint64)
	p.data = p.data[:n-1]
}

func execBitOrUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) | p.data[n-1].(uintptr)
	p.data = p.data[:n-1]
}

func execBitXorInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) ^ p.data[n-1].(int)
	p.data = p.data[:n-1]
}

func execBitXorInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) ^ p.data[n-1].(int8)
	p.data = p.data[:n-1]
}

func execBitXorInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) ^ p.data[n-1].(int16)
	p.data = p.data[:n-1]
}

func execBitXorInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) ^ p.data[n-1].(int32)
	p.data = p.data[:n-1]
}

func execBitXorInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) ^ p.data[n-1].(int64)
	p.data = p.data[:n-1]
}

func execBitXorUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) ^ p.data[n-1].(uint)
	p.data = p.data[:n-1]
}

func execBitXorUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) ^ p.data[n-1].(uint8)
	p.data = p.data[:n-1]
}

func execBitXorUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) ^ p.data[n-1].(uint16)
	p.data = p.data[:n-1]
}

func execBitXorUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) ^ p.data[n-1].(uint32)
	p.data = p.data[:n-1]
}

func execBitXorUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) ^ p.data[n-1].(uint64)
	p.data = p.data[:n-1]
}

func execBitXorUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) ^ p.data[n-1].(uintptr)
	p.data = p.data[:n-1]
}

func execBitAndNotInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) &^ p.data[n-1].(int)
	p.data = p.data[:n-1]
}

func execBitAndNotInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) &^ p.data[n-1].(int8)
	p.data = p.data[:n-1]
}

func execBitAndNotInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) &^ p.data[n-1].(int16)
	p.data = p.data[:n-1]
}

func execBitAndNotInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) &^ p.data[n-1].(int32)
	p.data = p.data[:n-1]
}

func execBitAndNotInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) &^ p.data[n-1].(int64)
	p.data = p.data[:n-1]
}

func execBitAndNotUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) &^ p.data[n-1].(uint)
	p.data = p.data[:n-1]
}

func execBitAndNotUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) &^ p.data[n-1].(uint8)
	p.data = p.data[:n-1]
}

func execBitAndNotUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) &^ p.data[n-1].(uint16)
	p.data = p.data[:n-1]
}

func execBitAndNotUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) &^ p.data[n-1].(uint32)
	p.data = p.data[:n-1]
}

func execBitAndNotUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) &^ p.data[n-1].(uint64)
	p.data = p.data[:n-1]
}

func execBitAndNotUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) &^ p.data[n-1].(uintptr)
	p.data = p.data[:n-1]
}

func execBitSHLInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) << toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execBitSHLInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) << toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execBitSHLInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) << toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execBitSHLInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) << toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execBitSHLInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) << toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execBitSHLUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) << toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execBitSHLUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) << toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execBitSHLUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) << toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execBitSHLUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) << toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execBitSHLUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) << toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execBitSHLUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) << toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execBitSHRInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) >> toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execBitSHRInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) >> toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execBitSHRInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) >> toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execBitSHRInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) >> toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execBitSHRInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) >> toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execBitSHRUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) >> toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execBitSHRUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) >> toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execBitSHRUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) >> toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execBitSHRUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) >> toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execBitSHRUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) >> toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execBitSHRUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) >> toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execLTInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) < p.data[n-1].(int)
	p.data = p.data[:n-1]
}

func execLTInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) < p.data[n-1].(int8)
	p.data = p.data[:n-1]
}

func execLTInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) < p.data[n-1].(int16)
	p.data = p.data[:n-1]
}

func execLTInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) < p.data[n-1].(int32)
	p.data = p.data[:n-1]
}

func execLTInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) < p.data[n-1].(int64)
	p.data = p.data[:n-1]
}

func execLTUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) < p.data[n-1].(uint)
	p.data = p.data[:n-1]
}

func execLTUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) < p.data[n-1].(uint8)
	p.data = p.data[:n-1]
}

func execLTUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) < p.data[n-1].(uint16)
	p.data = p.data[:n-1]
}

func execLTUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) < p.data[n-1].(uint32)
	p.data = p.data[:n-1]
}

func execLTUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) < p.data[n-1].(uint64)
	p.data = p.data[:n-1]
}

func execLTUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) < p.data[n-1].(uintptr)
	p.data = p.data[:n-1]
}

func execLTFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float32) < p.data[n-1].(float32)
	p.data = p.data[:n-1]
}

func execLTFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float64) < p.data[n-1].(float64)
	p.data = p.data[:n-1]
}

func execLTString(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(string) < p.data[n-1].(string)
	p.data = p.data[:n-1]
}

func execLEInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) <= p.data[n-1].(int)
	p.data = p.data[:n-1]
}

func execLEInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) <= p.data[n-1].(int8)
	p.data = p.data[:n-1]
}

func execLEInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) <= p.data[n-1].(int16)
	p.data = p.data[:n-1]
}

func execLEInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) <= p.data[n-1].(int32)
	p.data = p.data[:n-1]
}

func execLEInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) <= p.data[n-1].(int64)
	p.data = p.data[:n-1]
}

func execLEUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) <= p.data[n-1].(uint)
	p.data = p.data[:n-1]
}

func execLEUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) <= p.data[n-1].(uint8)
	p.data = p.data[:n-1]
}

func execLEUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) <= p.data[n-1].(uint16)
	p.data = p.data[:n-1]
}

func execLEUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) <= p.data[n-1].(uint32)
	p.data = p.data[:n-1]
}

func execLEUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) <= p.data[n-1].(uint64)
	p.data = p.data[:n-1]
}

func execLEUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) <= p.data[n-1].(uintptr)
	p.data = p.data[:n-1]
}

func execLEFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float32) <= p.data[n-1].(float32)
	p.data = p.data[:n-1]
}

func execLEFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float64) <= p.data[n-1].(float64)
	p.data = p.data[:n-1]
}

func execLEString(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(string) <= p.data[n-1].(string)
	p.data = p.data[:n-1]
}

func execGTInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) > p.data[n-1].(int)
	p.data = p.data[:n-1]
}

func execGTInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) > p.data[n-1].(int8)
	p.data = p.data[:n-1]
}

func execGTInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) > p.data[n-1].(int16)
	p.data = p.data[:n-1]
}

func execGTInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) > p.data[n-1].(int32)
	p.data = p.data[:n-1]
}

func execGTInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) > p.data[n-1].(int64)
	p.data = p.data[:n-1]
}

func execGTUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) > p.data[n-1].(uint)
	p.data = p.data[:n-1]
}

func execGTUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) > p.data[n-1].(uint8)
	p.data = p.data[:n-1]
}

func execGTUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) > p.data[n-1].(uint16)
	p.data = p.data[:n-1]
}

func execGTUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) > p.data[n-1].(uint32)
	p.data = p.data[:n-1]
}

func execGTUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) > p.data[n-1].(uint64)
	p.data = p.data[:n-1]
}

func execGTUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) > p.data[n-1].(uintptr)
	p.data = p.data[:n-1]
}

func execGTFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float32) > p.data[n-1].(float32)
	p.data = p.data[:n-1]
}

func execGTFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float64) > p.data[n-1].(float64)
	p.data = p.data[:n-1]
}

func execGTString(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(string) > p.data[n-1].(string)
	p.data = p.data[:n-1]
}

func execGEInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) >= p.data[n-1].(int)
	p.data = p.data[:n-1]
}

func execGEInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) >= p.data[n-1].(int8)
	p.data = p.data[:n-1]
}

func execGEInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) >= p.data[n-1].(int16)
	p.data = p.data[:n-1]
}

func execGEInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) >= p.data[n-1].(int32)
	p.data = p.data[:n-1]
}

func execGEInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) >= p.data[n-1].(int64)
	p.data = p.data[:n-1]
}

func execGEUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) >= p.data[n-1].(uint)
	p.data = p.data[:n-1]
}

func execGEUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) >= p.data[n-1].(uint8)
	p.data = p.data[:n-1]
}

func execGEUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) >= p.data[n-1].(uint16)
	p.data = p.data[:n-1]
}

func execGEUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) >= p.data[n-1].(uint32)
	p.data = p.data[:n-1]
}

func execGEUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) >= p.data[n-1].(uint64)
	p.data = p.data[:n-1]
}

func execGEUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) >= p.data[n-1].(uintptr)
	p.data = p.data[:n-1]
}

func execGEFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float32) >= p.data[n-1].(float32)
	p.data = p.data[:n-1]
}

func execGEFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float64) >= p.data[n-1].(float64)
	p.data = p.data[:n-1]
}

func execGEString(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(string) >= p.data[n-1].(string)
	p.data = p.data[:n-1]
}

func execEQInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) == p.data[n-1].(int)
	p.data = p.data[:n-1]
}

func execEQInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) == p.data[n-1].(int8)
	p.data = p.data[:n-1]
}

func execEQInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) == p.data[n-1].(int16)
	p.data = p.data[:n-1]
}

func execEQInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) == p.data[n-1].(int32)
	p.data = p.data[:n-1]
}

func execEQInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) == p.data[n-1].(int64)
	p.data = p.data[:n-1]
}

func execEQUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) == p.data[n-1].(uint)
	p.data = p.data[:n-1]
}

func execEQUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) == p.data[n-1].(uint8)
	p.data = p.data[:n-1]
}

func execEQUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) == p.data[n-1].(uint16)
	p.data = p.data[:n-1]
}

func execEQUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) == p.data[n-1].(uint32)
	p.data = p.data[:n-1]
}

func execEQUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) == p.data[n-1].(uint64)
	p.data = p.data[:n-1]
}

func execEQUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) == p.data[n-1].(uintptr)
	p.data = p.data[:n-1]
}

func execEQFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float32) == p.data[n-1].(float32)
	p.data = p.data[:n-1]
}

func execEQFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float64) == p.data[n-1].(float64)
	p.data = p.data[:n-1]
}

func execEQComplex64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(complex64) == p.data[n-1].(complex64)
	p.data = p.data[:n-1]
}

func execEQComplex128(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(complex128) == p.data[n-1].(complex128)
	p.data = p.data[:n-1]
}

func execEQString(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(string) == p.data[n-1].(string)
	p.data = p.data[:n-1]
}

func execNEInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) != p.data[n-1].(int)
	p.data = p.data[:n-1]
}

func execNEInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) != p.data[n-1].(int8)
	p.data = p.data[:n-1]
}

func execNEInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) != p.data[n-1].(int16)
	p.data = p.data[:n-1]
}

func execNEInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) != p.data[n-1].(int32)
	p.data = p.data[:n-1]
}

func execNEInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) != p.data[n-1].(int64)
	p.data = p.data[:n-1]
}

func execNEUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) != p.data[n-1].(uint)
	p.data = p.data[:n-1]
}

func execNEUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) != p.data[n-1].(uint8)
	p.data = p.data[:n-1]
}

func execNEUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) != p.data[n-1].(uint16)
	p.data = p.data[:n-1]
}

func execNEUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) != p.data[n-1].(uint32)
	p.data = p.data[:n-1]
}

func execNEUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) != p.data[n-1].(uint64)
	p.data = p.data[:n-1]
}

func execNEUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) != p.data[n-1].(uintptr)
	p.data = p.data[:n-1]
}

func execNEFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float32) != p.data[n-1].(float32)
	p.data = p.data[:n-1]
}

func execNEFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float64) != p.data[n-1].(float64)
	p.data = p.data[:n-1]
}

func execNEComplex64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(complex64) != p.data[n-1].(complex64)
	p.data = p.data[:n-1]
}

func execNEComplex128(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(complex128) != p.data[n-1].(complex128)
	p.data = p.data[:n-1]
}

func execNEString(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(string) != p.data[n-1].(string)
	p.data = p.data[:n-1]
}

func execLAndBool(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(bool) && p.data[n-1].(bool)
	p.data = p.data[:n-1]
}

func execLOrBool(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(bool) || p.data[n-1].(bool)
	p.data = p.data[:n-1]
}

func execNegInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = -p.data[n-1].(int)
}

func execNegInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = -p.data[n-1].(int8)
}

func execNegInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = -p.data[n-1].(int16)
}

func execNegInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = -p.data[n-1].(int32)
}

func execNegInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = -p.data[n-1].(int64)
}

func execNegUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = -p.data[n-1].(uint)
}

func execNegUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = -p.data[n-1].(uint8)
}

func execNegUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = -p.data[n-1].(uint16)
}

func execNegUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = -p.data[n-1].(uint32)
}

func execNegUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = -p.data[n-1].(uint64)
}

func execNegUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = -p.data[n-1].(uintptr)
}

func execNegFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = -p.data[n-1].(float32)
}

func execNegFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = -p.data[n-1].(float64)
}

func execNegComplex64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = -p.data[n-1].(complex64)
}

func execNegComplex128(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = -p.data[n-1].(complex128)
}

func execNotBool(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = !p.data[n-1].(bool)
}

func execBitNotInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = ^p.data[n-1].(int)
}

func execBitNotInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = ^p.data[n-1].(int8)
}

func execBitNotInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = ^p.data[n-1].(int16)
}

func execBitNotInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = ^p.data[n-1].(int32)
}

func execBitNotInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = ^p.data[n-1].(int64)
}

func execBitNotUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = ^p.data[n-1].(uint)
}

func execBitNotUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = ^p.data[n-1].(uint8)
}

func execBitNotUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = ^p.data[n-1].(uint16)
}

func execBitNotUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = ^p.data[n-1].(uint32)
}

func execBitNotUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = ^p.data[n-1].(uint64)
}

func execBitNotUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = ^p.data[n-1].(uintptr)
}
