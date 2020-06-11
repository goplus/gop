package bytecode

var builtinAddrOps = [...]func(i Instr, p *Context){
	(int(OpAddAssign) << bitsKind) | int(Int):        execAddAssignInt,
	(int(OpAddAssign) << bitsKind) | int(Int8):       execAddAssignInt8,
	(int(OpAddAssign) << bitsKind) | int(Int16):      execAddAssignInt16,
	(int(OpAddAssign) << bitsKind) | int(Int32):      execAddAssignInt32,
	(int(OpAddAssign) << bitsKind) | int(Int64):      execAddAssignInt64,
	(int(OpAddAssign) << bitsKind) | int(Uint):       execAddAssignUint,
	(int(OpAddAssign) << bitsKind) | int(Uint8):      execAddAssignUint8,
	(int(OpAddAssign) << bitsKind) | int(Uint16):     execAddAssignUint16,
	(int(OpAddAssign) << bitsKind) | int(Uint32):     execAddAssignUint32,
	(int(OpAddAssign) << bitsKind) | int(Uint64):     execAddAssignUint64,
	(int(OpAddAssign) << bitsKind) | int(Uintptr):    execAddAssignUintptr,
	(int(OpAddAssign) << bitsKind) | int(Float32):    execAddAssignFloat32,
	(int(OpAddAssign) << bitsKind) | int(Float64):    execAddAssignFloat64,
	(int(OpAddAssign) << bitsKind) | int(Complex64):  execAddAssignComplex64,
	(int(OpAddAssign) << bitsKind) | int(Complex128): execAddAssignComplex128,
	(int(OpAddAssign) << bitsKind) | int(String):     execAddAssignString,
	(int(OpSubAssign) << bitsKind) | int(Int):        execSubAssignInt,
	(int(OpSubAssign) << bitsKind) | int(Int8):       execSubAssignInt8,
	(int(OpSubAssign) << bitsKind) | int(Int16):      execSubAssignInt16,
	(int(OpSubAssign) << bitsKind) | int(Int32):      execSubAssignInt32,
	(int(OpSubAssign) << bitsKind) | int(Int64):      execSubAssignInt64,
	(int(OpSubAssign) << bitsKind) | int(Uint):       execSubAssignUint,
	(int(OpSubAssign) << bitsKind) | int(Uint8):      execSubAssignUint8,
	(int(OpSubAssign) << bitsKind) | int(Uint16):     execSubAssignUint16,
	(int(OpSubAssign) << bitsKind) | int(Uint32):     execSubAssignUint32,
	(int(OpSubAssign) << bitsKind) | int(Uint64):     execSubAssignUint64,
	(int(OpSubAssign) << bitsKind) | int(Uintptr):    execSubAssignUintptr,
	(int(OpSubAssign) << bitsKind) | int(Float32):    execSubAssignFloat32,
	(int(OpSubAssign) << bitsKind) | int(Float64):    execSubAssignFloat64,
	(int(OpSubAssign) << bitsKind) | int(Complex64):  execSubAssignComplex64,
	(int(OpSubAssign) << bitsKind) | int(Complex128): execSubAssignComplex128,
	(int(OpMulAssign) << bitsKind) | int(Int):        execMulAssignInt,
	(int(OpMulAssign) << bitsKind) | int(Int8):       execMulAssignInt8,
	(int(OpMulAssign) << bitsKind) | int(Int16):      execMulAssignInt16,
	(int(OpMulAssign) << bitsKind) | int(Int32):      execMulAssignInt32,
	(int(OpMulAssign) << bitsKind) | int(Int64):      execMulAssignInt64,
	(int(OpMulAssign) << bitsKind) | int(Uint):       execMulAssignUint,
	(int(OpMulAssign) << bitsKind) | int(Uint8):      execMulAssignUint8,
	(int(OpMulAssign) << bitsKind) | int(Uint16):     execMulAssignUint16,
	(int(OpMulAssign) << bitsKind) | int(Uint32):     execMulAssignUint32,
	(int(OpMulAssign) << bitsKind) | int(Uint64):     execMulAssignUint64,
	(int(OpMulAssign) << bitsKind) | int(Uintptr):    execMulAssignUintptr,
	(int(OpMulAssign) << bitsKind) | int(Float32):    execMulAssignFloat32,
	(int(OpMulAssign) << bitsKind) | int(Float64):    execMulAssignFloat64,
	(int(OpMulAssign) << bitsKind) | int(Complex64):  execMulAssignComplex64,
	(int(OpMulAssign) << bitsKind) | int(Complex128): execMulAssignComplex128,
	(int(OpQuoAssign) << bitsKind) | int(Int):        execQuoAssignInt,
	(int(OpQuoAssign) << bitsKind) | int(Int8):       execQuoAssignInt8,
	(int(OpQuoAssign) << bitsKind) | int(Int16):      execQuoAssignInt16,
	(int(OpQuoAssign) << bitsKind) | int(Int32):      execQuoAssignInt32,
	(int(OpQuoAssign) << bitsKind) | int(Int64):      execQuoAssignInt64,
	(int(OpQuoAssign) << bitsKind) | int(Uint):       execQuoAssignUint,
	(int(OpQuoAssign) << bitsKind) | int(Uint8):      execQuoAssignUint8,
	(int(OpQuoAssign) << bitsKind) | int(Uint16):     execQuoAssignUint16,
	(int(OpQuoAssign) << bitsKind) | int(Uint32):     execQuoAssignUint32,
	(int(OpQuoAssign) << bitsKind) | int(Uint64):     execQuoAssignUint64,
	(int(OpQuoAssign) << bitsKind) | int(Uintptr):    execQuoAssignUintptr,
	(int(OpQuoAssign) << bitsKind) | int(Float32):    execQuoAssignFloat32,
	(int(OpQuoAssign) << bitsKind) | int(Float64):    execQuoAssignFloat64,
	(int(OpQuoAssign) << bitsKind) | int(Complex64):  execQuoAssignComplex64,
	(int(OpQuoAssign) << bitsKind) | int(Complex128): execQuoAssignComplex128,
	(int(OpModAssign) << bitsKind) | int(Int):        execModAssignInt,
	(int(OpModAssign) << bitsKind) | int(Int8):       execModAssignInt8,
	(int(OpModAssign) << bitsKind) | int(Int16):      execModAssignInt16,
	(int(OpModAssign) << bitsKind) | int(Int32):      execModAssignInt32,
	(int(OpModAssign) << bitsKind) | int(Int64):      execModAssignInt64,
	(int(OpModAssign) << bitsKind) | int(Uint):       execModAssignUint,
	(int(OpModAssign) << bitsKind) | int(Uint8):      execModAssignUint8,
	(int(OpModAssign) << bitsKind) | int(Uint16):     execModAssignUint16,
	(int(OpModAssign) << bitsKind) | int(Uint32):     execModAssignUint32,
	(int(OpModAssign) << bitsKind) | int(Uint64):     execModAssignUint64,
	(int(OpModAssign) << bitsKind) | int(Uintptr):    execModAssignUintptr,
	(int(OpAndAssign) << bitsKind) | int(Int):        execAndAssignInt,
	(int(OpAndAssign) << bitsKind) | int(Int8):       execAndAssignInt8,
	(int(OpAndAssign) << bitsKind) | int(Int16):      execAndAssignInt16,
	(int(OpAndAssign) << bitsKind) | int(Int32):      execAndAssignInt32,
	(int(OpAndAssign) << bitsKind) | int(Int64):      execAndAssignInt64,
	(int(OpAndAssign) << bitsKind) | int(Uint):       execAndAssignUint,
	(int(OpAndAssign) << bitsKind) | int(Uint8):      execAndAssignUint8,
	(int(OpAndAssign) << bitsKind) | int(Uint16):     execAndAssignUint16,
	(int(OpAndAssign) << bitsKind) | int(Uint32):     execAndAssignUint32,
	(int(OpAndAssign) << bitsKind) | int(Uint64):     execAndAssignUint64,
	(int(OpAndAssign) << bitsKind) | int(Uintptr):    execAndAssignUintptr,
	(int(OpOrAssign) << bitsKind) | int(Int):         execOrAssignInt,
	(int(OpOrAssign) << bitsKind) | int(Int8):        execOrAssignInt8,
	(int(OpOrAssign) << bitsKind) | int(Int16):       execOrAssignInt16,
	(int(OpOrAssign) << bitsKind) | int(Int32):       execOrAssignInt32,
	(int(OpOrAssign) << bitsKind) | int(Int64):       execOrAssignInt64,
	(int(OpOrAssign) << bitsKind) | int(Uint):        execOrAssignUint,
	(int(OpOrAssign) << bitsKind) | int(Uint8):       execOrAssignUint8,
	(int(OpOrAssign) << bitsKind) | int(Uint16):      execOrAssignUint16,
	(int(OpOrAssign) << bitsKind) | int(Uint32):      execOrAssignUint32,
	(int(OpOrAssign) << bitsKind) | int(Uint64):      execOrAssignUint64,
	(int(OpOrAssign) << bitsKind) | int(Uintptr):     execOrAssignUintptr,
	(int(OpXorAssign) << bitsKind) | int(Int):        execXorAssignInt,
	(int(OpXorAssign) << bitsKind) | int(Int8):       execXorAssignInt8,
	(int(OpXorAssign) << bitsKind) | int(Int16):      execXorAssignInt16,
	(int(OpXorAssign) << bitsKind) | int(Int32):      execXorAssignInt32,
	(int(OpXorAssign) << bitsKind) | int(Int64):      execXorAssignInt64,
	(int(OpXorAssign) << bitsKind) | int(Uint):       execXorAssignUint,
	(int(OpXorAssign) << bitsKind) | int(Uint8):      execXorAssignUint8,
	(int(OpXorAssign) << bitsKind) | int(Uint16):     execXorAssignUint16,
	(int(OpXorAssign) << bitsKind) | int(Uint32):     execXorAssignUint32,
	(int(OpXorAssign) << bitsKind) | int(Uint64):     execXorAssignUint64,
	(int(OpXorAssign) << bitsKind) | int(Uintptr):    execXorAssignUintptr,
	(int(OpAndNotAssign) << bitsKind) | int(Int):     execAndNotAssignInt,
	(int(OpAndNotAssign) << bitsKind) | int(Int8):    execAndNotAssignInt8,
	(int(OpAndNotAssign) << bitsKind) | int(Int16):   execAndNotAssignInt16,
	(int(OpAndNotAssign) << bitsKind) | int(Int32):   execAndNotAssignInt32,
	(int(OpAndNotAssign) << bitsKind) | int(Int64):   execAndNotAssignInt64,
	(int(OpAndNotAssign) << bitsKind) | int(Uint):    execAndNotAssignUint,
	(int(OpAndNotAssign) << bitsKind) | int(Uint8):   execAndNotAssignUint8,
	(int(OpAndNotAssign) << bitsKind) | int(Uint16):  execAndNotAssignUint16,
	(int(OpAndNotAssign) << bitsKind) | int(Uint32):  execAndNotAssignUint32,
	(int(OpAndNotAssign) << bitsKind) | int(Uint64):  execAndNotAssignUint64,
	(int(OpAndNotAssign) << bitsKind) | int(Uintptr): execAndNotAssignUintptr,
	(int(OpLshAssign) << bitsKind) | int(Int):        execLshAssignInt,
	(int(OpLshAssign) << bitsKind) | int(Int8):       execLshAssignInt8,
	(int(OpLshAssign) << bitsKind) | int(Int16):      execLshAssignInt16,
	(int(OpLshAssign) << bitsKind) | int(Int32):      execLshAssignInt32,
	(int(OpLshAssign) << bitsKind) | int(Int64):      execLshAssignInt64,
	(int(OpLshAssign) << bitsKind) | int(Uint):       execLshAssignUint,
	(int(OpLshAssign) << bitsKind) | int(Uint8):      execLshAssignUint8,
	(int(OpLshAssign) << bitsKind) | int(Uint16):     execLshAssignUint16,
	(int(OpLshAssign) << bitsKind) | int(Uint32):     execLshAssignUint32,
	(int(OpLshAssign) << bitsKind) | int(Uint64):     execLshAssignUint64,
	(int(OpLshAssign) << bitsKind) | int(Uintptr):    execLshAssignUintptr,
	(int(OpRshAssign) << bitsKind) | int(Int):        execRshAssignInt,
	(int(OpRshAssign) << bitsKind) | int(Int8):       execRshAssignInt8,
	(int(OpRshAssign) << bitsKind) | int(Int16):      execRshAssignInt16,
	(int(OpRshAssign) << bitsKind) | int(Int32):      execRshAssignInt32,
	(int(OpRshAssign) << bitsKind) | int(Int64):      execRshAssignInt64,
	(int(OpRshAssign) << bitsKind) | int(Uint):       execRshAssignUint,
	(int(OpRshAssign) << bitsKind) | int(Uint8):      execRshAssignUint8,
	(int(OpRshAssign) << bitsKind) | int(Uint16):     execRshAssignUint16,
	(int(OpRshAssign) << bitsKind) | int(Uint32):     execRshAssignUint32,
	(int(OpRshAssign) << bitsKind) | int(Uint64):     execRshAssignUint64,
	(int(OpRshAssign) << bitsKind) | int(Uintptr):    execRshAssignUintptr,
	(int(OpInc) << bitsKind) | int(Int):              execIncInt,
	(int(OpInc) << bitsKind) | int(Int8):             execIncInt8,
	(int(OpInc) << bitsKind) | int(Int16):            execIncInt16,
	(int(OpInc) << bitsKind) | int(Int32):            execIncInt32,
	(int(OpInc) << bitsKind) | int(Int64):            execIncInt64,
	(int(OpInc) << bitsKind) | int(Uint):             execIncUint,
	(int(OpInc) << bitsKind) | int(Uint8):            execIncUint8,
	(int(OpInc) << bitsKind) | int(Uint16):           execIncUint16,
	(int(OpInc) << bitsKind) | int(Uint32):           execIncUint32,
	(int(OpInc) << bitsKind) | int(Uint64):           execIncUint64,
	(int(OpInc) << bitsKind) | int(Uintptr):          execIncUintptr,
	(int(OpInc) << bitsKind) | int(Float32):          execIncFloat32,
	(int(OpInc) << bitsKind) | int(Float64):          execIncFloat64,
	(int(OpInc) << bitsKind) | int(Complex64):        execIncComplex64,
	(int(OpInc) << bitsKind) | int(Complex128):       execIncComplex128,
	(int(OpDec) << bitsKind) | int(Int):              execDecInt,
	(int(OpDec) << bitsKind) | int(Int8):             execDecInt8,
	(int(OpDec) << bitsKind) | int(Int16):            execDecInt16,
	(int(OpDec) << bitsKind) | int(Int32):            execDecInt32,
	(int(OpDec) << bitsKind) | int(Int64):            execDecInt64,
	(int(OpDec) << bitsKind) | int(Uint):             execDecUint,
	(int(OpDec) << bitsKind) | int(Uint8):            execDecUint8,
	(int(OpDec) << bitsKind) | int(Uint16):           execDecUint16,
	(int(OpDec) << bitsKind) | int(Uint32):           execDecUint32,
	(int(OpDec) << bitsKind) | int(Uint64):           execDecUint64,
	(int(OpDec) << bitsKind) | int(Uintptr):          execDecUintptr,
	(int(OpDec) << bitsKind) | int(Float32):          execDecFloat32,
	(int(OpDec) << bitsKind) | int(Float64):          execDecFloat64,
	(int(OpDec) << bitsKind) | int(Complex64):        execDecComplex64,
	(int(OpDec) << bitsKind) | int(Complex128):       execDecComplex128,
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

func execQuoAssignInt(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int) /= p.data[n-2].(int)
	p.data = p.data[:n-2]
}

func execQuoAssignInt8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int8) /= p.data[n-2].(int8)
	p.data = p.data[:n-2]
}

func execQuoAssignInt16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int16) /= p.data[n-2].(int16)
	p.data = p.data[:n-2]
}

func execQuoAssignInt32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int32) /= p.data[n-2].(int32)
	p.data = p.data[:n-2]
}

func execQuoAssignInt64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int64) /= p.data[n-2].(int64)
	p.data = p.data[:n-2]
}

func execQuoAssignUint(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint) /= p.data[n-2].(uint)
	p.data = p.data[:n-2]
}

func execQuoAssignUint8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint8) /= p.data[n-2].(uint8)
	p.data = p.data[:n-2]
}

func execQuoAssignUint16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint16) /= p.data[n-2].(uint16)
	p.data = p.data[:n-2]
}

func execQuoAssignUint32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint32) /= p.data[n-2].(uint32)
	p.data = p.data[:n-2]
}

func execQuoAssignUint64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint64) /= p.data[n-2].(uint64)
	p.data = p.data[:n-2]
}

func execQuoAssignUintptr(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uintptr) /= p.data[n-2].(uintptr)
	p.data = p.data[:n-2]
}

func execQuoAssignFloat32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*float32) /= p.data[n-2].(float32)
	p.data = p.data[:n-2]
}

func execQuoAssignFloat64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*float64) /= p.data[n-2].(float64)
	p.data = p.data[:n-2]
}

func execQuoAssignComplex64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*complex64) /= p.data[n-2].(complex64)
	p.data = p.data[:n-2]
}

func execQuoAssignComplex128(i Instr, p *Context) {
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

func execAndAssignInt(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int) &= p.data[n-2].(int)
	p.data = p.data[:n-2]
}

func execAndAssignInt8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int8) &= p.data[n-2].(int8)
	p.data = p.data[:n-2]
}

func execAndAssignInt16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int16) &= p.data[n-2].(int16)
	p.data = p.data[:n-2]
}

func execAndAssignInt32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int32) &= p.data[n-2].(int32)
	p.data = p.data[:n-2]
}

func execAndAssignInt64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int64) &= p.data[n-2].(int64)
	p.data = p.data[:n-2]
}

func execAndAssignUint(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint) &= p.data[n-2].(uint)
	p.data = p.data[:n-2]
}

func execAndAssignUint8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint8) &= p.data[n-2].(uint8)
	p.data = p.data[:n-2]
}

func execAndAssignUint16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint16) &= p.data[n-2].(uint16)
	p.data = p.data[:n-2]
}

func execAndAssignUint32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint32) &= p.data[n-2].(uint32)
	p.data = p.data[:n-2]
}

func execAndAssignUint64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint64) &= p.data[n-2].(uint64)
	p.data = p.data[:n-2]
}

func execAndAssignUintptr(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uintptr) &= p.data[n-2].(uintptr)
	p.data = p.data[:n-2]
}

func execOrAssignInt(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int) |= p.data[n-2].(int)
	p.data = p.data[:n-2]
}

func execOrAssignInt8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int8) |= p.data[n-2].(int8)
	p.data = p.data[:n-2]
}

func execOrAssignInt16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int16) |= p.data[n-2].(int16)
	p.data = p.data[:n-2]
}

func execOrAssignInt32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int32) |= p.data[n-2].(int32)
	p.data = p.data[:n-2]
}

func execOrAssignInt64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int64) |= p.data[n-2].(int64)
	p.data = p.data[:n-2]
}

func execOrAssignUint(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint) |= p.data[n-2].(uint)
	p.data = p.data[:n-2]
}

func execOrAssignUint8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint8) |= p.data[n-2].(uint8)
	p.data = p.data[:n-2]
}

func execOrAssignUint16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint16) |= p.data[n-2].(uint16)
	p.data = p.data[:n-2]
}

func execOrAssignUint32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint32) |= p.data[n-2].(uint32)
	p.data = p.data[:n-2]
}

func execOrAssignUint64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint64) |= p.data[n-2].(uint64)
	p.data = p.data[:n-2]
}

func execOrAssignUintptr(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uintptr) |= p.data[n-2].(uintptr)
	p.data = p.data[:n-2]
}

func execXorAssignInt(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int) ^= p.data[n-2].(int)
	p.data = p.data[:n-2]
}

func execXorAssignInt8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int8) ^= p.data[n-2].(int8)
	p.data = p.data[:n-2]
}

func execXorAssignInt16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int16) ^= p.data[n-2].(int16)
	p.data = p.data[:n-2]
}

func execXorAssignInt32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int32) ^= p.data[n-2].(int32)
	p.data = p.data[:n-2]
}

func execXorAssignInt64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int64) ^= p.data[n-2].(int64)
	p.data = p.data[:n-2]
}

func execXorAssignUint(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint) ^= p.data[n-2].(uint)
	p.data = p.data[:n-2]
}

func execXorAssignUint8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint8) ^= p.data[n-2].(uint8)
	p.data = p.data[:n-2]
}

func execXorAssignUint16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint16) ^= p.data[n-2].(uint16)
	p.data = p.data[:n-2]
}

func execXorAssignUint32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint32) ^= p.data[n-2].(uint32)
	p.data = p.data[:n-2]
}

func execXorAssignUint64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint64) ^= p.data[n-2].(uint64)
	p.data = p.data[:n-2]
}

func execXorAssignUintptr(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uintptr) ^= p.data[n-2].(uintptr)
	p.data = p.data[:n-2]
}

func execAndNotAssignInt(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int) &^= p.data[n-2].(int)
	p.data = p.data[:n-2]
}

func execAndNotAssignInt8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int8) &^= p.data[n-2].(int8)
	p.data = p.data[:n-2]
}

func execAndNotAssignInt16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int16) &^= p.data[n-2].(int16)
	p.data = p.data[:n-2]
}

func execAndNotAssignInt32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int32) &^= p.data[n-2].(int32)
	p.data = p.data[:n-2]
}

func execAndNotAssignInt64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int64) &^= p.data[n-2].(int64)
	p.data = p.data[:n-2]
}

func execAndNotAssignUint(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint) &^= p.data[n-2].(uint)
	p.data = p.data[:n-2]
}

func execAndNotAssignUint8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint8) &^= p.data[n-2].(uint8)
	p.data = p.data[:n-2]
}

func execAndNotAssignUint16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint16) &^= p.data[n-2].(uint16)
	p.data = p.data[:n-2]
}

func execAndNotAssignUint32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint32) &^= p.data[n-2].(uint32)
	p.data = p.data[:n-2]
}

func execAndNotAssignUint64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint64) &^= p.data[n-2].(uint64)
	p.data = p.data[:n-2]
}

func execAndNotAssignUintptr(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uintptr) &^= p.data[n-2].(uintptr)
	p.data = p.data[:n-2]
}

func execLshAssignInt(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int) <<= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execLshAssignInt8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int8) <<= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execLshAssignInt16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int16) <<= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execLshAssignInt32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int32) <<= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execLshAssignInt64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int64) <<= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execLshAssignUint(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint) <<= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execLshAssignUint8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint8) <<= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execLshAssignUint16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint16) <<= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execLshAssignUint32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint32) <<= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execLshAssignUint64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint64) <<= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execLshAssignUintptr(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uintptr) <<= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execRshAssignInt(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int) >>= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execRshAssignInt8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int8) >>= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execRshAssignInt16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int16) >>= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execRshAssignInt32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int32) >>= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execRshAssignInt64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*int64) >>= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execRshAssignUint(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint) >>= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execRshAssignUint8(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint8) >>= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execRshAssignUint16(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint16) >>= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execRshAssignUint32(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint32) >>= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execRshAssignUint64(i Instr, p *Context) {
	n := len(p.data)
	*p.data[n-1].(*uint64) >>= toUint(p.data[n-2])
	p.data = p.data[:n-2]
}

func execRshAssignUintptr(i Instr, p *Context) {
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
