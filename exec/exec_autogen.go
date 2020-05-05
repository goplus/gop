
package exec

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
