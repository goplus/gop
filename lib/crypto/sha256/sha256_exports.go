package sha256

import "crypto/sha256"

// -----------------------------------------------------------------------------

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"_name":  "crypto/sha256",
	"new":    sha256.New,
	"new224": sha256.New224,

	"New":    sha256.New,
	"New224": sha256.New224,

	"BlockSize": sha256.BlockSize,
	"Size":      sha256.Size,
}

// -----------------------------------------------------------------------------
