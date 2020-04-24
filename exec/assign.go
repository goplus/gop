package exec

// -----------------------------------------------------------------------------

const (
	// OpAssign `=`
	OpAssign = iota
	// OpAddAssign `+=`
	OpAddAssign
	// OpSubAssign `-=`
	OpSubAssign
	// OpMulAssign `*=`
	OpMulAssign
	// OpDivAssign `/=`
	OpDivAssign
	// OpModAssign `%=`
	OpModAssign

	// OpBitAndAssign '&='
	OpBitAndAssign
	// OpBitOrAssign '|='
	OpBitOrAssign
	// OpBitXorAssign '^='
	OpBitXorAssign
	// OpBitAndNotAssign '&^='
	OpBitAndNotAssign
	// OpBitSHLAssign '<<='
	OpBitSHLAssign
	// OpBitSHRAssign '>>='
	OpBitSHRAssign

	// OpInc '++'
	OpInc
	// OpDec '--'
	OpDec
)

// -----------------------------------------------------------------------------
