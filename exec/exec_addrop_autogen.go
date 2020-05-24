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

var builtinAddrOps = [...]func(i Instr, p *Context){
	(int(OpAddAssign) << bitsKind) | int(Int):           execAddAssignInt,
	(int(OpAddAssign) << bitsKind) | int(Int8):          execAddAssignInt8,
	(int(OpAddAssign) << bitsKind) | int(Int16):         execAddAssignInt16,
	(int(OpAddAssign) << bitsKind) | int(Int32):         execAddAssignInt32,
	(int(OpAddAssign) << bitsKind) | int(Int64):         execAddAssignInt64,
	(int(OpAddAssign) << bitsKind) | int(Uint):          execAddAssignUint,
	(int(OpAddAssign) << bitsKind) | int(Uint8):         execAddAssignUint8,
	(int(OpAddAssign) << bitsKind) | int(Uint16):        execAddAssignUint16,
	(int(OpAddAssign) << bitsKind) | int(Uint32):        execAddAssignUint32,
	(int(OpAddAssign) << bitsKind) | int(Uint64):        execAddAssignUint64,
	(int(OpAddAssign) << bitsKind) | int(Uintptr):       execAddAssignUintptr,
	(int(OpAddAssign) << bitsKind) | int(Float32):       execAddAssignFloat32,
	(int(OpAddAssign) << bitsKind) | int(Float64):       execAddAssignFloat64,
	(int(OpAddAssign) << bitsKind) | int(Complex64):     execAddAssignComplex64,
	(int(OpAddAssign) << bitsKind) | int(Complex128):    execAddAssignComplex128,
	(int(OpAddAssign) << bitsKind) | int(String):        execAddAssignString,
	(int(OpSubAssign) << bitsKind) | int(Int):           execSubAssignInt,
	(int(OpSubAssign) << bitsKind) | int(Int8):          execSubAssignInt8,
	(int(OpSubAssign) << bitsKind) | int(Int16):         execSubAssignInt16,
	(int(OpSubAssign) << bitsKind) | int(Int32):         execSubAssignInt32,
	(int(OpSubAssign) << bitsKind) | int(Int64):         execSubAssignInt64,
	(int(OpSubAssign) << bitsKind) | int(Uint):          execSubAssignUint,
	(int(OpSubAssign) << bitsKind) | int(Uint8):         execSubAssignUint8,
	(int(OpSubAssign) << bitsKind) | int(Uint16):        execSubAssignUint16,
	(int(OpSubAssign) << bitsKind) | int(Uint32):        execSubAssignUint32,
	(int(OpSubAssign) << bitsKind) | int(Uint64):        execSubAssignUint64,
	(int(OpSubAssign) << bitsKind) | int(Uintptr):       execSubAssignUintptr,
	(int(OpSubAssign) << bitsKind) | int(Float32):       execSubAssignFloat32,
	(int(OpSubAssign) << bitsKind) | int(Float64):       execSubAssignFloat64,
	(int(OpSubAssign) << bitsKind) | int(Complex64):     execSubAssignComplex64,
	(int(OpSubAssign) << bitsKind) | int(Complex128):    execSubAssignComplex128,
	(int(OpMulAssign) << bitsKind) | int(Int):           execMulAssignInt,
	(int(OpMulAssign) << bitsKind) | int(Int8):          execMulAssignInt8,
	(int(OpMulAssign) << bitsKind) | int(Int16):         execMulAssignInt16,
	(int(OpMulAssign) << bitsKind) | int(Int32):         execMulAssignInt32,
	(int(OpMulAssign) << bitsKind) | int(Int64):         execMulAssignInt64,
	(int(OpMulAssign) << bitsKind) | int(Uint):          execMulAssignUint,
	(int(OpMulAssign) << bitsKind) | int(Uint8):         execMulAssignUint8,
	(int(OpMulAssign) << bitsKind) | int(Uint16):        execMulAssignUint16,
	(int(OpMulAssign) << bitsKind) | int(Uint32):        execMulAssignUint32,
	(int(OpMulAssign) << bitsKind) | int(Uint64):        execMulAssignUint64,
	(int(OpMulAssign) << bitsKind) | int(Uintptr):       execMulAssignUintptr,
	(int(OpMulAssign) << bitsKind) | int(Float32):       execMulAssignFloat32,
	(int(OpMulAssign) << bitsKind) | int(Float64):       execMulAssignFloat64,
	(int(OpMulAssign) << bitsKind) | int(Complex64):     execMulAssignComplex64,
	(int(OpMulAssign) << bitsKind) | int(Complex128):    execMulAssignComplex128,
	(int(OpDivAssign) << bitsKind) | int(Int):           execDivAssignInt,
	(int(OpDivAssign) << bitsKind) | int(Int8):          execDivAssignInt8,
	(int(OpDivAssign) << bitsKind) | int(Int16):         execDivAssignInt16,
	(int(OpDivAssign) << bitsKind) | int(Int32):         execDivAssignInt32,
	(int(OpDivAssign) << bitsKind) | int(Int64):         execDivAssignInt64,
	(int(OpDivAssign) << bitsKind) | int(Uint):          execDivAssignUint,
	(int(OpDivAssign) << bitsKind) | int(Uint8):         execDivAssignUint8,
	(int(OpDivAssign) << bitsKind) | int(Uint16):        execDivAssignUint16,
	(int(OpDivAssign) << bitsKind) | int(Uint32):        execDivAssignUint32,
	(int(OpDivAssign) << bitsKind) | int(Uint64):        execDivAssignUint64,
	(int(OpDivAssign) << bitsKind) | int(Uintptr):       execDivAssignUintptr,
	(int(OpDivAssign) << bitsKind) | int(Float32):       execDivAssignFloat32,
	(int(OpDivAssign) << bitsKind) | int(Float64):       execDivAssignFloat64,
	(int(OpDivAssign) << bitsKind) | int(Complex64):     execDivAssignComplex64,
	(int(OpDivAssign) << bitsKind) | int(Complex128):    execDivAssignComplex128,
	(int(OpModAssign) << bitsKind) | int(Int):           execModAssignInt,
	(int(OpModAssign) << bitsKind) | int(Int8):          execModAssignInt8,
	(int(OpModAssign) << bitsKind) | int(Int16):         execModAssignInt16,
	(int(OpModAssign) << bitsKind) | int(Int32):         execModAssignInt32,
	(int(OpModAssign) << bitsKind) | int(Int64):         execModAssignInt64,
	(int(OpModAssign) << bitsKind) | int(Uint):          execModAssignUint,
	(int(OpModAssign) << bitsKind) | int(Uint8):         execModAssignUint8,
	(int(OpModAssign) << bitsKind) | int(Uint16):        execModAssignUint16,
	(int(OpModAssign) << bitsKind) | int(Uint32):        execModAssignUint32,
	(int(OpModAssign) << bitsKind) | int(Uint64):        execModAssignUint64,
	(int(OpModAssign) << bitsKind) | int(Uintptr):       execModAssignUintptr,
	(int(OpBitAndAssign) << bitsKind) | int(Int):        execBitAndAssignInt,
	(int(OpBitAndAssign) << bitsKind) | int(Int8):       execBitAndAssignInt8,
	(int(OpBitAndAssign) << bitsKind) | int(Int16):      execBitAndAssignInt16,
	(int(OpBitAndAssign) << bitsKind) | int(Int32):      execBitAndAssignInt32,
	(int(OpBitAndAssign) << bitsKind) | int(Int64):      execBitAndAssignInt64,
	(int(OpBitAndAssign) << bitsKind) | int(Uint):       execBitAndAssignUint,
	(int(OpBitAndAssign) << bitsKind) | int(Uint8):      execBitAndAssignUint8,
	(int(OpBitAndAssign) << bitsKind) | int(Uint16):     execBitAndAssignUint16,
	(int(OpBitAndAssign) << bitsKind) | int(Uint32):     execBitAndAssignUint32,
	(int(OpBitAndAssign) << bitsKind) | int(Uint64):     execBitAndAssignUint64,
	(int(OpBitAndAssign) << bitsKind) | int(Uintptr):    execBitAndAssignUintptr,
	(int(OpBitOrAssign) << bitsKind) | int(Int):         execBitOrAssignInt,
	(int(OpBitOrAssign) << bitsKind) | int(Int8):        execBitOrAssignInt8,
	(int(OpBitOrAssign) << bitsKind) | int(Int16):       execBitOrAssignInt16,
	(int(OpBitOrAssign) << bitsKind) | int(Int32):       execBitOrAssignInt32,
	(int(OpBitOrAssign) << bitsKind) | int(Int64):       execBitOrAssignInt64,
	(int(OpBitOrAssign) << bitsKind) | int(Uint):        execBitOrAssignUint,
	(int(OpBitOrAssign) << bitsKind) | int(Uint8):       execBitOrAssignUint8,
	(int(OpBitOrAssign) << bitsKind) | int(Uint16):      execBitOrAssignUint16,
	(int(OpBitOrAssign) << bitsKind) | int(Uint32):      execBitOrAssignUint32,
	(int(OpBitOrAssign) << bitsKind) | int(Uint64):      execBitOrAssignUint64,
	(int(OpBitOrAssign) << bitsKind) | int(Uintptr):     execBitOrAssignUintptr,
	(int(OpBitXorAssign) << bitsKind) | int(Int):        execBitXorAssignInt,
	(int(OpBitXorAssign) << bitsKind) | int(Int8):       execBitXorAssignInt8,
	(int(OpBitXorAssign) << bitsKind) | int(Int16):      execBitXorAssignInt16,
	(int(OpBitXorAssign) << bitsKind) | int(Int32):      execBitXorAssignInt32,
	(int(OpBitXorAssign) << bitsKind) | int(Int64):      execBitXorAssignInt64,
	(int(OpBitXorAssign) << bitsKind) | int(Uint):       execBitXorAssignUint,
	(int(OpBitXorAssign) << bitsKind) | int(Uint8):      execBitXorAssignUint8,
	(int(OpBitXorAssign) << bitsKind) | int(Uint16):     execBitXorAssignUint16,
	(int(OpBitXorAssign) << bitsKind) | int(Uint32):     execBitXorAssignUint32,
	(int(OpBitXorAssign) << bitsKind) | int(Uint64):     execBitXorAssignUint64,
	(int(OpBitXorAssign) << bitsKind) | int(Uintptr):    execBitXorAssignUintptr,
	(int(OpBitAndNotAssign) << bitsKind) | int(Int):     execBitAndNotAssignInt,
	(int(OpBitAndNotAssign) << bitsKind) | int(Int8):    execBitAndNotAssignInt8,
	(int(OpBitAndNotAssign) << bitsKind) | int(Int16):   execBitAndNotAssignInt16,
	(int(OpBitAndNotAssign) << bitsKind) | int(Int32):   execBitAndNotAssignInt32,
	(int(OpBitAndNotAssign) << bitsKind) | int(Int64):   execBitAndNotAssignInt64,
	(int(OpBitAndNotAssign) << bitsKind) | int(Uint):    execBitAndNotAssignUint,
	(int(OpBitAndNotAssign) << bitsKind) | int(Uint8):   execBitAndNotAssignUint8,
	(int(OpBitAndNotAssign) << bitsKind) | int(Uint16):  execBitAndNotAssignUint16,
	(int(OpBitAndNotAssign) << bitsKind) | int(Uint32):  execBitAndNotAssignUint32,
	(int(OpBitAndNotAssign) << bitsKind) | int(Uint64):  execBitAndNotAssignUint64,
	(int(OpBitAndNotAssign) << bitsKind) | int(Uintptr): execBitAndNotAssignUintptr,
	(int(OpBitSHLAssign) << bitsKind) | int(Int):        execBitSHLAssignInt,
	(int(OpBitSHLAssign) << bitsKind) | int(Int8):       execBitSHLAssignInt8,
	(int(OpBitSHLAssign) << bitsKind) | int(Int16):      execBitSHLAssignInt16,
	(int(OpBitSHLAssign) << bitsKind) | int(Int32):      execBitSHLAssignInt32,
	(int(OpBitSHLAssign) << bitsKind) | int(Int64):      execBitSHLAssignInt64,
	(int(OpBitSHLAssign) << bitsKind) | int(Uint):       execBitSHLAssignUint,
	(int(OpBitSHLAssign) << bitsKind) | int(Uint8):      execBitSHLAssignUint8,
	(int(OpBitSHLAssign) << bitsKind) | int(Uint16):     execBitSHLAssignUint16,
	(int(OpBitSHLAssign) << bitsKind) | int(Uint32):     execBitSHLAssignUint32,
	(int(OpBitSHLAssign) << bitsKind) | int(Uint64):     execBitSHLAssignUint64,
	(int(OpBitSHLAssign) << bitsKind) | int(Uintptr):    execBitSHLAssignUintptr,
	(int(OpBitSHRAssign) << bitsKind) | int(Int):        execBitSHRAssignInt,
	(int(OpBitSHRAssign) << bitsKind) | int(Int8):       execBitSHRAssignInt8,
	(int(OpBitSHRAssign) << bitsKind) | int(Int16):      execBitSHRAssignInt16,
	(int(OpBitSHRAssign) << bitsKind) | int(Int32):      execBitSHRAssignInt32,
	(int(OpBitSHRAssign) << bitsKind) | int(Int64):      execBitSHRAssignInt64,
	(int(OpBitSHRAssign) << bitsKind) | int(Uint):       execBitSHRAssignUint,
	(int(OpBitSHRAssign) << bitsKind) | int(Uint8):      execBitSHRAssignUint8,
	(int(OpBitSHRAssign) << bitsKind) | int(Uint16):     execBitSHRAssignUint16,
	(int(OpBitSHRAssign) << bitsKind) | int(Uint32):     execBitSHRAssignUint32,
	(int(OpBitSHRAssign) << bitsKind) | int(Uint64):     execBitSHRAssignUint64,
	(int(OpBitSHRAssign) << bitsKind) | int(Uintptr):    execBitSHRAssignUintptr,
	(int(OpInc) << bitsKind) | int(Int):                 execIncInt,
	(int(OpInc) << bitsKind) | int(Int8):                execIncInt8,
	(int(OpInc) << bitsKind) | int(Int16):               execIncInt16,
	(int(OpInc) << bitsKind) | int(Int32):               execIncInt32,
	(int(OpInc) << bitsKind) | int(Int64):               execIncInt64,
	(int(OpInc) << bitsKind) | int(Uint):                execIncUint,
	(int(OpInc) << bitsKind) | int(Uint8):               execIncUint8,
	(int(OpInc) << bitsKind) | int(Uint16):              execIncUint16,
	(int(OpInc) << bitsKind) | int(Uint32):              execIncUint32,
	(int(OpInc) << bitsKind) | int(Uint64):              execIncUint64,
	(int(OpInc) << bitsKind) | int(Uintptr):             execIncUintptr,
	(int(OpInc) << bitsKind) | int(Float32):             execIncFloat32,
	(int(OpInc) << bitsKind) | int(Float64):             execIncFloat64,
	(int(OpInc) << bitsKind) | int(Complex64):           execIncComplex64,
	(int(OpInc) << bitsKind) | int(Complex128):          execIncComplex128,
	(int(OpDec) << bitsKind) | int(Int):                 execDecInt,
	(int(OpDec) << bitsKind) | int(Int8):                execDecInt8,
	(int(OpDec) << bitsKind) | int(Int16):               execDecInt16,
	(int(OpDec) << bitsKind) | int(Int32):               execDecInt32,
	(int(OpDec) << bitsKind) | int(Int64):               execDecInt64,
	(int(OpDec) << bitsKind) | int(Uint):                execDecUint,
	(int(OpDec) << bitsKind) | int(Uint8):               execDecUint8,
	(int(OpDec) << bitsKind) | int(Uint16):              execDecUint16,
	(int(OpDec) << bitsKind) | int(Uint32):              execDecUint32,
	(int(OpDec) << bitsKind) | int(Uint64):              execDecUint64,
	(int(OpDec) << bitsKind) | int(Uintptr):             execDecUintptr,
	(int(OpDec) << bitsKind) | int(Float32):             execDecFloat32,
	(int(OpDec) << bitsKind) | int(Float64):             execDecFloat64,
	(int(OpDec) << bitsKind) | int(Complex64):           execDecComplex64,
	(int(OpDec) << bitsKind) | int(Complex128):          execDecComplex128,
}

func execAddAssignInt(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int) += p.data[n-2].(int)
	p.data = p.data[:n-2]
}

func execAddAssignInt8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int8) += p.data[n-2].(int8)
	p.data = p.data[:n-2]
}

func execAddAssignInt16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int16) += p.data[n-2].(int16)
	p.data = p.data[:n-2]
}

func execAddAssignInt32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int32) += p.data[n-2].(int32)
	p.data = p.data[:n-2]
}

func execAddAssignInt64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int64) += p.data[n-2].(int64)
	p.data = p.data[:n-2]
}

func execAddAssignUint(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint) += p.data[n-2].(uint)
	p.data = p.data[:n-2]
}

func execAddAssignUint8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint8) += p.data[n-2].(uint8)
	p.data = p.data[:n-2]
}

func execAddAssignUint16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint16) += p.data[n-2].(uint16)
	p.data = p.data[:n-2]
}

func execAddAssignUint32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint32) += p.data[n-2].(uint32)
	p.data = p.data[:n-2]
}

func execAddAssignUint64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint64) += p.data[n-2].(uint64)
	p.data = p.data[:n-2]
}

func execAddAssignUintptr(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uintptr) += p.data[n-2].(uintptr)
	p.data = p.data[:n-2]
}

func execAddAssignFloat32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*float32) += p.data[n-2].(float32)
	p.data = p.data[:n-2]
}

func execAddAssignFloat64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*float64) += p.data[n-2].(float64)
	p.data = p.data[:n-2]
}

func execAddAssignComplex64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*complex64) += p.data[n-2].(complex64)
	p.data = p.data[:n-2]
}

func execAddAssignComplex128(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*complex128) += p.data[n-2].(complex128)
	p.data = p.data[:n-2]
}

func execAddAssignString(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*string) += p.data[n-2].(string)
	p.data = p.data[:n-2]
}

func execSubAssignInt(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int) -= p.data[n-2].(int)
	p.data = p.data[:n-2]
}

func execSubAssignInt8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int8) -= p.data[n-2].(int8)
	p.data = p.data[:n-2]
}

func execSubAssignInt16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int16) -= p.data[n-2].(int16)
	p.data = p.data[:n-2]
}

func execSubAssignInt32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int32) -= p.data[n-2].(int32)
	p.data = p.data[:n-2]
}

func execSubAssignInt64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int64) -= p.data[n-2].(int64)
	p.data = p.data[:n-2]
}

func execSubAssignUint(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint) -= p.data[n-2].(uint)
	p.data = p.data[:n-2]
}

func execSubAssignUint8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint8) -= p.data[n-2].(uint8)
	p.data = p.data[:n-2]
}

func execSubAssignUint16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint16) -= p.data[n-2].(uint16)
	p.data = p.data[:n-2]
}

func execSubAssignUint32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint32) -= p.data[n-2].(uint32)
	p.data = p.data[:n-2]
}

func execSubAssignUint64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint64) -= p.data[n-2].(uint64)
	p.data = p.data[:n-2]
}

func execSubAssignUintptr(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uintptr) -= p.data[n-2].(uintptr)
	p.data = p.data[:n-2]
}

func execSubAssignFloat32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*float32) -= p.data[n-2].(float32)
	p.data = p.data[:n-2]
}

func execSubAssignFloat64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*float64) -= p.data[n-2].(float64)
	p.data = p.data[:n-2]
}

func execSubAssignComplex64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*complex64) -= p.data[n-2].(complex64)
	p.data = p.data[:n-2]
}

func execSubAssignComplex128(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*complex128) -= p.data[n-2].(complex128)
	p.data = p.data[:n-2]
}

func execMulAssignInt(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int) *= p.data[n-2].(int)
	p.data = p.data[:n-2]
}

func execMulAssignInt8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int8) *= p.data[n-2].(int8)
	p.data = p.data[:n-2]
}

func execMulAssignInt16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int16) *= p.data[n-2].(int16)
	p.data = p.data[:n-2]
}

func execMulAssignInt32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int32) *= p.data[n-2].(int32)
	p.data = p.data[:n-2]
}

func execMulAssignInt64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int64) *= p.data[n-2].(int64)
	p.data = p.data[:n-2]
}

func execMulAssignUint(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint) *= p.data[n-2].(uint)
	p.data = p.data[:n-2]
}

func execMulAssignUint8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint8) *= p.data[n-2].(uint8)
	p.data = p.data[:n-2]
}

func execMulAssignUint16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint16) *= p.data[n-2].(uint16)
	p.data = p.data[:n-2]
}

func execMulAssignUint32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint32) *= p.data[n-2].(uint32)
	p.data = p.data[:n-2]
}

func execMulAssignUint64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint64) *= p.data[n-2].(uint64)
	p.data = p.data[:n-2]
}

func execMulAssignUintptr(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uintptr) *= p.data[n-2].(uintptr)
	p.data = p.data[:n-2]
}

func execMulAssignFloat32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*float32) *= p.data[n-2].(float32)
	p.data = p.data[:n-2]
}

func execMulAssignFloat64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*float64) *= p.data[n-2].(float64)
	p.data = p.data[:n-2]
}

func execMulAssignComplex64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*complex64) *= p.data[n-2].(complex64)
	p.data = p.data[:n-2]
}

func execMulAssignComplex128(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*complex128) *= p.data[n-2].(complex128)
	p.data = p.data[:n-2]
}

func execDivAssignInt(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int) /= p.data[n-2].(int)
	p.data = p.data[:n-2]
}

func execDivAssignInt8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int8) /= p.data[n-2].(int8)
	p.data = p.data[:n-2]
}

func execDivAssignInt16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int16) /= p.data[n-2].(int16)
	p.data = p.data[:n-2]
}

func execDivAssignInt32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int32) /= p.data[n-2].(int32)
	p.data = p.data[:n-2]
}

func execDivAssignInt64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int64) /= p.data[n-2].(int64)
	p.data = p.data[:n-2]
}

func execDivAssignUint(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint) /= p.data[n-2].(uint)
	p.data = p.data[:n-2]
}

func execDivAssignUint8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint8) /= p.data[n-2].(uint8)
	p.data = p.data[:n-2]
}

func execDivAssignUint16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint16) /= p.data[n-2].(uint16)
	p.data = p.data[:n-2]
}

func execDivAssignUint32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint32) /= p.data[n-2].(uint32)
	p.data = p.data[:n-2]
}

func execDivAssignUint64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint64) /= p.data[n-2].(uint64)
	p.data = p.data[:n-2]
}

func execDivAssignUintptr(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uintptr) /= p.data[n-2].(uintptr)
	p.data = p.data[:n-2]
}

func execDivAssignFloat32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*float32) /= p.data[n-2].(float32)
	p.data = p.data[:n-2]
}

func execDivAssignFloat64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*float64) /= p.data[n-2].(float64)
	p.data = p.data[:n-2]
}

func execDivAssignComplex64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*complex64) /= p.data[n-2].(complex64)
	p.data = p.data[:n-2]
}

func execDivAssignComplex128(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*complex128) /= p.data[n-2].(complex128)
	p.data = p.data[:n-2]
}

func execModAssignInt(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int) %= p.data[n-2].(int)
	p.data = p.data[:n-2]
}

func execModAssignInt8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int8) %= p.data[n-2].(int8)
	p.data = p.data[:n-2]
}

func execModAssignInt16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int16) %= p.data[n-2].(int16)
	p.data = p.data[:n-2]
}

func execModAssignInt32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int32) %= p.data[n-2].(int32)
	p.data = p.data[:n-2]
}

func execModAssignInt64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int64) %= p.data[n-2].(int64)
	p.data = p.data[:n-2]
}

func execModAssignUint(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint) %= p.data[n-2].(uint)
	p.data = p.data[:n-2]
}

func execModAssignUint8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint8) %= p.data[n-2].(uint8)
	p.data = p.data[:n-2]
}

func execModAssignUint16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint16) %= p.data[n-2].(uint16)
	p.data = p.data[:n-2]
}

func execModAssignUint32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint32) %= p.data[n-2].(uint32)
	p.data = p.data[:n-2]
}

func execModAssignUint64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint64) %= p.data[n-2].(uint64)
	p.data = p.data[:n-2]
}

func execModAssignUintptr(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uintptr) %= p.data[n-2].(uintptr)
	p.data = p.data[:n-2]
}

func execBitAndAssignInt(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int) &= p.data[n-2].(int)
	p.data = p.data[:n-2]
}

func execBitAndAssignInt8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int8) &= p.data[n-2].(int8)
	p.data = p.data[:n-2]
}

func execBitAndAssignInt16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int16) &= p.data[n-2].(int16)
	p.data = p.data[:n-2]
}

func execBitAndAssignInt32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int32) &= p.data[n-2].(int32)
	p.data = p.data[:n-2]
}

func execBitAndAssignInt64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int64) &= p.data[n-2].(int64)
	p.data = p.data[:n-2]
}

func execBitAndAssignUint(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint) &= p.data[n-2].(uint)
	p.data = p.data[:n-2]
}

func execBitAndAssignUint8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint8) &= p.data[n-2].(uint8)
	p.data = p.data[:n-2]
}

func execBitAndAssignUint16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint16) &= p.data[n-2].(uint16)
	p.data = p.data[:n-2]
}

func execBitAndAssignUint32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint32) &= p.data[n-2].(uint32)
	p.data = p.data[:n-2]
}

func execBitAndAssignUint64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint64) &= p.data[n-2].(uint64)
	p.data = p.data[:n-2]
}

func execBitAndAssignUintptr(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uintptr) &= p.data[n-2].(uintptr)
	p.data = p.data[:n-2]
}

func execBitOrAssignInt(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int) |= p.data[n-2].(int)
	p.data = p.data[:n-2]
}

func execBitOrAssignInt8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int8) |= p.data[n-2].(int8)
	p.data = p.data[:n-2]
}

func execBitOrAssignInt16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int16) |= p.data[n-2].(int16)
	p.data = p.data[:n-2]
}

func execBitOrAssignInt32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int32) |= p.data[n-2].(int32)
	p.data = p.data[:n-2]
}

func execBitOrAssignInt64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int64) |= p.data[n-2].(int64)
	p.data = p.data[:n-2]
}

func execBitOrAssignUint(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint) |= p.data[n-2].(uint)
	p.data = p.data[:n-2]
}

func execBitOrAssignUint8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint8) |= p.data[n-2].(uint8)
	p.data = p.data[:n-2]
}

func execBitOrAssignUint16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint16) |= p.data[n-2].(uint16)
	p.data = p.data[:n-2]
}

func execBitOrAssignUint32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint32) |= p.data[n-2].(uint32)
	p.data = p.data[:n-2]
}

func execBitOrAssignUint64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint64) |= p.data[n-2].(uint64)
	p.data = p.data[:n-2]
}

func execBitOrAssignUintptr(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uintptr) |= p.data[n-2].(uintptr)
	p.data = p.data[:n-2]
}

func execBitXorAssignInt(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int) ^= p.data[n-2].(int)
	p.data = p.data[:n-2]
}

func execBitXorAssignInt8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int8) ^= p.data[n-2].(int8)
	p.data = p.data[:n-2]
}

func execBitXorAssignInt16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int16) ^= p.data[n-2].(int16)
	p.data = p.data[:n-2]
}

func execBitXorAssignInt32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int32) ^= p.data[n-2].(int32)
	p.data = p.data[:n-2]
}

func execBitXorAssignInt64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int64) ^= p.data[n-2].(int64)
	p.data = p.data[:n-2]
}

func execBitXorAssignUint(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint) ^= p.data[n-2].(uint)
	p.data = p.data[:n-2]
}

func execBitXorAssignUint8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint8) ^= p.data[n-2].(uint8)
	p.data = p.data[:n-2]
}

func execBitXorAssignUint16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint16) ^= p.data[n-2].(uint16)
	p.data = p.data[:n-2]
}

func execBitXorAssignUint32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint32) ^= p.data[n-2].(uint32)
	p.data = p.data[:n-2]
}

func execBitXorAssignUint64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint64) ^= p.data[n-2].(uint64)
	p.data = p.data[:n-2]
}

func execBitXorAssignUintptr(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uintptr) ^= p.data[n-2].(uintptr)
	p.data = p.data[:n-2]
}

func execBitAndNotAssignInt(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int) &^= p.data[n-2].(int)
	p.data = p.data[:n-2]
}

func execBitAndNotAssignInt8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int8) &^= p.data[n-2].(int8)
	p.data = p.data[:n-2]
}

func execBitAndNotAssignInt16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int16) &^= p.data[n-2].(int16)
	p.data = p.data[:n-2]
}

func execBitAndNotAssignInt32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int32) &^= p.data[n-2].(int32)
	p.data = p.data[:n-2]
}

func execBitAndNotAssignInt64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int64) &^= p.data[n-2].(int64)
	p.data = p.data[:n-2]
}

func execBitAndNotAssignUint(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint) &^= p.data[n-2].(uint)
	p.data = p.data[:n-2]
}

func execBitAndNotAssignUint8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint8) &^= p.data[n-2].(uint8)
	p.data = p.data[:n-2]
}

func execBitAndNotAssignUint16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint16) &^= p.data[n-2].(uint16)
	p.data = p.data[:n-2]
}

func execBitAndNotAssignUint32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint32) &^= p.data[n-2].(uint32)
	p.data = p.data[:n-2]
}

func execBitAndNotAssignUint64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint64) &^= p.data[n-2].(uint64)
	p.data = p.data[:n-2]
}

func execBitAndNotAssignUintptr(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uintptr) &^= p.data[n-2].(uintptr)
	p.data = p.data[:n-2]
}

func execBitSHLAssignInt(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int) <<= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execBitSHLAssignInt8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int8) <<= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execBitSHLAssignInt16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int16) <<= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execBitSHLAssignInt32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int32) <<= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execBitSHLAssignInt64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int64) <<= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execBitSHLAssignUint(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint) <<= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execBitSHLAssignUint8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint8) <<= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execBitSHLAssignUint16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint16) <<= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execBitSHLAssignUint32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint32) <<= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execBitSHLAssignUint64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint64) <<= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execBitSHLAssignUintptr(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uintptr) <<= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execBitSHRAssignInt(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int) >>= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execBitSHRAssignInt8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int8) >>= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execBitSHRAssignInt16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int16) >>= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execBitSHRAssignInt32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int32) >>= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execBitSHRAssignInt64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int64) >>= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execBitSHRAssignUint(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint) >>= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execBitSHRAssignUint8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint8) >>= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execBitSHRAssignUint16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint16) >>= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execBitSHRAssignUint32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint32) >>= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execBitSHRAssignUint64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint64) >>= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execBitSHRAssignUintptr(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uintptr) >>= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execIncInt(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*int))++
	p.data = p.data[:n-1]
}

func execIncInt8(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*int8))++
	p.data = p.data[:n-1]
}

func execIncInt16(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*int16))++
	p.data = p.data[:n-1]
}

func execIncInt32(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*int32))++
	p.data = p.data[:n-1]
}

func execIncInt64(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*int64))++
	p.data = p.data[:n-1]
}

func execIncUint(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*uint))++
	p.data = p.data[:n-1]
}

func execIncUint8(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*uint8))++
	p.data = p.data[:n-1]
}

func execIncUint16(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*uint16))++
	p.data = p.data[:n-1]
}

func execIncUint32(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*uint32))++
	p.data = p.data[:n-1]
}

func execIncUint64(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*uint64))++
	p.data = p.data[:n-1]
}

func execIncUintptr(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*uintptr))++
	p.data = p.data[:n-1]
}

func execIncFloat32(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*float32))++
	p.data = p.data[:n-1]
}

func execIncFloat64(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*float64))++
	p.data = p.data[:n-1]
}

func execIncComplex64(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*complex64))++
	p.data = p.data[:n-1]
}

func execIncComplex128(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*complex128))++
	p.data = p.data[:n-1]
}

func execDecInt(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*int))--
	p.data = p.data[:n-1]
}

func execDecInt8(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*int8))--
	p.data = p.data[:n-1]
}

func execDecInt16(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*int16))--
	p.data = p.data[:n-1]
}

func execDecInt32(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*int32))--
	p.data = p.data[:n-1]
}

func execDecInt64(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*int64))--
	p.data = p.data[:n-1]
}

func execDecUint(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*uint))--
	p.data = p.data[:n-1]
}

func execDecUint8(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*uint8))--
	p.data = p.data[:n-1]
}

func execDecUint16(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*uint16))--
	p.data = p.data[:n-1]
}

func execDecUint32(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*uint32))--
	p.data = p.data[:n-1]
}

func execDecUint64(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*uint64))--
	p.data = p.data[:n-1]
}

func execDecUintptr(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*uintptr))--
	p.data = p.data[:n-1]
}

func execDecFloat32(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*float32))--
	p.data = p.data[:n-1]
}

func execDecFloat64(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*float64))--
	p.data = p.data[:n-1]
}

func execDecComplex64(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*complex64))--
	p.data = p.data[:n-1]
}

func execDecComplex128(i Instr, p *Context) {
	n := len(p.data)
	(*p.data[n-1].(*complex128))--
	p.data = p.data[:n-1]
}
