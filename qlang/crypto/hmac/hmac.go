package hmac

import "crypto/hmac"

// -----------------------------------------------------------------------------

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"new": hmac.New,
}

// -----------------------------------------------------------------------------
