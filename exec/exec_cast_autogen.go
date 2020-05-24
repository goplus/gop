package exec

var builtinCastOps = [...]func(i Instr, p *Context){
	(int(Int) << bitsKind) | int(Int8):             execIntFromInt8,
	(int(Int) << bitsKind) | int(Int16):            execIntFromInt16,
	(int(Int) << bitsKind) | int(Int32):            execIntFromInt32,
	(int(Int) << bitsKind) | int(Int64):            execIntFromInt64,
	(int(Int) << bitsKind) | int(Uint):             execIntFromUint,
	(int(Int) << bitsKind) | int(Uint8):            execIntFromUint8,
	(int(Int) << bitsKind) | int(Uint16):           execIntFromUint16,
	(int(Int) << bitsKind) | int(Uint32):           execIntFromUint32,
	(int(Int) << bitsKind) | int(Uint64):           execIntFromUint64,
	(int(Int) << bitsKind) | int(Uintptr):          execIntFromUintptr,
	(int(Int) << bitsKind) | int(Float32):          execIntFromFloat32,
	(int(Int) << bitsKind) | int(Float64):          execIntFromFloat64,
	(int(Int8) << bitsKind) | int(Int):             execInt8FromInt,
	(int(Int8) << bitsKind) | int(Int16):           execInt8FromInt16,
	(int(Int8) << bitsKind) | int(Int32):           execInt8FromInt32,
	(int(Int8) << bitsKind) | int(Int64):           execInt8FromInt64,
	(int(Int8) << bitsKind) | int(Uint):            execInt8FromUint,
	(int(Int8) << bitsKind) | int(Uint8):           execInt8FromUint8,
	(int(Int8) << bitsKind) | int(Uint16):          execInt8FromUint16,
	(int(Int8) << bitsKind) | int(Uint32):          execInt8FromUint32,
	(int(Int8) << bitsKind) | int(Uint64):          execInt8FromUint64,
	(int(Int8) << bitsKind) | int(Uintptr):         execInt8FromUintptr,
	(int(Int8) << bitsKind) | int(Float32):         execInt8FromFloat32,
	(int(Int8) << bitsKind) | int(Float64):         execInt8FromFloat64,
	(int(Int16) << bitsKind) | int(Int):            execInt16FromInt,
	(int(Int16) << bitsKind) | int(Int8):           execInt16FromInt8,
	(int(Int16) << bitsKind) | int(Int32):          execInt16FromInt32,
	(int(Int16) << bitsKind) | int(Int64):          execInt16FromInt64,
	(int(Int16) << bitsKind) | int(Uint):           execInt16FromUint,
	(int(Int16) << bitsKind) | int(Uint8):          execInt16FromUint8,
	(int(Int16) << bitsKind) | int(Uint16):         execInt16FromUint16,
	(int(Int16) << bitsKind) | int(Uint32):         execInt16FromUint32,
	(int(Int16) << bitsKind) | int(Uint64):         execInt16FromUint64,
	(int(Int16) << bitsKind) | int(Uintptr):        execInt16FromUintptr,
	(int(Int16) << bitsKind) | int(Float32):        execInt16FromFloat32,
	(int(Int16) << bitsKind) | int(Float64):        execInt16FromFloat64,
	(int(Int32) << bitsKind) | int(Int):            execInt32FromInt,
	(int(Int32) << bitsKind) | int(Int8):           execInt32FromInt8,
	(int(Int32) << bitsKind) | int(Int16):          execInt32FromInt16,
	(int(Int32) << bitsKind) | int(Int64):          execInt32FromInt64,
	(int(Int32) << bitsKind) | int(Uint):           execInt32FromUint,
	(int(Int32) << bitsKind) | int(Uint8):          execInt32FromUint8,
	(int(Int32) << bitsKind) | int(Uint16):         execInt32FromUint16,
	(int(Int32) << bitsKind) | int(Uint32):         execInt32FromUint32,
	(int(Int32) << bitsKind) | int(Uint64):         execInt32FromUint64,
	(int(Int32) << bitsKind) | int(Uintptr):        execInt32FromUintptr,
	(int(Int32) << bitsKind) | int(Float32):        execInt32FromFloat32,
	(int(Int32) << bitsKind) | int(Float64):        execInt32FromFloat64,
	(int(Int64) << bitsKind) | int(Int):            execInt64FromInt,
	(int(Int64) << bitsKind) | int(Int8):           execInt64FromInt8,
	(int(Int64) << bitsKind) | int(Int16):          execInt64FromInt16,
	(int(Int64) << bitsKind) | int(Int32):          execInt64FromInt32,
	(int(Int64) << bitsKind) | int(Uint):           execInt64FromUint,
	(int(Int64) << bitsKind) | int(Uint8):          execInt64FromUint8,
	(int(Int64) << bitsKind) | int(Uint16):         execInt64FromUint16,
	(int(Int64) << bitsKind) | int(Uint32):         execInt64FromUint32,
	(int(Int64) << bitsKind) | int(Uint64):         execInt64FromUint64,
	(int(Int64) << bitsKind) | int(Uintptr):        execInt64FromUintptr,
	(int(Int64) << bitsKind) | int(Float32):        execInt64FromFloat32,
	(int(Int64) << bitsKind) | int(Float64):        execInt64FromFloat64,
	(int(Uint) << bitsKind) | int(Int):             execUintFromInt,
	(int(Uint) << bitsKind) | int(Int8):            execUintFromInt8,
	(int(Uint) << bitsKind) | int(Int16):           execUintFromInt16,
	(int(Uint) << bitsKind) | int(Int32):           execUintFromInt32,
	(int(Uint) << bitsKind) | int(Int64):           execUintFromInt64,
	(int(Uint) << bitsKind) | int(Uint8):           execUintFromUint8,
	(int(Uint) << bitsKind) | int(Uint16):          execUintFromUint16,
	(int(Uint) << bitsKind) | int(Uint32):          execUintFromUint32,
	(int(Uint) << bitsKind) | int(Uint64):          execUintFromUint64,
	(int(Uint) << bitsKind) | int(Uintptr):         execUintFromUintptr,
	(int(Uint) << bitsKind) | int(Float32):         execUintFromFloat32,
	(int(Uint) << bitsKind) | int(Float64):         execUintFromFloat64,
	(int(Uint8) << bitsKind) | int(Int):            execUint8FromInt,
	(int(Uint8) << bitsKind) | int(Int8):           execUint8FromInt8,
	(int(Uint8) << bitsKind) | int(Int16):          execUint8FromInt16,
	(int(Uint8) << bitsKind) | int(Int32):          execUint8FromInt32,
	(int(Uint8) << bitsKind) | int(Int64):          execUint8FromInt64,
	(int(Uint8) << bitsKind) | int(Uint):           execUint8FromUint,
	(int(Uint8) << bitsKind) | int(Uint16):         execUint8FromUint16,
	(int(Uint8) << bitsKind) | int(Uint32):         execUint8FromUint32,
	(int(Uint8) << bitsKind) | int(Uint64):         execUint8FromUint64,
	(int(Uint8) << bitsKind) | int(Uintptr):        execUint8FromUintptr,
	(int(Uint8) << bitsKind) | int(Float32):        execUint8FromFloat32,
	(int(Uint8) << bitsKind) | int(Float64):        execUint8FromFloat64,
	(int(Uint16) << bitsKind) | int(Int):           execUint16FromInt,
	(int(Uint16) << bitsKind) | int(Int8):          execUint16FromInt8,
	(int(Uint16) << bitsKind) | int(Int16):         execUint16FromInt16,
	(int(Uint16) << bitsKind) | int(Int32):         execUint16FromInt32,
	(int(Uint16) << bitsKind) | int(Int64):         execUint16FromInt64,
	(int(Uint16) << bitsKind) | int(Uint):          execUint16FromUint,
	(int(Uint16) << bitsKind) | int(Uint8):         execUint16FromUint8,
	(int(Uint16) << bitsKind) | int(Uint32):        execUint16FromUint32,
	(int(Uint16) << bitsKind) | int(Uint64):        execUint16FromUint64,
	(int(Uint16) << bitsKind) | int(Uintptr):       execUint16FromUintptr,
	(int(Uint16) << bitsKind) | int(Float32):       execUint16FromFloat32,
	(int(Uint16) << bitsKind) | int(Float64):       execUint16FromFloat64,
	(int(Uint32) << bitsKind) | int(Int):           execUint32FromInt,
	(int(Uint32) << bitsKind) | int(Int8):          execUint32FromInt8,
	(int(Uint32) << bitsKind) | int(Int16):         execUint32FromInt16,
	(int(Uint32) << bitsKind) | int(Int32):         execUint32FromInt32,
	(int(Uint32) << bitsKind) | int(Int64):         execUint32FromInt64,
	(int(Uint32) << bitsKind) | int(Uint):          execUint32FromUint,
	(int(Uint32) << bitsKind) | int(Uint8):         execUint32FromUint8,
	(int(Uint32) << bitsKind) | int(Uint16):        execUint32FromUint16,
	(int(Uint32) << bitsKind) | int(Uint64):        execUint32FromUint64,
	(int(Uint32) << bitsKind) | int(Uintptr):       execUint32FromUintptr,
	(int(Uint32) << bitsKind) | int(Float32):       execUint32FromFloat32,
	(int(Uint32) << bitsKind) | int(Float64):       execUint32FromFloat64,
	(int(Uint64) << bitsKind) | int(Int):           execUint64FromInt,
	(int(Uint64) << bitsKind) | int(Int8):          execUint64FromInt8,
	(int(Uint64) << bitsKind) | int(Int16):         execUint64FromInt16,
	(int(Uint64) << bitsKind) | int(Int32):         execUint64FromInt32,
	(int(Uint64) << bitsKind) | int(Int64):         execUint64FromInt64,
	(int(Uint64) << bitsKind) | int(Uint):          execUint64FromUint,
	(int(Uint64) << bitsKind) | int(Uint8):         execUint64FromUint8,
	(int(Uint64) << bitsKind) | int(Uint16):        execUint64FromUint16,
	(int(Uint64) << bitsKind) | int(Uint32):        execUint64FromUint32,
	(int(Uint64) << bitsKind) | int(Uintptr):       execUint64FromUintptr,
	(int(Uint64) << bitsKind) | int(Float32):       execUint64FromFloat32,
	(int(Uint64) << bitsKind) | int(Float64):       execUint64FromFloat64,
	(int(Uintptr) << bitsKind) | int(Int):          execUintptrFromInt,
	(int(Uintptr) << bitsKind) | int(Int8):         execUintptrFromInt8,
	(int(Uintptr) << bitsKind) | int(Int16):        execUintptrFromInt16,
	(int(Uintptr) << bitsKind) | int(Int32):        execUintptrFromInt32,
	(int(Uintptr) << bitsKind) | int(Int64):        execUintptrFromInt64,
	(int(Uintptr) << bitsKind) | int(Uint):         execUintptrFromUint,
	(int(Uintptr) << bitsKind) | int(Uint8):        execUintptrFromUint8,
	(int(Uintptr) << bitsKind) | int(Uint16):       execUintptrFromUint16,
	(int(Uintptr) << bitsKind) | int(Uint32):       execUintptrFromUint32,
	(int(Uintptr) << bitsKind) | int(Uint64):       execUintptrFromUint64,
	(int(Uintptr) << bitsKind) | int(Float32):      execUintptrFromFloat32,
	(int(Uintptr) << bitsKind) | int(Float64):      execUintptrFromFloat64,
	(int(Float32) << bitsKind) | int(Int):          execFloat32FromInt,
	(int(Float32) << bitsKind) | int(Int8):         execFloat32FromInt8,
	(int(Float32) << bitsKind) | int(Int16):        execFloat32FromInt16,
	(int(Float32) << bitsKind) | int(Int32):        execFloat32FromInt32,
	(int(Float32) << bitsKind) | int(Int64):        execFloat32FromInt64,
	(int(Float32) << bitsKind) | int(Uint):         execFloat32FromUint,
	(int(Float32) << bitsKind) | int(Uint8):        execFloat32FromUint8,
	(int(Float32) << bitsKind) | int(Uint16):       execFloat32FromUint16,
	(int(Float32) << bitsKind) | int(Uint32):       execFloat32FromUint32,
	(int(Float32) << bitsKind) | int(Uint64):       execFloat32FromUint64,
	(int(Float32) << bitsKind) | int(Uintptr):      execFloat32FromUintptr,
	(int(Float32) << bitsKind) | int(Float64):      execFloat32FromFloat64,
	(int(Float64) << bitsKind) | int(Int):          execFloat64FromInt,
	(int(Float64) << bitsKind) | int(Int8):         execFloat64FromInt8,
	(int(Float64) << bitsKind) | int(Int16):        execFloat64FromInt16,
	(int(Float64) << bitsKind) | int(Int32):        execFloat64FromInt32,
	(int(Float64) << bitsKind) | int(Int64):        execFloat64FromInt64,
	(int(Float64) << bitsKind) | int(Uint):         execFloat64FromUint,
	(int(Float64) << bitsKind) | int(Uint8):        execFloat64FromUint8,
	(int(Float64) << bitsKind) | int(Uint16):       execFloat64FromUint16,
	(int(Float64) << bitsKind) | int(Uint32):       execFloat64FromUint32,
	(int(Float64) << bitsKind) | int(Uint64):       execFloat64FromUint64,
	(int(Float64) << bitsKind) | int(Uintptr):      execFloat64FromUintptr,
	(int(Float64) << bitsKind) | int(Float32):      execFloat64FromFloat32,
	(int(Complex64) << bitsKind) | int(Complex128): execComplex64FromComplex128,
	(int(Complex128) << bitsKind) | int(Complex64): execComplex128FromComplex64,
	(int(String) << bitsKind) | int(Int):           execStringFromInt,
	(int(String) << bitsKind) | int(Int8):          execStringFromInt8,
	(int(String) << bitsKind) | int(Int16):         execStringFromInt16,
	(int(String) << bitsKind) | int(Int32):         execStringFromInt32,
	(int(String) << bitsKind) | int(Int64):         execStringFromInt64,
	(int(String) << bitsKind) | int(Uint):          execStringFromUint,
	(int(String) << bitsKind) | int(Uint8):         execStringFromUint8,
	(int(String) << bitsKind) | int(Uint16):        execStringFromUint16,
	(int(String) << bitsKind) | int(Uint32):        execStringFromUint32,
	(int(String) << bitsKind) | int(Uint64):        execStringFromUint64,
	(int(String) << bitsKind) | int(Uintptr):       execStringFromUintptr,
}

func execIntFromInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int(p.data[n-1].(int8))
}

func execIntFromInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int(p.data[n-1].(int16))
}

func execIntFromInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int(p.data[n-1].(int32))
}

func execIntFromInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int(p.data[n-1].(int64))
}

func execIntFromUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int(p.data[n-1].(uint))
}

func execIntFromUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int(p.data[n-1].(uint8))
}

func execIntFromUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int(p.data[n-1].(uint16))
}

func execIntFromUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int(p.data[n-1].(uint32))
}

func execIntFromUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int(p.data[n-1].(uint64))
}

func execIntFromUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int(p.data[n-1].(uintptr))
}

func execIntFromFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int(p.data[n-1].(float32))
}

func execIntFromFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int(p.data[n-1].(float64))
}

func execInt8FromInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int8(p.data[n-1].(int))
}

func execInt8FromInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int8(p.data[n-1].(int16))
}

func execInt8FromInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int8(p.data[n-1].(int32))
}

func execInt8FromInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int8(p.data[n-1].(int64))
}

func execInt8FromUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int8(p.data[n-1].(uint))
}

func execInt8FromUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int8(p.data[n-1].(uint8))
}

func execInt8FromUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int8(p.data[n-1].(uint16))
}

func execInt8FromUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int8(p.data[n-1].(uint32))
}

func execInt8FromUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int8(p.data[n-1].(uint64))
}

func execInt8FromUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int8(p.data[n-1].(uintptr))
}

func execInt8FromFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int8(p.data[n-1].(float32))
}

func execInt8FromFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int8(p.data[n-1].(float64))
}

func execInt16FromInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int16(p.data[n-1].(int))
}

func execInt16FromInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int16(p.data[n-1].(int8))
}

func execInt16FromInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int16(p.data[n-1].(int32))
}

func execInt16FromInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int16(p.data[n-1].(int64))
}

func execInt16FromUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int16(p.data[n-1].(uint))
}

func execInt16FromUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int16(p.data[n-1].(uint8))
}

func execInt16FromUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int16(p.data[n-1].(uint16))
}

func execInt16FromUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int16(p.data[n-1].(uint32))
}

func execInt16FromUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int16(p.data[n-1].(uint64))
}

func execInt16FromUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int16(p.data[n-1].(uintptr))
}

func execInt16FromFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int16(p.data[n-1].(float32))
}

func execInt16FromFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int16(p.data[n-1].(float64))
}

func execInt32FromInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int32(p.data[n-1].(int))
}

func execInt32FromInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int32(p.data[n-1].(int8))
}

func execInt32FromInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int32(p.data[n-1].(int16))
}

func execInt32FromInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int32(p.data[n-1].(int64))
}

func execInt32FromUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int32(p.data[n-1].(uint))
}

func execInt32FromUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int32(p.data[n-1].(uint8))
}

func execInt32FromUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int32(p.data[n-1].(uint16))
}

func execInt32FromUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int32(p.data[n-1].(uint32))
}

func execInt32FromUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int32(p.data[n-1].(uint64))
}

func execInt32FromUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int32(p.data[n-1].(uintptr))
}

func execInt32FromFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int32(p.data[n-1].(float32))
}

func execInt32FromFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int32(p.data[n-1].(float64))
}

func execInt64FromInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int64(p.data[n-1].(int))
}

func execInt64FromInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int64(p.data[n-1].(int8))
}

func execInt64FromInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int64(p.data[n-1].(int16))
}

func execInt64FromInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int64(p.data[n-1].(int32))
}

func execInt64FromUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int64(p.data[n-1].(uint))
}

func execInt64FromUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int64(p.data[n-1].(uint8))
}

func execInt64FromUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int64(p.data[n-1].(uint16))
}

func execInt64FromUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int64(p.data[n-1].(uint32))
}

func execInt64FromUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int64(p.data[n-1].(uint64))
}

func execInt64FromUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int64(p.data[n-1].(uintptr))
}

func execInt64FromFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int64(p.data[n-1].(float32))
}

func execInt64FromFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = int64(p.data[n-1].(float64))
}

func execUintFromInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint(p.data[n-1].(int))
}

func execUintFromInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint(p.data[n-1].(int8))
}

func execUintFromInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint(p.data[n-1].(int16))
}

func execUintFromInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint(p.data[n-1].(int32))
}

func execUintFromInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint(p.data[n-1].(int64))
}

func execUintFromUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint(p.data[n-1].(uint8))
}

func execUintFromUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint(p.data[n-1].(uint16))
}

func execUintFromUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint(p.data[n-1].(uint32))
}

func execUintFromUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint(p.data[n-1].(uint64))
}

func execUintFromUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint(p.data[n-1].(uintptr))
}

func execUintFromFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint(p.data[n-1].(float32))
}

func execUintFromFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint(p.data[n-1].(float64))
}

func execUint8FromInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint8(p.data[n-1].(int))
}

func execUint8FromInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint8(p.data[n-1].(int8))
}

func execUint8FromInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint8(p.data[n-1].(int16))
}

func execUint8FromInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint8(p.data[n-1].(int32))
}

func execUint8FromInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint8(p.data[n-1].(int64))
}

func execUint8FromUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint8(p.data[n-1].(uint))
}

func execUint8FromUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint8(p.data[n-1].(uint16))
}

func execUint8FromUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint8(p.data[n-1].(uint32))
}

func execUint8FromUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint8(p.data[n-1].(uint64))
}

func execUint8FromUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint8(p.data[n-1].(uintptr))
}

func execUint8FromFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint8(p.data[n-1].(float32))
}

func execUint8FromFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint8(p.data[n-1].(float64))
}

func execUint16FromInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint16(p.data[n-1].(int))
}

func execUint16FromInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint16(p.data[n-1].(int8))
}

func execUint16FromInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint16(p.data[n-1].(int16))
}

func execUint16FromInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint16(p.data[n-1].(int32))
}

func execUint16FromInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint16(p.data[n-1].(int64))
}

func execUint16FromUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint16(p.data[n-1].(uint))
}

func execUint16FromUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint16(p.data[n-1].(uint8))
}

func execUint16FromUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint16(p.data[n-1].(uint32))
}

func execUint16FromUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint16(p.data[n-1].(uint64))
}

func execUint16FromUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint16(p.data[n-1].(uintptr))
}

func execUint16FromFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint16(p.data[n-1].(float32))
}

func execUint16FromFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint16(p.data[n-1].(float64))
}

func execUint32FromInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint32(p.data[n-1].(int))
}

func execUint32FromInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint32(p.data[n-1].(int8))
}

func execUint32FromInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint32(p.data[n-1].(int16))
}

func execUint32FromInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint32(p.data[n-1].(int32))
}

func execUint32FromInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint32(p.data[n-1].(int64))
}

func execUint32FromUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint32(p.data[n-1].(uint))
}

func execUint32FromUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint32(p.data[n-1].(uint8))
}

func execUint32FromUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint32(p.data[n-1].(uint16))
}

func execUint32FromUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint32(p.data[n-1].(uint64))
}

func execUint32FromUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint32(p.data[n-1].(uintptr))
}

func execUint32FromFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint32(p.data[n-1].(float32))
}

func execUint32FromFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint32(p.data[n-1].(float64))
}

func execUint64FromInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint64(p.data[n-1].(int))
}

func execUint64FromInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint64(p.data[n-1].(int8))
}

func execUint64FromInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint64(p.data[n-1].(int16))
}

func execUint64FromInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint64(p.data[n-1].(int32))
}

func execUint64FromInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint64(p.data[n-1].(int64))
}

func execUint64FromUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint64(p.data[n-1].(uint))
}

func execUint64FromUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint64(p.data[n-1].(uint8))
}

func execUint64FromUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint64(p.data[n-1].(uint16))
}

func execUint64FromUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint64(p.data[n-1].(uint32))
}

func execUint64FromUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint64(p.data[n-1].(uintptr))
}

func execUint64FromFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint64(p.data[n-1].(float32))
}

func execUint64FromFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uint64(p.data[n-1].(float64))
}

func execUintptrFromInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uintptr(p.data[n-1].(int))
}

func execUintptrFromInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uintptr(p.data[n-1].(int8))
}

func execUintptrFromInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uintptr(p.data[n-1].(int16))
}

func execUintptrFromInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uintptr(p.data[n-1].(int32))
}

func execUintptrFromInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uintptr(p.data[n-1].(int64))
}

func execUintptrFromUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uintptr(p.data[n-1].(uint))
}

func execUintptrFromUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uintptr(p.data[n-1].(uint8))
}

func execUintptrFromUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uintptr(p.data[n-1].(uint16))
}

func execUintptrFromUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uintptr(p.data[n-1].(uint32))
}

func execUintptrFromUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uintptr(p.data[n-1].(uint64))
}

func execUintptrFromFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uintptr(p.data[n-1].(float32))
}

func execUintptrFromFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = uintptr(p.data[n-1].(float64))
}

func execFloat32FromInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = float32(p.data[n-1].(int))
}

func execFloat32FromInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = float32(p.data[n-1].(int8))
}

func execFloat32FromInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = float32(p.data[n-1].(int16))
}

func execFloat32FromInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = float32(p.data[n-1].(int32))
}

func execFloat32FromInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = float32(p.data[n-1].(int64))
}

func execFloat32FromUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = float32(p.data[n-1].(uint))
}

func execFloat32FromUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = float32(p.data[n-1].(uint8))
}

func execFloat32FromUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = float32(p.data[n-1].(uint16))
}

func execFloat32FromUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = float32(p.data[n-1].(uint32))
}

func execFloat32FromUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = float32(p.data[n-1].(uint64))
}

func execFloat32FromUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = float32(p.data[n-1].(uintptr))
}

func execFloat32FromFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = float32(p.data[n-1].(float64))
}

func execFloat64FromInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = float64(p.data[n-1].(int))
}

func execFloat64FromInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = float64(p.data[n-1].(int8))
}

func execFloat64FromInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = float64(p.data[n-1].(int16))
}

func execFloat64FromInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = float64(p.data[n-1].(int32))
}

func execFloat64FromInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = float64(p.data[n-1].(int64))
}

func execFloat64FromUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = float64(p.data[n-1].(uint))
}

func execFloat64FromUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = float64(p.data[n-1].(uint8))
}

func execFloat64FromUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = float64(p.data[n-1].(uint16))
}

func execFloat64FromUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = float64(p.data[n-1].(uint32))
}

func execFloat64FromUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = float64(p.data[n-1].(uint64))
}

func execFloat64FromUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = float64(p.data[n-1].(uintptr))
}

func execFloat64FromFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = float64(p.data[n-1].(float32))
}

func execComplex64FromComplex128(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = complex64(p.data[n-1].(complex128))
}

func execComplex128FromComplex64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = complex128(p.data[n-1].(complex64))
}

func execStringFromInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = string(p.data[n-1].(int))
}

func execStringFromInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = string(p.data[n-1].(int8))
}

func execStringFromInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = string(p.data[n-1].(int16))
}

func execStringFromInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = string(p.data[n-1].(int32))
}

func execStringFromInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = string(p.data[n-1].(int64))
}

func execStringFromUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = string(p.data[n-1].(uint))
}

func execStringFromUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = string(p.data[n-1].(uint8))
}

func execStringFromUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = string(p.data[n-1].(uint16))
}

func execStringFromUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = string(p.data[n-1].(uint32))
}

func execStringFromUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = string(p.data[n-1].(uint64))
}

func execStringFromUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = string(p.data[n-1].(uintptr))
}
