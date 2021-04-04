package bytecode

import (
	"math/big"
	"reflect"
)

var builtinOps = [...]func(i Instr, p *Context){
	(int(BigInt) << bitsOperator) | int(OpQuo):           execQuoBigInt,
	(int(Int) << bitsOperator) | int(OpAdd):              execAddInt,
	(int(Int8) << bitsOperator) | int(OpAdd):             execAddInt8,
	(int(Int16) << bitsOperator) | int(OpAdd):            execAddInt16,
	(int(Int32) << bitsOperator) | int(OpAdd):            execAddInt32,
	(int(Int64) << bitsOperator) | int(OpAdd):            execAddInt64,
	(int(Uint) << bitsOperator) | int(OpAdd):             execAddUint,
	(int(Uint8) << bitsOperator) | int(OpAdd):            execAddUint8,
	(int(Uint16) << bitsOperator) | int(OpAdd):           execAddUint16,
	(int(Uint32) << bitsOperator) | int(OpAdd):           execAddUint32,
	(int(Uint64) << bitsOperator) | int(OpAdd):           execAddUint64,
	(int(Uintptr) << bitsOperator) | int(OpAdd):          execAddUintptr,
	(int(Float32) << bitsOperator) | int(OpAdd):          execAddFloat32,
	(int(Float64) << bitsOperator) | int(OpAdd):          execAddFloat64,
	(int(Complex64) << bitsOperator) | int(OpAdd):        execAddComplex64,
	(int(Complex128) << bitsOperator) | int(OpAdd):       execAddComplex128,
	(int(String) << bitsOperator) | int(OpAdd):           execAddString,
	(int(BigInt) << bitsOperator) | int(OpAdd):           execAddBigInt,
	(int(BigRat) << bitsOperator) | int(OpAdd):           execAddBigRat,
	(int(BigFloat) << bitsOperator) | int(OpAdd):         execAddBigFloat,
	(int(Int) << bitsOperator) | int(OpSub):              execSubInt,
	(int(Int8) << bitsOperator) | int(OpSub):             execSubInt8,
	(int(Int16) << bitsOperator) | int(OpSub):            execSubInt16,
	(int(Int32) << bitsOperator) | int(OpSub):            execSubInt32,
	(int(Int64) << bitsOperator) | int(OpSub):            execSubInt64,
	(int(Uint) << bitsOperator) | int(OpSub):             execSubUint,
	(int(Uint8) << bitsOperator) | int(OpSub):            execSubUint8,
	(int(Uint16) << bitsOperator) | int(OpSub):           execSubUint16,
	(int(Uint32) << bitsOperator) | int(OpSub):           execSubUint32,
	(int(Uint64) << bitsOperator) | int(OpSub):           execSubUint64,
	(int(Uintptr) << bitsOperator) | int(OpSub):          execSubUintptr,
	(int(Float32) << bitsOperator) | int(OpSub):          execSubFloat32,
	(int(Float64) << bitsOperator) | int(OpSub):          execSubFloat64,
	(int(Complex64) << bitsOperator) | int(OpSub):        execSubComplex64,
	(int(Complex128) << bitsOperator) | int(OpSub):       execSubComplex128,
	(int(BigInt) << bitsOperator) | int(OpSub):           execSubBigInt,
	(int(BigRat) << bitsOperator) | int(OpSub):           execSubBigRat,
	(int(BigFloat) << bitsOperator) | int(OpSub):         execSubBigFloat,
	(int(Int) << bitsOperator) | int(OpMul):              execMulInt,
	(int(Int8) << bitsOperator) | int(OpMul):             execMulInt8,
	(int(Int16) << bitsOperator) | int(OpMul):            execMulInt16,
	(int(Int32) << bitsOperator) | int(OpMul):            execMulInt32,
	(int(Int64) << bitsOperator) | int(OpMul):            execMulInt64,
	(int(Uint) << bitsOperator) | int(OpMul):             execMulUint,
	(int(Uint8) << bitsOperator) | int(OpMul):            execMulUint8,
	(int(Uint16) << bitsOperator) | int(OpMul):           execMulUint16,
	(int(Uint32) << bitsOperator) | int(OpMul):           execMulUint32,
	(int(Uint64) << bitsOperator) | int(OpMul):           execMulUint64,
	(int(Uintptr) << bitsOperator) | int(OpMul):          execMulUintptr,
	(int(Float32) << bitsOperator) | int(OpMul):          execMulFloat32,
	(int(Float64) << bitsOperator) | int(OpMul):          execMulFloat64,
	(int(Complex64) << bitsOperator) | int(OpMul):        execMulComplex64,
	(int(Complex128) << bitsOperator) | int(OpMul):       execMulComplex128,
	(int(BigInt) << bitsOperator) | int(OpMul):           execMulBigInt,
	(int(BigRat) << bitsOperator) | int(OpMul):           execMulBigRat,
	(int(BigFloat) << bitsOperator) | int(OpMul):         execMulBigFloat,
	(int(Int) << bitsOperator) | int(OpQuo):              execQuoInt,
	(int(Int8) << bitsOperator) | int(OpQuo):             execQuoInt8,
	(int(Int16) << bitsOperator) | int(OpQuo):            execQuoInt16,
	(int(Int32) << bitsOperator) | int(OpQuo):            execQuoInt32,
	(int(Int64) << bitsOperator) | int(OpQuo):            execQuoInt64,
	(int(Uint) << bitsOperator) | int(OpQuo):             execQuoUint,
	(int(Uint8) << bitsOperator) | int(OpQuo):            execQuoUint8,
	(int(Uint16) << bitsOperator) | int(OpQuo):           execQuoUint16,
	(int(Uint32) << bitsOperator) | int(OpQuo):           execQuoUint32,
	(int(Uint64) << bitsOperator) | int(OpQuo):           execQuoUint64,
	(int(Uintptr) << bitsOperator) | int(OpQuo):          execQuoUintptr,
	(int(Float32) << bitsOperator) | int(OpQuo):          execQuoFloat32,
	(int(Float64) << bitsOperator) | int(OpQuo):          execQuoFloat64,
	(int(Complex64) << bitsOperator) | int(OpQuo):        execQuoComplex64,
	(int(Complex128) << bitsOperator) | int(OpQuo):       execQuoComplex128,
	(int(BigRat) << bitsOperator) | int(OpQuo):           execQuoBigRat,
	(int(BigFloat) << bitsOperator) | int(OpQuo):         execQuoBigFloat,
	(int(Int) << bitsOperator) | int(OpMod):              execModInt,
	(int(Int8) << bitsOperator) | int(OpMod):             execModInt8,
	(int(Int16) << bitsOperator) | int(OpMod):            execModInt16,
	(int(Int32) << bitsOperator) | int(OpMod):            execModInt32,
	(int(Int64) << bitsOperator) | int(OpMod):            execModInt64,
	(int(Uint) << bitsOperator) | int(OpMod):             execModUint,
	(int(Uint8) << bitsOperator) | int(OpMod):            execModUint8,
	(int(Uint16) << bitsOperator) | int(OpMod):           execModUint16,
	(int(Uint32) << bitsOperator) | int(OpMod):           execModUint32,
	(int(Uint64) << bitsOperator) | int(OpMod):           execModUint64,
	(int(Uintptr) << bitsOperator) | int(OpMod):          execModUintptr,
	(int(BigInt) << bitsOperator) | int(OpMod):           execModBigInt,
	(int(Int) << bitsOperator) | int(OpAnd):              execAndInt,
	(int(Int8) << bitsOperator) | int(OpAnd):             execAndInt8,
	(int(Int16) << bitsOperator) | int(OpAnd):            execAndInt16,
	(int(Int32) << bitsOperator) | int(OpAnd):            execAndInt32,
	(int(Int64) << bitsOperator) | int(OpAnd):            execAndInt64,
	(int(Uint) << bitsOperator) | int(OpAnd):             execAndUint,
	(int(Uint8) << bitsOperator) | int(OpAnd):            execAndUint8,
	(int(Uint16) << bitsOperator) | int(OpAnd):           execAndUint16,
	(int(Uint32) << bitsOperator) | int(OpAnd):           execAndUint32,
	(int(Uint64) << bitsOperator) | int(OpAnd):           execAndUint64,
	(int(Uintptr) << bitsOperator) | int(OpAnd):          execAndUintptr,
	(int(BigInt) << bitsOperator) | int(OpAnd):           execAndBigInt,
	(int(Int) << bitsOperator) | int(OpOr):               execOrInt,
	(int(Int8) << bitsOperator) | int(OpOr):              execOrInt8,
	(int(Int16) << bitsOperator) | int(OpOr):             execOrInt16,
	(int(Int32) << bitsOperator) | int(OpOr):             execOrInt32,
	(int(Int64) << bitsOperator) | int(OpOr):             execOrInt64,
	(int(Uint) << bitsOperator) | int(OpOr):              execOrUint,
	(int(Uint8) << bitsOperator) | int(OpOr):             execOrUint8,
	(int(Uint16) << bitsOperator) | int(OpOr):            execOrUint16,
	(int(Uint32) << bitsOperator) | int(OpOr):            execOrUint32,
	(int(Uint64) << bitsOperator) | int(OpOr):            execOrUint64,
	(int(Uintptr) << bitsOperator) | int(OpOr):           execOrUintptr,
	(int(BigInt) << bitsOperator) | int(OpOr):            execOrBigInt,
	(int(Int) << bitsOperator) | int(OpXor):              execXorInt,
	(int(Int8) << bitsOperator) | int(OpXor):             execXorInt8,
	(int(Int16) << bitsOperator) | int(OpXor):            execXorInt16,
	(int(Int32) << bitsOperator) | int(OpXor):            execXorInt32,
	(int(Int64) << bitsOperator) | int(OpXor):            execXorInt64,
	(int(Uint) << bitsOperator) | int(OpXor):             execXorUint,
	(int(Uint8) << bitsOperator) | int(OpXor):            execXorUint8,
	(int(Uint16) << bitsOperator) | int(OpXor):           execXorUint16,
	(int(Uint32) << bitsOperator) | int(OpXor):           execXorUint32,
	(int(Uint64) << bitsOperator) | int(OpXor):           execXorUint64,
	(int(Uintptr) << bitsOperator) | int(OpXor):          execXorUintptr,
	(int(BigInt) << bitsOperator) | int(OpXor):           execXorBigInt,
	(int(Int) << bitsOperator) | int(OpAndNot):           execAndNotInt,
	(int(Int8) << bitsOperator) | int(OpAndNot):          execAndNotInt8,
	(int(Int16) << bitsOperator) | int(OpAndNot):         execAndNotInt16,
	(int(Int32) << bitsOperator) | int(OpAndNot):         execAndNotInt32,
	(int(Int64) << bitsOperator) | int(OpAndNot):         execAndNotInt64,
	(int(Uint) << bitsOperator) | int(OpAndNot):          execAndNotUint,
	(int(Uint8) << bitsOperator) | int(OpAndNot):         execAndNotUint8,
	(int(Uint16) << bitsOperator) | int(OpAndNot):        execAndNotUint16,
	(int(Uint32) << bitsOperator) | int(OpAndNot):        execAndNotUint32,
	(int(Uint64) << bitsOperator) | int(OpAndNot):        execAndNotUint64,
	(int(Uintptr) << bitsOperator) | int(OpAndNot):       execAndNotUintptr,
	(int(BigInt) << bitsOperator) | int(OpAndNot):        execAndNotBigInt,
	(int(Int) << bitsOperator) | int(OpLsh):              execLshInt,
	(int(Int8) << bitsOperator) | int(OpLsh):             execLshInt8,
	(int(Int16) << bitsOperator) | int(OpLsh):            execLshInt16,
	(int(Int32) << bitsOperator) | int(OpLsh):            execLshInt32,
	(int(Int64) << bitsOperator) | int(OpLsh):            execLshInt64,
	(int(Uint) << bitsOperator) | int(OpLsh):             execLshUint,
	(int(Uint8) << bitsOperator) | int(OpLsh):            execLshUint8,
	(int(Uint16) << bitsOperator) | int(OpLsh):           execLshUint16,
	(int(Uint32) << bitsOperator) | int(OpLsh):           execLshUint32,
	(int(Uint64) << bitsOperator) | int(OpLsh):           execLshUint64,
	(int(Uintptr) << bitsOperator) | int(OpLsh):          execLshUintptr,
	(int(BigInt) << bitsOperator) | int(OpLsh):           execLshBigInt,
	(int(Int) << bitsOperator) | int(OpRsh):              execRshInt,
	(int(Int8) << bitsOperator) | int(OpRsh):             execRshInt8,
	(int(Int16) << bitsOperator) | int(OpRsh):            execRshInt16,
	(int(Int32) << bitsOperator) | int(OpRsh):            execRshInt32,
	(int(Int64) << bitsOperator) | int(OpRsh):            execRshInt64,
	(int(Uint) << bitsOperator) | int(OpRsh):             execRshUint,
	(int(Uint8) << bitsOperator) | int(OpRsh):            execRshUint8,
	(int(Uint16) << bitsOperator) | int(OpRsh):           execRshUint16,
	(int(Uint32) << bitsOperator) | int(OpRsh):           execRshUint32,
	(int(Uint64) << bitsOperator) | int(OpRsh):           execRshUint64,
	(int(Uintptr) << bitsOperator) | int(OpRsh):          execRshUintptr,
	(int(BigInt) << bitsOperator) | int(OpRsh):           execRshBigInt,
	(int(Int) << bitsOperator) | int(OpLT):               execLTInt,
	(int(Int8) << bitsOperator) | int(OpLT):              execLTInt8,
	(int(Int16) << bitsOperator) | int(OpLT):             execLTInt16,
	(int(Int32) << bitsOperator) | int(OpLT):             execLTInt32,
	(int(Int64) << bitsOperator) | int(OpLT):             execLTInt64,
	(int(Uint) << bitsOperator) | int(OpLT):              execLTUint,
	(int(Uint8) << bitsOperator) | int(OpLT):             execLTUint8,
	(int(Uint16) << bitsOperator) | int(OpLT):            execLTUint16,
	(int(Uint32) << bitsOperator) | int(OpLT):            execLTUint32,
	(int(Uint64) << bitsOperator) | int(OpLT):            execLTUint64,
	(int(Uintptr) << bitsOperator) | int(OpLT):           execLTUintptr,
	(int(Float32) << bitsOperator) | int(OpLT):           execLTFloat32,
	(int(Float64) << bitsOperator) | int(OpLT):           execLTFloat64,
	(int(String) << bitsOperator) | int(OpLT):            execLTString,
	(int(BigInt) << bitsOperator) | int(OpLT):            execLTBigInt,
	(int(BigRat) << bitsOperator) | int(OpLT):            execLTBigRat,
	(int(BigFloat) << bitsOperator) | int(OpLT):          execLTBigFloat,
	(int(Int) << bitsOperator) | int(OpLE):               execLEInt,
	(int(Int8) << bitsOperator) | int(OpLE):              execLEInt8,
	(int(Int16) << bitsOperator) | int(OpLE):             execLEInt16,
	(int(Int32) << bitsOperator) | int(OpLE):             execLEInt32,
	(int(Int64) << bitsOperator) | int(OpLE):             execLEInt64,
	(int(Uint) << bitsOperator) | int(OpLE):              execLEUint,
	(int(Uint8) << bitsOperator) | int(OpLE):             execLEUint8,
	(int(Uint16) << bitsOperator) | int(OpLE):            execLEUint16,
	(int(Uint32) << bitsOperator) | int(OpLE):            execLEUint32,
	(int(Uint64) << bitsOperator) | int(OpLE):            execLEUint64,
	(int(Uintptr) << bitsOperator) | int(OpLE):           execLEUintptr,
	(int(Float32) << bitsOperator) | int(OpLE):           execLEFloat32,
	(int(Float64) << bitsOperator) | int(OpLE):           execLEFloat64,
	(int(String) << bitsOperator) | int(OpLE):            execLEString,
	(int(BigInt) << bitsOperator) | int(OpLE):            execLEBigInt,
	(int(BigRat) << bitsOperator) | int(OpLE):            execLEBigRat,
	(int(BigFloat) << bitsOperator) | int(OpLE):          execLEBigFloat,
	(int(Int) << bitsOperator) | int(OpGT):               execGTInt,
	(int(Int8) << bitsOperator) | int(OpGT):              execGTInt8,
	(int(Int16) << bitsOperator) | int(OpGT):             execGTInt16,
	(int(Int32) << bitsOperator) | int(OpGT):             execGTInt32,
	(int(Int64) << bitsOperator) | int(OpGT):             execGTInt64,
	(int(Uint) << bitsOperator) | int(OpGT):              execGTUint,
	(int(Uint8) << bitsOperator) | int(OpGT):             execGTUint8,
	(int(Uint16) << bitsOperator) | int(OpGT):            execGTUint16,
	(int(Uint32) << bitsOperator) | int(OpGT):            execGTUint32,
	(int(Uint64) << bitsOperator) | int(OpGT):            execGTUint64,
	(int(Uintptr) << bitsOperator) | int(OpGT):           execGTUintptr,
	(int(Float32) << bitsOperator) | int(OpGT):           execGTFloat32,
	(int(Float64) << bitsOperator) | int(OpGT):           execGTFloat64,
	(int(String) << bitsOperator) | int(OpGT):            execGTString,
	(int(BigInt) << bitsOperator) | int(OpGT):            execGTBigInt,
	(int(BigRat) << bitsOperator) | int(OpGT):            execGTBigRat,
	(int(BigFloat) << bitsOperator) | int(OpGT):          execGTBigFloat,
	(int(Int) << bitsOperator) | int(OpGE):               execGEInt,
	(int(Int8) << bitsOperator) | int(OpGE):              execGEInt8,
	(int(Int16) << bitsOperator) | int(OpGE):             execGEInt16,
	(int(Int32) << bitsOperator) | int(OpGE):             execGEInt32,
	(int(Int64) << bitsOperator) | int(OpGE):             execGEInt64,
	(int(Uint) << bitsOperator) | int(OpGE):              execGEUint,
	(int(Uint8) << bitsOperator) | int(OpGE):             execGEUint8,
	(int(Uint16) << bitsOperator) | int(OpGE):            execGEUint16,
	(int(Uint32) << bitsOperator) | int(OpGE):            execGEUint32,
	(int(Uint64) << bitsOperator) | int(OpGE):            execGEUint64,
	(int(Uintptr) << bitsOperator) | int(OpGE):           execGEUintptr,
	(int(Float32) << bitsOperator) | int(OpGE):           execGEFloat32,
	(int(Float64) << bitsOperator) | int(OpGE):           execGEFloat64,
	(int(String) << bitsOperator) | int(OpGE):            execGEString,
	(int(BigInt) << bitsOperator) | int(OpGE):            execGEBigInt,
	(int(BigRat) << bitsOperator) | int(OpGE):            execGEBigRat,
	(int(BigFloat) << bitsOperator) | int(OpGE):          execGEBigFloat,
	(int(Int) << bitsOperator) | int(OpEQ):               execEQInt,
	(int(Int8) << bitsOperator) | int(OpEQ):              execEQInt8,
	(int(Int16) << bitsOperator) | int(OpEQ):             execEQInt16,
	(int(Int32) << bitsOperator) | int(OpEQ):             execEQInt32,
	(int(Int64) << bitsOperator) | int(OpEQ):             execEQInt64,
	(int(Uint) << bitsOperator) | int(OpEQ):              execEQUint,
	(int(Uint8) << bitsOperator) | int(OpEQ):             execEQUint8,
	(int(Uint16) << bitsOperator) | int(OpEQ):            execEQUint16,
	(int(Uint32) << bitsOperator) | int(OpEQ):            execEQUint32,
	(int(Uint64) << bitsOperator) | int(OpEQ):            execEQUint64,
	(int(Uintptr) << bitsOperator) | int(OpEQ):           execEQUintptr,
	(int(Float32) << bitsOperator) | int(OpEQ):           execEQFloat32,
	(int(Float64) << bitsOperator) | int(OpEQ):           execEQFloat64,
	(int(Complex64) << bitsOperator) | int(OpEQ):         execEQComplex64,
	(int(Complex128) << bitsOperator) | int(OpEQ):        execEQComplex128,
	(int(String) << bitsOperator) | int(OpEQ):            execEQString,
	(int(BigInt) << bitsOperator) | int(OpEQ):            execEQBigInt,
	(int(BigRat) << bitsOperator) | int(OpEQ):            execEQBigRat,
	(int(BigFloat) << bitsOperator) | int(OpEQ):          execEQBigFloat,
	(int(Int) << bitsOperator) | int(OpNE):               execNEInt,
	(int(Int8) << bitsOperator) | int(OpNE):              execNEInt8,
	(int(Int16) << bitsOperator) | int(OpNE):             execNEInt16,
	(int(Int32) << bitsOperator) | int(OpNE):             execNEInt32,
	(int(Int64) << bitsOperator) | int(OpNE):             execNEInt64,
	(int(Uint) << bitsOperator) | int(OpNE):              execNEUint,
	(int(Uint8) << bitsOperator) | int(OpNE):             execNEUint8,
	(int(Uint16) << bitsOperator) | int(OpNE):            execNEUint16,
	(int(Uint32) << bitsOperator) | int(OpNE):            execNEUint32,
	(int(Uint64) << bitsOperator) | int(OpNE):            execNEUint64,
	(int(Uintptr) << bitsOperator) | int(OpNE):           execNEUintptr,
	(int(Float32) << bitsOperator) | int(OpNE):           execNEFloat32,
	(int(Float64) << bitsOperator) | int(OpNE):           execNEFloat64,
	(int(Complex64) << bitsOperator) | int(OpNE):         execNEComplex64,
	(int(Complex128) << bitsOperator) | int(OpNE):        execNEComplex128,
	(int(String) << bitsOperator) | int(OpNE):            execNEString,
	(int(BigInt) << bitsOperator) | int(OpNE):            execNEBigInt,
	(int(BigRat) << bitsOperator) | int(OpNE):            execNEBigRat,
	(int(BigFloat) << bitsOperator) | int(OpNE):          execNEBigFloat,
	(int(Bool) << bitsOperator) | int(OpLAnd):            execLAndBool,
	(int(Bool) << bitsOperator) | int(OpLOr):             execLOrBool,
	(int(Bool) << bitsOperator) | int(OpLNot):            execLNotBool,
	(int(Bool) << bitsOperator) | int(OpEQ):              execEQBool,
	(int(Bool) << bitsOperator) | int(OpNE):              execNEBool,
	(int(Int) << bitsOperator) | int(OpNeg):              execNegInt,
	(int(Int8) << bitsOperator) | int(OpNeg):             execNegInt8,
	(int(Int16) << bitsOperator) | int(OpNeg):            execNegInt16,
	(int(Int32) << bitsOperator) | int(OpNeg):            execNegInt32,
	(int(Int64) << bitsOperator) | int(OpNeg):            execNegInt64,
	(int(Uint) << bitsOperator) | int(OpNeg):             execNegUint,
	(int(Uint8) << bitsOperator) | int(OpNeg):            execNegUint8,
	(int(Uint16) << bitsOperator) | int(OpNeg):           execNegUint16,
	(int(Uint32) << bitsOperator) | int(OpNeg):           execNegUint32,
	(int(Uint64) << bitsOperator) | int(OpNeg):           execNegUint64,
	(int(Uintptr) << bitsOperator) | int(OpNeg):          execNegUintptr,
	(int(Float32) << bitsOperator) | int(OpNeg):          execNegFloat32,
	(int(Float64) << bitsOperator) | int(OpNeg):          execNegFloat64,
	(int(Complex64) << bitsOperator) | int(OpNeg):        execNegComplex64,
	(int(Complex128) << bitsOperator) | int(OpNeg):       execNegComplex128,
	(int(BigInt) << bitsOperator) | int(OpNeg):           execNegBigInt,
	(int(BigRat) << bitsOperator) | int(OpNeg):           execNegBigRat,
	(int(BigFloat) << bitsOperator) | int(OpNeg):         execNegBigFloat,
	(int(Int) << bitsOperator) | int(OpBitNot):           execBitNotInt,
	(int(Int8) << bitsOperator) | int(OpBitNot):          execBitNotInt8,
	(int(Int16) << bitsOperator) | int(OpBitNot):         execBitNotInt16,
	(int(Int32) << bitsOperator) | int(OpBitNot):         execBitNotInt32,
	(int(Int64) << bitsOperator) | int(OpBitNot):         execBitNotInt64,
	(int(Uint) << bitsOperator) | int(OpBitNot):          execBitNotUint,
	(int(Uint8) << bitsOperator) | int(OpBitNot):         execBitNotUint8,
	(int(Uint16) << bitsOperator) | int(OpBitNot):        execBitNotUint16,
	(int(Uint32) << bitsOperator) | int(OpBitNot):        execBitNotUint32,
	(int(Uint64) << bitsOperator) | int(OpBitNot):        execBitNotUint64,
	(int(Uintptr) << bitsOperator) | int(OpBitNot):       execBitNotUintptr,
	(int(BigInt) << bitsOperator) | int(OpBitNot):        execBitNotBigInt,
	(int(Slice) << bitsOperator) | int(OpEQ):             execRefTypeEQ,
	(int(Map) << bitsOperator) | int(OpEQ):               execRefTypeEQ,
	(int(Chan) << bitsOperator) | int(OpEQ):              execRefTypeEQ,
	(int(Ptr) << bitsOperator) | int(OpEQ):               execRefTypeEQ,
	(int(Slice) << bitsOperator) | int(OpNE):             execRefTypeNE,
	(int(Map) << bitsOperator) | int(OpNE):               execRefTypeNE,
	(int(Chan) << bitsOperator) | int(OpNE):              execRefTypeNE,
	(int(Ptr) << bitsOperator) | int(OpNE):               execRefTypeNE,
	(int(UnsafePointer) << bitsOperator) | int(OpEQ):     execRefTypeEQ,
	(int(UnsafePointer) << bitsOperator) | int(OpNE):     execRefTypeNE,
	(int(reflect.Func) << bitsOperator) | int(OpEQ):      execRefTypeEQ,
	(int(reflect.Func) << bitsOperator) | int(OpNE):      execRefTypeNE,
	(int(reflect.Interface) << bitsOperator) | int(OpEQ): execInterfaceEQ,
	(int(reflect.Interface) << bitsOperator) | int(OpNE): execInterfaceNE,
}

func execEQBool(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(bool) == p.data[n-1].(bool)
	p.data = p.data[:n-1]
}

func execNEBool(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(bool) != p.data[n-1].(bool)
	p.data = p.data[:n-1]
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

func execAddBigInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = new(big.Int).Add(p.data[n-2].(*big.Int), p.data[n-1].(*big.Int))
	p.data = p.data[:n-1]
}

func execAddBigRat(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = new(big.Rat).Add(p.data[n-2].(*big.Rat), p.data[n-1].(*big.Rat))
	p.data = p.data[:n-1]
}

func execAddBigFloat(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = new(big.Float).Add(p.data[n-2].(*big.Float), p.data[n-1].(*big.Float))
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

func execSubBigInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = new(big.Int).Sub(p.data[n-2].(*big.Int), p.data[n-1].(*big.Int))
	p.data = p.data[:n-1]
}

func execSubBigRat(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = new(big.Rat).Sub(p.data[n-2].(*big.Rat), p.data[n-1].(*big.Rat))
	p.data = p.data[:n-1]
}

func execSubBigFloat(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = new(big.Float).Sub(p.data[n-2].(*big.Float), p.data[n-1].(*big.Float))
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

func execMulBigInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = new(big.Int).Mul(p.data[n-2].(*big.Int), p.data[n-1].(*big.Int))
	p.data = p.data[:n-1]
}

func execMulBigRat(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = new(big.Rat).Mul(p.data[n-2].(*big.Rat), p.data[n-1].(*big.Rat))
	p.data = p.data[:n-1]
}

func execMulBigFloat(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = new(big.Float).Mul(p.data[n-2].(*big.Float), p.data[n-1].(*big.Float))
	p.data = p.data[:n-1]
}

func execQuoInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) / p.data[n-1].(int)
	p.data = p.data[:n-1]
}

func execQuoInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) / p.data[n-1].(int8)
	p.data = p.data[:n-1]
}

func execQuoInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) / p.data[n-1].(int16)
	p.data = p.data[:n-1]
}

func execQuoInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) / p.data[n-1].(int32)
	p.data = p.data[:n-1]
}

func execQuoInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) / p.data[n-1].(int64)
	p.data = p.data[:n-1]
}

func execQuoUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) / p.data[n-1].(uint)
	p.data = p.data[:n-1]
}

func execQuoUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) / p.data[n-1].(uint8)
	p.data = p.data[:n-1]
}

func execQuoUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) / p.data[n-1].(uint16)
	p.data = p.data[:n-1]
}

func execQuoUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) / p.data[n-1].(uint32)
	p.data = p.data[:n-1]
}

func execQuoUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) / p.data[n-1].(uint64)
	p.data = p.data[:n-1]
}

func execQuoUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) / p.data[n-1].(uintptr)
	p.data = p.data[:n-1]
}

func execQuoFloat32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float32) / p.data[n-1].(float32)
	p.data = p.data[:n-1]
}

func execQuoFloat64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(float64) / p.data[n-1].(float64)
	p.data = p.data[:n-1]
}

func execQuoComplex64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(complex64) / p.data[n-1].(complex64)
	p.data = p.data[:n-1]
}

func execQuoComplex128(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(complex128) / p.data[n-1].(complex128)
	p.data = p.data[:n-1]
}

func execQuoBigRat(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = new(big.Rat).Quo(p.data[n-2].(*big.Rat), p.data[n-1].(*big.Rat))
	p.data = p.data[:n-1]
}

func execQuoBigFloat(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = new(big.Float).Quo(p.data[n-2].(*big.Float), p.data[n-1].(*big.Float))
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

func execModBigInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = new(big.Int).Mod(p.data[n-2].(*big.Int), p.data[n-1].(*big.Int))
	p.data = p.data[:n-1]
}

func execAndInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) & p.data[n-1].(int)
	p.data = p.data[:n-1]
}

func execAndInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) & p.data[n-1].(int8)
	p.data = p.data[:n-1]
}

func execAndInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) & p.data[n-1].(int16)
	p.data = p.data[:n-1]
}

func execAndInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) & p.data[n-1].(int32)
	p.data = p.data[:n-1]
}

func execAndInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) & p.data[n-1].(int64)
	p.data = p.data[:n-1]
}

func execAndUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) & p.data[n-1].(uint)
	p.data = p.data[:n-1]
}

func execAndUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) & p.data[n-1].(uint8)
	p.data = p.data[:n-1]
}

func execAndUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) & p.data[n-1].(uint16)
	p.data = p.data[:n-1]
}

func execAndUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) & p.data[n-1].(uint32)
	p.data = p.data[:n-1]
}

func execAndUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) & p.data[n-1].(uint64)
	p.data = p.data[:n-1]
}

func execAndUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) & p.data[n-1].(uintptr)
	p.data = p.data[:n-1]
}

func execAndBigInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = new(big.Int).And(p.data[n-2].(*big.Int), p.data[n-1].(*big.Int))
	p.data = p.data[:n-1]
}

func execOrInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) | p.data[n-1].(int)
	p.data = p.data[:n-1]
}

func execOrInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) | p.data[n-1].(int8)
	p.data = p.data[:n-1]
}

func execOrInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) | p.data[n-1].(int16)
	p.data = p.data[:n-1]
}

func execOrInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) | p.data[n-1].(int32)
	p.data = p.data[:n-1]
}

func execOrInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) | p.data[n-1].(int64)
	p.data = p.data[:n-1]
}

func execOrUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) | p.data[n-1].(uint)
	p.data = p.data[:n-1]
}

func execOrUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) | p.data[n-1].(uint8)
	p.data = p.data[:n-1]
}

func execOrUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) | p.data[n-1].(uint16)
	p.data = p.data[:n-1]
}

func execOrUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) | p.data[n-1].(uint32)
	p.data = p.data[:n-1]
}

func execOrUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) | p.data[n-1].(uint64)
	p.data = p.data[:n-1]
}

func execOrUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) | p.data[n-1].(uintptr)
	p.data = p.data[:n-1]
}

func execOrBigInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = new(big.Int).Or(p.data[n-2].(*big.Int), p.data[n-1].(*big.Int))
	p.data = p.data[:n-1]
}

func execXorInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) ^ p.data[n-1].(int)
	p.data = p.data[:n-1]
}

func execXorInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) ^ p.data[n-1].(int8)
	p.data = p.data[:n-1]
}

func execXorInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) ^ p.data[n-1].(int16)
	p.data = p.data[:n-1]
}

func execXorInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) ^ p.data[n-1].(int32)
	p.data = p.data[:n-1]
}

func execXorInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) ^ p.data[n-1].(int64)
	p.data = p.data[:n-1]
}

func execXorUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) ^ p.data[n-1].(uint)
	p.data = p.data[:n-1]
}

func execXorUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) ^ p.data[n-1].(uint8)
	p.data = p.data[:n-1]
}

func execXorUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) ^ p.data[n-1].(uint16)
	p.data = p.data[:n-1]
}

func execXorUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) ^ p.data[n-1].(uint32)
	p.data = p.data[:n-1]
}

func execXorUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) ^ p.data[n-1].(uint64)
	p.data = p.data[:n-1]
}

func execXorUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) ^ p.data[n-1].(uintptr)
	p.data = p.data[:n-1]
}

func execXorBigInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = new(big.Int).Xor(p.data[n-2].(*big.Int), p.data[n-1].(*big.Int))
	p.data = p.data[:n-1]
}

func execAndNotInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) &^ p.data[n-1].(int)
	p.data = p.data[:n-1]
}

func execAndNotInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) &^ p.data[n-1].(int8)
	p.data = p.data[:n-1]
}

func execAndNotInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) &^ p.data[n-1].(int16)
	p.data = p.data[:n-1]
}

func execAndNotInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) &^ p.data[n-1].(int32)
	p.data = p.data[:n-1]
}

func execAndNotInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) &^ p.data[n-1].(int64)
	p.data = p.data[:n-1]
}

func execAndNotUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) &^ p.data[n-1].(uint)
	p.data = p.data[:n-1]
}

func execAndNotUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) &^ p.data[n-1].(uint8)
	p.data = p.data[:n-1]
}

func execAndNotUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) &^ p.data[n-1].(uint16)
	p.data = p.data[:n-1]
}

func execAndNotUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) &^ p.data[n-1].(uint32)
	p.data = p.data[:n-1]
}

func execAndNotUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) &^ p.data[n-1].(uint64)
	p.data = p.data[:n-1]
}

func execAndNotUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) &^ p.data[n-1].(uintptr)
	p.data = p.data[:n-1]
}

func execAndNotBigInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = new(big.Int).AndNot(p.data[n-2].(*big.Int), p.data[n-1].(*big.Int))
	p.data = p.data[:n-1]
}

func execLshInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) << toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execLshInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) << toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execLshInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) << toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execLshInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) << toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execLshInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) << toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execLshUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) << toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execLshUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) << toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execLshUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) << toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execLshUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) << toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execLshUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) << toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execLshUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) << toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execLshBigInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = new(big.Int).Lsh(p.data[n-2].(*big.Int), toUint(p.data[n-1]))
	p.data = p.data[:n-1]
}

func execRshInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int) >> toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execRshInt8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int8) >> toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execRshInt16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int16) >> toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execRshInt32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int32) >> toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execRshInt64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(int64) >> toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execRshUint(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint) >> toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execRshUint8(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint8) >> toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execRshUint16(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint16) >> toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execRshUint32(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint32) >> toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execRshUint64(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uint64) >> toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execRshUintptr(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(uintptr) >> toUint(p.data[n-1])
	p.data = p.data[:n-1]
}

func execRshBigInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = new(big.Int).Rsh(p.data[n-2].(*big.Int), toUint(p.data[n-1]))
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

func execLTBigInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(*big.Int).Cmp(p.data[n-1].(*big.Int)) < 0
	p.data = p.data[:n-1]
}

func execLTBigRat(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(*big.Rat).Cmp(p.data[n-1].(*big.Rat)) < 0
	p.data = p.data[:n-1]
}

func execLTBigFloat(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(*big.Float).Cmp(p.data[n-1].(*big.Float)) < 0
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

func execLEBigInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(*big.Int).Cmp(p.data[n-1].(*big.Int)) <= 0
	p.data = p.data[:n-1]
}

func execLEBigRat(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(*big.Rat).Cmp(p.data[n-1].(*big.Rat)) <= 0
	p.data = p.data[:n-1]
}

func execLEBigFloat(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(*big.Float).Cmp(p.data[n-1].(*big.Float)) <= 0
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

func execGTBigInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(*big.Int).Cmp(p.data[n-1].(*big.Int)) > 0
	p.data = p.data[:n-1]
}

func execGTBigRat(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(*big.Rat).Cmp(p.data[n-1].(*big.Rat)) > 0
	p.data = p.data[:n-1]
}

func execGTBigFloat(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(*big.Float).Cmp(p.data[n-1].(*big.Float)) > 0
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

func execGEBigInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(*big.Int).Cmp(p.data[n-1].(*big.Int)) >= 0
	p.data = p.data[:n-1]
}

func execGEBigRat(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(*big.Rat).Cmp(p.data[n-1].(*big.Rat)) >= 0
	p.data = p.data[:n-1]
}

func execGEBigFloat(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(*big.Float).Cmp(p.data[n-1].(*big.Float)) >= 0
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

func execEQBigInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(*big.Int).Cmp(p.data[n-1].(*big.Int)) == 0
	p.data = p.data[:n-1]
}

func execEQBigRat(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(*big.Rat).Cmp(p.data[n-1].(*big.Rat)) == 0
	p.data = p.data[:n-1]
}

func execEQBigFloat(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(*big.Float).Cmp(p.data[n-1].(*big.Float)) == 0
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

func execNEBigInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(*big.Int).Cmp(p.data[n-1].(*big.Int)) != 0
	p.data = p.data[:n-1]
}

func execNEBigRat(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(*big.Rat).Cmp(p.data[n-1].(*big.Rat)) != 0
	p.data = p.data[:n-1]
}

func execNEBigFloat(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-2] = p.data[n-2].(*big.Float).Cmp(p.data[n-1].(*big.Float)) != 0
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

func execLNotBool(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = !p.data[n-1].(bool)
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

func execNegBigInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = new(big.Int).Neg(p.data[n-1].(*big.Int))
}

func execNegBigRat(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = new(big.Rat).Neg(p.data[n-1].(*big.Rat))
}

func execNegBigFloat(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = new(big.Float).Neg(p.data[n-1].(*big.Float))
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

func execBitNotBigInt(i Instr, p *Context) {
	n := len(p.data)
	p.data[n-1] = new(big.Int).Not(p.data[n-1].(*big.Int))
}
